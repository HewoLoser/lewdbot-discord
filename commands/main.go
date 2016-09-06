package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strings"
)

var eightballResponses = []string{
	"Most definitely yes",
	"For sure",
	"As I see it, yes",
	"My sources say yes",
	"Yes",
	"Most likely",
	"Perhaps",
	"Maybe",
	"Not sure",
	"It is uncertain",
	"Ask me again later",
	"Don't count on it",
	"Probably not",
	"Very doubtful",
	"Most likely no",
	"Nope",
	"No",
	"My sources say no",
	"Dont even think about it",
	"Definitely no",
	"NO - It may cause disease contraction",
}

func ParseMessage(s *discordgo.Session, m *discordgo.MessageCreate, text string) (bool, string) {

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
	}

	command := strings.ToLower(text)
	reply := ""

	if strings.HasPrefix(command, "!8ball") {
		reply = eightball(text)

		return true, reply
	} else if channel.GuildID == "111928847846367232" || channel.GuildID == "135827109485608960" {
		if strings.HasPrefix(command, "!list") {

			reply = listRoles(s, channel.GuildID)

			return true, reply
		} else if strings.HasPrefix(command, "!subscribe") {
			if len(text) > 11 {
				reply = addRole(s, channel.GuildID, m.Author.ID, text[11:])
			} else {
				reply = "What are you subscribing to?~"
			}
			return true, reply
		} else if strings.HasPrefix(command, "!unsubscribe") {
			if len(text) > 13 {
				reply = removeRole(s, channel.GuildID, m.Author.ID, text[13:])
			} else {
				reply = "What are you unsubscribing from?~"
			}
			return true, reply
		}
	}

	return false, reply
}

func eightball(text string) string {
	answer := eightballResponses[rand.Intn(len(eightballResponses)-1)]

	if len(text) > 7 {
		question := text[7:]

		return fmt.Sprintf("*%s* **%s**", question, answer)
	}

	return answer
}

func listRoles(s *discordgo.Session, GuildID string) string {
	g, err := s.Guild(GuildID)
	if err != nil {
		fmt.Println(err)
	}

	var reply string

	for _, role := range g.Roles {
		fmt.Println(role)

		if role.Name == "@everyone" {
			continue
		}

		reply += fmt.Sprintf("%s | %d\n", role.Name, role.Position)
	}

	return reply
}

func addRole(s *discordgo.Session, GuildID string, UserID string, arg string) string {
	g, err := s.Guild(GuildID)
	if err != nil {
		fmt.Println(err)
	}

	exists, role := roleExists(g, arg)

	fmt.Println(arg, exists, role)

	if !exists {
		return "I can't find such group~"
		/*
			role, err := s.GuildRoleCreate(GuildID)
			if err != nil {
				fmt.Println(err)
				return "Failed to create role"
			}

			role.Name = arg
			role.Permissions = 37080064

			role, err = s.GuildRoleEdit(GuildID, role.ID, role.Name, role.Color, role.Hoist, role.Permissions)
			if err != nil {
				fmt.Println(err)
				return " "
			}
		*/
	}

	member, err := s.GuildMember(GuildID, UserID)
	if err != nil {
		fmt.Println(err)
	}

	for _, _role := range member.Roles {
		if _role == role.ID {
			return fmt.Sprintf("You're already subscribed to %s~", arg)
		}
	}

	member.Roles = append(member.Roles, role.ID)

	err = s.GuildMemberEdit(GuildID, UserID, member.Roles)
	if err != nil {
		fmt.Println(err)
		return "I can't touch that group dude, do it yourself~"
	}

	return fmt.Sprintf("You're now subscribed to %s~", arg)
}

func removeRole(s *discordgo.Session, GuildID string, UserID string, arg string) string {
	g, err := s.Guild(GuildID)
	if err != nil {
		fmt.Println(err)
	}

	exists, role := roleExists(g, arg)

	fmt.Println(arg, exists, role)

	if !exists {
		return "I can't find such group~"
	}

	member, err := s.GuildMember(GuildID, UserID)
	if err != nil {
		fmt.Println(err)
	}

	found := false
	pos := 0

	for i, _role := range member.Roles {
		if _role == role.ID {
			found = true
			pos = i
		}
	}

	if !found {
		return fmt.Sprintf("You're already not subscribed to %s~", arg)
	}

	member.Roles = append(member.Roles[:pos], member.Roles[pos+1:]...)

	err = s.GuildMemberEdit(GuildID, UserID, member.Roles)
	if err != nil {
		fmt.Println(err)
		return "I can't touch that group dude, do it yourself~"
	}

	return fmt.Sprintf("Unsubscribed from %s~", arg)
}

func roleExists(g *discordgo.Guild, name string) (bool, *discordgo.Role) {
	for _, role := range g.Roles {
		if role.Name == "@everyone" {
			continue
		}

		if role.Name == name {
			return true, role
		}

	}

	return false, nil // &discordgo.Role{}
}