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
	params := resend.SendEmailRequest{
		To: []string{"fooba@gmail.com"},
		Template: &resend.EmailTemplate{
			Id: "d27326b1-3311-4ace-baa2-47fae6e40f5a",
			Variables: map[string]any{
				"token": "foo",
			},
		},
	}

	sent, err := client.Emails.SendWithContext(ctx, &params)

	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)
}
