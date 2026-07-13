package dlc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcdlc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"
)

// mockMetaDlcDmsTable implements tccommon.ProviderMeta
type mockMetaDlcDmsTable struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDlcDmsTable) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDlcDmsTable{}

func newMockMetaDlcDmsTable() *mockMetaDlcDmsTable {
	return &mockMetaDlcDmsTable{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDmsTable(s string) *string {
	return &s
}

func ptrInt64DmsTable(i int64) *int64 {
	return &i
}

func ptrBoolDmsTable(b bool) *bool {
	return &b
}

// go test ./tencentcloud/services/dlc/ -run "TestResourceTencentCloudDlcDmsTable_Create" -v -count=1 -gcflags="all=-l"

// TestResourceTencentCloudDlcDmsTable_Create tests the Create flow: CreateDMSTable succeeds, then Read is called to refresh state
func TestResourceTencentCloudDlcDmsTable_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsTable().client, "UseDlcClient", dlcClient)

	var capturedCreateRequest *dlc.CreateDMSTableRequest
	patches.ApplyMethodFunc(dlcClient, "CreateDMSTableWithContext", func(_ context.Context, request *dlc.CreateDMSTableRequest) (*dlc.CreateDMSTableResponse, error) {
		capturedCreateRequest = request
		resp := dlc.NewCreateDMSTableResponse()
		resp.Response = &dlc.CreateDMSTableResponseParams{
			RequestId: ptrStringDmsTable("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeDMSTableWithContext for the Read call after Create
	patches.ApplyMethodFunc(dlcClient, "DescribeDMSTableWithContext", func(_ context.Context, request *dlc.DescribeDMSTableRequest) (*dlc.DescribeDMSTableResponse, error) {
		resp := dlc.NewDescribeDMSTableResponse()
		resp.Response = &dlc.DescribeDMSTableResponseParams{
			DbName: ptrStringDmsTable("tf_example_db"),
			Name:   ptrStringDmsTable("tf_example_table"),
			Type:   ptrStringDmsTable("EXTERNAL_TABLE"),
			Asset: &dlc.Asset{
				Name:        ptrStringDmsTable("tf_example_table"),
				Description: ptrStringDmsTable("tf example dlc dms table"),
				Owner:       ptrStringDmsTable("root"),
			},
			SchemaName: ptrStringDmsTable("default"),
			Retention:  ptrInt64DmsTable(0),
			Columns: []*dlc.DMSColumn{
				{
					Name:     ptrStringDmsTable("id"),
					Type:     ptrStringDmsTable("bigint"),
					Position: ptrInt64DmsTable(1),
				},
				{
					Name:     ptrStringDmsTable("name"),
					Type:     ptrStringDmsTable("string"),
					Position: ptrInt64DmsTable(2),
				},
			},
			Sds: &dlc.DMSSds{
				Location:     ptrStringDmsTable("cosn://tf-example-bucket/example/"),
				InputFormat:  ptrStringDmsTable("org.apache.hadoop.hive.ql.io.avro.AvroContainerInputFormat"),
				OutputFormat: ptrStringDmsTable("org.apache.hadoop.hive.ql.io.avro.AvroContainerOutputFormat"),
				SerdeLib:     ptrStringDmsTable("org.apache.hadoop.hive.serde2.avro.AvroSerDe"),
				SerdeParams:  []*dlc.KVPair{{Key: ptrStringDmsTable("serialization.format"), Value: ptrStringDmsTable("1")}},
			},
			RequestId: ptrStringDmsTable("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsTable()
	res := svcdlc.ResourceTencentCloudDlcDmsTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_name": "tf_example_db",
		"name":    "tf_example_table",
		"type":    "EXTERNAL_TABLE",
		"asset": []interface{}{
			map[string]interface{}{
				"name":        "tf_example_table",
				"description": "tf example dlc dms table",
				"owner":       "root",
			},
		},
		"columns": []interface{}{
			map[string]interface{}{"name": "id", "type": "bigint", "position": 1},
			map[string]interface{}{"name": "name", "type": "string", "position": 2},
		},
		"sds": []interface{}{
			map[string]interface{}{
				"location":      "cosn://tf-example-bucket/example/",
				"input_format":  "org.apache.hadoop.hive.ql.io.avro.AvroContainerInputFormat",
				"output_format": "org.apache.hadoop.hive.ql.io.avro.AvroContainerOutputFormat",
				"serde_lib":     "org.apache.hadoop.hive.serde2.avro.AvroSerDe",
				"serde_params": []interface{}{
					map[string]interface{}{"key": "serialization.format", "value": "1"},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tf_example_db#tf_example_table", d.Id())
	assert.NotNil(t, capturedCreateRequest)
	assert.NotNil(t, capturedCreateRequest.DbName)
	assert.Equal(t, "tf_example_db", *capturedCreateRequest.DbName)
	assert.NotNil(t, capturedCreateRequest.Name)
	assert.Equal(t, "tf_example_table", *capturedCreateRequest.Name)
}

// TestResourceTencentCloudDlcDmsTable_Read tests the Read flow when the table exists
func TestResourceTencentCloudDlcDmsTable_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsTable().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeDMSTableWithContext", func(_ context.Context, request *dlc.DescribeDMSTableRequest) (*dlc.DescribeDMSTableResponse, error) {
		resp := dlc.NewDescribeDMSTableResponse()
		resp.Response = &dlc.DescribeDMSTableResponseParams{
			DbName: ptrStringDmsTable("tf_example_db"),
			Name:   ptrStringDmsTable("tf_example_table"),
			Type:   ptrStringDmsTable("EXTERNAL_TABLE"),
			Asset: &dlc.Asset{
				Name:        ptrStringDmsTable("tf_example_table"),
				Description: ptrStringDmsTable("tf example dlc dms table"),
				Owner:       ptrStringDmsTable("root"),
			},
			SchemaName: ptrStringDmsTable("default"),
			Retention:  ptrInt64DmsTable(0),
			Columns: []*dlc.DMSColumn{
				{
					Name:     ptrStringDmsTable("id"),
					Type:     ptrStringDmsTable("bigint"),
					Position: ptrInt64DmsTable(1),
				},
			},
			RequestId: ptrStringDmsTable("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsTable()
	res := svcdlc.ResourceTencentCloudDlcDmsTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_name": "tf_example_db",
		"name":    "tf_example_table",
	})
	d.SetId("tf_example_db#tf_example_table")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tf_example_db#tf_example_table", d.Id())
	assert.Equal(t, "tf_example_db", d.Get("db_name").(string))
	assert.Equal(t, "tf_example_table", d.Get("name").(string))
	assert.Equal(t, "EXTERNAL_TABLE", d.Get("type").(string))
	assert.Equal(t, "default", d.Get("schema_name").(string))
}

// TestResourceTencentCloudDlcDmsTable_Read_NotFound tests the Read flow when the table does not exist (response is nil)
func TestResourceTencentCloudDlcDmsTable_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsTable().client, "UseDlcClient", dlcClient)

	// DescribeDmsTableById returns nil ret when retry fails with NonRetryableError (response is nil).
	// To simulate NotFound, we make DescribeDMSTableWithContext return a response with nil Response.
	patches.ApplyMethodFunc(dlcClient, "DescribeDMSTableWithContext", func(_ context.Context, request *dlc.DescribeDMSTableRequest) (*dlc.DescribeDMSTableResponse, error) {
		resp := dlc.NewDescribeDMSTableResponse()
		// Return an error so the retry in DescribeDmsTableById returns NonRetryableError and ret stays nil
		return resp, nil
	})

	// Since the mock returns a response with nil Response, DescribeDmsTableById will return NonRetryableError.
	// The Read method checks if respData == nil and calls d.SetId("").
	// But because the service returns NonRetryableError, the Read returns an error, not SetId("").
	// To properly test the d.SetId("") path, we need to mock DescribeDmsTableById directly.
	patches.ApplyMethodFunc(&svcdlc.DlcService{}, "DescribeDmsTableById", func(_ context.Context, dbName, name string) (*dlc.DescribeDMSTableResponseParams, error) {
		return nil, nil
	})

	meta := newMockMetaDlcDmsTable()
	res := svcdlc.ResourceTencentCloudDlcDmsTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_name": "tf_example_db",
		"name":    "tf_example_table",
	})
	d.SetId("tf_example_db#tf_example_table")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestResourceTencentCloudDlcDmsTable_Update tests the Update flow: AlterDMSTable succeeds, then Read is called
func TestResourceTencentCloudDlcDmsTable_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsTable().client, "UseDlcClient", dlcClient)

	var capturedAlterRequest *dlc.AlterDMSTableRequest
	patches.ApplyMethodFunc(dlcClient, "AlterDMSTableWithContext", func(_ context.Context, request *dlc.AlterDMSTableRequest) (*dlc.AlterDMSTableResponse, error) {
		capturedAlterRequest = request
		resp := dlc.NewAlterDMSTableResponse()
		resp.Response = &dlc.AlterDMSTableResponseParams{
			RequestId: ptrStringDmsTable("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeDMSTableWithContext for the Read call after Update
	patches.ApplyMethodFunc(dlcClient, "DescribeDMSTableWithContext", func(_ context.Context, request *dlc.DescribeDMSTableRequest) (*dlc.DescribeDMSTableResponse, error) {
		resp := dlc.NewDescribeDMSTableResponse()
		resp.Response = &dlc.DescribeDMSTableResponseParams{
			DbName:     ptrStringDmsTable("tf_example_db"),
			Name:       ptrStringDmsTable("tf_example_table_updated"),
			Type:       ptrStringDmsTable("EXTERNAL_TABLE"),
			SchemaName: ptrStringDmsTable("default"),
			Retention:  ptrInt64DmsTable(0),
			RequestId:  ptrStringDmsTable("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsTable()
	res := svcdlc.ResourceTencentCloudDlcDmsTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_name": "tf_example_db",
		"name":    "tf_example_table_updated",
		"type":    "EXTERNAL_TABLE",
	})
	// The old id was "tf_example_db#tf_example_table", simulating a rename from tf_example_table to tf_example_table_updated.
	// TestResourceDataRaw has no prior state, so HasChange("name") returns true, triggering the rename path:
	// CurrentName is taken from the id and Name from the config.
	d.SetId("tf_example_db#tf_example_table")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedAlterRequest)
	// The rename path sets CurrentName from the old id and Name from config
	assert.NotNil(t, capturedAlterRequest.CurrentName)
	assert.Equal(t, "tf_example_table", *capturedAlterRequest.CurrentName)
	assert.NotNil(t, capturedAlterRequest.Name)
	assert.Equal(t, "tf_example_table_updated", *capturedAlterRequest.Name)
	// Since name changed from tf_example_table to tf_example_table_updated, the id should be updated
	assert.Equal(t, "tf_example_db#tf_example_table_updated", d.Id())
}

// TestResourceTencentCloudDlcDmsTable_Delete tests the Delete flow: DropDMSTable succeeds
func TestResourceTencentCloudDlcDmsTable_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsTable().client, "UseDlcClient", dlcClient)

	var capturedDropRequest *dlc.DropDMSTableRequest
	patches.ApplyMethodFunc(dlcClient, "DropDMSTableWithContext", func(_ context.Context, request *dlc.DropDMSTableRequest) (*dlc.DropDMSTableResponse, error) {
		capturedDropRequest = request
		resp := dlc.NewDropDMSTableResponse()
		resp.Response = &dlc.DropDMSTableResponseParams{
			RequestId: ptrStringDmsTable("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsTable()
	res := svcdlc.ResourceTencentCloudDlcDmsTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_name":     "tf_example_db",
		"name":        "tf_example_table",
		"delete_data": true,
	})
	d.SetId("tf_example_db#tf_example_table")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedDropRequest)
	assert.NotNil(t, capturedDropRequest.DbName)
	assert.Equal(t, "tf_example_db", *capturedDropRequest.DbName)
	assert.NotNil(t, capturedDropRequest.Name)
	assert.Equal(t, "tf_example_table", *capturedDropRequest.Name)
}

// TestResourceTencentCloudDlcDmsTable_Schema tests the schema definition
func TestResourceTencentCloudDlcDmsTable_Schema(t *testing.T) {
	res := svcdlc.ResourceTencentCloudDlcDmsTable()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "db_name")
	assert.Contains(t, res.Schema, "name")
	assert.Contains(t, res.Schema, "asset")
	assert.Contains(t, res.Schema, "type")
	assert.Contains(t, res.Schema, "sds")
	assert.Contains(t, res.Schema, "columns")
	assert.Contains(t, res.Schema, "partition_keys")
	assert.Contains(t, res.Schema, "partitions")
	assert.Contains(t, res.Schema, "view_original_text")
	assert.Contains(t, res.Schema, "view_expanded_text")
	assert.Contains(t, res.Schema, "datasource_connection_name")
	assert.Contains(t, res.Schema, "delete_data")
	assert.Contains(t, res.Schema, "env_props")
	assert.Contains(t, res.Schema, "schema_name")
	assert.Contains(t, res.Schema, "retention")

	// db_name and name are Required
	assert.True(t, res.Schema["db_name"].Required)
	assert.True(t, res.Schema["name"].Required)

	// schema_name and retention are Computed only
	assert.True(t, res.Schema["schema_name"].Computed)
	assert.False(t, res.Schema["schema_name"].Optional)
	assert.True(t, res.Schema["retention"].Computed)
	assert.False(t, res.Schema["retention"].Optional)
}

// TestResourceTencentCloudDlcDmsTable_Create_EmptyResponse tests that Create returns an error when CreateDMSTable returns nil response
func TestResourceTencentCloudDlcDmsTable_Create_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsTable().client, "UseDlcClient", dlcClient)

	// Return a response with nil Response to trigger NonRetryableError
	patches.ApplyMethodFunc(dlcClient, "CreateDMSTableWithContext", func(_ context.Context, request *dlc.CreateDMSTableRequest) (*dlc.CreateDMSTableResponse, error) {
		resp := dlc.NewCreateDMSTableResponse()
		// Response is nil
		return resp, nil
	})

	meta := newMockMetaDlcDmsTable()
	res := svcdlc.ResourceTencentCloudDlcDmsTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_name": "tf_example_db",
		"name":    "tf_example_table",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestResourceTencentCloudDlcDmsTable_Delete_WithEnvProps tests that Delete correctly sends env_props
func TestResourceTencentCloudDlcDmsTable_Delete_WithEnvProps(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsTable().client, "UseDlcClient", dlcClient)

	var capturedDropRequest *dlc.DropDMSTableRequest
	patches.ApplyMethodFunc(dlcClient, "DropDMSTableWithContext", func(_ context.Context, request *dlc.DropDMSTableRequest) (*dlc.DropDMSTableResponse, error) {
		capturedDropRequest = request
		resp := dlc.NewDropDMSTableResponse()
		resp.Response = &dlc.DropDMSTableResponseParams{
			RequestId: ptrStringDmsTable("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsTable()
	res := svcdlc.ResourceTencentCloudDlcDmsTable()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"db_name":     "tf_example_db",
		"name":        "tf_example_table",
		"delete_data": true,
		"env_props": []interface{}{
			map[string]interface{}{"key": "env_key", "value": "env_value"},
		},
		"datasource_connection_name": "tf_example_connection",
	})
	d.SetId("tf_example_db#tf_example_table")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedDropRequest)
	assert.NotNil(t, capturedDropRequest.EnvProps)
	assert.Equal(t, "env_key", *capturedDropRequest.EnvProps.Key)
	assert.Equal(t, "env_value", *capturedDropRequest.EnvProps.Value)
	assert.NotNil(t, capturedDropRequest.DeleteData)
	assert.True(t, *capturedDropRequest.DeleteData)
	assert.NotNil(t, capturedDropRequest.DatasourceConnectionName)
	assert.Equal(t, "tf_example_connection", *capturedDropRequest.DatasourceConnectionName)
}
