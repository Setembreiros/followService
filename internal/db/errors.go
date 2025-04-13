package database

import "fmt"

type NotRelationshipCreatedError struct {
	followerId string
	followeeId any
}

func (e *NotRelationshipCreatedError) Error() string {
	errorMessage := fmt.Sprintf("No relationship created, %s -> %s", e.followerId, e.followeeId)
	return errorMessage
}

func NewNotRelationshipCreatedError(followerId, followeeId string) *NotRelationshipCreatedError {
	return &NotRelationshipCreatedError{
		followerId: followerId,
		followeeId: followeeId,
	}
}

type CleanDatabaseError struct{}

func (e *CleanDatabaseError) Error() string {
	errorMessage := fmt.Sprint("Unexpected error cleaning the database")
	return errorMessage
}

func NewCleanDatabaseError() *CleanDatabaseError {
	return &CleanDatabaseError{}
}
