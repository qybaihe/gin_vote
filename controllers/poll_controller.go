package controllers

import (
	"net/http"
	"time"
	"vote-demo/database"
	"vote-demo/models"

	"github.com/gin-gonic/gin"
)

// CreatePoll 创建新投票
func CreatePoll(c *gin.Context) {
	var input struct {
		Title       string        `json:"title" binding:"required"`
		Description string        `json:"description"`
		Type        string        `json:"type" binding:"required"`
		Options     []string      `json:"options" binding:"required,min=2"`
		EndTime     time.Time     `json:"end_time"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证投票类型
	if input.Type != models.PollTypeBinary && 
	   input.Type != models.PollTypeSingle && 
	   input.Type != models.PollTypeMulti {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的投票类型"})
		return
	}

	// 如果是二分选项类型，强制设置为两个选项：是/否
	if input.Type == models.PollTypeBinary {
		input.Options = []string{"是", "否"}
	}

	// 创建投票
	poll := models.Poll{
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		EndTime:     input.EndTime,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := database.DB.Create(&poll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建投票失败"})
		return
	}

	// 创建选项
	for _, optionText := range input.Options {
		option := models.Option{
			PollID:    poll.ID,
			Text:      optionText,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := database.DB.Create(&option).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建选项失败"})
			return
		}
	}

	// 重新查询完整的投票信息（包括选项）
	var result models.Poll
	database.DB.Preload("Options").First(&result, "id = ?", poll.ID)

	c.JSON(http.StatusCreated, result)
}

// GetPoll 获取投票详情
func GetPoll(c *gin.Context) {
	id := c.Param("id")

	var poll models.Poll
	if err := database.DB.Preload("Options").First(&poll, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "投票不存在"})
		return
	}

	c.JSON(http.StatusOK, poll)
}

// ListPolls 获取投票列表
func ListPolls(c *gin.Context) {
	var polls []models.Poll
	database.DB.Preload("Options").Find(&polls)

	c.JSON(http.StatusOK, polls)
}

// UpdatePoll 更新投票信息
func UpdatePoll(c *gin.Context) {
	id := c.Param("id")

	var poll models.Poll
	if err := database.DB.First(&poll, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "投票不存在"})
		return
	}

	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		EndTime     time.Time `json:"end_time"`
		IsActive    *bool     `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.Title != "" {
		updates["title"] = input.Title
	}
	
	if input.Description != "" {
		updates["description"] = input.Description
	}
	
	if !input.EndTime.IsZero() {
		updates["end_time"] = input.EndTime
	}
	
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if err := database.DB.Model(&poll).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新投票失败"})
		return
	}

	// 返回更新后的投票
	database.DB.Preload("Options").First(&poll, "id = ?", id)
	c.JSON(http.StatusOK, poll)
}

// DeletePoll 删除投票
func DeletePoll(c *gin.Context) {
	id := c.Param("id")

	var poll models.Poll
	if err := database.DB.First(&poll, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "投票不存在"})
		return
	}

	// 删除相关的投票记录
	database.DB.Where("poll_id = ?", id).Delete(&models.Vote{})
	
	// 删除相关的选项
	database.DB.Where("poll_id = ?", id).Delete(&models.Option{})
	
	// 删除投票
	database.DB.Delete(&poll)

	c.JSON(http.StatusOK, gin.H{"message": "投票已删除"})
}

// GetPollResults 获取投票结果
func GetPollResults(c *gin.Context) {
	id := c.Param("id")

	var poll models.Poll
	if err := database.DB.Preload("Options").First(&poll, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "投票不存在"})
		return
	}

	// 获取每个选项的投票数
	type OptionResult struct {
		ID    string `json:"id"`
		Text  string `json:"text"`
		Count int    `json:"count"`
	}

	var results []OptionResult
	for _, option := range poll.Options {
		var count int
		database.DB.Model(&models.Vote{}).Where("option_id = ?", option.ID).Count(&count)
		
		results = append(results, OptionResult{
			ID:    option.ID,
			Text:  option.Text,
			Count: count,
		})
	}

	// 获取总投票数
	var totalVotes int
	database.DB.Model(&models.Vote{}).Where("poll_id = ?", id).Count(&totalVotes)

	c.JSON(http.StatusOK, gin.H{
		"poll":       poll,
		"results":    results,
		"total_votes": totalVotes,
	})
} 