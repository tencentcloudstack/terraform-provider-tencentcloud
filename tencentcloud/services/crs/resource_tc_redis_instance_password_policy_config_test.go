package crs_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	redis_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"
)

type mockMetaForRedisPasswordPolicy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForRedisPasswordPolicy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForRedisPasswordPolicy{}

func newMockMetaForRedisPasswordPolicy() *mockMetaForRedisPasswordPolicy {
	return &mockMetaForRedisPasswordPolicy{client: &connectivity.TencentCloudClient{}}
}

func ptrStringRedisPwdPolicy(s string) *string { return &s }
func ptrBoolRedisPwdPolicy(b bool) *bool       { return &b }
func ptrInt64RedisPwdPolicy(v int64) *int64    { return &v }

func TestRedisInstancePasswordPolicyConfig_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisPasswordPolicy().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyInstancePasswordPolicy", func(request *redis_sdk.ModifyInstancePasswordPolicyRequest) (*redis_sdk.ModifyInstancePasswordPolicyResponse, error) {
		resp := redis_sdk.NewModifyInstancePasswordPolicyResponse()
		resp.Response = &redis_sdk.ModifyInstancePasswordPolicyResponseParams{
			RequestId: ptrStringRedisPwdPolicy("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicy", func(request *redis_sdk.DescribeInstancePasswordPolicyRequest) (*redis_sdk.DescribeInstancePasswordPolicyResponse, error) {
		resp := redis_sdk.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis_sdk.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: &redis_sdk.PasswordPolicy{
				Enabled:         ptrBoolRedisPwdPolicy(true),
				MinLetterCount:  ptrInt64RedisPwdPolicy(1),
				MinDigitCount:   ptrInt64RedisPwdPolicy(1),
				MinSpecialCount: ptrInt64RedisPwdPolicy(1),
				MinLength:       ptrInt64RedisPwdPolicy(8),
			},
			RequestId: ptrStringRedisPwdPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisPasswordPolicy()
	res := crs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-test1234",
		"password_policy": []interface{}{
			map[string]interface{}{
				"enabled":           true,
				"min_letter_count":  1,
				"min_digit_count":   1,
				"min_special_count": 1,
				"min_length":        8,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-test1234", d.Id())
}

func TestRedisInstancePasswordPolicyConfig_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisPasswordPolicy().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicy", func(request *redis_sdk.DescribeInstancePasswordPolicyRequest) (*redis_sdk.DescribeInstancePasswordPolicyResponse, error) {
		resp := redis_sdk.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis_sdk.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: &redis_sdk.PasswordPolicy{
				Enabled:         ptrBoolRedisPwdPolicy(true),
				MinLetterCount:  ptrInt64RedisPwdPolicy(2),
				MinDigitCount:   ptrInt64RedisPwdPolicy(2),
				MinSpecialCount: ptrInt64RedisPwdPolicy(2),
				MinLength:       ptrInt64RedisPwdPolicy(10),
			},
			RequestId: ptrStringRedisPwdPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisPasswordPolicy()
	res := crs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-test1234",
		"password_policy": []interface{}{
			map[string]interface{}{
				"enabled":           true,
				"min_letter_count":  1,
				"min_digit_count":   1,
				"min_special_count": 1,
				"min_length":        8,
			},
		},
	})
	d.SetId("crs-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-test1234", d.Get("instance_id"))

	policyList := d.Get("password_policy").([]interface{})
	assert.NotEmpty(t, policyList)
	policyMap := policyList[0].(map[string]interface{})
	assert.Equal(t, true, policyMap["enabled"])
	assert.Equal(t, 2, policyMap["min_letter_count"])
	assert.Equal(t, 2, policyMap["min_digit_count"])
	assert.Equal(t, 2, policyMap["min_special_count"])
	assert.Equal(t, 10, policyMap["min_length"])
}

func TestRedisInstancePasswordPolicyConfig_Read_Nil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisPasswordPolicy().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicy", func(request *redis_sdk.DescribeInstancePasswordPolicyRequest) (*redis_sdk.DescribeInstancePasswordPolicyResponse, error) {
		resp := redis_sdk.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis_sdk.DescribeInstancePasswordPolicyResponseParams{
			RequestId: ptrStringRedisPwdPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisPasswordPolicy()
	res := crs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-test1234",
		"password_policy": []interface{}{
			map[string]interface{}{
				"enabled":           true,
				"min_letter_count":  1,
				"min_digit_count":   1,
				"min_special_count": 1,
				"min_length":        8,
			},
		},
	})
	d.SetId("crs-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestRedisInstancePasswordPolicyConfig_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisPasswordPolicy().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyInstancePasswordPolicy", func(request *redis_sdk.ModifyInstancePasswordPolicyRequest) (*redis_sdk.ModifyInstancePasswordPolicyResponse, error) {
		resp := redis_sdk.NewModifyInstancePasswordPolicyResponse()
		resp.Response = &redis_sdk.ModifyInstancePasswordPolicyResponseParams{
			RequestId: ptrStringRedisPwdPolicy("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(redisClient, "DescribeInstancePasswordPolicy", func(request *redis_sdk.DescribeInstancePasswordPolicyRequest) (*redis_sdk.DescribeInstancePasswordPolicyResponse, error) {
		resp := redis_sdk.NewDescribeInstancePasswordPolicyResponse()
		resp.Response = &redis_sdk.DescribeInstancePasswordPolicyResponseParams{
			PasswordPolicy: &redis_sdk.PasswordPolicy{
				Enabled:         ptrBoolRedisPwdPolicy(false),
				MinLetterCount:  ptrInt64RedisPwdPolicy(1),
				MinDigitCount:   ptrInt64RedisPwdPolicy(1),
				MinSpecialCount: ptrInt64RedisPwdPolicy(1),
				MinLength:       ptrInt64RedisPwdPolicy(8),
			},
			RequestId: ptrStringRedisPwdPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisPasswordPolicy()
	res := crs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-test1234",
		"password_policy": []interface{}{
			map[string]interface{}{
				"enabled":           false,
				"min_letter_count":  1,
				"min_digit_count":   1,
				"min_special_count": 1,
				"min_length":        8,
			},
		},
	})
	d.SetId("crs-test1234")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestRedisInstancePasswordPolicyConfig_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisPasswordPolicy().client, "UseRedisClient", redisClient)

	meta := newMockMetaForRedisPasswordPolicy()
	res := crs.ResourceTencentCloudRedisInstancePasswordPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-test1234",
		"password_policy": []interface{}{
			map[string]interface{}{
				"enabled":           true,
				"min_letter_count":  1,
				"min_digit_count":   1,
				"min_special_count": 1,
				"min_length":        8,
			},
		},
	})
	d.SetId("crs-test1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
