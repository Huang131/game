# game

游戏排行榜系统
背景
你正在为一个多人在线游戏开发一个全服排行榜系统，该系统需要追踪玩家获得的积分，实时更新排名，并提供排行榜实时查询功能。
请根据个人技术特长 ，编程语言不限 ，优先使用 Go。
如无不便 ，请使用 github.com，gitlab.com 或 gitee.com 提交答案 ， 回复公开代码库地址即可。
一、基础功能实现
请实现一个基础的排行榜系统，包含以下功能：
1.支持更新玩家积分。
2.查询玩家当前排名。
3.获取前 N 名玩家的分数和名次。
4.查询自己名次前后共 N 名玩家的分数和名次。
排名规则如下：
1.分数从高到低进行排序。
2.如果多位玩家分数相同 ，则先得到该分数的玩家排在前面。
接口定义示例

```TypeScript
interface LeaderboardService {
    // 更新玩家分数
    updateScore(playerId: string, score: number, timestamp: number): void;
    // 获取玩家当前排名
    getPlayerRank(playerId: string): RankInfo;
    //获取排行榜前N名
    getTopN(n: number): Array<RankInfo>;
    // 获取玩家周边排名
    getPlayerRankRange(playerId: string, range: number): Array<RankInfo>;
}
```

二、系统设计
可以使用UML图或线框图辅助表达设计
**可靠性要求**
考虑到排行榜系统需要 7*24 小时运行，且需要确保数据的一致性和完整性，请说明你会如何设计系统来满足可靠性要求。

**性能要求**
考虑到排行榜系统需要实时查询和更新，且总玩家数量可达到百万级，请说明你会如何设计系统来满足性能要求。


三、游戏需求更改（选做）
假如游戏设计师想让获得相同成绩的玩家享有同等的荣誉，所以计划调整为采用"密集排名"的计算方式，请说明你会如何修改实现来满足新的游戏需求。
新的排名规则如下：
1.分数依旧从高到低排序。
2.对于分数相同的玩家，将获得完全相同的排名位次。
3.当出现新的不同分数时，该玩家的排名将在上一个排名基础上递增 1 位。
举例说明：
Plain Text - 玩家A：100分 → 排名第1 - 玩家B：100分 → 排名第1 - 玩家C：95分 → 排名第2 - 玩家D：95分 → 排名第2 - 玩家E：90分 → 排名第3

数据结构调整
新增分数到排名组的映射：map [int] RankGroup，记录该分数对应的所有玩家和起始排名
RankGroup 结构包含：分数、起始排名、玩家列表
排名计算逻辑
更新时维护分数有序列表（去重后按分数降序）
计算排名时：找到所有分数高于当前分数的唯一分数集合，其数量 + 1 即为当前排名
示例：分数列表 [100,9

// 密集排名计算
func (l *Leaderboard) getDenseRank(score int) int {
uniqueScores := make(map[int]bool)
for e := l.players.Front(); e != nil; e = e.Next() {
p := e.Value.(*Player)
if p.Score > score {
uniqueScores[p.Score] = true
}
}
return len(uniqueScores) + 1
}
