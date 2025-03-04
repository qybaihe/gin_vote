package controllers

import (
	"net/http"
	"time"
	"vote-demo/database"
	"vote-demo/models"

	"github.com/gin-gonic/gin"
)

// AddComment 添加评论
func AddComment(c *gin.Context) {
	pollID := c.Param("id")

	// 检查投票是否存在
	var poll models.Poll
	if err := database.DB.First(&poll, "id = ?", pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "投票不存在"})
		return
	}

	// 获取用户ID
	userID := c.GetHeader("User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供用户ID"})
		return
	}

	// 检查用户是否存在
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var input struct {
		Content  string `json:"content" binding:"required"`
		ParentID string `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果提供了父评论ID，检查父评论是否存在
	if input.ParentID != "" {
		var parentComment models.Comment
		if err := database.DB.First(&parentComment, "id = ?", input.ParentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "父评论不存在"})
			return
		}

		// 确保父评论属于同一个投票
		if parentComment.PollID != pollID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "父评论不属于该投票"})
			return
		}
	}

	// 创建评论
	comment := models.Comment{
		PollID:    pollID,
		UserID:    userID,
		Content:   input.Content,
		ParentID:  input.ParentID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	// 返回创建的评论，包括用户信息
	database.DB.Preload("User").First(&comment, "id = ?", comment.ID)

	c.JSON(http.StatusCreated, comment)
}

// GetPollComments 获取投票的评论
func GetPollComments(c *gin.Context) {
	pollID := c.Param("id")

	// 检查投票是否存在
	var poll models.Poll
	if err := database.DB.First(&poll, "id = ?", pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "投票不存在"})
		return
	}

	// 获取所有评论
	var comments []models.Comment
	database.DB.Where("poll_id = ?", pollID).
		Preload("User").
		Order("created_at DESC").
		Find(&comments)

	// 构建评论树
	var rootComments []models.Comment
	commentMap := make(map[string]*models.Comment)

	// 首先将所有评论放入map中
	for i := range comments {
		commentMap[comments[i].ID] = &comments[i]
	}

	// 然后构建评论树
	for i := range comments {
		comment := &comments[i]
		if comment.ParentID == "" {
			// 这是一个根评论
			rootComments = append(rootComments, *comment)
		} else {
			// 这是一个回复
			if parent, exists := commentMap[comment.ParentID]; exists {
				parent.Replies = append(parent.Replies, *comment)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": rootComments,
	})
}

// UpdateComment 更新评论
func UpdateComment(c *gin.Context) {
	commentID := c.Param("comment_id")

	// 检查评论是否存在
	var comment models.Comment
	if err := database.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查用户是否是评论的作者
	userID := c.GetHeader("User-ID")
	if userID == "" || userID != comment.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权更新此评论"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新评论
	updates := map[string]interface{}{
		"content":    input.Content,
		"updated_at": time.Now(),
	}

	if err := database.DB.Model(&comment).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新评论失败"})
		return
	}

	// 返回更新后的评论
	database.DB.Preload("User").First(&comment, "id = ?", commentID)
	c.JSON(http.StatusOK, comment)
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	commentID := c.Param("comment_id")

	// 检查评论是否存在
	var comment models.Comment
	if err := database.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查用户是否是评论的作者
	userID := c.GetHeader("User-ID")
	if userID == "" || userID != comment.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此评论"})
		return
	}

	// 删除所有回复
	database.DB.Where("parent_id = ?", commentID).Delete(&models.Comment{})

	// 删除评论
	database.DB.Delete(&comment)

	c.JSON(http.StatusOK, gin.H{"message": "评论已删除"})
} 