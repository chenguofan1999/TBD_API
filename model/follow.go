package model

import "fmt"

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

// InsertFollowRelation 用于测试
func InsertFollowRelation(followerID int, followedID int) {
	_, err := DB.Exec("insert INTO follow(follower_id,followed_id) values(?,?)", followerID, followedID)
	if err != nil {
		fmt.Printf("Insert data failed, err:%v", err)
		return
	}

	_, err = DB.Exec("update users set follower_num=follower_num+1 where user_id=?", followedID)
	if err != nil {
		fmt.Printf("Increment follower_num failed, err:%v", err)
		return
	}

	_, err = DB.Exec("update users set following_num=following_num+1 where user_id=?", followerID)
	if err != nil {
		fmt.Printf("Increment following_num failed, err:%v", err)
		return
	}

	fmt.Println(followerID, "follows", followedID)
}

func InsertFollowRelationByName(followerName string, followedName string) {
	var followerID int
	var followedID int

	row := DB.QueryRow("select user_id from users where username = ?", followerName)
	row.Scan(&followerID)

	row = DB.QueryRow("select user_id from users where username = ?", followedName)
	row.Scan(&followedID)

	InsertFollowRelation(followerID, followedID)
}

func QueryFollowersWithName(username string) []User {
	followers := make([]User, 0)

	follower_ids, err := DB.Query("select follower_id from users,follow where user_id = followed_id and username = ?", username)

	if err != nil {
		panic(err)
	}

	for follower_ids.Next() {
		var user User
		err = follower_ids.Scan(&user.UserID)
		row := DB.QueryRow("select username,bio,avatar_url,follower_num,following_num from users where user_id = ?", user.UserID)
		row.Scan(&user.Username, &user.Bio, &user.AvatarURL, &user.Followers, &user.Following)
		followers = append(followers, user)
	}

	return followers
}

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
