package user
import (
	"../grain"
	"fmt"
	"os"
//	"math"
)
type MJUser struct {

	Now [grain.Ho][grain.No]byte
	Public *UserPublic
	Down [grain.Ho][grain.No]byte
	See [grain.Ho][grain.No]byte
	Discard int
	ML *grain.MJList
	Over bool
	Rule *MJRule

	Uid  int
	Self  bool
}
func (self *MJUser) AddDown(gr *grain.MJGrain){
	self.Public.Down[self.Uid][gr.H][gr.N] += byte(gr.O)
}
func (self *MJUser) AddSee(gr *grain.MJGrain){
	self.Public.See[self.Uid][gr.H][gr.N] += byte(gr.O)
}
func (self *MJUser) DelSee(gr *grain.MJGrain){
	self.Public.See[self.Uid][gr.H][gr.N] -= byte(gr.O)
}
func (self *MJUser) Init(ML *grain.MJList,rule *MJRule,public *UserPublic,i int ){
	self.ML = ML
	self.Rule = rule
	self.Discard = -1
	self.GetDiscard(self.ReadNow(self.Now))
	self.Public = public
	self.Uid = i
	fmt.Println(self.Now,self.Discard)
	self.Self = false
}

func (self *MJUser)OutAna(se *AnalyInfo) (gr *grain.MJGrain){
//func (self *AnalyInfo)Out(r *MJRule,public *UserPublic) (gr *grain.MJGrain){
	var k []*SplitInfo
	var gs [4][]*grain.MJGrain
	for i,bit := range se.blockBit {
	//	bit := self.blockBit[i]
//		var k1,k2 int
		k1:=self.Rule.SplitStop(bit,se.block[i],nil)
		k2:=self.Rule.SplitStop(reversalByte(bit),reversalByte(se.block[i]),nil)

		s1 := self.Rule.GetSplitArrSum(k1)
		s2 := self.Rule.GetSplitArrSum(k2)
		if s1 == 0  || s2 == 0 {
			continue
		}
		c1 := len(k1)
		c2 := len(k2)
		var c int
		if c1 > c2 {
			c = c2
			k = k1[c:]
		}else{
			c = c1
			k = k2[c:]
		}
		for j:=0;j<c;j++{
			gs[k1[j].Val] = append(gs[k1[j].Val],&grain.MJGrain{H:se.i,N:k1[j].Block[0],O:1})
			gs[k2[j].Val] = append(gs[k2[j].Val],&grain.MJGrain{H:se.i,N:k2[j].Block[0],O:1})
		}
		if len(k) > 0 {
			for _,_k := range k {
				gs[_k.Val] = append(gs[_k.Val],&grain.MJGrain{H:se.i,N:_k.Block[0],O:1})
			}

			k = nil
		}
	}
	if len(gs[2]) != 0 {
		for i,g := range gs[2] {
			self.Public.Check1(g)
			SortGrain(gs[2],i)
		}
		return gs[2][0]
	}
	if len(gs[3]) != 0 {
		for i,g := range gs[3] {
			self.Public.Check1(g)
			gn := self.Now[g.H][g.N]
			if gn >1 && g.O == 0 {
				continue
			}
		//	g.O -= int(gn)
			SortGrain(gs[3],i)
		}
		return gs[3][0]
	}
	if len(gs[1]) != 0 {
		for i,g := range gs[1] {
			self.Public.Check1(g)
			SortGrain(gs[1],i)
		}
		return gs[1][0]
	}

	return nil

}
func (self *MJUser) ReadNow (Arr [grain.Ho][grain.No]byte) (analys []*AnalyInfo) {

//	var analys [grain.Ho]*AnalyInfo
	for i,n := range Arr {
		an := analyze(n)
		if an.num == 0 {
			if self.Discard <0 {
				self.Discard = i
			}
			continue
		}
		j:= len(analys)
		an.i = byte(i)
		analys = append(analys,an)
		SortAnaList(analys,j)
	}
//	for _,a := range analys {
//		fmt.Println(a)
//	}
	return analys

}
func (self *MJUser)GetDiscard (as []*AnalyInfo) *AnalyInfo {

	if self.Discard  < 0 {
		if as[1].num == as[2].num {
			if len(as[1].block) > len(as[2].block) {
				self.Discard = int(as[1].i)
				return as[1]
			}
		}
		self.Discard = int(as[2].i)
		return as[2]
	}else{
		for _,a := range as {
			if int(a.i) == self.Discard {
				return a
			}
		}
		return nil

	}


}

func (self MJUser) CheckOver (gr *grain.MJGrain) bool {
//	if self.Over {
//		return true
//	}
	self.Now[gr.H][gr.N]++
	ans := self.ReadNow(self.Now)
	if len(ans) == 3 {
		return false
	}
	isZ := 0
	for _,an:= range ans {
		isZ += an.Check(self.Rule)
	}
	if isZ == 1 {
		fmt.Println("over",self.Now,isZ,gr)
//		fmt.Println(self.Public.Down[self.Uid])
//		fmt.Println(ans[1])
//		var cmd string
//		fmt.Scanf("%s",&cmd)
		return true
	}
	return false
}

func (self *MJUser) SeeOut (gr *grain.MJGrain) (outgr *grain.MJGrain) {
	if self.Over {
		return nil
	}

	if self.CheckOver(gr) {
		self.AddDown(gr)
		self.Over = true
		fmt.Println(self.Now,self.Uid,"over pass")
		var key string
		fmt.Scanf("%s",&key)
		return nil
	}


	if int(gr.H) == self.Discard {
		return nil
	}
	N := self.Now[gr.H][gr.N]
//	fmt.Println(N)
	if N ==2 {
		gr.O = 3
		self.AddDown(gr)
		self.Now[gr.H][gr.N] = 0
		fmt.Println(self.Public.Down[self.Uid])
//		fmt.Println(self.Uid,self.Now)
		return self.Outs()
	}else if N ==3 {
		gr.O = 4
		self.AddDown(gr)
		self.Now[gr.H][gr.N] = 0
		fmt.Println(self.Public.Down[self.Uid])
//		fmt.Println(self.Uid,self.Now)
		if self.In() {
			return nil
		}
		return self.Outs()
	}
	return nil
}
func (self *MJUser) In () (b bool) {
	g := self.ML.Out()
	if g == nil {
		os.Exit(0)
	}
//	fmt.Println("in",g)
	b = self.CheckOver(g)
	g.O=1
	if !b {
		if self.Public.Down[self.Uid][g.H][g.N] == 3 {
			g.O = 3
			self.AddDown(g)
			return self.In()
		}
	//	n:= self.Now[g.H][g.N]
		if self.Now[g.H][g.N] == 3 {
			g.O = 4
			self.AddDown(g)
			self.Now[g.H][g.N] = 0
			return self.In()
		}
		self.Now[g.H][g.N]++
		return b
	}
	self.AddDown(g)
	fmt.Println("self over pass",self.Uid,self.Now,g)
	return b
}
func (self *MJUser) Outs () (gr *grain.MJGrain) {
	if self.Self {
		return self.OutSelf()
	}

	ans := self.ReadNow(self.Now)
	d := self.GetDiscard(ans)
	if d != nil {
		gr = &grain.MJGrain{H:d.i,N:d.block[0][0],O:1}
	}else{
		if len(ans) == 1 {
			gr =  self.OutAna(ans[0])
		}else{
			isC:=false
			if ans[1].num <3 {
				isC = true
				dans := self.ReadNow(self.Public.Down[self.Uid])
				for _,d := range dans {
					if d.i == ans[1].i  {
						isC = false
						break
					}
				}
			}
			if isC {
				gr = &grain.MJGrain{H:ans[1].i,N:ans[1].block[0][0],O:1}
			}else{
				a0 := ans[0].Check(self.Rule)
				if a0 < 0 {
					if self.DownGang(int(ans[0].i)) {
						return nil
					}
				}
				a1 := ans[1].Check(self.Rule)
				if a1 < 0 {
					if self.DownGang(int(ans[1].i)) {
						return nil
					}
				}
				if a1 < a0 {
					gr =  self.OutAna(ans[0])
				}else{
					gr =  self.OutAna(ans[1])
				}
			}
		}
	}
	gr.O=1
	self.Now[gr.H][gr.N]--
//	fmt.Println(self.Uid,self.Now,gr)
//	fmt.Println(self.Public.Down)
	self.AddSee(gr)
	return gr
}
func (self *MJUser) DownGang(I int) bool {
	isf := false
	for i,v := range self.Now[I] {
		if v == 4 {
			isf = true
			self.Now[I][i] = 0
			self.Public.Down[self.Uid][I][i] = 4
			break
		}
	}
	if !isf {
		panic(self)
	}
	return self.In()

}

func (self *MJUser) OutSelf () (gr *grain.MJGrain) {

	ans := self.ReadNow(self.Now)
	d := self.GetDiscard(ans)
	if d != nil {
		gr = &grain.MJGrain{H:d.i,N:d.block[0][0],O:1}
	}else{
		if len(ans) == 1 {
			gr =  self.OutAna(ans[0])
		}else{
			isC:=true
			if ans[1].num >3 {
				isC = false
			}else if ans[1].num == 3 {
				if self.Now[ans[1].i][ans[1].block[0][0]] == 3 {
					isC = false
				}
			}
			if isC {
				dans := self.ReadNow(self.Public.Down[self.Uid])
				for _,d := range dans {
					if d.i == ans[1].i  {
						isC = false
						break
					}
				}
			}

			if isC {
				gr = &grain.MJGrain{H:ans[1].i,N:ans[1].block[0][0],O:1}
			}else{
				a0 := ans[0].Check(self.Rule)
				if a0 < 0 {
					if self.DownGang(int(ans[0].i)) {
						return nil
					}
				}
				a1 := ans[1].Check(self.Rule)
				if a1 < 0 {
					if self.DownGang(int(ans[1].i)) {
						return nil
					}
				}
				if a1 < a0 {
					gr =  self.OutAna(ans[0])
				}else{
					gr =  self.OutAna(ans[1])
				}
			}
		}
	}
	fmt.Println("down")
	for i,d := range self.Public.Down {
		fmt.Println(i,d)
	}
	fmt.Println("see")
	for i,d := range self.Public.See {
		fmt.Println(i,d)
	}
	fmt.Println(ans[0].num,ans[1].num)
	fmt.Println("T",[][]int{[]int{0,1,2,3,4,5,6,7,8},[]int{0,1,2,3,4,5,6,7,8},[]int{0,1,2,3,4,5,6,7,8}})
	fmt.Println(self.Uid,self.Now,gr)
	var InH,InN int = -1,-1
	fmt.Scanf("%d %d\r",&InH,&InN)
	fmt.Println(InH,InN)
	if InH >0 && InH <3 && InN >0 && InN < 9 {
		gn :=self.Now[InH][InN]
		if gn >0{
			gr.H= byte(InH)
			gr.N = byte(InN)
		}
	}
	fmt.Println(gr)
	gr.O=1
	self.Now[gr.H][gr.N]--
//	fmt.Println(self.Uid,self.Now,gr)
//	fmt.Println(self.Public.Down)
	self.AddSee(gr)
	return gr
}
