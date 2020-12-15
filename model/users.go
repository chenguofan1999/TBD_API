package model

import (
	"errors"
	"fmt"
)

type User struct {
	UserID    int    `json:"userID" form:"userID"`
	Username  string `json:"username" form:"username"`
	Bio       string `json:"bio" form:"bio"`
	AvatarURL string `json:"avatar" form:"avatar"`
	Followers int    `json:"followerNum" form:"followerNum"`
	Following int    `json:"followingNum" form:"followingNum"`
}

// CreateUserTableIfNotExists Creates a Users Table If Not Exists
func CreateUserTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS users(
		user_id INT NOT NULL AUTO_INCREMENT,
		username VARCHAR(32) UNIQUE,
		password VARCHAR(32),
		bio VARCHAR(64) DEFAULT '',
		avatar_url VARCHAR(128) DEFAULT '',
		follower_num INT DEFAULT 0,
		following_num INT DEFAULT 0,
		PRIMARY KEY (user_id)
		); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create table failed", err)
		return
	}
	fmt.Println("create user table successed or it already exists")
}

// InsertUser is for test use
func InsertUser(username string, password string) error {
	if username == "" || password == "" {
		return errors.New("Invalid string")
	}

	result, err := DB.Exec("insert INTO users(username,password) values(?,?)", username, password)
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return errors.New("User exists")
	}
	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID
	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)
	}

	fmt.Println("Insert data id:", lastInsertID)
	return nil
}

func QueryUserWithName(username string) *User {
	user := new(User)

	row := DB.QueryRow("select user_id,username,bio,avatar_url,follower_num,following_num from users where username = ?", username)
	//注意一一对应
	err := row.Scan(&user.UserID, &user.Username, &user.Bio, &user.AvatarURL, &user.Followers, &user.Following)

	if err != nil {
		return nil
	}

	return user
}

func QueryPasswordByName(username string) string {
	row := DB.QueryRow("select password from users where username = ?", username)
	var password string
	if err := row.Scan(&password); err != nil {
		return ""
	}
	return password
}

func QueryUserIDWithName(username string) int {
	row := DB.QueryRow("select user_id from users where username = ?", username)
	var userID int
	if err := row.Scan(&userID); err != nil {
		return 0
	}
	return userID
}

func UpdateBio(userID int, newBio string) error {
	result, _ := DB.Exec("update users set bio = ? where user_id = ?", newBio, userID)

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user may not exist")
	}
	return nil
}

func UpdateAvatar(userID int, newAvatarURL string) error {
	result, _ := DB.Exec("update users set avatar_url = ? where user_id = ?", newAvatarURL, userID)

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user may not exist")
	}
	return nil
}
