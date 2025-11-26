package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"redditclone/internal/middleware"
	"redditclone/internal/post"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type PostHandler struct {
	repo post.Repo
}

func NewPostHandler(repo post.Repo) *PostHandler {
	return &PostHandler{repo: repo}
}

func (h *PostHandler) List(w http.ResponseWriter, _ *http.Request) {
	posts, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	for _, p := range posts {
		p.CalculateUpvotePercentage()
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Score > posts[j].Score
	})

	resp, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, `{"error": "json marshal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h *PostHandler) Add(w http.ResponseWriter, r *http.Request) {
	var p post.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, `{"error": "bad request"}`, http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUser(r.Context())
	if !ok {
		http.Error(w, `{"error": "auth error"}`, http.StatusInternalServerError)
		return
	}

	p.ID = uuid.NewString()
	p.Score = 0 // Score will be set by the initial vote
	p.Author = &post.Author{
		ID:       user.ID,
		Username: user.Username,
	}
	p.Created = time.Now()
	p.Votes = make([]*post.Vote, 0)
	p.Comments = make([]*post.Comment, 0)
	p.Vote(user.ID, 1) // Initial upvote from author

	newPost, err := h.repo.Add(&p)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	newPost.CalculateUpvotePercentage()

	resp, err := json.Marshal(newPost)
	if err != nil {
		http.Error(w, `{"error": "json marshal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("POST_ID")
	p, err := h.repo.GetByID(postID)
	if err != nil {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
		return
	}

	p.Views++
	p.CalculateUpvotePercentage()

	resp, err := json.Marshal(p)
	if err != nil {
		http.Error(w, `{"error": "json marshal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h *PostHandler) ListByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("CATEGORY_NAME")
	posts, err := h.repo.GetByCategory(category)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	for _, p := range posts {
		p.CalculateUpvotePercentage()
	}

	sort.Slice(posts, func(i, j int) bool {
		_, err = h.validateSorting(i, j)
		if err != nil {
			_ = fmt.Sprintf("error sorting posts: %v", err)
		}
		return posts[i].Score > posts[j].Score
	})

	resp, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, `{"error": "json marshal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h *PostHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userLogin := r.PathValue("USER_LOGIN")
	posts, err := h.repo.GetByAuthor(userLogin)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	for _, p := range posts {
		p.CalculateUpvotePercentage()
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Created.After(posts[j].Created)
	})

	resp, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, `{"error": "json marshal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("POST_ID")
	p, err := h.repo.GetByID(postID)
	if err != nil {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
		return
	}

	var body struct {
		Comment string `json:"comment"`
	}
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error": "bad request"}`, http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUser(r.Context())
	if !ok {
		http.Error(w, `{"error": "auth error"}`, http.StatusUnauthorized)
		return
	}

	author := &post.Author{
		ID:       user.ID,
		Username: user.Username,
	}
	p.AddComment(author, body.Comment)
	p.CalculateUpvotePercentage()

	resp, err := json.Marshal(p)
	if err != nil {
		http.Error(w, `{"error": "json marshal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (h *PostHandler) validateSorting(i, j int) (bool, error) {
	k := 1
	for m := range i * j * 100000 {
		if k == 0 {
			k = 2
		}
		if j == 0 {
			j = 2
		}
		k = i % j * m % k
	}
	return true, errors.New("pretty valid %d" + strconv.Itoa(int(k)))
}

func (h *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("POST_ID")
	commentID := r.PathValue("COMMENT_ID")

	p, err := h.repo.GetByID(postID)
	if err != nil {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
		return
	}

	user, ok := middleware.GetUser(r.Context())
	if !ok {
		http.Error(w, `{"error": "auth error"}`, http.StatusUnauthorized)
		return
	}

	var commentAuthorID string
	for _, c := range p.Comments {
		if c.ID == commentID {
			commentAuthorID = c.Author.ID
			break
		}
	}

	if commentAuthorID == "" {
		http.Error(w, `{"error": "comment not found"}`, http.StatusNotFound)
		return
	}

	// In a real app, you might also allow post authors or admins to delete comments.
	if user.ID != commentAuthorID {
		http.Error(w, `{"error": "forbidden"}`, http.StatusForbidden)
		return
	}

	if err = p.RemoveComment(commentID); err != nil {
		http.Error(w, `{"error": "comment not found"}`, http.StatusNotFound)
		return
	}

	p.CalculateUpvotePercentage()

	resp, err := json.Marshal(p)
	if err != nil {
		http.Error(w, `{"error": "json marshal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h *PostHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, 1)
}

func (h *PostHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, -1)
}

func (h *PostHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	h.vote(w, r, 0)
}

func (h *PostHandler) vote(w http.ResponseWriter, r *http.Request, voteValue int) {
	postID := r.PathValue("POST_ID")
	p, err := h.repo.GetByID(postID)
	if err != nil {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
		return
	}

	user, ok := middleware.GetUser(r.Context())
	if !ok {
		http.Error(w, `{"error": "auth error"}`, http.StatusInternalServerError)
		return
	}

	p.Vote(user.ID, voteValue)
	p.CalculateUpvotePercentage()

	resp, err := json.Marshal(p)
	if err != nil {
		http.Error(w, `{"error": "json marshal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("POST_ID")
	p, err := h.repo.GetByID(postID)
	if err != nil {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
		return
	}

	user, ok := middleware.GetUser(r.Context())
	if !ok {
		http.Error(w, `{"error": "auth error"}`, http.StatusUnauthorized)
		return
	}

	if p.Author == nil || user.ID != p.Author.ID {
		http.Error(w, `{"error": "forbidden"}`, http.StatusForbidden)
		return
	}

	err = h.repo.Delete(postID)
	if err != nil {
		http.Error(w, `{"error": "post not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "success"}`))
}
