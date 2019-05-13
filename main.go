package main

import (
  "fmt"
  "time"
  "os/exec"
  "log"

  "github.com/anjannath/systray"
  "hugo-control/assets"
  "hugo-control/config"
  "hugo-control/hugo"
  "github.com/kr/pretty"
)

var (
  siteSubmenus = make(map[string]*systray.MenuItem)
)

var CurrentSite config.Site
var CurrentConfig config.ConfigMulti

var FatalError string

func main() {
  FatalError = ""
  SetCurrentSite()
  systray.Run(onReady, onExit)
  go forever()
}

func SetCurrentSite(){
  var site_index int
  cfg, err := config.Read2()
  log.Printf("Eugene Config ERR: %# v", pretty.Formatter(err))
  log.Printf("Eugene Config: %# v", pretty.Formatter(cfg))

  if( err != nil){
    FatalError = "Can't read Config"
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


type MenuAction struct {
  start *systray.MenuItem
  stop  *systray.MenuItem
}

func renderMenu(){
  var e1 *systray.MenuItem
  var m0 *systray.MenuItem

  if ! config.ConfigFileExists() {
    e1 = systray.AddMenuItem("No config file","",0)
  }

  var d0 *systray.MenuItem
  d0 = systray.AddMenuItem(fmt.Sprintf("Test Eugene Config") , "", 0)

  cfg, err := config.Read()
  if err != nil {
    e1 = systray.AddMenuItem("Error reading config file","",0)
  } else {
    _ = e1

    if (cfg.LiveUrl != "") {
      m0 = systray.AddMenuItem(fmt.Sprintf("Open %s", cfg.LiveUrl ) , "", 0)
    } else {
      _ = m0
    }

    var toggleHugoServer *systray.MenuItem
    if hugo.HugoRunning(){
      toggleHugoServer = systray.AddMenuItem("stop lokale server", "", 0)
    } else {
      toggleHugoServer = systray.AddMenuItem("start lokale server", "", 0)
    }

    menuItemOpenConcept := systray.AddMenuItem(fmt.Sprintf("Open %s in conceptversie", cfg.SiteName), "", 0)

    go func() {
      for {
        //fmt.Printf("%v+\n", time.Now())
        time.Sleep(time.Second)
        if hugo.HugoRunning(){
          toggleHugoServer.SetTitle("Stop Server")
          menuItemOpenConcept.Enable()
        } else{
          toggleHugoServer.SetTitle("Start Server")
          menuItemOpenConcept.Disable()
        }
      }
    }()

    systray.AddSeparator()

    SM_SelectSite := systray.AddSubMenu("Select site")

    for _, site := range CurrentConfig.Sites {
      submenu := SM_SelectSite.AddSubMenuItem(site.Name,"", 0)
      siteSubmenus[site.Name] = submenu

    }



    systray.AddSeparator()
    exit := systray.AddMenuItem("Quit", "", 0)

    go func() {
      for {
        select {
        case <-d0.OnClickCh():

        case <-m0.OnClickCh():
          exec.Command("/usr/bin/open", cfg.LiveUrl).Output()
        case <-menuItemOpenConcept.OnClickCh():
          exec.Command("/usr/bin/open", "http://localhost:1313").Output()
        case <-toggleHugoServer.OnClickCh():
          if hugo.HugoRunning(){
            hugo.KillHugo();
          } else{
            hugo.StartHugo();
            menuItemOpenConcept.Show()
          }
        case <-exit.OnClickCh():
          systray.Quit()
          return
        }
      }
    }()
  }

}

func onReady() {
  fmt.Printf("OnReady: %v+\n", time.Now())
  systray.SetIcon(images.MonoData)

  if(FatalError != "") {
    systray.AddMenuItem(FatalError,"",0)
  } else {
    renderMenu()
  }
}

func onExit() {

}

func forever() {
  for {
    //fmt.Printf("%v+\n", time.Now())
    //time.Sleep(time.Second)
  }
}
