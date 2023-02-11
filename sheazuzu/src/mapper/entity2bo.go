package mapper

import (
	"sheazuzu/common/src/utils"
	"sheazuzu/sheazuzu/src/entity"
	"sheazuzu/sheazuzu/src/generated/sheazuzu"
)

func MatchDataToBo(data entity.MatchData) sheazuzu.MatchData {
	return sheazuzu.MatchData{
		AdditionalInformations: nil,
		AwayTeam:               utils.ToStringPtr(data.AwayTeam),
		Date:                   utils.ToStringPtr(data.Date),
		HomeTeam:               utils.ToStringPtr(data.HomeTeam),
		Id:                     utils.ToIntPtr(data.Id),
		MatchType:              utils.ToStringPtr(data.MatchType),
		Result:                 utils.ToStringPtr(data.Result),
	}
}
