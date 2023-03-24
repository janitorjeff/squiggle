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

func (t *Table) Render() string {
	var b strings.Builder

	b.WriteString("CREATE TABLE")

	if t.ifNotExists {
		b.WriteString(" IF NOT EXISTS")
	}

	b.WriteByte(' ')
	b.WriteString(t.name)
	b.WriteString(" (")

	var lines []string

	for _, c := range t.cols {
		lines = append(lines, c.render())
	}

	for _, cols := range t.unique {
		var unique strings.Builder
		unique.WriteString("\n\tUNIQUE(")
		for i, col := range cols {
			unique.WriteString(col.name)
			if len(cols)-1 > i {
				unique.WriteString(", ")
			}
		}
		unique.WriteByte(')')
		lines = append(lines, unique.String())
	}

	b.WriteString(strings.Join(lines, ","))
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

func (c *Col) render() string {
	var b strings.Builder

	b.WriteString("\n\t")
	b.WriteString(c.name)
	b.WriteByte(' ')
	b.WriteString(c.t.Render())

	if c.primary {
		b.WriteString(" PRIMARY KEY")
	}
	if c.autoIncrement {
		b.WriteString(" AUTOINCREMENT")
	}
	if c.notNull {
		b.WriteString(" NOT NULL")
	}
	if c.unique {
		b.WriteString(" UNIQUE")
	}
	if c.defaultVal != nil {
		b.WriteString(" DEFAULT ")
		b.WriteString(c.t.Cast(c.defaultVal))
	}

	if fc := c.foreignCol; fc != nil {
		b.WriteString(",\n\tFOREIGN KEY (")
		b.WriteString(c.name)
		b.WriteString(") REFERENCES ")
		b.WriteString(fc.table.name)
		b.WriteByte('(')
		b.WriteString(fc.name)
		b.WriteByte(')')
		if c.foreignCascade {
			b.WriteString(" ON DELETE CASCADE")
		}
	}

	return b.String()
}
