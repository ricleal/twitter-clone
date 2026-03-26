package memory

import (
	"fmt"

	memdb "github.com/hashicorp/go-memdb"
)

// userRecord is the internal storage format for users in go-memdb.
type userRecord struct {
	ID       string
	Username string
	Email    string
	Name     string
}

// tweetRecord is the internal storage format for tweets in go-memdb.
type tweetRecord struct {
	ID      string
	Content string
	UserID  string
}

// Table names used as keys throughout the memory store.
const (
	tableUsers  = "users"
	tableTweets = "tweets"
)

// NewDB creates a new in-memory database with the twitter-clone schema.
func NewDB() (*memdb.MemDB, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tableUsers: {
				Name: tableUsers,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"username": {
						Name:    "username",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Username"},
					},
					"email": {
						Name:    "email",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
				},
			},
			tableTweets: {
				Name: tableTweets,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
				},
			},
		},
	}
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create in-memory database: %w", err)
	}
	return db, nil
}
