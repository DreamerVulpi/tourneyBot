# TourneyBot

[![Go Report Card](https://goreportcard.com/badge/github.com/DreamerVulpi/tourneybot)](https://goreportcard.com/report/github.com/DreamerVulpi/tourneybot) [Русский раздел](#русский)

## English

<img style="padding: 10px" align="right" alt="TourneyBot logo" src="https://i.imgur.com/n9SG5IL.png" width="250">

TourneyBot is a project for tournament organizers on the [startgg](https://www.start.gg/) platform for the [Tekken 8](https://www.start.gg/game/tekken-8) game that helps solve the problem of interaction between players and organizers.

Using the open API [startgg](https://www.start.gg/) the bot receives data about a tournament in which different groups with different sets and participants participate. Then messages are sent to the participants of the tournament, which are located on the discord server of the organizer.

If you want to help the project, suggest ideas and developments in your [pull requests](https://github.com/DreamerVulpi/tourneybot/pulls).  

<br>

## Features

* Single and double elimination tournament formats are supported;
* Sending messages to all tournament participants every 5 minutes;
* Bot control and configuration of templates and variables via commands;
* Different invitations are sent depending on the type of match:

| A message with opponent's contacts  | A message with parameters to find a closed Tekken 8 lobby where the game will be played live |
| ------------- | ------------- |
|  ![example](https://i.imgur.com/tTVpPVX.png) | ![example](https://i.imgur.com/yeAIt6r.png)|
  
## Getting Started

### Installing

0. You need to get:
   * developer token for:
     * discord; [How get?](https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token)
     * startgg; [How get?](https://developer.start.gg/docs/authentication/)
   * discord application ID; [How get?](https://support-dev.discord.com/hc/en-us/articles/360028717192-Where-can-I-find-my-Application-Team-Server-ID)
   * guild ID (discord server ID); [How get?](https://support-dev.discord.com/hc/en-us/articles/360028717192-Where-can-I-find-my-Application-Team-Server-ID)
   * role ID for locales (currently supported is russian); [How get?](https://cybrancee.com/learn/knowledge-base/how-to-find-a-discord-role-id/)

    It is assumed that you have organized a tournament on the [startgg](https://www.start.gg/) platform, where participants fill out a registration form. In the form you need to provide the following information:
     * Gamer Tag;
     * Tekken ID;
     * Login to the organizer's Discord server (mandatory);
  
1. Download the finished build and create a ```config``` folder in the directory.
    * Create file ```config.toml```
    * Copy the template and fill the previously created file with it:

    ```toml
    [startgg]
    token = "your token"

    [discordbot]
    token = "Bot <your token>"
    guildID = "your id"
    appID = "your id"

    [roles]
    ru = "your id"
    ```

    * Create file ```tournament.toml```
    * Copy the template and fill the previously created file with it:

    ```toml
    [stream]
    area = "any"            # any | close
    language = "any"        # any | same
    crossplatform = true    # Enable: true | Disable: false
    connection = "any"      # any | "1"-"5"
    passCode = "1234"       # Min: "0000" | "Max: 9999"

    [rules]
    stage = "any"           # Name stage or any (check names in stages.go)
    format = 2              # FT (first N to win): 1-10
    rounds = 3              # 1-5
    duration = 60           # 30-99
    waiting = 10            # Time before disqualification in minutes: 1-any
    crossplatform = true    # Enable: true | Disable: false

    [logo]
    img = "your link to image"
    ```

2. Invite a bot to your discord server;
3. Set a role restriction on the use of commands on your discord server [How?](https://discord.com/blog/slash-commands-permissions-discord-apps-bots)

### Usage

1. Start the project;
2. Write a command ```/set-event link: <link to event your tournament>```
3. Start the tournament on [startgg](https://www.start.gg/);
4. Write a command ```/start-sending```
5. Enjoy the process!

All commands have a description and the necessary tips for their execution.

## Documentation

| Command  | Description |
| ------------- | ------------- |
| `/check`  | Check startgg, discord and bot variables |
| `/start-sending` | Start sending out invitations to tournament participants |
| `/stop-sending` | Stop sending invitations to tournament participants |
| `/set-event link:<link>` | Set an event in the bot configuration to retrieve all phaseGroups.  The event reference must include the path: `tournament/<tournament_name>/event/<event_name>` |
| `/edit-rules format:<[1-10]> stage:<name or any> rounds:<[1-5]> duration:<[30-99]> crossplatformplay:<true or false>` | Edit match rules |
| `/edit-stream-lobby area:<any or close> language:<any or same> conn:<any or [3-5]> crossplatformplay:<true or false> passcode:<[0000-9999]>` | Edit stream-lobby configurations |
| `/edit-logo-tournament url:<url>` | Edit the link to the tournament logo |

## Русский

<img style="padding: 10px" align="right" alt="TourneyBot logo" src="https://i.imgur.com/n9SG5IL.png" width="250">

TourneyBot проект для организаторов турнира платформы [startgg](https://www.start.gg/)  по игре [Tekken 8](https://www.start.gg/game/tekken-8) который помогает решить проблемы взаимодействия между игроками и организаторами.

Используя открытое API [startgg](https://www.start.gg/) бот получает данные о турнире в котором есть группы, в которых есть сеты, в которых есть участники. После отправляются сообщения участникам турнира, которые находятся на discord сервере организатора.

Если хотите помочь проекту, предлагайте идеи и разработки в ваших [пул реквестах](https://github.com/DreamerVulpi/tourneybot/pulls).  

<br>

## Особенности

* Поддержка форматов Single and double elimination;
* Отправка сообщений всем участникам турнира каждые 5 минут;
* Контролирование бота и изменение конфигурации шаблонов и переменных при помощи команд;
* В зависимости от типа матча рассылаются разные приглашения:

| Сообщение с контактами оппонента | Сообщение с параметрами для поиска закрытого лобби Tekken 8 где игра будет на стриме |
| ------------- | ------------- |
|  ![example](https://i.imgur.com/tTVpPVX.png) | ![example](https://i.imgur.com/yeAIt6r.png)|
  
## Начало работы

### Установка

0. Вам нужно получить:
   * токен разработчика для:
     * discord; [Как получить?](https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token)
     * startgg; [Как получить?](https://developer.start.gg/docs/authentication/)
   * discord application ID; [Как получить?](https://support-dev.discord.com/hc/en-us/articles/360028717192-Where-can-I-find-my-Application-Team-Server-ID)
   * guild ID (discord server ID); [Как получить?](https://support-dev.discord.com/hc/en-us/articles/360028717192-Where-can-I-find-my-Application-Team-Server-ID)
   * role ID для локалей (В настоящее время поддерживается русский); [Как получить?](https://cybrancee.com/learn/knowledge-base/how-to-find-a-discord-role-id/)

    Предполагается, что вы организовали турнир на платформе [startgg](https://www.start.gg/), участники которого заполняют регистрационную форму. В форме необходимо указать следующую информацию:
     * Gamer Tag;
     * Tekken ID;
     * Вход на дискорд сервер организатора (обязательно);
  
1. Загрузите готовую сборку и создайте в каталоге папку  ```config```.
    * Создайте файл ```config.toml```
    * Скопируйте шаблон и заполните им ранее созданный файл:

    ```toml
    [startgg]
    token = "ваш токен"

    [discordbot]
    token = "Bot <ваш токен>"
    guildID = "ваш id"
    appID = "ваш id"

    [roles]
    ru = "ваш id"
    ```

    * Создайте файл ```tournament.toml```
    * Скопируйте шаблон и заполните им ранее созданный файл:

    ```toml
    [stream]
    area = "any"            # any | close
    language = "any"        # any | same
    crossplatform = true    # Включена: true | Выключена: false
    connection = "any"      # any | "1"-"5"
    passCode = "1234"       # Мин: "0000" | "Макс: 9999"

    [rules]
    stage = "any"           # Имя локации или любое (any) (проверить имена в stages.go)
    format = 2              # FT (до N побед): 1-10
    rounds = 3              # 1-5
    duration = 60           # 30-99
    waiting = 10            # Время до дисквалификации в минутах: 1-any
    crossplatform = true    # Включена: true | Выключена: false

    [logo]
    img = "ваша ссылка на изображение"
    ```

2. Пригласите бота в ваш дискорд сервер;
3. Установите ограничение по ролям на использование команд на вашем сервере discord [Как?](https://discord.com/blog/slash-commands-permissions-discord-apps-bots)

### Использование

1. Запустите проект;
2. Напишите команду ```/установить-ивент link: <ссылка на ивент турнира>```
3. Запустите турнир на [startgg](https://www.start.gg/);
4. Напишите команду ```/начать-рассылку```
5. Наслаждайтесь процессом!

Все команды имеют описание и необходимые подсказки для их выполнения.

## Документация

| Команда  | Описание |
| ------------- | ------------- |
| `/проверка`  | Проверка startgg, discord and bot переменных |
| `/начать-рассылку` | Начните рассылать приглашения участникам турнира |
| `/остановить-рассылку` | Прекратите рассылать приглашения участникам турнира |
| `/установить-ивент link:<link>` | Установите событие в конфигурации бота для получения всех phaseGroups. Ссылка на событие должна содержать путь: `tournament/<название_турнира>/event/<название_ивента>` |
| `/редактировать-правила-матчей format:<[1-10]> stage:<name or any> rounds:<[1-5]> duration:<[30-99]> crossplatformplay:<true or false>` | Редактировать правила матчей |
| `/редактировать-стрим-лобби area:<any or close> language:<any or same> conn:<any or [3-5]> crossplatformplay:<true or false> passcode:<[0000-9999]>` | Редактировать конфигурацию лобби для стрима |
| `/редактировать-лого-турнира url:<url>` | Редактировать ссылку на логотип турнира |
