package cloudaudit_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcaudit "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cloudaudit"
)

func ptrStr(s string) *string { return &s }
func ptrI64(i int64) *int64   { return &i }

func cloudauditAuditSchema(t *testing.T) schema.Schema {
	d := svcaudit.NewCloudauditAuditDataSource()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

func newCloudauditAuditReadRequest(t *testing.T, auditName string) datasource.ReadRequest {
	s := cloudauditAuditSchema(t)

	attrTypes := map[string]tftypes.Type{
		"audit_name":           tftypes.String,
		"is_enable_cmq_notify": tftypes.Number,
		"read_write_attribute": tftypes.Number,
		"key_id":               tftypes.String,
		"audit_status":         tftypes.Number,
		"cos_region":           tftypes.String,
		"cmq_queue_name":       tftypes.String,
		"kms_alias":            tftypes.String,
		"kms_region":           tftypes.String,
		"is_enable_kms_encry":  tftypes.Number,
		"cos_bucket_name":      tftypes.String,
		"cmq_region":           tftypes.String,
		"log_file_prefix":      tftypes.String,
	}
	rawValues := map[string]tftypes.Value{
		"audit_name":           tftypes.NewValue(tftypes.String, auditName),
		"is_enable_cmq_notify": tftypes.NewValue(tftypes.Number, nil),
		"read_write_attribute": tftypes.NewValue(tftypes.Number, nil),
		"key_id":               tftypes.NewValue(tftypes.String, nil),
		"audit_status":         tftypes.NewValue(tftypes.Number, nil),
		"cos_region":           tftypes.NewValue(tftypes.String, nil),
		"cmq_queue_name":       tftypes.NewValue(tftypes.String, nil),
		"kms_alias":            tftypes.NewValue(tftypes.String, nil),
		"kms_region":           tftypes.NewValue(tftypes.String, nil),
		"is_enable_kms_encry":  tftypes.NewValue(tftypes.Number, nil),
		"cos_bucket_name":      tftypes.NewValue(tftypes.String, nil),
		"cmq_region":           tftypes.NewValue(tftypes.String, nil),
		"log_file_prefix":      tftypes.NewValue(tftypes.String, nil),
	}
	raw := tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, rawValues)

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

func setCloudauditAuditClient(d *svcaudit.CloudauditAuditDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/cloudaudit/ -run "TestCloudauditAudit" -v -count=1 -gcflags="all=-l"

// TestCloudauditAuditDataSource_Metadata validates Metadata.TypeName.
func TestCloudauditAuditDataSource_Metadata(t *testing.T) {
	d := svcaudit.NewCloudauditAuditDataSource()
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	assert.Equal(t, "tencentcloud_cloudaudit_audit", resp.TypeName)
}

// TestCloudauditAuditDataSource_Factory validates the registry factory.
func TestCloudauditAuditDataSource_Factory(t *testing.T) {
	got := svcaudit.NewCloudauditAuditDataSource()
	assert.NotNil(t, got)
}

// TestCloudauditAuditDataSource_ImplementsInterface is a compile-time guard.
func TestCloudauditAuditDataSource_ImplementsInterface(t *testing.T) {
	var _ datasource.DataSourceWithConfigure = &svcaudit.CloudauditAuditDataSource{}
}

// TestCloudauditAuditDataSource_Read_Success tests a successful read that
// writes all leaf fields into state and reports zero error diagnostics.
func TestCloudauditAuditDataSource_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	auditClient := &audit.Client{}
	patches.ApplyMethodReturn(client, "UseAuditClient", auditClient)

	var gotRequest *audit.DescribeAuditRequest
	patches.ApplyMethodFunc(auditClient, "DescribeAuditWithContext", func(_ context.Context, request *audit.DescribeAuditRequest) (*audit.DescribeAuditResponse, error) {
		gotRequest = request
		resp := &audit.DescribeAuditResponse{
			Response: &audit.DescribeAuditResponseParams{
				IsEnableCmqNotify:  ptrI64(1),
				ReadWriteAttribute: ptrI64(3),
				KeyId:              ptrStr("key-xxxx"),
				AuditStatus:        ptrI64(1),
				AuditName:          ptrStr("my-audit"),
				CosRegion:          ptrStr("ap-guangzhou"),
				CmqQueueName:       ptrStr("cmq-queue"),
				KmsAlias:           ptrStr("kms-alias"),
				KmsRegion:          ptrStr("ap-guangzhou"),
				IsEnableKmsEncry:   ptrI64(0),
				CosBucketName:      ptrStr("bucket-1234567890"),
				CmqRegion:          ptrStr("ap-guangzhou"),
				LogFilePrefix:      ptrStr("logs"),
			},
		}
		return resp, nil
	})

	d := svcaudit.NewCloudauditAuditDataSource().(*svcaudit.CloudauditAuditDataSource)
	setCloudauditAuditClient(d, client)

	req := newCloudauditAuditReadRequest(t, "my-audit")
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotRequest)
	assert.Equal(t, "my-audit", *gotRequest.AuditName)

	var state svcaudit.CloudauditAuditDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &state)...)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error reading state")

	assert.Equal(t, "my-audit", state.AuditName.ValueString())
	assert.Equal(t, int64(1), state.IsEnableCmqNotify.ValueInt64())
	assert.Equal(t, int64(3), state.ReadWriteAttribute.ValueInt64())
	assert.Equal(t, "key-xxxx", state.KeyId.ValueString())
	assert.Equal(t, int64(1), state.AuditStatus.ValueInt64())
	assert.Equal(t, "ap-guangzhou", state.CosRegion.ValueString())
	assert.Equal(t, "cmq-queue", state.CmqQueueName.ValueString())
	assert.Equal(t, "kms-alias", state.KmsAlias.ValueString())
	assert.Equal(t, "ap-guangzhou", state.KmsRegion.ValueString())
	assert.Equal(t, int64(0), state.IsEnableKmsEncry.ValueInt64())
	assert.Equal(t, "bucket-1234567890", state.CosBucketName.ValueString())
	assert.Equal(t, "ap-guangzhou", state.CmqRegion.ValueString())
	assert.Equal(t, "logs", state.LogFilePrefix.ValueString())
}

// TestCloudauditAuditDataSource_Read_APIError tests API error handling.
func TestCloudauditAuditDataSource_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	auditClient := &audit.Client{}
	patches.ApplyMethodReturn(client, "UseAuditClient", auditClient)
	patches.ApplyMethodFunc(auditClient, "DescribeAuditWithContext", func(_ context.Context, request *audit.DescribeAuditRequest) (*audit.DescribeAuditResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError.DescribeAuditError, Message=internal error")
	})

	d := svcaudit.NewCloudauditAuditDataSource().(*svcaudit.CloudauditAuditDataSource)
	setCloudauditAuditClient(d, client)

	req := newCloudauditAuditReadRequest(t, "my-audit")
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "DescribeAuditError")
}

// TestCloudauditAuditDataSource_Read_NotFound tests that a not-found error is
// surfaced as an explicit diagnostic rather than a silent empty state.
func TestCloudauditAuditDataSource_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	auditClient := &audit.Client{}
	patches.ApplyMethodReturn(client, "UseAuditClient", auditClient)
	patches.ApplyMethodFunc(auditClient, "DescribeAuditWithContext", func(_ context.Context, request *audit.DescribeAuditRequest) (*audit.DescribeAuditResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound.AuditNotExist, Message=audit not exist")
	})

	d := svcaudit.NewCloudauditAuditDataSource().(*svcaudit.CloudauditAuditDataSource)
	setCloudauditAuditClient(d, client)

	req := newCloudauditAuditReadRequest(t, "missing-audit")
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "not found")
}

// TestCloudauditAuditDataSource_Read_EmptyResponse tests that a nil response
// is treated as an API failure, not an empty result.
func TestCloudauditAuditDataSource_Read_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	auditClient := &audit.Client{}
	patches.ApplyMethodReturn(client, "UseAuditClient", auditClient)
	patches.ApplyMethodFunc(auditClient, "DescribeAuditWithContext", func(_ context.Context, request *audit.DescribeAuditRequest) (*audit.DescribeAuditResponse, error) {
		return &audit.DescribeAuditResponse{}, nil
	})

	d := svcaudit.NewCloudauditAuditDataSource().(*svcaudit.CloudauditAuditDataSource)
	setCloudauditAuditClient(d, client)

	req := newCloudauditAuditReadRequest(t, "my-audit")
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
}

// TestCloudauditAuditDataSource_Read_MissingInput tests that an empty
// audit_name is rejected before any API call.
func TestCloudauditAuditDataSource_Read_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	auditClient := &audit.Client{}
	patches.ApplyMethodReturn(client, "UseAuditClient", auditClient)
	apiCalled := false
	patches.ApplyMethodFunc(auditClient, "DescribeAuditWithContext", func(_ context.Context, request *audit.DescribeAuditRequest) (*audit.DescribeAuditResponse, error) {
		apiCalled = true
		return &audit.DescribeAuditResponse{}, nil
	})

	d := svcaudit.NewCloudauditAuditDataSource().(*svcaudit.CloudauditAuditDataSource)
	setCloudauditAuditClient(d, client)

	req := newCloudauditAuditReadRequest(t, "")
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing audit_name")
	assert.False(t, apiCalled)
}

// TestCloudauditAuditDataSource_Read_ClientNotConfigured tests the nil-client
// guard when configure has not run.
func TestCloudauditAuditDataSource_Read_ClientNotConfigured(t *testing.T) {
	d := svcaudit.NewCloudauditAuditDataSource().(*svcaudit.CloudauditAuditDataSource)

	req := newCloudauditAuditReadRequest(t, "my-audit")
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}
