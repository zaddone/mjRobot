package user
import (
	"../grain"
	"fmt"
	"os"
//	"math"
)
const (
	UN int = 13
//	All int = 14
)
type SplitInfo struct {
	Mapkey []byte
	Block []byte
	Val int
}
type MJRule struct {
	Rule map[string]int

}
func IntToByte(val []int) string {

	b := make([]byte,len(val))
	for i,v := range val {
		b[i] = 0
		for _i:=0; _i<v; _i++ {
			b[i] = b[i]<<1
			b[i]++
		}
	}
//	fmt.Println(b)
	return string(b)

}
func (self *MJRule) Init(){
	self.Rule = make(map[string]int)

	self.Rule[IntToByte([]int{4})] = -100
	self.Rule[IntToByte([]int{3})] = 0
	self.Rule[IntToByte([]int{2})] = 1
	self.Rule[IntToByte([]int{1})] = 2

//	self.Rule[IntToByte([]int{4,4})] = 0
//	self.Rule[IntToByte([]int{3,4})] = 0

//	self.Rule[IntToByte([]int{2,4})] = 1
//	self.Rule[IntToByte([]int{1,4})] = 2

//	self.Rule[IntToByte([]int{4,3})] = 0
//	self.Rule[IntToByte([]int{3,3})] = 0
	self.Rule[IntToByte([]int{2,3})] = 1
	self.Rule[IntToByte([]int{1,3})] = 2

//	self.Rule[IntToByte([]int{4,2})] = 1
//	self.Rule[IntToByte([]int{3,2})] = 1
	self.Rule[IntToByte([]int{2,2})] = 3
	self.Rule[IntToByte([]int{1,2})] = 3

//	self.Rule[IntToByte([]int{4,1})] = 2
//	self.Rule[IntToByte([]int{3,1})] = 2
	self.Rule[IntToByte([]int{2,1})] = 3
	self.Rule[IntToByte([]int{1,1})] = 3

}

func reversalByte(b []byte) (c []byte){

	le := len(b)
	c = make([]byte,le)
	le--
	for i,_b:= range b{
		c[le-i] = _b
	}
	return c
}

func  (self *MJRule) SplitStop(bit []byte,bitn []byte,in []*SplitInfo) []*SplitInfo {
	le := len(bit)
	if le==0 {
		return in
	}
	k := self.Rule[string(bit[0:1])]
	if k == 0 {
		return self.SplitStop(bit[1:],bitn[1:],in)
	}
	if le<3 {
	//	return append(in,self.Rule[string(bit)])
		return append(in,&SplitInfo{bit,bitn,self.Rule[string(bit)]})
	}

	var newb []byte
	var newbn []byte
	for _i,b := range bit[:3] {
		b = b>>1
		if b >0 {
			newb = append(newb,b)
			newbn = append(newbn,bitn[_i])
		}else{
			if len(newb)>0 {
			//	in = append(in, self.Rule[string(newb)])
				in = append(in, &SplitInfo{newb,newbn,self.Rule[string(newb)]})
				newb = nil
				newbn = nil
			}
		}
	}

	if len(newb) == 0 {
		k1:= self.SplitStop(bit[3:],bitn[3:],nil)
		k2:= self.SplitStop(reversalByte(bit[3:]),reversalByte(bitn[3:]),nil)

		s1 := self.GetSplitArrSum(k1)
		s2 := self.GetSplitArrSum(k2)
		if s1 > s2 {
			if s2>0 {
				in = append(in,k2...)
			}
		}else{
			if s1>0 {
				in = append(in,k1...)
			}
		}
		return in
	}

	return self.SplitStop(append(newb,bit[3:]...),append(newbn,bitn[3:]...),in)

}
func (self *MJRule)GetSplitArrSum(ks []*SplitInfo) (s int){

	if ks == nil {
		return 0
	}
	for _,_k:= range ks {
//		_k.Val = self.Rule[string(_k.Mapkey)]
		s+=_k.Val
	}
	return s

}
func  (self *MJRule) SplitOut(bit []byte) int {
	le := len(bit)
	if le==0 {
		return 0
	}
	k := self.Rule[string(bit[0:1])]
//	fmt.Println("Rule",bit[0:1],k)
	if k == 0 {
		return self.SplitOut(bit[1:])
	}
	if le<3 {
		return self.Rule[string(bit)]
	}

	k =0
	var newb []byte
	for _,b := range bit[0:3] {
		b = b>>1
		if b >0 {
			newb = append(newb,b)
		}else{
			if len(newb)>0 {
				k += self.Rule[string(newb)]
				newb = nil
			}
		}
	}
	if len(newb) == 0 {
		k1:= self.SplitOut(bit[3:])
		k2:= self.SplitOut(reversalByte(bit[3:]))
		if k1 > k2 {
			k+= k2
		}else{
			k+= k1
		}
		return k
	}

	k+= self.SplitOut(append(newb,bit[3:]...))
	return k
}

type UserPublic struct {
	Down [4][grain.Ho][grain.No]byte
	See [4][grain.Ho][grain.No]byte
}
func (self *UserPublic) Check1(gr *grain.MJGrain) {
	gr.O = 0
	for _,d := range self.Down {
		gr.O += int(d[gr.H][gr.N])
	}
	for _,s := range self.See {
		gr.O += int(s[gr.H][gr.N])
	}
}


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
}

type AnalyInfo struct {
	blockBit [][]byte

	block [][]byte
	blockNum []int
	sum int
	num int
	i byte
//	Cov float64
	n []byte
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

func (self *AnalyInfo)Check(r *MJRule) (ko int){

	ko = 0
	for _,bit := range self.blockBit {
	//	bit := self.blockBit[i]
//		var k1,k2 int
		k1:=r.SplitOut(bit)
		k2:=r.SplitOut(reversalByte(bit))
		if k1 == 0 || k2 == 0 {
			if k1<0 || k2 <0 {
				continue
			}
		}else{
			if k1 < 0 {
				return k1
			}
			if k2 < 0 {
				return k2
			}
		}
//		fmt.Println("check",k1,k2)
		if k1 > k2 {
			ko+= k2
		}else{
			ko+=k1
		}
//		if ko >1 {
//			return ko
//		}
	}
	return ko
}


func analyze (n [grain.No]byte) (ana *AnalyInfo)  {
	ana = new(AnalyInfo)
	ana.n = n[0:]
	var tmp []byte = nil
	var tmpbit []byte = nil
	var tmpNum int = 0
	for i,_n := range n {
		if _n == 0 {
			if tmp != nil {
				ana.block = append(ana.block,tmp)
				ana.blockNum = append(ana.blockNum,tmpNum)
				ana.blockBit = append(ana.blockBit,tmpbit)
				tmp= nil
				tmpNum = 0
				tmpbit = nil
			}
			continue
		}
		ana.num += int(_n)
		ana.sum += i

		tmpNum ++
		var bit byte
		for _i:=0;_i<int(_n);_i++ {
			bit = bit<<1
			bit++
	//		tmp = append(tmp,i)
		}
		tmp = append(tmp,byte(i))
		tmpbit = append(tmpbit,bit)
	}
	if len(tmp) > 0 {
		ana.block = append(ana.block,tmp)
		ana.blockNum = append(ana.blockNum,tmpNum)
		ana.blockBit = append(ana.blockBit,tmpbit)
	}
	return ana
}

func SortAnaList(as []*AnalyInfo,n int){

	if n==0 {
		return
	}
	_n := n-1
	if as[n].num >as[_n].num {
		as[n],as[_n] = as[_n],as[n]
		SortAnaList(as,_n)
	}

}
func SortGrain(as []*grain.MJGrain,n int){

	if n==0 {
		return
	}
	_n := n-1
	if as[n].O > as[_n].O {
		as[n],as[_n] = as[_n],as[n]
		SortGrain(as,_n)
	}

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
		fmt.Println(self.Public.Down[self.Uid])
//		fmt.Println(ans[1])
		var cmd string
		fmt.Scanf("%s",&cmd)
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
	fmt.Println("in",g)
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
	fmt.Println(self.Uid,self.Now,gr)
//	fmt.Println(self.Public.Down[self.Uid])
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

