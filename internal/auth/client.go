package auth

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"log"
	"net/http"

	"fmt"

	"os/exec"

	"bytes"

	"golang.org/x/oauth2"
)

const (
	addr            = "127.0.0.1:7310"
	startggAuthURL  = "https://start.gg/oauth/authorize"
	startggTokenURL = "https://api.start.gg/oauth/access_token"
	discordAuthURL  = "https://discord.com/api/oauth2/authorize"
	discordTokenURL = "https://discord.com/api/oauth2/token"
)

type Identity struct {
	ID          string
	Username    string
	Platform    string
	RawResponse map[string]interface{}
}

type AuthClient struct {
	Config     *oauth2.Config
	HTTPClient *http.Client
	TokenFile  string
}

func GetStartggOauth2() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("STARTGG_CLIENT_ID"),
		ClientSecret: os.Getenv("STARTGG_CLIENT_SECRET"),
		RedirectURL:  "http://127.0.0.1:7310/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  startggAuthURL,
			TokenURL: startggTokenURL,
		},
		Scopes: []string{"user.identity", "tournament.reporter"},
	}
}

func (_ *AuthClient) GetAccessToken(filename string) (string, error) {
	token, err := getTokenFromFile(filename)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func GetDiscordOauth2() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  "http://127.0.0.1:7310/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  discordAuthURL,
			TokenURL: discordTokenURL,
		},
		Scopes: []string{"identify"},
	}
}

func saveTokenToFile(filename string, token *oauth2.Token) error {
	bytes, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, 0600)
}

func getTokenFromFile(filename string) (*oauth2.Token, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	token := &oauth2.Token{}
	return token, json.Unmarshal(bytes, token)
}

func (ac *AuthClient) Init(ctx context.Context) error {
	token, err := getTokenFromFile(ac.TokenFile)
	if err != nil || !token.Valid() {
		codeChan := make(chan string)

		mux := http.NewServeMux()
		server := &http.Server{Addr: addr, Handler: mux}

		mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			fmt.Fprint(w, "Authurization is success! Back to programm.")
			codeChan <- code
		})

		// launch server for listening
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTP Server error: %v", err)
			}
		}()

		authURL := ac.Config.AuthCodeURL("state")
		log.Printf("Link to auth: %v", authURL)

		// open browser for os windows
		exec.Command("rundll32", "url.dll,FileProtocolHandler", authURL).Start()
		var code string
		select {
		case code = <-codeChan:
		case <-ctx.Done():
			server.Shutdown(context.Background())
			return ctx.Err()
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)

		token, err = ac.Config.Exchange(ctx, code)
		if err != nil {
			return err
		}
		saveTokenToFile(ac.TokenFile, token)
	}

	ac.HTTPClient = ac.Config.Client(ctx, token)
	return nil
}

func (ac *AuthClient) GetDiscordMe(ctx context.Context) (*Identity, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, err
	}
	resp, err := ac.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &Identity{ID: data.ID, Username: data.Username, Platform: "discord"}, nil
}

func (ac *AuthClient) GetStartGGMe(ctx context.Context) (*Identity, error) {
	query := `{"query": "query { currentUser { id name } }"}`

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.start.gg/gql/alpha", bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ac.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Data struct {
			CurrentUser struct {
				ID   json.Number `json:"id"`
				Name string      `json:"name"`
			} `json:"currentUser"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &Identity{
		ID:       response.Data.CurrentUser.ID.String(),
		Username: response.Data.CurrentUser.Name,
		Platform: "startgg",
	}, nil
}

func TestStartGGCall(ac *AuthClient) {
	ctx := context.Background()

	// Вызываем наш новый унифицированный метод
	me, err := ac.GetStartGGMe(ctx)
	if err != nil {
		log.Fatalf("Ошибка при проверке личности Start.gg: %v", err)
	}

	fmt.Printf("Успех! Мы зашли как: %s (ID: %s) на платформе %s\n",
		me.Username, me.ID, me.Platform)
}
