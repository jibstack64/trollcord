package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

const (
	CLOAK = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
)

var (
	discord *discordgo.Session
	token   string
	isBot   bool
)

func getChannelsOrGuild(mainGuildId *string, channels *[]*discordgo.Channel, restart *bool) error {
	msg := "you must provide channel id(s) / a server id."
	channelStrings := strings.Split(getInput("channel id(s) (seperate with ',') OR server id:", true, &msg), ",")
	err := loading("confirming channel(s) or guild...", func(finished *bool, err *error) {
		if len(channelStrings) == 1 {
			if guild, e := discord.Guild(channelStrings[0]); e != nil {
				if channel, e := discord.Channel(channelStrings[0]); e != nil {
					*err = e
					*finished = true
					return
				} else {
					*channels = make([]*discordgo.Channel, 1)
					*mainGuildId = channel.GuildID
					(*channels)[0] = channel
				}
			} else {
				ch, _ := discord.GuildChannels(guild.ID)
				*channels = ch
				*mainGuildId = guild.ID
			}
		} else {
			*channels = make([]*discordgo.Channel, len(channelStrings))
			for _, cs := range channelStrings {
				channel, e := discord.Channel(cs)
				if e != nil {
					*err = e
					*finished = true
					return
				} else {
					if *mainGuildId == "" {
						*mainGuildId = channel.GuildID
					} else {
						if *mainGuildId != channel.GuildID {
							*restart = true
							*err = discordgo.ErrNilState // random idk it works
							*finished = true
							return
						}
					}
					*channels = append(*channels, channel)
				}
			}
		}
		*finished = true
	})
	if *restart {
		errorPr("\nall channels must be in the same server; restarting.")
	}
	return err
}

func massSend(content string, channels []*discordgo.Channel, count int, pretty func(tracker, count int)) error {
	fmt.Print("\n")
	tracker := 1
	for i := 0; i < count; i++ {
		for c, channel := range channels {
			if channel.Type != discordgo.ChannelTypeGuildText && channel.Type != discordgo.ChannelTypeGuildNews {
				continue
			}
			_, err := discord.ChannelMessageSend(channel.ID, content)
			if err != nil {
				return err
			} else {
				if tracker > 1 {
					clearLine(1)
				}
				pretty(tracker, c)
			}
			tracker++
		}
	}

	if tracker == 0 {
		errorPr("sent no messages - missing perms or channel(s) are voice?")
	}

	return nil
}

func massPing() error {
	var mainGuildId string // main guild -> all channels are in!!!
	var channels []*discordgo.Channel
	var restart bool
	err := getChannelsOrGuild(&mainGuildId, &channels, &restart)
	if err != nil {
		return err
	}
	if restart {
		return massPing()
	}
	msg := "you must provide a count."
	countString := getInput("how many messages should be spammed (in each channel)?", true, &msg)
	count, err := strconv.Atoi(countString)
	if err != nil {
		errorPr(fmt.Sprintf("'%s' is not an integer; restarting.", countString))
		return massPing()
	} else if count <= 0 {
		errorPr("count cannot be/be below 0; restarting.")
		return massPing()
	}
	pingString := ""
	err = loading("generating mass ping...", func(finished *bool, err *error) {
		guildRoles, e := discord.GuildRoles(mainGuildId)
		*err = e
		for _, gR := range guildRoles {
			pingString += gR.Mention()
		}
		*finished = true
	})
	if err != nil {
		return err
	}

	return massSend(pingString, channels, count, func(tracker, count int) {
		if tracker%2 == 0 {
			SuccessColour.Printf("✨ sent %d messages...\n", tracker)
		} else {
			SuccessColour.Printf("🎉 sent %d messages...\n", tracker)
		}
	})
}

func webhookSpam() error {
	return nil
}

func textChannelSpam() error {
	var mainGuildId string
	var channels []*discordgo.Channel
	var restart bool
	err := getChannelsOrGuild(&mainGuildId, &channels, &restart)
	if err != nil {
		return err
	}
	if restart {
		return textChannelSpam()
	}

	msg := "you must provide a content value."
	content := getInput("message content:", true, &msg)
	if len(content) > 2000 {
		errorPr("content length cannot be above 2000; restarting.")
		return textChannelSpam()
	}

	faces := []string{
		"😤",
		"😠",
		"😖",
		"😡",
		"👿",
	}
	face := 0

	return massSend(content, channels, 10000, func(tracker, count int) {
		SuccessColour.Printf("%s spamming... (ctrl+c to stop)\n", faces[face])
		face++
		if face == len(faces) {
			face = 0
		}
	})
}

func serverDestroy() error {
	if !isBot {
		errorPr("how did we get here?")
	}
	return nil
}

func main() {
	// ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		message("\n\nexiting...\n")
		os.Exit(0)
	}()
	fmt.Println(color.GreenString(title())) // print title
	// get token
	token = getInput("enter your discord token:", true, nil)
	isBot = yesOrNo("is it a bot token?")
	if isBot {
		token = "Bot " + token
	}
	if cord, err := discordgo.New(token); err != nil {
		fatal("failed to initialise discord client.")
		return
	} else {
		discord = cord
		discord.UserAgent = CLOAK
		// loading slash
		err := loading("connecting to discord...", func(finished *bool, err *error) {
			_, e := discord.User("@me")
			*err = e
			*finished = true
		})
		if err != nil {
			fatal("\ninvalid token or failed connection.\n")
			return
		} else {
			success("\nsuccessfully connected.")
		}
	}
	// all sections
	options := []string{
		"mass pinger", "webhook spammer", "text channel spammer",
		"server destroyer",
	}
	if !isBot {
		options[len(options)-1] = ""
	}
	section := fromSelection("which tool do you wish to use?", options)
	var err error
	switch section {
	case 0:
		err = massPing()
	case 1:
		err = webhookSpam()
	case 2:
		err = textChannelSpam()
	case 3:
		err = serverDestroy()
	}
	if err != nil {
		fatal("\nan error occured: '" + err.Error() + "'.\n")
	} else {
		message("\nfinished successfully.\n")
	}
}
