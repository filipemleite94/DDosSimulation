package main

import(
	"fmt"
	"time"
    "encoding/gob"
    "bytes"
    "net"
)

type Node struct{
    ID int
    T time.Time
    IP string
    Port int
}

type Message struct{
    N Node
	Command string
}

var ServConn *net.UDPConn

func checkError(err error) bool{
	if err!=nil{
		fmt.Println("Erro: ", err)
        return true
	}
    return false
}

func printNode(str string, node Node){
    fmt.Println(str, node.ID, node.IP, node.Port, node.T)
}

func main(){
    ServerAddr,servErr := net.ResolveUDPAddr("udp", "127.0.0.1:50000")
    checkError(servErr)
    ServConn, servErr = net.ListenUDP("udp", ServerAddr)
    checkError(servErr)
    //fmt.Println("Listening in udp:", "127.0.0.1:50000")
    recBuffer := make([]byte, 1024)
    msg := Message{}
    for{
        n,addr,servErr := ServConn.ReadFromUDP(recBuffer)
        checkError(servErr)
        servErr = gob.NewDecoder(bytes.NewReader(recBuffer[:n])).Decode(&msg)
        checkError(servErr)
        t1:=time.Now().UnixNano()
        t2:=msg.N.T.UnixNano()
        fmt.Println(t1,t2,t1-t2,addr,msg.N.ID)
    }
}
