package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// ---- Acceptance Tests ----

func TestAccTencentCloudTeoCertificateConfigResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoCertificateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "mode", "sslcert"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.alias", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.cert_id", "EEIqXrZt"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.common_name"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.sign_algo", "RSA 2048"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.type", "managed"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_certificate_config.certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoCertificateConfig = testAccTeoZone + `

resource "tencentcloud_teo_ownership_verify" "ownership_verify" {
  domain = var.zone_name

  depends_on = [ tencentcloud_teo_zone.basic ]
}

resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
  zone_id     = tencentcloud_teo_zone.basic.id
  domain_name = "test.tf-teo.xyz"

  origin_info {
    origin      = "150.109.8.1"
    origin_type = "IP_DOMAIN"
  }

  depends_on = [ tencentcloud_teo_ownership_verify.ownership_verify ]
}

resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = format("test.%s", var.zone_name)
  mode    = "sslcert"
  zone_id = tencentcloud_teo_zone.basic.id

  server_cert_info {
    alias       = "terraform_test"
    cert_id     = "EEIqXrZt"
    common_name = var.zone_name
    //deploy_time = "2024-04-22T10:34:13Z"
    //expire_time = "2025-04-22T23:59:59Z"
    sign_algo   = "RSA 2048"
    type        = "managed"
  }

  depends_on = [tencentcloud_teo_acceleration_domain.acceleration_domain]
}

`

// ---- Mock helpers ----

// mockMeta implements tccommon.ProviderMeta for certificate config tests
type certConfigMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *certConfigMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &certConfigMockMeta{}

func newCertConfigMockMeta() *certConfigMockMeta {
	return &certConfigMockMeta{client: &connectivity.TencentCloudClient{}}
}

func certConfigPtrString(s string) *string {
	return &s
}

func certConfigPtrInt64(i int64) *int64 {
	return &i
}

// ---- Unit Tests (gomonkey mock) ----

// go test ./tencentcloud/services/teo/ -run "TestTeoCertificateConfig_" -v -count=1 -gcflags="all=-l"

// TestTeoCertificateConfig_Schema validates client_cert_info schema definition
func TestTeoCertificateConfig_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCertificateConfig()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "client_cert_info")

	clientCertInfo := res.Schema["client_cert_info"]
	assert.Equal(t, schema.TypeList, clientCertInfo.Type)
	assert.True(t, clientCertInfo.Optional)
	assert.True(t, clientCertInfo.Computed)
	assert.Equal(t, 1, clientCertInfo.MaxItems)

	// Check nested schema
	elem := clientCertInfo.Elem.(*schema.Resource)
	assert.Contains(t, elem.Schema, "switch")
	assert.Contains(t, elem.Schema, "cert_infos")

	// Check switch field
	switchField := elem.Schema["switch"]
	assert.Equal(t, schema.TypeString, switchField.Type)
	assert.True(t, switchField.Optional)
	assert.True(t, switchField.Computed)

	// Check cert_infos field
	certInfosField := elem.Schema["cert_infos"]
	assert.Equal(t, schema.TypeList, certInfosField.Type)
	assert.True(t, certInfosField.Optional)
	assert.True(t, certInfosField.Computed)

	// Check cert_infos nested schema
	certInfosElem := certInfosField.Elem.(*schema.Resource)
	assert.Contains(t, certInfosElem.Schema, "cert_id")
	assert.Contains(t, certInfosElem.Schema, "alias")
	assert.Contains(t, certInfosElem.Schema, "type")
	assert.Contains(t, certInfosElem.Schema, "expire_time")
	assert.Contains(t, certInfosElem.Schema, "deploy_time")
	assert.Contains(t, certInfosElem.Schema, "sign_algo")
	assert.Contains(t, certInfosElem.Schema, "status")

	// Check cert_id is Optional+Computed
	certIdField := certInfosElem.Schema["cert_id"]
	assert.Equal(t, schema.TypeString, certIdField.Type)
	assert.True(t, certIdField.Optional)
	assert.True(t, certIdField.Computed)

	// Check other fields are Computed
	for _, fieldName := range []string{"alias", "type", "expire_time", "deploy_time", "sign_algo", "status"} {
		field := certInfosElem.Schema[fieldName]
		assert.Equal(t, schema.TypeString, field.Type, fieldName+" should be TypeString")
		assert.True(t, field.Computed, fieldName+" should be Computed")
	}
}

// TestTeoCertificateConfig_CreateWithClientCertInfo tests creating resource with client_cert_info
func TestTeoCertificateConfig_CreateWithClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newCertConfigMockMeta().client, "UseTeoClient", teoClient)

	// Mock ModifyHostsCertificate
	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate",
		func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
			// Verify ClientCertInfo is set correctly
			assert.NotNil(t, request.ClientCertInfo)
			assert.Equal(t, "on", *request.ClientCertInfo.Switch)
			assert.Len(t, request.ClientCertInfo.CertInfos, 1)
			assert.Equal(t, "cert-abc123", *request.ClientCertInfo.CertInfos[0].CertId)

			resp := teov20220901.NewModifyHostsCertificateResponse()
			resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	// Mock DescribeAccelerationDomains for status check and read
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains",
		func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
			resp := teov20220901.NewDescribeAccelerationDomainsResponse()
			resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
				TotalCount: certConfigPtrInt64(1),
				AccelerationDomains: []*teov20220901.AccelerationDomain{
					{
						DomainName:   certConfigPtrString("test.example.com"),
						DomainStatus: certConfigPtrString("online"),
						Certificate: &teov20220901.AccelerationDomainCertificate{
							Mode: certConfigPtrString("sslcert"),
							List: []*teov20220901.CertificateInfo{
								{
									CertId:     certConfigPtrString("cert-abc123"),
									Alias:      certConfigPtrString("test-cert"),
									Type:       certConfigPtrString("managed"),
									ExpireTime: certConfigPtrString("2025-12-31T23:59:59Z"),
									DeployTime: certConfigPtrString("2024-01-01T00:00:00Z"),
									SignAlgo:   certConfigPtrString("RSA 2048"),
									Status:     certConfigPtrString("deployed"),
								},
							},
							ClientCertInfo: &teov20220901.MutualTLS{
								Switch: certConfigPtrString("on"),
								CertInfos: []*teov20220901.CertificateInfo{
									{
										CertId:     certConfigPtrString("cert-abc123"),
										Alias:      certConfigPtrString("client-cert-alias"),
										Type:       certConfigPtrString("upload"),
										ExpireTime: certConfigPtrString("2026-06-30T23:59:59Z"),
										DeployTime: certConfigPtrString("2024-06-01T00:00:00Z"),
										SignAlgo:   certConfigPtrString("ECDSA P-256"),
										Status:     certConfigPtrString("deployed"),
									},
								},
							},
						},
					},
				},
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	// Mock DescribeZones for zone name lookup
	patches.ApplyMethodFunc(teoClient, "DescribeZones",
		func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			resp := teov20220901.NewDescribeZonesResponse()
			resp.Response = &teov20220901.DescribeZonesResponseParams{
				TotalCount: certConfigPtrInt64(1),
				Zones: []*teov20220901.Zone{
					{
						ZoneId:   certConfigPtrString("zone-12345678"),
						ZoneName: certConfigPtrString("example.com"),
					},
				},
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newCertConfigMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"host":    "test.example.com",
		"mode":    "sslcert",
		"server_cert_info": []interface{}{
			map[string]interface{}{
				"cert_id": "cert-abc123",
			},
		},
		"client_cert_info": []interface{}{
			map[string]interface{}{
				"switch": "on",
				"cert_infos": []interface{}{
					map[string]interface{}{
						"cert_id": "cert-abc123",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#test.example.com", d.Id())

	// Verify client_cert_info is read back from API
	clientCertInfo := d.Get("client_cert_info").([]interface{})
	assert.Len(t, clientCertInfo, 1)
	clientCertInfoMap := clientCertInfo[0].(map[string]interface{})
	assert.Equal(t, "on", clientCertInfoMap["switch"])

	certInfos := clientCertInfoMap["cert_infos"].([]interface{})
	assert.Len(t, certInfos, 1)
	certInfoMap := certInfos[0].(map[string]interface{})
	assert.Equal(t, "cert-abc123", certInfoMap["cert_id"])
	assert.Equal(t, "client-cert-alias", certInfoMap["alias"])
	assert.Equal(t, "upload", certInfoMap["type"])
	assert.Equal(t, "2026-06-30T23:59:59Z", certInfoMap["expire_time"])
	assert.Equal(t, "2024-06-01T00:00:00Z", certInfoMap["deploy_time"])
	assert.Equal(t, "ECDSA P-256", certInfoMap["sign_algo"])
	assert.Equal(t, "deployed", certInfoMap["status"])
}

// TestTeoCertificateConfig_UpdateWithClientCertInfo tests updating resource with client_cert_info
func TestTeoCertificateConfig_UpdateWithClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newCertConfigMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate",
		func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
			assert.NotNil(t, request.ClientCertInfo)
			assert.Equal(t, "off", *request.ClientCertInfo.Switch)

			resp := teov20220901.NewModifyHostsCertificateResponse()
			resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains",
		func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
			resp := teov20220901.NewDescribeAccelerationDomainsResponse()
			resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
				TotalCount: certConfigPtrInt64(1),
				AccelerationDomains: []*teov20220901.AccelerationDomain{
					{
						DomainName:   certConfigPtrString("test.example.com"),
						DomainStatus: certConfigPtrString("online"),
						Certificate: &teov20220901.AccelerationDomainCertificate{
							Mode: certConfigPtrString("sslcert"),
							List: []*teov20220901.CertificateInfo{
								{CertId: certConfigPtrString("cert-abc123")},
							},
							ClientCertInfo: &teov20220901.MutualTLS{
								Switch: certConfigPtrString("off"),
							},
						},
					},
				},
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	patches.ApplyMethodFunc(teoClient, "DescribeZones",
		func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			resp := teov20220901.NewDescribeZonesResponse()
			resp.Response = &teov20220901.DescribeZonesResponseParams{
				TotalCount: certConfigPtrInt64(1),
				Zones: []*teov20220901.Zone{
					{ZoneId: certConfigPtrString("zone-12345678"), ZoneName: certConfigPtrString("example.com")},
				},
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newCertConfigMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"host":    "test.example.com",
		"mode":    "sslcert",
		"server_cert_info": []interface{}{
			map[string]interface{}{"cert_id": "cert-abc123"},
		},
		"client_cert_info": []interface{}{
			map[string]interface{}{"switch": "off"},
		},
	})
	d.SetId("zone-12345678#test.example.com")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	clientCertInfo := d.Get("client_cert_info").([]interface{})
	assert.Len(t, clientCertInfo, 1)
	clientCertInfoMap := clientCertInfo[0].(map[string]interface{})
	assert.Equal(t, "off", clientCertInfoMap["switch"])
}

// TestTeoCertificateConfig_ReadWithClientCertInfo tests reading client_cert_info from API response
func TestTeoCertificateConfig_ReadWithClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newCertConfigMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains",
		func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
			resp := teov20220901.NewDescribeAccelerationDomainsResponse()
			resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
				TotalCount: certConfigPtrInt64(1),
				AccelerationDomains: []*teov20220901.AccelerationDomain{
					{
						DomainName:   certConfigPtrString("test.example.com"),
						DomainStatus: certConfigPtrString("online"),
						Certificate: &teov20220901.AccelerationDomainCertificate{
							Mode: certConfigPtrString("sslcert"),
							List: []*teov20220901.CertificateInfo{
								{CertId: certConfigPtrString("server-cert-id"), SignAlgo: certConfigPtrString("RSA 2048")},
							},
							ClientCertInfo: &teov20220901.MutualTLS{
								Switch: certConfigPtrString("on"),
								CertInfos: []*teov20220901.CertificateInfo{
									{
										CertId: certConfigPtrString("client-cert-001"), Alias: certConfigPtrString("client-cert-1"),
										Type: certConfigPtrString("upload"), ExpireTime: certConfigPtrString("2025-12-31T23:59:59Z"),
										DeployTime: certConfigPtrString("2024-01-15T10:30:00Z"), SignAlgo: certConfigPtrString("RSA 4096"),
										Status: certConfigPtrString("deployed"),
									},
									{
										CertId: certConfigPtrString("client-cert-002"), Alias: certConfigPtrString("client-cert-2"),
										Type: certConfigPtrString("managed"), ExpireTime: certConfigPtrString("2026-06-30T23:59:59Z"),
										DeployTime: certConfigPtrString("2024-03-20T08:00:00Z"), SignAlgo: certConfigPtrString("ECDSA P-256"),
										Status: certConfigPtrString("deploying"),
									},
								},
							},
						},
					},
				},
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	patches.ApplyMethodFunc(teoClient, "DescribeZones",
		func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			resp := teov20220901.NewDescribeZonesResponse()
			resp.Response = &teov20220901.DescribeZonesResponseParams{
				TotalCount: certConfigPtrInt64(1),
				Zones: []*teov20220901.Zone{
					{ZoneId: certConfigPtrString("zone-12345678"), ZoneName: certConfigPtrString("example.com")},
				},
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newCertConfigMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"host":    "test.example.com",
	})
	d.SetId("zone-12345678#test.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clientCertInfo := d.Get("client_cert_info").([]interface{})
	assert.Len(t, clientCertInfo, 1)
	clientCertInfoMap := clientCertInfo[0].(map[string]interface{})
	assert.Equal(t, "on", clientCertInfoMap["switch"])

	certInfos := clientCertInfoMap["cert_infos"].([]interface{})
	assert.Len(t, certInfos, 2)

	cert1 := certInfos[0].(map[string]interface{})
	assert.Equal(t, "client-cert-001", cert1["cert_id"])
	assert.Equal(t, "deployed", cert1["status"])

	cert2 := certInfos[1].(map[string]interface{})
	assert.Equal(t, "client-cert-002", cert2["cert_id"])
	assert.Equal(t, "deploying", cert2["status"])
}

// TestTeoCertificateConfig_ReadWithNoClientCertInfo tests read when API returns no ClientCertInfo
func TestTeoCertificateConfig_ReadWithNoClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newCertConfigMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains",
		func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
			resp := teov20220901.NewDescribeAccelerationDomainsResponse()
			resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
				TotalCount: certConfigPtrInt64(1),
				AccelerationDomains: []*teov20220901.AccelerationDomain{
					{
						DomainName:   certConfigPtrString("test.example.com"),
						DomainStatus: certConfigPtrString("online"),
						Certificate: &teov20220901.AccelerationDomainCertificate{
							Mode: certConfigPtrString("eofreecert"),
						},
					},
				},
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	patches.ApplyMethodFunc(teoClient, "DescribeZones",
		func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			resp := teov20220901.NewDescribeZonesResponse()
			resp.Response = &teov20220901.DescribeZonesResponseParams{
				TotalCount: certConfigPtrInt64(1),
				Zones: []*teov20220901.Zone{
					{ZoneId: certConfigPtrString("zone-12345678"), ZoneName: certConfigPtrString("example.com")},
				},
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newCertConfigMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"host":    "test.example.com",
	})
	d.SetId("zone-12345678#test.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clientCertInfo := d.Get("client_cert_info").([]interface{})
	assert.Len(t, clientCertInfo, 0)
}

// TestTeoCertificateConfig_BackwardCompatibility tests existing config without client_cert_info still works
func TestTeoCertificateConfig_BackwardCompatibility(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newCertConfigMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate",
		func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
			assert.Nil(t, request.ClientCertInfo)

			resp := teov20220901.NewModifyHostsCertificateResponse()
			resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains",
		func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
			resp := teov20220901.NewDescribeAccelerationDomainsResponse()
			resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
				TotalCount: certConfigPtrInt64(1),
				AccelerationDomains: []*teov20220901.AccelerationDomain{
					{
						DomainName:   certConfigPtrString("test.example.com"),
						DomainStatus: certConfigPtrString("online"),
						Certificate: &teov20220901.AccelerationDomainCertificate{
							Mode: certConfigPtrString("sslcert"),
							List: []*teov20220901.CertificateInfo{
								{CertId: certConfigPtrString("cert-abc123")},
							},
						},
					},
				},
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	patches.ApplyMethodFunc(teoClient, "DescribeZones",
		func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			resp := teov20220901.NewDescribeZonesResponse()
			resp.Response = &teov20220901.DescribeZonesResponseParams{
				TotalCount: certConfigPtrInt64(1),
				Zones: []*teov20220901.Zone{
					{ZoneId: certConfigPtrString("zone-12345678"), ZoneName: certConfigPtrString("example.com")},
				},
				RequestId: certConfigPtrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newCertConfigMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"host":    "test.example.com",
		"mode":    "sslcert",
		"server_cert_info": []interface{}{
			map[string]interface{}{"cert_id": "cert-abc123"},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#test.example.com", d.Id())
	assert.Equal(t, "sslcert", d.Get("mode"))
}
