package main

import (
	"fmt"
	"os"

	"cloud.google.com/go/translate"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
)

func main() {
	Token := os.Getenv("DISCORD_BOT_TOKEN")

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageHandler)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	select {}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 自分のメッセージは無視
	if m.Author.ID == s.State.User.ID {
		return
	}

	jcid := os.Getenv("JAPANESE_CHANNEL_ID")
	ecid := os.Getenv("ENGLISH_CHANNEL_ID")
	var l, sl language.Tag
	var destinationChannelID string

	if m.ChannelID == jcid {
		destinationChannelID = ecid
		sl = language.Japanese
		l = language.English
	} else if m.ChannelID == ecid {
		destinationChannelID = jcid
		sl = language.English
		l = language.Japanese
	} else {
		// 対象チャンネル以外はスルー
		fmt.Println("OUT")
		return
	}

	// 翻訳処理
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		fmt.Println("Translation error:", err.Error())
		return
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{m.Content}, l, &translate.Options{
		Source: sl,
		Format: translate.Text,
	})
	if err != nil {
		fmt.Println("Translation error:", err.Error())
		return
	}

	translatedText := resp[0].Text
	// メッセージURLを生成
	messageURL := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", m.GuildID, m.ChannelID, m.ID)

	message := messageURL + "\n" + "> " + m.Content + "\n" + translatedText
	_, err = s.ChannelMessageSend(destinationChannelID, message)
	if err != nil {
		fmt.Println("Error sending message to Discord:", err.Error())
	}
}
