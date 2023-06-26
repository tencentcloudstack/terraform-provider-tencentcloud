package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-shanghai -sweep-run=tencentcloud_tcr_instance
	resource.AddTestSweepers("tencentcloud_tcr_instance", &resource.Sweeper{
		Name: "tencentcloud_tcr_instance",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := TCRService{client}

			instances, err := service.DescribeTCRInstances(ctx, "", nil)
			if err != nil {
				return err
			}

			for i := range instances {
				ins := instances[i]
				id := *ins.RegistryId
				name := *ins.RegistryName
				created, err := time.Parse(time.RFC3339, *ins.CreatedAt)
				if err != nil {
					created = time.Time{}
				}
				if isResourcePersist(name, &created) {
					continue
				}

				// Delete replicas

				// Delete replications first
				repRequest := tcr.NewDescribeReplicationInstancesRequest()
				repRequest.RegistryId = &id
				replicas, outErr := service.DescribeReplicationInstances(ctx, repRequest)

				if outErr != nil {
					return outErr
				}

				for i := range replicas {
					item := replicas[i]
					_ = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
						request := tcr.NewDeleteReplicationInstanceRequest()
						request.RegistryId = &id
						request.ReplicationRegistryId = item.ReplicationRegistryId
						request.ReplicationRegionId = item.ReplicationRegionId
						err := service.DeleteReplicationInstance(ctx, request)
						if err != nil {
							return retryError(err, tcr.INTERNALERROR_ERRORCONFLICT)
						}
						return nil
					})
				}

				// Delete Instance
				log.Printf("instance %s:%s will delete", id, name)
				err = service.DeleteTCRInstance(ctx, id, true)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudTcrInstanceResource_basic_and_update(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRInstance_basic,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "name", "testacctcrinstance1"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "instance_type", "basic"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "delete_bucket", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_instance.mytcr_instance", "internal_end_point"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_instance.mytcr_instance", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_instance.mytcr_instance", "public_domain"),
				),
			},
			{
				ResourceName:            "tencentcloud_tcr_instance.mytcr_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_bucket"},
			},
			{
				Config: testAccTCRInstance_basic_update_remark,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRInstanceExists("tencentcloud_tcr_instance.mytcr_instance"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "delete_bucket", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "open_public_operation", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "security_policy.#", "2"),
					// resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "security_policy.0.cidr_block", "192.168.1.1/24"),
					// resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "security_policy.1.cidr_block", "10.0.0.1/16"),
				),
			},
			{
				Config: testAccTCRInstance_basic_update_security,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRInstanceExists("tencentcloud_tcr_instance.mytcr_instance"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "open_public_operation", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "security_policy.#", "1"),
					// resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "security_policy.0.cidr_block", "192.168.1.1/24"),
				),
			},
			{
				Config: testAccTCRInstance_basic_update_instance_type,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRInstanceExists("tencentcloud_tcr_instance.mytcr_instance"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "instance_type", "premium"),
				),
			},
		},
	})
}

// Neet to fix because tcr deteleInstance api has issue
func TestAccTencentCloudNeedFixTcrInstanceResource_paypaid(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRInstance_paypaid,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "name", "paypaidtcrinstance"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "instance_type", "basic"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "registry_charge_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "instance_charge_type_prepaid_period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "instance_charge_type_prepaid_renew_flag", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_instance.mytcr_instance_paypaid", "expired_at"),
				),
			},
			{
				Config: testAccTCRInstance_update_paypaid_period,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "name", "paypaidtcrinstance"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "instance_type", "basic"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "registry_charge_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "instance_charge_type_prepaid_period", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance_paypaid", "instance_charge_type_prepaid_renew_flag", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_instance.mytcr_instance_paypaid", "expired_at"),
				),
			},
			{
				ResourceName:            "tencentcloud_tcr_instance.mytcr_instance_paypaid",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_bucket", "instance_charge_type_prepaid_period"},
			},
		},
	})
}

func TestAccTencentCloudTcrInstanceResource_replication(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRInstance_replica,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "name", "tfreplicas"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "replications.#", "2"),
				),
			},
			//{
			//	ResourceName:            "tencentcloud_tcr_instance.mytcr_instance",
			//	ImportState:             true,
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{"delete_bucket"},
			//},
			{
				Config: testAccTCRInstance_replica_update,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "name", "tfreplicas"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_instance.mytcr_instance", "replications.#", "3"),
				),
			},
		},
	})
}

func TestAccTencentCloudTcrInstanceResource_replica_set(t *testing.T) {
	inputs := []interface{}{
		map[string]interface{}{
			"region_id": 1,
			"sync_cos":  true,
		},
		map[string]interface{}{
			"region_id": 2,
		},
		map[string]interface{}{
			"region_id": 3,
			"sync_cos":  false,
		},
	}

	registries1 := []*tcr.ReplicationRegistry{
		{
			ReplicationRegistryId: helper.String("a"),
			RegistryId:            helper.String("x"),
			ReplicationRegionId:   helper.IntUint64(1),
		},
		{
			ReplicationRegistryId: helper.String("b"),
			RegistryId:            helper.String("x"),
			ReplicationRegionId:   helper.IntUint64(2),
		},
		{
			ReplicationRegistryId: helper.String("c"),
			RegistryId:            helper.String("x"),
			ReplicationRegionId:   helper.IntUint64(3),
		},
	}

	result1 := resourceTencentCloudTcrFillReplicas(inputs, registries1)
	expected1 := []interface{}{
		map[string]interface{}{
			"id":        "a",
			"region_id": 1,
			"sync_cos":  true,
		},
		map[string]interface{}{
			"id":        "b",
			"region_id": 2,
		},
		map[string]interface{}{
			"id":        "c",
			"region_id": 3,
			"sync_cos":  false,
		},
	}
	assert.Equalf(t, expected1, result1, "%s case 1 not equal, expected:\n%v\ngot: \n%v", t.Name(), expected1, result1)

	var registries2 []*tcr.ReplicationRegistry
	registries2Incr := []*tcr.ReplicationRegistry{
		{
			ReplicationRegistryId: helper.String("d"),
			RegistryId:            helper.String("x"),
			ReplicationRegionId:   helper.IntUint64(4),
		},
		{
			ReplicationRegistryId: helper.String("e"),
			RegistryId:            helper.String("x"),
			ReplicationRegionId:   helper.IntUint64(5),
		},
	}
	registries2 = append(registries2, registries1...)
	registries2 = append(registries2, registries2Incr...)
	result2 := resourceTencentCloudTcrFillReplicas(inputs, registries2)
	expected2 := []interface{}{
		map[string]interface{}{
			"id":        "a",
			"region_id": 1,
			"sync_cos":  true,
		},
		map[string]interface{}{
			"id":        "b",
			"region_id": 2,
		},
		map[string]interface{}{
			"id":        "c",
			"region_id": 3,
			"sync_cos":  false,
		},
		map[string]interface{}{
			"id":        "d",
			"region_id": 4,
		},
		map[string]interface{}{
			"id":        "e",
			"region_id": 5,
		},
	}

	assert.Equalf(t, expected2, result2, "%s case 2 not equal, expected:\n%v\ngot: \n%v", t.Name(), expected2, result2)
}

func testAccCheckTCRInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_instance" {
			continue
		}
		_, has, err := tcrService.DescribeTCRInstanceById(ctx, rs.Primary.ID)
		if has {
			return fmt.Errorf("TCR instance still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTCRInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("TCR instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("TCR instance id is not set")
		}

		tcrService := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, has, err := tcrService.DescribeTCRInstanceById(ctx, rs.Primary.ID)
		if !has {
			return fmt.Errorf("TCR instance %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTCRInstance_basic = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance1"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}`

const testAccTCRInstance_paypaid = `
resource "tencentcloud_tcr_instance" "mytcr_instance_paypaid" {
  name        = "paypaidtcrinstance"
  instance_type = "basic"
  delete_bucket = true
  registry_charge_type = 2
  instance_charge_type_prepaid_period = 1
  instance_charge_type_prepaid_renew_flag = 1
  tags ={
	test = "test"
  }
}`

const testAccTCRInstance_update_paypaid_period = `
resource "tencentcloud_tcr_instance" "mytcr_instance_paypaid" {
  name        = "paypaidtcrinstance"
  instance_type = "basic"
  delete_bucket = true
  registry_charge_type = 2
  instance_charge_type_prepaid_period = 2
  instance_charge_type_prepaid_renew_flag = 1
  tags ={
	test = "test"
  }
}`

const testAccTCRInstance_replica = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "tfreplicas"
  instance_type = "premium"
  delete_bucket = true

  replications {
    region_id = 1 # ap-guangzhou
  }
  replications {
    region_id = 5 # ap-hongkong
  }
}`

const testAccTCRInstance_replica_update = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "tfreplicas"
  instance_type = "premium"
  delete_bucket = true

  replications {
    region_id = 5 # ap-hongkong
  }
  replications {
    region_id = 8 # ap-beijing
  }
  replications {
    region_id = 16 #ap-chengdu
    syn_tag = true
  }
}`

const testAccTCRInstance_basic_update_remark = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance1"
  instance_type = "basic"
  delete_bucket = true
  open_public_operation = true
  security_policy {
    cidr_block = "192.168.1.1/24"
  }
  security_policy {
    cidr_block = "10.0.0.1/16"
  }

  tags ={
	test = "test"
  }
}`

const testAccTCRInstance_basic_update_security = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance1"
  instance_type = "basic"
  delete_bucket = true
  open_public_operation = true

  security_policy {
    cidr_block = "192.168.1.1/24"
  }

  tags ={
	test = "test"
  }
}
`

const testAccTCRInstance_basic_update_instance_type = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance1"
  instance_type = "premium"
  delete_bucket = true
  open_public_operation = true

  security_policy {
    cidr_block = "192.168.1.1/24"
  }

  tags ={
	test = "test"
  }
}
`
