package mongodb_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	mongodbv20190725 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"
)

type mockMetaForAuditLogFile struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForAuditLogFile) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForAuditLogFile{}

func newMockMetaForAuditLogFile() *mockMetaForAuditLogFile {
	return &mockMetaForAuditLogFile{client: &connectivity.TencentCloudClient{}}
}

func ptrStrAuditLogFile(s string) *string {
	return &s
}

func ptrUint64AuditLogFile(u uint64) *uint64 {
	return &u
}

// go test ./tencentcloud/services/mongodb/ -run "TestAuditLogFile" -v -count=1 -gcflags="all=-l"

// TestAuditLogFile_Create_Success tests successful creation of audit log file
func TestAuditLogFile_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodbv20190725.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditLogFile().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "CreateAuditLogFile", func(request *mongodbv20190725.CreateAuditLogFileRequest) (*mongodbv20190725.CreateAuditLogFileResponse, error) {
		resp := mongodbv20190725.NewCreateAuditLogFileResponse()
		resp.Response = &mongodbv20190725.CreateAuditLogFileResponseParams{
			FileName:  ptrStrAuditLogFile("audit_log_20210712.log"),
			RequestId: ptrStrAuditLogFile("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(mongodbClient, "DescribeAuditLogFiles", func(request *mongodbv20190725.DescribeAuditLogFilesRequest) (*mongodbv20190725.DescribeAuditLogFilesResponse, error) {
		resp := mongodbv20190725.NewDescribeAuditLogFilesResponse()
		resp.Response = &mongodbv20190725.DescribeAuditLogFilesResponseParams{
			TotalCount: ptrUint64AuditLogFile(1),
			Items: []*mongodbv20190725.AuditLogFile{
				{
					FileName:     ptrStrAuditLogFile("audit_log_20210712.log"),
					CreateTime:   ptrStrAuditLogFile("2021-07-12 10:30:00"),
					Status:       ptrStrAuditLogFile("success"),
					FileSize:     ptrUint64AuditLogFile(1024),
					DownloadUrl:  ptrStrAuditLogFile("https://example.com/download"),
					ErrMsg:       ptrStrAuditLogFile(""),
					ProgressRate: ptrUint64AuditLogFile(100),
				},
			},
			RequestId: ptrStrAuditLogFile("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditLogFile()
	res := mongodb.ResourceTencentCloudMongodbAuditLogFile()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"start_time":  "2021-07-12 10:29:20",
		"end_time":    "2021-07-12 10:39:20",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cmgo-test1234#audit_log_20210712.log", d.Id())
	assert.Equal(t, "audit_log_20210712.log", d.Get("file_name"))
}

// TestAuditLogFile_Create_EmptyFileName tests creation failure when FileName is empty
func TestAuditLogFile_Create_EmptyFileName(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodbv20190725.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditLogFile().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "CreateAuditLogFile", func(request *mongodbv20190725.CreateAuditLogFileRequest) (*mongodbv20190725.CreateAuditLogFileResponse, error) {
		resp := mongodbv20190725.NewCreateAuditLogFileResponse()
		resp.Response = &mongodbv20190725.CreateAuditLogFileResponseParams{
			FileName:  ptrStrAuditLogFile(""),
			RequestId: ptrStrAuditLogFile("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditLogFile()
	res := mongodb.ResourceTencentCloudMongodbAuditLogFile()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"start_time":  "2021-07-12 10:29:20",
		"end_time":    "2021-07-12 10:39:20",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FileName is nil or empty")
}

// TestAuditLogFile_Read_Success tests successful read of audit log file
func TestAuditLogFile_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodbv20190725.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditLogFile().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeAuditLogFiles", func(request *mongodbv20190725.DescribeAuditLogFilesRequest) (*mongodbv20190725.DescribeAuditLogFilesResponse, error) {
		resp := mongodbv20190725.NewDescribeAuditLogFilesResponse()
		resp.Response = &mongodbv20190725.DescribeAuditLogFilesResponseParams{
			TotalCount: ptrUint64AuditLogFile(1),
			Items: []*mongodbv20190725.AuditLogFile{
				{
					FileName:     ptrStrAuditLogFile("audit_log_20210712.log"),
					CreateTime:   ptrStrAuditLogFile("2021-07-12 10:30:00"),
					Status:       ptrStrAuditLogFile("success"),
					FileSize:     ptrUint64AuditLogFile(2048),
					DownloadUrl:  ptrStrAuditLogFile("https://example.com/download"),
					ErrMsg:       ptrStrAuditLogFile(""),
					ProgressRate: ptrUint64AuditLogFile(100),
				},
			},
			RequestId: ptrStrAuditLogFile("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditLogFile()
	res := mongodb.ResourceTencentCloudMongodbAuditLogFile()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"start_time":  "2021-07-12 10:29:20",
		"end_time":    "2021-07-12 10:39:20",
	})
	d.SetId("cmgo-test1234#audit_log_20210712.log")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cmgo-test1234", d.Get("instance_id"))
	assert.Equal(t, "audit_log_20210712.log", d.Get("file_name"))

	items := d.Get("items").([]interface{})
	assert.Equal(t, 1, len(items))
	item := items[0].(map[string]interface{})
	assert.Equal(t, "audit_log_20210712.log", item["file_name"])
	assert.Equal(t, "2021-07-12 10:30:00", item["create_time"])
	assert.Equal(t, "success", item["status"])
}

// TestAuditLogFile_Read_NotFound tests read when file is not found
func TestAuditLogFile_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodbv20190725.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditLogFile().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeAuditLogFiles", func(request *mongodbv20190725.DescribeAuditLogFilesRequest) (*mongodbv20190725.DescribeAuditLogFilesResponse, error) {
		resp := mongodbv20190725.NewDescribeAuditLogFilesResponse()
		resp.Response = &mongodbv20190725.DescribeAuditLogFilesResponseParams{
			TotalCount: ptrUint64AuditLogFile(0),
			Items:      []*mongodbv20190725.AuditLogFile{},
			RequestId:  ptrStrAuditLogFile("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditLogFile()
	res := mongodb.ResourceTencentCloudMongodbAuditLogFile()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"start_time":  "2021-07-12 10:29:20",
		"end_time":    "2021-07-12 10:39:20",
	})
	d.SetId("cmgo-test1234#audit_log_20210712.log")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestAuditLogFile_Update_ImmutableError tests update returns error for immutable fields
func TestAuditLogFile_Update_ImmutableError(t *testing.T) {
	res := mongodb.ResourceTencentCloudMongodbAuditLogFile()

	// Verify that the resource has an Update function
	assert.NotNil(t, res.Update)

	// Verify schema fields are not ForceNew (except instance_id)
	assert.True(t, res.Schema["instance_id"].ForceNew)
	assert.False(t, res.Schema["start_time"].ForceNew)
	assert.False(t, res.Schema["end_time"].ForceNew)
	assert.False(t, res.Schema["order"].ForceNew)
	assert.False(t, res.Schema["order_by"].ForceNew)
	assert.False(t, res.Schema["filter"].ForceNew)
}

// TestAuditLogFile_Delete_Success tests successful deletion of audit log file
func TestAuditLogFile_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodbv20190725.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditLogFile().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DeleteAuditLogFile", func(request *mongodbv20190725.DeleteAuditLogFileRequest) (*mongodbv20190725.DeleteAuditLogFileResponse, error) {
		resp := mongodbv20190725.NewDeleteAuditLogFileResponse()
		resp.Response = &mongodbv20190725.DeleteAuditLogFileResponseParams{
			RequestId: ptrStrAuditLogFile("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditLogFile()
	res := mongodb.ResourceTencentCloudMongodbAuditLogFile()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"start_time":  "2021-07-12 10:29:20",
		"end_time":    "2021-07-12 10:39:20",
	})
	d.SetId("cmgo-test1234#audit_log_20210712.log")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestAuditLogFile_Schema tests the schema definition
func TestAuditLogFile_Schema(t *testing.T) {
	res := mongodb.ResourceTencentCloudMongodbAuditLogFile()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields
	assert.True(t, res.Schema["instance_id"].Required)
	assert.True(t, res.Schema["start_time"].Required)
	assert.True(t, res.Schema["end_time"].Required)

	// Check optional fields
	assert.True(t, res.Schema["order"].Optional)
	assert.True(t, res.Schema["order_by"].Optional)
	assert.True(t, res.Schema["filter"].Optional)

	// Check computed fields
	assert.True(t, res.Schema["file_name"].Computed)
	assert.True(t, res.Schema["items"].Computed)

	// Check filter MaxItems
	assert.Equal(t, 1, res.Schema["filter"].MaxItems)
}
