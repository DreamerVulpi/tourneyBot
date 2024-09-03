package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	AuthToken             string
	GuildID               string
	AppID                 string
	Slug                  string
	templateInviteMessage = `
	Турнир: %v
	Следующий оппонент: %v
	Tekken ID: %v
	Дискорд: <@%v>
	Ссылка на check-in: %v
	
	*Это сообщение сгенерировано автоматически. Отвечать на него не нужно. В случае вопросов или помощи обращайтесь к помощникам организатора.*
	`
)

func SetAuthToken(token string) {
	AuthToken = token
}

func SetServerID(guildID string) {
	GuildID = guildID
}

func SetAppID(appID string) {
	AppID = appID
}

func SetSlug(slug string) {
	Slug = slug
}

func slug() bool {
	return len(Slug) > 0
}

func server() bool {
	return len(GuildID) > 0
}

func app() bool {
	return len(AppID) > 0
}

func Start() {
	session, err := discordgo.New(AuthToken)
	if err != nil {
		fmt.Println(err)
	}

	session, err = SetCommands(AppID, GuildID, session)
	if err != nil {
		fmt.Println(err)
	}

	session.AddHandler(func(
		s *discordgo.Session,
		m *discordgo.MessageCreate,
	) {
		sendMessage(s, m)
	})

	session.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		handlerCommands(s, i)
	})

	session.AddHandler(func(
		s *discordgo.Session,
		m *discordgo.MessageCreate,
	) {
		handlerInputs(s, m)
	})

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = session.Open()
	if err != nil {
		fmt.Println(err)
	}

	defer session.Close()

	fmt.Println("the bot is online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	session.Close()
}
