package main

import (
  "fmt"
  "time"
  "os/exec"

  "github.com/anjannath/systray"
  "hugo-control/assets"
  "hugo-control/config"
  "hugo-control/hugo"
)

var (
  menuTitles          = []string{"hugo 1", "hugo 2", "hugo 3"}
  submenus            = make(map[string]*systray.MenuItem)
  submenusToMenuItems = make(map[string]MenuAction)
)

func main() {
  systray.Run(onReady, onExit)
  go forever()
}

type MenuAction struct {
  start *systray.MenuItem
  stop  *systray.MenuItem
}

func onReady() {
  var e1 *systray.MenuItem
  var m0 *systray.MenuItem

  //serverRunning := make(chan bool)

  systray.SetIcon(images.MonoData)

  if ! config.ConfigFileExists() {
    e1 = systray.AddMenuItem("No config file","",0)
  }

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
        fmt.Printf("%v+\n", time.Now())
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
    exit := systray.AddMenuItem("Quit", "", 0)

    go func() {
      for {
        select {
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

func onExit() {

}

func forever() {
  for {
    fmt.Printf("%v+\n", time.Now())
    time.Sleep(time.Second)
  }
}
