package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

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

// go test ./tencentcloud/services/teo/ -run "TestSecurityClientAttester" -v -count=1 -gcflags="all=-l"

func TestSecurityClientAttester_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityClientAttester", func(request *teov20220901.DescribeSecurityClientAttesterRequest) (*teov20220901.DescribeSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewDescribeSecurityClientAttesterResponse()
		resp.Response = &teov20220901.DescribeSecurityClientAttesterResponseParams{
			TotalCount: ptrInt64SecurityClientAttester(2),
			ClientAttesters: []*teov20220901.ClientAttester{
				{
					Id:               ptrStrSecurityClientAttester("att-001"),
					Name:             ptrStrSecurityClientAttester("test-rce"),
					Type:             ptrStrSecurityClientAttester("CUSTOM"),
					AttesterSource:   ptrStrSecurityClientAttester("TC-RCE"),
					AttesterDuration: ptrStrSecurityClientAttester("60s"),
					TCRCEOption: &teov20220901.TCRCEOption{
						Channel: ptrStrSecurityClientAttester("channel-1"),
						Region:  ptrStrSecurityClientAttester("ap-beijing"),
					},
				},
				{
					Id:               ptrStrSecurityClientAttester("att-002"),
					Name:             ptrStrSecurityClientAttester("test-captcha"),
					Type:             ptrStrSecurityClientAttester("CUSTOM"),
					AttesterSource:   ptrStrSecurityClientAttester("TC-CAPTCHA"),
					AttesterDuration: ptrStrSecurityClientAttester("120s"),
					TCCaptchaOption: &teov20220901.TCCaptchaOption{
						CaptchaAppId: ptrStrSecurityClientAttester("captcha-app-123"),
						AppSecretKey: ptrStrSecurityClientAttester("secret-key-456"),
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
		"zone_id": "zone-test123",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":            "test-rce",
				"attester_source": "TC-RCE",
			},
		},
	})
	d.SetId("zone-test123#att-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clientAttesters := d.Get("client_attesters").([]interface{})
	assert.Equal(t, 1, len(clientAttesters))

	attesterMap := clientAttesters[0].(map[string]interface{})
	assert.Equal(t, "test-rce", attesterMap["name"])
	assert.Equal(t, "TC-RCE", attesterMap["attester_source"])
	assert.Equal(t, "60s", attesterMap["attester_duration"])
	assert.Equal(t, "att-001", attesterMap["id"])
	assert.Equal(t, "CUSTOM", attesterMap["type"])

	tcRceOptions := attesterMap["tc_rce_option"].([]interface{})
	assert.Equal(t, 1, len(tcRceOptions))
	tcRceOptionMap := tcRceOptions[0].(map[string]interface{})
	assert.Equal(t, "channel-1", tcRceOptionMap["channel"])
	assert.Equal(t, "ap-beijing", tcRceOptionMap["region"])

	clientAttesterIds := d.Get("client_attester_ids").([]interface{})
	assert.Equal(t, 1, len(clientAttesterIds))
	assert.Equal(t, "att-001", clientAttesterIds[0])
}

func TestSecurityClientAttester_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityClientAttester", func(request *teov20220901.DescribeSecurityClientAttesterRequest) (*teov20220901.DescribeSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewDescribeSecurityClientAttesterResponse()
		resp.Response = &teov20220901.DescribeSecurityClientAttesterResponseParams{
			TotalCount:      ptrInt64SecurityClientAttester(0),
			ClientAttesters: []*teov20220901.ClientAttester{},
			RequestId:       ptrStrSecurityClientAttester("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityClientAttester()
	res := teo.ResourceTencentCloudTeoSecurityClientAttester()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":            "test-rce",
				"attester_source": "TC-RCE",
			},
		},
	})
	d.SetId("zone-test123#att-999")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestSecurityClientAttester_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityClientAttester", func(request *teov20220901.CreateSecurityClientAttesterRequest) (*teov20220901.CreateSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewCreateSecurityClientAttesterResponse()
		resp.Response = &teov20220901.CreateSecurityClientAttesterResponseParams{
			ClientAttesterIds: []*string{
				ptrStrSecurityClientAttester("att-001"),
				ptrStrSecurityClientAttester("att-002"),
			},
			RequestId: ptrStrSecurityClientAttester("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityClientAttester", func(request *teov20220901.DescribeSecurityClientAttesterRequest) (*teov20220901.DescribeSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewDescribeSecurityClientAttesterResponse()
		resp.Response = &teov20220901.DescribeSecurityClientAttesterResponseParams{
			TotalCount: ptrInt64SecurityClientAttester(2),
			ClientAttesters: []*teov20220901.ClientAttester{
				{
					Id:               ptrStrSecurityClientAttester("att-001"),
					Name:             ptrStrSecurityClientAttester("test-rce"),
					Type:             ptrStrSecurityClientAttester("CUSTOM"),
					AttesterSource:   ptrStrSecurityClientAttester("TC-RCE"),
					AttesterDuration: ptrStrSecurityClientAttester("60s"),
					TCRCEOption: &teov20220901.TCRCEOption{
						Channel: ptrStrSecurityClientAttester("channel-1"),
						Region:  ptrStrSecurityClientAttester("ap-beijing"),
					},
				},
				{
					Id:               ptrStrSecurityClientAttester("att-002"),
					Name:             ptrStrSecurityClientAttester("test-eo-captcha"),
					Type:             ptrStrSecurityClientAttester("CUSTOM"),
					AttesterSource:   ptrStrSecurityClientAttester("TC-EO-CAPTCHA"),
					AttesterDuration: ptrStrSecurityClientAttester("120s"),
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
		"zone_id": "zone-test123",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":              "test-rce",
				"attester_source":   "TC-RCE",
				"attester_duration": "60s",
				"tc_rce_option": []interface{}{
					map[string]interface{}{
						"channel": "channel-1",
						"region":  "ap-beijing",
					},
				},
			},
			map[string]interface{}{
				"name":              "test-eo-captcha",
				"attester_source":   "TC-EO-CAPTCHA",
				"attester_duration": "120s",
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
	assert.NotEmpty(t, d.Id())

	idParts := d.Id()
	assert.Contains(t, idParts, "zone-test123")
	assert.Contains(t, idParts, "att-001")
	assert.Contains(t, idParts, "att-002")

	clientAttesterIds := d.Get("client_attester_ids").([]interface{})
	assert.Equal(t, 2, len(clientAttesterIds))
}

func TestSecurityClientAttester_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteSecurityClientAttester", func(request *teov20220901.DeleteSecurityClientAttesterRequest) (*teov20220901.DeleteSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewDeleteSecurityClientAttesterResponse()
		resp.Response = &teov20220901.DeleteSecurityClientAttesterResponseParams{
			RequestId: ptrStrSecurityClientAttester("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityClientAttester()
	res := teo.ResourceTencentCloudTeoSecurityClientAttester()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":            "test-rce",
				"attester_source": "TC-RCE",
			},
		},
		"client_attester_ids": []interface{}{
			"att-001",
			"att-002",
		},
	})
	d.SetId("zone-test123#att-001,att-002")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestSecurityClientAttester_Update_ModifySuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifySecurityClientAttester", func(request *teov20220901.ModifySecurityClientAttesterRequest) (*teov20220901.ModifySecurityClientAttesterResponse, error) {
		resp := teov20220901.NewModifySecurityClientAttesterResponse()
		resp.Response = &teov20220901.ModifySecurityClientAttesterResponseParams{
			RequestId: ptrStrSecurityClientAttester("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityClientAttester", func(request *teov20220901.DescribeSecurityClientAttesterRequest) (*teov20220901.DescribeSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewDescribeSecurityClientAttesterResponse()
		resp.Response = &teov20220901.DescribeSecurityClientAttesterResponseParams{
			TotalCount: ptrInt64SecurityClientAttester(1),
			ClientAttesters: []*teov20220901.ClientAttester{
				{
					Id:               ptrStrSecurityClientAttester("att-001"),
					Name:             ptrStrSecurityClientAttester("test-rce-updated"),
					Type:             ptrStrSecurityClientAttester("CUSTOM"),
					AttesterSource:   ptrStrSecurityClientAttester("TC-RCE"),
					AttesterDuration: ptrStrSecurityClientAttester("120s"),
					TCRCEOption: &teov20220901.TCRCEOption{
						Channel: ptrStrSecurityClientAttester("channel-2"),
						Region:  ptrStrSecurityClientAttester("ap-singapore"),
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
		"zone_id": "zone-test123",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":              "test-rce-updated",
				"attester_source":   "TC-RCE",
				"attester_duration": "120s",
				"tc_rce_option": []interface{}{
					map[string]interface{}{
						"channel": "channel-2",
						"region":  "ap-singapore",
					},
				},
			},
		},
		"client_attester_ids": []interface{}{
			"att-001",
		},
	})
	d.SetId("zone-test123#att-001")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	clientAttesters := d.Get("client_attesters").([]interface{})
	assert.Equal(t, 1, len(clientAttesters))
	attesterMap := clientAttesters[0].(map[string]interface{})
	assert.Equal(t, "test-rce-updated", attesterMap["name"])
}

func TestSecurityClientAttester_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoSecurityClientAttester()

	assert.NotNil(t, res)

	// Check zone_id field
	assert.Contains(t, res.Schema, "zone_id")
	zoneIdSchema := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneIdSchema.Type)
	assert.True(t, zoneIdSchema.Required)
	assert.True(t, zoneIdSchema.ForceNew)

	// Check client_attesters field
	assert.Contains(t, res.Schema, "client_attesters")
	clientAttestersSchema := res.Schema["client_attesters"]
	assert.Equal(t, schema.TypeList, clientAttestersSchema.Type)
	assert.True(t, clientAttestersSchema.Required)

	// Check client_attester_ids field
	assert.Contains(t, res.Schema, "client_attester_ids")
	clientAttesterIdsSchema := res.Schema["client_attester_ids"]
	assert.Equal(t, schema.TypeList, clientAttesterIdsSchema.Type)
	assert.True(t, clientAttesterIdsSchema.Computed)
}

func TestSecurityClientAttester_Read_EOCaptchaOption(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityClientAttester", func(request *teov20220901.DescribeSecurityClientAttesterRequest) (*teov20220901.DescribeSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewDescribeSecurityClientAttesterResponse()
		resp.Response = &teov20220901.DescribeSecurityClientAttesterResponseParams{
			TotalCount: ptrInt64SecurityClientAttester(1),
			ClientAttesters: []*teov20220901.ClientAttester{
				{
					Id:               ptrStrSecurityClientAttester("att-003"),
					Name:             ptrStrSecurityClientAttester("test-eo-captcha"),
					Type:             ptrStrSecurityClientAttester("CUSTOM"),
					AttesterSource:   ptrStrSecurityClientAttester("TC-EO-CAPTCHA"),
					AttesterDuration: ptrStrSecurityClientAttester("300s"),
					TCEOCaptchaOption: &teov20220901.TCEOCaptchaOption{
						CaptchaMode: ptrStrSecurityClientAttester("Adaptive"),
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
		"zone_id": "zone-test123",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":            "test-eo-captcha",
				"attester_source": "TC-EO-CAPTCHA",
			},
		},
	})
	d.SetId("zone-test123#att-003")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clientAttesters := d.Get("client_attesters").([]interface{})
	assert.Equal(t, 1, len(clientAttesters))
	attesterMap := clientAttesters[0].(map[string]interface{})
	assert.Equal(t, "TC-EO-CAPTCHA", attesterMap["attester_source"])

	tcEoCaptchaOptions := attesterMap["tc_eo_captcha_option"].([]interface{})
	assert.Equal(t, 1, len(tcEoCaptchaOptions))
	tcEoCaptchaOptionMap := tcEoCaptchaOptions[0].(map[string]interface{})
	assert.Equal(t, "Adaptive", tcEoCaptchaOptionMap["captcha_mode"])
}

func TestSecurityClientAttester_Read_CaptchaOption(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityClientAttester", func(request *teov20220901.DescribeSecurityClientAttesterRequest) (*teov20220901.DescribeSecurityClientAttesterResponse, error) {
		resp := teov20220901.NewDescribeSecurityClientAttesterResponse()
		resp.Response = &teov20220901.DescribeSecurityClientAttesterResponseParams{
			TotalCount: ptrInt64SecurityClientAttester(1),
			ClientAttesters: []*teov20220901.ClientAttester{
				{
					Id:               ptrStrSecurityClientAttester("att-004"),
					Name:             ptrStrSecurityClientAttester("test-captcha"),
					Type:             ptrStrSecurityClientAttester("CUSTOM"),
					AttesterSource:   ptrStrSecurityClientAttester("TC-CAPTCHA"),
					AttesterDuration: ptrStrSecurityClientAttester("180s"),
					TCCaptchaOption: &teov20220901.TCCaptchaOption{
						CaptchaAppId: ptrStrSecurityClientAttester("app-id-789"),
						AppSecretKey: ptrStrSecurityClientAttester("secret-abc"),
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
		"zone_id": "zone-test123",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":            "test-captcha",
				"attester_source": "TC-CAPTCHA",
			},
		},
	})
	d.SetId("zone-test123#att-004")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clientAttesters := d.Get("client_attesters").([]interface{})
	assert.Equal(t, 1, len(clientAttesters))
	attesterMap := clientAttesters[0].(map[string]interface{})
	assert.Equal(t, "TC-CAPTCHA", attesterMap["attester_source"])

	tcCaptchaOptions := attesterMap["tc_captcha_option"].([]interface{})
	assert.Equal(t, 1, len(tcCaptchaOptions))
	tcCaptchaOptionMap := tcCaptchaOptions[0].(map[string]interface{})
	assert.Equal(t, "app-id-789", tcCaptchaOptionMap["captcha_app_id"])
	assert.Equal(t, "secret-abc", tcCaptchaOptionMap["app_secret_key"])
}

func TestSecurityClientAttester_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityClientAttester().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityClientAttester", func(request *teov20220901.CreateSecurityClientAttesterRequest) (*teov20220901.CreateSecurityClientAttesterResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=internal error")
	})

	meta := newMockMetaForSecurityClientAttester()
	res := teo.ResourceTencentCloudTeoSecurityClientAttester()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"client_attesters": []interface{}{
			map[string]interface{}{
				"name":            "test-rce",
				"attester_source": "TC-RCE",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
}
