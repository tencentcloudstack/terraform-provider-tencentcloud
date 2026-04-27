package teo_test

import (
	"fmt"
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

// mockMetaForL7AccSetting implements tccommon.ProviderMeta
type mockMetaForL7AccSetting struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForL7AccSetting) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForL7AccSetting{}

func newMockMetaForL7AccSetting() *mockMetaForL7AccSetting {
	return &mockMetaForL7AccSetting{client: &connectivity.TencentCloudClient{}}
}

func ptrStrL7AccSetting(s string) *string {
	return &s
}

func ptrInt64L7AccSetting(i int64) *int64 {
	return &i
}

func ptrUint64L7AccSetting(i uint64) *uint64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestL7AccSettingZoneSetting" -v -count=1 -gcflags="all=-l"

// TestL7AccSettingZoneSetting_Read_Success tests Read populates zone_setting computed attribute
func TestL7AccSettingZoneSetting_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccSetting().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccSetting", func(request *teov20220901.DescribeL7AccSettingRequest) (*teov20220901.DescribeL7AccSettingResponse, error) {
		resp := teov20220901.NewDescribeL7AccSettingResponse()
		resp.Response = &teov20220901.DescribeL7AccSettingResponseParams{
			ZoneSetting: &teov20220901.ZoneConfigParameters{
				ZoneName: ptrStrL7AccSetting("test-zone.example.com"),
				ZoneConfig: &teov20220901.ZoneConfig{
					SmartRouting: &teov20220901.SmartRoutingParameters{
						Switch: ptrStrL7AccSetting("on"),
					},
					Cache: &teov20220901.CacheConfigParameters{
						FollowOrigin: &teov20220901.FollowOrigin{
							Switch:               ptrStrL7AccSetting("on"),
							DefaultCache:         ptrStrL7AccSetting("off"),
							DefaultCacheStrategy: ptrStrL7AccSetting("on"),
							DefaultCacheTime:     ptrInt64L7AccSetting(0),
						},
						NoCache: &teov20220901.NoCache{
							Switch: ptrStrL7AccSetting("off"),
						},
						CustomTime: &teov20220901.CacheConfigCustomTime{
							Switch:    ptrStrL7AccSetting("off"),
							CacheTime: ptrInt64L7AccSetting(2592000),
						},
					},
					MaxAge: &teov20220901.MaxAgeParameters{
						FollowOrigin: ptrStrL7AccSetting("on"),
						CacheTime:    ptrInt64L7AccSetting(600),
					},
					CacheKey: &teov20220901.CacheKeyConfigParameters{
						FullURLCache: ptrStrL7AccSetting("on"),
						IgnoreCase:   ptrStrL7AccSetting("off"),
						QueryString: &teov20220901.CacheKeyQueryString{
							Switch: ptrStrL7AccSetting("off"),
							Action: ptrStrL7AccSetting("includeCustom"),
							Values: []*string{ptrStrL7AccSetting("key1")},
						},
					},
					CachePrefresh: &teov20220901.CachePrefreshParameters{
						Switch:           ptrStrL7AccSetting("off"),
						CacheTimePercent: ptrInt64L7AccSetting(90),
					},
					OfflineCache: &teov20220901.OfflineCacheParameters{
						Switch: ptrStrL7AccSetting("on"),
					},
					Compression: &teov20220901.CompressionParameters{
						Switch:     ptrStrL7AccSetting("on"),
						Algorithms: []*string{ptrStrL7AccSetting("brotli"), ptrStrL7AccSetting("gzip")},
					},
					ForceRedirectHTTPS: &teov20220901.ForceRedirectHTTPSParameters{
						Switch:             ptrStrL7AccSetting("off"),
						RedirectStatusCode: ptrInt64L7AccSetting(302),
					},
					HSTS: &teov20220901.HSTSParameters{
						Switch:            ptrStrL7AccSetting("off"),
						Timeout:           ptrInt64L7AccSetting(0),
						IncludeSubDomains: ptrStrL7AccSetting("off"),
						Preload:           ptrStrL7AccSetting("off"),
					},
					TLSConfig: &teov20220901.TLSConfigParameters{
						Version:     []*string{ptrStrL7AccSetting("TLSv1"), ptrStrL7AccSetting("TLSv1.1"), ptrStrL7AccSetting("TLSv1.2"), ptrStrL7AccSetting("TLSv1.3")},
						CipherSuite: ptrStrL7AccSetting("loose-v2023"),
					},
					OCSPStapling: &teov20220901.OCSPStaplingParameters{
						Switch: ptrStrL7AccSetting("off"),
					},
					HTTP2: &teov20220901.HTTP2Parameters{
						Switch: ptrStrL7AccSetting("off"),
					},
					QUIC: &teov20220901.QUICParameters{
						Switch: ptrStrL7AccSetting("off"),
					},
					UpstreamHTTP2: &teov20220901.UpstreamHTTP2Parameters{
						Switch: ptrStrL7AccSetting("off"),
					},
					IPv6: &teov20220901.IPv6Parameters{
						Switch: ptrStrL7AccSetting("off"),
					},
					WebSocket: &teov20220901.WebSocketParameters{
						Switch:  ptrStrL7AccSetting("off"),
						Timeout: ptrInt64L7AccSetting(30),
					},
					PostMaxSize: &teov20220901.PostMaxSizeParameters{
						Switch:  ptrStrL7AccSetting("on"),
						MaxSize: ptrInt64L7AccSetting(838860800),
					},
					ClientIPHeader: &teov20220901.ClientIPHeaderParameters{
						Switch:     ptrStrL7AccSetting("off"),
						HeaderName: ptrStrL7AccSetting("X-Forwarded-For"),
					},
					ClientIPCountry: &teov20220901.ClientIPCountryParameters{
						Switch:     ptrStrL7AccSetting("off"),
						HeaderName: ptrStrL7AccSetting("EO-Client-IPCountry"),
					},
					Grpc: &teov20220901.GrpcParameters{
						Switch: ptrStrL7AccSetting("off"),
					},
					AccelerateMainland: &teov20220901.AccelerateMainlandParameters{
						Switch: ptrStrL7AccSetting("off"),
					},
					StandardDebug: &teov20220901.StandardDebugParameters{
						Switch:            ptrStrL7AccSetting("off"),
						AllowClientIPList: []*string{},
						Expires:           ptrStrL7AccSetting("1969-12-31T16:00:00Z"),
					},
				},
			},
			RequestId: ptrStrL7AccSetting("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForL7AccSetting()
	res := teo.ResourceTencentCloudTeoL7AccSetting()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-36bjhygh1bxe",
	})
	d.SetId("zone-36bjhygh1bxe")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify existing top-level attributes are still populated
	assert.Equal(t, "test-zone.example.com", d.Get("zone_name"))

	// Verify zone_setting computed attribute is populated
	zoneSetting := d.Get("zone_setting").([]interface{})
	assert.Equal(t, 1, len(zoneSetting), "zone_setting should have 1 element")

	zoneSettingMap := zoneSetting[0].(map[string]interface{})
	assert.Equal(t, "test-zone.example.com", zoneSettingMap["zone_name"])

	// Verify zone_setting.zone_config is populated
	zoneSettingZoneConfig := zoneSettingMap["zone_config"].([]interface{})
	assert.Equal(t, 1, len(zoneSettingZoneConfig), "zone_setting.zone_config should have 1 element")

	zoneConfigMap := zoneSettingZoneConfig[0].(map[string]interface{})
	assert.Equal(t, "on", zoneConfigMap["smart_routing"].([]interface{})[0].(map[string]interface{})["switch"])
	assert.Equal(t, "off", zoneConfigMap["accelerate_mainland"].([]interface{})[0].(map[string]interface{})["switch"])
	assert.Equal(t, "off", zoneConfigMap["quic"].([]interface{})[0].(map[string]interface{})["switch"])
}

// TestL7AccSettingZoneSetting_Read_APIError tests Read handles API error
func TestL7AccSettingZoneSetting_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccSetting().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccSetting", func(request *teov20220901.DescribeL7AccSettingRequest) (*teov20220901.DescribeL7AccSettingResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMetaForL7AccSetting()
	res := teo.ResourceTencentCloudTeoL7AccSetting()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
	})
	d.SetId("zone-invalid")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestL7AccSettingZoneSetting_Read_EmptyResponse tests Read with nil ZoneSetting
func TestL7AccSettingZoneSetting_Read_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForL7AccSetting().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccSetting", func(request *teov20220901.DescribeL7AccSettingRequest) (*teov20220901.DescribeL7AccSettingResponse, error) {
		resp := teov20220901.NewDescribeL7AccSettingResponse()
		resp.Response = &teov20220901.DescribeL7AccSettingResponseParams{
			ZoneSetting: nil,
			RequestId:   ptrStrL7AccSetting("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForL7AccSetting()
	res := teo.ResourceTencentCloudTeoL7AccSetting()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-36bjhygh1bxe",
	})
	d.SetId("zone-36bjhygh1bxe")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// Resource should be marked as removed when response is nil
	assert.Equal(t, "", d.Id())
}

func TestAccTencentCloudTeoL7AccSettingResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7AccSetting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_id", "zone-36bjhygh1bxe"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.accelerate_mainland.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.cache_time", "2592000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_strategy", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.no_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.action", "includeCustom"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.cache_time_percent", "90"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_country.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_header.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.algorithms.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.redirect_status_code", "302"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.grpc.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.include_sub_domains", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.preload", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.timeout", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ipv6.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.cache_time", "600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.follow_origin", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ocsp_stapling.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.offline_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.max_size", "838860800"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.quic.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.smart_routing.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.expires", "1969-12-31T16:00:00Z"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.cipher_suite", "loose-v2023"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.version.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.upstream_http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.timeout", "30"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoL7AccSettingUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_id", "zone-36bjhygh1bxe"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.accelerate_mainland.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.cache_time", "2592000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_strategy", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.no_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.action", "includeCustom"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.cache_time_percent", "90"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_country.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_header.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.algorithms.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.redirect_status_code", "302"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.grpc.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.include_sub_domains", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.preload", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.timeout", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ipv6.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.cache_time", "600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.follow_origin", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ocsp_stapling.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.offline_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.max_size", "838860800"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.quic.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.smart_routing.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.expires", "1969-12-31T16:00:00Z"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.cipher_suite", "loose-v2023"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.version.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.upstream_http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.timeout", "30"),
				),
			},
			{
				Config: testAccTeoL7AccSetting,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_id", "zone-36bjhygh1bxe"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.accelerate_mainland.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.cache_time", "2592000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.custom_time.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_strategy", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.default_cache_time", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.follow_origin.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache.0.no_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.action", "includeCustom"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_key.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.cache_time_percent", "90"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.cache_prefresh.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_country.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.client_ip_header.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.algorithms.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.compression.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.redirect_status_code", "302"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.force_redirect_https.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.grpc.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.include_sub_domains", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.preload", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.hsts.0.timeout", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ipv6.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.cache_time", "600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.max_age.0.follow_origin", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.ocsp_stapling.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.offline_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.max_size", "838860800"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.post_max_size.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.quic.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.smart_routing.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.standard_debug.0.expires", "1969-12-31T16:00:00Z"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.cipher_suite", "loose-v2023"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.tls_config.0.version.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.upstream_http2.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_setting.teo_l7_acc_setting", "zone_config.0.web_socket.0.timeout", "30"),
				),
			},
		},
	})
}

const testAccTeoL7AccSetting = `

resource "tencentcloud_teo_l7_acc_setting" "teo_l7_acc_setting" {
  zone_id = "zone-36bjhygh1bxe"
  zone_config {
    accelerate_mainland {
      switch = "off"
    }
    cache {
      custom_time {
        cache_time = 2592000
        switch     = "off"
      }
      follow_origin {
        default_cache          = "off"
        default_cache_strategy = "on"
        default_cache_time     = 0
        switch                 = "on"
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
      }
    }
    cache_prefresh {
      cache_time_percent = 90
      switch             = "off"
    }
    client_ip_country {
      switch      = "off"
    }
    client_ip_header {
      switch      = "off"
    }
    compression {
      algorithms = ["brotli", "gzip"]
      switch     = "on"
    }
    force_redirect_https {
      redirect_status_code = 302
      switch               = "off"
    }
    grpc {
      switch = "off"
    }
    hsts {
      include_sub_domains = "off"
      preload             = "off"
      switch              = "off"
      timeout             = 0
    }
    http2 {
      switch = "off"
    }
    ipv6 {
      switch = "off"
    }
    max_age {
      cache_time    = 600
      follow_origin = "on"
    }
    ocsp_stapling {
      switch = "off"
    }
    offline_cache {
      switch = "on"
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
    standard_debug {
      allow_client_ip_list = []
      expires              = "1969-12-31T16:00:00Z"
      switch               = "off"
    }
    tls_config {
      cipher_suite = "loose-v2023"
      version      = ["TLSv1", "TLSv1.1", "TLSv1.2", "TLSv1.3"]
    }
    upstream_http2 {
      switch = "off"
    }
    web_socket {
      switch  = "off"
      timeout = 30
    }
  }
}
`

const testAccTeoL7AccSettingUp = `

resource "tencentcloud_teo_l7_acc_setting" "teo_l7_acc_setting" {
  zone_id = "zone-36bjhygh1bxe"
  zone_config {
    accelerate_mainland {
      switch = "on"
    }
    cache {
      custom_time {
        cache_time = 2592000
        switch     = "off"
      }
      follow_origin {
        default_cache          = "off"
        default_cache_strategy = "on"
        default_cache_time     = 0
        switch                 = "on"
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
      }
    }
    cache_prefresh {
      cache_time_percent = 90
      switch             = "off"
    }
    client_ip_country {
      switch      = "off"
    }
    client_ip_header {
      switch      = "off"
    }
    compression {
      algorithms = ["brotli", "gzip"]
      switch     = "on"
    }
    force_redirect_https {
      redirect_status_code = 302
      switch               = "off"
    }
    grpc {
      switch = "off"
    }
    hsts {
      include_sub_domains = "off"
      preload             = "off"
      switch              = "off"
      timeout             = 0
    }
    http2 {
      switch = "off"
    }
    ipv6 {
      switch = "off"
    }
    max_age {
      cache_time    = 600
      follow_origin = "on"
    }
    ocsp_stapling {
      switch = "off"
    }
    offline_cache {
      switch = "on"
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
    standard_debug {
      allow_client_ip_list = []
      expires              = "1969-12-31T16:00:00Z"
      switch               = "off"
    }
    tls_config {
      cipher_suite = "loose-v2023"
      version      = ["TLSv1", "TLSv1.1", "TLSv1.2", "TLSv1.3"]
    }
    upstream_http2 {
      switch = "off"
    }
    web_socket {
      switch  = "off"
      timeout = 30
    }
  }
}
`
