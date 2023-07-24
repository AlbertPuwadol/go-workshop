package entity

type Data struct {
	ID           string `json:"_id"`
	Text         string `json:"text"`
	UserID       string `json:"user_id"`
	Likes        int    `json:"likes"`
	ParentThread string `json:"parent_thread"`
	RepostCount  int    `json:"repost_count"`
	User         Info   `json:"user"`
}

type Timeline struct {
	Data     []Data  `json:"data"`
	NextPage *string `json:"next_page"`
}
