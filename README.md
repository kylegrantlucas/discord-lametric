# discord-lametric

[![pipeline status](https://gitlab.com/kylegrantlucas/discord-lametric/badges/master/pipeline.svg)](https://gitlab.com/kylegrantlucas/discord-lametric/commits/master)

A simple daemon that listens to a Discord server/channels and publishes the messages to a LaMetric clock.

## Usage

### Go

`$ go get -u github.com/kylegrantlucas/discord-lametric`

`$ env DISCORD_EMAIL=xxxx DISCORD_PASSWORD=xxxx DISCORD_SERVER_ID=xxxx LAMETRIC_IP=xxxx LAMETRIC_API_KEY=xxx discord-lametric`

### Docker

`$ docker run -e DISCORD_EMAIL=xxxx -e DISCORD_PASSWORD=xxxx -e DISCORD_SERVER_ID=xxxx -e LAMETRIC_IP=xxxx -e LAMETRIC_API_KEY=xxx kylegrantlucas/discord-lametric`

### Docker Compose

```yaml
discord-lametric:
    container_name: discord-lametric
    image: kylegrantlucas/discord-lametric
    environment:
      - DISCORD_EMAIL=xxxx
      - DISCORD_PASSWORD=xxxx
      - LAMETRIC_IP=xxxx
      - LAMETRIC_API_KEY=xxxx
      - LAMETRIC_ICON_ID=xxxx
      - DISCORD_SERVER_ID=xxxxx
      - DISCORD_CHANNELS=xxxx,xxxx
    restart: unless-stopped
```

## Channel Selection

To limit the channels you listen to, simply pass teh environment variable `DISCORD_CHANNEL` with a comma seperate list of channels to listen on.

Example:

`DISCORD_CHANNEL=general,offtopic`

## LaMetric Icon

To specify a custom icon for the LaMetric notifications, set `LAMETRIC_ICON_ID`:

Example:

`LAMETRIC_ICON_ID=24240`