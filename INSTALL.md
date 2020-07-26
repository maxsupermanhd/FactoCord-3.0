# This file contains install instructions for Debian based systems
Headless factorio server is used in this instruction (from factorio.com) in `/home/factorio`.

Tested on Ubuntu 18.04.4 LTS, Ubuntu 18.04.4 LTS (Server), WSL1 Ubuntu

# Installing pre-built binaries

## Step 0
Configuring

Enter created directory `cd FactoCord-3.0/`

Copy `config-example.json` to `config.json` (`cp config-example.json config.json`)

Open `config.json` with any editor (ex. `nano config.json`)

Then in text editor you must set:
1. Your Discord token for the bot (discord_token)
2. ID of factorio channel for chatting (factorio_channel_id)
3. Launching parameters (flags to factorio executable) (launch_parameters)
4. Executable path (executable)
5. Admin IDs (for commands) (admin_ids)
6. Mod list .json file location (including the filename) (mod_list_location)

Make sure there are no comments before ] or }, because there's a bug in json5 library

## Step 1
Running

`./FactoCord-3.0`


# Installing from sources

## Step 0
Installing deps

Make sure system is up to date `sudo apt-get update -y && sudo apt-get upgrade -y`

Download go 1.12+ (`sudo apt install golang-go git -y`) (you may need to get it from the website, repos can be outdated)

Get go packages:

- [DiscordGo](https://github.com/bwmarrin/discordgo) `go get github.com/bwmarrin/discordgo`
- [json5](https://github.com/yosuke-furukawa/json5) `go get github.com/yosuke-furukawa/json5`

## Step 1
Cloning repo

`git clone https://github.com/maxsupermanhd/FactoCord-3.0.git`

## Step 2
Configuring

Enter created directory `cd FactoCord-3.0/`

Copy `config-example.json` to `config.json` (`cp config-example.json config.json`)

Open `config.json` with any editor (ex. `nano config.json`)

Then in text editor you must set:
1. Your Discord token for the bot (discord_token)
2. ID of factorio channel for chatting (factorio_channel_id)
3. Launching parameters (flags to factorio executable) (launch_parameters)
4. Executable path (executable)
5. Admin IDs (for commands) (admin_ids)
6. Mod list .json file location (including the filename) (mod_list_location)

Make sure there are no comments before ] or }, because there's a bug in json5 library

# Step 3
Building

`go build`

## Step 4
Running

`./FactoCord-3.0`

# Using scenario support
... will eventually disable achievements, but you will have nice and clear chat in Discord.
To have afk kicked people showed and be able to customize/potentially modify messages 
please use control.lua **addition** from repo root. 
If you don't want to be control.lua modified so hard, you can place it near and use 
a `require` to get discord sending function (wrapper) and have full functionality of FactoCord.
