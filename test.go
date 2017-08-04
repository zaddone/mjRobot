package main
import(
	"fmt"
	"math/rand"
	"time"
	"./robot"
)
func main(){
	Ro := new(robot.Robot)
	Ro.Inv = 2
	i := 0
	for {
		var now [3][9]int
		ra := rand.New(rand.NewSource(time.Now().Unix()))
		for i:=0;i<14;{
			h := 0//ra.Intn(2)
			n := ra.Intn(9)

			if now[h][n] > 2 {
				continue
			}
			now[h][n]+=1
			i+=1
		}
		Ro.Now = now
		k:=Ro.CheckFull(true)
		i++
		fmt.Printf("%d\r",i)

//		if k >0 {
			fmt.Println(now,k,i)
			var cmd string
			fmt.Scanf("%s",&cmd)
			i = 0
//		}
	}

}
