package cosvectors

import (
	"context"
	"errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	cosSDK "github.com/tencentyun/cos-go-sdk-v5"
)

func (s *CosVectorService) PutBucketPolicy(ctx context.Context, bucket, policy string) error {
	if err := validateCosVectorBucketName(bucket); err != nil {
		return err
	}
	if err := validateCosVectorPolicy(policy); err != nil {
		return err
	}
	if err := ratelimit.CheckContext(ctx, "PutVectorBucketPolicy"); err != nil {
		return s.apiError("PutVectorBucketPolicy", nil, err)
	}
	resp, err := s.client.Vector.PutVectorBucketPolicy(ctx, &cosSDK.PutVectorBucketPolicyOptions{VectorBucketName: bucket, Policy: policy})
	if err != nil || !cosVectorResponseOK(resp) {
		return s.apiError("PutVectorBucketPolicy", resp, err)
	}
	return nil
}

func (s *CosVectorService) GetBucketPolicy(ctx context.Context, bucket string) (string, error) {
	if err := validateCosVectorBucketName(bucket); err != nil {
		return "", err
	}
	if err := ratelimit.CheckContext(ctx, "GetVectorBucketPolicy"); err != nil {
		return "", s.apiError("GetVectorBucketPolicy", nil, err)
	}
	result, resp, err := s.client.Vector.GetVectorBucketPolicy(ctx, &cosSDK.GetVectorBucketPolicyOptions{VectorBucketName: bucket})
	if err != nil || !cosVectorResponseOK(resp) {
		return "", s.apiError("GetVectorBucketPolicy", resp, err)
	}
	if result == nil {
		return "", errors.New("GetVectorBucketPolicy returned no policy")
	}
	if _, err := normalizeCosVectorPolicy(result.Policy); err != nil {
		return "", errors.New("GetVectorBucketPolicy returned an invalid COS policy")
	}
	return result.Policy, nil
}

func (s *CosVectorService) DeleteBucketPolicy(ctx context.Context, bucket string) error {
	if err := validateCosVectorBucketName(bucket); err != nil {
		return err
	}
	if err := ratelimit.CheckContext(ctx, "DeleteVectorBucketPolicy"); err != nil {
		return s.apiError("DeleteVectorBucketPolicy", nil, err)
	}
	resp, err := s.client.Vector.DeleteVectorBucketPolicy(ctx, &cosSDK.DeleteVectorBucketPolicyOptions{VectorBucketName: bucket})
	if err != nil || !cosVectorResponseOK(resp) {
		return s.apiError("DeleteVectorBucketPolicy", resp, err)
	}
	return nil
}

func (s *CosVectorService) WaitBucketPolicy(ctx context.Context, bucket, expected string, exists bool) error {
	return waitCosVectorRead(ctx, func(ctx context.Context) (bool, error) {
		policy, err := s.GetBucketPolicy(ctx, bucket)
		return (exists && err == nil && cosVectorPoliciesEqual(policy, expected)) || (!exists && isCosVectorBucketNotFound(err)), err
	})
}
