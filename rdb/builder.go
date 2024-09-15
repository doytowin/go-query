/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
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

type QueryBuilder interface {
	BuildConditions() ([]string, []any)
}

func isValidValue(value reflect.Value) bool {
	return !value.IsNil()
}

func BuildWhereClause(query any) (string, []any) {
	conditions, args := buildConditions(query)
	if len(conditions) == 0 {
		return "", []any{}
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildConditions(query any) ([]string, []any) {
	if qb, ok := query.(QueryBuilder); ok {
		return qb.BuildConditions()
	}
	rtype := reflect.TypeOf(query)
	rvalue := reflect.ValueOf(query)
	if rtype.Kind() == reflect.Pointer {
		rtype = rtype.Elem()
		rvalue = rvalue.Elem()
	}
	args := make([]any, 0, rtype.NumField())
	conditions := make([]string, 0, rtype.NumField())

	registerFpByType(rtype)
	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		fpKey := buildFpKey(rtype, field)
		processor := fpMap[fpKey]
		if processor != nil {
			value := rvalue.FieldByName(field.Name)
			if isValidValue(value) {
				condition, arr := processor.Process(value.Elem())
				conditions = append(conditions, condition)
				args = append(args, arr...)
			}
		}
	}
	return conditions, args
}
