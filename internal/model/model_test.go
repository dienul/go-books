package model

import "testing"

func TestThicknessFromPages(t *testing.T) {
	tests := []struct {
		pages int
		want  string
	}{
		{pages: 1, want: "tipis"},
		{pages: 100, want: "tipis"},
		{pages: 101, want: "tebal"},
	}

	for _, test := range tests {
		if got := ThicknessFromPages(test.pages); got != test.want {
			t.Fatalf("ThicknessFromPages(%d) = %q, ingin %q", test.pages, got, test.want)
		}
	}
}
