package cdb_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"
)

// mockMeta implements tccommon.ProviderMeta for mysql_backup tests
type mockMysqlBackupMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMysqlBackupMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMysqlBackupMeta{}

func newMockMysqlBackupMeta() *mockMysqlBackupMeta {
	return &mockMysqlBackupMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrInt64(i int64) *int64 {
	return &i
}

func ptrUint64(i uint64) *uint64 {
	return &i
}

// TestMysqlBackupAttachmentCreate_Success tests successful backup creation
func TestMysqlBackupAttachmentCreate_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mysqlClient := &mysql.Client{}
	patches.ApplyMethodReturn(newMockMysqlBackupMeta().client, "UseMysqlClient", mysqlClient)

	patches.ApplyMethodFunc(mysqlClient, "CreateBackup", func(request *mysql.CreateBackupRequest) (*mysql.CreateBackupResponse, error) {
		assert.Equal(t, "cdb-test123", *request.InstanceId)
		assert.Equal(t, "physical", *request.BackupMethod)
		assert.Equal(t, "tf-test-backup", *request.ManualBackupName)

		resp := mysql.NewCreateBackupResponse()
		resp.Response = &mysql.CreateBackupResponseParams{
			BackupId:  ptrUint64(12345),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeBackups to return matching backup info for the Read call
	patches.ApplyMethodFunc(mysqlClient, "DescribeBackups", func(request *mysql.DescribeBackupsRequest) (*mysql.DescribeBackupsResponse, error) {
		resp := mysql.NewDescribeBackupsResponse()
		resp.Response = &mysql.DescribeBackupsResponseParams{
			TotalCount: ptrInt64(1),
			Items: []*mysql.BackupInfo{
				{
					BackupId:         ptrInt64(12345),
					Type:             ptrString("physical"),
					ManualBackupName: ptrString("tf-test-backup"),
					InstanceId:       ptrString("cdb-test123"),
					EncryptionFlag:   ptrString("off"),
					Way:              ptrString("manual"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMysqlBackupMeta()
	res := cdb.ResourceTencentCloudMysqlBackupAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":        "cdb-test123",
		"backup_method":      "physical",
		"manual_backup_name": "tf-test-backup",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "12345#cdb-test123", d.Id())
	assert.Equal(t, "cdb-test123", d.Get("instance_id"))
	assert.Equal(t, int64(12345), int64(d.Get("backup_id").(int)))
}

// TestMysqlBackupAttachmentCreate_LogicalWithDbTable tests logical backup with db/table list
func TestMysqlBackupAttachmentCreate_LogicalWithDbTable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mysqlClient := &mysql.Client{}
	patches.ApplyMethodReturn(newMockMysqlBackupMeta().client, "UseMysqlClient", mysqlClient)

	patches.ApplyMethodFunc(mysqlClient, "CreateBackup", func(request *mysql.CreateBackupRequest) (*mysql.CreateBackupResponse, error) {
		assert.Equal(t, "cdb-test456", *request.InstanceId)
		assert.Equal(t, "logical", *request.BackupMethod)
		assert.Len(t, request.BackupDBTableList, 2)
		assert.Equal(t, "db1", *request.BackupDBTableList[0].Db)
		assert.Equal(t, "tb1", *request.BackupDBTableList[0].Table)
		assert.Equal(t, "db2", *request.BackupDBTableList[1].Db)
		assert.Nil(t, request.BackupDBTableList[1].Table)

		resp := mysql.NewCreateBackupResponse()
		resp.Response = &mysql.CreateBackupResponseParams{
			BackupId:  ptrUint64(67890),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(mysqlClient, "DescribeBackups", func(request *mysql.DescribeBackupsRequest) (*mysql.DescribeBackupsResponse, error) {
		resp := mysql.NewDescribeBackupsResponse()
		resp.Response = &mysql.DescribeBackupsResponseParams{
			TotalCount: ptrInt64(1),
			Items: []*mysql.BackupInfo{
				{
					BackupId:   ptrInt64(67890),
					Type:       ptrString("logical"),
					InstanceId: ptrString("cdb-test456"),
					Way:        ptrString("manual"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMysqlBackupMeta()
	res := cdb.ResourceTencentCloudMysqlBackupAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "cdb-test456",
		"backup_method": "logical",
		"backup_db_table_list": []interface{}{
			map[string]interface{}{
				"database": "db1",
				"table":    "tb1",
			},
			map[string]interface{}{
				"database": "db2",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "67890#cdb-test456", d.Id())
}

// TestMysqlBackupAttachmentCreate_APIError tests API error handling during create
func TestMysqlBackupAttachmentCreate_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mysqlClient := &mysql.Client{}
	patches.ApplyMethodReturn(newMockMysqlBackupMeta().client, "UseMysqlClient", mysqlClient)

	patches.ApplyMethodFunc(mysqlClient, "CreateBackup", func(request *mysql.CreateBackupRequest) (*mysql.CreateBackupResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation, Message=Instance state not support backup")
	})

	meta := newMockMysqlBackupMeta()
	res := cdb.ResourceTencentCloudMysqlBackupAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "cdb-invalid",
		"backup_method": "physical",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FailedOperation")
}

// TestMysqlBackupAttachmentCreate_EmptyBackupId tests empty backup ID handling
func TestMysqlBackupAttachmentCreate_EmptyBackupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mysqlClient := &mysql.Client{}
	patches.ApplyMethodReturn(newMockMysqlBackupMeta().client, "UseMysqlClient", mysqlClient)

	patches.ApplyMethodFunc(mysqlClient, "CreateBackup", func(request *mysql.CreateBackupRequest) (*mysql.CreateBackupResponse, error) {
		resp := mysql.NewCreateBackupResponse()
		resp.Response = &mysql.CreateBackupResponseParams{
			BackupId:  nil,
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMysqlBackupMeta()
	res := cdb.ResourceTencentCloudMysqlBackupAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "cdb-test789",
		"backup_method": "physical",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BackupId is empty")
}

// TestMysqlBackupAttachmentRead_NotFound tests Read clears ID when backup not found
func TestMysqlBackupAttachmentRead_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mysqlClient := &mysql.Client{}
	patches.ApplyMethodReturn(newMockMysqlBackupMeta().client, "UseMysqlClient", mysqlClient)

	patches.ApplyMethodFunc(mysqlClient, "DescribeBackups", func(request *mysql.DescribeBackupsRequest) (*mysql.DescribeBackupsResponse, error) {
		resp := mysql.NewDescribeBackupsResponse()
		resp.Response = &mysql.DescribeBackupsResponseParams{
			TotalCount: ptrInt64(0),
			Items:      []*mysql.BackupInfo{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMysqlBackupMeta()
	res := cdb.ResourceTencentCloudMysqlBackupAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "cdb-test123",
		"backup_method": "physical",
	})
	d.SetId("99999#cdb-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestMysqlBackupAttachmentRead_Success tests Read populates fields
func TestMysqlBackupAttachmentRead_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mysqlClient := &mysql.Client{}
	patches.ApplyMethodReturn(newMockMysqlBackupMeta().client, "UseMysqlClient", mysqlClient)

	patches.ApplyMethodFunc(mysqlClient, "DescribeBackups", func(request *mysql.DescribeBackupsRequest) (*mysql.DescribeBackupsResponse, error) {
		assert.Equal(t, "cdb-test123", *request.InstanceId)
		resp := mysql.NewDescribeBackupsResponse()
		resp.Response = &mysql.DescribeBackupsResponseParams{
			TotalCount: ptrInt64(1),
			Items: []*mysql.BackupInfo{
				{
					BackupId:         ptrInt64(12345),
					Type:             ptrString("physical"),
					ManualBackupName: ptrString("tf-test-backup"),
					EncryptionFlag:   ptrString("on"),
					InstanceId:       ptrString("cdb-test123"),
					Way:              ptrString("manual"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMysqlBackupMeta()
	res := cdb.ResourceTencentCloudMysqlBackupAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "cdb-test123",
		"backup_method": "physical",
	})
	d.SetId("12345#cdb-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cdb-test123", d.Get("instance_id"))
	assert.Equal(t, int64(12345), int64(d.Get("backup_id").(int)))
	assert.Equal(t, "physical", d.Get("backup_method"))
	assert.Equal(t, "tf-test-backup", d.Get("manual_backup_name"))
	assert.Equal(t, "on", d.Get("encryption_flag"))
}

// TestMysqlBackupAttachmentDelete_Success tests successful backup deletion
func TestMysqlBackupAttachmentDelete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mysqlClient := &mysql.Client{}
	patches.ApplyMethodReturn(newMockMysqlBackupMeta().client, "UseMysqlClient", mysqlClient)

	patches.ApplyMethodFunc(mysqlClient, "DeleteBackup", func(request *mysql.DeleteBackupRequest) (*mysql.DeleteBackupResponse, error) {
		assert.Equal(t, "cdb-test123", *request.InstanceId)
		assert.Equal(t, int64(12345), *request.BackupId)

		resp := mysql.NewDeleteBackupResponse()
		resp.Response = &mysql.DeleteBackupResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMysqlBackupMeta()
	res := cdb.ResourceTencentCloudMysqlBackupAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "cdb-test123",
		"backup_method": "physical",
	})
	d.SetId("12345#cdb-test123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestMysqlBackupAttachmentDelete_APIError tests API error during deletion
func TestMysqlBackupAttachmentDelete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mysqlClient := &mysql.Client{}
	patches.ApplyMethodReturn(newMockMysqlBackupMeta().client, "UseMysqlClient", mysqlClient)

	patches.ApplyMethodFunc(mysqlClient, "DeleteBackup", func(request *mysql.DeleteBackupRequest) (*mysql.DeleteBackupResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Backup not found")
	})

	meta := newMockMysqlBackupMeta()
	res := cdb.ResourceTencentCloudMysqlBackupAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "cdb-test123",
		"backup_method": "physical",
	})
	d.SetId("12345#cdb-test123")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestMysqlBackupAttachment_Schema validates schema definition
func TestMysqlBackupAttachment_Schema(t *testing.T) {
	res := cdb.ResourceTencentCloudMysqlBackupAttachment()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "instance_id")
	assert.Contains(t, res.Schema, "backup_method")
	assert.Contains(t, res.Schema, "backup_db_table_list")
	assert.Contains(t, res.Schema, "manual_backup_name")
	assert.Contains(t, res.Schema, "encryption_flag")
	assert.Contains(t, res.Schema, "backup_id")

	// Check instance_id
	instanceId := res.Schema["instance_id"]
	assert.Equal(t, schema.TypeString, instanceId.Type)
	assert.True(t, instanceId.Required)
	assert.True(t, instanceId.ForceNew)

	// Check backup_method
	backupMethod := res.Schema["backup_method"]
	assert.Equal(t, schema.TypeString, backupMethod.Type)
	assert.True(t, backupMethod.Required)
	assert.True(t, backupMethod.ForceNew)

	// Check backup_id
	backupId := res.Schema["backup_id"]
	assert.Equal(t, schema.TypeInt, backupId.Type)
	assert.True(t, backupId.Computed)

	// Check backup_db_table_list
	dbTableList := res.Schema["backup_db_table_list"]
	assert.Equal(t, schema.TypeList, dbTableList.Type)
	assert.True(t, dbTableList.Optional)
	assert.True(t, dbTableList.ForceNew)
}

// TestMysqlBackupAttachment_IdFormat tests that ID uses correct format
func TestMysqlBackupAttachment_IdFormat(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mysqlClient := &mysql.Client{}
	patches.ApplyMethodReturn(newMockMysqlBackupMeta().client, "UseMysqlClient", mysqlClient)

	patches.ApplyMethodFunc(mysqlClient, "CreateBackup", func(request *mysql.CreateBackupRequest) (*mysql.CreateBackupResponse, error) {
		resp := mysql.NewCreateBackupResponse()
		resp.Response = &mysql.CreateBackupResponseParams{
			BackupId:  ptrUint64(42),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(mysqlClient, "DescribeBackups", func(request *mysql.DescribeBackupsRequest) (*mysql.DescribeBackupsResponse, error) {
		resp := mysql.NewDescribeBackupsResponse()
		resp.Response = &mysql.DescribeBackupsResponseParams{
			TotalCount: ptrInt64(1),
			Items: []*mysql.BackupInfo{
				{
					BackupId:   ptrInt64(42),
					Type:       ptrString("snapshot"),
					InstanceId: ptrString("cdb-snap"),
					Way:        ptrString("manual"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMysqlBackupMeta()
	res := cdb.ResourceTencentCloudMysqlBackupAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "cdb-snap",
		"backup_method": "snapshot",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "42#cdb-snap", d.Id())
	assert.Contains(t, d.Id(), "#")
}

func ptrString(s string) *string {
	return &s
}
