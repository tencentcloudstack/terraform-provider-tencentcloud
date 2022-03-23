package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testEmrClusterResourceKey = "tencentcloud_emr_cluster.emrrrr"

func TestAccTencentCloudEmrClusterResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEmrBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmrExists(testEmrClusterResourceKey),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "product_id", "4"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "display_strategy", "clusterList"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "vpc_settings.vpc_id", "vpc-i9zty99n"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "vpc_settings.subnet_id", "subnet-opwnh0sw"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "softwares.0", "zookeeper-3.6.1"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "support_ha", "0"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "instance_name", "emr-test-demo"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "login_settings.password", "tencent@cloud123"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "time_span", "3600"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "time_unit", "s"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "pay_mode", "0"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "placement.zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "placement.project_id", "0"),
					resource.TestCheckResourceAttrSet(testEmrClusterResourceKey, "instance_id"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "sg_id", "sg-qyy7jd2b"),
				),
			},
		},
	})
}

func testAccCheckEmrExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("emr cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("emr cluster id is not set")
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := EMRService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		instanceId := rs.Primary.ID
		clusters, err := service.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				clusters, err = service.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return nil
		}
		if len(clusters) <= 0 {
			return fmt.Errorf("emr cluster create fail")
		} else {
			log.Printf("[DEBUG]emr cluster  %s create  ok", rs.Primary.ID)
			return nil
		}

	}
}

const testEmrBasic = `
resource "tencentcloud_emr_cluster" "emrrrr" {
	product_id=4
	display_strategy="clusterList"
	vpc_settings={
	  vpc_id="vpc-i9zty99n"
	  subnet_id:"subnet-opwnh0sw"
	}
	softwares=[
	  "zookeeper-3.6.1",
  ]
	support_ha=0
	instance_name="emr-test-demo"
	resource_spec {
	  master_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.S2"
		storage_type=5
		root_size=50
	  }
	  core_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.S2"
		storage_type=5
		root_size=50
	  }
	  master_count=1
	  core_count=2
	}
	login_settings={
	  password="tencent@cloud123"
	}
	time_span=3600
	time_unit="s"
	pay_mode=0
	placement={
	  zone="ap-guangzhou-3"
	  project_id=0
	}
	sg_id="sg-qyy7jd2b"
  }
`
