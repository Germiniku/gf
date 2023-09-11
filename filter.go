package gf

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
)

const (
	TAG_FILTER = "filter"
	SEP        = ";"
	SEP_COL    = "."
	COL        = "COL"
	REQ        = "REQUIRED"
	ZERO       = "ZERO"
	OPR        = "OPR"
)

type Column struct {
	Value          reflect.Value
	Name           string
	Required       bool
	AllowZeroValue bool
	Operator       string
}

func Filter(p interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tx := db
		joinTables, exprs, err := parse(p)
		if err != nil {
			tx.AddError(fmt.Errorf("parse error:%v", err))
			return tx
		}
		for _, table := range joinTables {
			tx = tx.Joins(table)
		}
		tx = tx.Clauses(exprs...)
		return tx
	}
}

func parse(p interface{}) ([]string, []clause.Expression, error) {
	exprs := []clause.Expression{}
	joinTables := []string{}
	elem := reflect.ValueOf(p).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)
		filter, ok := field.Tag.Lookup(TAG_FILTER)
		if !ok {
			continue
		}
		settings := schema.ParseTagSetting(filter, SEP)
		column := &Column{Value: elem.Field(i)}
		if column.Value.Kind() == reflect.Ptr {
			column.Value = reflect.Indirect(column.Value)
		}
		for key, val := range settings {
			switch key {
			case COL:
				column.Name = val
			case REQ:
				column.Required = true
			case ZERO:
				column.AllowZeroValue = true
			case OPR:
				column.Operator = val
			}
		}
		if column.Name == "" {
			column.Name = strcase.ToSnake(field.Name)
		}
		if column.Operator == "" {
			column.Operator = "eq"
		}
		if !column.Value.IsValid() {
			continue
		}
		if isZero(column.Value) {
			if column.Required {
				return nil, nil, fmt.Errorf("column %s required,empty", column.Name)
			}
			if !column.AllowZeroValue {
				continue
			}
		}
		joinTable := parseJoinTable(column.Name)
		if joinTable != "" {
			joinTables = append(joinTables, joinTable)
		}
		expr, err := buildClause(column.Name, column.Operator, column.Value.Interface())
		if err != nil {
			return nil, nil, err
		}
		exprs = append(exprs, expr)
	}
	return joinTables, exprs, nil
}

func parseJoinTable(col string) string {
	parts := strings.Split(col, SEP_COL)
	if len(parts) == 1 {
		return ""
	}
	return parts[1]
}

func buildClause(column, operator string, value interface{}) (clause.Expression, error) {
	switch operator {
	case "eq":
		return Eq{Column: column, Value: value}, nil
	case "neq":
		return Neq{Column: column, Value: value}, nil
	case "gt":
		return Gt{Column: column, Value: value}, nil
	case "gte":
		return Gte{Column: column, Value: value}, nil
	case "lt":
		return Lt{Column: column, Value: value}, nil
	case "lte":
		return Lte{Column: column, Value: value}, nil
	case "in":
		return In{Column: column, Value: value}, nil
	case "!in":
		return NIn{Column: column, Value: value}, nil
	case "like", "contains":
		return Contains{Column: column, Value: value}, nil
	case "!like", "!contains":
		return NContains{Column: column, Value: value}, nil
	case "any":
		return Any{Column: column, Value: value}, nil
	case "overlap":
		return Overlap{Column: column, Value: value}, nil
	}
	return nil, fmt.Errorf("operator %s is not support", operator)
}

func isZero(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Uint8:
		return false
	default:
		return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
	}
}
