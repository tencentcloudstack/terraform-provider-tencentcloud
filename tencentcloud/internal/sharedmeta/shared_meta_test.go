package sharedmeta

import (
	"sync"
	"testing"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

// TestSharedMeta_SetGet verifies the basic write/read semantics.
func TestSharedMeta_SetGet(t *testing.T) {
	t.Cleanup(ResetSharedMetaForTest)
	ResetSharedMetaForTest()

	if got := GetSharedMeta(); got != nil {
		t.Fatalf("expected nil before SetSharedMeta, got %p", got)
	}

	c := &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
	SetSharedMeta(c)

	got := GetSharedMeta()
	if got != c {
		t.Fatalf("expected pointer %p, got %p", c, got)
	}

	if got.Region != "ap-guangzhou" {
		t.Fatalf("expected region ap-guangzhou, got %q", got.Region)
	}
}

// TestSharedMeta_LastWriteWins verifies that multiple writes leave the most
// recent value (e.g. acceptance tests that reuse the singleton).
func TestSharedMeta_LastWriteWins(t *testing.T) {
	t.Cleanup(ResetSharedMetaForTest)
	ResetSharedMetaForTest()

	first := &connectivity.TencentCloudClient{Region: "ap-shanghai"}
	second := &connectivity.TencentCloudClient{Region: "ap-beijing"}

	SetSharedMeta(first)
	SetSharedMeta(second)

	if got := GetSharedMeta(); got != second {
		t.Fatalf("expected the most recent client, got %+v", got)
	}
}

// TestSharedMeta_ConcurrentSetGet verifies concurrency safety: many
// goroutines calling Set/Get concurrently must not race (go test -race
// would catch a violation).
func TestSharedMeta_ConcurrentSetGet(t *testing.T) {
	t.Cleanup(ResetSharedMetaForTest)
	ResetSharedMetaForTest()

	const goroutines = 32
	const iterations = 200

	clients := make([]*connectivity.TencentCloudClient, goroutines)
	for i := range clients {
		clients[i] = &connectivity.TencentCloudClient{Region: "r"}
	}

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := 0; i < goroutines; i++ {
		idx := i
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				SetSharedMeta(clients[idx])
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = GetSharedMeta()
			}
		}()
	}

	wg.Wait()

	if GetSharedMeta() == nil {
		t.Fatalf("expected non-nil after concurrent writes")
	}
}

// TestSharedMeta_ResetForTest verifies that the test helper clears state.
func TestSharedMeta_ResetForTest(t *testing.T) {
	SetSharedMeta(&connectivity.TencentCloudClient{Region: "ap-guangzhou"})
	ResetSharedMetaForTest()

	if got := GetSharedMeta(); got != nil {
		t.Fatalf("expected nil after reset, got %p", got)
	}
}
