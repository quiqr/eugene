package config

/*
type Config struct {
  HugoDir string
  SiteName string
  LiveUrl  string
}
*/

type ConfigMulti struct {
    Current_Site int
    Sites []Site
}

type Site struct {
  Name string
  Hugo_Src_Dir string
  Hugo_Output_Dir string
  Live_Url string
  Publishing_Command string
}
