package robot
import(
//	"fmt"
)
type Bl struct {
	I int
	block [][]int
	num int
	ext int
	weight int
}
func (self *Bl) Init(no []int) {
	self.block = nil
	var tmparr []int
	for j,n := range no {
		if n == 0 {
			if len(tmparr)> 0 {
				self.block = append(self.block,tmparr)
				tmparr = nil
			}
			continue
		}
		self.num += n
		tmparr = append(tmparr,j)
	}
	if len(tmparr)> 0 {
		self.block = append(self.block,tmparr)
		tmparr = nil
	}
}
func AddWeightOut(dan []int,na []int,w int,out []int) (int,[]int) {
	da := make([]int,len(dan))
	copy(da,dan)
	le := len(na)
	if le == 0 {
		return w,out
	}
	if le < 3{
		for _,_i := range na{
			if da[_i] <3 {
				w += (3- da[_i])
				if da[_i] == 1 {
					out = append(out,_i)
				}
			}
		}
		return w,out
	}
	n :=da[na[0]]
	if n >2 {
		return AddWeightOut(da,na[1:],w,out)
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
				if da[na[i]] == 1 {
					out = append(out,na[i])
				}
			}
		}
	}
//	fmt.Println(start)
	na = na[start:]
	w1,out1:= AddWeightOut(da,na,w,nil)
	len_na := len(na)
	if len_na < 3 {
		return w1,append(out,out1...)
	}
	rna := make([]int,len_na)
	for i,n := range na {
		rna[len_na-i-1] = n
	}
	w2,out2 := AddWeightOut(da,rna,w,nil)

	if w1 < w2 {
		return w1,append(out,out1...)
	}else{
		return w2,append(out,out2...)
	}

}
func GetWeightValOut(ns []int,tmpIs []int) (int,[]int) {
//	rn1 := make([]int,len(ns))
//	copy(rn1,ns)
	_w1,_out1:= AddWeightOut(ns,tmpIs,0,nil)
	if _w1 == 0{
		return _w1,_out1
	}
	l := len(tmpIs)
	if l <3 {
		return _w1,_out1
	}
	rs := make([]int,l)
	for i,_s := range tmpIs {
		rs[l-i-1] = _s
	}
//	rn2 := make([]int,len(ns))
//	copy(rn2,ns)
	_w2,_out2 := AddWeightOut(ns,rs,0,nil)
	if _w2 < _w1 {
		return _w2,_out2
	}
	return _w1,_out1

}
func (self *Bl) GetWeightOut(no []int) (w int,outg []*MJGrain) {
	w =0
	for _,s := range self.block {
		_w,_out := GetWeightValOut(no,s)
		w+= _w
		i :=0
		for _,_s := range s {
			i += no[_s]
		}
		for _,o := range _out {
			outg = append(outg,&MJGrain{H:self.I,N:o,O:i})
		}
	}
	return w,outg
}
func (self Bl) UpdateRow(no []int) *Bl {
	self.block = nil
	var tmparr []int
	for j,n := range no {
		if n == 0 {
			if len(tmparr)> 0 {
				self.block = append(self.block,tmparr)
				tmparr = nil
			}
			continue
		}
		tmparr = append(tmparr,j)
	}
	if len(tmparr)> 0 {
		self.block = append(self.block,tmparr)
		tmparr = nil
	}
	return &self
}

type Covs struct {
	Ro  *Robot
	blocks []*Bl
}
func (self *Covs)SortCovNum(i int){
	if i == 0 {
		return
	}
	I := i-1
	if self.blocks[I].num > self.blocks[i].num {
		self.blocks[I],self.blocks[i] = self.blocks[i],self.blocks[I]
		self.SortCovNum(I)
	}
}
func (self *Covs) SameOuts(isf int) (g *MJGrain) {
	var gs []*MJGrain
	var w []int = make([]int,len(self.blocks))
	for i,c := range self.blocks {
		w[i],_ = c.GetWeightOut(self.Ro.Now[c.I][0:])
	}
	I:=0
	ki :=2
	if isf < 0 {
		ki --
	}
	for _,c := range self.blocks {
		for _,b := range c.block {
			n :=0
			for _,j := range b {
				n +=self.Ro.Now[c.I][j]
			}
			for _,j := range b{
				_n :=self.Ro.Now[c.I][j]
				if _n > ki {
					continue
				}
				g:= &MJGrain{H:c.I,N:j,O:0,S:n}
				self.GetOutsValOther(g,w)
				if g.O >0 {
					gs = append(gs,g)
					SortdGrains(gs,I)
					I++
				}
			}
		}
	}
	if I ==0 {
		return nil
	}
	if I == 1 {
		return gs[0]
	}
	fi := gs[0].O
	for i,g := range gs {
		if g.O != fi {
			I = i
			break
		}
		g.O = 0
		self.Ro.Public.GetCountGrain(g)
		SortdGrains(gs,i)
	}
	if I == 1 {
		return  gs[0]
	}
	fi = gs[0].O
	for i,g := range gs {
		if g.O != fi {
			I = i
			break
		}
		g.O = g.S
		SortGrains(gs,i)
	}
	if I == 1 {
		return  gs[0]
	}
	fi = gs[0].O
	for i,g := range gs {
		if g.O != fi {
			I = i
			break
		}
		g.O = self.Ro.Now[g.H][g.N]
		SortGrains(gs,i)
	}
	return  gs[0]
}
func (self *Covs) SameOuts1(isf int) (g *MJGrain) {
	var gs []*MJGrain
	I := 0
	min:=0
	for _,c := range self.blocks {
		lb := len(c.block)
		min += lb
		if lb == 0 {
			continue
		}
		_,outg := c.GetWeightOut(self.Ro.Now[c.I][0:])
		if len(outg) == 0 {
			continue
		}
	//	n := 0
	//	for _,g := range outg {
	//		n += self.Ro.Now[g.H][g.N]
	//	}
		for _,g :=  range outg {
			n := self.Ro.Now[g.H][g.N]
			if n == 3 {
				isOut := false
				for _,do := range self.Ro.Public.See {
					if do[g.H][g.N] >0 {
						isOut = true
						break
					}
				}
				if !isOut {
					continue
				}
			}
	//		g.S = n
//			n--
	//		g.O = n
	//		self.GetOutsVal(g)
			if isf <0 {
				if n == 1 {
					gs = append(gs,g)
					SortGrains(gs,I)
					I ++
				}
			}else{
				gs = append(gs,g)
				SortGrains(gs,I)
				I ++
			}
		}
	}
	if min < 2 {
		return nil
	}
	if I == 0 {
		return nil
	}
	if I == 1 {
		return gs[0]
	}

	fi := gs[0].O
	for i,g := range gs {
		if g.O != fi {
			I = i
			break
		}
	}
	if I == 1 {
		return  gs[0]
	}
	gs = gs[:I]
	for i,g := range gs {
		g.O = self.Ro.Now[g.H][g.N]
		self.Ro.Public.GetCountGrain(g)
		SortdGrains(gs,i)
	}

	fi = gs[0].O
	for i,g := range gs {
		if g.O != fi {
			I = i
			break
		}
	}
	if I == 1 {
		return  gs[0]
	}
	for i,g := range gs {
		g.O = 0
		self.GetOutsVal(g)
		SortdGrains(gs,i)
	}
	return  gs[0]
//	return nil

}
//func (self *Covs) SameOutsbak() (g *MJGrain) {
//
//	var gs  []*MJGrain = nil
//	le := 0
//	for _,c := range self.blocks {
//	//	c.UpdateRow(self.Ro.Now[c.I])
//		if len(c.block) == 0 {
//			continue
//		}
//		w := c.GetWeight(self.Ro.Now[c.I][0:])
//		for _,b := range c.block{
//			for _,j := range b {
//		//		n = self.Ro.Now[c.I][j]
//				self.Ro.Now[c.I][j]--
//				bc := c.UpdateRow(self.Ro.Now[c.I][0:])
//				w1 := bc.GetWeight(self.Ro.Now[c.I][0:])
//				wd := w - w1
//				if wd > 0  {
//					gs = append(gs,&MJGrain{H:c.I,N:j,O:len(b)})
//					SortGrains(gs,le)
//					le++
//				}
//				self.Ro.Now[c.I][j]++
//			}
//		}
//	}
//	if len(gs) == 0 {
//		return nil
//	}
//	lo := gs[0].O
//	for i,g := range gs {
//		if g.O > lo {
//			break
//		}
//		self.GetOutsVal(g)
//		SortdGrains(gs,i)
//	}
//	return gs[0]
//}
func (self *Covs) GetOutsValOther(g *MJGrain,W []int) {
	self.Ro.Now[g.H][g.N]--
	o := 0
	for i,c := range self.blocks {
		if len(c.block ) == 0 {
			continue
		}
		for j,_n := range self.Ro.Now[c.I]{
			if _n == 4 {
				continue
			}
			_n = 4-_n
			for _,os := range self.Ro.Public.Down{
				_n -= os[c.I][j]
			}
			for _,os := range self.Ro.Public.See{
				_n -= os[c.I][j]
			}
			if _n <= 0 {
				continue
			}
			self.Ro.Now[c.I][j]++
			_c := c.UpdateRow(self.Ro.Now[c.I][0:])
			w1,_ := _c.GetWeightOut(self.Ro.Now[c.I][0:])
		//	if w1 < W[i] {
		//		o++
			dw := W[i] - w1
			if dw >0 {
				o += dw
			}
			self.Ro.Now[c.I][j]--
		}
	}
	self.Ro.Now[g.H][g.N]++
	g.O = o
}
func (self *Covs) GetOutsVal(g *MJGrain) {
	self.Ro.Now[g.H][g.N]--
	o := 0
	for _,c := range self.blocks {
		_c := c.UpdateRow(self.Ro.Now[c.I][0:])
		w :=len(_c.block)
		if w == 0 {
			continue
		}
		for j,_n := range self.Ro.Now[c.I]{
			if _n == 4 {
				continue
			}
			self.Ro.Now[c.I][j]++
			bc := _c.UpdateRow(self.Ro.Now[c.I][0:])
			if w > len(bc.block) {
				_n = 4-_n
				for _,os := range self.Ro.Public.Down{
					_n -= os[c.I][j]
				}
				for _,os := range self.Ro.Public.See{
					_n -= os[c.I][j]
				}
				if _n >0 {
					o += _n
				}
			}
			self.Ro.Now[c.I][j]--
		}
	}
	self.Ro.Now[g.H][g.N]++
	g.O = o
}
