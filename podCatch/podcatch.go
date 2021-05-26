package podcatch
import (
  "fmt"
  . "podcatch/structs"
  "github.com/pelletier/go-toml"
  "encoding/xml"
  "io"
  "io/ioutil"
  "net/http"
  "regexp"
  "os"
  "log"
  "strings"
  id3 "github.com/mikkyang/id3-go"
  v2 "github.com/mikkyang/id3-go/v2"
)
var Version string = "0.1"
var Settings Settings
var Podcasts map[string]Podcast = make(map[string]Podcast)
var donefile string
var podcatchdir string
var homedir string
var dbdir string
func Start(){
  fmt.Printf("Starting PodCatch Version : %s...\r\n", Version )
  getHomeDirs()
  getSettings()
  getPodcasts()
}
func getHomeDirs(){
  h, err := os.UserHomeDir()
  if err != nil {
    log.Fatal( err )
  }
  homedir = h
  podcatchdir = h + "/.podcatch/"
}
func getSettings(){
  settings := podcatchdir + "settings.toml"
  if !checkFileExists(settings){
    pwd, err := os.Getwd()
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println("Copying default settings.toml to user dir.")
    ok, err := copyFile(pwd + "/defaults/settings.toml",settings)
    if err != nil {
      log.Fatal(err)
    }
    if ok > 0 {
      fmt.Println("Copied.")
    }
  }
  content, err := ioutil.ReadFile(settings)
  if err != nil {
    log.Fatal(err)
  }
  e := toml.Unmarshal(content,&Settings)
  if e != nil {
    log.Fatal(err)
  }
  Settings.Directory = strings.Replace(Settings.Directory,"~",homedir,1)
  dbdir = Settings.Directory + ".db/"
  os.Mkdir(dbdir,0755)

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
  pcs := podcatchdir + "podcasts.toml"
  if !checkFileExists(pcs){
    pwd, err := os.Getwd()
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println("Copying default podcasts.toml to user dir.")
    ok, err := copyFile(pwd + "/defaults/podcasts.toml",pcs)
    if err != nil {
      log.Fatal(err)
    }
    if ok > 0 {
      fmt.Println("Copied.")
    }
  }
  content, err := ioutil.ReadFile(podcatchdir + "podcasts.toml")
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
        addId3(podcast.Name, item,dir + "/" + filename)
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
  db := dbdir + "complete"
  if checkCreate(db) {
    if len(donefile) < 1 {
      content, err := ioutil.ReadFile(db)
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
func addId3(artist string, item Item, file string) {
  fmt.Printf("Saving ID3 to %s\r\n",file)
  mp3File, err := id3.Open(file)
  if err != nil {
    log.Fatal(err)
  }
  defer mp3File.Close()
  if mp3File.Artist() == "" {
    mp3File.SetArtist(artist)
  }
  if mp3File.Album() == "" {
    mp3File.SetAlbum(artist)
  }
  if mp3File.Title() == "" {
    mp3File.SetTitle(item.Title)
  }
  if len(mp3File.Comments()) == 0 {
    ft := v2.V23FrameTypeMap["COMM"]
    textFrame := v2.NewTextFrame(ft, item.Description)
    mp3File.AddFrames(textFrame)
  }
}
func markAsReceived(item Item)  {
  db := dbdir + "complete"
  checkCreate(db)
  file, err := os.OpenFile(db, os.O_APPEND|os.O_WRONLY, 0755)
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
  db := dbdir + "error"
  checkCreate(db)
  file, err := os.OpenFile(db, os.O_APPEND|os.O_WRONLY, 0755)
  if err != nil {
    log.Println(err)
  }
  defer file.Close()
  content := fmt.Sprintf("%s\r\n%s",item.Title, item.Media.URL)
  if _, err := file.WriteString(content); err != nil {
    log.Fatal(err)
  }
}
func checkFileExists(file string) bool {
  if _, err := os.Stat(file); err == nil {
    // fmt.Println("Exists")
    // exists
    return true
  } else if os.IsNotExist(err) {
      // fmt.Println("Not Exists")
    // not exists
    return false
  } else {
    // fmt.Println("Maybe Exists, Maybe Not")
    return false
  // Schrodinger: file may or may not exist. See err for details.
  // Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
  }
  return false
}
func checkCreate(file string) bool {
  if checkFileExists(file) {
    return true
  } else {
    if createFile(file) {
      return true
    }
  }
  return false
}
func createFile(file string) bool {
  f, err := os.Create(file)
  if err != nil {
    log.Fatal(err)
    return false
  }
  defer f.Close()
  return true
}
func copyFile(src, dst string) (int64, error) {
  sourceFileStat, err := os.Stat(src)
  if err != nil {
    return 0, err
  }
  if !sourceFileStat.Mode().IsRegular() {
    return 0, fmt.Errorf("%s is not a regular file", src)
  }
  source, err := os.Open(src)
  if err != nil {
    return 0, err
  }
  defer source.Close()
  destination, err := os.Create(dst)
  if err != nil {
    return 0, err
  }
  defer destination.Close()
  nBytes, err := io.Copy(destination, source)
  return nBytes, err
}
