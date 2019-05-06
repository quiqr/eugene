
task :default => :build

desc "build hugo control"
task :build do
  system("go build")
  system("kill -9 `ps ax | grep hugo-control | grep -v grep | cut -d \" \" -f1`")
  system("cp ./hugo-control ./dist/Hugo\\ Control.app/Contents/MacOS/")
  system("open ./dist/Hugo\\ Control.app")
end


