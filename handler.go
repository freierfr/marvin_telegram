package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/protobuf/proto"
	tele "gopkg.in/telebot.v3"

	"marvin/pb"
	"marvin/utils"
)

func handleMessage(c tele.Context, bot *tele.Bot, isGroupMessage bool) error {
	if !utils.IsAllowedUser(c.Sender().ID) {
		reply := fmt.Sprintf("You are not allowed to use this bot. Your user ID is %d", c.Sender().ID)
		c.Send(reply)
		return nil
	}

	c.Notify(tele.Typing)

	client := openai.NewClient(utils.GetConfig("OPENAI_API_KEY"))
	ctx := context.Background()
	receivedMessage := c.Text()

	redis_client := utils.ConnectRedis()
	redis_key := fmt.Sprintf("telegram:%d", c.Chat().ID)
	previous_messages, err := redis_client.LRange(ctx, redis_key, 0, -1).Result()
	if err != nil {
		log.Printf("redis set error: %v\n", err)
		c.Send(err.Error())
		panic(err)
	}
	redis_client.Expire(ctx, redis_key, (60*60)*time.Second)

	msg := &pb.Message{
		Type:         pb.MessageType_FROM_USER,
		Message:      receivedMessage,
		FromUserId:   c.Sender().ID,
		FromUsername: c.Sender().Username,
	}
	payload, _ := proto.Marshal(msg)
	err = redis_client.RPush(ctx, redis_key, payload).Err()
	// err := redis_client.Set(ctx, redis_key, "bar", (60*60)*time.Second).Err()
	if err != nil {
		log.Printf("redis set error: %v\n", err)
		c.Send(err.Error())
		panic(err)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleAssistant,
			Content: utils.GetConfig("BOT_PROMPT"),
		},
	}

	for _, previous_message := range previous_messages {
		tmp_message := &pb.Message{}
		if err := proto.Unmarshal([]byte(previous_message), tmp_message); err != nil {
			log.Printf("redis read error: %v\n", err)
			c.Send(err.Error())
			panic(err)
		}
		// check type of message and assign correct role
		if tmp_message.Type == pb.MessageType_FROM_USER {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: previous_message,
			})
		} else if tmp_message.Type == pb.MessageType_FROM_ASSISTANT {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: previous_message,
			})
		}

	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: receivedMessage,
	})

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    utils.GetConfig("OPENAI_MODEL"),
			Messages: messages,
		},
	)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		c.Send(err.Error())
		return err
	}

	msg.Type = pb.MessageType_FROM_ASSISTANT
	msg.Message = resp.Choices[0].Message.Content
	msg.FromUserId = 0
	msg.FromUsername = "Marvin"
	payload, _ = proto.Marshal(msg)
	err = redis_client.RPush(ctx, redis_key, payload).Err()

	c.Send(resp.Choices[0].Message.Content)

	return err
}

func main() {
	godotenv.Load()

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	pref := tele.Settings{
		Token: utils.GetConfig("TELEGRAM_TOKEN"),
		Poller: &tele.Webhook{
			Listen:      listenAddr,
			Endpoint:    &tele.WebhookEndpoint{PublicURL: utils.GetConfig("WEBHOOK_URL")},
			DropUpdates: true,
		},
		Verbose: true,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.SetCommands([]tele.Command{
		{
			Text:        "image",
			Description: "generate an image from the provided description, ex: /image a dog playing with a ball",
		},
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		// detect if it's a group message
		if c.Chat().Type != tele.ChatPrivate {
			if strings.HasPrefix(c.Text(), utils.GetConfig("GROUP_TRIGGER_KEYWORD")) {
				return handleMessage(c, b, true)
			}
		} else {
			return handleMessage(c, b, false)
		}

		return nil
	})

	b.Handle("/image", func(c tele.Context) error {
		c.Notify(tele.Typing)

		client := openai.NewClient(utils.GetConfig("OPENAI_API_KEY"))
		ctx := context.Background()

		reqUrl := openai.ImageRequest{
			Prompt:         c.Message().Payload,
			Size:           openai.CreateImageSize1024x1024,
			ResponseFormat: openai.CreateImageResponseFormatURL,
			N:              1,
		}

		respUrl, err := client.CreateImage(ctx, reqUrl)
		if err != nil {
			log.Printf("Image creation error: %v\n", err)
			c.Send(err.Error())
			panic(err)
		}
		fmt.Println(respUrl.Data[0].URL)
		p := &tele.Photo{File: tele.FromURL(respUrl.Data[0].URL)}
		err = c.Send(p)
		if err != nil {
			log.Printf("Image sending error: %v\n", err)
			c.Send(err.Error())
			panic(err)
		}
		return nil
	})

	b.Start()
}
