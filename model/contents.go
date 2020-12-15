package model

import (
	"errors"
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

// InsertContent :...
func InsertContent(authorName string, title string, text string, imageURLs []string) (int, error) {
	author := QueryUserWithName(authorName)
	if author == nil {
		return 0, errors.New("no such user")
	}

	authorID := author.UserID
	result, err := DB.Exec("insert into contents(author_id,content_title,content_text,create_time) values(?,?,?,?)", authorID, title, text, time.Now().Unix())
	if err != nil {
		return 0, errors.New("create text content error")
	}

	contentID, _ := result.LastInsertId() //获取插入数据的自增ID

	// imageURLs 是一个 imageURL 的切片
	for _, imageURL := range imageURLs {
		_, err := DB.Exec("insert into images(content_id,image_url) values(?,?)", int(contentID), imageURL)
		if err != nil {
			return 0, errors.New("insert image failed")
		}

		_, err = DB.Exec("update contents set image_num=image_num+1 where content_id = ?", int(contentID))
		if err != nil {
			return 0, errors.New("insert image failed")
		}
	}

	return int(contentID), nil
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

// QueryContentWithContentID Query a Content With ContentID
func QueryContentWithContentID(contentID int) *Content {
	content := new(Content)
	row := DB.QueryRow("select content_id,content_title,content_text,create_time,username,bio,avatar_url from contents,users where author_id = user_id and content_id = ?", contentID)

	err := row.Scan(&content.ContentID, &content.Title, &content.Text, &content.Time, &content.Author.Username, &content.Author.Bio, &content.Author.AvatarURL)
	if err != nil {
		panic(err)
	}

	imageRows, err := DB.Query("select image_url from images where content_id = ?", contentID)
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
	return content
}
