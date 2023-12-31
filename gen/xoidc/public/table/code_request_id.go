//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var CodeRequestID = newCodeRequestIDTable("public", "code_request_id", "")

type codeRequestIDTable struct {
	postgres.Table

	// Columns
	Code       postgres.ColumnString
	RequestID  postgres.ColumnString
	CreateTime postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type CodeRequestIDTable struct {
	codeRequestIDTable

	EXCLUDED codeRequestIDTable
}

// AS creates new CodeRequestIDTable with assigned alias
func (a CodeRequestIDTable) AS(alias string) *CodeRequestIDTable {
	return newCodeRequestIDTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new CodeRequestIDTable with assigned schema name
func (a CodeRequestIDTable) FromSchema(schemaName string) *CodeRequestIDTable {
	return newCodeRequestIDTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new CodeRequestIDTable with assigned table prefix
func (a CodeRequestIDTable) WithPrefix(prefix string) *CodeRequestIDTable {
	return newCodeRequestIDTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new CodeRequestIDTable with assigned table suffix
func (a CodeRequestIDTable) WithSuffix(suffix string) *CodeRequestIDTable {
	return newCodeRequestIDTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newCodeRequestIDTable(schemaName, tableName, alias string) *CodeRequestIDTable {
	return &CodeRequestIDTable{
		codeRequestIDTable: newCodeRequestIDTableImpl(schemaName, tableName, alias),
		EXCLUDED:           newCodeRequestIDTableImpl("", "excluded", ""),
	}
}

func newCodeRequestIDTableImpl(schemaName, tableName, alias string) codeRequestIDTable {
	var (
		CodeColumn       = postgres.StringColumn("code")
		RequestIDColumn  = postgres.StringColumn("request_id")
		CreateTimeColumn = postgres.TimestampzColumn("create_time")
		allColumns       = postgres.ColumnList{CodeColumn, RequestIDColumn, CreateTimeColumn}
		mutableColumns   = postgres.ColumnList{CreateTimeColumn}
	)

	return codeRequestIDTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Code:       CodeColumn,
		RequestID:  RequestIDColumn,
		CreateTime: CreateTimeColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
