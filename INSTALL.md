# This file contains install instructions for Ubuntu 18.04
In instruction used headless factorio copy (from website) in /home/factorio

# Step 0
Installing deps

`sudo apt-get update -y && sudo apt-get upgrade -y`
`sudo apt install golang-go git -y`
`go get github.com/bwmarrin/discordgo`
`go get github.com/joho/godotenv`
`go get github.com/hpcloud/tail/...`

# Step 1
Cloning repo

`git clone https://github.com/maxsupermanhd/FactoCord-3.0.git`

# Step 2
Configuring

`cd FactoCord-3.0/`
`nano .envexample`

Then in text editor you must set:
1. Your Discord token for the bot (DiscordToken)
2. Id of factorio channel for chatting (FactorioChannelID)
3. Launching parameters (flags to factorio executable) (LaunchParameters)
4. Executable path (Executable)
5. Admin IDs (for commands) (AdminIDs)
6. Mod list .json file location (include filename) (ModListLocation)

Then rename `.envexample` to `.env` (`mv .envexample .env`)

# Step 3
Building

`go build`

# Step 4
Running

`./FactoCord`

