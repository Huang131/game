package leaderboard

import "time"

// Player 表示玩家信息，包含唯一标识、分数和得分时间戳（用于处理同分情况）
type Player struct {
	PlayerID  string    // 玩家唯一ID
	Score     int       // 玩家当前分数  排名降序
	Timestamp time.Time // 得分时间戳,时间戳早的靠前
}

// RankInfo 表示排名信息，包含玩家ID、分数和具体排名
type RankInfo struct {
	PlayerID string `json:"playerId"` // 玩家唯一ID
	Score    int    `json:"score"`    // 玩家分数
	Rank     int    `json:"rank"`     // 玩家排名（从1开始）
}

type LeaderboardService interface {
	UpdateScore(playerID string, score int, timestamp time.Time) // 更新分数
	GetPlayerRank(playerID string) (RankInfo, bool)              // 获取个人排名
	GetTopN(n int) []RankInfo                                    // 获取前N名
	GetPlayerRankRange(playerID string, rangeN int) []RankInfo   // 获取周边排名
}

// max 函数用于返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min 函数用于返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
