package model

import (
	"database/sql"
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

func QueryContentsWithName(authorName string) ([]Content, error) {
	fmt.Println("Querying contents with name")

	rows, err := DB.Query("select content_id,content_title,content_text,create_time,username,bio,avatar_url from contents,users where author_id = user_id and username = ? order by content_id desc", authorName)
	if err != nil {
		return []Content{}, err
	}

	return GetContentsFromRows(rows)
}

func QueryPublicContents(num int) ([]Content, error) {
	fmt.Println("Querying public contents")

	rows, err := DB.Query("select content_id,content_title,content_text,create_time,username,bio,avatar_url from contents,users where author_id = user_id order by content_id desc limit ?", num)
	if err != nil {
		return []Content{}, err
	}

	return GetContentsFromRows(rows)
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

// DeleteContentWithContentID 删除一条内容，级联删除其所有评论
func DeleteContentWithContentID(contentID int) error {
	_, err := DB.Exec(`delete from contents where content_id = ?`, contentID)
	if err != nil {
		return errors.New("Content May Not Exist")
	}
	return nil
}

// GetContentsOfFollowingUsersWithUserID 获取指定ID的用户关注的用户的内容
func GetContentsOfFollowingUsersWithUserID(userID int) ([]Content, error) {
	fmt.Println("Querying contents of following users")

	rows, err := DB.Query("select content_id,content_title,content_text,create_time,username,bio,avatar_url  from contents,users where author_id = user_id and author_id in (select followed_id from follow where follower_id = ?)", userID)
	if err != nil {
		return []Content{}, err
	}

	return GetContentsFromRows(rows)
}

// GetContentsFromRows 是一个辅助函数
func GetContentsFromRows(rows *sql.Rows) ([]Content, error) {
	contents := make([]Content, 0)

	for rows.Next() {
		var content Content
		err := rows.Scan(&content.ContentID, &content.Title, &content.Text, &content.Time, &content.Author.Username, &content.Author.Bio, &content.Author.AvatarURL)
		if err != nil {
			return contents, err
		}

		imageRows, err := DB.Query("select image_url from images where content_id = ?", content.ContentID)
		if err != nil {
			return contents, err
		}

		imageURLs := make([]string, 0)
		for imageRows.Next() {
			var imageURL string
			err = imageRows.Scan(&imageURL)
			if err != nil {
				return contents, err
			}

			imageURLs = append(imageURLs, imageURL)
		}
		content.Images = imageURLs
		contents = append(contents, content)
	}
	return contents, nil
}
