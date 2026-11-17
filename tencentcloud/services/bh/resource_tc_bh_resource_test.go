package bh_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bh"
)

func TestAccTencentCloudBhResourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_resource.example", "id"),
				),
			},
			{
				Config: testAccBhResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_resource.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_bh_resource.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBhResource = `
resource "tencentcloud_bh_resource" "example" {
  deploy_region    = "ap-guangzhou"
  vpc_id           = "vpc-q1of50wz"
  subnet_id        = "subnet-7uhvm46o"
  resource_edition = "standard"
  resource_node    = 20
  time_unit        = "m"
  time_span        = "1"
  pay_mode         = 1
  auto_renew_flag  = 1
  deploy_zone      = "ap-guangzhou-6"
  cidr_block       = "192.168.11.0/24"
  vpc_cidr_block   = "192.168.0.0/16"
}
`

const testAccBhResourceUpdate = `
resource "tencentcloud_bh_resource" "example" {
  deploy_region    = "ap-guangzhou"
  vpc_id           = "vpc-q1of50wz"
  subnet_id        = "subnet-7uhvm46o"
  resource_edition = "standard"
  resource_node    = 10
  time_unit        = "m"
  time_span        = "1"
  pay_mode         = 1
  auto_renew_flag  = 0
  deploy_zone      = "ap-guangzhou-6"
  cidr_block       = "192.168.11.0/24"
  vpc_cidr_block   = "192.168.0.0/16"
}
`

type mockMetaBhResource struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBhResource) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBhResource{}

func newMockMetaBhResource() *mockMetaBhResource {
	return &mockMetaBhResource{client: &connectivity.TencentCloudClient{}}
}

func ptrStringBhRes(s string) *string {
	return &s
}

func ptrUint64BhRes(i uint64) *uint64 {
	return &i
}

// buildBhResourceData builds a *schema.ResourceData with the given old state
// (key->string) and new raw config (key->interface{}) so that d.HasChange
// returns true for fields that differ between old and new.
func buildBhResourceData(t *testing.T, oldState map[string]string, newConfig map[string]interface{}) *schema.ResourceData {
	res := bh.ResourceTencentCloudBhResource()
	state := &terraform.InstanceState{
		ID:         "bh-saas-abc123",
		Attributes: oldState,
	}
	rawConfig := terraform.NewResourceConfigRaw(newConfig)
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)
	return d
}

// mockDescribeResourcesForAccess patches both DescribeResourcesWithContext
// (used by update wait loops) and DescribeResources (used by the final Read via
// the BhService) to return a Resource whose WebAccess/ClientAccess match the
// target values.
func mockDescribeResourcesForAccess(patches *gomonkey.Patches, bhClient *bhv20230418.Client, webAccess, clientAccess uint64) {
	resourceItem := &bhv20230418.Resource{
		ResourceId:   ptrStringBhRes("bh-saas-abc123"),
		WebAccess:    ptrUint64BhRes(webAccess),
		ClientAccess: ptrUint64BhRes(clientAccess),
	}
	patches.ApplyMethodFunc(bhClient, "DescribeResourcesWithContext", func(_ context.Context, request *bhv20230418.DescribeResourcesRequest) (*bhv20230418.DescribeResourcesResponse, error) {
		resp := bhv20230418.NewDescribeResourcesResponse()
		resp.Response = &bhv20230418.DescribeResourcesResponseParams{
			ResourceSet: []*bhv20230418.Resource{resourceItem},
			TotalCount:  ptrUint64BhRes(1),
			RequestId:   ptrStringBhRes("fake-request-id"),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(bhClient, "DescribeResources", func(request *bhv20230418.DescribeResourcesRequest) (*bhv20230418.DescribeResourcesResponse, error) {
		resp := bhv20230418.NewDescribeResourcesResponse()
		resp.Response = &bhv20230418.DescribeResourcesResponseParams{
			ResourceSet: []*bhv20230418.Resource{resourceItem},
			TotalCount:  ptrUint64BhRes(1),
			RequestId:   ptrStringBhRes("fake-request-id"),
		}
		return resp, nil
	})
}

func TestBhResource_UpdateWebAccessEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bhClient := &bhv20230418.Client{}
	patches.ApplyMethodReturn(newMockMetaBhResource().client, "UseBhV20230418Client", bhClient)

	enableCalled := false
	patches.ApplyMethodFunc(bhClient, "EnableWebAccessWithContext", func(_ context.Context, request *bhv20230418.EnableWebAccessRequest) (*bhv20230418.EnableWebAccessResponse, error) {
		assert.NotNil(t, request.ResourceId)
		assert.Equal(t, "bh-saas-abc123", *request.ResourceId)
		enableCalled = true
		resp := bhv20230418.NewEnableWebAccessResponse()
		resp.Response = &bhv20230418.EnableWebAccessResponseParams{
			RequestId: ptrStringBhRes("fake-request-id-enable"),
		}
		return resp, nil
	})

	mockDescribeResourcesForAccess(patches, bhClient, 1, 1)

	oldState := map[string]string{
		"id":            "bh-saas-abc123",
		"web_access":    "0",
		"client_access": "1",
	}
	newConfig := map[string]interface{}{
		"web_access": 1,
	}
	d := buildBhResourceData(t, oldState, newConfig)

	meta := newMockMetaBhResource()
	err := bh.ResourceTencentCloudBhResource().Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, enableCalled)
}

func TestBhResource_UpdateWebAccessDisable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bhClient := &bhv20230418.Client{}
	patches.ApplyMethodReturn(newMockMetaBhResource().client, "UseBhV20230418Client", bhClient)

	disableCalled := false
	patches.ApplyMethodFunc(bhClient, "DisableWebAccessWithContext", func(_ context.Context, request *bhv20230418.DisableWebAccessRequest) (*bhv20230418.DisableWebAccessResponse, error) {
		assert.NotNil(t, request.ResourceId)
		assert.Equal(t, "bh-saas-abc123", *request.ResourceId)
		disableCalled = true
		resp := bhv20230418.NewDisableWebAccessResponse()
		resp.Response = &bhv20230418.DisableWebAccessResponseParams{
			RequestId: ptrStringBhRes("fake-request-id-disable"),
		}
		return resp, nil
	})

	mockDescribeResourcesForAccess(patches, bhClient, 0, 1)

	oldState := map[string]string{
		"id":            "bh-saas-abc123",
		"web_access":    "1",
		"client_access": "1",
	}
	newConfig := map[string]interface{}{
		"web_access": 0,
	}
	d := buildBhResourceData(t, oldState, newConfig)

	meta := newMockMetaBhResource()
	err := bh.ResourceTencentCloudBhResource().Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, disableCalled)
}

func TestBhResource_UpdateClientAccessEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bhClient := &bhv20230418.Client{}
	patches.ApplyMethodReturn(newMockMetaBhResource().client, "UseBhV20230418Client", bhClient)

	enableCalled := false
	patches.ApplyMethodFunc(bhClient, "EnableClientTcpAccessWithContext", func(_ context.Context, request *bhv20230418.EnableClientTcpAccessRequest) (*bhv20230418.EnableClientTcpAccessResponse, error) {
		assert.NotNil(t, request.ResourceId)
		assert.Equal(t, "bh-saas-abc123", *request.ResourceId)
		enableCalled = true
		resp := bhv20230418.NewEnableClientTcpAccessResponse()
		resp.Response = &bhv20230418.EnableClientTcpAccessResponseParams{
			RequestId: ptrStringBhRes("fake-request-id-enable"),
		}
		return resp, nil
	})

	mockDescribeResourcesForAccess(patches, bhClient, 1, 1)

	oldState := map[string]string{
		"id":            "bh-saas-abc123",
		"web_access":    "1",
		"client_access": "0",
	}
	newConfig := map[string]interface{}{
		"client_access": 1,
	}
	d := buildBhResourceData(t, oldState, newConfig)

	meta := newMockMetaBhResource()
	err := bh.ResourceTencentCloudBhResource().Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, enableCalled)
}

func TestBhResource_UpdateClientAccessDisable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bhClient := &bhv20230418.Client{}
	patches.ApplyMethodReturn(newMockMetaBhResource().client, "UseBhV20230418Client", bhClient)

	disableCalled := false
	patches.ApplyMethodFunc(bhClient, "DisableClientTcpAccessWithContext", func(_ context.Context, request *bhv20230418.DisableClientTcpAccessRequest) (*bhv20230418.DisableClientTcpAccessResponse, error) {
		assert.NotNil(t, request.ResourceId)
		assert.Equal(t, "bh-saas-abc123", *request.ResourceId)
		disableCalled = true
		resp := bhv20230418.NewDisableClientTcpAccessResponse()
		resp.Response = &bhv20230418.DisableClientTcpAccessResponseParams{
			RequestId: ptrStringBhRes("fake-request-id-disable"),
		}
		return resp, nil
	})

	mockDescribeResourcesForAccess(patches, bhClient, 1, 0)

	oldState := map[string]string{
		"id":            "bh-saas-abc123",
		"web_access":    "1",
		"client_access": "1",
	}
	newConfig := map[string]interface{}{
		"client_access": 0,
	}
	d := buildBhResourceData(t, oldState, newConfig)

	meta := newMockMetaBhResource()
	err := bh.ResourceTencentCloudBhResource().Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, disableCalled)
}
