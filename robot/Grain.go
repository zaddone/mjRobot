package robot
type MJGrain struct {
	H int // 0 1 2
	N int //0--8

	O int
	S int
}
func SortGrains(gs []*MJGrain,i int){
	if i == 0 {
		return
	}
	I := i-1
	if gs[I].O > gs[i].O {
		gs[I] , gs[i] = gs[i] , gs[I]
		SortGrains(gs,I)
	}

}
func SortdGrains(gs []*MJGrain,i int){
	if i == 0 {
		return
	}
	I := i-1
	if gs[I].O < gs[i].O {
		gs[I] , gs[i] = gs[i] , gs[I]
		SortdGrains(gs,I)
	}

}
type UserPublic struct {
	Down [4][3][9]int
	See [4][3][9]int
}
func (self *UserPublic) GetCountGrain(gr *MJGrain) {
//	gr.O = 0

	for _,d := range self.Down {
		gr.O += d[gr.H][gr.N]
	}
	for _,d := range self.See {
		gr.O += d[gr.H][gr.N]
	}
}
