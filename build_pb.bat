@echo off
::for /r %%d in (.) do cd %%d && cd && cd %cd%
::for /r %%d in (.) do @echo %%d

::for /r %%d in (.) do cd %%d && protoc --go_out=plugins=protorpc:. *.proto

protoc --go_out=plugins=protorpc:. Msg/*.proto
protoc --go_out=plugins=protorpc:. msg.pb/*.proto

pause
