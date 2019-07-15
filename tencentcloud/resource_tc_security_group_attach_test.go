package tencentcloud

/*import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudSecurityGroupAttach_basic(t *testing.T) {
	var sgId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupAttachDestroy(&sgId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupAttachConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_security_group_attach.foo"),
					testAccCheckSecurityGroupAttachExists("tencentcloud_security_group_attach.foo", &sgId),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.foo", "cvm_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.foo", "mysql_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.foo", "clb_ids.#", "1"),

				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupAttach_cvm(t *testing.T) {
	var sgId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupAttachDestroy(&sgId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupAttachConfigCvm(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupAttachExists("tencentcloud_security_group_attach.cvm", &sgId),
					testAccCheckTencentCloudDataSourceID("tencentcloud_security_group_attach.cvm"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.cvm", "cvm_ids.#", "1"),
				),
			},

			{
				Config: testAccSecurityGroupAttachConfigCvmUpdateAdd(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupAttachExists("tencentcloud_security_group_attach.cvm", &sgId),
					testAccCheckTencentCloudDataSourceID("tencentcloud_security_group_attach.cvm"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.cvm", "cvm_ids.#", "2"),
				),
			},

			{
				Config: testAccSecurityGroupAttachConfigCvmUpdateDelete(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupAttachExists("tencentcloud_security_group_attach.cvm", &sgId),
					testAccCheckTencentCloudDataSourceID("tencentcloud_security_group_attach.cvm"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.cvm", "cvm_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupAttach_mysql(t *testing.T) {
	var sgId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupAttachDestroy(&sgId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupAttachConfigMysql(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupAttachExists("tencentcloud_security_group_attach.mysql_sherlokyang", &sgId),
					testAccCheckTencentCloudDataSourceID("tencentcloud_security_group_attach.mysql_sherlokyang"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.mysql_sherlokyang", "mysql_ids.#", "1"),
				),
			},

			{
				Config: testAccSecurityGroupAttachConfigMysqlUpdate1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupAttachExists("tencentcloud_security_group_attach.mysql_sherlokyang", &sgId),
					testAccCheckTencentCloudDataSourceID("tencentcloud_security_group_attach.mysql_sherlokyang"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.mysql_sherlokyang", "mysql_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroupAttach_clb(t *testing.T) {
	var sgId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupAttachDestroy(&sgId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupAttachConfigClb(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupAttachExists("tencentcloud_security_group.sg_attach_test_clb", &sgId),
					testAccCheckTencentCloudDataSourceID("tencentcloud_security_group_attach.clb"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.clb", "clb_ids.#", "1"),
				),
			},

			{
				Config: testAccSecurityGroupAttachConfigClbUpdateAdd(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupAttachExists("tencentcloud_security_group.sg_attach_test_clb", &sgId),
					testAccCheckTencentCloudDataSourceID("tencentcloud_security_group_attach.clb"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.clb", "clb_ids.#", "2"),
				),
			},

			{
				Config: testAccSecurityGroupAttachConfigClbUpdateDelete(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupAttachExists("tencentcloud_security_group.sg_attach_test_clb", &sgId),
					testAccCheckTencentCloudDataSourceID("tencentcloud_security_group_attach.clb"),
					resource.TestCheckResourceAttr("tencentcloud_security_group_attach.clb", "clb_ids.#", "0"),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupAttachExists(n string, sgId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		sgAttachId := rs.Primary.ID

		if sgAttachId == "" {
			return fmt.Errorf("no security group attach ID is set")
		}

		*sgId = sgAttachId

		client := testAccProvider.Meta().(*TencentCloudClient)

		service := VpcService{client: client.apiV3Conn}

		_, has, err := service.DescribeSecurityGroup(context.TODO(), sgAttachId)
		if err != nil {
			return err
		}

		if has == 0 {
			return fmt.Errorf("security group not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckSecurityGroupAttachDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn

		service := VpcService{client: client}

		sg, has, err := service.DescribeSecurityGroup(context.TODO(), *id)
		if err != nil {
			return err
		}

		// if security group does not exist, all instances should be destroyed
		if has == 0 {
			return nil
		}

		cvmService := CvmService{client: client}
		cvmIns, err := cvmService.DescribeBySecurityGroups(context.TODO(), *sg.SecurityGroupId)
		if err != nil {
			return err
		}
		if len(cvmIns) != 0 {
			return fmt.Errorf("cvm still attach security group %s", *id)
		}

		mysqlService := MysqlService{client: client}
		mysqlIns, err := mysqlService.DescribeDBInstancesBySecurityGroup(context.TODO(), *sg.SecurityGroupId)
		if err != nil {
			return err
		}
		if len(mysqlIns) != 0 {
			return fmt.Errorf("mysql still attach security group %s", *id)
		}

		clbService := ClbService{client: client}
		clbIns, err := clbService.DescribeLoadBalances(context.TODO(), nil, sg.SecurityGroupId)
		if err != nil {
			return err
		}
		if len(clbIns) != 0 {
			return fmt.Errorf("cloud load balance still attach security group %s", *id)
		}

		return nil
	}
}

func testAccSecurityGroupAttachConfigCvm() string {
	return `
data "tencentcloud_instance_types" "sg_attach_test_cvm" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_security_group" "sg_attach_test_cvm" {
	name        = "sg_attach_test_cvm"
	description = "sg_attach_test_cvm"
}

resource "tencentcloud_vpc" "sg_attach_test_cvm" {
	cidr_block = "10.0.0.0/16"
	name       = "sg_attach_test_cvm"
}

resource "tencentcloud_subnet" "sg_attach_test_cvm" {
	vpc_id = "${tencentcloud_vpc.sg_attach_test_cvm.id}"
	name              = "sg-attach-test-cvm"
    cidr_block        = "10.0.1.0/24"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_instance" "sg_attach_test_cvm" {
	instance_name     = "sg-attach-test-cvm"
	availability_zone = "ap-guangzhou-3"
	image_id          = "img-6rrx0ymd"
	instance_type     = "${data.tencentcloud_instance_types.sg_attach_test_cvm.instance_types.0.instance_type}"
	password          = "123abcDEF,._+-&="

    system_disk_type = "CLOUD_PREMIUM"
    system_disk_size = 50

	vpc_id    = "${tencentcloud_vpc.sg_attach_test_cvm.id}"
	subnet_id = "${tencentcloud_subnet.sg_attach_test_cvm.id}"
}


resource "tencentcloud_security_group_attach" "cvm" {
	security_group_id = "${tencentcloud_security_group.sg_attach_test_cvm.id}"
	cvm_ids = ["${tencentcloud_instance.sg_attach_test_cvm.id}"]
	seni_ids = []
	mysql_ids = []
	clb_ids = []
}`
}

func testAccSecurityGroupAttachConfigCvmUpdateAdd() string {
	return `
data "tencentcloud_instance_types" "sg_attach_test_cvm" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_security_group" "sg_attach_test_cvm" {
	name        = "sg_attach_test_cvm"
	description = "sg_attach_test_cvm"
}

resource "tencentcloud_vpc" "sg_attach_test_cvm" {
	cidr_block = "10.0.0.0/16"
	name       = "sg_attach_test_cvm"
}

resource "tencentcloud_subnet" "sg_attach_test_cvm" {
	vpc_id = "${tencentcloud_vpc.sg_attach_test_cvm.id}"
	name              = "sg-attach-test-cvm"
    cidr_block        = "10.0.1.0/24"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_instance" "sg_attach_test_cvm" {
	instance_name     = "sg-attach-test-cvm"
	availability_zone = "ap-guangzhou-3"
	image_id          = "img-6rrx0ymd"
	instance_type     = "${data.tencentcloud_instance_types.sg_attach_test_cvm.instance_types.0.instance_type}"
	password          = "123abcDEF,._+-&="

    system_disk_type = "CLOUD_PREMIUM"
    system_disk_size = 50

	vpc_id    = "${tencentcloud_vpc.sg_attach_test_cvm.id}"
	subnet_id = "${tencentcloud_subnet.sg_attach_test_cvm.id}"
}

resource "tencentcloud_instance" "sg_attach_test_cvm_2" {
	instance_name     = "sg-attach-test-cvm"
	availability_zone = "ap-guangzhou-3"
	image_id          = "img-6rrx0ymd"
	instance_type     = "${data.tencentcloud_instance_types.sg_attach_test_cvm.instance_types.0.instance_type}"
	password          = "123abcDEF,._+-&="

    system_disk_type = "CLOUD_PREMIUM"
    system_disk_size = 50

	vpc_id    = "${tencentcloud_vpc.sg_attach_test_cvm.id}"
	subnet_id = "${tencentcloud_subnet.sg_attach_test_cvm.id}"
}


resource "tencentcloud_security_group_attach" "cvm" {
	security_group_id = "${tencentcloud_security_group.sg_attach_test_cvm.id}"
	cvm_ids = [
	  "${tencentcloud_instance.sg_attach_test_cvm.id}",
      "${tencentcloud_instance.sg_attach_test_cvm_2.id}"
	]
	seni_ids = []
	mysql_ids = []
	clb_ids = []
}`
}

func testAccSecurityGroupAttachConfigCvmUpdateDelete() string {
	return `
data "tencentcloud_instance_types" "sg_attach_test_cvm" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_security_group" "sg_attach_test_cvm" {
	name        = "sg_attach_test_cvm"
	description = "sg_attach_test_cvm"
}

resource "tencentcloud_vpc" "sg_attach_test_cvm" {
	cidr_block = "10.0.0.0/16"
	name       = "sg_attach_test_cvm"
}

resource "tencentcloud_subnet" "sg_attach_test_cvm" {
	vpc_id = "${tencentcloud_vpc.sg_attach_test_cvm.id}"
	name              = "sg-attach-test-cvm"
    cidr_block        = "10.0.1.0/24"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_instance" "sg_attach_test_cvm" {
	instance_name     = "sg-attach-test-cvm"
	availability_zone = "ap-guangzhou-3"
	image_id          = "img-6rrx0ymd"
	instance_type     = "${data.tencentcloud_instance_types.sg_attach_test_cvm.instance_types.0.instance_type}"
	password          = "123abcDEF,._+-&="

    system_disk_type = "CLOUD_PREMIUM"
    system_disk_size = 50

	vpc_id    = "${tencentcloud_vpc.sg_attach_test_cvm.id}"
	subnet_id = "${tencentcloud_subnet.sg_attach_test_cvm.id}"
}

resource "tencentcloud_instance" "sg_attach_test_cvm_2" {
	instance_name     = "sg-attach-test-cvm"
	availability_zone = "ap-guangzhou-3"
	image_id          = "img-6rrx0ymd"
	instance_type     = "${data.tencentcloud_instance_types.sg_attach_test_cvm.instance_types.0.instance_type}"
	password          = "123abcDEF,._+-&="

    system_disk_type = "CLOUD_PREMIUM"
    system_disk_size = 50

	vpc_id    = "${tencentcloud_vpc.sg_attach_test_cvm.id}"
	subnet_id = "${tencentcloud_subnet.sg_attach_test_cvm.id}"
}


resource "tencentcloud_security_group_attach" "cvm" {
	security_group_id = "${tencentcloud_security_group.sg_attach_test_cvm.id}"
	cvm_ids = []
	seni_ids = []
	mysql_ids = []
	clb_ids = []
}`
}

func testAccSecurityGroupAttachConfigMysql() string {
	return `
resource "tencentcloud_security_group" "sg_attach_test_mysql_sherlokyang" {
  name        = "sg_attach_test_mysql_sherlokyang"
  description = "sg_attach_test_mysql_sherlokyang"
}

resource "tencentcloud_mysql_instance" "sg_attach_test_mysql_sherlokyang" {
  internet_service = 0
  engine_version = "5.7"
  root_password = "243mh1F0PBHdFDjbllZhhw=="
  availability_zone = "ap-guangzhou-3"
  instance_name = "sg_attach_test_mysql_sherlokyang"
  mem_size = 1000
  volume_size = 50
  intranet_port = 3306
}


resource "tencentcloud_security_group_attach" "mysql_sherlokyang" {
  security_group_id = "${tencentcloud_security_group.sg_attach_test_mysql_sherlokyang.id}"
  cvm_ids = []
  seni_ids = []
  mysql_ids = ["${tencentcloud_mysql_instance.sg_attach_test_mysql_sherlokyang.id}"]
  clb_ids = []
}`
}

func testAccSecurityGroupAttachConfigMysqlUpdate1() string {
	return `
resource "tencentcloud_security_group" "sg_attach_test_mysql_sherlokyang" {
  name        = "sg_attach_test_mysql_sherlokyang"
  description = "sg_attach_test_mysql_sherlokyang"
}

resource "tencentcloud_mysql_instance" "sg_attach_test_mysql_sherlokyang" {
  internet_service = 0
  engine_version = "5.7"
  root_password = "243mh1F0PBHdFDjbllZhhw=="
  availability_zone = "ap-guangzhou-3"
  instance_name = "sg_attach_test_mysql_sherlokyang"
  mem_size = 1000
  volume_size = 50
  intranet_port = 3306
}


resource "tencentcloud_security_group_attach" "mysql_sherlokyang" {
  security_group_id = "${tencentcloud_security_group.sg_attach_test_mysql_sherlokyang.id}"
  cvm_ids = []
  seni_ids = []
  mysql_ids = []
  clb_ids = []
}`
}

func testAccSecurityGroupAttachConfigClb() string {
	return `
resource "tencentcloud_security_group" "sg_attach_test_clb" {
  name        = "sg_attach_test_clb"
  description = "sg_attach_test_clb"
}


resource "tencentcloud_lb" "sg_attach_test_clb" {
  type       = "OPEN"
  forward    = "APPLICATION"
  name       = "sg_attach_test_clb"
}


resource "tencentcloud_security_group_attach" "clb" {
  security_group_id = "${tencentcloud_security_group.sg_attach_test_clb.id}"
  cvm_ids = []
  seni_ids = []
  mysql_ids = []
  clb_ids = ["${tencentcloud_lb.sg_attach_test_clb.id}"]
}
`
}

func testAccSecurityGroupAttachConfigClbUpdateAdd() string {
	return `
resource "tencentcloud_security_group" "sg_attach_test_clb" {
  name        = "sg_attach_test_clb"
  description = "sg_attach_test_clb"
}


resource "tencentcloud_lb" "sg_attach_test_clb" {
  type       = "OPEN"
  forward    = "APPLICATION"
  name       = "sg_attach_test_clb"
}

resource "tencentcloud_lb" "sg_attach_test_clb_2" {
  type       = "OPEN"
  forward    = "APPLICATION"
  name       = "sg_attach_test_clb_2"
}


resource "tencentcloud_security_group_attach" "clb" {
  security_group_id = "${tencentcloud_security_group.sg_attach_test_clb.id}"
  cvm_ids = []
  seni_ids = []
  mysql_ids = []
  clb_ids = [
    "${tencentcloud_lb.sg_attach_test_clb.id}",
    "${tencentcloud_lb.sg_attach_test_clb_2.id}"
  ]
}
`
}

func testAccSecurityGroupAttachConfigClbUpdateDelete() string {
	return `
resource "tencentcloud_security_group" "sg_attach_test_clb" {
  name        = "sg_attach_test_clb"
  description = "sg_attach_test_clb"
}


resource "tencentcloud_lb" "sg_attach_test_clb" {
  type       = "OPEN"
  forward    = "APPLICATION"
  name       = "sg_attach_test_clb"
}

resource "tencentcloud_lb" "sg_attach_test_clb_2" {
  type       = "OPEN"
  forward    = "APPLICATION"
  name       = "sg_attach_test_clb_2"
}


resource "tencentcloud_security_group_attach" "clb" {
  security_group_id = "${tencentcloud_security_group.sg_attach_test_clb.id}"
  cvm_ids = []
  seni_ids = []
  mysql_ids = []
  clb_ids = []
}
`
}

func testAccSecurityGroupAttachConfig() string {
	return `data "tencentcloud_instance_types" "sg_attach_test_types" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_security_group" "sg_attach_test" {
  name        = "sg_attach_test"
  description = "sg_attach_test"
}

resource "tencentcloud_vpc" "sg_attach_test" {
  cidr_block = "10.0.0.0/16"
  name       = "sg-attach-test"
}

resource "tencentcloud_subnet" "sg_attach_test_subnet" {
  vpc_id = "${tencentcloud_vpc.sg_attach_test.id}"
  name              = "sg-attach-test"
  cidr_block        = "10.0.1.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_instance" "sg_attach_test-cvm" {
  instance_name     = "sg-attach-test"
  availability_zone = "ap-guangzhou-3"
  image_id          = "img-6rrx0ymd"
  instance_type     = "${data.tencentcloud_instance_types.sg_attach_test_types.instance_types.0.instance_type}"
  password          = "243mh1F0PBHdFDjbllZhhw=="

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 50

  vpc_id    = "${tencentcloud_vpc.sg_attach_test.id}"
  subnet_id = "${tencentcloud_subnet.sg_attach_test_subnet.id}"
}


resource "tencentcloud_mysql_instance" "sg_attach_test_cdb" {
  internet_service = 0
  engine_version = "5.7"
  root_password = "243mh1F0PBHdFDjbllZhhw=="
  availability_zone = "ap-guangzhou-3"
  instance_name = "myTestMysql"
  mem_size = 1000
  volume_size = 50
  vpc_id = "${tencentcloud_vpc.sg_attach_test.id}"
  subnet_id = "${tencentcloud_subnet.sg_attach_test_subnet.id}"
  intranet_port = 3306
}


resource "tencentcloud_lb" "sg_attach_test_clb" {
  type       = "INTERNAL"
  forward    = "APPLICATION"
  name       = "sg_attach_test_clb"
  vpc_id     = "${tencentcloud_vpc.sg_attach_test.id}"
}


resource "tencentcloud_security_group_attach" "foo" {
  security_group_id = "${tencentcloud_security_group.sg_attach_test.id}"
  cvm_ids = ["${tencentcloud_instance.sg_attach_test-cvm.id}"]
  seni_ids = []
  mysql_ids = ["${tencentcloud_mysql_instance.sg_attach_test_cdb.id}"]
  clb_ids = ["${tencentcloud_lb.sg_attach_test_clb.id}"]
}`
}*/
