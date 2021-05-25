package jsonplaceholder

import (
	"context"
	"fmt"
	"sync"
)

type Post struct {
	UserID   int        `json:"userId"`
	ID       int        `json:"id"`
	Title    string     `json:"title"`
	Body     string     `json:"body"`
	Comments []*Comment `json:"comments"`
}

type PostService interface {
	// List returns all resources
	List(ctx context.Context) ([]*Post, error)

	// Get return the resource matching given ID
	Get(ctx context.Context, postID uint64) (*Post, error)

	// SearchByUserID returns post written by given user ID
	SearchByUserID(ctx context.Context, userID uint64) ([]*Post, error)
}

type postService struct {
	c            *Client
	resourceName Resource
}

func (p *postService) List(ctx context.Context) ([]*Post, error) {
	var posts []*Post
	if err := p.c.fetchMultiple(ctx, p.resourceName, &posts); err != nil {
		return posts, nil
	}

	var wg sync.WaitGroup
	wg.Add(len(posts))
	for _, pst := range posts {
		go func(post *Post) {
			defer wg.Done()
			p.processSubResources(ctx, post)
		}(pst)
	}

	wg.Wait()
	return posts, nil
}

func (p *postService) Get(ctx context.Context, postID uint64) (*Post, error) {
	if postID == 0 {
		return nil, fmt.Errorf("invalid postID")
	}

	var post *Post
	if err := p.c.fetch(ctx, p.resourceName, postID, &post); err != nil {
		return post, nil
	}

	p.processSubResources(ctx, post)
	return post, nil
}

func (p *postService) SearchByUserID(ctx context.Context, userID uint64) ([]*Post, error) {
	if userID == 0 {
		return nil, fmt.Errorf("invalid userID")
	}

	var posts []*Post
	if err := p.c.fetchSub(ctx, UserResource, userID, p.resourceName, &posts); err != nil {
		return posts, err
	}

	var wg sync.WaitGroup
	wg.Add(len(posts))
	for _, pst := range posts {
		go func(post *Post) {
			defer wg.Done()
			p.processSubResources(ctx, post)
		}(pst)
	}

	wg.Wait()
	return posts, nil
}

func newPostService(c *Client) PostService {
	return &postService{c: c, resourceName: PostResource}
}

func (p *postService) processSubResources(ctx context.Context, post *Post) {
	// Init empty slice
	post.Comments = make([]*Comment, 0)

	// Fetch comments for given post
	var comments []*Comment
	if err := p.c.fetchSub(ctx, p.resourceName, uint64(post.ID), CommentResource, &comments); err == nil {
		post.Comments = comments
	}
}
