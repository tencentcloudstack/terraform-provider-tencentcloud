package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTsfInstancesAttachmentResource_basic -v
func TestAccTencentCloudTsfInstancesAttachmentResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfInstancesAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfInstancesAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfInstancesAttachmentExists("tencentcloud_tsf_instances_attachment.instances_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_instances_attachment.instances_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_instances_attachment.instances_attachment", "cluster_id", "cluster-ym9mxm3a"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_instances_attachment.instances_attachment", "instance_id", "ins-acfxjkeg"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_instances_attachment.instances_attachment", "instance_import_mode", "M"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_instances_attachment.instances_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfInstancesAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_instances_attachment" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		instanceId := idSplit[1]

		res, err := service.DescribeTsfInstancesAttachmentById(ctx, clusterId, instanceId)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == "ResourceNotFound.GroupNotExist" {
				return nil
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf instances attachment %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfInstancesAttachmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		instanceId := idSplit[1]

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfInstancesAttachmentById(ctx, clusterId, instanceId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf instances attachment %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfInstancesAttachment = `

resource "tencentcloud_tsf_instances_attachment" "instances_attachment" {
	cluster_id = "cluster-ym9mxm3a"
	instance_id = "ins-acfxjkeg"
	# os_name = ""
	# image_id = ""
	# password = ""
	# key_id = ""
	# sg_id = ""
	instance_import_mode = "M"
  }

`
