package teo_test

import (
	"context"
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

func TestAccTencentCloudTeoConfigGroupVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoConfigGroupVersion,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_teo_config_group_version.teo_config_group_version", "id"),
				resource.TestCheckResourceAttr("tencentcloud_teo_config_group_version.teo_config_group_version", "zone_id", "zone-2xkazzl8yf6k"),
				resource.TestCheckResourceAttr("tencentcloud_teo_config_group_version.teo_config_group_version", "group_id", "cg-3lchxitnb5pb"),
				resource.TestCheckResourceAttr("tencentcloud_teo_config_group_version.teo_config_group_version", "description", "test version"),
				resource.TestCheckResourceAttrSet("tencentcloud_teo_config_group_version.teo_config_group_version", "content"),
			),
		}},
	})
}

const testAccTeoConfigGroupVersion = `

resource "tencentcloud_teo_config_group_version" "teo_config_group_version" {
  content = <<EOT
{
  "FormatVersion": "1.0",
  "ZoneConfig": {
    "SmartRouting": {
      "Switch": "off"
    },
    "Cache": {
      "NoCache": {
        "Switch": "off"
      },
      "FollowOrigin": {
        "Switch": "on",
        "DefaultCache": "on",
        "DefaultCacheStrategy": "on",
        "DefaultCacheTime": 0
      },
      "CustomTime": {
        "Switch": "off",
        "CacheTime": 2592000
      }
    },
    "MaxAge": {
      "FollowOrigin": "on",
      "CacheTime": 600
    },
    "CacheKey": {
      "FullURLCache": "on",
      "QueryString": {
        "Switch": "off",
        "Action": "includeCustom"
      },
      "IgnoreCase": "off"
    },
    "CachePrefresh": {
      "Switch": "off",
      "CacheTimePercent": 90
    },
    "OfflineCache": {
      "Switch": "on"
    },
    "Compression": {
      "Switch": "off",
      "Algorithms": [
        "brotli",
        "gzip"
      ]
    },
    "ForceRedirectHTTPS": {
      "Switch": "off",
      "RedirectStatusCode": 302
    },
    "HSTS": {
      "Switch": "off",
      "Timeout": 0,
      "IncludeSubDomains": "off",
      "Preload": "off"
    },
    "TLSConfig": {
      "Version": [
        "TLSv1",
        "TLSv1.1",
        "TLSv1.2",
        "TLSv1.3"
      ],
      "CipherSuite": "loose-v2023"
    },
    "OCSPStapling": {
      "Switch": "off"
    },
    "HTTP2": {
      "Switch": "off"
    },
    "QUIC": {
      "Switch": "off"
    },
    "UpstreamHTTP2": {
      "Switch": "off"
    },
    "IPv6": {
      "Switch": "off"
    },
    "WebSocket": {
      "Switch": "off",
      "Timeout": 30
    },
    "PostMaxSize": {
      "Switch": "on",
      "MaxSize": 838860800
    },
    "ClientIPHeader": {
      "Switch": "off"
    },
    "ClientIPCountry": {
      "Switch": "off"
    },
    "gRPC": {
      "Switch": "off"
    },
    "NetworkErrorLogging": {
      "Switch": "off"
    },
    "AccelerateMainland": {
      "Switch": "off"
    },
    "StandardDebug": {
      "Switch": "off",
      "AllowClientIPList": [
        "1.14.231.0/24",
        "1.194.255.0/24"
      ],
      "Expires": "2025-09-01T12:45:37Z"
    }
  },
  "Rules": [
    {
      "RuleName": "Web Acceleration - cdn.defaultsetting.cn",
      "Branches": [
        {
          "Condition": "$${http.request.host} in ['cdn.defaultsetting.cn']",
          "Actions": [
            {
              "Name": "Cache",
              "CacheParameters": {
                "CustomTime": {
                  "Switch": "on",
                  "IgnoreCacheControl": "off",
                  "CacheTime": 2592000
                }
              }
            },
            {
              "Name": "CacheKey",
              "CacheKeyParameters": {
                "FullURLCache": "on",
                "QueryString": {
                  "Switch": "off"
                },
                "IgnoreCase": "off"
              }
            }
          ],
          "SubRules": [
            {
              "Branches": [
                {
                  "Condition": "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']",
                  "Actions": [
                    {
                      "Name": "Cache",
                      "CacheParameters": {
                        "NoCache": {
                          "Switch": "on"
                        }
                      }
                    }
                  ]
                }
              ]
            },
            {
              "Branches": [
                {
                  "Condition": "$${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']",
                  "Actions": [
                    {
                      "Name": "MaxAge",
                      "MaxAgeParameters": {
                        "FollowOrigin": "off",
                        "CacheTime": 3600
                      }
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "RuleName": "Web Acceleration - pages.migraine.com.cn",
      "Branches": [
        {
          "Condition": "$${http.request.host} in ['pages.migraine.com.cn']",
          "Actions": [
            {
              "Name": "Cache",
              "CacheParameters": {
                "CustomTime": {
                  "Switch": "on",
                  "IgnoreCacheControl": "off",
                  "CacheTime": 0
                }
              }
            },
            {
              "Name": "CacheKey",
              "CacheKeyParameters": {
                "FullURLCache": "on",
                "QueryString": {
                  "Switch": "off"
                },
                "IgnoreCase": "off"
              }
            }
          ],
          "SubRules": [
            {
              "Branches": [
                {
                  "Condition": "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']",
                  "Actions": [
                    {
                      "Name": "Cache",
                      "CacheParameters": {
                        "NoCache": {
                          "Switch": "on"
                        }
                      }
                    }
                  ]
                }
              ]
            },
            {
              "Branches": [
                {
                  "Condition": "$${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']",
                  "Actions": [
                    {
                      "Name": "MaxAge",
                      "MaxAgeParameters": {
                        "FollowOrigin": "off",
                        "CacheTime": 3600
                      }
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "RuleName": "Skip Pages Authentication Verification",
      "Branches": [
        {
          "Condition": "$${http.request.host} in ['pages.migraine.com.cn']",
          "Actions": [
            {
              "Name": "ModifyRequestHeader",
              "ModifyRequestHeaderParameters": {
                "HeaderActions": [
                  {
                    "Action": "add",
                    "Name": "X-SKIP-TOKEN",
                    "Value": "eop-1022"
                  }
                ]
              }
            }
          ]
        }
      ]
    },
    {
      "RuleName": "Regional Origin Pull",
      "Branches": [
        {
          "Condition": "$${http.request.host} in ['abc.migraine.com.cn']",
          "SubRules": [
            {
              "Branches": [
                {
                  "Condition": "$${http.request.ip.country} in ['Asia']",
                  "Actions": [
                    {
                      "Name": "ModifyOrigin",
                      "ModifyOriginParameters": {
                        "OriginType": "IPDomain",
                        "Origin": "1.2.3.4",
                        "OriginProtocol": "follow",
                        "HTTPOriginPort": 80,
                        "HTTPSOriginPort": 443
                      }
                    }
                  ]
                }
              ]
            },
            {
              "Branches": [
                {
                  "Condition": "$${http.request.ip.country} in ['Africa']",
                  "Actions": [
                    {
                      "Name": "ModifyOrigin",
                      "ModifyOriginParameters": {
                        "OriginType": "IPDomain",
                        "Origin": "3.4.5.6",
                        "OriginProtocol": "follow",
                        "HTTPOriginPort": 80,
                        "HTTPSOriginPort": 443
                      }
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
EOT
  description = "test version"
  group_id    = "cg-3lchxitnb5pb"
  zone_id     = "zone-2xkazzl8yf6k"
}
`

// --- gomonkey mock unit tests for source_version ---
// go test ./tencentcloud/services/teo/ -run "TestTeoConfigGroupVersion" -v -count=1 -gcflags="all=-l"

// mockMetaForConfigGroupVersion implements tccommon.ProviderMeta
type mockMetaForConfigGroupVersion struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForConfigGroupVersion) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForConfigGroupVersion{}

func newMockMetaForConfigGroupVersion() *mockMetaForConfigGroupVersion {
	return &mockMetaForConfigGroupVersion{client: &connectivity.TencentCloudClient{}}
}

func ptrStrConfigGroupVersion(s string) *string {
	return &s
}

// TestTeoConfigGroupVersion_Create_WithSourceVersion tests Create with source_version set
func TestTeoConfigGroupVersion_Create_WithSourceVersion(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigGroupVersion().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateConfigGroupVersionWithContext", func(_ context.Context, request *teov20220901.CreateConfigGroupVersionRequest) (*teov20220901.CreateConfigGroupVersionResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "cg-abcdefgh", *request.GroupId)
		assert.NotNil(t, request.Content)
		assert.NotNil(t, request.SourceVersion)
		assert.Equal(t, "ver-source123", *request.SourceVersion)
		resp := teov20220901.NewCreateConfigGroupVersionResponse()
		resp.Response = &teov20220901.CreateConfigGroupVersionResponseParams{
			VersionId: ptrStrConfigGroupVersion("ver-new123456"),
			RequestId: ptrStrConfigGroupVersion("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeConfigGroupVersionDetail", func(request *teov20220901.DescribeConfigGroupVersionDetailRequest) (*teov20220901.DescribeConfigGroupVersionDetailResponse, error) {
		resp := teov20220901.NewDescribeConfigGroupVersionDetailResponse()
		resp.Response = &teov20220901.DescribeConfigGroupVersionDetailResponseParams{
			ConfigGroupVersionInfo: &teov20220901.ConfigGroupVersionInfo{
				VersionId:     ptrStrConfigGroupVersion("ver-new123456"),
				VersionNumber: ptrStrConfigGroupVersion("1"),
				GroupId:       ptrStrConfigGroupVersion("cg-abcdefgh"),
				GroupType:     ptrStrConfigGroupVersion("l7_acceleration"),
				Description:   ptrStrConfigGroupVersion("test version"),
				SourceVersion: ptrStrConfigGroupVersion("ver-source123"),
				Status:        ptrStrConfigGroupVersion("inactive"),
				CreateTime:    ptrStrConfigGroupVersion("2025-01-01T00:00:00Z"),
			},
			Content:   ptrStrConfigGroupVersion("{}"),
			RequestId: ptrStrConfigGroupVersion("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForConfigGroupVersion()
	res := teo.ResourceTencentCloudTeoConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":        "zone-12345678",
		"group_id":       "cg-abcdefgh",
		"content":        "{}",
		"source_version": "ver-source123",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678"+tccommon.FILED_SP+"cg-abcdefgh"+tccommon.FILED_SP+"ver-new123456", d.Id())
	assert.Equal(t, "ver-source123", d.Get("source_version"))
}

// TestTeoConfigGroupVersion_Create_WithoutSourceVersion tests Create without source_version
func TestTeoConfigGroupVersion_Create_WithoutSourceVersion(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigGroupVersion().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateConfigGroupVersionWithContext", func(_ context.Context, request *teov20220901.CreateConfigGroupVersionRequest) (*teov20220901.CreateConfigGroupVersionResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "cg-abcdefgh", *request.GroupId)
		assert.Nil(t, request.SourceVersion)
		resp := teov20220901.NewCreateConfigGroupVersionResponse()
		resp.Response = &teov20220901.CreateConfigGroupVersionResponseParams{
			VersionId: ptrStrConfigGroupVersion("ver-new789012"),
			RequestId: ptrStrConfigGroupVersion("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeConfigGroupVersionDetail", func(request *teov20220901.DescribeConfigGroupVersionDetailRequest) (*teov20220901.DescribeConfigGroupVersionDetailResponse, error) {
		resp := teov20220901.NewDescribeConfigGroupVersionDetailResponse()
		resp.Response = &teov20220901.DescribeConfigGroupVersionDetailResponseParams{
			ConfigGroupVersionInfo: &teov20220901.ConfigGroupVersionInfo{
				VersionId:     ptrStrConfigGroupVersion("ver-new789012"),
				VersionNumber: ptrStrConfigGroupVersion("2"),
				GroupId:       ptrStrConfigGroupVersion("cg-abcdefgh"),
				GroupType:     ptrStrConfigGroupVersion("l7_acceleration"),
				Description:   ptrStrConfigGroupVersion("test version"),
				Status:        ptrStrConfigGroupVersion("inactive"),
				CreateTime:    ptrStrConfigGroupVersion("2025-01-01T00:00:00Z"),
			},
			Content:   ptrStrConfigGroupVersion("{}"),
			RequestId: ptrStrConfigGroupVersion("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForConfigGroupVersion()
	res := teo.ResourceTencentCloudTeoConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-12345678",
		"group_id": "cg-abcdefgh",
		"content":  "{}",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678"+tccommon.FILED_SP+"cg-abcdefgh"+tccommon.FILED_SP+"ver-new789012", d.Id())
	assert.Equal(t, "", d.Get("source_version"))
}

// TestTeoConfigGroupVersion_Read_WithNilSourceVersion tests Read when SourceVersion is nil
func TestTeoConfigGroupVersion_Read_WithNilSourceVersion(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigGroupVersion().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeConfigGroupVersionDetail", func(request *teov20220901.DescribeConfigGroupVersionDetailRequest) (*teov20220901.DescribeConfigGroupVersionDetailResponse, error) {
		resp := teov20220901.NewDescribeConfigGroupVersionDetailResponse()
		resp.Response = &teov20220901.DescribeConfigGroupVersionDetailResponseParams{
			ConfigGroupVersionInfo: &teov20220901.ConfigGroupVersionInfo{
				VersionId:     ptrStrConfigGroupVersion("ver-new123456"),
				VersionNumber: ptrStrConfigGroupVersion("1"),
				GroupId:       ptrStrConfigGroupVersion("cg-abcdefgh"),
				GroupType:     ptrStrConfigGroupVersion("l7_acceleration"),
				Description:   ptrStrConfigGroupVersion("test version"),
				Status:        ptrStrConfigGroupVersion("active"),
				CreateTime:    ptrStrConfigGroupVersion("2025-01-01T00:00:00Z"),
			},
			Content:   ptrStrConfigGroupVersion("{}"),
			RequestId: ptrStrConfigGroupVersion("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForConfigGroupVersion()
	res := teo.ResourceTencentCloudTeoConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-12345678",
		"group_id": "cg-abcdefgh",
		"content":  "{}",
	})
	d.SetId("zone-12345678" + tccommon.FILED_SP + "cg-abcdefgh" + tccommon.FILED_SP + "ver-new123456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678"+tccommon.FILED_SP+"cg-abcdefgh"+tccommon.FILED_SP+"ver-new123456", d.Id())
	assert.Equal(t, "", d.Get("source_version"))
	assert.Equal(t, "ver-new123456", d.Get("version_id"))
	assert.Equal(t, "active", d.Get("status"))
}

// TestTeoConfigGroupVersion_Schema validates source_version schema definition
func TestTeoConfigGroupVersion_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoConfigGroupVersion()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "source_version")

	sourceVersion := res.Schema["source_version"]
	assert.Equal(t, schema.TypeString, sourceVersion.Type)
	assert.True(t, sourceVersion.Optional)
	assert.False(t, sourceVersion.Required)
	assert.True(t, sourceVersion.ForceNew)
}
