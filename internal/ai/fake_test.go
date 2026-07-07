package ai

import (
	"context"
	"testing"
)

func TestNewFakeProvider(t *testing.T) {
	p := NewFakeProvider()
	if p.Name() != "fake" {
		t.Fatalf("expected name fake, got %s", p.Name())
	}
	if p.Model() != "fake" {
		t.Fatalf("expected model fake, got %s", p.Model())
	}

	findings, err := p.Review(context.Background(), "review this")
	if err != nil {
		t.Fatalf("Review error: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 canned findings, got %d", len(findings))
	}
	if findings[0].ID != "F001" {
		t.Fatalf("expected F001, got %s", findings[0].ID)
	}
}
