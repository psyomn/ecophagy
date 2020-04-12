package main

import (
	"strings"
	"testing"
)

func TestMinifyGenrePyDict(t *testing.T) {
	line := `[{'id': 123, 'name': 'potato'}, {'id': 321, 'name': 'patata'}]`
	actual := minifyNamePyDict(line)
	expected := `potato,patata`

	if actual != expected {
		t.Log("expected: ", expected)
		t.Log("actual: ", actual)
		t.Fatal()
	}
}

func TestManyNamesPyDict(t *testing.T) {
	data := `[{'name': 'Procirep', 'id': 311},
                  {'name': 'Constellation Productions', 'id': 590},
                  {'name': 'France 3 Cinéma', 'id': 591},
                  {'name': 'Claudie Ossard Productions', 'id': 592},
                  {'name': 'Eurimages', 'id': 850},
                  {'name': 'MEDIA Programme of the European Union', 'id': 851},
                  {'name': 'Cofimage 5', 'id': 1871},
                  {'name': 'Televisión Española (TVE)', 'id': 6639},
                  {'name': 'Tele München Fernseh Produktionsgesellschaft (TMG)', 'id': 7237},
                  {'name': "Club d'Investissement Média", 'id': 8170},
                  {'name': 'Canal+ España', 'id': 9335},
                  {'name': 'Elías Querejeta Producciones Cinematográficas S.L.', 'id': 12009},
                  {'name': 'Centre National de la Cinématographie (CNC)', 'id': 18367},
                  {'name': 'Victoires Productions', 'id': 25020},
                  {'name': 'Constellation', 'id': 25021},
                  {'name': 'Lumière Pictures', 'id': 25129},
                  {'name': 'Canal+', 'id': 47532},
                  {'name': 'Studio Image', 'id': 68517},
                  {'name': 'Cofimage 4', 'id': 79437},
                  {'name': 'Ossane', 'id': 79438},
                  {'name': 'Phoenix Images', 'id': 79439}]`

	names := minifyNamePyDict(data)

	if len(strings.Split(names, ",")) != 21 {
		t.Fatal("wrong number of entries")
	}
}

func TestEntryParserSimple(t *testing.T) {
	type tc struct {
		json     string
		expected string
		cursor   int
	}

	tests := []tc{
		{`[{'name': "Loew's Incorporated", 'id': 31892}]`, `Loew's Incorporated`, 30},
		{`[{'name': "potatoland", 'id': 31892}]`, `potatoland`, 21},
		{`[{'name': ""I can't believe "you" write like this potatoland"", 'id': 31892}]`,
			`I can't believe "you" write like this potatoland`, 61},
		{`[{'name': 'potatoland', 'id': 31892}]`, `potatoland`, 21},
		{`[{'id': 31892, 'name': "Loew's Incorporated", 'nope': 'nope'}]`, `Loew's Incorporated`, 43},
	}

	for _, test := range tests {
		expected := test.expected
		expectedCursor := test.cursor
		actual, index := parseEntry(test.json)

		if actual != expected {
			t.Fatalf("expected: %v, actual: %v", expected, actual)
		}

		if index != expectedCursor {
			t.Fatalf("expected: %v, was: %v", index, expectedCursor)
		}
	}
}
