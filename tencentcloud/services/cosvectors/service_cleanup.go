package cosvectors

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// DeleteBucketAfterCleanup 仅等待确证非空且无可见索引的异步回收。
func (s *CosVectorService) DeleteBucketAfterCleanup(ctx context.Context, bucket string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	delay := time.Second
	var lastErr error
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("waiting for COS bucket cleanup (last error: %v): %w", lastErr, err)
		}
		err := s.DeleteBucket(ctx, bucket)
		if err == nil || isCosVectorBucketNotFound(err) {
			return nil
		}
		var apiErr *cosVectorBucketError
		if !errors.As(err, &apiErr) || !apiErr.bucketNotEmpty {
			return err
		}
		indexes, listErr := s.ListIndexes(ctx, bucket)
		if listErr != nil {
			return listErr
		}
		if len(indexes) != 0 {
			return err
		}
		lastErr = err
		timer := time.NewTimer(delay/2 + time.Duration(rand.Int63n(int64(delay/2))))
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("waiting for COS bucket cleanup (last error: %v): %w", lastErr, ctx.Err())
		case <-timer.C:
		}
		if delay < 16*time.Second {
			delay *= 2
		}
	}
}
