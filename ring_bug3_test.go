package main

import "testing"

// TestRemoveNodeClearsStaleEntries verifies that RemoveNode does not retain
// references to removed entries in the underlying ring slice, which would
// prevent garbage collection of the node name strings.
func TestRemoveNodeClearsStaleEntries(t *testing.T) {
	s := NewService()
	s.AddNode("keep", 50)
	s.AddNode("remove", 50)

	totalBefore := len(s.ring) // should be 100
	s.RemoveNode("remove")
	totalAfter := len(s.ring) // should be 50

	if totalAfter != 50 {
		t.Fatalf("expected 50 entries after removal, got %d", totalAfter)
	}

	// The backing array should not retain stale entries beyond len.
	// After proper cleanup, cap should equal len (new allocation) or
	// stale slots should be zeroed. Check that the backing array tail
	// does not still reference the removed node name.
	fullSlice := s.ring[:cap(s.ring)]
	staleCount := 0
	for i := totalAfter; i < len(fullSlice); i++ {
		if fullSlice[i].node != "" {
			staleCount++
		}
	}
	if staleCount > 0 {
		t.Fatalf("memory leak: %d stale entries in backing array still reference removed node (cap=%d, len=%d)",
			staleCount, cap(s.ring), totalAfter)
	}
	_ = totalBefore
}
