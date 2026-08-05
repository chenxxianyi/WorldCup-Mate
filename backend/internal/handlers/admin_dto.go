package handlers

// Admin DTOs (ADM-01). All admin write endpoints use explicit DTOs;
// core updates must never accept a bare map[string]interface{}.

type TeamCreateInput struct {
	Name         string  `json:"name" binding:"required"`
	NameEn       string  `json:"name_en"`
	FIFACode     *string `json:"fifa_code"`
	ExternalCode *string `json:"external_code"`
	TeamType     string  `json:"team_type"`
	FlagURL      string  `json:"flag_url"`
	Continent    string  `json:"continent"`
	Country      string  `json:"country"`
	Venue        string  `json:"venue"`
	GroupID      *uint   `json:"group_id"`
	Coach        string  `json:"coach"`
	Description  string  `json:"description"`
}

// TeamUpdateInput: pointer fields allow partial updates; a nil field is
// left unchanged (except FIFACode/GroupID where nil means "clear").
type TeamUpdateInput struct {
	Name         *string `json:"name"`
	NameEn       *string `json:"name_en"`
	FIFACode     *string `json:"fifa_code"`
	ExternalCode *string `json:"external_code"`
	TeamType     *string `json:"team_type"`
	FlagURL      *string `json:"flag_url"`
	Continent    *string `json:"continent"`
	Country      *string `json:"country"`
	Venue        *string `json:"venue"`
	GroupID      *uint   `json:"group_id"`
	Coach        *string `json:"coach"`
	Description  *string `json:"description"`
}

type GroupCreateInput struct {
	Name  string `json:"name" binding:"required"`
	Stage string `json:"stage"`
}

type GroupUpdateInput struct {
	Name  *string `json:"name"`
	Stage *string `json:"stage"`
}

type CityCreateInput struct {
	Name     string `json:"name" binding:"required"`
	NameEn   string `json:"name_en"`
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
}

type CityUpdateInput struct {
	Name     *string `json:"name"`
	NameEn   *string `json:"name_en"`
	Country  *string `json:"country"`
	Timezone *string `json:"timezone"`
}

type StadiumCreateInput struct {
	Name     string `json:"name" binding:"required"`
	NameEn   string `json:"name_en"`
	CityID   uint   `json:"city_id" binding:"required"`
	Capacity int    `json:"capacity"`
}

type StadiumUpdateInput struct {
	Name     *string `json:"name"`
	NameEn   *string `json:"name_en"`
	CityID   *uint   `json:"city_id"`
	Capacity *int    `json:"capacity"`
}

type MatchUpdateInput struct {
	MatchNo         *int    `json:"match_no"`
	HomeTeamID      *uint   `json:"home_team_id"`
	AwayTeamID      *uint   `json:"away_team_id"`
	GroupID         *uint   `json:"group_id"`
	Stage           *string `json:"stage"`
	StadiumID       *uint   `json:"stadium_id"`
	CityID          *uint   `json:"city_id"`
	KickoffTimeUTC  *string `json:"kickoff_time_utc"`
	ImportanceLevel *int    `json:"importance_level"`
	RecommendTag    *string `json:"recommend_tag"`
	CompetitionID   *uint   `json:"competition_id"`
	Season          *int    `json:"season"`
	Matchday        *int    `json:"matchday"`
}

// StandingUpdateInput: manual standings correction with an audit reason.
type StandingUpdateInput struct {
	Position       *int   `json:"position"`
	Points         *int   `json:"points"`
	Played         *int   `json:"played"`
	Won            *int   `json:"won"`
	Drawn          *int   `json:"drawn"`
	Lost           *int   `json:"lost"`
	GoalsFor       *int   `json:"goals_for"`
	GoalsAgainst   *int   `json:"goals_against"`
	GoalDifference *int   `json:"goal_difference"`
	Reason         string `json:"reason"`
}

// CompetitionUpdateInput: explicit DTO for admin competition updates
// (ADM-01 convention — never a bare map). Status is validated so a typo
// cannot silently hide a league from the public list.
type CompetitionUpdateInput struct {
	Name      *string `json:"name"`
	NameEn    *string `json:"name_en"`
	Country   *string `json:"country"`
	LogoURL   *string `json:"logo_url"`
	Format    *string `json:"format"`
	Status    *string `json:"status"`
	Season    *int    `json:"season"`
	SortOrder *int    `json:"sort_order"`
}
