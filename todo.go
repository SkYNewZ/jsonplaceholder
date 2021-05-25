package jsonplaceholder

import (
	"context"
	"fmt"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoService interface {
	// List returns all resources
	List(ctx context.Context) ([]*Todo, error)

	// Get return the resource matching given ID
	Get(ctx context.Context, todoID uint64) (*Todo, error)

	// SearchByUserID returns todos written by given user ID
	SearchByUserID(ctx context.Context, userID uint64) ([]*Todo, error)
}

type todoService struct {
	c            *Client
	resourceName Resource
}

func (t *todoService) List(ctx context.Context) ([]*Todo, error) {
	var todos []*Todo
	err := t.c.fetchMultiple(ctx, t.resourceName, &todos)
	return todos, err
}

func (t *todoService) Get(ctx context.Context, todoID uint64) (*Todo, error) {
	if todoID == 0 {
		return nil, fmt.Errorf("invalid todoID")
	}

	var todo *Todo
	err := t.c.fetch(ctx, t.resourceName, todoID, &todo)
	return todo, err
}

func (t *todoService) SearchByUserID(ctx context.Context, userID uint64) ([]*Todo, error) {
	if userID == 0 {
		return nil, fmt.Errorf("invalid userID")
	}

	var todos []*Todo
	err := t.c.fetchSub(ctx, UserResource, userID, t.resourceName, &todos)
	return todos, err

}

func newTodoService(c *Client) TodoService {
	return &todoService{c: c, resourceName: TodoResource}
}
