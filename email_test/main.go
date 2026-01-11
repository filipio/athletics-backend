package main

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func main() {
	ctx := context.TODO()
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	params := &resend.SendEmailRequest{
		From:    "contact@info.lekkoatletawka.eu",
		To:      []string{"filip.juza.2000@gmail.com"},
		Subject: "hello world",
		Html:    "<p>it works!</p>",
		ReplyTo: "filip.juza.2000@gmail.com",
	}

	sent, err := client.Emails.SendWithContext(ctx, params)

	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)
}
