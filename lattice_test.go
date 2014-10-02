package kagome

import (
	//"fmt"
	"testing"
)

func TestLatticeBuild01(t *testing.T) {
	la := newLattice()
	la.build("となりのトトロ")
	//fmt.Println(la)
	la.forward()
	la.backward()
}
