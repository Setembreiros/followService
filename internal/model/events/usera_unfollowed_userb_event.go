package events

type UserAUnfollowedUserBEvent struct {
	FollowerID string `json:"followerId"`
	FolloweeID string `json:"followeeId"`
}
