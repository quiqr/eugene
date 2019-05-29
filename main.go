package main

import (
  "bytes"
  "strings"
  "fmt"
  "time"
  "os/exec"
  "log"

  "github.com/anjannath/systray"
  //"../systray_sm_fork_pim"

  "eugene/assets"
  "eugene/config"
  "eugene/hugo"
  "github.com/kr/pretty"
)

const AppVersion = "0.0.4"

var (
  siteSubmenus = make(map[string]*systray.MenuItem)
  menuItemLiveUrl *systray.MenuItem
  menuItemPublish *systray.MenuItem
  menuItemToggleHugoServer *systray.MenuItem
  menuItemOpenConcept *systray.MenuItem
  menuItemExit *systray.MenuItem
  menuSelectSite *systray.MenuItem
  menuToggles *systray.MenuItem
  menuConfig *systray.MenuItem
)

func main() {
  config.SetCurrentSite()
  log.Printf("Eugene Config: %# v", pretty.Formatter(config.CurrentSite))
  systray.Run(onReady, onExit)
}

func onReady() {
  fmt.Printf("OnReady: %v+\n", time.Now())
  systray.SetIcon(images.EugeneMonoData)

  if(config.FatalError != "") {
    systray.AddMenuItem(config.FatalError,"",0)
  } else {
    renderCMSMenu()
  }

  renderFooterMenu()
}

func onExit() {
  if hugo.HugoRunning(){
    hugo.KillHugo();
  }
}

func renderCMSMenu(){

  setCurrentSiteMenu()
  systray.AddSeparator()
  stagingMenu()
  switchSitesMenu()
  systray.AddSeparator()
  togglesMenu()
  configMenu()
  systray.AddSeparator()
  listenToServer()
  handleCMSMenuClicks()
}

func renderFooterMenu(){
  systray.AddMenuItem(fmt.Sprintf("Version: %s", AppVersion), "", 0)
  menuItemExit = systray.AddMenuItem("Quit", "", 0)
  handleFooterMenuClicks()
}

func updateCMSMenu() {
  //open Live Url
  if (config.CurrentSite.Live_Url != "") {
   menuItemLiveUrl.SetTitle(fmt.Sprintf("Open %s", config.CurrentSite.Live_Url ))
  }

  //open concept versie
  menuItemOpenConcept.SetTitle(fmt.Sprintf("Open %s in conceptversie", config.CurrentSite.Name))

  //open Live Url
  log.Println(config.CurrentSite.Live_Publishing_Command)
  if (config.CurrentSite.Live_Url != "" && config.CurrentSite.Live_Publishing_Command != "") {
   menuItemPublish.SetTitle(fmt.Sprintf("Publish to %s", config.CurrentSite.Live_Url ))
   menuItemPublish.Enable()
  } else{
   menuItemPublish.SetTitle(fmt.Sprintf("Conf not complete for Publishing"))
   menuItemPublish.Disable()
  }
}

func stagingMenu(){


}

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

  //open Live Url
  if (config.CurrentSite.Live_Url != "" && config.CurrentSite.Live_Publishing_Command != "") {
   menuItemPublish = systray.AddMenuItem(fmt.Sprintf("Publish to %s", config.CurrentSite.Live_Url ) , "", 0)
  } else{
   menuItemPublish = systray.AddMenuItem(fmt.Sprintf("Config not complete for Publishing") , "", 0)
  }

  //open concept versie
  menuItemOpenConcept = systray.AddMenuItem(fmt.Sprintf("Open %s in conceptversie", config.CurrentSite.Name), "", 0)
}

func configMenu(){
  menuConfig = systray.AddSubMenu("Configuration")
  menuOpenConfigFile := menuConfig.AddSubMenuItem("Open configuration file","", 0)

  go func() {
    for {
      select {

      case <-menuOpenConfigFile.OnClickCh():
        var errOut bytes.Buffer
        c := exec.Command("/usr/bin/open", config.File2())
        c.Dir = config.CurrentSite.Live_Hugo_Output_Dir
        c.Stderr = &errOut
        out, err := c.Output()
        outStr := strings.TrimSpace(string(out))
        if err != nil {
          err = fmt.Errorf("open: error=%q stderr=%s", err, string(errOut.Bytes()))
        }
        log.Printf("open Result: %# v", pretty.Formatter(outStr))
      }
    }
  }()
}

func togglesMenu(){
  menuToggles = systray.AddSubMenu("Toggles")

  menuToggleShowConcept := menuToggles.AddSubMenuItem("Show draft items","", 0)

  go func() {
    for {
      select {

      case <-menuToggleShowConcept.OnClickCh():

        if(config.ShowDraftItems){
          config.ShowDraftItems = false
          menuToggleShowConcept.SetTitle("Hide draft items")

        } else {
          config.ShowDraftItems = true
          menuToggleShowConcept.SetTitle("Show draft items")
        }
        if hugo.HugoRunning(){
          hugo.KillHugo();
        }
      }
    }
  }()
}

func switchSitesMenu(){
  if(len(config.CurrentConfig.Sites)>1){
    menuSelectSite = systray.AddSubMenu("Switch site")
    for _, site := range config.CurrentConfig.Sites {
      tmpSubmenuItem := menuSelectSite.AddSubMenuItem(site.Name,"", 0)
      siteSubmenus[site.Name] = tmpSubmenuItem
    }

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
          updateCMSMenu()
        }
      }(sName, sMenuItem)
    }
  }
}

func listenToServer(){
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
}

func gitCommand(args ...string){
  var errOut bytes.Buffer
  c := exec.Command("/usr/bin/git", args...)
  c.Dir = config.CurrentSite.Live_Hugo_Output_Dir
  c.Stderr = &errOut
  out, err := c.Output()
  outStr := strings.TrimSpace(string(out))
  if err != nil {
    err = fmt.Errorf("git: error=%q stderr=%s", err, string(errOut.Bytes()))
  }
  log.Printf("git args: %# v", pretty.Formatter(args))
  log.Printf("Publish Result: %# v", pretty.Formatter(outStr))
}

func handleCMSMenuClicks(){
  go func() {
    for {
      select {

      case <-menuItemLiveUrl.OnClickCh():
        exec.Command("/usr/bin/open", config.CurrentSite.Live_Url).Output()
      case <-menuItemPublish.OnClickCh():

        //open Live Url
        if (config.CurrentSite.Live_Url != "" && config.CurrentSite.Live_Publishing_Command != "") {
          gitCommand("add", ".")
          gitCommand("commit", "-m", "Published with Hugo Control", "-a")
          gitCommand("push")
        }

      case <-menuItemOpenConcept.OnClickCh():
        exec.Command("/usr/bin/open", "http://localhost:1313").Output()
      case <-menuItemToggleHugoServer.OnClickCh():
        if hugo.HugoRunning(){
          hugo.KillHugo();
        } else{
          hugo.StartHugo();
          menuItemOpenConcept.Show()
        }
      }
    }
  }()
}

func handleFooterMenuClicks(){
  go func() {
    for {
      select {
      case <-menuItemExit.OnClickCh():
        systray.Quit()
        return
      }
    }
  }()
}

