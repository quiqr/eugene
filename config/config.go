package config

import (
  "os"
  "path"
  "log"
  "path/filepath"
  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"
)

const DefaultDir string = "~/.bitbar-hugo"
const yamlFile = "config.yml"
const yamlFile2 = "eugene-conf.yml"
var CurrentSite Site
var CurrentConfig ConfigMulti
var FatalError = ""

func SetCurrentSite(){
  var site_index int
  cfg, err := Read2()
  //log.Printf("Eugene Config ERR: %# v", pretty.Formatter(err))
  //log.Printf("Eugene Config: %# v", pretty.Formatter(cfg))

  if( err != nil){
    FatalError = "Can't read configfile"
  }

  if(len(cfg.Sites) == 0){
    FatalError = "No sites configured"
    return
  }

  if(cfg.Current_Site >= len(cfg.Sites)){
    site_index = 0
  } else{
    site_index = cfg.Current_Site
  }

  CurrentConfig = cfg
  CurrentSite = cfg.Sites[site_index]

  /*
  log.Printf("Eugene num sites: %d", len(cfg.Sites))
  log.Printf("Eugene current site index from config: %d", (cfg.Current_Site+1))
  log.Printf("Eugene current Site: %# v", pretty.Formatter(CurrentSite))
  */
}




func Dir() string {
  cfgPath, _ := homedir.Expand(DefaultDir)
  return path.Clean(cfgPath)
}

func File() string {
  return path.Clean(filepath.Join(Dir(), yamlFile))
}

func File2() string {
  return path.Clean(filepath.Join(Dir(), yamlFile2))
}

func ConfigDirExists() bool {
  dir := Dir()
  if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
    return true
  } else {
    return false
  }
}

func ConfigFileExists() bool {
  file := File()
  if stat, err := os.Stat(file); err == nil && stat.Mode().IsRegular() {
    return true
  } else {
    return false
  }
}

func FindSiteIndexByName(name string) int {

  for index, site := range CurrentConfig.Sites {
    //log.Println(site)
    if(name == site.Name){
      return index
    }
  }

  return -1
}

func SetCurrentSiteIndexByName(name string){
  newIndex := FindSiteIndexByName(name)
  log.Println(newIndex)
  if(newIndex >= 0){

    CurrentConfig.Current_Site = newIndex
    CurrentSite = CurrentConfig.Sites[newIndex]
    //  log.Println(newIndex)
    //   SetCurrentSite()
  }
}

// EnsureConfigDir creates a configDir() if it doesn't already exist
func EnsureConfigDir() error {
  dir := Dir()
  if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
    return nil
  }
  err := os.Mkdir(dir, 0700)
  if err != nil {
    return err
  }
  return nil
}

// Read config from the specified dir returning a slice of OpenFaaS instances
func Read2() (ConfigMulti, error) {

  viper.SetConfigName("eugene-config")
  viper.SetConfigFile(File2())
  viper.SetConfigType("yaml")

  err1 := viper.ReadInConfig()
  if  err1 != nil {
    log.Printf("err: %s\n", err1)
  }

  var config ConfigMulti
  if err := viper.Unmarshal(&config); err != nil {
    log.Printf("err2: %s\n", err)
    return config, err
  }

  return config, nil
}

/*func Write() (error) {
}
*/

// Read config from the specified dir returning a slice of OpenFaaS instances
func Read() (Config, error) {

  viper.SetDefault("hugo_src_dir", "")
  viper.SetDefault("site_name", "website")
  viper.SetDefault("live_url", "")

  viper.SetConfigName("config")
  viper.SetConfigFile(File())

  err := viper.ReadInConfig()

  viper.SetConfigType("yaml")

  if err == nil {
    return Config{
      HugoDir: viper.Get("hugo_src_dir").(string),
      SiteName: viper.Get("site_name").(string),
      LiveUrl: viper.Get("live_url").(string),
    }, nil
  } else {
    return Config{}, err
  }
}
