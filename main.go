package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"./commands"
	"./commands/admin"
	"./support"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

// Running is the boolean that tells if the server is running or not
var Running bool
var Stopped *bool

var FactoCord_version string

// Pipe is an WriteCloser interface
var Pipe io.WriteCloser

// Session is a discordgo session
var Session *discordgo.Session

func main() {
	support.Config.LoadEnv()
	Running = false
	admin.R = &Running
	admin.Stopped = false
	Stopped = &admin.Stopped

	FactoCord_version, _ = getFactoCordVersion()
	commands.Version = FactoCord_version

	// Do not exit the app on this error.
	if err := os.Remove("factorio.log"); err != nil {
		fmt.Println("Factorio.log doesn't exist, continuing anyway")
	}

	logging, err := os.OpenFile("factorio.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to open factorio.log\nDetails: %s", time.Now(), err))
	}

	mwriter := io.MultiWriter(logging, os.Stdout)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for {
			// If the process is already running DO NOT RUN IT AGAIN
			// Also do NOT start server if the stop command was used.
			if !Running && !*Stopped {
				Running = true
				cmd := exec.Command(support.Config.Executable, support.Config.LaunchParameters...)
				cmd.Stderr = os.Stderr
				cmd.Stdout = mwriter
				Pipe, err = cmd.StdinPipe()
				if err != nil {
					support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to execute cmd.StdinPipe()\nDetails: %s", time.Now(), err))
				}

				err := cmd.Start()

				if err != nil {
					support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to start the server\nDetails: %s", time.Now(), err))
				}
				Session.ChannelMessageSend(support.Config.FactorioChannelID, "The factorio server was started!")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		Console := bufio.NewReader(os.Stdin)
		for {
			line, _, err := Console.ReadLine()
			if err != nil {
				support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to read the input to pass as input to the console\nDetails: %s", time.Now(), err))
			}
			_, err = io.WriteString(Pipe, fmt.Sprintf("%s\n", line))
			if err != nil {
				support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to pass input to the console\nDetails: %s", time.Now(), err))
			}
		}
	}()

	go func() {
		// Wait 10 seconds on start up before continuing
		time.Sleep(10 * time.Second)

		for {
			support.CacheDiscordMembers(Session)
			//sleep for 4 hours (caches every 4 hours)
			time.Sleep(4 * time.Hour)
		}
	}()
	discord()
}

func discord() {
	// No hard coding the token }:<
	discordToken := support.Config.DiscordToken
	commands.RegisterCommands()
	admin.P = &Pipe
	fmt.Println("Starting bot..")
	bot, err := discordgo.New("Bot " + discordToken)
	Session = bot
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to create the Discord session\nDetails: %s", time.Now(), err))
		return
	}

	err = bot.Open()

	if err != nil {
		fmt.Println("error opening connection,", err)
		support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to connect to Discord\nDetails: %s", time.Now(), err))
		return
	}

	bot.AddHandler(messageCreate)
	bot.AddHandlerOnce(support.Chat)
	time.Sleep(3 * time.Second)
	bot.UpdateStatus(0, "FactoCord V"+FactoCord_version)
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, support.Config.Prefix) {
		//command := strings.Split(m.Content[1:len(m.Content)], " ")
		//name := strings.ToLower(command[0])
		input := strings.Replace(m.Content, support.Config.Prefix, "", -1)
		log.Print("[" + m.Author.Username + "] Command: " + input)
		commands.RunCommand(input, s, m)
		return
	}

	if m.ChannelID == support.Config.FactorioChannelID {
		log.Print("[" + m.Author.Username + "] Says: " + m.Content)
		// Pipes normal chat allowing it to be seen ingame
		_, err := io.WriteString(Pipe, fmt.Sprintf("[Discord] <%s>: %s\r\n", m.Author.Username, m.ContentWithMentionsReplaced()))
		if err != nil {
			support.ErrorLog(fmt.Errorf("%s: An error occurred when attempting to pass Discord chat to in-game\nDetails: %s", time.Now(), err))
		}
		return
	}
}
