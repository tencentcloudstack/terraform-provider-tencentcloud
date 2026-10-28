package teo_test

import (
	"context"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	teosdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStringTeoApplyFreeCertificate(s string) *string {
	return &s
}

// teoApplyFreeCertificateSchema returns the framework action schema used to
// build the raw config value for InvokeRequest.
func teoApplyFreeCertificateSchema(t *testing.T) action_schema.Schema {
	a := teo.NewTeoApplyFreeCertificate()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newTeoApplyFreeCertificateInvokeRequest builds an action.InvokeRequest
// carrying the given zone_id, domain, and verification_method as raw tftypes values.
func newTeoApplyFreeCertificateInvokeRequest(t *testing.T, zoneId, domain, verificationMethod string) action.InvokeRequest {
	s := teoApplyFreeCertificateSchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":             tftypes.String,
			"domain":              tftypes.String,
			"verification_method": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"zone_id":             tftypes.NewValue(tftypes.String, zoneId),
		"domain":              tftypes.NewValue(tftypes.String, domain),
		"verification_method": tftypes.NewValue(tftypes.String, verificationMethod),
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// newTeoApplyFreeCertificateInvokeRequestNullField builds an InvokeRequest
// where the specified field is null.
func newTeoApplyFreeCertificateInvokeRequestNullField(t *testing.T, nullField string) action.InvokeRequest {
	s := teoApplyFreeCertificateSchema(t)

	zoneIdVal := tftypes.NewValue(tftypes.String, "zone-123")
	domainVal := tftypes.NewValue(tftypes.String, "example.com")
	methodVal := tftypes.NewValue(tftypes.String, "dns_challenge")

	switch nullField {
	case "zone_id":
		zoneIdVal = tftypes.NewValue(tftypes.String, nil)
	case "domain":
		domainVal = tftypes.NewValue(tftypes.String, nil)
	case "verification_method":
		methodVal = tftypes.NewValue(tftypes.String, nil)
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":             tftypes.String,
			"domain":              tftypes.String,
			"verification_method": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"zone_id":             zoneIdVal,
		"domain":              domainVal,
		"verification_method": methodVal,
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setTeoApplyFreeCertificateClient injects a mock connectivity client into
// the action instance via the Configure method.
func setTeoApplyFreeCertificateClient(a *teo.TeoApplyFreeCertificate, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoApplyFreeCertificate" -v -count=1 -gcflags="all=-l"

// TestTeoApplyFreeCertificate_Invoke_Success tests a successful invoke
func TestTeoApplyFreeCertificate_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotRequest *teosdk.ApplyFreeCertificateRequest
	patches.ApplyMethodFunc(teoClient, "ApplyFreeCertificateWithContext", func(_ context.Context, request *teosdk.ApplyFreeCertificateRequest) (*teosdk.ApplyFreeCertificateResponse, error) {
		gotRequest = request
		resp := teosdk.NewApplyFreeCertificateResponse()
		resp.Response = &teosdk.ApplyFreeCertificateResponseParams{
			DnsVerification: &teosdk.DnsVerification{
				Subdomain:   ptrStringTeoApplyFreeCertificate("_acme-challenge"),
				RecordType:  ptrStringTeoApplyFreeCertificate("CNAME"),
				RecordValue: ptrStringTeoApplyFreeCertificate("xxx.teo-validation.com"),
			},
			RequestId: ptrStringTeoApplyFreeCertificate("fake-request-id"),
		}
		return resp, nil
	})

	commonName := "example.com"
	patches.ApplyMethodFunc(teoClient, "CheckFreeCertificateVerificationWithContext", func(_ context.Context, request *teosdk.CheckFreeCertificateVerificationRequest) (*teosdk.CheckFreeCertificateVerificationResponse, error) {
		resp := teosdk.NewCheckFreeCertificateVerificationResponse()
		resp.Response = &teosdk.CheckFreeCertificateVerificationResponseParams{
			CommonName:         &commonName,
			SignatureAlgorithm: ptrStringTeoApplyFreeCertificate("RSA 2048"),
			ExpireTime:         ptrStringTeoApplyFreeCertificate("2025-12-31T23:59:59Z"),
		}
		return resp, nil
	})

	aImpl := teo.NewTeoApplyFreeCertificate()
	a := aImpl.(*teo.TeoApplyFreeCertificate)
	setTeoApplyFreeCertificateClient(a, client)

	req := newTeoApplyFreeCertificateInvokeRequest(t, "zone-123", "example.com", "dns_challenge")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, "zone-123", *gotRequest.ZoneId)
	assert.Equal(t, "example.com", *gotRequest.Domain)
	assert.Equal(t, "dns_challenge", *gotRequest.VerificationMethod)
}

// TestTeoApplyFreeCertificate_Invoke_ApplyError tests API error handling
func TestTeoApplyFreeCertificate_Invoke_ApplyError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ApplyFreeCertificateWithContext", func(_ context.Context, request *teosdk.ApplyFreeCertificateRequest) (*teosdk.ApplyFreeCertificateResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{
			Code:    "LimitExceeded",
			Message: "apply quota exceeded",
		}
	})

	patches.ApplyMethodFunc(teoClient, "CheckFreeCertificateVerificationWithContext", func(_ context.Context, request *teosdk.CheckFreeCertificateVerificationRequest) (*teosdk.CheckFreeCertificateVerificationResponse, error) {
		return teosdk.NewCheckFreeCertificateVerificationResponse(), nil
	})

	aImpl := teo.NewTeoApplyFreeCertificate()
	a := aImpl.(*teo.TeoApplyFreeCertificate)
	setTeoApplyFreeCertificateClient(a, client)

	req := newTeoApplyFreeCertificateInvokeRequest(t, "zone-123", "example.com", "dns_challenge")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "applying TEO free certificate")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "LimitExceeded")
}

// TestTeoApplyFreeCertificate_Invoke_VerificationTimeout tests verification
// polling timeout when CheckFreeCertificateVerification always returns pending
func TestTeoApplyFreeCertificate_Invoke_VerificationTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch ReadRetryTimeout to a very small value so that the polling timeout
	// (3 * ReadRetryTimeout) finishes quickly in the test.
	patches.ApplyGlobalVar(&tccommon.ReadRetryTimeout, 1*time.Second)

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ApplyFreeCertificateWithContext", func(_ context.Context, request *teosdk.ApplyFreeCertificateRequest) (*teosdk.ApplyFreeCertificateResponse, error) {
		resp := teosdk.NewApplyFreeCertificateResponse()
		resp.Response = &teosdk.ApplyFreeCertificateResponseParams{
			DnsVerification: &teosdk.DnsVerification{
				Subdomain:   ptrStringTeoApplyFreeCertificate("_acme-challenge"),
				RecordType:  ptrStringTeoApplyFreeCertificate("CNAME"),
				RecordValue: ptrStringTeoApplyFreeCertificate("xxx.teo-validation.com"),
			},
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "CheckFreeCertificateVerificationWithContext", func(_ context.Context, request *teosdk.CheckFreeCertificateVerificationRequest) (*teosdk.CheckFreeCertificateVerificationResponse, error) {
		// Always return empty CommonName (verification still pending)
		resp := teosdk.NewCheckFreeCertificateVerificationResponse()
		resp.Response = &teosdk.CheckFreeCertificateVerificationResponseParams{}
		return resp, nil
	})

	aImpl := teo.NewTeoApplyFreeCertificate()
	a := aImpl.(*teo.TeoApplyFreeCertificate)
	setTeoApplyFreeCertificateClient(a, client)

	req := newTeoApplyFreeCertificateInvokeRequest(t, "zone-123", "example.com", "dns_challenge")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "checking TEO free certificate verification")
}

// TestTeoApplyFreeCertificate_Invoke_MissingInput tests missing required inputs
// are rejected before any API call
func TestTeoApplyFreeCertificate_Invoke_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	apiCalled := false
	patches.ApplyMethodFunc(teoClient, "ApplyFreeCertificateWithContext", func(_ context.Context, request *teosdk.ApplyFreeCertificateRequest) (*teosdk.ApplyFreeCertificateResponse, error) {
		apiCalled = true
		return teosdk.NewApplyFreeCertificateResponse(), nil
	})
	patches.ApplyMethodFunc(teoClient, "CheckFreeCertificateVerificationWithContext", func(_ context.Context, request *teosdk.CheckFreeCertificateVerificationRequest) (*teosdk.CheckFreeCertificateVerificationResponse, error) {
		resp := teosdk.NewCheckFreeCertificateVerificationResponse()
		resp.Response = &teosdk.CheckFreeCertificateVerificationResponseParams{
			CommonName: ptrStringTeoApplyFreeCertificate("example.com"),
		}
		return resp, nil
	})

	aImpl := teo.NewTeoApplyFreeCertificate()

	testCases := []struct {
		name      string
		nullField string
	}{
		{"null zone_id", "zone_id"},
		{"null domain", "domain"},
		{"null verification_method", "verification_method"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiCalled = false
			a := aImpl.(*teo.TeoApplyFreeCertificate)
			setTeoApplyFreeCertificateClient(a, client)

			req := newTeoApplyFreeCertificateInvokeRequestNullField(t, tc.nullField)
			resp := &action.InvokeResponse{}
			a.Invoke(context.Background(), req, resp)

			assert.True(t, resp.Diagnostics.HasError())
			assert.False(t, apiCalled, "API should not be called when required input is missing")
		})
	}
}

// TestTeoApplyFreeCertificate_Invoke_ClientNotConfigured tests the nil-client guard
func TestTeoApplyFreeCertificate_Invoke_ClientNotConfigured(t *testing.T) {
	aImpl := teo.NewTeoApplyFreeCertificate()
	a := aImpl.(*teo.TeoApplyFreeCertificate)

	req := newTeoApplyFreeCertificateInvokeRequest(t, "zone-123", "example.com", "dns_challenge")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestTeoApplyFreeCertificate_MetadataAndSchema validates metadata and schema definition
func TestTeoApplyFreeCertificate_MetadataAndSchema(t *testing.T) {
	aImpl := teo.NewTeoApplyFreeCertificate()
	a := aImpl.(*teo.TeoApplyFreeCertificate)

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_apply_free_certificate", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "zone_id")
	assert.Contains(t, attrs, "domain")
	assert.Contains(t, attrs, "verification_method")

	assert.True(t, attrs["zone_id"].IsRequired())
	assert.True(t, attrs["domain"].IsRequired())
	assert.True(t, attrs["verification_method"].IsRequired())

	_, ok := attrs["zone_id"].(action_schema.StringAttribute)
	assert.True(t, ok, "zone_id should be a StringAttribute")
	_, ok = attrs["domain"].(action_schema.StringAttribute)
	assert.True(t, ok, "domain should be a StringAttribute")
	_, ok = attrs["verification_method"].(action_schema.StringAttribute)
	assert.True(t, ok, "verification_method should be a StringAttribute")
}
