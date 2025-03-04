package controllers

import (
	"net/http"
	"time"
	"vote-demo/database"
	"vote-demo/models"

	"github.com/gin-gonic/gin"
)

// AddOption 添加选项
func AddOption(c *gin.Context) {
	pollID := c.Param("id")

	// 检查投票是否存在
	var poll models.Poll
	if err := database.DB.First(&poll, "id = ?", pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "投票不存在"})
		return
	}

	// 检查投票是否为二分类型，二分类型不允许添加选项
	if poll.Type == models.PollTypeBinary {
		c.JSON(http.StatusBadRequest, gin.H{"error": "二分类型投票不允许添加选项"})
		return
	}

	var input struct {
		Text string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建新选项
	option := models.Option{
		PollID:    pollID,
		Text:      input.Text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&option).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建选项失败"})
		return
	}

	c.JSON(http.StatusCreated, option)
}

// UpdateOption 更新选项
func UpdateOption(c *gin.Context) {
	optionID := c.Param("option_id")

	var option models.Option
	if err := database.DB.First(&option, "id = ?", optionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "选项不存在"})
		return
	}

	// 检查投票是否为二分类型，二分类型不允许修改选项
	var poll models.Poll
	if err := database.DB.First(&poll, "id = ?", option.PollID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票信息失败"})
		return
	}

	if poll.Type == models.PollTypeBinary {
		c.JSON(http.StatusBadRequest, gin.H{"error": "二分类型投票不允许修改选项"})
		return
	}

	var input struct {
		Text string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新选项
	updates := map[string]interface{}{
		"text":       input.Text,
		"updated_at": time.Now(),
	}

	if err := database.DB.Model(&option).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新选项失败"})
		return
	}

	// 返回更新后的选项
	database.DB.First(&option, "id = ?", optionID)
	c.JSON(http.StatusOK, option)
}

// DeleteOption 删除选项
func DeleteOption(c *gin.Context) {
	optionID := c.Param("option_id")

	var option models.Option
	if err := database.DB.First(&option, "id = ?", optionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "选项不存在"})
		return
	}

	// 检查投票是否为二分类型，二分类型不允许删除选项
	var poll models.Poll
	if err := database.DB.First(&poll, "id = ?", option.PollID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取投票信息失败"})
		return
	}

	if poll.Type == models.PollTypeBinary {
		c.JSON(http.StatusBadRequest, gin.H{"error": "二分类型投票不允许删除选项"})
		return
	}

	// 检查剩余选项数量，至少保留两个选项
	var count int
	database.DB.Model(&models.Option{}).Where("poll_id = ?", option.PollID).Count(&count)
	if count <= 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "投票至少需要两个选项"})
		return
	}

	// 删除相关的投票记录
	database.DB.Where("option_id = ?", optionID).Delete(&models.Vote{})
	
	// 删除选项
	database.DB.Delete(&option)

	c.JSON(http.StatusOK, gin.H{"message": "选项已删除"})
} 