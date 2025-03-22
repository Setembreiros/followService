package follow_user_artist

type UserPairFollowRelation struct {
	FollowerID string `json:"followerId"`
	FolloweeID string `json:"followeeId"`
}

type UserAFollowedUserBEvent struct {
	FollowerID string `json:"followerId"`
	FolloweeID string `json:"followeeId"`
}
