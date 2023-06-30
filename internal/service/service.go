package service

import (
	"context"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/ricleal/twitter-clone/internal/openapi"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) Create(ctx context.Context, u openapi.User) error {
	user := User{
		Username: u.Username,
		Email:    string(u.Email),
	}
	if u.Name != nil {
		user.Name = *u.Name
	}
	return s.repo.Create(ctx, user)
}

func (s *UserService) FindAll(ctx context.Context) ([]openapi.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var openapiUsers []openapi.User
	for _, u := range users {
		openapiUsers = append(openapiUsers, openapi.User{
			Id:       &u.ID,
			Username: u.Username,
			Email:    openapi_types.Email(u.Email),
			Name:     &u.Name,
		})
	}

	return openapiUsers, nil
}

func (s *UserService) FindByID(ctx context.Context, id string) (*openapi.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, nil
	}

	return &openapi.User{
		Id:       &u.ID,
		Username: u.Username,
		Email:    openapi_types.Email(u.Email),
		Name:     &u.Name,
	}, nil
}

type TweetService struct {
	repo TweetRepository
}

func NewTweetService(repo TweetRepository) *TweetService {
	return &TweetService{repo}
}

func (s *TweetService) Create(ctx context.Context, t openapi.Tweet) error {
	tweet := Tweet{
		Content: t.Content,
		UserID:  t.User,
	}
	return s.repo.Create(ctx, tweet)
}

func (s *TweetService) FindAll(ctx context.Context) ([]openapi.Tweet, error) {
	tweets, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var openapiTweets []openapi.Tweet
	for _, t := range tweets {
		openapiTweets = append(openapiTweets, openapi.Tweet{
			Id:      &t.ID,
			Content: t.Content,
			User:    t.UserID,
		})
	}

	return openapiTweets, nil
}

func (s *TweetService) FindByID(ctx context.Context, id string) (*openapi.Tweet, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, nil
	}

	return &openapi.Tweet{
		Id:      &t.ID,
		Content: t.Content,
		User:    t.UserID,
	}, nil
}
