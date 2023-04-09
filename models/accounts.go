package models

// Account struct represents a user
// THIS CONTAINS PASSWORD, THIS SHOULD NEVER BE DELIVERED THROUGH NETWORK
// TODO Store password hash instead
type Account struct {
	DbId           `json:"id"`
	UserId         `json:"user_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PermissionInfo `json:"permission"`
}

type PermissionLevel int8
type DbId int
type UserId string // 이거 진짜 쓸라나? DbId로만 퉁칠 수 있긴 한데

// PermissionInfo 유저 권한에 따른 추가 정보, 권한 레벨은 무조건 있어야 함
// StudentInfo, TeacherInfo, AdminInfo, Unknown이 아래 인터페이스를 구현함.
type PermissionInfo interface {
	GetLevel() PermissionLevel
}

var (
	UNKNOWN PermissionLevel = 0
	STUDENT PermissionLevel = 1
	TEACHER PermissionLevel = 2
	ADMIN   PermissionLevel = 3
)

type Unknown struct{}

func (info Unknown) GetLevel() PermissionLevel {
	return UNKNOWN
}

type StudentInfo struct {
	SchoolId    `json:"school_id"`
	Timetable   `json:"timetable"`
	Grade       int      `json:"grade"`  //학년
	Class       int      `json:"class"`  //반
	Number      int      `json:"number"` //번호
	ChecklistId DbId     `json:"checklist_id"`
	Friends     []UserId `json:"friends"`
}

func (info StudentInfo) GetLevel() PermissionLevel {
	return STUDENT
}

type TeacherInfo struct {
	SchoolId `json:"school_id"`
}

func (info TeacherInfo) GetLevel() PermissionLevel {
	return TEACHER
}

type AdminInfo struct{}

func (info AdminInfo) GetLevel() PermissionLevel {
	return ADMIN
}
