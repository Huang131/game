package main

import (
	"fmt"
	"game/leaderboard"
	"time"
)

func main() {
	// 初始化链表和跳表排行榜
	linkedListLeaderboard := leaderboard.NewLeaderboardLinkedList()
	skipListLeaderboard := leaderboard.NewLeaderboardSkipList()

	// 定义测试数据
	testData := []struct {
		playerID  string
		score     int
		timestamp time.Time
	}{
		{"player1", 100, time.Now().Add(-time.Hour)},
		{"player2", 90, time.Now()},
		{"player3", 100, time.Now()},
		{"player4", 80, time.Now()},
	}

	// 更新分数
	for _, data := range testData {
		linkedListLeaderboard.UpdateScore(data.playerID, data.score, data.timestamp)
		skipListLeaderboard.UpdateScore(data.playerID, data.score, data.timestamp)
	}

	// 测试 GetPlayerRank
	playerID := "player2"
	linkedListRank, linkedListOk := linkedListLeaderboard.GetPlayerRank(playerID)
	skipListRank, skipListOk := skipListLeaderboard.GetPlayerRank(playerID)
	if linkedListOk != skipListOk || (linkedListOk && linkedListRank != skipListRank) {
		fmt.Printf("GetPlayerRank 结果不一致: 链表 %+v, 跳表 %+v\n", linkedListRank, skipListRank)
	} else {
		fmt.Printf("GetPlayerRank 结果一致: %+v\n", linkedListRank)
	}

	// 测试 GetTopN
	n := 2
	linkedListTopN := linkedListLeaderboard.GetTopN(n)
	skipListTopN := skipListLeaderboard.GetTopN(n)
	if !equalRankInfoSlices(linkedListTopN, skipListTopN) {
		fmt.Printf("GetTopN 结果不一致: 链表 %+v, 跳表 %+v\n", linkedListTopN, skipListTopN)
	} else {
		fmt.Printf("GetTopN 结果一致: %+v\n", linkedListTopN)
	}

	// 测试 GetPlayerRankRange
	rangeN := 1
	linkedListRankRange := linkedListLeaderboard.GetPlayerRankRange(playerID, rangeN)
	skipListRankRange := skipListLeaderboard.GetPlayerRankRange(playerID, rangeN)
	if !equalRankInfoSlices(linkedListRankRange, skipListRankRange) {
		fmt.Printf("GetPlayerRankRange 结果不一致: 链表 %+v, 跳表 %+v\n", linkedListRankRange, skipListRankRange)
	} else {
		fmt.Printf("GetPlayerRankRange 结果一致: %+v\n", linkedListRankRange)
	}
}

// equalRankInfoSlices 用于比较两个 RankInfo 切片是否相等
func equalRankInfoSlices(a, b []leaderboard.RankInfo) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
