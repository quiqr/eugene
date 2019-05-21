package config

type ConfigMulti struct {
    Current_Site int
    Sites []Site
}

type Site struct {
  Name string
  Hugo_Src_Dir string
  Live_Hugo_Output_Dir string
  Live_Url string
  Live_Publishing_Command string
}
