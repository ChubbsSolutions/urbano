package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/ChubbsSolutions/urbano/objects"
	"github.com/codegangsta/cli"
	"github.com/mailgun/mailgun-go"
	"github.com/ttacon/chalk"
)

//Define constants
const version string = "0.2"
const author string = "Chubbs Solutions"
const email string = "urbano@chubbs.solutions"
const appName string = "urbano"
const appDescription string = "Get a fresh Urban  Dictionary word in your inbox"

//MailgunPublicAPIKey key for the mail service.
var MailgunPublicAPIKey = os.Getenv("MAILGUN_PUBLIC_API_KEY")

//MailgunPrivateAPIKey key for the mail service.
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
					os.Exit(-1)
				}

				err = displayNumbers(word)
				if err != nil {
					fmt.Println(chalk.Red, err)
					os.Exit(-1)
				}

				if len(c.Args()) == 1 {
					recipient := c.Args()[0]
					subject := fmt.Sprintf("Urban Dictionary Word of the day for %s", t.Format("Jan 02, 2006"))
					err = emailWord(word, subject, recipient)
					if err != nil {
						fmt.Println(chalk.Red, err)
						os.Exit(-1)
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
					os.Exit(-1)
				}

				err = displayNumbers(word)
				if err != nil {
					fmt.Println(chalk.Red, err)
					os.Exit(-1)
				}
			},
		},
	}

	app.Run(os.Args)

}

//getNewWord gets a random UD word
func getNewWord() (objects.WordData, error) {
	var UDURL = "http://api.urbandictionary.com/v0/random"
	wd := objects.WordDataSlice{}
	var word objects.WordData

	resp, err := http.Get(UDURL)
	if err != nil {
		return word, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		return word, err
	}

	err = json.Unmarshal([]byte(string(data)), &wd)
	if err != nil {
		return word, err
	}

	tu := 600

	for _, element := range wd.List {
		if element.ThumbsUp > tu {
			word = element
		}
	}
	return word, nil
}

func displayNumbers(word objects.WordData) error {

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
