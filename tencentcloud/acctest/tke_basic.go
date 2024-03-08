package acctest

const (
	TkeDefaultZone        = "ap-guangzhou-3"
	TkeDefaultVpcCidr     = "172.16.0.0/16"
	TkeDefaultSubnetCidr1 = "172.16.0.0/20"
	TkeDefaultSubnetCidr2 = "172.16.16.0/20"
	TkeDefaultImgId       = "img-2lr9q49h"
)

// todo tke cluster
const TkeResourceKubernetesClusterTestCase = TkeResourceVpcTestCase + TkeResourceSecurityGroupTestCase + TkeDatasourceInstanceTypeTestCase + TkeDefaultVariable + `
resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                                     = local.vpc_id
  cluster_cidr                               = var.tke_cidr_a.0
  cluster_max_pod_num                        = 32
  cluster_name                               = "` + DefaultTkeClusterName + `"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  cluster_internet_domain                    = "tf.cluster-internet.com"
  cluster_intranet                           = true
  cluster_intranet_domain                    = "tf.cluster-intranet.com"
  cluster_version                            = "1.22.5"
  cluster_os                                 = "tlinux2.2(tkernel3)x86_64"
  cluster_level								 = "L5"
  auto_upgrade_cluster_level				 = true
  cluster_intranet_subnet_id                 = local.subnet_id1
  cluster_internet_security_group               = local.sg_id
  managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = local.final_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.subnet_id1
    img_id                     = var.default_img_id
    security_group_ids         = [local.sg_id]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
      file_system = "ext3"
	  auto_format_and_mount = "true"
	  mount_target = "/var/lib/docker"
      disk_partition = "/dev/sdb1"
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"

  tags = {
    "test" = "test"
  }

  unschedulable = 0

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
  extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]
}
`

// sg
const TkeResourceSecurityGroupTestCase = DefaultInstanceVariable + `
resource "tencentcloud_security_group" "example" {
  name        = "tf_tke_sg_test"
  description = "sg test"
}

locals {
  sg_id  = tencentcloud_security_group.example.id
  sg_id2 = tencentcloud_security_group.example.id
}
`

//InstanceType
const TkeDatasourceInstanceTypeTestCase = DefaultInstanceVariable + `
data "tencentcloud_instance_types" "ins_type" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }

  cpu_core_count = 2
  memory_size    = 2
}

locals {
  type1 = [
    for i in data.tencentcloud_instance_types.ins_type.instance_types : i
    if lookup(i, "instance_charge_type") == "POSTPAID_BY_HOUR"
  ]
  type2      = [for i in data.tencentcloud_instance_types.ins_type.instance_types : i]
  final_type = concat(local.type1, local.type2)[0].instance_type
}
`

// vpc
const TkeResourceVpcTestCase = DefaultInstanceVariable + `
resource "tencentcloud_vpc" "tke_vpc" {
  name       = "tf_tke_vpc_test"
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_subnet" "tke_subnet1" {
  name              = "tf_tke_subnet_test1"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr1
  is_multicast      = false
}

resource "tencentcloud_subnet" "tke_subnet2" {
  name              = "tf_tke_subnet_test2"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr2
  is_multicast      = false
}

locals {
  vpc_id     = tencentcloud_vpc.tke_vpc.id
  subnet_id1 = tencentcloud_subnet.tke_subnet1.id
  subnet_id2 = tencentcloud_subnet.tke_subnet2.id
}
`

const TkeDefaultVariable = `
//variable "availability_zone" {
//  default = "` + TkeDefaultZone + `"
//}

variable "vpc_cidr" {
  default = "` + TkeDefaultVpcCidr + `"
}

variable "subnet_cidr1" {
  default = "` + TkeDefaultSubnetCidr1 + `"
}

variable "subnet_cidr2" {
  default = "` + TkeDefaultSubnetCidr2 + `"
}

variable "tke_cidr_a" {
  default = [
    "10.31.0.0/23",
    "10.31.2.0/24",
    "10.31.3.0/24",
    "10.31.16.0/24",
    "10.31.32.0/24"
  ]
}

variable "tke_cidr_b" {
  default = [
    "172.18.0.0/20",
    "172.18.16.0/21",
    "172.18.24.0/21",
    "172.18.32.0/20",
    "172.18.48.0/20"
  ]
}

variable "tke_cidr_c" {
  default = [
    "192.168.0.0/18",
    "192.168.64.0/19",
    "192.168.96.0/20",
    "192.168.112.0/21",
    "192.168.120.0/21"
  ]
}

variable "default_img_id" {
  default = "` + TkeDefaultImgId + `"
}

variable "default_project" {
  default = ` + DefaultProjectId + `
}

variable "default_img_id" {
  default = "` + DefaultTkeOSImageId + `"
}

variable "default_img" {
  default = "` + DefaultTkeOSImageName + `"
}
`
