package locale

type InviteMessage struct {
	Title                    string
	Description              string
	MessageHeader            string
	Nickname                 string
	TekkenID                 string
	Discord                  string
	CheckIn                  string
	Warning                  string
	SettingsHeader           string
	Format                   string
	FormatDescription        string
	FT                       string
	Stage                    string
	AnyStage                 string
	Rounds                   string
	Duration                 string
	DurationCount            string
	Crossplatform            string
	CrossplatformStatusTrue  string
	CrossplatformStatusFalse string
}

type StreamLobbyMessage struct {
	Title                    string
	Description              string
	MessageHeader            string
	Warning                  string
	ParamsHeader             string
	Area                     string
	AnyArea                  string
	CloseArea                string
	Language                 string
	AnyLanguage              string
	SameLanguage             string
	TypeConnection           string
	AnyConnection            string
	Crossplatform            string
	CrossplatformStatusTrue  string
	CrossplatformStatusFalse string
	Passcode                 string
	PasscodeTemplate         string
}

type lang struct {
	InviteMessage      InviteMessage
	StreamLobbyMessage StreamLobbyMessage
}
