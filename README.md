# AwesomeBot

A telegram bot written in go, based on [telebot](https://github.com/tucnak/telebot)

[![Development chat](https://asafniv.me/bucket/telegramchat.png)](https://t.me/aw420dev)

Running
-
- Put your build configuration in `config` (Telegram bot token, Sentry DSN)
- Run `./build.sh`
- The binary will be in `out/awesomeBot`, You should not publish this binary anywhere, it contains your secret token!

Configuration
-
AwesomeBot stores configuration files in a folder named "AwesomeBot" in your OSes appropriate directory:

- macOS: `~/Library/AwesomeBot`
- Linux: `~/.config/AwesomeBot`
- Windows: Should be in your user's directory (`C:\Users\yourname\AwesomeBot`)

The files:

`Database/User` - A user database used to look up user IDs from usernames (for `/ban @username`) - it collects user info from every message it sees (you shouldn't edit it)

`helpmsg` - The help message sent by `/help`

Commands
-

`/ban` - Ban a user

`/kick` - Kick a user

`/purge` - Delete a range of messages

`/delete` - Delete a message

`/song` - Download a song from YouTube as an AAC file

`/mid` - Get a message's ID

`/cid` - Get a chat's ID

`/id` - Get a user's ID

`/help` - Send the help message (only works in PM)

PRs are welcome!
-
