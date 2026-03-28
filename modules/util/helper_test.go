package util

import "testing"

func TestCheckExt(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
	}{
		{"site.html", false},
		{"picture.jpg", true},
		{"olderSite.htm", true},
		{"file.exe", false},
	}

	for _, tc := range testCases {
		result := CheckExt(tc.path)
		if result != tc.expected {
			t.Errorf("case %s: expected %v but got %v", tc.path, tc.expected, result)
		}
	}
}
