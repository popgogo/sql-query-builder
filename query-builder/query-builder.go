package querybuilder

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	TableName    string
	Fields       []string
	Conditions   []Condition
	OrConditions []Condition
	Relations    []Relation
	CTEs         []CTE
}

type Condition struct {
	Field    string
	Operator string
	Value    interface{}
}

type Relation struct {
	Table      string
	ForeignKey string
	PrimaryKey string
}

type CTE struct {
	Name       string
	QueryParts *QueryBuilder
}

func NewQueryBuilder(tableName string) *QueryBuilder {
	return &QueryBuilder{
		TableName:    tableName,
		Fields:       []string{},
		Conditions:   []Condition{},
		OrConditions: []Condition{},
		Relations:    []Relation{},
		CTEs:         []CTE{},
	}
}

func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.Fields = append(qb.Fields, fields...)
	return qb
}

func (qb *QueryBuilder) Where(field, operator string, value interface{}) *QueryBuilder {
	qb.Conditions = append(qb.Conditions, Condition{Field: field, Operator: operator, Value: value})
	return qb
}

func (qb *QueryBuilder) OrWhere(field, operator string, value interface{}) *QueryBuilder {
	qb.OrConditions = append(qb.OrConditions, Condition{Field: field, Operator: operator, Value: value})
	return qb
}

func (qb *QueryBuilder) Join(table, foreignKey, primaryKey string) *QueryBuilder {
	qb.Relations = append(qb.Relations, Relation{Table: table, ForeignKey: foreignKey, PrimaryKey: primaryKey})
	return qb
}

func (qb *QueryBuilder) AddCTE(name string, queryBuilder *QueryBuilder) *QueryBuilder {
	qb.CTEs = append(qb.CTEs, CTE{Name: name, QueryParts: queryBuilder})
	return qb
}

func (qb *QueryBuilder) BuildQuery() (string, []interface{}) {
	var query strings.Builder
	var args []interface{}
	placeholderIndex := 1

	if len(qb.CTEs) > 0 {
		cteParts := []string{}
		for _, cte := range qb.CTEs {
			cteQuery, cteArgs := cte.QueryParts.BuildQuery()
			cteParts = append(cteParts, fmt.Sprintf("%s AS (%s)", cte.Name, cteQuery))
			args = append(args, cteArgs...)
		}
		query.WriteString(fmt.Sprintf("WITH %s ", strings.Join(cteParts, ", ")))
	}

	query.WriteString(fmt.Sprintf("SELECT %s FROM %s", strings.Join(qb.Fields, ", "), qb.TableName))

	if len(qb.Relations) > 0 {
		for _, relation := range qb.Relations {
			query.WriteString(fmt.Sprintf(" JOIN %s ON %s.%s = %s.%s", relation.Table, qb.TableName, relation.PrimaryKey, relation.Table, relation.ForeignKey))
		}
	}

	if len(qb.Conditions) > 0 || len(qb.OrConditions) > 0 {
		query.WriteString(" WHERE ")
		conditions := []string{}
		for _, condition := range qb.Conditions {
			conditions = append(conditions, fmt.Sprintf("%s %s $%d", condition.Field, condition.Operator, placeholderIndex))
			args = append(args, condition.Value)
			placeholderIndex++
		}
		if len(qb.OrConditions) > 0 {
			orConditions := []string{}
			for _, condition := range qb.OrConditions {
				orConditions = append(orConditions, fmt.Sprintf("%s %s $%d", condition.Field, condition.Operator, placeholderIndex))
				args = append(args, condition.Value)
				placeholderIndex++
			}
			conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(orConditions, " OR ")))
		}
		query.WriteString(strings.Join(conditions, " AND "))
	}

	return query.String(), args
}
