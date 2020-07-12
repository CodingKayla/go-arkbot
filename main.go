package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gorcon/rcon"
)

type botClient struct {
	console *rcon.Conn
}

type user struct {
	name string
	id   string
}

type userList struct {
	users []user
}

func (u userList) sendToChat(b *botClient) error {
	userNames := []string{}
	for _, onlinePlayer := range u.users {
		userNames = append(userNames, onlinePlayer.name)
	}

	message := fmt.Sprintf("%d Connected Player(s): %s", len(u.users), strings.Join(userNames, ","))

	_, err := b.sendChat(message)
	if err != nil {
		return err
	}

	return nil

}

func (b botClient) listPlayers() (userList, error) {
	players := userList{}

	response, err := b.console.Execute("listplayers")
	if err != nil {
		return players, err
	}

	cleanResponse := strings.Replace(response, "\n", "@#", -1)
	cleanResponse = strings.Replace(cleanResponse, " @# ", "", -1)
	cleanResponse = strings.Replace(cleanResponse, "@#", "", 1)
	splitResponse := strings.Split(cleanResponse, "@#")

	for _, player := range splitResponse {
		userRegex := regexp.MustCompile(`(?m)(^[0-9]{1,2}\.) ([\S]+), ([0-9]+)`)
		userParts := userRegex.FindAllStringSubmatch(player, -1)

		foundUser := user{
			name: userParts[0][2],
			id:   userParts[0][3],
		}

		players.users = append(players.users, foundUser)
	}

	return players, nil
}

func (b botClient) sendChat(message string) (string, error) {
	response, err := b.console.Execute(fmt.Sprintf("ServerChat %s", message))
	if err != nil {
		return response, err
	}

	return response, nil
}

func (b botClient) getChatBuffer() (string, error) {
	response, err := b.console.Execute("GetChat")
	if err != nil {
		return response, err
	}

	return response, nil
}

func newBot(host string, password string) (*botClient, error) {
	var client botClient

	conn, err := rcon.Dial(host, password)
	if err != nil {
		return &client, err
	}

	client.console = conn
	return &client, nil
}

func (b *botClient) close() {
	b.console.Close()
}

func main() {

	host := os.Getenv("ARKHOST")
	pass := os.Getenv("ARKPASS")

	if host == "" || pass == "" {
		log.Fatal("ARKHOST and ARKPASS env variables must be set.")
	}

	bot, err := newBot(host, pass)
	if err != nil {
		log.Fatal(err)
	}

	// players, err := bot.listPlayers()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = players.sendToChat(bot)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	chat, err := bot.getChatBuffer()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(chat)

	defer bot.close()

}
