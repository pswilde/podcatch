# PodCatch - a simple Podcast downloader

PodCatch is a _very_ simple Podcast downloader written in GoLand built to help me download podcasts
and store them in my media directory so [Jellyfin](https://jellyfin.org/) can
index them and thus can be listened to using Jellyfin's various apps.  
Due to how PodCatch is configurable, I'm sure it would have other uses too.

### Why PodCatch over other podcast downloaders?
I tried a few other podcast downloaders and largely they were fine, in fact
PodCatch definitely borrows some nice ideas from other podcast downloaders. My
personal issue with the others is that the filename they download as is not always
"friendly"; that is to say, if there's no ID3 data containing the podcast title,
trying to organise a list of podcasts with a _GUID_.mp3 filename in Jellyfin is not easy.  
To tackle this issue, PodCatch looks up the title of the Podcast from its RSS feed,
strips it of invalid characters, and uses that name as the file name. This allows for a
much easier task when organising podcasts.

### Install
Very simple install:
```git clone this repo
cd podcatch
go build -o podcatch-bin
