package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTatCommandResource_basic -v
func TestAccTencentCloudTatCommandResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTatCommandeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTatCommand,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTatCommandExists("tencentcloud_tat_command.command"),
					resource.TestCheckResourceAttr("tencentcloud_tat_command.command", "username", "root"),
					resource.TestCheckResourceAttr("tencentcloud_tat_command.command", "command_name", "ls"),
					resource.TestCheckResourceAttr("tencentcloud_tat_command.command", "content", "ls"),
					resource.TestCheckResourceAttr("tencentcloud_tat_command.command", "description", "shell desc"),
					resource.TestCheckResourceAttr("tencentcloud_tat_command.command", "command_type", "SHELL"),
					resource.TestCheckResourceAttr("tencentcloud_tat_command.command", "working_directory", "/root"),
					resource.TestCheckResourceAttr("tencentcloud_tat_command.command", "timeout", "50"),
				),
			},
			{
				ResourceName:      "tencentcloud_tat_command.command",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTatCommandeDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TatService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tat_command" {
			continue
		}

		command, err := service.DescribeTatCommand(ctx, rs.Primary.ID)
		if command != nil {
			return fmt.Errorf("tat command %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTatCommandExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TatService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		command, err := service.DescribeTatCommand(ctx, rs.Primary.ID)
		if command == nil {
			return fmt.Errorf("tat command %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTatCommand = `

resource "tencentcloud_tat_command" "command" {
	username          = "root"
	command_name      = "ls"
	content           = "ls"
	description       = "shell desc"
	command_type      = "SHELL"
	working_directory = "/root"
	timeout = 50
}

`
