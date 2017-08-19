package ai
import(
	"strconv"
	"fmt"
)
type MJGrain struct {
	H int // 0 1 2
	N int //0--8

	O []int
//	O1 int
}
func SortGrain(grs []*MJGrain,i int,o int) {

	if i == 0 {
		return
	}
	I := i-1
	if grs[i].O[o] < grs[I].O[o] {
		grs[i],grs[I] = grs[I],grs[i]
	}
	SortGrain(grs,I,o)

}
func SortdGrain(grs []*MJGrain,i int,o int) {

	if i == 0 {
		return
	}
	I := i-1
	if grs[i].O[o] > grs[I].O[o] {
		grs[i],grs[I] = grs[I],grs[i]
	}
	SortdGrain(grs,I,o)

}
func SortGrainArrs(grs []*MJGrain,o int,isd... bool ) (gr *MJGrain) {

	gr = grs[0]
	if o >= len(gr.O) || len(grs)== 1 || o >= len(isd) {
		return gr
	}
	if isd[o] {
		for i,_ := range grs {
			SortdGrain(grs,i,o)
		}
	}else{
		for i,_ := range grs {
			SortGrain(grs,i,o)
		}
	}
	if grs[0].O[o] != grs[1].O[o] {
		return grs[0]
	}
	_i :=2
	for i,g := range grs[2:]{
		if g.O[o] != grs[0].O[o] {
			_i += i
			break
		}
	}
	grs = grs[:_i]
	o ++
	return SortGrainArrs(grs,o,isd...)

}
func SortGrainArr(grs []*MJGrain,isd bool,o int) (gr *MJGrain) {

	gr = grs[0]
	if o >= len(gr.O) || len(grs)== 1 {
		return gr
	}
	for i,_ := range grs {
		if isd {
			SortdGrain(grs,i,o)
		}else{
			SortGrain(grs,i,o)
		}
	}
	if grs[0].O[o] != grs[1].O[o] {
		return grs[0]
	}
	_i :=2
	for i,g := range grs[2:]{
		if g.O[o] != grs[0].O[o] {
			_i += i
			break
		}
	}
	grs = grs[:_i]
	o ++
	return SortGrainArr(grs,isd,o)

}


type UserPublic struct {
	Now [3][9]int
	See [3][9]int
	Down [3][9]int
	Inv  int
	IsF  bool
	LastSee int
}
func (self *UserPublic) Init(){
	self.IsF = false
	self.Inv = -1
	self.LastSee = -1
	for i,d := range self.See {
		for j,_ := range d {
			self.See[i][j] = 0
			self.Down[i][j] = 0
			self.Now[i][j] = 0
		}
	}
}
func (self *UserPublic) SetNow( str []string) (yu []int){
	for i,d := range self.Now{
		for j,_ := range d {
			self.Now[i][j] = 0
		}
	}
	for _,s := range str {
		n,err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		if n > -1 && n < 27 {
			self.Now[n/9][n%9]++
		}else{
			yu = append(yu,n)
		}
	}

	return yu
}
type Ai struct {
	Public [4]*UserPublic
	BaseVal int
	NowVal int
	LastPublicId  int
	LastForData [2]string
	Blocks [3]*Block
	BlocksMap map[int]*Block
}
func (self *Ai) Init() {
	for i,p := range self.Public {
		if p == nil {
			p = new(UserPublic)
		}
		p.Init()
		self.Public[i] = p
	}
//	var d [2]string
//	self.LastForData = d
	self.LastPublicId = -1
	self.BaseVal = 0
	self.BlocksMap = make(map[int]*Block)

}
func (self *Ai) SetLastUid(uid int) {
	self.LastPublicId = uid
}
func (self *Ai) InitNow(str []string){
	self.Public[3].SetNow(str)
	for i:=0;i<3;i++{
//	for i,_ := range self.Blocks {
		bl := self.BlocksMap[i]
		if bl == nil {
			bl = new(Block)
			bl.Init(self.Public[3].Now[i][:],i)
			self.Blocks[i] = bl
			self.BlocksMap[i] = bl
			SortBlock(self.Blocks[:],i)
		}else{
	//		bl.Init(self.Public[3].Now[i][0:],i)
			bl.Update()
		}
//		if bl.Num == 0 {
//			continue
//		}
	//	self.Blocks = append(self.Blocks,bl)
	}
//	self.GetDiscard()
	if self.Public[3].Inv == -1 {
		self.Public[3].Inv = self.GetDiscard()
	}
}
func (self *Ai) updateBlocksWeight(add bool,i,j,n int ){

	if add {
		self.Public[3].Now[i][j] += n
	}else {
		self.Public[3].Now[i][j] -= n
	}
	self.BlocksMap[i].Update()

}
func (self *Ai) SetBaseVal(n int) {
	self.BaseVal = n
}
func (self *Ai) GetFullVal( n int) int {
	if n ==0 || self.BaseVal == 0 {
		return 0
	}
	return (n - self.NowVal )/self.BaseVal
}
func (self *Ai) SetLastSee(uid int,n int) {
	self.Public[uid].LastSee = n
	self.LastPublicId = uid
}
func (self *Ai) GetLastSee() int {
	return self.Public[self.LastPublicId].LastSee
}
func (self *Ai) SetLastForData(n int,s string){
	self.LastForData[n] = s
}
func (self *Ai) GetLastForData() [2]string{
	return self.LastForData
}
func (self *Ai) SetNowVal(n int){
	self.NowVal = n
}
func (self *Ai) SetValid(nid int,n int){
	self.Public[nid].Inv = n
}
func (self *Ai) SetPublicDown(nid int,n int,val int){
	self.Public[nid].Down[n/9][n%9] = val
}
func (self *Ai) SetPublicSee(nid int,n int,val int){
	self.Public[nid].Down[n/9][n%9] += val
}
func (self *Ai) GetDiscard() int{

	if self.Blocks[0].Num != self.Blocks[1].Num {
		return self.Blocks[0].I
	}
	if self.Blocks[0].Numi > self.Blocks[1].Numi {
		return self.Blocks[0].I
	}else if self.Blocks[0].Numi < self.Blocks[1].Numi {
		return self.Blocks[1].I
	}
	if len(self.Blocks[0].Bl) > len(self.Blocks[1].Bl) {
		return self.Blocks[0].I
	}else if len(self.Blocks[0].Bl) < len(self.Blocks[1].Bl) {
		return self.Blocks[1].I
	}
	return self.Blocks[0].I

}
func (self *Ai) SelfSee(isf int,cOut int) (int) {
	if cOut == -1 {
		f := self.checkFull(true)
		if f>0 {
			if f >= isf {
				return 33
			}
		}
	}
	for i,no := range self.Public[3].Now {
		if i == self.Public[3].Inv {
			continue
		}
		for j,n := range no {
			if n == 0 {
				continue
			}
			if n == 4 {
				self.Public[3].Down[i][j] = 4
				return 32
			}
			_n := self.Public[3].Down[i][j]
			if _n == 0 {
				continue
			}
			if n > 1 {
				self.Public[3].Down[i][j] = 0
				continue
			}

			gr1 := self.CheckJiao(isf)
			if len(gr1) == 0 {
				self.Public[3].Down[i][j] = 3
				return 32
			}
			self.updateBlocksWeight(false,i,j,n)
			gr2 := self.CheckJiao(isf)
			if len(gr2) >0 {
				self.Public[3].Down[i][j] = 3
				return 32
			}
		}
	}
	return cOut

}
func (self *Ai) CheckAllJiao(isf int) (jiao []*MJGrain) {
	ji :=0
	for _,bl := range self.Blocks {
		for _,b := range bl.Bl{
			for _,j := range b {
			//	n :=bl.No[j]
			//	if  n>2 {
			//		continue
			//	}
				bl.No[j]--
				bl.Update()
				gs := self.CheckJiao(isf)
				bl.No[j]++
				lgs := len(gs)
				if lgs == 0 {
					continue
				}
				n1 :=0
				for _,g:= range gs {
					n1 += 4-self.GetPublicNum(g,true)
				}
				gr := &MJGrain{H:bl.I,N:j,O:[]int{0,0,0,0,0}}
				gr.O[0] = lgs
				gr.O[1] = n1
				gr.O[2] = 0
				if self.GetPublicNum(gr,false) >0 {
					gr.O[2] = 1
				}
				gr.O[3] =gs[0].O[0]
				gr.O[4] = bl.No[j]
				jiao = append(jiao,gr)
				ji++
			}
		}
		bl.Update()
	}
	return jiao
}
func (self *Ai) EasyOuts() (out int) {
	return 0
}
func (self *Ai) Outs(isf int,No30 bool) (out int){
	out = -1
	InvBl := self.BlocksMap[self.Public[3].Inv]
	var gr *MJGrain = nil
	defer func(){
		if gr != nil {
			fmt.Println(self.Public[3].Now,isf)
			self.Public[3].See[gr.H][gr.N]++
			self.Public[3].Now[gr.H][gr.N]--
		}
	}()
	if InvBl.Num >0 {
		G:
		for _,bl := range self.Blocks[0].Bl {
			for _,j := range bl {
				gr = &MJGrain{H:self.Blocks[0].I,N:j}
				out =  self.Blocks[0].I *9 + j
				break G
			}
		}
	}
	if No30 {
		out = self.SelfSee(isf,out)
	}
	if out > -1 {
		fmt.Println("out -1")
		return out
	}
	// jiao
	jiao := self.CheckAllJiao(isf)
	if len(jiao) >0 {
		gr = SortGrainArrs(jiao,0,true,true,true,true,false)
		fmt.Println("out 1")
		return gr.H*9+gr.N
	}
	var grs []*MJGrain = nil
	for _,bl := range self.Blocks {
		if bl.Num ==0 {
			continue
		}
		grs = append(grs,bl.RunFind(self)...)
	}
	if len(grs) > 0 {
		gr = SortGrainArrs(grs,0,true,true,false,false)
		str := ""
		for _,g := range grs {
			str = fmt.Sprintf("%s %d",str,g.H*9+g.N)
		}
		fmt.Println("out 2",str)
		return gr.H*9 +gr.N
	}
	for _,bl := range self.Blocks {
		for _,b := range bl.Bl {
			for _,j := range b {
				g :=&MJGrain{H:bl.I,N:j,O:[]int{bl.No[j],0}}
				g.O[1] = self.GetPublicNum(g,false)
				grs = append(grs,g)
			}
		}
	}
	gr = SortGrainArrs(grs,0,false,true)
	fmt.Println("out 3")
	return gr.H*9 +gr.N


}
func (self *Ai) OutDiscard(num int) (o []int){

	self.Public[3].Inv = self.GetDiscard()
	bl := self.BlocksMap[self.Public[3].Inv]
	for _,bs := range bl.Bl{
		for _,j:= range bs {
			n := self.Public[3].Now[bl.I][j]
			for i:=0;i<n;i++{
				self.Public[3].Now[bl.I][j]--
				num --
				o = append(o,bl.I*9 + j)
			}
		}
	}
	bl.Update()
	for i :=num; i>0; i-- {
		out := self.Outs(-1,false)
		o = append(o,out)
//		self.BlocksMap[]
	}
	return o

}
//func (self *Ai) SeeOut(v int,isf int) int{
//	return 0
//}

func (self *Ai) SeeOuts(v int,key []int,isf int) int{

	i:=v/9
	j := v%9
	ns :=0
	defer func(){
		fmt.Println(v,isf)
		self.Public[3].Down[i][j] = ns
	}()
	ks := [4]int{-1,-1,-1,-1}
	for _,k := range key {
		ks[3-k] = k
	}
	for _,k := range ks {
		if k == 3 {
			self.Public[3].Now[i][j]++
			self.BlocksMap[i].Update()
			f := self.checkFull(false)
			if f == 0 {
				return k
			}
			if f>0 && f >= isf {
				return k
			}
			continue
		}else if k == 2 {
			ns = 3
			return k

			gr1 := self.CheckJiao(isf)
			if len(gr1) == 0 {
				ns = 3
				return k
			}
			self.Public[3].Now[i][j]=0
			self.BlocksMap[i].Update()
			gr2 := self.CheckJiao(isf)
			if len(gr2)>0 {
				ns = 3
				return k
			}
			continue
		}else if k == 1 {
			gr1 := self.CheckJiao(isf)
			if len(gr1) == 0 {
				ns = 2
				return k
			}
			self.Public[3].Now[i][j]=0
			self.BlocksMap[i].Update()
			gr2 := self.CheckAllJiao(isf)
			lg2 := len(gr2)
			if lg2 == 0 {
			//	ns = 2
			//	return k
				continue
			}
			for gi:=0;gi <lg2;gi++{
				SortdGrain(gr2,gi,0)
			}
			if gr2[0].O[0] >= len(gr1) {
				ns = 2
				return k
			}
			for gi:=0;gi < lg2;gi++{
				SortdGrain(gr2,gi,1)
			}
			sum1 :=0
			for _,g := range gr1 {
				sum1 += 4 - self.GetPublicNum(g,true)
			}
			if gr2[0].O[1] >= sum1 {
				ns = 2
				return k
			}
			for gi:=0;gi < lg2;gi++{
				SortdGrain(gr2,gi,3)
			}
			if gr2[0].O[3] >= gr1[0].O[0] {
				ns = 2
				return k
			}
		//	self.Public[3].Now[i][j]=2
		//	if len(gr2) >0{
		//		sum1 :=0
		//		for _,g := range gr1 {
		//			sum1 += 4 - self.GetPublicNum(g,true)
		//		}
		//		sum2 :=0
		//		for _,g := range gr2 {
		//			sum2 += 4 - self.GetPublicNum(g,true)
		//		}
		//		if sum1 <= sum2 || gr2[0].O[0] > gr1[0].O[0] {
		//			ns = 2
		//			return k
		//		}
		//	}
			continue
		}else if k == 0 {
			return 0
		}
	}
	return -1

}
func (self *Ai) CheckJiao(isf int) (gr []*MJGrain) {

	I :=0
	for _,bl := range self.Blocks{
		if bl.I == self.Public[3].Inv {
			if len(bl.Bl) > 0 {
				return nil
			}
			continue
		}
		for i,n := range bl.No {
			if n == 4 {
				continue
			}
			bl.No[i]++
			bl.Update()
		//	self.updateBlocksWeight(true,bl.I,i)
			fn := self.checkFull(false)

			if fn > 0 && fn >= isf {
				g :=&MJGrain{H:bl.I,N:i,O:[]int{fn}}
			//	ns :=self.GetPublicNum(g,true)
			//	if ns == 4 {
			//		continue
			//	}
			//	g.O = append(g.O,ns)
				gr = append(gr,g)
				SortdGrain(gr,I,0)
				I ++
			}
			bl.No[i]--

		}
		bl.Update()
	}
	return gr

}
func (self *Ai) GetPublicNum(gr *MJGrain,isNow bool) (n int) {

	if isNow {
		n = self.Public[3].Now[gr.H][gr.N]

	}else{
		n = 0
	}
	for _,pu := range self.Public {
		n += pu.Down[gr.H][gr.N]
		n += pu.See[gr.H][gr.N]
	}
	return n

}
func (self *Ai) checkFull(Max bool) int {
	w :=0
	dui := 0
	dan := 0
	qing:= 0
	f := 1
	for _,bl := range self.Blocks {
//		bl.Init(self.Public[3].Now[bl.I][0:],bl.I)
		if bl.I == self.Public[3].Inv && len(bl.Bl) > 0 {
			return 0
		}
		w += bl.GetWeight()
		for _,b := range bl.Bl {
			for _,j := range b {
				n :=bl.No[j]
				if n == 2 {
					dui ++
				}else if n == 4 {
					dui += 2
					f += 1
				}else if n == 1 {
					dan ++
				}
			}
		}
		if bl.Num >0 {
			qing++
		}else{
			for _,n := range self.Public[3].Down[bl.I] {
				if n >0 {
					qing++
					break
				}
			}
		}
		
	}
	if w != 1 && dui != 7 {
		return 0
	}
	if qing == 1 || dui ==7 {
		f = 4
	}
	if dan == 0 && dui ==1 {
		f *= 2
	}
	for i,no := range self.Public[3].Down {
		for j,n := range no {
			if n >2 {
				f ++
			}else if n == 2 {
				if self.Public[3].Now[i][j] >0 {
					f ++
				}
			}
		}
	}
	if Max {
		for _,p := range self.Public[:3]{
			if p.IsF {
				f += f
			}
		}
	}
	return f

}
