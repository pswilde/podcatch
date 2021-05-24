package structs
import (
  "encoding/xml"
)
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
