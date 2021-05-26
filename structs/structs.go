package structs
import (
  "encoding/xml"
)
type Settings struct {
  Directory string
  Limit int
}
var DefaultSettings string = `
  Directory = "~/podcasts/"
  Limit = 10
`
var DefaultPodcasts string = `
[HelloInternet]
  Name = "Hello Internet"
  URL = "http://www.hellointernet.fm/podcast?format=rss"

[NSTAAF]
  Name = "No Such Thing as a Fish"
  URL = "https://audioboom.com/channels/2399216.rss"
`
type Podcast struct {
  URL string
  Name string
  Directory string
  RSS Rss
}
type Rss struct {
   XMLName     xml.Name `xml:"rss"`
   Version     string   `xml:"version,attr"`
   Channel     Channel  `xml:"channel"`
   Description string   `xml:"description"`
   Title       string   `xml:"title"`
   Link        string   `xml:"link"`
}
type Channel struct {
   XMLName     xml.Name `xml:"channel"`
   Title       string   `xml:"title"`
   Link        string   `xml:"link"`
   Description string   `xml:"description"`
   Items       []Item   `xml:"item"`
}
type Item struct {
   XMLName  xml.Name `xml:"item"`
   Title       string `xml:"title"`
   Episode     string `xml:"episode"`
   Link        string `xml:"link"`
   Description string `xml:"description"`
   PubDate     string `xml:"pubdate"`
   Guid        string `xml:"guid"`
   Media       Media  `xml:"enclosure"`
}
type Media struct {
  XMLName xml.Name `xml:"enclosure"`
  URL   string     `xml:"url,attr"`
  Type  string     `xml:"type,attr"`
}
