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

	//第一步创建主go程监听clinet循环连接
	linner,err:=net.Listen("tcp","127.0.0.1:8881")
	  if err!=nil{
	  	fmt.Println("Listen is err",err)
		  return
	  }
	//defer linner.Close()

	//启动广播bordcaster
	go bordcaster()

	//第二步循环创建套接字连接
	for  {
		conn,err:=linner.Accept()
		if err!=nil{
			fmt.Println("Accept is err",err)
			continue
		}
		//defer conn.Close()
	//启动handler
	go handler(conn)
	}


}

func bordcaster()  {
	fmt.Println("bordcaster is running")
	Clinet:=make(map[string]Clinets)                   //Clinets结构体类型的map
	for  {                                             //持续监听
		select {
		case cli := <-enteing:                         //将用户进来的信息写入Clinet结构体
			Clinet[cli.addr]=cli
		case cli := <-leaving:                         //将用户离开的信息删除，关闭通道
			delete(Clinet,cli.addr)
			close(cli.ch)
		case msg := <- masving:                        //将用户发送的信息写入ch通道里  然后广播消息
			for _, cli := range Clinet{                //循环所有用户
				cli.ch <- msg                          //把信息写入所有用户的通道里
			}
		case cli := <- list:                           //当前用户
			for _,v:=range Clinet{                     //扫描所有用户信息
				cli.ch<-v.name                         //把所有用户信息写入当前拉取用户的信息里
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
	cli:=Clinets{name:name,addr:addr,ch:ch}                            //用户进来就记录信息
	go CilinetWrite(conn,ch)                                           //写入信息到clinet
	cli.ch<-"["+cli.name+"]: welcome to you "                          //通知自己   cli代表是自己
    masving<-"["+cli.name+"]: age wo"                                  //广播通知
	enteing<-cli                                                       //把信息写入通道
	var hasData=make(chan struct{})
	//超时强踢功能
	go func() {
		select {
		   case <-hasData:                                             //如果读到值就刷新  重新计时
		   case <-time.After(60*time.Second):
		   	conn.Close()

		}
	}()
	input:=bufio.NewScanner(conn)                                      //读取信息的扫描仪
	for input.Scan(){                                                  //持续扫描               一直阻塞扫描  断开才会跳出方法往下走
		hasData<- struct{}{}                                          //超时强踢   如果用户有输入就赋值
        inputData:=input.Text()
		if inputData=="#list"{                                       //如果内容等于#list命令的话就把客户端信息写入list     #list拉取当前所有在线用户
			list<-cli
		}else if len(inputData) >8 && inputData[0:7]=="#rename"{       //更改名字
            cli.name=inputData[8:]
            rename<-cli
            cli.ch<-"rename is running"
		}else{
			msg:="["+cli.name+"]:"+input.Text()                            //扫描到内容
			masving<-msg                                                   //写入广播通道
		}
	}

	leaving<-cli
	masving<-"["+cli.name+"]: time out"                                    //离开时通知
}

func CilinetWrite(conn net.Conn,ch chan string)  {
	for msg:=range ch{                                                       //循环通道写入客户端
		fmt.Fprintln(conn,msg)
	}
}