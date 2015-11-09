thrift-0.9.3.exe  -gen go *.thrift
xcopy /s /e /i /y .\gen-go\TestGo\Rpc\* .\ &&rd /s /q .\gen-go\