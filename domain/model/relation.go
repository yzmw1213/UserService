package model

// Relation ユーザーのフォロー関係
type Relation struct {
	FollowerUserID uint32
	FollowedUserID uint32
}
