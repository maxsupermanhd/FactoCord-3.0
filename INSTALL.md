# This file contains install instructions for Debian based systems
In instruction used headless factorio copy (from website) in `/home/factorio`.

Tested on Ubuntu 18.04.4 LTS, Ubuntu 18.04.4 LTS (Server)

# Step 0
Installing deps

Make sure system is up to date `sudo apt-get update -y && sudo apt-get upgrade -y`

Download go 1.8+ (`sudo apt install golang-go git -y`)

Get go packages:

- [DiscordGo](https://github.com/bwmarrin/discordgo) `go get github.com/bwmarrin/discordgo`
- [godotenv](https://github.com/joho/godotenv/) `go get github.com/joho/godotenv`
- [tails](https://github.com/hpcloud/tail) `go get github.com/hpcloud/tail/...`

# Step 1
Cloning repo

`git clone https://github.com/maxsupermanhd/FactoCord-3.0.git`

# Step 2
Configuring

Enter created directory `cd FactoCord-3.0/`

Open `.envexample` with any editor (ex. `nano .envexample`)

Then in text editor you must set:
1. Your Discord token for the bot (DiscordToken)
2. Id of factorio channel for chatting (FactorioChannelID)
3. Launching parameters (flags to factorio executable) (LaunchParameters)
4. Executable path (Executable)
5. Admin IDs (for commands) (AdminIDs)
6. Mod list .json file location (include filename) (ModListLocation)

Then rename (or copy) `.envexample` to `.env` (`cp .envexample .env`)

# Step 3
Building

`go build`

# Step 4
Running

`./FactoCord`

