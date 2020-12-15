package api

import (
	"net/http"
	"strconv"
	"tbd/model"

	"github.com/gin-gonic/gin"
)

// GetCommentsByContentIDandFilter get comments by contentID and filter
func GetCommentsByContentIDandFilter(c *gin.Context) {
	// 获取 contentID
	// 如果参数不能被转换为整型 ID, 返回400
	contentID, err := strconv.Atoi(c.Query("contentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}

	// 如果 Content 不存在, 返回404
	content := model.QueryContentWithContentID(contentID)
	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not found",
		})
		return
	}

	// 获取参数 filterReply, 默认过滤 reply
	filterReply := c.DefaultQuery("filterReply", "true")
	if filterReply == "true" {
		comments := model.QueryCommentsWithContentID(contentID)
		c.JSON(http.StatusOK, comments)
	} else if filterReply == "false" {
		comments := model.QueryAllCommentsWithContentID(contentID)
		c.JSON(http.StatusOK, comments)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
	}

}

// GetCommentByCommentID Get a Comment By CommentID
func GetCommentByCommentID(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("commentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}

	comment := model.QueryCommentWithCommentID(commentID)
	if comment == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not found",
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// GetRepliesByCommentID Get Replies By CommentID
func GetRepliesByCommentID(c *gin.Context) {
	// 获取 commentID
	// 如果参数不能被转换为整型 ID, 返回400
	commentID, err := strconv.Atoi(c.Param("commentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}

	comments := model.QueryRepliesWithCommentID(commentID)
	c.JSON(http.StatusOK, comments)
}

type postCommentFormat struct {
	ContentID int    `json:"contentID" form:"contentID"`
	Text      string `json:"text" form:"text"`
}

// PostComment : 当前用户发布一条评论(评论一条content)
func PostComment(c *gin.Context) {
	// 得到登录用户名
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)

	var input postCommentFormat
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "bind error"})
		return
	}

	if err := model.InsertComment(loginUserName, input.ContentID, input.Text); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type postReplyFormat struct {
	ReplyTo int    `json:"replyTo" form:"replyTo"`
	Text    string `json:"text" form:"text"`
}

// PostReply : 当前用户发表一条回复(回复一条 Comment)
func PostReply(c *gin.Context) {
	// 得到登录用户名
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)

	var input postReplyFormat
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "bind error"})
		return
	}
	if err := model.InsertReply(loginUserName, input.ReplyTo, input.Text); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DeleteComment : 删除Comment，在 URL 的 Path 中写入 commentID
func DeleteComment(c *gin.Context) {
	// 得到登录用户名
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	loginUserName := GetNameByToken(tokenString)

	// 获取 URL 中 commentID 并验证格式是否正确
	commentID, err := strconv.Atoi(c.Param("commentID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request"})
		return
	}

	// 验证 comment 是否存在，存在即获取
	comment := model.QueryCommentWithCommentID(commentID)
	if comment == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Bad Request"})
		return
	}

	// 验证是否为该用户所发
	if comment.Creator.Username != loginUserName {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}

	// 执行删除
	err = model.DeleteCommentWithCommentID(commentID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})

}
