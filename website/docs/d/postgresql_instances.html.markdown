---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_instances"
sidebar_current: "docs-tencentcloud-datasource-postgresql_instances"
description: |-
  Use this data source to query PostgreSQL instances
---

# tencentcloud_postgresql_instances

Use this data source to query PostgreSQL instances

## Example Usage

### Query all postgresql instances

```hcl
data "tencentcloud_postgresql_instances" "example" {}
```

### Query postgresql instances by filters

```hcl
data "tencentcloud_postgresql_instances" "example" {
  id         = "postgres-gngyhl9d"
  name       = "tf-example"
  project_id = "1235143"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) ID of the postgresql instance to be query.
* `name` - (Optional, String) Name of the postgresql instance to be query.
* `project_id` - (Optional, String) Project ID of the postgresql instance to be query.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `db_instance_set` - Instance details set.
  * `app_id` - User `AppId`.
  * `auto_renew` - Auto-renewal or not:
<li>`0`: manual renewal</li>
<li>`1`: auto-renewal</li>
Default value: 0.
  * `create_time` - Instance creation time.
  * `db_charset` - Instance character set, which currently supports only:
<li>UTF8</li>
<li>LATIN1</li>.
  * `db_engine_config` - Configuration information for the database engine, and the configuration format is as follows:.
{`$key1`:`$value1`, `$key2`:`$value2`}
Supported engines include:.
mssql_compatible engine:.
<li>migrationMode: specifies the database mode. optional parameter. valid values: single-db (single-database schema) and multi-db (multiple database schemas). defaults to single-db.</li>.
<li>defaultLocale: specifies the sorting area rule, an optional parameter that cannot be modified after initialization. default value is en_US. valid values include:.
`af_ZA`, `sq_AL`, `ar_DZ`, `ar_BH`, `ar_EG`, `ar_IQ`, `ar_JO`, `ar_KW`, `ar_LB`, `ar_LY`, `ar_MA`, `ar_OM`, `ar_QA`, `ar_SA`, `ar_SY`, `ar_TN`, `ar_AE`, `ar_YE`, `hy_AM`, `az_Cyrl_AZ`, `az_Latn_AZ`, `eu_ES`, `be_BY`, `bg_BG`, `ca_ES`, `zh_HK`, `zh_MO`, `zh_CN`, `zh_SG`, `zh_TW`, `hr_HR`, `cs_CZ`, `da_DK`, `nl_BE`, `nl_NL`, `en_AU`, `en_BZ`, `en_CA`, `en_IE`, `en_JM`, `en_NZ`, `en_PH`, `en_ZA`, `en_TT`, `en_GB`, `en_US`, `en_ZW`, `et_EE`, `fo_FO`, `fa_IR`, `fi_FI`, `fr_BE`, `fr_CA`, `fr_FR`, `fr_LU`, `fr_MC`, `fr_CH`, `mk_MK`, `ka_GE`, `de_AT`, `de_DE`, `de_LI`, `de_LU`, `de_CH`, `el_GR`, `gu_IN`, `he_IL`, `hi_IN`, `hu_HU`, `is_IS`, `id_ID`, `it_IT`, `it_CH`, `ja_JP`, `kn_IN`, `kok_IN`, `ko_KR`, `ky_KG`, `lv_LV`, `lt_LT`, `ms_BN`, `ms_MY`, `mr_IN`, `mn_MN`, `nb_NO`, `nn_NO`, `pl_PL`, `pt_BR`, `pt_PT`, `pa_IN`, `ro_RO`, `ru_RU`, `sa_IN`, `sr_Cyrl_RS`, `sr_Latn_RS`, `sk_SK`, `sl_SI`, `es_AR`, `es_BO`, `es_CL`, `es_CO`, `es_CR`, `es_DO`, `es_EC`, `es_SV`, `es_GT`, `es_HN`, `es_MX`, `es_NI`, `es_PA`, `es_PY`,`es_PE`, `es_PR`, `es_ES`, `es_TRADITIONAL`, `es_UY`, `es_VE`, `sw_KE`, `sv_FI`, `sv_SE`, `tt_RU`, `te_IN`, `th_TH`, `tr_TR`, `uk_UA`, `ur_IN`, `ur_PK`, `uz_Cyrl_UZ`, `uz_Latn_UZ`, `vi_VN`.</li>
<li>serverCollationName: Sorting rule name, an optional parameter, which cannot be modified after initialization, its default value is sql_latin1_general_cp1_ci_as, and its valid values include: `bbf_unicode_general_ci_as`, `bbf_unicode_cp1_ci_as`, `bbf_unicode_CP1250_ci_as`, `bbf_unicode_CP1251_ci_as`, `bbf_unicode_cp1253_ci_as`, `bbf_unicode_cp1254_ci_as`, `bbf_unicode_cp1255_ci_as`, `bbf_unicode_cp1256_ci_as`, `bbf_unicode_cp1257_ci_as`, `bbf_unicode_cp1258_ci_as`, `bbf_unicode_cp874_ci_as`, `sql_latin1_general_cp1250_ci_as`, `sql_latin1_general_cp1251_ci_as`, `sql_latin1_general_cp1_ci_as`, `sql_latin1_general_cp1253_ci_as`, `sql_latin1_general_cp1254_ci_as`, `sql_latin1_general_cp1255_ci_as`, `sql_latin1_general_cp1256_ci_as`, `sql_latin1_general_cp1257_ci_as`, `sql_latin1_general_cp1258_ci_as`, `chinese_prc_ci_as`, `cyrillic_general_ci_as`, `finnish_swedish_ci_as`, `french_ci_as`, `japanese_ci_as`, `korean_wansung_ci_as`, `latin1_general_ci_as`, `modern_spanish_ci_as`, `polish_ci_as`, `thai_ci_as`, `traditional_spanish_ci_as`, `turkish_ci_as`, `ukrainian_ci_as`, and `vietnamese_ci_as`.</li>.
  * `db_engine` - Database engine, which supports:
<li>`postgresql`: tencentdb for postgresql</li>.
<li>`mssql_compatible`: specifies mssql compatible - tencentdb for PostgreSQL.</li>.
Default value: `postgresql`.
  * `db_instance_class` - Purchasable specification ID.
  * `db_instance_cpu` - Number of assigned CPUs.
  * `db_instance_id` - Instance ID.
  * `db_instance_memory` - Assigned instance memory size in GB.
  * `db_instance_name` - Instance name.
  * `db_instance_net_info` - Instance network connection information.
    * `address` - DNS domain name.
    * `ip` - Ip.
    * `net_type` - Network type. 1: inner (private network address), 2: public (public network address).
    * `port` - Connection port address.
    * `protocol_type` - Specifies the protocol type to connect to the database. currently supported: postgresql, mssql (mssql compatible syntax).
    * `status` - Network connection status. Valid values: `initing` (never enabled before), `opened` (enabled), `closed` (disabled), `opening` (enabling), `closing` (disabling).
    * `subnet_id` - Subnet ID.
    * `vpc_id` - VPC ID. specifies the ID of the virtual private cloud.
  * `db_instance_status` - Instance status, including: `applying` (applying), `init` (to be initialized), `initing` (initializing), `running` (running), `limited run` (restricted operation), `isolating` (isolating), `isolated` (isolated), `disisolating` (de-isolating), `recycling` (recycling), `recycled` (recycled), `job running` (task executing), `offline` (offline), `migrating` (migrating), `expanding` (scaling out), `waitSwitch` (waiting to switch), `switching` (switching), `readonly` (readonly), `restarting` (restarting), `network changing` (network modification in progress), `upgrading` (kernel version upgrading), `audit-switching` (audit status changing), `primary-switching` (primary-secondary switching), `offlining` (offline), `deployment changing` (modify az), `cloning` (restoring data), `parameter modifying` (parameter modification in progress), `log-switching` (log status change), `restoring` (recovering), and `expanding` (scaling out).
  * `db_instance_storage` - Assigned instance storage capacity in GB.
  * `db_instance_type` - Instance type, which includes:
<li>primary: primary instance </li>
<li>readonly: read-only instance</li>
<li>guard: disaster recovery instance</li>
<li>temp: temporary instance</li>.
  * `db_instance_version` - Instance version. Valid value: `standard` (dual-server high-availability; one-primary-one-standby).
  * `db_kernel_version` - PostgreSQL kernel version, such as v12.7_r1.8. version information can be obtained from [DescribeDBVersions](https://www.tencentcloud.comom/document/api/409/89018?from_cn_redirect=1).
  * `db_major_version` - PostgreSQL major version number. specifies the version information that can be obtained from the [DescribeDBVersions](https://www.tencentcloud.comom/document/api/409/89018?from_cn_redirect=1) api. currently supports major versions 10, 11, 12, 13, 14, 15.
  * `db_node_set` - Instance node information
Note: This field may return null, indicating that no valid values can be obtained.
    * `role` - Node type. Valid values:
`Primary`;
`Standby`.
    * `zone` - AZ where the node resides, such as ap-guangzhou-1.
  * `db_version` - Number of the major PostgreSQL community version and minor version, such as 12.4, which can be queried by the [DescribeDBVersions](https://intl.cloud.tencent.com/document/api/409/89018?from_cn_redirect=1) API.
  * `deletion_protection` - Specifies whether to enable deletion protection for the instance. valid values as follows:.
-Specifies whether to enable deletion protection. valid values: true (enable deletion protection).
-Specifies whether to disable deletion protection. valid values: false (disable deletion protection).
  * `expanded_cpu` - Number of cpu cores that have been elastically scaled out.
  * `expire_time` - Instance expiration time.
  * `is_support_tde` - Whether the instance supports TDE data encryption.
<Li>0: not supported</li>.
<Li>1: supported.</li>.
Default value: 0

For TDE data encryption, see [overview of data transparent encryption](https://www.tencentcloud.comom/document/product/409/71748?from_cn_redirect=1).
  * `isolated_time` - Instance isolation time.
  * `master_db_instance_id` - Primary instance information. returned only when the instance is a read-only instance.
  * `network_access_list` - Network access list of the instance (this field has been deprecated)
Note: this field may return `null`, indicating that no valid values can be obtained.
    * `resource_id` - Network resource id, instance id, or RO group id.
    * `resource_type` - Resource type. valid values: 1 (instance), 2 (RO group).
    * `subnet_id` - Subnet ID.
    * `vip6` - IPv6 Address.
    * `vip` - IPv4 Address.
    * `vpc_id` - VPC ID. specifies the ID of the virtual private cloud.
    * `vpc_status` - Network status. valid values: 1-applying, 2-active, 3-deleting, 4-deleted.
    * `vport` - Specifies the access port.
  * `offline_time` - Decommissioning time.
  * `pay_type` - Billing mode:
<li>prepaid: monthly subscription, prepaid</li>
<li>postpaid: pay-as-you-go, postpaid</li>.
  * `project_id` - Project ID.
  * `read_only_instance_num` - Specifies the number of read-only instances.
  * `region` - Instance region such as ap-guangzhou, which corresponds to the`Region` field in `RegionSet`.
  * `root_user` - Instance root account name, default value is `root`.
  * `status_in_readonly_group` - Describes the state of the read-only instance in the read-only group.
  * `subnet_id` - VPC subnet ID, such as subnet-51lcif9y. effective VPC subnet ids can be obtained by logging in to the console or calling the api [DescribeSubnets](https://www.tencentcloud.comom/document/api/215/15784?from_cn_redirect=1) to acquire the unSubnetId field in the api return.
  * `support_ipv6` - Whether the instance supports IPv6:
<li>`0`: no</li>
<li>`1`: yes</li>
Default value: 0.
  * `tag_list` - Describes the Tag information associated with the instance.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `type` - Machine type.
  * `uid` - Instance `Uid`.
  * `update_time` - Last updated time of the instance attribute.
  * `vpc_id` - vpc ID, such as vpc-e6w23k31. a valid vpc ID can be obtained by logging in to the console to query or by calling the api [DescribeVpcs](https://www.tencentcloud.comom/document/api/215/15778?from_cn_redirect=1) and acquiring the unVpcId field in the api return.
  * `zone` - Instance AZ such as ap-guangzhou-3, which corresponds to the `Zone` field of `ZoneSet`.
* `instance_list` - (**Deprecated**) It has been deprecated from version 1.82.64. Please use `db_instance_set` instead. A list of postgresql instances. Each element contains the following attributes.
  * `auto_renew_flag` - Auto renew flag.
  * `availability_zone` - Availability zone.
  * `charge_type` - Pay type of the postgresql instance.
  * `charset` - Charset of the postgresql instance.
  * `create_time` - Create time of the postgresql instance.
  * `db_kernel_version` - PostgreSQL kernel version number.
  * `db_major_version` - PostgreSQL major version number.
  * `engine_version` - Version of the postgresql database engine.
  * `id` - ID of the postgresql instance.
  * `memory` - Memory size(in GB).
  * `name` - Name of the postgresql instance.
  * `private_access_ip` - IP address for private access.
  * `private_access_port` - Port for private access.
  * `project_id` - Project id, default value is 0.
  * `public_access_host` - Host for public access.
  * `public_access_port` - Port for public access.
  * `public_access_switch` - Indicates whether to enable the access to an instance from public network or not.
  * `root_user` - Instance root account name, default value is `root`.
  * `storage` - Volume size(in GB).
  * `subnet_id` - ID of subnet.
  * `tags` - The available tags within this postgresql.
  * `vpc_id` - ID of VPC.


