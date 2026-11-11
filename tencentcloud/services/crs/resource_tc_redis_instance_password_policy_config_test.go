package crs_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svccrs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"
)

// mockMetaForRedisInstancePasswordPolicyConfig implements tccommon.ProviderMeta
type mockMetaForRedisInstancePasswordPolicyConfig struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForRedisInstancePasswordPolicyConfig) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForRedisInstancePasswordPolicyConfig{}

func newMockMetaForRedisInstancePasswordPolicyConfig() *mockMetaForRedisInstancePasswordPolicyConfig {
	return &mockMetaForRedisInstancePasswordPolicyConfig{client: &connectivity.TencentCloudClient{}}
}

func ptrBoolPasswordPolicy(v bool) *bool {
	return &v
}

func ptrInt64PasswordPolicy(v int64) *int64 {
	return &v
}

func ptrStrPasswordPolicy(s string) *string {
	return &s
}

// go test ./tencentcloud/services/crs/ -run "TestRedisInstancePasswordPolicyConfig" -v -count=1 -gcflags="all=-l"

// TestRedisInstancePasswordPolicyConfig_Schema tests the schema definition
func TestRedisInstancePasswordPolicyConfig_Schema(t *testing.T) {
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Importer)

	// Required fields
	assert.Contains(t, res.Schema, "instance_id")
	assert.True(t, res.Schema["instance_id"].Required)
	assert.True(t, res.Schema["instance_id"].ForceNew)

	assert.Contains(t, res.Schema, "enabled")
	assert.True(t, res.Schema["enabled"].Required)

	// Optional + Computed fields
	for _, field := range []string{"min_letter_count", "min_digit_count", "min_special_count", "min_length"} {
		assert.Contains(t, res.Schema, field)
		assert.True(t, res.Schema[field].Optional)
		assert.True(t, res.Schema[field].Computed)
	}
}

// TestRedisInstancePasswordPolicyConfig_Read_Success tests Read populates fields from PasswordPolicy
func TestRedisInstancePasswordPolicyConfig_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisInstancePasswordPolicyConfig().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicyWithContext", func(ctx context.Context, request *redis.DescribeInstancePasswordPolicyRequest) (*redis.DescribeInstancePasswordPolicyResponse, error) {
		assert.Equal(t, "crs-cqdfdzvt", *request.InstanceId)
		resp := redis.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: &redis.PasswordPolicy{
				Enabled:         ptrBoolPasswordPolicy(true),
				MinLetterCount:  ptrInt64PasswordPolicy(2),
				MinDigitCount:   ptrInt64PasswordPolicy(3),
				MinSpecialCount: ptrInt64PasswordPolicy(1),
				MinLength:       ptrInt64PasswordPolicy(12),
			},
			RequestId: ptrStrPasswordPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisInstancePasswordPolicyConfig()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-cqdfdzvt",
		"enabled":     true,
	})
	d.SetId("crs-cqdfdzvt")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-cqdfdzvt", d.Get("instance_id"))
	assert.Equal(t, true, d.Get("enabled"))
	assert.Equal(t, 2, d.Get("min_letter_count"))
	assert.Equal(t, 3, d.Get("min_digit_count"))
	assert.Equal(t, 1, d.Get("min_special_count"))
	assert.Equal(t, 12, d.Get("min_length"))
}

// TestRedisInstancePasswordPolicyConfig_Read_NilPasswordPolicy tests Read clears state when PasswordPolicy is nil
func TestRedisInstancePasswordPolicyConfig_Read_NilPasswordPolicy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisInstancePasswordPolicyConfig().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicyWithContext", func(ctx context.Context, request *redis.DescribeInstancePasswordPolicyRequest) (*redis.DescribeInstancePasswordPolicyResponse, error) {
		resp := redis.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: nil,
			RequestId:      ptrStrPasswordPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisInstancePasswordPolicyConfig()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-cqdfdzvt",
		"enabled":     true,
	})
	d.SetId("crs-cqdfdzvt")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestRedisInstancePasswordPolicyConfig_Read_InstanceNotExists tests Read returns error when instance not found
func TestRedisInstancePasswordPolicyConfig_Read_InstanceNotExists(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisInstancePasswordPolicyConfig().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicyWithContext", func(ctx context.Context, request *redis.DescribeInstancePasswordPolicyRequest) (*redis.DescribeInstancePasswordPolicyResponse, error) {
		return nil, sdkErrors.NewTencentCloudSDKError("ResourceNotFound.InstanceNotExists", "Instance not found", "fake-request-id")
	})

	meta := newMockMetaForRedisInstancePasswordPolicyConfig()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-cqdfdzvt",
		"enabled":     true,
	})
	d.SetId("crs-cqdfdzvt")

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestRedisInstancePasswordPolicyConfig_Create_Success tests Create sets ID and delegates to Update
func TestRedisInstancePasswordPolicyConfig_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisInstancePasswordPolicyConfig().client, "UseRedisClient", redisClient)

	// Mock ModifyInstancePasswordPolicyWithContext for Update
	patches.ApplyMethodFunc(redisClient, "ModifyInstancePasswordPolicyWithContext", func(ctx context.Context, request *redis.ModifyInstancePasswordPolicyRequest) (*redis.ModifyInstancePasswordPolicyResponse, error) {
		assert.Equal(t, "crs-cqdfdzvt", *request.InstanceId)
		assert.NotNil(t, request.PasswordPolicy)
		assert.Equal(t, true, *request.PasswordPolicy.Enabled)
		resp := redis.NewModifyInstancePasswordPolicyResponse()
		resp.Response = &redis.ModifyInstancePasswordPolicyResponseParams{
			RequestId: ptrStrPasswordPolicy("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeInstancePasswordPolicyWithContext for Read
	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicyWithContext", func(ctx context.Context, request *redis.DescribeInstancePasswordPolicyRequest) (*redis.DescribeInstancePasswordPolicyResponse, error) {
		resp := redis.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: &redis.PasswordPolicy{
				Enabled:         ptrBoolPasswordPolicy(true),
				MinLetterCount:  ptrInt64PasswordPolicy(1),
				MinDigitCount:   ptrInt64PasswordPolicy(1),
				MinSpecialCount: ptrInt64PasswordPolicy(1),
				MinLength:       ptrInt64PasswordPolicy(8),
			},
			RequestId: ptrStrPasswordPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisInstancePasswordPolicyConfig()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":       "crs-cqdfdzvt",
		"enabled":           true,
		"min_letter_count":  1,
		"min_digit_count":   1,
		"min_special_count": 1,
		"min_length":        8,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-cqdfdzvt", d.Id())
	assert.Equal(t, true, d.Get("enabled"))
	assert.Equal(t, 8, d.Get("min_length"))
}

// TestRedisInstancePasswordPolicyConfig_Update_Success tests Update calls ModifyInstancePasswordPolicy and refreshes
func TestRedisInstancePasswordPolicyConfig_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisInstancePasswordPolicyConfig().client, "UseRedisClient", redisClient)

	// Mock ModifyInstancePasswordPolicyWithContext
	patches.ApplyMethodFunc(redisClient, "ModifyInstancePasswordPolicyWithContext", func(ctx context.Context, request *redis.ModifyInstancePasswordPolicyRequest) (*redis.ModifyInstancePasswordPolicyResponse, error) {
		assert.Equal(t, "crs-cqdfdzvt", *request.InstanceId)
		assert.NotNil(t, request.PasswordPolicy)
		assert.Equal(t, false, *request.PasswordPolicy.Enabled)
		resp := redis.NewModifyInstancePasswordPolicyResponse()
		resp.Response = &redis.ModifyInstancePasswordPolicyResponseParams{
			RequestId: ptrStrPasswordPolicy("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeInstancePasswordPolicyWithContext for Read
	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicyWithContext", func(ctx context.Context, request *redis.DescribeInstancePasswordPolicyRequest) (*redis.DescribeInstancePasswordPolicyResponse, error) {
		resp := redis.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: &redis.PasswordPolicy{
				Enabled:         ptrBoolPasswordPolicy(false),
				MinLetterCount:  ptrInt64PasswordPolicy(2),
				MinDigitCount:   ptrInt64PasswordPolicy(2),
				MinSpecialCount: ptrInt64PasswordPolicy(2),
				MinLength:       ptrInt64PasswordPolicy(16),
			},
			RequestId: ptrStrPasswordPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisInstancePasswordPolicyConfig()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":       "crs-cqdfdzvt",
		"enabled":           false,
		"min_letter_count":  2,
		"min_digit_count":   2,
		"min_special_count": 2,
		"min_length":        16,
	})
	d.SetId("crs-cqdfdzvt")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, false, d.Get("enabled"))
	assert.Equal(t, 16, d.Get("min_length"))
}

// TestRedisInstancePasswordPolicyConfig_Delete_NoOp tests Delete is a no-op
func TestRedisInstancePasswordPolicyConfig_Delete_NoOp(t *testing.T) {
	meta := newMockMetaForRedisInstancePasswordPolicyConfig()
	res := svccrs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-cqdfdzvt",
		"enabled":     true,
	})
	d.SetId("crs-cqdfdzvt")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-cqdfdzvt", d.Id())
}
