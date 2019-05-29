
task :default => :build

desc "build Eugene"
task :build do
  system("go build")
  system("kill -9 `ps ax | grep eugene | grep -v grep | cut -d \" \" -f1`")
  system("rm -Rf ./dist/Eugene.app ")
  system("cp -a ./macos/Eugene.app dist/Eugene.app")
  system("cp -a ./eugene ./dist/Eugene.app/Contents/MacOS/")
  system("open ./dist/Eugene.app")
end


