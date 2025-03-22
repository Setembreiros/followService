package model

type UserPairRelationship struct {
	FollowerID string `json:"followerId"`
	FolloweeID string `json:"followeeId"`
}
