<p align="center">FactoCord 2.0 - a Continuation of the <a href="https://github.com/FactoKit/FactoCord">Factocord</a> - Factorio to Discord bridge bot for Linux</p>
<p align="center">
<a href="https://goreportcard.com/report/github.com/thecmdradama/FactoCord-2.0"><img src="https://goreportcard.com/badge/github.com/thecmdradama/FactoCord-2.0" alt="Go Report Card"></a>
</p>

# Running
*Make sure you have your .env file in the same directory as the executable/binary, you can use .envexample the template*

For those unfamiliar with linux, files starting with "." are hidden. 

There are two ways of starting FactoCord

1. Using the start.sh bash script (bash start.sh or ./start.sh) (make sure you chmod +x the script first)
2. Manually running the binary (./FactoCord) Note: File must be in CamelCase... 

The script leverages screen.  You can attach to the actual Factorio instance leveraging that tool.

# Installing as a service

To install FactoCord as a service so that it can run on startup, you can use the provided service.sh

*Note you must run service.sh as root/sudo to install it as a service*

Example of running service.sh:
`./service.sh factorio /home/facotrio/factocord/`


# Compiling

`Requires go 1.8 or above`

FactoCord uses the following packages:

- [DiscordGo](https://github.com/bwmarrin/discordgo)
- [godotenv](https://github.com/joho/godotenv/)
- [tails](https://github.com/hpcloud/tail)

You will need to add these lib as go get (most likely as root or sudo):

- `go get github.com/bwmarrin/discordgo`
- `go get github.com/joho/godotenv`
- `go get github.com/hpcloud/tail/...`

To compile:
ensure the local copy is a git initialized.  For Linux I do :
`git clone https://github.com/suprnova74/FactoCord-2.0 FactoCord`
`git init FactoCord`

and then run `go build`

to grab the appropriate version information, the local repository needs to be git initialized or build will fail.


# Error reporting

When FactoCord encounters an error will log to error.log within the same directory as itself.

If you are having an issue make sure to check the error.log to see what the problem is.

If you are unable to solve the issue yourself, please post an issue containing the error.log and I will review and attempt to solve what the problem is.

# Update-Log

Added rudimentary translator function.  It will take any translation from Smoogle bot sent to the server's channel and send it to in-game.   In-game requests for translations are currently not possible via this method as Smoogle ignores any bot output.
Merged many features from credomane / FactoCord fork


# Screenshots

<p><img src="https://i.imgur.com/nrPMtBK.png" alt="list of commands"></p>
<p><img src="http://i.imgur.com/dztOTrk.png" alt="in-game chat being sent to discord, notice how you can mention discord members"></p>
<p><img src="http://i.imgur.com/Npl0vBb.png" alt="discord chat being sent to in-game"></p>


Special thanks again to FMCore for creating the initial FactoCord project. https://github.com/FactoKit/FactoCord


