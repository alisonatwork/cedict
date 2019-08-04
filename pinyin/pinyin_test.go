package pinyin

import "testing"

func TestNumberToMarkNoSuffix(t *testing.T) {
	expected := "nǐ hǎo"
	got := NumberToMark("ni3 hao3")
	if got != expected {
		t.Errorf("Expected: %s, got: %s", expected, got)
	}
}

func TestNumberToMarkSuffix(t *testing.T) {
	expected := "wǎng zhàn"
	got := NumberToMark("wang3 zhan4")
	if got != expected {
		t.Errorf("Expected: %s, got: %s", expected, got)
	}
}

func TestNumberToMarkUmlaut(t *testing.T) {
	expected := "lǜ dòu lǜ chá"
	got := NumberToMark("lv4 dou4 lu:4 cha2")
	if got != expected {
		t.Errorf("Expected: %s, got: %s", expected, got)
	}
}
