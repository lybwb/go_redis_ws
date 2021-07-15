package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	RedisClient *redis.Pool
)

func InitRedis(host string, auth string, db int) error {

	RedisClient = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   4000,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host, redis.DialPassword(auth), redis.DialDatabase(db))
			if nil != err {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	rd := RedisClient.Get()
	defer rd.Close()

	c, err := redis.Dial("tcp", host, redis.DialPassword(auth), redis.DialDatabase(db))
	defer c.Close()
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return err
	}
	fmt.Println("Connect to redis ok")
	return nil

}

func IsConnError(err error) bool {
	var needNewConn bool

	if err == nil {
		return false
	}

	if err == io.EOF {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "connect: connection refused") {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "connection closed") {
		needNewConn = true
	}
	return needNewConn
}

// 在pool加入TestOnBorrow方法来去除扫描坏连接
func Redo(command string, opt ...interface{}) (interface{}, error) {
	if RedisClient == nil {
		return "", errors.New("error,redis client is null")
	}
	rd := RedisClient.Get()
	defer rd.Close()

	var conn redis.Conn
	var err error
	var maxretry = 3
	var needNewConn bool

	resp, err := rd.Do(command, opt...)
	needNewConn = IsConnError(err)
	if needNewConn == false {
		return resp, err
	} else {
		conn, err = RedisClient.Dial()
	}

	for index := 0; index < maxretry; index++ {
		if conn == nil && index+1 > maxretry {
			return resp, err
		}
		if conn == nil {
			conn, err = RedisClient.Dial()
		}
		if err != nil {
			continue
		}

		resp, err := conn.Do(command, opt...)
		needNewConn = IsConnError(err)
		if needNewConn == false {
			return resp, err
		} else {
			conn, err = RedisClient.Dial()
		}
	}

	conn.Close()
	return "", errors.New("redis error")
}

type SubscribeCallback func(topicMap sync.Map, topic, msg string)

type Subscriber struct {
	client   redis.PubSubConn
	Ws       *WSHub //websocket
	cbMap    sync.Map
	CallBack interface {
		OnReceive(SubscribeCallback)
	}
}

var fnSubReceived SubscribeCallback

func (c *Subscriber) OnReceive(cb SubscribeCallback) {
	fnSubReceived = cb
}

func (c *Subscriber) Init(ws *WSHub) {

	conn := RedisClient.Get()

	c.client = redis.PubSubConn{conn}
	c.Ws = ws
	go func() {
		for {
			log.Println("redis wait...")
			switch res := c.client.Receive().(type) {
			case redis.Message:
				fmt.Printf("receive:%#v\n", res)
				topic := res.Channel
				message := string(res.Data)
				fnSubReceived(c.cbMap, topic, message)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", res.Channel, res.Kind, res.Count)
			case error:
				log.Println("error handle", res)
				if IsConnError(res) {
					conn, err := RedisClient.Dial()
					if err != nil {
						log.Printf("err=%s\n", err)
					}
					c.client = redis.PubSubConn{conn}
				}
				continue
			}
		}
	}()

}

func (c *Subscriber) Close() {
	err := c.client.Close()
	if err != nil {
		log.Println("redis close error.")
	}
}

func (c *Subscriber) Subscribe(channel interface{}, clientid string) {
	err := c.client.Subscribe(channel)
	if err != nil {
		log.Println("redis Subscribe error.", err)
	}
	c.cbMap.Store(clientid, channel.(string))
}

func (c *Subscriber) PSubscribe(channel interface{}, clientid string) {
	err := c.client.PSubscribe(channel)
	if err != nil {
		log.Println("redis PSubscribe error.", err)
	}

	c.cbMap.Store(clientid, channel.(string))
}
