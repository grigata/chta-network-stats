package main

import "testing"

func TestReadOptions(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		web   bool
		count int
	}{
		{"defaults", nil, false, 100}, {"console count", []string{"500"}, false, 500},
		{"web defaults", []string{"web"}, true, 100}, {"web count", []string{"WEB", "250"}, true, 250},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := readOptions(test.args)
			if err != nil {
				t.Fatal(err)
			}
			if got.Web != test.web || got.BlockCount != test.count {
				t.Fatalf("got %+v", got)
			}
		})
	}
}

func TestReadOptionsRejectsInvalidInput(t *testing.T) {
	for _, args := range [][]string{{"0"}, {"-1"}, {"web", "nope"}, {"web", "1", "extra"}} {
		if _, err := readOptions(args); err == nil {
			t.Fatalf("expected error for %v", args)
		}
	}
}
