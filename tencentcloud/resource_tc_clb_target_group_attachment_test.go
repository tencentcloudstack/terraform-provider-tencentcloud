package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

const clbTargetGroupAttachment = "tencentcloud_clb_target_group_attachment.group"

func TestAccTencentClbTargetGroupAttachmentResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbTargetGroupAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroupAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupAttachmentExists(clbTargetGroupAttachment),
					resource.TestCheckResourceAttrSet(clbTargetGroupAttachment, "clb_id"),
					resource.TestCheckResourceAttrSet(clbTargetGroupAttachment, "listener_id"),
					resource.TestCheckResourceAttrSet(clbTargetGroupAttachment, "target_group_id"),
				),
			},
			{
				ResourceName:      clbTargetGroupAttachment,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentClbTargetGroupAttachmentHttpResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbTargetGroupAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroupAttachmentHttp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupAttachmentExists(clbTargetGroupAttachment),
					resource.TestCheckResourceAttrSet(clbTargetGroupAttachment, "clb_id"),
					resource.TestCheckResourceAttrSet(clbTargetGroupAttachment, "listener_id"),
					resource.TestCheckResourceAttrSet(clbTargetGroupAttachment, "target_group_id"),
				),
			},
		},
	})
}

func testAccCheckClbTargetGroupAttachmentDestroy(s *terraform.State) error {
	var (
		clbService = ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		targetInfos []*clb.TargetGroupInfo
		err         error
	)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_target_group_attachment" {
			continue
		}

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 4 {
			return fmt.Errorf("CLB target group attachment id is clb_id#listener_id#target_group_id#rule_id(only required for 7 layer CLB)")
		}

		targetInfos, err = clbService.DescribeTargetGroups(ctx, ids[0], nil)
		if err != nil {
			return err
		}
		for _, info := range targetInfos {
			for _, rule := range info.AssociatedRule {
				var originLocationId string
				originClbId := *rule.LoadBalancerId
				originListenerId := *rule.ListenerId
				if rule.LocationId != nil {
					originLocationId = *rule.LocationId
				}
				if originListenerId == ids[1] && originClbId == ids[2] && originLocationId == ids[3] {
					return fmt.Errorf("rule association target group instance still exist. [targetGroupId=%s, listenerId=%s, cldId=%s, ruleId=%s]",
						ids[0], ids[1], ids[2], ids[3])
				}
			}
		}
		return nil
	}
	return nil
}

func testAccCheckClbTargetGroupAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var (
			logId = getLogId(contextNil)
			ctx   = context.WithValue(context.TODO(), logIdKey, logId)
		)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("CLB target group attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("CLB target group attachment id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 4 {
			return fmt.Errorf("CLB target group attachment id is clb_id#listener_id#target_group_id#rule_id(only required for 7 layer CLB)")
		}

		has, err := clbService.DescribeAssociateTargetGroups(ctx, ids)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("CLB target group attachment not exist")
		}

		return nil
	}
}

const testAccClbTargetGroupAttachmentHttp = `
resource "tencentcloud_vpc" "foo" {
	name       = "guagua-ci-temp-test"
	cidr_block = "10.0.0.0/16"
  }

  resource "tencentcloud_clb_instance" "clb_basic" {
	network_type = "OPEN"
	clb_name     = "tf-clb-attach-basic"
	vpc_id       = tencentcloud_vpc.foo.id
  }

  resource "tencentcloud_clb_listener" "listener_basic" {
	clb_id        = tencentcloud_clb_instance.clb_basic.id
	port          = 1
	protocol      = "HTTP"
	listener_name = "listener_basic"
  }

  resource "tencentcloud_clb_listener_rule" "rule_basic" {
	clb_id              = tencentcloud_clb_instance.clb_basic.id
	listener_id         = tencentcloud_clb_listener.listener_basic.listener_id
	domain              = "abc.com"
	url                 = "/"
	session_expire_time = 30
	scheduler           = "WRR"
	target_type         = "TARGETGROUP"
  }

  resource "tencentcloud_clb_target_group" "test"{
	  target_group_name = "test-target-keep-1"
	  vpc_id            = tencentcloud_vpc.foo.id
  }

  resource "tencentcloud_clb_target_group_attachment" "group" {
	  clb_id          = tencentcloud_clb_instance.clb_basic.id
	  listener_id     = tencentcloud_clb_listener.listener_basic.listener_id
	  rule_id         = tencentcloud_clb_listener_rule.rule_basic.rule_id
	  target_group_id = tencentcloud_clb_target_group.test.id 
  }
`

const testAccClbTargetGroupAttachment = `
  resource "tencentcloud_vpc" "foo" {
	name       = "guagua-ci-temp-test"
	cidr_block = "10.0.0.0/16"
  }
  
  resource "tencentcloud_clb_instance" "clb_open" {
	network_type              = "OPEN"
	clb_name                  = "tf-clb-update-open"
	vpc_id                    = tencentcloud_vpc.foo.id
	project_id                = 0
	target_region_info_region = "ap-guangzhou"
	target_region_info_vpc_id = tencentcloud_vpc.foo.id
	tags = {
	  test = "test"
	}
  }
  
  resource "tencentcloud_clb_listener" "TCP_listener" {
	clb_id                     = tencentcloud_clb_instance.clb_open.id
	listener_name              = "test_listener"
	port                       = 80
	protocol                   = "TCP"
	health_check_switch        = true
	health_check_time_out      = 2
	health_check_interval_time = 5
	health_check_health_num    = 3
	health_check_unhealth_num  = 3
	session_expire_time        = 30
	scheduler                  = "WRR"
	target_type = "TARGETGROUP"
  }
  
  resource "tencentcloud_clb_target_group" "test"{
	target_group_name = "test-target-keep-1"
	vpc_id = tencentcloud_vpc.foo.id
  }
  
  resource "tencentcloud_clb_target_group_attachment" "group" {
	  clb_id          = tencentcloud_clb_instance.clb_open.id
	  listener_id     = tencentcloud_clb_listener.TCP_listener.listener_id
	  target_group_id = tencentcloud_clb_target_group.test.id 
  }
`
