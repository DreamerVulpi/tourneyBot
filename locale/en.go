package locale

var En = lang{
	InviteMessage: InviteMessage{
		Title:                    "Tournament **%v**",
		Description:              "An invitation to the tournament, including all the necessary information.\n\n*This message was automatically generated. There is no need to reply. If you have any questions or need assistance, please contact one of the organizers' assistants.*",
		MessageHeader:            "Your opponent's data",
		Nickname:                 "**Nickname**",
		TekkenID:                 "**Tekken ID**",
		Discord:                  "**Discord**",
		CheckIn:                  "**Link to check-in**",
		Warning:                  "You have %v min. to check-in before you are automatically disqualified.",
		SettingsHeader:           "Settings according to the rules",
		Format:                   "**Format**",
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
		Warning:                  "You have %v min. to check-in before you are automatically disqualified.",
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
	},
	ViewDataMessage: ViewDataMessage{
		Title:               "Check data",
		Description:         "A slug is made of two parts, the tournament name and the event name. The format is this:\n*tournament/<tournament_name>/event/<event_name>*",
		MessageRulesHeader:  "Rules matches",
		MessageStreamHeader: "Stream lobby data",
	},
	ErrorMessage: ErrorMessage{
		Input:   "Your input data isn't correct",
		Respond: "Сan't respond on message",
	},
	ResponseMessage: ResponseMessage{
		Stopng: "Stopping...",
		Stopd:  "Stopped.",
	},
}
