package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/lunny/nodb"
	"github.com/lunny/nodb/config"
	tb "gopkg.in/tucnak/telebot.v2"
)

/*
create - 创建抽奖活动
stop - 关闭抽奖活动. 如 /stop [抽奖活动唯一ID]
join - 参与抽奖活动. 如 /join [抽奖活动唯一ID]
roll - 开奖, 仅发起者可开启. 如 /roll [抽奖活动唯一ID]
partner - 查看参与者列表, 如 /partner [抽奖活动唯一ID]
lucky_list - 查看中奖列表, 如 /lucky_list [抽奖活动唯一ID]
*/
var bot *tb.Bot
var ldb *nodb.DB
var conf *BotConfig

func readConfig() {
	yamlFile, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Println("配置文件内容:", conf)
}

func createBot() *tb.Bot {
	httpTimeout := conf.Http.Timeout
	transport := &http.Transport{}

	if conf.Http.Proxy != "" {
		transport.Proxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse(conf.Http.Proxy)
		}
	}

	client := &http.Client{
		Timeout:   time.Duration(httpTimeout) * time.Second,
		Transport: transport,
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  conf.Bot.Token,
		Poller: &tb.LongPoller{Timeout: time.Duration(conf.Bot.Poller) * time.Second},
		Client: client,
	})

	if err != nil {
		log.Fatal(err)
		return nil
	}
	return bot
}

func initDB() *nodb.DB {

	cfg := new(config.Config)
	cfg.DataDir = "./db"
	dbs, err := nodb.Open(cfg)
	if err != nil {
		fmt.Printf("nodb: error opening db: %v", err)
	}

	db, _ := dbs.Select(0)
	// dbDir := "./db"
	// db, err := db.OpenDB(dbDir)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// db.Create(TableLottery)
	// db.Create(TableLuckyUser)
	// db.Create(TablePartner)
	return db
}

func main() {
	readConfig()
	ldb = initDB()
	bot = createBot()

	// bot.Handle(tb.OnText, ListenCreated)
	bot.Handle("/start", Help)
	bot.Handle("/create", CreateLottery)
	bot.Handle("/lotteries", MyLottery)
	bot.Handle("/show", ShowLottery)
	bot.Handle("/stop", StopLottery)
	bot.Handle("/delete", DeleteLottery)
	bot.Handle("/join", JoinLottery)
	bot.Handle("/roll", RollLottery)
	bot.Handle("/partner", ShowPartner)
	bot.Handle("/lucky_list", LuckyList)
	bot.Handle("/help", Help)

	log.Println("抽奖 bot 启动成功!")

	bot.Start()
}
