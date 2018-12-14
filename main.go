package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/doctype/steam"
	"github.com/fatih/color"
)

type Friend struct {
	SteamID uint64
	Name    string
	Tags    []string
}

type MessageTo struct {
	DestinationTags []string
	Message         string
}

type Friends struct {
	Friends  []Friend
	Messages []MessageTo
}

type Config struct {
	Username string
	Password string
}

func addFriends(apiKey string, steamID steam.SteamID, reader *bufio.Reader) {
	apiSession := steam.NewSessionWithAPIKey(apiKey)

	friends, err := apiSession.GetFriends(steamID)
	if err != nil {
		color.Red("Failed to retrieve the friend list.")
	}

	var addedFriends []Friend
	for _, friend := range friends {
		friendSteamID := strconv.FormatUint(friend.SteamID, 10)
		summaries, err := apiSession.GetPlayerSummaries(friendSteamID)
		friendSummary := summaries[0]
		if err != nil {
			color.Red("Failed to get profile for ", friend.SteamID)
		} else {
			fmt.Printf("Friend %s (%d) \n", friendSummary.PersonaName, friend.SteamID)
			color.HiCyan("Add (y/N): ")
			confirm, _ := reader.ReadString('\n')
			if strings.TrimSpace(confirm) == "y" {
				color.HiCyan("Enter tags (comma separated) (optional): ")
				tagsInput, _ := reader.ReadString('\n')
				tags := strings.Split(tagsInput, ",")
				for i, tag := range tags {
					tags[i] = strings.TrimSpace(tag)
				}
				addedFriend := Friend{SteamID: friend.SteamID, Name: friendSummary.PersonaName, Tags: tags}
				addedFriends = append(addedFriends, addedFriend)
			}
		}
	}
	var messageTo []MessageTo
	messageTo = append(messageTo, MessageTo{DestinationTags: make([]string, 1), Message: "NS2"})
	friendsData := Friends{Friends: addedFriends, Messages: messageTo}
	addedFriendsJson, _ := json.MarshalIndent(friendsData, "", "\t")
	err = ioutil.WriteFile("friends.json", addedFriendsJson, 0644)
	if err != nil {
		color.Red("Failed writing JSON file")
		fmt.Println(err)
	} else {
		color.Green("Successfully saved to friends.json")
		color.Yellow("Add messages to the corresponding JSON key and rerun this program.")
	}

	fmt.Print("\nPress any key to exit...")
	reader.ReadString('\n')
}

func sendMessages(session *steam.Session, steamID steam.SteamID, reader *bufio.Reader) {
	jsonFile, err := os.Open("friends.json")
	if err != nil {
		color.Red("Error opening friends.json")
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonContents, _ := ioutil.ReadAll(jsonFile)
	var friends Friends
	json.Unmarshal([]byte(jsonContents), &friends)

	if err = session.ChatLogin(""); err != nil {
		log.Fatal(err)
	}
	defer session.ChatLogoff()

	var sent bool
	var count int
	for _, friend := range friends.Friends {
		sent = false
		for _, friendTag := range friend.Tags {
			for _, message := range friends.Messages {
				for _, destinationTag := range message.DestinationTags {
					if (friendTag == destinationTag || destinationTag == "") && !sent {
						summaries, _ := session.GetPlayerSummaries(strconv.FormatUint(friend.SteamID, 10))
						if summaries[0].PersonaState != 0 {
							fmt.Printf("Sending \"%s\" to %s (matching tag: %s)\n", message.Message, friend.Name, destinationTag)
							session.ChatSendMessage(steam.SteamID(friend.SteamID), message.Message, steam.MessageTypeSayText)
						}
						sent = true
						count++
						break
					}
				}
			}
		}
	}

	color.Green(fmt.Sprintf("Successfully sent to %d friends", count))

	fmt.Print("\nPress any key to exit...")
	reader.ReadString('\n')
}

func main() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		color.Red("Cannot open config.json")
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonContents, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal([]byte(jsonContents), &config)
	username := config.Username
	password := config.Password

	session := steam.NewSession(&http.Client{}, "")

	fmt.Printf("Trying to logging as %s...\n", username)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter SteamGuard code: ")
	twoFactorCodeInput, _ := reader.ReadString('\n')
	twoFactorCode := strings.TrimSpace(twoFactorCodeInput)

	err = session.LoginTwoFactorCode(username, password, twoFactorCode)
	if err != nil {
		color.Red("Login failed")
		fmt.Println(err)
	} else {
		color.Green("Login successful")
	}

	apiKey, err := session.GetWebAPIKey()
	if err != nil {
		color.Red("Cannot get the API key.")
	}

	if err = session.ChatLogin(""); err != nil {
		log.Fatal(err)
	}
	defer session.ChatLogoff()

	steamID := steam.SteamID(session.GetSteamID())
	fmt.Println("")
	fmt.Println("[1] Generate friends.json")
	fmt.Println("[2] Send messages")
	fmt.Println("")

	for {
		color.HiCyan("Enter option number: ")
		chosenInput, _ := reader.ReadString('\n')
		chosen := strings.TrimSpace(chosenInput)

		if chosen == "1" {
			addFriends(apiKey, steamID, reader)
			break
		} else if chosen == "2" {
			sendMessages(session, steamID, reader)
			break
		} else {
			color.Red("Invalid option")
		}
	}

}
