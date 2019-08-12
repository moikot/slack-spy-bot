# Slack Spy Bot

[![Build Status](https://travis-ci.com/moikot/slack-spy-bot.svg?branch=master)](https://travis-ci.com/moikot/slack-spy-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/moikot/slack-spy-bot)](https://goreportcard.com/report/github.com/moikot/slack-spy-bot)
[![Coverage Status](https://coveralls.io/repos/github/moikot/slack-spy-bot/badge.svg?branch=master)](https://coveralls.io/github/moikot/slack-spy-bot?branch=master)

Notifies when a user goes online or offline.

## How to run

You can run it as a Docker container on a Google VM.

```shell
docker run -d -e BOT_TOKEN=[token] moikot/slack-spy-bot
```

If you've got Golang environment and Dep, you can build it from source.

```shell
git clone git@github.com:moikot/slack-spy-bot.git

cd slack-spy-bot

dep ensure -vendor-only

export GOOGLE_APPLICATION_CREDENTIALS=[credentials]; \
  export BOT_TOKEN=[token]; \
  go run .
```

# How to use

Assuming that you've already added the bot to your Slack applications, 
and it successfully connected to your Slack, you should be able to issue the 
following commands:

* To start spying on a user
  ```
  spy-on @User
  ```
* To stop spying on a user
  ```
  spy-off @User
  ```
