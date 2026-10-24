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

func TestAccTencentCloudTeoSecurityClientAttesterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoSecurityClientAttester,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_client_attester.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_client_attester.example", "client_attesters.0.name", "tf-example-attester"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_client_attester.example", "client_attesters.0.attester_source", "TC-RCE"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_client_attester.example", "client_attesters.0.attester_duration", "300s"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_client_attester.example", "client_attesters.0.id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_client_attester.example", "client_attesters.0.type"),
				),
			},
			{
				Config: testAccTeoSecurityClientAttesterUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_security_client_attester.example", "client_attesters.0.name", "tf-example-attester-updated"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_client_attester.example", "client_attesters.0.attester_duration", "600s"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_security_client_attester.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoSecurityClientAttester = `
resource "tencentcloud_teo_security_client_attester" "example" {
  zone_id = "zone-3fkff38fyw8s"

  client_attesters {
    name              = "tf-example"
    attester_source   = "TC-RCE"
    attester_duration = "300s"

    tc_rce_option {
      channel = "12399223"
      region  = "ap-beijing"
    }
  }
}
`

const testAccTeoSecurityClientAttesterUpdate = `
resource "tencentcloud_teo_security_client_attester" "example" {
  zone_id = "zone-3fkff38fyw8s"

  client_attesters {
    name              = "tf-example0-update"
    attester_source   = "TC-RCE"
    attester_duration = "600s"

    tc_rce_option {
      channel = "12399223"
      region  = "ap-beijing"
    }
  }
}
`

// --- gomonkey mock unit tests for new client attester fields ---
// go test ./tencentcloud/services/teo/ -run "TestTeoSecurityClientAttester" -v -count=1 -gcflags="all=-l"

// mockMetaForSecurityClientAttester implements tccommon.ProviderMeta
type mockMetaForSecurityClientAttester struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForSecurityClientAttester) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForSecurityClientAttester{}

func newMockMetaForSecurityClientAttester() *mockMetaForSecurityClientAttester {
	return &mockMetaForSecurityClientAttester{client: &connectivity.TencentCloudClient{}}
}

func ptrStrSecurityClientAttester(s string) *string {
	return &s
}

func ptrInt64SecurityClientAttester(i int64) *int64 {
	return &i
}

// TestTeoSecurityClientAttester_Create_WithUsageLimit tests Create with usage_limit and max_usage_count set
func TestTeoSecurityClientAttester_Create_WithUsageLimit(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityClientAttesterWithContext", func(_ context.Context, request *teov20220901.CreateSecurityClientAttesterRequest) (*teov20220901.CreateSecurityClientAttesterResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Len(t, request.ClientAttesters, 1)
		attester := request.ClientAttesters[0]
		assert.NotNil(t, attester.UsageLimit)
		assert.Equal(t, "on", *attester.UsageLimit)
		assert.NotNil(t, attester.MaxUsageCount)
		assert.Equal(t, int64(100), *attester.MaxUsageCount)
		resp := teov20220901.NewCreateSecurityClientAttesterResponse()
		resp.Response = &teov20220901.CreateSecurityClientAttesterResponseParams{
			ClientAttesterIds: []*string{ptrStrSecurityClientAttester("attest-abcdefgh")},
			RequestId:         ptrStrSecurityClientAttester("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityClientAttester", func(request *teov20220901.DescribeSecurityClientAttesterRequest) (*teov20220901.DescribeSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewDescribeSecurityClientAttesterResponse()
		resp.Response = &teov20220901.DescribeSecurityClientAttesterResponseParams{
			TotalCount: ptrInt64SecurityClientAttester(1),
			ClientAttesters: []*teov20220901.ClientAttester{
				{
					Id:               ptrStrSecurityClientAttester("attest-abcdefgh"),
					Name:             ptrStrSecurityClientAttester("tf-example"),
					Type:             ptrStrSecurityClientAttester("CUSTOM"),
					AttesterSource:   ptrStrSecurityClientAttester("TC-EO-CAPTCHA"),
					AttesterDuration: ptrStrSecurityClientAttester("300s"),
					UsageLimit:       ptrStrSecurityClientAttester("on"),
					MaxUsageCount:    ptrInt64SecurityClientAttester(100),
					TCEOCaptchaOption: &teov20220901.TCEOCaptchaOption{
						CaptchaMode: ptrStrSecurityClientAttester("Invisible"),
					},
				},
			},
			RequestId: ptrStrSecurityClientAttester("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityClientAttester()
	res := teo.ResourceTencentCloudTeoSecurityClientAttester()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":            "tf-example",
				"attester_source": "TC-EO-CAPTCHA",
				"usage_limit":     "on",
				"max_usage_count": 100,
				"tc_eo_captcha_option": []interface{}{
					map[string]interface{}{
						"captcha_mode": "Invisible",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678"+tccommon.FILED_SP+"attest-abcdefgh", d.Id())
	attesters := d.Get("client_attesters").([]interface{})
	assert.Len(t, attesters, 1)
	attesterMap := attesters[0].(map[string]interface{})
	assert.Equal(t, "on", attesterMap["usage_limit"])
	assert.Equal(t, 100, attesterMap["max_usage_count"])
}

// TestTeoSecurityClientAttester_Create_WithMobileCaptchaMode tests Create with mobile_captcha_mode set
func TestTeoSecurityClientAttester_Create_WithMobileCaptchaMode(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityClientAttesterWithContext", func(_ context.Context, request *teov20220901.CreateSecurityClientAttesterRequest) (*teov20220901.CreateSecurityClientAttesterResponse, error) {
		assert.Len(t, request.ClientAttesters, 1)
		attester := request.ClientAttesters[0]
		assert.NotNil(t, attester.TCEOCaptchaOption)
		assert.NotNil(t, attester.TCEOCaptchaOption.MobileCaptchaMode)
		assert.Equal(t, "Visible", *attester.TCEOCaptchaOption.MobileCaptchaMode)
		resp := teov20220901.NewCreateSecurityClientAttesterResponse()
		resp.Response = &teov20220901.CreateSecurityClientAttesterResponseParams{
			ClientAttesterIds: []*string{ptrStrSecurityClientAttester("attest-ijklmnop")},
			RequestId:         ptrStrSecurityClientAttester("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityClientAttester", func(request *teov20220901.DescribeSecurityClientAttesterRequest) (*teov20220901.DescribeSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewDescribeSecurityClientAttesterResponse()
		resp.Response = &teov20220901.DescribeSecurityClientAttesterResponseParams{
			TotalCount: ptrInt64SecurityClientAttester(1),
			ClientAttesters: []*teov20220901.ClientAttester{
				{
					Id:             ptrStrSecurityClientAttester("attest-ijklmnop"),
					Name:           ptrStrSecurityClientAttester("tf-example"),
					Type:           ptrStrSecurityClientAttester("CUSTOM"),
					AttesterSource: ptrStrSecurityClientAttester("TC-EO-CAPTCHA"),
					TCEOCaptchaOption: &teov20220901.TCEOCaptchaOption{
						CaptchaMode:       ptrStrSecurityClientAttester("Adaptive"),
						MobileCaptchaMode: ptrStrSecurityClientAttester("Visible"),
					},
				},
			},
			RequestId: ptrStrSecurityClientAttester("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityClientAttester()
	res := teo.ResourceTencentCloudTeoSecurityClientAttester()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":            "tf-example",
				"attester_source": "TC-EO-CAPTCHA",
				"tc_eo_captcha_option": []interface{}{
					map[string]interface{}{
						"captcha_mode":        "Adaptive",
						"mobile_captcha_mode": "Visible",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678"+tccommon.FILED_SP+"attest-ijklmnop", d.Id())
	attesters := d.Get("client_attesters").([]interface{})
	assert.Len(t, attesters, 1)
	attesterMap := attesters[0].(map[string]interface{})
	eoList := attesterMap["tc_eo_captcha_option"].([]interface{})
	assert.Len(t, eoList, 1)
	eoMap := eoList[0].(map[string]interface{})
	assert.Equal(t, "Visible", eoMap["mobile_captcha_mode"])
}

// TestTeoSecurityClientAttester_Create_WithoutNewFields tests Create without new fields (backward compatibility)
func TestTeoSecurityClientAttester_Create_WithoutNewFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityClientAttesterWithContext", func(_ context.Context, request *teov20220901.CreateSecurityClientAttesterRequest) (*teov20220901.CreateSecurityClientAttesterResponse, error) {
		assert.Len(t, request.ClientAttesters, 1)
		attester := request.ClientAttesters[0]
		assert.Nil(t, attester.UsageLimit)
		assert.Nil(t, attester.MaxUsageCount)
		resp := teov20220901.NewCreateSecurityClientAttesterResponse()
		resp.Response = &teov20220901.CreateSecurityClientAttesterResponseParams{
			ClientAttesterIds: []*string{ptrStrSecurityClientAttester("attest-qrstuvwx")},
			RequestId:         ptrStrSecurityClientAttester("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityClientAttester", func(request *teov20220901.DescribeSecurityClientAttesterRequest) (*teov20220901.DescribeSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewDescribeSecurityClientAttesterResponse()
		resp.Response = &teov20220901.DescribeSecurityClientAttesterResponseParams{
			TotalCount: ptrInt64SecurityClientAttester(1),
			ClientAttesters: []*teov20220901.ClientAttester{
				{
					Id:               ptrStrSecurityClientAttester("attest-qrstuvwx"),
					Name:             ptrStrSecurityClientAttester("tf-example"),
					Type:             ptrStrSecurityClientAttester("CUSTOM"),
					AttesterSource:   ptrStrSecurityClientAttester("TC-RCE"),
					AttesterDuration: ptrStrSecurityClientAttester("300s"),
				},
			},
			RequestId: ptrStrSecurityClientAttester("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityClientAttester()
	res := teo.ResourceTencentCloudTeoSecurityClientAttester()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":              "tf-example",
				"attester_source":   "TC-RCE",
				"attester_duration": "300s",
				"tc_rce_option": []interface{}{
					map[string]interface{}{
						"channel": "12399223",
						"region":  "ap-beijing",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678"+tccommon.FILED_SP+"attest-qrstuvwx", d.Id())
	attesters := d.Get("client_attesters").([]interface{})
	assert.Len(t, attesters, 1)
	attesterMap := attesters[0].(map[string]interface{})
	assert.Equal(t, "", attesterMap["usage_limit"])
	assert.Equal(t, 0, attesterMap["max_usage_count"])
}

// TestTeoSecurityClientAttester_Schema validates new fields schema definition
func TestTeoSecurityClientAttester_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoSecurityClientAttester()
	assert.NotNil(t, res)

	clientAttesters := res.Schema["client_attesters"]
	assert.NotNil(t, clientAttesters.Elem)
	attesterSchema := clientAttesters.Elem.(*schema.Resource).Schema

	usageLimit, ok := attesterSchema["usage_limit"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, usageLimit.Type)
	assert.True(t, usageLimit.Optional)
	assert.False(t, usageLimit.Required)

	maxUsageCount, ok := attesterSchema["max_usage_count"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeInt, maxUsageCount.Type)
	assert.True(t, maxUsageCount.Optional)
	assert.False(t, maxUsageCount.Required)

	eoOpt, ok := attesterSchema["tc_eo_captcha_option"]
	assert.True(t, ok)
	eoSchema := eoOpt.Elem.(*schema.Resource).Schema
	mobileCaptchaMode, ok := eoSchema["mobile_captcha_mode"]
	assert.True(t, ok)
	assert.Equal(t, schema.TypeString, mobileCaptchaMode.Type)
	assert.True(t, mobileCaptchaMode.Optional)
	assert.False(t, mobileCaptchaMode.Required)
}
