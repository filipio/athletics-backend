package email

import (
	"os"
	"sync"

	"github.com/resend/resend-go/v3"
)

var client *resend.Client
var clientOnce sync.Once

func GetClient() *resend.Client {
	clientOnce.Do(func() {
		apiKey := os.Getenv("RESEND_API_KEY")
		client = resend.NewClient(apiKey)
	})
	return client
}
