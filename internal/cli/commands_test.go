package cli

import (
	"reflect"
	"testing"
)

func TestExpandYolo(t *testing.T) {
	cases := []struct {
		name    string
		in      []string
		wantOut []string
		wantHit bool
	}{
		{
			name:    "replaces yolo",
			in:      []string{"--resume", "--yolo"},
			wantOut: []string{"--resume", "--dangerously-skip-permissions"},
			wantHit: true,
		},
		{
			name:    "no yolo present",
			in:      []string{"--resume", "-p", "hi"},
			wantOut: []string{"--resume", "-p", "hi"},
			wantHit: false,
		},
		{
			name:    "multiple yolo occurrences",
			in:      []string{"--yolo", "--yolo"},
			wantOut: []string{"--dangerously-skip-permissions", "--dangerously-skip-permissions"},
			wantHit: true,
		},
		{
			name:    "empty input",
			in:      []string{},
			wantOut: []string{},
			wantHit: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, hit := expandYolo(c.in)
			if !reflect.DeepEqual(out, c.wantOut) {
				t.Errorf("out = %#v, want %#v", out, c.wantOut)
			}
			if hit != c.wantHit {
				t.Errorf("hit = %v, want %v", hit, c.wantHit)
			}
		})
	}
}
