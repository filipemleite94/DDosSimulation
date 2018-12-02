@echo off
rem recompile all the iportant classes
@echo on
go build Server.go
go build HostileServer.go
go build Client.go
g++ CountFrequency.cpp -o CountFrequency.exe
@echo off
rem Kill clients that were not opened
@echo on
Taskkill /IM Client.exe /F
@echo off
rem Start the target server
@echo on
start RunForrestAux.bat
@echo off
rem Start the botnet server use 'L' for list and 'T' for tree
@echo on
start Server.exe L
timeout 2
@echo off
rem Start the client process in background, the variables there are its ID,
rem number of packages sent per milisecond and the average delay before 
rem transmitting the information for the next layer
@echo on
FOR /L %%G IN (1,1,100) DO start /B Client.exe %%G 4 100
@echo off
rem Wait 25 seconds and then kill all the Client process running in the background
@echo on
timeout 25
Taskkill /IM Client.exe /F
echo FIM
@echo off
rem Pay atention that if no attack is given to the botnet server's 
rem console in less than the previous time, no attack will happen.
rem To do the attack simple write attack1 or attack2 in the console.
rem To stop the botnet server, just write down exit in its console
rem To stop the target server use Ctrl+C
@echo on