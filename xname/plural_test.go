package xname

import (
	"strings"
	"testing"

	"go.olapie.com/x/xname/internal/plurals"
)

var replacements = map[string]string{
	"star":        "stars",
	"STAR":        "STARS",
	"Star":        "Stars",
	"bus":         "buses",
	"fish":        "fish",
	"mouse":       "mice",
	"query":       "queries",
	"ability":     "abilities",
	"agency":      "agencies",
	"movie":       "movies",
	"archive":     "archives",
	"index":       "indices",
	"wife":        "wives",
	"safe":        "saves",
	"half":        "halves",
	"move":        "moves",
	"salesperson": "salespeople",
	"person":      "people",
	"spokesman":   "spokesmen",
	"man":         "men",
	"woman":       "women",
	"basis":       "bases",
	"diagnosis":   "diagnoses",
	"diagnosis_a": "diagnosis_as",
	"datum":       "data",
	"medium":      "media",
	"stadium":     "stadia",
	"analysis":    "analyses",
	"node_child":  "node_children",
	"child":       "children",
	"experience":  "experiences",
	"day":         "days",
	"comment":     "comments",
	"foobar":      "foobars",
	"newsletter":  "newsletters",
	"old_news":    "old_news",
	"news":        "news",
	"series":      "series",
	"species":     "species",
	"quiz":        "quizzes",
	"perspective": "perspectives",
	"ox":          "oxen",
	"photo":       "photos",
	"buffalo":     "buffaloes",
	"tomato":      "tomatoes",
	"dwarf":       "dwarves",
	"elf":         "elves",
	"information": "information",
	"equipment":   "equipment",
	"criterion":   "criteria",
}

// storage is used to restore the state of the global variables
// on each test execution, to ensure no global state pollution
type storage struct {
	singulars    RegularSlice
	plurals      RegularSlice
	irregulars   IrregularSlice
	uncountables []string
}

var backup = storage{}

func init() {
	AddIrregular("criterion", "criteria")
	copy(backup.singulars, plurals.SingularReplacements)
	copy(backup.plurals, plurals.PluralReplacements)
	copy(backup.irregulars, plurals.IrregularReplacements)
	copy(backup.uncountables, plurals.UncountableReplacements)
}

func restore() {
	copy(plurals.SingularReplacements, backup.singulars)
	copy(plurals.PluralReplacements, backup.plurals)
	copy(plurals.IrregularReplacements, backup.irregulars)
	copy(plurals.UncountableReplacements, backup.uncountables)
}

func TestPlural(t *testing.T) {
	for key, value := range replacements {
		if v := Plural(strings.ToUpper(key)); v != strings.ToUpper(value) {
			t.Errorf("%v's plural should be %v, but got %v", strings.ToUpper(key), strings.ToUpper(value), v)
		}

		if v := Plural(Title(key)); v != Title(value) {
			t.Errorf("%v's plural should be %v, but got %v", Title(key), Title(value), v)
		}

		if v := Plural(key); v != value {
			t.Errorf("%v's plural should be %v, but got %v", key, value, v)
		}
	}
}

func TestSingular(t *testing.T) {
	for key, value := range replacements {
		if v := Singular(strings.ToUpper(value)); v != strings.ToUpper(key) {
			t.Errorf("%v's singular should be %v, but got %v", strings.ToUpper(value), strings.ToUpper(key), v)
		}

		if v := Singular(Title(value)); v != Title(key) {
			t.Errorf("%v's singular should be %v, but got %v", Title(value), strings.Title(key), v)
		}

		if v := Singular(value); v != key {
			t.Errorf("%v's singular should be %v, but got %v", value, key, v)
		}
	}
}

func TestAddPlural(t *testing.T) {
	defer restore()
	ln := len(plurals.PluralReplacements)
	AddPlural("", "")
	if ln+1 != len(plurals.PluralReplacements) {
		t.Errorf("Expected len %d, got %d", ln+1, len(plurals.PluralReplacements))
	}
}

func TestAddSingular(t *testing.T) {
	defer restore()
	ln := len(plurals.SingularReplacements)
	AddSingular("", "")
	if ln+1 != len(plurals.SingularReplacements) {
		t.Errorf("Expected len %d, got %d", ln+1, len(plurals.SingularReplacements))
	}
}

func TestAddIrregular(t *testing.T) {
	defer restore()
	ln := len(plurals.IrregularReplacements)
	AddIrregular("", "")
	if ln+1 != len(plurals.IrregularReplacements) {
		t.Errorf("Expected len %d, got %d", ln+1, len(plurals.IrregularReplacements))
	}
}

func TestAddUncountable(t *testing.T) {
	defer restore()
	ln := len(plurals.UncountableReplacements)
	AddUncountable("", "")
	if ln+2 != len(plurals.UncountableReplacements) {
		t.Errorf("Expected len %d, got %d", ln+2, len(plurals.UncountableReplacements))
	}
}

func TestGetPluralReplacements(t *testing.T) {
	replacements := GetPluralReplacements()
	if len(replacements) != len(plurals.PluralReplacements) {
		t.Errorf("Expected len %d, got %d", len(replacements), len(plurals.PluralReplacements))
	}
}

func TestGetSingularReplacements(t *testing.T) {
	singular := GetSingularReplacements()
	if len(singular) != len(plurals.SingularReplacements) {
		t.Errorf("Expected len %d, got %d", len(singular), len(plurals.SingularReplacements))
	}
}

func TestGetIrregularReplacements(t *testing.T) {
	irregular := GetIrregularReplacements()
	if len(irregular) != len(plurals.IrregularReplacements) {
		t.Errorf("Expected len %d, got %d", len(irregular), len(plurals.IrregularReplacements))
	}
}

func TestGetUncountableReplacements(t *testing.T) {
	uncountables := GetUncountableReplacements()
	if len(uncountables) != len(plurals.UncountableReplacements) {
		t.Errorf("Expected len %d, got %d", len(uncountables), len(plurals.UncountableReplacements))
	}
}

func TestSetPlural(t *testing.T) {
	defer restore()
	SetPluralReplacements(RegularSlice{{}, {}})
	if len(plurals.PluralReplacements) != 2 {
		t.Errorf("Expected len 2, got %d", len(plurals.PluralReplacements))
	}
}

func TestSetSingular(t *testing.T) {
	defer restore()
	SetSingularReplacements(RegularSlice{{}, {}})
	if len(plurals.SingularReplacements) != 2 {
		t.Errorf("Expected len 2, got %d", len(plurals.SingularReplacements))
	}
}

func TestSetIrregular(t *testing.T) {
	defer restore()
	SetIrregularReplacements(IrregularSlice{{}, {}})
	if len(plurals.IrregularReplacements) != 2 {
		t.Errorf("Expected len 2, got %d", len(plurals.IrregularReplacements))
	}
}

func TestSetUncountable(t *testing.T) {
	defer restore()
	SetUncountableReplacements([]string{"", ""})
	if len(plurals.UncountableReplacements) != 2 {
		t.Errorf("Expected len 2, got %d", len(plurals.UncountableReplacements))
	}
}
