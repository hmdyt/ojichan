package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hmdyt/ojichan/ojichan"
)

var (
	TOKEN           string
	CHANNEL_ID_LIST []string
	MYSELF_URL      string

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

	if channel_ids_str := os.Getenv("OJICHAN_DISCORD_CHANNEL_ID"); channel_ids_str == "" {
		log.Fatal("環境変数 OJICHAN_DISCORD_CHANNEL_ID が設定されていません")
	} else if channel_ids := strings.Split(channel_ids_str, ","); len(channel_ids) == 0 {
		log.Fatal("環境変数 OJICHAN_DISCORD_CHANNEL_ID が不正です")
	} else {
		CHANNEL_ID_LIST = channel_ids
	}

	MYSELF_URL = os.Getenv("OJICHAN_MYSELF_URL")
}

func main() {
	go startHttpServer()
	go awake()

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

	if !isValidChannel(m.ChannelID) {
		return
	}

	if isIncludeReactionWord(m.Content) {
		log.Printf("おじさんがメッセージに気づきました: %s", m.Content)
		ojiMessage := ojichan.Say(getName(m.Author))
		s.ChannelMessageSend(m.ChannelID, ojiMessage)
	}
}

func isValidChannel(channelID string) bool {
	for _, id := range CHANNEL_ID_LIST {
		if channelID == id {
			return true
		}
	}
	return false
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

func startHttpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.ListenAndServe(":8080", nil)
}

func awake() {
	for {
		http.Get(MYSELF_URL)
		time.Sleep(1 * time.Minute)
	}
}
