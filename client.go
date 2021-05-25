package jsonplaceholder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	apiURL = "https://jsonplaceholder.typicode.com"
)

type Resource string

const (
	PostResource    Resource = "posts"
	CommentResource Resource = "comments"
	AlbumResource   Resource = "albums"
	PhotoResource   Resource = "photos"
	TodoResource    Resource = "todos"
	UserResource    Resource = "users"
)

// Client https://jsonplaceholder.typicode.com/ client
type Client struct {
	httpClient *http.Client

	Post    PostService
	Comment CommentService
	Album   AlbumService
	Photo   PhotoService
	Todo    TodoService
	User    UserService
}

// New creates a new Client
func New(opts ...Option) *Client {
	c := new(Client)
	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	c.Post = newPostService(c)
	c.Comment = newCommentService(c)
	c.Photo = newPhotoService(c)
	c.Album = newAlbumService(c)
	c.Todo = newTodoService(c)
	c.User = newUserService(c)
	return c
}

func (c *Client) fetchMultiple(ctx context.Context, resource Resource, data interface{}) error {
	url := makeURL(resource)
	return c.doFetch(ctx, url, data)
}

func (c *Client) fetchSub(ctx context.Context, resource Resource, resourceID uint64, subResource Resource, data interface{}) error {
	url := makeURL(resource, resourceID, subResource)
	return c.doFetch(ctx, url, data)
}

func (c *Client) fetch(ctx context.Context, resource Resource, resourceID uint64, data interface{}) error {
	url := makeURL(resource, resourceID)
	return c.doFetch(ctx, url, data)
}

// TODO: handle POST/PUT/DELETE with body
func (c *Client) doFetch(ctx context.Context, url string, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request return an error code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	return nil
}

func makeURL(params ...interface{}) string {
	var url strings.Builder
	url.WriteString(apiURL)
	for _, param := range params {
		url.WriteRune('/')
		url.WriteString(fmt.Sprint(param))
	}

	return url.String()
}
