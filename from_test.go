package callbag_test

import (
	"reflect"
	"testing"

	cb "github.com/alinz/go-callbag"
)

func TestFromIterBasic(t *testing.T) {
	source := cb.FromIter("1", "2", "3", "4")

	expectedProcess := []string{
		"greets",
		"1",
		"2",
		"3",
		"4",
		"end",
	}

	i := 0

	var talkback cb.Source

	source(cb.NewGreets(func(p cb.Payload) {
		switch v := p.(type) {
		case cb.Greets:
			if expectedProcess[i] != "greets" {
				t.Errorf("expected Greets, got %s", reflect.TypeOf(v))
			}
			i++
			talkback = v.Source()
			talkback(cb.NewData(nil))
		case cb.Data:
			val := v.Value().(string)
			if val != expectedProcess[i] {
				t.Errorf("expected Data to be %s, got %s", expectedProcess[i], val)
			}
			i++
			talkback(cb.NewData(nil))
		case cb.Terminate:
			if expectedProcess[i] != "end" {
				t.Errorf("expected end but got %s", expectedProcess[i])
			}
			i++
		}
	}))

	if i < len(expectedProcess) {
		t.Errorf("missing %s", expectedProcess[i])
	}
}

func TestFromRangeBlowupStack(t *testing.T) {

	max := 1000000
	source := cb.FromRange(0, max)
	i := 0

	var talkback cb.Source
	source(cb.NewGreets(func(p cb.Payload) {
		switch v := p.(type) {
		case cb.Greets:
			talkback = v.Source()
			talkback(cb.NewData(nil))
		case cb.Data:
			i++
			talkback(cb.NewData(nil))
		}
	}))

	if i != max {
		t.Errorf("failed to reach %d, got until %d", max, i)
	}
}
