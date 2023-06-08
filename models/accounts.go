package models

import (
	"errors"
)
import "github.com/google/uuid"

// Account struct represents a user
type Account struct {
	DbId           `json:"id"`
	UserId         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Password       []byte
	PermissionInfo `json:"permission"`
}

type FlatAccount struct {
	DbId
	UserId   []byte
	Name     string
	Email    string
	Password []byte
	PermissionLevel
	SchoolId
	TimeTableEntries  []byte
	TimeTableIsPublic bool
	Grade             int
	Class             int
	Number            int
	ChecklistId       DbId
	Friends           []byte
}

func (flatAccount FlatAccount) Restore() (Account, error) {
	var info PermissionInfo
	var err error

	timetableEntries := Int64ArrayToDbIdArray(BytesToInt64Array(flatAccount.TimeTableEntries))
	friends := BytesToUUIDArray(flatAccount.Friends)

	switch flatAccount.PermissionLevel {
	case STUDENT:
		timetable := Timetable{
			Entries:  timetableEntries,
			IsPublic: flatAccount.TimeTableIsPublic,
		}
		info = StudentInfo{
			SchoolId:    flatAccount.SchoolId,
			Timetable:   timetable,
			Grade:       flatAccount.Grade,
			Class:       flatAccount.Class,
			Number:      flatAccount.Number,
			ChecklistId: flatAccount.ChecklistId,
			Friends:     friends,
		}
	case TEACHER:
		info = TeacherInfo{
			SchoolId: flatAccount.SchoolId,
		}
	case ADMIN:
		info = AdminInfo{}
	default:
		err = errors.New("unknown Permission Level")
	}
	if err != nil {
		return Account{}, err
	}

	userId, err := uuid.FromBytes(flatAccount.UserId)
	if err != nil {
		return Account{}, err
	}
	return Account{
		DbId:           flatAccount.DbId,
		UserId:         userId,
		Name:           flatAccount.Name,
		Email:          flatAccount.Email,
		Password:       flatAccount.Password,
		PermissionInfo: info,
	}, nil
}

func (account Account) ToSql() (FlatAccount, error) {
	var toReturn FlatAccount
	var err error
	toReturn.DbId = account.DbId
	toReturn.UserId = account.UserId[:]
	toReturn.Name = account.Name
	toReturn.Email = account.Email
	toReturn.Password = account.Password
	toReturn.PermissionLevel = account.PermissionInfo.GetLevel()

	switch account.PermissionInfo.GetLevel() {
	case STUDENT:
		info := account.PermissionInfo.(StudentInfo)
		toReturn.SchoolId = info.SchoolId
		toReturn.TimeTableEntries = Int64ArrayToBytes(DbIdArrayToInt64Array(info.Timetable.Entries))
		toReturn.TimeTableIsPublic = info.Timetable.IsPublic
		toReturn.Grade = info.Grade
		toReturn.Class = info.Class
		toReturn.Number = info.Number
		toReturn.ChecklistId = info.ChecklistId
		toReturn.Friends = UuidArrayToBytes(info.Friends)
	case TEACHER:
		info := account.PermissionInfo.(TeacherInfo)
		toReturn.SchoolId = info.SchoolId
	case ADMIN:
		break
	default:
		err = errors.New("unknown Permission Level")
	}

	if err != nil {
		return FlatAccount{}, err
	}

	return toReturn, nil
}

type PermissionLevel int8
type DbId int64

// PermissionInfo 유저 권한에 따른 추가 정보, 권한 레벨은 무조건 있어야 함
// StudentInfo, TeacherInfo, AdminInfo, Unknown이 아래 인터페이스를 구현함.
type PermissionInfo interface {
	GetLevel() PermissionLevel
}

var (
	STUDENT PermissionLevel = 1
	TEACHER PermissionLevel = 2
	ADMIN   PermissionLevel = 3
)

type StudentInfo struct {
	SchoolId    `json:"school_id"`
	Timetable   `json:"timetable"`
	Grade       int         `json:"grade"`  //학년
	Class       int         `json:"class"`  //반
	Number      int         `json:"number"` //번호
	ChecklistId DbId        `json:"checklist_id"`
	Friends     []uuid.UUID `json:"friends"`
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
