package ojichan

import (
	"log"
	"os"
	"strconv"

	"github.com/greymd/ojichat/generator"
)

var (
	// 絵文字/顔文字の最大連続数 [default: 4].
	EMOJI_NUM int
	// 句読点挿入頻度レベル [min:0, max:3] [default: 0].
	UNCTUATION_LEVEL int
)

func init() {
	if tmp := os.Getenv("OJICHAN_EMOJI_NUM"); tmp == "" {
		log.Fatal("環境変数 OJICHAN_EMOJI_NUM が設定されていません")
	} else if num, err := strconv.Atoi(tmp); err != nil {
		log.Fatalf("環境変数 OJICHAN_EMOJI_NUM は int 型である必要があります: %v", err)
	} else {
		EMOJI_NUM = num
	}

	if tmp := os.Getenv("OJICHAN_UNCTUATION_LEVEL"); tmp == "" {
		log.Fatal("環境変数 OJICHAN_UNCTUATION_LEVEL が設定されていません")
	} else if num, err := strconv.Atoi(tmp); err != nil {
		log.Fatalf("環境変数 OJICHAN_UNCTUATION_LEVEL は int 型である必要があります: %v", err)
	} else {
		UNCTUATION_LEVEL = num
	}

}

func Say(targetName string) string {
	config := generator.Config{
		TargetName:       targetName,
		EmojiNum:         EMOJI_NUM,
		PunctuationLevel: UNCTUATION_LEVEL,
	}

	message, err := generator.Start(config)
	if err != nil {
		log.Fatalf("おじさんの言葉が出ないようです: %v", err)
	}

	return message
}
