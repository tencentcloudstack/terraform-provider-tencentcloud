package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudClsCosShipper_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsCosShipper,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsCosShipperExists("tencentcloud_cls_cos_shipper.shipper"),
					resource.TestCheckResourceAttr("tencentcloud_cls_cos_shipper.shipper", "shipper_name", "tf-shipper-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_cos_shipper.shipper",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClsCosShipperExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS Cos Shipper][Exists] check: CLS Cos Shipper %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS Cos Shipper][Exists] check: CLS Cos Shipper id is not set")
		}
		clsService := ClsService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := clsService.DescribeClsCosShipperById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("[CHECK][CLS Cos Shipper][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsCosShipper = `
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-shipper-test"
  tags        = {
    "test" = "test"
  }
}

resource "tencentcloud_cls_topic" "topic" {
  auto_split           = true
  logset_id            = tencentcloud_cls_logset.logset.id
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test"
  }
  topic_name           = "tf-shipper-test"
}

resource "tencentcloud_cls_cos_shipper" "shipper" {
  bucket       = "preset-scf-bucket-1308919341"
  interval     = 300
  max_size     = 200
  partition    = "/%Y/%m/%d/%H/"
  prefix       = "ap-guangzhou-fffsasad-1649734752"
  shipper_name = "tf-shipper-test"
  topic_id     = tencentcloud_cls_topic.topic.id

  compress {
    format = "lzop"
  }

  content {
    format = "json"

    json {
      enable_tag  = true
      meta_fields = [
        "__FILENAME__",
        "__SOURCE__",
        "__TIMESTAMP__",
      ]
    }
  }
}


`
