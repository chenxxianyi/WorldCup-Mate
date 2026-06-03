package services

import (
	"sort"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"
	"worldcup-mate/internal/repositories"
)

func GetStandingsByGroupID(groupID uint) ([]models.GroupStanding, error) {
	standings, err := repositories.GetStandingsByGroupID(groupID)
	if err != nil {
		return nil, err
	}

	// If no standings exist, generate initial standings from teams in the group
	if len(standings) == 0 {
		teams, err := repositories.GetTeamsByGroupID(groupID)
		if err != nil {
			return nil, err
		}

		for i, team := range teams {
			standings = append(standings, models.GroupStanding{
				GroupID:            groupID,
				TeamID:             team.ID,
				Team:               team,
				Played:             0,
				Won:                0,
				Drawn:              0,
				Lost:               0,
				GoalsFor:           0,
				GoalsAgainst:       0,
				GoalDifference:     0,
				Points:             0,
				Rank:               i + 1,
				QualificationStatus: "unknown",
			})
		}
	}

	return standings, nil
}

func GetAllStandings() ([]models.GroupStanding, error) {
	return repositories.GetAllStandings()
}

func GetBestThird() ([]models.GroupStanding, error) {
	standings, err := repositories.GetBestThird()
	if err != nil {
		return nil, err
	}

	// If no third-place standings exist, generate initial data from each group's third team
	if len(standings) == 0 {
		groups, err := repositories.ListGroups()
		if err != nil {
			return nil, err
		}

		for _, group := range groups {
			// Get all teams in this group
			teams, err := repositories.GetTeamsByGroupID(group.ID)
			if err != nil || len(teams) < 3 {
				continue
			}
			// Third team (index 2) is the third-place team
			thirdTeam := teams[2]
			standings = append(standings, models.GroupStanding{
				GroupID:            group.ID,
				Group:              group,
				TeamID:             thirdTeam.ID,
				Team:               thirdTeam,
				Played:             0,
				Won:                0,
				Drawn:              0,
				Lost:               0,
				GoalsFor:           0,
				GoalsAgainst:       0,
				GoalDifference:     0,
				Points:             0,
				Rank:               3,
				QualificationStatus: "unknown",
			})
		}
	}

	return standings, nil
}

func RecalculateGroupStanding(groupID uint) error {
	matches, _, err := repositories.ListMatches(repositories.MatchFilter{
		GroupID:  groupID,
		Status:   "finished",
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		return err
	}

	stats := make(map[uint]*models.GroupStanding)

	for _, m := range matches {
		if m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		home, ok := stats[m.HomeTeamID]
		if !ok {
			home = &models.GroupStanding{GroupID: groupID, TeamID: m.HomeTeamID}
			stats[m.HomeTeamID] = home
		}
		away, ok := stats[m.AwayTeamID]
		if !ok {
			away = &models.GroupStanding{GroupID: groupID, TeamID: m.AwayTeamID}
			stats[m.AwayTeamID] = away
		}

		home.Played++
		away.Played++
		home.GoalsFor += *m.HomeScore
		home.GoalsAgainst += *m.AwayScore
		away.GoalsFor += *m.AwayScore
		away.GoalsAgainst += *m.HomeScore

		if *m.HomeScore > *m.AwayScore {
			home.Won++
			home.Points += 3
			away.Lost++
		} else if *m.HomeScore < *m.AwayScore {
			away.Won++
			away.Points += 3
			home.Lost++
		} else {
			home.Drawn++
			away.Drawn++
			home.Points++
			away.Points++
		}
	}

	standings := make([]*models.GroupStanding, 0, len(stats))
	for _, s := range stats {
		s.GoalDifference = s.GoalsFor - s.GoalsAgainst
		standings = append(standings, s)
	}

	sort.Slice(standings, func(i, j int) bool {
		if standings[i].Points != standings[j].Points {
			return standings[i].Points > standings[j].Points
		}
		if standings[i].GoalDifference != standings[j].GoalDifference {
			return standings[i].GoalDifference > standings[j].GoalDifference
		}
		if standings[i].GoalsFor != standings[j].GoalsFor {
			return standings[i].GoalsFor > standings[j].GoalsFor
		}
		return standings[i].TeamID < standings[j].TeamID
	})

	for i, s := range standings {
		s.Rank = i + 1
		switch i {
		case 0, 1:
			s.QualificationStatus = "qualified"
		case 2:
			s.QualificationStatus = "possible"
		default:
			s.QualificationStatus = "eliminated"
		}
		if err := repositories.UpsertStanding(s); err != nil {
			return err
		}
	}

	return nil
}

func RecalculateBestThird() error {
	groups, err := repositories.ListGroups()
	if err != nil {
		return err
	}

	var thirdPlaceTeams []models.GroupStanding
	for _, g := range groups {
		standings, err := repositories.GetStandingsByGroupID(g.ID)
		if err != nil {
			continue
		}
		for _, s := range standings {
			if s.Rank == 3 {
				thirdPlaceTeams = append(thirdPlaceTeams, s)
			}
		}
	}

	sort.Slice(thirdPlaceTeams, func(i, j int) bool {
		if thirdPlaceTeams[i].Points != thirdPlaceTeams[j].Points {
			return thirdPlaceTeams[i].Points > thirdPlaceTeams[j].Points
		}
		if thirdPlaceTeams[i].GoalDifference != thirdPlaceTeams[j].GoalDifference {
			return thirdPlaceTeams[i].GoalDifference > thirdPlaceTeams[j].GoalDifference
		}
		if thirdPlaceTeams[i].GoalsFor != thirdPlaceTeams[j].GoalsFor {
			return thirdPlaceTeams[i].GoalsFor > thirdPlaceTeams[j].GoalsFor
		}
		return thirdPlaceTeams[i].TeamID < thirdPlaceTeams[j].TeamID
	})

	for i, s := range thirdPlaceTeams {
		if i < 8 {
			s.QualificationStatus = "qualified"
		} else {
			s.QualificationStatus = "eliminated"
		}
		_ = database.DB.Model(&models.GroupStanding{}).Where("id = ?", s.ID).
			Update("qualification_status", s.QualificationStatus).Error
	}

	return nil
}
