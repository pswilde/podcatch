package structs
import (
  "encoding/xml"
)
type Settings struct {
  Directory string
  Limit int
}
type Podcast struct {
  URL string
  Name string
  Directory string
  RSS Rss
}
type NFO struct {
  XMLName      xml.Name `xml:"podcast"`
  Title       string    `xml:"title"`
  Outline     string    `xml:"outline"`
  Aired       string    `xml:"aired"`
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
