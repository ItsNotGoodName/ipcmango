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

var TimeSeeds = newTimeSeedsTable("dahua", "time_seeds", "")

type timeSeedsTable struct {
	postgres.Table

	// Columns
	Seed     postgres.ColumnInteger
	CameraID postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type TimeSeedsTable struct {
	timeSeedsTable

	EXCLUDED timeSeedsTable
}

// AS creates new TimeSeedsTable with assigned alias
func (a TimeSeedsTable) AS(alias string) *TimeSeedsTable {
	return newTimeSeedsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new TimeSeedsTable with assigned schema name
func (a TimeSeedsTable) FromSchema(schemaName string) *TimeSeedsTable {
	return newTimeSeedsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new TimeSeedsTable with assigned table prefix
func (a TimeSeedsTable) WithPrefix(prefix string) *TimeSeedsTable {
	return newTimeSeedsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new TimeSeedsTable with assigned table suffix
func (a TimeSeedsTable) WithSuffix(suffix string) *TimeSeedsTable {
	return newTimeSeedsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newTimeSeedsTable(schemaName, tableName, alias string) *TimeSeedsTable {
	return &TimeSeedsTable{
		timeSeedsTable: newTimeSeedsTableImpl(schemaName, tableName, alias),
		EXCLUDED:       newTimeSeedsTableImpl("", "excluded", ""),
	}
}

func newTimeSeedsTableImpl(schemaName, tableName, alias string) timeSeedsTable {
	var (
		SeedColumn     = postgres.IntegerColumn("seed")
		CameraIDColumn = postgres.IntegerColumn("camera_id")
		allColumns     = postgres.ColumnList{SeedColumn, CameraIDColumn}
		mutableColumns = postgres.ColumnList{SeedColumn, CameraIDColumn}
	)

	return timeSeedsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Seed:     SeedColumn,
		CameraID: CameraIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
