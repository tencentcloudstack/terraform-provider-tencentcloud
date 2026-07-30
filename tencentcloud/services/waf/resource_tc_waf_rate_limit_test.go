package waf_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/waf"
)

type mockMetaWafRateLimit struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaWafRateLimit) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaWafRateLimit{}

func newMockMetaWafRateLimit() *mockMetaWafRateLimit {
	return &mockMetaWafRateLimit{client: &connectivity.TencentCloudClient{}}
}

func ptrStringWRL(s string) *string {
	return &s
}

func ptrInt64WRL(i int64) *int64 {
	return &i
}

func ptrBoolWRL(b bool) *bool {
	return &b
}

func TestWafRateLimit_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	wafClient := &wafv20180125.Client{}
	patches.ApplyMethodReturn(newMockMetaWafRateLimit().client, "UseWafV20180125Client", wafClient)

	patches.ApplyMethodFunc(wafClient, "CreateRateLimitV2WithContext", func(_ context.Context, request *wafv20180125.CreateRateLimitV2Request) (*wafv20180125.CreateRateLimitV2Response, error) {
		resp := wafv20180125.NewCreateRateLimitV2Response()
		resp.Response = &wafv20180125.CreateRateLimitV2ResponseParams{
			LimitRuleID: ptrInt64WRL(12345),
			Domain:      ptrStringWRL("example.com"),
			RequestId:   ptrStringWRL("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(wafClient, "DescribeRateLimitsV2WithContext", func(_ context.Context, request *wafv20180125.DescribeRateLimitsV2Request) (*wafv20180125.DescribeRateLimitsV2Response, error) {
		resp := wafv20180125.NewDescribeRateLimitsV2Response()
		resp.Response = &wafv20180125.DescribeRateLimitsV2ResponseParams{
			Total:     ptrInt64WRL(1),
			RequestId: ptrStringWRL("fake-request-id"),
			RateLimits: []*wafv20180125.LimitRuleV2{
				{
					LimitRuleID:   ptrInt64WRL(12345),
					Name:          ptrStringWRL("test-rule"),
					Priority:      ptrInt64WRL(100),
					Status:        ptrInt64WRL(1),
					LimitObject:   ptrStringWRL("Domain"),
					LimitStrategy: ptrInt64WRL(1),
					LimitWindow: &wafv20180125.LimitWindow{
						Second: ptrInt64WRL(10),
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaWafRateLimit()
	res := waf.ResourceTencentCloudWafRateLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"domain":         "example.com",
		"name":           "test-rule",
		"priority":       100,
		"status":         1,
		"limit_object":   "Domain",
		"limit_strategy": 1,
		"limit_window": []interface{}{
			map[string]interface{}{
				"second":      10,
				"minute":      0,
				"hour":        0,
				"quota_share": false,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "example.com#12345", d.Id())
}

func TestWafRateLimit_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	wafClient := &wafv20180125.Client{}
	patches.ApplyMethodReturn(newMockMetaWafRateLimit().client, "UseWafV20180125Client", wafClient)

	patches.ApplyMethodFunc(wafClient, "DescribeRateLimitsV2WithContext", func(_ context.Context, request *wafv20180125.DescribeRateLimitsV2Request) (*wafv20180125.DescribeRateLimitsV2Response, error) {
		resp := wafv20180125.NewDescribeRateLimitsV2Response()
		resp.Response = &wafv20180125.DescribeRateLimitsV2ResponseParams{
			Total:     ptrInt64WRL(1),
			RequestId: ptrStringWRL("fake-request-id"),
			RateLimits: []*wafv20180125.LimitRuleV2{
				{
					LimitRuleID:   ptrInt64WRL(12345),
					Name:          ptrStringWRL("test-rule"),
					Priority:      ptrInt64WRL(100),
					Status:        ptrInt64WRL(1),
					LimitObject:   ptrStringWRL("Domain"),
					LimitStrategy: ptrInt64WRL(1),
					LimitWindow: &wafv20180125.LimitWindow{
						Second: ptrInt64WRL(10),
						Minute: ptrInt64WRL(100),
					},
					BlockPage: ptrInt64WRL(0),
					ObjectSrc: ptrInt64WRL(0),
					Order:     ptrInt64WRL(0),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaWafRateLimit()
	res := waf.ResourceTencentCloudWafRateLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"domain":         "example.com",
		"name":           "test-rule",
		"priority":       100,
		"status":         1,
		"limit_object":   "Domain",
		"limit_strategy": 1,
		"limit_window": []interface{}{
			map[string]interface{}{
				"second":      10,
				"minute":      0,
				"hour":        0,
				"quota_share": false,
			},
		},
	})
	d.SetId("example.com#12345")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "example.com#12345", d.Id())
	assert.Equal(t, "test-rule", d.Get("name"))
	assert.Equal(t, 100, d.Get("priority"))
	assert.Equal(t, 1, d.Get("status"))
	assert.Equal(t, "Domain", d.Get("limit_object"))
	assert.Equal(t, 1, d.Get("limit_strategy"))
}

func TestWafRateLimit_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	wafClient := &wafv20180125.Client{}
	patches.ApplyMethodReturn(newMockMetaWafRateLimit().client, "UseWafV20180125Client", wafClient)

	patches.ApplyMethodFunc(wafClient, "DescribeRateLimitsV2WithContext", func(_ context.Context, request *wafv20180125.DescribeRateLimitsV2Request) (*wafv20180125.DescribeRateLimitsV2Response, error) {
		resp := wafv20180125.NewDescribeRateLimitsV2Response()
		resp.Response = &wafv20180125.DescribeRateLimitsV2ResponseParams{
			Total:      ptrInt64WRL(0),
			RequestId:  ptrStringWRL("fake-request-id"),
			RateLimits: []*wafv20180125.LimitRuleV2{},
		}
		return resp, nil
	})

	meta := newMockMetaWafRateLimit()
	res := waf.ResourceTencentCloudWafRateLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"domain":         "example.com",
		"name":           "test-rule",
		"priority":       100,
		"status":         1,
		"limit_object":   "Domain",
		"limit_strategy": 1,
		"limit_window": []interface{}{
			map[string]interface{}{
				"second":      10,
				"minute":      0,
				"hour":        0,
				"quota_share": false,
			},
		},
	})
	d.SetId("example.com#12345")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestWafRateLimit_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	wafClient := &wafv20180125.Client{}
	patches.ApplyMethodReturn(newMockMetaWafRateLimit().client, "UseWafV20180125Client", wafClient)

	patches.ApplyMethodFunc(wafClient, "UpdateRateLimitV2WithContext", func(_ context.Context, request *wafv20180125.UpdateRateLimitV2Request) (*wafv20180125.UpdateRateLimitV2Response, error) {
		resp := wafv20180125.NewUpdateRateLimitV2Response()
		resp.Response = &wafv20180125.UpdateRateLimitV2ResponseParams{
			LimitRuleID: ptrInt64WRL(12345),
			Domain:      ptrStringWRL("example.com"),
			RequestId:   ptrStringWRL("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(wafClient, "DescribeRateLimitsV2WithContext", func(_ context.Context, request *wafv20180125.DescribeRateLimitsV2Request) (*wafv20180125.DescribeRateLimitsV2Response, error) {
		resp := wafv20180125.NewDescribeRateLimitsV2Response()
		resp.Response = &wafv20180125.DescribeRateLimitsV2ResponseParams{
			Total:     ptrInt64WRL(1),
			RequestId: ptrStringWRL("fake-request-id"),
			RateLimits: []*wafv20180125.LimitRuleV2{
				{
					LimitRuleID:   ptrInt64WRL(12345),
					Name:          ptrStringWRL("updated-rule"),
					Priority:      ptrInt64WRL(200),
					Status:        ptrInt64WRL(1),
					LimitObject:   ptrStringWRL("Domain"),
					LimitStrategy: ptrInt64WRL(0),
					LimitWindow: &wafv20180125.LimitWindow{
						Minute: ptrInt64WRL(200),
					},
					QuotaShare: ptrBoolWRL(true),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaWafRateLimit()
	res := waf.ResourceTencentCloudWafRateLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"domain":         "example.com",
		"name":           "updated-rule",
		"priority":       200,
		"status":         1,
		"limit_object":   "Domain",
		"limit_strategy": 0,
		"quota_share":    true,
		"limit_window": []interface{}{
			map[string]interface{}{
				"second":      0,
				"minute":      200,
				"hour":        0,
				"quota_share": false,
			},
		},
	})
	d.SetId("example.com#12345")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "example.com#12345", d.Id())
}

func TestWafRateLimit_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	wafClient := &wafv20180125.Client{}
	patches.ApplyMethodReturn(newMockMetaWafRateLimit().client, "UseWafV20180125Client", wafClient)

	patches.ApplyMethodFunc(wafClient, "DeleteRateLimitsV2WithContext", func(_ context.Context, request *wafv20180125.DeleteRateLimitsV2Request) (*wafv20180125.DeleteRateLimitsV2Response, error) {
		resp := wafv20180125.NewDeleteRateLimitsV2Response()
		resp.Response = &wafv20180125.DeleteRateLimitsV2ResponseParams{
			RequestId: ptrStringWRL("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaWafRateLimit()
	res := waf.ResourceTencentCloudWafRateLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"domain":         "example.com",
		"name":           "test-rule",
		"priority":       100,
		"status":         1,
		"limit_object":   "Domain",
		"limit_strategy": 1,
		"limit_window": []interface{}{
			map[string]interface{}{
				"second":      10,
				"minute":      0,
				"hour":        0,
				"quota_share": false,
			},
		},
	})
	d.SetId("example.com#12345")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
