<p align="center">FactoCord 3.0 - a Continuation of the <a href="https://github.com/thecmdradama/FactoCord-2.0">Factocord2</a> - Factorio to Discord bridge bot for Linux</p>
<p align="center">
<a href="https://goreportcard.com/report/github.com/thecmdradama/FactoCord-2.0"><img src="https://goreportcard.com/badge/github.com/thecmdradama/FactoCord-2.0" alt="Go Report Card"></a>
</p>

# Compiling

`Requires go 1.8 or above`

FactoCord uses the following packages:

- [DiscordGo](https://github.com/bwmarrin/discordgo)
- [godotenv](https://github.com/joho/godotenv/)
- [tails](https://github.com/hpcloud/tail)

You will need to add these lib as go get:

- `go get github.com/bwmarrin/discordgo`
- `go get github.com/joho/godotenv`
- `go get github.com/hpcloud/tail/...`

To compile just do `go build`

for more detailed instructions see INSTALL.md file

# Error reporting

When FactoCord3 encounters an error will log to error.log within the same directory as itself.

If you are having an issue make sure to check the error.log to see what the problem is.

If you are unable to solve the issue yourself, please post an issue containing the error.log and I will review and attempt to solve what the problem is.


# Screenshots

<p><img src="https://i.imgur.com/nrPMtBK.png" alt="list of commands"></p>
<p><img src="http://i.imgur.com/dztOTrk.png" alt="in-game chat being sent to discord, notice how you can mention discord members"></p>
<p><img src="http://i.imgur.com/Npl0vBb.png" alt="discord chat being sent to in-game"></p>


Special thanks again to thecmdradama for creating the FactoCord2 project. https://github.com/FactoKit/FactoCord-2.0
Special thanks again to FMCore for creating the initial FactoCord project. https://github.com/FactoKit/FactoCord
