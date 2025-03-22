package database

//go:generate mockgen -source=database.go -destination=test/mock/database.go

type Database struct {
	Client DatabaseClient
}

type DatabaseClient interface {
}

func NewDatabase(client DatabaseClient) *Database {
	return &Database{
		Client: client,
	}
}
