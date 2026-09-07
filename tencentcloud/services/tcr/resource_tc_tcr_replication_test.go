package tcr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tcrv20190924 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcr"
)

type mockMetaTcrReplication struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTcrReplication) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTcrReplication{}

func newMockMetaTcrReplication() *mockMetaTcrReplication {
	return &mockMetaTcrReplication{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTcrReplication(s string) *string {
	return &s
}

func ptrBoolTcrReplication(b bool) *bool {
	return &b
}

// TestAccTencentCloudTcrReplicationResource_basic is the existing acceptance test
func TestAccTencentCloudTcrReplicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTcrReplication,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_tcr_replication.example", "id"),
			),
		}},
	})
}

const testAccTcrReplication = `
resource "tencentcloud_tcr_replication" "example" {
  source_registry_id      = "tcr-9q9h1nof"
  destination_registry_id = "tcr-jtih9ngc"
  rule {
    name           = "tf-example"
    dest_namespace = ""
    override       = true
    deletion       = true
    filters {
      type  = "name"
      value = "tf-example/**"
    }
  }

  destination_region_id = 1
  description           = "remark."
}
`

// TestTcrReplicationUpdate_Description tests in-place update of description
func TestTcrReplicationUpdate_Description(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tcrClient := &tcrv20190924.Client{}
	patches.ApplyMethodReturn(newMockMetaTcrReplication().client, "UseTCRClient", tcrClient)

	patches.ApplyMethodFunc(tcrClient, "ModifyReplicationWithContext", func(ctx context.Context, request *tcrv20190924.ModifyReplicationRequest) (*tcrv20190924.ModifyReplicationResponse, error) {
		assert.Equal(t, "tcr-9q9h1nof", *request.SourceRegistryId)
		assert.Equal(t, "tf-example", *request.RuleName)
		assert.Equal(t, "updated-remark.", *request.Description)
		resp := tcrv20190924.NewModifyReplicationResponse()
		resp.Response = &tcrv20190924.ModifyReplicationResponseParams{
			RequestId: ptrStringTcrReplication("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tcrClient, "DescribeReplicationPolicies", func(request *tcrv20190924.DescribeReplicationPoliciesRequest) (*tcrv20190924.DescribeReplicationPoliciesResponse, error) {
		resp := tcrv20190924.NewDescribeReplicationPoliciesResponse()
		resp.Response = &tcrv20190924.DescribeReplicationPoliciesResponseParams{
			ReplicationPolicyInfoList: []*tcrv20190924.ReplicationPolicyInfo{
				{
					Name:         ptrStringTcrReplication("tf-example"),
					DestResource: ptrStringTcrReplication(""),
					Override:     ptrBoolTcrReplication(true),
					Description:  ptrStringTcrReplication("updated-remark."),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaTcrReplication()
	res := tcr.ResourceTencentCloudTcrReplication()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"source_registry_id":      "tcr-9q9h1nof",
		"destination_registry_id": "tcr-jtih9ngc",
		"rule": []interface{}{
			map[string]interface{}{
				"name":           "tf-example",
				"dest_namespace": "",
				"override":       true,
				"deletion":       true,
				"filters": []interface{}{
					map[string]interface{}{
						"type":  "name",
						"value": "tf-example/**",
					},
				},
			},
		},
		"destination_region_id": 1,
		"description":           "updated-remark.",
	})
	d.SetId("tcr-9q9h1nof#tf-example")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTcrReplicationUpdate_Filters tests in-place update of rule filters
func TestTcrReplicationUpdate_Filters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tcrClient := &tcrv20190924.Client{}
	patches.ApplyMethodReturn(newMockMetaTcrReplication().client, "UseTCRClient", tcrClient)

	patches.ApplyMethodFunc(tcrClient, "ModifyReplicationWithContext", func(ctx context.Context, request *tcrv20190924.ModifyReplicationRequest) (*tcrv20190924.ModifyReplicationResponse, error) {
		assert.Equal(t, "tcr-9q9h1nof", *request.SourceRegistryId)
		assert.Equal(t, "tf-example", *request.RuleName)
		assert.NotNil(t, request.Rule)
		assert.Len(t, request.Rule.Filters, 2)
		assert.Equal(t, "name", *request.Rule.Filters[0].Type)
		assert.Equal(t, "tf-example/**", *request.Rule.Filters[0].Value)
		assert.Equal(t, "tag", *request.Rule.Filters[1].Type)
		assert.Equal(t, "latest", *request.Rule.Filters[1].Value)
		resp := tcrv20190924.NewModifyReplicationResponse()
		resp.Response = &tcrv20190924.ModifyReplicationResponseParams{
			RequestId: ptrStringTcrReplication("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tcrClient, "DescribeReplicationPolicies", func(request *tcrv20190924.DescribeReplicationPoliciesRequest) (*tcrv20190924.DescribeReplicationPoliciesResponse, error) {
		resp := tcrv20190924.NewDescribeReplicationPoliciesResponse()
		resp.Response = &tcrv20190924.DescribeReplicationPoliciesResponseParams{
			ReplicationPolicyInfoList: []*tcrv20190924.ReplicationPolicyInfo{
				{
					Name:         ptrStringTcrReplication("tf-example"),
					DestResource: ptrStringTcrReplication(""),
					Override:     ptrBoolTcrReplication(true),
					Filters: []*tcrv20190924.PolicyFilter{
						{
							Type:  ptrStringTcrReplication("name"),
							Value: ptrStringTcrReplication("tf-example/**"),
						},
						{
							Type:  ptrStringTcrReplication("tag"),
							Value: ptrStringTcrReplication("latest"),
						},
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaTcrReplication()
	res := tcr.ResourceTencentCloudTcrReplication()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"source_registry_id":      "tcr-9q9h1nof",
		"destination_registry_id": "tcr-jtih9ngc",
		"rule": []interface{}{
			map[string]interface{}{
				"name":           "tf-example",
				"dest_namespace": "",
				"override":       true,
				"deletion":       true,
				"filters": []interface{}{
					map[string]interface{}{
						"type":  "name",
						"value": "tf-example/**",
					},
					map[string]interface{}{
						"type":  "tag",
						"value": "latest",
					},
				},
			},
		},
		"destination_region_id": 1,
		"description":           "remark.",
	})
	d.SetId("tcr-9q9h1nof#tf-example")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTcrReplicationUpdate_APIError tests update with API error (retry exhausted)
func TestTcrReplicationUpdate_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tcrClient := &tcrv20190924.Client{}
	patches.ApplyMethodReturn(newMockMetaTcrReplication().client, "UseTCRClient", tcrClient)

	patches.ApplyMethodFunc(tcrClient, "ModifyReplicationWithContext", func(ctx context.Context, request *tcrv20190924.ModifyReplicationRequest) (*tcrv20190924.ModifyReplicationResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=Internal error")
	})

	meta := newMockMetaTcrReplication()
	res := tcr.ResourceTencentCloudTcrReplication()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"source_registry_id":      "tcr-9q9h1nof",
		"destination_registry_id": "tcr-jtih9ngc",
		"rule": []interface{}{
			map[string]interface{}{
				"name":           "tf-example",
				"dest_namespace": "",
				"override":       true,
				"deletion":       true,
				"filters": []interface{}{
					map[string]interface{}{
						"type":  "name",
						"value": "tf-example/**",
					},
				},
			},
		},
		"destination_region_id": 1,
		"description":           "updated-remark.",
	})
	d.SetId("tcr-9q9h1nof#tf-example")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
}
