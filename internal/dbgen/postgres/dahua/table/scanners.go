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

var Scanners = newScannersTable("dahua", "scanners", "")

type scannersTable struct {
	postgres.Table

	// Columns
	ID           postgres.ColumnInteger
	FullComplete postgres.ColumnBool
	FullCursor   postgres.ColumnTimestampz
	FullEpoch    postgres.ColumnTimestampz
	QuickCursor  postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ScannersTable struct {
	scannersTable

	EXCLUDED scannersTable
}

// AS creates new ScannersTable with assigned alias
func (a ScannersTable) AS(alias string) *ScannersTable {
	return newScannersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ScannersTable with assigned schema name
func (a ScannersTable) FromSchema(schemaName string) *ScannersTable {
	return newScannersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ScannersTable with assigned table prefix
func (a ScannersTable) WithPrefix(prefix string) *ScannersTable {
	return newScannersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ScannersTable with assigned table suffix
func (a ScannersTable) WithSuffix(suffix string) *ScannersTable {
	return newScannersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newScannersTable(schemaName, tableName, alias string) *ScannersTable {
	return &ScannersTable{
		scannersTable: newScannersTableImpl(schemaName, tableName, alias),
		EXCLUDED:      newScannersTableImpl("", "excluded", ""),
	}
}

func newScannersTableImpl(schemaName, tableName, alias string) scannersTable {
	var (
		IDColumn           = postgres.IntegerColumn("id")
		FullCompleteColumn = postgres.BoolColumn("full_complete")
		FullCursorColumn   = postgres.TimestampzColumn("full_cursor")
		FullEpochColumn    = postgres.TimestampzColumn("full_epoch")
		QuickCursorColumn  = postgres.TimestampzColumn("quick_cursor")
		allColumns         = postgres.ColumnList{IDColumn, FullCompleteColumn, FullCursorColumn, FullEpochColumn, QuickCursorColumn}
		mutableColumns     = postgres.ColumnList{IDColumn, FullCursorColumn, FullEpochColumn, QuickCursorColumn}
	)

	return scannersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:           IDColumn,
		FullComplete: FullCompleteColumn,
		FullCursor:   FullCursorColumn,
		FullEpoch:    FullEpochColumn,
		QuickCursor:  QuickCursorColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
