package model

import (
	"errors"
	"fmt"
)

type Follow struct {
	FollowerID int `json:"followerID" form:"followerID"`
	FollowedID int `json:"followedID" form:"followedID"`
}

func CreateFollowTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS follow(
		follower_id INT,
		followed_id INT,
		PRIMARY KEY (follower_id, followed_id),
		FOREIGN KEY (follower_id) REFERENCES users(user_id),
		FOREIGN KEY (followed_id) REFERENCES users(user_id)
		); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create Follow table failed", err)
		return
	}
	fmt.Println("create Follow table successed or it already exists")
}

// InsertFollowRelation ：用户ID为 followerID 的用户 follow 用户ID为 followedID 的用户
func InsertFollowRelation(followerID int, followedID int) error {
	_, err := DB.Exec("insert INTO follow(follower_id,followed_id) values(?,?)", followerID, followedID)
	if err != nil {
		return errors.New("No such user, or following already")
	}

	if followerID == followedID {
		return errors.New("You can't follow yourself")
	}

	DB.Exec("update users set follower_num=follower_num+1 where user_id=?", followedID)
	DB.Exec("update users set following_num=following_num+1 where user_id=?", followerID)

	fmt.Println(followerID, "follows", followedID)
	return nil
}

// InsertFollowRelationByName : 用户名为 followerName 的用户 follow 用户名为 followedName 的用户
func InsertFollowRelationByName(followerName string, followedName string) error {
	followerID := QueryUserIDWithName(followerName)
	followedID := QueryUserIDWithName(followedName)
	return InsertFollowRelation(followerID, followedID)
}

// DeleteFollowRelation ：用户ID为 followerID 的用户 unfollow 用户ID为 followedID 的用户
func DeleteFollowRelation(followerID int, followedID int) error {
	result, err := DB.Exec("delete from follow where follower_id = ? and followed_id = ?", followerID, followedID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("No such user, or not following originally")
	}

	_, err = DB.Exec("update users set follower_num=follower_num-1 where user_id=?", followedID)
	if err != nil {
		return err
	}

	_, err = DB.Exec("update users set following_num=following_num-1 where user_id=?", followerID)
	if err != nil {
		return err
	}

	fmt.Println(followerID, "unfollows", followedID)
	return nil
}

// DeleteFollowRelationByName : 用户名为 followerName 的用户 unfollow 用户名为 followedName 的用户
func DeleteFollowRelationByName(followerName string, followedName string) error {
	var followerID int
	var followedID int

	row := DB.QueryRow("select user_id from users where username = ?", followerName)
	row.Scan(&followerID)

	row = DB.QueryRow("select user_id from users where username = ?", followedName)
	row.Scan(&followedID)

	return DeleteFollowRelation(followerID, followedID)
}

// QueryFollowersWithName : 根据用户名查询关注者
func QueryFollowersWithName(username string) []User {
	followers := make([]User, 0)

	followerIDs, err := DB.Query("select follower_id from users,follow where user_id = followed_id and username = ?", username)

	if err != nil {
		panic(err)
	}

	for followerIDs.Next() {
		var user User
		err = followerIDs.Scan(&user.UserID)
		row := DB.QueryRow("select username,bio,avatar_url,follower_num,following_num from users where user_id = ?", user.UserID)
		row.Scan(&user.Username, &user.Bio, &user.AvatarURL, &user.Followers, &user.Following)
		followers = append(followers, user)
	}

	return followers
}

// QueryFollowingWithName : 根据用户名查询TA关注的人
func QueryFollowingWithName(username string) []User {
	followers := make([]User, 0)

	followerIDs, err := DB.Query("select followed_id from users,follow where user_id = follower_id and username = ?", username)

	if err != nil {
		panic(err)
	}

	for followerIDs.Next() {
		var user User
		err = followerIDs.Scan(&user.UserID)
		row := DB.QueryRow("select username,bio,avatar_url,follower_num,following_num from users where user_id = ?", user.UserID)
		row.Scan(&user.Username, &user.Bio, &user.AvatarURL, &user.Followers, &user.Following)
		followers = append(followers, user)
	}

	return followers
}

func QueryHasFollowed(followerID int, followedID int) (bool, error) {
	// check followerID valid
	var temp int
	row := DB.QueryRow("select user_id from users where user_id = ?", followerID)
	err := row.Scan(&temp)
	if err != nil {
		return false, errors.New("no such user")
	}

	// check followedID valid
	row = DB.QueryRow("select user_id from users where user_id = ?", followedID)
	err = row.Scan(&temp)
	if err != nil {
		return false, errors.New("no such user")
	}

	row = DB.QueryRow("select 1 from follow where follower_id=? and followed_id=?", followerID, followedID)
	err = row.Scan(&temp)
	if err != nil {
		return false, nil
	}

	return true, nil
}
