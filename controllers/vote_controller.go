package controllers

import (
	"net/http"
	"time"
	"vote-demo/database"
	"vote-demo/models"

	"github.com/gin-gonic/gin"
)

// CastVote 进行投票
func CastVote(c *gin.Context) {
	pollID := c.Param("id")

	// 检查投票是否存在
	var poll models.Poll
	if err := database.DB.First(&poll, "id = ?", pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "投票不存在"})
		return
	}

	// 检查投票是否已结束
	if !poll.EndTime.IsZero() && poll.EndTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "投票已结束"})
		return
	}

	// 检查投票是否活跃
	if !poll.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "投票已关闭"})
		return
	}

	// 获取用户ID（在实际应用中，这应该从认证中间件获取）
	userID := c.GetHeader("User-ID")
	if userID == "" {
		// 为了演示，如果没有提供用户ID，我们创建一个临时用户
		user := models.User{
			Username:  "anonymous_" + time.Now().Format("20060102150405"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时用户失败"})
			return
		}
		userID = user.ID
	}

	var input struct {
		OptionIDs []string `json:"option_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证选项是否属于该投票
	for _, optionID := range input.OptionIDs {
		var option models.Option
		if err := database.DB.Where("id = ? AND poll_id = ?", optionID, pollID).First(&option).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "选项不存在或不属于该投票"})
			return
		}
	}

	// 根据投票类型验证选项数量
	switch poll.Type {
	case models.PollTypeBinary, models.PollTypeSingle:
		if len(input.OptionIDs) != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "该投票类型只允许选择一个选项"})
			return
		}
	case models.PollTypeMulti:
		// 多选类型允许多个选项，无需额外验证
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的投票类型"})
		return
	}

	// 检查用户是否已经投过票
	var existingVotes []models.Vote
	database.DB.Where("poll_id = ? AND user_id = ?", pollID, userID).Find(&existingVotes)

	// 如果用户已经投过票，根据投票类型处理
	if len(existingVotes) > 0 {
		// 对于单选和二分类型，删除之前的投票
		if poll.Type == models.PollTypeBinary || poll.Type == models.PollTypeSingle {
			for _, vote := range existingVotes {
				database.DB.Delete(&vote)
			}
		} else {
			// 对于多选类型，检查是否重复投票
			for _, optionID := range input.OptionIDs {
				for _, vote := range existingVotes {
					if vote.OptionID == optionID {
						c.JSON(http.StatusBadRequest, gin.H{"error": "您已经为该选项投过票"})
						return
					}
				}
			}
		}
	}

	// 创建投票记录
	var votes []models.Vote
	for _, optionID := range input.OptionIDs {
		vote := models.Vote{
			PollID:    pollID,
			OptionID:  optionID,
			UserID:    userID,
			CreatedAt: time.Now(),
		}
		if err := database.DB.Create(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "投票失败"})
			return
		}
		votes = append(votes, vote)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "投票成功",
		"votes":   votes,
	})
}

// GetUserVotes 获取用户在特定投票中的投票记录
func GetUserVotes(c *gin.Context) {
	pollID := c.Param("id")
	userID := c.GetHeader("User-ID")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供用户ID"})
		return
	}

	var votes []models.Vote
	database.DB.Where("poll_id = ? AND user_id = ?", pollID, userID).Find(&votes)

	// 获取选项详情
	var options []models.Option
	var optionIDs []string
	for _, vote := range votes {
		optionIDs = append(optionIDs, vote.OptionID)
	}

	if len(optionIDs) > 0 {
		database.DB.Where("id IN (?)", optionIDs).Find(&options)
	}

	c.JSON(http.StatusOK, gin.H{
		"votes":   votes,
		"options": options,
	})
} 