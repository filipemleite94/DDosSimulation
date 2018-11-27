go build Server.go
go build HostileServer.go
go build Client.go
g++ CountFrequency.cpp -o CountFrequency.exe
Taskkill /IM Client.exe /F
start RunForrestAux.bat
start Server.exe L
timeout 2
FOR /L %%G IN (1,1,100) DO start /B Client.exe %%G 4 100
timeout 25
Taskkill /IM Client.exe /F
echo FIM
