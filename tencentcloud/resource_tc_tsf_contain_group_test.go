package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudNeedFixTsfContainGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfContainGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfContainGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfContainGroupExists("tencentcloud_tsf_contain_group.contain_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_contain_group.contain_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_contain_group.contain_group", "application_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_tsf_contain_group.contain_group", "namespace_id", ""),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_contain_group.contain_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfContainGroupDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_contain_group" {
			continue
		}

		res, err := service.DescribeTsfLaneRuleById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf ContainGroup %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfContainGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfLaneRuleById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf ContainGroup %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfContainGroup = `

resource "tencentcloud_tsf_contain_group" "contain_group" {
  application_id = ""
  namespace_id = ""
  group_name = ""
  instance_num = 
  access_type = 
  protocol_ports {
		protocol = ""
		port = 
		target_port = 
		node_port = 

  }
  cluster_id = ""
  cpu_limit = ""
  mem_limit = ""
  group_comment = ""
  update_type = 
  update_ivl = 
  cpu_request = ""
  mem_request = ""
  group_resource_type = ""
  subnet_id = ""
  agent_cpu_request = ""
  agent_cpu_limit = ""
  agent_mem_request = ""
  agent_mem_limit = ""
  istio_cpu_request = ""
  istio_cpu_limit = ""
  istio_mem_request = ""
  istio_mem_limit = ""
}

`
