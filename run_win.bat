set GOPATH=%CD%

call go install github.com/nightdeveloper/smartcenter/main

call "bin/main.exe"
