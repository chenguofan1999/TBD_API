package model

import (
	"database/sql"
	"fmt"
	"time"
)

// CreateCommentTableIfNotExists Creates a Contents Table If Not Exists
func CreateCommentTableIfNotExists() {
	sql := `CREATE TABLE IF NOT EXISTS comments(
		comment_id INT NOT NULL AUTO_INCREMENT,
		commenter_id INT,
		content_id INT,
		reply_to INT DEFAULT NULL,
		comment_text TEXT,
		create_time BIGINT,
		PRIMARY KEY (comment_id),
		FOREIGN KEY (commenter_id) REFERENCES users(user_id),
		FOREIGN KEY (content_id) REFERENCES contents(content_id),
		FOREIGN KEY (reply_to) REFERENCES comments(comment_id)
		); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("Create comment table failed", err)
		return
	}
	fmt.Println("Create comment table successed or it already exists")
}

// InsertComment is for test use
func InsertComment(username string, contentID int, text string) {
	commenterID := QueryUserWithName(username).UserID
	result, err := DB.Exec("insert into comments(commenter_id,content_id,comment_text,create_time) values(?,?,?,?)", commenterID, contentID, text, time.Now().Unix())
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return
	}

	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID
	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)
		return
	}
	fmt.Println("Insert comment id:", lastInsertID)
}

// InsertReply is for test use
func InsertReply(username string, replyTo int, text string) {
	commenterID := QueryUserWithName(username).UserID
	contentID := QueryContentIDwithCommentID(replyTo)
	result, err := DB.Exec("insert into comments(commenter_id,content_id,reply_to,comment_text,create_time) values(?,?,?,?,?)", commenterID, contentID, replyTo, text, time.Now().Unix())
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return
	}

	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID
	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)
		return
	}
	fmt.Println("Insert comment id:", lastInsertID, "(a reply)")
}

func QueryContentIDwithCommentID(commentID int) int {
	var contentID int
	row := DB.QueryRow("select content_id from comments where comment_id = ?", commentID)
	row.Scan(&contentID)
	return contentID
}

/*
type MiniUser struct {
	Username  string `json:"username" form:"username"`
	Bio       string `json:"bio" form:"bio"`
	AvatarURL string `json:"avatar_url" form:"avatar_url"`
}
*/

type Comment struct {
	CommentID int      `json:"commentID" form:"commentID"`
	ContentID int      `json:"contentID" form:"contentID"`
	Text      string   `json:"text" form:"text"`
	Time      int64    `json:"createTime" form:"createTime"`
	Commenter MiniUser `json:"commenter" form:"commenter"`
}

func QueryCommentsWithContentID(contentID int) []Comment {
	comments := make([]Comment, 0)
	rows, err := DB.Query(`select comment_id,content_id,reply_to,comment_text,create_time,username,bio,avatar_url
		from comments,users where commenter_id = user_id and content_id = ?`, contentID)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var comment Comment
		var nullableReplyTo sql.NullInt32
		err = rows.Scan(&comment.CommentID, &comment.ContentID, &nullableReplyTo, &comment.Text, &comment.Time,
			&comment.Commenter.Username, &comment.Commenter.Bio, &comment.Commenter.AvatarURL)
		if err != nil {
			panic(err)
		}

		// reply_to = NULL, 是条评论
		if nullableReplyTo.Valid == false {
			comments = append(comments, comment)
		}
	}
	return comments
}

type Reply struct {
	ReplyID int      `json:"replyID" form:"replyID"`
	ReplyTo int      `json:"replyTo" form:"replyTo"`
	Text    string   `json:"text" form:"text"`
	Time    int64    `json:"createTime" form:"createTime"`
	Replier MiniUser `json:"Replier" form:"Replier"`
}

func QueryReplyWithID(replyTo int) []Reply {
	replies := make([]Reply, 0)
	rows, err := DB.Query(`select comment_id,reply_to,comment_text,create_time,username,bio,avatar_url
		from comments,users where commenter_id=user_id and reply_to = ?`, replyTo)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var reply Reply
		err = rows.Scan(&reply.ReplyID, &reply.ReplyTo, &reply.Text, &reply.Time, &reply.Replier.Username, &reply.Replier.Bio, &reply.Replier.AvatarURL)
		if err != nil {
			panic(err)
		}
		replies = append(replies, reply)
	}

	return replies
}
