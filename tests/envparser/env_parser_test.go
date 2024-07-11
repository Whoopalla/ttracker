package envparser_test

import (
	"os"
	"test/ttracker/pkg/envparser"
	"testing"
)

type testEntry struct {
	input    any
	expected any
}

func TestIsCorrectExtension_Test(t *testing.T) {
	entries := []struct {
		input    string
		expected bool
	}{
		{"a", false},
		{"abcd", false},
		{"abcde", false},
		{".env", false},
		{".en", false},
		{".envi", false},

		{"a.env", true},
		{"ab.env", true},
		{"aaa.env", true},
		{".env.env", true},
	}

	for _, e := range entries {
		if result := envparser.IsCorrectExtension(string(e.input)); result != e.expected {
			t.Fatalf("Input: |%v| Expected: %v but was: %v", e.input, e.expected, result)
		}
	}
}

func TestSetEnvironment(t *testing.T) {
	const dotEnvFileContent = "2+2=4\nRose=Bud\nCharley=Chaplyn"
	const dotEnvFileName = "config.env"
	f, err := os.CreateTemp(os.TempDir(), "123")
	os.Rename(f.Name(), dotEnvFileName)
	_, err = f.Write([]byte(dotEnvFileContent))

	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dotEnvFileName)
	err = envparser.Load(dotEnvFileName)
	if err != nil {
		t.Fatal(err)
	}

	if os.Getenv("2+2") != "4" || os.Getenv("Rose") != "Bud" || os.Getenv("Charley") != "Chaplyn" {
		t.Fatal()
	}
}
