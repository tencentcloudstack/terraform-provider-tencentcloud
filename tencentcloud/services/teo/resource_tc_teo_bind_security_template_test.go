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

// go test ./tencentcloud/services/teo/ -run "TestTeoBindSecurityTemplate_ReadSuccess" -v -count=1 -gcflags="all=-l"
// TestTeoBindSecurityTemplate_ReadSuccess tests Read uses DescribeSecurityTemplateBindings to find the binding and sets status and message.
func TestTeoBindSecurityTemplate_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindSecTpl().client, "UseTeoV20220901Client", teoClient)

	// Mock DescribeSecurityTemplateBindings to return a binding with message.
	patches.ApplyMethodFunc(teoClient, "DescribeSecurityTemplateBindings", func(request *teov20220901.DescribeSecurityTemplateBindingsRequest) (*teov20220901.DescribeSecurityTemplateBindingsResponse, error) {
		resp := teov20220901.NewDescribeSecurityTemplateBindingsResponse()
		resp.Response = &teov20220901.DescribeSecurityTemplateBindingsResponseParams{
			SecurityTemplate: []*teov20220901.SecurityTemplateBinding{
				{
					TemplateId: ptrBindSecTplString("temp-7dr7dm78"),
					TemplateScope: []*teov20220901.TemplateScope{
						{
							ZoneId: ptrBindSecTplString("zone-39quuimqg8r6"),
							EntityStatus: []*teov20220901.EntityStatus{
								{
									Entity:  ptrBindSecTplString("aaa.makn.cn"),
									Status:  ptrBindSecTplString("online"),
									Message: ptrBindSecTplString("config applied successfully"),
								},
							},
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
	assert.Equal(t, "config applied successfully", d.Get("message"))
}

// TestTeoBindSecurityTemplate_ReadSuccessNoMessage tests Read sets status but message is nil in API response.
func TestTeoBindSecurityTemplate_ReadSuccessNoMessage(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindSecTpl().client, "UseTeoV20220901Client", teoClient)

	// Mock DescribeSecurityTemplateBindings to return a binding without message.
	patches.ApplyMethodFunc(teoClient, "DescribeSecurityTemplateBindings", func(request *teov20220901.DescribeSecurityTemplateBindingsRequest) (*teov20220901.DescribeSecurityTemplateBindingsResponse, error) {
		resp := teov20220901.NewDescribeSecurityTemplateBindingsResponse()
		resp.Response = &teov20220901.DescribeSecurityTemplateBindingsResponseParams{
			SecurityTemplate: []*teov20220901.SecurityTemplateBinding{
				{
					TemplateId: ptrBindSecTplString("temp-7dr7dm78"),
					TemplateScope: []*teov20220901.TemplateScope{
						{
							ZoneId: ptrBindSecTplString("zone-39quuimqg8r6"),
							EntityStatus: []*teov20220901.EntityStatus{
								{
									Entity:  ptrBindSecTplString("aaa.makn.cn"),
									Status:  ptrBindSecTplString("online"),
									Message: nil,
								},
							},
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

	assert.Equal(t, "zone-39quuimqg8r6", d.Get("zone_id"))
	assert.Equal(t, "temp-7dr7dm78", d.Get("template_id"))
	assert.Equal(t, "aaa.makn.cn", d.Get("entity"))
	assert.Equal(t, "online", d.Get("status"))
	// message should be empty string (zero value) since API returned nil
	assert.Equal(t, "", d.Get("message"))
}

// TestTeoBindSecurityTemplate_ReadNotFound tests Read clears id when binding not found.
func TestTeoBindSecurityTemplate_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindSecTpl().client, "UseTeoV20220901Client", teoClient)

	// Mock DescribeSecurityTemplateBindings to return a template that does not bind to the entity.
	patches.ApplyMethodFunc(teoClient, "DescribeSecurityTemplateBindings", func(request *teov20220901.DescribeSecurityTemplateBindingsRequest) (*teov20220901.DescribeSecurityTemplateBindingsResponse, error) {
		resp := teov20220901.NewDescribeSecurityTemplateBindingsResponse()
		resp.Response = &teov20220901.DescribeSecurityTemplateBindingsResponseParams{
			SecurityTemplate: []*teov20220901.SecurityTemplateBinding{
				{
					TemplateId: ptrBindSecTplString("temp-7dr7dm78"),
					TemplateScope: []*teov20220901.TemplateScope{
						{
							ZoneId: ptrBindSecTplString("zone-39quuimqg8r6"),
							EntityStatus: []*teov20220901.EntityStatus{
								{
									Entity: ptrBindSecTplString("other.makn.cn"),
									Status: ptrBindSecTplString("online"),
								},
							},
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

// TestTeoBindSecurityTemplate_ReadEmptySecurityTemplate tests Read clears id when API returns empty SecurityTemplate list.
func TestTeoBindSecurityTemplate_ReadEmptySecurityTemplate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindSecTpl().client, "UseTeoV20220901Client", teoClient)

	// Mock DescribeSecurityTemplateBindings to return empty SecurityTemplate list.
	patches.ApplyMethodFunc(teoClient, "DescribeSecurityTemplateBindings", func(request *teov20220901.DescribeSecurityTemplateBindingsRequest) (*teov20220901.DescribeSecurityTemplateBindingsResponse, error) {
		resp := teov20220901.NewDescribeSecurityTemplateBindingsResponse()
		resp.Response = &teov20220901.DescribeSecurityTemplateBindingsResponseParams{
			SecurityTemplate: []*teov20220901.SecurityTemplateBinding{},
			RequestId:        ptrBindSecTplString("fake-request-id"),
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

// TestTeoBindSecurityTemplate_ReadWithMessage tests the message attribute is set correctly during Read.
func TestTeoBindSecurityTemplate_ReadWithMessage(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaBindSecTpl().client, "UseTeoV20220901Client", teoClient)

	// Mock DescribeSecurityTemplateBindings to return a binding with a failure message.
	patches.ApplyMethodFunc(teoClient, "DescribeSecurityTemplateBindings", func(request *teov20220901.DescribeSecurityTemplateBindingsRequest) (*teov20220901.DescribeSecurityTemplateBindingsResponse, error) {
		resp := teov20220901.NewDescribeSecurityTemplateBindingsResponse()
		resp.Response = &teov20220901.DescribeSecurityTemplateBindingsResponseParams{
			SecurityTemplate: []*teov20220901.SecurityTemplateBinding{
				{
					TemplateId: ptrBindSecTplString("temp-7dr7dm78"),
					TemplateScope: []*teov20220901.TemplateScope{
						{
							ZoneId: ptrBindSecTplString("zone-39quuimqg8r6"),
							EntityStatus: []*teov20220901.EntityStatus{
								{
									Entity:  ptrBindSecTplString("aaa.makn.cn"),
									Status:  ptrBindSecTplString("fail"),
									Message: ptrBindSecTplString("domain not found in zone"),
								},
							},
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

	assert.Equal(t, "fail", d.Get("status"))
	assert.Equal(t, "domain not found in zone", d.Get("message"))
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
	assert.Contains(t, res.Schema, "message")

	assert.True(t, res.Schema["zone_id"].Required)
	assert.True(t, res.Schema["entity"].Required)
	assert.True(t, res.Schema["template_id"].Required)

	assert.True(t, res.Schema["status"].Computed)
	assert.True(t, res.Schema["message"].Computed)
}
