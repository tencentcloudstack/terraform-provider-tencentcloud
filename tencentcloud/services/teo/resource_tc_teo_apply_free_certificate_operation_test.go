package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaForApplyFreeCertificate struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForApplyFreeCertificate) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForApplyFreeCertificate{}

func newMockMetaForApplyFreeCertificate() *mockMetaForApplyFreeCertificate {
	return &mockMetaForApplyFreeCertificate{client: &connectivity.TencentCloudClient{}}
}

func ptrStrApplyFreeCert(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestApplyFreeCertificate" -v -count=1 -gcflags="all=-l"

// TestApplyFreeCertificate_Create_DnsVerification tests Create with dns_challenge verification method
func TestApplyFreeCertificate_Create_DnsVerification(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForApplyFreeCertificate().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ApplyFreeCertificate", func(request *teov20220901.ApplyFreeCertificateRequest) (*teov20220901.ApplyFreeCertificateResponse, error) {
		resp := teov20220901.NewApplyFreeCertificateResponse()
		resp.Response = &teov20220901.ApplyFreeCertificateResponseParams{
			DnsVerification: &teov20220901.DnsVerification{
				Subdomain:   ptrStrApplyFreeCert("_acme-challenge.www"),
				RecordType:  ptrStrApplyFreeCert("CNAME"),
				RecordValue: ptrStrApplyFreeCert("abc123.acme.edgeone.com"),
			},
			RequestId: ptrStrApplyFreeCert("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForApplyFreeCertificate()
	res := teo.ResourceTencentCloudTeoApplyFreeCertificateOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-abc123",
		"domain":              "www.example.com",
		"verification_method": "dns_challenge",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-abc123#www.example.com", d.Id())

	dnsVerification := d.Get("dns_verification").([]interface{})
	assert.Equal(t, 1, len(dnsVerification))
	dnsMap := dnsVerification[0].(map[string]interface{})
	assert.Equal(t, "_acme-challenge.www", dnsMap["subdomain"])
	assert.Equal(t, "CNAME", dnsMap["record_type"])
	assert.Equal(t, "abc123.acme.edgeone.com", dnsMap["record_value"])
}

// TestApplyFreeCertificate_Create_FileVerification tests Create with http_challenge verification method
func TestApplyFreeCertificate_Create_FileVerification(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForApplyFreeCertificate().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ApplyFreeCertificate", func(request *teov20220901.ApplyFreeCertificateRequest) (*teov20220901.ApplyFreeCertificateResponse, error) {
		resp := teov20220901.NewApplyFreeCertificateResponse()
		resp.Response = &teov20220901.ApplyFreeCertificateResponseParams{
			FileVerification: &teov20220901.FileVerification{
				Path:    ptrStrApplyFreeCert("/.well-known/teo-verification/z12h416twn.txt"),
				Content: ptrStrApplyFreeCert("teo-verification-content-abc123"),
			},
			RequestId: ptrStrApplyFreeCert("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForApplyFreeCertificate()
	res := teo.ResourceTencentCloudTeoApplyFreeCertificateOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-abc123",
		"domain":              "www.example.com",
		"verification_method": "http_challenge",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-abc123#www.example.com", d.Id())

	fileVerification := d.Get("file_verification").([]interface{})
	assert.Equal(t, 1, len(fileVerification))
	fileMap := fileVerification[0].(map[string]interface{})
	assert.Equal(t, "/.well-known/teo-verification/z12h416twn.txt", fileMap["path"])
	assert.Equal(t, "teo-verification-content-abc123", fileMap["content"])
}

// TestApplyFreeCertificate_Create_NilResponse tests Create handles nil response
func TestApplyFreeCertificate_Create_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForApplyFreeCertificate().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ApplyFreeCertificate", func(request *teov20220901.ApplyFreeCertificateRequest) (*teov20220901.ApplyFreeCertificateResponse, error) {
		resp := teov20220901.NewApplyFreeCertificateResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaForApplyFreeCertificate()
	res := teo.ResourceTencentCloudTeoApplyFreeCertificateOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-abc123",
		"domain":              "www.example.com",
		"verification_method": "dns_challenge",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "response is nil")
}

// TestApplyFreeCertificate_Read tests Read is a no-op
func TestApplyFreeCertificate_Read(t *testing.T) {
	meta := newMockMetaForApplyFreeCertificate()
	res := teo.ResourceTencentCloudTeoApplyFreeCertificateOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-abc123",
		"domain":              "www.example.com",
		"verification_method": "dns_challenge",
	})
	d.SetId("zone-abc123#www.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

// TestApplyFreeCertificate_Delete tests Delete is a no-op
func TestApplyFreeCertificate_Delete(t *testing.T) {
	meta := newMockMetaForApplyFreeCertificate()
	res := teo.ResourceTencentCloudTeoApplyFreeCertificateOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-abc123",
		"domain":              "www.example.com",
		"verification_method": "dns_challenge",
	})
	d.SetId("zone-abc123#www.example.com")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
