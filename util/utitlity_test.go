package util

import (
	"testing"
	"reflect"
)

func TestIsEquals(t *testing.T) {
	tests := []struct {
		input     string
		testInput string
		expected  bool
	}{
		{
			input:     "Hello",
			testInput: "Hello",
			expected:  true,
		},
		{
			input:     "Hello",
			testInput: "hello",
			expected:  false,
		}}

	for testNumber, test := range tests {
		actual := IsEquals(test.input, test.testInput)
		if actual != test.expected {
			t.Errorf("Test No %d: Expected: %v but Got %v", testNumber+1, test.expected, actual)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		data     string
		expected  bool
	}{
		{
			data:     "",
			expected:  true,
		},
		{
			data:     "Hello",
			expected:  false,
		}}

	for testNumber, test := range tests {
		actual := IsEmpty(test.data)
		if actual != test.expected {
			t.Errorf("Test No %d: Expected: %v but Got %v", testNumber+1, test.expected, actual)
		}
	}
}


func TestHashAndSalt(t *testing.T) {
	tests := []struct {
		password string
		expected string
	}{
		{
			password: "pass123!",
			expected: "string",
		},
		{
			password: "Hello",
			expected: "",
		}}

	for testNumber, test := range tests {
		actual := HashAndSalt([]byte(test.password))
		if reflect.TypeOf(actual) != reflect.TypeOf(test.expected) {
			t.Errorf("Test No %d: Expected : %v but Got: %v", testNumber+1, test.expected, actual)
		}
	}
}

func TestComparePasswords(t *testing.T) {
	tests := []struct {
		hashedPassword string
		plainPassword  string
		expected       bool
	}{
		{
			hashedPassword: "$2a$04$LHqVJqStaZowVVSJ/otrW./FXJ65Os.Ei5R0tMcgFjVnESmPvhs3a",
			plainPassword:  "pass123!",
			expected:       true,
		},
		{
			hashedPassword: "Hello",
			plainPassword:  "pass123!",
			expected:       false,
		}}

	for testNumber, test := range tests {
		actual := ComparePasswords(test.hashedPassword, test.plainPassword)
		if actual != test.expected {
			t.Errorf("Test No %d: Expected: %v but Got: %v", testNumber+1, test.expected, actual)
		}
	}
}