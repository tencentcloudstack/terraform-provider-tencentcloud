package dlc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	dlcsvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"
)

type mockMetaDlcDmsDatabase struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDlcDmsDatabase) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDlcDmsDatabase{}

func newMockMetaDlcDmsDatabase() *mockMetaDlcDmsDatabase {
	return &mockMetaDlcDmsDatabase{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDlcDms(s string) *string {
	return &s
}

func ptrInt64DlcDms(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/dlc/ -run "TestDlcDmsDatabase" -v -count=1 -gcflags="all=-l"

// TestDlcDmsDatabase_Create tests the Create operation
func TestDlcDmsDatabase_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsDatabase().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "CreateDMSDatabaseWithContext", func(ctx context.Context, request *dlc.CreateDMSDatabaseRequest) (*dlc.CreateDMSDatabaseResponse, error) {
		assert.Equal(t, "mydb", *request.Name)
		assert.Equal(t, "myschema", *request.SchemaName)
		assert.Equal(t, "conn1", *request.DatasourceConnectionName)
		assert.Equal(t, "cosn://bucket/path", *request.Location)

		resp := dlc.NewCreateDMSDatabaseResponse()
		resp.Response = &dlc.CreateDMSDatabaseResponseParams{
			RequestId: ptrStringDlcDms("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dlcClient, "DescribeDMSDatabaseWithContext", func(ctx context.Context, request *dlc.DescribeDMSDatabaseRequest) (*dlc.DescribeDMSDatabaseResponse, error) {
		resp := dlc.NewDescribeDMSDatabaseResponse()
		resp.Response = &dlc.DescribeDMSDatabaseResponseParams{
			Name:       ptrStringDlcDms("mydb"),
			SchemaName: ptrStringDlcDms("myschema"),
			Location:   ptrStringDlcDms("cosn://bucket/path"),
			Asset: &dlc.Asset{
				Id:          ptrInt64DlcDms(1),
				Name:        ptrStringDlcDms("asset-name"),
				Guid:        ptrStringDlcDms("guid-001"),
				Catalog:     ptrStringDlcDms("catalog"),
				DataVersion: ptrInt64DlcDms(2),
				Params: []*dlc.KVPair{
					{Key: ptrStringDlcDms("k"), Value: ptrStringDlcDms("v")},
				},
			},
			RequestId: ptrStringDlcDms("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsDatabase()
	res := dlcsvc.ResourceTencentCloudDlcDmsDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                       "mydb",
		"schema_name":                "myschema",
		"datasource_connection_name": "conn1",
		"location":                   "cosn://bucket/path",
		"asset": []interface{}{
			map[string]interface{}{
				"name":    "asset-name",
				"catalog": "catalog",
				"params": []interface{}{
					map[string]interface{}{
						"key":   "k",
						"value": "v",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mydb#myschema#conn1", d.Id())
	assert.Equal(t, "mydb", d.Get("name").(string))
	assert.Equal(t, "myschema", d.Get("schema_name").(string))
	assert.Equal(t, "conn1", d.Get("datasource_connection_name").(string))
	assert.Equal(t, "cosn://bucket/path", d.Get("location").(string))
}

// TestDlcDmsDatabase_Read tests the Read operation
func TestDlcDmsDatabase_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsDatabase().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeDMSDatabaseWithContext", func(ctx context.Context, request *dlc.DescribeDMSDatabaseRequest) (*dlc.DescribeDMSDatabaseResponse, error) {
		assert.Equal(t, "mydb", *request.Name)
		assert.Equal(t, "myschema", *request.SchemaName)
		assert.Equal(t, "conn1", *request.DatasourceConnectionName)

		resp := dlc.NewDescribeDMSDatabaseResponse()
		resp.Response = &dlc.DescribeDMSDatabaseResponseParams{
			Name:       ptrStringDlcDms("mydb"),
			SchemaName: ptrStringDlcDms("myschema"),
			Location:   ptrStringDlcDms("cosn://bucket/path"),
			Asset: &dlc.Asset{
				Id:           ptrInt64DlcDms(1),
				Name:         ptrStringDlcDms("asset-name"),
				Guid:         ptrStringDlcDms("guid-001"),
				Catalog:      ptrStringDlcDms("catalog"),
				Description:  ptrStringDlcDms("desc"),
				Owner:        ptrStringDlcDms("owner"),
				OwnerAccount: ptrStringDlcDms("owner-account"),
				DataVersion:  ptrInt64DlcDms(2),
				CreateTime:   ptrStringDlcDms("2024-01-01 00:00:00"),
				ModifiedTime: ptrStringDlcDms("2024-01-02 00:00:00"),
				DatasourceId: ptrInt64DlcDms(3),
				Params: []*dlc.KVPair{
					{Key: ptrStringDlcDms("pk"), Value: ptrStringDlcDms("pv")},
				},
				BizParams: []*dlc.KVPair{
					{Key: ptrStringDlcDms("bk"), Value: ptrStringDlcDms("bv")},
				},
				PermValues: []*dlc.KVPair{
					{Key: ptrStringDlcDms("permk"), Value: ptrStringDlcDms("permv")},
				},
			},
			RequestId: ptrStringDlcDms("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsDatabase()
	res := dlcsvc.ResourceTencentCloudDlcDmsDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                       "",
		"schema_name":                "",
		"datasource_connection_name": "",
		"location":                   "",
	})
	d.SetId("mydb#myschema#conn1")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mydb#myschema#conn1", d.Id())
	assert.Equal(t, "mydb", d.Get("name").(string))
	assert.Equal(t, "myschema", d.Get("schema_name").(string))
	assert.Equal(t, "conn1", d.Get("datasource_connection_name").(string))
	assert.Equal(t, "cosn://bucket/path", d.Get("location").(string))

	assetList := d.Get("asset").([]interface{})
	assert.Len(t, assetList, 1)
	assetMap := assetList[0].(map[string]interface{})
	assert.Equal(t, "asset-name", assetMap["name"].(string))
	assert.Equal(t, "guid-001", assetMap["guid"].(string))
	assert.Equal(t, 1, assetMap["id"].(int))
	assert.Equal(t, 2, assetMap["data_version"].(int))

	paramsList := assetMap["params"].([]interface{})
	assert.Len(t, paramsList, 1)
	paramsMap := paramsList[0].(map[string]interface{})
	assert.Equal(t, "pk", paramsMap["key"].(string))
	assert.Equal(t, "pv", paramsMap["value"].(string))
}

// TestDlcDmsDatabase_Read_NotFound tests Read when database is not found
func TestDlcDmsDatabase_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsDatabase().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeDMSDatabaseWithContext", func(ctx context.Context, request *dlc.DescribeDMSDatabaseRequest) (*dlc.DescribeDMSDatabaseResponse, error) {
		resp := dlc.NewDescribeDMSDatabaseResponse()
		resp.Response = &dlc.DescribeDMSDatabaseResponseParams{
			RequestId: ptrStringDlcDms("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsDatabase()
	res := dlcsvc.ResourceTencentCloudDlcDmsDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                       "",
		"schema_name":                "",
		"datasource_connection_name": "",
		"location":                   "",
	})
	d.SetId("notfound#schema#conn")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestDlcDmsDatabase_Update tests the Update operation
func TestDlcDmsDatabase_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsDatabase().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "AlterDMSDatabaseWithContext", func(ctx context.Context, request *dlc.AlterDMSDatabaseRequest) (*dlc.AlterDMSDatabaseResponse, error) {
		assert.Equal(t, "mydb", *request.CurrentName)
		assert.Equal(t, "myschema", *request.SchemaName)
		assert.Equal(t, "conn1", *request.DatasourceConnectionName)
		assert.Equal(t, "cosn://bucket/new", *request.Location)
		assert.NotNil(t, request.Asset)

		resp := dlc.NewAlterDMSDatabaseResponse()
		resp.Response = &dlc.AlterDMSDatabaseResponseParams{
			RequestId: ptrStringDlcDms("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dlcClient, "DescribeDMSDatabaseWithContext", func(ctx context.Context, request *dlc.DescribeDMSDatabaseRequest) (*dlc.DescribeDMSDatabaseResponse, error) {
		resp := dlc.NewDescribeDMSDatabaseResponse()
		resp.Response = &dlc.DescribeDMSDatabaseResponseParams{
			Name:       ptrStringDlcDms("mydb"),
			SchemaName: ptrStringDlcDms("myschema"),
			Location:   ptrStringDlcDms("cosn://bucket/new"),
			Asset: &dlc.Asset{
				Id:   ptrInt64DlcDms(1),
				Name: ptrStringDlcDms("asset-name-new"),
			},
			RequestId: ptrStringDlcDms("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsDatabase()
	res := dlcsvc.ResourceTencentCloudDlcDmsDatabase()

	// Build prior state where location=cosn://bucket/old, then change to
	// cosn://bucket/new in the new config so d.HasChange("location") is true
	// while the immutable id fields stay unchanged.
	state := &terraform.InstanceState{
		ID: "mydb#myschema#conn1",
		Attributes: map[string]string{
			"id":                         "mydb#myschema#conn1",
			"name":                       "mydb",
			"schema_name":                "myschema",
			"datasource_connection_name": "conn1",
			"location":                   "cosn://bucket/old",
			"delete_data":                "false",
			"cascade":                    "false",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":                       "mydb",
		"schema_name":                "myschema",
		"datasource_connection_name": "conn1",
		"location":                   "cosn://bucket/new",
		"asset": []interface{}{
			map[string]interface{}{
				"name": "asset-name-new",
			},
		},
	})

	diff, err := res.Diff(nil, state, rawConfig, meta)
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mydb#myschema#conn1", d.Id())
	assert.Equal(t, "cosn://bucket/new", d.Get("location").(string))
}

// TestDlcDmsDatabase_Update_ImmutableChange tests that changing an immutable
// argument during update returns an error.
func TestDlcDmsDatabase_Update_ImmutableChange(t *testing.T) {
	res := dlcsvc.ResourceTencentCloudDlcDmsDatabase()

	state := &terraform.InstanceState{
		ID: "mydb#myschema#conn1",
		Attributes: map[string]string{
			"id":                         "mydb#myschema#conn1",
			"name":                       "mydb",
			"schema_name":                "myschema",
			"datasource_connection_name": "conn1",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":                       "newdb",
		"schema_name":                "myschema",
		"datasource_connection_name": "conn1",
	})

	diff, err := res.Diff(nil, state, rawConfig, newMockMetaDlcDmsDatabase())
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	updateErr := res.Update(d, newMockMetaDlcDmsDatabase())
	assert.Error(t, updateErr)
	assert.Contains(t, updateErr.Error(), "immutable")
}

// TestDlcDmsDatabase_Delete tests the Delete operation
func TestDlcDmsDatabase_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsDatabase().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DropDMSDatabaseWithContext", func(ctx context.Context, request *dlc.DropDMSDatabaseRequest) (*dlc.DropDMSDatabaseResponse, error) {
		assert.Equal(t, "mydb", *request.Name)
		assert.Equal(t, "conn1", *request.DatasourceConnectionName)
		assert.True(t, *request.DeleteData)
		assert.True(t, *request.Cascade)

		resp := dlc.NewDropDMSDatabaseResponse()
		resp.Response = &dlc.DropDMSDatabaseResponseParams{
			RequestId: ptrStringDlcDms("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsDatabase()
	res := dlcsvc.ResourceTencentCloudDlcDmsDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                       "mydb",
		"schema_name":                "myschema",
		"datasource_connection_name": "conn1",
		"delete_data":                true,
		"cascade":                    true,
	})
	d.SetId("mydb#myschema#conn1")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestDlcDmsDatabase_Delete_Defaults tests Delete with default false values
func TestDlcDmsDatabase_Delete_Defaults(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsDatabase().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DropDMSDatabaseWithContext", func(ctx context.Context, request *dlc.DropDMSDatabaseRequest) (*dlc.DropDMSDatabaseResponse, error) {
		assert.Equal(t, "mydb", *request.Name)
		assert.Equal(t, "conn1", *request.DatasourceConnectionName)
		assert.False(t, *request.DeleteData)
		assert.False(t, *request.Cascade)

		resp := dlc.NewDropDMSDatabaseResponse()
		resp.Response = &dlc.DropDMSDatabaseResponseParams{
			RequestId: ptrStringDlcDms("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsDatabase()
	res := dlcsvc.ResourceTencentCloudDlcDmsDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                       "mydb",
		"schema_name":                "myschema",
		"datasource_connection_name": "conn1",
	})
	d.SetId("mydb#myschema#conn1")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestDlcDmsDatabase_Import tests import by reading a resource set via composite ID
func TestDlcDmsDatabase_Import(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcDmsDatabase().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeDMSDatabaseWithContext", func(ctx context.Context, request *dlc.DescribeDMSDatabaseRequest) (*dlc.DescribeDMSDatabaseResponse, error) {
		assert.Equal(t, "impdb", *request.Name)
		assert.Equal(t, "impschema", *request.SchemaName)
		assert.Equal(t, "impconn", *request.DatasourceConnectionName)

		resp := dlc.NewDescribeDMSDatabaseResponse()
		resp.Response = &dlc.DescribeDMSDatabaseResponseParams{
			Name:       ptrStringDlcDms("impdb"),
			SchemaName: ptrStringDlcDms("impschema"),
			Location:   ptrStringDlcDms("cosn://bucket/imp"),
			RequestId:  ptrStringDlcDms("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcDmsDatabase()
	res := dlcsvc.ResourceTencentCloudDlcDmsDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                       "",
		"schema_name":                "",
		"datasource_connection_name": "",
		"location":                   "",
	})
	d.SetId("impdb#impschema#impconn")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "impdb#impschema#impconn", d.Id())
	assert.Equal(t, "impdb", d.Get("name").(string))
	assert.Equal(t, "impschema", d.Get("schema_name").(string))
	assert.Equal(t, "impconn", d.Get("datasource_connection_name").(string))
}
