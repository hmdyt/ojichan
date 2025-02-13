package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/hmdyt/ojichan/ojichan"
)

var (
	TOKEN      string
	CHANNEL_ID string

	REACTION_WORDS = []string{
		"おじ",
		"oji",
		"オジ",
		"おっさん",
	}
)

func init() {
	TOKEN = os.Getenv("OJICHAN_DISCORD_TOKEN")
	if TOKEN == "" {
		log.Fatal("環墿変数 OJICHAN_DISCORD_TOKEN が設定されていません")
	}

	CHANNEL_ID = os.Getenv("OJICHAN_DISCORD_CHANNEL_ID")
	if CHANNEL_ID == "" {
		log.Fatal("環境変数 OJICHAN_DISCORD_CHANNEL_ID が設定されていません")
	}
}

func main() {
	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatalf("セッションの作成に失敗しました: %v", err)
	} else {
		defer dg.Close()
	}

	dg.AddHandler(handler)
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	err = dg.Open()
	if err != nil {
		log.Fatalf("セッションの開始に失敗しました: %v", err)
	} else {
		log.Println("おじさん、起床👨")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID != CHANNEL_ID {
		return
	}

	if isIncludeReactionWord(m.Content) {
		log.Printf("おじさんがメッセージに気づきました: %s", m.Content)
		ojiMessage := ojichan.Say(getName(m.Author))
		s.ChannelMessageSend(m.ChannelID, ojiMessage)
	}
}

func isIncludeReactionWord(content string) bool {
	for _, word := range REACTION_WORDS {
		if strings.Contains(content, word) {
			return true
		}
	}
	return false
}

func getName(u *discordgo.User) string {
	switch {
	case u == nil:
		return "名無し"
	case u.GlobalName != "":
		return u.GlobalName
	case u.Username != "":
		return u.Username
	default:
		return "名無し"
	}
}
