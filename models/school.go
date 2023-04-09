package models

type School struct {
	ID         DbId `json:"id"`
	SchoolId   `json:"school_id"`
	RegionId   `json:"region_id"`
	SchoolName string `json:"school_name"`
	RegionName string `json:"region_name"`
}

type SchoolId string
type RegionId string
