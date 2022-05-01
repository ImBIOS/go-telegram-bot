package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Handler(rw http.ResponseWriter, req *http.Request) {

	// Create our web hook request body type instance
	body := &webHookReqBody{}

	// Decodes the incoming request into our cutom webhookreqbody type
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		log.Printf("An error occured (Handler)")
		log.Panic(err)
		return
	}

	// If the command /joke is recieved call the sendReply function
	if strings.ToLower(body.Message.Text) == "/joke" {
		err := sendReply(body.Message.Chat.ID)
		if err != nil {
			log.Panic(err)
			return
		}
	}
}

type webHookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func sendReply(chatID int64) error {
	fmt.Println("sendReply called")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("botToken")

	// calls the joke fetcher fucntion and gets a random joke from the API
	text, err := jokeFetcher()
	if err != nil {
		return err
	}

	//Creates an instance of our custom sendMessageReqBody Type
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   text,
	}

	// Convert our custom type into json format
	reqBytes, err := json.Marshal(reqBody)

	if err != nil {
		return err
	}

	// Make a request to send our message using the POST method to the telegram bot API
	resp, err := http.Post(
		"https://api.telegram.org/bot"+token+"/"+"sendMessage",
		"application/json",
		bytes.NewBuffer(reqBytes),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + resp.Status)
	}

	return err
}

func jokeFetcher() (string, error) {
	resp, err := http.Get("http://api.icndb.com/jokes/random")
	c := &joke{}
	if err != nil {
		return "", err
	}
	err = json.NewDecoder(resp.Body).Decode(c)
	return c.Value.Joke, err
}

type joke struct {
	Value struct {
		Joke string `json:"joke"`
	} `json:"value"`
}

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}
