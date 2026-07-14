package reader

import "testing"

func TestResolveStatus(t *testing.T) {
	tests := []struct {
		name    string
		meta    metadata
		want    Status
		wantErr bool
	}{
		{"defaults to published", metadata{}, StatusPublished, false},
		{"explicit published", metadata{Status: StatusPublished}, StatusPublished, false},
		{"unlisted", metadata{Status: StatusUnlisted}, StatusUnlisted, false},
		{"draft", metadata{Status: StatusDraft}, StatusDraft, false},
		{"legacy draft flag", metadata{IsDraft: true}, StatusDraft, false},
		{"status wins over legacy flag", metadata{Status: StatusPublished, IsDraft: true}, StatusPublished, false},
		{"unknown status errors", metadata{Status: "nonsense"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.meta.resolveStatus()
			if (err != nil) != tt.wantErr {
				t.Fatalf("resolveStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("resolveStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsListed(t *testing.T) {
	tests := []struct {
		status Status
		want   bool
	}{
		{StatusPublished, true},
		{StatusUnlisted, false},
		{StatusDraft, true},
	}
	for _, tt := range tests {
		if got := (Story{Status: tt.status}).IsListed(); got != tt.want {
			t.Errorf("IsListed() for %q = %v, want %v", tt.status, got, tt.want)
		}
	}
}
