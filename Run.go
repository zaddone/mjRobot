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
	u := GetAIUser(str[0],(nc==5 || nc == 6))
	if nid < 3 {
		n,err := strconv.Atoi(str[3])
		if err != nil {
			panic(err)
		}
		lsn := n<<8 + nid
		if u.LastSee == lsn {
//			fmt.Println("is have",nid,n)
			return []byte{1}
		}
		fmt.Println("k->",nid,n)
		u.Public.See[nid][n/9][n%9]++
		Test[nid][n/9][n%9]++
		u.LastSee = lsn
//		u.NowU ++
	//	if nid == 0 {
	//		u.NowU = 1
	//	}else{
	//	}
		return []byte{1}
	}else{
		if nid == 18 {
			fmt.Println(str)
			u.BaseVal = nc
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
		if nid == 3 {
			u.UpdateNow(str[4:])
		}else if nid  == 4 {
			u.InitNow(str[4:])
		}
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
			n :=u.SeeOut(u.LastSee>>8,(nowVal - u.NowVal )/u.BaseVal)
			sends = fmt.Sprintf("%d %d",nid,n)
		}else if nc == 8 {
			n :=u.Outs((nowVal - u.NowVal)/u.BaseVal)
			Test[3][n/9][n%9]++
			sends = fmt.Sprintf("%d %d",nid,n)

		}else if nc == 9 {
			n := u.SeeSelf((nowVal - u.NowVal)/u.BaseVal)
		//	g :=u.CheckGang()
		//	u.Public.Down[u.Uid][g.H][g.N]+=u.Now[g.H][g.N]
			sends = fmt.Sprintf("%d %d",nc,n)
		}
		fmt.Println("sends:",sends)
//		u.NowU = 0
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
			for i,t := range Test {
				fmt.Println(i,t)
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
