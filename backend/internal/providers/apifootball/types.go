package apifootball

type SquadResponse struct {
	Get      string       `json:"get"`
	Errors   interface{}  `json:"errors"`
	Results  int          `json:"results"`
	Response []SquadEntry `json:"response"`
}

type SquadEntry struct {
	Team    Team     `json:"team"`
	Players []Player `json:"players"`
}

type Team struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type TeamSearchResponse struct {
	Get      string            `json:"get"`
	Errors   interface{}       `json:"errors"`
	Results  int               `json:"results"`
	Response []TeamSearchEntry `json:"response"`
}

type TeamSearchEntry struct {
	Team  SearchTeam `json:"team"`
	Venue Venue      `json:"venue"`
}

type SearchTeam struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Country  string `json:"country"`
	National bool   `json:"national"`
	Logo     string `json:"logo"`
}

type Venue struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type Player struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Number   int    `json:"number"`
	Position string `json:"position"`
	Photo    string `json:"photo"`
}

type FixtureLineupsResponse struct {
	Get      string               `json:"get"`
	Errors   interface{}          `json:"errors"`
	Results  int                  `json:"results"`
	Response []FixtureLineupEntry `json:"response"`
}

type FixtureLineupEntry struct {
	Team        Team                     `json:"team"`
	Formation   string                   `json:"formation"`
	Coach       FixtureLineupCoach       `json:"coach"`
	StartXI     []FixtureLineupPlayerRow `json:"startXI"`
	Substitutes []FixtureLineupPlayerRow `json:"substitutes"`
}

type FixtureLineupCoach struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

type FixtureLineupPlayerRow struct {
	Player FixtureLineupPlayer `json:"player"`
}

type FixtureLineupPlayer struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Number int    `json:"number"`
	Pos    string `json:"pos"`
	Grid   string `json:"grid"`
}

type FixtureSearchResponse struct {
	Get      string               `json:"get"`
	Errors   interface{}          `json:"errors"`
	Results  int                  `json:"results"`
	Response []FixtureSearchEntry `json:"response"`
}

type FixtureSearchEntry struct {
	Fixture FixtureSearchFixture `json:"fixture"`
	Teams   FixtureSearchTeams   `json:"teams"`
}

type FixtureSearchFixture struct {
	ID   int64  `json:"id"`
	Date string `json:"date"`
}

type FixtureSearchTeams struct {
	Home Team `json:"home"`
	Away Team `json:"away"`
}
