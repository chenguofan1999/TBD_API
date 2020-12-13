package model

import "fmt"

type User struct {
	UserID    int    `json:"UserID" form:"UserID"`
	Username  string `json:"username" form:"username"`
	Password  string `json:"password" form:"password"`
	Bio       string `json:"bio" form:"bio"`
	AvatarURL string `json:"avatar_url" form:"avatar_url"`
	Followers int    `json:"followers" form:"followers"`
	Following int    `json:"following" form:"following"`
}

// CreateUserTableIfNotExists Creates a Users Table If Not Exists
func CreateUserTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS users(
		user_id INT NOT NULL AUTO_INCREMENT,
		username VARCHAR(32) UNIQUE,
		password VARCHAR(32),
		bio VARCHAR(64) DEFAULT '',
		avatar_url VARCHAR(128) DEFAULT '',
		followerCount INT DEFAULT 0,
		followingCount INT DEFAULT 0,
		PRIMARY KEY (user_id)
		); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create table failed", err)
		return
	}
	fmt.Println("create user table successed or it already exists")
}

func InsertUser(user User) {
	result, err := DB.Exec("insert INTO users(username,password) values(?,?)", user.Username, user.Password)
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return
	}
	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID
	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)
		return
	}
	fmt.Println("Insert data id:", lastInsertID)
}

func QueryWithName(username string) *User {
	user := new(User)

	row := DB.QueryRow("select user_id,username,bio,avatar_url,followerCount,followingCount from users where username = ?", username)
	//注意一一对应
	err := row.Scan(&user.UserID, &user.Username, &user.Bio, &user.AvatarURL, &user.Followers, &user.Following)

	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return nil
	}
	fmt.Println("Single row data:", *user)

	return user
}
