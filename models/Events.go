package models

type Events struct {
	ID       DbId `json:"id"`
	SchoolId `json:"school_id"`
	Month    int          `json:"month"`
	Events   []EventEntry `json:"events"`
}

type EventEntry struct {
	Date               string         `json:"date"`
	DateKind           string         `json:"date_kind"` // 대전 대신고 기준 이 필드는 "휴업일", "공휴일", 또는 "해당 없음" 임. 방학도 휴업일임.
	EventName          string         `json:"event_name"`
	EventContents      string         `json:"event_contents"`
	FirstGradeAttends  AttendanceType `json:"first_grade"`
	SecondGradeAttends AttendanceType `json:"second_grade"`
	ThirdGradeAttends  AttendanceType `json:"third_grade"`
	ModifiedDate       string         `json:"modified_date"`
}

type AttendanceType int8

// Neis API는 특정 행사 참여 여부를 Y, *, N으로 구분하는데, Y는 참여, *은 해당 학년 없음, N은 참여 안함
// IGNORED 는 해당 학년이 없음을 의미함
const (
	NO      = -1
	IGNORED = 0
	YES     = 1
)

func (attendanceType AttendanceType) Attends() bool {
	return attendanceType == YES
}
