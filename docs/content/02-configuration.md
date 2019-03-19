---
title: Configuration
menu: true
weight: 2
---
Learn how to configure supper 

You can specify the configuration of how supper function in a `.yaml` configuration file.
The configuration file can be located in different places according to operating system and precendence.
Here is an overview:

* Global configuration
  - Unix: `/etc/supper/supper.yaml`
  - Windows: `%HOMEPATH%\AppData\Roaming\Supper\supper.yaml`
* Local configuration
  - `.supper.yaml` found in current working directory

The local configuration has precendence over the global configuration. If supper is run with 
the `--config` argument, then that configuration file is used above all else.

## Example
Below the default configuration for supper is shown. Most of the configuration is
self-explanatory, however some details will be covered in the following sections.
```yaml
# Supper configuration file

# Satisfy the following languages when downloading subtitles
languages:
  - en
  - es
  - de

# Path to store application logs
logfile: /var/log/supper/supper.log

# Download only hearing impaired subtitles
impared: false

# Bind web server to port
port: 5670

# Base path for reverse proxy
proxypath: "/"

# Movie collection configuration
movies:
  # Directory to store movie collection
  directory: /media/movies

  # Template to use for renaming movies
  template: >
    {{ .Movie }} ({{ .Year }})/
    {{ .Movie }} ({{ .Year }}) {{ .Quality }}

# TV show collection configuration
tvshows:
  # Directory to store TV shows
  directory: /media/tvshows

  # Template to use for renaming TV shows
  template: >
    {{ .TVShow }}/Season {{ .Season | pad }}/
    {{ .TVShow }} - S{{ .Season | pad }}E{{ .Episode | pad }} - {{ .Name }}

# Plugins are run after downloading a subtitle. The plugin is a simple shell
# command which is given the .srt file path in the SUBTITLE environment variable
plugins:
  # - name: my-plugin-name
  #   exec: echo $SUBTITLE
```

## Templates
Templates are used to rename movie and TV series into folder/file names. The templating
scheme uses the golang templating language and is highly customizable. You may define
subfolders in your templating scheme using the path seperator `/`.

### Movies
The following directives are available for movie templates:

| Directive   | Description                       | Example           |
| :-------:   | :-------------------------------: | :---------------: |
| `.Movie`    | The name of the movie             | `Inception`       |
| `.Year`     | The release year of the movie     | `2010`            |
| `.Quality`  | Quality of the movie release      | `720p`            |
| `.Codec`    | Codec of the movie release        | `h264`            |
| `.Source`   | Source of the movie release       | `BluRay`          |
| `.Group`    | Release group of the movie        | N/A               |

**Example:**
```handlebars
{{ .Movie }} ({{ .Year }})/{{ .Movie }} ({{ .Year }}) {{ .Quality }}
```


## TV shows
The following directovies are available for tv show templates:

| Directive   | Description                       | Example           |
| :-------:   | :-------------------------------: | :---------------: |
| `.TVShow`   | The name of the TV show           | `Game of Thrones` |
| `.Name`     | The name of the episode           | `Pilot`           |
| `.Season`   | Season number                     | `1`               |
| `.Episode`  | Episode number                    | `1`               |
| `.Quality`  | Quality of the movie release      | `720p`            |
| `.Codec`    | Codec of the movie release        | `h264`            |
| `.Source`   | Source of the movie release       | `BluRay`          |
| `.Group`    | Release group of the movie        | N/A               |

**Example:**
```handlebars
{{ .TVShow }}/Season {{ .Season | pad }}/
{{ .TVShow }} - S{{ .Season | pad }}E{{ .Episode | pad }} - {{ .Name }}
```

## Template Functions
You can utilize template functions to manipulate with the data in your templating schemes

#### `pad`
Pads the number with zeros to make up exactly two characters total. Useful for padding season and episode numbers.

**Example** 
```handlebars
{{ .Season | pad }}
```
will output `01` for season one instead of just `1`.