package entity

import (
	"fmt"
	"time"
)

type UserID int64
type LeaderCardIdIDType int64
type UserTokenType string

type User struct {
	Id           UserID             `db:"id"`
	Name         string             `db:"name"`
	Token        UserTokenType      `db:"token"`
	LeaderCardId LeaderCardIdIDType `db:"leader_card_id"`
	CreatedAt    time.Time          `db:"created_at"`
	UpdatedAt    time.Time          `db:"updated_at"`
}

type UserValidationError struct {
	MemberName string
}

func (e *UserValidationError) Error() string {
	return fmt.Sprintf("user validation error: %s", e.MemberName)
}

func (u *User) ValidateNotEmpty() error {
	// TODO: `validate:required` でチェックしたほうが良いかもしれない

	var uZeroValue = User{}
	if u.Id == uZeroValue.Id {
		// MySQL の AUTO_INCREMENT は 1 スタート
		return &UserValidationError{MemberName: "Id"}
	}
	if u.Name == uZeroValue.Name {
		return &UserValidationError{MemberName: "Name"}
	}
	if u.Token == uZeroValue.Token {
		return &UserValidationError{MemberName: "Token"}
	}
	// Allow Leader Card ID to be zero value
	// if u.LeaderCardId == uZeroValue.LeaderCardId {
	// 	return &UserValidationError{MemberName: "LeaderCardId"}
	// }
	if u.CreatedAt == uZeroValue.CreatedAt {
		return &UserValidationError{MemberName: "LeaderCardId"}
	}
	if u.UpdatedAt == uZeroValue.UpdatedAt {
		return &UserValidationError{MemberName: "Modified"}
	}
	return nil
}
