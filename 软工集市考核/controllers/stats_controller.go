package controllers

import (
	"net/http"
	"time"
	"vote-demo/database"
	"vote-demo/models"

	"github.com/gin-gonic/gin"
)

// GetPollStats 获取投票的详细统计信息
func GetPollStats(c *gin.Context) {
	pollID := c.Param("id")

	// 检查投票是否存在
	var poll models.Poll
	if err := database.DB.Preload("Options").First(&poll, "id = ?", pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "投票不存在"})
		return
	}

	// 获取投票总数
	var totalVotes int
	database.DB.Model(&models.Vote{}).Where("poll_id = ?", pollID).Count(&totalVotes)

	// 获取参与投票的用户数
	var uniqueUsers []string
	rows, err := database.DB.Model(&models.Vote{}).
		Where("poll_id = ?", pollID).
		Select("DISTINCT user_id").
		Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		rows.Scan(&userID)
		uniqueUsers = append(uniqueUsers, userID)
	}

	// 获取每个选项的投票数和百分比
	type OptionStat struct {
		ID        string  `json:"id"`
		Text      string  `json:"text"`
		Count     int     `json:"count"`
		Percentage float64 `json:"percentage"`
	}

	var optionStats []OptionStat
	for _, option := range poll.Options {
		var count int
		database.DB.Model(&models.Vote{}).Where("option_id = ?", option.ID).Count(&count)
		
		percentage := 0.0
		if totalVotes > 0 {
			percentage = float64(count) / float64(totalVotes) * 100
		}
		
		optionStats = append(optionStats, OptionStat{
			ID:         option.ID,
			Text:       option.Text,
			Count:      count,
			Percentage: percentage,
		})
	}

	// 获取投票的时间分布
	type TimeDistribution struct {
		Hour  int `json:"hour"`
		Count int `json:"count"`
	}

	var timeDistribution []TimeDistribution
	for i := 0; i < 24; i++ {
		timeDistribution = append(timeDistribution, TimeDistribution{
			Hour:  i,
			Count: 0,
		})
	}

	rows, err = database.DB.Model(&models.Vote{}).
		Where("poll_id = ?", pollID).
		Select("strftime('%H', created_at) as hour, COUNT(*) as count").
		Group("hour").
		Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var hour string
			var count int
			rows.Scan(&hour, &count)
			
			hourInt := 0
			if h, err := time.Parse("15", hour); err == nil {
				hourInt = h.Hour()
			}
			
			if hourInt >= 0 && hourInt < 24 {
				timeDistribution[hourInt].Count = count
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"poll":              poll,
		"total_votes":       totalVotes,
		"unique_voters":     len(uniqueUsers),
		"option_stats":      optionStats,
		"time_distribution": timeDistribution,
	})
}

// GetTrendingPolls 获取热门投票
func GetTrendingPolls(c *gin.Context) {
	// 获取过去7天内的热门投票（按投票数排序）
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	
	type PollWithVoteCount struct {
		models.Poll
		VoteCount int `json:"vote_count"`
	}
	
	var trendingPolls []PollWithVoteCount
	
	// 查询过去7天内有投票记录的投票
	rows, err := database.DB.Table("votes").
		Select("poll_id, COUNT(*) as vote_count").
		Where("created_at > ?", sevenDaysAgo).
		Group("poll_id").
		Order("vote_count DESC").
		Limit(10).
		Rows()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门投票失败"})
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var pollID string
		var voteCount int
		rows.Scan(&pollID, &voteCount)
		
		var poll models.Poll
		if err := database.DB.Preload("Options").First(&poll, "id = ?", pollID).Error; err == nil {
			trendingPolls = append(trendingPolls, PollWithVoteCount{
				Poll:      poll,
				VoteCount: voteCount,
			})
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"trending_polls": trendingPolls,
	})
}

// GetUserStats 获取用户的投票统计信息
func GetUserStats(c *gin.Context) {
	userID := c.Param("id")
	
	// 检查用户是否存在
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	
	// 获取用户参与的投票数
	var participatedPollCount int
	rows, err := database.DB.Model(&models.Vote{}).
		Where("user_id = ?", userID).
		Select("DISTINCT poll_id").
		Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}
	defer rows.Close()
	
	var participatedPollIDs []string
	for rows.Next() {
		var pollID string
		rows.Scan(&pollID)
		participatedPollIDs = append(participatedPollIDs, pollID)
	}
	participatedPollCount = len(participatedPollIDs)
	
	// 获取用户的投票总数
	var totalVoteCount int
	database.DB.Model(&models.Vote{}).Where("user_id = ?", userID).Count(&totalVoteCount)
	
	// 获取用户最近的投票记录
	var recentVotes []models.Vote
	database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(10).
		Find(&recentVotes)
	
	// 获取这些投票记录对应的投票和选项信息
	var recentVoteDetails []gin.H
	for _, vote := range recentVotes {
		var poll models.Poll
		var option models.Option
		
		database.DB.First(&poll, "id = ?", vote.PollID)
		database.DB.First(&option, "id = ?", vote.OptionID)
		
		recentVoteDetails = append(recentVoteDetails, gin.H{
			"vote_id":     vote.ID,
			"poll_id":     vote.PollID,
			"poll_title":  poll.Title,
			"option_id":   vote.OptionID,
			"option_text": option.Text,
			"created_at":  vote.CreatedAt,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"user":                  user,
		"participated_polls":    participatedPollCount,
		"total_votes":           totalVoteCount,
		"recent_vote_details":   recentVoteDetails,
	})
} 