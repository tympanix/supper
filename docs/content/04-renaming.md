---
title: Renaming
menu: true
weight: 4
---
How to rename and organize media

You can organize and rename individual files or whole directories with supper

### Examples:
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