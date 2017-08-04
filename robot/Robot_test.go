package robot
import(
	"testing"
	"fmt"
	"math/rand"
	"time"
)
func Test_Robot(t *testing.T){
	Ro := new(Robot)
	Ro.Inv = 2
	for{
		var now [3][9]int
		ra := rand.New(rand.NewSource(time.Now().Unix()))
		for i:=0;i<15;{
			h := ra.Intn(2)
			n := ra.Intn(9)
			if now[h][n] == 4 {
				continue
			}
			now[h][n]++
			i++
		}
		Ro.Now = now
		k := Ro.CheckFull(true)
		fmt.Println(now,k)

		var cmd string
		fmt.Scanf("%s",&cmd)
	}
}
