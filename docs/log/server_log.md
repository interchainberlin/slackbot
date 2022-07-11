# Instruction on logs service of slackbot/pooltoy

Pooltoy's log is forwarded to `#intrachain-logs` chanel in slack.

There are four services running:
- [pooltoy](https://github.com/interchainio/pooltoy)
- [slackbot](https://github.com/interchainberlin/slackbot)
- slacktail
- pooltoylog
  
  cli for pooltoylog: `journalctl /root/slackbot/slackbot -f  | nc -u 127.0.0.1 9999`

pooltoy runs the pooltoy chain, Journald tracks pooltoy's log. pooltoylog fetches these logs and sends them to slacktail server. slacktail further sends them to `#intrachain-logs` chanel in slack.



