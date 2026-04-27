package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// mockMeta implements tccommon.ProviderMeta
type mockMetaApplicationProxy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaApplicationProxy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaApplicationProxy{}

func newMockMetaApplicationProxy() *mockMetaApplicationProxy {
	return &mockMetaApplicationProxy{client: &connectivity.TencentCloudClient{}}
}

func ptrStringApplicationProxy(s string) *string {
	return &s
}

func ptrInt64ApplicationProxy(v int64) *int64 {
	return &v
}

func ptrUint64ApplicationProxy(v uint64) *uint64 {
	return &v
}

// go test ./tencentcloud/services/teo/ -run "TestApplicationProxy_Delete" -v -count=1 -gcflags="all=-l"

// TestApplicationProxy_Delete_Success_AlreadyOffline tests Delete when proxy is already offline
func TestApplicationProxy_Delete_Success_AlreadyOffline(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaApplicationProxy().client, "UseTeoClient", teoClient)

	// Mock DescribeApplicationProxies for the read function - return offline status
	patches.ApplyMethodFunc(teoClient, "DescribeApplicationProxies", func(request *teov20220901.DescribeApplicationProxiesRequest) (*teov20220901.DescribeApplicationProxiesResponse, error) {
		resp := teov20220901.NewDescribeApplicationProxiesResponse()
		resp.Response = &teov20220901.DescribeApplicationProxiesResponseParams{
			ApplicationProxies: []*teov20220901.ApplicationProxy{
				{
					ZoneId:             ptrStringApplicationProxy("zone-1234567890"),
					ProxyId:            ptrStringApplicationProxy("proxy-abcdefghij"),
					ProxyName:          ptrStringApplicationProxy("test-proxy"),
					ProxyType:          ptrStringApplicationProxy("instance"),
					PlatType:           ptrStringApplicationProxy("domain"),
					SecurityType:       ptrInt64ApplicationProxy(1),
					AccelerateType:     ptrInt64ApplicationProxy(1),
					SessionPersistTime: ptrUint64ApplicationProxy(600),
					Status:             ptrStringApplicationProxy("offline"),
				},
			},
			RequestId: ptrStringApplicationProxy("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DeleteApplicationProxy
	patches.ApplyMethodFunc(teoClient, "DeleteApplicationProxy", func(request *teov20220901.DeleteApplicationProxyRequest) (*teov20220901.DeleteApplicationProxyResponse, error) {
		resp := teov20220901.NewDeleteApplicationProxyResponse()
		resp.Response = &teov20220901.DeleteApplicationProxyResponseParams{
			RequestId: ptrStringApplicationProxy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaApplicationProxy()
	res := svcteo.ResourceTencentCloudTeoApplicationProxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":              "zone-1234567890",
		"proxy_id":             "proxy-abcdefghij",
		"proxy_name":           "test-proxy",
		"proxy_type":           "instance",
		"plat_type":            "domain",
		"security_type":        1,
		"accelerate_type":      1,
		"session_persist_time": 600,
	})
	d.SetId("zone-1234567890#proxy-abcdefghij")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestApplicationProxy_Delete_Success_SetOffline tests Delete when proxy needs to be set offline first
func TestApplicationProxy_Delete_Success_SetOffline(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaApplicationProxy().client, "UseTeoClient", teoClient)

	callCount := 0
	// Mock DescribeApplicationProxies - first call returns online, second returns offline
	patches.ApplyMethodFunc(teoClient, "DescribeApplicationProxies", func(request *teov20220901.DescribeApplicationProxiesRequest) (*teov20220901.DescribeApplicationProxiesResponse, error) {
		callCount++
		resp := teov20220901.NewDescribeApplicationProxiesResponse()
		status := "offline"
		if callCount == 1 {
			status = "online"
		}
		resp.Response = &teov20220901.DescribeApplicationProxiesResponseParams{
			ApplicationProxies: []*teov20220901.ApplicationProxy{
				{
					ZoneId:             ptrStringApplicationProxy("zone-1234567890"),
					ProxyId:            ptrStringApplicationProxy("proxy-abcdefghij"),
					ProxyName:          ptrStringApplicationProxy("test-proxy"),
					ProxyType:          ptrStringApplicationProxy("instance"),
					PlatType:           ptrStringApplicationProxy("domain"),
					SecurityType:       ptrInt64ApplicationProxy(1),
					AccelerateType:     ptrInt64ApplicationProxy(1),
					SessionPersistTime: ptrUint64ApplicationProxy(600),
					Status:             ptrStringApplicationProxy(status),
				},
			},
			RequestId: ptrStringApplicationProxy("fake-request-id"),
		}
		return resp, nil
	})

	// Mock ModifyApplicationProxyStatus
	patches.ApplyMethodFunc(teoClient, "ModifyApplicationProxyStatus", func(request *teov20220901.ModifyApplicationProxyStatusRequest) (*teov20220901.ModifyApplicationProxyStatusResponse, error) {
		resp := teov20220901.NewModifyApplicationProxyStatusResponse()
		resp.Response = &teov20220901.ModifyApplicationProxyStatusResponseParams{
			RequestId: ptrStringApplicationProxy("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DeleteApplicationProxy
	patches.ApplyMethodFunc(teoClient, "DeleteApplicationProxy", func(request *teov20220901.DeleteApplicationProxyRequest) (*teov20220901.DeleteApplicationProxyResponse, error) {
		resp := teov20220901.NewDeleteApplicationProxyResponse()
		resp.Response = &teov20220901.DeleteApplicationProxyResponseParams{
			RequestId: ptrStringApplicationProxy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaApplicationProxy()
	res := svcteo.ResourceTencentCloudTeoApplicationProxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":              "zone-1234567890",
		"proxy_id":             "proxy-abcdefghij",
		"proxy_name":           "test-proxy",
		"proxy_type":           "instance",
		"plat_type":            "domain",
		"security_type":        1,
		"accelerate_type":      1,
		"session_persist_time": 600,
	})
	d.SetId("zone-1234567890#proxy-abcdefghij")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestApplicationProxy_Delete_APIError tests Delete handles API error
func TestApplicationProxy_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaApplicationProxy().client, "UseTeoClient", teoClient)

	// Mock DescribeApplicationProxies - return offline status so we skip ModifyApplicationProxyStatus
	patches.ApplyMethodFunc(teoClient, "DescribeApplicationProxies", func(request *teov20220901.DescribeApplicationProxiesRequest) (*teov20220901.DescribeApplicationProxiesResponse, error) {
		resp := teov20220901.NewDescribeApplicationProxiesResponse()
		resp.Response = &teov20220901.DescribeApplicationProxiesResponseParams{
			ApplicationProxies: []*teov20220901.ApplicationProxy{
				{
					ZoneId:             ptrStringApplicationProxy("zone-1234567890"),
					ProxyId:            ptrStringApplicationProxy("proxy-abcdefghij"),
					ProxyName:          ptrStringApplicationProxy("test-proxy"),
					ProxyType:          ptrStringApplicationProxy("instance"),
					PlatType:           ptrStringApplicationProxy("domain"),
					SecurityType:       ptrInt64ApplicationProxy(1),
					AccelerateType:     ptrInt64ApplicationProxy(1),
					SessionPersistTime: ptrUint64ApplicationProxy(600),
					Status:             ptrStringApplicationProxy("offline"),
				},
			},
			RequestId: ptrStringApplicationProxy("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DeleteApplicationProxy to return error
	patches.ApplyMethodFunc(teoClient, "DeleteApplicationProxy", func(request *teov20220901.DeleteApplicationProxyRequest) (*teov20220901.DeleteApplicationProxyResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Proxy not found")
	})

	meta := newMockMetaApplicationProxy()
	res := svcteo.ResourceTencentCloudTeoApplicationProxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":              "zone-1234567890",
		"proxy_id":             "proxy-abcdefghij",
		"proxy_name":           "test-proxy",
		"proxy_type":           "instance",
		"plat_type":            "domain",
		"security_type":        1,
		"accelerate_type":      1,
		"session_persist_time": 600,
	})
	d.SetId("zone-1234567890#proxy-abcdefghij")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestApplicationProxy_Delete_VerifyDGetParams verifies that zone_id and proxy_id from d.Get() are correctly passed to DeleteApplicationProxyRequest
func TestApplicationProxy_Delete_VerifyDGetParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaApplicationProxy().client, "UseTeoClient", teoClient)

	// Mock DescribeApplicationProxies for the read function - return offline status
	patches.ApplyMethodFunc(teoClient, "DescribeApplicationProxies", func(request *teov20220901.DescribeApplicationProxiesRequest) (*teov20220901.DescribeApplicationProxiesResponse, error) {
		resp := teov20220901.NewDescribeApplicationProxiesResponse()
		resp.Response = &teov20220901.DescribeApplicationProxiesResponseParams{
			ApplicationProxies: []*teov20220901.ApplicationProxy{
				{
					ZoneId:             ptrStringApplicationProxy("zone-from-dget"),
					ProxyId:            ptrStringApplicationProxy("proxy-from-dget"),
					ProxyName:          ptrStringApplicationProxy("test-proxy"),
					ProxyType:          ptrStringApplicationProxy("instance"),
					PlatType:           ptrStringApplicationProxy("domain"),
					SecurityType:       ptrInt64ApplicationProxy(1),
					AccelerateType:     ptrInt64ApplicationProxy(1),
					SessionPersistTime: ptrUint64ApplicationProxy(600),
					Status:             ptrStringApplicationProxy("offline"),
				},
			},
			RequestId: ptrStringApplicationProxy("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DeleteApplicationProxy - verify the request contains correct ZoneId and ProxyId from d.Get()
	patches.ApplyMethodFunc(teoClient, "DeleteApplicationProxy", func(request *teov20220901.DeleteApplicationProxyRequest) (*teov20220901.DeleteApplicationProxyResponse, error) {
		assert.NotNil(t, request.ZoneId, "ZoneId should not be nil")
		assert.NotNil(t, request.ProxyId, "ProxyId should not be nil")
		assert.Equal(t, "zone-from-dget", *request.ZoneId, "ZoneId in DeleteApplicationProxyRequest should match d.Get(\"zone_id\")")
		assert.Equal(t, "proxy-from-dget", *request.ProxyId, "ProxyId in DeleteApplicationProxyRequest should match d.Get(\"proxy_id\")")

		resp := teov20220901.NewDeleteApplicationProxyResponse()
		resp.Response = &teov20220901.DeleteApplicationProxyResponseParams{
			RequestId: ptrStringApplicationProxy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaApplicationProxy()
	res := svcteo.ResourceTencentCloudTeoApplicationProxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":              "zone-from-dget",
		"proxy_id":             "proxy-from-dget",
		"proxy_name":           "test-proxy",
		"proxy_type":           "instance",
		"plat_type":            "domain",
		"security_type":        1,
		"accelerate_type":      1,
		"session_persist_time": 600,
	})
	// Intentionally set a different ID to prove d.Get() is used instead of d.Id()
	d.SetId("zone-from-id#proxy-from-id")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestApplicationProxy_Delete_EmptyZoneId tests Delete returns error when zone_id is empty
func TestApplicationProxy_Delete_EmptyZoneId(t *testing.T) {
	meta := newMockMetaApplicationProxy()
	res := svcteo.ResourceTencentCloudTeoApplicationProxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":              "",
		"proxy_id":             "proxy-abcdefghij",
		"proxy_name":           "test-proxy",
		"proxy_type":           "instance",
		"plat_type":            "domain",
		"security_type":        1,
		"accelerate_type":      1,
		"session_persist_time": 600,
	})
	d.SetId("#proxy-abcdefghij")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zone_id is required")
}

// TestApplicationProxy_Delete_EmptyProxyId tests Delete returns error when proxy_id is empty
func TestApplicationProxy_Delete_EmptyProxyId(t *testing.T) {
	meta := newMockMetaApplicationProxy()
	res := svcteo.ResourceTencentCloudTeoApplicationProxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":              "zone-1234567890",
		"proxy_id":             "",
		"proxy_name":           "test-proxy",
		"proxy_type":           "instance",
		"plat_type":            "domain",
		"security_type":        1,
		"accelerate_type":      1,
		"session_persist_time": 600,
	})
	d.SetId("zone-1234567890#")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "proxy_id is required")
}

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_teo_zone
	resource.AddTestSweepers("tencentcloud_teo_application_proxy", &resource.Sweeper{
		Name: "tencentcloud_teo_application_proxy",
		F:    testSweepApplicationProxy,
	})
}

func testSweepApplicationProxy(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(region)
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svcteo.NewTeoService(client)

	for {
		proxy, err := service.DescribeTeoApplicationProxy(ctx, "", "")
		if err != nil {
			return err
		}

		if proxy == nil {
			return nil
		}

		err = service.DeleteTeoApplicationProxyById(ctx, *proxy.ZoneId, *proxy.ProxyId)
		if err != nil {
			return err
		}
	}
}

// go test -i; go test -test.run TestAccTencentCloudTeoApplicationProxy_basic -v
func TestAccTencentCloudTeoApplicationProxy_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckApplicationProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoApplicationProxy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationProxyExists("tencentcloud_teo_application_proxy.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_application_proxy.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "accelerate_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "security_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "plat_type", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "proxy_name", "test-instance"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "proxy_type", "instance"),
					resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "session_persist_time", "2400"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_application_proxy.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckApplicationProxyDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_application_proxy" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]

		agents, err := service.DescribeTeoApplicationProxy(ctx, zoneId, proxyId)
		if agents != nil {
			return fmt.Errorf("zone ApplicationProxy %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckApplicationProxyExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		agents, err := service.DescribeTeoApplicationProxy(ctx, zoneId, proxyId)
		if agents == nil {
			return fmt.Errorf("zone ApplicationProxy %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoApplicationProxy = testAccTeoZone + `

resource "tencentcloud_teo_application_proxy" "basic" {
  zone_id = tencentcloud_teo_zone.basic.id

  accelerate_type      = 1
  security_type        = 1
  plat_type            = "domain"
  proxy_name           = "test-instance"
  proxy_type           = "instance"
  session_persist_time = 2400
}

`
