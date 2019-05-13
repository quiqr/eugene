package main

import (
  "fmt"
  "time"
  "os/exec"
  "log"

  //"github.com/anjannath/systray"
  "../systray_sm_fork_pim"

  "hugo-control/assets"
  "hugo-control/config"
  "hugo-control/hugo"
  "github.com/kr/pretty"
)

var (
  siteSubmenus = make(map[string]*systray.MenuItem)
  menuItemLiveUrl *systray.MenuItem
  menuItemToggleHugoServer *systray.MenuItem
  menuItemOpenConcept *systray.MenuItem
)

//var CurrentSite config.Site
//var CurrentConfig config.ConfigMulti

//var FatalError string

func main() {
  config.SetCurrentSite()
  log.Printf("Eugene Config: %# v", pretty.Formatter(config.CurrentSite))
  systray.Run(onReady, onExit)
}

/*
func SetCurrentSite(){
  var site_index int
  cfg, err := config.Read2()
  //log.Printf("Eugene Config ERR: %# v", pretty.Formatter(err))

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

  log.Printf("Eugene num sites: %d", len(cfg.Sites))
  log.Printf("Eugene current site index from config: %d", (cfg.Current_Site+1))
  log.Printf("Eugene current Site: %# v", pretty.Formatter(CurrentSite))
}
*/

func setCurrentSiteMenu(){

  //open Live Url
  if (config.CurrentSite.Live_Url != "") {
   menuItemLiveUrl = systray.AddMenuItem(fmt.Sprintf("Open %s", config.CurrentSite.Live_Url ) , "", 0)
  }

  //start Server Concept
  //stop Server Concept
  if hugo.HugoRunning(){
    menuItemToggleHugoServer = systray.AddMenuItem("stop lokale server", "", 0)
  } else {
    menuItemToggleHugoServer = systray.AddMenuItem("start lokale server", "", 0)
  }

  //open concept versie
  menuItemOpenConcept = systray.AddMenuItem(fmt.Sprintf("Open %s in conceptversie", config.CurrentSite.Name), "", 0)

  systray.AddSeparator()




}

func updateSiteMenu() {
  //open Live Url
  if (config.CurrentSite.Live_Url != "") {
   menuItemLiveUrl.SetTitle(fmt.Sprintf("Open %s", config.CurrentSite.Live_Url ))
  }

  //open concept versie
  menuItemOpenConcept.SetTitle(fmt.Sprintf("Open %s in conceptversie", config.CurrentSite.Name))
}

func renderMenu(){

  setCurrentSiteMenu()

  menuSelectSite := systray.AddSubMenu("Switch site")

  for _, site := range config.CurrentConfig.Sites {
    tmpSubmenuItem := menuSelectSite.AddSubMenuItem(site.Name,"", 0)
    siteSubmenus[site.Name] = tmpSubmenuItem
  }

  systray.AddSeparator()
  menuItemExit := systray.AddMenuItem("Quit", "", 0)

  for sName, sMenuItem := range siteSubmenus {

    go func(name string, siteMenuitem *systray.MenuItem) {
      for {
        <-siteMenuitem.OnClickCh()
        log.Println("Selecting %s", name)
        config.FindSiteIndexByName(name)
        config.SetCurrentSiteIndexByName(name)
        if hugo.HugoRunning(){
          hugo.KillHugo();
        }
        updateSiteMenu()
      }
    }(sName, sMenuItem)
  }

  go func() {
    for {
      time.Sleep(time.Second)
      if hugo.HugoRunning(){
        menuItemToggleHugoServer.SetTitle("Stop Server")
        menuItemOpenConcept.Enable()
      } else{
        menuItemToggleHugoServer.SetTitle("Start Server")
        menuItemOpenConcept.Disable()
      }
    }
  }()


  go func() {
    for {
      select {

      case <-menuItemLiveUrl.OnClickCh():
        exec.Command("/usr/bin/open", config.CurrentSite.Live_Url).Output()
      case <-menuItemOpenConcept.OnClickCh():
        exec.Command("/usr/bin/open", "http://localhost:1313").Output()
      case <-menuItemToggleHugoServer.OnClickCh():
        if hugo.HugoRunning(){
          hugo.KillHugo();
        } else{
          hugo.StartHugo();
          menuItemOpenConcept.Show()
        }
      case <-menuItemExit.OnClickCh():
        systray.Quit()
        return
      }
    }
  }()

}

func onReady() {
  fmt.Printf("OnReady: %v+\n", time.Now())
  systray.SetIcon(images.MonoData)

  if(config.FatalError != "") {
    systray.AddMenuItem(config.FatalError,"",0)
  } else {
    renderMenu()
  }
}

func onExit() {
  if hugo.HugoRunning(){
    hugo.KillHugo();
  }
}
