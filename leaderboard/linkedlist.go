package leaderboard

import (
	"container/list"
	"sync"
	"time"
)

// LeaderboardLinkedList 排行榜（链表实现）
type LeaderboardLinkedList struct {
	mu        sync.RWMutex
	players   *list.List               // 按Score降序、Timestamp升序排列的双向链表
	playerMap map[string]*list.Element // 玩家ID到链表节点的映射
}

func NewLeaderboardLinkedList() *LeaderboardLinkedList {
	return &LeaderboardLinkedList{
		players:   list.New(),
		playerMap: make(map[string]*list.Element),
	}
}

// UpdateScore 更新分数（链表插入）
// 如果玩家已经存在于排行榜中，先删除旧记录，然后插入新记录，保持链表的有序性
func (l *LeaderboardLinkedList) UpdateScore(playerID string, score int, timestamp time.Time) {
	// 加锁，防止并发需改
	l.mu.Lock()
	defer l.mu.Unlock()

	// 删除旧记录
	if elem, exists := l.playerMap[playerID]; exists {
		l.players.Remove(elem)
		delete(l.playerMap, playerID)
	}

	newPlayer := &Player{
		PlayerID:  playerID,
		Score:     score,
		Timestamp: timestamp,
	}
	// 遍历链表,链表按分数降序、时间戳升序排列
	for e := l.players.Front(); e != nil; e = e.Next() {
		p := e.Value.(*Player)
		if p.Score > score || (p.Score == score && p.Timestamp.Before(timestamp)) {
			continue
		}
		l.players.InsertBefore(newPlayer, e)
		l.playerMap[playerID] = e.Prev()
		return
	}

	// 遍历完链表都没有找到合适的插入位置，将新节点插入到链表尾部
	l.players.PushBack(newPlayer)
	l.playerMap[playerID] = l.players.Back()
}

// GetPlayerRank 获取玩家排名 链表遍历计算
// 如果玩家存在于排行榜中，返回其排名信息和 true；否则返回空的排名信息和 false
func (l *LeaderboardLinkedList) GetPlayerRank(playerID string) (RankInfo, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 检查玩家是否存在于排行榜中
	if elem, exists := l.playerMap[playerID]; exists {
		rank := 1
		// 遍历链表，统计排名
		for e := l.players.Front(); e != elem; e = e.Next() {
			rank++
		}
		p := elem.Value.(*Player)

		return RankInfo{
			PlayerID: playerID,
			Score:    p.Score,
			Rank:     rank,
		}, true
	}

	// 玩家不存在，返回空的排名信息和 false
	return RankInfo{}, false
}

// GetTopN 获取TopN 链表头部遍历
func (l *LeaderboardLinkedList) GetTopN(n int) []RankInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var res []RankInfo
	i := 0
	// 遍历链表，获取前 N 名玩家信息
	for e := l.players.Front(); e != nil && i < n; e = e.Next() {
		p := e.Value.(*Player)
		res = append(res, RankInfo{p.PlayerID, p.Score, i + 1})
		i++
	}
	return res
}

// GetPlayerRankRange 获取周边排名（链表二次遍历）
// 返回指定玩家前后各 rangeN 名玩家的排名信息
func (l *LeaderboardLinkedList) GetPlayerRankRange(playerID string, rangeN int) []RankInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()
	// 检查玩家是否存在于排行榜中
	if elem, exists := l.playerMap[playerID]; exists {
		currentRank := 1
		// 遍历链表，统计指定玩家的排名
		for e := l.players.Front(); e != elem; e = e.Next() {
			currentRank++
		}

		start := max(1, currentRank-rangeN)             // 计算排名范围的起始位置
		end := min(l.players.Len(), currentRank+rangeN) // 计算排名范围的结束位置

		var res []RankInfo
		i := 0
		// 遍历链表，获取排名范围内的玩家信息
		for e := l.players.Front(); e != nil && i < end; e = e.Next() {
			if i+1 >= start {
				p := e.Value.(*Player)
				res = append(res, RankInfo{p.PlayerID, p.Score, i + 1})
			}
			i++
		}
		return res
	}
	// 玩家不存在
	return nil
}
