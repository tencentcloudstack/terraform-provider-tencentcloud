package tencentcloud

import (
	"os"
	"regexp"
	"time"
)

/*
---------------------------------------------------
If you want to run through the test cases,
the following must be changed to your resource id.
---------------------------------------------------
*/

var appid string = os.Getenv("TENCENTCLOUD_APPID")
var ownerUin string = os.Getenv("TENCENTCLOUD_OWNER_UIN")

const (
	keepResource    = "keep"
	defaultResource = "Default"
)

var persistResource = regexp.MustCompile("^(keep|Default)")

// Check if resource should persist instead of recycled
func isResourcePersist(name string, createdTime *time.Time) bool {
	createdWithin30Minutes := false
	if createdTime != nil {
		createdWithin30Minutes = createdTime.Add(time.Minute * 30).After(time.Now())
	}
	return persistResource.MatchString(name) || createdWithin30Minutes
}

// vpn
const defaultVpnDataSource = `
data "tencentcloud_vpn_gateways" "foo" {
  name = "keep-vpn-gw"
}

data "tencentcloud_vpn_connections" "conns" {
  name = "keep-vpn-conn"
}
`

// clb
const (
	defaultSshCertificate  = "vYSQkJ3K"
	defaultSshCertificateB = "vYVlNIhW"
)

const (
	defaultRegion      = "ap-guangzhou"
	defaultVpcId       = "vpc-86v957zb"
	defaultVpcCidr     = "172.16.0.0/16"
	defaultVpcCidrLess = "172.16.0.0/18"

	defaultCvmAZone    = "ap-guangzhou-7"
	defaultCvmVpcId    = "vpc-l0dw94uh"
	defaultCvmSubnetId = "subnet-ccj2qg0m"

	defaultAZone          = "ap-guangzhou-3"
	defaultSubnetId       = "subnet-enm92y0m"
	defaultSubnetCidr     = "172.16.0.0/20"
	defaultSubnetCidrLess = "172.16.0.0/22"

	defaultInsName       = "tf-ci-test"
	defaultInsNameUpdate = "tf-ci-test-update"

	defaultDayuBgp    = "bgp-000006mq"
	defaultDayuBgpMul = "bgp-0000008o"
	defaultDayuBgpIp  = "bgpip-00000294"
	defaultDayuNet    = "net-0000007e"

	defaultGaapProxyId             = "link-ljb08m2l"
	defaultGaapProxyId2            = "link-8lpyo88p"
	defaultGaapSecurityPolicyId    = "sp-05t5q92x"
	defaultGaapRealserverDomainId1 = "rs-qs0h6wxp"
	defaultGaapRealserverDomain1   = "github.com"
	defaultGaapRealserverDomainId2 = "rs-qcygnwpd"
	defaultGaapRealserverDomain2   = "www.github.com"
	defaultGaapRealserverIpId1     = "rs-24e1ol23"
	defaultGaapRealserverIp1       = "119.29.29.35"
	defaultGaapRealserverIpId2     = "rs-70qzt26p"
	defaultGaapRealserverIp2       = "1.1.1.5"

	defaultSecurityGroup  = "sg-ijato2x1"
	defaultSecurityGroup2 = "sg-51rgzop1"

	defaultProjectId   = "1250480"
	defaultDayuBgpIdV2 = "bgpip-000004x0"
	defaultDayuBgpIpV2 = "119.28.217.253"

	defaultTkeOSImageId   = "img-2lr9q49h"
	defaultTkeOSImageName = "tlinux2.2(tkernel3)x86_64"
	defaultEMRSgId        = "sg-mag8k2fj"
)

//ckafka
const (
	defaultKafkaInstanceId = "ckafka-vv7wpvae"
	defaultKafkaVpcId      = "vpc-68vi2d3h"
	defaultKafkaSubnetId   = "subnet-ob6clqwk"
)

const defaultKafkaVariable = `
variable "instance_id" {
  default = "` + defaultKafkaInstanceId + `"
}
variable "vpc_id" {
  default = "` + defaultKafkaVpcId + `"
}
variable "subnet_id" {
  default = "` + defaultKafkaSubnetId + `"
}
`

// Tke Exclusive Network Environment
const (
	tkeExclusiveVpcName   = "keep_tke_exclusive_vpc"
	defaultTkeClusterId   = "cls-ely08ic4"
	defaultTkeClusterName = "keep-tke-cluster"
)

/*
---------------------------------------------------
The following are common test case used as templates.
---------------------------------------------------
*/

const defaultVpcVariable = `
variable "instance_name" {
  default = "` + defaultInsName + `"
}

variable "instance_name_update" {
  default = "` + defaultInsNameUpdate + `"
}

variable "availability_region" {
  default = "` + defaultRegion + `"
}

variable "availability_zone" {
  default = "` + defaultAZone + `"
}

variable "availability_cvm_zone" {
  default = "` + defaultCvmAZone + `"
}

variable "cvm_vpc_id" {
  default = "` + defaultCvmVpcId + `"
}

variable "cvm_subnet_id" {
  default = "` + defaultCvmSubnetId + `"
}

variable "vpc_id" {
  default = "` + defaultVpcId + `"
}

variable "vpc_cidr" {
  default = "` + defaultVpcCidr + `"
}

variable "vpc_cidr_less" {
  default = "` + defaultVpcCidrLess + `"
}

variable "subnet_id" {
  default = "` + defaultSubnetId + `"
}

variable "subnet_cidr" {
  default = "` + defaultSubnetCidr + `"
}

variable "subnet_cidr_less" {
  default = "` + defaultSubnetCidrLess + `"
}
`

const defaultInstanceVariable = defaultVpcVariable + `
data "tencentcloud_availability_zones" "default" {
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S6"]
  }
  filter {
    name   = "zone"
    values = [var.availability_cvm_zone]
  }
  cpu_core_count = 2
  memory_size    = 2
}
`
const defaultAzVariable = `
variable "default_az" {
  default = "ap-guangzhou-3"
}
`

const defaultImages = `
variable "default_img_id" {
  default = "` + defaultTkeOSImageId + `"
}

variable "default_img" {
  default = "` + defaultTkeOSImageName + `"
}
`

// default VPC/Subnet datasource
const defaultVpcSubnets = defaultAzVariable + `

data "tencentcloud_vpc_subnets" "gz3" {
  availability_zone = var.default_az
  is_default = true
}

locals {
  vpc_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
}`

const defaultSecurityGroupData = `
data "tencentcloud_security_groups" "internal" {
  name = "default"
}

locals {
  # local.sg_id
  sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
}
`

const (
	defaultMySQLName = "keep_preset_mysql"
)

// ref with `local.mysql_id`
const CommonPresetMysql = `

variable "availability_zone" {
  default = "` + defaultAZone + `"
}
variable "region" {
  default = "` + defaultRegion + `"
}

data "tencentcloud_mysql_instance" "mysql" {
  instance_name = "` + defaultMySQLName + `"
}

locals {
  mysql_id = data.tencentcloud_mysql_instance.mysql.instance_list.0.mysql_id
}
`

// SQLServer
const defaultSQLServerName = "keep-preset_sqlserver"
const defaultSQLServerDB = "keep_sqlserver_db"
const defaultSQLServerAccount = "keep_sqlserver_account"

const CommonPresetSQLServer = `
data "tencentcloud_sqlserver_instances" "sqlserver" {
  name = "` + defaultSQLServerName + `"
}

locals {
  # local.sqlserver_id
  sqlserver_id = data.tencentcloud_sqlserver_instances.sqlserver.instance_list.0.id
  sqlserver_db = "` + defaultSQLServerDB + `"
}
`

const CommonPresetSQLServerAccount = CommonPresetSQLServer + `
data "tencentcloud_sqlserver_accounts" "test"{
  instance_id = local.sqlserver_id
  name = "` + defaultSQLServerAccount + `"
}

locals {
  # local.sqlserver_id
  sqlserver_account = data.tencentcloud_sqlserver_accounts.test.list.0.name
}
`

const testAccSqlserverAZ = `
data "tencentcloud_availability_zones_by_product" "zone" {
  product = "sqlserver"
}

locals {
  # local.az, local.az1
  az = data.tencentcloud_availability_zones_by_product.zone.zones[0].name
  az1 = data.tencentcloud_availability_zones_by_product.zone.zones[1].name
}
`

const instanceCommonTestCase = defaultInstanceVariable + `
resource "tencentcloud_instance" "default" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.1.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 10
  vpc_id                     = var.cvm_vpc_id
  subnet_id                  = var.cvm_subnet_id
}
`

// End of SQLServer

// PostgreSQL

const defaultPGSQLName = "keep-postgresql"
const CommonPresetPGSQL = `
data "tencentcloud_postgresql_instances" "foo" {
  name = "` + defaultPGSQLName + `"
}

locals {
  pgsql_id = data.tencentcloud_postgresql_instances.foo.instance_list.0.id
}
`

// End of PostgreSQL

const defaultCVMName = "keep-cvm"
const presetCVM = `
data "tencentcloud_instances" "instance" {
  instance_name = "` + defaultCVMName + `"
}

locals {
  cvm_id = data.tencentcloud_instances.instance.instance_list.0.instance_id
  cvm_az = "` + defaultAZone + `"
}
`

const userInfoData = `
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
  uin = data.tencentcloud_user_info.info.uin
  owner_uin = data.tencentcloud_user_info.info.owner_uin
}
`

const mysqlInstanceCommonTestCase = defaultVpcVariable + `
resource "tencentcloud_mysql_instance" "default" {
  mem_size = 1000
  volume_size = 25
  instance_name = var.instance_name
  engine_version = "5.7"
  root_password = "0153Y474"
  availability_zone = var.availability_zone
  force_delete = true
}
`
const mysqlInstanceHighPerformanceTestCase = defaultVpcVariable + `
resource "tencentcloud_mysql_instance" "default" {
  mem_size = 1000
  volume_size = 50
  instance_name = var.instance_name
  engine_version = "5.7"
  root_password = "0153Y474"
  availability_zone = var.availability_zone
  force_delete = true
}
`

const mysqlInstanceHighPerformancePrepaidTestCase = defaultVpcVariable + `
resource "tencentcloud_mysql_instance" "default" {
  mem_size = 1000
  volume_size = 50
  pay_type = 0
  instance_name = var.instance_name
  engine_version = "5.7"
  root_password = "0153Y474"
  availability_zone = var.availability_zone
  force_delete = true
}
`

const defaultSCFCosBucket = `
data "tencentcloud_user_info" "info" {}

data "tencentcloud_cos_buckets" "buckets" {
  bucket_prefix = "preset-scf-bucket-${data.tencentcloud_user_info.info.app_id}"
}

locals {
  bucket_name = data.tencentcloud_cos_buckets.buckets.bucket_list.0.bucket
  bucket_url = data.tencentcloud_cos_buckets.buckets.bucket_list.0.cos_bucket_url
}
`

const defaultScfNamespace = "preset-scf-namespace"

const defaultFileSystemName = "keep_preset_cfs"

const defaultFileSystem = `
data "tencentcloud_cfs_file_systems" "fs" {
  name = "` + defaultFileSystemName + `"
}

# doesn't support datasource for now
variable "mount_id" {
  default = "cfs-iobiaxtj"
}

locals {
  cfs = data.tencentcloud_cfs_file_systems.fs.file_system_list.0
  cfs_id = local.cfs.file_system_id
}`

const defaultCamVariables = `
variable "cam_role_basic" {
  default = "keep-cam-role"
}

variable "cam_policy_basic" {
  default = "keep-cam-policy"
}

variable "cam_group_basic" {
  default = "keep-cam-group"
}

variable "cam_user_basic" {
  default = "keep-cam-user"
}
`

// TCR Service
const defaultTCRInstanceName = "keep-tcr-instance"
const defaultTCRNamespace = "keep-tcr-namespace"
const defaultTCRRepoName = "keep-tcr-repo"

const defaultTCRInstanceVar = `
variable "tcr_name" {
  default = "` + defaultTCRInstanceName + `"
}

variable "tcr_namespace" {
  default = "` + defaultTCRNamespace + `"
}

variable "tcr_repo" {
  default = "` + defaultTCRRepoName + `"
}
`

const defaultTCRInstanceData = defaultTCRInstanceVar + `
data "tencentcloud_tcr_instances" "tcr" {
  name = var.tcr_name
}

locals {
  tcr_id = data.tencentcloud_tcr_instances.tcr.instance_list.0.id
}
`

// End of TCR Service

// TcaPlus DB

const defaultTcaPlusClusterName = "keep-tcaplus-cluster"
const defaultTcaPlusClusterTableGroup = "keep_table_group"
const defaultTcaPlusClusterTable = "keep_players"
const defaultTcaPlusVar = `
variable "tcaplus_cluster" {
  default = "` + defaultTcaPlusClusterName + `"
}

variable "tcaplus_table_group" {
  default = "` + defaultTcaPlusClusterTableGroup + `"
}

variable "tcaplus_table" {
  default = "` + defaultTcaPlusClusterTable + `"
}
`
const defaultTcaPlusData = defaultTcaPlusVar + `
data "tencentcloud_tcaplus_clusters" "tcaplus" {
  cluster_name = var.tcaplus_cluster
}

data "tencentcloud_tcaplus_tablegroups" "group" {
  cluster_id = data.tencentcloud_tcaplus_clusters.tcaplus.list.0.cluster_id
  tablegroup_name = var.tcaplus_table_group
}

locals {
  tcaplus_id = data.tencentcloud_tcaplus_clusters.tcaplus.list.0.cluster_id
  tcaplus_table_group = var.tcaplus_table_group
  tcaplus_table_group_id = data.tencentcloud_tcaplus_tablegroups.group.list.0.tablegroup_id
  tcaplus_table = var.tcaplus_table
}
`

// End of TcaPlus DB

// TKE Service
const TkeInstanceType = `
data "tencentcloud_instance_types" "ins_type" {
  availability_zone = "` + defaultCvmAZone + `"
  cpu_core_count    = 1
  memory_size       = 2
}

locals {
  type1 = [for i in data.tencentcloud_instance_types.ins_type.instance_types: i if lookup(i, "instance_charge_type") == "POSTPAID_BY_HOUR"][0].instance_type
  type2 = [for i in data.tencentcloud_instance_types.ins_type.instance_types: i if lookup(i, "instance_charge_type") == "POSTPAID_BY_HOUR"][1].instance_type
}
`

const TkeExclusiveNetwork = `
data "tencentcloud_vpc_instances" "vpc" {
  name = "` + tkeExclusiveVpcName + `"
}

data "tencentcloud_vpc_subnets" "subnet" {
  vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}

locals {
  vpc_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
}
`

const TkeDataSource = `
data "tencentcloud_kubernetes_clusters" "tke" {
  cluster_name = "` + defaultTkeClusterName + `"
}

locals {
  cluster_id = data.tencentcloud_kubernetes_clusters.tke.list.0.cluster_id
}
`

// End of TKE Service
