package gf

import (
	"fmt"
	"gorm.io/gorm/clause"
	"reflect"
)

type Eq = clause.Eq

type Neq = clause.Neq

type Gt = clause.Gt
type Lt = clause.Lt

type Lte = clause.Lte

type Gte = clause.Gte

type Contains clause.Eq

type NContains clause.Eq

type In clause.Eq

type NIn clause.Eq
type Any clause.Eq

type Overlap clause.Eq

func (c Contains) Build(builder clause.Builder) {
	builder.WriteQuoted(c.Column)
	builder.WriteString(" LIKE ")
	builder.AddVar(builder, fmt.Sprintf("%%%s%%", c.Value))
}

func (nc NContains) Build(builder clause.Builder) {
	builder.WriteQuoted(nc.Column)
	builder.WriteString(" NOT LIKE ")
	builder.AddVar(builder, fmt.Sprintf("%%%s%%", nc.Value))
}

func (in In) Build(builder clause.Builder) {
	var values []interface{}
	v := reflect.ValueOf(in.Value)
	for i := 0; i < v.Len(); i++ {
		values = append(values, v.Index(i).Interface())
	}
	c := clause.IN{
		Column: in.Column,
		Values: values,
	}
	c.Build(builder)
}

func (nin NIn) Build(builder clause.Builder) {
	var values []interface{}
	v := reflect.ValueOf(nin.Value)
	for i := 0; i < v.Len(); i++ {
		values = append(values, v.Index(i).Interface())
	}
	c := clause.IN{
		Column: nin.Column,
		Values: values,
	}
	c.NegationBuild(builder)
}

func (any Any) Build(builder clause.Builder) {
	builder.AddVar(builder, any.Value)
	builder.WriteString(" = ANY(")
	builder.WriteQuoted(any.Column)
	builder.WriteString(")")
}

func (overlap Overlap) Build(builder clause.Builder) {
	var values []interface{}
	v := reflect.ValueOf(overlap.Value)
	for i := 0; i < v.Len(); i++ {
		values = append(values, v.Index(i).Interface())
	}
	switch len(values) {
	case 0:
		builder.WriteString(" 1 = 1")
	default:
		builder.WriteQuoted(overlap.Column)
		builder.WriteString(" && ARRAY[")
		builder.AddVar(builder, values...)
		builder.WriteString("]")
	}
}
