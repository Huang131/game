package leaderboard

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const numPlayers = 1000

func BenchmarkLinkedListUpdateScore(b *testing.B) {
	lb := NewLeaderboardLinkedList()
	for i := 0; i < b.N; i++ {
		playerID := fmt.Sprintf("player%d", i%numPlayers)
		score := rand.Intn(1000)
		timestamp := time.Now()
		lb.UpdateScore(playerID, score, timestamp)
	}
}

func BenchmarkSkipListUpdateScore(b *testing.B) {
	lb := NewLeaderboardSkipList()
	for i := 0; i < b.N; i++ {
		playerID := fmt.Sprintf("player%d", i%numPlayers)
		score := rand.Intn(1000)
		timestamp := time.Now()
		lb.UpdateScore(playerID, score, timestamp)
	}
}

func BenchmarkLinkedListGetPlayerRank(b *testing.B) {
	lb := NewLeaderboardLinkedList()
	for i := 0; i < numPlayers; i++ {
		playerID := fmt.Sprintf("player%d", i)
		score := rand.Intn(1000)
		timestamp := time.Now()
		lb.UpdateScore(playerID, score, timestamp)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playerID := fmt.Sprintf("player%d", rand.Intn(numPlayers))
		lb.GetPlayerRank(playerID)
	}
}

func BenchmarkSkipListGetPlayerRank(b *testing.B) {
	lb := NewLeaderboardSkipList()
	for i := 0; i < numPlayers; i++ {
		playerID := fmt.Sprintf("player%d", i)
		score := rand.Intn(1000)
		timestamp := time.Now()
		lb.UpdateScore(playerID, score, timestamp)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playerID := fmt.Sprintf("player%d", rand.Intn(numPlayers))
		lb.GetPlayerRank(playerID)
	}
}

func BenchmarkLinkedListGetTopN(b *testing.B) {
	lb := NewLeaderboardLinkedList()
	for i := 0; i < numPlayers; i++ {
		playerID := fmt.Sprintf("player%d", i)
		score := rand.Intn(1000)
		timestamp := time.Now()
		lb.UpdateScore(playerID, score, timestamp)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb.GetTopN(10)
	}
}

func BenchmarkSkipListGetTopN(b *testing.B) {
	lb := NewLeaderboardSkipList()
	for i := 0; i < numPlayers; i++ {
		playerID := fmt.Sprintf("player%d", i)
		score := rand.Intn(1000)
		timestamp := time.Now()
		lb.UpdateScore(playerID, score, timestamp)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb.GetTopN(10)
	}
}

func BenchmarkLinkedListGetPlayerRankRange(b *testing.B) {
	lb := NewLeaderboardLinkedList()
	for i := 0; i < numPlayers; i++ {
		playerID := fmt.Sprintf("player%d", i)
		score := rand.Intn(1000)
		timestamp := time.Now()
		lb.UpdateScore(playerID, score, timestamp)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playerID := fmt.Sprintf("player%d", rand.Intn(numPlayers))
		lb.GetPlayerRankRange(playerID, 5)
	}
}

func BenchmarkSkipListGetPlayerRankRange(b *testing.B) {
	lb := NewLeaderboardSkipList()
	for i := 0; i < numPlayers; i++ {
		playerID := fmt.Sprintf("player%d", i)
		score := rand.Intn(1000)
		timestamp := time.Now()
		lb.UpdateScore(playerID, score, timestamp)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playerID := fmt.Sprintf("player%d", rand.Intn(numPlayers))
		lb.GetPlayerRankRange(playerID, 5)
	}
}

func TestLeaderboardLinkedList_DenseRank(t *testing.T) {
	leaderboard := NewLeaderboardLinkedList()

	// 更新玩家积分
	leaderboard.UpdateScore("A", 100, time.Now())
	leaderboard.UpdateScore("B", 100, time.Now().Add(-time.Second))
	leaderboard.UpdateScore("C", 95, time.Now())
	leaderboard.UpdateScore("D", 95, time.Now().Add(-time.Second))
	leaderboard.UpdateScore("E", 90, time.Now())

	// 测试 GetPlayerRank 方法
	testsGetPlayerRank := []struct {
		playerID string
		expected RankInfo
	}{
		{"A", RankInfo{PlayerID: "A", Score: 100, Rank: 1}},
		{"B", RankInfo{PlayerID: "B", Score: 100, Rank: 1}},
		{"C", RankInfo{PlayerID: "C", Score: 95, Rank: 2}},
		{"D", RankInfo{PlayerID: "D", Score: 95, Rank: 2}},
		{"E", RankInfo{PlayerID: "E", Score: 90, Rank: 3}},
	}

	for _, tt := range testsGetPlayerRank {
		result := leaderboard.GetDensePlayerRank(tt.playerID)
		if result.Rank != tt.expected.Rank || result.Score != tt.expected.Score {
			t.Errorf("GetDensePlayerRank(%s) = %+v; want %+v", tt.playerID, result, tt.expected)
		}
	}

	// 测试 GetTopN 方法
	topN := 3
	expectedTopN := []RankInfo{
		{PlayerID: "B", Score: 100, Rank: 1},
		{PlayerID: "A", Score: 100, Rank: 1},
		{PlayerID: "D", Score: 95, Rank: 2},
	}
	resultTopN := leaderboard.GetDenseTopN(topN)
	if len(resultTopN) != len(expectedTopN) {
		t.Errorf("GetDenseTopN(%d) length = %d; want %d", topN, len(resultTopN), len(expectedTopN))
	}
	for i := range resultTopN {
		if resultTopN[i].Rank != expectedTopN[i].Rank || resultTopN[i].Score != expectedTopN[i].Score {
			t.Errorf("GetDenseTopN(%d)[%d] = %+v; want %+v", topN, i, resultTopN[i], expectedTopN[i])
		}
	}

	// 测试 GetPlayerRankRange 方法
	rangeN := 1
	playerID := "C"
	expectedRange := []RankInfo{
		{PlayerID: "B", Score: 100, Rank: 1},
		{PlayerID: "A", Score: 100, Rank: 1},
		{PlayerID: "D", Score: 95, Rank: 2},
		{PlayerID: "C", Score: 95, Rank: 2},
		{PlayerID: "E", Score: 90, Rank: 3},
	}
	resultRange := leaderboard.GetDensePlayerRankRange(playerID, rangeN)
	if len(resultRange) != len(expectedRange) {
		t.Errorf("GetDensePlayerRankRange(%s, %d) length = %d; want %d", playerID, rangeN, len(resultRange), len(expectedRange))
	}
	for i := range resultRange {
		if resultRange[i].Rank != expectedRange[i].Rank || resultRange[i].Score != expectedRange[i].Score {
			t.Errorf("GetDensePlayerRankRange(%s, %d)[%d] = %+v; want %+v", playerID, rangeN, i, resultRange[i], expectedRange[i])
		}
	}
}
