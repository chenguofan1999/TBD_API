package model

import (
	"fmt"
	"time"
)

type MiniUser struct {
	Username  string `json:"username" form:"username"`
	Bio       string `json:"bio" form:"bio"`
	AvatarURL string `json:"avatar" form:"avatar"`
}

type Content struct {
	ContentID int      `json:"contentID" form:"contentID"`
	Title     string   `json:"title" form:"title"`
	Text      string   `json:"text" form:"text"`
	Time      int64    `json:"createTime" form:"createTime"`
	Author    MiniUser `json:"author" form:"author"`
	Images    []string `json:"images" form:"images"`
}

// CreateContentTableIfNotExists Creates a Contents Table If Not Exists
func CreateContentTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS contents(
		content_id INT NOT NULL AUTO_INCREMENT,
		author_id INT,
		content_title TINYTEXT,
		content_text TEXT,
		create_time BIGINT,
		like_num INT DEFAULT 0,
		comment_num INT DEFAULT 0,
		image_num INT DEFAULT 0,
		PRIMARY KEY (content_id),
		FOREIGN KEY (author_id) REFERENCES users(user_id)
		); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create content table failed", err)
		return
	}
	fmt.Println("Create content table successed or it already exists")
}

func CreateImageTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS images(
		image_id INT NOT NULL AUTO_INCREMENT,
		image_url VARCHAR(256),
		content_id INT NOT NULL,
		PRIMARY KEY (image_id),
		FOREIGN KEY (content_id) REFERENCES contents(content_id)
		); `
	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create image table failed", err)
		return
	}
	fmt.Println("Create image table successed or it already exists")
}

// InsertContent is for test use
func InsertTextContent(authorName string, title string, text string) {
	author := QueryUserWithName(authorName)
	author_id := author.UserID

	result, err := DB.Exec("insert into contents(author_id,content_title,content_text,create_time) values(?,?,?,?)", author_id, title, text, time.Now().Unix())
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return
	}

	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID
	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)
		return
	}
	fmt.Println("Insert content id:", lastInsertID)
}

func QueryContentsWithName(authorName string) []Content {
	contents := make([]Content, 0)
	rows, err := DB.Query("select content_id,content_title,content_text,create_time,username,bio,avatar_url from contents,users where author_id = user_id and username = ?", authorName)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var content Content
		err = rows.Scan(&content.ContentID, &content.Title, &content.Text, &content.Time, &content.Author.Username, &content.Author.Bio, &content.Author.AvatarURL)
		if err != nil {
			panic(err)
		}

		imageRows, err := DB.Query("select image_url from images where content_id = ?", content.ContentID)
		if err != nil {
			panic(err)
		}

		imageURLs := make([]string, 0)
		for imageRows.Next() {
			var imageURL string
			err = imageRows.Scan(&imageURL)
			if err != nil {
				panic(err)
			}

			imageURLs = append(imageURLs, imageURL)
		}
		content.Images = imageURLs
		contents = append(contents, content)
	}
	return contents
}
