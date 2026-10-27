package postgresql_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
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

func ptrStrAuditLog(s string) *string {
	return &s
}

func ptrUint64AuditLog(v uint64) *uint64 {
	return &v
}

// go test ./tencentcloud/services/postgresql/ -run "TestPostgresAuditLogFile" -v -count=1 -gcflags="all=-l"

// TestPostgresAuditLogFile_Create_Success tests Create function
func TestPostgresAuditLogFile_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditLogFile().client, "UsePostgresqlClient", pgClient)

	callCount := 0
	patches.ApplyMethodFunc(pgClient, "DescribeAuditLogFiles", func(request *postgresql.DescribeAuditLogFilesRequest) (*postgresql.DescribeAuditLogFilesResponse, error) {
		callCount++
		resp := postgresql.NewDescribeAuditLogFilesResponse()
		if callCount == 1 {
			resp.Response = &postgresql.DescribeAuditLogFilesResponseParams{
				TotalCount: ptrUint64AuditLog(0),
				Items:      []*postgresql.AuditLogFile{},
				RequestId:  ptrStrAuditLog("fake-request-id"),
			}
		} else {
			resp.Response = &postgresql.DescribeAuditLogFilesResponseParams{
				TotalCount: ptrUint64AuditLog(1),
				Items: []*postgresql.AuditLogFile{
					{
						FileName:    ptrStrAuditLog("audit_20260325.csv"),
						Status:      ptrStrAuditLog("success"),
						FileSize:    ptrUint64AuditLog(10),
						CreateTime:  ptrStrAuditLog("2026-03-25 00:05:00"),
						DownloadUrl: ptrStrAuditLog("https://example.com/download"),
						Progress:    ptrUint64AuditLog(100),
						FinishTime:  ptrStrAuditLog("2026-03-25 00:06:00"),
					},
				},
				RequestId: ptrStrAuditLog("fake-request-id"),
			}
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(pgClient, "CreateAuditLogFile", func(request *postgresql.CreateAuditLogFileRequest) (*postgresql.CreateAuditLogFileResponse, error) {
		assert.Equal(t, "postgres-12345", *request.InstanceId)
		assert.Equal(t, "2026-03-25 00:00:00", *request.StartTime)
		assert.Equal(t, "2026-03-25 01:00:00", *request.EndTime)
		assert.Equal(t, "postgres", *request.Product)
		resp := postgresql.NewCreateAuditLogFileResponse()
		resp.Response = &postgresql.CreateAuditLogFileResponseParams{
			RequestId: ptrStrAuditLog("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditLogFile()
	res := svcpostgresql.ResourceTencentCloudPostgresAuditLogFile()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "postgres-12345",
		"start_time":  "2026-03-25 00:00:00",
		"end_time":    "2026-03-25 01:00:00",
		"product":     "postgres",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Contains(t, d.Id(), "postgres-12345")
	assert.Contains(t, d.Id(), "audit_20260325.csv")
}

// TestPostgresAuditLogFile_Read_Success tests Read function
func TestPostgresAuditLogFile_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditLogFile().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeAuditLogFiles", func(request *postgresql.DescribeAuditLogFilesRequest) (*postgresql.DescribeAuditLogFilesResponse, error) {
		assert.Equal(t, "postgres-12345", *request.InstanceId)
		assert.Equal(t, "audit_20260325.csv", *request.FileName)
		resp := postgresql.NewDescribeAuditLogFilesResponse()
		resp.Response = &postgresql.DescribeAuditLogFilesResponseParams{
			TotalCount: ptrUint64AuditLog(1),
			Items: []*postgresql.AuditLogFile{
				{
					FileName:    ptrStrAuditLog("audit_20260325.csv"),
					Status:      ptrStrAuditLog("success"),
					FileSize:    ptrUint64AuditLog(10),
					CreateTime:  ptrStrAuditLog("2026-03-25 00:05:00"),
					DownloadUrl: ptrStrAuditLog("https://example.com/download"),
					ErrMsg:      ptrStrAuditLog(""),
					Progress:    ptrUint64AuditLog(100),
					FinishTime:  ptrStrAuditLog("2026-03-25 00:06:00"),
				},
			},
			RequestId: ptrStrAuditLog("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditLogFile()
	res := svcpostgresql.ResourceTencentCloudPostgresAuditLogFile()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "postgres-12345",
		"start_time":  "2026-03-25 00:00:00",
		"end_time":    "2026-03-25 01:00:00",
		"product":     "postgres",
	})
	d.SetId("postgres-12345#audit_20260325.csv")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "postgres-12345", d.Get("instance_id"))
	assert.Equal(t, "audit_20260325.csv", d.Get("file_name"))
	assert.Equal(t, "success", d.Get("status"))
	assert.Equal(t, 10, d.Get("file_size"))
	assert.Equal(t, "2026-03-25 00:05:00", d.Get("create_time"))
	assert.Equal(t, "https://example.com/download", d.Get("download_url"))
	assert.Equal(t, 100, d.Get("progress"))
	assert.Equal(t, "2026-03-25 00:06:00", d.Get("finish_time"))
}

// TestPostgresAuditLogFile_Read_NotFound tests Read when file not found
func TestPostgresAuditLogFile_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditLogFile().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DescribeAuditLogFiles", func(request *postgresql.DescribeAuditLogFilesRequest) (*postgresql.DescribeAuditLogFilesResponse, error) {
		resp := postgresql.NewDescribeAuditLogFilesResponse()
		resp.Response = &postgresql.DescribeAuditLogFilesResponseParams{
			TotalCount: ptrUint64AuditLog(0),
			Items:      []*postgresql.AuditLogFile{},
			RequestId:  ptrStrAuditLog("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditLogFile()
	res := svcpostgresql.ResourceTencentCloudPostgresAuditLogFile()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "postgres-12345",
		"start_time":  "2026-03-25 00:00:00",
		"end_time":    "2026-03-25 01:00:00",
		"product":     "postgres",
	})
	d.SetId("postgres-12345#audit_20260325.csv")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestPostgresAuditLogFile_Delete_Success tests Delete function
func TestPostgresAuditLogFile_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditLogFile().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DeleteAuditLogFile", func(request *postgresql.DeleteAuditLogFileRequest) (*postgresql.DeleteAuditLogFileResponse, error) {
		assert.Equal(t, "postgres-12345", *request.InstanceId)
		assert.Equal(t, "postgres", *request.Product)
		assert.Equal(t, "audit_20260325.csv", *request.FileName)
		resp := postgresql.NewDeleteAuditLogFileResponse()
		resp.Response = &postgresql.DeleteAuditLogFileResponseParams{
			RequestId: ptrStrAuditLog("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForAuditLogFile()
	res := svcpostgresql.ResourceTencentCloudPostgresAuditLogFile()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "postgres-12345",
		"start_time":  "2026-03-25 00:00:00",
		"end_time":    "2026-03-25 01:00:00",
		"product":     "postgres",
	})
	d.SetId("postgres-12345#audit_20260325.csv")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestPostgresAuditLogFile_Delete_Error tests Delete with API error
func TestPostgresAuditLogFile_Delete_Error(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	pgClient := &postgresql.Client{}
	patches.ApplyMethodReturn(newMockMetaForAuditLogFile().client, "UsePostgresqlClient", pgClient)

	patches.ApplyMethodFunc(pgClient, "DeleteAuditLogFile", func(request *postgresql.DeleteAuditLogFileRequest) (*postgresql.DeleteAuditLogFileResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=Internal error")
	})

	meta := newMockMetaForAuditLogFile()
	res := svcpostgresql.ResourceTencentCloudPostgresAuditLogFile()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "postgres-12345",
		"start_time":  "2026-03-25 00:00:00",
		"end_time":    "2026-03-25 01:00:00",
		"product":     "postgres",
	})
	d.SetId("postgres-12345#audit_20260325.csv")

	err := res.Delete(d, meta)
	assert.Error(t, err)
}
