package user
import (
	"../grain"
	"fmt"
)
type MJAI struct {
	Now [grain.Ho][grain.No]byte
	Public *UserPublic
	Rule *MJRule
	Discard int
	Uid int
	LastSee int
	LastData string
	Nc int
	Num int
}
func (self *MJAI) Init(rule *MJRule,public *UserPublic,i int ){
	self.Rule = rule
	self.Discard = -1
	self.Public = public
//	self.GetDiscard(self.ReadNow(self.Now))
	self.Uid = i
	self.Nc = 0
	self.Num = 0
//	fmt.Println(self.Now,self.Discard)
//	self.Self = true
}
func (self MJAI) GangTest(g *grain.MJGrain) bool {
//	return 0 > self.GetCheckDiff()

	N :=self.Now[g.H][g.N]
	self.Now[g.H][g.N]=0
	v0 := self.GetCheckDiff()
	self.Now[g.H][g.N]=N
	v1 := self.GetCheckDiff()
	return v1>v0
}

func (self *MJAI)GetDiscard (as []*AnalyInfo) *AnalyInfo {

	if self.Discard  < 0 {
		if as[1].num == as[2].num {
			if as[1].sum > as[2].sum {
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
func (self *MJAI) AddDown(gr *grain.MJGrain){
	self.Public.Down[self.Uid][gr.H][gr.N] += byte(gr.O)
}
func (self *MJAI) AddSee(gr *grain.MJGrain){
	self.Public.See[self.Uid][gr.H][gr.N] += byte(gr.O)
}
//func (self *MJAI) DelSee(gr *grain.MJGrain){
//	self.Public.See[self.Uid][gr.H][gr.N] -= byte(gr.O)
//}
func (self *MJAI) ReadNow (Arr [grain.Ho][grain.No]byte) (analys []*AnalyInfo) {

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
func (self *MJAI) OutDiscard (num int) (outs []int) {
	ans := self.ReadNow(self.Now)
	d := self.GetDiscard(ans)
	v :=0
	fmt.Println(self.Now,num)
	if d != nil {
		b:=int(d.i*9)
		fmt.Println(b,d.i,num)
		for i,n := range self.Now[d.i]{
//			for _i :=0;_i<int(n);_i++{
				if n==0 {
					continue
				}
				outs = append(outs,b+i)
				v++
				if v == num {
//					fmt.Println(outs)
					return outs
				}
//			}
		}
//		fmt.Println(outs)
		return outs
	}
	return nil
}

func (self *MJAI) FindOutGrain(ans []*AnalyInfo) (gr *grain.MJGrain) {
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
//		panic(12)
		return nil
	}
	if le ==1 {
		return gs[0]
	}
	return self.NeedThink2(gs,v)
}
func (self *MJAI) NowDiffToSee ( n [3][9]byte)  {
	var _n,i,j,d int
	for i=0;i<3;i++{
		for j=0;j<9;j++{
			_n = int(self.Now[i][j])
			if _n == 0 {
				continue
			}
			d = _n - int(n[i][j])
			if d==0 || d>1 {
				continue
			}
			self.AddSee(&grain.MJGrain{H:byte(i),N:byte(j),O:d})
		}
	}
}
func (self *MJAI) Outs () int {
	ans := self.ReadNow(self.Now)
	d := self.GetDiscard(ans)
	var gr *grain.MJGrain
	if d != nil {
		gr = &grain.MJGrain{H:d.i,N:d.block[0][0],O:1}
	//	self.Now[gr.H][gr.N]--
	}else{
	//	if len(ans) > 1 && ans[1].num < 3 {
	//		gr = &grain.MJGrain{H:ans[1].i,N:ans[1].block[0][0],O:1}
	//	}else{
			gr = self.FindOutGrain(ans)
			if gr == nil {
				return -1
			}
	//	}
	}
//	self.AddSee(gr)
	return int(gr.H*9+gr.N)
}
func (self *MJAI) GetCheckDiff() (v int) {

	ans := self.ReadNow(self.Now)
	v = 0
	for _,an:= range ans {
		v += an.Check(self.Rule)
	}
//	fmt.Println(self.Uid,"diff",v)
	return v

}
func (self *MJAI) SeeOut (n int) (int) {
	gr := &grain.MJGrain{H:byte(n/9),N:byte(n%9),O:1}
	N := self.Now[gr.H][gr.N]
//	fmt.Println(N)
	if N ==2 {
		if self.GangTest(gr) {
			gr.O = 2
			self.AddDown(gr)
			self.Now[gr.H][gr.N] = 0
	//		fmt.Println(self.Public.Down[self.Uid])
//			fmt.Println(self.Uid,self.Now)
			return 1
		}
	}else if N ==3 {
//		return 2
//		if self.GangTest(gr) {
//			gr.O = 4
//			self.AddDown(gr)
//			self.Now[gr.H][gr.N] = 0
//	//		fmt.Println(self.Public.Down[self.Uid])
////			fmt.Println(self.Uid,self.Now)
//	//		if self.In() {
//	//			return nil
//	//		}
//			return 2
////			return self.Outs()
//		}else{
			gr.O = 3
			self.AddDown(gr)
			self.Now[gr.H][gr.N] = 1
		//	fmt.Println(self.Public.Down[self.Uid])
			return 2
//			return self.Outs()

//		}
	}
	return 0
}
func (self *MJAI) CheckDadui() bool {
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
func (self MJAI) Think1(gr *grain.MJGrain) {
	self.Now[gr.H][gr.N] --
	gr.O= self.GetCheckDiff()
}
func (self MJAI) Think2(gr *grain.MJGrain) (v int) {
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
func (self *MJAI) NeedThink2(grs []*grain.MJGrain,val int) (gr *grain.MJGrain) {
//	fmt.Println(grs)
	var Tmp [grain.Ho][grain.No]byte
	j :=0
//	fmt.Println("---0")
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
	//	fmt.Println(g)
		j++
	}
	if len(grN)==0{
	//	fmt.Println(val)
	//	fmt.Println(self)
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
//	fmt.Println("---1")
	lastO := gr.O
	for _,g := range grN {
		if lastO != g.O {
			break
		}
	//	fmt.Println(g)
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
//		fmt.Println("j2")
	//	return gr2[0]
		lastO := gr2[0].O
		for _i,_g := range gr2{
			if lastO != _g.O {
				break
			}
		//	fmt.Println(_g)
			self.Public.Check1(_g)
			SortGrain(gr2,_i)
		}
		return gr2[0]
	}
	if j1>0 {
	//	fmt.Println("j1")
		lastO := gr1[0].O
		for _i,_g := range gr1{
			if lastO != _g.O {
				break
			}
		//	fmt.Println(_g)
			self.Public.Check1(_g)
			SortGrain(gr1,_i)
		}
		return gr1[0]
	}

	panic(11)
	return nil
}
