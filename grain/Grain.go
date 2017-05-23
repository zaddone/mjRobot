package grain
import (
	"math/rand"
	"time"
)
const (
	Con int = 108
	Ho int = 3
	No int = 9
)
type MJGrain struct {
	H byte // 0 1 2
	N byte //0--8

	O int
}
type MJList struct {
	List [Con]*MJGrain //108
	Count int
}
func (self *MJList) Out() *MJGrain {
	self.Count --
	if self.Count<0 {
		return nil
	}
	return self.List[self.Count]
}
func (self *MJList) Init(){
	self.Count = Con
//	self.List = make([]*MJGrain,Co)
	ra := rand.New(rand.NewSource(time.Now().Unix()))
	var MJ [Ho][No]int
	for i:=0;i<Con;{
		h:= byte(ra.Intn(Ho))
		n:= byte(ra.Intn(No))
		if MJ[h][n] == 4 {
			continue
		}
		MJ[h][n] ++
//		fmt.Println(h,n,MJ[h][n])
		self.List[i]=&MJGrain{H:h,N:n}
		i++

	}

}
