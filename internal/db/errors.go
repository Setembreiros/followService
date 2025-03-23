package database

import "fmt"

type CleanDatabaseError struct{}

type NotRelationshipCreatedError struct {
	followerId string
	followeeId string
}

type NotRelationshipDeletedError struct {
	followerId string
	followeeId string
}

func NewCleanDatabaseError() *CleanDatabaseError {
	return &CleanDatabaseError{}
}

func NewNotRelationshipCreatedError(followerId, followeeId string) *NotRelationshipCreatedError {
	return &NotRelationshipCreatedError{
		followerId: followerId,
		followeeId: followeeId,
	}
}

func NewNotRelationshipDeletedError(followerId, followeeId string) *NotRelationshipDeletedError {
	return &NotRelationshipDeletedError{
		followerId: followerId,
		followeeId: followeeId,
	}
}

func (e *CleanDatabaseError) Error() string {
	errorMessage := "Unexpected error cleaning the database"
	return errorMessage
}

func (e *NotRelationshipCreatedError) Error() string {
	errorMessage := fmt.Sprintf("No relationship created, %s -> %s", e.followerId, e.followeeId)
	return errorMessage
}

func (e *NotRelationshipDeletedError) Error() string {
	errorMessage := fmt.Sprintf("No relationship deleted, %s -> %s", e.followerId, e.followeeId)
	return errorMessage
}
