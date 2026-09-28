package mongodb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	mongodb_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"
)

func ptrStringMongodbRestoreDbInstance(s string) *string {
	return &s
}

func mongodbRestoreDbInstanceSchema(t *testing.T) action_schema.Schema {
	a := mongodb.NewMongodbRestoreDbInstance()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// timeoutsObjType builds the tftypes.Object type for the timeouts block.
func timeoutsObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"invoke": tftypes.String,
		},
	}
}

// nullTimeoutsValue is the null value for an unconfigured timeouts block.
func nullTimeoutsValue() tftypes.Value {
	return tftypes.NewValue(timeoutsObjType(), nil)
}

// collectionObjType builds the tftypes.Object type for a collections block
// element.
func collectionObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"old_collection": tftypes.String,
			"new_collection": tftypes.String,
		},
	}
}

// databaseObjType builds the tftypes.Object type for a databases block element.
func databaseObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"db":          tftypes.String,
			"collections": tftypes.List{ElementType: collectionObjType()},
		},
	}
}

// newCollectionValue builds a tftypes.Value for one collections element.
func newCollectionValue(oldCol, newCol string) tftypes.Value {
	return tftypes.NewValue(collectionObjType(), map[string]tftypes.Value{
		"old_collection": tftypes.NewValue(tftypes.String, oldCol),
		"new_collection": tftypes.NewValue(tftypes.String, newCol),
	})
}

// newDatabaseValue builds a tftypes.Value for one databases element.
func newDatabaseValue(db string, collections []tftypes.Value) tftypes.Value {
	colList := tftypes.NewValue(tftypes.List{ElementType: collectionObjType()}, nil)
	if len(collections) > 0 {
		colList = tftypes.NewValue(tftypes.List{ElementType: collectionObjType()}, collections)
	}
	return tftypes.NewValue(databaseObjType(), map[string]tftypes.Value{
		"db":          tftypes.NewValue(tftypes.String, db),
		"collections": colList,
	})
}

// newMongodbRestoreDbInstanceInvokeRequest builds an InvokeRequest whose
// timeouts block is left unconfigured, so the action's built-in default (15m)
// applies.
func newMongodbRestoreDbInstanceInvokeRequest(t *testing.T, instanceId, restoreTime string, databases []tftypes.Value) action.InvokeRequest {
	return newMongodbRestoreDbInstanceInvokeRequestWithTimeouts(t, instanceId, restoreTime, databases, nullTimeoutsValue())
}

// newMongodbRestoreDbInstanceInvokeRequestWithTimeouts builds an InvokeRequest
// carrying the given raw value for the timeouts block.
func newMongodbRestoreDbInstanceInvokeRequestWithTimeouts(t *testing.T, instanceId, restoreTime string, databases []tftypes.Value, timeoutsVal tftypes.Value) action.InvokeRequest {
	s := mongodbRestoreDbInstanceSchema(t)

	dbListType := tftypes.List{ElementType: databaseObjType()}
	var dbListVal tftypes.Value
	if len(databases) == 0 {
		dbListVal = tftypes.NewValue(dbListType, nil)
	} else {
		dbListVal = tftypes.NewValue(dbListType, databases)
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"instance_id":  tftypes.String,
			"restore_time": tftypes.String,
			"databases":    dbListType,
			"timeouts":     timeoutsObjType(),
		},
	}, map[string]tftypes.Value{
		"instance_id":  tftypes.NewValue(tftypes.String, instanceId),
		"restore_time": tftypes.NewValue(tftypes.String, restoreTime),
		"databases":    dbListVal,
		"timeouts":     timeoutsVal,
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

func setMongodbRestoreDbInstanceClient(a *mongodb.MongodbRestoreDbInstance, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/mongodb/ -run "TestMongodbRestoreDbInstance" -v -count=1 -gcflags="all=-l"

// TestMongodbRestoreDbInstance_Invoke_Success tests a successful invoke
func TestMongodbRestoreDbInstance_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseMongodbClient", mongodbClient)

	flowId := int64(12345)
	var gotRequest *mongodb_sdk.RestoreDBInstanceRequest
	patches.ApplyMethodFunc(mongodbClient, "RestoreDBInstanceWithContext", func(_ context.Context, request *mongodb_sdk.RestoreDBInstanceRequest) (*mongodb_sdk.RestoreDBInstanceResponse, error) {
		gotRequest = request
		resp := mongodb_sdk.NewRestoreDBInstanceResponse()
		resp.Response = &mongodb_sdk.RestoreDBInstanceResponseParams{
			FlowId:    &flowId,
			RequestId: ptrStringMongodbRestoreDbInstance("fake-request-id"),
		}
		return resp, nil
	})

	successStatus := "success"
	patches.ApplyMethodFunc(mongodbClient, "DescribeAsyncRequestInfo", func(request *mongodb_sdk.DescribeAsyncRequestInfoRequest) (*mongodb_sdk.DescribeAsyncRequestInfoResponse, error) {
		resp := mongodb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &mongodb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status: &successStatus,
		}
		return resp, nil
	})

	aImpl := mongodb.NewMongodbRestoreDbInstance()
	a := aImpl.(*mongodb.MongodbRestoreDbInstance)
	setMongodbRestoreDbInstanceClient(a, client)

	req := newMongodbRestoreDbInstanceInvokeRequest(t,
		"cmgo-xxxxxxxx",
		"2024-09-01 12:00:00",
		[]tftypes.Value{
			newDatabaseValue("db1", []tftypes.Value{
				newCollectionValue("col_old_1", "col_new_1"),
				newCollectionValue("col_old_2", "col_new_2"),
			}),
		},
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, "cmgo-xxxxxxxx", *gotRequest.InstanceId)
	assert.Equal(t, "2024-09-01 12:00:00", *gotRequest.RestoreTime)
	assert.Len(t, gotRequest.Databases, 1)
	assert.Equal(t, "db1", *gotRequest.Databases[0].Db)
	assert.Len(t, gotRequest.Databases[0].Collections, 2)
	assert.Equal(t, "col_old_1", *gotRequest.Databases[0].Collections[0].OldCollection)
	assert.Equal(t, "col_new_1", *gotRequest.Databases[0].Collections[0].NewCollection)
	assert.Equal(t, "col_old_2", *gotRequest.Databases[0].Collections[1].OldCollection)
	assert.Equal(t, "col_new_2", *gotRequest.Databases[0].Collections[1].NewCollection)
}

// TestMongodbRestoreDbInstance_Invoke_APIError tests API error handling
func TestMongodbRestoreDbInstance_Invoke_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "RestoreDBInstanceWithContext", func(_ context.Context, request *mongodb_sdk.RestoreDBInstanceRequest) (*mongodb_sdk.RestoreDBInstanceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameterValue, Message=restore time out of backup retention")
	})

	aImpl := mongodb.NewMongodbRestoreDbInstance()
	a := aImpl.(*mongodb.MongodbRestoreDbInstance)
	setMongodbRestoreDbInstanceClient(a, client)

	req := newMongodbRestoreDbInstanceInvokeRequest(t,
		"cmgo-invalid",
		"2024-09-01 12:00:00",
		[]tftypes.Value{
			newDatabaseValue("db1", []tftypes.Value{
				newCollectionValue("col_old", "col_new"),
			}),
		},
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "restoring mongodb instance")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "InvalidParameterValue")
}

// TestMongodbRestoreDbInstance_Invoke_AsyncTaskFailed tests async task failure
func TestMongodbRestoreDbInstance_Invoke_AsyncTaskFailed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseMongodbClient", mongodbClient)

	flowId := int64(99999)
	patches.ApplyMethodFunc(mongodbClient, "RestoreDBInstanceWithContext", func(_ context.Context, request *mongodb_sdk.RestoreDBInstanceRequest) (*mongodb_sdk.RestoreDBInstanceResponse, error) {
		resp := mongodb_sdk.NewRestoreDBInstanceResponse()
		resp.Response = &mongodb_sdk.RestoreDBInstanceResponseParams{
			FlowId: &flowId,
		}
		return resp, nil
	})

	failedStatus := "failed"
	patches.ApplyMethodFunc(mongodbClient, "DescribeAsyncRequestInfo", func(request *mongodb_sdk.DescribeAsyncRequestInfoRequest) (*mongodb_sdk.DescribeAsyncRequestInfoResponse, error) {
		resp := mongodb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &mongodb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status: &failedStatus,
		}
		return resp, nil
	})

	aImpl := mongodb.NewMongodbRestoreDbInstance()
	a := aImpl.(*mongodb.MongodbRestoreDbInstance)
	setMongodbRestoreDbInstanceClient(a, client)

	req := newMongodbRestoreDbInstanceInvokeRequest(t,
		"cmgo-xxxxxxxx",
		"2024-09-01 12:00:00",
		[]tftypes.Value{
			newDatabaseValue("db1", []tftypes.Value{
				newCollectionValue("col_old", "col_new"),
			}),
		},
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "waiting for mongodb restore_db_instance task")
}

// TestMongodbRestoreDbInstance_Invoke_MissingInput tests missing required inputs
// are rejected before any API call
func TestMongodbRestoreDbInstance_Invoke_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseMongodbClient", mongodbClient)

	apiCalled := false
	patches.ApplyMethodFunc(mongodbClient, "RestoreDBInstanceWithContext", func(_ context.Context, request *mongodb_sdk.RestoreDBInstanceRequest) (*mongodb_sdk.RestoreDBInstanceResponse, error) {
		apiCalled = true
		return mongodb_sdk.NewRestoreDBInstanceResponse(), nil
	})

	aImpl := mongodb.NewMongodbRestoreDbInstance()
	a := aImpl.(*mongodb.MongodbRestoreDbInstance)
	setMongodbRestoreDbInstanceClient(a, client)

	// empty instance_id
	req := newMongodbRestoreDbInstanceInvokeRequest(t,
		"",
		"2024-09-01 12:00:00",
		[]tftypes.Value{
			newDatabaseValue("db1", []tftypes.Value{
				newCollectionValue("col_old", "col_new"),
			}),
		},
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing instance_id")
	assert.False(t, apiCalled)

	// empty restore_time
	req = newMongodbRestoreDbInstanceInvokeRequest(t,
		"cmgo-xxxxxxxx",
		"",
		[]tftypes.Value{
			newDatabaseValue("db1", []tftypes.Value{
				newCollectionValue("col_old", "col_new"),
			}),
		},
	)
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing restore_time")
	assert.False(t, apiCalled)

	// null databases
	req = newMongodbRestoreDbInstanceInvokeRequest(t,
		"cmgo-xxxxxxxx",
		"2024-09-01 12:00:00",
		nil,
	)
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing databases")
	assert.False(t, apiCalled)
}

// TestMongodbRestoreDbInstance_Invoke_ClientNotConfigured tests the nil-client guard
func TestMongodbRestoreDbInstance_Invoke_ClientNotConfigured(t *testing.T) {
	aImpl := mongodb.NewMongodbRestoreDbInstance()
	a := aImpl.(*mongodb.MongodbRestoreDbInstance)

	req := newMongodbRestoreDbInstanceInvokeRequest(t,
		"cmgo-xxxxxxxx",
		"2024-09-01 12:00:00",
		[]tftypes.Value{
			newDatabaseValue("db1", []tftypes.Value{
				newCollectionValue("col_old", "col_new"),
			}),
		},
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestMongodbRestoreDbInstance_MetadataAndSchema validates metadata and schema definition
func TestMongodbRestoreDbInstance_MetadataAndSchema(t *testing.T) {
	aImpl := mongodb.NewMongodbRestoreDbInstance()
	a := aImpl.(*mongodb.MongodbRestoreDbInstance)

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_mongodb_restore_db_instance", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "instance_id")
	assert.Contains(t, attrs, "restore_time")
	assert.True(t, attrs["instance_id"].IsRequired())
	assert.True(t, attrs["restore_time"].IsRequired())

	instanceIdAttr, ok := attrs["instance_id"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, instanceIdAttr.Required)

	blocks := schemaResp.Schema.Blocks
	assert.Contains(t, blocks, "databases")

	databasesBlock, ok := blocks["databases"].(action_schema.ListNestedBlock)
	assert.True(t, ok)

	dbAttrs := databasesBlock.NestedObject.Attributes
	assert.Contains(t, dbAttrs, "db")
	dbAttr, ok := dbAttrs["db"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, dbAttr.Required)

	dbBlocks := databasesBlock.NestedObject.Blocks
	assert.Contains(t, dbBlocks, "collections")

	collectionsBlock, ok := dbBlocks["collections"].(action_schema.ListNestedBlock)
	assert.True(t, ok)

	colAttrs := collectionsBlock.NestedObject.Attributes
	assert.Contains(t, colAttrs, "old_collection")
	assert.Contains(t, colAttrs, "new_collection")
	oldColAttr, ok := colAttrs["old_collection"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, oldColAttr.Required)
	newColAttr, ok := colAttrs["new_collection"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, newColAttr.Required)

	assert.Contains(t, blocks, "timeouts")
	timeoutsBlock, ok := blocks["timeouts"].(action_schema.SingleNestedBlock)
	assert.True(t, ok)
	// The block is provided by the official
	// terraform-plugin-framework-timeouts module, so it carries the custom
	// timeouts type plus a duration validator on `invoke`.
	assert.IsType(t, actiontimeouts.Type{}, timeoutsBlock.CustomType)
	invokeAttr, ok := timeoutsBlock.Attributes["invoke"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, invokeAttr.Optional)
	assert.NotEmpty(t, invokeAttr.Validators)
}

func TestMongodbRestoreDbInstance_Invoke_MissingNestedInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseMongodbClient", mongodbClient)

	apiCalled := false
	patches.ApplyMethodFunc(mongodbClient, "RestoreDBInstanceWithContext", func(_ context.Context, request *mongodb_sdk.RestoreDBInstanceRequest) (*mongodb_sdk.RestoreDBInstanceResponse, error) {
		apiCalled = true
		return mongodb_sdk.NewRestoreDBInstanceResponse(), nil
	})

	aImpl := mongodb.NewMongodbRestoreDbInstance()
	a := aImpl.(*mongodb.MongodbRestoreDbInstance)
	setMongodbRestoreDbInstanceClient(a, client)

	// databases block with empty collections list
	dbWithNoCols := tftypes.NewValue(databaseObjType(), map[string]tftypes.Value{
		"db":          tftypes.NewValue(tftypes.String, "db1"),
		"collections": tftypes.NewValue(tftypes.List{ElementType: collectionObjType()}, nil),
	})
	req := newMongodbRestoreDbInstanceInvokeRequest(t,
		"cmgo-xxxxxxxx",
		"2024-09-01 12:00:00",
		[]tftypes.Value{dbWithNoCols},
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing collections")
	assert.False(t, apiCalled)
}

// TestMongodbRestoreDbInstance_Invoke_InvalidTimeout verifies that a malformed
// timeouts.invoke value is rejected before the restore API is called. The
// official terraform-plugin-framework-timeouts Value.Invoke surfaces the parse
// failure through a diagnostic.
func TestMongodbRestoreDbInstance_Invoke_InvalidTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseMongodbClient", mongodbClient)

	apiCalled := false
	patches.ApplyMethodFunc(mongodbClient, "RestoreDBInstanceWithContext", func(_ context.Context, request *mongodb_sdk.RestoreDBInstanceRequest) (*mongodb_sdk.RestoreDBInstanceResponse, error) {
		apiCalled = true
		return mongodb_sdk.NewRestoreDBInstanceResponse(), nil
	})

	aImpl := mongodb.NewMongodbRestoreDbInstance()
	a := aImpl.(*mongodb.MongodbRestoreDbInstance)
	setMongodbRestoreDbInstanceClient(a, client)

	req := newMongodbRestoreDbInstanceInvokeRequestWithTimeouts(t,
		"cmgo-xxxxxxxx",
		"2024-09-01 12:00:00",
		[]tftypes.Value{
			newDatabaseValue("db1", []tftypes.Value{
				newCollectionValue("col_old", "col_new"),
			}),
		},
		tftypes.NewValue(timeoutsObjType(), map[string]tftypes.Value{
			"invoke": tftypes.NewValue(tftypes.String, "not-a-duration"),
		}),
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Timeout Cannot Be Parsed")
	assert.False(t, apiCalled)
}

// TestMongodbRestoreDbInstance_Invoke_CustomTimeout verifies that a configured
// timeouts.invoke value really bounds the async polling: the mocked task stays
// in "running" forever, so the invoke can only return (quickly) if the 1ms
// custom timeout - rather than the 15m default - was applied.
func TestMongodbRestoreDbInstance_Invoke_CustomTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseMongodbClient", mongodbClient)

	flowId := int64(77777)
	patches.ApplyMethodFunc(mongodbClient, "RestoreDBInstanceWithContext", func(_ context.Context, request *mongodb_sdk.RestoreDBInstanceRequest) (*mongodb_sdk.RestoreDBInstanceResponse, error) {
		resp := mongodb_sdk.NewRestoreDBInstanceResponse()
		resp.Response = &mongodb_sdk.RestoreDBInstanceResponseParams{
			FlowId: &flowId,
		}
		return resp, nil
	})

	runningStatus := "running"
	patches.ApplyMethodFunc(mongodbClient, "DescribeAsyncRequestInfo", func(request *mongodb_sdk.DescribeAsyncRequestInfoRequest) (*mongodb_sdk.DescribeAsyncRequestInfoResponse, error) {
		resp := mongodb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &mongodb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status: &runningStatus,
		}
		return resp, nil
	})

	aImpl := mongodb.NewMongodbRestoreDbInstance()
	a := aImpl.(*mongodb.MongodbRestoreDbInstance)
	setMongodbRestoreDbInstanceClient(a, client)

	req := newMongodbRestoreDbInstanceInvokeRequestWithTimeouts(t,
		"cmgo-xxxxxxxx",
		"2024-09-01 12:00:00",
		[]tftypes.Value{
			newDatabaseValue("db1", []tftypes.Value{
				newCollectionValue("col_old", "col_new"),
			}),
		},
		tftypes.NewValue(timeoutsObjType(), map[string]tftypes.Value{
			"invoke": tftypes.NewValue(tftypes.String, "1ms"),
		}),
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "waiting for mongodb restore_db_instance task")
}
