package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudEmrClusterResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEmrBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmrExists(testTkeClusterResourceKey),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "product_id", "4"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "display_strategy", "clusterList"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "vpc_settings.vpc_id", "vpc-i9zty99n"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "vpc_settings.subnet_id", "subnet-opwnh0sw"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "softwares.0", "zookeeper-3.6.1"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "support_ha", "0"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "instance_name", "emr-test"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.master_resource_spec.mem_size", "8192"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.master_resource_spec.cpu", "4"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.master_resource_spec.disk_size", "100"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.master_resource_spec.disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.master_resource_spec.spec", "CVM.S2"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.master_resource_spec.storage_type", "5"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.core_resource_spec.mem_size", "8192"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.core_resource_spec.cpu", "4"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.core_resource_spec.disk_size", "100"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.core_resource_spec.disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.core_resource_spec.spec", "CVM.S2"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.core_resource_spec.storage_type", "5"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.master_count", "1"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "resource_spec.core_count", "2"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "login_settings.password", "tencent@cloud123"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "time_span", "3600"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "time_unit", "s"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "pay_mode", "0"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "placement.zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "placement.project_id", "0"),
					resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "instance_id"),
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
		params := make(map[string]interface{})
		instances := make([]string, 0)
		instances = append(instances, instanceId)
		params["instance_ids"] = instances

		clusters, err := service.DescribeInstances(ctx, params)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				clusters, err = service.DescribeInstances(ctx, params)
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
	  }
	  core_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.S2"
		storage_type=5
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
  }
`
