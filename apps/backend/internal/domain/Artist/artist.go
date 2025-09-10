package artist

// Artist represents an artist in the system.
type Artist struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Bio       string `json:"bio" db:"bio"`
	AvatarURL string `json:"avatar_url" db:"avatar_url"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}
