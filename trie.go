//  Copyright (c) 2014 ikawaha.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package kagome

// Any type implements Trie interface may be used as a dictionary.
type Trie interface {
	FindString(string) (id int, ok bool)               // search a dictionary by a keyword.
	CommonPrefixSearchString(string) (ids, lens []int) // finds keywords sharing common prefix in a dictionary.
}
