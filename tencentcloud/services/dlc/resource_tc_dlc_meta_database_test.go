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

type mockMetaDlcMetaDatabase struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDlcMetaDatabase) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDlcMetaDatabase{}

func newMockMetaDlcMetaDatabase() *mockMetaDlcMetaDatabase {
	return &mockMetaDlcMetaDatabase{client: &connectivity.TencentCloudClient{}}
}

func ptrStrMetaDb(s string) *string {
	return &s
}

func ptrBoolMetaDb(b bool) *bool {
	return &b
}

// go test ./tencentcloud/services/dlc/ -run "TestResourceTencentCloudDlcMetaDatabase" -v -count=1 -gcflags="all=-l"

func TestResourceTencentCloudDlcMetaDatabaseCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcMetaDatabase().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.CreateMetaDatabaseRequest
	patches.ApplyMethodFunc(dlcClient, "CreateMetaDatabaseWithContext", func(_ context.Context, request *dlc.CreateMetaDatabaseRequest) (*dlc.CreateMetaDatabaseResponse, error) {
		capturedRequest = request
		resp := dlc.NewCreateMetaDatabaseResponse()
		resp.Response = &dlc.CreateMetaDatabaseResponseParams{
			BatchId:   ptrStrMetaDb("batch-123"),
			TaskIdSet: []*string{ptrStrMetaDb("task-1"), ptrStrMetaDb("task-2")},
			RequestId: ptrStrMetaDb("fake-request-id-create"),
		}
		return resp, nil
	})

	// mock DescribeDatabase for async polling after create
	patches.ApplyMethodFunc(dlcClient, "DescribeDatabaseWithContext", func(_ context.Context, request *dlc.DescribeDatabaseRequest) (*dlc.DescribeDatabaseResponse, error) {
		resp := dlc.NewDescribeDatabaseResponse()
		resp.Response = &dlc.DescribeDatabaseResponseParams{
			DatabaseInfo: &dlc.DatabaseResponseInfo{
				DatabaseName:        ptrStrMetaDb("tf_example_db"),
				Comment:             ptrStrMetaDb("tf example meta database"),
				CreateTime:          ptrStrMetaDb("1700000000"),
				ModifiedTime:        ptrStrMetaDb("1700000001"),
				Location:            ptrStrMetaDb("cosn://example/path"),
				UserAlias:           ptrStrMetaDb("user_alias"),
				UserSubUin:          ptrStrMetaDb("100000000001"),
				DatabaseId:          ptrStrMetaDb("db-123"),
				CatalogName:         ptrStrMetaDb("DataLakeCatalog"),
				CatalogType:         ptrStrMetaDb("DLC"),
				IsInformationSchema: ptrBoolMetaDb(false),
				Properties: []*dlc.Property{
					{
						Key:   ptrStrMetaDb("k1"),
						Value: ptrStrMetaDb("v1"),
					},
				},
				GovernPolicy: &dlc.DataGovernPolicy{
					RuleType:     ptrStrMetaDb("Customize"),
					GovernEngine: ptrStrMetaDb("engine_name"),
				},
			},
			RequestId: ptrStrMetaDb("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
		"comment":       "tf example meta database",
		"govern_policy": []interface{}{
			map[string]interface{}{
				"rule_type":     "Customize",
				"govern_engine": "engine_name",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.MetaDatabaseInfo)
	assert.Equal(t, "tf_example_db", *capturedRequest.MetaDatabaseInfo.DatabaseName)
	assert.Equal(t, "tf example meta database", *capturedRequest.MetaDatabaseInfo.Comment)
	assert.NotNil(t, capturedRequest.GovernPolicy)
	assert.Equal(t, "Customize", *capturedRequest.GovernPolicy.RuleType)

	assert.Equal(t, "tf_example_db", d.Id())
	assert.Equal(t, "tf_example_db", d.Get("database_name").(string))
	assert.Equal(t, "batch-123", d.Get("batch_id").(string))
	assert.Equal(t, "cosn://example/path", d.Get("location").(string))
	assert.Equal(t, "DataLakeCatalog", d.Get("catalog_name").(string))
	assert.False(t, d.Get("is_information_schema").(bool))
}

func TestResourceTencentCloudDlcMetaDatabaseCreateWithDatasourceConnectionName(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcMetaDatabase().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.CreateMetaDatabaseRequest
	patches.ApplyMethodFunc(dlcClient, "CreateMetaDatabaseWithContext", func(_ context.Context, request *dlc.CreateMetaDatabaseRequest) (*dlc.CreateMetaDatabaseResponse, error) {
		capturedRequest = request
		resp := dlc.NewCreateMetaDatabaseResponse()
		resp.Response = &dlc.CreateMetaDatabaseResponseParams{
			BatchId:   ptrStrMetaDb("batch-456"),
			TaskIdSet: []*string{ptrStrMetaDb("task-3")},
			RequestId: ptrStrMetaDb("fake-request-id-create"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dlcClient, "DescribeDatabaseWithContext", func(_ context.Context, request *dlc.DescribeDatabaseRequest) (*dlc.DescribeDatabaseResponse, error) {
		resp := dlc.NewDescribeDatabaseResponse()
		resp.Response = &dlc.DescribeDatabaseResponseParams{
			DatabaseInfo: &dlc.DatabaseResponseInfo{
				DatabaseName: ptrStrMetaDb("tf_example_db"),
			},
			RequestId: ptrStrMetaDb("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name":              "tf_example_db",
		"datasource_connection_name": "custom_connection",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.DatasourceConnectionName)
	assert.Equal(t, "custom_connection", *capturedRequest.DatasourceConnectionName)

	assert.Equal(t, "custom_connection#tf_example_db", d.Id())
	assert.Equal(t, "custom_connection", d.Get("datasource_connection_name").(string))
}

func TestResourceTencentCloudDlcMetaDatabaseCreateNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcMetaDatabase().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "CreateMetaDatabaseWithContext", func(_ context.Context, request *dlc.CreateMetaDatabaseRequest) (*dlc.CreateMetaDatabaseResponse, error) {
		resp := dlc.NewCreateMetaDatabaseResponse()
		resp.Response = &dlc.CreateMetaDatabaseResponseParams{
			BatchId: nil,
		}
		return resp, nil
	})

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

func TestResourceTencentCloudDlcMetaDatabaseRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcMetaDatabase().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.DescribeDatabaseRequest
	patches.ApplyMethodFunc(dlcClient, "DescribeDatabaseWithContext", func(_ context.Context, request *dlc.DescribeDatabaseRequest) (*dlc.DescribeDatabaseResponse, error) {
		capturedRequest = request
		resp := dlc.NewDescribeDatabaseResponse()
		resp.Response = &dlc.DescribeDatabaseResponseParams{
			DatabaseInfo: &dlc.DatabaseResponseInfo{
				DatabaseName:        ptrStrMetaDb("tf_example_db"),
				Comment:             ptrStrMetaDb("tf example meta database"),
				CreateTime:          ptrStrMetaDb("1700000000"),
				ModifiedTime:        ptrStrMetaDb("1700000001"),
				Location:            ptrStrMetaDb("cosn://example/path"),
				UserAlias:           ptrStrMetaDb("user_alias"),
				UserSubUin:          ptrStrMetaDb("100000000001"),
				DatabaseId:          ptrStrMetaDb("db-123"),
				CatalogName:         ptrStrMetaDb("DataLakeCatalog"),
				CatalogType:         ptrStrMetaDb("DLC"),
				IsInformationSchema: ptrBoolMetaDb(false),
				Properties: []*dlc.Property{
					{
						Key:   ptrStrMetaDb("k1"),
						Value: ptrStrMetaDb("v1"),
					},
				},
				GovernPolicy: &dlc.DataGovernPolicy{
					RuleType:     ptrStrMetaDb("Customize"),
					GovernEngine: ptrStrMetaDb("engine_name"),
				},
			},
			RequestId: ptrStrMetaDb("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("tf_example_db")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "tf_example_db", *capturedRequest.DatabaseName)
	assert.Nil(t, capturedRequest.DatasourceConnectionName)

	assert.Equal(t, "tf_example_db", d.Id())
	assert.Equal(t, "tf_example_db", d.Get("database_name").(string))
	assert.Equal(t, "cosn://example/path", d.Get("location").(string))
	assert.Equal(t, "1700000000", d.Get("create_time").(string))
	assert.Equal(t, "100000000001", d.Get("user_sub_uin").(string))
	assert.Equal(t, "DLC", d.Get("catalog_type").(string))
	assert.False(t, d.Get("is_information_schema").(bool))

	properties := d.Get("properties").([]interface{})
	assert.Len(t, properties, 1)
	propMap := properties[0].(map[string]interface{})
	assert.Equal(t, "k1", propMap["key"])
	assert.Equal(t, "v1", propMap["value"])

	governPolicy := d.Get("govern_policy").([]interface{})
	assert.Len(t, governPolicy, 1)
	gpMap := governPolicy[0].(map[string]interface{})
	assert.Equal(t, "Customize", gpMap["rule_type"])
	assert.Equal(t, "engine_name", gpMap["govern_engine"])
}

func TestResourceTencentCloudDlcMetaDatabaseReadCompositeId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcMetaDatabase().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.DescribeDatabaseRequest
	patches.ApplyMethodFunc(dlcClient, "DescribeDatabaseWithContext", func(_ context.Context, request *dlc.DescribeDatabaseRequest) (*dlc.DescribeDatabaseResponse, error) {
		capturedRequest = request
		resp := dlc.NewDescribeDatabaseResponse()
		resp.Response = &dlc.DescribeDatabaseResponseParams{
			DatabaseInfo: &dlc.DatabaseResponseInfo{
				DatabaseName: ptrStrMetaDb("tf_example_db"),
			},
			RequestId: ptrStrMetaDb("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("custom_connection#tf_example_db")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "tf_example_db", *capturedRequest.DatabaseName)
	assert.NotNil(t, capturedRequest.DatasourceConnectionName)
	assert.Equal(t, "custom_connection", *capturedRequest.DatasourceConnectionName)

	assert.Equal(t, "custom_connection", d.Get("datasource_connection_name").(string))
}

func TestResourceTencentCloudDlcMetaDatabaseReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcMetaDatabase().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeDatabaseWithContext", func(_ context.Context, request *dlc.DescribeDatabaseRequest) (*dlc.DescribeDatabaseResponse, error) {
		resp := dlc.NewDescribeDatabaseResponse()
		resp.Response = &dlc.DescribeDatabaseResponseParams{
			DatabaseInfo: nil,
			RequestId:    ptrStrMetaDb("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("tf_example_db")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestResourceTencentCloudDlcMetaDatabaseDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcMetaDatabase().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.DeleteMetaDatabaseRequest
	patches.ApplyMethodFunc(dlcClient, "DeleteMetaDatabaseWithContext", func(_ context.Context, request *dlc.DeleteMetaDatabaseRequest) (*dlc.DeleteMetaDatabaseResponse, error) {
		capturedRequest = request
		resp := dlc.NewDeleteMetaDatabaseResponse()
		resp.Response = &dlc.DeleteMetaDatabaseResponseParams{
			BatchId:   ptrStrMetaDb("batch-del"),
			TaskIdSet: []*string{ptrStrMetaDb("task-del-1")},
			RequestId: ptrStrMetaDb("fake-request-id-delete"),
		}
		return resp, nil
	})

	// mock DescribeDatabase for async polling - return nil DatabaseInfo to signal deleted
	patches.ApplyMethodFunc(dlcClient, "DescribeDatabaseWithContext", func(_ context.Context, request *dlc.DescribeDatabaseRequest) (*dlc.DescribeDatabaseResponse, error) {
		resp := dlc.NewDescribeDatabaseResponse()
		resp.Response = &dlc.DescribeDatabaseResponseParams{
			DatabaseInfo: nil,
			RequestId:    ptrStrMetaDb("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("custom_connection#tf_example_db")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "tf_example_db", *capturedRequest.DatabaseName)
	assert.NotNil(t, capturedRequest.DatasourceConnectionName)
	assert.Equal(t, "custom_connection", *capturedRequest.DatasourceConnectionName)
}

func TestResourceTencentCloudDlcMetaDatabaseDeleteBareId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcMetaDatabase().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.DeleteMetaDatabaseRequest
	patches.ApplyMethodFunc(dlcClient, "DeleteMetaDatabaseWithContext", func(_ context.Context, request *dlc.DeleteMetaDatabaseRequest) (*dlc.DeleteMetaDatabaseResponse, error) {
		capturedRequest = request
		resp := dlc.NewDeleteMetaDatabaseResponse()
		resp.Response = &dlc.DeleteMetaDatabaseResponseParams{
			BatchId:   ptrStrMetaDb("batch-del"),
			TaskIdSet: []*string{ptrStrMetaDb("task-del-1")},
			RequestId: ptrStrMetaDb("fake-request-id-delete"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dlcClient, "DescribeDatabaseWithContext", func(_ context.Context, request *dlc.DescribeDatabaseRequest) (*dlc.DescribeDatabaseResponse, error) {
		resp := dlc.NewDescribeDatabaseResponse()
		resp.Response = &dlc.DescribeDatabaseResponseParams{
			DatabaseInfo: nil,
			RequestId:    ptrStrMetaDb("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("tf_example_db")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "tf_example_db", *capturedRequest.DatabaseName)
	assert.Nil(t, capturedRequest.DatasourceConnectionName)
}

func TestResourceTencentCloudDlcMetaDatabaseUpdateImmutableChanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
		"comment":       "new comment",
	})
	d.SetId("tf_example_db")

	// Patch HasChange to simulate a change in comment
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "comment"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "argument `comment` cannot be changed")
}

func TestResourceTencentCloudDlcMetaDatabaseUpdateNoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcMetaDatabase().client, "UseDlcClient", dlcClient)

	// mock DescribeDatabase for the Read call after Update (no changes)
	patches.ApplyMethodFunc(dlcClient, "DescribeDatabaseWithContext", func(_ context.Context, request *dlc.DescribeDatabaseRequest) (*dlc.DescribeDatabaseResponse, error) {
		resp := dlc.NewDescribeDatabaseResponse()
		resp.Response = &dlc.DescribeDatabaseResponseParams{
			DatabaseInfo: &dlc.DatabaseResponseInfo{
				DatabaseName: ptrStrMetaDb("tf_example_db"),
			},
			RequestId: ptrStrMetaDb("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
	})
	d.SetId("tf_example_db")

	// Patch HasChange to simulate no changes
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestResourceTencentCloudDlcMetaDatabaseSchema(t *testing.T) {
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "database_name")
	assert.Contains(t, res.Schema, "datasource_connection_name")
	assert.Contains(t, res.Schema, "comment")
	assert.Contains(t, res.Schema, "govern_policy")
	assert.Contains(t, res.Schema, "smart_policy")
	assert.Contains(t, res.Schema, "batch_id")
	assert.Contains(t, res.Schema, "task_id_set")
	assert.Contains(t, res.Schema, "properties")
	assert.Contains(t, res.Schema, "create_time")
	assert.Contains(t, res.Schema, "modified_time")
	assert.Contains(t, res.Schema, "location")
	assert.Contains(t, res.Schema, "user_alias")
	assert.Contains(t, res.Schema, "user_sub_uin")
	assert.Contains(t, res.Schema, "database_id")
	assert.Contains(t, res.Schema, "catalog_name")
	assert.Contains(t, res.Schema, "catalog_type")
	assert.Contains(t, res.Schema, "is_information_schema")

	dbName := res.Schema["database_name"]
	assert.Equal(t, schema.TypeString, dbName.Type)
	assert.True(t, dbName.Required)
	assert.True(t, dbName.ForceNew)

	governPolicy := res.Schema["govern_policy"]
	assert.Equal(t, schema.TypeList, governPolicy.Type)
	assert.True(t, governPolicy.Optional)
	assert.Equal(t, 1, governPolicy.MaxItems)

	smartPolicy := res.Schema["smart_policy"]
	assert.Equal(t, schema.TypeList, smartPolicy.Type)
	assert.True(t, smartPolicy.Optional)
	assert.Equal(t, 1, smartPolicy.MaxItems)

	batchId := res.Schema["batch_id"]
	assert.Equal(t, schema.TypeString, batchId.Type)
	assert.True(t, batchId.Computed)

	isInfoSchema := res.Schema["is_information_schema"]
	assert.Equal(t, schema.TypeBool, isInfoSchema.Type)
	assert.True(t, isInfoSchema.Computed)

	assert.NotNil(t, res.Importer)
}

func TestResourceTencentCloudDlcMetaDatabaseCreateFullSmartPolicy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaDlcMetaDatabase().client, "UseDlcClient", dlcClient)

	var capturedRequest *dlc.CreateMetaDatabaseRequest
	patches.ApplyMethodFunc(dlcClient, "CreateMetaDatabaseWithContext", func(_ context.Context, request *dlc.CreateMetaDatabaseRequest) (*dlc.CreateMetaDatabaseResponse, error) {
		capturedRequest = request
		resp := dlc.NewCreateMetaDatabaseResponse()
		resp.Response = &dlc.CreateMetaDatabaseResponseParams{
			BatchId:   ptrStrMetaDb("batch-789"),
			TaskIdSet: []*string{ptrStrMetaDb("task-9")},
			RequestId: ptrStrMetaDb("fake-request-id-create"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(dlcClient, "DescribeDatabaseWithContext", func(_ context.Context, request *dlc.DescribeDatabaseRequest) (*dlc.DescribeDatabaseResponse, error) {
		resp := dlc.NewDescribeDatabaseResponse()
		resp.Response = &dlc.DescribeDatabaseResponseParams{
			DatabaseInfo: &dlc.DatabaseResponseInfo{
				DatabaseName: ptrStrMetaDb("tf_example_db"),
			},
			RequestId: ptrStrMetaDb("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaDlcMetaDatabase()
	res := svcdlc.ResourceTencentCloudDlcMetaDatabase()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"database_name": "tf_example_db",
		"smart_policy": []interface{}{
			map[string]interface{}{
				"base_info": []interface{}{
					map[string]interface{}{
						"uin":         "100000000001",
						"policy_type": "TABLE",
						"catalog":     "DataLakeCatalog",
					},
				},
				"policy": []interface{}{
					map[string]interface{}{
						"inherit": "true",
						"resources": []interface{}{
							map[string]interface{}{
								"attribution_type":    "CATALOG",
								"resource_type":       "spark-sql",
								"name":                "engine_name",
								"status":              1,
								"resource_group_name": "rg_name",
								"favor": []interface{}{
									map[string]interface{}{
										"priority": 5,
										"catalog":  "DataLakeCatalog",
									},
								},
								"resource_conf": []interface{}{
									map[string]interface{}{
										"parallelism": 4,
									},
								},
							},
						},
						"written": []interface{}{
							map[string]interface{}{
								"written_enable": "enable",
								"advance_policy": []interface{}{
									map[string]interface{}{
										"compact_enable":         "enable",
										"delete_enable":          "enable",
										"min_input_files":        5,
										"target_file_size_bytes": 1024,
										"sort_orders": []interface{}{
											map[string]interface{}{
												"column":         "col1",
												"sort_direction": "asc",
												"null_order":     "first",
											},
										},
									},
								},
							},
						},
						"lifecycle": []interface{}{
							map[string]interface{}{
								"lifecycle_enable":     "enable",
								"expired_field":        "dt",
								"expired_field_format": "yyyy-MM-dd",
								"expiration":           7,
								"drop_table":           true,
							},
						},
						"index": []interface{}{
							map[string]interface{}{
								"index_enable": "enable",
							},
						},
						"change_table": []interface{}{
							map[string]interface{}{
								"data_retention_time": 30,
							},
						},
						"table_expiration": []interface{}{
							map[string]interface{}{
								"enabled":    true,
								"expiration": 30,
							},
						},
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.SmartPolicy)
	assert.NotNil(t, capturedRequest.SmartPolicy.BaseInfo)
	assert.Equal(t, "100000000001", *capturedRequest.SmartPolicy.BaseInfo.Uin)
	assert.Equal(t, "TABLE", *capturedRequest.SmartPolicy.BaseInfo.PolicyType)
	assert.Equal(t, "DataLakeCatalog", *capturedRequest.SmartPolicy.BaseInfo.Catalog)

	policy := capturedRequest.SmartPolicy.Policy
	assert.NotNil(t, policy)
	assert.Equal(t, "true", *policy.Inherit)
	assert.Len(t, policy.Resources, 1)
	assert.Equal(t, "CATALOG", *policy.Resources[0].AttributionType)
	assert.Equal(t, int64(1), *policy.Resources[0].Status)
	assert.Equal(t, "rg_name", *policy.Resources[0].ResourceGroupName)
	assert.Len(t, policy.Resources[0].Favor, 1)
	assert.Equal(t, int64(5), *policy.Resources[0].Favor[0].Priority)
	assert.NotNil(t, policy.Resources[0].ResourceConf)
	assert.Equal(t, int64(4), *policy.Resources[0].ResourceConf.Parallelism)

	assert.NotNil(t, policy.Written)
	assert.Equal(t, "enable", *policy.Written.WrittenEnable)
	assert.NotNil(t, policy.Written.AdvancePolicy)
	assert.Equal(t, "enable", *policy.Written.AdvancePolicy.CompactEnable)
	assert.Equal(t, int64(5), *policy.Written.AdvancePolicy.MinInputFiles)
	assert.Len(t, policy.Written.AdvancePolicy.SortOrders, 1)
	assert.Equal(t, "col1", *policy.Written.AdvancePolicy.SortOrders[0].Column)

	assert.NotNil(t, policy.Lifecycle)
	assert.Equal(t, "enable", *policy.Lifecycle.LifecycleEnable)
	assert.Equal(t, "dt", *policy.Lifecycle.ExpiredField)
	assert.Equal(t, int64(7), *policy.Lifecycle.Expiration)
	assert.True(t, *policy.Lifecycle.DropTable)

	assert.NotNil(t, policy.Index)
	assert.Equal(t, "enable", *policy.Index.IndexEnable)

	assert.NotNil(t, policy.ChangeTable)
	assert.Equal(t, int64(30), *policy.ChangeTable.DataRetentionTime)

	assert.NotNil(t, policy.TableExpiration)
	assert.True(t, *policy.TableExpiration.Enabled)
	assert.Equal(t, uint64(30), *policy.TableExpiration.Expiration)
}
