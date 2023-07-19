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

// cos
const (
	defaultCosCertificateName         = "keep-cos-domain-cert"
	defaultCosCertificateBucketPrefix = "keep-cert-test"
	defaultCosCertDomainName          = "mikatong.com"
)

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

	defaultGaapProxyId              = "link-ljb08m2l"
	defaultGaapProxyId2             = "link-8lpyo88p"
	defaultGaapSecurityPolicyId     = "sp-5lqp4l77"
	defaultGaapRealserverDomainId1  = "rs-qs0h6wxp"
	defaultGaapRealserverDomain1    = "github.com"
	defaultGaapRealserverDomainId2  = "rs-qcygnwpd"
	defaultGaapRealserverDomain2    = "www.github.com"
	defaultGaapRealserverIpId1      = "rs-24e1ol23"
	defaultGaapRealserverIp1        = "119.29.29.35"
	defaultGaapRealserverIpId2      = "rs-70qzt26p"
	defaultGaapRealserverIp2        = "1.1.1.5"
	defaultHttpsDomainCertificateId = "cert-crg2aynt"

	defaultSecurityGroup = "sg-ijato2x1"

	defaultProjectId = "1250480"

	defaultTkeOSImageId   = "img-2lr9q49h"
	defaultTkeOSImageName = "tlinux2.2(tkernel3)x86_64"
)

// Project
const defaultProjectVariable = `
variable "default_project" {
  default = ` + defaultProjectId + `
}
`

// EMR
const (
	defaultEMRVpcId    = defaultVpcId
	defaultEMRSubnetId = defaultSubnetId
	defaultEMRSgId     = "sg-694qit0p"
)

const defaultEMRVariable = `
variable "vpc_id" {
  default = "` + defaultEMRVpcId + `"
}
variable "subnet_id" {
  default = "` + defaultEMRSubnetId + `"
}
variable "sg_id" {
  default = "` + defaultEMRSgId + `"
}
`

// cvm-image
const (
	defaultCvmId  = "ins-8oqqya08"
	defaultDiskId = "disk-5jjrs2lm"
	defaultSnapId = "snap-8f2updnb"
)

const defaultCvmImageVariable = `
variable "cvm_id" {
  default = "` + defaultCvmId + `"
}
variable "disk_id" {
  default = "` + defaultDiskId + `"
}
variable "snap_id" {
  default = "` + defaultSnapId + `"
}
`

// cvm-modification
const defaultCommonCvmId = "ins-cr2rfq78"
const defaultCvmModificationVariable = `
variable "cvm_id" {
  default = "` + defaultCommonCvmId + `"
}
`

// cvm-reboot
const defaultRebootCvmId = "ins-f9jr4bd2"
const defaultRebootCvmVariable = `
variable "cvm_id" {
  default = "` + defaultRebootCvmId + `"
}
`

// AS
const defaultAsVariable = `
variable "availability_zone" {
  default = "` + defaultCvmAZone + `"
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
	defaultTkeClusterType = "tke"
	defaultPrometheusId   = "prom-1lspn8sw"
	defaultTemplateId     = "temp-gqunlvo1"
	clusterPrometheusId   = "prom-g261hacc"
	tkeClusterIdAgent     = "cls-9ae9qo9k"
	tkeClusterTypeAgent   = "eks"
	defaultAgentId        = "agent-q3zy8gt8"
)

// Cloud monitoring grafana visualization
const (
	defaultGrafanaVpcId      = "vpc-391sv4w3"
	defaultGrafanaSubnetId   = "subnet-ljyn7h30"
	defaultGrafanaInstanceId = "grafana-dp2hnnfa"
	defaultGrafanaReceiver   = "Consumer-nfyxuzmbmq"
	defaultGrafanaPlugin     = "grafana-clock-panel"
	defaultGrafanaVersion    = "1.2.0"
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

variable "sg_id" {
  default = "` + defaultSecurityGroup + `"
}

variable "subnet_cidr" {
  default = "` + defaultSubnetCidr + `"
}

variable "subnet_cidr_less" {
  default = "` + defaultSubnetCidrLess + `"
}
`

const fixedTagVariable = `
variable "fixed_tags" {
  default = {
    fixed_resource: "do_not_remove"
  }
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

const defaultSecurityGroupData = fixedTagVariable + `
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
const defaultPubSQLServerName = "keep-publish-instance"
const defaultSubSQLServerName = "keep-subscribe-instance"
const defaultSQLServerDB = "keep_sqlserver_db"
const defaultSQLServerPubSubDB = "keep_pubsub_db"
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
  sqlserver_pwd = data.tencentcloud_sqlserver_accounts.test.list.0.name
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

const CommonPubSubSQLServer = `
data "tencentcloud_sqlserver_instances" "pub_sqlserver" {
  name = "` + defaultPubSQLServerName + `"
}
data "tencentcloud_sqlserver_instances" "sub_sqlserver" {
  name = "` + defaultSubSQLServerName + `"
}

locals {
  pub_sqlserver_id = data.tencentcloud_sqlserver_instances.pub_sqlserver.instance_list.0.id
  sub_sqlserver_id = data.tencentcloud_sqlserver_instances.sub_sqlserver.instance_list.0.id
  sqlserver_pubsub_db = "` + defaultSQLServerPubSubDB + `"
}
`

const instanceCommonTestCase = defaultInstanceVariable + `
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
const defaultPGOperationName = "keep-pg-operation"
const OperationPresetPGSQL = `
data "tencentcloud_postgresql_instances" "foo" {
  name = "` + defaultPGOperationName + `"
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
  cvm_private_ip = data.tencentcloud_instances.instance.instance_list.0.private_ip
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
const defaultTCRInstanceName = "keep-tcr-instance-sh"
const defaultTCRInstanceId = "tcr-800uf6j9"
const defaultTCRNamespace = "keep-tcr-namespace-sh"
const defaultTCRRepoName = "keep-tcr-repo-sh"
const defaultTCRSSL = "zjUMifFK"

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
const defaultTcaPlusTdrClusterName = "keep_tdr_tcaplus_cluster"
const defaultTcaPlusTdrClusterTableGroup = "keep_tdr_table_group"
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

variable "tcaplus_tcr_cluster" {
  default = "` + defaultTcaPlusTdrClusterName + `"
}

variable "tcaplus_tcr_table_group" {
  default = "` + defaultTcaPlusTdrClusterTableGroup + `"
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
  availability_zone = "` + defaultCvmAZone + `"
  cpu_core_count    = 2
  exclude_sold_out  = true
}

locals {
  ins_az = "` + defaultCvmAZone + `"
  type1 = [for i in data.tencentcloud_instance_types.ins_type.instance_types: i if lookup(i, "instance_charge_type") == "POSTPAID_BY_HOUR"]
  type2 = [for i in data.tencentcloud_instance_types.ins_type.instance_types: i]
  final_type = concat(local.type1, local.type2)[0].instance_type
}
`

const TkeExclusiveNetwork = defaultAzVariable + `
data "tencentcloud_vpc_instances" "vpc" {
  name = "` + tkeExclusiveVpcName + `"
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
  cluster_name = "` + defaultTkeClusterName + `"
}

locals {
  cluster_id = data.tencentcloud_kubernetes_clusters.tke.list.0.cluster_id
}
`

// End of TKE Service

// MongoDB
const (
	defaultMongoDBVPCId    = "vpc-rwj54lfr"
	defaultMongoDBSubnetId = "subnet-nyt57zps"
)
const defaultMongoDBSecurityGroupId = "sg-if748odn"
const DefaultMongoDBSpec = `
data "tencentcloud_mongodb_zone_config" "zone_config" {
  available_zone = "ap-guangzhou-6"
}

variable "engine_versions" {
  default = {
    "3.6": "MONGO_36_WT",
    "4.0": "MONGO_40_WT",
    "4.2": "MONGO_42_WT"
  }
}
variable "sg_id" {
  default = "` + defaultMongoDBSecurityGroupId + `"
}
variable "vpc_id" {
  default = "` + defaultMongoDBVPCId + `"
}
variable "subnet_id" {
  default = "` + defaultMongoDBSubnetId + `"
}

locals {
  filtered_spec = [for i in data.tencentcloud_mongodb_zone_config.zone_config.list: i if lookup(i, "machine_type") == "HIO10G" && lookup(i, "engine_version") != "3.2"]
  spec = concat(local.filtered_spec, data.tencentcloud_mongodb_zone_config.zone_config.list)
  machine_type = local.spec.0.machine_type
  cluster_type = local.spec.0.cluster_type
  memory = local.spec.0.memory / 1024
  volume = local.spec.0.min_storage / 1024
  engine_version = lookup(var.engine_versions, local.spec.0.engine_version)
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

// TEO

const (
	defaultZoneName    = "tf-teo-t.xyz"
	defaultZoneId      = "zone-2a1u0y616jz6"
	defaultPolicyId    = "11581"
	applicationProxyId = "sid-2f8xqhmf32dt"
)

// End of TEO

// TCM

const (
	defaultMeshClusterId = "cls-9ae9qo9k"
	defaultMeshVpcId     = "vpc-pyyv5k3v"
	defaultMeshSubnetId  = "subnet-06i8auk6"
)

// End of TCM

// DCDB
const (
	defaultDcdbInstanceId    = "tdsqlshard-lgz66iqr"
	defaultDcdbInstanceName  = "keep-dcdb-test"
	defaultDcdbInsVpcId      = "vpc-4owdpnwr"
	defaultDcdbInsIdSubnetId = "subnet-qylstu34"
	defaultDcdbSGId          = "sg-ijato2x1"
	defaultDcdbSGName        = "default"
)

// ref with `local.dcdb_id`
const CommonPresetDcdb = `

variable "availability_zone" {
  default = "` + defaultAZone + `"
}
variable "region" {
  default = "` + defaultRegion + `"
}

data "tencentcloud_dcdb_instances" "dcdb" {
  search_name = "instancename"
  search_key = "` + defaultDcdbInstanceName + `"
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
	defaultRegionSes = "ap-hongkong"
)

// End of SES
// MARIADB
const (
	defaultMariadbSubnetId        = "subnet-jdi5xn22"
	defaultMariadbVpcId           = "vpc-k1t8ickr"
	defaultMariadbSecurityGroupId = "sg-7kpsbxdb"

	defaultMariadbInstanceSubnetId = "subnet-4w4twlf4"
	defaultMariadbInstanceVpcId    = "vpc-9m66acml"
)

// End of MARIADB
// PTS
const (
	defaultPtsProjectId  = "project-45vw7v82"
	defaultScenarioId    = "scenario-gb5ix8m2"
	defaultScenarioIdJob = "scenario-22q19f3k"
	defaultPtsNoticeId   = "notice-tj75hgqj"
)

// End of PTS

// CSS
const (
	defaultCSSLiveType        = "PullLivePushLive"
	defaultCSSDomainName      = "177154.push.tlivecloud.com"
	defaultCSSStreamName      = defaultCSSPrefix + "test_stream_name"
	defaultCSSAppName         = "live"
	defaultCSSOperator        = "tf_admin"
	defaultCSSPrefix          = "tf_css_"
	defaultCSSPlayDomainName  = "test122.jingxhu.top"
	defaultCSSPushDomainName  = "177154.push.tlivecloud.com"
	defaultCSSBindingCertName = "keep_ssl_css_domain_test"
)

// End of CSS

// TAT
const (
	defaultInstanceId = "ins-881b1c8w"
	defaultCommandId  = "cmd-rxbs7f5z"
)

// End of TAT

// TDCPG
const (
	defaultTdcpgClusterId      = "tdcpg-m5e26fi8"
	defaultTdcpgClusterName    = "keep-tdcpg-test"
	defaultTdcpgPayMode        = "POSTPAID_BY_HOUR"
	defaultTdcpgInstanceId     = "tdcpg-ins-fc0e5kes"
	defaultTdcpgInstanceName   = "keep-tdcpg-instance-test"
	defaultTdcpgZone           = "ap-guangzhou-3"
	defaultTdcpgTestNamePrefix = "tf-tdcpg-"
)

// End of TDCPG

// DBBRAIN
const (
	defaultDbBrainsagId      = "sag-01z37l4g"
	defaultDbBrainInstanceId = "cdb-fitq5t9h"
)

// End of DBBRAIN

// RUM
const (
	defaultRumInstanceId = "rum-pasZKEI3RLgakj"
)

// End of RUM

// DTS
const (
	defaultDTSJobId = "dts-r5gpejpe"
)

// End of DTS

// TEM
const (
	defaultEnvironmentId = "en-85377m6j"
	defaultApplicationId = "app-joqr9bd5"
	defaultTemVpcId      = "vpc-4owdpnwr"
	defaultTemSubnetId   = "subnet-c1l35990"
	defaultLogsetId      = "33aaf0ae-6163-411b-a415-9f27450f68db"
	defaultTopicId       = "88735a07-bea4-4985-8763-e9deb6da4fad"
)

// End of TEM

// CI
const (
	defaultCiBucket  = "terraform-ci-1308919341"
	defaultStyleName = "terraform_test"
)

// End of CI

// Cynosdb
const (
	defaultCynosdbClusterId         = "cynosdbmysql-bws8h88b"
	defaultCynosdbClusterInstanceId = "cynosdbmysql-ins-afqx1hy0"
	defaultCynosdbSecurityGroup     = "sg-baxfiao5"
)

const CommonCynosdb = `

variable "cynosdb_cluster_id" {
  default = "` + defaultCynosdbClusterId + `"
}
variable "cynosdb_cluster_instance_id" {
  default = "` + defaultCynosdbClusterInstanceId + `"
}
variable "cynosdb_cluster_security_group_id" {
  default = "` + defaultCynosdbSecurityGroup + `"
}
`

// End of Cynosdb

// TSF
const (
	defaultNamespaceId        = "namespace-aemrg36v"
	defaultTsfApplicationId   = "application-a24x29xv"
	defaultTsfClustId         = "cluster-vwgj5e6y"
	defaultTsfGroupId         = "group-yrjkln9v"
	defaultTsfGateway         = "gw-ins-lvdypq5k"
	defaultTsfDestNamespaceId = "namespace-aemrg36v"
	defaultTsfConfigId        = "dcfg-y54wzk3a"
	defaultTsfApiId           = "api-j03q029a"
	defaultTsfGWGroupId       = "group-vzd97zpy"
	defaultTsfFileConfigId    = "dcfg-f-ab6l9x5y"
	defaultTsfImageId         = "img-7r9vq8wd"
	defaultTsfGWNamespaceId   = "namespace-vwgo38wy"
)

// End of TSF

// CBS
const defaultCbsBackupDiskId = "disk-r69pg9vw"

const CbsBackUp = `
variable "cbs_backup_disk_id" {
  default = "` + defaultCbsBackupDiskId + `"
}
`

// End of CBS

// CRS
const (
	defaultCrsInstanceId     = "crs-jf4ico4v"
	defaultCrsVpcId          = "vpc-4owdpnwr"
	defaultCrsSubnetId       = "subnet-4o0zd840"
	defaultCrsSecurityGroups = "sg-edmur627"
	defaultCrsGroupId        = "crs-rpl-orfiwmn5"
)

const DefaultCrsVar = `
variable "vpc_id" {
  default = "` + defaultCrsVpcId + `"
}
variable "subnet_id" {
  default = "` + defaultCrsSubnetId + `"
}
`

// End of CRS

const (
	defaultLighthouseInstanceId   = "lhins-hwe21u91"
	defaultLighthoustDiskId       = "lhdisk-do4p4hz6"
	defaultLighthouseBackupDiskId = "lhdisk-cwodsc4q"
	defaultLighthouseBackUpId     = "lhbak-bpum3ygx"
	defaultLighthouseSnapshotId   = "lhsnap-9jaw9m17"
)

const DefaultLighthoustVariables = `
variable "lighthouse_id" {
  default = "` + defaultLighthouseInstanceId + `"
}

variable "lighthouse_disk_id" {
  default = "` + defaultLighthoustDiskId + `"
}

variable "lighthouse_backup_disk_id" {
  default = "` + defaultLighthouseBackupDiskId + `"
}

variable "lighthouse_backup_id" {
  default = "` + defaultLighthouseBackUpId + `"
}

variable "lighthouse_snapshot_id" {
  default = "` + defaultLighthouseSnapshotId + `"
}
`

// TSE
const (
	defaultEngineResourceSpec = "spec-qvj6k7t4q"
	defaultTseVpcId           = "vpc-4owdpnwr"
	defaultTseSubnetId        = "subnet-dwj7ipnc"
)

// End of TSE

// ES
const (
	defaultEsInstanceId    = "es-5wn36he6"
	defaultEsSecurityGroup = "sg-edmur627"
)

const DefaultEsVariables = `
variable "instance_id" {
  default = "` + defaultEsInstanceId + `"
}

variable "security_group_id" {
  default = "` + defaultEsSecurityGroup + `"
}
`

// End of TSE
