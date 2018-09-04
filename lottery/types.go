package main

const TableLottery = "lotteries"
const TableLuckyUser = "lucky_users"
const TablePartner = "partners"
const LotteryPrefix = "lottery:"
const DescPrefix = "descriptions"

// BotConfig 机器人配置
type BotConfig struct {
	Http struct {
		Proxy   string `yaml:proxy`
		Timeout int    `yaml:timeout`
	}
	Bot struct {
		Token  string `yaml:token`
		Poller int    `yaml:poller`
	}
}

// Lottery 抽奖活动
type Lottery struct {
	ID          string `json:id`
	Description string `json:description`
	Owner       string `json:owner`
	Quota       int    `json:quota`
	Created     int64  `json:created`
	IsRoll      bool   `json:is_roll`
}

// LuckyUser 中奖名单
type LuckyUser struct {
	UserName string `json:username`
	LID      string `json:lid`
	Created  int    `json:created`
}

// Partner 参与者列表
type Partner struct {
	LID   string   `json:lid`
	Title string   `json:title`
	Users []string `json:users`
}
