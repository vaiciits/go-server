package handlers

import (
	"database/sql"
	"go-server/models"
)

func HandleGetPosts(nfd int, db *sql.DB) {
	rows, err := db.Query(`
		SELECT p.id, u.username, p.content, 
		       (SELECT COUNT(*) FROM likes l WHERE l.post_id = p.id) as likes, 
		       p.created_at 
		FROM posts p 
		JOIN users u ON p.user_id = u.id 
		ORDER BY p.id DESC LIMIT 10
	`)
	if err != nil {
		SendResponse(nfd, 500, map[string]string{"error": err.Error()})
		return
	}
	if rows.Err() != nil {
		SendResponse(nfd, 500, map[string]string{"error": rows.Err().Error()})
		return
	}
	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Author, &p.Content, &p.Likes, &p.CreatedAt); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	SendResponse(nfd, 200, posts)
}
