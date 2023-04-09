package models

// CafeteriaMenu struct represents a cafeteria menu
type CafeteriaMenu struct {
	ID       DbId `json:"id"`
	SchoolId `json:"school_id"`
	MealName string      `json:"meal_name"`
	Date     string      `json:"date"`
	Items    []MenuEntry `json:"items"`
}

// MenuEntry struct represents an item in a cafeteria menu
type MenuEntry struct {
	Name      string        `json:"name"`
	Allergies []AllergyType `json:"allergy,omitempty"`
	Contents  string        `json:"contents"`
}

type AllergyType int8

var (
	난류   = 1
	우유   = 2
	메밀   = 3
	땅콩   = 4
	대두   = 5
	밀    = 6
	고등어  = 7
	게    = 8
	새우   = 9
	돼지고기 = 10
	복숭아  = 11
	토마토  = 12
	아황산염 = 13
	호두   = 14
	닭고기  = 15
	쇠고기  = 16
	오징어  = 17
	조개류  = 18
)
