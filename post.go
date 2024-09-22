package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID         int
	UserID     int
	Title      string
	Content    string
	CreatedAt  time.Time
	Categories []string
	Likes      int
	Dislikes   int
}

func CreatePost(db *sql.DB, userID int, title, content string, categories []string) (*Post, error) {
	query := `INSERT INTO posts (user_id, title, content, created_at) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, userID, title, content, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Insert categories
	for _, category := range categories {
		_, err = db.Exec("INSERT INTO post_categories (post_id, category) VALUES (?, ?)", id, category)
		if err != nil {
			return nil, err
		}
	}

	return &Post{
		ID:         int(id),
		UserID:     userID,
		Title:      title,
		Content:    content,
		CreatedAt:  time.Now(),
		Categories: categories,
	}, nil
}

func GetPostByID(db *sql.DB, id int) (*Post, error) {
	query := `SELECT id, user_id, title, content, created_at, 
              (SELECT COUNT(*) FROM post_likes WHERE post_id = posts.id) as likes,
              (SELECT COUNT(*) FROM post_dislikes WHERE post_id = posts.id) as dislikes
              FROM posts WHERE id = ?`
	row := db.QueryRow(query, id)

	post := &Post{}
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes)
	if err != nil {
		return nil, err
	}

	// Fetch categories
	rows, err := db.Query("SELECT category FROM post_categories WHERE post_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		post.Categories = append(post.Categories, category)
	}

	return post, nil
}

func GetAllPosts(db *sql.DB) ([]*Post, error) {
	query := `SELECT id, user_id, title, content, created_at, 
              (SELECT COUNT(*) FROM post_likes WHERE post_id = posts.id) as likes,
              (SELECT COUNT(*) FROM post_dislikes WHERE post_id = posts.id) as dislikes
              FROM posts ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *Post) AddLike(db *sql.DB, userID int) error {
	_, err := db.Exec("INSERT INTO post_likes (post_id, user_id) VALUES (?, ?)", p.ID, userID)
	return err
}

func (p *Post) AddDislike(db *sql.DB, userID int) error {
	_, err := db.Exec("INSERT INTO post_dislikes (post_id, user_id) VALUES (?, ?)", p.ID, userID)
	return err
}
