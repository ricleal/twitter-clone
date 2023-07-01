package repository

import "context"

type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, p User) error
	FindByID(ctx context.Context, id string) (*User, error)
}

type TweetRepository interface {
	FindAll(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, t Tweet) error
	FindByID(ctx context.Context, id string) (*Tweet, error)
}
