package dlc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcdlc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"
)

type mockMetaDlcInternalTable struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDlcInternalTable) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDlcInternalTable{}

func newMockMetaDlcInternalTable() *mockMetaDlcInternalTable {
	return &mockMetaDlcInternalTable{client: &connectivity.TencentCloudClient{}}
}

func ptrStrDlcIt(s string) *string {
	return &s
}

func ptrBoolDlcIt(b bool) *bool {
	return &b
}

func ptrInt64DlcIt(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/dlc/ -run "TestDlcInternalTable" -v -count=1 -gcflags="all=-l"

// TestDlcInternalTable_CreateThenRead tests create followed by read
func TestDlcInternalTable_CreateThenRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcInternalTable().client, "UseDlcClient", dlcClient)

	sqlText := "CREATE TABLE tf_example_db.tf_example_table (id bigint, name string)"
	patches.ApplyMethodFunc(dlcClient, "GenerateInternalTableWithContext", func(_ context.Context, request *dlc.GenerateInternalTableRequest) (*dlc.GenerateInternalTableResponse, error) {
		assert.NotNil(t, request.TableBaseInfo)
		assert.Equal(t, "tf_example_db", *request.TableBaseInfo.DatabaseName)
		assert.Equal(t, "tf_example_table", *request.TableBaseInfo.TableName)
		assert.Len(t, request.Columns, 2)
		assert.Equal(t, "id", *request.Columns[0].Name)
		assert.Equal(t, "bigint", *request.Columns[0].Type)

		resp := dlc.NewGenerateInternalTableResponse()
		resp.Response = &dlc.GenerateInternalTableResponseParams{
			Execution: &dlc.Execution{
				SQL: ptrStrDlcIt(sqlText),
			},
			IsTIcebergSql: ptrBoolDlcIt(false),
			RequestId:     ptrStrDlcIt("fake-request-id-create"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dlcClient, "DescribeTableWithContext", func(_ context.Context, request *dlc.DescribeTableRequest) (*dlc.DescribeTableResponse, error) {
		assert.Equal(t, "tf_example_db", *request.DatabaseName)
		assert.Equal(t, "tf_example_table", *request.TableName)

		resp := dlc.NewDescribeTableResponse()
		resp.Response = &dlc.DescribeTableResponseParams{
			Table: &dlc.TableResponseInfo{
				TableBaseInfo: &dlc.TableBaseInfo{
					DatabaseName: ptrStrDlcIt("tf_example_db"),
					TableName:    ptrStrDlcIt("tf_example_table"),
					TableComment: ptrStrDlcIt("test table"),
					Type:         ptrStrDlcIt("TABLE"),
					TableFormat:  ptrStrDlcIt("hive"),
				},
				Columns: []*dlc.Column{
					{
						Name:    ptrStrDlcIt("id"),
						Type:    ptrStrDlcIt("bigint"),
						Comment: ptrStrDlcIt("id column"),
					},
					{
						Name:    ptrStrDlcIt("name"),
						Type:    ptrStrDlcIt("string"),
						Comment: ptrStrDlcIt("name column"),
					},
				},
				Location:         ptrStrDlcIt("cosn://tf-bucket/path"),
				CreateTime:       ptrStrDlcIt("1700000000000"),
				ModifiedTime:     ptrStrDlcIt("1700000001000"),
				InputFormat:      ptrStrDlcIt("org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat"),
				StorageSize:      ptrInt64DlcIt(1024),
				RecordCount:      ptrInt64DlcIt(100),
				InputFormatShort: ptrStrDlcIt("Parquet"),
			},
			RequestId: ptrStrDlcIt("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcInternalTable()
	res := svcdlc.ResourceTencentCloudDlcInternalTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
		"table_name":    "tf_example_table",
		"table_comment": "test table",
		"type":          "TABLE",
		"table_format":  "hive",
		"columns": []interface{}{
			map[string]interface{}{
				"name":    "id",
				"type":    "bigint",
				"comment": "id column",
			},
			map[string]interface{}{
				"name":    "name",
				"type":    "string",
				"comment": "name column",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tf_example_db#tf_example_table", d.Id())
	assert.Equal(t, sqlText, d.Get("execution").(string))
	assert.Equal(t, false, d.Get("is_t_iceberg_sql").(bool))
	assert.Equal(t, "cosn://tf-bucket/path", d.Get("location").(string))
	assert.Equal(t, "1700000000000", d.Get("create_time").(string))
	assert.Equal(t, "1700000001000", d.Get("modified_time").(string))
	assert.Equal(t, 1024, d.Get("storage_size").(int))
	assert.Equal(t, 100, d.Get("record_count").(int))
}

// TestDlcInternalTable_Update tests update of mutable fields
func TestDlcInternalTable_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcInternalTable().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.AlterTableCommentRequest
	patches.ApplyMethodFunc(dlcClient, "AlterTableCommentWithContext", func(_ context.Context, request *dlc.AlterTableCommentRequest) (*dlc.AlterTableCommentResponse, error) {
		capturedRequest = request
		resp := dlc.NewAlterTableCommentResponse()
		resp.Response = &dlc.AlterTableCommentResponseParams{
			RequestId: ptrStrDlcIt("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dlcClient, "DescribeTableWithContext", func(_ context.Context, request *dlc.DescribeTableRequest) (*dlc.DescribeTableResponse, error) {
		resp := dlc.NewDescribeTableResponse()
		resp.Response = &dlc.DescribeTableResponseParams{
			Table: &dlc.TableResponseInfo{
				TableBaseInfo: &dlc.TableBaseInfo{
					DatabaseName: ptrStrDlcIt("tf_example_db"),
					TableName:    ptrStrDlcIt("tf_example_table"),
					TableComment: ptrStrDlcIt("updated comment"),
					Type:         ptrStrDlcIt("TABLE"),
					TableFormat:  ptrStrDlcIt("hive"),
				},
				Columns: []*dlc.Column{
					{
						Name: ptrStrDlcIt("id"),
						Type: ptrStrDlcIt("bigint"),
					},
				},
			},
			RequestId: ptrStrDlcIt("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcInternalTable()
	res := svcdlc.ResourceTencentCloudDlcInternalTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
		"table_name":    "tf_example_table",
		"table_comment": "updated comment",
		"type":          "TABLE",
		"table_format":  "hive",
		"columns": []interface{}{
			map[string]interface{}{
				"name": "id",
				"type": "bigint",
			},
		},
	})
	d.SetId("tf_example_db#tf_example_table")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.TableBaseInfo)
	assert.Equal(t, "updated comment", *capturedRequest.TableBaseInfo.TableComment)
	assert.Equal(t, "tf_example_db", d.Id())
	assert.Equal(t, "updated comment", d.Get("table_comment").(string))
}

// TestDlcInternalTable_Delete tests delete
func TestDlcInternalTable_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcInternalTable().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.DeleteTableRequest
	patches.ApplyMethodFunc(dlcClient, "DeleteTableWithContext", func(_ context.Context, request *dlc.DeleteTableRequest) (*dlc.DeleteTableResponse, error) {
		capturedRequest = request
		resp := dlc.NewDeleteTableResponse()
		resp.Response = &dlc.DeleteTableResponseParams{
			RequestId: ptrStrDlcIt("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcInternalTable()
	res := svcdlc.ResourceTencentCloudDlcInternalTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
		"table_name":    "tf_example_table",
		"columns": []interface{}{
			map[string]interface{}{
				"name": "id",
				"type": "bigint",
			},
		},
	})
	d.SetId("tf_example_db#tf_example_table")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.TableBaseInfo)
	assert.Equal(t, "tf_example_db", *capturedRequest.TableBaseInfo.DatabaseName)
	assert.Equal(t, "tf_example_table", *capturedRequest.TableBaseInfo.TableName)
}

// TestDlcInternalTable_EmptyCreateResult tests that empty create result returns error
func TestDlcInternalTable_EmptyCreateResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcInternalTable().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "GenerateInternalTableWithContext", func(_ context.Context, request *dlc.GenerateInternalTableRequest) (*dlc.GenerateInternalTableResponse, error) {
		resp := dlc.NewGenerateInternalTableResponse()
		resp.Response = &dlc.GenerateInternalTableResponseParams{
			RequestId: ptrStrDlcIt("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcInternalTable()
	res := svcdlc.ResourceTencentCloudDlcInternalTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
		"table_name":    "tf_example_table",
		"columns": []interface{}{
			map[string]interface{}{
				"name": "id",
				"type": "bigint",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Execution is nil")
}

// TestDlcInternalTable_TableNotFound tests that read clears id when table not found
func TestDlcInternalTable_TableNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcInternalTable().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeTableWithContext", func(_ context.Context, request *dlc.DescribeTableRequest) (*dlc.DescribeTableResponse, error) {
		resp := dlc.NewDescribeTableResponse()
		resp.Response = &dlc.DescribeTableResponseParams{
			RequestId: ptrStrDlcIt("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcInternalTable()
	res := svcdlc.ResourceTencentCloudDlcInternalTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
		"table_name":    "tf_example_table",
		"columns": []interface{}{
			map[string]interface{}{
				"name": "id",
				"type": "bigint",
			},
		},
	})
	d.SetId("tf_example_db#tf_example_table")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Empty(t, d.Id())
}

// TestDlcInternalTable_Schema validates schema definition
func TestDlcInternalTable_Schema(t *testing.T) {
	res := svcdlc.ResourceTencentCloudDlcInternalTable()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "database_name")
	assert.Contains(t, res.Schema, "table_name")
	assert.Contains(t, res.Schema, "datasource_connection_name")
	assert.Contains(t, res.Schema, "table_comment")
	assert.Contains(t, res.Schema, "type")
	assert.Contains(t, res.Schema, "table_format")
	assert.Contains(t, res.Schema, "user_alias")
	assert.Contains(t, res.Schema, "user_sub_uin")
	assert.Contains(t, res.Schema, "db_govern_policy_is_disable")
	assert.Contains(t, res.Schema, "govern_policy")
	assert.Contains(t, res.Schema, "smart_policy")
	assert.Contains(t, res.Schema, "primary_keys")
	assert.Contains(t, res.Schema, "columns")
	assert.Contains(t, res.Schema, "partitions")
	assert.Contains(t, res.Schema, "properties")
	assert.Contains(t, res.Schema, "upsert_keys")

	assert.Contains(t, res.Schema, "location")
	assert.Contains(t, res.Schema, "modified_time")
	assert.Contains(t, res.Schema, "create_time")
	assert.Contains(t, res.Schema, "input_format")
	assert.Contains(t, res.Schema, "storage_size")
	assert.Contains(t, res.Schema, "record_count")
	assert.Contains(t, res.Schema, "map_materialized_view_name")
	assert.Contains(t, res.Schema, "heat_value")
	assert.Contains(t, res.Schema, "input_format_short")
	assert.Contains(t, res.Schema, "execution")
	assert.Contains(t, res.Schema, "is_t_iceberg_sql")

	databaseName := res.Schema["database_name"]
	assert.True(t, databaseName.Required)
	assert.True(t, databaseName.ForceNew)

	tableName := res.Schema["table_name"]
	assert.True(t, tableName.Required)
	assert.True(t, tableName.ForceNew)

	columns := res.Schema["columns"]
	assert.True(t, columns.Required)
	assert.True(t, columns.ForceNew)

	properties := res.Schema["properties"]
	assert.True(t, properties.Optional)
	assert.True(t, properties.ForceNew)

	primaryKeys := res.Schema["primary_keys"]
	assert.True(t, primaryKeys.Optional)
	assert.True(t, primaryKeys.ForceNew)

	upsertKeys := res.Schema["upsert_keys"]
	assert.True(t, upsertKeys.Optional)
	assert.True(t, upsertKeys.ForceNew)

	tableComment := res.Schema["table_comment"]
	assert.True(t, tableComment.Optional)
	assert.False(t, tableComment.ForceNew)

	location := res.Schema["location"]
	assert.True(t, location.Computed)

	storageSize := res.Schema["storage_size"]
	assert.True(t, storageSize.Computed)
}

// TestDlcInternalTable_CreateWithSmartPolicy tests create with smart policy
func TestDlcInternalTable_CreateWithSmartPolicy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcInternalTable().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.GenerateInternalTableRequest
	patches.ApplyMethodFunc(dlcClient, "GenerateInternalTableWithContext", func(_ context.Context, request *dlc.GenerateInternalTableRequest) (*dlc.GenerateInternalTableResponse, error) {
		capturedRequest = request
		resp := dlc.NewGenerateInternalTableResponse()
		resp.Response = &dlc.GenerateInternalTableResponseParams{
			Execution: &dlc.Execution{
				SQL: ptrStrDlcIt("CREATE TABLE"),
			},
			IsTIcebergSql: ptrBoolDlcIt(true),
			RequestId:     ptrStrDlcIt("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dlcClient, "DescribeTableWithContext", func(_ context.Context, request *dlc.DescribeTableRequest) (*dlc.DescribeTableResponse, error) {
		resp := dlc.NewDescribeTableResponse()
		resp.Response = &dlc.DescribeTableResponseParams{
			Table: &dlc.TableResponseInfo{
				TableBaseInfo: &dlc.TableBaseInfo{
					DatabaseName: ptrStrDlcIt("tf_example_db"),
					TableName:    ptrStrDlcIt("tf_example_table"),
				},
				Columns: []*dlc.Column{
					{
						Name: ptrStrDlcIt("id"),
						Type: ptrStrDlcIt("bigint"),
					},
				},
			},
			RequestId: ptrStrDlcIt("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcInternalTable()
	res := svcdlc.ResourceTencentCloudDlcInternalTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
		"table_name":    "tf_example_table",
		"columns": []interface{}{
			map[string]interface{}{
				"name": "id",
				"type": "bigint",
			},
		},
		"smart_policy": []interface{}{
			map[string]interface{}{
				"base_info": []interface{}{
					map[string]interface{}{
						"uin":         "100012345678",
						"policy_type": "Table",
						"catalog":     "DataLakeCatalog",
						"database":    "tf_example_db",
						"table":       "tf_example_table",
						"app_id":      "1300123456",
					},
				},
				"policy": []interface{}{
					map[string]interface{}{
						"inherit": "false",
						"table_expiration": []interface{}{
							map[string]interface{}{
								"enabled":    true,
								"expiration": 90,
							},
						},
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest.TableBaseInfo.SmartPolicy)
	assert.NotNil(t, capturedRequest.TableBaseInfo.SmartPolicy.BaseInfo)
	assert.Equal(t, "100012345678", *capturedRequest.TableBaseInfo.SmartPolicy.BaseInfo.Uin)
	assert.NotNil(t, capturedRequest.TableBaseInfo.SmartPolicy.Policy)
	assert.NotNil(t, capturedRequest.TableBaseInfo.SmartPolicy.Policy.TableExpiration)
	assert.Equal(t, true, *capturedRequest.TableBaseInfo.SmartPolicy.Policy.TableExpiration.Enabled)
	assert.Equal(t, uint64(90), *capturedRequest.TableBaseInfo.SmartPolicy.Policy.TableExpiration.Expiration)
}

// TestDlcInternalTable_APIError tests Create when API returns error
func TestDlcInternalTable_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcInternalTable().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "GenerateInternalTableWithContext", func(_ context.Context, request *dlc.GenerateInternalTableRequest) (*dlc.GenerateInternalTableResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Database not found")
	})

	meta := newMockMetaDlcInternalTable()
	res := svcdlc.ResourceTencentCloudDlcInternalTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "invalid_db",
		"table_name":    "tf_example_table",
		"columns": []interface{}{
			map[string]interface{}{
				"name": "id",
				"type": "bigint",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}
