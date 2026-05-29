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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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

// mockMetaClbTargetGroup implements tccommon.ProviderMeta for unit tests
type mockMetaClbTargetGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaClbTargetGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaClbTargetGroup{}

func newMockMetaClbTargetGroup() *mockMetaClbTargetGroup {
	return &mockMetaClbTargetGroup{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/clb/ -run "TestUnitClbTargetGroup_CreateWithSnatEnable" -v -count=1 -gcflags="all=-l"
func TestUnitClbTargetGroup_CreateWithSnatEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaClbTargetGroup().client, "UseClbClient", clbClient)

	// Patch CreateTargetGroup to return success
	patches.ApplyMethodFunc(clbClient, "CreateTargetGroup", func(request *clb.CreateTargetGroupRequest) (*clb.CreateTargetGroupResponse, error) {
		assert.NotNil(t, request.SnatEnable)
		assert.True(t, *request.SnatEnable)

		resp := &clb.CreateTargetGroupResponse{}
		resp.Response = &clb.CreateTargetGroupResponseParams{
			TargetGroupId: helper.String("lbtg-test123"),
			RequestId:     helper.String("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeTargetGroupList for the Read call after Create
	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := &clb.DescribeTargetGroupListResponse{}
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: helper.IntUint64(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   helper.String("lbtg-test123"),
					TargetGroupName: helper.String("test-snat"),
					VpcId:           helper.String("vpc-xxxxxx"),
					Port:            helper.IntUint64(80),
					TargetGroupType: helper.String("v2"),
					Protocol:        helper.String("TCP"),
					SnatEnable:      helper.Bool(true),
				},
			},
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaClbTargetGroup()
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
	assert.Equal(t, "lbtg-test123", d.Id())
	assert.Equal(t, true, d.Get("snat_enable").(bool))
}

// go test ./tencentcloud/services/clb/ -run "TestUnitClbTargetGroup_CreateWithoutSnatEnable" -v -count=1 -gcflags="all=-l"
func TestUnitClbTargetGroup_CreateWithoutSnatEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaClbTargetGroup().client, "UseClbClient", clbClient)

	// Patch CreateTargetGroup - SnatEnable should not be set
	patches.ApplyMethodFunc(clbClient, "CreateTargetGroup", func(request *clb.CreateTargetGroupRequest) (*clb.CreateTargetGroupResponse, error) {
		assert.Nil(t, request.SnatEnable)

		resp := &clb.CreateTargetGroupResponse{}
		resp.Response = &clb.CreateTargetGroupResponseParams{
			TargetGroupId: helper.String("lbtg-test456"),
			RequestId:     helper.String("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeTargetGroupList for the Read call after Create
	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := &clb.DescribeTargetGroupListResponse{}
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: helper.IntUint64(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   helper.String("lbtg-test456"),
					TargetGroupName: helper.String("test-no-snat"),
					VpcId:           helper.String("vpc-xxxxxx"),
					Port:            helper.IntUint64(80),
					TargetGroupType: helper.String("v2"),
					Protocol:        helper.String("TCP"),
					SnatEnable:      helper.Bool(false),
				},
			},
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaClbTargetGroup()
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
	assert.Equal(t, "lbtg-test456", d.Id())
	assert.Equal(t, false, d.Get("snat_enable").(bool))
}

// go test ./tencentcloud/services/clb/ -run "TestUnitClbTargetGroup_UpdateSnatEnable" -v -count=1 -gcflags="all=-l"
func TestUnitClbTargetGroup_UpdateSnatEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaClbTargetGroup().client, "UseClbClient", clbClient)

	// Patch ModifyTargetGroupAttribute
	patches.ApplyMethodFunc(clbClient, "ModifyTargetGroupAttribute", func(request *clb.ModifyTargetGroupAttributeRequest) (*clb.ModifyTargetGroupAttributeResponse, error) {
		assert.NotNil(t, request.SnatEnable)
		assert.True(t, *request.SnatEnable)
		assert.Equal(t, "lbtg-test789", *request.TargetGroupId)

		resp := &clb.ModifyTargetGroupAttributeResponse{}
		resp.Response = &clb.ModifyTargetGroupAttributeResponseParams{
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeTargetGroupList for the Read call after Update
	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := &clb.DescribeTargetGroupListResponse{}
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: helper.IntUint64(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   helper.String("lbtg-test789"),
					TargetGroupName: helper.String("test-update-snat"),
					VpcId:           helper.String("vpc-xxxxxx"),
					Port:            helper.IntUint64(80),
					TargetGroupType: helper.String("v2"),
					Protocol:        helper.String("TCP"),
					SnatEnable:      helper.Bool(true),
				},
			},
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaClbTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"snat_enable": true,
	})
	d.SetId("lbtg-test789")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, true, d.Get("snat_enable").(bool))
}

// go test ./tencentcloud/services/clb/ -run "TestUnitClbTargetGroup_ReadSnatEnable" -v -count=1 -gcflags="all=-l"
func TestUnitClbTargetGroup_ReadSnatEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clbClient := &clb.Client{}
	patches.ApplyMethodReturn(newMockMetaClbTargetGroup().client, "UseClbClient", clbClient)

	// Patch DescribeTargetGroupList
	patches.ApplyMethodFunc(clbClient, "DescribeTargetGroupList", func(request *clb.DescribeTargetGroupListRequest) (*clb.DescribeTargetGroupListResponse, error) {
		resp := &clb.DescribeTargetGroupListResponse{}
		resp.Response = &clb.DescribeTargetGroupListResponseParams{
			TotalCount: helper.IntUint64(1),
			TargetGroupSet: []*clb.TargetGroupInfo{
				{
					TargetGroupId:   helper.String("lbtg-testread"),
					TargetGroupName: helper.String("test-read-snat"),
					VpcId:           helper.String("vpc-xxxxxx"),
					Port:            helper.IntUint64(80),
					TargetGroupType: helper.String("v2"),
					Protocol:        helper.String("TCP"),
					SnatEnable:      helper.Bool(true),
				},
			},
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaClbTargetGroup()
	res := localclb.ResourceTencentCloudClbTargetGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("lbtg-testread")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, true, d.Get("snat_enable").(bool))
	assert.Equal(t, "test-read-snat", d.Get("target_group_name").(string))
}
