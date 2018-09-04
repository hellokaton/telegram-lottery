# Telegram Lottery Bot

一个 telegram 上的抽奖 bot，给群主提供随机分配奖励的功能，主要用于筛选谁中了奖。

## 部署

代码为 Go 源码，在本地需交叉编译为其他平台二进制包，修改配置文件即可。

```yaml
http:
  proxy: 
  timeout: 20

bot:
  token: YOUR_TOKEN
  poller: 10
```

## 使用

**创建抽奖活动**

```bash
/create 10
本次奖励发放 10 台 iPhoneX，名额有限，速速参与。
```

这里的 `10` 是中奖名额数，换行后的文本是抽奖活动的描述。
确认后机器人会回复一个抽奖的唯一序号，如 `IeOybFIBcWBQBUQC`

**参与抽奖**

```bash
/join 
```

****

## TODO

1. 查看我发起的抽奖活动
2. 查看群组内的抽奖活动
3. 查看某个抽奖活动详情
