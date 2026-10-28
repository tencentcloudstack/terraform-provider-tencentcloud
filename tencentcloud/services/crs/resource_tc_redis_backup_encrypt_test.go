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

// mockMetaForRedisBackupEncrypt implements tccommon.ProviderMeta
type mockMetaForRedisBackupEncrypt struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForRedisBackupEncrypt) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForRedisBackupEncrypt{}

func newMockMetaForRedisBackupEncrypt() *mockMetaForRedisBackupEncrypt {
	return &mockMetaForRedisBackupEncrypt{client: &connectivity.TencentCloudClient{}}
}

func ptrBoolBackupEncrypt(v bool) *bool {
	return &v
}

func ptrInt64BackupEncrypt(v int64) *int64 {
	return &v
}

func ptrStrBackupEncrypt(s string) *string {
	return &s
}

// go test ./tencentcloud/services/crs/ -run "TestRedisBackupEncrypt" -v -count=1 -gcflags="all=-l"

// TestRedisBackupEncrypt_Schema tests the schema definition
func TestRedisBackupEncrypt_Schema(t *testing.T) {
	res := svccrs.ResourceTencentCloudRedisBackupEncrypt()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Importer)

	// Required fields
	assert.Contains(t, res.Schema, "instance_id")
	assert.True(t, res.Schema["instance_id"].Required)
	assert.True(t, res.Schema["instance_id"].ForceNew)

	assert.Contains(t, res.Schema, "backup_encrypt")
	assert.True(t, res.Schema["backup_encrypt"].Required)
	assert.Equal(t, schema.TypeBool, res.Schema["backup_encrypt"].Type)

	// Optional + Computed fields
	for _, field := range []string{"encrypt_key_type", "kms_key_id", "kms_region"} {
		assert.Contains(t, res.Schema, field)
		assert.True(t, res.Schema[field].Optional)
		assert.True(t, res.Schema[field].Computed)
	}
}

// TestRedisBackupEncrypt_Read_Success tests Read populates fields from DescribeBackupEncryptConfig
func TestRedisBackupEncrypt_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(ctx context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		assert.Equal(t, "crs-cqdfdzvt", *request.InstanceId)
		resp := redis.NewDescribeBackupEncryptConfigResponse()
		resp.Response = &redis.DescribeBackupEncryptConfigResponseParams{
			BackupEncrypt:  ptrBoolBackupEncrypt(true),
			EncryptKeyType: ptrInt64BackupEncrypt(1),
			KmsKeyId:       ptrStrBackupEncrypt("xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx"),
			KmsRegion:      ptrStrBackupEncrypt("ap-guangzhou"),
			RequestId:      ptrStrBackupEncrypt("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisBackupEncrypt()
	res := svccrs.ResourceTencentCloudRedisBackupEncrypt()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-cqdfdzvt",
	})
	d.SetId("crs-cqdfdzvt")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-cqdfdzvt", d.Id())
	assert.Equal(t, true, d.Get("backup_encrypt"))
	assert.Equal(t, 1, d.Get("encrypt_key_type"))
	assert.Equal(t, "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx", d.Get("kms_key_id"))
	assert.Equal(t, "ap-guangzhou", d.Get("kms_region"))
}

// TestRedisBackupEncrypt_Read_EmptyConfig tests Read clears state when response is nil
func TestRedisBackupEncrypt_Read_EmptyConfig(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(ctx context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		resp := redis.NewDescribeBackupEncryptConfigResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaForRedisBackupEncrypt()
	res := svccrs.ResourceTencentCloudRedisBackupEncrypt()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-cqdfdzvt",
	})
	d.SetId("crs-cqdfdzvt")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestRedisBackupEncrypt_Read_InstanceNotExists tests Read returns error when instance not found
func TestRedisBackupEncrypt_Read_InstanceNotExists(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(ctx context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		return nil, sdkErrors.NewTencentCloudSDKError("ResourceNotFound.InstanceNotExists", "Instance not found", "fake-request-id")
	})

	meta := newMockMetaForRedisBackupEncrypt()
	res := svccrs.ResourceTencentCloudRedisBackupEncrypt()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "crs-cqdfdzvt",
	})
	d.SetId("crs-cqdfdzvt")

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestRedisBackupEncrypt_Create_Success tests Create sets ID and delegates to Update
func TestRedisBackupEncrypt_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	// Mock ModifyBackupEncryptWithContext for Update
	patches.ApplyMethodFunc(redisClient, "ModifyBackupEncryptWithContext", func(ctx context.Context, request *redis.ModifyBackupEncryptRequest) (*redis.ModifyBackupEncryptResponse, error) {
		assert.Equal(t, "crs-cqdfdzvt", *request.InstanceId)
		assert.NotNil(t, request.BackupEncrypt)
		assert.Equal(t, true, *request.BackupEncrypt)
		resp := redis.NewModifyBackupEncryptResponse()
		resp.Response = &redis.ModifyBackupEncryptResponseParams{
			RequestId: ptrStrBackupEncrypt("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeBackupEncryptConfigWithContext for Read
	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(ctx context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		resp := redis.NewDescribeBackupEncryptConfigResponse()
		resp.Response = &redis.DescribeBackupEncryptConfigResponseParams{
			BackupEncrypt:  ptrBoolBackupEncrypt(true),
			EncryptKeyType: ptrInt64BackupEncrypt(1),
			KmsKeyId:       ptrStrBackupEncrypt("xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx"),
			KmsRegion:      ptrStrBackupEncrypt("ap-guangzhou"),
			RequestId:      ptrStrBackupEncrypt("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisBackupEncrypt()
	res := svccrs.ResourceTencentCloudRedisBackupEncrypt()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":      "crs-cqdfdzvt",
		"backup_encrypt":   true,
		"encrypt_key_type": 1,
		"kms_key_id":       "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx",
		"kms_region":       "ap-guangzhou",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-cqdfdzvt", d.Id())
	assert.Equal(t, true, d.Get("backup_encrypt"))
	assert.Equal(t, 1, d.Get("encrypt_key_type"))
	assert.Equal(t, "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx", d.Get("kms_key_id"))
	assert.Equal(t, "ap-guangzhou", d.Get("kms_region"))
}

// TestRedisBackupEncrypt_Update_Success tests Update calls ModifyBackupEncrypt and refreshes
func TestRedisBackupEncrypt_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	redisClient := &redis.Client{}
	patches.ApplyMethodReturn(newMockMetaForRedisBackupEncrypt().client, "UseRedisClient", redisClient)

	// Mock ModifyBackupEncryptWithContext
	patches.ApplyMethodFunc(redisClient, "ModifyBackupEncryptWithContext", func(ctx context.Context, request *redis.ModifyBackupEncryptRequest) (*redis.ModifyBackupEncryptResponse, error) {
		assert.Equal(t, "crs-cqdfdzvt", *request.InstanceId)
		assert.NotNil(t, request.BackupEncrypt)
		assert.Equal(t, false, *request.BackupEncrypt)
		resp := redis.NewModifyBackupEncryptResponse()
		resp.Response = &redis.ModifyBackupEncryptResponseParams{
			RequestId: ptrStrBackupEncrypt("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeBackupEncryptConfigWithContext for Read
	patches.ApplyMethodFunc(redisClient, "DescribeBackupEncryptConfigWithContext", func(ctx context.Context, request *redis.DescribeBackupEncryptConfigRequest) (*redis.DescribeBackupEncryptConfigResponse, error) {
		resp := redis.NewDescribeBackupEncryptConfigResponse()
		resp.Response = &redis.DescribeBackupEncryptConfigResponseParams{
			BackupEncrypt:  ptrBoolBackupEncrypt(false),
			EncryptKeyType: ptrInt64BackupEncrypt(0),
			KmsKeyId:       nil,
			KmsRegion:      nil,
			RequestId:      ptrStrBackupEncrypt("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRedisBackupEncrypt()
	res := svccrs.ResourceTencentCloudRedisBackupEncrypt()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "crs-cqdfdzvt",
		"backup_encrypt": false,
	})
	d.SetId("crs-cqdfdzvt")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, false, d.Get("backup_encrypt"))
}

// TestRedisBackupEncrypt_Delete_NoOp tests Delete is a no-op
func TestRedisBackupEncrypt_Delete_NoOp(t *testing.T) {
	meta := newMockMetaForRedisBackupEncrypt()
	res := svccrs.ResourceTencentCloudRedisBackupEncrypt()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":    "crs-cqdfdzvt",
		"backup_encrypt": true,
	})
	d.SetId("crs-cqdfdzvt")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "crs-cqdfdzvt", d.Id())
}
