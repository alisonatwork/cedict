package lookup

import (
	"reflect"
	"strings"
	"testing"
)

func TestBuildReturnsErrOnBadInput(t *testing.T) {
	input := "random"

	_, err := Build(strings.NewReader(input))
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestBuildIgnoresComments(t *testing.T) {
	input := "# just a comment"

	lookup, err := Build(strings.NewReader(input))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(lookup.Simplified) > 0 || len(lookup.Traditional) > 0 {
		t.Errorf("Expected empty lookup table, got %v", lookup)
	}
}

func TestBuildIgnoresEmptyLines(t *testing.T) {
	input := "\n\n\n"

	lookup, err := Build(strings.NewReader(input))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(lookup.Simplified) > 0 || len(lookup.Traditional) > 0 {
		t.Errorf("Expected empty lookup table, got %v", lookup)
	}
}

func TestBuildHandlesDuplicateDefinitions(t *testing.T) {
	input := "森 森 [Sen1] /Mori (Japanese surname)/\n森 森 [sen1] /forest/"
	key := "森"

	lookup, err := Build(strings.NewReader(input))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(lookup.Simplified) != 1 {
		t.Errorf("Expected one entry in Simplified lookup, got %v", lookup.Simplified)
	}
	if len(lookup.Simplified[key]) != 2 {
		t.Errorf("Expected two definitions in Simplified[%s], got %v", key, lookup.Simplified)
	}
}

func TestBuildParsesSimplified(t *testing.T) {
	input := "你好 你好 [ni3 hao3] /hello/hi/"
	expected := Entry{
		Simplified:  "你好",
		Traditional: "你好",
		Pinyin:      "ni3 hao3",
		Definitions: []string{"hello", "hi"},
	}
	key := expected.Simplified

	lookup, err := Build(strings.NewReader(input))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(lookup.Simplified) != 1 {
		t.Errorf("Expected one entry in Simplified lookup, got %v", lookup.Simplified)
	}
	if len(lookup.Simplified[key]) != 1 {
		t.Errorf("Expected one definition in Simplified[%s], got %v", key, lookup.Simplified)
	}
	if !reflect.DeepEqual(*lookup.Simplified[key][0], expected) {
		t.Errorf("Expected Simplified[%s]->%v, got %v", key, expected, *lookup.Simplified[key][0])
	}
}

func TestBuildParsesTraditional(t *testing.T) {
	input := "麵條 面条 [mian4 tiao2] /noodles/"
	expected := Entry{
		Simplified:  "面条",
		Traditional: "麵條",
		Pinyin:      "mian4 tiao2",
		Definitions: []string{"noodles"},
	}
	key := expected.Traditional

	lookup, err := Build(strings.NewReader(input))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(lookup.Traditional) != 1 {
		t.Errorf("Expected one entry in Traditional lookup, got %v", lookup.Traditional)
	}
	if len(lookup.Traditional[key]) != 1 {
		t.Errorf("Expected one definition in Traditional[%s], got %v", key, lookup.Traditional)
	}
	if !reflect.DeepEqual(*lookup.Traditional[key][0], expected) {
		t.Errorf("Expected Traditional[%s]->%v, got %v", key, expected, *lookup.Traditional[key][0])
	}
}
