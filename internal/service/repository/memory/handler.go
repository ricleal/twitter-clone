package memory

import (
	"sync"

	"github.com/ricleal/twitter-clone/internal/service/repository"
)

type Handler struct {
	Tweets []repository.Tweet
	Users  []repository.User
	m      sync.Mutex
}
