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
	. "github.com/doytowin/goooqo/core"
	"reflect"
	"strings"
)

var whereId = " WHERE id = ?"
var emMap = make(map[string]*metadata)

type metadata struct {
	TableName   string
	columnMetas []ColumnMetadata
}

type EntityMetadata[E Entity] struct {
	metadata
	ColStr          string
	fieldsWithoutId []string
	createStr       string
	placeholders    string
	updateStr       string
}

func (em *EntityMetadata[E]) buildArgs(entity E) []any {
	args := make([]any, len(em.fieldsWithoutId))
	rv := reflect.ValueOf(entity)
	for i, col := range em.fieldsWithoutId {
		value := rv.FieldByName(col)
		args[i] = ReadValue(value)
	}
	return args
}

func (em *EntityMetadata[E]) buildSelect(query Query) (string, []any) {
	whereClause, args := BuildWhereClause(query)
	s := "SELECT " + em.ColStr + " FROM " + em.TableName + whereClause
	s += BuildSortClause(query.GetSort())
	if query.NeedPaging() {
		s = BuildPageClause(&s, query.CalcOffset(), query.GetPageSize())
	}
	return s, args
}

func (em *EntityMetadata[E]) buildSelectById() string {
	return "SELECT " + em.ColStr + " FROM " + em.TableName + whereId
}

func (em *EntityMetadata[E]) buildCount(query Query) (string, []any) {
	whereClause, args := BuildWhereClause(query)
	sqlStr := "SELECT count(0) FROM " + em.TableName + whereClause
	return sqlStr, args
}

func (em *EntityMetadata[E]) buildDeleteById() string {
	return "DELETE FROM " + em.TableName + whereId
}

func (em *EntityMetadata[E]) buildDelete(query any) (string, []any) {
	whereClause, args := BuildWhereClause(query)
	sqlStr := "DELETE FROM " + em.TableName + whereClause
	return sqlStr, args
}

func (em *EntityMetadata[E]) buildCreate(entity E) (string, []any) {
	return em.createStr, em.buildArgs(entity)
}

func (em *EntityMetadata[E]) buildCreateMulti(entities []E) (string, []any) {
	args := make([]any, 0, len(entities)*len(em.fieldsWithoutId))
	for _, entity := range entities {
		args = append(args, em.buildArgs(entity)...)
	}
	createStr := em.createStr + strings.Repeat(", "+em.placeholders, len(entities)-1)
	return createStr, args
}

func (em *EntityMetadata[E]) buildUpdate(entity E) (string, []any) {
	args := em.buildArgs(entity)
	args = append(args, entity.GetId())
	return em.updateStr, args
}

func (em *EntityMetadata[E]) buildPatch(entity E, extra int) (string, []any) {
	args := make([]any, 0, len(em.fieldsWithoutId)+extra)
	sqlStr := "UPDATE " + em.TableName + " SET "

	rv := reflect.ValueOf(entity)
	for _, col := range em.fieldsWithoutId {
		value := rv.FieldByName(col)
		v := ReadValue(value)
		if v != nil {
			sqlStr += ConvertToColumnCase(col) + " = ?, "
			args = append(args, v)
		}
	}
	return sqlStr[0 : len(sqlStr)-2], args
}

func (em *EntityMetadata[E]) buildPatchById(entity E) (string, []any) {
	sqlStr, args := em.buildPatch(entity, 1)
	sqlStr = sqlStr + whereId
	args = append(args, entity.GetId())
	return sqlStr, args
}

func (em *EntityMetadata[E]) buildPatchByQuery(entity E, query Query) (string, []any) {
	whereClause, argsQ := BuildWhereClause(query)
	patchClause, argsE := em.buildPatch(entity, len(argsQ))

	args := append(argsE, argsQ...)
	sqlStr := patchClause + whereClause

	return sqlStr, args
}

func buildEntityMetadata[E RdbEntity]() EntityMetadata[E] {
	entity := *new(E)
	entityType := reflect.TypeOf(entity)
	columnMetas := BuildColumnMetas(entityType)

	columns := make([]string, len(columnMetas))
	columnsWithoutId := make([]string, 0, len(columnMetas))
	fieldsWithoutId := make([]string, 0, len(columnMetas))

	for i, md := range columnMetas {
		columns[i] = md.ColumnName
		if !md.IsId {
			fieldsWithoutId = append(fieldsWithoutId, md.Field.Name)
			columnsWithoutId = append(columnsWithoutId, md.ColumnName)
		}
	}

	tableName := entity.GetTableName()

	placeholders := "(?" + strings.Repeat(", ?", len(columnsWithoutId)-1) + ")"
	createStr := "INSERT INTO " + tableName +
		" (" + strings.Join(columnsWithoutId, ", ") + ") " +
		"VALUES " + placeholders

	set := make([]string, len(columnsWithoutId))
	for i, col := range columnsWithoutId {
		set[i] = col + " = ?"
	}
	updateStr := "UPDATE " + tableName + " SET " + strings.Join(set, ", ") + whereId

	emMap[entityType.Name()] = &metadata{
		TableName:   tableName,
		columnMetas: columnMetas,
	}
	return EntityMetadata[E]{
		metadata:        *emMap[entityType.Name()],
		ColStr:          strings.Join(columns, ", "),
		fieldsWithoutId: fieldsWithoutId,
		createStr:       createStr,
		placeholders:    placeholders,
		updateStr:       updateStr,
	}
}
