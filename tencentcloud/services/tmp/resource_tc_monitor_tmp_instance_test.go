package tmp_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorInstance_basic -v
func TestAccTencentCloudMonitorInstance_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMonInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("tencentcloud_monitor_tmp_instance.example"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.example", "instance_name", "tf-tmp-instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.example", "data_retention_time", "30"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "zone"),
				),
			},
			{
				Config: testInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("tencentcloud_monitor_tmp_instance.example"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.example", "instance_name", "tf-tmp-instance-update"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.example", "data_retention_time", "90"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "zone"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMonInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_instance" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		instance, err := service.DescribeMonitorTmpInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance != nil {
			status := strconv.FormatInt(*instance.InstanceStatus, 10)
			if strings.Contains("5,6,8,9", status) {
				return nil
			}
			return fmt.Errorf("instance %s still exists: %v", rs.Primary.ID, *instance.InstanceStatus)
		}
	}

	return nil
}

func testAccCheckInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := service.DescribeMonitorTmpInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil || *instance.InstanceStatus != 2 {
			return fmt.Errorf("instance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testInstance_basic = `
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_monitor_tmp_instance" "example" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 30
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraform"
  }
}
`

const testInstance_update = `
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_monitor_tmp_instance" "example" {
  instance_name       = "tf-tmp-instance-update"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 90
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraformUpdate"
  }
}
`
