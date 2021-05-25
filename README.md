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
```
cd podcatch
chmod +x install.sh
./install.sh
podcatch
```

### Setup
Settings and Podcasts lists are automatically created in the  `~/.podcatch/` directory.
These files are in the TOML format, and the `podcasts.toml` file uses a `map[string]Podcast` interface.  
The automatically created files contain some example podcasts from which you should be able to understand the
required layout.

### Future plans
I'd really like to work on consuming more of the RSS data (i.e. description) and store alongside the MP3
with the hopes that Jellyfin can parse that to supply episode overviews and more data about it.  
I have tried working with this a little by creating an NFO file, but this has not worked as yet.

## Thanks for reading!
