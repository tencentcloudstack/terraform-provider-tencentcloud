package clb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"

	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_clb_target_group
	resource.AddTestSweepers("tencentcloud_clb_target_group", &resource.Sweeper{
		Name: "tencentcloud_clb_target_group",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := localclb.NewClbService(client)

			tgs, err := service.DescribeTargetGroups(ctx, "", nil)
			if err != nil {
				return err
			}

			// add scanning resources
			var resources, nonKeepResources []*tccommon.ResourceInstance
			for _, v := range tgs {
				if !tccommon.CheckResourcePersist(*v.TargetGroupName, *v.CreatedTime) {
					nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
						Id:   *v.TargetGroupId,
						Name: *v.TargetGroupName,
					})
				}
				resources = append(resources, &tccommon.ResourceInstance{
					Id:         *v.TargetGroupId,
					Name:       *v.TargetGroupName,
					CreateTime: *v.CreatedTime,
				})
			}
			tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateTargetGroup")

			for i := range tgs {
				tg := tgs[i]
				created := tccommon.ParseTimeFromCommonLayout(tg.CreatedTime)
				if tcacctest.IsResourcePersist(*tg.TargetGroupName, &created) {
					continue
				}
				log.Printf("%s will be remvoed", *tg.TargetGroupName)
				err = service.DeleteTarget(ctx, *tg.TargetGroupId)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudClbTargetGroup_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbTargetGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group.test", "target_group_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group.test", "vpc_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbInstanceTargetGroup(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstanceTargetGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.target_group"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_name", "tgt_grp_test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "port", "33"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.bind_ip", "10.0.0.4"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.port", "33"),
				),
			},
			{
				Config: testAccClbInstanceTargetGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.target_group"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_name", "tgt_grp_test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "port", "44"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.bind_ip", "10.0.0.4"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.port", "44"),
				),
			},
		},
	})
}

func testAccCheckClbTargetGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_target_group" {
			continue
		}
		time.Sleep(5 * time.Second)
		filters := map[string]string{}
		targetGroupInfos, err := clbService.DescribeTargetGroups(ctx, rs.Primary.ID, filters)
		if len(targetGroupInfos) > 0 && err == nil {
			return fmt.Errorf("[CHECK][CLB target group][Destroy] check: CLB target group still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbTargetGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB target group][Exists] check: CLB target group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB target group][Exists] check: CLB target group id is not set")
		}
		clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		filters := map[string]string{}
		targetGroupInfos, err := clbService.DescribeTargetGroups(ctx, rs.Primary.ID, filters)
		if err != nil {
			return err
		}
		if len(targetGroupInfos) == 0 {
			return fmt.Errorf("[CHECK][CLB target group][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbTargetGroup_basic = `
resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "qwe"
}
`

const testAccClbInstanceTargetGroup = `
resource "tencentcloud_clb_target_group" "target_group" {
    target_group_name = "tgt_grp_test"
    port              = 33
    vpc_id            = "vpc-4owdpnwr"
    target_group_instances {
      bind_ip = "172.16.16.95"
      port = 18800
    }
}
`

const testAccClbInstanceTargetGroupUpdate = `
resource "tencentcloud_clb_target_group" "target_group" {
    target_group_name = "tgt_grp_test"
    port              = 44
	vpc_id            = "vpc-4owdpnwr"
    target_group_instances {
      bind_ip = "172.16.16.95"
      port = 18800
    }
}
`

// ---- Unit Tests (gomonkey mock) ----

// mockMetaForClbTargetGroup implements tccommon.ProviderMeta
type mockMetaForClbTargetGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForClbTargetGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForClbTargetGroup{}

func newMockMetaForClbTargetGroup() *mockMetaForClbTargetGroup {
	return &mockMetaForClbTargetGroup{client: &connectivity.TencentCloudClient{}}
}

func ptrStringClbTargetGroup(s string) *string {
	return &s
}

func ptrBoolClbTargetGroup(b bool) *bool {
	return &b
}

func ptrUint64ClbTargetGroup(v uint64) *uint64 {
	return &v
}

// go test ./tencentcloud/services/clb/ -run "TestClbTargetGroup_" -v -count=1 -gcflags="all=-l"

// TestClbTargetGroup_Create_WithSnatEnable tests that snat_enable is passed to CreateTargetGroup API
func TestClbTargetGroup_Create_WithSnatEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaForClbTargetGroup().client, "UseClbClient", clbClient)

	var capturedRequest *clb.CreateTargetGroupRequest
	patches.ApplyMethodFunc(clbClient, "CreateTargetGroup", func(request *clb.CreateTargetGroupRequest) (*clb.CreateTargetGroupResponse, error) {
		capturedRequest = request
		resp := clb.NewCreateTargetGroupResponse()
		resp.Response = &clb.CreateTargetGroupResponseParams{
			TargetGroupId: ptrStringClbTargetGroup("lbtg-test12345"),
		}
		return resp, nil
	})

	// Mock DescribeTargetGroupList for the Read call after Create
	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := clb.NewDescribeTargetGroupListResponse()
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: ptrUint64ClbTargetGroup(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   ptrStringClbTargetGroup("lbtg-test12345"),
					TargetGroupName: ptrStringClbTargetGroup("test-snat"),
					VpcId:           ptrStringClbTargetGroup("vpc-xxxxxx"),
					Port:            ptrUint64ClbTargetGroup(80),
					TargetGroupType: ptrStringClbTargetGroup("v2"),
					Protocol:        ptrStringClbTargetGroup("TCP"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaForClbTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"target_group_name": "test-snat",
		"vpc_id":            "vpc-xxxxxx",
		"port":              80,
		"type":              "v2",
		"protocol":          "TCP",
		"snat_enable":       true,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lbtg-test12345", d.Id())

	// Verify SnatEnable was passed to the API
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.SnatEnable)
	assert.Equal(t, true, *capturedRequest.SnatEnable)
}

// TestClbTargetGroup_Create_WithoutSnatEnable tests that snat_enable is not passed when not specified
func TestClbTargetGroup_Create_WithoutSnatEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaForClbTargetGroup().client, "UseClbClient", clbClient)

	var capturedRequest *clb.CreateTargetGroupRequest
	patches.ApplyMethodFunc(clbClient, "CreateTargetGroup", func(request *clb.CreateTargetGroupRequest) (*clb.CreateTargetGroupResponse, error) {
		capturedRequest = request
		resp := clb.NewCreateTargetGroupResponse()
		resp.Response = &clb.CreateTargetGroupResponseParams{
			TargetGroupId: ptrStringClbTargetGroup("lbtg-test67890"),
		}
		return resp, nil
	})

	// Mock DescribeTargetGroupList for the Read call after Create
	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := clb.NewDescribeTargetGroupListResponse()
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: ptrUint64ClbTargetGroup(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   ptrStringClbTargetGroup("lbtg-test67890"),
					TargetGroupName: ptrStringClbTargetGroup("test-no-snat"),
					VpcId:           ptrStringClbTargetGroup("vpc-xxxxxx"),
					Port:            ptrUint64ClbTargetGroup(80),
					TargetGroupType: ptrStringClbTargetGroup("v2"),
					Protocol:        ptrStringClbTargetGroup("TCP"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaForClbTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"target_group_name": "test-no-snat",
		"vpc_id":            "vpc-xxxxxx",
		"port":              80,
		"type":              "v2",
		"protocol":          "TCP",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lbtg-test67890", d.Id())

	// Verify SnatEnable was NOT passed to the API
	assert.NotNil(t, capturedRequest)
	assert.Nil(t, capturedRequest.SnatEnable)
}

// TestClbTargetGroup_Update_SnatEnable tests that snat_enable change triggers ModifyTargetGroupAttribute
func TestClbTargetGroup_Update_SnatEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaForClbTargetGroup().client, "UseClbClient", clbClient)

	var capturedModifyRequest *clb.ModifyTargetGroupAttributeRequest
	patches.ApplyMethodFunc(clbClient, "ModifyTargetGroupAttribute", func(request *clb.ModifyTargetGroupAttributeRequest) (*clb.ModifyTargetGroupAttributeResponse, error) {
		capturedModifyRequest = request
		resp := clb.NewModifyTargetGroupAttributeResponse()
		resp.Response = &clb.ModifyTargetGroupAttributeResponseParams{}
		return resp, nil
	})

	// Mock DescribeTargetGroupList for the Read call after Update
	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := clb.NewDescribeTargetGroupListResponse()
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: ptrUint64ClbTargetGroup(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   ptrStringClbTargetGroup("lbtg-update123"),
					TargetGroupName: ptrStringClbTargetGroup("test-update-snat"),
					VpcId:           ptrStringClbTargetGroup("vpc-xxxxxx"),
					Port:            ptrUint64ClbTargetGroup(80),
					TargetGroupType: ptrStringClbTargetGroup("v2"),
					Protocol:        ptrStringClbTargetGroup("TCP"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaForClbTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()

	// Only set snat_enable to simulate a change on that field only.
	// Do NOT set immutable fields (vpc_id, full_listen_switch, ip_version)
	// to avoid triggering the immutable field check.
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"target_group_name": "test-update-snat",
		"port":              80,
		"snat_enable":       true,
	})
	d.SetId("lbtg-update123")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify ModifyTargetGroupAttribute was called with SnatEnable
	assert.NotNil(t, capturedModifyRequest)
	assert.NotNil(t, capturedModifyRequest.SnatEnable)
	assert.Equal(t, true, *capturedModifyRequest.SnatEnable)
}

// TestClbTargetGroup_Schema_SnatEnable tests that snat_enable is defined correctly in schema
func TestClbTargetGroup_Schema_SnatEnable(t *testing.T) {
	res := localclb.ResourceTencentCloudClbTargetGroup()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "snat_enable")

	snatEnableSchema := res.Schema["snat_enable"]
	assert.Equal(t, schema.TypeBool, snatEnableSchema.Type)
	assert.True(t, snatEnableSchema.Optional)
	assert.False(t, snatEnableSchema.Required)
	assert.False(t, snatEnableSchema.Computed)
}
