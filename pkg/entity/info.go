package entity

type Info struct {
	ID              string `json:"_id"`
	DisplayName     string `json:"display_name"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profile_image_url"`
	Description     string `json:"description"`
	Follower        int    `json:"follower"`
	Following       int    `json:"following"`
}
