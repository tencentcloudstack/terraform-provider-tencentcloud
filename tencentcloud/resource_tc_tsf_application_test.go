package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApplicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplication,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_application.application", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_application.application",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplication = `

resource "tencentcloud_tsf_application" "application" {
  application_name = "my-app"
  application_type = "C"
  microservice_type = "M"
  application_desc = "This is my application"
  application_log_config = ""
  application_resource_type = ""
  application_runtime_type = "Java"
  program_id = "p-123456"
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
  program_id_list = 
}

`
