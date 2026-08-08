package api_test

import (
	"context"
	"fmt"
	"time"

	"github.com/LeonardoSilva-10/brdocs/api"
)

func ExampleClient_ConsultaCNPJ() {
	client := api.NewClient(
		api.WithSerproCredentials("consumerKey", "consumerSecret"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = client
	_ = ctx

	fmt.Println("Client created successfully")
	// Output:
	// Client created successfully
}
