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

var ClientGrantTypes = newClientGrantTypesTable("public", "client_grant_types", "")

type clientGrantTypesTable struct {
	postgres.Table

	// Columns
	ClientID  postgres.ColumnInteger
	GrantType postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ClientGrantTypesTable struct {
	clientGrantTypesTable

	EXCLUDED clientGrantTypesTable
}

// AS creates new ClientGrantTypesTable with assigned alias
func (a ClientGrantTypesTable) AS(alias string) *ClientGrantTypesTable {
	return newClientGrantTypesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ClientGrantTypesTable with assigned schema name
func (a ClientGrantTypesTable) FromSchema(schemaName string) *ClientGrantTypesTable {
	return newClientGrantTypesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ClientGrantTypesTable with assigned table prefix
func (a ClientGrantTypesTable) WithPrefix(prefix string) *ClientGrantTypesTable {
	return newClientGrantTypesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ClientGrantTypesTable with assigned table suffix
func (a ClientGrantTypesTable) WithSuffix(suffix string) *ClientGrantTypesTable {
	return newClientGrantTypesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newClientGrantTypesTable(schemaName, tableName, alias string) *ClientGrantTypesTable {
	return &ClientGrantTypesTable{
		clientGrantTypesTable: newClientGrantTypesTableImpl(schemaName, tableName, alias),
		EXCLUDED:              newClientGrantTypesTableImpl("", "excluded", ""),
	}
}

func newClientGrantTypesTableImpl(schemaName, tableName, alias string) clientGrantTypesTable {
	var (
		ClientIDColumn  = postgres.IntegerColumn("client_id")
		GrantTypeColumn = postgres.StringColumn("grant_type")
		allColumns      = postgres.ColumnList{ClientIDColumn, GrantTypeColumn}
		mutableColumns  = postgres.ColumnList{}
	)

	return clientGrantTypesTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ClientID:  ClientIDColumn,
		GrantType: GrantTypeColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
