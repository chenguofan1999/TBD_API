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
