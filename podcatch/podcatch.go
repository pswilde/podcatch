 package podcatch
import (
  "fmt"
  . "podcatch/structs"
  "github.com/pelletier/go-toml"
  "encoding/xml"
  "io/ioutil"
  "net/http"
  "regexp"
  "os"
  "log"
  "strings"
)
var Version string = "0.1"
var Settings Settings
var Podcasts map[string]Podcast = make(map[string]Podcast)
var donefile string
func Start(){
  fmt.Printf("Starting PodCatch Version : %s...\r\n", Version )
  getSettings()
  getPodcasts()
}
func getSettings(){
  content, err := ioutil.ReadFile("settings.toml")
  if err != nil {
    log.Fatal(err)
  }
  e := toml.Unmarshal(content,&Settings)
  if e != nil {
    log.Fatal(err)
  }
}
func getPodcasts(){
  if len(Podcasts) == 0 {
    getPodcastFiles()
  }
  for shortname,podcast := range Podcasts {
    podcast.Directory = shortname
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
  count := 0
  for _,item := range podcast.RSS.Channel.Items {
    if count >= Settings.Limit {
      break
    }
    if !podcastDownloaded(item){
      fmt.Printf("Downloading '%s %s' from : %s.\r\n", item.Episode, item.Title, item.Media.URL)
      re := regexp.MustCompile(`[^0-9a-zA-Z-_]+`)
      filename := item.Episode + re.ReplaceAllString(item.Title,"_") + ".mp3"
      dir := Settings.Directory + podcast.Directory
      err := os.Mkdir(dir, 0777)
      if err != nil && err.Error() != fmt.Sprintf("mkdir %s: file exists",dir){
        log.Fatal(err)
      }
      ok := downloadMp3(item.Media.URL, dir + "/" + filename)
      if ok {
        // createNFO(item, strings.Replace(dir + "/" + filename,".mp3",".nfo",1))
        markAsReceived(item)
      } else {
        markAsErrored(item)
      }
    } else {
      fmt.Printf("Skipping '%s' - already downloaded\r\n", item.Title)
    }
    count = count + 1
  }
}
func podcastDownloaded(item Item) bool {
  if len(donefile) < 1 {
    content, err := ioutil.ReadFile(".db/complete")
    if err != nil {
      log.Fatal(err)
    }
    donefile = string(content)
  }
  if strings.Contains(donefile,item.Title){
    return true
  }
  if strings.Contains(donefile,item.Media.URL){
    return true
  }
  return false
}
func downloadMp3(url string, file string) bool {
  ok := false
  resp, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()
  data, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }
  err = ioutil.WriteFile(file, data, 0775)
  if err != nil {
    log.Fatal(err)
  }
  ok = true
  return ok
}
func createNFO(item Item, file string) {
  fmt.Printf("Saving NFO file %s",file)
  var nfo NFO
  nfo.Title = item.Title
  nfo.Outline = item.Description
  nfo.Aired = item.PubDate
  data, err := xml.Marshal(nfo)
  if err != nil {
    log.Fatal(err)
  }
  err = ioutil.WriteFile(file, data, 0775)
  if err != nil {
    log.Fatal(err)
  }
}
func markAsReceived(item Item)  {
  os.Mkdir(".db", 0777)
  file, err := os.OpenFile(".db/complete", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
  if err != nil {
    log.Println(err)
  }
  defer file.Close()
  content := fmt.Sprintf("%s - %s\r\n",item.Title, item.Media.URL)
  if _, err := file.WriteString(content); err != nil {
    log.Fatal(err)
  }
}
func markAsErrored(item Item)  {
  os.Mkdir(".db", 0777)
  file, err := os.OpenFile(".db/error", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
  if err != nil {
    log.Println(err)
  }
  defer file.Close()
  content := fmt.Sprintf("%s\r\n%s",item.Title, item.Media.URL)
  if _, err := file.WriteString(content); err != nil {
    log.Fatal(err)
  }
}
