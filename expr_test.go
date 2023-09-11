package gf

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestExpr(t *testing.T) {
	db, _ := gorm.Open(tests.DummyDialector{}, nil)
	results := []struct {
		Clauses []clause.Interface
		Result  string
		Vars    []interface{}
	}{
		{
			[]clause.Interface{
				clause.Where{
					Exprs: []clause.Expression{
						Eq{Column: clause.PrimaryColumn, Value: "123456"},
						clause.Or(Neq{Column: "name", Value: "gf"}),
					}}},
			"WHERE `users`.`id` = ? OR `name` <> ?", []interface{}{"123456", "gf"},
		},
	}
	for idx, result := range results {
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			var (
				buildNames    []string
				buildNamesMap = map[string]bool{}
				user, _       = schema.Parse(&tests.User{}, &sync.Map{}, db.NamingStrategy)
				stmt          = gorm.Statement{DB: db, Table: user.Table, Schema: user, Clauses: make(map[string]clause.Clause)}
			)
			for _, clause := range result.Clauses {
				if _, ok := buildNamesMap[clause.Name()]; !ok {
					buildNames = append(buildNames, clause.Name())
					buildNamesMap[clause.Name()] = true
				}
				stmt.AddClause(clause)
			}
			stmt.Build(buildNames...)
			if strings.TrimSpace(stmt.SQL.String()) != result.Result {
				t.Errorf("SQL expects %v got %v", result.Result, stmt.SQL.String())
			}
			if !reflect.DeepEqual(stmt.Vars, result.Vars) {
				t.Errorf("Vars expect %+v got %v", stmt.Vars, result.Vars)
			}
		})
	}
}
