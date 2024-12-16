package main

import "testing"

func TestCleanInput(t *testing.T) {
    cases := []struct { input string; expected []string } {
        {
            input: "  hello   world   ",
            expected: []string { "hello", "world" },
        },
        {
            input: "1234 this IS    a TeSt!",
            expected: []string { "1234", "this", "is", "a", "test!" },
        },
        {
            input: "A B C 1 2 3",
            expected: []string { "a", "b", "c", "1", "2", "3" },
        },
        {
            input: "                          ",
            expected: []string {},
        },
    }

    for _, c := range cases {
        result := cleanInput(c.input)
        for i := range result {
            if result[i] != c.expected[i] {
                t.Errorf("clean input failed with '%s': %s != %s", c.input, result[i], c.expected[i])
            }
        }
    }
}
