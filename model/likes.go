package model

import (
	"errors"
	"fmt"
)

// CreateLikeTableIfNotExists Create Like Table If Not Exists
func CreateLikeTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS likes(
		user_id INT,
		content_id INT,
		PRIMARY KEY (user_id, content_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (content_id) REFERENCES contents(content_id) ON DELETE CASCADE
		)ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create likes table failed", err)
		return
	}
	fmt.Println("create likes table successed or it already exists")
}

// InsertLikeRelation Inserts Like Relation
func InsertLikeRelation(userID int, contentID int) error {
	_, err := DB.Exec("insert into likes(user_id,content_id) values(?,?)", userID, contentID)
	if err != nil {
		return errors.New("user not exists or content not exists or duplicate like")
	}

	DB.Exec("update contents set like_num=like_num+1 where content_id=?", contentID)
	return nil
}

// InsertLikeRelationUsingName Insert Like Relation Using Name
func InsertLikeRelationUsingName(username string, contentID int) error {
	userID := QueryUserIDWithName(username)
	return InsertLikeRelation(userID, contentID)
}

func DeleteLikeRelation(userID int, contentID int) error {
	result, _ := DB.Exec("delete from likes where user_id = ? and content_id = ?", userID, contentID)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("Did not like it")
	}

	DB.Exec("update contents set like_num=like_num-1 where content_id=?", contentID)
	return nil

}

func QueryLikedContentsWithUserID(userID int) ([]Content, error) {
	fmt.Println("Querying liked contents with userID")

	rows, err := DB.Query("select contents.content_id,content_title,content_text,create_time,username,bio,avatar_url from likes,contents,users where likes.content_id = contents.content_id and author_id = users.user_id and likes.user_id = ? order by create_time desc", userID)
	if err != nil {
		return []Content{}, err
	}

	return GetContentsFromRows(rows)
}

func QueryHasLiked(userID int, contentID int) (bool, error) {
	// check userID valid
	var temp int
	row := DB.QueryRow("select user_id from users where user_id = ?", userID)
	err := row.Scan(&temp)
	if err != nil {
		return false, errors.New("no such user")
	}

	// check contentID valid
	row = DB.QueryRow("select content_id from contents where content_id = ?", contentID)
	err = row.Scan(&temp)
	if err != nil {
		return false, errors.New("no such content")
	}

	row = DB.QueryRow("select 1 from likes where user_id = ? and content_id = ?", userID, contentID)
	err = row.Scan(&temp)
	if err != nil {
		return false, nil
	}

	return true, nil
}
