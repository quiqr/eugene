package hugo

import (
  "hugo-control/config"
  "os"
  "os/exec"
  "fmt"
  "path"
  "path/filepath"
)

var hugoServer = exec.Command(HugoBinPath(), "server","-D", "-s", HugoDir() )

func HugoBinExists() bool {

  hugofile := path.Clean(filepath.Join(config.Dir(), "hugo"))

  if stat, err := os.Stat(hugofile); err == nil && stat.Mode().IsRegular() {
    return true
  } else {
    return false
  }
}

func HugoDir() string {
  cfg,_ := config.Read()
  return cfg.HugoDir
}

func HugoBinPath() string {
  return path.Clean(filepath.Join(config.Dir(), "hugo"))
}

func StartHugo() {
  hugoServer = exec.Command(HugoBinPath(), "server","-D", "-s", HugoDir())
  hugoServer.Start()
}

func KillHugo()  {

  if err := hugoServer.Process.Kill(); err != nil {
    fmt.Println("failed to kill process: ")
  }
}

func HugoPid() string {
  out, _ := exec.Command("bash", "-c", "/bin/ps ax | /usr/bin/grep \"bitbar-hugo\\/hugo\"| /usr/bin/grep -v grep | /usr/bin/head -n1 | /usr/bin/cut -d\" \" -f 1").Output()

  pid := string(out)
  if fmt.Sprintf("%s", pid) != "<nil>" {
    return pid
  } else {
    return ""
  }
}

func HugoRunning() bool {
  if HugoPid() != "" {
    return true
  } else {
    return false
  }
}




