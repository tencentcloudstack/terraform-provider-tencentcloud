package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_instance.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithSmallInstanceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.foo"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "system_disk_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "system_disk_type", "CLOUD_SSD"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.0.data_disk_type", "CLOUD_SSD"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.0.data_disk_size", "100"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_prepaid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "tencentcloud_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithSmallInstanceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.foo"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "system_disk_type", "CLOUD_SSD"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_sg(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_instance.sg",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithSecurityGroup(`["${tencentcloud_security_group.my_sg1.id}"]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.sg"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.sg"),
					resource.TestCheckResourceAttr("tencentcloud_instance.sg", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_security_group.my_sg1", "id"),
					resource.TestCheckResourceAttr("tencentcloud_instance.sg", "security_groups.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.sg_rule_1", "type", "ingress"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.sg_rule_1", "port_range", "80,8080"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.sg_rule_2", "type", "ingress"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_rule.sg_rule_2", "port_range", "3000"),
				),
			},
			{
				Config: testAccInstanceConfigWithSecurityGroup(`[
					"${tencentcloud_security_group.my_sg1.id}",
					"${tencentcloud_security_group.my_sg2.id}"
				]`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.sg"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.sg"),
					resource.TestCheckResourceAttr("tencentcloud_instance.sg", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_security_group.my_sg2", "id"),
					resource.TestCheckResourceAttr("tencentcloud_instance.sg", "security_groups.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_network(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_instance.network",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithInternet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.network"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.network"),
					resource.TestCheckResourceAttr("tencentcloud_instance.network", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.network", "public_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.network", "private_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_network_no_public_ip(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_instance.network",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithInternetNoPublicIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.network"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.network"),
					resource.TestCheckResourceAttr("tencentcloud_instance.network", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr("tencentcloud_instance.network", "public_ip", ""),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.network", "private_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_vpc(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_instance.vpc_ins",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithVPC,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.vpc_ins"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.vpc_ins"),
					resource.TestCheckResourceAttr("tencentcloud_instance.vpc_ins", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.vpc_ins", "private_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.vpc_ins", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.vpc_ins", "subnet_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_keypair(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_instance.login",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithKeyPair("tf_acc_test_key1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.login"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.login"),
					resource.TestCheckResourceAttr("tencentcloud_instance.login", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.login", "key_name"),
				),
			},
			{
				Config: testAccInstanceConfigWithKeyPair("tf_acc_test_key2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.login"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.login"),
					resource.TestCheckResourceAttr("tencentcloud_instance.login", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.login", "key_name"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_imageIdChanged(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_instance.hello",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithImageIdChanged(
					"img-8toqc6s3",
					"testpwd123",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.hello"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.hello"),
					resource.TestCheckResourceAttr("tencentcloud_instance.hello", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr("tencentcloud_instance.hello", "password", "testpwd123"),
				),
			},
			{
				Config: testAccInstanceConfigWithImageIdChanged(
					"img-oikl1tzv",
					"testpwd1234",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.hello"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.hello"),
					resource.TestCheckResourceAttr("tencentcloud_instance.hello", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr("tencentcloud_instance.hello", "image_id", "img-oikl1tzv"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_passwordChanged(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_instance.hello",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithInstanceNameChanged("tf_testing_1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.hello"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.hello"),
					resource.TestCheckResourceAttr("tencentcloud_instance.hello", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr("tencentcloud_instance.hello", "instance_name", "tf_testing_1"),
				),
			},
			{
				Config: testAccInstanceConfigWithInstanceNameChanged("tf_testing_2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.hello"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.hello"),
					resource.TestCheckResourceAttr("tencentcloud_instance.hello", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr("tencentcloud_instance.hello", "instance_name", "tf_testing_2"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_password(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: "tencentcloud_instance.login",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				//Config: testAccInstanceConfigWithPassword,
				Config: testAccInstanceConfigWithPasswordChanged("TF_test_123"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.login"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.login"),
					resource.TestCheckResourceAttr("tencentcloud_instance.login", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.login", "password"),
				),
			},
			{
				Config: testAccInstanceConfigWithPasswordChanged("TF_test_123456"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_instance.login"),
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.login"),
					resource.TestCheckResourceAttr("tencentcloud_instance.login", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.login", "password"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_hostname(t *testing.T) {
	id := "tencentcloud_instance.hostname"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: id,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithHostname,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "private_ip"),
					resource.TestCheckResourceAttrSet(id, "hostname"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstance_projectId(t *testing.T) {
	id := "tencentcloud_instance.project_id"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		IDRefreshName: id,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigWithProjectId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "private_ip"),
					resource.TestCheckResourceAttrSet(id, "project_id"),
				),
			},
		},
	})
}

func testAccCheckTencentCloudInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		provider := testAccProvider
		// Ignore if Meta is empty, this can happen for validation providers
		if provider.Meta() == nil {
			return fmt.Errorf("Provider Meta is nil")
		}

		client := provider.Meta().(*TencentCloudClient).commonConn
		instanceIds := []string{
			rs.Primary.ID,
		}
		_, err := waitInstanceReachTargetStatus(client, instanceIds, "RUNNING")
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	return testAccCheckInstanceDestroyWithProvider(s, testAccProvider)
}

func testAccCheckInstanceDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*TencentCloudClient).commonConn

	var instanceIds []string
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_instance" {
			continue
		}
		instanceId := rs.Primary.ID
		instanceIds = append(instanceIds, instanceId)
	}

	// NOTE,
	// for prepaid instances, STOPPED means terminated process is done
	// for postpaid instances, not found is expected
	_, err := waitInstanceReachTargetStatus(client, instanceIds, "STOPPED")
	if err != nil {
		if err.Error() == instanceNotFoundErrorMsg(instanceIds) {
			return nil
		}
		return err
	}
	return nil
}

func testAccInstanceConfigWithKeyPair(keyname string) string {
	return fmt.Sprintf(
		`
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_key_pair" "my_key" {
  key_name = "%s"
}

resource "tencentcloud_instance" "login" {
  instance_name = "terraform_automation_test_kuruk"
  availability_zone = "ap-guangzhou-3"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  key_name = "${tencentcloud_key_pair.my_key.id}"
  system_disk_type = "CLOUD_SSD"
}
`,
		keyname,
	)
}

func testAccInstanceConfigWithInstanceNameChanged(name string) string {
	return fmt.Sprintf(
		`
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_instance" "hello" {
  instance_name = "%s"
  availability_zone = "ap-guangzhou-3"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  system_disk_type = "CLOUD_SSD"
}
`,
		name,
	)
}

func testAccInstanceConfigWithImageIdChanged(imageId, password string) string {
	return fmt.Sprintf(
		`
data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_instance" "hello" {
  instance_name = "tf_test_reset_instance"
  availability_zone = "ap-guangzhou-3"
  image_id      = "%s"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  password      = "%s"
  system_disk_type = "CLOUD_SSD"
}
`,
		imageId,
		password,
	)
}

const testAccInstanceConfigWithSmallInstanceType = `
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 2
}

resource "tencentcloud_instance" "foo" {
  instance_name = "terraform_automation_test_kuruk"
  availability_zone = "ap-guangzhou-3"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"

  system_disk_type = "CLOUD_SSD"
  data_disks = [
    {
      data_disk_type = "CLOUD_SSD"
      data_disk_size = 100
    },
    {
      data_disk_type = "CLOUD_SSD"
      data_disk_size = 100
    }
  ]
  disable_security_service = true
  disable_monitor_service = true
}
`

const testAccInstanceConfigChargeTypePrepaid = `
data "tencentcloud_image" "myimage" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

resource "tencentcloud_instance" "foo" {
  instance_name = "terraform_ci_test"
  availability_zone = "ap-guangzhou-3"
  image_id      = "${data.tencentcloud_image.myimage.image_id}"
  instance_type = "S2.SMALL1"
  system_disk_type = "CLOUD_SSD"
  instance_charge_type = "PREPAID"
  instance_charge_type_prepaid_period = 1
}
`

const testAccInstanceConfigWithInternet = `
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_instance" "network" {
  instance_name = "terraform_automation_test_kuruk"
  availability_zone = "ap-guangzhou-3"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out = 1
  allocate_public_ip = true
  system_disk_type = "CLOUD_SSD"
}
`

const testAccInstanceConfigWithInternetNoPublicIP = `
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_instance" "network" {
  instance_name = "terraform_automation_test_kuruk"
  availability_zone = "ap-guangzhou-3"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
  internet_max_bandwidth_out = 1
  allocate_public_ip = false
  system_disk_type = "CLOUD_SSD"
}
`

func testAccInstanceConfigWithPasswordChanged(pwd string) string {
	return fmt.Sprintf(
		`
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_instance" "login" {
  instance_name = "terraform_automation_test_kuruk"
  availability_zone = "ap-guangzhou-3"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  internet_max_bandwidth_out = 1
  password = "%s"
  system_disk_type = "CLOUD_SSD"
}
`,
		pwd,
	)
}

const testAccInstanceConfigWithVPC = `
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 2
}

resource "tencentcloud_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_vpc_test"
}

resource "tencentcloud_subnet" "my_subnet" {
  vpc_id = "${tencentcloud_vpc.my_vpc.id}"
  availability_zone = "ap-guangzhou-3"
  name              = "tf_test_subnet"
  cidr_block        = "10.0.2.0/24"
}

resource "tencentcloud_instance" "vpc_ins" {
  instance_name = "terraform_automation_test_kuruk_vpc"
  availability_zone = "ap-guangzhou-3"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  vpc_id = "${tencentcloud_vpc.my_vpc.id}"
  subnet_id = "${tencentcloud_subnet.my_subnet.id}"
  system_disk_type = "CLOUD_SSD"
}
`

const testAccInstanceConfigWithHostname = `
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 2
}

data "tencentcloud_availability_zones" "my_favorate_zones" {
	name = "ap-guangzhou-3"
}

resource "tencentcloud_instance" "hostname" {
  instance_name = "terraform_automation_test_kuruk"
  availability_zone     = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  hostname      = "testing"

  system_disk_type = "CLOUD_SSD"
}
`

const testAccInstanceConfigWithProjectId = `
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 2
}

data "tencentcloud_availability_zones" "my_favorate_zones" {
	name = "ap-guangzhou-3"
}

resource "tencentcloud_instance" "project_id" {
  instance_name = "terraform_automation_test_kuruk"
  availability_zone     = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  project_id    = 0

  system_disk_type = "CLOUD_SSD"
}
`

func testAccInstanceConfigWithSecurityGroup(rule string) string {
	return fmt.Sprintf(`
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 2
}

resource "tencentcloud_security_group" "my_sg1" {
  name = "tf_test_sg_name"
  description = "tf_test_sg_desc"
}

resource "tencentcloud_security_group_rule" "sg_rule_1" {
  security_group_id = "${tencentcloud_security_group.my_sg1.id}"
  type = "ingress"
  cidr_ip = "0.0.0.0/0"
  ip_protocol = "tcp"
  port_range = "80,8080"
  policy = "accept"
}

resource "tencentcloud_security_group" "my_sg2" {
  name = "tf_test_sg_name"
  description = "tf_test_sg_desc"
}

resource "tencentcloud_security_group_rule" "sg_rule_2" {
  security_group_id = "${tencentcloud_security_group.my_sg2.id}"
  type = "ingress"
  cidr_ip = "0.0.0.0/0"
  ip_protocol = "tcp"
  port_range = "3000"
  policy = "accept"
}

resource "tencentcloud_instance" "sg" {
  instance_name = "terraform_automation_test_kuruk_sg"
  availability_zone = "ap-guangzhou-3"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  system_disk_type = "CLOUD_SSD"

  internet_max_bandwidth_out = 1

  security_groups = %s
}
`,
		rule,
	)
}
