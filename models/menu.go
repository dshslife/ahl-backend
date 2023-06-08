package models

import "github.com/vishalkuo/bimap"

// CafeteriaMenu struct represents a cafeteria menu
type CafeteriaMenu struct {
	ID       DbId `json:"id"`
	SchoolId `json:"school_id"`
	MealName string `json:"meal_name"`
	Date     string `json:"date"`
	Contents string `json:"items"`
}

type AllergyType int8

var Allergies = bimap.NewBiMap[string, int8]()

func InitAllergies() {
	Allergies.Insert("난류", 1)
	Allergies.Insert("우유", 2)
	Allergies.Insert("메밀", 3)
	Allergies.Insert("땅콩", 4)
	Allergies.Insert("대두", 5)
	Allergies.Insert("밀", 6)
	Allergies.Insert("고등어", 7)
	Allergies.Insert("게", 8)
	Allergies.Insert("새우", 9)
	Allergies.Insert("돼지고", 10)
	Allergies.Insert("복숭아", 11)
	Allergies.Insert("토마토", 12)
	Allergies.Insert("아황산", 13)
	Allergies.Insert("호두", 14)
	Allergies.Insert("닭고기", 15)
	Allergies.Insert("쇠고기", 16)
	Allergies.Insert("오징어", 17)
	Allergies.Insert("조개류", 18)
	Allergies.MakeImmutable()
}
