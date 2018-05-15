# Supper 
**A blazingly fast multimedia manager**

[![Build Status](https://semaphoreci.com/api/v1/tympanix/supper/branches/master/shields_badge.svg)](https://semaphoreci.com/tympanix/supper)
[![Go Report Card](https://goreportcard.com/badge/github.com/tympanix/supper)](https://goreportcard.com/report/github.com/tympanix/supper)
[![codecov](https://codecov.io/gh/tympanix/supper/branch/master/graph/badge.svg)](https://codecov.io/gh/tympanix/supper)

## Features
 - [x] Download subtitles for movies & TV shows
 - [x] Rename and organize your media collection
 - [x] Custom renaming templates
 - [x] Extract media from archives (zip/rar)
 - [x] Web interface to manage your media collection

## Help
```
A blazingly fast multimedia manager

Usage:
  supper [command]

Available Commands:
  help        Help about any command
  rename      Rename and process media files
  subtitle    Download subtitles for media
  version     Print the version number and exit
  web         Listen and serve the web application

Flags:
      --config string    load config file at specified path
      --dry              test run command without any effects
      --force            overwrite media files on conflicts
  -h, --help             help for supper
      --logfile string   store application logs in specified path
      --strict           exit the application on any error
  -v, --verbose          enable verbose logging
      --version          show the application version and exit

Use "supper [command] --help" for more information about a command.
```

## Credit
[FileBot](https://www.filebot.net) for inspiration (or frustration with it being slow and clunky, sorry)
