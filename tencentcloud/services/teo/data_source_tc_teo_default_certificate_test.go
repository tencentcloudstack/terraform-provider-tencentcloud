package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoDefaultCertificateDataSource" -v -count=1 -gcflags="all=-l"

// TestTeoDefaultCertificateDataSource_ReadSuccess tests successful read with certificate data
func TestTeoDefaultCertificateDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDefaultCertificates", func(request *teov20220901.DescribeDefaultCertificatesRequest) (*teov20220901.DescribeDefaultCertificatesResponse, error) {
		resp := teov20220901.NewDescribeDefaultCertificatesResponse()
		resp.Response = &teov20220901.DescribeDefaultCertificatesResponseParams{
			TotalCount: ptrInt64(1),
			DefaultServerCertInfo: []*teov20220901.DefaultServerCertInfo{
				{
					CertId:        ptrString("cert-abc123"),
					Alias:         ptrString("my-cert"),
					Type:          ptrString("default"),
					ExpireTime:    ptrString("2025-12-31T00:00:00Z"),
					EffectiveTime: ptrString("2024-01-01T00:00:00Z"),
					CommonName:    ptrString("example.com"),
					SubjectAltName: []*string{
						ptrString("www.example.com"),
						ptrString("api.example.com"),
					},
					Status:   ptrString("deployed"),
					Message:  ptrString(""),
					SignAlgo: ptrString("RSA"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoDefaultCertificate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	certInfo := d.Get("default_server_cert_info").([]interface{})
	assert.Len(t, certInfo, 1)
	certMap := certInfo[0].(map[string]interface{})
	assert.Equal(t, "cert-abc123", certMap["cert_id"])
	assert.Equal(t, "my-cert", certMap["alias"])
	assert.Equal(t, "default", certMap["type"])
	assert.Equal(t, "deployed", certMap["status"])
	assert.Equal(t, "RSA", certMap["sign_algo"])
	assert.Equal(t, "example.com", certMap["common_name"])
}

// TestTeoDefaultCertificateDataSource_ReadEmpty tests read with no certificates
func TestTeoDefaultCertificateDataSource_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDefaultCertificates", func(request *teov20220901.DescribeDefaultCertificatesRequest) (*teov20220901.DescribeDefaultCertificatesResponse, error) {
		resp := teov20220901.NewDescribeDefaultCertificatesResponse()
		resp.Response = &teov20220901.DescribeDefaultCertificatesResponseParams{
			TotalCount:            ptrInt64(0),
			DefaultServerCertInfo: []*teov20220901.DefaultServerCertInfo{},
			RequestId:             ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoDefaultCertificate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-empty",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	certInfo := d.Get("default_server_cert_info").([]interface{})
	assert.Len(t, certInfo, 0)
}

// TestTeoDefaultCertificateDataSource_Schema validates schema definition
func TestTeoDefaultCertificateDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoDefaultCertificate()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "default_server_cert_info")
	assert.Contains(t, res.Schema, "result_output_file")

	// zone_id is Optional (not Required)
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Optional)

	// filters is Optional
	filters := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filters.Type)
	assert.True(t, filters.Optional)

	// default_server_cert_info is Computed
	certInfo := res.Schema["default_server_cert_info"]
	assert.Equal(t, schema.TypeList, certInfo.Type)
	assert.True(t, certInfo.Computed)

	// result_output_file is Optional
	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)
}

func ptrInt64(n int64) *int64 {
	return &n
}
