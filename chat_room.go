package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Clinets struct {
	name string
	addr string
	ch   chan string
}
var enteing =make(chan Clinets,0)
var leaving =make(chan Clinets,0)
var masving =make(chan string,0)
var list    =make(chan Clinets,0)
var rename  =make(chan Clinets,0)

func main()  {


	linner,err:=net.Listen("tcp","127.0.0.1:8881")
	  if err!=nil{
	  	fmt.Println("Listen is err",err)
		  return
	  }
	//defer linner.Close()


	go bordcaster()


	for  {
		conn,err:=linner.Accept()
		if err!=nil{
			fmt.Println("Accept is err",err)
			continue
		}
		//defer conn.Close()
	
	go handler(conn)
	}


}

func bordcaster()  {
	fmt.Println("bordcaster is running")
	Clinet:=make(map[string]Clinets)
	for  {
		select {
		case cli := <-enteing:
			Clinet[cli.addr]=cli
		case cli := <-leaving:
			delete(Clinet,cli.addr)
			close(cli.ch)
		case msg := <- masving:
			for _, cli := range Clinet{
				cli.ch <- msg
			}
		case cli := <- list:
			for _,v:=range Clinet{
				cli.ch<-v.name
			}
		case cli := <- rename:
			Clinet[cli.addr]=cli
		}
	}

}

func handler(conn net.Conn)  {
	fmt.Println(conn.RemoteAddr().String(),"net.conn is running")

	name:=conn.RemoteAddr().String()
	addr:=conn.RemoteAddr().String()
	ch  :=make(chan string)
	cli:=Clinets{name:name,addr:addr,ch:ch}
	go CilinetWrite(conn,ch)
	cli.ch<-"["+cli.name+"]: welcome to you "
    masving<-"["+cli.name+"]: age wo"
	enteing<-cli
	var hasData=make(chan struct{})
	//超时强踢功能
	go func() {
		select {
		   case <-hasData:
		   case <-time.After(60*time.Second):
		   	conn.Close()

		}
	}()
	input:=bufio.NewScanner(conn)
	for input.Scan(){
		hasData<- struct{}{}
        inputData:=input.Text()
		if inputData=="#list"{
			list<-cli
		}else if len(inputData) >8 && inputData[0:7]=="#rename"{
            cli.name=inputData[8:]
            rename<-cli
            cli.ch<-"rename is running"
		}else{
			msg:="["+cli.name+"]:"+input.Text()
			masving<-msg
		}
	}

	leaving<-cli
	masving<-"["+cli.name+"]: time out"
}

func CilinetWrite(conn net.Conn,ch chan string)  {
	for msg:=range ch{
		fmt.Fprintln(conn,msg)
	}
}