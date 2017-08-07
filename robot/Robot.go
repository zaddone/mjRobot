package robot
import(
	"fmt"
	"strconv"
//	"math"
)
type Robot struct {
	Now [3][9]int
	Invalid [3]bool
	Inv  int

	Public *UserPublic
	LastSee int
	LastForData [2]string
//	Now1 [3][]int
	Uid int
	NowVal int
	NowU [3]int
	BaseVal int
}
func (self *Robot) Init(p *UserPublic,u int){
	self.Uid = u
	self.Public = p
	self.Inv = -1
	self.BaseVal = 5
	self.NowU = [3]int{0,0,0}
}
func (self *Robot) UpdateNow(str []string){
	var i,j int
	for _,s := range str {
		n,err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		i = n/9
		j = n%9
		self.Now[i][j]++
	}
}
func (self *Robot) SetValid(i int) {
//	self.Invalid = [3]bool{false,false,false}
//	self.Invalid[i] = true
	self.Inv = i
}
func (self *Robot) InitNow(str []string){
	var now [3][9]int
	for _,s := range str {
		n,err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		now[n/9][n%9]++
	}
	self.Now = now
//	self.SetValid(valid)
//	self.Uid = 3
	if self.Inv == -1 {
		self.SetValid(self.GetDiscard(true))
	}

}
func (self Robot) CheckGang() *MJGrain {
	for i,na := range self.Now{
		if self.Inv == i {
			continue
		}
		for j,n := range na {
			_n:= self.Public.Down[self.Uid][i][j]
			if _n >0 {
				n += _n+1
			}
			if n == 0 {
				continue
			}
			if n >= 4 {
				return &MJGrain{H:i,N:j}
			}
		}
	}
	return nil
}
func (self Robot) SeeSelf(isf int) int {
	isF := self.CheckFull(true)
	if isF>0 { //&& isF + isf >= 0 {
		return 3+30
	}
	for i,no := range self.Now  {
		if i == self.Inv {
			continue
		}
		for j,n := range no {
			if n == 0 {
				continue
			}
			if n  == 4 || self.Public.Down[self.Uid][i][j]  >0 {
		//		self.Now[i][j] -= n
		//			gs,_ := self.IsTing(isf)
		//			if gs
		//		self.Now[i][j] += n
				self.Public.Down[self.Uid][i][j]  = 3
				self.Public.See[self.Uid][i][j] = 1
				return 2+30
			}
		}
	}
	return self.Outs(isf,false)
}
func (self Robot) SeeOut(v int,isf int) int {

	i:=v/9
	j := v%9
	self.Now[i][j]++
	isF := self.CheckFull(false)
	self.Now[i][j]--
	if isF >0 {
		if isF + isf >=0 {
			return 3
		}else  {
			da := 0
			for _,no := range self.Now {
				for _,n := range no {
					if n == 1 {
						da++
					}
				}
			}
			if da < 2 {
				return 3
			}
		}
	}
	fmt.Println("g-c",i,j,v)
	gs,l := self.IsTing(isf,false)
	l = len(gs)*4 - l
	if l > 0 {

		if self.Now[i][j] == 3 {
			self.Public.Down[self.Uid][i][j]=3
			return 2

//			self.Now[i][j] = 0
//			_gs,_l := self.IsTing(isf,false)
//			lgs := len(_gs)
//			if lgs >= len(gs) || lgs*4-_l >= l  {
//				self.Public.Down[self.Uid][i][j]=3
//				return 2
//			}
//			self.Now[i][j] = 3
		}
		if self.Now[i][j] > 1 {
			self.Now[i][j]-=2
			for _i,no := range self.Now {
				if _i == self.Inv {
					continue
				}
				for _j,_ := range no {
					self.Now[_i][_j] --
					_gs,_ := self.IsTing(isf,false)
					self.Now[_i][_j] ++
					if len(_gs)>0 {
						self.Public.Down[self.Uid][i][j]=2
						return 1
					}
				}
			}
		}
		return 0
	}
	if self.Now[i][j] == 2 {
		for _,n := range self.Now[self.Inv] {
			if n >0 {
				self.Public.Down[self.Uid][i][j]=2
				return 1
			}
		}
		c:=&Bl{I:i}
		c.Init(self.Now[i][0:])
		_,out1 := c.GetWeightOut(self.Now[i][0:])
		self.Now[i][j] = 0
		
		cb := c.UpdateRow(self.Now[i][0:])
		_,out2:= cb.GetWeightOut(self.Now[i][0:])
		if len(out2) - len(out1) < 2 {
			self.Public.Down[self.Uid][i][j]=2
			return 1
		}
	}else if self.Now[i][j] == 3 {
		self.Public.Down[self.Uid][i][j]=3
		return 2
	}
	return 0
}
//func (self *Robot) Outd(isf int) int {
//	cs:= &Covs{Ro:self}
//	for _i,no := range self.Now {
//		c := &Bl{I:_i}
//		c.Init(no[0:])
//		lc := len(cs.blocks)
//		cs.blocks = append(cs.blocks,c)
//		cs.SortCovNum(lc)
//	}
//}
func (self *Robot) Outs(isf int,No30 bool) int {
	cs:= &Covs{Ro:self}
	for _i,no := range self.Now {
		if _i == self.Inv {
			for j,n := range no {
				if n != 0 {
					fmt.Println("out -1")
					return self.Inv*9+j
				}
			}
		}
	}
	if No30 {
		isF := self.CheckFull(true)
		if isF>0 { //&& isF + isf >= 0 {
			fmt.Println("full 33")
			return 33
		}
	}
	for _i,no := range self.Now {
		c := &Bl{I:_i}
		c.Init(no[0:])
		for _j,n := range self.Public.Down[self.Uid][_i]{
			if n > 0 {
				if No30{
					if self.Now[_i][_j] >0 {
						fmt.Println(self.Public.Down[self.Uid][_i])
						return 32
					}
				}
				c.num+=n
			}
			if No30{
				if self.Now[_i][_j] == 4 {
					return 32
				}
			}
		}
		lc := len(cs.blocks)
		cs.blocks = append(cs.blocks,c)
		cs.SortCovNum(lc)
		
	}
	if isf < 0 {
		if cs.blocks[0].num < 3 {
			for _j,n := range self.Now[cs.blocks[0].I] {
				if n == 0 {
					continue
				}
				fmt.Println("out 1")
				return cs.blocks[0].I*9+_j
			}
		}
	}
	if self.Inv != -1 {
		var gsTing []*MJGrain
		TL :=0
		for _,c := range cs.blocks {
			for _,b := range c.block {
				for _,j := range b {
					self.Now[c.I][j] --
					gs,l := self.IsTing(0,false)
					lg := len(gs)
					l = lg*4 - l
					self.Now[c.I][j] ++
					if lg >0 { //&& n < k {
						p := &MJGrain{H:c.I,N:j}
						if isf < 0 {
							p.O = gs[0].O
						}else{
							p.O = l
						}
						gsTing = append(gsTing,p)
						SortdGrains(gsTing,TL)
						TL ++
					}
				}
			}
		}
		if TL >0 {
			fmt.Println("out 0")
			if TL == 1 {
				return gsTing[0].H*9+gsTing[0].N
			}
			var IT int
			if isf < 0 {
				IT= gsTing[0].O
				for i,g := range gsTing {
					if g.O != IT {
					//	TL = i
						break
					}
					g.O = self.Now[g.H][g.N]
					SortGrains(gsTing,i)
				}
				IT = gsTing[0].O
				for i,g := range gsTing {
					if g.O != IT {
						break
					}
					g.O = 0
					self.Public.GetCountGrain(g)
					SortdGrains(gsTing,i)
				}
			}else{
				IT= gsTing[0].O
				for i,g := range gsTing {
					if g.O != IT {
						break
					}
					g.O = 0
					self.Public.GetCountGrain(g)
					SortdGrains(gsTing,i)
				}
				IT = gsTing[0].O
				for i,g := range gsTing {
					if g.O != IT {
					//	TL = i
						break
					}
					g.O = self.Now[g.H][g.N]
					SortGrains(gsTing,i)
				}
	
			}
			for _,g := range gsTing {
				if self.Now[g.H][g.N] > 2 {
					continue
				}
				return g.H*9 +g.N
			}
		}
	}
	g :=cs.SameOuts(isf)
	if g != nil{
		fmt.Println("out 3")
		return g.H*9 + g.N
	}
	var dangs [3][]*MJGrain
	for _,_c:= range cs.blocks {
		for _,b := range _c.block {
			for _,_b := range b{
				n :=self.Now[_c.I][_b]
				if  n == 0 {
					continue
				}
				g :=&MJGrain{H:_c.I,N:_b,O:n}
				self.Public.GetCountGrain(g)
				n --
				le := len(dangs[n])
				dangs[n] = append(dangs[n],g)
				SortdGrains(dangs[n],le)
			}
		}
	}
	for _,das := range dangs {
		for _,d:= range das {
			fmt.Println("out 2")
			return d.H*9+d.N
		}
	}
	panic(9)
	return -1
}
func (self Robot) CheckG(gr *MJGrain)  {
	self.Now[gr.H][gr.N]--
	maxw := self.GetWeights()
	gr.O = 0
	for i,no := range self.Now {
		if i == self.Inv {
			continue
		}
		for j,n := range no {
			n = 4-n
			for _,nd := range self.Public.Down {
				n-=nd[i][j]
			}
			for _,nd := range self.Public.See {
				n-=nd[i][j]
			}
			if n <= 0 {
				continue
			}
			self.Now[i][j]++
			if self.GetWeights() < maxw{
				gr.O+= n
			}
			self.Now[i][j]--
		}
	}
}
func (self *Robot) GetDiscard(isNeed bool) int {
	cov:= &Covs{Ro:self}

	var tmparr []int
//	G:
	for i,no := range self.Now {
		c := &Bl{I:i,num:0}
		for j,n := range no {
			if n ==0 {
				if len(tmparr)>0{
					c.block = append(c.block,tmparr)
					tmparr = nil
				}
				continue
			}
			if n >2 {
				c.num +=2
			}
			c.num += n
			c.ext ++
			tmparr = append(tmparr,j)
		}
		if len(tmparr)>0{
			c.block = append(c.block,tmparr)
			tmparr = nil
		}
//		if c.num < 3 {
//			return i
//		}
		cov.blocks = append(cov.blocks,c)
		cov.SortCovNum(i)
	}
	if !isNeed {
		if cov.blocks[0].num >3 {
			return -1
		}
	}
	if cov.blocks[0].num == cov.blocks[1].num {
		if cov.blocks[0].ext < cov.blocks[1].ext {
			return cov.blocks[1].I
		}
		if cov.blocks[0].ext == cov.blocks[1].ext {
			if len(cov.blocks[0].block) < len(cov.blocks[1].block) {
				return cov.blocks[1].I
			}
		}
	}
	return cov.blocks[0].I
}
func (self *Robot) OutDiscard(num int) (o []int) {
	self.Inv =self.GetDiscard(true)
	if self.Inv != -1 {
		for _j,n := range self.Now[self.Inv] {
			if n != 0 {
				for i :=0;i<n;i++{
					self.Now[self.Inv][_j]--
					o = append(o,self.Inv*9+_j)
					if len(o) == num {
						break
					}
				}
			}
		}
	}
	k := num - len(o)
	if k >0 {
		for i:=0;i<k;i++{
			out :=self.Outs(-1,false)
			o= append(o,out)
			self.Now[out/9][out%9]--
		}
	}
	return o
}
func (self Robot) IsTing(isF int,isMax bool) ( gs []*MJGrain,sum int) {
	sum = 0
	I :=0
	for i,no := range self.Now {
		if i == self.Inv {
			for _,n := range no {
				if n >0 {
					return nil,0
				}
			}
			continue
		}
		for j,_ := range no {
			self.Now[i][j]++
			isf := self.CheckFull(isMax)
			self.Now[i][j]--
			if isf>0 && isf+isF >= 0 {
				g :=&MJGrain{H:i,N:j,O:self.Now[i][j],S:isf}
				self.Public.GetCountGrain(g)
				sum+=g.O
				g.O = isf
				gs = append(gs,g)
				SortdGrains(gs,I)
				I++
			}
		}
	}
	return gs,sum
}
func GetWeightVal(ns []int,tmpIs []int) int {
//	rn1 := make([]int,len(ns))
//	copy(rn1,ns)
	_w1:= AddWeight(ns,tmpIs,0)
	if _w1 == 0{
		return _w1
	}
	l := len(tmpIs)
	if l <3 {
		return _w1
	}
	rs := make([]int,l)
	for i,_s := range tmpIs {
		rs[l-i-1] = _s
	}
//	rn2 := make([]int,len(ns))
//	copy(rn2,ns)
	_w2 := AddWeight(ns,rs,0)
	if _w2 < _w1 {
		return _w2
	}
	return _w1

}
func (self Robot) CheckFullTest() int {
	var weight int = 0
	var tmpIs []int
	for i,no := range self.Now {
		if i == self.Inv {
			continue
		}
		for j,n := range no {
			if n == 0 {
				if len(tmpIs) != 0 {
					weight += GetWeightVal(no[0:],tmpIs)
					tmpIs = nil
				}
				continue
			}else {
				tmpIs = append(tmpIs,j)
			}
		}
		if len(tmpIs) >0 {
			weight += GetWeightVal(no[0:],tmpIs)
			tmpIs = nil
		}
	}
	return weight
//	fmt.Println(weight)
}
func (self Robot) CheckFull(Max bool) int {

	var weight int = 0
	var  dui int
	var tmpIs []int
	qing :=0
	dan :=0
	for i,no := range self.Now {
		if i == self.Inv {
			continue
		}
		q :=0
		for j,n := range no {
			if n == 0 {
				if len(tmpIs) != 0 {
					weight += GetWeightVal(no[0:],tmpIs)
					tmpIs = nil
				}
				continue
			}else {
				q++
				tmpIs = append(tmpIs,j)
				if n == 2 {
					dui++
				}else if n == 4 {
					dui+=2
				}else if n == 1 {
					dan++
				}
			}
		}
		if q >0 {
			qing++
			if len(tmpIs) >0 {
				weight += GetWeightVal(no[0:],tmpIs)
				tmpIs = nil
			}
		}
	}

	f := 0
	if weight == 1 {
		f=1
		if dan == 0 {
			f+=1
		}
		if qing == 1 {
			f*=4
		}
		if Max {
			f*=2
		}
		if self.Public != nil {
		for i,no := range self.Public.Down[self.Uid] {
			for j,n := range no {
				if n > 2 {
					f ++
				}else if n ==2 {
					if self.Now[i][j] >0 {
						f ++
					}
				}
			}
		}
		}
		for _,no := range self.Now {
			for _,n := range no {
				if n == 0 {
					continue
				}
				if n == 4 {
					f++
				}
//				n+=self.Public.Down[self.Uid][i][j]
			}
		}
		if Max{
//			if self.NowU
			nu :=0
			for _,n := range self.NowU {
				if n >= 0 {
					nu ++
				}
			}
			return f*nu
		}else{
			return f
		}
	}
	if dui == 7 {
		f = 4
		if qing == 1 {
			f*=4
		}
		if Max {
			f*=2
		}
		for _,no := range self.Now {
			for _,n := range no {
				if n == 4 {
					f++
				}
			}
		}
		return f
	}
	return f

}
func (self Robot) Exchange(v int) (out int){
	i := v/9
	if self.Invalid[i]{
		return v
	}
	j := v%9
	self.Now[i][j] ++
	for _j,n := range self.Now[self.Inv] {
		if n != 0 {
			return self.Inv*9+_j
		}
	}
	w1:= self.GetWeights()
	minw := w1
	var tmp []*MJGrain
	for _i,no := range self.Now {
		if _i == self.Inv {
			continue
		}
		for _j,n := range no {
			if n == 0 {
				continue
			}
			self.Now[_i][_j] --
			p := &MJGrain{H:_i,N:_j}
			_w := self.GetWeights()
			if _w < minw {
				tmp = []*MJGrain{p}
				minw = _w
			}else if _w == minw {
				tmp = append(tmp,p)
			}
			self.Now[_i][_j] ++
		}
	}
	le := len(tmp)
	if le ==0 {
		panic("tmp  == 0")
	}
	if le == 1 {
		return tmp[0].H*9+tmp[0].N
	}
	for i,gr := range tmp {
		self.Public.GetCountGrain(gr)
		SortdGrains(tmp,i)
	}
	return tmp[0].H*9+tmp[0].N

}
func GetWeight(ns []int) (int) {
	var n int

	var is [][]int
	var tmpIs []int
	var w int = 0

	for i:=0;i<9;i++ {
		n = ns[i]
		if n == 0 {
			if len(tmpIs) != 0 {
				is = append(is,tmpIs)
				tmpIs = nil
			}
			continue
		}else {
			tmpIs = append(tmpIs,i)
			if n>2 {
				continue
			}
			w+=3-n
		}
	}
	if len(tmpIs) != 0 {
		is = append(is,tmpIs)
		tmpIs = nil
	}

	var w1 int = 0
	for _,s := range is {
		w1 += GetWeightVal(ns,s)
	}
	return w+w1
}
func AddWeight(dan []int,na []int,w int) int {
	da := make([]int,len(dan))
	copy(da,dan)
	le := len(na)
	if le == 0 {
		return w
	}
	if le < 3{
		for _,_i := range na{
			if da[_i] <3 {
				w += (3- da[_i])
			}
		}
		return w
	}
	n :=da[na[0]]
	if n >2 {
		return AddWeight(da,na[1:],w)
	}
	var start int = 3
	isb := false
	for i:=2;i>=0;i--{
		da[na[i]] --
		if da[na[i]] == 0 {
			if !isb {
				isb = true
			}
			continue
		}else {
			if !isb {
				start = i
			}else{
				w += 3-da[na[i]]
			}
		}
	}
//	fmt.Println(start)
	na = na[start:]
	w1:= AddWeight(da,na,w)
	len_na := len(na)
	if len_na < 3 {
		return w1
	}
	rna := make([]int,len_na)
	for i,n := range na {
		rna[len_na-i-1] = n
	}
	w2 := AddWeight(da,rna,w)

	if w1 < w2 {
		return w1
	}else{
		return w2
	}

}
func (self *Robot) GetWeights() (w int){
	w = 0
	for _i,no := range self.Now {
		if self.Inv == _i {
			continue
		}
		w += GetWeight(no[0:])
	}
	return w
}
