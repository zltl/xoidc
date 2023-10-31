package m

import "github.com/bwmarrin/snowflake"

type UserID snowflake.ID

func (u UserID) Int64() int64 {
	return u.Int64()
}

func (u UserID) String() string {
	return snowflake.ID(u).Base64()
}

func UserIDFromInt64(i int64) UserID {
	return UserID(snowflake.ID(i))
}

func ParseUserID(uid string) (UserID, error) {
	if uid == "" {
		return UserID(0), nil
	}
	id, err := snowflake.ParseBase64(uid)
	return UserID(id), err
}
