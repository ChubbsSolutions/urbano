package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ChubbsSolutions/urbano/objects"
	"github.com/codegangsta/cli"
	"github.com/mailgun/mailgun-go"
	"github.com/ttacon/chalk"
)

//Define constants
const version string = "0.3"
const author string = "Chubbs Solutions"
const email string = "urbano@chubbs.solutions"
const appName string = "urbano"
const appDescription string = "Get a fresh Urban  Dictionary word in your inbox"

//MailgunPublicAPIKey key for the mail service.
var MailgunPublicAPIKey = os.Getenv("MAILGUN_PUBLIC_API_KEY")

//MailgunDomain key for the mail service.
var MailgunDomain = os.Getenv("MAILGUN_DOMAIN")

//Evaluate options on main
func main() {
	t := time.Now()
	app := cli.NewApp()
	app.Name = appName
	app.Usage = appDescription
	app.Email = email
	app.Author = author
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name:      "send",
			ShortName: "s",
			Usage:     "Get and Send a new word by email",
			Action: func(c *cli.Context) {

				word, err := getNewWord()
				if err != nil {
					fmt.Println(chalk.Red, "Error getting the word")
					return
				}

				err = displayWord(word)
				if err != nil {
					fmt.Println(chalk.Red, err)
					return
				}

				if len(c.Args()) == 1 {
					recipient := c.Args()[0]
					subject := fmt.Sprintf("Urban Dictionary Word of the day for %s", t.Format("Jan 02, 2006"))
					err = emailWord(word, subject, recipient)
					if err != nil {
						fmt.Println(chalk.Red, err)
						return
					}

				}
			},
		},
		{
			Name:      "display",
			ShortName: "d",
			Usage:     "Display a new word",
			Action: func(c *cli.Context) {
				word, err := getNewWord()
				if err != nil {
					fmt.Println(chalk.Red, err)
					return
				}

				err = displayWord(word)
				if err != nil {
					fmt.Println(chalk.Red, err)
					return
				}
			},
		},
		{
			Name:      "word",
			ShortName: "w",
			Usage:     "Define a word",
			Action: func(c *cli.Context) {
				if len(c.Args()) != 1 {
					fmt.Println("Usage: urbano word|w word")
					return
				}
				word, err := getWordDefinition(c.Args()[0])
				if fmt.Sprintf("%s", err) == "NOTFOUND" {
					fmt.Printf("Sorry, %s is not on Urban Dictionary.\n", c.Args()[0])
					return
				}
				if err != nil {
					fmt.Println(chalk.Red, err)
					return
				}

				err = displayWord(word)
				if err != nil {
					fmt.Println(chalk.Red, err)
					return
				}
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(chalk.Red, err)
		return

	}

}

//getNewWord gets a random UD word
func getNewWord() (objects.WordData, error) {
	var UDURL = "http://api.urbandictionary.com/v0/random"
	wd := objects.WordDataSlice{}
	var word objects.WordData
	var good = false
	tu := 30000

	for good == false {
		resp, err := http.Get(UDURL)
		if err != nil {
			return word, err
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		data, _ := ioutil.ReadAll(resp.Body)
		if err != nil {
			return word, err
		}

		err = json.Unmarshal([]byte(string(data)), &wd)
		if err != nil {
			return word, err
		}

		for _, element := range wd.List {
			if element.ThumbsUp > tu {
				word = element
				good = true
			}
		}
	}
	return word, nil
}

//getWordDefinition gets a random UD word
func getWordDefinition(wordToDefine string) (objects.WordData, error) {
	var UDURL = "http://api.urbandictionary.com/v0/define?term=" + strings.Replace(wordToDefine, " ", "", -1)
	wd := objects.WordDataSlice{}
	var word objects.WordData

	resp, err := http.Get(UDURL)
	if err != nil {
		return word, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	data, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return word, err
	}

	err = json.Unmarshal([]byte(string(data)), &wd)
	if err != nil {
		return word, err
	}

	for _, element := range wd.List {
		if element.ThumbsUp > word.ThumbsUp {
			word = element
		}
	}
	if word.Definition == "" {
		return word, errors.New("NOTFOUND")
	}
	return word, nil
}

//displayWord displaying word.
func displayWord(word objects.WordData) error {

	fmt.Print(chalk.Cyan, "Word of the day: ", word.Word, "   ++Thumbs Up: ", word.ThumbsUp, "   --Thumbs Down: ", word.ThumbsDown, "\n\n")
	fmt.Print(chalk.Green, word.Definition, "\n\n")
	fmt.Print(chalk.Blue, "Example: ", word.Example, "\n\n")
	fmt.Print(chalk.Yellow, "Courtesy of ", word.Author, "\n\n\n")
	fmt.Print(chalk.Black, "Brought to you by Chubbs Solutions.")
	return nil
}

func emailWord(word objects.WordData, subject, recipient string) error {

	if MailgunPublicAPIKey == "" || MailgunDomain == "" {
		return errors.New("Please set the MAILGUN_PUBLIC_API_KEY and MAILGUN_DOMAIN variables.")
	}

	sender := fmt.Sprintf("donotreply@%s", MailgunDomain)

	body := fmt.Sprintf("Word of the day: %s   ++Thumbs Up: %v   --Thumbs Down: %v\n\n", word.Word, word.ThumbsUp, word.ThumbsDown)
	body += fmt.Sprintf("%s\n\n", word.Definition)
	body += fmt.Sprintf("Example: %s\n\n", word.Example)
	body += fmt.Sprintf("Courtesy of %s\n\n\n", word.Author)
	body += fmt.Sprintf("Brought to you by Chubbs Solutions.")

	gun := mailgun.NewMailgun(MailgunDomain, MailgunPublicAPIKey, "")
	m := mailgun.NewMessage(sender, subject, body, recipient)

	_, _, err := gun.Send(m)
	if err != nil {
		return err
	}
	// fmt.Printf("Response ID: %s\n", id)
	// fmt.Printf("Message from server: %s\n", response)

	return nil
}
