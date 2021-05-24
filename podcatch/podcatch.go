 package podcatch
import (
  "fmt"
  . "podcatch/structs"
  "github.com/pelletier/go-toml"
  "encoding/xml"
  "io/ioutil"
  "net/http"
  "regexp"
  "log"
)
var Version string = "0.1"
var Podcasts map[string]Podcast = make(map[string]Podcast)
func Start(){
  fmt.Printf("Starting PodCatch Version : %s...\r\n", Version )
  getPodcasts()
}
func getPodcasts(){
  if len(Podcasts) == 0 {
    getPodcastFiles()
  }
  for shortname,podcast := range Podcasts {
    fmt.Println(shortname)
    fmt.Printf("Checking RSS for %s...\r\n", podcast.Name)
    podcast.RSS = getRSS(podcast)
    downloadCasts(podcast)
  }
}
func getPodcastFiles() {
  content, err := ioutil.ReadFile("podcasts.toml")
  if err != nil {
    log.Fatal(err)
  }
  e := toml.Unmarshal(content,&Podcasts)
  if e != nil {
    log.Fatal(err)
  }
  fmt.Printf("Found %d podcasts.\r\n",len(Podcasts))
}
func getRSS(podcast Podcast) Rss {
  resp, err := http.Get(podcast.URL)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()
  html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
  return parseRSS(podcast,html)
}
func parseRSS(podcast Podcast, rssxml []byte) Rss {
  var rss Rss
  e := xml.Unmarshal(rssxml,&rss)
  if e != nil {
    log.Fatal(e)
  }
  return rss
}
func downloadCasts(podcast Podcast) {
  fmt.Println(podcast.RSS.Version)
  for _,item := range podcast.RSS.Channel.Items {
    fmt.Printf("Downloading '%s' from : %s.\r\n", item.Title, item.Media.URL)
    re := regexp.MustCompile(`[^0-9a-zA-Z-_]+`)
    filename := re.ReplaceAllString(item.Title,"_") + ".mp3"
    fmt.Println(filename)
    // fmt.Printf("%s - %s - %s\r\n", item.Title, item.Media.URL, item.Media.Type)
  }
}
