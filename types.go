package squiggle

import (
	"fmt"
)

type SQLType int

// List of types taken from: https://www.sqlite.org/datatype3.html
const (
	BOOLEAN SQLType = iota

	INT
	INTEGER
	TINYINT
	SMALLINT
	MEDIUMINT
	BIGINT
	UNSIGNEDBIGINT
	INT2
	INT8

	CHARACTER20
	VARCHAR255
	VARYINGCHARACTER255
	NCHAR55
	NATIVECHARACTER70
	NVARCHAR100
	TEXT
	CLOB

	BLOB

	REAL
	DOUBLE
	DOUBLEPRECISION
	FLOAT

	NUMERIC
	DECIMAL105
	DATE
	DATETIME
)

func (t SQLType) Render() string {
	switch t {
	case BOOLEAN:
		return "BOOLEAN"

	case INT:
		return "INT"
	case INTEGER:
		return "INTEGER"
	case TINYINT:
		return "TINYINT"
	case SMALLINT:
		return "SMALLINT"
	case MEDIUMINT:
		return "MEDIUMINT"
	case BIGINT:
		return "BIGINT"
	case UNSIGNEDBIGINT:
		return "UNSIGNED BIG INT"
	case INT2:
		return "INT2"
	case INT8:
		return "INT8"

	case CHARACTER20:
		return "CHARACTER(20)"
	case VARCHAR255:
		return "VARCHAR(255)"
	case VARYINGCHARACTER255:
		return "VARYING CHARACTER(255)"
	case NCHAR55:
		return "NCHAR(55)"
	case NATIVECHARACTER70:
		return "NATIVE CHARACTER(70)"
	case NVARCHAR100:
		return "NVARCHAR(100)"
	case TEXT:
		return "TEXT"
	case CLOB:
		return "CLOB"

	case BLOB:
		return "BLOB"

	case REAL:
		return "REAL"
	case DOUBLE:
		return "DOUBLE"
	case DOUBLEPRECISION:
		return "DOUBLE PRECISION"
	case FLOAT:
		return "FLOAT"

	case NUMERIC:
		return "NUMERIC"
	case DECIMAL105:
		return "DECIMAL(10,5)"
	case DATE:
		return "DATE"
	case DATETIME:
		return "DATETIME"
	default:
		panic("unknown type")
	}
}

func (t SQLType) Cast(val any) string {
	switch t {
	case BOOLEAN:
		switch val.(type) {
		case bool:
			if val.(bool) {
				return "TRUE"
			}
			return "FALSE"
		default:
			panic("only bool accepted")
		}

	case INT, INTEGER, TINYINT, SMALLINT, MEDIUMINT, BIGINT, INT2, INT8, NUMERIC:
		switch val.(type) {
		case int, int8, int16, int32, int64:
			return fmt.Sprint(val)
		default:
			panic("only int types accepted")
		}

	case UNSIGNEDBIGINT:
		switch val.(type) {
		case uint, uint8, uint16, uint32, uint64:
			return fmt.Sprint(val)
		default:
			panic("only uint types accepted")
		}

	case CHARACTER20, VARCHAR255, VARYINGCHARACTER255, NCHAR55, NATIVECHARACTER70, NVARCHAR100, TEXT, CLOB:
		switch val.(type) {
		case string:
			return val.(string)
		default:
			panic("only string type accepted")
		}

	case BLOB:
		switch val.(type) {
		case []byte:
			// TODO
			return "BLOB"
		default:
			panic("only byte slice accepted")
		}

	case REAL, DOUBLE, DOUBLEPRECISION, FLOAT, DECIMAL105:
		switch val.(type) {
		case float32, float64:
			return fmt.Sprint(val)
		default:
			panic("only float types accepted")
		}

	case DATE:
		// TODO
		fallthrough
	case DATETIME:
		// TODO
		fallthrough

	default:
		panic("unknown type")
	}
}
