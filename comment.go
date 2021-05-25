package jsonplaceholder

import "context"

type Comment struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

type CommentService interface {
	// List returns all resources
	List(ctx context.Context) ([]*Comment, error)

	// Get return the resource matching given ID
	Get(ctx context.Context, postID uint64) (*Comment, error)
}

type commentService struct {
	c            *Client
	resourceName Resource
}

func (c *commentService) List(ctx context.Context) ([]*Comment, error) {
	var comments []*Comment
	err := c.c.fetchMultiple(ctx, c.resourceName, &comments)
	return comments, err
}

func (c *commentService) Get(ctx context.Context, commentID uint64) (*Comment, error) {
	var comment *Comment
	err := c.c.fetch(ctx, c.resourceName, commentID, &comment)
	return comment, err
}

func newCommentService(c *Client) CommentService {
	return &commentService{c: c, resourceName: CommentResource}
}
