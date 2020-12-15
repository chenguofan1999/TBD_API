package model

import (
	"database/sql"
	"fmt"
	"time"
)

// GeneralComment 可同时表示 Comment 或 Reply
type GeneralComment struct {
	CommentID int      `json:"commentID" form:"commentID"`
	ContentID int      `json:"contentID" form:"contentID"`
	IsReply   bool     `json:"isReply" form:"isReply"`
	ReplyTo   int      `json:"replyTo" form:"replyTo"`
	Text      string   `json:"text" form:"text"`
	Time      int64    `json:"createTime" form:"createTime"`
	Creator   MiniUser `json:"creator" form:"creator"`
}

/*
type MiniUser struct {
	Username  string `json:"username" form:"username"`
	Bio       string `json:"bio" form:"bio"`
	AvatarURL string `json:"avatar_url" form:"avatar_url"`
}


// Comment means only the direct comment of a content
type Comment struct {
	CommentID int      `json:"commentID" form:"commentID"`
	ContentID int      `json:"contentID" form:"contentID"`
	Text      string   `json:"text" form:"text"`
	Time      int64    `json:"createTime" form:"createTime"`
	Commenter MiniUser `json:"commenter" form:"commenter"`
}

// Reply is a reply to a comment or another Reply
type Reply struct {
	CommentID int      `json:"commentID" form:"commentID"`
	ReplyTo   int      `json:"replyTo" form:"replyTo"`
	Text      string   `json:"text" form:"text"`
	Time      int64    `json:"createTime" form:"createTime"`
	Replier   MiniUser `json:"replier" form:"replier"`
}
*/

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

// QueryContentIDwithCommentID 根据 commentID 得到对应的 content 的 contentID
func QueryContentIDwithCommentID(commentID int) int {
	var contentID int
	row := DB.QueryRow("select content_id from comments where comment_id = ?", commentID)
	row.Scan(&contentID)
	return contentID
}

// QueryRepliesWithCommentID 得到一条 comment 的全部回复
func QueryRepliesWithCommentID(commentID int) []GeneralComment {
	replies := make([]GeneralComment, 0)
	rows, err := DB.Query(`select comment_id,content_id,reply_to,comment_text,create_time,username,bio,avatar_url
		from comments,users where commenter_id=user_id and reply_to = ?`, commentID)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var reply GeneralComment
		err = rows.Scan(&reply.CommentID, &reply.ContentID, &reply.ReplyTo, &reply.Text, &reply.Time, &reply.Creator.Username, &reply.Creator.Bio, &reply.Creator.AvatarURL)
		if err != nil {
			panic(err)
		}

		reply.IsReply = true
		replies = append(replies, reply)
	}

	return replies
}

//QueryCommentWithCommentID 根据 commentID 得到 comment
func QueryCommentWithCommentID(commentID int) *GeneralComment {
	comment := new(GeneralComment)
	var nullableReplyTo sql.NullInt32
	row := DB.QueryRow(`select comment_id,content_id,reply_to,comment_text,create_time,username,bio,avatar_url
	from comments,users where commenter_id = user_id and comment_id = ?`, commentID)

	err := row.Scan(&comment.CommentID, &comment.ContentID, &nullableReplyTo, &comment.Text, &comment.Time,
		&comment.Creator.Username, &comment.Creator.Bio, &comment.Creator.AvatarURL)
	if err != nil {
		panic(err)
	}

	if nullableReplyTo.Valid {
		comment.IsReply = true
		comment.ReplyTo = int(nullableReplyTo.Int32)
	} else {
		comment.IsReply = false
		comment.ReplyTo = 0
	}

	return comment
}

// QueryAllCommentsWithContentID 根据 contentID 得到全部的 comment，包含 comment 的回复
func QueryAllCommentsWithContentID(contentID int) []GeneralComment {
	comments := make([]GeneralComment, 0)
	rows, err := DB.Query(`select comment_id,content_id,reply_to,comment_text,create_time,username,bio,avatar_url
	from comments,users where commenter_id = user_id and content_id = ?`, contentID)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var comment GeneralComment
		var nullableReplyTo sql.NullInt32

		err = rows.Scan(&comment.CommentID, &comment.ContentID, &nullableReplyTo, &comment.Text, &comment.Time,
			&comment.Creator.Username, &comment.Creator.Bio, &comment.Creator.AvatarURL)
		if err != nil {
			panic(err)
		}

		if nullableReplyTo.Valid {
			comment.IsReply = true
			comment.ReplyTo = int(nullableReplyTo.Int32)
		} else {
			comment.IsReply = false
			comment.ReplyTo = 0
		}
		comments = append(comments, comment)
	}

	return comments
}

// QueryCommentsWithContentID 根据 contentID 得到 comment，不包含 comment 的回复
func QueryCommentsWithContentID(contentID int) []GeneralComment {
	comments := make([]GeneralComment, 0)
	rows, err := DB.Query(`select comment_id,content_id,reply_to,comment_text,create_time,username,bio,avatar_url
		from comments,users where commenter_id = user_id and content_id = ?`, contentID)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var comment GeneralComment
		var nullableReplyTo sql.NullInt32
		err = rows.Scan(&comment.CommentID, &comment.ContentID, &nullableReplyTo, &comment.Text, &comment.Time,
			&comment.Creator.Username, &comment.Creator.Bio, &comment.Creator.AvatarURL)
		if err != nil {
			panic(err)
		}

		// reply_to = NULL, 是条评论
		if nullableReplyTo.Valid == false {
			comment.IsReply = false
			comment.ReplyTo = 0
			comments = append(comments, comment)
		}
	}
	return comments
}
