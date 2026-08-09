package projects

import (
	"fmt"
	"slices"
	"testing"

	"github.com/pierrre/assert"
)

func ExampleList() {
	for _, p := range List {
		fmt.Println(p)
	}
	// Output:
	// assert
	// cellauto
	// compare
	// di
	// errors
	// file-duplicate
	// file-random
	// geohash
	// githubhook
	// go-cache-prog
	// go-libs
	// go-stuff
	// image-viewer
	// langton
	// mandelbrot
	// pretty
	// unlimited-channel
	// vld
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    []Project
		wantErr string
	}{
		{
			name: "single",
			data: "a",
			want: []Project{"a"},
		},
		{
			name: "multiple",
			data: "a\nb\nc\n",
			want: []Project{"a", "b", "c"},
		},
		{
			name: "blank lines ignored",
			data: "\na\n\nb\n\n",
			want: []Project{"a", "b"},
		},
		{
			name: "whitespace trimmed",
			data: " a \n\tb\t\n",
			want: []Project{"a", "b"},
		},
		{
			name: "trailing newline",
			data: "a\nb",
			want: []Project{"a", "b"},
		},
		{
			name:    "empty",
			data:    "",
			wantErr: "no projects",
		},
		{
			name:    "only blank lines",
			data:    "\n\n",
			wantErr: "no projects",
		},
		{
			name:    "duplicate",
			data:    "a\na\n",
			wantErr: `duplicate project "a" at line 2`,
		},
		{
			name:    "not sorted",
			data:    "b\na\n",
			wantErr: `project "a" at line 2 is not sorted (previous "b")`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse([]byte(tt.data))
			if tt.wantErr != "" {
				assert.ErrorEqual(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.SliceEqual(t, got, tt.want)
		})
	}
}

func TestMustParse(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ps := MustParse([]byte("a\nb\n"))
		assert.SliceEqual(t, ps, []Project{"a", "b"})
	})
	t.Run("panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustParse([]byte("a\na\n"))
		})
	})
}

func TestList(t *testing.T) {
	assert.SliceNotEmpty(t, List)
	assert.True(t, slices.IsSorted(List))
	uniq := slices.Compact(slices.Clone(List))
	assert.SliceLen(t, uniq, len(List))
}
