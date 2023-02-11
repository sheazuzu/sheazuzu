package entity

import "github.com/jinzhu/gorm"

type MatchData struct {
	AdditionalInformation AdditionalInformation `gorm:"foreignKey:additional;association_foreignKey:id"`
	AwayTeam              string
	Date                  string
	HomeTeam              string
	Id                    int `gorm:"column:id;primary_key:yes"`
	MatchType             string
	Result                string
}

type AdditionalInformation struct {
	gorm.Model
	Additional  string
	Information string
}
