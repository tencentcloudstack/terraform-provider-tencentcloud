package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

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
	ctx := context.WithValue(context.TODO(), "logId", logId)

	sharedClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(TencentCloudClient)

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
		if !strings.HasPrefix(instanceName, defaultInsName) {
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
					resource.TestCheckResourceAttr(id, "instance_type", "S2.SMALL2"),
				),
			},
			{
				ResourceName:            id,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_monitor_service", "disable_security_service", "hostname", "password", "allocate_public_ip"},
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
					resource.TestCheckResourceAttr(id, "system_disk_size", "50"),
					resource.TestCheckResourceAttr(id, "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_size", "100"),
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
				Config: testAccTencentCloudInstanceWithNetwork("false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckNoResourceAttr(id, "public_ip"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithNetwork("true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "public_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceWithPrivateIP(t *testing.T) {
	t.Parallel()

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
				Config: testAccTencentCloudInstanceWithKeyPair("tf_acc_test_key1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "key_name"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithKeyPair("tf_acc_test_key2"),
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

func testAccCheckTencentCloudInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

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
	ctx := context.WithValue(context.TODO(), "logId", logId)
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
		if instance != nil {
			return fmt.Errorf("cvm instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccTencentCloudInstanceBasic = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}
`

const testAccTencentCloudInstanceModifyInstanceType = defaultInstanceVariable + `
data "tencentcloud_instance_types" "new_type" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }

  cpu_core_count = 1
  memory_size    = 2
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.new_type.instance_types.0.instance_type
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}
`

const testAccTencentCloudInstanceWithDataDisk = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type

  system_disk_type = "CLOUD_PREMIUM"
  data_disks {
    data_disk_type       = "CLOUD_PREMIUM"
    data_disk_size       = 100
    delete_with_instance = true
  }

  data_disks {
    data_disk_type       = "CLOUD_PREMIUM"
    data_disk_size       = 100
    delete_with_instance = true
  }

  disable_security_service = true
  disable_monitor_service  = true
}
`

func testAccTencentCloudInstanceWithNetwork(hasPublicIp string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out = 1
  allocate_public_ip         = %s
  system_disk_type           = "CLOUD_PREMIUM"
}
`,
		hasPublicIp,
	)
}

const testAccTencentCloudInstanceWithPrivateIP = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  private_ip        = "172.16.0.130"
}
`

func testAccTencentCloudInstanceWithKeyPair(keyName string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_key_pair" "foo" {
  key_name = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  key_name          = tencentcloud_key_pair.foo.id
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
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  internet_max_bandwidth_out = 1
  password                   = "%s"
  system_disk_type           = "CLOUD_PREMIUM"
}
`,
		password,
	)
}

func testAccTencentCloudInstanceWithName(instanceName string) string {
	return fmt.Sprintf(
		defaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name     = "%s"
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
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
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
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
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  internet_max_bandwidth_out = 1
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
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
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
  availability_zone  = data.tencentcloud_availability_zones.default.zones.0.name
  image_id           = data.tencentcloud_images.default.images.0.image_id
  instance_type      = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type   = "CLOUD_PREMIUM"
  placement_group_id = tencentcloud_placement_group.foo.id
}
`

const testAccTencentCloudInstanceWithSpotpaid = defaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name        = var.instance_name
  availability_zone    = data.tencentcloud_availability_zones.default.zones.0.name
  image_id             = data.tencentcloud_images.default.images.0.image_id
  instance_type        = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  hostname             = var.instance_name
  system_disk_type     = "CLOUD_PREMIUM"
  instance_charge_type = "SPOTPAID"
  spot_instance_type   = "ONE-TIME"
  spot_max_price       = "0.5"
}
`
