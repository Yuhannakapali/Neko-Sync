package anime

// Anime represents an anime content item.
type Anime struct {
	ID          string `json:"id" db:"id"`
	NoOfEpisode int    `json:"no_of_episode" db:"no_of_episode"`
	Season      string `json:"season" db:"season"`
	Year        string `json:"year" db:"year"`
	Status      string `json:"status" db:"status"`
	Rating      string `json:"rating" db:"rating"`
}
