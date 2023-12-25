package acctest

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

var Appid string = os.Getenv("TENCENTCLOUD_APPID")
var OwnerUin string = os.Getenv("TENCENTCLOUD_OWNER_UIN")

const (
	KeepResource    = "keep"
	DefaultResource = "Default"
)

var PersistResource = regexp.MustCompile("^(keep|Default)")

// Check if resource should persist instead of recycled
func IsResourcePersist(name string, createdTime *time.Time) bool {
	createdWithin30Minutes := false
	if createdTime != nil {
		createdWithin30Minutes = createdTime.Add(time.Minute * 30).After(time.Now())
	}
	return PersistResource.MatchString(name) || createdWithin30Minutes
}

// vpn
const DefaultVpnDataSource = `
data "tencentcloud_vpn_gateways" "foo" {
  name = "keep-vpn-gw"
}

data "tencentcloud_vpn_connections" "conns" {
  name = "keep-vpn-conn"
}
`

// cos
const (
	DefaultCosCertificateName         = "keep-cos-domain-cert"
	DefaultCosCertificateBucketPrefix = "keep-cert-test"
	DefaultCosCertDomainName          = "mikatong.com"
)

// clb
const (
	DefaultSshCertificate  = "vYSQkJ3K"
	DefaultSshCertificateB = "vYVlNIhW"
)

const (
	DefaultRegion      = "ap-guangzhou"
	DefaultVpcId       = "vpc-86v957zb"
	DefaultVpcCidr     = "172.16.0.0/16"
	DefaultVpcCidrLess = "172.16.0.0/18"

	DefaultCvmAZone                 = "ap-guangzhou-7"
	DefaultCvmInternationalZone     = "ap-guangzhou-3"
	DefaultCvmVpcId                 = "vpc-l0dw94uh"
	DefaultCvmSubnetId              = "subnet-ccj2qg0m"
	DefaultCvmTestingAZone          = "ap-guangzhou-2"
	DefaultCvmTestingVpcId          = "vpc-701bm52d"
	DefaultCvmTestingSubnetId       = "subnet-1q62lj3m"
	DefaultCvmTestingImgId          = "img-eb30mz89"
	DefaultCvmInternationalVpcId    = "vpc-m7ryq37p"
	DefaultCvmInternationalSubnetId = "subnet-lwrsb7a0"

	DefaultAZone          = "ap-guangzhou-3"
	DefaultSubnetId       = "subnet-enm92y0m"
	DefaultSubnetCidr     = "172.16.0.0/20"
	DefaultSubnetCidrLess = "172.16.0.0/22"

	DefaultInsName       = "tf-ci-test"
	DefaultInsNameUpdate = "tf-ci-test-update"

	DefaultDayuBgp    = "bgp-000006mq"
	DefaultDayuBgpMul = "bgp-0000008o"
	DefaultDayuBgpIp  = "bgpip-00000294"
	DefaultDayuNet    = "net-0000007e"

	DefaultGaapProxyId              = "link-ljb08m2l"
	DefaultGaapProxyId2             = "link-8lpyo88p"
	DefaultGaapSecurityPolicyId     = "sp-5lqp4l77"
	DefaultGaapRealserverDomainId1  = "rs-qs0h6wxp"
	DefaultGaapRealserverDomain1    = "github.com"
	DefaultGaapRealserverDomainId2  = "rs-qcygnwpd"
	DefaultGaapRealserverDomain2    = "www.github.com"
	DefaultGaapRealserverIpId1      = "rs-24e1ol23"
	DefaultGaapRealserverIp1        = "119.29.29.35"
	DefaultGaapRealserverIpId2      = "rs-70qzt26p"
	DefaultGaapRealserverIp2        = "1.1.1.5"
	DefaultHttpsDomainCertificateId = "cert-crg2aynt"

	DefaultSecurityGroup = "sg-ijato2x1"

	DefaultProjectId = "1250480"

	DefaultTkeOSImageId   = "img-2lr9q49h"
	DefaultTkeOSImageName = "tlinux2.2(tkernel3)x86_64"
)

// Project
const DefaultProjectVariable = `
variable "default_project" {
  default = ` + DefaultProjectId + `
}
`

// EMR
const (
	DefaultEMRVpcId    = DefaultVpcId
	DefaultEMRSubnetId = DefaultSubnetId
	DefaultEMRSgId     = "sg-694qit0p"
)

const DefaultEMRVariable = `
variable "vpc_id" {
  default = "` + DefaultEMRVpcId + `"
}
variable "subnet_id" {
  default = "` + DefaultEMRSubnetId + `"
}
variable "sg_id" {
  default = "` + DefaultEMRSgId + `"
}
`

// cvm-image
const (
	DefaultCvmId  = "ins-8oqqya08"
	DefaultDiskId = "disk-5jjrs2lm"
	DefaultSnapId = "snap-8f2updnb"
)

const DefaultCvmImageVariable = `
variable "cvm_id" {
  default = "` + DefaultCvmId + `"
}
variable "disk_id" {
  default = "` + DefaultDiskId + `"
}
variable "snap_id" {
  default = "` + DefaultSnapId + `"
}
`

// cvm-modification
const DefaultCommonCvmId = "ins-cr2rfq78"
const DefaultCvmModificationVariable = `
variable "cvm_id" {
  default = "` + DefaultCommonCvmId + `"
}
`

// cvm-reboot
const DefaultRebootCvmId = "ins-f9jr4bd2"
const DefaultRebootCvmVariable = `
variable "cvm_id" {
  default = "` + DefaultRebootCvmId + `"
}
`

// AS
const DefaultAsVariable = `
variable "availability_zone" {
  default = "` + DefaultCvmAZone + `"
}

data "tencentcloud_instance_types" "default" {
  filter {
    name   = "zone"
    values = [var.availability_zone]
  }
  cpu_core_count = 2
  exclude_sold_out = true
}
`

// ckafka
const (
	DefaultKafkaInstanceId            = "ckafka-vv7wpvae"
	DefaultKafkaVpcId                 = "vpc-njbzmzyd"
	DefaultKafkaSubnetId              = "subnet-2txtpql8"
	DefaultKafkaInternationalVpcId    = "vpc-ereuklyj"
	DefaultKafkaInternationalSubnetId = "subnet-e7rvxfx2"
)

const DefaultKafkaVariable = `
variable "instance_id" {
  default = "` + DefaultKafkaInstanceId + `"
}
variable "vpc_id" {
  default = "` + DefaultKafkaVpcId + `"
}
variable "subnet_id" {
  default = "` + DefaultKafkaSubnetId + `"
}
variable "international_vpc_id" {
  default = "` + DefaultKafkaInternationalVpcId + `"
}
variable "international_subnet_id" {
  default = "` + DefaultKafkaInternationalSubnetId + `"
}
`

// Tke Exclusive Network Environment
const (
	TkeExclusiveVpcName   = "keep_tke_exclusive_vpc"
	DefaultTkeClusterId   = "cls-ely08ic4"
	DefaultTkeClusterName = "keep-tke-cluster"
	DefaultTkeClusterType = "tke"
	DefaultPrometheusId   = "prom-1lspn8sw"
	DefaultTemplateId     = "temp-gqunlvo1"
	ClusterPrometheusId   = "prom-g261hacc"
	TkeClusterIdAgent     = "cls-9ae9qo9k"
	TkeClusterTypeAgent   = "eks"
	DefaultAgentId        = "agent-q3zy8gt8"
)

// Cloud monitoring grafana visualization
const (
	DefaultGrafanaVpcId                 = "vpc-391sv4w3"
	DefaultGrafanaSubnetId              = "subnet-ljyn7h30"
	DefaultInternationalGrafanaVpcId    = "vpc-dg21ckzx"
	DefaultInternationalGrafanaSubnetId = "subnet-i5lq9vy4"
	DefaultGrafanaInstanceId            = "grafana-dp2hnnfa"
	DefaultGrafanaReceiver              = "Consumer-nfyxuzmbmq"
	DefaultGrafanaPlugin                = "grafana-clock-panel"
	DefaultGrafanaVersion               = "1.2.0"
)

/*
---------------------------------------------------
The following are common test case used as templates.
---------------------------------------------------
*/

const DefaultVpcVariable = `
variable "instance_name" {
  default = "` + DefaultInsName + `"
}

variable "instance_name_update" {
  default = "` + DefaultInsNameUpdate + `"
}

variable "availability_region" {
  default = "` + DefaultRegion + `"
}

variable "availability_zone" {
  default = "` + DefaultAZone + `"
}

variable "availability_cvm_zone" {
  default = "` + DefaultCvmAZone + `"
}

variable "availability_cvm_international_zone" {
  default = "` + DefaultCvmInternationalZone + `"
}

variable "availability_cvm_testing_zone" {
  default = "` + DefaultCvmTestingAZone + `"
}

variable "cvm_testing_vpc_id" {
  default = "` + DefaultCvmTestingVpcId + `"
}

variable "cvm_testing_subnet_id" {
  default = "` + DefaultCvmTestingSubnetId + `"
}

variable "cvm_testing_image_id" {
  default = "` + DefaultCvmTestingImgId + `"
}

variable "cvm_vpc_id" {
  default = "` + DefaultCvmVpcId + `"
}

variable "cvm_subnet_id" {
  default = "` + DefaultCvmSubnetId + `"
}

variable "cvm_international_vpc_id" {
  default = "` + DefaultCvmInternationalVpcId + `"
}

variable "cvm_international_subnet_id" {
  default = "` + DefaultCvmInternationalSubnetId + `"
}


variable "vpc_id" {
  default = "` + DefaultVpcId + `"
}

variable "vpc_cidr" {
  default = "` + DefaultVpcCidr + `"
}

variable "vpc_cidr_less" {
  default = "` + DefaultVpcCidrLess + `"
}

variable "subnet_id" {
  default = "` + DefaultSubnetId + `"
}

variable "sg_id" {
  default = "` + DefaultSecurityGroup + `"
}

variable "subnet_cidr" {
  default = "` + DefaultSubnetCidr + `"
}

variable "subnet_cidr_less" {
  default = "` + DefaultSubnetCidrLess + `"
}
`

const FixedTagVariable = `
variable "fixed_tags" {
  default = {
    fixed_resource: "do_not_remove"
  }
}
`

const DefaultInstanceVariable = DefaultVpcVariable + `
data "tencentcloud_availability_zones" "default" {
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_images" "testing" {
  image_type = ["PUBLIC_IMAGE"]
}

data "tencentcloud_instance_types" "default" {
  filter {
    name   = "zone"
    values = [var.availability_cvm_zone]
  }
  filter {
    name   = "instance-family"
    values = ["S1", "S2", "S3", "S4", "S5", "SR1", "SA1", "SA2"]
  }
  cpu_core_count = 2
  exclude_sold_out = true
}
`
const DefaultAzVariable = `
variable "default_az" {
  default = "ap-guangzhou-3"
}

variable "default_az7" {
  default = "ap-guangzhou-7"
}
`

const DefaultImages = `
variable "default_img_id" {
  default = "` + DefaultTkeOSImageId + `"
}

variable "default_img" {
  default = "` + DefaultTkeOSImageName + `"
}
`

// Default VPC/Subnet datasource
const DefaultVpcSubnets = DefaultAzVariable + `

data "tencentcloud_vpc_subnets" "gz3" {
  availability_zone = var.default_az
  is_default = true
}

locals {
  vpc_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.gz3.instance_list.0.subnet_id
}`

const DefaultSecurityGroupData = FixedTagVariable + `
data "tencentcloud_security_groups" "internal" {
  name = "default"
  tags = var.fixed_tags
}

data "tencentcloud_security_groups" "exclusive" {
  name = "test_preset_sg"
}

locals {
  # local.sg_id
  sg_id = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
  sg_id2 = data.tencentcloud_security_groups.exclusive.security_groups.0.security_group_id
}
`

const (
	DefaultMySQLName = "keep_preset_mysql"
)

// ref with `local.mysql_id`
const CommonPresetMysql = `

variable "availability_zone" {
  default = "` + DefaultAZone + `"
}
variable "region" {
  default = "` + DefaultRegion + `"
}

data "tencentcloud_mysql_instance" "mysql" {
  instance_name = "` + DefaultMySQLName + `"
}

locals {
  mysql_id = data.tencentcloud_mysql_instance.mysql.instance_list.0.mysql_id
}
`

// SQLServer
const DefaultSQLServerName = "keep-preset_sqlserver"
const DefaultPubSQLServerName = "keep-publish-instance"
const DefaultSubSQLServerName = "keep-subscribe-instance"
const DefaultSQLServerDB = "keep_sqlserver_db"
const DefaultSQLServerPubSubDB = "keep_pubsub_db"
const DefaultSQLServerPubDB = "keep_pub_db"
const DefaultSQLServerSubDB = "keep_sub_db"
const DefaultSQLServerAccount = "keep_sqlserver_account"

const CommonPresetSQLServer = `
data "tencentcloud_sqlserver_instances" "sqlserver" {
  name = "` + DefaultSQLServerName + `"
}

locals {
  # local.sqlserver_id
  sqlserver_id = data.tencentcloud_sqlserver_instances.sqlserver.instance_list.0.id
  sqlserver_db = "` + DefaultSQLServerDB + `"
}
`

const CommonPresetSQLServerAccount = CommonPresetSQLServer + `
data "tencentcloud_sqlserver_accounts" "test"{
  instance_id = local.sqlserver_id
  name = "` + DefaultSQLServerAccount + `"
}

locals {
  # local.sqlserver_id
  sqlserver_account = data.tencentcloud_sqlserver_accounts.test.list.0.name
  sqlserver_pwd = data.tencentcloud_sqlserver_accounts.test.list.0.name
}
`

const TestAccSqlserverAZ = `
data "tencentcloud_availability_zones_by_product" "zone" {
  product = "sqlserver"
}

locals {
  # local.az, local.az1
  az = data.tencentcloud_availability_zones_by_product.zone.zones[0].name
  az1 = data.tencentcloud_availability_zones_by_product.zone.zones[1].name
}
`

const CommonPubSubSQLServer = `
data "tencentcloud_sqlserver_instances" "pub_sqlserver" {
  name = "` + DefaultPubSQLServerName + `"
}
data "tencentcloud_sqlserver_instances" "sub_sqlserver" {
  name = "` + DefaultSubSQLServerName + `"
}

locals {
  pub_sqlserver_id = data.tencentcloud_sqlserver_instances.pub_sqlserver.instance_list.0.id
  sub_sqlserver_id = data.tencentcloud_sqlserver_instances.sub_sqlserver.instance_list.0.id
  sqlserver_pubsub_db = "` + DefaultSQLServerPubSubDB + `"
  sqlserver_pub_db = "` + DefaultSQLServerPubDB + `"
  sqlserver_sub_db = "` + DefaultSQLServerSubDB + `"
}
`

const InstanceCommonTestCase = DefaultInstanceVariable + `
resource "tencentcloud_instance" "default" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
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
const DefaultPGOperationName = "keep-pg-operation"
const OperationPresetPGSQL = `
data "tencentcloud_postgresql_instances" "foo" {
  name = "` + DefaultPGOperationName + `"
}

data "tencentcloud_postgresql_readonly_groups" "ro_groups" {
  filters {
	name = "db-master-instance-id"
	values = [data.tencentcloud_postgresql_instances.foo.instance_list.0.id]
  }
  order_by = "CreateTime"
  order_by_type = "asc"
}

locals {
  pgsql_id = data.tencentcloud_postgresql_instances.foo.instance_list.0.id
  pgrogroup_id = data.tencentcloud_postgresql_readonly_groups.ro_groups.read_only_group_list.0.read_only_group_id
}
`
const DefaultPGSQLName = "keep-postgresql"
const CommonPresetPGSQL = `
data "tencentcloud_postgresql_instances" "foo" {
  name = "` + DefaultPGSQLName + `"
}

locals {
  pgsql_id = data.tencentcloud_postgresql_instances.foo.instance_list.0.id
}
`

// End of PostgreSQL

const DefaultCVMName = "keep-cvm"
const PresetCVM = `
data "tencentcloud_instances" "instance" {
  instance_name = "` + DefaultCVMName + `"
}

locals {
  cvm_id = data.tencentcloud_instances.instance.instance_list.0.instance_id
  cvm_az = "` + DefaultAZone + `"
  cvm_private_ip = data.tencentcloud_instances.instance.instance_list.0.private_ip
}
`

const UserInfoData = `
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
  uin = data.tencentcloud_user_info.info.uin
  owner_uin = data.tencentcloud_user_info.info.owner_uin
}
`

const DefaultSCFCosBucket = `
data "tencentcloud_user_info" "info" {}

data "tencentcloud_cos_buckets" "buckets" {
  bucket_prefix = "preset-scf-bucket-${data.tencentcloud_user_info.info.app_id}"
}

locals {
  bucket_name = data.tencentcloud_cos_buckets.buckets.bucket_list.0.bucket
  bucket_url = data.tencentcloud_cos_buckets.buckets.bucket_list.0.cos_bucket_url
}
`

const DefaultScfNamespace = "preset-scf-namespace"

const DefaultFileSystemName = "keep_preset_cfs"

const DefaultFileSystem = `
data "tencentcloud_cfs_file_systems" "fs" {
  name = "` + DefaultFileSystemName + `"
}

# doesn't support datasource for now
variable "mount_id" {
  default = "cfs-iobiaxtj"
}

locals {
  cfs = data.tencentcloud_cfs_file_systems.fs.file_system_list.0
  cfs_id = local.cfs.file_system_id
}`

const DefaultCamVariables = `
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
const DefaultTCRInstanceName = "keep-tcr-instance-sh"
const DefaultTCRInstanceId = "tcr-aoz8mxoz"
const DefaultTCRNamespace = "keep-tcr-namespace-sh"
const DefaultTCRRepoName = "keep-tcr-repo-sh"
const DefaultTCRSSL = "zjUMifFK"

const DefaultTCRInstanceVar = `
variable "tcr_name" {
  default = "` + DefaultTCRInstanceName + `"
}

variable "tcr_namespace" {
  default = "` + DefaultTCRNamespace + `"
}

variable "tcr_repo" {
  default = "` + DefaultTCRRepoName + `"
}
`

const DefaultTCRInstanceData = DefaultTCRInstanceVar + `
data "tencentcloud_tcr_instances" "tcr" {
  name = var.tcr_name
}

locals {
  tcr_id = data.tencentcloud_tcr_instances.tcr.instance_list.0.id
}
`

// End of TCR Service

// TcaPlus DB

const DefaultTcaPlusClusterName = "keep-tcaplus-cluster"
const DefaultTcaPlusClusterTableGroup = "keep_table_group"
const DefaultTcaPlusTdrClusterName = "keep_tdr_tcaplus_cluster"
const DefaultTcaPlusTdrClusterTableGroup = "keep_tdr_table_group"
const DefaultTcaPlusClusterTable = "keep_players"
const DefaultTcaPlusVar = `
variable "tcaplus_cluster" {
  default = "` + DefaultTcaPlusClusterName + `"
}

variable "tcaplus_table_group" {
  default = "` + DefaultTcaPlusClusterTableGroup + `"
}

variable "tcaplus_table" {
  default = "` + DefaultTcaPlusClusterTable + `"
}

variable "tcaplus_tcr_cluster" {
  default = "` + DefaultTcaPlusTdrClusterName + `"
}

variable "tcaplus_tcr_table_group" {
  default = "` + DefaultTcaPlusTdrClusterTableGroup + `"
}
`
const DefaultTcaPlusData = DefaultTcaPlusVar + `
data "tencentcloud_tcaplus_clusters" "tcaplus" {
  cluster_name = var.tcaplus_cluster
}

data "tencentcloud_tcaplus_tablegroups" "group" {
  cluster_id = data.tencentcloud_tcaplus_clusters.tcaplus.list.0.cluster_id
  tablegroup_name = var.tcaplus_table_group
}

data "tencentcloud_tcaplus_clusters" "tdr_tcaplus" {
  cluster_name = "keep_tdr_tcaplus_cluster"
}
  
data "tencentcloud_tcaplus_tablegroups" "tdr_group" {
  cluster_id = data.tencentcloud_tcaplus_clusters.tdr_tcaplus.list.0.cluster_id
  tablegroup_name = "keep_tdr_table_group"
}

locals {
  tcaplus_id = data.tencentcloud_tcaplus_clusters.tcaplus.list.0.cluster_id
  tcr_tcaplus_id = data.tencentcloud_tcaplus_clusters.tdr_tcaplus.list.0.cluster_id
  tcaplus_table_group = var.tcaplus_table_group
  tcaplus_table_group_id = data.tencentcloud_tcaplus_tablegroups.group.list.0.tablegroup_id
  tcr_tcaplus_table_group_id = data.tencentcloud_tcaplus_tablegroups.tdr_group.list.0.tablegroup_id
  tcaplus_table = var.tcaplus_table
}
`

// End of TcaPlus DB

// TKE Service

// List sample CIDRs to avoid conflict when running multiple cluster testcase parallel
const TkeCIDRs = `
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
`

const TkeDefaultNodeInstanceVar = `
variable "ins_type" {
  default = "SA2.LARGE8"
}
`

// @deprecated. Avoid using this because it may return diff results
const TkeInstanceType = `
data "tencentcloud_instance_types" "ins_type" {
  availability_zone = "` + DefaultCvmAZone + `"
  cpu_core_count    = 2
  exclude_sold_out  = true
}

locals {
  ins_az = "` + DefaultCvmAZone + `"
  type1 = [for i in data.tencentcloud_instance_types.ins_type.instance_types: i if lookup(i, "instance_charge_type") == "POSTPAID_BY_HOUR"]
  type2 = [for i in data.tencentcloud_instance_types.ins_type.instance_types: i]
  final_type = concat(local.type1, local.type2)[0].instance_type
}
`

const TkeExclusiveNetwork = DefaultAzVariable + `
data "tencentcloud_vpc_instances" "vpc" {
  name = "` + TkeExclusiveVpcName + `"
}

data "tencentcloud_vpc_subnets" "subnet" {
  vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}

data "tencentcloud_instance_types" "default" {
	filter {
	  name   = "zone"
	  values = [var.default_az]
	}
  filter {
    name   = "instance-charge-type"
    values = ["POSTPAID_BY_HOUR"]
  }
	cpu_core_count = 2
	exclude_sold_out = true
}

locals {
  vpc_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
  scale_instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
}
`

const TkeDataSource = `
data "tencentcloud_kubernetes_clusters" "tke" {
  cluster_name = "` + DefaultTkeClusterName + `"
}

locals {
  cluster_id = data.tencentcloud_kubernetes_clusters.tke.list.0.cluster_id
}
`

// End of TKE Service

// MongoDB
const (
	DefaultMongoDBVPCId    = "vpc-rwj54lfr"
	DefaultMongoDBSubnetId = "subnet-nyt57zps"
)
const DefaultMongoDBSecurityGroupId = "sg-if748odn"
const DefaultMongoDBSpec = `
data "tencentcloud_mongodb_zone_config" "zone_config" {
  available_zone = "ap-guangzhou-6"
}

data "tencentcloud_security_group" "foo" {
  name = "default"
}

variable "engine_versions" {
  default = {
    "3.6": "MONGO_36_WT",
    "4.0": "MONGO_40_WT",
    "4.2": "MONGO_42_WT"
    "4.4": "MONGO_44_WT"
  }
}
variable "sg_id" {
  default = "` + DefaultMongoDBSecurityGroupId + `"
}
variable "vpc_id" {
  default = "` + DefaultMongoDBVPCId + `"
}
variable "subnet_id" {
  default = "` + DefaultMongoDBSubnetId + `"
}

locals {
  filtered_spec = [for i in data.tencentcloud_mongodb_zone_config.zone_config.list: i if lookup(i, "machine_type") == "HIO10G" && lookup(i, "engine_version") != "3.2"]
  spec = concat(local.filtered_spec, data.tencentcloud_mongodb_zone_config.zone_config.list)
  machine_type = local.spec.0.machine_type
  cluster_type = local.spec.0.cluster_type
  memory = local.spec.0.memory / 1024
  volume = local.spec.0.min_storage / 1024
  engine_version = lookup(var.engine_versions, local.spec.0.engine_version)
  security_group_id = data.tencentcloud_security_group.foo.id
}

locals {
  filtered_sharding_spec = [for i in data.tencentcloud_mongodb_zone_config.zone_config.list: i if lookup(i, "cluster_type") == "SHARD" && lookup(i, "min_replicate_set_num") > 0 && lookup(i, "machine_type") == "HIO10G" && lookup(i, "engine_version") != "3.2"]
  sharding_spec = concat(local.filtered_sharding_spec, [for i in data.tencentcloud_mongodb_zone_config.zone_config.list: i if lookup(i, "cluster_type") == "SHARD" && lookup(i, "min_replicate_set_num") > 0])
  sharding_machine_type = local.sharding_spec.0.machine_type
  sharding_memory = local.sharding_spec.0.memory / 1024
  sharding_volume = local.sharding_spec.0.min_storage / 1024
  sharding_engine_version = lookup(var.engine_versions, local.sharding_spec.0.engine_version)
}
`

// End of MongoDB

// TCM

const (
	DefaultMeshClusterId = "cls-9ae9qo9k"
	DefaultMeshVpcId     = "vpc-pyyv5k3v"
	DefaultMeshSubnetId  = "subnet-06i8auk6"
)

// End of TCM

// DCDB
const (
	DefaultDcdbInstanceId    = "tdsqlshard-lgz66iqr"
	DefaultDcdbInstanceName  = "keep-dcdb-test"
	DefaultDcdbInsVpcId      = "vpc-4owdpnwr"
	DefaultDcdbInsIdSubnetId = "subnet-qylstu34"
	DefaultDcdbSGId          = "sg-ijato2x1"
	DefaultDcdbSGName        = "default"
)

// ref with `local.dcdb_id`
const CommonPresetDcdb = `

variable "availability_zone" {
  default = "` + DefaultAZone + `"
}
variable "region" {
  default = "` + DefaultRegion + `"
}

data "tencentcloud_dcdb_instances" "dcdb" {
  search_name = "instancename"
  search_key = "` + DefaultDcdbInstanceName + `"
}

locals {
  dcdb_id = data.tencentcloud_dcdb_instances.dcdb.list.0.instance_id
}
`

// ref with `local.redis_id`
const CommonPresetRedis = `
locals {
  redis_id = "crs-jf4ico4v"
  redis_name = "Keep-terraform"
}
`

// End of DCDB
// SES
const (
	DefaultRegionSes = "ap-hongkong"
)

// End of SES
// MARIADB
const (
	DefaultMariadbSubnetId        = "subnet-jdi5xn22"
	DefaultMariadbVpcId           = "vpc-k1t8ickr"
	DefaultMariadbSecurityGroupId = "sg-7kpsbxdb"

	DefaultMariadbInstanceSubnetId = "subnet-4w4twlf4"
	DefaultMariadbInstanceVpcId    = "vpc-9m66acml"
)

// End of MARIADB
// PTS
const (
	DefaultPtsProjectId  = "project-45vw7v82"
	DefaultScenarioId    = "scenario-gb5ix8m2"
	DefaultScenarioIdJob = "scenario-22q19f3k"
	DefaultPtsNoticeId   = "notice-tj75hgqj"
)

// End of PTS

// CSS
const (
	DefaultCSSLiveType        = "PullLivePushLive"
	DefaultCSSDomainName      = "177154.push.tlivecloud.com"
	DefaultCSSStreamName      = DefaultCSSPrefix + "test_stream_name"
	DefaultCSSAppName         = "live"
	DefaultCSSOperator        = "tf_admin"
	DefaultCSSPrefix          = "tf_css_"
	DefaultCSSPlayDomainName  = "test122.jingxhu.top"
	DefaultCSSPushDomainName  = "177154.push.tlivecloud.com"
	DefaultCSSBindingCertName = "keep_ssl_css_domain_test"
)

// End of CSS

// TAT
const (
	DefaultInstanceId = "ins-881b1c8w"
	DefaultCommandId  = "cmd-rxbs7f5z"
)

// End of TAT

// TDCPG
const (
	DefaultTdcpgClusterId      = "tdcpg-m5e26fi8"
	DefaultTdcpgClusterName    = "keep-tdcpg-test"
	DefaultTdcpgPayMode        = "POSTPAID_BY_HOUR"
	DefaultTdcpgInstanceId     = "tdcpg-ins-fc0e5kes"
	DefaultTdcpgInstanceName   = "keep-tdcpg-instance-test"
	DefaultTdcpgZone           = "ap-guangzhou-3"
	DefaultTdcpgTestNamePrefix = "tf-tdcpg-"
)

// End of TDCPG

// DBBRAIN
const (
	DefaultDbBrainsagId      = "sag-01z37l4g"
	DefaultDbBrainInstanceId = "cdb-fitq5t9h"
)

// End of DBBRAIN

// RUM
const (
	DefaultRumInstanceId = "rum-pasZKEI3RLgakj"
	DefaultRumProjectId  = "131407"
)

// End of RUM

// DTS
const (
	DefaultDTSJobId = "dts-r5gpejpe"
)

// End of DTS

// TEM
const (
	DefaultEnvironmentId = "en-85377m6j"
	DefaultApplicationId = "app-joqr9bd5"
	DefaultTemVpcId      = "vpc-4owdpnwr"
	DefaultTemSubnetId   = "subnet-c1l35990"
	DefaultLogsetId      = "33aaf0ae-6163-411b-a415-9f27450f68db"
	DefaultTopicId       = "88735a07-bea4-4985-8763-e9deb6da4fad"
)

// End of TEM

// CI
const (
	DefaultCiBucket  = "terraform-ci-1308919341"
	DefaultStyleName = "terraform_test"
)

// End of CI

// Cynosdb
const (
	DefaultCynosdbClusterId         = "cynosdbmysql-bws8h88b"
	DefaultCynosdbClusterInstanceId = "cynosdbmysql-ins-afqx1hy0"
	DefaultCynosdbSecurityGroup     = "sg-baxfiao5"
)

const CommonCynosdb = `

variable "cynosdb_cluster_id" {
  default = "` + DefaultCynosdbClusterId + `"
}
variable "cynosdb_cluster_instance_id" {
  default = "` + DefaultCynosdbClusterInstanceId + `"
}
variable "cynosdb_cluster_security_group_id" {
  default = "` + DefaultCynosdbSecurityGroup + `"
}
`

// End of Cynosdb

// TSF
const (
	DefaultNamespaceId         = "namespace-aemrg36v"
	DefaultTsfApplicationId    = "application-a24x29xv"
	DefaultTsfClustId          = "cluster-vwgj5e6y"
	DefaultTsfGroupId          = "group-yrjkln9v"
	DefaultTsfGateway          = "gw-ins-lvdypq5k"
	DefaultTsfDestNamespaceId  = "namespace-aemrg36v"
	DefaultTsfConfigId         = "dcfg-y54wzk3a"
	DefaultTsfApiId            = "api-j03q029a"
	DefaultTsfGWGroupId        = "group-vzd97zpy"
	DefaultTsfFileConfigId     = "dcfg-f-ab6l9x5y"
	DefaultTsfImageId          = "img-7r9vq8wd"
	DefaultTsfGWNamespaceId    = "namespace-vwgo38wy"
	DefaultTsfContainerGroupId = "group-y43x5jpa"
	DefaultTsfpodName          = "keep-terraform-7f4874bc5c-w75q4"
)

// End of TSF

// CBS
const DefaultCbsBackupDiskId = "disk-r69pg9vw"

const CbsBackUp = `
variable "cbs_backup_disk_id" {
  default = "` + DefaultCbsBackupDiskId + `"
}
`

// End of CBS

// CRS
const (
	DefaultCrsInstanceId     = "crs-jf4ico4v"
	DefaultCrsVpcId          = "vpc-4owdpnwr"
	DefaultCrsSubnetId       = "subnet-4o0zd840"
	DefaultCrsSecurityGroups = "sg-edmur627"
	DefaultCrsGroupId        = "crs-rpl-orfiwmn5"
)

const DefaultCrsVar = `
variable "vpc_id" {
  default = "` + DefaultCrsVpcId + `"
}
variable "subnet_id" {
  default = "` + DefaultCrsSubnetId + `"
}
`

// End of CRS

const (
	DefaultLighthouseInstanceId   = "lhins-g4bwdjbf"
	DefaultLighthoustDiskId       = "lhdisk-do4p4hz6"
	DefaultLighthouseBackupDiskId = "lhdisk-cwodsc4q"
	DefaultLighthouseBackUpId     = "lhbak-bpum3ygx"
	DefaultLighthouseSnapshotId   = "lhsnap-o2mvsvk3"
)

const DefaultLighthoustVariables = `
variable "lighthouse_id" {
  default = "` + DefaultLighthouseInstanceId + `"
}

variable "lighthouse_disk_id" {
  default = "` + DefaultLighthoustDiskId + `"
}

variable "lighthouse_backup_disk_id" {
  default = "` + DefaultLighthouseBackupDiskId + `"
}

variable "lighthouse_backup_id" {
  default = "` + DefaultLighthouseBackUpId + `"
}

variable "lighthouse_snapshot_id" {
  default = "` + DefaultLighthouseSnapshotId + `"
}
`

// TSE
const (
	DefaultEngineResourceSpec = "spec-qvj6k7t4q"
	DefaultTseVpcId           = "vpc-4owdpnwr"
	DefaultTseSubnetId        = "subnet-dwj7ipnc"
	DefaultTseGatewayId       = "gateway-ddbb709b"
	DefaultTseCertId          = "vYSQkJ3K"
)

const DefaultTseVar = `
variable "gateway_id" {
  default = "` + DefaultTseGatewayId + `"
}
variable "cert_id" {
  default = "` + DefaultTseCertId + `"
}
`

// End of TSE

// ES
const (
	DefaultEsInstanceId    = "es-5wn36he6"
	DefaultEsSecurityGroup = "sg-edmur627"
	DefaultEsLogstash      = "ls-kru90fkz"
)

const DefaultEsVariables = `
variable "instance_id" {
  default = "` + DefaultEsInstanceId + `"
}

variable "security_group_id" {
  default = "` + DefaultEsSecurityGroup + `"
}

variable "logstash_id" {
  default = "` + DefaultEsLogstash + `"
}
`

// End of TSE

// Clickhouse
const DefaultClickhouseInstanceId = "cdwch-pcap78rz"

const DefaultClickhouseVariables = `
variable "instance_id" {
  default = "` + DefaultClickhouseInstanceId + `"
}
`

// End of Clickhouse

// CLB
const ClbTargetEniTestCase = InstanceCommonTestCase + `
resource "tencentcloud_eni" "clb_eni_target" {
  name        = "ci-test-eni"
  vpc_id      = var.cvm_vpc_id
  subnet_id   = var.cvm_subnet_id
  description = "clb eni backend"
  ipv4_count  = 1
}

resource "tencentcloud_eni_attachment" "foo" {
  eni_id      = tencentcloud_eni.clb_eni_target.id
  instance_id = tencentcloud_instance.default.id
}
`

//End of Clb

// MPS
const (
	DefaultMpsScheduleId   = 24685
	DefaultMpsScheduleName = "keep_mps_schedule_001"
)

//End of MPS
