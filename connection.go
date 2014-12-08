//  Copyright (c) 2014 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

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
