# XKCD Bot

A XKCD Telegram Bot written in Go. Available at [@xkcdGO_bot](https://telegram.me/xkcdGO_bot).

## Features

```
/current        - Send the current comic
/comic [num]    - Send comic #num. If not specified, sends a random comic
/subscribe      - Subscribes to XKCD, receiving new comics when available.
/unsubscribe    - To stop receiving new comics.
```

## Installation

`git clone https://github.com/iurimateus/xkcdbot.git`  
`cd xkcdbot`  
`go run *.go` or `go build -o xkcd-bot *.go`

## Usage

An environment variable named `TELEGRAM_TOKEN` must be set (see
[@BotFather](https://telegram.me/BotFather)).  

Intended to run as a service. For systemd, see this example for a unit
file:

```
[Service]
Environment="TELEGRAM_TOKEN=000000000:xxxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxxx"
WorkingDirectory=<directory-to-binary>

Type=simple

ExecStart=<complete-path-to-binary>
Restart=always
RestartSec=5
TimeoutStopSec=5
SyslogIdentifier=xkcd-bot

[Install]
WantedBy=multi-user.target
```

### Misc

Currently, comics and subscribed chats are stored in json files and its contents reside
in-memory when the bot is running. Comic file is saved to the disk periodically
according to changes and users file is saved upon users request to subscribe/unsubscribe.
