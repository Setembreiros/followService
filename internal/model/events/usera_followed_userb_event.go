package events

type UserAFollowedUserBEvent struct {
	FollowerID string `json:"followerId"`
	FolloweeID string `json:"followeeId"`
}
