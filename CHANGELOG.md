## 1.37.0 (Unreleased)

FEATURES:
* **New Resource**: `tencentcloud_postgresql_instance`
* **New Data Source**: `tencentcloud_postgresql_instances`
* **New Data Source**: `tencentcloud_postgresql_speccodes`
* **New Data Source**: `tencentcloud_sqlserver_zone_config`

## 1.36.1 (June 12, 2020)

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` add new argument `labels`.
* Resource: `tencentcloud_kubernetes_as_scaling_group` add new argument `labels`.
* Resource: `tencentcloud_cos_bucket` add new arguments `encryption_algorithm` and `versioning_enable`.

## 1.36.0 (June 08, 2020)

FEATURES:

* **New Data Source**: `tencentcloud_availability_regions`

ENHANCEMENTS: 

* Data Source: `tencentcloud_redis_instances` add new argument `charge_type` to support prepaid type.
* Resource: `tencentcloud_redis_instance` add new argument `charge_type`, `prepaid_period` and `force_delete` to support prepaid type.
* Resource: `tencentcloud_mysql_instance` add new argument `force_delete` to support soft deletion.
* Resource: `tencentcloud_mysql_readonly_instance` add new argument `force_delete` to support soft deletion.

BUG FIXES:

* Resource: `tencentcloud_instance` fix `allocate_public_ip` inconsistency when eip is attached to the cvm.

DEPRECATED:
* Data Source: `tencentcloud_mysql_instances`: optional argument `pay_type` is no longer supported, replace by `charge_type`.
* Resource: `tencentcloud_mysql_instance`: optional arguments `pay_type` and `period` are no longer supported, replace by `charge_type` and `prepaid_period`.
* Resource: `tencentcloud_mysql_readonly_instance`: optional arguments `pay_type` and `period` are no longer supported, replace by `charge_type` and `prepaid_period`.
* Resource: `tencentcloud_tcaplus_group` replace by `tencentcloud_tcaplus_tablegroup`
* Data Source: `tencentcloud_tcaplus_groups` replace by `tencentcloud_tcaplus_tablegroups`
* Resource: `tencentcloud_tcaplus_tablegroup`,`tencentcloud_tcaplus_idl` and `tencentcloud_tcaplus_table`  arguments `group_id`/`group_name`  replace by `tablegroup_id`/`tablegroup_name`
* Data Source: `tencentcloud_tcaplus_groups`,`tencentcloud_tcaplus_idls` and `tencentcloud_tcaplus_tables` arguments `group_id`/`group_name`  replace by `tablegroup_id`/`tablegroup_name`

## 1.35.1 (June 02, 2020)

ENHANCEMENTS: 

* Resource: `tencentcloud_as_scaling_config`, `tencentcloud_eip` and `tencentcloud_kubernetes_cluster` remove the validate function of `internet_max_bandwidth_out`.
* Resource: `tencentcloud_vpn_gateway` update available value of `bandwidth`.

## 1.35.0 (June 01, 2020)

FEATURES:

* **New Data Source**: `tencentcloud_elasticsearch_instances`
* **New Resource**: `tencentcloud_elasticsearch_instance`

## 1.34.0 (May 28, 2020)

ENHANCEMENTS: 

* upgrade terraform-plugin-sdk

## 1.33.2 (May 25, 2020)

DEPRECATED:
* Data Source: `tencentcloud_tcaplus_applications` replace by `tencentcloud_tcaplus_clusters`,optional arguments `app_id` and `app_name` are no longer supported, replace by `cluster_id` and `cluster_name`
* Data Source: `tencentcloud_tcaplus_zones` replace by `tencentcloud_tcaplus_groups`,optional arguments `app_id`,`zone_id` and `zone_name` are no longer supported, replace by `cluster_id`,`group_id` and `cluster_name`
* Data Source: `tencentcloud_tcaplus_tables` optional arguments `app_id` and `zone_id` are no longer supported, replace by `cluster_id` and `group_id`
* Data Source: `tencentcloud_tcaplus_idls`: optional argument `app_id` is no longer supported, replace by `cluster_id`.
* Resource: `tencentcloud_tcaplus_application` replace by `tencentcloud_tcaplus_cluster`,input argument `app_name` is no longer supported, replace by `cluster_name`
* Resource: `tencentcloud_tcaplus_zone` replace by `tencentcloud_tcaplus_group`, input arguments `app_id` and `zone_name` are no longer supported, replace by `cluster_id` and `group_name`
* Resource: `tencentcloud_tcaplus_idl` input arguments `app_id` and `zone_id` are no longer supported, replace by `cluster_id` and `group_id`
* Resource: `tencentcloud_tcaplus_table` input arguments `app_id`and `zone_id` are no longer supported, replace by `cluster_id` and `group_id`
* Resource: `tencentcloud_redis_instance`: optional argument `type` is no longer supported, replace by `type_id`.
* Data Source: `tencentcloud_redis_instances`: output argument `type` is no longer supported, replace by `type_id`.
* Data Source: `tencentcloud_redis_zone_config`: output argument `type` is no longer supported, replace by `type_id`.

## 1.33.1 (May 22, 2020)

ENHANCEMENTS: 

* Data Source: `tencentcloud_redis_instances` add new argument `type_id`, `redis_shard_num`, `redis_replicas_num`
* Data Source: `tencentcloud_redis_zone_config` add output argument `type_id` and new output argument `type_id`, `redis_shard_nums`, `redis_replicas_nums`
* Data Source: `tencentcloud_ccn_instances` add new type `VPNGW` for field `instance_type`
* Data Source: `tencentcloud_vpn_gateways` add new type `CCN` for field `type`
* Resource: `tencentcloud_redis_instance` add new argument `type_id`, `redis_shard_num`, `redis_replicas_num`
* Resource: `tencentcloud_ccn_attachment` add new type `CNN_INSTANCE_TYPE_VPNGW` for field `instance_type`
* Resource: `tencentcloud_vpn_gateway` add new type `CCN` for field `type`

BUG FIXES:

* Resource: `tencentcloud_cdn_domain` fix `https_config` inconsistency after apply([#413](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/413)).

DEPRECATED:

* Resource: `tencentcloud_redis_instance`: optional argument `type` is no longer supported, replace by `type_id`.
* Data Source: `tencentcloud_redis_instances`: output argument `type` is no longer supported, replace by `type_id`.
* Data Source: `tencentcloud_redis_zone_config`: output argument `type` is no longer supported, replace by `type_id`.

## 1.33.0 (May 18, 2020)

FEATURES:

* **New Data Source**: `tencentcloud_monitor_policy_conditions`
* **New Data Source**: `tencentcloud_monitor_data`
* **New Data Source**: `tencentcloud_monitor_product_event`
* **New Data Source**: `tencentcloud_monitor_binding_objects`
* **New Data Source**: `tencentcloud_monitor_policy_groups`
* **New Data Source**: `tencentcloud_monitor_product_namespace`
* **New Resource**: `tencentcloud_monitor_policy_group`
* **New Resource**: `tencentcloud_monitor_binding_object`
* **New Resource**: `tencentcloud_monitor_binding_receiver`

ENHANCEMENTS: 

* Data Source: `tencentcloud_instances` add new output argument `instance_charge_type_prepaid_renew_flag`.
* Data Source: `tencentcloud_cbs_storages` add new output argument `prepaid_renew_flag`.
* Data Source: `tencentcloud_cbs_storages` add new output argument `charge_type`.
* Resource: `tencentcloud_instance` support update with argument `instance_charge_type_prepaid_renew_flag`.
* Resource: `tencentcloud_cbs_storage` add new argument `force_delete`.
* Resource: `tencentcloud_cbs_storage` add new argument `charge_type`.
* Resource: `tencentcloud_cbs_storage` add new argument `prepaid_renew_flag`.
* Resource: `tencentcloud_cdn_domain` add new argument `full_url_cache`([#405](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/405)).

DEPRECATED:

* Resource: `tencentcloud_cbs_storage`: optional argument `period` is no longer supported.

## 1.32.1 (April 30, 2020)

ENHANCEMENTS: 

* Resource: `tencentcloud_ccn_attachment` add new argument `ccn_uin`.
* Resource: `tencentcloud_instance` add new argument `force_delete`.

BUG FIXES:

* Resource: `tencentcloud_scf_function` fix update `zip_file`.

## 1.32.0 (April 20, 2020)

FEATURES:

* **New Resource**: `tencentcloud_kubernetes_cluster_attachment`([#285](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/285)).

ENHANCEMENTS: 

* Resource: `tencentcloud_cdn_domain` add new attribute `cname`([#395](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/395)).

BUG FIXES:

* Resource: `tencentcloud_cos_bucket_object` mark the object as destroyed when the object not exist.

## 1.31.2 (April 17, 2020)

ENHANCEMENTS: 

* Resource: `tencentcloud_cbs_storage` support modify `tags`.

## 1.31.1 (April 14, 2020)

BUG FIXES: 

* Resource: `tencentcloud_keypair` fix bug when trying to destroy resources containing CVM and key pair([#375](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/375)).
* Resource: `tencentcloud_clb_attachment` fix bug when trying to destroy multiple attachments in the array. 
* Resource: `tencentcloud_cam_group_membership` fix bug when trying to destroy multiple users in the array. 

ENHANCEMENTS:

* Resource: `tencentcloud_mysql_account` add new argument `host`([#372](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/372)).
* Resource: `tencentcloud_mysql_account_privilege` add new argument `account_host`([#372](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/372)).
* Resource: `tencentcloud_mysql_privilege` add new argument `account_host`([#372](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/372)).
* Resource: `tencentcloud_mysql_readonly_instance` check master monitor data before create([#379](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/379)).
* Resource: `tencentcloud_tcaplus_application` remove the pull password from server. 
* Resource: `tencentcloud_instance` support import `allocate_public_ip`([#382](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/382)).
* Resource: `tencentcloud_redis_instance` add two redis types.
* Data Source: `tencentcloud_vpc_instances` add new argument `cidr_block`,`tag_key` ([#378](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/378)).
* Data Source: `tencentcloud_vpc_route_tables` add new argument `tag_key`,`vpc_id`,`association_main` ([#378](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/378)).
* Data Source: `tencentcloud_vpc_subnets` add new argument `cidr_block`,`tag_key`,`is_remote_vpc_snat` ([#378](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/378)).
* Data Source: `tencentcloud_mysql_zone_config` and `tencentcloud_redis_zone_config` remove region check.

## 1.31.0 (April 07, 2020)

FEATURES:

* **New Resource**: `tencentcloud_cdn_domain`

ENHANCEMENTS:

* Data Source: `tencentcloud_cam_users` add new argument `user_id`.
* Resource: `tencentcloud_vpc` add retry logic.

BUG FIXES: 

* Resource: `tencentcloud_instance` fix timeout error when modify password.

## 1.30.7 (March 31, 2020)

BUG FIXES: 

* Resource: `tencentcloud_kubernetes_as_scaling_group` set a value to argument `key_ids` cause error .

## 1.30.6 (March 30, 2020)

ENHANCEMENTS:

* Resource: `tencentcloud_tcaplus_idl` add new argument `zone_id`. 
* Resource: `tencentcloud_cam_user` add new argument `force_delete`.([#354](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/354))
* Data Source: `tencentcloud_vpc_subnets` add new argument `vpc_id`. 

## 1.30.5 (March 19, 2020)

BUG FIXES: 

* Resource: `tencentcloud_key_pair` will be replaced when `public_key` contains comment.
* Resource: `tencentcloud_scf_function` upload local file error.

ENHANCEMENTS:

* Resource: `tencentcloud_scf_function` runtime support nodejs8.9 and nodejs10.15. 

## 1.30.4 (March 10, 2020)

BUG FIXES:

* Resource: `tencentcloud_cam_policy` fix read nil issue when the resource is not exist.([#344](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/#344)).
* Resource: `tencentcloud_key_pair` will be replaced when the end of `public_key` contains spaces([#343](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/343)).
* Resource: `tencentcloud_scf_function` fix trigger does not support cos_region.

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` add new attributes `cluster_os_type`,`cluster_internet`,`cluster_intranet`,`managed_cluster_internet_security_policies` and `cluster_intranet_subnet_id`.


## 1.30.3 (February 24, 2020)

BUG FIXES:

* Resource: `tencentcloud_instance` fix that classic network does not support([#339](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/339)).

## 1.30.2 (February 17, 2020)

ENHANCEMENTS:

* Data Source: `tencentcloud_cam_policies` add new attribute `policy_id`.
* Data Source: `tencentcloud_cam_groups` add new attribute `group_id`.

## 1.30.1 (January 21, 2020)

BUG FIXES:

* Resource: `tencentcloud_dnat` fix `elastic_port` and `internal_port` type error.
* Resource: `tencentcloud_vpn_gateway` fix `state` type error.
* Resource: `tencentcloud_dayu_ddos_policy` fix that `white_ips` and `black_ips` can not be updated.
* Resource: `tencentcloud_dayu_l4_rule` fix that rule parameters can not be updated.

ENHANCEMENTS:

* Data Source: `tencentcloud_key_pairs` support regular expression search by name.

## 1.30.0 (January 14, 2020)

FEATURES:

* **New Data Source**: `tencentcloud_dayu_cc_http_policies`
* **New Data Source**: `tencentcloud_dayu_cc_https_policies`
* **New Data Source**: `tencentcloud_dayu_ddos_policies`
* **New Data Source**: `tencentcloud_dayu_ddos_policy_attachments`
* **New Data Source**: `tencentcloud_dayu_ddos_policy_cases`
* **New Data Source**: `tencentcloud_dayu_l4_rules`
* **New Data Source**: `tencentcloud_dayu_l7_rules`
* **New Resource**: `tencentcloud_dayu_cc_http_policy`
* **New Resource**: `tencentcloud_dayu_cc_https_policy`
* **New Resource**: `tencentcloud_dayu_ddos_policy`
* **New Resource**: `tencentcloud_dayu_ddos_policy_attachment`
* **New Resource**: `tencentcloud_dayu_ddos_policy_case`
* **New Resource**: `tencentcloud_dayu_l4_rule`
* **New Resource**: `tencentcloud_dayu_l7_rule`

BUG FIXES:

* gaap: optimize gaap describe: when describe resource by id but get more than 1 resources, return error directly instead of using the first result 
* Resource: `tencentcloud_eni_attachment` fix detach may failed.
* Resource: `tencentcloud_instance` remove the tag that be added by as attachment automatically([#300](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/300)).
* Resource: `tencentcloud_clb_listener` fix `sni_switch` type error([#297](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/297)).
* Resource: `tencentcloud_vpn_gateway` shows argument `prepaid_renew_flag` has changed when applied again([#298](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/298)).
* Resource: `tencentcloud_clb_instance` fix the bug that instance id is not set in state file([#303](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/303)).
* Resource: `tencentcloud_vpn_gateway` that is postpaid charge type cannot be deleted normally([#312](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/312)).
* Resource: `tencentcloud_vpn_gateway` add `InternalError` SDK error to triggering the retry process.
* Resource: `tencentcloud_vpn_gateway` fix read nil issue when the resource is not exist.
* Resource: `tencentcloud_clb_listener_rule` fix unclear error message of SSL type error.
* Resource: `tencentcloud_ha_vip_attachment` fix read nil issue when the resource is not exist.
* Data Source: `tencentcloud_security_group` fix `project_id` type error.
* Data Source: `tencentcloud_security_groups` fix `project_id` filter not works([#303](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/314)).

## 1.29.0 (January 06, 2020)

FEATURES:

* **New Data Source**: `tencentcloud_gaap_domain_error_pages`
* **New Resource**: `tencentcloud_gaap_domain_error_page`

ENHANCEMENTS:
* Data Source: `tencentcloud_vpc_instances` add new optional argument `is_default`.
* Data Source: `tencentcloud_vpc_subnets` add new optional argument `availability_zone`,`is_default`.

BUG FIXES:
* Resource: `tencentcloud_redis_instance` field security_groups are id list, not name list([#291](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/291)).

## 1.28.0 (December 25, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_cbs_snapshot_policies`
* **New Resource**: `tencentcloud_cbs_snapshot_policy_attachment`

ENHANCEMENTS:

* doc: rewrite website index
* Resource: `tencentcloud_instance` support modifying instance type([#251](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/251)).
* Resource: `tencentcloud_gaap_http_domain` add new optional argument `realserver_certificate_ids`.
* Data Source: `tencentcloud_gaap_http_domains` add new output argument `realserver_certificate_ids`.

DEPRECATED:

* Resource: `tencentcloud_gaap_http_domain`: optional argument `realserver_certificate_id` is no longer supported.
* Data Source: `tencentcloud_gaap_http_domains`: output argument `realserver_certificate_id` is no longer supported.

## 1.27.0 (December 17, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_tcaplus_applications`
* **New Data Source**: `tencentcloud_tcaplus_zones`
* **New Data Source**: `tencentcloud_tcaplus_tables`
* **New Data Source**: `tencentcloud_tcaplus_idls`
* **New Resource**: `tencentcloud_tcaplus_application`
* **New Resource**: `tencentcloud_tcaplus_zone`
* **New Resource**: `tencentcloud_tcaplus_idl`
* **New Resource**: `tencentcloud_tcaplus_table`

ENHANCEMENTS:

* Resource: `tencentcloud_mongodb_instance` support more instance type([#241](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/241)).
* Resource: `tencentcloud_kubernetes_cluster` support more instance type([#237](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/237)).

BUG FIXES:

* Fix bug that resource `tencentcloud_instance` delete error when instance launch failed.
* Fix bug that resource `tencentcloud_security_group` read error when response is InternalError.
* Fix bug that the type of `cluster_type` is wrong in data source `tencentcloud_mongodb_instances`([#242](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/242)).
* Fix bug that resource `tencentcloud_eip` unattach error([#233](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/233)).
* Fix bug that terraform read nil attachment resource when the attached resource of attachment resource is removed of resource CLB and CAM.
* Fix doc example error of resource `tencentcloud_nat_gateway`.

DEPRECATED:

* Resource: `tencentcloud_eip`: optional argument `applicable_for_clb` is no longer supported.

## 1.26.0 (December 09, 2019)

FEATURES:

* **New Resource**: `tencentcloud_mysql_privilege`([#223](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/223)).
* **New Resource**: `tencentcloud_kubernetes_as_scaling_group`([#202](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/202)).

ENHANCEMENTS:

* Resource: `tencentcloud_gaap_layer4_listener` support import.
* Resource: `tencentcloud_gaap_http_rule` support import.
* Resource: `tencentcloud_gaap_security_rule` support import.
* Resource: `tencentcloud_gaap_http_domain` add new optional argument `client_certificate_ids`.
* Resource: `tencentcloud_gaap_layer7_listener` add new optional argument `client_certificate_ids`.
* Data Source: `tencentcloud_gaap_http_domains` add new output argument `client_certificate_ids`.
* Data Source: `tencentcloud_gaap_layer7_listeners` add new output argument `client_certificate_ids`.

DEPRECATED:

* Resource: `tencentcloud_gaap_http_domain`: optional argument `client_certificate_id` is no longer supported.
* Resource: `tencentcloud_gaap_layer7_listener`: optional argument `client_certificate_id` is no longer supported.
* Resource: `tencentcloud_mysql_account_privilege` replaced by `tencentcloud_mysql_privilege`.
* Data Source: `tencentcloud_gaap_http_domains`: output argument `client_certificate_id` is no longer supported.
* Data Source: `tencentcloud_gaap_layer7_listeners`: output argument `client_certificate_id` is no longer supported.

BUG FIXES:

* Fix bug that resource `tencentcloud_clb_listener` 's unchangeable `health_check_switch`([#235](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/235)).
* Fix bug that resource `tencentcloud_clb_instance` read nil and report error.
* Fix example errors of resource `tencentcloud_cbs_snapshot_policy` and data source `tencentcloud_dnats`.

## 1.25.2 (December 04, 2019)

BUG FIXES:
* Fixed bug that the validator of cvm instance type is incorrect.

## 1.25.1 (December 03, 2019)

ENHANCEMENTS:
* Optimized error message of validators.

BUG FIXES:
* Fixed bug that the type of `state` is incorrect in data source `tencentcloud_nat_gateways`([#226](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/226)).
* Fixed bug that the value of `cluster_max_pod_num` is incorrect in resource `tencentcloud_kubernetes_cluster`([#228](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/228)).


## 1.25.0 (December 02, 2019)

ENHANCEMENTS:

* Resource: `tencentcloud_instance` support `SPOTPAID` instance. Thanks to @LipingMao ([#209](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/209)).
* Resource: `tencentcloud_vpn_gateway` add argument `prepaid_renew_flag` and `prepaid_period` to support prepaid VPN gateway instance creation.

BUG FIXES:
* Fixed bugs that update operations on `tencentcloud_cam_policy` do not work.
* Fixed bugs that filters on `tencentcloud_cam_users` do not work.

DEPRECATED:
 * Data Source: `tencentcloud_cam_user_policy_attachments`:`policy_type` is no longer supported.
 * Data Source: `tencentcloud_cam_group_policy_attachments`:`policy_type` is no longer supported.

## 1.24.1 (November 26, 2019)

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` add support for `PREPAID` instance type. Thanks to @woodylic ([#204](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/204)).
* Resource: `tencentcloud_cos_bucket` add optional argument tags
* Data Source: `tencentcloud_cos_buckets` add optional argument tags

BUG FIXES:
* Fixed docs issues of `tencentcloud_nat_gateway`

## 1.24.0 (November 20, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_ha_vips`
* **New Data Source**: `tencentcloud_ha_vip_eip_attachments`
* **New Resource**: `tencentcloud_ha_vip`
* **New Resource**: `tencentcloud_ha_vip_eip_attachment`

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` cluster_os add new support: `centos7.6x86_64` and `ubuntu18.04.1 LTSx86_64` 
* Resource: `tencentcloud_nat_gateway` add computed argument `created_time`.

BUG FIXES:

* Fixed docs issues of CAM, DNAT and NAT_GATEWAY
* Fixed query issue that paged-query was not supported in data source `tencentcloud_dnats`
* Fixed query issue that filter `address_ip` was set incorrectly in data source `tencentcloud_eips`

## 1.23.0 (November 14, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_images`
* **New Data Source**: `tencentcloud_vpn_gateways`
* **New Data Source**: `tencentcloud_customer_gateways`
* **New Data Source**: `tencentcloud_vpn_connections`
* **New Resource**: `tencentcloud_vpn_gateway`
* **New Resource**: `tencentcloud_customer_gateway`
* **New Resource**: `tencentcloud_vpn_connection`
* **Provider TencentCloud**: add `security_token` argument

ENHANCEMENTS:

* All api calls now using api3.0
* Resource: `tencentcloud_eip` add optional argument `tags`.
* Data Source: `tencentcloud_eips` add optional argument `tags`.

BUG FIXES:

* Fixed docs of CAM

## 1.22.0 (November 05, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_cfs_file_systems`
* **New Data Source**: `tencentcloud_cfs_access_groups`
* **New Data Source**: `tencentcloud_cfs_access_rules`
* **New Data Source**: `tencentcloud_scf_functions`
* **New Data Source**: `tencentcloud_scf_namespaces`
* **New Data Source**: `tencentcloud_scf_logs`
* **New Resource**: `tencentcloud_cfs_file_system`
* **New Resource**: `tencentcloud_cfs_access_group`
* **New Resource**: `tencentcloud_cfs_access_rule`
* **New Resource**: `tencentcloud_scf_function`
* **New Resource**: `tencentcloud_scf_namespace`

## 1.21.2 (October 29, 2019)

BUG FIXES:

* Resource: `tencentcloud_gaap_realserver` add ip/domain exists check
* Resource: `tencentcloud_kubernetes_cluster` add error handling logic and optional argument `tags`.
* Resource: `tencentcloud_kubernetes_scale_worker` add error handling logic.
* Data Source: `tencentcloud_kubernetes_clusters` add optional argument `tags`.

## 1.21.1 (October 23, 2019)

ENHANCEMENTS:

* Updated golang to version 1.13.x

BUG FIXES:

* Fixed docs of CAM

## 1.21.0 (October 15, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_cam_users`
* **New Data Source**: `tencentcloud_cam_groups`
* **New Data Source**: `tencentcloud_cam_policies`
* **New Data Source**: `tencentcloud_cam_roles`
* **New Data Source**: `tencentcloud_cam_user_policy_attachments`
* **New Data Source**: `tencentcloud_cam_group_policy_attachments`
* **New Data Source**: `tencentcloud_cam_role_policy_attachments`
* **New Data Source**: `tencentcloud_cam_group_memberships`
* **New Data Source**: `tencentcloud_cam_saml_providers`
* **New Data Source**: `tencentcloud_reserved_instance_configs`
* **New Data Source**: `tencentcloud_reserved_instances`
* **New Resource**: `tencentcloud_cam_user`
* **New Resource**: `tencentcloud_cam_group`
* **New Resource**: `tencentcloud_cam_role`
* **New Resource**: `tencentcloud_cam_policy`
* **New Resource**: `tencentcloud_cam_user_policy_attachment`
* **New Resource**: `tencentcloud_cam_group_policy_attachment`
* **New Resource**: `tencentcloud_cam_role_policy_attachment`
* **New Resource**: `tencentcloud_cam_group_membership`
* **New Resource**: `tencentcloud_cam_saml_provider`
* **New Resource**: `tencentcloud_reserved_instance`

ENHANCEMENTS:

* Resource: `tencentcloud_gaap_http_domain` support import
* Resource: `tencentcloud_gaap_layer7_listener` support import

BUG FIXES:

* Resource: `tencentcloud_gaap_http_domain` fix sometimes can't enable realserver auth

## 1.20.1 (October 08, 2019)

ENHANCEMENTS:

* Data Source: `tencentcloud_availability_zones` refactor logic with api3.0 .
* Data Source: `tencentcloud_as_scaling_groups` add optional argument `tags` and attribute `tags` for `scaling_group_list`.
* Resource: `tencentcloud_eip` add optional argument `type`, `anycast_zone`, `internet_service_provider`, etc.
* Resource: `tencentcloud_as_scaling_group` add optional argument `tags`.

BUG FIXES:

* Data Source: `tencentcloud_gaap_http_domains` set response `certificate_id`, `client_certificate_id`, `realserver_auth`, `basic_auth` and `gaap_auth` default value when they are nil.
* Resource: `tencentcloud_gaap_http_domain` set response `certificate_id`, `client_certificate_id`, `realserver_auth`, `basic_auth` and `gaap_auth` default value when they are nil.

## 1.20.0 (September 24, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_eips`
* **New Data Source**: `tencentcloud_instances`
* **New Data Source**: `tencentcloud_key_pairs`
* **New Data Source**: `tencentcloud_placement_groups`
* **New Resource**: `tencentcloud_placement_group`

ENHANCEMENTS:

* Data Source: `tencentcloud_redis_instances` add optional argument `tags`.
* Data Source: `tencentcloud_mongodb_instances` add optional argument `tags`.
* Data Source: `tencentcloud_instance_types` add optional argument `availability_zone` and `gpu_core_count`.
* Data Source: `tencentcloud_gaap_http_rules` add optional argument `forward_host` and attributes `forward_host` in `rules`.
* Resource: `tencentcloud_redis_instance` add optional argument `tags`.
* Resource: `tencentcloud_mongodb_instance` add optional argument `tags`.
* Resource: `tencentcloud_mongodb_sharding_instance` add optional argument `tags`.
* Resource: `tencentcloud_instance` add optional argument `placement_group_id`.
* Resource: `tencentcloud_eip` refactor logic with api3.0 .
* Resource: `tencentcloud_eip_association` refactor logic with api3.0 .
* Resource: `tencentcloud_key_pair` refactor logic with api3.0 .
* Resource: `tencentcloud_gaap_http_rule` add optional argument `forward_host`.

BUG FIXES:
* Resource: `tencentcloud_mysql_instance`: miss argument `availability_zone` causes the instance to be recreated.

DEPRECATED:

* Data Source: `tencentcloud_eip` replaced by `tencentcloud_eips`.

## 1.19.0 (September 19, 2019)

FEATURES:

* **New Resource**: `tencentcloud_security_group_lite_rule`.

ENHANCEMENTS:

* Data Source: `tencentcloud_security_groups`: add optional argument `tags`.
* Data Source: `tencentcloud_security_groups`: add optional argument `result_output_file` and new attributes `ingress`, `egress` for list `security_groups`.
* Resource: `tencentcloud_security_group`: add optional argument `tags`.
* Resource: `tencentcloud_as_scaling_config`: internet charge type support `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.

BUG FIXES:
* Resource: `tencentcloud_clb_listener_rule`: fix unclear description and errors in example.
* Resource: `tencentcloud_instance`: fix hostname is not work.

## 1.18.1 (September 17, 2019)

FEATURES:

* **Update Data Source**: `tencentcloud_vpc_instances` add optional argument `tags`
* **Update Data Source**: `tencentcloud_vpc_subnets` add optional argument `tags`
* **Update Data Source**: `tencentcloud_route_tables` add optional argument `tags`
* **Update Resource**: `tencentcloud_vpc` add optional argument `tags`
* **Update Resource**: `tencentcloud_subnet` add optional argument `tags`
* **Update Resource**: `tencentcloud_route_table` add optional argument `tags`

ENHANCEMENTS:

* Data Source:`tencentcloud_kubernetes_clusters`  support pull out authentication information for cluster access too.
* Resource:`tencentcloud_kubernetes_cluster`  support pull out authentication information for cluster access.

BUG FIXES:

* Resource: `tencentcloud_mysql_instance`: when the mysql is abnormal state, read the basic information report error

DEPRECATED:

* Data Source: `tencentcloud_kubernetes_clusters`:`container_runtime` is no longer supported. 

## 1.18.0 (September 10, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_ssl_certificates`
* **New Data Source**: `tencentcloud_dnats`
* **New Data Source**: `tencentcloud_nat_gateways`
* **New Resource**: `tencentcloud_ssl_certificate`
* **Update Resource**: `tencentcloud_clb_redirection` add optional argument `is_auto_rewrite`
* **Update Resource**: `tencentcloud_nat_gateway` , add more configurable items.
* **Update Resource**: `tencentcloud_nat` , add more configurable items.

DEPRECATED:
* Data Source: `tencentcloud_nats` replaced by `tencentcloud_nat_gateways`.

## 1.17.0 (September 04, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_gaap_proxies`
* **New Data Source**: `tencentcloud_gaap_realservers`
* **New Data Source**: `tencentcloud_gaap_layer4_listeners`
* **New Data Source**: `tencentcloud_gaap_layer7_listeners`
* **New Data Source**: `tencentcloud_gaap_http_domains`
* **New Data Source**: `tencentcloud_gaap_http_rules`
* **New Data Source**: `tencentcloud_gaap_security_policies`
* **New Data Source**: `tencentcloud_gaap_security_rules`
* **New Data Source**: `tencentcloud_gaap_certificates`
* **New Resource**: `tencentcloud_gaap_proxy`
* **New Resource**: `tencentcloud_gaap_realserver`
* **New Resource**: `tencentcloud_gaap_layer4_listener`
* **New Resource**: `tencentcloud_gaap_layer7_listener`
* **New Resource**: `tencentcloud_gaap_http_domain`
* **New Resource**: `tencentcloud_gaap_http_rule`
* **New Resource**: `tencentcloud_gaap_certificate`
* **New Resource**: `tencentcloud_gaap_security_policy`
* **New Resource**: `tencentcloud_gaap_security_rule`

## 1.16.3 (August 30, 2019)

BUG FIXIES:

* Resource: `tencentcloud_kubernetes_cluster`: cgi error retry.
* Resource: `tencentcloud_kubernetes_scale_worker`: cgi error retry.

## 1.16.2 (August 28, 2019)

BUG FIXIES:

* Resource: `tencentcloud_instance`: fixed cvm data disks missing computed.
* Resource: `tencentcloud_mysql_backup_policy`: `backup_model` remove logical backup support. 
* Resource: `tencentcloud_mysql_instance`: `tags` adapt to the new official api.

## 1.16.1 (August 27, 2019)

ENHANCEMENTS:
* `tencentcloud_instance`: refactor logic with api3.0 .

## 1.16.0 (August 20, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_kubernetes_clusters`
* **New Resource**: `tencentcloud_kubernetes_scale_worker`
* **New Resource**: `tencentcloud_kubernetes_cluster`

DEPRECATED:
* Data Source: `tencentcloud_container_clusters` replaced by `tencentcloud_kubernetes_clusters`.
* Data Source: `tencentcloud_container_cluster_instances` replaced by `tencentcloud_kubernetes_clusters`.
* Resource: `tencentcloud_container_cluster` replaced by `tencentcloud_kubernetes_cluster`.
* Resource: `tencentcloud_container_cluster_instance` replaced by `tencentcloud_kubernetes_scale_worker`.

## 1.15.2 (August 14, 2019)

ENHANCEMENTS:

* `tencentcloud_as_scaling_group`: fixed issue that binding scaling group to load balancer does not work.
* `tencentcloud_clb_attachements`: rename `rewrite_source_rule_id` with `source_rule_id` and rename `rewrite_target_rule_id` with `target_rule_id`.

## 1.15.1 (August 13, 2019)

ENHANCEMENTS:

* `tencentcloud_instance`: changed `image_id` property to ForceNew ([#78](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/78))
* `tencentcloud_instance`: improved with retry ([#82](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/82))
* `tencentcloud_cbs_storages`: improved with retry ([#82](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/82))
* `tencentcloud_clb_instance`: bug fixed and improved with retry ([#37](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/37))

## 1.15.0 (August 07, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_clb_instances`
* **New Data Source**: `tencentcloud_clb_listeners`
* **New Data Source**: `tencentcloud_clb_listener_rules`
* **New Data Source**: `tencentcloud_clb_attachments`
* **New Data Source**: `tencentcloud_clb_redirections`
* **New Resource**: `tencentcloud_clb_instance`
* **New Resource**: `tencentcloud_clb_listener`
* **New Resource**: `tencentcloud_clb_listener_rule`
* **New Resource**: `tencentcloud_clb_attachment`
* **New Resource**: `tencentcloud_clb_redirection`

DEPRECATED:
* Resource: `tencentcloud_lb` replaced by `tencentcloud_clb_instance`.
* Resource: `tencentcloud_alb_server_attachment` replaced by `tencentcloud_clb_attachment`.

## 1.14.1 (August 05, 2019)

BUG FIXIES:

* resource/tencentcloud_security_group_rule: fixed security group rule id is not compatible with previous version.

## 1.14.0 (July 30, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_security_groups`
* **New Data Source**: `tencentcloud_mongodb_instances`
* **New Data Source**: `tencentcloud_mongodb_zone_config`
* **New Resource**: `tencentcloud_mongodb_instance`
* **New Resource**: `tencentcloud_mongodb_sharding_instance`
* **Update Resource**: `tencentcloud_security_group_rule` add optional argument `description`

DEPRECATED:
* Data Source: `tencnetcloud_security_group` replaced by `tencentcloud_security_groups`

ENHANCEMENTS:
* Refactoring security_group logic with api3.0

## 1.13.0 (July 23, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_dc_gateway_instances`
* **New Data Source**: `tencentcloud_dc_gateway_ccn_routes`
* **New Resource**: `tencentcloud_dc_gateway`
* **New Resource**: `tencentcloud_dc_gateway_ccn_route`

## 1.12.0 (July 16, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_dc_instances`
* **New Data Source**: `tencentcloud_dcx_instances`
* **New Resource**: `tencentcloud_dcx`
* **UPDATE Resource**:`tencentcloud_mysql_instance` and `tencentcloud_mysql_readonly_instance` completely delete instance. 

BUG FIXIES:

* resource/tencentcloud_instance: fixed issue when data disks set as delete_with_instance not works.
* resource/tencentcloud_instance: if managed public_ip manually, please don't define `allocate_public_ip` ([#62](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/62)).
* resource/tencentcloud_eip_association: fixed issue when instances were manually deleted ([#60](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/60)).
* resource/tencentcloud_mysql_readonly_instance:remove an unsupported property `gtid`

## 1.11.0 (July 02, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_ccn_instances`
* **New Data Source**: `tencentcloud_ccn_bandwidth_limits`
* **New Resource**: `tencentcloud_ccn`
* **New Resource**: `tencentcloud_ccn_attachment`
* **New Resource**: `tencentcloud_ccn_bandwidth_limit`

## 1.10.0 (June 27, 2019)

ENHANCEMENTS:

* Refactoring vpc logic with api3.0
* Refactoring cbs logic with api3.0

FEATURES:
* **New Data Source**: `tencentcloud_vpc_instances`
* **New Data Source**: `tencentcloud_vpc_subnets`
* **New Data Source**: `tencentcloud_vpc_route_tables`
* **New Data Source**: `tencentcloud_cbs_storages`
* **New Data Source**: `tencentcloud_cbs_snapshots`
* **New Resource**: `tencentcloud_route_table_entry`
* **New Resource**: `tencentcloud_cbs_snapshot_policy`
* **Update Resource**: `tencentcloud_vpc` , add more configurable items.
* **Update Resource**: `tencentcloud_subnet` , add more configurable items.
* **Update Resource**: `tencentcloud_route_table`, add more configurable items.
* **Update Resource**: `tencentcloud_cbs_storage`, add more configurable items.
* **Update Resource**: `tencentcloud_instance`: add optional argument `tags`.
* **Update Resource**: `tencentcloud_security_group_rule`: add optional argument `source_sgid`.
 
DEPRECATED:
* Data Source: `tencentcloud_vpc` replaced by `tencentcloud_vpc_instances`.
* Data Source: `tencentcloud_subnet` replaced by  `tencentcloud_vpc_subnets`.
* Data Source: `tencentcloud_route_table` replaced by `tencentcloud_vpc_route_tables`.
* Resource: `tencentcloud_route_entry` replaced by `tencentcloud_route_table_entry`.

## 1.9.1 (June 24, 2019)

BUG FIXIES:

* data/tencentcloud_instance: fixed vpc ip is in use error when re-creating with private ip ([#46](https://github.com/terraform-providers/terraform-provider-tencentcloud/issues/46)).

## 1.9.0 (June 18, 2019)

ENHANCEMENTS:

* update to `v0.12.1` Terraform SDK version

BUG FIXIES:

* data/tencentcloud_security_group: `project_id` remote API return sometime is string type.
* resource/tencentcloud_security_group: just like `data/tencentcloud_security_group`

## 1.8.0 (June 11, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_as_scaling_configs`
* **New Data Source**: `tencentcloud_as_scaling_groups`
* **New Data Source**: `tencentcloud_as_scaling_policies`
* **New Resource**: `tencentcloud_as_scaling_config`
* **New Resource**: `tencentcloud_as_scaling_group`
* **New Resource**: `tencentcloud_as_attachment`
* **New Resource**: `tencentcloud_as_scaling_policy`
* **New Resource**: `tencentcloud_as_schedule`
* **New Resource**: `tencentcloud_as_lifecycle_hook`
* **New Resource**: `tencentcloud_as_notification`

## 1.7.0 (May 23, 2019)

FEATURES:
* **New Data Source**: `tencentcloud_redis_zone_config`
* **New Data Source**: `tencentcloud_redis_instances`
* **New Resource**: `tencentcloud_redis_instance`
* **New Resource**: `tencentcloud_redis_backup_config`

ENHANCEMENTS:

* resource/tencentcloud_instance: Add `hostname`, `project_id`, `delete_with_instance` argument.
* Update tencentcloud-sdk-go to better support redis api.

## 1.6.0 (May 15, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_cos_buckets`
* **New Data Source**: `tencentcloud_cos_bucket_object`
* **New Resource**: `tencentcloud_cos_bucket`
* **New Resource**: `tencentcloud_cos_bucket_object`

ENHANCEMENTS:

* Add the framework of auto generating terraform docs

## 1.5.0 (April 26, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_mysql_backup_list`
* **New Data Source**: `tencentcloud_mysql_zone_config`
* **New Data Source**: `tencentcloud_mysql_parameter_list`
* **New Data Source**: `tencentcloud_mysql_instance`
* **New Resource**: `tencentcloud_mysql_backup_policy`
* **New Resource**: `tencentcloud_mysql_account`
* **New Resource**: `tencentcloud_mysql_account_privilege`
* **New Resource**: `tencentcloud_mysql_instance`
* **New Resource**: `tencentcloud_mysql_readonly_instance`

ENHANCEMENTS:

* resource/tencentcloud_subnet: `route_table_id` now is an optional argument

## 1.4.0 (April 12, 2019)

ENHANCEMENTS:

* data/tencentcloud_image: add `image_name` attribute to this data source.
* resource/tencentcloud_instance: data disk count limit now is upgrade from 1 to 10, as API has supported more disks.
* resource/tencentcloud_instance: PREPAID instance now can be deleted, but still have some limit in API.

BUG FIXIES:

* resource/tencentcloud_instance: `allocate_public_ip` doesn't work properly when it is set to false.

## 1.3.0 (March 12, 2019)

FEATURES:

* **New Resource**: `tencentcloud_lb` ([#3](https://github.com/terraform-providers/terraform-provider-scaffolding/issues/3))

ENHANCEMENTS:

* resource/tencentcloud_instance: Add `user_data_raw` argument ([#4](https://github.com/terraform-providers/terraform-provider-scaffolding/issues/4))

## 1.2.2 (September 28, 2018)

BUG FIXES:

* resource/tencentcloud_cbs_storage: make name to be required ([#25](https://github.com/tencentyun/terraform-provider-tencentcloud/issues/25))
* resource/tencentcloud_instance: support user data and private ip

## 1.2.0 (April 3, 2018)

FEATURES:

* **New Resource**: `tencentcloud_container_cluster`
* **New Resource**: `tencentcloud_container_cluster_instance`
* **New Data Source**: `tencentcloud_container_clusters`
* **New Data Source**: `tencentcloud_container_cluster_instances`

## 1.1.0 (March 9, 2018)

FEATURES:

* **New Resource**: `tencentcloud_eip`
* **New Resource**: `tencentcloud_eip_association`
* **New Data Source**: `tencentcloud_eip`
* **New Resource**: `tencentcloud_nat_gateway`
* **New Resource**: `tencentcloud_dnat`
* **New Data Source**: `tencentcloud_nats`
* **New Resource**: `tencentcloud_cbs_snapshot`
* **New Resource**: `tencentcloud_alb_server_attachment`

## 1.0.0 (January 19, 2018)

FEATURES:

### CVM

RESOURCES:

* instance create
* instance read
* instance update
    * reset instance
    * reset password
    * update instance name
    * update security groups
* instance delete
* key pair create
* key pair read
* key pair delete

DATA SOURCES:

* image read
* instance\_type read
* zone read

### VPC

RESOURCES:

* vpc create
* vpc read
* vpc update (update name)
* vpc delete
* subnet create
* subnet read
* subnet update (update name)
* subnet delete
* security group create
* security group read
* security group update (update name, description)
* security group delete
* security group rule create
* security group rule read
* security group rule delete
* route table create
* route table read
* route table update (update name)
* route table delete
* route entry create
* route entry read
* route entry delete

DATA SOURCES:

* vpc read
* subnet read
* security group read
* route table read

### CBS

RESOURCES:

* storage create
* storage read
* storage update (update name)
* storage attach
* storage detach
