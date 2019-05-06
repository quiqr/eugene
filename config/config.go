package config

import (
  "os"
  "path"
  "path/filepath"
  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"
)

const DefaultDir string = "~/.bitbar-hugo"
const yamlFile = "config.yml"

func Dir() string {
  cfgPath, _ := homedir.Expand(DefaultDir)
  return path.Clean(cfgPath)
}

func File() string {
  return path.Clean(filepath.Join(Dir(), yamlFile))
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
