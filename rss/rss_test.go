package rss

import "testing"

func TestGenerateFeed(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		want    string
		wantErr bool
	}{
		{"base test", 5, "", false},
		{"base test", 3000, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateFeed(4)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
