package main
import (
	"fmt"
	"net"
	"flag"
	"runtime"
//	"time"
	"strings"
	"strconv"
	"./user"
)
var (
	Port = flag.String("p",":3333","Port")
	Conn = flag.Int("c",4,"conn")
	AIUser map[string]*user.MJAI
	Rule *user.MJRule
)
func GetAIUser(key string,isUpdate bool) (u *user.MJAI) {
	u = AIUser[key]

	if u == nil || isUpdate  {
		u = new(user.MJAI)
		u.Init(Rule,new(user.UserPublic),3)
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
	fmt.Println(data)
	if nc ==5 || nc == 50 || nc == 6 || nc == 60 {
		for _,s := range str[3:] {
			n,err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			i:=n/9
			j:=n%9
			u.Now[i][j]++
		}
		if (len(str[3:])==1){
			return []byte{1}
		}
		if nc == 5 || nc == 50 {
			outs := u.OutDiscard(3)
			sends := str[1]
			for _,o := range outs {
				sends = fmt.Sprintf("%s %d",sends,o)
			}
			fmt.Println("sends:",sends)
			return []byte(sends)

		}else if nc == 6  || nc == 60 {
	//		if (len(str[3:])==1)return []byte{1}
			u.Discard = -1
			ans := u.ReadNow(u.Now)
			a := u.GetDiscard(ans)
			sends := fmt.Sprintf("6 %d",u.Discard)
			fmt.Println("sends:",sends,a)
			return []byte(sends)
		}
	}
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
		u.Public.See[nid][n/9][n%9]++
		u.LastSee = lsn
		return []byte{1}
	}else{
		if nid == 3 {
			n,err := strconv.Atoi(str[3])
			if err != nil {
				panic(err)
			}
			u.Now[n/9][n%9]++
		}else if nid  == 4 {
			var now [3][9]byte
			for _,s := range str[3:] {
				n,err := strconv.Atoi(s)
				if err != nil {
					panic(err)
				}
				i:=n/9
				j:=n%9
				now[i][j]++
			}
			u.Now = now
		}
		if nc == 7  {
			n :=u.SeeOut(u.LastSee>>8)
			sends := fmt.Sprintf("4 %d",n)
			fmt.Println("sends:",sends)
			return []byte(sends)
		}else if nc == 8 {
			n :=u.Outs()
			sends := fmt.Sprintf("3 %d",n)
			fmt.Println("sends:",sends)
			return []byte(sends)

		}
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
		fmt.Println(cmd)
	}

}

func main(){
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	Rule = new(user.MJRule)
	Rule.Init()
	AIUser = make(map[string]*user.MJAI)
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
