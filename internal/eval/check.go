// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package eval

import (
	"errors"
	"fmt"
	"strings"
)

//!+Check

// Check check vars
func (v Var) Check(vars map[Var]bool) error {
	s := string(v)
	if strings.HasSuffix(s, "'") && strings.HasPrefix(s, "'") {
		return nil
	}
	if _, ok := vars[v]; ok {
		return nil
	}
	return errors.New(errCodeMissParam)
}

func (literal) Check(vars map[Var]bool) error {
	return nil
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/><≥≤≡≠∩∪%", b.op) {
		return errors.New(errCodeWithoutFun)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (s section) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("∈∉", s.op) {
		return errors.New(errCodeWithoutFun)
	}
	if err := s.x.Check(vars); err != nil {
		return err
	}
	for _, arg := range s.y {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

func (c call) Check(vars map[Var]bool) error {
	// 参数无法固定
	// artily, ok := numParams[c.fn]
	// if !ok {
	// 	return fmt.Errorf("unknown function %q", c.fn)
	// }
	// if len(c.args) != artily {
	// 	return fmt.Errorf("call to %s has %d args, want %d",
	// 		c.fn, len(c.args), artily)
	// }
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

// var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

//!-Check
