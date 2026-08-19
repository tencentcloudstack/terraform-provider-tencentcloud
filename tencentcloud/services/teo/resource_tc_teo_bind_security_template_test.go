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

func TestAccTencentCloudTeoBindSecurityTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoBindSecurityTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_bind_security_template.teo_bind_security_template", "zone_id", "zone-39quuimqg8r6"),
					resource.TestCheckResourceAttr("tencentcloud_teo_bind_security_template.teo_bind_security_template", "template_id", "temp-7dr7dm78"),
					resource.TestCheckResourceAttr("tencentcloud_teo_bind_security_template.teo_bind_security_template", "entity", "aaa.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_bind_security_template.teo_bind_security_template", "operate", "unbind-use-default"),
					resource.TestCheckResourceAttr("tencentcloud_teo_bind_security_template.teo_bind_security_template", "status", "online"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_bind_security_template.teo_bind_security_template",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"operate",
				},
			},
		},
	})
}

const testAccTeoBindSecurityTemplate = `

resource "tencentcloud_teo_bind_security_template" "teo_bind_security_template" {
  operate     = "unbind-use-default"
  template_id = "temp-7dr7dm78"
  zone_id     = "zone-39quuimqg8r6"
  entity 	  = "aaa.makn.cn"
  over_write  = false
}

`

// mockMetaBindSecTpl implements tccommon.ProviderMeta
type mockMetaBindSecTpl struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBindSecTpl) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBindSecTpl{}

func newMockMetaBindSecTpl() *mockMetaBindSecTpl {
	return &mockMetaBindSecTpl{client: &connectivity.TencentCloudClient{}}
}

func ptrBindSecTplString(s string) *string {
	return &s
}

func ptrBindSecTplInt64(n int64) *int64 {
	return &n
}

// go test ./tencentcloud/services/teo/ -run "TestTeoBindSecurityTemplate_ReadSuccess" -v -count=1 -gcflags="all=-l"
// TestTeoBindSecurityTemplate_ReadSuccess tests Read uses DescribeZones + DescribeWebSecurityTemplates to find the binding and sets status.
func TestTeoBindSecurityTemplate_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindSecTpl().client, "UseTeoV20220901Client", teoClient)

	// Mock DescribeZones to return one zone.
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrBindSecTplInt64(1),
			Zones: []*teov20220901.Zone{
				{
					ZoneId: ptrBindSecTplString("zone-39quuimqg8r6"),
				},
			},
			RequestId: ptrBindSecTplString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeWebSecurityTemplates to return a template bound to the entity.
	patches.ApplyMethodFunc(teoClient, "DescribeWebSecurityTemplates", func(request *teov20220901.DescribeWebSecurityTemplatesRequest) (*teov20220901.DescribeWebSecurityTemplatesResponse, error) {
		resp := teov20220901.NewDescribeWebSecurityTemplatesResponse()
		resp.Response = &teov20220901.DescribeWebSecurityTemplatesResponseParams{
			TotalCount: ptrBindSecTplInt64(1),
			SecurityPolicyTemplates: []*teov20220901.SecurityPolicyTemplateInfo{
				{
					ZoneId:       ptrBindSecTplString("zone-39quuimqg8r6"),
					TemplateId:   ptrBindSecTplString("temp-7dr7dm78"),
					TemplateName: ptrBindSecTplString("default template"),
					BindDomains: []*teov20220901.BindDomainInfo{
						{
							Domain: ptrBindSecTplString("aaa.makn.cn"),
							ZoneId: ptrBindSecTplString("zone-39quuimqg8r6"),
							Status: ptrBindSecTplString("online"),
						},
					},
				},
			},
			RequestId: ptrBindSecTplString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBindSecTpl()
	res := teo.ResourceTencentCloudTeoBindSecurityTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-39quuimqg8r6",
		"template_id": "temp-7dr7dm78",
		"entity":      "aaa.makn.cn",
	})
	d.SetId("zone-39quuimqg8r6#temp-7dr7dm78#aaa.makn.cn")

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), "test-log-id", d, meta)
	_ = ctx

	err := res.Read(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, "zone-39quuimqg8r6", d.Get("zone_id"))
	assert.Equal(t, "temp-7dr7dm78", d.Get("template_id"))
	assert.Equal(t, "aaa.makn.cn", d.Get("entity"))
	assert.Equal(t, "online", d.Get("status"))
}

// TestTeoBindSecurityTemplate_ReadNotFound tests Read clears id when binding not found.
func TestTeoBindSecurityTemplate_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindSecTpl().client, "UseTeoV20220901Client", teoClient)

	// Mock DescribeZones to return one zone.
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrBindSecTplInt64(1),
			Zones: []*teov20220901.Zone{
				{
					ZoneId: ptrBindSecTplString("zone-39quuimqg8r6"),
				},
			},
			RequestId: ptrBindSecTplString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeWebSecurityTemplates to return a template that does not bind to the entity.
	patches.ApplyMethodFunc(teoClient, "DescribeWebSecurityTemplates", func(request *teov20220901.DescribeWebSecurityTemplatesRequest) (*teov20220901.DescribeWebSecurityTemplatesResponse, error) {
		resp := teov20220901.NewDescribeWebSecurityTemplatesResponse()
		resp.Response = &teov20220901.DescribeWebSecurityTemplatesResponseParams{
			TotalCount: ptrBindSecTplInt64(1),
			SecurityPolicyTemplates: []*teov20220901.SecurityPolicyTemplateInfo{
				{
					ZoneId:       ptrBindSecTplString("zone-39quuimqg8r6"),
					TemplateId:   ptrBindSecTplString("temp-7dr7dm78"),
					TemplateName: ptrBindSecTplString("default template"),
					BindDomains: []*teov20220901.BindDomainInfo{
						{
							Domain: ptrBindSecTplString("other.makn.cn"),
							ZoneId: ptrBindSecTplString("zone-39quuimqg8r6"),
							Status: ptrBindSecTplString("online"),
						},
					},
				},
			},
			RequestId: ptrBindSecTplString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBindSecTpl()
	res := teo.ResourceTencentCloudTeoBindSecurityTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-39quuimqg8r6",
		"template_id": "temp-7dr7dm78",
		"entity":      "aaa.makn.cn",
	})
	d.SetId("zone-39quuimqg8r6#temp-7dr7dm78#aaa.makn.cn")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoBindSecurityTemplate_ReadNoZone tests Read clears id when no zone exists.
func TestTeoBindSecurityTemplate_ReadNoZone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindSecTpl().client, "UseTeoV20220901Client", teoClient)

	// Mock DescribeZones to return no zone.
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrBindSecTplInt64(0),
			Zones:      []*teov20220901.Zone{},
			RequestId:  ptrBindSecTplString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBindSecTpl()
	res := teo.ResourceTencentCloudTeoBindSecurityTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-39quuimqg8r6",
		"template_id": "temp-7dr7dm78",
		"entity":      "aaa.makn.cn",
	})
	d.SetId("zone-39quuimqg8r6#temp-7dr7dm78#aaa.makn.cn")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoBindSecurityTemplate_ReadBatchZones tests Read correctly handles more than 100 zone ids by batching DescribeWebSecurityTemplates calls.
func TestTeoBindSecurityTemplate_ReadBatchZones(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindSecTpl().client, "UseTeoV20220901Client", teoClient)

	// Build 150 zones; the target zone is placed in the second batch (index 120) to verify batching works.
	zones := make([]*teov20220901.Zone, 0, 150)
	targetZoneId := "zone-target-120"
	for i := 0; i < 150; i++ {
		zid := "zone-" + itoaBindSecTpl(i)
		if i == 120 {
			zid = targetZoneId
		}
		zones = append(zones, &teov20220901.Zone{
			ZoneId: ptrBindSecTplString(zid),
		})
	}

	callCount := 0
	// Mock DescribeZones to return 150 zones in two pages (100 + 50).
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		offset := 0
		if request.Offset != nil {
			offset = int(*request.Offset)
		}
		limit := 100
		if request.Limit != nil {
			limit = int(*request.Limit)
		}
		end := offset + limit
		if end > len(zones) {
			end = len(zones)
		}
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrBindSecTplInt64(int64(len(zones))),
			Zones:      zones[offset:end],
			RequestId:  ptrBindSecTplString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeWebSecurityTemplates: only the second batch (which contains the target zone) returns the binding.
	patches.ApplyMethodFunc(teoClient, "DescribeWebSecurityTemplates", func(request *teov20220901.DescribeWebSecurityTemplatesRequest) (*teov20220901.DescribeWebSecurityTemplatesResponse, error) {
		callCount++
		for _, zid := range request.ZoneIds {
			if zid != nil && *zid == targetZoneId {
				resp := teov20220901.NewDescribeWebSecurityTemplatesResponse()
				resp.Response = &teov20220901.DescribeWebSecurityTemplatesResponseParams{
					TotalCount: ptrBindSecTplInt64(1),
					SecurityPolicyTemplates: []*teov20220901.SecurityPolicyTemplateInfo{
						{
							ZoneId:       ptrBindSecTplString(targetZoneId),
							TemplateId:   ptrBindSecTplString("temp-7dr7dm78"),
							TemplateName: ptrBindSecTplString("default template"),
							BindDomains: []*teov20220901.BindDomainInfo{
								{
									Domain: ptrBindSecTplString("aaa.makn.cn"),
									ZoneId: ptrBindSecTplString(targetZoneId),
									Status: ptrBindSecTplString("online"),
								},
							},
						},
					},
					RequestId: ptrBindSecTplString("fake-request-id"),
				}
				return resp, nil
			}
		}
		resp := teov20220901.NewDescribeWebSecurityTemplatesResponse()
		resp.Response = &teov20220901.DescribeWebSecurityTemplatesResponseParams{
			TotalCount:              ptrBindSecTplInt64(0),
			SecurityPolicyTemplates: []*teov20220901.SecurityPolicyTemplateInfo{},
			RequestId:               ptrBindSecTplString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBindSecTpl()
	res := teo.ResourceTencentCloudTeoBindSecurityTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     targetZoneId,
		"template_id": "temp-7dr7dm78",
		"entity":      "aaa.makn.cn",
	})
	d.SetId(targetZoneId + "#temp-7dr7dm78#aaa.makn.cn")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "online", d.Get("status"))
	assert.GreaterOrEqual(t, callCount, 2, "DescribeWebSecurityTemplates should be called at least twice for 150 zones")
}

// itoaBindSecTpl converts an int to its decimal string representation without using strconv to keep the test self-contained.
func itoaBindSecTpl(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

// TestTeoBindSecurityTemplate_Schema validates schema definition.
func TestTeoBindSecurityTemplate_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoBindSecurityTemplate()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "entity")
	assert.Contains(t, res.Schema, "template_id")
	assert.Contains(t, res.Schema, "operate")
	assert.Contains(t, res.Schema, "over_write")
	assert.Contains(t, res.Schema, "status")

	assert.True(t, res.Schema["zone_id"].Required)
	assert.True(t, res.Schema["entity"].Required)
	assert.True(t, res.Schema["template_id"].Required)

	assert.True(t, res.Schema["status"].Computed)
}
