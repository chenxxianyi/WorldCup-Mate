package services

import (
	"errors"
	"fmt"
	"strconv"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type officialVenue struct {
	City    string
	Stadium string
}

type officialCityMeta struct {
	Country  string
	Timezone string
}

// FIFA public calendar API, checked 2026-06-04:
// https://api.fifa.com/api/v3/calendar/matches?language=en&count=200&idCompetition=17&from=2026-06-01&to=2026-07-31
var officialVenueByExternalID = map[int64]officialVenue{
	537327: {City: "Mexico City", Stadium: "Mexico City Stadium"},                       // Match 1
	537328: {City: "Guadalajara", Stadium: "Guadalajara Stadium"},                       // Match 2
	537333: {City: "Toronto", Stadium: "Toronto Stadium"},                               // Match 3
	537345: {City: "Los Angeles", Stadium: "Los Angeles Stadium"},                       // Match 4
	537340: {City: "Boston", Stadium: "Boston Stadium"},                                 // Match 5
	537346: {City: "Vancouver", Stadium: "BC Place Vancouver"},                          // Match 6
	537339: {City: "New Jersey", Stadium: "New York/New Jersey Stadium"},                // Match 7
	537334: {City: "San Francisco Bay Area", Stadium: "San Francisco Bay Area Stadium"}, // Match 8
	537352: {City: "Philadelphia", Stadium: "Philadelphia Stadium"},                     // Match 9
	537351: {City: "Houston", Stadium: "Houston Stadium"},                               // Match 10
	537357: {City: "Dallas", Stadium: "Dallas Stadium"},                                 // Match 11
	537358: {City: "Monterrey", Stadium: "Monterrey Stadium"},                           // Match 12
	537370: {City: "Miami", Stadium: "Miami Stadium"},                                   // Match 13
	537369: {City: "Atlanta", Stadium: "Atlanta Stadium"},                               // Match 14
	537364: {City: "Los Angeles", Stadium: "Los Angeles Stadium"},                       // Match 15
	537363: {City: "Seattle", Stadium: "Seattle Stadium"},                               // Match 16
	537391: {City: "New Jersey", Stadium: "New York/New Jersey Stadium"},                // Match 17
	537392: {City: "Boston", Stadium: "Boston Stadium"},                                 // Match 18
	537397: {City: "Kansas City", Stadium: "Kansas City Stadium"},                       // Match 19
	537398: {City: "San Francisco Bay Area", Stadium: "San Francisco Bay Area Stadium"}, // Match 20
	537410: {City: "Toronto", Stadium: "Toronto Stadium"},                               // Match 21
	537409: {City: "Dallas", Stadium: "Dallas Stadium"},                                 // Match 22
	537403: {City: "Houston", Stadium: "Houston Stadium"},                               // Match 23
	537404: {City: "Mexico City", Stadium: "Mexico City Stadium"},                       // Match 24
	537329: {City: "Atlanta", Stadium: "Atlanta Stadium"},                               // Match 25
	537335: {City: "Los Angeles", Stadium: "Los Angeles Stadium"},                       // Match 26
	537336: {City: "Vancouver", Stadium: "BC Place Vancouver"},                          // Match 27
	537330: {City: "Guadalajara", Stadium: "Guadalajara Stadium"},                       // Match 28
	537341: {City: "Philadelphia", Stadium: "Philadelphia Stadium"},                     // Match 29
	537342: {City: "Boston", Stadium: "Boston Stadium"},                                 // Match 30
	537347: {City: "San Francisco Bay Area", Stadium: "San Francisco Bay Area Stadium"}, // Match 31
	537348: {City: "Seattle", Stadium: "Seattle Stadium"},                               // Match 32
	537353: {City: "Toronto", Stadium: "Toronto Stadium"},                               // Match 33
	537354: {City: "Kansas City", Stadium: "Kansas City Stadium"},                       // Match 34
	537359: {City: "Houston", Stadium: "Houston Stadium"},                               // Match 35
	537360: {City: "Monterrey", Stadium: "Monterrey Stadium"},                           // Match 36
	537372: {City: "Miami", Stadium: "Miami Stadium"},                                   // Match 37
	537371: {City: "Atlanta", Stadium: "Atlanta Stadium"},                               // Match 38
	537365: {City: "Los Angeles", Stadium: "Los Angeles Stadium"},                       // Match 39
	537366: {City: "Vancouver", Stadium: "BC Place Vancouver"},                          // Match 40
	537394: {City: "New Jersey", Stadium: "New York/New Jersey Stadium"},                // Match 41
	537393: {City: "Philadelphia", Stadium: "Philadelphia Stadium"},                     // Match 42
	537399: {City: "Dallas", Stadium: "Dallas Stadium"},                                 // Match 43
	537400: {City: "San Francisco Bay Area", Stadium: "San Francisco Bay Area Stadium"}, // Match 44
	537411: {City: "Boston", Stadium: "Boston Stadium"},                                 // Match 45
	537412: {City: "Toronto", Stadium: "Toronto Stadium"},                               // Match 46
	537405: {City: "Houston", Stadium: "Houston Stadium"},                               // Match 47
	537406: {City: "Guadalajara", Stadium: "Guadalajara Stadium"},                       // Match 48
	537343: {City: "Miami", Stadium: "Miami Stadium"},                                   // Match 49
	537344: {City: "Atlanta", Stadium: "Atlanta Stadium"},                               // Match 50
	537337: {City: "Vancouver", Stadium: "BC Place Vancouver"},                          // Match 51
	537338: {City: "Seattle", Stadium: "Seattle Stadium"},                               // Match 52
	537331: {City: "Mexico City", Stadium: "Mexico City Stadium"},                       // Match 53
	537332: {City: "Monterrey", Stadium: "Monterrey Stadium"},                           // Match 54
	537356: {City: "Philadelphia", Stadium: "Philadelphia Stadium"},                     // Match 55
	537355: {City: "New Jersey", Stadium: "New York/New Jersey Stadium"},                // Match 56
	537362: {City: "Dallas", Stadium: "Dallas Stadium"},                                 // Match 57
	537361: {City: "Kansas City", Stadium: "Kansas City Stadium"},                       // Match 58
	537349: {City: "Los Angeles", Stadium: "Los Angeles Stadium"},                       // Match 59
	537350: {City: "San Francisco Bay Area", Stadium: "San Francisco Bay Area Stadium"}, // Match 60
	537395: {City: "Boston", Stadium: "Boston Stadium"},                                 // Match 61
	537396: {City: "Toronto", Stadium: "Toronto Stadium"},                               // Match 62
	537368: {City: "Seattle", Stadium: "Seattle Stadium"},                               // Match 63
	537367: {City: "Vancouver", Stadium: "BC Place Vancouver"},                          // Match 64
	537374: {City: "Houston", Stadium: "Houston Stadium"},                               // Match 65
	537373: {City: "Guadalajara", Stadium: "Guadalajara Stadium"},                       // Match 66
	537413: {City: "New Jersey", Stadium: "New York/New Jersey Stadium"},                // Match 67
	537414: {City: "Philadelphia", Stadium: "Philadelphia Stadium"},                     // Match 68
	537402: {City: "Kansas City", Stadium: "Kansas City Stadium"},                       // Match 69
	537401: {City: "Dallas", Stadium: "Dallas Stadium"},                                 // Match 70
	537407: {City: "Miami", Stadium: "Miami Stadium"},                                   // Match 71
	537408: {City: "Atlanta", Stadium: "Atlanta Stadium"},                               // Match 72
	537417: {City: "Los Angeles", Stadium: "Los Angeles Stadium"},                       // Match 73
	537415: {City: "Boston", Stadium: "Boston Stadium"},                                 // Match 74
	537418: {City: "Monterrey", Stadium: "Monterrey Stadium"},                           // Match 75
	537423: {City: "Houston", Stadium: "Houston Stadium"},                               // Match 76
	537416: {City: "New Jersey", Stadium: "New York/New Jersey Stadium"},                // Match 77
	537424: {City: "Dallas", Stadium: "Dallas Stadium"},                                 // Match 78
	537425: {City: "Mexico City", Stadium: "Mexico City Stadium"},                       // Match 79
	537426: {City: "Atlanta", Stadium: "Atlanta Stadium"},                               // Match 80
	537421: {City: "San Francisco Bay Area", Stadium: "San Francisco Bay Area Stadium"}, // Match 81
	537422: {City: "Seattle", Stadium: "Seattle Stadium"},                               // Match 82
	537419: {City: "Toronto", Stadium: "Toronto Stadium"},                               // Match 83
	537420: {City: "Los Angeles", Stadium: "Los Angeles Stadium"},                       // Match 84
	537429: {City: "Vancouver", Stadium: "BC Place Vancouver"},                          // Match 85
	537427: {City: "Miami", Stadium: "Miami Stadium"},                                   // Match 86
	537430: {City: "Kansas City", Stadium: "Kansas City Stadium"},                       // Match 87
	537428: {City: "Dallas", Stadium: "Dallas Stadium"},                                 // Match 88
	537375: {City: "Philadelphia", Stadium: "Philadelphia Stadium"},                     // Match 89
	537376: {City: "Houston", Stadium: "Houston Stadium"},                               // Match 90
	537377: {City: "New Jersey", Stadium: "New York/New Jersey Stadium"},                // Match 91
	537378: {City: "Mexico City", Stadium: "Mexico City Stadium"},                       // Match 92
	537379: {City: "Dallas", Stadium: "Dallas Stadium"},                                 // Match 93
	537380: {City: "Seattle", Stadium: "Seattle Stadium"},                               // Match 94
	537381: {City: "Atlanta", Stadium: "Atlanta Stadium"},                               // Match 95
	537382: {City: "Vancouver", Stadium: "BC Place Vancouver"},                          // Match 96
	537383: {City: "Boston", Stadium: "Boston Stadium"},                                 // Match 97
	537384: {City: "Los Angeles", Stadium: "Los Angeles Stadium"},                       // Match 98
	537385: {City: "Miami", Stadium: "Miami Stadium"},                                   // Match 99
	537386: {City: "Kansas City", Stadium: "Kansas City Stadium"},                       // Match 100
	537387: {City: "Dallas", Stadium: "Dallas Stadium"},                                 // Match 101
	537388: {City: "Atlanta", Stadium: "Atlanta Stadium"},                               // Match 102
	537389: {City: "Miami", Stadium: "Miami Stadium"},                                   // Match 103
	537390: {City: "New Jersey", Stadium: "New York/New Jersey Stadium"},                // Match 104
}

var officialCityMetaByName = map[string]officialCityMeta{
	"Atlanta":                {Country: "USA", Timezone: "America/New_York"},
	"Boston":                 {Country: "USA", Timezone: "America/New_York"},
	"Dallas":                 {Country: "USA", Timezone: "America/Chicago"},
	"Guadalajara":            {Country: "Mexico", Timezone: "America/Mexico_City"},
	"Houston":                {Country: "USA", Timezone: "America/Chicago"},
	"Kansas City":            {Country: "USA", Timezone: "America/Chicago"},
	"Los Angeles":            {Country: "USA", Timezone: "America/Los_Angeles"},
	"Mexico City":            {Country: "Mexico", Timezone: "America/Mexico_City"},
	"Miami":                  {Country: "USA", Timezone: "America/New_York"},
	"Monterrey":              {Country: "Mexico", Timezone: "America/Monterrey"},
	"New Jersey":             {Country: "USA", Timezone: "America/New_York"},
	"Philadelphia":           {Country: "USA", Timezone: "America/New_York"},
	"San Francisco Bay Area": {Country: "USA", Timezone: "America/Los_Angeles"},
	"Seattle":                {Country: "USA", Timezone: "America/Los_Angeles"},
	"Toronto":                {Country: "Canada", Timezone: "America/Toronto"},
	"Vancouver":              {Country: "Canada", Timezone: "America/Vancouver"},
}

var officialCityDisplayNameByName = map[string]string{
	"Atlanta":                "亚特兰大",
	"Boston":                 "波士顿",
	"Dallas":                 "达拉斯",
	"Guadalajara":            "瓜达拉哈拉",
	"Houston":                "休斯敦",
	"Kansas City":            "堪萨斯城",
	"Los Angeles":            "洛杉矶",
	"Mexico City":            "墨西哥城",
	"Miami":                  "迈阿密",
	"Monterrey":              "蒙特雷",
	"New Jersey":             "新泽西",
	"Philadelphia":           "费城",
	"San Francisco Bay Area": "旧金山湾区",
	"Seattle":                "西雅图",
	"Toronto":                "多伦多",
	"Vancouver":              "温哥华",
}

var officialStadiumDisplayNameByName = map[string]string{
	"Atlanta Stadium":                "亚特兰大体育场",
	"BC Place Vancouver":             "温哥华 BC Place",
	"Boston Stadium":                 "波士顿体育场",
	"Dallas Stadium":                 "达拉斯体育场",
	"Guadalajara Stadium":            "瓜达拉哈拉体育场",
	"Houston Stadium":                "休斯敦体育场",
	"Kansas City Stadium":            "堪萨斯城体育场",
	"Los Angeles Stadium":            "洛杉矶体育场",
	"Mexico City Stadium":            "墨西哥城体育场",
	"Miami Stadium":                  "迈阿密体育场",
	"Monterrey Stadium":              "蒙特雷体育场",
	"New York/New Jersey Stadium":    "纽约/新泽西体育场",
	"Philadelphia Stadium":           "费城体育场",
	"San Francisco Bay Area Stadium": "旧金山湾区体育场",
	"Seattle Stadium":                "西雅图体育场",
	"Toronto Stadium":                "多伦多体育场",
}

func ensureVenueForExternalMatch(externalID int64) (uint, uint, error) {
	venue, ok := officialVenueByExternalID[externalID]
	if !ok {
		return ensureFallbackVenue()
	}
	return ensureOfficialVenue(venue)
}

func ensureOfficialVenue(venue officialVenue) (uint, uint, error) {
	if venue.City == "" || venue.Stadium == "" {
		return ensureFallbackVenue()
	}

	cityID, err := ensureOfficialCity(venue.City)
	if err != nil {
		return 0, 0, err
	}

	stadiumID, err := upsertOfficialStadium(venue.Stadium, cityID)
	if err != nil {
		return 0, 0, err
	}
	return cityID, stadiumID, nil
}

func ensureOfficialCity(name string) (uint, error) {
	meta, ok := officialCityMetaByName[name]
	if !ok || meta.Country == "" || meta.Timezone == "" {
		return 0, fmt.Errorf("missing official city metadata: %s", name)
	}
	displayName := officialCityDisplayName(name)

	var cityID uint
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var city models.City
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("name_en = ?", name).
			Order("id ASC").
			First(&city).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			city = models.City{Name: displayName, NameEn: name, Country: meta.Country, Timezone: meta.Timezone}
			if err := tx.Create(&city).Error; err != nil {
				return err
			}
			cityID = city.ID
			return nil
		}
		if err != nil {
			return err
		}
		if err := tx.Model(&city).Updates(map[string]any{
			"name":     displayName,
			"country":  meta.Country,
			"timezone": meta.Timezone,
		}).Error; err != nil {
			return err
		}
		cityID = city.ID
		return nil
	})
	return cityID, err
}

func upsertOfficialStadium(name string, cityID uint) (uint, error) {
	stadiumName := officialStadiumDisplayName(name)
	var stadiumID uint
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var stadium models.Stadium
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("name_en = ?", name).
			Order("id ASC").
			First(&stadium).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			stadium = models.Stadium{Name: stadiumName, NameEn: name, CityID: cityID}
			if err := tx.Create(&stadium).Error; err != nil {
				return err
			}
			stadiumID = stadium.ID
			return nil
		}
		if err != nil {
			return err
		}
		if err := tx.Model(&stadium).Updates(map[string]any{
			"name":    stadiumName,
			"city_id": cityID,
		}).Error; err != nil {
			return err
		}
		stadiumID = stadium.ID
		return nil
	})
	return stadiumID, err
}

func officialCityDisplayName(name string) string {
	if displayName := officialCityDisplayNameByName[name]; displayName != "" {
		return displayName
	}
	return name
}

func officialStadiumDisplayName(name string) string {
	if displayName := officialStadiumDisplayNameByName[name]; displayName != "" {
		return displayName
	}
	return name
}

func ApplyOfficialVenueMappings() (int64, error) {
	var updated int64
	for externalID, venue := range officialVenueByExternalID {
		cityID, stadiumID, err := ensureOfficialVenue(venue)
		if err != nil {
			return updated, fmt.Errorf("official venue for match %d: %w", externalID, err)
		}

		result := database.DB.Model(&models.Match{}).
			Where("external_provider = ? AND external_id = ?", syncProviderFootballData, strconv.FormatInt(externalID, 10)).
			Updates(map[string]any{"city_id": cityID, "stadium_id": stadiumID})
		if result.Error != nil {
			return updated, result.Error
		}
		updated += result.RowsAffected
	}
	return updated, nil
}
