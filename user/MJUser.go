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
//	Down [grain.Ho][grain.No]byte
//	See [grain.Ho][grain.No]byte
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
	self.Public = public
	self.GetDiscard(self.ReadNow(self.Now))
	self.Uid = i
	fmt.Println(self.Now,self.Discard)
	self.Self = true
}

func (self *MJUser) ReadNow (Arr [grain.Ho][grain.No]byte) (analys []*AnalyInfo) {

//	var analys [grain.Ho]*AnalyInfo
	for i,n := range Arr {
		an := analyze(n)
		for _,_n := range self.Public.Down[self.Uid][i]{
			an.num+=int(_n)
		}
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
	}
	return nil


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
		isZ += an.CheckExt(self.Rule)
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
		if self.GangTest(gr) {
			gr.O = 3
			self.AddDown(gr)
			self.Now[gr.H][gr.N] = 0
			fmt.Println(self.Public.Down[self.Uid])
//			fmt.Println(self.Uid,self.Now)
			return self.Outs()
		}
	}else if N ==3 {
		if self.GangTest(gr) {
			gr.O = 4
			self.AddDown(gr)
			self.Now[gr.H][gr.N] = 0
			fmt.Println(self.Public.Down[self.Uid])
//			fmt.Println(self.Uid,self.Now)
			if self.In() {
				return nil
			}
			return self.Outs()
		}else{
			gr.O = 3
			self.AddDown(gr)
			self.Now[gr.H][gr.N] = 1
			fmt.Println(self.Public.Down[self.Uid])
			return self.Outs()

		}
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
	self.Over = b
	g.O=1
	if !b {
		if self.Public.Down[self.Uid][g.H][g.N] == 3 {
			if self.GangTest(g) {
		//		g.O = 1
				self.AddDown(g)
				return self.In()
			}
		}
	//	n:= self.Now[g.H][g.N]
		if self.Now[g.H][g.N] == 3 {
			if self.GangTest(g) {
				g.O = 4
				self.AddDown(g)
				self.Now[g.H][g.N] = 0
				return self.In()
			}
		}
		self.Now[g.H][g.N]++
		return b
	}
	var key string
	fmt.Scanf("%s",&key)
	self.AddDown(g)
	fmt.Println("self over pass",self.Uid,self.Now,g)
	return b
}
func (self *MJUser) Outs () (gr *grain.MJGrain) {
	ans := self.ReadNow(self.Now)
	d := self.GetDiscard(ans)
	if d != nil {
		gr = &grain.MJGrain{H:d.i,N:d.block[0][0],O:1}
	}else{
		a0 := ans[0].Check(self.Rule)
		if a0 < 0 {
			if self.DownGang(int(ans[0].i)) {
				return nil
			}
		}
		if len(ans) > 1 {
			if ans[1].num < 3 {
				gr = &grain.MJGrain{H:ans[1].i,N:ans[1].block[0][0],O:1}
			}else{
				a1 := ans[1].Check(self.Rule)
				if a1 < 0 {
					if self.DownGang(int(ans[1].i)) {
						return nil
					}
				}
			}
		}
		gr = self.FindOutGrain(ans)
		if gr == nil {
			panic("9")
			return nil
		}
	}
	if self.Self && d == nil {
		self.OutSelf(gr)
	}
	self.Now[gr.H][gr.N]--
	gr.O=1
	if self.Public.Down[self.Uid][gr.H][gr.N] == 3 {

		self.AddDown(gr)
		if self.In() {
			return nil
		}
		return self.Outs()
	}else{
		self.AddSee(gr)
	}
	return gr
}
func (self *MJUser) OutSelf(gr *grain.MJGrain) {
	fmt.Println("down")
	for i,d := range self.Public.Down {
		fmt.Println(i,d)
	}
	fmt.Println("see")
	for i,d := range self.Public.See {
		fmt.Println(i,d)
	}
//	fmt.Println(ans[0].num,ans[1].num)
	fmt.Println("T",[][]int{[]int{0,1,2,3,4,5,6,7,8},[]int{0,1,2,3,4,5,6,7,8},[]int{0,1,2,3,4,5,6,7,8}})
	fmt.Println(self.Uid,self.Now,gr)
	var InH,InN int = -1,-1
	fmt.Scanf("%d %d\r",&InH,&InN)
	fmt.Println(InH,InN)
	if InH >=0 && InH <3 && InN >=0 && InN < 9 {
		gn :=self.Now[InH][InN]
		if gn >0{
			gr.H= byte(InH)
			gr.N = byte(InN)
		}
	}
	fmt.Println(gr)
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

func (self *MJUser) FindOutGrain(ans []*AnalyInfo) (gr *grain.MJGrain) {
	var gs []*grain.MJGrain
	var v int
	var k []*SplitInfo
	for _i:=len(ans)-1;_i>=0;_i--{
		an:=ans[_i]
		for i,a := range an.blockBit {
			k1:=self.Rule.SplitStop(a,an.block[i],nil)
			k2:=self.Rule.SplitStop(reversalByte(a),reversalByte(an.block[i]),nil)
			s1 := self.Rule.GetSplitArrSum(k1)
			s2 := self.Rule.GetSplitArrSum(k2)
			if s1 == 0  || s2 == 0 {
				continue
			}
			c1 := len(k1)
			c2 := len(k2)
			var c int
			if c1 == c2 {
				v+=s2
				c = c2
			}else if c1 > c2 {
				v+=s2
				c = c2
				k = k1[c:]
			}else{
				v+=c1
				c = c1
				k = k2[c:]
			}

			for j:=0;j<c;j++{
				for _,b:= range k1[j].Block {
					gs = append(gs,&grain.MJGrain{H:an.i,N:b,O:j})
					SortdGrain(gs,len(gs)-1)
				}
				for _,b:= range k2[j].Block {
					gs = append(gs,&grain.MJGrain{H:an.i,N:b,O:j})
					SortdGrain(gs,len(gs)-1)
				}
			}
			if len(k) > 0 {
				for j,_k := range k {
					for _,b:= range _k.Block {
						gs = append(gs,&grain.MJGrain{H:an.i,N:b,O:j})
						SortdGrain(gs,len(gs)-1)
					}
				}
				k = nil
			}

		}
	}
	le := len(gs)
	if le == 0 {
		panic(12)
		return nil
	}
	if le ==1 {
		return gs[0]
	}
	return self.NeedThink2(gs,v)
}
func (self *MJUser) NeedThink2(grs []*grain.MJGrain,val int) (gr *grain.MJGrain) {
//	fmt.Println(grs)
	var Tmp [grain.Ho][grain.No]byte
	j :=0
	fmt.Println("---0")
	var grN []*grain.MJGrain
	for _,g := range grs {
		if Tmp[g.H][g.N] >0 {
			continue
		}
		Tmp[g.H][g.N]++
//		if g.O > val {
//			continue
//		}
		self.Think1(g)
		grN = append(grN,g)
		SortdGrain(grN,j)
		fmt.Println(g)
		j++
	}
	if len(grN)==0{
		fmt.Println(val)
		fmt.Println(self)
		panic(10)
	}
	gr = grN[0]
	if len(grN) == 1 {
		return gr
	}

	var gr1 []*grain.MJGrain
	var gr2 []*grain.MJGrain
	j2 := 0
	j1 :=0
	fmt.Println("---1")
	lastO := gr.O
	for _,g := range grN {
		if lastO != g.O {
			break
		}
		fmt.Println(g)
		mv:= self.Think2(g)
		if mv==0 {
			g.O =int( self.Now[g.H][g.N])*10
			n0 := int(g.N)-1
			if n0 >= 0 {
				g.O+=int( self.Now[g.H][n0])
			}
			n1 := g.N+1
			if n1 <= 8 {
				g.O+=int( self.Now[g.H][n1])
			}

			gr1 = append(gr1,g)
			SortdGrain(gr1,j1)
			j1++
		}else{
			g.O = mv
			gr2 = append(gr2,g)
			SortGrain(gr2,j2)
			j2++
		}

	}
	if j2>0 {
		fmt.Println("j2")
	//	return gr2[0]
		lastO := gr2[0].O
		for _i,_g := range gr2{
			if lastO != _g.O {
				break
			}
			fmt.Println(_g)
			self.Public.Check1(_g)
			SortGrain(gr2,_i)
		}
		return gr2[0]
	}
	if j1>0 {
		fmt.Println("j1")
		lastO := gr1[0].O
		for _i,_g := range gr1{
			if lastO != _g.O {
				break
			}
			fmt.Println(_g)
			self.Public.Check1(_g)
			SortGrain(gr1,_i)
		}
		return gr1[0]
	}

	panic(11)
	return nil
}
func (self MJUser) Think2(gr *grain.MJGrain) (v int) {
	self.Now[gr.H][gr.N] --
	var i,j int
	for i  = 0;i<grain.Ho;i++ {
		if i == self.Discard {
			continue
		}
		for j=0;j<grain.No;j++{
			self.Now[i][j] ++
			val := self.GetCheckDiff()
			self.Now[i][j] --
			if val == 1 {
		//		fmt.Println(val)
				g := &grain.MJGrain{byte(i),byte(j),4-int(self.Now[i][j])}
				self.Public.Check2(g)
				if g.O >0 && self.CheckDadui() {
					v+=4
				}else{
					v+=g.O
				}
			}
		//	if val >= gr.O {
		//		continue
		//	}
		}
	}
	return v
}
func (self MJUser) Think1(gr *grain.MJGrain) {
	self.Now[gr.H][gr.N] --
	gr.O= self.GetCheckDiff()
}
func (self *MJUser) CheckDadui() bool {
	ans := self.ReadNow(self.Now)
	for _,an:= range ans {
		for i,a := range an.blockBit{
			in := self.Rule.SplitStop(a,an.block[i],nil)
			for _,n := range in {
				if len(n.Block) == 3 && n.Val == 0 {
					return false
				}
			}
		}
	}
	return true
}
func (self *MJUser) GetCheckDiff() (v int) {

	ans := self.ReadNow(self.Now)
	v = 0
	for _,an:= range ans {
		v += an.Check(self.Rule)
	}
//	fmt.Println(self.Uid,"diff",v)
	return v

}
func (self MJUser) GangTest(g *grain.MJGrain) bool {
//	return 0 > self.GetCheckDiff()

	N :=self.Now[g.H][g.N]
	self.Now[g.H][g.N]=0
	v0 := self.GetCheckDiff()
	self.Now[g.H][g.N]=N
	v1 := self.GetCheckDiff()
	return v1>v0
}
