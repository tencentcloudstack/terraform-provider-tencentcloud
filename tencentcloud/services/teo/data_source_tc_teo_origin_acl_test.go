package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// mockMetaOriginAclDS implements tccommon.ProviderMeta
type mockMetaOriginAclDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaOriginAclDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaOriginAclDS{}

func newMockMetaOriginAclDS() *mockMetaOriginAclDS {
	return &mockMetaOriginAclDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStringOriginAclDS(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestTeoOriginAclDS" -v -count=1 -gcflags="all=-l"

// TestTeoOriginAclDS_ReadWithOriginACLFamily tests data source Read with OriginACLFamily in response
func TestTeoOriginAclDS_ReadWithOriginACLFamily(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginAclDS().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginACLWithContext", func(_ context.Context, request *teov20220901.DescribeOriginACLRequest) (*teov20220901.DescribeOriginACLResponse, error) {
		resp := teov20220901.NewDescribeOriginACLResponse()
		resp.Response = &teov20220901.DescribeOriginACLResponseParams{
			OriginACLInfo: &teov20220901.OriginACLInfo{
				Status:          ptrStringOriginAclDS("online"),
				L7Hosts:         []*string{ptrStringOriginAclDS("example1.com")},
				L4ProxyIds:      []*string{ptrStringOriginAclDS("sid-abc123")},
				OriginACLFamily: ptrStringOriginAclDS("mlc"),
			},
		}
		return resp, nil
	})

	meta := newMockMetaOriginAclDS()
	res := teo.DataSourceTencentCloudTeoOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())

	originAclInfo := d.Get("origin_acl_info").([]interface{})
	assert.Len(t, originAclInfo, 1)
	infoMap := originAclInfo[0].(map[string]interface{})
	assert.Equal(t, "mlc", infoMap["origin_acl_family"].(string))
}

// TestTeoOriginAclDS_ReadWithNilOriginACLFamily tests data source Read when OriginACLFamily is nil
func TestTeoOriginAclDS_ReadWithNilOriginACLFamily(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginAclDS().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginACLWithContext", func(_ context.Context, request *teov20220901.DescribeOriginACLRequest) (*teov20220901.DescribeOriginACLResponse, error) {
		resp := teov20220901.NewDescribeOriginACLResponse()
		resp.Response = &teov20220901.DescribeOriginACLResponseParams{
			OriginACLInfo: &teov20220901.OriginACLInfo{
				Status:     ptrStringOriginAclDS("online"),
				L7Hosts:    []*string{ptrStringOriginAclDS("example1.com")},
				L4ProxyIds: []*string{ptrStringOriginAclDS("sid-abc123")},
			},
		}
		return resp, nil
	})

	meta := newMockMetaOriginAclDS()
	res := teo.DataSourceTencentCloudTeoOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())

	originAclInfo := d.Get("origin_acl_info").([]interface{})
	assert.Len(t, originAclInfo, 1)
	infoMap := originAclInfo[0].(map[string]interface{})
	assert.Equal(t, "", infoMap["origin_acl_family"].(string))
}

// TestTeoOriginAclDS_Schema tests the data source schema definition
func TestTeoOriginAclDS_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoOriginAcl()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "origin_acl_info")

	originAclInfoSchema := res.Schema["origin_acl_info"]
	assert.Equal(t, schema.TypeList, originAclInfoSchema.Type)

	elemRes := originAclInfoSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "origin_acl_family")

	originACLFamily := elemRes.Schema["origin_acl_family"]
	assert.Equal(t, schema.TypeString, originACLFamily.Type)
	assert.True(t, originACLFamily.Computed)
}

// TestAccTencentCloudTeoOriginAclDataSource_basic is the acceptance test
func TestAccTencentCloudTeoOriginAclDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoOriginAclDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_origin_acl.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_origin_acl.example", "zone_id"),
			),
		}},
	})
}

const testAccTeoOriginAclDataSource = `
data "tencentcloud_teo_origin_acl" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
`
