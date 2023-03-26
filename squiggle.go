package squiggle

import (
	"strings"
)

type Col struct {
	name string
	t    SQLType

	defaultVal any

	primary       bool
	autoIncrement bool
	notNull       bool
	unique        bool

	foreignCol     *Col
	foreignCascade bool

	table *Table
}

type Table struct {
	name        string
	ifNotExists bool
	cols        []*Col
	unique      [][]*Col
}

func NewTable(name string) *Table {
	return &Table{
		name: name,
	}
}

func (t *Table) IfNotExists() *Table {
	t.ifNotExists = true
	return t
}

func (t *Table) Unique(cols ...*Col) *Table {
	for _, col := range cols {
		if col.table != t {
			panic("wrong table")
		}
	}
	t.unique = append(t.unique, cols)
	return t
}

func (t *Table) col(name string, sqlT SQLType) *Col {
	for _, col := range t.cols {
		if col.name == name {
			panic("duplicate column name '" + name + "'")
		}
	}
	col := &Col{
		name:  name,
		t:     sqlT,
		table: t,
	}
	t.cols = append(t.cols, col)
	return col
}

func (t *Table) Int(name string) *Col {
	return t.col(name, INT)
}

func (t *Table) Integer(name string) *Col {
	return t.col(name, INTEGER)
}

func (t *Table) Bool(name string) *Col {
	return t.col(name, BOOLEAN)
}

func (t *Table) VarChar(name string) *Col {
	return t.col(name, VARCHAR255)
}

const (
	renderCreateTable = "CREATE TABLE"
	renderIfNotExists = " IF NOT EXISTS"
	renderUniqueCols  = "\n\tUNIQUE("
)

func (t *Table) alloc() int {
	// 1 for space, 2 for space paren, 3 for newline closing paren semicolon
	// which are present at the end
	allocLen := len(renderCreateTable) + 1 + 2 + len(t.name) + 3
	if t.ifNotExists {
		allocLen += len(renderIfNotExists)
	}
	for _, c := range t.cols {
		// 1 for comma, which is almost always present
		allocLen += c.alloc() + 1
	}
	for _, cols := range t.unique {
		// 1 for closing paren
		allocLen += len(renderUniqueCols) + 1
		for _, col := range cols {
			// 2 for comma space, which is likely present
			allocLen += len(col.name) + 2
		}
	}
	return allocLen
}

func (t *Table) Render() string {
	var b strings.Builder

	b.Grow(t.alloc())

	b.WriteString(renderCreateTable)

	if t.ifNotExists {
		b.WriteString(renderIfNotExists)
	}

	b.WriteByte(' ')
	b.WriteString(t.name)
	b.WriteString(" (")

	for i, c := range t.cols {
		b.WriteString(c.render())
		if i != len(t.cols)-1 || len(t.unique) != 0 {
			b.WriteByte(',')
		}
	}

	for _, cols := range t.unique {
		b.WriteString(renderUniqueCols)
		for i, col := range cols {
			b.WriteString(col.name)
			if len(cols)-1 > i {
				b.WriteString(", ")
			}
		}
		b.WriteByte(')')
	}

	b.WriteString("\n);")

	return b.String()
}

func (c *Col) Primary() *Col {
	c.primary = true
	return c
}

func (c *Col) Auto() *Col {
	c.autoIncrement = true
	return c
}

func (c *Col) NotNull() *Col {
	c.notNull = true
	return c
}

func (c *Col) Unique() *Col {
	c.unique = true
	return c
}

func (c *Col) Default(val any) *Col {
	c.defaultVal = val
	return c
}

func (c *Col) Foreign(fc *Col) *Col {
	c.foreignCol = fc
	return c
}

func (c *Col) Cascade() *Col {
	if c.foreignCol == nil {
		panic("no foreign column")
	}
	c.foreignCascade = true
	return c
}

func (c *Col) Ok() *Table {
	return c.table
}

const (
	renderPrimaryKey      = " PRIMARY KEY"
	renderAutoIncrement   = " AUTOINCREMENT"
	renderNotNull         = " NOT NULL"
	renderUnique          = " UNIQUE"
	renderDefault         = " DEFAULT "
	renderForeignKey      = ",\n\tFOREIGN KEY ("
	renderReferences      = ") REFERENCES "
	renderOnDeleteCascade = " ON DELETE CASCADE"
)

func (c *Col) alloc() int {
	// 6 is the type's length, arbitralily picked
	allocLen := 2 + len(c.name) + 1 + 6
	if c.primary {
		allocLen += len(renderPrimaryKey)
	}
	if c.autoIncrement {
		allocLen += len(renderAutoIncrement)
	}
	if c.notNull {
		allocLen += len(renderNotNull)
	}
	if c.unique {
		allocLen += len(renderUnique)
	}
	if c.defaultVal != nil {
		// The 5 is an abitrarily selected number that seemed sane and
		// represents the actual default value
		allocLen += len(renderDefault) + 5
	}
	if c.foreignCol != nil {
		allocLen += len(renderForeignKey) +
			len(renderReferences) +
			2 + // parens
			len(c.name) +
			len(c.foreignCol.table.name) +
			len(c.foreignCol.name)

		if c.foreignCascade {
			allocLen += len(renderOnDeleteCascade)
		}
	}

	return allocLen
}

func (c *Col) render() string {
	var b strings.Builder

	b.Grow(c.alloc())

	b.WriteString("\n\t")
	b.WriteString(c.name)
	b.WriteByte(' ')
	b.WriteString(c.t.Render())

	if c.primary {
		b.WriteString(renderPrimaryKey)
	}
	if c.autoIncrement {
		b.WriteString(renderAutoIncrement)
	}
	if c.notNull {
		b.WriteString(renderNotNull)
	}
	if c.unique {
		b.WriteString(renderUnique)
	}
	if c.defaultVal != nil {
		b.WriteString(renderDefault)
		b.WriteString(c.t.Cast(c.defaultVal))
	}

	if fc := c.foreignCol; fc != nil {
		b.WriteString(renderForeignKey)
		b.WriteString(c.name)
		b.WriteString(renderReferences)
		b.WriteString(fc.table.name)
		b.WriteByte('(')
		b.WriteString(fc.name)
		b.WriteByte(')')
		if c.foreignCascade {
			b.WriteString(renderOnDeleteCascade)
		}
	}

	return b.String()
}
