set GOPATH=%CD%

del bin/main.exe

call go install github.com/nightdeveloper/smartcenter/main

call "bin/main.exe"
