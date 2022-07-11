# slackbot
slackbot for [pooltoy](https://github.com/interchainio/pooltoy) "translates" certain types of requests from slack channels into a pooltoy chain querys/transactions.

## Instruction on running slackbot

Make a file `.env` , add the following:
```
export VERIFICATION_TOKEN=<YOUR_VERIFICATION_TOKEN>
export PORT=8080
```
then `source .env`

and start with `go run main.go`.

## Instruction on using pooltoy in slack channel
`/brrr @slack_username [emoji]`: mint emoji

`/send @slack_uername [emoji]`: send emoji to a user

`/balance @slack_username`: query balance

`/til-brrr @slack_username`: query the next earliest allowed `/brrr` time.

