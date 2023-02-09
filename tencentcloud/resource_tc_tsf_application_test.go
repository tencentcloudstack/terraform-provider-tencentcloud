package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationConfigResource_basic -v
func TestAccTencentCloudTsfApplicationResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApplicationExists("tencentcloud_tsf_application.application"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_application.application", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "application_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "application_type", "V"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "microservice_type", "N"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "application_desc", "application desc"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.0.name", "service-name"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.ports.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.ports.0.target_port", "9090"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.ports.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.health_check.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.health_check.0.path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_application.application",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfApplicationDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_application" {
			continue
		}

		res, err := service.DescribeTsfApplicationById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf application %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfApplicationExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfApplicationById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf application %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfApplication = `

resource "tencentcloud_tsf_application" "application" {
	application_name = "terraform-test"
	application_type = "V"
	microservice_type = "N"
	application_desc = "application desc"
	application_runtime_type = ""
	# application_remark_name = "remark-name"
	program_id = ""
  	service_config_list {
		name = "service-name"
		ports {
			target_port = 9090
			protocol = "TCP"
		}
		health_check {
			path = "/"
		}

  	}
	ignore_create_image_repository = true
	# program_id_list =

	tags = {
		"createdBy" = "terraform"
	}
}

`
