package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	tb "gopkg.in/tucnak/telebot.v2"
)

var markdownOption = &tb.SendOptions{
	ParseMode: "Markdown",
}

// CreateLottery 创建一次抽奖活动
// 创建奖励必须是隐私模式
func CreateLottery(m *tb.Message) {
	if !m.Private() {
		return
	}
	values := strings.Split(m.Text, "\n")
	if len(values) != 2 {
		bot.Send(m.Sender, "*请输入正确的创建抽奖活动格式*", markdownOption)
		return
	}

	number := strings.Split(values[0], " ")[1]

	log.Println("strings: ", number, values[1])

	quota, err := strconv.Atoi(number)

	if err != nil {
		bot.Send(m.Sender, "*请输入正确的中奖人数数字格式*", markdownOption)
		return
	}

	lid := UUID()

	// 保存数据库
	key := LotteryPrefix + lid
	lottery := Lottery{
		ID:          lid,
		Description: values[1],
		Owner:       m.Sender.Username,
		Quota:       quota,
		Created:     time.Now().Unix(),
	}

	data, _ := json.Marshal(lottery)
	value := string(data)

	log.Printf("key: %s, value: %s\n", key, value)

	ldb.Set([]byte(key), []byte(value))

	msg := "活动创建成功\n抽奖序号为: " + lid + "\n将该序号发送给参与者输入 /join [序号] 即可参与"

	bot.Send(m.Sender, msg)
}

// MyLottery 我发起的所有抽奖活动
func MyLottery(m *tb.Message) {

}

// ShowLottery 显示抽奖活动详情
func ShowLottery(m *tb.Message) {

}

// StopLottery 暂停某个抽奖活动
func StopLottery(m *tb.Message) {
	checked(m)
}

// DeleteLottery 删除某个抽奖活动
func DeleteLottery(m *tb.Message) {
	checked(m)
}

// JoinLottery 参与某次抽奖活动
func JoinLottery(m *tb.Message) {
	if flag, lottery := checked(m); flag {
		log.Println("lottery:", lottery)
		if lottery.IsRoll {
			bot.Send(m.Sender, "*该活动已经开奖，您无法参与*", markdownOption)
			return
		}
		key := TablePartner + lottery.ID

		ldb.LPush([]byte(key), []byte(m.Sender.Username))

		bot.Send(m.Sender, "*参与成功，等待开奖吧*", markdownOption)
	}

	// 1. 检测是否已开奖
	// 2. 检测当前用户是否已经参与
	// 3. 加入参与者列表
}

// RollLottery 开奖
func RollLottery(m *tb.Message) {
	checked(m)
	// 2. 检测当前用户是否是创建者
	// 3. 检测抽奖是否已经开启
	// 4. 随机分配 N 个名额, 如果参与人数 <= 分配人数, 则取总参与人数
}

// ShowPartner 显示参与者
// 显示参与者必须是隐私模式，创建者除外
func ShowPartner(m *tb.Message) {
	checked(m)
}

// LuckyList 显示中奖列表
// 显示中奖列表必须是隐私模式，创建者除外
func LuckyList(m *tb.Message) {
	checked(m)
}

// Help 显示帮助信息
func Help(m *tb.Message) {
	bot.Send(m.Sender, "你好，欢迎使用本机器人\n\n 点击 /help - 获取帮助\n")
}

var base64Table = [64]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
}

// UUID 生成一个UUID
func UUID() string {
	v4, _ := uuid.NewV4()
	src := v4.Bytes()
	var dst = make([]byte, 16)

	for i, v := range src {
		dst[i] = base64Table[v>>2]
	}
	return string(dst)
}

func checked(m *tb.Message) (bool, *Lottery) {
	p := &Lottery{}
	if m.Payload == "" {
		bot.Send(m.Sender, "*请输入正确的指令*", markdownOption)
		return false, p
	}
	log.Printf("收到 [%s] 的消息[%d] %s\n", m.Sender.Username, m.Chat.ID, m.Text)
	key := LotteryPrefix + m.Payload
	value, err := ldb.Get([]byte(key))
	if err != nil || len(value) == 0 {
		bot.Send(m.Sender, "*不存在该抽奖活动*", markdownOption)
		return false, p
	}

	err = json.Unmarshal(value, p)
	if err != nil {
		log.Println("error:", err)
		return false, p
	}
	return true, p
}
