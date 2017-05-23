package main
import (
	"fmt"
	"./user"
	"./grain"
	"os"
	"math/rand"
)
var (
	ML *grain.MJList
	Users []*user.MJUser
	Rule  *user.MJRule
	Public *user.UserPublic
)
func InitUserAll () {
	ML = new(grain.MJList)
	Public = new(user.UserPublic)
//	DeskTop = new(user.Desktop)
	ML.Init()
	Rule = new(user.MJRule)
	Rule.Init()
	Users = make([]*user.MJUser,4)
	begin:=0
	for i,_ := range Users {
		Users[i]= new(user.MJUser)
		begin = ML.Count - user.UN
		if i==0 {
			begin  -=1
		}
		n :=0
		for j := begin;j<ML.Count;j++ {
			ml :=ML.List[j]
			Users[i].Now[ml.H][ml.N]++
			n++
		}
		Users[i].Init(ML,Rule,Public,i)
		ML.Count = begin
	}
	Users[rand.Intn(4)].Self = true

}
func NumAdd(i int) (j int) {
	j = i + 1
	if j == 4 {
		return 0
	}
	if Users[j].Over {
		j = NumAdd(j)
		if i == j {
			os.Exit(i)
		}
		return j
	}
	return j
}
func Quan(gr *grain.MJGrain,b,I int) int {

//	fmt.Println(I,Users[I].Now,Users[I].Discard)
//	fmt.Println(b,I,gr)
	if b == I {
//		fmt.Println("quan over")
//		DeskTop.See[b][gr.H][gr.N] ++
		I = NumAdd(I)
		if Users[I].In() {
			fmt.Println("self",b)
			return Quan(nil,I,I)
		}
		g := Users[I].Outs()
		if g ==nil {
			fmt.Println("self",b)
			return Quan(nil,I,I)
		}
		return Quan(g,I,NumAdd(I))
	}
	g := Users[I].SeeOut(gr)
//	var cmd string
//	fmt.Scanf("%s",&cmd)
	if g == nil {
		fmt.Println(b,I)
		return Quan(gr,b,NumAdd(I))
	}
//	DeskTop.See[b][gr.H][gr.N] ++
	gr.O = 1
	Users[b].DelSee(gr)

	return Quan(g,I,NumAdd(I))

}
func main(){
	InitUserAll()
	b:=0
	g := Users[b].Outs()
	if g ==  nil {
		return
	}
//	fmt.Println(g)
//	g := ML.Out()
	b = Quan(g,b,b+1)
	b++

}
