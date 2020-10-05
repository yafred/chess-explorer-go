package stat

// Player ... a profile from Chess.com
type Player struct {
	PlayerID   int    `json:"player_id"`
	ID         string `json:"@id"`
	URL        string `json:"url"`
	Username   string `json:"username"`
	Followers  int    `json:"followers"`
	Country    string `json:"country"`
	LastOnline int    `json:"last_online"`
	Joined     int    `json:"joined"`
	Status     string `json:"status"`
	IsStreamer bool   `json:"is_streamer"`
}

// GamesContainer ... a list of Games from Chess.com
type GamesContainer struct {
	Games []struct {
		URL         string `json:"url"`
		Pgn         string `json:"pgn"`
		TimeControl string `json:"time_control"`
		EndTime     int    `json:"end_time"`
		Rated       bool   `json:"rated"`
		Fen         string `json:"fen"`
		TimeClass   string `json:"time_class"`
		Rules       string `json:"rules"`
		White       struct {
			Rating   int    `json:"rating"`
			Result   string `json:"result"`
			ID       string `json:"@id"`
			Username string `json:"username"`
		} `json:"white"`
		Black struct {
			Rating   int    `json:"rating"`
			Result   string `json:"result"`
			ID       string `json:"@id"`
			Username string `json:"username"`
		} `json:"black"`
	} `json:"games"`
}

// ArchivesContainer ... a list of available archives from Chess.com
type ArchivesContainer struct {
	Archives []string `json:"archives"`
}
