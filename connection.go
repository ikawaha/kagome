package kagome

type ConnectionTable struct {
	Row, Col int
	Vec      []int16
}

func (ct *ConnectionTable) At(row, col int) int16 {
	return ct.Vec[ct.Col*row+col]
}
