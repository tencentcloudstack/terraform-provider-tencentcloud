package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_ccn", &resource.Sweeper{
		Name: "tencentcloud_ccn",
		F:    testSweepCcnInstance,
	})
}

func testSweepCcnInstance(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(*TencentCloudClient)

	vpcService := VpcService{
		client: client.apiV3Conn,
	}

	instances, err := vpcService.DescribeCcns(ctx, "", "")
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {
		instanceId := v.ccnId
		instanceName := v.name

		now := time.Now()

		createTime := stringTotime(v.createTime)
		interval := now.Sub(createTime).Minutes()
		if instanceName != "" {
			if strings.HasPrefix(instanceName, keepResource) || strings.HasPrefix(instanceName, defaultResource) {
				continue
			}
		}

		// less than 30 minute, not delete
		if needProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = vpcService.DeleteCcn(ctx, instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudCcnV3Basic(t *testing.T) {
	t.Parallel()
	keyName := "tencentcloud_ccn.main"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCcnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnExists(keyName),
					resource.TestCheckResourceAttr(keyName, "name", "ci-temp-test-ccn"),
					resource.TestCheckResourceAttr(keyName, "description", "ci-temp-test-ccn-des"),
					resource.TestCheckResourceAttr(keyName, "instance_count", "0"),
					resource.TestCheckResourceAttr(keyName, "qos", "AG"),
					resource.TestCheckResourceAttrSet(keyName, "state"),
					resource.TestCheckResourceAttrSet(keyName, "create_time"),
				),
			},
			{
				ResourceName:      keyName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add tag
			{
				Config: testAccCcn_multiTags("master"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCcnExists(keyName),
					resource.TestCheckResourceAttr(keyName, "tags.role", "master"),
				),
			},
			// update tag
			{
				Config: testAccCcn_multiTags("master-version2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCcnExists(keyName),
					resource.TestCheckResourceAttr(keyName, "tags.role", "master-version2"),
				),
			},
			// remove tag
			{
				Config: testAccCcnConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCcnExists(keyName),
					resource.TestCheckNoResourceAttr(keyName, "tags.role"),
				),
			},
		},
	})
}

func TestAccTencentCloudCcnV3Update(t *testing.T) {
	t.Parallel()
	keyName := "tencentcloud_ccn.main"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCcnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnExists(keyName),
					resource.TestCheckResourceAttr(keyName, "name", "ci-temp-test-ccn"),
					resource.TestCheckResourceAttr(keyName, "description", "ci-temp-test-ccn-des"),
					resource.TestCheckResourceAttr(keyName, "instance_count", "0"),
					resource.TestCheckResourceAttr(keyName, "qos", "AG"),
					resource.TestCheckResourceAttrSet(keyName, "state"),
					resource.TestCheckResourceAttrSet(keyName, "create_time"),
				),
			},
			{
				Config: testAccCcnConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCcnExists(keyName),
					resource.TestCheckResourceAttr(keyName, "name", "ci-temp-test-ccn-update"),
					resource.TestCheckResourceAttr(keyName, "description", "ci-temp-test-ccn-des-update"),
					resource.TestCheckResourceAttr(keyName, "instance_count", "0"),
					resource.TestCheckResourceAttr(keyName, "qos", "AG"),
					resource.TestCheckResourceAttrSet(keyName, "state"),
					resource.TestCheckResourceAttrSet(keyName, "create_time"),
				),
			},
		},
	})
}

func testAccCheckCcnExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeCcn(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has > 0 {
			return nil
		}
		return fmt.Errorf("ccn not exists.")
	}
}

func testAccCheckCcnDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ccn" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, has, err := service.DescribeCcn(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has == 0 {
			return nil
		}
		return fmt.Errorf("ccn not delete ok")
	}
	return nil
}

func testAccCcn_multiTags(value string) string {
	return fmt.Sprintf(
		`
resource tencentcloud_ccn main {
	name        = "ci-temp-test-ccn"
	description = "ci-temp-test-ccn-des"
	qos         = "AG"
	tags = {
		role = "%s"
	}
}
`, value)
}

const testAccCcnConfig = `
resource tencentcloud_ccn main {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
}
`

const testAccCcnConfigUpdate = `
resource tencentcloud_ccn main {
  name                 = "ci-temp-test-ccn-update"
  description          = "ci-temp-test-ccn-des-update"
  qos                  = "AG"
}
`
