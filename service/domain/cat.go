package domain

import "context"

// Cat is domain object
type Cat struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

// Cats is domain object array
type Cats []*Cat

// CatRepository is infrastructure interface
type CatClient interface {
	Search(ctx context.Context) (Cats, error)
}
