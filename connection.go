package kagome

// ConnectionTable represents a connection matrix of morphs.
type ConnectionTable struct {
	Row, Col int
	Vec      []int16
}

// At returns the connection cost of matrix[row, col].
func (ct *ConnectionTable) At(row, col int) int16 {
	return ct.Vec[ct.Col*row+col]
}
