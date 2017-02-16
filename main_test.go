package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var shatests = []struct {
	content string
	hash    string
}{
	{
		"foo", "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae",
	},
	{
		"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	},
}

func TestComputeSHA256(t *testing.T) {
	for i, test := range shatests {
		tempfile, err := ioutil.TempFile("", "psha-test-")
		if err != nil {
			t.Fatalf("TempFile returned %v", err)
		}

		_, err = io.WriteString(tempfile, test.content)
		if err != nil {
			t.Fatalf("WriteString() %v", err)
		}

		err = tempfile.Close()
		if err != nil {
			t.Fatalf("Close() %v", err)
		}

		hash, err := computeSHA256(tempfile.Name())
		if err != nil {
			t.Fatalf("computeSHA256() %v", err)
		}

		if hash != test.hash {
			t.Errorf("test %d failed: want %v, got %v", i, test.hash, hash)
		}

		if err = os.Remove(tempfile.Name()); err != nil {
			t.Fatalf("Remove() %v", err)
		}
	}
}

func BenchmarkComputeSHA256(b *testing.B) {
	tempfile, err := ioutil.TempFile("", "psha-test-")
	if err != nil {
		b.Fatalf("TempFile returned %v", err)
	}

	_, err = io.WriteString(tempfile, "Aachen")
	if err != nil {
		b.Fatalf("WriteString() %v", err)
	}

	err = tempfile.Close()
	if err != nil {
		b.Fatalf("Close() %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hash, err := computeSHA256(tempfile.Name())
		if err != nil {
			b.Fatalf("computeSHA256() %v", err)
		}

		if hash != "6aa8d75d4bfe6065abe8dd7de4cf23bcab2daa596a56b1f119a3413257a231d5" {
			b.Errorf("test %d failed", i)
		}

	}

	if err = os.Remove(tempfile.Name()); err != nil {
		b.Fatalf("Remove() %v", err)
	}
}
