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
