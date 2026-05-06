package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
	input    string
	expected []string
}{
	{
		input:    "  this is a teST ",
		expected: []string{"this", "is", "a", "test"},
	},

	{
		input:    "CharizArd pikachu BulBasAur  ",
		expected: []string{"charizard", "pikachu", "bulbasaur"},
	},
	{
		input:    "heLLo how aRe  you ",
		expected: []string{"hello", "how", "are", "you"},
	},
}

for _, c := range cases {
	actual := cleanInput(c.input)
	if len(actual) != len(c.expected){
		t.Errorf("length mismatch")
		continue
	}

	for i := range actual {
		word := actual[i]
		expectedWord := c.expected[i]
		if word != expectedWord {
			t.Errorf("actual output does not match expected output")
		}
	}
}

}