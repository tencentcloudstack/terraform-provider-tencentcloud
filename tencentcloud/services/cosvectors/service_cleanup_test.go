package cosvectors

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	cosSDK "github.com/tencentyun/cos-go-sdk-v5"
)

func TestCosVectorBucketCleanupPendingRecovery(t *testing.T) {
	s, m := newCosVectorBucketTestService(t)
	bucket := "tf-cleanup-1250000000"
	if _, err := s.CreateBucket(context.Background(), bucket, nil); err != nil {
		t.Fatal(err)
	}
	m.mu.Lock()
	m.cleanupPending = 1
	m.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.DeleteBucketAfterCleanup(ctx, bucket); err != nil {
		t.Fatal(err)
	}
	if m.count() != 0 || m.callCount("DeleteVectorBucket") != 2 || m.callCount("ListIndexes") != 1 {
		t.Fatal("asynchronous cleanup not confirmed safely")
	}
}

func TestCosVectorBucketCleanupPendingDeadline(t *testing.T) {
	s, m := newCosVectorBucketTestService(t)
	bucket := "tf-cleanup-timeout-1250000000"
	if _, err := s.CreateBucket(context.Background(), bucket, nil); err != nil {
		t.Fatal(err)
	}
	m.mu.Lock()
	m.cleanupPending = 100
	m.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if err := s.DeleteBucketAfterCleanup(ctx, bucket); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("cleanup deadline lost: %v", err)
	}
	if m.count() != 1 || m.callCount("DeleteVectorBucket") != 1 {
		t.Fatal("deadline caused repeated deletion or lost bucket")
	}
}

func TestCosVectorBucketCleanupConflictClassification(t *testing.T) {
	s, m := newCosVectorBucketTestService(t)
	bucket := "tf-cleanup-conflict-1250000000"
	if _, err := s.CreateBucket(context.Background(), bucket, nil); err != nil {
		t.Fatal(err)
	}
	m.setFailure("DeleteVectorBucket", 409)
	if err := s.DeleteBucketAfterCleanup(context.Background(), bucket); err == nil {
		t.Fatal("generic conflict ignored")
	}
	if m.callCount("DeleteVectorBucket") != 1 || m.callCount("ListIndexes") != 0 {
		t.Fatal("unrecognized409 was retried")
	}
	m.setFailure("DeleteVectorBucket", 0)
	m.mu.Lock()
	m.cleanupPending = 1
	m.mu.Unlock()
	if _, err := s.CreateIndex(context.Background(), &cosSDK.CreateIndexOptions{VectorBucketName: bucket, IndexName: "live-index", DataType: "float32", Dimension: 3, DistanceMetric: "cosine"}); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteBucketAfterCleanup(context.Background(), bucket); err == nil {
		t.Fatal("live index bucket deleted")
	}
	for _, action := range []string{"DeleteIndex", "PutVectorBucketPolicy"} {
		if m.callCount(action) != 0 {
			t.Fatal("cleanup waiter modified child resources")
		}
	}
	resp := &cosSDK.Response{Response: &http.Response{StatusCode: 409, Header: make(http.Header)}}
	for _, action := range []string{"CreateVectorBucket", "DeleteIndex", "GetVectorBucket"} {
		err := s.apiError(action, resp, &cosSDK.VectorErrorResponse{Code: "ConflictException", Message: "bucket not empty", Response: resp.Response})
		var apiErr *cosVectorBucketError
		if !errors.As(err, &apiErr) || apiErr.bucketNotEmpty {
			t.Fatal("non-delete bucket action treated as cleanup pending")
		}
	}
}

func TestCosVectorBucketCleanupListPermissionFailure(t *testing.T) {
	s, m := newCosVectorBucketTestService(t)
	bucket := "tf-cleanup-permission-1250000000"
	if _, err := s.CreateBucket(context.Background(), bucket, nil); err != nil {
		t.Fatal(err)
	}
	m.mu.Lock()
	m.cleanupPending = 2
	m.mu.Unlock()
	m.setFailure("ListIndexes", 403)
	if err := s.DeleteBucketAfterCleanup(context.Background(), bucket); err == nil {
		t.Fatal("list permission failure ignored")
	}
	if m.count() != 1 || m.callCount("DeleteVectorBucket") != 1 {
		t.Fatal("permission failure retried or removed bucket")
	}
}
