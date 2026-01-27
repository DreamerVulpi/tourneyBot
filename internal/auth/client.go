package auth

import (
	"context"
	"encoding/json"
	"os"

	"log"
	"net/http"

	"fmt"

	"os/exec"

	"bytes"
	"io"

	"golang.org/x/oauth2"
)

const (
	addr     = "127.0.0.1:7310"
	authURL  = "https://start.gg/oauth/authorize"
	tokenURL = "https://api.start.gg/oauth/access_token"
)

func NewOauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("STARTGG_CLIENT_ID"),
		ClientSecret: os.Getenv("STARTGG_CLIENT_SECRET"),
		RedirectURL:  "http://127.0.0.1:7310/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{"user.identity"},
	}
}

func saveTokenToFile(token *oauth2.Token) error {
	bytes, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = os.WriteFile("token.json", bytes, 0600)
	if err != nil {
		return err
	}

	log.Println("saveTokenToFile: saved file to token.json")
	return nil
}

func getTokenFromFile() (*oauth2.Token, error) {
	bytes, err := os.ReadFile("token.json")
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	err = json.Unmarshal(bytes, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func NewClient(ctx context.Context) (*http.Client, error) {
	conf := NewOauth2Config()

	token, err := getTokenFromFile()
	if err != nil || !token.Valid() {
		codeChan := make(chan string)
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			authCode := r.URL.Query().Get("code")
			codeChan <- authCode
			fmt.Fprint(w, "Success! Back to programm.")
		})

		// launch server for listening
		go func() {
			if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		}()

		authURL := conf.AuthCodeURL("state")
		log.Printf("Link to auth: %v", authURL)

		// open browser
		exec.Command("rundll32", "url.dll,FileProtocolHandler", authURL).Start()
		receivedCode := <-codeChan

		token, err = conf.Exchange(ctx, receivedCode)
		if err != nil {
			return nil, err
		}

		saveTokenToFile(token)

	}

	return conf.Client(ctx, token), nil
}

func TestStartGGCall(client *http.Client) {
	// GraphQL запрос: узнать ID и имя текущего пользователя
	query := `{"query": "query { currentUser { id name } }"}`

	endpoint := "https://api.start.gg/gql/alpha"

	resp, err := client.Post(endpoint, "application/json", bytes.NewBufferString(query))
	if err != nil {
		log.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка сервера: %s\n", string(body))
		return
	}

	fmt.Printf("Ответ от start.gg: %s\n", string(body))
}
