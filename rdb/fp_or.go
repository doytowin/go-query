/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"reflect"
	"strings"
)

type fpBasicArrayByOr struct {
	fpSuffix FieldProcessor
}

func buildFpBasicArrayByOr(fieldName string) FieldProcessor {
	return &fpBasicArrayByOr{fpSuffix: buildFpSuffix(strings.TrimSuffix(fieldName, "Or"))}
}

func (fp *fpBasicArrayByOr) Process(value reflect.Value) (string, []any) {
	var args, arr []any
	conditions := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		conditions[i], arr = fp.fpSuffix.Process(value.Index(i))
		args = append(args, arr...)
	}
	return fpForOr.connect(conditions), args
}

type fpStructArrayByOr struct {
	fpForAnd FieldProcessor
}

func buildFpStructArrayByOr() FieldProcessor {
	return &fpStructArrayByOr{fpForAnd}
}

func (fp *fpStructArrayByOr) Process(value reflect.Value) (condition string, args []any) {
	conditions := make([]string, value.Len())
	var arr []any
	for i := 0; i < value.Len(); i++ {
		conditions[i], arr = fp.fpForAnd.Process(value.Index(i))
		args = append(args, arr...)
	}
	return fpForOr.connect(conditions), args
}
