package memory

import (
	"sync"

	"github.com/ricleal/twitter-clone/internal/service/repository"
)

// Handler is a mock repository handler that stores data in memory.
type Handler struct {
	Tweets []repository.Tweet
	Users  []repository.User
	m      sync.Mutex
}
