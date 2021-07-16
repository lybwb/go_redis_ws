package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/emr"
)

type AuthCredential struct {
	AccessToken         string    `json:"access_token"`
	AccessTokenExpireTs time.Time `json:"access_token_expire_ts"`
}

type User struct {
	UserID        int               `json:"uid"`
	Email         string            `json:"email"`
	UserName      string            `json:"user_name"`
	PhoneNumber   string            `json:"phone_number"`
	UserType      string            `json:"user_type"`
	OwnerUID      int               `json:"owner_uid"`
	OwnerEmail    string            `json:"owner_email"`
	Organization  string            `json:"organization"`
	PlatformLogo  string            `json:"platform_logo"`
	AvatarSrc     string            `json:"avatar_src"`
	Disable       bool              `json:"disable"`
	TsCreated     string            `json:"ts_created"`
	Collaborators []Collaborator    `json:"collaborators"` // users that collaborates this user
	ServiceList   map[string]string `json:"service_list"`

	Community *Community `json:"community,omitempty"`
}

type UserInfo struct {
	UserID       int     `json:"uid"`
	Email        string  `json:"email"`
	UserName     string  `json:"user_name"`
	PhoneNumber  *string `json:"phone_number,omitempty"`
	UserType     *string `json:"user_type,omitempty"`
	Organization *string `json:"organization,omitempty"`
	PlatformLogo *string `json:"platform_logo,omitempty"`
	AvatarSrc    *string `json:"avatar_src,omitempty"`
	Disable      *bool   `json:"disable,omitempty"`
	TsCreated    *string `json:"ts_created,omitempty"`
}

type UserDetails struct {
	UserInfo
	OwnerUID      *int           `json:"owner_uid,omitempty"`
	Owner         *UserInfo      `json:"owner,omitempty"`
	Collaboration *Collaboration `json:"collaboration,omitempty"`
	Community     *Community     `json:"community,omitempty"`
}

type Community struct {
	CommunityCID  int     `json:"cid"`
	CommunityID   string  `json:"community_id"`
	CommunityName string  `json:"community_name"`
	Description   *string `json:"description,omitempty"`
}

type Collaboration struct {
	UserID        int                `json:"uid"`
	OwnerUID      int                `json:"owner_uid"`
	NickName      string             `json:"nickname"`
	Collaborators *[]Collaborator    `json:"collaborators,omitempty"`
	ServiceList   *map[string]string `json:"service_list,omitempty"`
}

type Collaborator struct {
	UserID   int    `json:"uid"`
	NickName string `json:"nickname"`
}
type ResourceActions struct {
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}

type ActionResources struct {
	Action    string   `json:"action"`
	Resources []string `json:"resources"`
}

type PolicyStatment struct {
	Sid       string   `json:"sid"`
	Effect    string   `json:"effect"`
	Actions   []string `json:"actions"`
	Resources []string `json:"resources"`
}

type Policy struct {
	Version    string           `json:"version"`
	Statements []PolicyStatment `json:"statements"`
}

type AWSTempToken struct {
	AccessKey          string    `json:"access_key"`
	SecretAccessKey    string    `json:"secret_access_key"`
	SessionToken       string    `json:"session_token"`
	ExpireAt           time.Time `json:"expire_at"`
	ExpireAfterSeconds int64     `json:"expire_after_seconds"`
}

type SignUpRequest struct {
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	Email       string    `json:"email"`
	CompanySize int       `json:"company_size"`
	CompanyName string    `json:"company_name"`
	Status      int       `json:"status"`
	GoThroughBy int       `json:"gothrough_by"`
	TsCreated   time.Time `json:"ts_created"`
}

type UserGroup struct {
	GroupID   string            `json:"group_id"`
	GroupName string            `json:"group_name"`
	OwnerID   int               `json:"owner_id"`
	Members   []UserGroupMember `json:"members"`
}

type UserGroupMember struct {
	GroupID  string `json:"group_id"`
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

type UserGroupUpdateForm struct {
	UserGroup
}

// InvitedUser used to store invite request
type InvitedUser struct {
	Name              string `json:"name"`
	Email             string `json:"email" binding:"required"`
	VerificationToken string `json:"verification_token" binding:"required"`
}

// InviteInfo used to store invite information
type InviteInfo struct {
	InvitedUsers []InvitedUser `json:"invited_users" binding:"required,dive"`
	UserType     string        `json:"user_type" binding:"required"`
}

// InvitedUser used to store single user info
type InvitedUserDetails struct {
	InvitedUser
	InviterUID         int       `json:"inviter_uid"`
	InviterUserDetails *User     `json:"inviter_user_details"`
	UserType           string    `json:"user_type"`
	PlatformLogo       string    `json:"platform_logo"`
	InvitationTime     time.Time `json:"invitation_time"`
}

type EmailSt struct {
	EmailAddress string            `json:"email_address"`
	Subject      string            `json:"subject"`
	Body         string            `json:"body"`
	Attachments  map[string]string `json:"attachments"`
}

type NotifyCategory struct {
	Category    string `json:"category"`
	NotifyType  int    `json:"notify_type"`
	Description string `json:"description"`
}

type NotifyMessage struct {
	MessageID  string `json:"message_id"`
	Category   string `json:"category"`
	NotifyType int    `json:"notify_type"`
	CreatorUID int    `json:"creator_uid"`
	Message    string `json:"ts_publish"`
}

type NotifyUserMessage struct {
	MessageID        string    `json:"message_id"`
	Message          string    `json:"message"`
	Category         string    `json:"category"`
	NotifyType       int       `json:"notify_type"`
	SenderUID        int       `json:"sender_uid"`
	SenderName       string    `json:"sender_name"`
	Uid              int       `json:"uid"`
	Read             bool      `json:"read"`
	TsMessagePublish time.Time `json:"ts_message_publish"`
	TsRead           time.Time `json:"ts_read"`
}

type UserNotifyUnread struct {
	UnreadNoticeCount            int  `json:"unread_notice_count"`
	UnreadMessageCount           int  `json:"unread_message_count"`
	UnhandledContentRequestCount *int `json:"unhandled_content_request_count,omitempty"`
}

type NotifyFilterForm struct {
	FilterTrash     bool               `json:"filter_trash"`
	FilterRead      bool               `json:"filter_read"`
	FilterUnread    bool               `json:"filter_unread"`
	FilterSearchKey bool               `json:"filter_search_key"`
	SearchKey       string             `json:"search_key"`
	NotifyType      int                `json:"notify_type,required"`
	LastNotify      *NotifyUserMessage `json:"last_notify,required"`
}

type NotifyFilterResult struct {
	Notifices          []NotifyUserMessage `json:"notifies"`
	SenderAvatarSrcMap map[int]string      `json:"sender_avatar_src_map"`
	NotifyType         int                 `json:"notify_type"`
}

type PointLatLng struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type EnvelopeLatLng struct {
	NorthLat float64 `json:"north_lat"`
	SouthLat float64 `json:"south_lat"`
	EastLng  float64 `json:"east_lng"`
	WestLng  float64 `json:"west_lng"`
}

type CenterLatLng struct {
	CenterLat float64 `json:"center_lat"`
	CenterLng float64 `json:"center_lng"`
}

type EnvelopeCenter struct {
	EnvelopeLatLng
	CenterLatLng
}

// MapboxDataset
/*
Example:
{
	dataset: "7fde68f4-d5f1-40d1-9af0-7c47ba4f1d59",
	otherdata: "{"access_token":"Mk4orbHJOie2ZwhbsTRZbCCFgE6AbtAFkQqg0Ei-Ptg=","id":1125,"total_layer_row_count":787029}"
}
*/
type MapboxDataset struct {
	Dataset   string `json:"dataset"`
	OtherData string `json:"otherdata"`
}

type RequestVectorData struct {
	ID                       int    `json:"id"`
	PageURL                  string `json:"page_url"`
	SubmittedTime            int    `json:"submitted_time"`
	Organization             string `json:"organization"`
	PositionTitle            string `json:"position_title"`
	FirstName                string `json:"firstname"`
	LastName                 string `json:"lastname"`
	Phone                    string `json:"phone"`
	Email                    string `json:"email"`
	DataType                 string `json:"data_type"`
	Purpose                  string `json:"purpose"`
	UsageDescription         string `json:"usage_description"`
	AreaOfInterest           string `json:"area_of_interest"`
	AOIPath                  string `json:"aoi_path"`
	AdditionalAreaOfInterest string `json:"additional_areas_of_interest"`
	ReferralPartner          string `json:"referral_partner"`
	UploadShapeFileOfRegion  string `json:"upload_shapefile_of_region"`
	Status                   string `json:"status"`
	ReviewedBy               string `json:"reviewed_by"`
}

type AwsPEMKey struct {
	KeyID      string `json:"key_id"`
	PrivateKey string `json:"private_key"`
}

type Point struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type GeoCodingSearchFilter struct {
	Key         string `json:"key" binding:"required"`
	TopLeft     Point  `json:"top_left"`
	BottomRight Point  `json:"bottom_right"`
}

type FilterLocation struct {
	Wkt     string `json:"wkt"`
	Point   Point  `json:"location"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type GeoCodingFilterTarget struct {
	Status bool              `json:"status"`
	Result *[]FilterLocation `json:"result"`
}

type CatetitleDefault struct {
	Catetitle string  `json:"catetitle"`
	Color     string  `json:"color"`
	ZIndex    int     `json:"z_index"`
	Opacity   float32 `json:"opacity"`
	Width     int     `json:"width"`
}

type CatetitleInfo struct {
	CatetitleDefault
	GeomType   string `json:"geom_type"`
	RenderType string `json:"render_type"`
}

type CatetitleDetails struct {
	CatetitleDefault
	GeomTypes  []string `json:"geom_types"`
	RenderType string   `json:"render_type"`
}

type CatetitleGeomType struct {
	Catetitle string `json:"catetitle"`
	GeomType  string `json:"geom_type"`
}

type LayerAttributeFeatures struct {
	LayerID           int                       `json:"layer_id"`
	LayerName         string                    `json:"layer_name"`
	TotalFeatureCount int                       `json:"total_feature_count"`
	Offset            int                       `json:"offset"`
	Size              int                       `json:"size"`
	AttributeDataType map[string]string         `json:"attribute_data_type"`
	Features          *[]map[string]interface{} `json:"features"`
}

type CodingStatement struct {
	UserID    int    `json:"uid"`
	Statement string `json:"statement"`
	TsCreated string `json:"ts_created"`
}

type CodingStatementCatetitle struct {
	Catetitle string           `json:"catetitle"`
	GeomTypes *json.RawMessage `json:"geom_types"`
}

type CodingStatementCatetitles struct {
	CodingStatement
	CenterLat  float64                    `json:"center_lat"`
	CenterLng  float64                    `json:"center_lng"`
	Catetitles []CodingStatementCatetitle `json:"catetitles"`
}

type CodingStatementHistory struct {
	Statement     string `json:"statement"`
	TsCreated     string `json:"ts_created"`
	ErrCode       int    `json:"errcode"`
	ErrMesg       string `json:"errmesg"`
	ExecLatencyMs int64  `json:"exec_latency_ms"`
}

type CodingStatementHistoryList struct {
	UserID     int                      `json:"uid"`
	Statements []CodingStatementHistory `json:"statements"`
}

type CodingStatementExecuteResult struct {
	Statement     string                    `json:"statement"`
	ErrCode       int                       `json:"errcode"`
	ErrMesg       string                    `json:"errmesg"`
	ExecLatencyMs int64                     `json:"exec_latency_ms"`
	Offset        int                       `json:"offset"`
	Size          int                       `json:"size"`
	TotalRowCount int                       `json:"total_row_count"`
	VisibleFields []string                  `json:"visible_fields"`
	ResultData    *[]map[string]interface{} `json:"result_data"`
}

type SparkClusterBase struct {
	Cid       int    `json:"cid"`
	ClusterID string `json:"cluster_id" required:"true"`
}

type SparkCluster struct {
	SparkClusterBase
	Name             string                `json:"name"`
	TargetType       string                `json:"target_type"`
	TargetID         string                `json:"target_id"`
	TargetHost       string                `json:"target_host"`
	Zone             string                `json:"zone"`
	RuntimeVersion   string                `json:"runtime_version"`
	Runtime          string                `json:"runtime"`
	NodeType         string                `json:"node_type"`
	NodeNum          int64                 `json:"node_num"`
	NodeCpu          int64                 `json:"node_cpu"`
	NodeMem          int64                 `json:"node_mem"`
	Extend           string                `json:"extend"`
	OwnerUID         int                   `json:"owner_uid"`
	UserID           int                   `json:"uid"`
	Status           string                `json:"status"`
	Disable          bool                  `json:"disable"`
	TsCreated        *time.Time            `json:"ts_created"`
	EmrClusterDetail *emr.Cluster          `json:"emr_cluster_detail"`
	InstanceGroups   []*SparkInstanceGroup `json:"instance_groups"`
	EC2Instance      *ec2.Instance         `json:"ec2_instance"`
}

type SparkClusterListOutput struct {
	SparkClusters []*SparkCluster `json:"spark_clusters"`
	Offset        int             `json:"offset"`
	Size          int             `json:"size"`
	TotalSize     int             `json:"total_size"`
}

type SparkInstanceGroup struct {
	ID                     string             `json:"id"`
	Name                   string             `json:"name"`
	InstanceGroupType      string             `json:"instance_group_type"`
	InstanceType           string             `json:"instance_type"`
	RunningInstanceCount   int64              `json:"running_instance_count"`
	RequestedInstanceCount int64              `json:"requested_instance_count"`
	Status                 string             `json:"status"`
	EmrInstanceGroupDetail *emr.InstanceGroup `json:"emr_instance_group_detail"`
}

type ModifySparkInstanceGroup struct {
	InstanceGroupID string `json:"instance_group_id" binding:"required"`
	InstanceCount   int64  `json:"instance_count" binding:"required"`
}

type ClusterReqTaskResponse struct {
	ErrCode           int
	ErrMsg            string
	EmrCluster        *emr.Cluster
	EmrInstanceGroups []*emr.InstanceGroup
	EC2Instance       *ec2.Instance
}

type ClusterReqTask struct {
	Cluster            *SparkCluster
	NeedInstanceGroups bool
	ResChan            chan *ClusterReqTaskResponse
	RateLimitChan      *chan int
}

type QuotaMetricUsage struct {
	UserID   int         `json:"uid"`
	OwnerUID int         `json:"owner_uid"`
	Metrics  interface{} `json:"metrics"`
}

type StatisticMonthlyMetricUsage struct {
	UserID   int         `json:"uid"`
	OwnerUID int         `json:"owner_uid"`
	Month    string      `json:"month"`
	Metrics  interface{} `json:"metrics"`
}

type StatisticClusterOperationHistory struct {
	UserID           int         `json:"uid"`
	OwnerUID         int         `json:"owner_uid"`
	Timestamp        string      `json:"timestamp"`
	Offset           int         `json:"offset"`
	Size             int         `json:"size"`
	OperationHistory interface{} `json:"operation_history"`
}

type MapRequestHubSpotResult struct {
	Requests []RequestVectorData `json:"requests"`
	After    string              `json:"after"`
}

type RasterLayer struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Watermark string `json:"watermark"`
}

type RasterLayerDetails struct {
	RasterLayerID   int             `json:"raster_layer_id"`
	RasterLayerName string          `json:"raster_layer_name"`
	Title           string          `json:"title"`
	Watermark       string          `json:"watermark"`
	URL             string          `json:"url"`
	FileID          *int            `json:"file_id,omitempty"`
	FilePath        *string         `json:"file_path,omitempty"`
	UserID          *int            `json:"uid"`
	IsDeleted       *bool           `json:"is_deleted"`
	RasterType      *string         `json:"raster_type"`
	Center          *CenterLatLng   `json:"center,omitempty"`
	Envelope        *EnvelopeLatLng `json:"envelope,omitempty"`
	S3Source        *string         `json:"s3source"`
	TsCreated       *string         `json:"ts_created"`
	// shared layer
	IsShared                 *bool                     `json:"is_shared,omitempty"`
	HasAccess                *bool                     `json:"has_access,omitempty"`
	RasterLayerAccessDetails *RasterLayerAccessDetails `json:"raster_layer_access_details,omitempty"`
}

type RasterLayerAccessDetails struct {
	RasterLayerID   int    `json:"raster_layer_id"`
	RasterLayerName string `json:"raster_layer_name"`
	UserID          int    `json:"uid"`
	FileID          int    `json:"file_id,omitempty"`
	FilePath        string `json:"file_path,omitempty"`
	TsCreated       string `json:"ts_created"`
	Mode            string `json:"mode"`
	Auth            int    `json:"auth"`
}

type LayerCreateIndexRequestMessage struct {
	Version   int    `json:"version"`
	RequestID string `json:"request_id"`
	LayerID   int    `json:"layer_id"`
	LayerUID  int    `json:"layer_uid"`
	DataTable string `json:"data_table"`
	FieldName string `json:"field_name"`
	UserID    int    `json:"uid"`
	TsCreated string `json:"ts_created"`
}
type LayerCreateIndexRequest struct {
	ID        int    `json:"id"`
	LayerID   int    `json:"layer_id"`
	DataTable string `json:"data_table"`
	FieldName string `json:"field_name"`
	UserID    int    `json:"uid"`
	Progress  int    `json:"progress"`
	TsCreated string `json:"ts_created"`
	ErrCode   int    `json:"errcode"`
	ErrMesg   string `json:"errmesg"`
}

type TableDefault struct {
	SchemaName string `json:"schema_name"`
	TableName  string `json:"table_name"`
}

type TableRowCount struct {
	TableDefault
	RowCount int `json:"row_count"`
}

type LayerTableFieldDetails struct {
	FieldName          string                   `json:"field_name"`
	FieldType          string                   `json:"field_type"`
	IndexName          *string                  `json:"index_name,omitempty"`
	IndexType          *string                  `json:"index_type,omitempty"`
	CreateIndexStatus  *string                  `json:"create_index_status,omitempty"`
	CreateIndexRequest *LayerCreateIndexRequest `json:"create_index_request,omitempty"`
}

type LayerCreateRequest struct {
	UserID              int                 `json:"uid"`
	Geofile             string              `json:"geofile"`
	GeofileETag         string              `json:"geofile_etag"`
	FilePath            string              `json:"file_path"`
	TableName           string              `json:"table_name"`
	IsLargeFeatureSplit bool                `json:"is_large_feature_split"`
	Catetitles          []CatetitleGeomType `json:"catetitles"`
	Envelope            EnvelopeLatLng      `json:"envelope"`
	Center              CenterLatLng        `json:"center"`
}

type LayerDetails struct {
	LayerName           string          `json:"layer_name"`
	LayerID             int             `json:"layer_id"`
	FileID              *int            `json:"file_id,omitempty"`
	FilePath            *string         `json:"file_path,omitempty"`
	UserID              *int            `json:"uid,omitempty"`
	Envelope            *EnvelopeLatLng `json:"envelope,omitempty"`
	Center              *CenterLatLng   `json:"center,omitempty"`
	Catetitles          interface{}     `json:"catetitles,omitempty"`
	CatetitlesDefault   []CatetitleInfo `json:"catetitles_default,omitempty"`
	CatetitlesCustomize []CatetitleInfo `json:"catetitles_customize,omitempty"` // relate to project
	IsLargeFeatureSplit *bool           `json:"is_large_feature_split,omitempty"`
	GeoFile             *string         `json:"geofile,omitempty"`
	GeoFileETag         *string         `json:"geofile_etag,omitempty"`
	DataTable           *string         `json:"data_table,omitempty"`
	FeatureCount        *int            `json:"feature_count,omitempty"`
	DataView            *string         `json:"data_view,omitempty"`
	DataViewStatement   *string         `json:"data_view_statement,omitempty"`
	CatetitleField      *string         `json:"catetitle_field,omitempty"` // relate to project
	IsDeleted           *bool           `json:"is_deleted,omitempty"`
	SortIndex           *int            `json:"sort_index,omitempty"` // sort-index for project
	// shared layer
	IsShared           *bool               `json:"is_shared,omitempty"`
	HasAccess          *bool               `json:"has_access,omitempty"`
	LayerAccessDetails *LayerAccessDetails `json:"layer_access_details,omitempty"`
}

type LayerAccessDetails struct {
	LayerID           int          `json:"layer_id"`
	LayerName         string       `json:"layer_name"` //come from shared_layer
	UserID            int          `json:"uid"`
	FileID            int          `json:"file_id"`
	FilePath          string       `json:"file_path"`
	AccessID          string       `json:"access_id"`
	TsCreated         string       `json:"ts_created"`
	IsAOILimit        bool         `json:"is_aoi_limit"`
	AOIIDs            []int        `json:"aoi_ids"`
	AOIMergedGeom     string       `json:"aoi_merged_geom"`
	AOIMergedGeomLen  int          `json:"aoi_merged_geom_len"`
	Center            CenterLatLng `json:"center"` //come from shared_layer
	FeatureCount      int          `json:"feature_count"`
	DataView          string       `json:"data_view"`
	DataViewStatement string       `json:"data_vew_statement"`
	Mode              string       `json:"mode"`
	Auth              int          `json:"auth"`
}

type LayerRequest struct {
	Progress  int        `json:"progress"`
	GeoFile   string     `json:"geofile"`
	LayerFile string     `json:"layerfile"`
	RequestID int64      `json:"request_id"`
	TsCreated *time.Time `json:"ts_created"`
	Status    string     `json:"status"`
	ErrCode   int        `json:"errcode"`
	ErrMesg   string     `json:"errmesg"`
}

type LayerAccessModify struct {
	IsDataViewModified         bool
	IsDataViewHasValidGeometry bool
	LayerAccessDetails
}

type ReviewTaskDetails struct {
	ReviewID   string        `json:"review_id"`
	TaskID     int           `json:"task_id"`
	TaskStatus string        `json:"task_status"`
	UserID     int           `json:"uid"`
	Geom       *string       `json:"geom,omitempty"`
	Center     *CenterLatLng `json:"center,omitempty"`
	TsAcquired *string       `json:"ts_acquired,omitempty"`
}

type ReviewCommentDetails struct {
	ReviewID  string  `json:"review_id"`
	CommentID int     `json:"comment_id"`
	TaskID    int     `json:"task_id"`
	UserID    int     `json:"uid"`
	Geom      *string `json:"geom,omitempty"`
	Comment   *string `json:"comment,omitempty"`
	TsUpdated *string `json:"ts_updated,omitempty"`
}

type ReviewProjectMountDetails struct {
	ID                int    `json:"id"`
	ReviewID          string `json:"review_id"`
	ProjectID         string `json:"project_id"`
	UserID            int    `json:"uid"`
	QueryViewName     string `json:"query_view_name"`
	QueryStatement    string `json:"query_statement"`
	QueryStatementMD5 string `json:"query_statement_md5"`
	QueryStatementLen int    `json:"query_statement_len"`
	TsCreated         string `json:"ts_created"`
}

type ReviewDetails struct {
	ReviewID       string          `json:"review_id"`
	ReviewName     string          `json:"review_name"`
	ProjectID      string          `json:"project_id"`
	UserID         int             `json:"uid"`
	Description    string          `json:"description"`
	TaskLayerID    int             `json:"task_layer_id"`
	CommentLayerID int             `json:"comment_layer_id"`
	ReviewStatus   string          `json:"review_status"`
	TsCreated      string          `json:"ts_created"`
	ProjectDetails *ProjectDetails `json:"project_details,omitempty"`
	// task layer and comment layer
	TaskLayerDetails    *LayerDetails `json:"task_layer_details,omitempty"`
	CommentLayerDetails *LayerDetails `json:"comment_layer_details,omitempty"`

	// active task
	ActiveTasks []*ReviewTaskDetails `json:"active_tasks"`

	// review project mount
	ReviewProjectMountDetails *ReviewProjectMountDetails `json:"review_project_mount_details,omitempty"`
}

type ReviewListOutput struct {
	Reviews   []*ReviewDetails `json:"reviews"`
	Offset    int              `json:"offset"`
	Size      int              `json:"size"`
	TotalSize int              `json:"total_size"`
}

type ProjectCustomizeTools struct {
	GeoCodingAllow    bool            `json:"allow_geocoding"`
	GeoCodingSupport  *bool           `json:"geocoding_support,omitempty"`
	GeoCodingEnvelope *EnvelopeLatLng `json:"geocoding_envelope,omitempty"`
}

type ProjectDetails struct {
	ProjectName        string                 `json:"project_name"`
	ProjectID          string                 `json:"project_id"`
	FileID             *int                   `json:"file_id,omitempty"`
	FilePath           *string                `json:"file_path,omitempty"`
	UserID             *int                   `json:"uid,omitempty"`
	Description        *string                `json:"description,omitempty"`
	IsSaved            *bool                  `json:"is_saved,omitempty"`
	AllowEdit          *bool                  `json:"allow_edit,omitempty"`
	CustomizeTools     *ProjectCustomizeTools `json:"customize_tools,omitempty"`
	RasterLayers       *[]RasterLayerDetails  `json:"raster_layers,omitempty"`
	LayerList          *[]LayerDetails        `json:"layer_list,omitempty"`
	SortedLayerIDs     *[]int                 `json:"sorted_layer_ids,omitempty"`
	TotalLayerRowCount *int                   `json:"total_layer_row_count,omitempty"`
	TsCreated          *string                `json:"ts_created,omitempty"`
	// Catetitles         interface{}            `json:"catetitles,omitempty"`

	// shared project
	IsShared             *bool                 `json:"is_shared,omitempty"`
	HasAccess            *bool                 `json:"has_access,omitempty"`
	ProjectAccessDetails *ProjectAccessDetails `json:"project_access_details,omitempty"`
}

type ProjectAccessDetails struct {
	ProjectID       string `json:"project_id"`
	UserID          int    `json:"uid"`
	FileID          int    `json:"file_id"`
	FilePath        string `json:"file_path"`
	AccessID        string `json:"access_id"`
	ContentAccessID string `json:"content_access_id"`
	TsCreated       string `json:"ts_created"`
	Mode            string `json:"mode"`
	Auth            int    `json:"auth"`
}

type ProjectListOutput struct {
	SharedProjects  []ProjectDetails `json:"shared_projects"`
	PrivateProjects []ProjectDetails `json:"private_projects"`
}

type MountDetails struct {
	MountID           int             `json:"id"`
	MountName         string          `json:"name"`
	ProjectID         string          `json:"project_id"`
	UserID            int             `json:"uid"`
	QueryViewName     string          `json:"query_view_name"`
	QueryStatement    string          `json:"query_statement"`
	QueryStatementMD5 string          `json:"query_statement_md5"`
	QueryStatementLen int             `json:"query_statement_len"`
	ProjectDetails    *ProjectDetails `json:"project_details,omitempty"`
}
type MountedProjectDetails struct {
	MountID           int     `json:"id"`
	QueryStatementLen *int    `json:"query_statement_len,omitempty"`
	QueryStatementMD5 *string `json:"query_statement_md5,omitempty"`
	TsUpdated         string  `json:"open_time"`
	ProjectDetails
}
type MountedProjectListOutput struct {
	GlobalRasterLayers []*RasterLayer           `json:"global_raster_layers"`
	MountedProjects    []*MountedProjectDetails `json:"opened_projects"`
}

type LayerListOutput struct {
	SharedLayers  []LayerDetails `json:"shared_layers"`
	PrivateLayers []LayerDetails `json:"private_layers"`
}

type ContentPurposeDetails struct {
	PurposeID   int    `json:"purpose_id"`
	PurposeName string `json:"purpose_name"`
	UserID      int    `json:"uid"`
}

type ContentPreviewDetails struct {
	PreviewMountID    int    `json:"id"`
	PreviewName       string `json:"preview_name"`
	ContentID         string `json:"content_id"`
	ProjectID         string `json:"project_id"`
	UserID            int    `json:"uid"`
	QueryViewName     string `json:"query_view_name"`
	QueryStatement    string `json:"query_statement"`
	QueryStatementMD5 string `json:"query_statement_md5"`
	QueryStatementLen int    `json:"query_statement_len"`
	TsCreated         string `json:"ts_created"`

	// Project
	ProjectDetails *ProjectDetails `json:"project_details,omitempty"`
}

type ContentDetails struct {
	ContentID           string                   `json:"content_id"`
	ContentName         string                   `json:"content_name"`
	ContentType         string                   `json:"content_type"`
	PreviewImageURL     string                   `json:"preview_image_url"`
	Description         string                   `json:"description"`
	TsCreated           string                   `json:"ts_created"`
	AllowShare          *bool                    `json:"allow_share"`
	HasAccess           *bool                    `json:"has_access,omitempty"`
	UserID              *int                     `json:"uid,omitempty"`
	UserName            *string                  `json:"owner_user_name,omitempty"`
	UserAvatarSrc       *string                  `json:"avatar_src,omitempty"`
	PublicLevel         *int                     `json:"public_level,omitempty"`
	SharedCount         *int                     `json:"shared_count,omitempty"`
	PurposeList         *[]ContentPurposeDetails `json:"purposes,omitempty"`
	Metadata            *string                  `json:"metadata,omitempty"`
	ProjectID           *string                  `json:"project_id,omitempty"`
	ProjectDetails      *ProjectDetails          `json:"project_details,omitempty"`
	TermsOfUse          *string                  `json:"terms_of_use,omitempty"`
	SortPriority        *int                     `json:"sort_priority,omitempty"`
	CommunityVisibility *string                  `json:"community_visibility,omitempty"`
	CommunityIDs        *[]string                `json:"community_ids,omitempty"`
	UserDetails         *UserDetails             `json:"user_details,omitempty"`

	// Preview
	ContentPreviewDetails *ContentPreviewDetails `json:"content_preview_details,omitempty"`
}

type ContentListOutput struct {
	MarketContents  []ContentDetails `json:"market_contents"`
	SharedContents  []ContentDetails `json:"shared_contents"`
	PrivateContents []ContentDetails `json:"private_contents"`
}

type ContentDownloadFile struct {
	FileID      uint64  `json:"file_id"`
	FileName    string  `json:"file_name"`
	FileSize    uint64  `json:"file_size"`
	DataType    string  `json:"data_type"`
	DownloadURL string  `json:"download_url"`
	GeoFile     *string `json:"geofile,omitempty"`
	GeoFileETag *string `json:"geofile_etag,omitempty"`
}

type ContentDownloadFiles struct {
	RefreshTime      string                `json:"refresh_time"`
	ValidDurationSec int64                 `json:"valid_duration_sec"`
	DownloadFiles    []ContentDownloadFile `json:"files"`
}

type ContentAccessDetails struct {
	UserID     int    `json:"uid"`
	UserName   string `json:"user_name"`
	NickName   string `json:"nickname,omitempty"`
	UserEmail  string `json:"email"`
	ContentID  string `json:"content_id,omitempty"`
	AccessID   string `json:"access_id,omitempty"`
	Disable    *bool  `json:"disable,omitempty"`
	SharedBy   *int   `json:"shared_by,omitempty"`
	DisabledBy *int   `json:"disabled_by,omitempty"`
}

type ContentShareInfo struct {
	ContentAccesses []ContentAccessDetails `json:"accesses"`
	AvailableUsers  []ContentAccessDetails `json:"available_users"`
}

type ContentRequestExtendDetails struct {
	RequestID          *string `json:"request_id"`
	AOIID              *int    `json:"aoi_id"`
	Description        *string `json:"description,omitempty"`
	SubmittedTimeStamp *string `json:"submitted_ts,omitempty"`
	UserID             *int    `json:"uid,omitempty"`
	ContentOwnerUID    *int    `json:"content_owner_uid,omitempty"`
	ContentOwnerEmail  *string `json:"content_owner_email,omitempty"`
}

type ContentRequestExtend struct {
	SubmissionID *int `json:"submission_id,omitempty"`
}

type ContentRequestDetails struct {
	RequestID          string                      `json:"request_id"`
	RequestType        string                      `json:"request_type"`
	ContentName        string                      `json:"content_name"`
	ContentID          string                      `json:"content_id"`
	ContentType        *string                     `json:"content_type"`
	UserName           string                      `json:"user_name"`
	Email              string                      `json:"email"`
	Description        string                      `json:"description"`
	Status             string                      `json:"status"`
	TsCreated          string                      `json:"ts_created"`
	TsLastModified     string                      `json:"ts_last_modified"`
	TypeFilter         *string                     `json:"type_filter,omitempty"`
	HasAOI             bool                        `json:"has_aoi"`
	AOIID              int                         `json:"aoi_id"`
	AOIFile            *string                     `json:"aoi_file,omitempty"`
	AOIGeom            *string                     `json:"aoi_geom,omitempty"`
	AOIGeomLen         int                         `json:"aoi_geom_len"`
	AOIGeomSimplifyLen int                         `json:"aoi_geom_simplify_len"`
	Purpose            string                      `json:"purpose,omitempty"`
	Extend             ContentRequestExtend        `json:"extend"`
	ExtendDetails      ContentRequestExtendDetails `json:"details"`
}

type ContentRequestListOutput struct {
	ContentRequests []ContentRequestDetails `json:"content_requests"`
	UnhandledCount  int                     `json:"unhandled_count"`
	Offset          int                     `json:"offset"`
	Size            int                     `json:"size"`
	TotalSize       int                     `json:"total_size"`
}

type ContentManager struct {
	UserID    int    `json:"uid"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	NickName  string `json:"nickname"`
	AvatarSrc string `json:"avatar_src"`
}

type FileExtend struct {
	ProjectID     *string `json:"project_id,omitempty"`
	LayerID       *int    `json:"layer_id,omitempty"`
	RasterLayerID *int    `json:"raster_layer_id,omitempty"`
}

type FileDetails struct {
	FileID         int        `json:"file_id,omitempty"`
	UserID         int        `json:"uid"`
	FileName       string     `json:"file_name"`
	ParentPath     string     `json:"parent_path"`
	FilePath       string     `json:"file_path"`
	FileType       string     `json:"file_type"`
	IsShared       bool       `json:"is_shared"`
	StorageType    string     `json:"storage_type"`
	TsCreated      string     `json:"ts_created"`
	TsLastModified string     `json:"ts_last_modified"`
	Extend         FileExtend `json:"extend"`
	ProcessStatus  int        `json:"process_status"`
	// details for project or layer
	ProjectDetails           *ProjectDetails           `json:"project_details,omitempty"`
	ProjectAccessDetails     *ProjectAccessDetails     `json:"project_access_details,omitempty"`
	LayerDetails             *LayerDetails             `json:"layer_details,omitempty"`
	LayerAccessDetails       *LayerAccessDetails       `json:"layer_access_details,omitempty"`
	RasterLayerDetails       *RasterLayerDetails       `json:"raster_layer_details,omitempty"`
	RasterLayerAccessDetails *RasterLayerAccessDetails `json:"raster_layer_access_details,omitempty"`
}

type FileAccessUser struct {
	UserID    int     `json:"uid"`
	Email     string  `json:"email"`
	AvatarSrc string  `json:"avatar_src"`
	UserName  string  `json:"user_name"`
	NickName  *string `json:"nickname,omitempty"`
}
type FileAccessDetails struct {
	FileAccessUser
	Mode                 string                `json:"mode"`
	Auth                 int                   `json:"auth"`
	FileID               *int                  `json:"file_id,omitempty"`
	FileType             *string               `json:"file_type,omitempty"`
	Extend               *FileExtend           `json:"extend,omitempty"`
	ProjectAccessDetails *ProjectAccessDetails `json:"project_access_details,omitempty"`
	LayerAccessDetails   *LayerAccessDetails   `json:"layer_access_details,omitempty"`
}

type FsFileOutput struct {
	FileID         *int   `json:"file_id,omitempty"`
	FilePath       string `json:"file_path"`
	FileSize       int64  `json:"file_size"`
	IsShared       bool   `json:"is_shared"`
	StorageType    string `json:"storage_type"` // s3|database
	TsLastModified string `json:"ts_last_modified"`
}

type FsAsyncOperationStatusStat struct {
	Status    string `json:"status"`
	TaskCount int    `json:"task_count"`
}

type FsListAsyncOperationsOutput struct {
	OperationID string                        `json:"operation_id"` // transfer int64 to string for gin missing precision
	Operation   string                        `json:"operation"`
	StatusStat  []*FsAsyncOperationStatusStat `json:"status_statistics"`
}

type FsListObjectsOutput struct {
	Prefix                  string         `json:"prefix"`
	Files                   []FsFileOutput `json:"files"`
	Folders                 []string       `json:"folders"`
	IsTruncated             bool           `json:"is_truncated"`
	S3NextContinuationToken *string        `json:"s3_next_continuation_token,omitempty"`
	DBLastItem              *string        `json:"db_last_item,omitempty"`
	KeyCount                int            `json:"key_count"`
	// Async operation
	AsyncOperation *FsListAsyncOperationsOutput `json:"async_operation,omitempty"`
}

type FileCopyTaskResult struct {
	ErrCode     int
	ErrMesg     string
	IsAsyncCopy bool
	FileCopy    *FileCopyInput
	Index       int
}

type FileCopyTask struct {
	FileCopy          *FileCopyInput
	ResultChan        chan *FileCopyTaskResult
	ParallelLimitChan *chan int
	Index             int
}

type FileMoveTaskResult struct {
	ErrCode     int
	ErrMesg     string
	IsAsyncMove bool
	FileMove    *FileMoveInput
	Index       int
}

type FileMoveTask struct {
	FileMove          *FileMoveInput
	ResultChan        chan *FileMoveTaskResult
	ParallelLimitChan *chan int
	Index             int
}

type FileDeleteTaskResult struct {
	ErrCode       int
	ErrMesg       string
	IsAsyncDelete bool
	FileDelete    *FileDeleteInput
	Index         int
}

type FileDeleteTask struct {
	FileDelete        *FileDeleteInput
	ResultChan        chan *FileDeleteTaskResult
	ParallelLimitChan *chan int
	Index             int
}

type FileOperationTaskDetails struct {
	TaskID         string `json:"task_id"`
	OperationID    int64  `json:"operation_id"`
	Operation      string `json:"operation"`
	UserID         int64  `json:"uid"`
	InodeID        int64  `json:"inode_id"`
	FileType       string `json:"file_type"`
	SourcePath     string `json:"source_path"`
	TargetPath     string `json:"target_path"`
	SourceDir      string `json:"source_dir"`
	TargetDir      string `json:"target_dir"`
	TaskMessage    string `json:"task_message"`
	IsSyncProcess  bool   `json:"is_sync_process"`
	Status         string `json:"status"`
	Progress       int    `json:"progress"`
	ErrCode        int    `json:"errcode"`
	ErrMesg        string `json:"errmesg"`
	TsLastModified string `json:"ts_last_modified"`
	TsCreated      string `json:"ts_created"`
}

type FileCopyInput struct {
	SourceEXFSPath string `json:"source_exfs_path" binding:"required"`
	TargetEXFSDir  string `json:"target_exfs_dir" binding:"required"`
	FileSize       int    `json:"file_size" binding:"required"`
	StorageType    string `json:"storage_type" binding:"required"`
	// for async process
	SyncProcessMaxFileSize int64
	// for operation task
	Task *FileOperationTaskDetails `json:"file_operation_task_details,omitempty"`
}

type FileMoveInput struct {
	SourceEXFSPath string `json:"source_exfs_path" binding:"required"`
	TargetEXFSPath string `json:"target_exfs_path" binding:"required"`
	FileSize       int    `json:"file_size" binding:"required"`
	StorageType    string `json:"storage_type" binding:"required"`
	// for async process
	SyncProcessMaxFileSize int64
	// for operation task
	Task *FileOperationTaskDetails `json:"file_operation_task_details,omitempty"`
}

type FileDeleteInput struct {
	EXFSPath    string `json:"exfs_path" binding:"required"`
	StorageType string `json:"storage_type" binding:"required"`
	// for operation task
	Task *FileOperationTaskDetails `json:"file_operation_task_details,omitempty"`
}

type FileDeleteAccessInput struct {
	UserID int `json:"uid"`
}
type FileNewAccessInput struct {
	UserID    int    `json:"uid" binding:"required"`
	Mode      string `json:"mode" binding:"required"`
	Auth      int    `json:"auth"`
	TargetDir string `json:"target_dir,omitempty"`
}
type FileShareUpdateInput struct {
	EXFSPath     string                  `json:"exfs_path" binding:"required"`
	DeleteAccess []FileDeleteAccessInput `json:"delete_access"`
	NewAccess    []FileNewAccessInput    `json:"new_access"`
	StorageType  string                  `json:"storage_type" binding:"required"`
}

type FailedTask struct {
	FilePath string `json:"file_path"`
	Reason   string `json:"reason"`
}

type DirShareDetails struct {
	Path             string       `json:"path"`
	TsCreated        *time.Time   `json:"ts_created"`
	TsFinished       *time.Time   `json:"ts_finished"`
	TaskID           int64        `json:"task_id"`
	ParentTaskID     int64        `json:"parent_task_id"`
	AccessorUID      int          `json:"uid"`
	UserName         string       `json:"user_name"`
	Email            string       `json:"email"`
	AvatarSrc        string       `json:"avatar_src"`
	TargetPath       string       `json:"target_path"`
	Operation        string       `json:"operation"`
	TotalFileCount   int          `json:"total_file_count"`
	SuccessFileCount int          `json:"success_file_count"`
	FailedFileCount  int          `json:"failed_file_count"`
	Status           string       `json:"status"`
	FailedList       []FailedTask `json:"failed_task"`
}
type FileShareInfo struct {
	FileDetails     *FileDetails         `json:"file_details"`
	DirShareDetails []*DirShareDetails   `json:"dir_share_details"`
	S3ObjectDetails *S3ObjectDetails     `json:"s3_object_details"`
	FileAccesses    []*FileAccessDetails `json:"file_accesses"`
	AvailableUsers  []*FileAccessUser    `json:"available_users"`
}

type InferDBFileNameInput struct {
	NewFileName string // input
	FitFileName string // output
}

type InferDBFileNameMapInput struct {
	UserID     int
	ParentPath string
	InferMap   map[int]*InferDBFileNameInput // map<layer_id, {new_filename, fit_filename}>
}

type S3ObjectDetails struct {
	S3ObjectID int    `json:"oid"`
	UserID     int    `json:"uid"`
	S3Key      string `json:"s3_key"`
	S3ETag     string `json:"s3_etag"`
	IsDeleted  bool   `json:"is_deleted"`
	TsCreated  string `json:"ts_created"`
	// shared object
	IsShared              *bool                  `json:"is_shared,omitempty"`
	HasAccess             *bool                  `json:"has_access,omitempty"`
	S3ObjectAccessDetails *S3ObjectAccessDetails `json:"s3_object_access_details,omitempty"`
}

type S3ObjectAccessDetails struct {
	S3ObjectID int    `json:"oid"`
	UserID     int    `json:"uid"`
	S3Key      string `json:"s3_key"`
	TsCreated  string `json:"ts_created"`
	Mode       string `json:"mode"`
	Auth       int    `json:"auth"`
}

type Application struct {
	AppID       string `json:"app_id"`
	AppName     string `json:"app_name"`
	AppIconURL  string `json:"app_icon_url"`
	Description string `json:"description"`
	TsCreated   string `json:"ts_created"`
}

type ReviewCommentTemplate struct {
	TemplateID int    `json:"template_id"`
	ReviewID   string `json:"review_id"`
	UserID     int    `json:"uid"`
	Hotkey     string `json:"hotkey"`
	Comment    string `json:"comment"`
	IsGlobal   bool   `json:"is_global"`
}

type ReviewTask struct {
	ReviewID   string        `json:"review_id"`
	TaskID     int           `json:"task_id"`
	TaskStatus *string       `json:"task_status,omitempty"`
	Geom       string        `json:"geom,omitempty"`
	Center     *CenterLatLng `json:"center,omitempty"`
	TsAcquired *time.Time    `json:"ts_acquired,omitempty"`
}

type ReviewComment struct {
	ReviewID  string     `json:"review_id"`
	CommentID string     `json:"comment_id"`
	TaskID    int        `json:"task_id"`
	Geom      string     `json:"geom,omitempty"`
	Comment   string     `json:"comment,omitempty"`
	TsUpdated *time.Time `json:"ts_updated,omitempty"`
}
