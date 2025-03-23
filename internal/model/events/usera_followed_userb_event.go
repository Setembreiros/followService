package events

type UserAFollowedUserBEvent struct {
	FollowerID string `json:"followerId"`
	FolloweeID string `json:"followeeId"`
}

type UserAUnfollowedUserBEvent struct {
	FollowerID string `json:"followerId"`
	FolloweeID string `json:"followeeId"`
}
