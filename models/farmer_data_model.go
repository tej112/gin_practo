package models

type Farmer struct {
	Pk                    string  `json:"pk,omitempty"`
	Name                  string  `json:"name" binding:"required"`
	Contact_num           string  `json:"contact_num" binding:"required"`
	Profile_pic           string  `json:"profile_pic" default:""`
	Who_are_you           string  `json:"who_are_you" binding:"required" default:"Farmer"`
	Language_preference   string  `json:"language_preference" binding:"required" default:"en"`
	Profile_complete_perc float32 `json:"profile_complete_perc,omitempty" default:"80"`
	Created_at            string  `json:"created_at,omitempty"`
	Last_updated_at       string  `json:"last_updated_at,omitempty"`
}

type Hash struct {
	Pk                  string `json:"pk,omitempty"`
	Language_preference string `json:"language_preference,omitempty"`
}
