# This file contains install instructions for Debian based systems
Headless factorio server is used in this instruction (from factorio.com) in `/home/factorio`.

Tested on Ubuntu 18.04.4 LTS, Ubuntu 18.04.4 LTS (Server)

# Step 0
Installing deps

Make sure system is up to date `sudo apt-get update -y && sudo apt-get upgrade -y`

Download go 1.8+ (`sudo apt install golang-go git -y`)

Get go packages:

- [DiscordGo](https://github.com/bwmarrin/discordgo) `go get github.com/bwmarrin/discordgo`
- [godotenv](https://github.com/joho/godotenv/) `go get github.com/joho/godotenv`

# Step 1
Cloning repo

`git clone https://github.com/maxsupermanhd/FactoCord-3.0.git`

# Step 2
Configuring

Enter created directory `cd FactoCord-3.0/`

rename (or copy) `.envexample` to `.env` (`cp .envexample .env`)

Open `.env` with any editor (ex. `nano .env`)

Then in text editor you must set:
1. Your Discord token for the bot (DiscordToken)
2. ID of factorio channel for chatting (FactorioChannelID)
3. Launching parameters (flags to factorio executable) (LaunchParameters)
4. Executable path (Executable)
5. Admin IDs (for commands) (AdminIDs)
6. Mod list .json file location (including the filename) (ModListLocation)


# Step 3
Building

`go build`

# Step 4
Running

`./FactoCord`

# Using scenario support
... will eventually disable achievements, but you will have nice and clear chat in Discord.
To have afk kicked people showed and be able to customize/potentially modify messages 
please use control.lua **addition** from repo root. 
If you don't want to be control.lua modified so hard, you can place it near and use 
a `require` to get discord sending function (wrapper) and have full functionality of FactoCord.
