package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

//Tokenは環境変数
var (
	Token             = "TOKEN"
	BotName           = "CLIENT_ID"
	stopBot           = make(chan bool)
	vcsession         *discordgo.VoiceConnection
	HelloGo           = "!hellogo"
	ChannelVoiceJoin  = "!vcjoin"
	ChannelVoiceLeave = "!vcleave"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	discord, err := discordgo.New()
	discord.Token = Token
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}

	discord.AddHandler(onMessageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening...")
	<-stopBot
	return
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("Error getting channel: ", err)
		return
	}

	voiceChannelID := os.Getenv("VOICE_CHANNEL_ID")

	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	switch {
	case strings.HasPrefix(m.Content, "!hellogo"):
		sendMessage(s, c, "Hello Go!")

	case strings.HasPrefix(m.Content, "!vcjoin"):
		guildChannels, _ := s.GuildChannels(c.GuildID)
		var sendText string
		for _, a := range guildChannels {
			sendText += fmt.Sprintf("%vチャンネルの%v(IDは%v)\n", a.Type, a.Name, a.ID)
		}
		sendMessage(s, c, sendText) //チャンネル名、ID、タイプをBOTが呟く

		vcsession, _ = s.ChannelVoiceJoin(c.GuildID, voiceChannelID, false, false)
		vcsession.AddHandler(onVoiceReceived)
	case strings.HasPrefix(m.Content, "!vcleave"):
		vcsession.Disconnect()
	}
}

func onVoiceReceived(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
	log.Print("ʕ◔ϖ◔ʔ")
}

func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	_, err := s.ChannelMessageSend(c.ID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}
