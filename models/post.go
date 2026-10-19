package models

type Post struct {
	ID        int64  `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	Likes     int    `json:"likes"`
	CreatedAt string `json:"createdAt"`
}
