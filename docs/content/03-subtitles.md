---
title: Subtitles
menu: true
weight: 3
---
How to download subtitles with supper

You can download subtitles with supper for individual files of whole directories easily.

### Flags
`--lang|-l`:
Which language(s) subtitles will be downloaded in. Multiple languages
can be specified using the flag multiple times. If no languages are specified
those is the configuration file will be sued by default.

`--modified|-m`: Only download subtitles for media modified since the given duration.
Durations are specified using numbers and letters (e.g. `2d12h30m`) 

`--score|-s`: Specifies a minimum score for any subtitle to be downloaded. If a subtitle is
not available satisfying the minimum score, no subtitle will be downloaded. The score of
a subtitle is determined by the likelyhood of the subtitle being synchronized to the media
according to an internal scoring algorithm. Values are given in percent (withput the percent sign)

`--imapired|-i`: Only download hearing impaired subtitles. By default Supper will not consider
or download hearing impaired subtitles. This flagg reverses this behaviour.

`--limit|-l`: Limit the number of media to process. A default values is applied as a safeguard
agains accidental filepaths. Therefore this flag must be specified for large quantaties of media.
Specifying a negative number will disable the limit.

To see all applicable flags see `supper sub --help`.

## Languages:
Here is a list of the supported languages. The language `tag` specified in the table
below can be used as an argument to the `--lang|-l` flag.

| Language         | Tag  | Language         | Tag  | Language         | Tag  |
| ---              | ---  | ---              | ---  | ---              | ---  |
| Albanian         | `sq` | Greek            | `el` | Portuguese       | `pt` |
| Latvian          | `lv` | Hebrew           | `he` | Portuguese       | `pt` |  
| Arabic           | `ar` | Hindi            | `hi` | Romanian         | `ro` |  
| Armenian         | `hy` | Hungarian        | `hu` | Russian          | `ru` |  
| Azerbaijani      | `az` | Icelandic        | `is` | Serbian          | `sr` |  
| Bangla           | `bn` | Indonesian       | `id` | Slovak           | `sk` |  
| Bulgarian        | `bg` | Italian          | `it` | Slovenian        | `sl` |  
| Catalan          | `ca` | Japanese         | `ja` | Spanish          | `es` |  
| Chinese          | `zh` | Korean           | `ko` | Swahili          | `sw` |  
| Croatian         | `hr` | Lithuanian       | `lt` | Swedish          | `sv` |  
| Czech            | `cs` | Macedonian       | `mk` | Tamil            | `ta` |  
| Danish           | `da` | Malay            | `ms` | Telugu           | `te` |  
| Dutch            | `nl` | Malayalam        | `ml` | Thai             | `th` |  
| English          | `en` | Mongolian        | `mn` | Turkish          | `tr` |  
| Finnish          | `fi` | Norwegian        | `no` | Ukrainian        | `uk` |  
| French           | `fr` | Persian          | `fa` | Urdu             | `ur` |  
| Georgian         | `ka` | Persian          | `fa` | Vietnamese       | `vi` |  
| German           | `de` | Polish           | `pl` |                  |      | 


## Examples:
Download subtitles for all media in the `/media/movies` folder in english, german and spanish for files 
added or modified within the last 24 hours:
```bash
supper sub -l en -l de -l es -m 24h /media/movies
```

Download english subtitles for the file `tvshow.mp4` if a subtitle can be found with a score higher or equal to 75%:
```bash
supper sub -l en -s 75 tvshow.mp4
```

Download and overwrite existing english subtitles for all media in `/media/tvshows`
```bash
supper sub -l en --force /media/tvshows
```