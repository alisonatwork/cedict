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

func TestMatchMatchesSimplified(t *testing.T) {
	input := "椀 碗 [wan3] /variant of 碗[wan3]/\n碗 碗 [wan3] /bowl/cup/CL:隻|只[zhi1],個|个[ge4]/\n麵條 面条 [mian4 tiao2] /noodles/"

	lookup, err := Build(strings.NewReader(input))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	matches := lookup.Match("A碗面条")
	if len(matches) != 4 {
		t.Errorf("Expected 4 matches, got %v", len(matches))
	}
	if matches[0].Simplified != "A" {
		t.Errorf("Expected first match not found, got %v", matches[0])
	}
	if matches[1].Simplified != "碗" {
		t.Errorf("Expected second match got  碗, got %v %v", matches[1].Traditional, matches[1].Simplified)
	}
	if matches[2].Simplified != "碗" {
		t.Errorf("Expected third match got 碗 碗, got %v %v", matches[2].Traditional, matches[2].Simplified)
	}
	if matches[3].Simplified != "面条" {
		t.Errorf("Expected fourth match got 面条, got %v", matches[3].Simplified)
	}
}

func TestMatchMatchesTraditional(t *testing.T) {
	input := "椀 碗 [wan3] /variant of 碗[wan3]/\n碗 碗 [wan3] /bowl/cup/CL:隻|只[zhi1],個|个[ge4]/\n麵條 面条 [mian4 tiao2] /noodles/"

	lookup, err := Build(strings.NewReader(input))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	matches := lookup.Match("A椀麵條")
	if len(matches) != 3 {
		t.Errorf("Expected 3 matches, got %v", len(matches))
	}
	if matches[0].Traditional != "A" {
		t.Errorf("Expected first match not found, got %v", matches[0])
	}
	if matches[1].Traditional != "椀" {
		t.Errorf("Expected second match got 椀, got %v", matches[1].Traditional)
	}
	if matches[2].Traditional != "麵條" {
		t.Errorf("Expected third match got 麵條, got %v", matches[2].Traditional)
	}
}

func TestMatchCollatesUnknown(t *testing.T) {
	input := "面 面 [mian4] /face/side/surface/aspect/top/classifier for flat surfaces such as drums, mirrors, flags etc/\n麪 面 [mian4] /variant of 麵|面[mian4]/\n麵 面 [mian4] /flour/noodles/(of food) soft (not crunchy)/(slang) (of a person) ineffectual/spineless/"

	lookup, err := Build(strings.NewReader(input))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	matches := lookup.Match("biangbiang面")
	if len(matches) != 4 {
		t.Errorf("Expected 4 matches, got %v", len(matches))
	}
	if matches[0].Simplified != "biangbiang" {
		t.Errorf("Expected first match not found, got %v", matches[0])
	}
}
