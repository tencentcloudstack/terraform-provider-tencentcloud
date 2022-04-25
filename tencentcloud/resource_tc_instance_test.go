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
	resource.AddTestSweepers("tencentcloud_instance", &resource.Sweeper{
		Name: "tencentcloud_instance",
		F:    testSweepCvmInstance,
	})
}

func testSweepCvmInstance(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(*TencentCloudClient)

	cvmService := CvmService{
		client: client.apiV3Conn,
	}

	instances, err := cvmService.DescribeInstanceByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {
		instanceId := *v.InstanceId
		instanceName := *v.InstanceName
		now := time.Now()
		createTime := stringTotime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()

		if strings.HasPrefix(instanceName, keepResource) || strings.HasPrefix(instanceName, defaultResource) {
			continue
		}

		if needProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = cvmService.DeleteInstance(ctx, instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudInstanceBasic(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "private_ip"),
					resource.TestCheckResourceAttrSet(id, "vpc_id"),
					resource.TestCheckResourceAttrSet(id, "subnet_id"),
					resource.TestCheckResourceAttrSet(id, "project_id"),
				),
			},
			{
				Config: testAccTencentCloudInstanceModifyInstanceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "instance_type"),
				),
			},
			{
				ResourceName:            id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_monitor_service", "disable_security_service", "hostname", "password", "force_delete"},
			},
		},
	})
}

func TestAccTencentCloudInstanceWithDataDisk(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithDataDisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "system_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_snapshot_id", ""),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_size", "100"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithDataDiskUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "system_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "150"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_snapshot_id", ""),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_size", "150"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithNetwork(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithNetworkFalse("false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckNoResourceAttr(id, "public_ip"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithNetwork("true", 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "internet_max_bandwidth_out", "5"),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "public_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithPrivateIP(t *testing.T) {

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithPrivateIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithKeyPair(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithKeyPair("key_pair_0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "key_name"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithKeyPair("key_pair_1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "key_name"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithPassword(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithPassword("TF_test_123"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "password"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithPassword("TF_test_123456"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "password"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithImageLogin(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithImageLogin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "keep_image_login", "true"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithName(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithName(defaultInsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "instance_name", defaultInsName),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithName(defaultInsNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "instance_name", defaultInsNameUpdate),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithHostname(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithHostname,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "hostname", defaultInsName),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithSecurityGroup(t *testing.T) {
	t.Parallel()

	instanceId := "tencentcloud_instance.foo"
	securitygroupId := "tencentcloud_security_group.foo"
	securitygroupRuleFooId := "tencentcloud_security_group_rule.foo"
	securitygroupRuleBarId := "tencentcloud_security_group_rule.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: instanceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithSecurityGroup(`[tencentcloud_security_group.foo.id]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(instanceId),
					testAccCheckTencentCloudInstanceExists(instanceId),
					resource.TestCheckResourceAttr(instanceId, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(instanceId, "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet(securitygroupId, "id"),
					resource.TestCheckResourceAttr(securitygroupRuleFooId, "type", "ingress"),
					resource.TestCheckResourceAttr(securitygroupRuleFooId, "port_range", "80,8080"),
					resource.TestCheckResourceAttr(securitygroupRuleBarId, "type", "ingress"),
					resource.TestCheckResourceAttr(securitygroupRuleBarId, "port_range", "3000"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithSecurityGroup(`[
					tencentcloud_security_group.foo.id,
					tencentcloud_security_group.bar.id
				]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(instanceId),
					testAccCheckTencentCloudInstanceExists(instanceId),
					resource.TestCheckResourceAttr(instanceId, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(instanceId, "security_groups.#", "2"),
					resource.TestCheckResourceAttrSet(securitygroupId, "id"),
					resource.TestCheckResourceAttr(securitygroupRuleFooId, "type", "ingress"),
					resource.TestCheckResourceAttr(securitygroupRuleFooId, "port_range", "80,8080"),
					resource.TestCheckResourceAttr(securitygroupRuleBarId, "type", "ingress"),
					resource.TestCheckResourceAttr(securitygroupRuleBarId, "port_range", "3000"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithTags(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithTags(`{
					"hello" = "world"
					"happy" = "hour"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "tags.hello", "world"),
					resource.TestCheckResourceAttr(id, "tags.happy", "hour"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithTags(`{
					"hello" = "hello"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "tags.hello", "hello"),
					resource.TestCheckNoResourceAttr(id, "tags.happy"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithPlacementGroup(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithPlacementGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "placement_group_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithSpotpaid(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithSpotpaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstancePostpaidToPrepaid(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		IDRefreshName: id,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstancePostPaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
			{
				Config: testAccTencentCloudInstanceBasicToPrepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_period", "1"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
				),
			},
		},
	})
}

func testAccCheckTencentCloudInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cvm instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cvm instance id is not set")
		}

		cvmService := CvmService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				instance, err = cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("cvm instance id is not found")
		}
		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cvmService := CvmService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_instance" {
			continue
		}

		instance, err := cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				instance, err = cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance != nil && *instance.InstanceState != CVM_STATUS_SHUTDOWN && *instance.InstanceState != CVM_STATUS_TERMINATING {
			return fmt.Errorf("cvm instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccTencentCloudInstanceBasic = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = var.cvm_vpc_id
  subnet_id         = var.cvm_subnet_id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}
`

const testAccTencentCloudInstancePostPaid = `
data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "` + defaultInsName + `"
  availability_zone = "` + defaultAZone + `"
  image_id          = "` + defaultTkeOSImageId + `"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
}
`

const testAccTencentCloudInstanceBasicToPrepaid = `
data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "` + defaultInsName + `"
  availability_zone = "` + defaultAZone + `"
  image_id          = "` + defaultTkeOSImageId + `"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  instance_charge_type       = "PREPAID"
  instance_charge_type_prepaid_period = 1
  instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
}
`

const testAccTencentCloudInstanceModifyInstanceType = defaultInstanceVariable + `
data "tencentcloud_instance_types" "new_type" {
	availability_zone = var.availability_cvm_zone
  
	cpu_core_count = 2
	memory_size    = 2
  }

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.new_type.instance_types.0.instance_type
  vpc_id            = var.cvm_vpc_id
  subnet_id         = var.cvm_subnet_id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}
`

const testAccTencentCloudInstanceWithDataDisk = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.1.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 100

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 100
    delete_with_instance  = true
	// encrypt = true
  } 
   
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 100
    # data_disk_snapshot_id = "snap-nvzu3dmh"
    delete_with_instance  = true
  }

  disable_security_service = true
  disable_monitor_service  = true
}
`

const testAccTencentCloudInstanceWithDataDiskUpdate = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.1.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 100

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
	// encrypt = true
  } 
   
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    # data_disk_snapshot_id = "snap-nvzu3dmh"
    delete_with_instance  = true
  }

  disable_security_service = true
  disable_monitor_service  = true
}
`

func testAccTencentCloudInstanceWithNetworkFalse(hasPublicIp string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  allocate_public_ip         = %s
  system_disk_type           = "CLOUD_PREMIUM"
}
`,
		hasPublicIp,
	)
}

func testAccTencentCloudInstanceWithNetwork(hasPublicIp string, maxBandWidthOut int64) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  internet_max_bandwidth_out = %d
  allocate_public_ip         = %s
  system_disk_type           = "CLOUD_PREMIUM"
}
`,
		maxBandWidthOut, hasPublicIp,
	)
}

const testAccTencentCloudInstanceWithPrivateIP = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  vpc_id            = var.cvm_vpc_id
  subnet_id         = var.cvm_subnet_id
  private_ip        = "10.0.0.123"
}
`

func testAccTencentCloudInstanceWithKeyPair(keyName string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_key_pair" "key_pair_0" {
  key_name = "key_pair_0"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

resource "tencentcloud_key_pair" "key_pair_1" {
  key_name = "key_pair_1"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCzwYE6KI8uULEvSNA2k1tlsLtMDe+x1Saw6yL3V1mk9NFws0K2BshYqsnP/BlYiGZv/Nld5xmGoA9LupOcUpyyGGSHZdBrMx1Dz9ajewe7kGowRWwwMAHTlzh9+iqeg/v6P5vW6EwK4hpGWgv06vGs3a8CzfbHu1YRbZAO/ysp3ymdL+vGvw/vzC0T+YwPMisn9wFD5FTlJ+Em6s9PzxqR/41t4YssmCwUV78ZoYL8CyB0emuB8wALvcXbdUVxMxpBEHd5U6ZP5+HPxU2WFbWqiFCuErLIZRuxFw8L/Ot+JOyNnadN1XU4crYDX5cML1i/ExXKVIDoBaLtgAJOpyeP"
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  key_name          = tencentcloud_key_pair.%s.id
  system_disk_type  = "CLOUD_PREMIUM"
}
`,
		keyName,
	)
}

func testAccTencentCloudInstanceWithPassword(password string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  password                   = "%s"
  system_disk_type           = "CLOUD_PREMIUM"
}
`,
		password,
	)
}

const testAccTencentCloudInstanceWithImageLogin = defaultInstanceVariable + `
data "tencentcloud_images" "zoo" {
  image_type = ["PRIVATE_IMAGE"]
  os_name    = "centos"
}
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.zoo.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  keep_image_login 			 = true
  system_disk_type           = "CLOUD_PREMIUM"
}
`

func testAccTencentCloudInstanceWithName(instanceName string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name     = "%s"
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
}
`,
		instanceName,
	)
}

const testAccTencentCloudInstanceWithHostname = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  hostname          = var.instance_name
  system_disk_type  = "CLOUD_PREMIUM"
}
`

func testAccTencentCloudInstanceWithSecurityGroup(ids string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_security_group" "foo" {
  name        = var.instance_name
  description = var.instance_name
}

resource "tencentcloud_security_group_rule" "foo" {
  security_group_id = tencentcloud_security_group.foo.id
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}

resource "tencentcloud_security_group" "bar" {
  name        = var.instance_name
  description = var.instance_name
}

resource "tencentcloud_security_group_rule" "bar" {
  security_group_id = tencentcloud_security_group.bar.id
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "3000"
  policy            = "accept"
}

resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  security_groups            = %s
}
`,
		ids,
	)
}

func testAccTencentCloudInstanceWithTags(tags string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"

  tags = %s
}
`,
		tags,
	)
}

const testAccTencentCloudInstanceWithPlacementGroup = defaultInstanceVariable + `
resource "tencentcloud_placement_group" "foo" {
  name = var.instance_name
  type = "HOST"
}

resource "tencentcloud_instance" "foo" {
  instance_name      = var.instance_name
  availability_zone  = var.availability_cvm_zone
  image_id           = data.tencentcloud_images.default.images.0.image_id
  instance_type      = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type   = "CLOUD_PREMIUM"
  placement_group_id = tencentcloud_placement_group.foo.id
}
`

const testAccTencentCloudInstanceWithSpotpaid = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name        = var.instance_name
  availability_zone    = var.availability_cvm_zone
  image_id             = data.tencentcloud_images.default.images.0.image_id
  instance_type        = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type     = "CLOUD_PREMIUM"
  instance_charge_type = "SPOTPAID"
  spot_instance_type   = "ONE-TIME"
  spot_max_price       = "0.5"
}
`
