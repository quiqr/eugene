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
//  "github.com/kr/pretty"
)

var (
  siteSubmenus = make(map[string]*systray.MenuItem)
  liveUrlMenuItem *systray.MenuItem
  toggleHugoServerMenuItem *systray.MenuItem
)

//var CurrentSite config.Site
//var CurrentConfig config.ConfigMulti

//var FatalError string

func main() {
  config.SetCurrentSite()
  systray.Run(onReady, onExit)
  //go forever()
}

/*
func SetCurrentSite(){
  var site_index int
  cfg, err := config.Read2()
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

  log.Printf("Eugene num sites: %d", len(cfg.Sites))
  log.Printf("Eugene current site index from config: %d", (cfg.Current_Site+1))
  log.Printf("Eugene current Site: %# v", pretty.Formatter(CurrentSite))
}
*/

func setCurrentSiteMenu(){

  //open Live Url
  if (config.CurrentSite.Live_Url != "") {
   liveUrlMenuItem = systray.AddMenuItem(fmt.Sprintf("Open %s", config.CurrentSite.Live_Url ) , "", 0)
  }

  //start Server Concept

  //stop Server Concept

  //open concept versie


}

func updateSiteMenu() {
  //open Live Url
  if (config.CurrentSite.Live_Url != "") {
   liveUrlMenuItem.SetTitle(fmt.Sprintf("Open %s", config.CurrentSite.Live_Url ))
  }

  //start Server Concept

  //stop Server Concept

  //open concept versie
}


func renderMenu(){

  //var liveUrlMenuItem *systray.MenuItem
  //d0 = systray.AddMenuItem(fmt.Sprintf("Test Eugene Config") , "", 0)

  //cfg, err := config.Read()

  setCurrentSiteMenu()
  updateSiteMenu()

  if hugo.HugoRunning(){
    toggleHugoServerMenuItem = systray.AddMenuItem("stop lokale server", "", 0)
  } else {
    toggleHugoServerMenuItem = systray.AddMenuItem("start lokale server", "", 0)
  }

  menuItemOpenConcept := systray.AddMenuItem(fmt.Sprintf("Open %s in conceptversie", config.CurrentSite.Name), "", 0)

  go func() {
    for {
      time.Sleep(time.Second)
      if hugo.HugoRunning(){
        toggleHugoServerMenuItem.SetTitle("Stop Server")
        menuItemOpenConcept.Enable()
      } else{
        toggleHugoServerMenuItem.SetTitle("Start Server")
        menuItemOpenConcept.Disable()
      }
    }
  }()

  systray.AddSeparator()

  menuSelectSite := systray.AddSubMenu("Select site")

  for _, site := range config.CurrentConfig.Sites {
    tmpSubmenuItem := menuSelectSite.AddSubMenuItem(site.Name,"", 0)
    siteSubmenus[site.Name] = tmpSubmenuItem
  }

  systray.AddSeparator()
  exit := systray.AddMenuItem("Quit", "", 0)

  for sName, sMenuItem := range siteSubmenus {

    go func(name string, siteMenuitem *systray.MenuItem) {
      for {
        <-siteMenuitem.OnClickCh()
        log.Println("Selecting %s", name)
        config.FindSiteIndexByName(name)
        config.SetCurrentSiteIndexByName(name)
        updateSiteMenu()
      }
    }(sName, sMenuItem)
  }

  go func() {
    for {
      select {
        //case <-d0.OnClickCh():

      case <-liveUrlMenuItem.OnClickCh():
        exec.Command("/usr/bin/open", config.CurrentSite.Live_Url).Output()
      case <-menuItemOpenConcept.OnClickCh():
        exec.Command("/usr/bin/open", "http://localhost:1313").Output()
      case <-toggleHugoServerMenuItem.OnClickCh():
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

}

func forever() {
  for {
    //fmt.Printf("%v+\n", time.Now())
    //time.Sleep(time.Second)
  }
}
