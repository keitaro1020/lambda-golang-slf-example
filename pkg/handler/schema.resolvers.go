package handler

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/domain"
	"github.com/keitaro1020/lambda-golang-slf-practice/scripts/graphql/generated"
)

func (r *catResolver) ID(ctx context.Context, obj *domain.Cat) (string, error) {
	return string(obj.ID), nil
}

func (r *queryResolver) Cats(ctx context.Context, first int64) ([]*domain.Cat, error) {
	return r.app.GetCats(ctx, first)
}

func (r *queryResolver) Cat(ctx context.Context, id string) (*domain.Cat, error) {
	return r.app.GetCat(ctx, domain.CatID(id))
}

// Cat returns generated.CatResolver implementation.
func (r *Resolver) Cat() generated.CatResolver { return &catResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type catResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
