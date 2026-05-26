package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -test.run TestAccTencentCloudTeoZoneSetting_basic -v
func TestAccTencentCloudTeoZoneSetting_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZoneSetting,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneSettingExists("tencentcloud_teo_zone_setting.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_zone_setting.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.cache.0.cache_time", "10"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.follow_origin.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.follow_origin.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.no_cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache.0.no_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.query_string.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.query_string.0.action", "includeCustom"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_key.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_prefresh.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_prefresh.0.percent", "90"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "cache_prefresh.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "client_ip_header.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "client_ip_header.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "compression.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "compression.0.algorithms.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "compression.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "force_redirect.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "force_redirect.0.redirect_status_code", "302"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "force_redirect.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.0.include_sub_domains", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.0.preload", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.hsts.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.http2", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.ocsp_stapling", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "https.0.tls_version.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "ipv6.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "ipv6.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "max_age.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "max_age.0.follow_origin", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "max_age.0.max_age_time", "600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "offline_cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "offline_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "origin.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "origin.0.origin_pull_protocol", "follow"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "post_max_size.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "post_max_size.0.max_size", "838860800"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "post_max_size.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "quic.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "quic.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "smart_routing.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "upstream_http2.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "upstream_http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "web_socket.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone_setting.basic", "web_socket.0.switch", "off"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_zone_setting.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckZoneSettingExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		setting, err := service.DescribeTeoZoneSetting(ctx, rs.Primary.ID)
		if setting == nil {
			return fmt.Errorf("zone setting %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoZoneSetting = testAccTeoZone + `

resource "tencentcloud_teo_zone_setting" "basic" {
  zone_id = tencentcloud_teo_zone.basic.id

  cache {
    cache {
      switch     = "on"
      cache_time = 10
    }
    follow_origin {
      switch = "off"
    }
    no_cache {
      switch = "off"
    }
  }

  cache_key {
    full_url_cache = "on"
    ignore_case    = "off"

    query_string {
      action = "includeCustom"
      switch = "off"
      value  = []
    }
  }

  cache_prefresh {
    percent = 90
    switch  = "on"
  }

  client_ip_header {
    switch = "off"
  }

  compression {
    algorithms = [
      "brotli",
      "gzip",
    ]
    switch = "on"
  }

  force_redirect {
    redirect_status_code = 302
    switch               = "off"
  }

  https {
    http2         = "on"
    ocsp_stapling = "off"
    tls_version   = [
      "TLSv1",
      "TLSv1.1",
      "TLSv1.2",
      "TLSv1.3",
    ]

    hsts {
      include_sub_domains = "off"
      max_age             = 0
      preload             = "off"
      switch              = "off"
    }
  }

  ipv6 {
    switch = "off"
  }

  max_age {
    follow_origin = "on"
    max_age_time  = 600
  }

  offline_cache {
    switch = "on"
  }

  origin {
    backup_origins       = []
    origin_pull_protocol = "follow"
    origins              = []
  }

  post_max_size {
    max_size = 838860800
    switch   = "on"
  }

  quic {
    switch = "off"
  }

  smart_routing {
    switch = "off"
  }

  upstream_http2 {
    switch = "off"
  }

  web_socket {
    switch  = "off"
    timeout = 30
  }
}
`

// mockMetaForZoneSetting implements tccommon.ProviderMeta
type mockMetaForZoneSetting struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForZoneSetting) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForZoneSetting{}

func newMockMetaForZoneSetting() *mockMetaForZoneSetting {
	return &mockMetaForZoneSetting{client: &connectivity.TencentCloudClient{}}
}

func ptrStrZoneSetting(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestZoneSettingJITVideoProcess" -v -count=1 -gcflags="all=-l"

// TestZoneSettingJITVideoProcess_Read_SwitchOn tests Read populates jit_video_process with switch=on
func TestZoneSettingJITVideoProcess_Read_SwitchOn(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForZoneSetting().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZoneSetting", func(request *teov20220901.DescribeZoneSettingRequest) (*teov20220901.DescribeZoneSettingResponse, error) {
		resp := teov20220901.NewDescribeZoneSettingResponse()
		resp.Response = &teov20220901.DescribeZoneSettingResponseParams{
			ZoneSetting: &teov20220901.ZoneSetting{
				JITVideoProcess: &teov20220901.JITVideoProcess{
					Switch: ptrStrZoneSetting("on"),
				},
			},
			RequestId: ptrStrZoneSetting("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForZoneSetting()
	res := svcteo.ResourceTencentCloudTeoZoneSetting()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	jitVideoProcess := d.Get("jit_video_process").([]interface{})
	assert.Equal(t, 1, len(jitVideoProcess))
	jitMap := jitVideoProcess[0].(map[string]interface{})
	assert.Equal(t, "on", jitMap["switch"])
}

// TestZoneSettingJITVideoProcess_Read_SwitchOff tests Read populates jit_video_process with switch=off
func TestZoneSettingJITVideoProcess_Read_SwitchOff(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForZoneSetting().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZoneSetting", func(request *teov20220901.DescribeZoneSettingRequest) (*teov20220901.DescribeZoneSettingResponse, error) {
		resp := teov20220901.NewDescribeZoneSettingResponse()
		resp.Response = &teov20220901.DescribeZoneSettingResponseParams{
			ZoneSetting: &teov20220901.ZoneSetting{
				JITVideoProcess: &teov20220901.JITVideoProcess{
					Switch: ptrStrZoneSetting("off"),
				},
			},
			RequestId: ptrStrZoneSetting("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForZoneSetting()
	res := svcteo.ResourceTencentCloudTeoZoneSetting()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	jitVideoProcess := d.Get("jit_video_process").([]interface{})
	assert.Equal(t, 1, len(jitVideoProcess))
	jitMap := jitVideoProcess[0].(map[string]interface{})
	assert.Equal(t, "off", jitMap["switch"])
}

// TestZoneSettingJITVideoProcess_Read_Nil tests Read handles nil JITVideoProcess
func TestZoneSettingJITVideoProcess_Read_Nil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForZoneSetting().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZoneSetting", func(request *teov20220901.DescribeZoneSettingRequest) (*teov20220901.DescribeZoneSettingResponse, error) {
		resp := teov20220901.NewDescribeZoneSettingResponse()
		resp.Response = &teov20220901.DescribeZoneSettingResponseParams{
			ZoneSetting: &teov20220901.ZoneSetting{
				JITVideoProcess: nil,
			},
			RequestId: ptrStrZoneSetting("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForZoneSetting()
	res := svcteo.ResourceTencentCloudTeoZoneSetting()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	jitVideoProcess := d.Get("jit_video_process").([]interface{})
	assert.Equal(t, 0, len(jitVideoProcess))
}

// TestZoneSettingJITVideoProcess_Update tests Update sends JITVideoProcess in request
func TestZoneSettingJITVideoProcess_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForZoneSetting().client, "UseTeoClient", teoClient)

	var capturedRequest *teov20220901.ModifyZoneSettingRequest
	patches.ApplyMethodFunc(teoClient, "ModifyZoneSettingWithContext", func(ctx context.Context, request *teov20220901.ModifyZoneSettingRequest) (*teov20220901.ModifyZoneSettingResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewModifyZoneSettingResponse()
		resp.Response = &teov20220901.ModifyZoneSettingResponseParams{
			RequestId: ptrStrZoneSetting("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeZoneSetting", func(request *teov20220901.DescribeZoneSettingRequest) (*teov20220901.DescribeZoneSettingResponse, error) {
		resp := teov20220901.NewDescribeZoneSettingResponse()
		resp.Response = &teov20220901.DescribeZoneSettingResponseParams{
			ZoneSetting: &teov20220901.ZoneSetting{
				JITVideoProcess: &teov20220901.JITVideoProcess{
					Switch: ptrStrZoneSetting("on"),
				},
			},
			RequestId: ptrStrZoneSetting("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForZoneSetting()
	res := svcteo.ResourceTencentCloudTeoZoneSetting()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"jit_video_process": []interface{}{
			map[string]interface{}{
				"switch": "on",
			},
		},
	})
	d.SetId("zone-test123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.JITVideoProcess)
	assert.Equal(t, "on", *capturedRequest.JITVideoProcess.Switch)
}
