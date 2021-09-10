// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gnark DO NOT EDIT

package cs

import (
	"errors"
	"fmt"

	"github.com/consensys/gnark/backend/hint"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// ErrUnsatisfiedConstraint can be generated when solving a R1CS
var ErrUnsatisfiedConstraint = errors.New("constraint is not satisfied")

type hintFunction func(input []fr.Element) fr.Element

type solution struct {
	values          []fr.Element
	solved          []bool
	nbSolved        int
	mHintsFunctions map[hint.ID]hintFunction
}

func newSolution(nbWires int, hintFunctions []hint.Function) (solution, error) {
	s := solution{
		values:          make([]fr.Element, nbWires),
		solved:          make([]bool, nbWires),
		mHintsFunctions: make(map[hint.ID]hintFunction, len(hintFunctions)+2),
	}

	s.mHintsFunctions = make(map[hint.ID]hintFunction, len(hintFunctions)+2)
	s.mHintsFunctions[hint.IsZero] = powModulusMinusOne
	s.mHintsFunctions[hint.IthBit] = ithBit

	for i := 0; i < len(hintFunctions); i++ {
		if _, ok := s.mHintsFunctions[hintFunctions[i].ID]; ok {
			return solution{}, fmt.Errorf("duplicate hint function with id %d", uint32(hintFunctions[i].ID))
		}
		f, ok := hintFunctions[i].F.(hintFunction)
		if !ok {
			return solution{}, fmt.Errorf("invalid hint function signature with id %d", uint32(hintFunctions[i].ID))
		}
		s.mHintsFunctions[hintFunctions[i].ID] = f
	}

	return s, nil
}

func (s *solution) set(id int, value fr.Element) {
	if s.solved[id] {
		panic("solving the same wire twice should never happen.")
	}
	s.values[id] = value
	s.solved[id] = true
	s.nbSolved++
}

func (s *solution) isValid() bool {
	return s.nbSolved == len(s.values)
}
