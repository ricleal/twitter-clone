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

var dbSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"users": {
			Name: "users",
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
		"tweets": {
			Name: "tweets",
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

// NewDB creates a new in-memory database with the twitter-clone schema.
func NewDB() (*memdb.MemDB, error) {
	db, err := memdb.NewMemDB(dbSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create in-memory database: %w", err)
	}
	return db, nil
}
