package ai
import(
//	"strconv"
)
type Block struct {
	Bl [][]int
	Num int
	Numi int
	No []int
	I  int
	W  int
	We []int
}
//func (self *Block) FindWeight(grs []*MJGrain) {
func (self *Block) RunFind(a *Ai ) []*MJGrain {
	if len(self.We) == 0 {
		return nil
	}
	var grs []*MJGrain
	for _,i := range self.We {
		if self.No[i] >2 {
			continue
		}
		g :=&MJGrain{H:self.I,N:i,O:[]int{2,0,self.No[i],0}}
		if a.GetPublicNum(g,false) >0 {
			g.O[1] =1
		}else {
			g.O[1] = 0
		}
		grs = append(grs,g)
	}
	var No [9]int
	copy(No[0:],self.No)
	for i,n := range No {
		if n == 4 {
			continue
		}
		if n == 0 {
			isj := false
			i_ := i+1
			if i_<9 {
				if No[i_] > 0 {
					isj = true
				}
			}
			if !isj {
				i_  = i-1
				if i_ >= 0 {
					if No[i_] > 0 {
						isj = true
					}

				}
			}
			if !isj {
				continue
			}
		}
		_n := a.GetPublicNum(&MJGrain{H:self.I,N:i},true)
		if _n == 4 {
			continue
		}
		_n = 4 - _n
		No[i] ++
		tmpBlock := new(Block)
		tmpBlock.Init(No[0:],self.I)
	//	if len(self.We) >= len(tmpBlock.We) {
			tmpBlock.Find(a,grs,0,_n*10)
	//	}
		No[i] --
	}
	return grs
}
func  (self *Block)Find(a *Ai,grs []*MJGrain,level int,num int) {
//	if len(self.We) == 0 {
//		return
//	}
	var NewGrs []*MJGrain
//	if len(grs) > len(self.We){
		for _,g := range grs {
			isH := false
			for _,n := range self.We {
				if n == g.N {
					isH = true
					NewGrs = append(NewGrs,g)
					break
				}
			}
			if !isH {
	//			g.O[0] = level
//				g.O[2] += num
				if g.O[0] > level  {
					g.O[0] = level
					g.O[3] = num
				}else if g.O[0] == level {
					g.O[3] += num
				}
			}
		}
		if len(NewGrs) ==  0 {
			return
		}
//	}else{
//		NewGrs = grs
//	}
	level ++
	if level > 1 {
		return
	}
	var No [9]int
	copy(No[0:],self.No)
	for i,n := range No {
		if n == 4 {
			continue
		}
		if n == 0 {
			isj := false
			i_ := i+1
			if i_<9 {
				if No[i_] > 0 {
					isj = true
				}
			}
			if !isj {
				i_  = i-1
				if i_ >= 0 {
					if No[i_] > 0 {
						isj = true
					}

				}
			}
			if !isj {
				continue
			}
		}
		_n := a.GetPublicNum(&MJGrain{H:self.I,N:i},true)
		if _n == 4 {
			continue
		}
		_n = 4 - _n
		No[i] ++
		tmpBlock := new(Block)
		tmpBlock.Init(No[0:],self.I)
	//	if len(self.We) >= len(tmpBlock.We) {
			tmpBlock.Find(a,NewGrs,level,_n+num)
	//	}
		No[i] --
	}

}
func (self *Block) Update() {
	self.Num = 0
	self.Numi = 0
	self.Bl = nil

	var tmparr []int = nil

	for j,n := range self.No {
		if n == 0 {
			if (len(tmparr)>0){
				self.Bl = append(self.Bl,tmparr)
				tmparr = nil
			}
			continue
		}
		self.Num += n
		self.Numi ++
		tmparr = append(tmparr,j)
	}
	if (len(tmparr)>0){
		self.Bl = append(self.Bl,tmparr)
		tmparr = nil
	}

	self.SetWeight()
}
func (self *Block) Init(no []int,i int) {

	self.I = i
	self.No = no

	self.Num = 0
	self.Numi = 0
	self.Bl = nil

	var tmparr []int = nil

	for j,n := range no {
		if n == 0 {
			if (len(tmparr)>0){
				self.Bl = append(self.Bl,tmparr)
				tmparr = nil
			}
			continue
		}
		self.Num += n
		self.Numi ++
		tmparr = append(tmparr,j)
	}
	if (len(tmparr)>0){
		self.Bl = append(self.Bl,tmparr)
		tmparr = nil
	}

	self.SetWeight()
}
func (self *Block) GetWeight() int {

	return self.W

}
func (self *Block) SetWeight() {

	self.W = 0
	self.We = nil
	var No [9]int
	copy(No[0:],self.No)

	for _,bl := range self.Bl {
		w,e := GetWaightVal(No,bl)
		self.W += w
		self.We = append(self.We,e...)
	}

}
func WaightVal(No [9]int,bl []int,w int,o []int) (int,[]int) {
	le := len(bl)
	if le == 0 {
		return w,o
	}
	if le < 3 {
		for _,i := range bl {
			if No[i] < 3 {
				w += (3 - No[i])
				//if No[i] == 1 {
					o = append(o,i)
				//}
			}
		}
		return w,o
	}
	n:= No[bl[0]]
	if n > 2 {
		return WaightVal(No,bl[1:],w,o)
	}
	var start int = 3
	isb := false
	for i:=2;i>=0;i--{
		No[bl[i]] --
		if No[bl[i]] == 0 {
			if !isb {
				isb = true
			}
			continue
		}else {
			if !isb {
				start = i
			}else{
				w += 3 - No[bl[i]]
				//if No[i] == 1 {
					o = append(o,bl[i])
				//}
			}
		}
	}

	bl = bl[start:]
	w1,o:= WaightVal(No,bl,w,o)
	len_bl := len(bl)
	if len_bl < 3 {
		return w1,o
	}
	rbl := make([]int,len_bl)
	for i,n := range bl {
		rbl[len_bl - i - 1] = n
	}
	os := make([]int,len(o))
	copy(os,o)
	w2,os := WaightVal(No,rbl,w,os)

	if w2 < w1 {
		return w2,os
	}
	if w2 == w1 {
		var no []int = nil
		for _,_o := range os {
			isS := false
			for _,o_ := range o {
				if _o == o_ {
					isS = true
					break
				}
			}
			if !isS {
				no = append(no,_o)
			}
		}
		o = append(o,no...)

	}
	return w1,o

}
func GetWaightVal(No [9]int, bl []int) (int,[]int) {

	_w1,o1 := WaightVal(No,bl,0,nil)
	if  _w1 == 0 {
		return 0,nil
	}
	l := len(bl)
	if l < 3 {
		return _w1,o1
	}

	rs := make([]int,l)
	for i,_s := range bl {
		rs[l-i-1] = _s
	}
	_w2,o2 := WaightVal(No,rs,0,nil)
	if _w2 < _w1 {
		return _w2,o2
	}

	if _w2 == _w1 {
		var no []int = nil
		for _,_o := range o2 {
			isS := false
			for _,o_ := range o1 {
				if _o == o_ {
					isS = true
					break
				}
			}
			if !isS {
				no = append(no,_o)
			}
		}
		o1 = append(o1,no...)
	}
	return _w1,o1

}
func SortBlock(bs []*Block,i int) {
	if i == 0 {
		return
	}
	I := i-1
	if bs[i].Num < bs[I].Num {
		bs[i],bs[I] = bs[I],bs[i]
	}
	SortBlock(bs,I)
}
