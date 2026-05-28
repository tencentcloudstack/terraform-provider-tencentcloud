package sqlserver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	sqlserverv20180328 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sqlserver"
)

// mockMetaForDbInstanceSslConfig implements tccommon.ProviderMeta
type mockMetaForDbInstanceSslConfig struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForDbInstanceSslConfig) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForDbInstanceSslConfig{}

func newMockMetaForDbInstanceSslConfig() *mockMetaForDbInstanceSslConfig {
	return &mockMetaForDbInstanceSslConfig{client: &connectivity.TencentCloudClient{}}
}

func ptrStrSslConfig(s string) *string {
	return &s
}

func ptrUint64SslConfig(v uint64) *uint64 {
	return &v
}

func ptrInt64SslConfig(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/sqlserver/ -run "TestDbInstanceSslConfig" -v -count=1 -gcflags="all=-l"

// TestDbInstanceSslConfig_Read_Success tests Read populates fields from SSLConfig
func TestDbInstanceSslConfig_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sqlserverClient := &sqlserverv20180328.Client{}
	patches.ApplyMethodReturn(newMockMetaForDbInstanceSslConfig().client, "UseSqlserverClient", sqlserverClient)

	patches.ApplyMethodFunc(sqlserverClient, "DescribeDBInstancesAttribute", func(request *sqlserverv20180328.DescribeDBInstancesAttributeRequest) (*sqlserverv20180328.DescribeDBInstancesAttributeResponse, error) {
		resp := sqlserverv20180328.NewDescribeDBInstancesAttributeResponse()
		resp.Response = &sqlserverv20180328.DescribeDBInstancesAttributeResponseParams{
			InstanceId: ptrStrSslConfig("mssql-gy1lc54f"),
			SSLConfig: &sqlserverv20180328.SSLConfig{
				Encryption:        ptrStrSslConfig("enable"),
				SSLValidityPeriod: ptrStrSslConfig("2026-12-31 23:59:59"),
				SSLValidity:       ptrUint64SslConfig(1),
				IsKMS:             ptrInt64SslConfig(0),
			},
			RequestId: ptrStrSslConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDbInstanceSslConfig()
	res := sqlserver.ResourceTencentCloudSqlserverDbInstanceSslConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "mssql-gy1lc54f",
		"encryption":  "enable",
	})
	d.SetId("mssql-gy1lc54f")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mssql-gy1lc54f", d.Get("instance_id"))
	assert.Equal(t, "enable", d.Get("encryption"))
	assert.Equal(t, "2026-12-31 23:59:59", d.Get("ssl_validity_period"))
	assert.Equal(t, 1, d.Get("ssl_validity"))
	assert.Equal(t, 0, d.Get("is_kms"))
}

// TestDbInstanceSslConfig_Read_NilSSLConfig tests Read handles nil SSLConfig
func TestDbInstanceSslConfig_Read_NilSSLConfig(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sqlserverClient := &sqlserverv20180328.Client{}
	patches.ApplyMethodReturn(newMockMetaForDbInstanceSslConfig().client, "UseSqlserverClient", sqlserverClient)

	patches.ApplyMethodFunc(sqlserverClient, "DescribeDBInstancesAttribute", func(request *sqlserverv20180328.DescribeDBInstancesAttributeRequest) (*sqlserverv20180328.DescribeDBInstancesAttributeResponse, error) {
		resp := sqlserverv20180328.NewDescribeDBInstancesAttributeResponse()
		resp.Response = &sqlserverv20180328.DescribeDBInstancesAttributeResponseParams{
			InstanceId: ptrStrSslConfig("mssql-gy1lc54f"),
			SSLConfig:  nil,
			RequestId:  ptrStrSslConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDbInstanceSslConfig()
	res := sqlserver.ResourceTencentCloudSqlserverDbInstanceSslConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "mssql-gy1lc54f",
		"encryption":  "enable",
	})
	d.SetId("mssql-gy1lc54f")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mssql-gy1lc54f", d.Get("instance_id"))
}

// TestDbInstanceSslConfig_Read_NilResponse tests Read handles error response
func TestDbInstanceSslConfig_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sqlserverClient := &sqlserverv20180328.Client{}
	patches.ApplyMethodReturn(newMockMetaForDbInstanceSslConfig().client, "UseSqlserverClient", sqlserverClient)

	patches.ApplyMethodFunc(sqlserverClient, "DescribeDBInstancesAttribute", func(request *sqlserverv20180328.DescribeDBInstancesAttributeRequest) (*sqlserverv20180328.DescribeDBInstancesAttributeResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Instance not found")
	})

	meta := newMockMetaForDbInstanceSslConfig()
	res := sqlserver.ResourceTencentCloudSqlserverDbInstanceSslConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "mssql-gy1lc54f",
		"encryption":  "enable",
	})
	d.SetId("mssql-gy1lc54f")

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestDbInstanceSslConfig_Create_Success tests Create sets ID and delegates to Update
func TestDbInstanceSslConfig_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sqlserverClient := &sqlserverv20180328.Client{}
	patches.ApplyMethodReturn(newMockMetaForDbInstanceSslConfig().client, "UseSqlserverClient", sqlserverClient)

	// Mock ModifyDBInstanceSSLWithContext
	patches.ApplyMethodFunc(sqlserverClient, "ModifyDBInstanceSSLWithContext", func(ctx context.Context, request *sqlserverv20180328.ModifyDBInstanceSSLRequest) (*sqlserverv20180328.ModifyDBInstanceSSLResponse, error) {
		assert.Equal(t, "enable", *request.Type)
		resp := sqlserverv20180328.NewModifyDBInstanceSSLResponse()
		resp.Response = &sqlserverv20180328.ModifyDBInstanceSSLResponseParams{
			FlowId:    ptrUint64SslConfig(12345),
			RequestId: ptrStrSslConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeFlowStatus for polling
	patches.ApplyMethodFunc(sqlserverClient, "DescribeFlowStatus", func(request *sqlserverv20180328.DescribeFlowStatusRequest) (*sqlserverv20180328.DescribeFlowStatusResponse, error) {
		resp := sqlserverv20180328.NewDescribeFlowStatusResponse()
		status := int64(0) // SQLSERVER_TASK_SUCCESS
		resp.Response = &sqlserverv20180328.DescribeFlowStatusResponseParams{
			Status:    &status,
			RequestId: ptrStrSslConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeDBInstancesAttribute for Read
	patches.ApplyMethodFunc(sqlserverClient, "DescribeDBInstancesAttribute", func(request *sqlserverv20180328.DescribeDBInstancesAttributeRequest) (*sqlserverv20180328.DescribeDBInstancesAttributeResponse, error) {
		resp := sqlserverv20180328.NewDescribeDBInstancesAttributeResponse()
		resp.Response = &sqlserverv20180328.DescribeDBInstancesAttributeResponseParams{
			InstanceId: ptrStrSslConfig("mssql-gy1lc54f"),
			SSLConfig: &sqlserverv20180328.SSLConfig{
				Encryption:        ptrStrSslConfig("enable"),
				SSLValidityPeriod: ptrStrSslConfig("2026-12-31 23:59:59"),
				SSLValidity:       ptrUint64SslConfig(1),
				IsKMS:             ptrInt64SslConfig(0),
			},
			RequestId: ptrStrSslConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDbInstanceSslConfig()
	res := sqlserver.ResourceTencentCloudSqlserverDbInstanceSslConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "mssql-gy1lc54f",
		"encryption":  "enable",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mssql-gy1lc54f", d.Id())
	assert.Equal(t, "enable", d.Get("encryption"))
}

// TestDbInstanceSslConfig_Update_Success tests Update calls ModifyDBInstanceSSL and polls encryption status
func TestDbInstanceSslConfig_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sqlserverClient := &sqlserverv20180328.Client{}
	patches.ApplyMethodReturn(newMockMetaForDbInstanceSslConfig().client, "UseSqlserverClient", sqlserverClient)

	// Mock ModifyDBInstanceSSLWithContext
	patches.ApplyMethodFunc(sqlserverClient, "ModifyDBInstanceSSLWithContext", func(ctx context.Context, request *sqlserverv20180328.ModifyDBInstanceSSLRequest) (*sqlserverv20180328.ModifyDBInstanceSSLResponse, error) {
		assert.Equal(t, "disable", *request.Type)
		resp := sqlserverv20180328.NewModifyDBInstanceSSLResponse()
		resp.Response = &sqlserverv20180328.ModifyDBInstanceSSLResponseParams{
			FlowId:    ptrUint64SslConfig(12345),
			RequestId: ptrStrSslConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeFlowStatus for polling
	patches.ApplyMethodFunc(sqlserverClient, "DescribeFlowStatus", func(request *sqlserverv20180328.DescribeFlowStatusRequest) (*sqlserverv20180328.DescribeFlowStatusResponse, error) {
		resp := sqlserverv20180328.NewDescribeFlowStatusResponse()
		status := int64(0) // SQLSERVER_TASK_SUCCESS
		resp.Response = &sqlserverv20180328.DescribeFlowStatusResponseParams{
			Status:    &status,
			RequestId: ptrStrSslConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeDBInstancesAttribute for Read
	patches.ApplyMethodFunc(sqlserverClient, "DescribeDBInstancesAttribute", func(request *sqlserverv20180328.DescribeDBInstancesAttributeRequest) (*sqlserverv20180328.DescribeDBInstancesAttributeResponse, error) {
		resp := sqlserverv20180328.NewDescribeDBInstancesAttributeResponse()
		resp.Response = &sqlserverv20180328.DescribeDBInstancesAttributeResponseParams{
			InstanceId: ptrStrSslConfig("mssql-gy1lc54f"),
			SSLConfig: &sqlserverv20180328.SSLConfig{
				Encryption:  ptrStrSslConfig("disable"),
				SSLValidity: ptrUint64SslConfig(0),
				IsKMS:       ptrInt64SslConfig(0),
			},
			RequestId: ptrStrSslConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDbInstanceSslConfig()
	res := sqlserver.ResourceTencentCloudSqlserverDbInstanceSslConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "mssql-gy1lc54f",
		"encryption":  "disable",
	})
	d.SetId("mssql-gy1lc54f")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "disable", d.Get("encryption"))
}

// TestDbInstanceSslConfig_Delete_NoOp tests Delete is a no-op
func TestDbInstanceSslConfig_Delete_NoOp(t *testing.T) {
	meta := newMockMetaForDbInstanceSslConfig()
	res := sqlserver.ResourceTencentCloudSqlserverDbInstanceSslConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "mssql-gy1lc54f",
		"encryption":  "enable",
	})
	d.SetId("mssql-gy1lc54f")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mssql-gy1lc54f", d.Id())
}

// TestDbInstanceSslConfig_Schema tests the schema definition
func TestDbInstanceSslConfig_Schema(t *testing.T) {
	res := sqlserver.ResourceTencentCloudSqlserverDbInstanceSslConfig()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Importer)

	// Required fields
	assert.Contains(t, res.Schema, "instance_id")
	assert.True(t, res.Schema["instance_id"].Required)
	assert.True(t, res.Schema["instance_id"].ForceNew)

	assert.Contains(t, res.Schema, "encryption")
	assert.True(t, res.Schema["encryption"].Required)

	// Optional fields
	assert.Contains(t, res.Schema, "is_kms")
	assert.True(t, res.Schema["is_kms"].Optional)

	assert.Contains(t, res.Schema, "cmk_id")
	assert.True(t, res.Schema["cmk_id"].Optional)

	assert.Contains(t, res.Schema, "cmk_region")
	assert.True(t, res.Schema["cmk_region"].Optional)

	// Computed fields
	assert.Contains(t, res.Schema, "ssl_validity_period")
	assert.True(t, res.Schema["ssl_validity_period"].Computed)

	assert.Contains(t, res.Schema, "ssl_validity")
	assert.True(t, res.Schema["ssl_validity"].Computed)
}
