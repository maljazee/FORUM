package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Content   string
	CreatedAt time.Time
	Likes     int
	Dislikes  int
}

func CreateComment(db *sql.DB, postID, userID int, content string) (*Comment, error) {
	query := `INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, postID, userID, content, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Comment{
		ID:        int(id),
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

func GetCommentsByPostID(db *sql.DB, postID int) ([]*Comment, error) {
	query := `SELECT c.id, c.user_id, c.content, c.created_at,
              (SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id) as likes,
              (SELECT COUNT(*) FROM comment_dislikes WHERE comment_id = c.id) as dislikes
              FROM comments c WHERE c.post_id = ? ORDER BY c.created_at DESC`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{PostID: postID}
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Likes, &comment.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (c *Comment) AddLike(db *sql.DB, userID int) error {
	_, err := db.Exec("INSERT INTO comment_likes (comment_id, user_id) VALUES (?, ?)", c.ID, userID)
	return err
}

func (c *Comment) AddDislike(db *sql.DB, userID int) error {
	_, err := db.Exec("INSERT INTO comment_dislikes (comment_id, user_id) VALUES (?, ?)", c.ID, userID)
	return err
}

func (c *Comment) Delete(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM comments WHERE id = ?", c.ID)
	return err
}
