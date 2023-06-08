package models

type School struct {
	ID              DbId `json:"id"`
	SchoolId        `json:"school_id"`
	RegionId        `json:"region_id"`
	SchoolName      string `json:"school_name"`
	RegionName      string `json:"region_name"`
	SchoolEmailOnly bool   `json:"school_email_only"`
	SchoolEmail     string `json:"school_email"`
}

type SchoolId string
type RegionId string
