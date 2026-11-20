package cosvectors

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	cosSDK "github.com/tencentyun/cos-go-sdk-v5"
)

var cosVectorIndexNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)
var cosVectorMetadataKeyPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]{0,62}$`)

func validateCosVectorIndexName(name string) error {
	if len(name) < 3 || len(name) > 63 || strings.Contains(name, "..") || !cosVectorIndexNamePattern.MatchString(name) {
		return errors.New("index_name must be 3-63 lowercase letters, digits, dots or hyphens, starting and ending with a letter or digit, without '..'")
	}
	return nil
}

func parseCosVectorIndexImportID(id, region string) (string, string, error) {
	parts := strings.Split(id, "/index/")
	if len(parts) == 2 && strings.HasPrefix(parts[0], "qcs:") {
		bucket, err := parseCosVectorBucketImportID(parts[0], region)
		if err != nil {
			return "", "", err
		}
		return bucket, parts[1], validateCosVectorIndexName(parts[1])
	}
	parts = strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", errors.New("import requires bucket-name/index-name or a matching index QCS")
	}
	if err := validateCosVectorBucketName(parts[0]); err != nil {
		return "", "", err
	}
	return parts[0], parts[1], validateCosVectorIndexName(parts[1])
}

func (s *CosVectorService) validIndexQCS(qcs, bucket, index string) bool {
	gotBucket, gotIndex, err := parseCosVectorIndexImportID(qcs, s.region)
	return strings.HasPrefix(qcs, "qcs:") && err == nil && gotBucket == bucket && gotIndex == index
}

func validateCosVectorIndexOptions(opt *cosSDK.CreateIndexOptions) error {
	if opt == nil {
		return errors.New("index configuration is required")
	}
	if err := validateCosVectorBucketName(opt.VectorBucketName); err != nil {
		return err
	}
	if err := validateCosVectorIndexName(opt.IndexName); err != nil {
		return err
	}
	if opt.DataType != "float32" || opt.Dimension < 1 || opt.Dimension > 4096 || (opt.DistanceMetric != "cosine" && opt.DistanceMetric != "euclidean") {
		return errors.New("index requires data_type=float32, dimension between 1 and 4096, and distance_metric cosine or euclidean")
	}
	if opt.MetadataConfiguration != nil {
		keys := opt.MetadataConfiguration.NonFilterableMetadataKeys
		if len(keys) < 1 || len(keys) > 10 {
			return errors.New("non_filterable_metadata_keys must contain 1-10 distinct non-empty keys")
		}
		seen := make(map[string]bool, len(keys))
		for _, key := range keys {
			if !cosVectorMetadataKeyPattern.MatchString(key) || seen[key] {
				return errors.New("metadata keys must be distinct and match ^[a-zA-Z_][a-zA-Z0-9_]{0,62}$")
			}
			seen[key] = true
		}
	}
	return nil
}

func (s *CosVectorService) CreateIndex(ctx context.Context, opt *cosSDK.CreateIndexOptions) (string, error) {
	if err := validateCosVectorIndexOptions(opt); err != nil {
		return "", err
	}
	if err := ratelimit.CheckContext(ctx, "CreateIndex"); err != nil {
		return "", s.apiError("CreateIndex", nil, err)
	}
	result, resp, err := s.client.Vector.CreateIndex(ctx, opt)
	if err != nil || !cosVectorResponseOK(resp) {
		return "", s.apiError("CreateIndex", resp, err)
	}
	if result == nil || !s.validIndexQCS(result.IndexQcs, opt.VectorBucketName, opt.IndexName) {
		return "", s.apiError("CreateIndex", resp, errors.New("invalid create response"))
	}
	return result.IndexQcs, nil
}

func (s *CosVectorService) GetIndex(ctx context.Context, bucket, name string) (*cosSDK.IndexInfo, error) {
	if err := validateCosVectorBucketName(bucket); err != nil {
		return nil, err
	}
	if err := validateCosVectorIndexName(name); err != nil {
		return nil, err
	}
	if err := ratelimit.CheckContext(ctx, "GetIndex"); err != nil {
		return nil, s.apiError("GetIndex", nil, err)
	}
	result, resp, err := s.client.Vector.GetIndex(ctx, &cosSDK.GetIndexOptions{VectorBucketName: bucket, IndexName: name})
	if err != nil || !cosVectorResponseOK(resp) {
		return nil, s.apiError("GetIndex", resp, err)
	}
	if result == nil || result.Index == nil || result.Index.VectorBucketName != bucket || result.Index.IndexName != name || !s.validIndexQCS(result.Index.IndexQcs, bucket, name) {
		return nil, errors.New("GetIndex returned missing or mismatched index identity")
	}
	index := result.Index
	if index.MetadataConfiguration != nil && len(index.MetadataConfiguration.NonFilterableMetadataKeys) == 0 {
		index.MetadataConfiguration = nil
	}
	if err := validateCosVectorIndexOptions(&cosSDK.CreateIndexOptions{VectorBucketName: bucket, IndexName: name, DataType: index.DataType, Dimension: index.Dimension, DistanceMetric: index.DistanceMetric, MetadataConfiguration: index.MetadataConfiguration}); err != nil {
		return nil, errors.New("GetIndex returned an unsupported index configuration")
	}
	return index, nil
}

func (s *CosVectorService) DeleteIndex(ctx context.Context, bucket, name string) error {
	if err := validateCosVectorBucketName(bucket); err != nil {
		return err
	}
	if err := validateCosVectorIndexName(name); err != nil {
		return err
	}
	if err := ratelimit.CheckContext(ctx, "DeleteIndex"); err != nil {
		return s.apiError("DeleteIndex", nil, err)
	}
	resp, err := s.client.Vector.DeleteIndex(ctx, &cosSDK.DeleteIndexOptions{VectorBucketName: bucket, IndexName: name})
	if err != nil || !cosVectorResponseOK(resp) {
		return s.apiError("DeleteIndex", resp, err)
	}
	return nil
}

// ListIndexes 先验证所有分页与资源身份，避免无效分页时部分误删。
func (s *CosVectorService) ListIndexes(ctx context.Context, bucket string) ([]cosSDK.IndexBrief, error) {
	if err := validateCosVectorBucketName(bucket); err != nil {
		return nil, err
	}
	var indexes []cosSDK.IndexBrief
	seenTokens, seenNames := map[string]bool{}, map[string]bool{}
	token := ""
	for {
		if err := ratelimit.CheckContext(ctx, "ListIndexes"); err != nil {
			return nil, s.apiError("ListIndexes", nil, err)
		}
		result, resp, err := s.client.Vector.ListIndexes(ctx, &cosSDK.ListIndexesOptions{VectorBucketName: bucket, MaxResults: 100, NextToken: token})
		if err != nil || !cosVectorResponseOK(resp) {
			return nil, s.apiError("ListIndexes", resp, err)
		}
		if result == nil || result.Indexes == nil {
			return nil, errors.New("ListIndexes returned a missing indexes array")
		}
		for _, index := range result.Indexes {
			if index.VectorBucketName != bucket || !s.validIndexQCS(index.IndexQcs, bucket, index.IndexName) {
				return nil, errors.New("ListIndexes returned mismatched resource identity")
			}
			if !seenNames[index.IndexName] {
				indexes = append(indexes, index)
				seenNames[index.IndexName] = true
			}
		}
		if result.NextToken == "" {
			return indexes, nil
		}
		if seenTokens[result.NextToken] {
			return nil, errors.New("ListIndexes repeated a pagination token")
		}
		seenTokens[result.NextToken] = true
		token = result.NextToken
	}
}

func (s *CosVectorService) WaitIndex(ctx context.Context, bucket, name string, exists bool) (*cosSDK.IndexInfo, error) {
	var index *cosSDK.IndexInfo
	err := waitCosVectorRead(ctx, func(ctx context.Context) (bool, error) {
		got, err := s.GetIndex(ctx, bucket, name)
		index = got
		return (exists && err == nil) || (!exists && isCosVectorBucketNotFound(err)), err
	})
	return index, err
}

// waitCosVectorRead 仅用于可见性确认，不重发修改请求。
func waitCosVectorRead(ctx context.Context, read func(context.Context) (bool, error)) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	delay := 250 * time.Millisecond
	var lastErr error
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("waiting for COS vector resource (last error: %v): %w", lastErr, err)
		}
		done, err := read(ctx)
		if done {
			return nil
		}
		if err != nil && !isCosVectorBucketNotFound(err) && !isCosVectorBucketRetryable(err) {
			return err
		}
		lastErr = err
		timer := time.NewTimer(delay/2 + time.Duration(rand.Int63n(int64(delay/2))))
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("waiting for COS vector resource (last error: %v): %w", lastErr, ctx.Err())
		case <-timer.C:
		}
		if delay < 2*time.Second {
			delay *= 2
		}
	}
}

func (s *CosVectorService) EmptyBucket(ctx context.Context, bucket string) error {
	indexes, err := s.ListIndexes(ctx, bucket)
	if isCosVectorBucketNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, index := range indexes {
		if err := s.DeleteIndex(ctx, bucket, index.IndexName); err != nil && !isCosVectorBucketNotFound(err) {
			return err
		}
		if _, err := s.WaitIndex(ctx, bucket, index.IndexName, false); err != nil {
			return err
		}
	}
	return nil
}
