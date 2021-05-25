package jsonplaceholder

import (
	"context"
	"fmt"
	"sync"
)

type Album struct {
	UserID int      `json:"userId"`
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Photos []*Photo `json:"photos"`
}

type AlbumService interface {
	// List returns all resources
	List(ctx context.Context) ([]*Album, error)

	// Get return the resource matching given ID
	Get(ctx context.Context, albumID uint64) (*Album, error)

	// SearchByUserID return albums owned by given user
	SearchByUserID(ctx context.Context, userID uint64) ([]*Album, error)
}

type albumService struct {
	c            *Client
	resourceName Resource
}

func (a *albumService) List(ctx context.Context) ([]*Album, error) {
	var albums []*Album
	if err := a.c.fetchMultiple(ctx, a.resourceName, &albums); err != nil {
		return albums, err
	}

	var wg sync.WaitGroup
	wg.Add(len(albums))
	for _, alb := range albums {
		go func(album *Album) {
			defer wg.Done()
			a.processSubResources(ctx, album)
		}(alb)
	}

	wg.Wait()
	return albums, nil
}

func (a *albumService) Get(ctx context.Context, albumID uint64) (*Album, error) {
	var album *Album
	if err := a.c.fetch(ctx, a.resourceName, albumID, &album); err != nil {
		return album, err
	}

	a.processSubResources(ctx, album)
	return album, nil
}

func (a *albumService) SearchByUserID(ctx context.Context, userID uint64) ([]*Album, error) {
	if userID == 0 {
		return nil, fmt.Errorf("invalid userID")
	}

	var albums []*Album
	if err := a.c.fetchSub(ctx, UserResource, userID, a.resourceName, &albums); err != nil {
		return albums, nil
	}

	var wg sync.WaitGroup
	wg.Add(len(albums))
	for _, alb := range albums {
		go func(album *Album) {
			defer wg.Done()
			a.processSubResources(ctx, album)
		}(alb)
	}

	wg.Wait()
	return albums, nil
}

func (a *albumService) processSubResources(ctx context.Context, album *Album) {
	album.Photos = make([]*Photo, 0)

	var photos []*Photo
	if err := a.c.fetchSub(ctx, a.resourceName, uint64(album.ID), PhotoResource, &photos); err == nil {
		album.Photos = photos
	}
}

func newAlbumService(c *Client) AlbumService {
	return &albumService{c: c, resourceName: AlbumResource}
}
