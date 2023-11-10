package main

import (
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
	"os"
	"bufio"
	"time"
)

var globalQuery string

func Handler(bot *openwechat.Bot, msg *openwechat.Message) {
    // 发送自动回复消息
    if _, err := msg.ReplyText("我只是个反馈机器人，无法正面回复你的消息，有需要回复关注公众号LT莱特电竞"); err != nil {
        log.Printf("回复消息失败：%v\n", err)
    }

    if strings.Contains(globalQuery, "特定条件") {
        // 获取当前用户
        self, err := bot.GetCurrentUser()
        if err != nil {
            log.Printf("获取当前用户失败：%v\n", err)
            return
        }

        // 获取所有好友
        friends, err := self.Friends()
        if err != nil {
            log.Printf("获取好友列表失败：%v\n", err)
            return
        }

        for _, friend := range friends {
            if friend.RemarkName == "国服车队联系员-梦乐" {
                // 向找到的好友发送消息
                if _, err := friend.SendText("来自机器人的特定消息"); err != nil {
                    log.Printf("发送消息失败：%v\n", err)
                }
                return // 找到后发送消息并退出函数
            }
        }
        log.Println("找不到指定联系人")
    }
}

func MonitorOutFile(filePath string) {
	for {
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("无法打开文件: %v", err)
			time.Sleep(1 * time.Minute) // 等待后重试
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// 假设.query值在文件中明确标记
			if strings.HasPrefix(line, "query:") {
				globalQuery = strings.TrimPrefix(line, "query:")
				break
			}
		}

		file.Close()
		time.Sleep(1 * time.Minute) // 每分钟检查一次
	}
}

func Run() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数，这里使用了匿名函数来适配Handler函数的签名
	bot.MessageHandler = func(msg *openwechat.Message) {
		Handler(bot, msg) // 修改这里，传入bot实例
	}

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	// 执行热登录，这里传入了正确的参数，假设是true，您需要根据实际情况调整
	err := bot.HotLogin(reloadStorage, true)
	if err != nil {
		log.Println("热登录失败，尝试普通登录")
		if err = bot.Login(); err != nil {
			log.Printf("登录失败: %v\n", err)
			return
		}
	}

	// 启动文件监控协程，传入正确的文件路径
	go MonitorOutFile("../../nohup.out") // 根据实际文件路径进行修改

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}

func main() {
	Run()
}
