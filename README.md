# Slack Spy Bot

Notifies when a user goes online and offline.

## How to run

You can run it as a Docker container on amd64 or arm32/64 platforms.

```shell
docker run -d -e SPY_BOT_TOKEN=[token] moikot/slack-spy-bot
```

If you've got Golang environment and Dep, you can build it from source.

```shell
git clone git@github.com:moikot/slack-spy-bot.git
cd slack-spy-bot
dep ensure -vendor-only
export GOOGLE_APPLICATION_CREDENTIALS=[credentials]; \
  export SPY_BOT_TOKEN=[token]; \
  go run .

```

# How to use

Assuming that you've already added the bot to your Slack applications, 
and it successfully connected you Slack, you should be able to issue the 
following commands:

* To start spying on a user
  ```
  spy-on @User
  ```
* To stop spying on a user
  ```
  spy-off @User
  ```