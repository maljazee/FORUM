package handlers

import (
	"forum/database"
	"net/http"
	"strconv"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	content := r.FormValue("content")
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if content == "" {
		http.Error(w, "Comment content is required", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	if err != nil {
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+strconv.Itoa(postID), http.StatusSeeOther)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(`
		SELECT c.id, c.user_id, u.username, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
	`, postID)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []struct {
		ID        int
		UserID    int
		Username  string
		Content   string
		CreatedAt string
	}

	for rows.Next() {
		var comment struct {
			ID        int
			UserID    int
			Username  string
			Content   string
			CreatedAt string
		}
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.Username, &comment.Content, &comment.CreatedAt)
		if err != nil {
			http.Error(w, "Error scanning comments", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	// Here you would typically render the comments using a template
	// For this example, we'll just write the number of comments
	w.Write([]byte("Number of comments: " + strconv.Itoa(len(comments))))
}

