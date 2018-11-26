package main

import(
	"fmt"
	"os"
	"time"
    "io/ioutil"
    "encoding/gob"
    "bytes"
    "net"
    "strings"
    "strconv"
    "bufio"
)


//Structs, global variables and general functions--------------------
//-------------------------------------------------------------------
type Node struct{
    ID int
    T time.Time
    IP string
    Port int
}

type State struct{
    LastID int
    MapKeys map[int]Node
}

type Message struct{
    N Node
	Command string
}

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

var state State
var ServConn *net.UDPConn
var lastConnectedID int
var firstConnectedID int

//Deal with user input-----------------------------------------------
//-------------------------------------------------------------------
var ch chan string
	
func readInput() {
	// Non-blocking async routine to listen for terminal input
	reader:=bufio.NewReader(os.Stdin)
	for {
		text, _, _ := reader.ReadLine()
		ch<-string(text)
	}
}

func dealWithUserInput(){
	for{
		select {
			case x, valid := <-ch:
				if valid {
                    if x == "exit" {
                        writeState()
                        os.Exit(0)
                    }else{
                        if x=="attack1" {
                            msg := Message{Node{0, time.Now().UTC().Add(10000*time.Millisecond), "127.0.0.1", 50000}, "attack1"}
                            sendMsg(msg, connect(state.MapKeys[firstConnectedID].IP, state.MapKeys[firstConnectedID].Port))
                        }else if x=="attack2" {
                            msg := Message{Node{0, time.Now().UTC().Add(10000*time.Millisecond), "127.0.0.1", 50000}, "attack2"}
                            sendMsg(msg, connect(state.MapKeys[firstConnectedID].IP, state.MapKeys[firstConnectedID].Port))
                        }
                    }
				}else{
					fmt.Println("Channel closed!")
				}
			default:
				time.Sleep(time.Second*1)
		}
	}
}

//State serializer part----------------------------------------------
//-------------------------------------------------------------------
func printState(){
    fmt.Println("Printing Map: ")
    fmt.Println(state.LastID)
    for key,value := range state.MapKeys {
        printNode(strconv.Itoa(key)+":", value)
    }
}

func initState(){
    lastConnectedID = 0
    ch = make(chan string)
    if _, err := os.Stat("ServerStatus"); os.IsNotExist(err) {
        state = State{0, make(map[int]Node)}
        fmt.Println("Does not exist file yet")
        printState()
    }else{
        byteArray,err := ioutil.ReadFile("ServerStatus")
        checkError(err)
        err=gob.NewDecoder(bytes.NewReader(byteArray)).Decode(&state)
        checkError(err)
        printState()
    }
}

func writeState(){
    var f *os.File
    var err error
    f,err=os.OpenFile("ServerStatus", os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0)
    var buf bytes.Buffer
    e:= gob.NewEncoder(&buf)
    err=e.Encode(state)
    checkError(err)
    _,err=f.Write(buf.Bytes())
    checkError(err)
    f.Close()
}

//Web part-----------------------------------------------------------
//-------------------------------------------------------------------
func doServerJob(conn *net.UDPConn) (Message,net.Addr){
    recBuffer := make([]byte, 1024)
	n,addr,servErr := conn.ReadFromUDP(recBuffer)
	checkError(servErr)
	msg := Message{}
    servErr = gob.NewDecoder(bytes.NewReader(recBuffer[:n])).Decode(&msg)
    checkError(servErr)
    printNode("Receiving Message "+addr.String()+" "+msg.Command, msg.N)
    return msg,addr
}

func getByteArray(msg Message)[]byte{
    var buf bytes.Buffer
    e:= gob.NewEncoder(&buf)
    var err error
    err = e.Encode(msg)
    checkError(err)
    printNode("Encoding Message "+msg.Command, msg.N)
    return buf.Bytes()
}

func sendMsg(msg Message, conn *net.UDPConn){
    fmt.Println("Sending Message To "+conn.RemoteAddr().String())
    _,err := conn.Write(getByteArray(msg))
    checkError(err)
}

func sendMsgTo(msg Message, conn *net.UDPConn, addr net.Addr){
    fmt.Println("Sending Message To "+addr.String())
    _,err := conn.WriteTo(getByteArray(msg),addr)
    checkError(err)
}

func initServer(){
	ServerAddr,servErr := net.ResolveUDPAddr("udp", "127.0.0.1:10000")
    checkError(servErr)
    /* Now listen at selected port */
    ServConn, servErr = net.ListenUDP("udp", ServerAddr)
    checkError(servErr)
    fmt.Println("Listening in udp:", "127.0.0.1:10000")
}

func runServer(){
    for{
        msg,addr:=doServerJob(ServConn)
        takeAction(msg, ServConn, addr)
    }
}

func connect(ip string, port int)*net.UDPConn{
    remoteAddr,clientErr := net.ResolveUDPAddr("udp",ip+":"+strconv.Itoa(port))
    checkError(clientErr)
    LocalAddr, clientErr := net.ResolveUDPAddr("udp", "127.0.0.1:0")
    checkError(clientErr)
    conn, clientErr := net.DialUDP("udp", LocalAddr, remoteAddr)
    checkError(clientErr)
    return conn
}

//Mensage interpreter part-------------------------------------------
//-------------------------------------------------------------------
func takeAction(msg Message, conn *net.UDPConn, addr net.Addr){
    id:=msg.N.ID
    if(msg.Command=="Join"){
        strs:=strings.Split(addr.String(), ":")
        port:=msg.N.Port
        if(id==0){
            state.LastID++
            id = state.LastID
            state.MapKeys[id] = Node{id,time.Now().UTC(), strs[0], port}
            msg.N = state.MapKeys[id]
            msg.Command = "Welcome"
            sendMsgTo(msg, conn, addr)
            printNode("Just joined:", msg.N)
        }else{
            state.MapKeys[id] = Node{id,state.MapKeys[id].T, strs[0], port}
            printNode("Joined back:", msg.N)
        }
        if(lastConnectedID>0){
            msg.N = state.MapKeys[lastConnectedID]
        }else{
            msg.N = state.MapKeys[id]
            firstConnectedID = id
        }
        printNode("With Parent:", msg.N)
        msg.Command = "Node, this is your father"
        sendMsgTo(msg, conn, addr)
        lastConnectedID = id
    }else{
        //Return time for election
        msg.N = state.MapKeys[id]
        sendMsgTo(msg, conn, addr)
    }
}

func main(){
    fmt.Println("Initializing state")
    initState()
    fmt.Println("Initializing server")
    initServer()
    go readInput()
    go dealWithUserInput()
    runServer()
}
