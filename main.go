package main
import (
	"fmt"
	"net"
	"flag"
	"runtime"
//	"time"
	"strings"
	"strconv"
//	"./robot"
	"./ai"
)
type Robot interface{
	Init()
	SetLastSee(nid int,n int)
	GetLastSee() int
	Outs(isf int,No30 bool) int
	OutDiscard(num int) (o []int)
	GetDiscard() int
	SeeOuts(v int,key []int,isf int) int
	SetLastForData(n int,s string)
	GetLastForData() [2]string
	SetNowVal(n int)
	SetValid(nid int,n int)
	SetPublicDown(nid int,n int,val int)
	SetPublicSee(nid int,n int,val int)
	InitNow(str []string)
	GetFullVal(n int) int
	SetBaseVal(n int)
//	SetLastUid(uid int)
}
var (
	Port = flag.String("p",":3334","Port")
	Conn = flag.Int("c",4,"conn")
	AIUser map[string]Robot
)
func GetAIUser(key string,isUpdate bool) (u Robot) {
	u = AIUser[key]

	if u == nil || isUpdate  {
		u = new(ai.Ai)
		u.Init()
		AIUser[key] = u

	}
	return u
}
func DataHandle(data string) []byte {
	if len(data) == 1 {
		return []byte{1}
	}
	str := strings.Split(data," ")
	nid,err := strconv.Atoi(str[1])
	if err != nil {
		panic(err)
	}
	nc,err := strconv.Atoi(str[2])
	if err != nil {
		panic(err)
	}
	u := GetAIUser(str[0],(nid <90 &&(nc==5 || nc == 6)))
	la := u.GetLastForData()
	if la[0] == data {
		return []byte(la[1])
	}
	fmt.Println(data)
	u.SetLastForData(0,data)
	u.SetLastForData(1,string([]byte{1}))
	if nid < 3 {
		if nc >0 {
			if nc == 1 {
		//		fmt.Println("fc->",nid,nc)
		//		if u.NowU[nid] >=0 {
		//			u.NowU[nid] = - (u.NowU[nid]+1)
		//		}
//				u.Public.Down[nid][n/9][n%9]= nc
			}else{
				n,err := strconv.Atoi(str[3])
				if err != nil {
					panic(err)
				}
				fmt.Println("nc->",nid,n,nc)
				if nc  == 100 {
			//		u.NowU[nid] = n
				}else{
					u.SetPublicDown(nid,n,nc)
					//u.Public.Down[nid][n/9][n%9]= nc

				}
			}
		}else{
			n,err := strconv.Atoi(str[3])
			if err != nil {
				panic(err)
			}
		//	if u.NowU[nid] < 0 {
		//		u.NowU[nid] = (- u.NowU[nid]) -1
		//	}
	//		if u.LastSee == lsn {
//	//			fmt.Println("is have",nid,n)
	//			return []byte{1}
	//		}
			fmt.Println("k->",nid,n)
//			u.Public.See[nid][n/9][n%9]++
			u.SetPublicDown(nid,n,1)
		//	Test[nid][n/9][n%9]++
			u.SetLastSee(nid,n)
		}
//		u.NowU ++
	//	if nid == 0 {
	//		u.NowU = 1
	//	}else{
	//	}
		return []byte{1}
	}else{
		if nid == 90 {
			fmt.Println("90=====",str)
			u.SetBaseVal(nc)
			return []byte{1}
		}else if nid == 91 {
			fmt.Println("91=====",str)
		//	u.BaseVal = nc
			return []byte{1}
		}
//		strTmp := strings.Split(str[3],"_")

		nowVal,err := strconv.Atoi(str[3])
		if err != nil {
			panic(err)
		}
		fmt.Println(str[4:])
		var sends string
		if nc == 7 {
			var ns []int
			for _,s := range str[4:] {
				n,err := strconv.Atoi(s)
				if err != nil {
					panic(err)
				}
				ns = append(ns,n)
			}
			n := u.SeeOuts(u.GetLastSee(),ns,u.GetFullVal(nowVal))
			sends = fmt.Sprintf("%d %d",nid,n)
		}else{
			u.InitNow(str[4:])
			if nc == 4 {
				return []byte{1}
			}
			if nc == 5 {
				outs := u.OutDiscard(3)
				sends = str[1]
				for _,o := range outs {
					sends = fmt.Sprintf("%s %d",sends,o)
				}
			}else if nc == 6 {
				a := u.GetDiscard()
				u.SetValid(3,a)
				sends = fmt.Sprintf("6 %d",a)
				u.SetNowVal(nowVal)
//			}else if nc == 7  {
//				n :=u.SeeOut(u.GetLastSee(),u.GetFullVal(nowVal))
//				sends = fmt.Sprintf("%d %d",nid,n)
			}else if nc == 8 {
				n :=u.Outs(u.GetFullVal(nowVal),true)
				sends = fmt.Sprintf("%d %d",nid,n)

			}else if nc == 9 {
				n :=u.Outs(u.GetFullVal(nowVal),false)
				sends = fmt.Sprintf("%d %d",nc,n)
			}
		}
		fmt.Println("sends:",sends)
		u.SetLastForData(1,sends)

		return []byte(sends)
	}
	return []byte{1}
}
func ReadData(conn net.Conn) error {
	var db [1024]byte
	n,err := conn.Read(db[0:])
	if err != nil {
		return err
	}
	_,err = conn.Write(DataHandle(string(db[:n])))
	return err
}
func syncHandle(conn net.Conn){
	for {
		err := ReadData(conn)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	defer conn.Close()

}
func console(){
	var cmd string
	for{
		fmt.Scanf("%s",&cmd)
		if cmd == "test" {
	//		for t,u := range AIUser {
	//			fmt.Println(t)
	//			fmt.Println("down")
		//		for i,d := range u.Public.Down{
		//			fmt.Println(i,"d",d)
		//	//		fmt.Println(i,"s",u.Public.See[i])
		//		}
	//		}
		}
		fmt.Println(cmd)
	}

}

func main(){

//	var now [3][9]int
//	tmp := now[1][:]
//	fmt.Println(now,tmp)
//	now[1][0] = 9
//	fmt.Println(now,tmp)

	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	AIUser = make(map[string]Robot)
	go console()
	listener,err := net.Listen("tcp",*Port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	conn_chan := make(chan net.Conn)
	for i:=0;i<*Conn;i++{
		go func(){
			for conn:= range conn_chan{
				syncHandle(conn)
			}
		}()
	}
	for {
		conn,err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
//		conn.SetDeadline(time.Now().Add(30*time.Second))
		conn_chan <- conn
	}

}
