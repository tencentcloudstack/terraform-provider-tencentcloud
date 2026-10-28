package kms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	kms_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/kms"
)

func ptrStringKmsInstanceOp(s string) *string {
	return &s
}

func kmsInstanceSchema(t *testing.T) action_schema.Schema {
	a := kms.NewKmsInstance()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// prepaidObjType builds the tftypes.Object type for the instance_charge_prepaid block.
func prepaidObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"period":     tftypes.Number,
			"renew_flag": tftypes.Number,
		},
	}
}

// tagObjType builds the tftypes.Object type for a tags block element.
func tagObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"tag_key":   tftypes.String,
			"tag_value": tftypes.String,
		},
	}
}

// newPrepaidValue builds a tftypes.Value for the instance_charge_prepaid block.
func newPrepaidValue(period, renewFlag int64) tftypes.Value {
	return tftypes.NewValue(prepaidObjType(), map[string]tftypes.Value{
		"period":     tftypes.NewValue(tftypes.Number, period),
		"renew_flag": tftypes.NewValue(tftypes.Number, renewFlag),
	})
}

// nullPrepaidValue is the null value for an unconfigured instance_charge_prepaid block.
func nullPrepaidValue() tftypes.Value {
	return tftypes.NewValue(prepaidObjType(), nil)
}

// newTagValue builds a tftypes.Value for one tags element.
func newTagValue(tagKey, tagValue string) tftypes.Value {
	return tftypes.NewValue(tagObjType(), map[string]tftypes.Value{
		"tag_key":   tftypes.NewValue(tftypes.String, tagKey),
		"tag_value": tftypes.NewValue(tftypes.String, tagValue),
	})
}

// newKmsInstanceInvokeRequest builds an action.InvokeRequest carrying
// the given raw config values.
func newKmsInstanceInvokeRequest(
	t *testing.T,
	instanceChargeType, productVersion string,
	prepaidVal tftypes.Value,
	tags []tftypes.Value,
) action.InvokeRequest {
	s := kmsInstanceSchema(t)

	tagsListType := tftypes.List{ElementType: tagObjType()}
	var tagsVal tftypes.Value
	if tags == nil {
		tagsVal = tftypes.NewValue(tagsListType, nil)
	} else {
		tagsVal = tftypes.NewValue(tagsListType, tags)
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"instance_charge_type":    tftypes.String,
			"product_version":         tftypes.String,
			"instance_charge_prepaid": prepaidObjType(),
			"tags":                    tagsListType,
		},
	}, map[string]tftypes.Value{
		"instance_charge_type":    tftypes.NewValue(tftypes.String, instanceChargeType),
		"product_version":         tftypes.NewValue(tftypes.String, productVersion),
		"instance_charge_prepaid": prepaidVal,
		"tags":                    tagsVal,
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setKmsInstanceClient injects a mock connectivity client into the
// action instance via the Configure method.
func setKmsInstanceClient(a *kms.KmsInstance, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/kms/ -run "TestKmsInstance" -v -count=1 -gcflags="all=-l"

// TestKmsInstance_Invoke_Success tests a successful invoke
func TestKmsInstance_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	kmsClient := &kms_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseKmsClient", kmsClient)

	var gotRequest *kms_sdk.CreateKmsInstanceRequest
	patches.ApplyMethodFunc(kmsClient, "CreateKmsInstanceWithContext", func(_ context.Context, request *kms_sdk.CreateKmsInstanceRequest) (*kms_sdk.CreateKmsInstanceResponse, error) {
		gotRequest = request
		assert.Equal(t, "PREPAID", *request.InstanceChargeType)
		assert.Equal(t, "professional", *request.ProductVersion)
		assert.NotNil(t, request.ClientToken)
		assert.NotEmpty(t, *request.ClientToken)
		assert.NotNil(t, request.InstanceChargePrepaid)
		assert.Equal(t, uint64(1), *request.InstanceChargePrepaid.Period)
		assert.Equal(t, uint64(1), *request.InstanceChargePrepaid.RenewFlag)
		assert.Len(t, request.Tags, 1)
		assert.Equal(t, "env", *request.Tags[0].TagKey)
		assert.Equal(t, "prod", *request.Tags[0].TagValue)

		resp := kms_sdk.NewCreateKmsInstanceResponse()
		resp.Response = &kms_sdk.CreateKmsInstanceResponseParams{
			InstanceId: ptrStringKmsInstanceOp("kms-abc123"),
			DealName:   ptrStringKmsInstanceOp("20261028110001"),
			BillId:     ptrStringKmsInstanceOp("bill-001"),
			RequestId:  ptrStringKmsInstanceOp("fake-request-id"),
		}
		return resp, nil
	})

	a := kms.NewKmsInstance().(*kms.KmsInstance)
	setKmsInstanceClient(a, client)

	req := newKmsInstanceInvokeRequest(t, "PREPAID", "professional",
		newPrepaidValue(1, 1),
		[]tftypes.Value{newTagValue("env", "prod")},
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotRequest)
}

// TestKmsInstance_Invoke_MinimalInput tests that optional fields are left
// unset (nil) on the request when they are not configured.
func TestKmsInstance_Invoke_MinimalInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	kmsClient := &kms_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseKmsClient", kmsClient)

	var gotRequest *kms_sdk.CreateKmsInstanceRequest
	patches.ApplyMethodFunc(kmsClient, "CreateKmsInstanceWithContext", func(_ context.Context, request *kms_sdk.CreateKmsInstanceRequest) (*kms_sdk.CreateKmsInstanceResponse, error) {
		gotRequest = request
		resp := kms_sdk.NewCreateKmsInstanceResponse()
		resp.Response = &kms_sdk.CreateKmsInstanceResponseParams{
			InstanceId: ptrStringKmsInstanceOp("kms-min001"),
		}
		return resp, nil
	})

	a := kms.NewKmsInstance().(*kms.KmsInstance)
	setKmsInstanceClient(a, client)

	req := newKmsInstanceInvokeRequest(t, "PREPAID", "",
		nullPrepaidValue(),
		nil,
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotRequest)
	assert.Equal(t, "PREPAID", *gotRequest.InstanceChargeType)
	assert.Nil(t, gotRequest.ProductVersion)
	assert.NotNil(t, gotRequest.ClientToken)
	assert.NotEmpty(t, *gotRequest.ClientToken)
	assert.Nil(t, gotRequest.InstanceChargePrepaid)
	assert.Empty(t, gotRequest.Tags)
}

// TestKmsInstance_Invoke_APIError tests API error handling
func TestKmsInstance_Invoke_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	kmsClient := &kms_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseKmsClient", kmsClient)

	patches.ApplyMethodFunc(kmsClient, "CreateKmsInstanceWithContext", func(_ context.Context, request *kms_sdk.CreateKmsInstanceRequest) (*kms_sdk.CreateKmsInstanceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=LimitExceeded, Message=instance quota exceeded")
	})

	a := kms.NewKmsInstance().(*kms.KmsInstance)
	setKmsInstanceClient(a, client)

	req := newKmsInstanceInvokeRequest(t, "PREPAID", "",
		nullPrepaidValue(),
		nil,
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "creating kms_instance")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "LimitExceeded")
}

// TestKmsInstance_Invoke_EmptyResponse tests that an empty response
// returns an error diagnostic.
func TestKmsInstance_Invoke_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	kmsClient := &kms_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseKmsClient", kmsClient)

	patches.ApplyMethodFunc(kmsClient, "CreateKmsInstanceWithContext", func(_ context.Context, request *kms_sdk.CreateKmsInstanceRequest) (*kms_sdk.CreateKmsInstanceResponse, error) {
		resp := kms_sdk.NewCreateKmsInstanceResponse()
		resp.Response = &kms_sdk.CreateKmsInstanceResponseParams{
			InstanceId: ptrStringKmsInstanceOp(""),
		}
		return resp, nil
	})

	a := kms.NewKmsInstance().(*kms.KmsInstance)
	setKmsInstanceClient(a, client)

	req := newKmsInstanceInvokeRequest(t, "PREPAID", "",
		nullPrepaidValue(),
		nil,
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "creating kms_instance")
}

// TestKmsInstance_Invoke_MissingInput tests missing required input
// is rejected before any API call
func TestKmsInstance_Invoke_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	kmsClient := &kms_sdk.Client{}
	patches.ApplyMethodReturn(client, "UseKmsClient", kmsClient)

	apiCalled := false
	patches.ApplyMethodFunc(kmsClient, "CreateKmsInstanceWithContext", func(_ context.Context, request *kms_sdk.CreateKmsInstanceRequest) (*kms_sdk.CreateKmsInstanceResponse, error) {
		apiCalled = true
		return kms_sdk.NewCreateKmsInstanceResponse(), nil
	})

	a := kms.NewKmsInstance().(*kms.KmsInstance)
	setKmsInstanceClient(a, client)

	req := newKmsInstanceInvokeRequest(t, "", "",
		nullPrepaidValue(),
		nil,
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing instance_charge_type")
	assert.False(t, apiCalled)
}

// TestKmsInstance_Invoke_ClientNotConfigured tests the nil-client guard
func TestKmsInstance_Invoke_ClientNotConfigured(t *testing.T) {
	a := kms.NewKmsInstance().(*kms.KmsInstance)

	req := newKmsInstanceInvokeRequest(t, "PREPAID", "",
		nullPrepaidValue(),
		nil,
	)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestKmsInstance_MetadataAndSchema validates metadata and schema definition
func TestKmsInstance_MetadataAndSchema(t *testing.T) {
	a := kms.NewKmsInstance().(*kms.KmsInstance)

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_kms_instance", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "instance_charge_type")
	assert.Contains(t, attrs, "product_version")
	assert.NotContains(t, attrs, "client_token")
	assert.True(t, attrs["instance_charge_type"].IsRequired())
	assert.True(t, attrs["product_version"].IsOptional())

	blocks := schemaResp.Schema.Blocks
	assert.Contains(t, blocks, "instance_charge_prepaid")
	assert.Contains(t, blocks, "tags")

	prepaidBlock, ok := blocks["instance_charge_prepaid"].(action_schema.SingleNestedBlock)
	assert.True(t, ok)
	assert.Contains(t, prepaidBlock.Attributes, "period")
	assert.Contains(t, prepaidBlock.Attributes, "renew_flag")
	assert.True(t, prepaidBlock.Attributes["period"].IsRequired())
	assert.True(t, prepaidBlock.Attributes["renew_flag"].IsRequired())

	tagsBlock, ok := blocks["tags"].(action_schema.ListNestedBlock)
	assert.True(t, ok)
	tagAttrs := tagsBlock.NestedObject.Attributes
	assert.Contains(t, tagAttrs, "tag_key")
	assert.Contains(t, tagAttrs, "tag_value")
	assert.True(t, tagAttrs["tag_key"].IsRequired())
	assert.True(t, tagAttrs["tag_value"].IsRequired())
}
