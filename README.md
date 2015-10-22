# urbano

[![Build Status](https://travis-ci.org/ChubbsSolutions/urbano.png)](https://travis-ci.org/ChubbsSolutions/urbano)

Get a new random Urban Dictionary word.

##Download the software

[Download](https://github.com/ChubbsSolutions/urbano/releases) the latest version of urbano for all major platforms.

##Compile and run the source

Requires Go 1.5 or newer (earlier versions untested). Remember to set the GOPATH variable.

```
git clone https://github.com/ChubbsSolutions/urbano
cd urbano
go get
go run urbano.go
```

##Usage

```NAME:
   urbano - Get a fresh Urban  Dictionary word in your inbox

USAGE:
   urbano [global options] command [command options] [arguments...]

VERSION:
   0.1

AUTHOR(S):
	Chubbs Solutions <urbano@chubbs.solutions>

COMMANDS:
   send, s	Get and Send a new word by email
   display, d	Display a new word
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

##Email

Urbano can now email you the word. You just need to set a Mailgun API key and domain. Here's how you would do it in Linux:

```
export MAILGUN_PUBLIC_API_KEY=example-key
export MAILGUN_DOMAIN=example-domain.com
```
Just add an email address as an argument.


##About

Crafted with :heart: in Indiana by [Chubbs Solutions] (http://chubbs.solutions).
