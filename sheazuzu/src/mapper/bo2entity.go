package mapper

import (
	"sheazuzu/common/src/utils"
	"sheazuzu/sheazuzu/src/entity"
	"sheazuzu/sheazuzu/src/generated/sheazuzu"
)

func BoToMatchData(data sheazuzu.MatchData) entity.MatchData {
	return entity.MatchData{
		AdditionalInformation: entity.AdditionalInformation{},
		AwayTeam:              utils.ToString(data.AwayTeam),
		Date:                  utils.ToString(data.Date),
		HomeTeam:              utils.ToString(data.HomeTeam),
		Id:                    utils.ToInt(data.Id),
		MatchType:             utils.ToString(data.MatchType),
		Result:                utils.ToString(data.Result),
	}
}
