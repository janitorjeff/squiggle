package squiggle

import (
	"testing"
)

var scopesRender = `CREATE TABLE Scopes (
	id INT PRIMARY KEY AUTOINCREMENT
);`

var placeRender = `CREATE TABLE IF NOT EXISTS CommandTTSPlaceSettings (
	place INT PRIMARY KEY,
	FOREIGN KEY (place) REFERENCES Scopes(id) ON DELETE CASCADE,
	subonly BOOLEAN NOT NULL DEFAULT FALSE
);`

var personRender = `CREATE TABLE IF NOT EXISTS CommandTTSPersonSettings (
	person INT NOT NULL,
	place INT NOT NULL,
	voice VARCHAR(255) NOT NULL,
	UNIQUE(person, place)
);`

func testTable(t *testing.T, table *Table, expected string) {
	if r := table.Render(); r != expected {
		t.Fatalf("incorrect table output for scopes got:\n\n%s\n\nexpected:\n\n%s", r, expected)
	}
}

func TestTableRender(t *testing.T) {
	scopes := NewTable("Scopes")
	id := scopes.Int("id").Primary().Auto()

	placeSettings := NewTable("CommandTTSPlaceSettings").IfNotExists().
		Int("place").Primary().Foreign(id).Cascade().Ok().
		Bool("subonly").NotNull().Default(false).Ok()

	personSettings := NewTable("CommandTTSPersonSettings").IfNotExists()
	person := personSettings.Int("person").NotNull()
	place := personSettings.Int("place").NotNull()
	personSettings.Unique(person, place)
	personSettings.VarChar("voice").NotNull()

	testTable(t, scopes, scopesRender)
	testTable(t, placeSettings, placeRender)
	testTable(t, personSettings, personRender)
}
