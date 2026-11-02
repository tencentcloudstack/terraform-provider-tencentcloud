package ssl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	sslSdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	providerSsl "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssl"
)

func ptrStringSslUploadCertInstance(s string) *string { return &s }
func ptrUint64SslUploadCertInstance(n uint64) *uint64 { return &n }
func ptrInt64SslUploadCertInstance(n int64) *int64    { return &n }

// sslUploadCertificateInstanceSchema returns the framework action schema used
// to build the raw config value for InvokeRequest.
func sslUploadCertificateInstanceSchema(t *testing.T) action_schema.Schema {
	a := providerSsl.NewSslUploadCertificateInstance()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// sslUploadTimeoutsObjType builds the tftypes.Object type for the timeouts
// block.
func sslUploadTimeoutsObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"invoke": tftypes.String,
		},
	}
}

// sslUploadNullTimeoutsValue is the null value for an unconfigured timeouts
// block.
func sslUploadNullTimeoutsValue() tftypes.Value {
	return tftypes.NewValue(sslUploadTimeoutsObjType(), nil)
}

// sslUploadRegionObjType builds the tftypes.Object type for a
// resource_types_regions block element.
func sslUploadRegionObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"resource_type": tftypes.String,
			"regions":       tftypes.Set{ElementType: tftypes.String},
		},
	}
}

// newSslUploadRegionValue builds a tftypes.Value for one
// resource_types_regions element.
func newSslUploadRegionValue(resourceType string, regions []string) tftypes.Value {
	regionVals := make([]tftypes.Value, 0, len(regions))
	for _, r := range regions {
		regionVals = append(regionVals, tftypes.NewValue(tftypes.String, r))
	}
	var regionsVal tftypes.Value
	if len(regionVals) == 0 {
		regionsVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	} else {
		regionsVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, regionVals)
	}
	var resourceTypeVal tftypes.Value
	if resourceType == "" {
		resourceTypeVal = tftypes.NewValue(tftypes.String, nil)
	} else {
		resourceTypeVal = tftypes.NewValue(tftypes.String, resourceType)
	}
	return tftypes.NewValue(sslUploadRegionObjType(), map[string]tftypes.Value{
		"resource_type": resourceTypeVal,
		"regions":       regionsVal,
	})
}

// newSslUploadCertificateInstanceInvokeRequest builds an InvokeRequest whose
// timeouts block is left unconfigured, so the action's built-in default (15m)
// applies.
func newSslUploadCertificateInstanceInvokeRequest(t *testing.T, oldCertificateId string, resourceTypes []string, publicKey, privateKey string) action.InvokeRequest {
	return newSslUploadCertificateInstanceInvokeRequestFull(t, oldCertificateId, resourceTypes, publicKey, privateKey, nil, nil, sslUploadNullTimeoutsValue())
}

// newSslUploadCertificateInstanceInvokeRequestFull builds an InvokeRequest
// carrying the given raw values, including the optional certificate_id, the
// resource_types_regions block and the timeouts block.
func newSslUploadCertificateInstanceInvokeRequestFull(t *testing.T, oldCertificateId string, resourceTypes []string, publicKey, privateKey string, certificateId *string, regions []tftypes.Value, timeoutsVal tftypes.Value) action.InvokeRequest {
	s := sslUploadCertificateInstanceSchema(t)

	resourceTypesSetType := tftypes.Set{ElementType: tftypes.String}
	resourceTypeVals := make([]tftypes.Value, 0, len(resourceTypes))
	for _, rt := range resourceTypes {
		resourceTypeVals = append(resourceTypeVals, tftypes.NewValue(tftypes.String, rt))
	}
	var resourceTypesVal tftypes.Value
	if len(resourceTypeVals) == 0 {
		resourceTypesVal = tftypes.NewValue(resourceTypesSetType, nil)
	} else {
		resourceTypesVal = tftypes.NewValue(resourceTypesSetType, resourceTypeVals)
	}

	var certificateIdVal tftypes.Value
	if certificateId == nil {
		certificateIdVal = tftypes.NewValue(tftypes.String, nil)
	} else {
		certificateIdVal = tftypes.NewValue(tftypes.String, *certificateId)
	}

	regionsListType := tftypes.List{ElementType: sslUploadRegionObjType()}
	var regionsVal tftypes.Value
	if len(regions) == 0 {
		regionsVal = tftypes.NewValue(regionsListType, nil)
	} else {
		regionsVal = tftypes.NewValue(regionsListType, regions)
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"old_certificate_id":      tftypes.String,
			"resource_types":          resourceTypesSetType,
			"certificate_public_key":  tftypes.String,
			"certificate_private_key": tftypes.String,
			"certificate_id":          tftypes.String,
			"resource_types_regions":  regionsListType,
			"timeouts":                sslUploadTimeoutsObjType(),
		},
	}, map[string]tftypes.Value{
		"old_certificate_id":      tftypes.NewValue(tftypes.String, oldCertificateId),
		"resource_types":          resourceTypesVal,
		"certificate_public_key":  tftypes.NewValue(tftypes.String, publicKey),
		"certificate_private_key": tftypes.NewValue(tftypes.String, privateKey),
		"certificate_id":          certificateIdVal,
		"resource_types_regions":  regionsVal,
		"timeouts":                timeoutsVal,
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setSslUploadCertificateInstanceClient injects a mock connectivity client
// into the action instance via the Configure method.
func setSslUploadCertificateInstanceClient(a *providerSsl.SslUploadCertificateInstance, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/ssl/ -run "TestSslUploadCertificateInstance" -v -count=1 -gcflags="all=-l"

// TestSslUploadCertificateInstance_Invoke_Success tests a successful invoke
func TestSslUploadCertificateInstance_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	sslClient := &sslSdk.Client{}
	patches.ApplyMethodReturn(client, "UseSSLCertificateClient", sslClient)

	var gotRequest *sslSdk.UploadUpdateCertificateInstanceRequest
	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateInstanceWithContext", func(_ context.Context, request *sslSdk.UploadUpdateCertificateInstanceRequest) (*sslSdk.UploadUpdateCertificateInstanceResponse, error) {
		gotRequest = request
		resp := sslSdk.NewUploadUpdateCertificateInstanceResponse()
		resp.Response = &sslSdk.UploadUpdateCertificateInstanceResponseParams{
			DeployRecordId: ptrUint64SslUploadCertInstance(12345),
			DeployStatus:   ptrInt64SslUploadCertInstance(1),
			RequestId:      ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecordWithContext", func(_ context.Context, request *sslSdk.DescribeHostUploadUpdateRecordRequest) (*sslSdk.DescribeHostUploadUpdateRecordResponse, error) {
		resp := sslSdk.NewDescribeHostUploadUpdateRecordResponse()
		resp.Response = &sslSdk.DescribeHostUploadUpdateRecordResponseParams{
			TotalCount: ptrUint64SslUploadCertInstance(1),
			DeployRecordList: []*sslSdk.UploadUpdateRecordInfo{
				{
					Id:        ptrUint64SslUploadCertInstance(12345),
					OldCertId: ptrStringSslUploadCertInstance("f8xnynbd"),
					Status:    ptrUint64SslUploadCertInstance(1),
				},
			},
			RequestId: ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	a := providerSsl.NewSslUploadCertificateInstance()
	setSslUploadCertificateInstanceClient(a.(*providerSsl.SslUploadCertificateInstance), client)

	certificateId := "new-cert-id"
	req := newSslUploadCertificateInstanceInvokeRequestFull(t,
		"f8xnynbd",
		[]string{"clb"},
		"public-key-content",
		"private-key-content",
		&certificateId,
		[]tftypes.Value{newSslUploadRegionValue("clb", []string{"ap-guangzhou"})},
		sslUploadNullTimeoutsValue())
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, "f8xnynbd", *gotRequest.OldCertificateId)
	assert.Equal(t, []*string{ptrStringSslUploadCertInstance("clb")}, gotRequest.ResourceTypes)
	// when certificate_id is set together with key content, certificate_id
	// takes precedence and the key content is cleared
	assert.Equal(t, "", *gotRequest.CertificatePublicKey)
	assert.Equal(t, "", *gotRequest.CertificatePrivateKey)
	assert.Equal(t, "new-cert-id", *gotRequest.CertificateId)
	assert.Len(t, gotRequest.ResourceTypesRegions, 1)
	assert.Equal(t, "clb", *gotRequest.ResourceTypesRegions[0].ResourceType)
	assert.Equal(t, []*string{ptrStringSslUploadCertInstance("ap-guangzhou")}, gotRequest.ResourceTypesRegions[0].Regions)
}

// TestSslUploadCertificateInstance_Invoke_SuccessWithCertificateId tests a
// successful invoke where the new certificate is given by certificate_id only
func TestSslUploadCertificateInstance_Invoke_SuccessWithCertificateId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	sslClient := &sslSdk.Client{}
	patches.ApplyMethodReturn(client, "UseSSLCertificateClient", sslClient)

	var gotRequest *sslSdk.UploadUpdateCertificateInstanceRequest
	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateInstanceWithContext", func(_ context.Context, request *sslSdk.UploadUpdateCertificateInstanceRequest) (*sslSdk.UploadUpdateCertificateInstanceResponse, error) {
		gotRequest = request
		resp := sslSdk.NewUploadUpdateCertificateInstanceResponse()
		resp.Response = &sslSdk.UploadUpdateCertificateInstanceResponseParams{
			DeployRecordId: ptrUint64SslUploadCertInstance(12345),
			DeployStatus:   ptrInt64SslUploadCertInstance(1),
			RequestId:      ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecordWithContext", func(_ context.Context, request *sslSdk.DescribeHostUploadUpdateRecordRequest) (*sslSdk.DescribeHostUploadUpdateRecordResponse, error) {
		resp := sslSdk.NewDescribeHostUploadUpdateRecordResponse()
		resp.Response = &sslSdk.DescribeHostUploadUpdateRecordResponseParams{
			TotalCount: ptrUint64SslUploadCertInstance(1),
			DeployRecordList: []*sslSdk.UploadUpdateRecordInfo{
				{
					Id:     ptrUint64SslUploadCertInstance(12345),
					Status: ptrUint64SslUploadCertInstance(1),
				},
			},
			RequestId: ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	a := providerSsl.NewSslUploadCertificateInstance()
	setSslUploadCertificateInstanceClient(a.(*providerSsl.SslUploadCertificateInstance), client)

	certificateId := "new-cert-id"
	req := newSslUploadCertificateInstanceInvokeRequestFull(t,
		"f8xnynbd",
		[]string{"clb"},
		"",
		"",
		&certificateId,
		nil,
		sslUploadNullTimeoutsValue())
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, "new-cert-id", *gotRequest.CertificateId)
	// key content must be cleared when certificate_id is set
	assert.Equal(t, "", *gotRequest.CertificatePublicKey)
	assert.Equal(t, "", *gotRequest.CertificatePrivateKey)
}

// TestSslUploadCertificateInstance_Invoke_CreateRetry tests that a zero
// DeployRecordId is retried until the task is created
func TestSslUploadCertificateInstance_Invoke_CreateRetry(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	sslClient := &sslSdk.Client{}
	patches.ApplyMethodReturn(client, "UseSSLCertificateClient", sslClient)

	callCount := 0
	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateInstanceWithContext", func(_ context.Context, request *sslSdk.UploadUpdateCertificateInstanceRequest) (*sslSdk.UploadUpdateCertificateInstanceResponse, error) {
		callCount++
		resp := sslSdk.NewUploadUpdateCertificateInstanceResponse()
		if callCount < 2 {
			// first call: task creation still in progress
			resp.Response = &sslSdk.UploadUpdateCertificateInstanceResponseParams{
				DeployRecordId: ptrUint64SslUploadCertInstance(0),
				DeployStatus:   ptrInt64SslUploadCertInstance(0),
				RequestId:      ptrStringSslUploadCertInstance("fake-request-id"),
			}
		} else {
			// second call: task created successfully
			resp.Response = &sslSdk.UploadUpdateCertificateInstanceResponseParams{
				DeployRecordId: ptrUint64SslUploadCertInstance(12345),
				DeployStatus:   ptrInt64SslUploadCertInstance(1),
				RequestId:      ptrStringSslUploadCertInstance("fake-request-id"),
			}
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecordWithContext", func(_ context.Context, request *sslSdk.DescribeHostUploadUpdateRecordRequest) (*sslSdk.DescribeHostUploadUpdateRecordResponse, error) {
		resp := sslSdk.NewDescribeHostUploadUpdateRecordResponse()
		resp.Response = &sslSdk.DescribeHostUploadUpdateRecordResponseParams{
			TotalCount: ptrUint64SslUploadCertInstance(1),
			DeployRecordList: []*sslSdk.UploadUpdateRecordInfo{
				{
					Id:     ptrUint64SslUploadCertInstance(12345),
					Status: ptrUint64SslUploadCertInstance(1),
				},
			},
			RequestId: ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	a := providerSsl.NewSslUploadCertificateInstance()
	setSslUploadCertificateInstanceClient(a.(*providerSsl.SslUploadCertificateInstance), client)

	req := newSslUploadCertificateInstanceInvokeRequest(t,
		"f8xnynbd",
		[]string{"clb"},
		"public-key-content",
		"private-key-content")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.GreaterOrEqual(t, callCount, 2)
}

// TestSslUploadCertificateInstance_Invoke_DeployTimeout tests the error
// surfaced when the deploy record does not complete within the invoke timeout
func TestSslUploadCertificateInstance_Invoke_DeployTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	sslClient := &sslSdk.Client{}
	patches.ApplyMethodReturn(client, "UseSSLCertificateClient", sslClient)

	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateInstanceWithContext", func(_ context.Context, request *sslSdk.UploadUpdateCertificateInstanceRequest) (*sslSdk.UploadUpdateCertificateInstanceResponse, error) {
		resp := sslSdk.NewUploadUpdateCertificateInstanceResponse()
		resp.Response = &sslSdk.UploadUpdateCertificateInstanceResponseParams{
			DeployRecordId: ptrUint64SslUploadCertInstance(12345),
			DeployStatus:   ptrInt64SslUploadCertInstance(1),
			RequestId:      ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecordWithContext", func(_ context.Context, request *sslSdk.DescribeHostUploadUpdateRecordRequest) (*sslSdk.DescribeHostUploadUpdateRecordResponse, error) {
		resp := sslSdk.NewDescribeHostUploadUpdateRecordResponse()
		resp.Response = &sslSdk.DescribeHostUploadUpdateRecordResponseParams{
			TotalCount: ptrUint64SslUploadCertInstance(1),
			DeployRecordList: []*sslSdk.UploadUpdateRecordInfo{
				{
					Id:     ptrUint64SslUploadCertInstance(12345),
					Status: ptrUint64SslUploadCertInstance(0),
				},
			},
			RequestId: ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	a := providerSsl.NewSslUploadCertificateInstance()
	setSslUploadCertificateInstanceClient(a.(*providerSsl.SslUploadCertificateInstance), client)

	req := newSslUploadCertificateInstanceInvokeRequestFull(t,
		"f8xnynbd",
		[]string{"clb"},
		"public-key-content",
		"private-key-content",
		nil,
		nil,
		tftypes.NewValue(sslUploadTimeoutsObjType(), map[string]tftypes.Value{
			"invoke": tftypes.NewValue(tftypes.String, "1s"),
		}))
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "waiting for ssl upload certificate instance deploy record")
}

// TestSslUploadCertificateInstance_Invoke_DeployInProgressThenDone tests that
// the action polls the deploy record until it completes instead of failing on
// the first in-progress response
func TestSslUploadCertificateInstance_Invoke_DeployInProgressThenDone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	sslClient := &sslSdk.Client{}
	patches.ApplyMethodReturn(client, "UseSSLCertificateClient", sslClient)

	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateInstanceWithContext", func(_ context.Context, request *sslSdk.UploadUpdateCertificateInstanceRequest) (*sslSdk.UploadUpdateCertificateInstanceResponse, error) {
		resp := sslSdk.NewUploadUpdateCertificateInstanceResponse()
		resp.Response = &sslSdk.UploadUpdateCertificateInstanceResponseParams{
			DeployRecordId: ptrUint64SslUploadCertInstance(12345),
			DeployStatus:   ptrInt64SslUploadCertInstance(1),
			RequestId:      ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	statusCallCount := 0
	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecordWithContext", func(_ context.Context, request *sslSdk.DescribeHostUploadUpdateRecordRequest) (*sslSdk.DescribeHostUploadUpdateRecordResponse, error) {
		statusCallCount++
		status := uint64(0)
		if statusCallCount >= 2 {
			status = 1
		}
		resp := sslSdk.NewDescribeHostUploadUpdateRecordResponse()
		resp.Response = &sslSdk.DescribeHostUploadUpdateRecordResponseParams{
			TotalCount: ptrUint64SslUploadCertInstance(1),
			DeployRecordList: []*sslSdk.UploadUpdateRecordInfo{
				{
					Id:     ptrUint64SslUploadCertInstance(12345),
					Status: &status,
				},
			},
			RequestId: ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	a := providerSsl.NewSslUploadCertificateInstance()
	setSslUploadCertificateInstanceClient(a.(*providerSsl.SslUploadCertificateInstance), client)

	req := newSslUploadCertificateInstanceInvokeRequestFull(t,
		"f8xnynbd",
		[]string{"clb"},
		"public-key-content",
		"private-key-content",
		nil,
		nil,
		tftypes.NewValue(sslUploadTimeoutsObjType(), map[string]tftypes.Value{
			"invoke": tftypes.NewValue(tftypes.String, "30s"),
		}))
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.GreaterOrEqual(t, statusCallCount, 2, "the deploy record should be polled more than once")
}

// TestSslUploadCertificateInstance_Invoke_DeployFailedStatus tests that an
// unexpected deploy record status (e.g. 2) is treated as a failure
func TestSslUploadCertificateInstance_Invoke_DeployFailedStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	sslClient := &sslSdk.Client{}
	patches.ApplyMethodReturn(client, "UseSSLCertificateClient", sslClient)

	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateInstanceWithContext", func(_ context.Context, request *sslSdk.UploadUpdateCertificateInstanceRequest) (*sslSdk.UploadUpdateCertificateInstanceResponse, error) {
		resp := sslSdk.NewUploadUpdateCertificateInstanceResponse()
		resp.Response = &sslSdk.UploadUpdateCertificateInstanceResponseParams{
			DeployRecordId: ptrUint64SslUploadCertInstance(12345),
			DeployStatus:   ptrInt64SslUploadCertInstance(1),
			RequestId:      ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(sslClient, "DescribeHostUploadUpdateRecordWithContext", func(_ context.Context, request *sslSdk.DescribeHostUploadUpdateRecordRequest) (*sslSdk.DescribeHostUploadUpdateRecordResponse, error) {
		resp := sslSdk.NewDescribeHostUploadUpdateRecordResponse()
		resp.Response = &sslSdk.DescribeHostUploadUpdateRecordResponseParams{
			TotalCount: ptrUint64SslUploadCertInstance(1),
			DeployRecordList: []*sslSdk.UploadUpdateRecordInfo{
				{
					Id:     ptrUint64SslUploadCertInstance(12345),
					Status: ptrUint64SslUploadCertInstance(2),
				},
			},
			RequestId: ptrStringSslUploadCertInstance("fake-request-id"),
		}
		return resp, nil
	})

	a := providerSsl.NewSslUploadCertificateInstance()
	setSslUploadCertificateInstanceClient(a.(*providerSsl.SslUploadCertificateInstance), client)

	req := newSslUploadCertificateInstanceInvokeRequest(t,
		"f8xnynbd",
		[]string{"clb"},
		"public-key-content",
		"private-key-content")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "unexpected status 2")
}

// TestSslUploadCertificateInstance_Invoke_APIError tests API error handling
func TestSslUploadCertificateInstance_Invoke_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	sslClient := &sslSdk.Client{}
	patches.ApplyMethodReturn(client, "UseSSLCertificateClient", sslClient)

	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateInstanceWithContext", func(_ context.Context, request *sslSdk.UploadUpdateCertificateInstanceRequest) (*sslSdk.UploadUpdateCertificateInstanceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound.CertificateNotFound, Message=certificate not found")
	})

	a := providerSsl.NewSslUploadCertificateInstance()
	setSslUploadCertificateInstanceClient(a.(*providerSsl.SslUploadCertificateInstance), client)

	req := newSslUploadCertificateInstanceInvokeRequest(t,
		"f8xnynbd",
		[]string{"clb"},
		"public-key-content",
		"private-key-content")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "uploading ssl certificate instance")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "ResourceNotFound.CertificateNotFound")
}

// TestSslUploadCertificateInstance_Invoke_MissingInput tests missing required
// inputs are rejected before any API call
func TestSslUploadCertificateInstance_Invoke_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	sslClient := &sslSdk.Client{}
	patches.ApplyMethodReturn(client, "UseSSLCertificateClient", sslClient)

	apiCalled := false
	patches.ApplyMethodFunc(sslClient, "UploadUpdateCertificateInstanceWithContext", func(_ context.Context, request *sslSdk.UploadUpdateCertificateInstanceRequest) (*sslSdk.UploadUpdateCertificateInstanceResponse, error) {
		apiCalled = true
		return sslSdk.NewUploadUpdateCertificateInstanceResponse(), nil
	})

	a := providerSsl.NewSslUploadCertificateInstance()
	setSslUploadCertificateInstanceClient(a.(*providerSsl.SslUploadCertificateInstance), client)

	// empty old_certificate_id
	req := newSslUploadCertificateInstanceInvokeRequest(t,
		"",
		[]string{"clb"},
		"public-key-content",
		"private-key-content")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing old_certificate_id")
	assert.False(t, apiCalled)

	// public key only, without private key or certificate_id
	req = newSslUploadCertificateInstanceInvokeRequest(t,
		"f8xnynbd",
		[]string{"clb"},
		"public-key-content",
		"")
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing certificate")
	assert.False(t, apiCalled)

	// private key only, without public key or certificate_id
	req = newSslUploadCertificateInstanceInvokeRequest(t,
		"f8xnynbd",
		[]string{"clb"},
		"",
		"private-key-content")
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing certificate")
	assert.False(t, apiCalled)

	// neither certificate_id nor key content
	req = newSslUploadCertificateInstanceInvokeRequest(t,
		"f8xnynbd",
		[]string{"clb"},
		"",
		"")
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing certificate")
	assert.False(t, apiCalled)

	// null resource_types
	req = newSslUploadCertificateInstanceInvokeRequest(t,
		"f8xnynbd",
		nil,
		"public-key-content",
		"private-key-content")
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing resource_types")
	assert.False(t, apiCalled)
}

// TestSslUploadCertificateInstance_Invoke_ClientNotConfigured tests the
// nil-client guard
func TestSslUploadCertificateInstance_Invoke_ClientNotConfigured(t *testing.T) {
	a := providerSsl.NewSslUploadCertificateInstance()

	req := newSslUploadCertificateInstanceInvokeRequest(t,
		"f8xnynbd",
		[]string{"clb"},
		"public-key-content",
		"private-key-content")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestSslUploadCertificateInstance_MetadataAndSchema validates metadata and
// schema definition
func TestSslUploadCertificateInstance_MetadataAndSchema(t *testing.T) {
	a := providerSsl.NewSslUploadCertificateInstance()

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_ssl_upload_certificate_instance", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "old_certificate_id")
	assert.Contains(t, attrs, "resource_types")
	assert.Contains(t, attrs, "certificate_public_key")
	assert.Contains(t, attrs, "certificate_private_key")
	assert.Contains(t, attrs, "certificate_id")

	assert.True(t, attrs["old_certificate_id"].IsRequired())
	assert.True(t, attrs["resource_types"].IsRequired())
	assert.True(t, attrs["certificate_public_key"].IsOptional())
	assert.True(t, attrs["certificate_private_key"].IsOptional())
	assert.True(t, attrs["certificate_id"].IsOptional())

	resourceTypesAttr, ok := attrs["resource_types"].(action_schema.SetAttribute)
	assert.True(t, ok)
	assert.Equal(t, types.StringType, resourceTypesAttr.ElementType)

	publicKeyAttr, ok := attrs["certificate_public_key"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.False(t, publicKeyAttr.Required)
	// WriteOnly must NOT be set on action attributes: the framework
	// hard-codes the write-only client capability to false for action
	// schemas, so a non-null write-only value always fails validation.
	assert.False(t, publicKeyAttr.WriteOnly)

	privateKeyAttr, ok := attrs["certificate_private_key"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.False(t, privateKeyAttr.Required)
	assert.False(t, privateKeyAttr.WriteOnly)

	blocks := schemaResp.Schema.Blocks
	assert.Contains(t, blocks, "resource_types_regions")
	assert.Contains(t, blocks, "timeouts")
}
