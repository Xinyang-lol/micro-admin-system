package pb

type Empty struct{}

type IDRequest struct {
	Id int64 `json:"id"`
}

type UserIDRequest struct {
	UserId int64 `json:"user_id"`
}

type IDResponse struct {
	Id int64 `json:"id"`
}

type ListRequest struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Keyword  string `json:"keyword"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

type PermissionRequest struct {
	UserId     int64  `json:"user_id"`
	Permission string `json:"permission"`
}

type PermissionResponse struct {
	Allowed bool `json:"allowed"`
}

type UserInfo struct {
	Id          int64    `json:"id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Status      int32    `json:"status"`
	DeptId      int64    `json:"dept_id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type User struct {
	Id        int64   `json:"id"`
	Username  string  `json:"username"`
	Nickname  string  `json:"nickname"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Status    int32   `json:"status"`
	DeptId    int64   `json:"dept_id"`
	RoleIds   []int64 `json:"role_ids"`
	CreatedAt string  `json:"created_at"`
}

type UserSaveRequest struct {
	Id       int64   `json:"id"`
	Username string  `json:"username"`
	Nickname string  `json:"nickname"`
	Password string  `json:"password"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Status   int32   `json:"status"`
	DeptId   int64   `json:"dept_id"`
	RoleIds  []int64 `json:"role_ids"`
}

type UserStatusRequest struct {
	Id     int64 `json:"id"`
	Status int32 `json:"status"`
}

type UserPasswordRequest struct {
	Id       int64  `json:"id"`
	Password string `json:"password"`
}

type UserRolesRequest struct {
	Id      int64   `json:"id"`
	RoleIds []int64 `json:"role_ids"`
}

type UserListResponse struct {
	Items []*User `json:"items"`
	Total int64   `json:"total"`
}

type Role struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	Code    string  `json:"code"`
	Remark  string  `json:"remark"`
	Status  int32   `json:"status"`
	MenuIds []int64 `json:"menu_ids"`
}

type RoleSaveRequest struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Remark string `json:"remark"`
	Status int32  `json:"status"`
}

type RoleMenusRequest struct {
	Id      int64   `json:"id"`
	MenuIds []int64 `json:"menu_ids"`
}

type RoleListResponse struct {
	Items []*Role `json:"items"`
	Total int64   `json:"total"`
}

type Menu struct {
	Id         int64   `json:"id"`
	ParentId   int64   `json:"parent_id"`
	Name       string  `json:"name"`
	Path       string  `json:"path"`
	Component  string  `json:"component"`
	Permission string  `json:"permission"`
	Icon       string  `json:"icon"`
	Type       int32   `json:"type"`
	Sort       int32   `json:"sort"`
	Children   []*Menu `json:"children"`
}

type MenuSaveRequest struct {
	Id         int64  `json:"id"`
	ParentId   int64  `json:"parent_id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Icon       string `json:"icon"`
	Type       int32  `json:"type"`
	Sort       int32  `json:"sort"`
}

type MenuTreeResponse struct {
	Items []*Menu `json:"items"`
}

type Dept struct {
	Id       int64   `json:"id"`
	ParentId int64   `json:"parent_id"`
	Name     string  `json:"name"`
	Sort     int32   `json:"sort"`
	Status   int32   `json:"status"`
	Children []*Dept `json:"children"`
}

type DeptSaveRequest struct {
	Id       int64  `json:"id"`
	ParentId int64  `json:"parent_id"`
	Name     string `json:"name"`
	Sort     int32  `json:"sort"`
	Status   int32  `json:"status"`
}

type DeptTreeResponse struct {
	Items []*Dept `json:"items"`
}

type DeviceEmpty struct{}

type DeviceIDRequest struct {
	Id int64 `json:"id"`
}

type DeviceIDResponse struct {
	Id int64 `json:"id"`
}

type DeviceListRequest struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Keyword  string `json:"keyword"`
	TypeId   int64  `json:"type_id"`
	Status   string `json:"status"`
}

type Device struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	TypeId    int64  `json:"type_id"`
	TypeName  string `json:"type_name"`
	Status    string `json:"status"`
	Location  string `json:"location"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
}

type DeviceSaveRequest struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	TypeId   int64  `json:"type_id"`
	Status   string `json:"status"`
	Location string `json:"location"`
	Remark   string `json:"remark"`
}

type DeviceListResponse struct {
	Items []*Device `json:"items"`
	Total int64     `json:"total"`
}

type DeviceStatisticsResponse struct {
	Total     int64            `json:"total"`
	Online    int64            `json:"online"`
	Offline   int64            `json:"offline"`
	Repair    int64            `json:"repair"`
	TypeStats map[string]int64 `json:"type_stats"`
}

type DeviceType struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Remark string `json:"remark"`
}

type DeviceTypeSaveRequest struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Remark string `json:"remark"`
}

type DeviceTypeListResponse struct {
	Items []*DeviceType `json:"items"`
}

type FileEmpty struct{}

type FileIDRequest struct {
	Id int64 `json:"id"`
}

type FileIDResponse struct {
	Id int64 `json:"id"`
}

type UploadFileMetaRequest struct {
	OriginalName string `json:"original_name"`
	StoredName   string `json:"stored_name"`
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
	UploaderId   int64  `json:"uploader_id"`
}

type FileListRequest struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Keyword  string `json:"keyword"`
}

type FileInfo struct {
	Id           int64  `json:"id"`
	OriginalName string `json:"original_name"`
	StoredName   string `json:"stored_name"`
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	ContentType  string `json:"content_type"`
	UploaderId   int64  `json:"uploader_id"`
	CreatedAt    string `json:"created_at"`
}

type FileListResponse struct {
	Items []*FileInfo `json:"items"`
	Total int64       `json:"total"`
}
