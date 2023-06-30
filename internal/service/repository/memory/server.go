package memory

import "github.com/ricleal/twitter-clone/internal/service"

type Server struct {
	Tweets []service.Tweet
	Users  []service.User
}

func New() *Server {
	return &Server{}
}
