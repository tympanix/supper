---
title: Subtitles
menu: true
weight: 3
---
How to download subtitles with supper

You can download subtitles with supper for individual files of whole directories easily

## Examples:
Download subtitles for all media in the `/media/movies` folder in english, german and spanish for files 
added or modified in the last 24 hours:
```bash
supper sub -len -lde -les -m 24h /media/movies
```

Download english subtitles for the file `tvshow.mp4` if the subtitles has a score better than 75%:
```bash
supper sub -len -s 75 tvshow.mp4
```

Download and overwrite existing english subtitles for all media in `/media/tvshows`
```bash
supper sub -len --force /media/tvshows
```