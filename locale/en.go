package locale

var En = Lang{
	InviteMessage: InviteMessage{
		Title:                    "Tournament **%v**",
		Description:              "An invitation to the tournament, including all the necessary information.\n\n*This message was automatically generated. There is no need to reply. If you have any questions or need assistance, please contact one of the organizers' assistants.*",
		MessageHeader:            "Your opponent's data",
		Nickname:                 "**Nickname**",
		GameID:                   "**Game ID**",
		Discord:                  "**Discord**",
		CheckIn:                  "**Link to check-in**",
		Warning:                  "You have %v min. to check-in before you are automatically disqualified.",
		SettingsHeader:           "Settings according to the rules",
		StandardFormat:           "**Format**",
		FinalsFormat:             "**Finals format**",
		FormatDescription:        "** (First to %v win)**",
		FT:                       "FT%v",
		Stage:                    "**Stage**",
		AnyStage:                 "Randomly selected **ALWAYS** if the opponent did not continue the set by pressing the \"Rematch\" button",
		Rounds:                   "**Rounds in 1 match**",
		Duration:                 "**Time in 1 round**",
		DurationCount:            "%v seconds",
		Crossplatform:            "**Crossplatform game**",
		CrossplatformStatusTrue:  "Enable",
		CrossplatformStatusFalse: "Disable",
	},
	StreamLobbyMessage: StreamLobbyMessage{
		Title:                    "Tournament **%v**",
		Description:              "An invitation to a live broadcast match. You need to go to the lobby below and wait for the organizer's team on the stream for further actions.\n\n*This message was automatically generated. There is no need to reply. If you have any questions or need assistance, please contact one of the organizers' assistants.*",
		MessageHeader:            "**Link to check-in**",
		Warning:                  "You have %v minutes to register before you will be automatically disqualified (meaning from the very first message received in one stage of the tournament).",
		ParamsHeader:             "**Params for searching lobby**",
		Area:                     "**Area**",
		AnyArea:                  "Any",
		CloseArea:                "Close to Me",
		Language:                 "**Language**",
		AnyLanguage:              "Any",
		SameLanguage:             "Same as Me",
		TypeConnection:           "**Type connection**",
		AnyConnection:            "Any",
		Crossplatform:            "**Crossplatform game**",
		CrossplatformStatusTrue:  "Enable",
		CrossplatformStatusFalse: "Disable",
		Passcode:                 "**Passcode**",
		PasscodeTemplate:         "```%v```",
		StreamLink:               "**Link to stream**",
	},
	ViewDataMessage: ViewDataMessage{
		Title:               "Check data",
		Description:         "A slug is made of two parts, the tournament name and the event name. The format is this:\n*tournament/<tournament_name>/event/<event_name>*",
		MessageRulesHeader:  "Rules matches",
		MessageStreamHeader: "Stream lobby data",
		LogoTournament:      "Logo tournament",
	},
	ErrorMessage: ErrorMessage{
		Input:   "Your input data isn't correct",
		Respond: "Сan't respond on message",
		NoData:  "N/D",
	},
	ResponseMessage: ResponseMessage{
		Starting:  "Starting sending...",
		InProcess: "In process...",
		Stopping:  "Stopping...",
		Stopped:   "Stopped.",
	},
}
