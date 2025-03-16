package database

//go:generate mockgen -source=database.go -destination=mock/database.go

type TableAttributes struct {
	Name          string
	AttributeType string
}

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
