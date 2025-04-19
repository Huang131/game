package leaderboard

// getDenseRanks 计算密集排名
func (l *LeaderboardLinkedList) getDenseRanks() map[string]int {
	ranks := make(map[string]int)
	currentRank := 0
	prevScore := -1

	for e := l.players.Front(); e != nil; e = e.Next() {
		p := e.Value.(*Player)
		if p.Score != prevScore {
			currentRank = currentRank + 1
			prevScore = p.Score
		}
		ranks[p.PlayerID] = currentRank
	}
	return ranks
}

// GetDensePlayerRank 获取玩家排名 链表遍历计算
// 如果玩家存在于排行榜中，返回其排名信息和 true；否则返回空的排名信息和 false
func (l *LeaderboardLinkedList) GetDensePlayerRank(playerID string) RankInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 检查玩家是否存在于排行榜中
	if elem, exists := l.playerMap[playerID]; exists {
		ranks := l.getDenseRanks()
		p := elem.Value.(*Player)

		return RankInfo{
			PlayerID: playerID,
			Score:    p.Score,
			Rank:     ranks[playerID],
		}
	}

	// 玩家不存在，返回空的排名信息
	return RankInfo{}
}

// GetDenseTopN 获取TopN 链表头部遍历
func (l *LeaderboardLinkedList) GetDenseTopN(n int) []RankInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()

	ranks := l.getDenseRanks()
	var res []RankInfo
	i := 0
	// 遍历链表，获取前 N 名玩家信息
	for e := l.players.Front(); e != nil && i < n; e = e.Next() {
		p := e.Value.(*Player)
		res = append(res, RankInfo{p.PlayerID, p.Score, ranks[p.PlayerID]})
		i++
	}
	return res
}

// GetDensePlayerRankRange 获取周边排名（链表二次遍历）
// 返回指定玩家前后各 rangeN 名玩家的排名信息
func (l *LeaderboardLinkedList) GetDensePlayerRankRange(playerID string, rangeN int) []RankInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()
	// 检查玩家是否存在于排行榜中
	if _, exists := l.playerMap[playerID]; exists {
		ranks := l.getDenseRanks()
		currentRank := ranks[playerID]

		start := max(1, currentRank-rangeN)             // 计算排名范围的起始位置
		end := min(l.players.Len(), currentRank+rangeN) // 计算排名范围的结束位置

		var res []RankInfo
		rankCount := 0
		// 遍历链表，获取排名范围内的玩家信息
		for e := l.players.Front(); e != nil; e = e.Next() {
			p := e.Value.(*Player)
			playerRank := ranks[p.PlayerID]
			if playerRank >= start && playerRank <= end {
				res = append(res, RankInfo{p.PlayerID, p.Score, playerRank})
				rankCount++
			}
		}
		return res
	}
	return nil
}
