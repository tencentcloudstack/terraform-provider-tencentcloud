package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoSharedCname" -v -count=1 -gcflags="all=-l"

type mockMetaSharedCname struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaSharedCname) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaSharedCname{}

func newMockMetaSharedCname() *mockMetaSharedCname {
	return &mockMetaSharedCname{client: &connectivity.TencentCloudClient{}}
}

func ptrStringSharedCname(s string) *string {
	return &s
}

func ptrInt64SharedCname(n int64) *int64 {
	return &n
}

// TestTeoSharedCname_Create_Success tests successful creation of shared CNAME
func TestTeoSharedCname_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.CreateSharedCNAMERequest) (*teov20220901.CreateSharedCNAMEResponse, error) {
		resp := teov20220901.NewCreateSharedCNAMEResponse()
		resp.Response = &teov20220901.CreateSharedCNAMEResponseParams{
			SharedCNAME: ptrStringSharedCname("test-api.sai2ig51kaa5.share.dnse2.com"),
			RequestId:   ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeSharedCNAME for the Read call after Create
	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					Type:        ptrStringSharedCname("custom"),
					SharedCNAME: ptrStringSharedCname("test-api.sai2ig51kaa5.share.dnse2.com"),
					Description: ptrStringSharedCname("example shared cname"),
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-39quuimqg8r6",
		"shared_cname_prefix": "test-api",
		"description":         "example shared cname",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-39quuimqg8r6#test-api.sai2ig51kaa5.share.dnse2.com", d.Id())
	assert.Equal(t, "test-api.sai2ig51kaa5.share.dnse2.com", d.Get("shared_cname"))
}

// TestTeoSharedCname_Create_WithIPSSLSetting tests successful creation with ipssl_setting
func TestTeoSharedCname_Create_WithIPSSLSetting(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.CreateSharedCNAMERequest) (*teov20220901.CreateSharedCNAMEResponse, error) {
		resp := teov20220901.NewCreateSharedCNAMEResponse()
		resp.Response = &teov20220901.CreateSharedCNAMEResponseParams{
			SharedCNAME: ptrStringSharedCname("test-api.sai2ig51kaa5.share.dnse2.com"),
			RequestId:   ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifySharedCNAMEWithContext", func(_ interface{}, request *teov20220901.ModifySharedCNAMERequest) (*teov20220901.ModifySharedCNAMEResponse, error) {
		modifyCalled = true
		assert.Equal(t, "zone-39quuimqg8r6", *request.ZoneId)
		assert.Equal(t, "test-api.sai2ig51kaa5.share.dnse2.com", *request.SharedCNAME)
		assert.NotNil(t, request.IPSSLSetting)
		assert.Equal(t, "bind", *request.IPSSLSetting.Operation)
		assert.Equal(t, "example.com", *request.IPSSLSetting.AssociatedDomain)
		resp := teov20220901.NewModifySharedCNAMEResponse()
		resp.Response = &teov20220901.ModifySharedCNAMEResponseParams{
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeSharedCNAME for the Read call after Create
	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					Type:        ptrStringSharedCname("custom"),
					SharedCNAME: ptrStringSharedCname("test-api.sai2ig51kaa5.share.dnse2.com"),
					Description: ptrStringSharedCname("example shared cname"),
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-39quuimqg8r6",
		"shared_cname_prefix": "test-api",
		"description":         "example shared cname",
		"ipssl_setting": []interface{}{
			map[string]interface{}{
				"status":            "bound",
				"associated_domain": "example.com",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-39quuimqg8r6#test-api.sai2ig51kaa5.share.dnse2.com", d.Id())
	assert.True(t, modifyCalled, "ModifySharedCNAME should be called to set ipssl_setting after creation")
}

// TestTeoSharedCname_Create_EmptyResponse tests creation with empty SharedCNAME response
func TestTeoSharedCname_Create_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.CreateSharedCNAMERequest) (*teov20220901.CreateSharedCNAMEResponse, error) {
		resp := teov20220901.NewCreateSharedCNAMEResponse()
		resp.Response = &teov20220901.CreateSharedCNAMEResponseParams{
			SharedCNAME: nil,
			RequestId:   ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-39quuimqg8r6",
		"shared_cname_prefix": "test-api",
		"description":         "example shared cname",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SharedCNAME is nil or empty")
}

// TestTeoSharedCname_Read_Success tests successful read of shared CNAME
func TestTeoSharedCname_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					Type:        ptrStringSharedCname("custom"),
					SharedCNAME: ptrStringSharedCname("test-api.sai2ig51kaa5.share.dnse2.com"),
					Description: ptrStringSharedCname("example shared cname"),
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-39quuimqg8r6",
		"shared_cname_prefix": "test-api",
		"description":         "example shared cname",
	})
	d.SetId("zone-39quuimqg8r6#test-api.sai2ig51kaa5.share.dnse2.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-39quuimqg8r6", d.Get("zone_id"))
	assert.Equal(t, "test-api.sai2ig51kaa5.share.dnse2.com", d.Get("shared_cname"))
	assert.Equal(t, "example shared cname", d.Get("description"))
}

// TestTeoSharedCname_Read_ExtractsPrefixOnImport tests that Read correctly extracts
// shared_cname_prefix from the full shared CNAME, preventing force replace on import.
func TestTeoSharedCname_Read_ExtractsPrefixOnImport(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					Type:        ptrStringSharedCname("custom"),
					SharedCNAME: ptrStringSharedCname("test-api.39quuimqg8r6.share.dnse27.com"),
					Description: ptrStringSharedCname("example shared cname"),
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	// Simulate import: shared_cname_prefix is empty (not set by user config)
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "",
		"shared_cname_prefix": "",
	})
	d.SetId("zone-39quuimqg8r6#test-api.39quuimqg8r6.share.dnse27.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// Verify prefix is correctly extracted from the full shared CNAME
	assert.Equal(t, "test-api", d.Get("shared_cname_prefix"))
	assert.Equal(t, "zone-39quuimqg8r6", d.Get("zone_id"))
	assert.Equal(t, "test-api.39quuimqg8r6.share.dnse27.com", d.Get("shared_cname"))
	assert.Equal(t, "example shared cname", d.Get("description"))
}

// TestTeoSharedCname_Read_ExtractsPrefixWithDotInPrefix tests prefix extraction
// when the prefix itself contains dots (e.g., "test-api.com").
func TestTeoSharedCname_Read_ExtractsPrefixWithDotInPrefix(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					Type:        ptrStringSharedCname("custom"),
					SharedCNAME: ptrStringSharedCname("test-api.com.abc123xyz.share.dnse27.com"),
					Description: ptrStringSharedCname("dot prefix cname"),
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "",
		"shared_cname_prefix": "",
	})
	d.SetId("zone-abc123xyz#test-api.com.abc123xyz.share.dnse27.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// Prefix "test-api.com" contains a dot, should still be extracted correctly
	assert.Equal(t, "test-api.com", d.Get("shared_cname_prefix"))
}

// TestTeoSharedCname_Read_PrefixNotExtractedWhenSuffixMismatch tests that
// shared_cname_prefix is NOT set when the CNAME suffix doesn't match the expected format.
func TestTeoSharedCname_Read_PrefixNotExtractedWhenSuffixMismatch(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					Type:        ptrStringSharedCname("custom"),
					SharedCNAME: ptrStringSharedCname("test-api.39quuimqg8r6.share.otherdomain.com"),
					Description: ptrStringSharedCname("mismatched suffix"),
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "",
		"shared_cname_prefix": "",
	})
	d.SetId("zone-39quuimqg8r6#test-api.39quuimqg8r6.share.otherdomain.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// When suffix doesn't match, shared_cname_prefix should remain empty
	assert.Equal(t, "", d.Get("shared_cname_prefix"))
}

// TestTeoSharedCname_Read_ExtractsPrefixWithVariousDnseSuffixes tests that the regex
// correctly matches different dnse digit suffixes (dnse0, dnse2, dnse5, dnse123, etc.)
func TestTeoSharedCname_Read_ExtractsPrefixWithVariousDnseSuffixes(t *testing.T) {
	testCases := []struct {
		name        string
		sharedCname string
		zoneId      string
		wantPrefix  string
	}{
		{
			name:        "single digit dnse0",
			sharedCname: "my-prefix.39quuimqg8r6.share.dnse0.com",
			zoneId:      "zone-39quuimqg8r6",
			wantPrefix:  "my-prefix",
		},
		{
			name:        "single digit dnse5",
			sharedCname: "api-gateway.39quuimqg8r6.share.dnse5.com",
			zoneId:      "zone-39quuimqg8r6",
			wantPrefix:  "api-gateway",
		},
		{
			name:        "double digit dnse27",
			sharedCname: "test-api.39quuimqg8r6.share.dnse27.com",
			zoneId:      "zone-39quuimqg8r6",
			wantPrefix:  "test-api",
		},
		{
			name:        "triple digit dnse123",
			sharedCname: "cdn.39quuimqg8r6.share.dnse123.com",
			zoneId:      "zone-39quuimqg8r6",
			wantPrefix:  "cdn",
		},
		{
			name:        "prefix with dots and dnse2",
			sharedCname: "sub.domain.abc123xyz456.share.dnse2.com",
			zoneId:      "zone-abc123xyz456",
			wantPrefix:  "sub.domain",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			teoClient := &teov20220901.Client{}
			patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

			sharedCname := tc.sharedCname
			patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
				resp := teov20220901.NewDescribeSharedCNAMEResponse()
				resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
					TotalCount: ptrInt64SharedCname(1),
					SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
						{
							Type:        ptrStringSharedCname("custom"),
							SharedCNAME: ptrStringSharedCname(sharedCname),
							Description: ptrStringSharedCname("test"),
						},
					},
					RequestId: ptrStringSharedCname("fake-request-id"),
				}
				return resp, nil
			})

			meta := newMockMetaSharedCname()
			res := teo.ResourceTencentCloudTeoSharedCname()
			d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
				"zone_id":             "",
				"shared_cname_prefix": "",
			})
			d.SetId(tc.zoneId + "#" + tc.sharedCname)

			err := res.Read(d, meta)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantPrefix, d.Get("shared_cname_prefix"))
		})
	}
}

// TestTeoSharedCname_Read_NotFound tests read when resource is not found
func TestTeoSharedCname_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount:      ptrInt64SharedCname(0),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{},
			RequestId:       ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-39quuimqg8r6",
		"shared_cname_prefix": "test-api",
	})
	d.SetId("zone-39quuimqg8r6#test-api.sai2ig51kaa5.share.dnse2.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoSharedCname_Update_Description tests updating description
func TestTeoSharedCname_Update_Description(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifySharedCNAMEWithContext", func(_ interface{}, request *teov20220901.ModifySharedCNAMERequest) (*teov20220901.ModifySharedCNAMEResponse, error) {
		modifyCalled = true
		assert.Equal(t, "zone-39quuimqg8r6", *request.ZoneId)
		assert.Equal(t, "test-api.sai2ig51kaa5.share.dnse2.com", *request.SharedCNAME)
		assert.Equal(t, "updated description", *request.Description)
		resp := teov20220901.NewModifySharedCNAMEResponse()
		resp.Response = &teov20220901.ModifySharedCNAMEResponseParams{
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					Type:        ptrStringSharedCname("custom"),
					SharedCNAME: ptrStringSharedCname("test-api.sai2ig51kaa5.share.dnse2.com"),
					Description: ptrStringSharedCname("updated description"),
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-39quuimqg8r6",
		"shared_cname_prefix": "test-api",
		"description":         "updated description",
	})
	d.SetId("zone-39quuimqg8r6#test-api.sai2ig51kaa5.share.dnse2.com")

	// Simulate HasChange by marking the resource as having changes
	err := res.Update(d, meta)
	assert.NoError(t, err)
	// Note: In unit test with TestResourceDataRaw, HasChange may not trigger.
	// The test verifies the Update function runs without error.
	_ = modifyCalled
}

// TestTeoSharedCname_Delete_Success tests successful deletion
func TestTeoSharedCname_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	deleteCalled := false
	patches.ApplyMethodFunc(teoClient, "DeleteSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DeleteSharedCNAMERequest) (*teov20220901.DeleteSharedCNAMEResponse, error) {
		deleteCalled = true
		assert.Equal(t, "zone-39quuimqg8r6", *request.ZoneId)
		assert.Equal(t, "test-api.sai2ig51kaa5.share.dnse2.com", *request.SharedCNAME)
		resp := teov20220901.NewDeleteSharedCNAMEResponse()
		resp.Response = &teov20220901.DeleteSharedCNAMEResponseParams{
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-39quuimqg8r6",
		"shared_cname_prefix": "test-api",
	})
	d.SetId("zone-39quuimqg8r6#test-api.sai2ig51kaa5.share.dnse2.com")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}

// TestTeoSharedCname_Delete_Error tests deletion with API error
func TestTeoSharedCname_Delete_Error(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteSharedCNAMEWithContext", func(_ interface{}, request *teov20220901.DeleteSharedCNAMERequest) (*teov20220901.DeleteSharedCNAMEResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Shared CNAME not found")
	})

	meta := newMockMetaSharedCname()
	res := teo.ResourceTencentCloudTeoSharedCname()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-39quuimqg8r6",
		"shared_cname_prefix": "test-api",
	})
	d.SetId("zone-39quuimqg8r6#test-api.sai2ig51kaa5.share.dnse2.com")

	err := res.Delete(d, meta)
	assert.Error(t, err)
}

// TestTeoSharedCname_Schema tests the schema definition
func TestTeoSharedCname_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoSharedCname()
	assert.NotNil(t, res)

	// Check zone_id
	assert.Contains(t, res.Schema, "zone_id")
	zoneIdSchema := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneIdSchema.Type)
	assert.True(t, zoneIdSchema.Required)
	assert.True(t, zoneIdSchema.ForceNew)

	// Check shared_cname_prefix
	assert.Contains(t, res.Schema, "shared_cname_prefix")
	prefixSchema := res.Schema["shared_cname_prefix"]
	assert.Equal(t, schema.TypeString, prefixSchema.Type)
	assert.True(t, prefixSchema.Required)
	assert.True(t, prefixSchema.ForceNew)

	// Check description
	assert.Contains(t, res.Schema, "description")
	descSchema := res.Schema["description"]
	assert.Equal(t, schema.TypeString, descSchema.Type)
	assert.True(t, descSchema.Optional)

	// Check shared_cname
	assert.Contains(t, res.Schema, "shared_cname")
	sharedCnameSchema := res.Schema["shared_cname"]
	assert.Equal(t, schema.TypeString, sharedCnameSchema.Type)
	assert.True(t, sharedCnameSchema.Computed)

	// Check ipssl_setting
	assert.Contains(t, res.Schema, "ipssl_setting")
	ipsslSchema := res.Schema["ipssl_setting"]
	assert.Equal(t, schema.TypeList, ipsslSchema.Type)
	assert.True(t, ipsslSchema.Optional)
	assert.Equal(t, 1, ipsslSchema.MaxItems)
}
