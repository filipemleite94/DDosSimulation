go build Server.go
go build HostileServer.go
go build Client.go
g++ CountFrequency.cpp -o CountFrequency.exe

start RunForrestAux.bat
start Server.exe 
timeout 2
FOR /L %%G IN (1,1,100) DO start /B Client.exe %%G
timeout 25
Taskkill /IM Client.exe /F
echo FIM
