package leaderboard

import (
	"math/rand"
	"sync"
	"time"
)

const (
	MaxLevel = 16   // 最大跳表层数
	P        = 0.25 // 概率因子
)

// Node 结构体表示跳表的节点，包含玩家信息、分数、时间戳和各层的前向指针
type Node struct {
	player    *Player   // 玩家信息
	score     int       // 玩家分数
	timestamp time.Time // 玩家得分时间戳
	forward   []*Node   // 各层的前向指针
}

// LeaderboardSkipList 结构体表示使用跳表实现的排行榜
// 包含读写锁用于并发控制，跳表头节点，当前最大层数，以及玩家ID到跳表节点的映射
type LeaderboardSkipList struct {
	mu        sync.RWMutex
	header    *Node            // 头节点
	level     int              // 当前最大层数
	playerMap map[string]*Node // 玩家ID到节点的映射
}

func NewLeaderboardSkipList() *LeaderboardSkipList {
	header := &Node{forward: make([]*Node, MaxLevel)}
	return &LeaderboardSkipList{
		header:    header,
		level:     1,
		playerMap: make(map[string]*Node),
	}
}

// 生成随机层数
func randomLevel() int {
	level := 1
	for ; level < MaxLevel && rand.Float32() < P; level++ {
	}
	return level
}

// UpdateScore 更新分数,跳表插入
// 如果玩家已经存在于排行榜中，先删除旧记录，然后插入新记录，保持跳表的有序性
func (l *LeaderboardSkipList) UpdateScore(playerID string, score int, timestamp time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 删除旧记录
	if node, exists := l.playerMap[playerID]; exists {
		l.deleteNode(node)
		delete(l.playerMap, playerID)
	}

	// 创建新节点
	newNode := &Node{
		player:    &Player{playerID, score, timestamp},
		score:     score,
		timestamp: timestamp,
		forward:   make([]*Node, randomLevel()),
	}

	// 查找插入位置
	update := make([]*Node, MaxLevel) // 记录每一层在插入新节点时，需要更新其 forward 指针的前一个节点
	current := l.header
	for i := l.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && (current.forward[i].score > score || (current.forward[i].score == score && current.forward[i].timestamp.Before(timestamp))) {
			current = current.forward[i]
		}
		update[i] = current
	}

	// 插入新节点并更新各层指针
	for i := 0; i < len(newNode.forward) && update[i] != nil; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	// 更新当前最大层数
	if len(newNode.forward) > l.level {
		l.level = len(newNode.forward)
	}

	l.playerMap[playerID] = newNode
}

// 删除节点（内部使用）
// 从跳表的最高层开始，逐层查找要删除的节点，记录每一层需要更新的前一个节点在 update 切片中。
// 遍历每一层，将前一个节点的 forward 指针指向要删除节点的下一个节点，从而将该节点从跳表中移除。
func (l *LeaderboardSkipList) deleteNode(node *Node) {
	update := make([]*Node, MaxLevel)
	current := l.header
	for i := l.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i] != node {
			current = current.forward[i]
		}
		update[i] = current
	}

	for i := 0; i < l.level; i++ {
		if update[i].forward[i] == node {
			update[i].forward[i] = node.forward[i]
		}
	}
}

// GetPlayerRank 获取玩家排名（跳表遍历计算）
func (l *LeaderboardSkipList) GetPlayerRank(playerID string) (RankInfo, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 检查玩家是否存在于排行榜中
	if node, exists := l.playerMap[playerID]; exists {
		rank := 1
		current := l.header.forward[0]
		// 遍历跳表底层链表，统计排名
		for current != nil && current != node {
			rank++
			current = current.forward[0]
		}
		return RankInfo{playerID, node.score, rank}, true // 返回排名信息
	}
	return RankInfo{}, false
}

// GetTopN 获取TopN（跳表底层遍历）
func (l *LeaderboardSkipList) GetTopN(n int) []RankInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var res []RankInfo
	current := l.header.forward[0]
	i := 0
	// 从跳表的底层链表开始遍历，添加前n名玩家信息
	for current != nil && i < n {
		res = append(res, RankInfo{
			PlayerID: current.player.PlayerID,
			Score:    current.score,
			Rank:     i + 1,
		})
		current = current.forward[0]
		i++
	}
	return res
}

// GetPlayerRankRange 获取周边排名（跳表底层遍历+排名计算）
func (l *LeaderboardSkipList) GetPlayerRankRange(playerID string, rangeN int) []RankInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 检查玩家是否存在于排行榜中
	if node, exists := l.playerMap[playerID]; exists {
		// 计算当前排名
		rank := 1
		current := l.header.forward[0]
		for current != nil && current != node {
			rank++
			current = current.forward[0]
		}

		start := max(1, rank-rangeN)                 // 计算排名范围的起始位置
		end := min(l.getTotalPlayers(), rank+rangeN) // 计算排名范围的结束位置

		var res []RankInfo
		current = l.header.forward[0]
		// 遍历跳表底层链表，获取排名范围内的玩家信息
		for i := 0; current != nil && i < end; current = current.forward[0] {
			if i+1 >= start {
				res = append(res, RankInfo{
					PlayerID: current.player.PlayerID,
					Score:    current.score,
					Rank:     i + 1,
				})
			}
			i++
		}
		return res
	}
	return nil
}

// 获取总玩家数
func (l *LeaderboardSkipList) getTotalPlayers() int {
	count := 0
	current := l.header.forward[0]
	for current != nil {
		count++
		current = current.forward[0]
	}
	return count
}
