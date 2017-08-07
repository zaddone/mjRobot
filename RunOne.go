package main
import (
	"fmt"
	"net"
	"flag"
	"runtime"
//	"time"
	"strings"
	"strconv"
	"./robot"
)
var (
	Port = flag.String("p",":3333","Port")
	Conn = flag.Int("c",4,"conn")
	AIUser map[string]*robot.Robot
	Test [4][3][9]int
)
func GetAIUser(key string,isUpdate bool) (u *robot.Robot) {
	u = AIUser[key]

	if u == nil || isUpdate  {
		u = new(robot.Robot)
		u.Init(new(robot.UserPublic),3)
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
	if u.LastForData[0] == data {
		return []byte(u.LastForData[1])
	}
	fmt.Println(data)
	u.LastForData[0] = data
	u.LastForData[1] = string([]byte{1})
	if nid < 3 {
		if nc >0 {
			if nc == 1 {
				fmt.Println("fc->",nid,nc)
				if u.NowU[nid] >=0 {
					u.NowU[nid] = - (u.NowU[nid]+1)
				}
//				u.Public.Down[nid][n/9][n%9]= nc
			}else{
				n,err := strconv.Atoi(str[3])
				if err != nil {
					panic(err)
				}
				fmt.Println("nc->",nid,n,nc)
				if nc  == 100 {
					u.NowU[nid] = n
				}else{
					u.Public.Down[nid][n/9][n%9]= nc

				}
			}
		}else{
			n,err := strconv.Atoi(str[3])
			if err != nil {
				panic(err)
			}
			if u.NowU[nid] < 0 {
				u.NowU[nid] = (- u.NowU[nid]) -1
			}
	//		if u.LastSee == lsn {
//	//			fmt.Println("is have",nid,n)
	//			return []byte{1}
	//		}
			fmt.Println("k->",nid,n)
			u.Public.See[nid][n/9][n%9]++
		//	Test[nid][n/9][n%9]++
			u.LastSee = n
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
			u.BaseVal = nc
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
//		baseVal,err := strconv.Atoi(strTmp[1])
//		if err != nil {
//			panic(err)
//		}
		fmt.Println("nid",nid,nc,nowVal,u.NowVal,u.BaseVal)
		fmt.Println(str[4:])
		u.InitNow(str[4:])
//		if nid == 3 {
//			u.UpdateNow(str[4:])
//		}else if nid  == 4 {
//			u.InitNow(str[4:])
//		}
	//	u.InitNow(str[3:])
		var sends string
		if nc == 5 {
			outs := u.OutDiscard(3)
			sends = str[1]
			for _,o := range outs {
				sends = fmt.Sprintf("%s %d",sends,o)
			}

			u.NowVal = nowVal
		}else if nc == 6 {
			a := u.GetDiscard(true)
			u.SetValid(a)
			sends = fmt.Sprintf("6 %d",a)
			u.NowVal = nowVal
		}else if nc == 7  {
			n :=u.SeeOut(u.LastSee,(nowVal - u.NowVal )/u.BaseVal)
			sends = fmt.Sprintf("%d %d",nid,n)
		}else if nc == 8 {
			n :=u.Outs((nowVal - u.NowVal)/u.BaseVal,true)
			if n >= 0 && n <27  {
				u.Public.See[u.Uid][n/9][n%9] ++
			}
//			Test[3][n/9][n%9]++
			sends = fmt.Sprintf("%d %d",nid,n)

		}else if nc == 9 {
			n :=u.Outs((nowVal - u.NowVal)/u.BaseVal,false)
//			n := u.SeeSelf((nowVal - u.NowVal)/u.BaseVal)
		//	g :=u.CheckGang()
		//	u.Public.Down[u.Uid][g.H][g.N]+=u.Now[g.H][g.N]
			sends = fmt.Sprintf("%d %d",nc,n)
		}
		fmt.Println("sends:",sends)
//		u.NowU = 0
		u.LastForData[1] = sends

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
			for t,u := range AIUser {
				fmt.Println(t)
				fmt.Println("down")
				for i,d := range u.Public.Down{
					fmt.Println(i,"d",d)
			//		fmt.Println(i,"s",u.Public.See[i])
				}
			}
		}
		fmt.Println(cmd)
	}

}

func main(){
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	AIUser = make(map[string]*robot.Robot)
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
