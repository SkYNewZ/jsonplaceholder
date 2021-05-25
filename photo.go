package jsonplaceholder

import "context"

type Photo struct {
	AlbumID      int    `json:"albumId"`
	ID           int    `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type PhotoService interface {
	// List returns all resources
	List(ctx context.Context) ([]*Photo, error)

	// Get return the resource matching given ID
	Get(ctx context.Context, photoID uint64) (*Photo, error)
}

type photoService struct {
	c            *Client
	resourceName Resource
}

func (p *photoService) List(ctx context.Context) ([]*Photo, error) {
	var photos []*Photo
	err := p.c.fetchMultiple(ctx, p.resourceName, &photos)
	return photos, err
}

func (p *photoService) Get(ctx context.Context, photoID uint64) (*Photo, error) {
	var photo *Photo
	err := p.c.fetch(ctx, p.resourceName, photoID, &photo)
	return photo, err
}

func newPhotoService(c *Client) PhotoService {
	return &photoService{c: c, resourceName: PhotoResource}
}
