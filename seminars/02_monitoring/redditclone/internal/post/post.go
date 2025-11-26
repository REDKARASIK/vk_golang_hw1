package post

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Comment struct {
	ID      string    `json:"id"`
	Author  *Author   `json:"author"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

type Vote struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

type Post struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	URL              string     `json:"url,omitempty"`
	Author           *Author    `json:"author"`
	Category         string     `json:"category"`
	Score            int        `json:"score"`
	Votes            []*Vote    `json:"votes"`
	Comments         []*Comment `json:"comments"`
	Created          time.Time  `json:"created"`
	Views            int        `json:"views"`
	Type             string     `json:"type"`
	Text             string     `json:"text,omitempty"`
	UpvotePercentage int        `json:"upvotePercentage"`
}

type Repo interface {
	GetAll() ([]*Post, error)
	GetByID(id string) (*Post, error)
	GetByCategory(category string) ([]*Post, error)
	GetByAuthor(authorUsername string) ([]*Post, error)
	Add(post *Post) (*Post, error)
	Delete(id string) error
}

type MemoryRepo struct {
	mu    sync.RWMutex
	posts []*Post
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		posts: make([]*Post, 0),
	}
}

func (r *MemoryRepo) GetAll() ([]*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.posts, nil
}

func (r *MemoryRepo) GetByID(id string) (*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.posts {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("post not found")
}

func (r *MemoryRepo) GetByCategory(category string) ([]*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*Post
	for _, p := range r.posts {
		if p.Category == category {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

func (r *MemoryRepo) GetByAuthor(authorUsername string) ([]*Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []*Post
	for _, p := range r.posts {
		if p.Author != nil && p.Author.Username == authorUsername {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

func (r *MemoryRepo) Add(post *Post) (*Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.posts = append(r.posts, post)
	return post, nil
}

func (r *MemoryRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, p := range r.posts {
		if p.ID == id {
			r.posts = append(r.posts[:i], r.posts[i+1:]...)
			return nil
		}
	}
	return errors.New("post not found")
}

func (p *Post) Vote(userID string, vote int) {
	existingVoteIndex := -1
	for i, v := range p.Votes {
		if v.User == userID {
			existingVoteIndex = i
			break
		}
	}

	if existingVoteIndex != -1 {
		existingVote := p.Votes[existingVoteIndex]
		p.Score -= existingVote.Vote
		if vote == 0 {
			p.Votes = append(p.Votes[:existingVoteIndex], p.Votes[existingVoteIndex+1:]...)
		} else {
			p.Score += vote
			existingVote.Vote = vote
		}
	} else if vote != 0 {
		p.Score += vote
		p.Votes = append(p.Votes, &Vote{User: userID, Vote: vote})
	}
}

func (p *Post) AddComment(author *Author, body string) {
	p.Comments = append(p.Comments, &Comment{
		ID:      uuid.NewString(),
		Author:  author,
		Body:    body,
		Created: time.Now(),
	})
}

func (p *Post) RemoveComment(commentID string) error {
	for i, c := range p.Comments {
		if c.ID == commentID {
			p.Comments = append(p.Comments[:i], p.Comments[i+1:]...)
			return nil
		}
	}
	return errors.New("comment not found")
}

func (p *Post) CalculateUpvotePercentage() {
	upvotes := 0
	downvotes := 0
	for _, v := range p.Votes {
		if v.Vote == 1 {
			upvotes++
		} else if v.Vote == -1 {
			downvotes++
		}
	}
	totalVotes := upvotes + downvotes
	p.UpvotePercentage = (upvotes * 100) / totalVotes
}
