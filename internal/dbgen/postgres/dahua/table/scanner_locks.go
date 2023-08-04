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

var ScannerLocks = newScannerLocksTable("dahua", "scanner_locks", "")

type scannerLocksTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnInteger
	UUID      postgres.ColumnString
	CreatedAt postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ScannerLocksTable struct {
	scannerLocksTable

	EXCLUDED scannerLocksTable
}

// AS creates new ScannerLocksTable with assigned alias
func (a ScannerLocksTable) AS(alias string) *ScannerLocksTable {
	return newScannerLocksTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ScannerLocksTable with assigned schema name
func (a ScannerLocksTable) FromSchema(schemaName string) *ScannerLocksTable {
	return newScannerLocksTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ScannerLocksTable with assigned table prefix
func (a ScannerLocksTable) WithPrefix(prefix string) *ScannerLocksTable {
	return newScannerLocksTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ScannerLocksTable with assigned table suffix
func (a ScannerLocksTable) WithSuffix(suffix string) *ScannerLocksTable {
	return newScannerLocksTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newScannerLocksTable(schemaName, tableName, alias string) *ScannerLocksTable {
	return &ScannerLocksTable{
		scannerLocksTable: newScannerLocksTableImpl(schemaName, tableName, alias),
		EXCLUDED:          newScannerLocksTableImpl("", "excluded", ""),
	}
}

func newScannerLocksTableImpl(schemaName, tableName, alias string) scannerLocksTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		UUIDColumn      = postgres.StringColumn("uuid")
		CreatedAtColumn = postgres.TimestampzColumn("created_at")
		allColumns      = postgres.ColumnList{IDColumn, UUIDColumn, CreatedAtColumn}
		mutableColumns  = postgres.ColumnList{IDColumn, UUIDColumn, CreatedAtColumn}
	)

	return scannerLocksTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		UUID:      UUIDColumn,
		CreatedAt: CreatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
