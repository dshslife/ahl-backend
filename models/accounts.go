package models

import "errors"

// Account struct represents a user
type Account struct {
	DbId           `json:"id"`
	UserId         `json:"user_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       []byte
	PermissionInfo `json:"permission"`
}

type SqlAccount struct {
	DbId
	UserId
	Name     string
	Email    string
	Password []byte
	PermissionLevel
	SchoolId
	Timetable
	Grade       int
	Class       int
	Number      int
	ChecklistId DbId
	Friends     []UserId
}

func (sqlAccount SqlAccount) Finalize() (Account, error) {
	var info PermissionInfo
	var err error
	switch sqlAccount.PermissionLevel {
	case UNKNOWN:
		info = Unknown{}
	case STUDENT:
		info = StudentInfo{
			SchoolId:    sqlAccount.SchoolId,
			Timetable:   sqlAccount.Timetable,
			Grade:       sqlAccount.Grade,
			Class:       sqlAccount.Class,
			Number:      sqlAccount.Number,
			ChecklistId: sqlAccount.ChecklistId,
			Friends:     sqlAccount.Friends,
		}
	case TEACHER:
		info = TeacherInfo{
			SchoolId: sqlAccount.SchoolId,
		}
	case ADMIN:
		info = AdminInfo{}
	default:
		err = errors.New("unknown Permission Level")
	}
	if err != nil {
		return Account{}, err
	}
	return Account{
		DbId:           sqlAccount.DbId,
		UserId:         sqlAccount.UserId,
		Name:           sqlAccount.Name,
		Email:          sqlAccount.Email,
		Password:       sqlAccount.Password,
		PermissionInfo: info,
	}, nil
}

func (account Account) ToSql() (SqlAccount, error) {
	var toReturn SqlAccount
	var err error
	toReturn.DbId = account.DbId
	toReturn.UserId = account.UserId
	toReturn.Name = account.Name
	toReturn.Email = account.Email
	toReturn.Password = account.Password
	toReturn.PermissionLevel = account.PermissionInfo.GetLevel()

	switch account.PermissionInfo.GetLevel() {
	case UNKNOWN:
		break
	case STUDENT:
		info := account.PermissionInfo.(StudentInfo)
		toReturn.SchoolId = info.SchoolId
		toReturn.Timetable = info.Timetable
		toReturn.Grade = info.Grade
		toReturn.Class = info.Class
		toReturn.Number = info.Number
		toReturn.ChecklistId = info.ChecklistId
		toReturn.Friends = info.Friends
	case TEACHER:
		info := account.PermissionInfo.(TeacherInfo)
		toReturn.SchoolId = info.SchoolId
	case ADMIN:
		break
	default:
		err = errors.New("unknown Permission Level")
	}

	if err != nil {
		return SqlAccount{}, err
	}

	return toReturn, nil
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
