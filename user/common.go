package user
import (
	"../grain"
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

