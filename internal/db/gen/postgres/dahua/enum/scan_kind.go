//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/v2/postgres"

var ScanKind = &struct {
	Full   postgres.StringExpression
	Manual postgres.StringExpression
}{
	Full:   postgres.NewEnumValue("full"),
	Manual: postgres.NewEnumValue("manual"),
}
