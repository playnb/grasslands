@echo off
::for /r %%d in (.) do cd %%d && cd && cd %cd%
::for /r %%d in (.) do @echo %%d

::for /r %%d in (.) do cd %%d && protoc --go_out=plugins=protorpc:. *.proto

SET DIR1= %GOPATH%/src/github.com/playnb/mustang/msg.pb/
CD %DIR1%
protoc --go_out=plugins=protorpc:. *.proto

pause
