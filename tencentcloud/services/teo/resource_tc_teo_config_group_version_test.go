package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

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

// go test ./tencentcloud/services/teo/ -run "TestTeoConfigGroupVersion" -v -count=1 -gcflags="all=-l"

// TestTeoConfigGroupVersion_Read_Success tests Read populates version_id and other computed output fields
func TestTeoConfigGroupVersion_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigGroupVersion().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeConfigGroupVersionDetail", func(request *teov20220901.DescribeConfigGroupVersionDetailRequest) (*teov20220901.DescribeConfigGroupVersionDetailResponse, error) {
		resp := teov20220901.NewDescribeConfigGroupVersionDetailResponse()
		resp.Response = &teov20220901.DescribeConfigGroupVersionDetailResponseParams{
			ConfigGroupVersionInfo: &teov20220901.ConfigGroupVersionInfo{
				VersionId:     ptrStrConfigGroupVersion("ver-2kplomhisdcb"),
				VersionNumber: ptrStrConfigGroupVersion("1"),
				GroupId:       ptrStrConfigGroupVersion("cg-3lchxitnb5pb"),
				GroupType:     ptrStrConfigGroupVersion("l7_acceleration"),
				Description:   ptrStrConfigGroupVersion("test version"),
				Status:        ptrStrConfigGroupVersion("inactive"),
				CreateTime:    ptrStrConfigGroupVersion("2025-01-01T00:00:00Z"),
			},
			Content:   ptrStrConfigGroupVersion(`{"FormatVersion":"1.0"}`),
			RequestId: ptrStrConfigGroupVersion("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForConfigGroupVersion()
	res := teo.ResourceTencentCloudTeoConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-2kazzl8yf6k",
		"group_id": "cg-3lchxitnb5pb",
		"content":  `{"FormatVersion":"1.0"}`,
	})
	compositeId := "zone-2kazzl8yf6k" + tccommon.FILED_SP + "cg-3lchxitnb5pb" + tccommon.FILED_SP + "ver-2kplomhisdcb"
	d.SetId(compositeId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, compositeId, d.Id())
	assert.Equal(t, "ver-2kplomhisdcb", d.Get("version_id"))
	assert.Equal(t, "1", d.Get("version_number"))
	assert.Equal(t, "cg-3lchxitnb5pb", d.Get("group_id"))
	assert.Equal(t, "l7_acceleration", d.Get("group_type"))
	assert.Equal(t, "test version", d.Get("description"))
	assert.Equal(t, "inactive", d.Get("status"))
	assert.Equal(t, "2025-01-01T00:00:00Z", d.Get("create_time"))
	assert.Equal(t, `{"FormatVersion":"1.0"}`, d.Get("content"))
	assert.Equal(t, "zone-2kazzl8yf6k", d.Get("zone_id"))
}

// TestTeoConfigGroupVersion_Read_NilConfigGroupVersionInfo tests Read handles nil ConfigGroupVersionInfo without panic
func TestTeoConfigGroupVersion_Read_NilConfigGroupVersionInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigGroupVersion().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeConfigGroupVersionDetail", func(request *teov20220901.DescribeConfigGroupVersionDetailRequest) (*teov20220901.DescribeConfigGroupVersionDetailResponse, error) {
		resp := teov20220901.NewDescribeConfigGroupVersionDetailResponse()
		resp.Response = &teov20220901.DescribeConfigGroupVersionDetailResponseParams{
			ConfigGroupVersionInfo: nil,
			Content:                nil,
			RequestId:              ptrStrConfigGroupVersion("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForConfigGroupVersion()
	res := teo.ResourceTencentCloudTeoConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-2kazzl8yf6k",
		"group_id": "cg-3lchxitnb5pb",
		"content":  `{"FormatVersion":"1.0"}`,
	})
	compositeId := "zone-2kazzl8yf6k" + tccommon.FILED_SP + "cg-3lchxitnb5pb" + tccommon.FILED_SP + "ver-2kplomhisdcb"
	d.SetId(compositeId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// id is preserved because respData != nil (only ConfigGroupVersionInfo is nil)
	assert.Equal(t, compositeId, d.Id())
	// version_id is not set (empty string default)
	assert.Equal(t, "", d.Get("version_id"))
}

// TestTeoConfigGroupVersion_Read_EmptyResponse tests Read clears id after [CRUD] logging when API returns nil response
func TestTeoConfigGroupVersion_Read_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigGroupVersion().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeConfigGroupVersionDetail", func(request *teov20220901.DescribeConfigGroupVersionDetailRequest) (*teov20220901.DescribeConfigGroupVersionDetailResponse, error) {
		// Return a response whose Response is nil, so service.DescribeTeoConfigGroupVersionById returns nil
		resp := &teov20220901.DescribeConfigGroupVersionDetailResponse{}
		return resp, nil
	})

	meta := newMockMetaForConfigGroupVersion()
	res := teo.ResourceTencentCloudTeoConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-2kazzl8yf6k",
		"group_id": "cg-3lchxitnb5pb",
		"content":  `{"FormatVersion":"1.0"}`,
	})
	compositeId := "zone-2kazzl8yf6k" + tccommon.FILED_SP + "cg-3lchxitnb5pb" + tccommon.FILED_SP + "ver-2kplomhisdcb"
	d.SetId(compositeId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// id is cleared because respData == nil, after [CRUD] log preserves the scene
	assert.Equal(t, "", d.Id())
}

// TestTeoConfigGroupVersion_Read_BrokenId tests Read returns error when composite id is broken
func TestTeoConfigGroupVersion_Read_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigGroupVersion().client, "UseTeoV20220901Client", teoClient)

	meta := newMockMetaForConfigGroupVersion()
	res := teo.ResourceTencentCloudTeoConfigGroupVersion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-2kazzl8yf6k",
		"group_id": "cg-3lchxitnb5pb",
		"content":  `{"FormatVersion":"1.0"}`,
	})
	d.SetId("zone-2kazzl8yf6k" + tccommon.FILED_SP + "cg-3lchxitnb5pb")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}
