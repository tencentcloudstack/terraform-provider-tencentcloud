package sts_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	stsv20180813 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sts"
)

type mockMetaAssumeRole struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaAssumeRole) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaAssumeRole{}

func newMockMetaAssumeRole() *mockMetaAssumeRole {
	return &mockMetaAssumeRole{client: &connectivity.TencentCloudClient{}}
}

func ptrStringAssumeRole(s string) *string {
	return &s
}

func ptrInt64AssumeRole(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/sts/ -run "TestAssumeRoleOperation" -v -count=1 -gcflags="all=-l"

// TestAssumeRoleOperation_Success tests Create with required parameters only
func TestAssumeRoleOperation_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	stsClient := &stsv20180813.Client{}
	patches.ApplyMethodReturn(newMockMetaAssumeRole().client, "UseStsClient", stsClient)

	patches.ApplyMethodFunc(stsClient, "AssumeRoleWithContext", func(ctx context.Context, request *stsv20180813.AssumeRoleRequest) (*stsv20180813.AssumeRoleResponse, error) {
		assert.Equal(t, "qcs::cam::uin/100000000001:roleName/testRoleName", *request.RoleArn)
		assert.Equal(t, "test-session", *request.RoleSessionName)
		resp := stsv20180813.NewAssumeRoleResponse()
		resp.Response = &stsv20180813.AssumeRoleResponseParams{
			Credentials: &stsv20180813.Credentials{
				Token:        ptrStringAssumeRole("fake-token"),
				TmpSecretId:  ptrStringAssumeRole("fake-secret-id"),
				TmpSecretKey: ptrStringAssumeRole("fake-secret-key"),
			},
			ExpiredTime: ptrInt64AssumeRole(1700000000),
			Expiration:  ptrStringAssumeRole("2023-11-14T22:13:20Z"),
			RequestId:   ptrStringAssumeRole("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaAssumeRole()
	res := sts.ResourceTencentCloudStsAssumeRoleOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"role_arn":          "qcs::cam::uin/100000000001:roleName/testRoleName",
		"role_session_name": "test-session",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	credentials := d.Get("credentials").([]interface{})
	assert.Len(t, credentials, 1)
	credMap := credentials[0].(map[string]interface{})
	assert.Equal(t, "fake-token", credMap["token"])
	assert.Equal(t, "fake-secret-id", credMap["tmp_secret_id"])
	assert.Equal(t, "fake-secret-key", credMap["tmp_secret_key"])

	assert.Equal(t, 1700000000, d.Get("expired_time").(int))
	assert.Equal(t, "2023-11-14T22:13:20Z", d.Get("expiration").(string))
}

// TestAssumeRoleOperation_AllParams tests Create with all parameters
func TestAssumeRoleOperation_AllParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	stsClient := &stsv20180813.Client{}
	patches.ApplyMethodReturn(newMockMetaAssumeRole().client, "UseStsClient", stsClient)

	patches.ApplyMethodFunc(stsClient, "AssumeRoleWithContext", func(ctx context.Context, request *stsv20180813.AssumeRoleRequest) (*stsv20180813.AssumeRoleResponse, error) {
		assert.Equal(t, "qcs::cam::uin/100000000001:roleName/testRoleName", *request.RoleArn)
		assert.Equal(t, "test-session", *request.RoleSessionName)
		assert.NotNil(t, request.DurationSeconds)
		assert.Equal(t, uint64(7200), *request.DurationSeconds)
		assert.NotNil(t, request.Policy)
		assert.Equal(t, "test-policy", *request.Policy)
		assert.NotNil(t, request.ExternalId)
		assert.Equal(t, "ext-id", *request.ExternalId)
		assert.NotNil(t, request.SourceIdentity)
		assert.Equal(t, "source-id", *request.SourceIdentity)
		assert.NotNil(t, request.SerialNumber)
		assert.Equal(t, "serial-123", *request.SerialNumber)
		assert.NotNil(t, request.TokenCode)
		assert.Equal(t, "123456", *request.TokenCode)
		assert.Len(t, request.Tags, 1)
		assert.Equal(t, "env", *request.Tags[0].Key)
		assert.Equal(t, "test", *request.Tags[0].Value)

		resp := stsv20180813.NewAssumeRoleResponse()
		resp.Response = &stsv20180813.AssumeRoleResponseParams{
			Credentials: &stsv20180813.Credentials{
				Token:        ptrStringAssumeRole("fake-token-2"),
				TmpSecretId:  ptrStringAssumeRole("fake-secret-id-2"),
				TmpSecretKey: ptrStringAssumeRole("fake-secret-key-2"),
			},
			ExpiredTime: ptrInt64AssumeRole(1700003600),
			Expiration:  ptrStringAssumeRole("2023-11-14T23:13:20Z"),
			RequestId:   ptrStringAssumeRole("fake-request-id-2"),
		}
		return resp, nil
	})

	meta := newMockMetaAssumeRole()
	res := sts.ResourceTencentCloudStsAssumeRoleOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"role_arn":          "qcs::cam::uin/100000000001:roleName/testRoleName",
		"role_session_name": "test-session",
		"duration_seconds":  7200,
		"policy":            "test-policy",
		"external_id":       "ext-id",
		"source_identity":   "source-id",
		"serial_number":     "serial-123",
		"token_code":        "123456",
		"tags": []interface{}{
			map[string]interface{}{
				"key":   "env",
				"value": "test",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
}

// TestAssumeRoleOperation_APIError tests Create when API returns error
func TestAssumeRoleOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	stsClient := &stsv20180813.Client{}
	patches.ApplyMethodReturn(newMockMetaAssumeRole().client, "UseStsClient", stsClient)

	patches.ApplyMethodFunc(stsClient, "AssumeRoleWithContext", func(ctx context.Context, request *stsv20180813.AssumeRoleRequest) (*stsv20180813.AssumeRoleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound.RoleNotFound, Message=Role not found")
	})

	meta := newMockMetaAssumeRole()
	res := sts.ResourceTencentCloudStsAssumeRoleOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"role_arn":          "qcs::cam::uin/100000000001:roleName/invalidRole",
		"role_session_name": "test-session",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestAssumeRoleOperation_Read tests Read is no-op
func TestAssumeRoleOperation_Read(t *testing.T) {
	res := sts.ResourceTencentCloudStsAssumeRoleOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"role_arn":          "qcs::cam::uin/100000000001:roleName/testRoleName",
		"role_session_name": "test-session",
	})
	d.SetId("test-id")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestAssumeRoleOperation_Delete tests Delete is no-op
func TestAssumeRoleOperation_Delete(t *testing.T) {
	res := sts.ResourceTencentCloudStsAssumeRoleOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"role_arn":          "qcs::cam::uin/100000000001:roleName/testRoleName",
		"role_session_name": "test-session",
	})
	d.SetId("test-id")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestAssumeRoleOperation_Schema validates schema definition
func TestAssumeRoleOperation_Schema(t *testing.T) {
	res := sts.ResourceTencentCloudStsAssumeRoleOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "role_arn")
	assert.Contains(t, res.Schema, "role_session_name")
	assert.Contains(t, res.Schema, "duration_seconds")
	assert.Contains(t, res.Schema, "policy")
	assert.Contains(t, res.Schema, "external_id")
	assert.Contains(t, res.Schema, "tags")
	assert.Contains(t, res.Schema, "source_identity")
	assert.Contains(t, res.Schema, "serial_number")
	assert.Contains(t, res.Schema, "token_code")
	assert.Contains(t, res.Schema, "credentials")
	assert.Contains(t, res.Schema, "expired_time")
	assert.Contains(t, res.Schema, "expiration")

	roleArn := res.Schema["role_arn"]
	assert.Equal(t, schema.TypeString, roleArn.Type)
	assert.True(t, roleArn.Required)
	assert.True(t, roleArn.ForceNew)

	roleSessionName := res.Schema["role_session_name"]
	assert.Equal(t, schema.TypeString, roleSessionName.Type)
	assert.True(t, roleSessionName.Required)
	assert.True(t, roleSessionName.ForceNew)

	durationSeconds := res.Schema["duration_seconds"]
	assert.Equal(t, schema.TypeInt, durationSeconds.Type)
	assert.True(t, durationSeconds.Optional)
	assert.True(t, durationSeconds.ForceNew)

	policy := res.Schema["policy"]
	assert.Equal(t, schema.TypeString, policy.Type)
	assert.True(t, policy.Optional)
	assert.True(t, policy.ForceNew)

	credentials := res.Schema["credentials"]
	assert.Equal(t, schema.TypeList, credentials.Type)
	assert.True(t, credentials.Computed)
}
