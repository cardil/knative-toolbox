/*
 Copyright 2024 The Knative Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package strings

// NewSet creates a new Set from given strings.
func NewSet(ins ...string) Set {
	s := Set{}
	s.All(ins)
	return s
}

// Set is basic implementation of string based set.
type Set struct {
	set map[string]bool
}

// Add adds a string to the Set.
func (s *Set) Add(in string) {
	if s.set == nil {
		s.set = make(map[string]bool)
	}
	s.set[in] = true
}

// All adds all elements of a slice to the Set.
func (s *Set) All(ins []string) {
	for _, in := range ins {
		s.Add(in)
	}
}

// Equal checks if given Set is equal to the other given Set.
func (s Set) Equal(other Set) bool {
	if len(s.set) != len(other.set) {
		return false
	}
	for elem := range s.set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

// Contains checks if element exists within the Set.
func (s Set) Contains(elem string) bool {
	_, ok := s.set[elem]
	return ok
}

// Slice coverts a Set to a []string.
func (s Set) Slice() []string {
	sl := make([]string, 0, len(s.set))
	for elem := range s.set {
		sl = append(sl, elem)
	}
	return sl
}

// Len gets the length of the Set.
func (s Set) Len() int {
	return len(s.set)
}
