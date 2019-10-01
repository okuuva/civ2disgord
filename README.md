# civ2disgord

A little script that helps Civilization VI and Discord understand each other by converting turn notifications from
Civilization VI Play by Cloud into the format used by Discord webhook bot. Not only that, but it can do it in a way
that pings the user with their Discord handle instead of just shouting the Steam nick.

## Building

If you want to build the standalone script, run
```
go build main.go logging.go checkers.go cmdline.go -o civ2disgord
```

For running it in AWS Lambda, run
```
go build awslambda.go logging.go checkers.go cmdline.go -o civ2disgord
```

## Running
When called, the script converts the JSON message from Civ VI PbC into a Discord webhook message and sends it to the
said webhook. The message is provided by CLI argument. Mappings (Steam nick -> Discord ID, Game name -> webhook url)
can be read from either YAML config file or environment variables. Sample configuration file can be found from
[config.yml](./civ2disgord/config.yml).

Env variable keys need to be base64 encoded in order to support all the funky character people have on their Steam nicks. 
In order to make the conversion less painful, you can convert valid config file into environment variables by providing
config file and running the script with `--to-env` flag. The output is in .env file format and is printed to the stdout.
Check `civ2disgord --help` for all the CLI flags.

Currently the standalone script only supports oneshot conversions, i.e. it needs to be called separately each time
you need to convert and send messages to Discord. Daemon mode is on the TODO-list, pull-requests welcome.

## TODO
* Daemon mode for the standalone script
* Google App Engine script similar to awslambda.go
* makefile for easier builds
