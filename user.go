package jsonplaceholder

import (
	"context"
	"fmt"
	"sync"
)

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Address  *Address `json:"address"`
	Phone    string   `json:"phone"`
	Website  string   `json:"website"`
	Company  *Company `json:"company"`
	Albums   []*Album `json:"albums"`
	Todos    []*Todo  `json:"todos"`
	Posts    []*Post  `json:"posts"`
}

type Geo struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Address struct {
	Street  string `json:"street"`
	Suite   string `json:"suite"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
	Geo     *Geo   `json:"geo"`
}

type Company struct {
	Name        string `json:"name"`
	Catchphrase string `json:"catchPhrase"`
	Bs          string `json:"bs"`
}

type UserService interface {
	// List returns all resources
	List(ctx context.Context) ([]*User, error)

	// Get return the resource matching given ID
	Get(ctx context.Context, userID uint64) (*User, error)
}

type userService struct {
	c            *Client
	resourceName Resource
}

func (u *userService) List(ctx context.Context) ([]*User, error) {
	var users []*User
	if err := u.c.fetchMultiple(ctx, u.resourceName, &users); err != nil {
		return users, err
	}

	// Fetch user's sub resources
	var wg sync.WaitGroup
	wg.Add(len(users))
	for _, usr := range users {
		go func(user *User) {
			defer wg.Done()
			u.processSubResources(ctx, user)
		}(usr)
	}

	wg.Wait()
	return users, nil
}

func (u *userService) Get(ctx context.Context, userID uint64) (*User, error) {
	if userID == 0 {
		return nil, fmt.Errorf("invalid userID")
	}

	var user *User
	if err := u.c.fetch(ctx, u.resourceName, userID, &user); err != nil {
		return user, err
	}

	u.processSubResources(ctx, user)
	return user, nil
}

func newUserService(c *Client) UserService {
	return &userService{c: c, resourceName: UserResource}
}

func (u *userService) processSubResources(ctx context.Context, usr *User) {
	var wg sync.WaitGroup

	// Init empty slices
	usr.Albums = make([]*Album, 0)
	usr.Todos = make([]*Todo, 0)
	usr.Posts = make([]*Post, 0)

	// Albums
	wg.Add(1)
	go func(user *User) {
		defer wg.Done()
		if albums, err := u.c.Album.SearchByUserID(ctx, uint64(user.ID)); err == nil {
			user.Albums = albums
		}
	}(usr)

	// Todos
	wg.Add(1)
	go func(user *User) {
		defer wg.Done()
		if todos, err := u.c.Todo.SearchByUserID(ctx, uint64(user.ID)); err == nil {
			user.Todos = todos
		}
	}(usr)

	// Posts
	wg.Add(1)
	go func(user *User) {
		defer wg.Done()
		if posts, err := u.c.Post.SearchByUserID(ctx, uint64(user.ID)); err == nil {
			user.Posts = posts
		}
	}(usr)

	wg.Wait()
}
