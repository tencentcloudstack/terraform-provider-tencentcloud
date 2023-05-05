package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationResource_basic -v
func TestAccTencentCloudTsfApplicationResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApplicationExists("tencentcloud_tsf_application.application"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_application.application", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "application_desc", "This is my application"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "application_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "application_resource_type", "DEF"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "application_runtime_type", "Java"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "application_type", "C"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "ignore_create_image_repository", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "microservice_type", "M"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.0.name", "my-service"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.0.health_check.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.0.health_check.0.path", "/health"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.0.ports.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.0.ports.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_application.application", "service_config_list.0.ports.0.target_port", "8080"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tsf_application.application",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckTsfApplicationDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_application_release_config" {
			continue
		}

		res, err := service.DescribeTsfApplicationById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf Application %s still exists", rs.Primary.ID)
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
			return fmt.Errorf("tsf Application %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfApplication = `

resource "tencentcloud_tsf_application" "application" {
	application_name = "terraform-test"
	application_type = "C"
	microservice_type = "M"
	application_desc = "This is my application"
	application_runtime_type = "Java"
	service_config_list {
	  name = "my-service"
	  ports {
		target_port = 8080
		protocol = "HTTP"
	  }
	  health_check {
		path = "/health"
	  }
	}
	ignore_create_image_repository = true
}

`
