package memory

import "github.com/ricleal/twitter-clone/internal/service/repository"

type Server struct {
	Tweets []repository.Tweet
	Users  []repository.User
}

func New() *Server {
	return &Server{}
}
