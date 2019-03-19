---
title: Renaming
menu: true
weight: 4
---
How to rename and organize media

You can organize and rename individual files or whole directories with supper

### Flags
`--action|-a`: The action to perform when renaming media. Can be one of
`move`, `symlink`, `hardlink` or `copy`. Default operation is to hardlink
files

`--extract|-x`: Additionally extract media from archives (zip/rar).

`--movies'-m`: Only rename movies

`--subtitles|-s`: Only rename subtitles

`--tvshows|-t`: Only rename tv shows

To see all applicable flags see: `supper ren --help`

### Examples
Rename all media in the `/media/downloads` folder using hardlink (default action): 
```bash
supper ren /media/downloads
```

Rename only movies in the `/media/downloads` folder using copy:
```bash
supper ren --action copy --movies /media/downloads
```

Rename all media in the `/media/downloads` folder and extract media from archives (rar/zip):
```bash
supper ren --extract /media/downloads
```