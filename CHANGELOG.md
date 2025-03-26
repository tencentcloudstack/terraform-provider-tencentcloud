## 1.81.176(March 24 , 2025)

ENHANCEMENTS:

* resource/tencentcloud_cos_bucket_domain_certificate_attachment: support `cert_id` params ([#3235](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3235))

## 1.81.175(March 21 , 2025)

FEATURES:

* **New Data Source:** `tencentcloud_cdwpg_instances` ([#3232](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3232))
* **New Data Source:** `tencentcloud_cdwpg_log` ([#3232](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3232))
* **New Data Source:** `tencentcloud_cdwpg_nodes` ([#3232](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3232))
* **New Data Source:** `tencentcloud_sqlserver_collation_time_zone` ([#3228](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3228))
* **New Resource:** `tencentcloud_route_table_entry_config` ([#3217](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3217))

ENHANCEMENTS:

* resource/tencentcloud_cls_cos_shipper: support new params ([#3231](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3231))
* resource/tencentcloud_dts_migrate_service: Optimize deletion logic ([#3230](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3230))
* resource/tencentcloud_dts_sync_job: Optimize deletion logic ([#3230](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3230))
* resource/tencentcloud_dts_sync_job: Optimize update logic ([#3230](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3230))
* resource/tencentcloud_emr_cluster: handle error when the element in pre_executed_file_settings is nil ([#3219](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3219))
* resource/tencentcloud_instance: support `user_data` and `user_data_raw` update ([#3227](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3227))
* resource/tencentcloud_kubernetes_cluster_endpoint: add `kube_config`, `kube_config_intranet` fields output ([#3225](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3225))
* resource/tencentcloud_mongodb_instance_transparent_data_encryption: fix import ([#3222](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3222))
* resource/tencentcloud_placement_group: support `affinity` and `tags` ([#3229](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3229))
* resource/tencentcloud_postgresql_instance: support update `charge_type`, `period`, `auto_renew_flag` ([#3223](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3223))
* resource/tencentcloud_reserved_instance: update cvm sdk ([#3227](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3227))
* resource/tencentcloud_serverless_hbase_instance: Change `tags` from `TypeList` to `TypeSet` ([#3233](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3233))
* resource/tencentcloud_sqlserver_instance: Support time_zone field. ([#3228](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3228))
* resource/tencentcloud_sqlserver_readonly_instance: Support time_zone field. ([#3228](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3228))
* resource/tencentcloud_tcr_customized_domain: update doc ([#3224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3224))

## 1.81.174(March 14 , 2025)

FEATURES:

* **New Resource:** `tencentcloud_nat_gateway_flow_monitor` ([#3215](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3215))
* **New Resource:** `tencentcloud_sqlserver_wan_ip_config` ([#3198](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3198))

ENHANCEMENTS:

* resource/tencentcloud_clb_listener_rule: support `health_source_ip_type` param ([#3208](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3208))
* resource/tencentcloud_cos_bucket: update doc ([#3214](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3214))
* resource/tencentcloud_kubernetes_cluster_endpoint: optimize code logic and documentation ([#3207](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3207))
* resource/tencentcloud_sqlserver_basic_instance: add new params ([#3198](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3198))
* resource/tencentcloud_sqlserver_general_cloud_instance: add new params ([#3198](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3198))
* resource/tencentcloud_sqlserver_instance: add new params ([#3198](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3198))
* resource/tencentcloud_ssl_certificates: fix crash while no return `CertificatePublicKey` ([#3216](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3216))

## 1.81.173(March 13 , 2025)

FEATURES:

* **New Resource:** `tencentcloud_mongodb_readonly_instance` ([#3200](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3200))

ENHANCEMENTS:

* resource/tencentcloud_cynosdb_cluster: optimize ro group sg ([#3203](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3203))
* resource/tencentcloud_elasticsearch_instance: Support for public_access fields. ([#3205](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3205))
* resource/tencentcloud_kubernetes_addon: update doc ([#3204](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3204))
* resource/tencentcloud_mongodb_sharding_instance: support update mongos_memory ([#3199](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3199))
* resource/tencentcloud_postgresql_instance: support upgrade version ([#3201](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3201))

## 1.81.172(March 7 , 2025)

FEATURES:

* **New Resource:** `tencentcloud_as_load_balancer` ([#3186](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3186))
* **New Resource:** `tencentcloud_teo_l7_acc_rule` ([#3185](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3185))

ENHANCEMENTS:

* data-source/tencentcloud_vpc_instances: support `common_assistant_cidr` and `container_assistant_cidr` ([#3194](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3194))
* resource/tencentcloud_cbs_storage: fix prepaid_period update ([#3187](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3187))
* resource/tencentcloud_elasticsearch_instance: Modify kibana_public_access settings when creating ([#3191](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3191))
* resource/tencentcloud_emr_cluster: support scene_name ([#3189](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3189))
* resource/tencentcloud_mysql_instance: update `device_type` description ([#3193](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3193))
* resource/tencentcloud_serverless_hbase_instance: support update instance_name and tags ([#3188](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3188))
* resource/tencentcloud_teo_l7_acc_rule: Add retry query ([#3192](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3192))
* resource/tencentcloud_teo_l7_acc_setting: Add retry query ([#3196](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3196))

## 1.81.171(March 6 , 2025)

ENHANCEMENTS:

* resource/tencentcloud_cfs_access_rule: retry the FailedOperation.PgroupIsUpdating error ([#3134](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3134))
* resource/tencentcloud_redis_instance: Restore the replica_zone_id logic ([#3184](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3184))

## 1.81.170(March 3 , 2025)

FEATURES:

* **New Resource:** `tencentcloud_cynosdb_backup_config` ([#3181](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3181))
* **New Resource:** `tencentcloud_dasb_asset_sync_job_operation` ([#3171](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3171))
* **New Resource:** `tencentcloud_postgresql_instance_ssl_config` ([#3183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3183))
* **New Resource:** `tencentcloud_teo_l4_proxy_rule` ([#3149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3149))
* **New Resource:** `tencentcloud_teo_l7_acc_setting` ([#3169](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3169))

ENHANCEMENTS:

* data-source/tencentcloud_dnspod_record_list: add `sub_domains` param ([#3165](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3165))
* resource/tencentcloud_accept_join_share_unit_invitation_operation: optimize unique ID ([#3174](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3174))
* resource/tencentcloud_cfs_auto_snapshot_policy: update doc ([#3176](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3176))
* resource/tencentcloud_cos_bucket: optimize `acl_body` verification logic ([#3177](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3177))
* resource/tencentcloud_dc_gateway_ccn_routes: add new params ([#3178](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3178))
* resource/tencentcloud_dc_internet_address_config: update doc ([#3179](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3179))
* resource/tencentcloud_dc_share_dcx_config: update doc ([#3180](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3180))
* resource/tencentcloud_mysql_readonly_instance: modify doc. ([#3172](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3172))
* resource/tencentcloud_postgresql_instance: Support window period switching. ([#3167](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3167))
* resource/tencentcloud_postgresql_readonly_instance: Support window period switching. ([#3167](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3167))
* resource/tencentcloud_redis_instance: Optimize the change issue of the `replica_zone_ids` field in a single availability zone. ([#3160](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3160))
* resource/tencentcloud_reject_join_share_unit_invitation_operation: optimize unique ID ([#3174](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3174))
* resource/tencentcloud_scf_function: update doc ([#3182](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3182))
* resource/tencentcloud_tse_cngw_service: modify field properties ([#3173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3173))
* resource/tencentcloud_vpc_ipv6_eni_address: Offline this document ([#3175](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3175))
* resource/tencentcloud_waf_custom_rule: update doc ([#3166](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3166))
* resource/tencentcloud_waf_custom_white_rule: update doc ([#3166](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3166))

## 1.81.169 (February 27 , 2025)

FEATURES:

* **New Data Source:** `tencentcloud_emr_job_status_detail` ([#3144](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3144))
* **New Data Source:** `tencentcloud_emr_service_node_infos` ([#3144](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3144))
* **New Resource:** `tencentcloud_clb_customized_config_attachment` ([#3163](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3163))
* **New Resource:** `tencentcloud_clb_customized_config_v2` ([#3162](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3162))
* **New Resource:** `tencentcloud_emr_auto_scale_strategy` ([#3144](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3144))
* **New Resource:** `tencentcloud_emr_deploy_yarn_operation` ([#3144](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3144))
* **New Resource:** `tencentcloud_emr_yarn` ([#3144](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3144))
* **New Resource:** `tencentcloud_postgresql_instance_network_access` ([#3004](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3004))
* **New Resource:** `tencentcloud_postgresql_parameters` ([#3146](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3146))

ENHANCEMENTS:

* provider: OIDC auth support set `provider_id` ([#3152](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3152))
* resource/tencentcloud_as_scaling_config: add `disaster_recover_group_ids` ([#3147](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3147))
* resource/tencentcloud_cfw_edge_policy: Update resource fields and code logic ([#3142](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3142))
* resource/tencentcloud_cfw_nat_policy: Update resource fields and code logic ([#3142](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3142))
* resource/tencentcloud_cos_bucket_policy: update doc ([#3155](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3155))
* resource/tencentcloud_dcx_extra_config: update field properties and doc ([#3143](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3143))
* resource/tencentcloud_eip_public_address_adjust: update doc ([#3154](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3154))
* resource/tencentcloud_tag: optimize the creation retry logic ([#3156](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3156))
* resource/tencentcloud_tdmq_rabbitmq_vip_instance: Add `pay_mode`, `cluster_version`, `public_access_endpoint` and `vpcs` ([#3158](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3158))
* resource/tencentcloud_tdmq_rabbitmq_vip_instance: update doc ([#3145](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3145))
* resource/tencentcloud_teo_rule_engine: update doc ([#3148](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3148))
* resource/tencentcloud_vpc_bandwidth_package_attachment: Mark `resource_type` as `Computed` ([#3159](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3159))
* resource/tencentcloud_waf_anti_fake: update doc ([#3161](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3161))
* resource/tencentcloud_waf_anti_info_leak: update doc ([#3161](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3161))
* resource/tencentcloud_waf_cc: add new params and update doc ([#3161](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3161))
* resource/tencentcloud_waf_cc_auto_status: update doc ([#3161](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3161))
* resource/tencentcloud_waf_custom_rule: add new params and update doc ([#3161](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3161))
* resource/tencentcloud_waf_custom_white_rule: add new params and update doc ([#3161](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3161))
* resource/waf: Update WAF resource query logic ([#3141](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3141))

## 1.81.168 (February 18 , 2025)

FEATURES:

* **New Data Source:** `tencentcloud_serverless_hbase_instances` ([#3140](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3140))
* **New Resource:** `tencentcloud_mqtt_instance` ([#3135](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3135))
* **New Resource:** `tencentcloud_mqtt_instance_public_endpoint` ([#3135](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3135))
* **New Resource:** `tencentcloud_mqtt_topic` ([#3135](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3135))

ENHANCEMENTS:

* resource/tencentcloud_clb_instance: update doc ([#3137](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3137))
* resource/tencentcloud_clb_target_group_instance_attachment: Query module update returns weight field ([#3138](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3138))
* resource/tencentcloud_clb_target_group_instance_attachment: operation increases status query ([#3139](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3139))
* resource/tencentcloud_invite_organization_member_operation: update doc ([#3136](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3136))
* resource/tencentcloud_serverless_hbase_instance: fixed a delay in status updates ([#3140](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3140))

## 1.81.167 (February 14 , 2025)

FEATURES:

* **New Data Source:** `tencentcloud_mongodb_instance_urls` ([#3124](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3124))

ENHANCEMENTS:

* data-source/tencentcloud_eb_bus: fix `connection_briefs` and `target_briefs` write to state failed ([#3105](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3105))
* datasource/tencentcloud_instance_types: support `cbs_filter` and `cbs_configs` ([#3118](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3118))
* datasource/tencentcloud_ssl_certificates: update doc ([#3121](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3121))
* resource/tencentcloud_as_scaling_config: fix the issue of changing password ([#3122](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3122))
* resource/tencentcloud_as_scaling_config: update `instance_types` to limit the number of input parameters ([#3123](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3123))
* resource/tencentcloud_cam_saml_provider: update doc ([#3120](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3120))
* resource/tencentcloud_cbs_disk_backup: update resource creation error message ([#3129](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3129))
* resource/tencentcloud_cbs_storage: fix `disk_backup_quota` update failed ([#3096](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3096))
* resource/tencentcloud_ckafka_topic: fix the issue with default values for params ([#3107](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3107))
* resource/tencentcloud_clb_security_group_attachment: update doc ([#3119](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3119))
* resource/tencentcloud_dc_gateway: update code ([#3091](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3091))
* resource/tencentcloud_dc_instance: update field properties and doc ([#3130](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3130))
* resource/tencentcloud_eb_put_events: update doc ([#3125](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3125))
* resource/tencentcloud_eip: fix of creating `BANDWIDTH_PREPAID_BY_MONTH` eip ([#3111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3111))
* resource/tencentcloud_vpc: fix the issue where field `assistant_cidrs` cannot be modified correctly ([#3102](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3102))
* resource/tencentcloud_vpc_flow_log: support CCN type ([#3106](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3106))
* resource/tencentcloud_waf_cc: Update doc ([#3127](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3127))
* resource/tencentcloud_waf_ip_access_control_v2: Fix the issue of unable to query properly ([#3127](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3127))

## 1.81.166 (February 11 , 2025)

FEATURES:

* **New Resource:** `tencentcloud_monitor_tmp_multiple_writes_list` ([#3115](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3115))

ENHANCEMENTS:

* datasource/tencentcloud_security_groups: update doc ([#3117](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3117))
* resource/tencentcloud_clb_attachment: fix the diff in the `targets` field ([#3110](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3110))
* resource/tencentcloud_instance: `cam_role_name` support update ([#3114](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3114))
* resource/tencentcloud_kubernetes_charts: update doc ([#3116](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3116))
* resource/tencentcloud_monitor_tmp_multiple_writes: resource will be deprecated in version v1.81.166 ([#3115](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3115))

## 1.81.165 (February 7 , 2025)

ENHANCEMENTS:

* resource/tencentcloud_cdn_domain: `header_mode` support `set` ([#3113](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3113))
* resource/tencentcloud_monitor_tmp_instance: set creation wait ([#3109](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3109))

## 1.81.164 (January 24 , 2025)

ENHANCEMENTS:

* resource/tencentcloud_clb_listener: support update `listener_name` ([#3103](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3103))
* resource/tencentcloud_cos_bucket: support routing_rules ([#3108](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3108))

## 1.81.163 (January 22 , 2025)

FEATURES:

* **New Data Source:** `tencentcloud_cdc_dedicated_clusters` ([#3099](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3099))

ENHANCEMENTS:

* datasource/tencentcloud_ha_vips: fix the issue of failed filtering queries ([#3098](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3098))
* resource/tencentcloud_cam_group_membership: update doc ([#3093](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3093))
* resource/tencentcloud_dcx: resource creation increases waiting time ([#3092](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3092))
* resource/tencentcloud_identity_center_user: params user_name unsupport change ([#3094](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3094))
* resource/tencentcloud_instance: update the length limit of `instance_name` ([#3101](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3101))
* resource/tencentcloud_monitor_tmp_exporter_integration: update SDK information ([#3100](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3100))
* resource/tencentcloud_organization_org_member_email: fix the issue of update function execution failure ([#3095](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3095))
* resource/tencentcloud_postgresql_instance: update query interface limit count ([#3084](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3084))
* resource/tencentcloud_tat_command: update code and doc ([#3083](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3083))

## 1.81.162 (January 21 , 2025)

ENHANCEMENTS:

* resource/tencentcloud_clb_customized_config: update code ([#3085](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3085))
* resource/tencentcloud_clb_listener_rule: support `multi_cert_info` ([#3082](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3082))
* resource/tencentcloud_clb_target_group_attachment: update code ([#3090](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3090))
* resource/tencentcloud_cos_bucket: update website param ([#3088](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3088))
* resource/tencentcloud_instance: support local disks ([#3089](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3089))
* resource/tencentcloud_nat_gateway: remove max limit of eips ([#3097](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3097))

## 1.81.161 (January 16 , 2025)

ENHANCEMENTS:

* resource/tencentcloud_clb_listener: support `multi_cert_info` ([#3081](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3081))
* resource/tencentcloud_clb_listener_rule: update code ([#3080](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3080))
* resource/tencentcloud_vpn_customer_gateway: support `bgp_asn` ([#3073](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3073))

## 1.81.160 (January 14 , 2025)

FEATURES:

* **New Resource:** `tencentcloud_cls_web_callback` ([#3068](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3068))

ENHANCEMENTS:

* datasource/tencentcloud_cbs_snapshots: update code ([#3078](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3078))
* resource/tencentcloud_clb_listener: support `snat_enable` ([#3076](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3076))
* resource/tencentcloud_clb_listener_rule: params `domains` support modify ([#3079](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3079))
* resource/tencentcloud_cls_logset: update code ([#3074](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3074))
* resource/tencentcloud_identity_center_role_assignment: process failed task status ([#3064](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3064))
* resource/tencentcloud_instance: create support idempotent ([#3075](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3075))
* resource/tencentcloud_instance: retry reading cbs ([#3072](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3072))
* resource/tencentcloud_route_table_entry: `next_type` support new type ([#3077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3077))

## 1.81.159 (January 9 , 2025)

ENHANCEMENTS:

* resource/tencentcloud_as_start_instance_refresh: support `fail_process` ([#3071](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3071))
* resource/tencentcloud_cbs_snapshot: update code and support `tags` ([#3070](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3070))
* resource/tencentcloud_clb_target_group_instance_attachment: update function support retry ([#3069](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3069))
* resource/tencentcloud_kubernetes_cluster: update code ([#3066](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3066))

## 1.81.158 (January 8 , 2025)

FEATURES:

* **New Resource:** `tencentcloud_mongodb_instance_params` ([#3055](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3055))
* **New Resource:** `tencentcloud_private_dns_extend_end_point` ([#3057](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3057))

ENHANCEMENTS:

* resource/tencentcloud_ckafka_instance: support param `elastic_bandwidth_switch` ([#3051](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3051))
* resource/tencentcloud_cos_bucket: update `website` params ([#3061](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3061))
* resource/tencentcloud_emr_cluster: update read ([#3054](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3054))

## 1.81.157 (January 3 , 2025)

ENHANCEMENTS:

* resource/tencentcloud_images: update code ([#3052](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3052))
* resource/tencentcloud_sqlserver_readonly_instance: support import ([#3050](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3050))
* resource/tencentcloud_vpc_end_point: update `end_point_vip` params ([#3053](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3053))

## 1.81.156 (December 31 , 2024)

ENHANCEMENTS:

* resource/tencentcloud_as_scaling_group: support `replace_mode`, `desired_capacity_sync_with_max_min_size` params ([#3048](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3048))
* resource/tencentcloud_cam_role: support update `session_duration` ([#3049](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3049))
* resource/tencentcloud_kms_key: support check `pending_delete_window_in_days` params ([#3047](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3047))

## 1.81.155 (December 27 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_role_configuration_provisionings` ([#3046](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3046))
* **New Resource:** `tencentcloud_provision_role_configuration_operation` ([#3046](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3046))

ENHANCEMENTS:

* provider: STS auth support retry ([#3044](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3044))
* resource/tencentcloud_clb_listener_rule: support `health_check_port` ([#3045](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3045))
* resource/tencentcloud_eni_attachment: support retry ([#3043](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3043))
* resource/tencentcloud_vpn_gateway: support `bgp_asn` ([#3041](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3041))

## 1.81.154 (December 24 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_monitor_tmp_multiple_writes` ([#3036](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3036))

ENHANCEMENTS:

* data-source/tencentcloud_dnspod_record_list: add param `instance_list` ([#3035](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3035))
* provider: update provider auth ([#3039](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3039))
* resource/tencentcloud_clb_listener_default_domain: support retry ([#3038](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3038))
* resource/tencentcloud_clb_listener_rule: support retry ([#3038](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3038))
* resource/tencentcloud_kubernetes_node_pool: support CBM for `instance_type` ([#3034](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3034))

## 1.81.153 (December 24 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_monitor_tmp_multiple_writes` ([#3036](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3036))

ENHANCEMENTS:

* data-source/tencentcloud_dnspod_record_list: add param `instance_list` ([#3035](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3035))
* provider: update provider auth ([#3039](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3039))
* resource/tencentcloud_clb_listener_default_domain: support retry ([#3038](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3038))
* resource/tencentcloud_clb_listener_rule: support retry ([#3038](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3038))
* resource/tencentcloud_kubernetes_node_pool: support CBM for `instance_type` ([#3034](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3034))

## 1.81.152 (December 20 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_organization_org_share_unit_members` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))
* **New Data Source:** `tencentcloud_organization_org_share_unit_resources` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))
* **New Data Source:** `tencentcloud_organization_org_share_units` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))
* **New Resource:** `tencentcloud_accept_join_share_unit_invitation_operation` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))
* **New Resource:** `tencentcloud_organization_org_share_unit_resource` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))
* **New Resource:** `tencentcloud_reject_join_share_unit_invitation_operation` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))

ENHANCEMENTS:

* provider: support `allowed_account_ids` and `forbidden_account_ids` ([#3030](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3030))
* resource/tencentcloud_as_scaling_policy: add nil check ([#3032](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3032))
* resource/tencentcloud_clb_listener_rule: fix error with `domains` ([#3028](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3028))
* resource/tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment: fix code panic ([#3031](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3031))


## 1.81.151 (December 20 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_organization_org_share_unit_members` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))
* **New Data Source:** `tencentcloud_organization_org_share_unit_resources` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))
* **New Data Source:** `tencentcloud_organization_org_share_units` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))
* **New Resource:** `tencentcloud_accept_join_share_unit_invitation_operation` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))
* **New Resource:** `tencentcloud_organization_org_share_unit_resource` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))
* **New Resource:** `tencentcloud_reject_join_share_unit_invitation_operation` ([#3029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3029))

ENHANCEMENTS:

* provider: support `allowed_account_ids` and `forbidden_account_ids` ([#3030](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3030))
* resource/tencentcloud_as_scaling_policy: add nil check ([#3032](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3032))
* resource/tencentcloud_clb_listener_rule: fix error with `domains` ([#3028](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3028))
* resource/tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment: fix code panic ([#3031](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3031))

## 1.81.150 (December 18 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_serverless_hbase_instance` ([#3002](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3002))

ENHANCEMENTS:

* resource/tencentcloud_instance: support `disk_name` ([#3024](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2996))
* resource/tencentcloud_as_lifecycle_hook: update create function ([#3024](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3024))
* resource/tencentcloud_cam_role_permission_boundary_attachment: fix code panic ([#3026](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3026))
* resource/tencentcloud_dasb_bind_device_resource: support `read`, `update` and `delete` ([#3022](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3022))
* resource/tencentcloud_mongodb_instance: remove RequiredWith on `availability_zone_list` and `hidden_zone` ([#3025](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3025))
* resource/tencentcloud_mongodb_sharding_instance: remove RequiredWith on `availability_zone_list` and `hidden_zone` ([#3025](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3025))

DEPRECATED:

* resource: `tencentcloud_lite_hbase_instance` replaced by `tencentcloud_serverless_hbase_instance`

## 1.81.149 (December 16 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_classic_elastic_public_ipv6s` ([#3021](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3021))
* **New Data Source:** `tencentcloud_elastic_public_ipv6s` ([#3021](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3021))
* **New Resource:** `tencentcloud_classic_elastic_public_ipv6` ([#3021](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3021))
* **New Resource:** `tencentcloud_elastic_public_ipv6` ([#3021](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3021))
* **New Resource:** `tencentcloud_elastic_public_ipv6_attachment` ([#3021](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3021))

ENHANCEMENTS:

* resource/tencentcloud_cbs_storage: update `tag` ([#3020](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3020))
* resource/tencentcloud_instance: update unit test ([#3019](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3019))

## 1.81.148 (December 13 , 2024)

ENHANCEMENTS:

* resource/tencentcloud_clb_instance: spport `eip_address_id` ([#3009](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3009))
* resource/tencentcloud_cynosdb_cluster: support `instance_init_infos` ([#3015](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3015))
* resource/tencentcloud_kubernetes_scale_worker: support `tags` ([#3010](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3010))
* resource/tencentcloud_monitor_tmp_tke_cluster_agent: support `open_default_record` field. ([#3018](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3018))
* resource/tencentcloud_teo_rule_engine: Add enumeration to the `target` field. ([#3008](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3008))

## 1.81.147 (December 11 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_subdomain_validate_status` ([#3005](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3005))
* **New Resource:** `tencentcloud_subdomain_validate_txt_value_operation` ([#3005](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3005))

## 1.81.146 (December 9 , 2024)

ENHANCEMENTS:

* resource/tencentcloud_cls_alarm_notice: add new params ([#3000](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/3000))
* resource/tencentcloud_kubernetes_node_pool: spport `timeouts` for create and update ([#2998](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2998))

## 1.81.145 (December 3 , 2024)

ENHANCEMENTS:

* datasource/tencentcloud_cam_role_detail: Update doc ([#2987](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2987))
* datasource/tencentcloud_cam_sub_accounts: Update doc ([#2987](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2987))
* resource/tencentcloud_cam_role: Add computed attribute `role_arn` ([#2990](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2990))
* resource/tencentcloud_ccn: Modify the description of `qos` ([#2989](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2989))
* resource/tencentcloud_cls_alarm: support `multi_conditions` ([#2985](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2985))
* resource/tencentcloud_identity_center_group: `group_type` support optional ([#2992](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2992))
* resource/tencentcloud_identity_center_scim_credential: support param `credential_secret` ([#2994](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2994))
* resource/tencentcloud_kubernetes_node_pool: add `auto_update_instance_tags` params ([#2991](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2991))
* resource/tencentcloud_kubernetes_node_pool: update doc and code ([#2986](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2986))
* resource/tencentcloud_mysql_backup_policy: `backup_model` support `snapshot` ([#2993](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2993))
* resource/tencentcloud_tdmq_rabbitmq_vip_instance: Update doc ([#2995](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2995))

## 1.81.144 (December 1 , 2024)

ENHANCEMENTS:

* resource/tencentcloud_as_scaling_group: The asg ignores the order of `forward_balancer_ids` ([#2984](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2984))
* resource/tencentcloud_as_start_instance_refresh: Add `timeouts` for refresh action ([#2984](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2984))

## 1.81.143 (November 29 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_cls_cloud_product_log_task` ([#2976](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2976))
* **New Resource:** `tencentcloud_cls_notice_content` ([#2981](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2981))
* **New Resource:** `tencentcloud_scf_custom_domain` ([#2983](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2983))

ENHANCEMENTS:

* resource/tencentcloud_as_scaling_config: support cdc ([#2980](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2980))
* resource/tencentcloud_clb_listener_rule: support param `oauth` ([#2978](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2978))
* resource/tencentcloud_kubernetes_node_pool: add `wait_node_ready`, `scale_tolerance` params ([#2979](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2979))
* resource/tencentcloud_vpn_connection: add `negotiation_type`, `bgp_config`, `health_check_config` params ([#2982](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2982))

## 1.81.142 (November 23 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_reserve_ip_address` ([#2972](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2972))

## 1.81.141 (November 22 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_cam_role_detail` ([#2970](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2970))
* **New Data Source:** `tencentcloud_cam_sub_accounts` ([#2970](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2970))

ENHANCEMENTS:

* datasource/tencentcloud_kubernetes_cluster_common_names: Update code ([#2970](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2970))
* resource/tencentcloud_cynosdb_cluster: suppor `password` modification ([#2968](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2968))
* resource/tencentcloud_kubernetes_node_pool: Fix the issue where `node_os` cannot be modified ([#2964](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2964))
* resource/tencentcloud_mysql_readonly_instance: Supports creating read-only instances across regions. ([#2959](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2959))
* resource/tencentcloud_redis_instance: Processing the change from a single-AZ instance to a multi-AZ instance. ([#2961](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2961))
* resource/tencentcloud_vpc_flow_log: support modifying `tags` ([#2963](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2963))

## 1.81.140 (November 15 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment` ([#2962](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2962))
* **New Resource:** `tencentcloud_waf_ip_access_control_v2` ([#2946](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2946))

ENHANCEMENTS:

* provider: Adapt to data type for `tccli` ([#2944](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2944))
* resource/tencentcloud_clb_instance: update doc ([#2956](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2956))
* resource/tencentcloud_clb_listener_rule: support set `GRPC` and `GRPCS` ([#2952](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2952))
* resource/tencentcloud_dnspod_record: support refresh state while destroy outside of Terraform ([#2957](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2957))
* resource/tencentcloud_emr_cluster: support param pre_executed_file_settings ([#2960](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2960))
* resource/tencentcloud_kubernetes_auth_attachment: Update the logic of the read function ([#2953](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2953))
* resource/tencentcloud_kubernetes_cluster: support set `instance_delete_mode` ([#2949](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2949))
* resource/tencentcloud_kubernetes_cluster: update `exist_instance` params ([#2958](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2958))
* resource/tencentcloud_kubernetes_native_node_pool: support field `machine_type`. ([#2951](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2951))
* resource/tencentcloud_nat_gateway: add `stock_public_ip_addresses_bandwidth_out` fields ([#2955](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2955))
* resource/tencentcloud_scf_function: Update doc ([#2945](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2945))

## 1.81.139 (November 8 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_identity_center_scim_credential` ([#2950](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2950))
* **New Resource:** `tencentcloud_identity_center_scim_credential_status` ([#2950](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2950))
* **New Resource:** `tencentcloud_identity_center_scim_synchronization_status` ([#2950](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2950))
* **New Resource:** `tencentcloud_private_dns_end_point` ([#2948](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2948))
* **New Resource:** `tencentcloud_private_dns_forward_rule` ([#2948](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2948))

## 1.81.138 (November 6 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_events_audit_track` ([#2931](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2931))
* **New Resource:** `tencentcloud_kubernetes_cluster_master_attachment` ([#2926](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2926))
* **New Resource:** `tencentcloud_tcss_image_registry` ([#2935](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2935))
* **New Resource:** `tencentcloud_vpc_notify_routes` ([#2941](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2941))

ENHANCEMENTS:

* resource/tencentcloud_as_lifecycle_hook: add `lifecycle_transition_type` ([#2943](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2943))
* resource/tencentcloud_as_start_instance_refresh: add `max_surge` ([#2942](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2942))
* resource/tencentcloud_clb_listener: support `h2c_switch` params ([#2939](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2939))
* resource/tencentcloud_kubernetes_node_pool: compatible `NodePoolQueryFailed` error when deleted node pool ([#2940](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2940))
* resource/tencentcloud_scf_function: field `triggers.type` suuport `http` and `cls` ([#2938](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2938))
* resource/tencentcloud_ses_email_address: support set `smtp_password` ([#2937](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2937))
* resource/tencentcloud_ssl_check_certificate_domain_verification_operation: Add retry and customize the maximum timeout ([#2936](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2936))

## 1.81.137 (November 1 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_cvm_action_timer` ([#2874](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2874))
* **New Resource:** `tencentcloud_subscribe_private_zone_service` ([#2929](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2929))

ENHANCEMENTS:

* datasource/tencentcloud_ccn_route_table_input_policies: Update unit test ([#2915](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2915))
* datasource/tencentcloud_gaap_http_domains: add param `is_default_server` ([#2927](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2927))
* datasource/tencentcloud_gaap_layer7_listeners: add params `group_id`, `tls_support_versions` and `tls_ciphers` ([#2927](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2927))
* datasource/tencentcloud_gaap_proxy_detail: add param `is_support_tls_choice` ([#2927](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2927))
* datasource/tencentcloud_gaap_proxy_statistics: update `metric_names` ([#2927](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2927))
* resource/tencentcloud_audit_track: add `storage_account_id`, `storage_app_id` fields ([#2930](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2930))
* resource/tencentcloud_cam_role: update the verification rules for field `document` ([#2917](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2917))
* resource/tencentcloud_emr_cluster: Support multi_disks. ([#2919](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2919))
* resource/tencentcloud_gaap_http_domain: add params `group_id` and `is_default_server` ([#2927](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2927))
* resource/tencentcloud_gaap_layer7_listener: add params `group_id`, `tls_support_versions` and `tls_ciphers` ([#2927](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2927))
* resource/tencentcloud_kubernetes_cluster: support `cluster_os` field to update and adjust its set logic ([#2918](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2918))
* resource/tencentcloud_kubernetes_cluster: support field `resource_delete_options` for deleting CBS ([#2916](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2916))
* resource/tencentcloud_postgresql_clone_db_instance: add `dedicated_cluster_id` params ([#2920](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2920))
* resource/tencentcloud_postgresql_instance: update field properties for `backup_plan` ([#2921](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2921))
* resource/tencentcloud_sqlserver_instance: deprecated field `ha_type` ([#2923](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2923))
* resource/tencentcloud_sqlserver_instance: update field `multi_zones` ([#2922](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2922))

## 1.81.136 (October 25 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_ccn_route_table_input_policies` ([#2910](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2910))

ENHANCEMENTS:

* resource/tencentcloud_ccn_attachment: fix attachment with ccn_uin ([#2913](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2913))
* resource/tencentcloud_kubernetes_node_pool: support `annotations` field ([#2909](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2909))
* resource/tencentcloud_kubernetes_scale_worker: Update read function ([#2908](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2908))
* resource/tencentcloud_private_dns_zone: Support unbind all VPCs bound to the current private dns zone. ([#2914](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2914))

BUG FIXES:

* resource/tencentcloud_cynosdb_cluster_databases: Fix the problem of matching errors caused by fuzzy query. ([#2912](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2912))

## 1.81.135 (October 23 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_organization_nodes` ([#2906](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2906))
* **New Resource:** `tencentcloud_open_identity_center_operation` ([#2906](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2906))

ENHANCEMENTS:

* provider/Update auth ([#2887](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2887))
* resource/kubernetes_native_node_pool: adjust the `labels` field to Set type ([#2900](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2900))
* resource/tencentcloud_cam_policy: Update error message ([#2902](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2902))
* resource/tencentcloud_ccn_attachment: fix attachment with ccn_uin ([#2901](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2901))
* resource/tencentcloud_emr_cluster: support emr horizontal expansion ([#2904](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2904))
* resource/tencentcloud_kubernetes_auth_attachment: increase the timeout of the read retry ([#2900](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2900))
* resource/tencentcloud_kubernetes_cluster: adjust TKE tags logic to support existing native nodes ([#2896](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2896))
* resource/tencentcloud_tcr_tag_retention_rule: adjust modify logic ([#2905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2905))
* resource/tencentcloud_waf_modify_access_period:  This resource has been deprecated in Terraform TencentCloud provider version 1.81.135. ([#2903](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2903))

## 1.81.134 (October 18 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_identity_center_groups` ([#2894](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2894))
* **New Data Source:** `tencentcloud_identity_center_role_configurations` ([#2894](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2894))
* **New Data Source:** `tencentcloud_identity_center_users` ([#2894](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2894))
* **New Resource:** `tencentcloud_clb_listener_default_domain` ([#2899](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2899))
* **New Resource:** `tencentcloud_sg_rule` ([#2886](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2886))

ENHANCEMENTS:

* datasource/tencentcloud_enis: add `ipv6s` ([#2898](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2898))
* provider: add new auth for `cam_role_name` ([#2892](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2892))
* resource/tencentcloud_cos_bucket: `website` add new `redirect_all_requests_to` params ([#2897](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2897))
* resource/tencentcloud_kubernetes_scale_worker: Fix import function ([#2891](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2891))
* resource/tencentcloud_nat_gateway: adapt to scenarios where `nat_product_version` is 2 ([#2895](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2895))
* resource/tencentcloud_teo_rule_engine: support rule_priority field. ([#2893](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2893))

## 1.81.133 (October 14 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_postgresql_clone_db_instance` ([#2888](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2888))

ENHANCEMENTS:

* resource/tencentcloud_cdn_domain: Support origin_company field and add computed to the message field. ([#2883](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2883))
* resource/tencentcloud_kubernetes_scale_worker: add params `create_result_output_file` ([#2885](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2885))
* resource/tencentcloud_lite_hbase_instance: fix delete state refresh ([#2889](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2889))
* resource/tencentcloud_security_group_rule_set: update doc ([#2890](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2890))

## 1.81.132 (October 11 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_ssl_check_certificate_domain_verification_operation` ([#2871](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2871))

ENHANCEMENTS:

* ccn/Update sdk version ([#2882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2882))
* resource/tencentcloud_as_scaling_config: Supports `enhanced_automation_tools_service` params. ([#2875](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2875))
* resource/tencentcloud_ccn_routes: fix ccn routes read ([#2877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2877))
* resource/tencentcloud_instance: support `disable_automation_service` params ([#2873](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2873))
* resource/tencentcloud_mysql_cls_log_attachment: Update md doc ([#2879](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2879))

## 1.81.131 (October 1 , 2024)

ENHANCEMENTS:

* resource/tencentcloud_security_group_rule: Supports setting the `ip_protocol` parameter to `ALL`. ([#2870](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2870))

## 1.81.130 (September 30 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_Fimage_from_family` ([#2869](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2869))
* **New Resource:** `tencentcloud_postgresql_apply_parameter_template_operation` ([#2867](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2867))
* **New Resource:** `tencentcloud_teo_security_ip_group` ([#2867](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2864))
* **New Resource:** `tencentcloud_teo_function` ([#2865](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2865))
* **New Resource:** `tencentcloud_teo_function_rule` ([#2865](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2865))
* **New Resource:** `tencentcloud_teo_function_rule_priority` ([#2865](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2865))
* **New Resource:** `tencentcloud_teo_function_runtime_environment` ([#2865](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2865))

ENHANCEMENTS:

* resource/tencentcloud_as_scaling_config: support `image_family` params ([#2869](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2869))
* resource/tencentcloud_cos_bucket: fix the issue where acl_body a cannot be modified ([#2868](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2868))
* resource/tencentcloud_image: support `image_family` params ([#2869](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2869))
* resource/tencentcloud_kubernetes_cluster_attachment: support param `security_groups` ([#2866](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2866))
* resource/tencentcloud_kubernetes_cluster_attachment: support param `taints` of worker_config ([#2866](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2866))
* resource/tencentcloud_kubernetes_scale_worker:  support `taints` parameter ([#2859](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2859))

## 1.81.129 (September 29 , 2024)

ENHANCEMENTS:

* resource/tencentcloud_cos_bucket: support acl for cdc ([#2860](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2860))
* resource/tencentcloud_redis_instance: support force_delete to postpaid instance ([#2861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2861))
* resource/tencentcloud_redis_startup_instance_operation: adjust startup status logic ([#2862](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2862))

BUG FIXES:

* resource/tencentcloud_kubernetes_native_node_pool: fix node pool creating timeout ([#2858](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2858))
* resource/tencentcloud_security_group_rule: fix delete rule failed ([#2863](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2863))

## 1.81.128 (September 27 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_audit_events` ([#2857](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2857))
* **New Resource:** `tencentcloud_redis_log_delivery` ([#2853](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2853))

ENHANCEMENTS:

* data_source/tencentcloud_enis: support `cdc_id` params ([#2855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2855))
* data_source/tencentcloud_vpc_subnets: support `cdc_id` params ([#2855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2855))
* resource/tencentcloud_cos_bucket: support SSE-KMS encryption ([#2848](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2848))
* resource/tencentcloud_eip: support `cdc_id` params ([#2855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2855))
* resource/tencentcloud_kubernetes_node_pool:  support delete `taints` and `labels` params ([#2837](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2837))
* resource/tencentcloud_kubernetes_scale_worker: Lift the upper limit of 100 ([#2850](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2850))
* resource/tencentcloud_monitor_binding_object: update monitor region map ([#2856](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2856))
* resource/tencentcloud_organization_org_member: support `tags` params ([#2852](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2852))
* resource/tencentcloud_organization_org_node: support `tags` params ([#2852](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2852))
* resource/tencentcloud_vpc_end_point: support `cdc_id` params ([#2855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2855))
* resource/tencentcloud_vpc_end_point_service: support `cdc_id` params ([#2855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2855))

## 1.81.127 (September 26 , 2024)

ENHANCEMENTS:

* resource/tencentcloud_eip: support `cdc_id` params ([#2849](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2849))
* resource/tencentcloud_teo_certificate_config: update teo_certificate_config doc ([#2847](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2847))

## 1.81.126 (September 25 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_postgresql_dedicated_clusters` ([#2845](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2845))
* **New Resource:** `tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment` ([#2842](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2842))
* **New Resource:** `tencentcloud_kubernetes_log_config` ([#2840](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2840))

ENHANCEMENTS:

* resource/tencentcloud_cdn_domain: Add more origin type ([#2843](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2843))
* resource/tencentcloud_cos_bucket: COS supports CDC scenarios ([#2841](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2841))
* datasource/tencentcloud_kubernetes_scale_worker: adjust the message when creating instances all failed ([#2839](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2839))
* resource/tencentcloud_kubernetes_cluster: support `cdc_id` and `pre_start_user_script` parameter ([#2835](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2835))
* datasource/tencentcloud_kubernetes_clusters: support `cdc_id` parameter ([#2835](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2835))
* resource/tencentcloud_organization_instance: support param `root_node_name`. ([#2844](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2844))
* resource/tencentcloud_postgresql_instance: support `dedicated_cluster_id` params ([#2845](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2845))
* resource/tencentcloud_postgresql_readonly_instance: support `dedicated_cluster_id` params ([#2845](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2845))

## 1.81.125 (September 23 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_invite_organization_member_operation` ([#2838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2838))

ENHANCEMENTS:

* resource/tencentcloud_cos_bucket: support `cdc_id` params ([#2805](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2805))

## 1.81.124 (September 20 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_lite_hbase_instances` ([#2836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2836))
* **New Data Source:** `tencentcloud_redis_clusters` ([#2761](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2761))
* **New Resource:** `tencentcloud_lite_hbase_instance` ([#2836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2836))

ENHANCEMENTS:

* resource/tencentcloud_redis_instance: support cdc ([#2761](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2761))

## 1.81.123 (September 14 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_kubernetes_health_check_policy` ([#2826](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2826))

ENHANCEMENTS:

* datasource/tencentcloud_ssl_certificates:  add `owner_uin` and `validity_period` params ([#2832](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2832))
* resource/tencentcloud_ccn_route_table_associate_instance_config:  Update document ([#2829](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2829))
* resource/tencentcloud_cos_bucket_inventory: Fix null pointer issue ([#2819](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2819))
* resource/tencentcloud_elasticsearch_instance:  support `cos_backup` ([#2822](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2822))
* resource/tencentcloud_instance:  Update and delete module code logic ([#2830](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2830))
* resource/tencentcloud_thpc_workspaces:  Optimized code ([#2824](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2824))

BUG FIXES:

* resource/tencentcloud_kubernetes_cluster:  Support to update argument `cluster_level` and `auto_upgrade_cluster_level` ([#2828](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2828))

## 1.81.122 (September 13 , 2024)

ENHANCEMENTS:

* resource/tencentcloud_ccn_route_table_associate_instance_config:  Support assume role ([#2823](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2823))

## 1.81.121 (September 11 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_as_start_instance_refresh` ([#2814](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2814))
* **New Resource:** `tencentcloud_thpc_workspaces` ([#2813](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2813))

ENHANCEMENTS:

* resource/tencentcloud_gaap_layer4_listener: support udp listener enable `health_check` ([#2816](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2816))
* resource/tencentcloud_vpc:  Fix the issue where field assistant_cidrs cannot be edited ([#2817](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2817))

## 1.81.120 (September 6 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_postgresql_account_privileges` ([#2807](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2807))
* **New Resource:** `tencentcloud_batch_apply_account_baselines` ([#2803](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2803))
* **New Resource:** `tencentcloud_postgresql_account` ([#2807](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2807))

ENHANCEMENTS:

* resource/tencentcloud_as_scaling_group: Support update `configuration_id` ([#2811](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2811))
* resource/tencentcloud_cdn_domain: Add more origin type ([#2812](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2812))
* resource/tencentcloud_clb_attachment: support `SRV` type ([#2806](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2806))
* resource/tencentcloud_identity_center_role_configuration_permission_policy_attachment: support `role_policy_name` ([#2809](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2809))
* resource/tencentcloud_identity_center_user: read add retry ([#2809](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2809))
* resource/tencentcloud_private_dns_record: Optimize the logic for deleting resources ([#2810](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2810))
* resource/tencentcloud_route_table_entry: fix creation failure ([#2808](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2808))

## 1.81.119 (August 30 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_cdwdoris_instances` ([#2799](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2799))
* **New Data Source:** `tencentcloud_organization_services` ([#2792](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2792))
* **New Resource:** `tencentcloud_cdwdoris_instance` ([#2799](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2799))
* **New Resource:** `tencentcloud_cdwdoris_workload_group` ([#2799](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2799))
* **New Resource:** `tencentcloud_identity_center_external_saml_identity_provider` ([#2795](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2795))
* **New Resource:** `tencentcloud_identity_center_group` ([#2795](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2795))
* **New Resource:** `tencentcloud_identity_center_role_assignment` ([#2795](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2795))
* **New Resource:** `tencentcloud_identity_center_role_configuration` ([#2795](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2795))
* **New Resource:** `tencentcloud_identity_center_role_configuration_permission_policy_attachment` ([#2795](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2795))
* **New Resource:** `tencentcloud_identity_center_user` ([#2795](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2795))
* **New Resource:** `tencentcloud_identity_center_user_group_attachment` ([#2795](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2795))
* **New Resource:** `tencentcloud_identity_center_user_sync_provisioning` ([#2795](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2795))
* **New Resource:** `tencentcloud_organization_service_assign` ([#2792](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2792))

ENHANCEMENTS:

* datasource/tencentcloud_ssl_certificates: Check overclocking support retry. ([#2801](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2801))
* datasource/tencentcloud_ssl_describe_certificate: Check overclocking support retry. ([#2801](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2801))
* resource/tencentcloud_clb_listener_rule: support `domains` ([#2789](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2789))
* resource/tencentcloud_emr_cluster: add `auto_renew` param. ([#2798](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2798))
* resource/tencentcloud_instance: support delete prepaid disk. ([#2794](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2794))
* resource/tencentcloud_kubernetes_node_pool: support `instance_name_style` param ([#2791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2791))
* tencentcloud_mysql_readonly_instance: Read-only instance creation supports setting read-only groups. ([#2800](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2800))

BUG FIXES:

* resource/tencentcloud_kubernetes_addon_attachment: support multiple resources scenario. ([#2797](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2797))

## 1.81.118 (August 23 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_mysql_cls_log_attachment` ([#2780](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2780))

ENHANCEMENTS:

* datasource/tencentcloud_as_scaling_configs: add `version_number`. ([#2786](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2786))
* provider: support tke cam role auth ([#2785](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2785))
* resource/tencentcloud_kubernetes_node_pool: add computed to field `tags`. ([#2781](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2781))
* resource/tencentcloud_mysql_instance: Optimize code logic ([#2752](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2752))
* resource/tencentcloud_mysql_instance: support mysql `engine_type`. ([#2779](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2779))

## 1.81.117 (August 16 , 2024)

FEATURES:

* **New Data Source:** `tencentcloud_cls_logsets` ([#2775](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2775))

ENHANCEMENTS:

* provider: support cam role name auth ([#2767](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2767))
* resource/tencentcloud_clb_attachment: support domain and url params ([#2776](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2776))
* resource/tencentcloud_postgresql_instance: support delete protection ([#2777](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2777))
* resource/tencentcloud_vpn_gateway: Optimize VPN gateway change issue. ([#2778](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2778))

## 1.81.116 (August 14 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_mysql_ssl` ([#2687](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2687))

ENHANCEMENTS:

* resource/tencentcloud_as_lifecycle_hook: support import and parameter `lifecycle_command` ([#2772](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2772))
* resource/tencentcloud_as_scaling_group: support `health_check_type` and `lb_health_check_grace_period` ([#2772](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2772))
* resource/tencentcloud_dnspod_record: compatible with record value ends with dot ([#2773](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2773))
* resource/tencentcloud_mongodb_instance: support `maintenance_start` and `maintenance_end` ([#2765](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2765))

BUG FIXES:

* resource/tencentcloud_dnspod_record: fix weight auto update to 0 when update value ([#2770](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2770))
* resource/tencentcloud_kubernetes_cluster_attachment: fix the issue that param `unschedulable` cannot work ([#2764](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2764))
* resource/tencentcloud_kubernetes_scale_worker: fix param `data_disk` of worker_config cannot work ([#2769](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2769))
* resource/tencentcloud_kubernetes_scale_worker: fix the cluster instance paging logic ([#2774](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2774))

## 1.81.115 (August 9 , 2024)

FEATURES:

* **New Resource:** `tencentcloud_vpc_private_nat_gateway` ([#2763](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2763))

ENHANCEMENTS:

* resource/tencentcloud_emr_cluster: fix read need_master_wan ([#2768](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2768))
* resource/tencentcloud_kubernetes_cluster: support param `ignore_service_cidr_conflict` ([#2756](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2756))

## 1.81.114 (August 2 , 2024)

ENHANCEMENTS:

* data_source/tencentcloud_instances: add params `uuid` ([#2757](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2757))
* resource/tencentcloud_cls_data_transform: Update read interface ([#2750](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2750))
* resource/tencentcloud_elasticsearch_instance: Update disk_type params ([#2751](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2751))
* resource/tencentcloud_kubernetes_cluster_attachment: support `image_id` ([#2749](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2749))
* resource/tencentcloud_vpn_gateway: add default value for `prepaid_period` while import ([#2759](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2759))

## 1.81.113 (July 26, 2024)

FEATURES:

* **New Data Source:** `tencentcloud_cdc_dedicated_cluster_hosts` ([#2737](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2737))
* **New Data Source:** `tencentcloud_cdc_dedicated_cluster_instance_types` ([#2737](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2737))
* **New Data Source:** `tencentcloud_cdc_dedicated_cluster_orders` ([#2737](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2737))
* **New Resource:** `tencentcloud_cdc_dedicated_cluster` ([#2737](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2737))
* **New Resource:** `tencentcloud_cdc_site` ([#2737](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2737))

ENHANCEMENTS:

* provider: add SAML, OIDC for STS client ([#2742](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2742))
* resource/tencentcloud_ccn: support route_ecmp_flag and route_overlap_flag ([#2746](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2746))

BUG FIXES:

* resource/tencentcloud_kubernetes_cluster_attachment: fix the issue triggered by `labels` ([#2745](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2745))

## 1.81.112 (July 24, 2024)

ENHANCEMENTS:

* data_source/tencentcloud_cbs_storages: add params `dedicated_cluster_id` ([#2715](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2715))
* data_source/tencentcloud_cbs_storages_set: add params `dedicated_cluster_id` ([#2715](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2715))
* data_source/tencentcloud_clb_instances: add params `dedicated_cluster_id` ([#2720](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2720))
* data_source/tencentcloud_cvm_chc_hosts: deprecated `host_ips` param ([#2719](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2719))
* data_source/tencentcloud_instances: add params `dedicated_cluster_id` ([#2719](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2719))
* provider: support securityToken for credentials ([#2739](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2739))
* resource/tencentcloud_cbs_storage: add params `dedicated_cluster_id` ([#2715](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2715))
* resource/tencentcloud_cbs_storage_set: add params `dedicated_cluster_id` ([#2715](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2715))
* resource/tencentcloud_clb_instance: add params `dedicated_cluster_id` ([#2720](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2720))
* resource/tencentcloud_cls_topic: add `is_web_tracking`, `extends` params ([#2721](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2721))
* resource/tencentcloud_cvm_launch_template: deprecated `host_ips` param ([#2719](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2719))
* resource/tencentcloud_cynosdb_cluster: add `slave_zone` param ([#2711](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2711))
* resource/tencentcloud_dasb_user: update param phone Adapter country area code ([#2734](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2734))
* resource/tencentcloud_instance: add param `system_disk_resize_online` ([#2740](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2740))
* resource/tencentcloud_instance: add params `dedicated_cluster_id` ([#2719](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2719))

## 1.81.111 (July 19, 2024)

ENHANCEMENTS:

* resource/tencentcloud_kubernetes_scale_worker: throw error while all instances create failed ([#2735](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2735))

## 1.81.110 (July 17, 2024)

FEATURES:

* **New Resource:** `tencentcloud_ccn_route_table_input_policies` ([#2730](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2730))
* **New Resource:** `tencentcloud_ccn_route_table_selection_policies` ([#2730](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2730))

## 1.81.109 (July 17, 2024)

FEATURES:

* **New Data Source:** `tencentcloud_ccn_routes` ([#2730](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2730))
* **New Data Source:** `tencentcloud_monitor_tmp_instances` ([#2731](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2731))
* **New Resource:** `tencentcloud_ccn_route_table` ([#2730](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2730))
* **New Resource:** `tencentcloud_ccn_route_table_associate_instance_config` ([#2730](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2730))
* **New Resource:** `tencentcloud_ccn_route_table_broadcast_policies` ([#2730](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2730))
* **New Resource:** `tencentcloud_ccn_route_table_input_policies` ([#2730](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2730))
* **New Resource:** `tencentcloud_ccn_route_table_selection_policies` ([#2730](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2730))
* **New Resource:** `tencentcloud_kubernetes_addon_config` ([#2725](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2725))

ENHANCEMENTS:

* resource/tencentcloud_clb_attachment: support `CCN` backend ([#2729](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2729))
* resource/tencentcloud_mongodb_instance: add describe retry ([#2727](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2727))
* resource/tencentcloud_mongodb_sharding_instance: add describe retry ([#2727](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2727))
* resource/tencentcloud_mongodb_standby_instance: add describe retry ([#2727](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2727))
* resource/tencentcloud_route_table_entry: add computed `route_item_id` ([#2732](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2732))

## 1.81.108 (July 12, 2024)

ENHANCEMENTS:

* data_source/tencentcloud_subnet: add params `cdc_id` ([#2717](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2717))
* data_source/tencentcloud_vpc_subnets: add params `cdc_id` ([#2717](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2717))
* resource/tencentcloud_clb_instance: support update `project_id` params. ([#2705](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2705))
* resource/tencentcloud_kubernetes_scale_worker: support retry ([#2723](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2723))
* resource/tencentcloud_monitor_tmp_alert_group: add retry operator ([#2718](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2718))
* resource/tencentcloud_rum_project: add retry operator ([#2718](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2718))
* resource/tencentcloud_subnet: add params `cdc_id` ([#2717](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2717))

## 1.81.107 (July 5, 2024)

ENHANCEMENTS:

* resource/tencentcloud_apm_instance: add retry operator ([#2707](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2707))
* resource/tencentcloud_elasticsearch_instance: support es kibana switch ([#2708](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2708))
* resource/tencentcloud_monitor_tmp_instance: update document ([#2704](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2704))
* resource/tencentcloud_postgresql_readonly_group: support import ([#2709](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2709))
* resource/tencentcloud_tdmq_rabbitmq_user: supports return `max_connections`, `max_channels` params; supports import this resource ([#2703](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2703))
* resource/tencentcloud_tdmq_rabbitmq_vip_instance: update document ([#2703](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2703))
* resource/tencentcloud_tdmq_rabbitmq_virtual_host: supports set `trace_flag` params; supports import this resource ([#2703](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2703))

## 1.81.106 (June 28, 2024)

FEATURES:

* **New Resource:** `tencentcloud_kubernetes_native_node_pool` ([#2658](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2658))
* **New Resource:** `tencentcloud_mongodb_instance_backup_rule` ([#2692](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2692))
* **New Resource:** `tencentcloud_mongodb_instance_transparent_data_encryption` ([#2692](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2692))
* **New Resource:** `tencentcloud_teo_realtime_log_delivery` ([#2697](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2697))

ENHANCEMENTS:

* resource/tencentcloud_clb_instance: open ipv6 clb support set `subnet_id` ([#2696](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2696))
* resource/tencentcloud_cls_alarm: supports set `syntax_rule` ([#2699](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2699))
* resource/tencentcloud_tdmq_instance: Fix the issue of returning multiple return values ([#2690](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2690))
* resource/tencentcloud_tdmq_namespace: Fix the issue of returning multiple return values ([#2690](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2690))
* resource/tencentcloud_tdmq_namespace_role_attachment: Fix the issue of returning multiple return values ([#2690](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2690))
* resource/tencentcloud_tdmq_role: Fix the issue of returning multiple return values ([#2690](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2690))
* resource/tencentcloud_tdmq_topic: Fix the issue of returning multiple return values ([#2690](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2690))
* resource/tencentcloud_trocket_rocketmq_instance: support more specifications of param `sku_code` for create instance ([#2689](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2689))

## 1.81.105 (June 17, 2024)

ENHANCEMENTS:

* resource/tencentcloud_teo_acceleration_domain: Add `origin_protocol`, `http_origin_port`, `https_origin_port` and `ipv6_status` fields ([#2685](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2685))

## 1.81.104 (June 14, 2024)

ENHANCEMENTS:

* resource/tencentcloud_clb_instance: Change the record list type to a set to solve the order problem. ([#2680](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2680))
* resource/tencentcloud_postgresql_instance: Fix the issue of failed instance creation read ([#2682](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2682))

BUG FIXES:

* resource/tencentcloud_mysql_instance: If the id is empty, retry with the same clienttoken for 10 minutes ([#2684](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2684))
* resource/tencentcloud_user_info: fix `owner_uin` ([#2683](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2683))

## 1.81.103 (June 13, 2024)

ENHANCEMENTS:

* resource/tencentcloud_clb_instance: Support ipv6 return value `ipv6_mode` and `address_ipv6` ([#2675](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2675))
* resource/tencentcloud_kubernetes_cluster_endpoint: Support ForceNew for `cluster_id` ([#2681](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2681))

## 1.81.102 (June 11, 2024)

ENHANCEMENTS:

* resource/tencentcloud_cynosdb_cluster: modify tag. ([#2674](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2674))
* resource/tencentcloud_teo_rule_engine: Change action to Optional ([#2672](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2672))
* resource/tencentcloud_user_info: Support retry ([#2671](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2671))

## 1.81.101 (June 5, 2024)

ENHANCEMENTS:

* resource/tencentcloud_postgresql_instance: support params db_major_version create resource ([#2618](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2618))

BUG FIXES:

* resource/tencentcloud_cfs_file_system: fix null fs_id ([#2669](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2669))

## 1.81.100 (June 4, 2024)

ENHANCEMENTS:

* resource/tencentcloud_cls_index: optimize example ([#2667](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2667))
* resource/tencentcloud_tdmq_rocketmq_vip_instance: support `ip_rules` params ([#2659](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2659))

## 1.81.99 (May 28, 2024)

ENHANCEMENTS:

* resource/tencentcloud_organization_org_manage_policy: Supplementary documentation ([#2656](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2656))
* resource/tencentcloud_organization_org_manage_policy_config: Supplementary documentation ([#2656](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2656))
* resource/tencentcloud_organization_org_manage_policy_target: Supplementary documentation ([#2656](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2656))

## 1.81.98 (May 24, 2024)

ENHANCEMENTS:

* resource/tencentcloud_cls_alarm: Support `alarm_level` params and fix import issues. ([#2653](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2653))
* resource/tencentcloud_kubernetes_scale_worker: Add `pre_start_user_script` and `user_script` fields ([#2651](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2651))

## 1.81.97 (May 23, 2024)

ENHANCEMENTS:

* tencentcloud_postgresql_instance: support change pg password regardless user ([#2652](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2652))

## 1.81.96 (May 22, 2024)

FEATURES:

* **New Resource:** `tencentcloud_teo_l4_proxy` ([#2617](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2617))

ENHANCEMENTS:

* resource/tencentcloud_teo_origin_group: Incompatible changes ([#2617](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2617))
* tencentcloud_cam_role: `session_duration` increases computed familiarity ([#2641](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2641))
* tencentcloud_mysql_dr_instance: modify doc ([#2647](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2647))

## 1.81.95 (May 17, 2024)

ENHANCEMENTS:

* datasource/tencentcloud_kubernetes_cluster_instances: querying cluster node information limit supports 100. ([#2637](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2637))
* resource/tencentcloud_cynosdb_readonly_instance: improvement id empty situation ([#2636](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2636))

## 1.81.94 (May 15, 2024)

ENHANCEMENTS:

* resource/tencentcloud_ckafka_instance: remove `public-network` deprecated and validate minimum value is 3 ([#2631](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2631))
* resource/tencentcloud_instance: Update the limit of period and charge type ([#2630](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2630))
* resource/tencentcloud_redis_instance: Optimize Availability Zone ([#2634](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2634))

## 1.81.93 (May 13, 2024)

FEATURES:

* **New Resource:** `tencentcloud_kubernetes_addon` ([#2624](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2624))
* **New Resource:** `tencentcloud_tdmq_topic_with_full_id` ([#2625](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2625))

ENHANCEMENTS:

* resource/tencentcloud_clb_attachment: optimization documentation ([#2622](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2622))
* resource/tencentcloud_mongodb_instance: support update security group ([#2628](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2628))
* resource/tencentcloud_mongodb_sharding_instance: support update security group ([#2628](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2628))
* resource/tencentcloud_mongodb_standby_instance: support update security group ([#2628](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2628))

## 1.81.92 (May 6, 2024)

ENHANCEMENTS:

* resource/tencentcloud_clb_listener: Support not checking healthcheck parameters ([#2598](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2598))
* resource/tencentcloud_monitor_tmp_alert_rule: Fix array order issue ([#2620](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2620))
* resource/tencentcloud_mysql_privilege: Fix the problem of multiple values ​​in fields database, table, column ([#2619](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2619))

## 1.81.91 (April 30, 2024)

ENHANCEMENTS:

* resource/tencentcloud_ckafka_route: Updated ckafka route documentation ([#2614](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2614))

## 1.81.90 (April 29, 2024)

FEATURES:

* **New Resource:** `tencentcloud_mysql_dr_instance` ([#2596](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2596))
* **New Resource:** `tencentcloud_organization_org_manage_policy` ([#2604](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2604))
* **New Resource:** `tencentcloud_organization_org_manage_policy_config` ([#2604](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2604))
* **New Resource:** `tencentcloud_organization_org_manage_policy_target` ([#2604](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2604))

ENHANCEMENTS:

* datasource/tencentcloud_instances: add output `os_name` ([#2608](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2608))
* resource/tencentcloud_cls_index: add dynamic_index params ([#2606](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2606))
* resource/tencentcloud_instance: add output `memory`, `os_name`, `cpu` ([#2608](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2608))
* resource/tencentcloud_postgresql_instance: Support field `cpu` ([#2605](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2605))
* resource/tencentcloud_postgresql_readonly_instance: Support field `cpu` ([#2605](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2605))

## 1.81.89 (April 23, 2024)

ENHANCEMENTS:

* resource/tencentcloud_gaap_proxy: fix gaap time layout ([#2597](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2597))
* resource/tencentcloud_postgresql_readonly_group: Support the computed param `ip` of the `net_info_list` ([#2595](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2595))

## 1.81.88 (April 18, 2024)

ENHANCEMENTS:

* resource/tencentcloud_scf_function: return function_id params ([#2590](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2590))
* resource/tencentcloud_tcr_instance: Support `region_name` parma ([#2593](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2593))

## 1.81.87 (April 7, 2024)

FEATURES:

* **New Resource:** `tencentcloud_tse_cngw_network_access_control` ([#2581](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2581))

ENHANCEMENTS:

* datasource/tencentcloud_mariadb_db_instances: Support `vip`, `vport`, `internet_domain`, `internet_ip`, `internet_port` field ([#2583](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2583))
* resource/tencentcloud_kubernetes_cluster: Field `cluster_subnet_id` supports ForceNew ([#2578](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2578))
* resource/tencentcloud_vod_adaptive_dynamic_streaming_template: add param `segment_type` ([#2582](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2582))
* resource/tencentcloud_vod_procedure_template: update param `adaptive_dynamic_streaming_task_list` and `review_audio_video_task` ([#2582](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2582))

## 1.81.86 (March 28, 2024)

FEATURES:

* **New Resource:** `tencentcloud_vod_event_config` ([#2575](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2575))

## 1.81.85 (March 27, 2024)

FEATURES:

* **New Resource:** `tencentcloud_eni_ipv4_address` ([#2574](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2574))
* **New Resource:** `tencentcloud_eni_ipv6_address` ([#2573](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2573))

## 1.81.84 (March 25, 2024)

ENHANCEMENTS:

* datasource/tencentcloud_ckafka_topics: Update ckafka topics docs ([#2570](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2570))
* resource/tencentcloud_vod_adaptive_dynamic_streaming_template: Add vcrf, gop, preserve_hdr_switch, codec_tag. Adjust resource unique id to subAppId#templateId. ([#2569](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2569))
* resource/tencentcloud_vod_image_sprite_template: Add format, type. Adjust resource unique id to subAppId#templateId. ([#2569](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2569))
* resource/tencentcloud_vod_procedure_template: Update params media_process_task, add ai_analysis_task, ai_recognition_task, review_audio_video_task, type. Adjust resource unique id to subAppId#templateId. ([#2569](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2569))
* resource/tencentcloud_vod_sample_snapshot_template: Adjust resource unique id to subAppId#templateId. ([#2569](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2569))
* resource/tencentcloud_vod_snapshot_by_time_offset_template: Add type. Adjust resource unique id to subAppId#templateId. ([#2569](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2569))

## 1.81.83 (March 20, 2024)

FEATURES:

* **New Resource:** `tencentcloud_dasb_resource` ([#2566](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2566))

ENHANCEMENTS:

* resource/tencentcloud_instance: Output `Uuid` of cvm instance ([#2568](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2568))
* resource/tencentcloud_private_dns_record: params `zone_id` Add ForceNew ([#2567](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2567))

## 1.81.82 (March 15, 2024)

ENHANCEMENTS:

* resource/tencentcloud_dasb_user_group_members: Fix errors in querying user group members ([#2564](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2564))

## 1.81.81 (March 13, 2024)

FEATURES:

* **New Resource:** `tencentcloud_tse_cngw_network` ([#2526](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2526))
* **New Resource:** `tencentcloud_tse_cngw_strategy` ([#2526](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2526))
* **New Resource:** `tencentcloud_tse_cngw_strategy_bind_group` ([#2526](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2526))

ENHANCEMENTS:

* datasource/tencentcloud_kubernetes_cluster_instances: optimize cluster instance queries ([#2560](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2560))
* resource/tencentcloud_tse_cngw_group: Add `group_id` field to output. ([#2526](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2526))

## 1.81.80 (March 11, 2024)

BUG FIXES:

* resource/tencentcloud_ckafka_instance: fix `charge_type` has change after import. ([#2557](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2557))

## 1.81.79 (March 8, 2024)

FEATURES:

* **New Resource:** `tencentcloud_vod_sample_snapshot_template` ([#2553](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2553))
* **New Resource:** `tencentcloud_vod_transcode_template` ([#2553](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2553))
* **New Resource:** `tencentcloud_vod_watermark_template` ([#2553](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2553))

ENHANCEMENTS:

* resource/tencentcloud_api_gateway_api_key_attachment: Remove private information from documents. ([#2556](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2556))
* resource/tencentcloud_private_dns_zone: Add ForceNew attribute to domain parameter. ([#2555](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2555))

BUG FIXES:

* resource/tencentcloud_kubernetes_node_pool: Fix the problem of creating node pool tags field. ([#2554](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2554))

## 1.81.78 (March 3, 2024)

ENHANCEMENTS:

* resource/tencentcloud_mysql_instance: Fixed the problem of clienttoken duplication caused by creation retry. ([#2549](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2549))

BUG FIXES:

* resource/tencentcloud_instance: fix the issue of invalid parameters when creating CVM. ([#2551](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2551))
* resource/tencentcloud_monitor_tmp_tke_cluster_agent: Fix the error reported when modifying the associated cluster. ([#2550](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2550))

## 1.81.77 (March 1, 2024)

FEATURES:

* **New Resource:** `tencentcloud_csip_risk_center` ([#2498](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2498))

ENHANCEMENTS:

* resource/tencentcloud_cls_topic: support hot_period, describes params ([#2545](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2545))
* resource/tencentcloud_vod_procedure_template: fix api must set SubAppId ([#2542](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2542))
* resource/tencentcloud_vod_sub_application: fix status update problem ([#2542](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2542))

BUG FIXES:

* datasource/tencentcloud_tcmq_topic: fix `topic_list.topic_id` value ([#2546](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2546))

## 1.81.76 (February 28, 2024)

FEATURES:

* **New Data Source:** `tencentcloud_private_dns_private_zone_list` ([#2539](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2539))

ENHANCEMENTS:

* datasource/tencentcloud_kubernetes_clusters: Fix `lan_ip` field return value ([#2541](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2541))
* resource/tencentcloud_kubernetes_node_pool: add `pre_start_user_script` field ([#2543](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2543))

BUG FIXES:

* resource/tencentcloud_clb_target_group_attachments: Fixed the problem of not supporting four-layer LB ([#2489](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2489))

## 1.81.75 (February 23, 2024)

ENHANCEMENTS:

* resource/tencentcloud_kubernetes_cluster_endpoint: Endpoint resources increase sleep time ([#2538](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2538))
* resource/tencentcloud_mongodb_instance_account: support import ([#2536](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2536))

## 1.81.74 (February 21, 2024)

ENHANCEMENTS:

* datasource/tencentcloud_kubernetes_clusters: add `kube_config_file_prefix` field ([#2530](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2530))
* resource/tencentcloud_kubernetes_cluster: Optimize the change problem of `cluster_internet_security_group` and `cluster_intranet_subnet_id` fields ([#2533](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2533))
* resource/tencentcloud_mongodb_instance: support params `add_node_list` and `remove_node_list` ([#2531](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2531))
* resource/tencentcloud_monitor_tmp_exporter_integration: Documentation example added. ([#2527](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2527))

## 1.81.73 (February 5, 2024)

FEATURES:

* **New Data Source:** `tencentcloud_organization_org_share_area` ([#2520](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2520))
* **New Resource:** `tencentcloud_api_gateway_update_service` ([#2518](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2518))
* **New Resource:** `tencentcloud_organization_org_share_unit` ([#2521](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2521))
* **New Resource:** `tencentcloud_organization_org_share_unit_member` ([#2521](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2521))

ENHANCEMENTS:

* resource/tencentcloud_clb_attachments: support param `eni_ip` ([#2512](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2512))
* resource/tencentcloud_emr_cluster: update `login_settings` description ([#2515](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2515))
* resource/tencentcloud_kubernetes_cluster: Optimize `cluster_internet` and `cluster_intranet` fields ([#2511](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2511))
* resource/tencentcloud_kubernetes_node_pool: update `system_disk_type` and `disk_type` field description ([#2514](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2514))
* resource/tencentcloud_lighthouse_instance: support params `public_addresses` and `private_addresses` ([#2508](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2508))
* resource/tencentcloud_mps_schedule: Add `resource_id` field ([#2513](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2513))
* resource/tencentcloud_mps_schedules: Add `resource_id` field ([#2513](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2513))
* resource/tencentcloud_mysql_instance: Support configuration modification when active/standby switchover occurs. ([#2516](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2516))
* resource/tencentcloud_tdmq_rabbitmq_vip_instance: Adapt to imported resources ([#2510](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2510))

## 1.81.72 (January 30, 2024)

ENHANCEMENTS:

* resource/tencentcloud_ccn_attachment: support tke cluster addon modify ([#2507](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2507))
* resource/tencentcloud_instance: fix private ip release problem ([#2509](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2509))
* resource/tencentcloud_mysql_instance: Optimize the availability zone problem when modifying the configuration after active/standby switchover. ([#2502](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2502))

## 1.81.71 (January 26, 2024)

ENHANCEMENTS:

* datasource/tencentcloud_ssl_certificates: Filter out invalid certificates without reporting errors directly ([#2500](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2500))
* resource/tencentcloud_dnspod_record: fix modify dnspod record weight ([#2506](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2506))
* resource/tencentcloud_instance: wait cvm private ip release when instance destroyed ([#2503](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2503))
* resource/tencentcloud_private_dns_zone: Increase the waiting time for instance creation ([#2504](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2504))

## 1.81.70 (January 24, 2024)

FEATURES:

* **New Resource:** `tencentcloud_monitor_tmp_alert_group` ([#2487](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2487))

ENHANCEMENTS:

* resource/tencentcloud_clb_attachment: support cross-domain binding, `target` status writing ([#2497](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2497))
* resource/tencentcloud_clb_instance: Remove `vip` restrictions ([#2493](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2493))
* resource/tencentcloud_clb_listener: Support UDP health detection and handle the problem of COSTOM returning to default value ([#2496](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2496))
* resource/tencentcloud_mysql_instance: `availability_zone` modification supports error reporting ([#2495](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2495))
* resource/tencentcloud_security_groups: Supports more than one hundred data returns ([#2494](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2494))

## 1.81.69 (January 19, 2024)

ENHANCEMENTS:

* data_source/tencentcloud_customer_gateways: Optimization parameter error ([#2491](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2491))
* resource/tencentcloud_cfw_edge_policy: Add return parameters ([#2481](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2481))
* resource/tencentcloud_cfw_nat_policy: Add return parameters ([#2481](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2481))
* resource/tencentcloud_cfw_vpc_policy: Add return parameters ([#2481](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2481))
* resource/tencentcloud_clb_instance: Support set specified `vip` ([#2485](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2485))
* resource/tencentcloud_vpn_connection: fix description for `ipsec_pfs_dh_group` ([#2486](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2486))

## 1.81.68 (January 17, 2024)

BUG FIXES:

* resource/tencentcloud_kubernetes_scale_worker: fix DescribeClusters return parameter `Property` parsing code ([#2480](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2480))

## 1.81.67 (January 17, 2024)

FEATURES:

* **New Data Source:** `tencentcloud_clickhouse_instance_nodes` ([#2483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2483))

ENHANCEMENTS:

* resource/tencentcloud_monitor_grafana_instance: support `auto_voucher` field ([#2473](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2473))
* resource/tencentcloud_vpc_end_point: Support set `security_groups_ids` ([#2482](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2482))

## 1.81.66 (January 15, 2024)

ENHANCEMENTS:

* resource/tencentcloud_instance: `instance_charge_type` support `UNDERWRITE` ([#2478](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2478))
* resource/tencentcloud_instance: fix tag conflict ([#2464](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2464))
* resource/tencentcloud_mysql_instance: Modify configuration to support window period switching ([#2472](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2472))

## 1.81.65 (January 12, 2024)

ENHANCEMENTS:

* resource/resource_tc_redis_instance: Increase `mem_size` optional value ([#2459](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2459))
* resource/target_group_attachments: support parallel ([#2467](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2467))
* resource/tencentcloud_clb_listener_rule: support parallel ([#2467](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2467))
* resource/tencentcloud_kubernetes_cluster: add `cluster_internet` and `cluster_intranet` configuration check ([#2469](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2469))

## 1.81.64 (January 5, 2024)

FEATURES:

* **New Resource:** `tencentcloud_tdmq_subscription` ([#2451](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2451))

ENHANCEMENTS:

* resource/tencentcloud_dcx: Add waiting tunnel ready logic ([#2460](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2460))
* resource/tencentcloud_emr_cluster: support import ([#2461](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2461))
* resource/tencentcloud_waf_clb_domain: Adapt to more scenarios of `bot_status` ([#2458](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2458))
* resource/tencentcloud_waf_saas_domain: Adapt to more scenarios of `bot_status` ([#2458](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2458))
* tencentcloud_clb_target_group_attachments: Support two-way binding of clb and targetGroup ([#2454](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2454))

BUG FIXES:

* tencentcloud_clb_redirection: Fix null pointer exception ([#2457](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2457))

## 1.81.63 (January 3, 2024)

FEATURES:

* **New Data Source:** `tencentcloud_oceanus_job_events` ([#2452](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2452))
* **New Data Source:** `tencentcloud_oceanus_meta_table` ([#2452](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2452))

ENHANCEMENTS:

* datasource/tencentcloud_cos_buckets: add `abort_incomplete_multipart_upload` field ([#2449](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2449))
* resource/tencentcloud_cos_bucket: add `abort_incomplete_multipart_upload` field ([#2449](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2449))
* resource/tencentcloud_dcx: Support import ([#2448](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2448))
* resource/tencentcloud_dnspod_domain_instance: Add computed `slave_dns` param ([#2450](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2450))
* resource/tencentcloud_ssm_ssh_key_pair_secret: Fix ssh_key_name problem ([#2444](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2444))
* tencentcloud_mysql_instance: Support upgrade switching during window period ([#2429](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2429))
* tencentcloud_mysql_readonly_instance: Fix read-only instance `slave_deploy_mode` problem ([#2429](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2429))
* tencentcloud_redis_instance: Support upgrade switching during window period ([#2429](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2429))

## 1.81.62 (December 29, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_oceanus_job_events` ([#2442](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2442))

ENHANCEMENTS:

* datasource/tencentcloud_tdmq_pro_instances: support params `create_time` and `tags` ([#2431](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2431))
* resource/tencentcloud_apm_instance: support param `pay_mode` ([#2432](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2432))
* resource/tencentcloud_cat_task_set: support param `node_ip_type` ([#2433](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2433))
* resource/tencentcloud_lighthouse_instance: support import ([#2443](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2443))
* resource/tencentcloud_mariadb_hour_db_instance: Fix tag problem ([#2434](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2434))
* resource/tencentcloud_mariadb_hour_db_instance: Fix zone sorting problem ([#2435](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2435))
* resource/tencentcloud_ssm_ssh_key_pair_secret: Fix kms_key_id problem ([#2440](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2440))

## 1.81.61 (December 22, 2023)

ENHANCEMENTS:

* resource/tencentcloud_bi_datasource_cloud: support param cluster_id ([#2422](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2422))
* resource/tencentcloud_bi_embed_token_apply: support param ticket_num ([#2422](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2422))
* resource/tencentcloud_ssl_update_certificate_instance_operation: Support upload ca ([#2421](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2421))

## 1.81.60 (December 20, 2023)

FEATURES:

* **New Resource:** `tencentcloud_vpc_peer_connect_accept_operation` ([#2415](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2415))
* **New Resource:** `tencentcloud_vpc_peer_connect_manager` ([#2415](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2415))
* **New Resource:** `tencentcloud_vpc_peer_connect_reject_operation` ([#2415](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2415))

ENHANCEMENTS:

* resource/tencentcloud_cos_bucket: Support `acl_body` import. ([#2417](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2417))
* resource/tencentcloud_scf_function: Update `runtime` description and add example with triggers. ([#2412](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2412))
* resource/tencentcloud_vpn_connection: Support creating routing type vpn connection. ([#2409](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2409))
* tencentcloud_ssl_describe_certificate: Support param `company_type` ([#2418](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2418))

BUG FIXES:

* resource/tencentcloud_dnspod_record: Fix `mx` parameter passing problem. ([#2413](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2413))
* resource/tencentcloud_lighthouse_disk: update the completion condition of disk isolation ([#2414](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2414))

## 1.81.59 (December 15, 2023)

FEATURES:

* **New Resource:** `tencentcloud_tdmq_professional_cluster` ([#2403](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2403))

ENHANCEMENTS:

* resource/tencentcloud_clb_listener: Add parameters `health_source_ip_type`, `session_type` and `keepalive_enable`. ([#2405](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2405))
* resource/tencentcloud_clb_listener_rule: Add parameter `quic`. ([#2406](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2406))

## 1.81.58 (December 13, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_sqlserver_desc_ha_log` ([#2394](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2394))
* **New Resource:** `tencentcloud_clb_target_group_attachments` ([#2398](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2398))
* **New Resource:** `tencentcloud_sqlserver_instance_ssl` ([#2394](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2394))

ENHANCEMENTS:

* data_source/tencentcloud_sqlserver_ins_attribute: Support query `ssl_config` param. ([#2394](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2394))
* resource/tencentcloud_mysql_instance: Fix the availability zone problem after active/standby switchover. ([#2401](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2401))
* resource/tencentcloud_ssm_product_secret: Support create ssm secret for tdsql-c-mysql. ([#2399](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2399))

## 1.81.57 (December 11, 2023)

ENHANCEMENTS:

* resource/tencentcloud_monitor_alarm_policy: Support `filter`, `group_by` field. ([#2390](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2390))
* resource/tencentcloud_private_dns_zone_vpc_attachment: Asynchronous operation adaptation. ([#2391](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2391))

## 1.81.56 (December 8, 2023)

FEATURES:

* **New Data Source:** `resource/tencentcloud_dlc_describe_data_engine_events` ([#2387](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2387))
* **New Resource:** `tencentcloud_postgresql_instance_ha_config` ([#2388](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2388))

ENHANCEMENTS:

* resource/tencentcloud_sqlserver_basic_instance: Support set `collect` param. ([#2383](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2383))
* resource/tencentcloud_sqlserver_general_cloud_instance: Support more db version. ([#2383](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2383))

## 1.81.55 (December 6, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_clickhouse_instance_shards` ([#2382](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2382))
* **New Data Source:** `tencentcloud_clickhouse_spec` ([#2382](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2382))
* **New Data Source:** `tencentcloud_emr_auto_scale_records` ([#2375](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2375))
* **New Resource:** `tencentcloud_clickhouse_keyval_config` ([#2377](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2377))
* **New Resource:** `tencentcloud_clickhouse_xml_config` ([#2377](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2377))
* **New Resource:** `tencentcloud_gaap_custom_header` ([#2378](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2378))
* **New Resource:** `tencentcloud_gaap_proxy_group` ([#2376](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2376))

ENHANCEMENTS:

* resource/tencentcloud_gaap_global_domain: Support param status ([#2378](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2378))
* resource/tencentcloud_gaap_http_domain: Support update param domain ([#2378](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2378))
* resource/tencentcloud_redis_instance: Support upgrade version ([#2380](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2380))

BUG FIXES:

* resource/tencentcloud_vpc: Fix modify `assistant_cidrs` request ([#2356](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2356))

## 1.81.54 (December 4, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_elasticsearch_describe_index_list` ([#2371](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2371))
* **New Data Source:** `tencentcloud_organization_members` ([#2370](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2370))
* **New Resource:** `tencentcloud_oceanus_folder` ([#2372](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2372))

## 1.81.53 (December 1, 2023)

ENHANCEMENTS:

* datasource/tencentcloud_scf_function: add container_image_accelerate, image_port, dns_cache and intranet_config parameters ([#2308](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2308))
* resource/tencentcloud_dayu_ddos_policy_v2: Support param water_print_config ([#2367](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2367))
* resource/tencentcloud_kubernetes_cluster: set up the cluster and import it into Terraform ([#2339](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2339))
* resource/tencentcloud_scf_function: add container_image_accelerate, image_port, dns_cache and intranet_config parameters ([#2308](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2308))
* tencentcloud_kubernetes_cluster: Support `eni_subnet_ids` updates. ([#2346](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2346))

BUG FIXES:

* resource/tencentcloud_organization_org_member: Fix parameter types and update logic ([#2366](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2366))
* resource/tencentcloud_organization_org_member_email: Fix parameter types ([#2366](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2366))
* resource/tencentcloud_scf_function: fix the issue where only one of multiple tag key updates takes effect ([#2308](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2308))

## 1.81.52 (November 29, 2023)

FEATURES:

* **New Resource:** `tencentcloud_antiddos_cc_black_white_ip` ([#2352](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2352))
* **New Resource:** `tencentcloud_antiddos_cc_precision_policy` ([#2352](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2352))
* **New Resource:** `tencentcloud_antiddos_packet_filter_config` ([#2352](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2352))
* **New Resource:** `tencentcloud_antiddos_port_acl_config` ([#2352](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2352))
* **New Resource:** `tencentcloud_tse_waf_domains` ([#2348](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2348))
* **New Resource:** `tencentcloud_tse_waf_protection` ([#2348](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2348))

ENHANCEMENTS:

* resource/tencentcloud_clb_target_group_attachment: Support async task limit ([#2357](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2357))
* tencentcloud_postgresql_readonly_instance: Optimize isolation operations. ([#2355](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2355))

BUG FIXES:

* resource/tencentcloud_kubernetes_node_pool: Add verification when modifying the maximum value and expected value synchronously. ([#2354](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2354))

## 1.81.51 (November 27, 2023)

FEATURES:

* **New Resource:** `tencentcloud_route_table_association` ([#2345](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2345))

ENHANCEMENTS:

* tencentcloud_clb_instance: Support set `vip_isp`. ([#2344](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2344))

## 1.81.50 (November 25, 2023)

FEATURES:

* **New Resource:** `tencentcloud_antiddos_ddos_geo_ip_block_config` ([#2342](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2342))
* **New Resource:** `tencentcloud_antiddos_ddos_speed_limit_config` ([#2342](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2342))
* **New Resource:** `tencentcloud_antiddos_default_alarm_threshold` ([#2342](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2342))
* **New Resource:** `tencentcloud_antiddos_ip_alarm_threshold_config` ([#2342](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2342))
* **New Resource:** `tencentcloud_antiddos_scheduling_domain_user_name` ([#2342](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2342))
* **New Resource:** `tencentcloud_waf_modify_access_period` ([#2340](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2340))

ENHANCEMENTS:

* resource/tencentcloud_kubernetes_addon_attachment: Support `raw_values` and `raw_values_type`. ([#2327](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2327))
* tencentcloud_waf_clb_instance: Support set `bot_management` and `api_security`. ([#2340](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2340))
* tencentcloud_waf_saas_instance: Support set `bot_management` and `api_security`. ([#2340](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2340))

## 1.81.49 (November 22, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_antiddos_overview_cc_trend` ([#2330](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2330))
* **New Data Source:** `tencentcloud_antiddos_overview_ddos_event_list` ([#2330](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2330))
* **New Data Source:** `tencentcloud_antiddos_overview_ddos_trend` ([#2330](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2330))
* **New Data Source:** `tencentcloud_antiddos_overview_index` ([#2330](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2330))
* **New Data Source:** `tencentcloud_antiddos_pending_risk_info` ([#2330](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2330))
* **New Resource:** `tencentcloud_antiddos_ddos_black_white_ip` ([#2330](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2330))
* **New Resource:** `tencentcloud_dasb_bind_device_account_password` ([#2319](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2319))
* **New Resource:** `tencentcloud_dasb_bind_device_account_private_key` ([#2319](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2319))
* **New Resource:** `tencentcloud_dasb_reset_user` ([#2319](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2319))
* **New Resource:** `tencentcloud_dasb_user_group` ([#2319](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2319))
* **New Resource:** `tencentcloud_mysql_database` ([#2332](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2332))
* **New Resource:** `tencentcloud_waf_cc_auto_status` ([#2331](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2331))
* **New Resource:** `tencentcloud_waf_cc_session` ([#2331](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2331))
* **New Resource:** `tencentcloud_waf_ip_access_control` ([#2331](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2331))

ENHANCEMENTS:

* provider: `shared_credentials_dir` Support set Home path, `~/path`. ([#2333](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2333))
* resource/tencentcloud_postgresql_instance: Optimize isolation operations ([#2316](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2316))

## 1.81.48 (November 20, 2023)

ENHANCEMENTS:

* provider: Update some interface input permissions; Add profile authentication ([#2302](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2302))
* resource/tencentcloud_cam_role: Support set `session_duration`. ([#2320](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2320))

## 1.81.47 (November 17, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_dlc_describe_updatable_data_engines` ([#2314](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2314))
* **New Data Source:** `tencentcloud_emr_cvm_quota` ([#2311](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2311))
* **New Resource:** `tencentcloud_dlc_update_data_engine_config_operation` ([#2314](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2314))
* **New Resource:** `tencentcloud_dnspod_domain_lock` ([#2312](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2312))

ENHANCEMENTS:

* resource/tencentcloud_private_dns_zone: Update Field Properties ([#2310](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2310))
* resource/tencentcloud_private_dns_zone_vpc_attachment: Update Field Properties ([#2310](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2310))
* tencentcloud_kubernetes_node_pool: add Deprecated to the security_group_ids parameter description ([#2309](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2309))

## 1.81.46 (November 15, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_waf_user_clb_regions` ([#2301](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2301))
* **New Data Source:** `tencentcloud_wedata_data_source_list` ([#2299](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2299))
* **New Data Source:** `tencentcloud_wedata_data_source_without_info` ([#2299](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2299))
* **New Resource:** `tencentcloud_css_pull_stream_task_restart` ([#2291](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2291))
* **New Resource:** `tencentcloud_css_start_stream_monitor` ([#2291](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2291))
* **New Resource:** `tencentcloud_dnspod_snapshot_config` ([#2306](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2306))
* **New Resource:** `tencentcloud_waf_auto_deny_rules` ([#2298](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2298))
* **New Resource:** `tencentcloud_waf_cc` ([#2301](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2301))
* **New Resource:** `tencentcloud_waf_module_status` ([#2298](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2298))
* **New Resource:** `tencentcloud_waf_protection_mode` ([#2298](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2298))
* **New Resource:** `tencentcloud_waf_web_shell` ([#2298](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2298))
* **New Resource:** `tencentcloud_wedata_datasource` ([#2299](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2299))
* **New Resource:** `tencentcloud_wedata_integration_offline_task` ([#2299](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2299))
* **New Resource:** `tencentcloud_wedata_integration_realtime_task` ([#2299](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2299))
* **New Resource:** `tencentcloud_wedata_integration_task_node` ([#2299](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2299))

ENHANCEMENTS:

* datasource/tencentcloud_scf_function: add async_run_enable ([#2300](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2300))
* resource/tencentcloud_clb_instance: Support set `sla_type` field. ([#2304](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2304))
* resource/tencentcloud_redis_instance: Fields `redis_shard_num`, `redis_replicas_num`, `mem_size` add enumeration values and value checks ([#2305](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2305))
* resource/tencentcloud_scf_function: add async_run_enable ([#2300](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2300))
* tencentcloud_ckafka_datahub_task: Optimize `input_value_type` and `input_value` empty string processing ([#2296](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2296))
* tencentcloud_kubernetes_node_pool: modify the `security_group_ids` in the example plug-in tke nodepool to `orderly_security_group_ids` ([#2297](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2297))
* tencentcloud_waf_clb_domain: Support set cls status ([#2301](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2301))
* tencentcloud_waf_saas_domain: Support set cls status ([#2301](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2301))
* tencentcloud_waf_saas_domain: Support set waf status ([#2298](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2298))

BUG FIXES:

* tencentcloud_mysql_readonly_instance: Fix the problem of error reporting when modifying configuration ([#2303](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2303))

## 1.81.45 (November 10, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_cam_group_user_account` ([#2292](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2292))
* **New Resource:** `tencentcloud_dasb_acl` ([#2286](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2286))
* **New Resource:** `tencentcloud_dasb_bind_device_resource` ([#2286](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2286))
* **New Resource:** `tencentcloud_dasb_cmd_template` ([#2286](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2286))
* **New Resource:** `tencentcloud_dasb_device` ([#2286](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2286))
* **New Resource:** `tencentcloud_dasb_device_account` ([#2286](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2286))
* **New Resource:** `tencentcloud_dasb_device_group` ([#2286](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2286))
* **New Resource:** `tencentcloud_dasb_device_group_members` ([#2286](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2286))
* **New Resource:** `tencentcloud_dasb_resource` ([#2286](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2286))
* **New Resource:** `tencentcloud_dasb_user` ([#2286](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2286))
* **New Resource:** `tencentcloud_dasb_user_group_members` ([#2286](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2286))

ENHANCEMENTS:

* tencentcloud_ccn_bandwidth_limit: Support recover bandwidth limit when run `terraform destroy` ([#2290](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2290))


## 1.81.44 (November 8, 2023)

ENHANCEMENTS:

* resource/tencentcloud_mongodb_instance: support set `security_groups` state. ([#2285](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2285))
* resource/tencentcloud_mongodb_sharding_instance: support set `security_groups` state. ([#2285](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2285))
* tencentcloud_kubernetes_node_pool: `tags` field supports update operation ([#2287](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2287))
* tencentcloud_kubernetes_node_pool: adjust `system_disk_size` validate range ([#2287](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2287))

BUG FIXES:

* tencentcloud_dnspod_record: Fix modify dnspod record remark ([#2288](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2288))

## 1.81.43 (November 6, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_monitor_alarm_notice_callbacks` ([#2282](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2282))
* **New Data Source:** `tencentcloud_monitor_alarm_policy` ([#2282](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2282))
* **New Resource:** `tencentcloud_dnspod_custom_line` ([#2280](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2280))
* **New Resource:** `tencentcloud_dnspod_download_snapshot_operation` ([#2279](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2279))
* **New Resource:** `tencentcloud_dnspod_modify_domain_owner_operation` ([#2279](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2279))
* **New Resource:** `tencentcloud_dnspod_modify_record_group_operation` ([#2279](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2279))
* **New Resource:** `tencentcloud_monitor_alarm_all_namespaces` ([#2282](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2282))
* **New Resource:** `tencentcloud_monitor_alarm_monitor_type` ([#2282](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2282))
* **New Resource:** `tencentcloud_monitor_tmp_regions` ([#2282](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2282))
* **New Resource:** `tencentcloud_wedata_baseline` ([#2260](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2260))
* **New Resource:** `tencentcloud_wedata_dq_rule` ([#2260](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2260))
* **New Resource:** `tencentcloud_wedata_function` ([#2260](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2260))
* **New Resource:** `tencentcloud_wedata_resource` ([#2260](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2260))
* **New Resource:** `tencentcloud_wedata_script` ([#2260](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2260))

ENHANCEMENTS:

* resource/tencentcloud_mongodb_instance: Support `security_groups` (which `engine_version` is `MONGO_40_WT`) ([#2281](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2281))
* resource/tencentcloud_mongodb_sharding_instance: Support `security_groups` (which `engine_version` is `MONGO_40_WT`) ([#2281](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2281))
* resource/tencentcloud_mongodb_standby_instance: Support `security_groups` (which `engine_version` is `MONGO_40_WT`) ([#2281](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2281))
* tencentcloud_dnspod_record: Support set `remark`. ([#2279](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2279))

## 1.81.42 (November 3, 2023)

FEATURES:

* **New Resource:** `tencentcloud_dlc_update_row_filter_operation` ([#2270](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2270))
* **New Resource:** `tencentcloud_dnspod_domain_alias` ([#2273](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2273))

ENHANCEMENTS:

* tencentcloud_instance: Support set create timeout. ([#2276](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2276))
* tencentcloud_mysql_backup_policy: Support instance binlog setting. ([#2272](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2272))
* tencentcloud_nat_gateway: Support create standard nat ([#2275](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2275))

## 1.81.41 (November 1, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_css_backup_stream` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Data Source:** `tencentcloud_css_monitor_report` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Data Source:** `tencentcloud_css_pad_templates` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Data Source:** `tencentcloud_css_pull_stream_task_status` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Data Source:** `tencentcloud_css_stream_monitor_list` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Data Source:** `tencentcloud_css_time_shift_record_detail` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Data Source:** `tencentcloud_css_time_shift_stream_list` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Data Source:** `tencentcloud_css_watermarks` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Data Source:** `tencentcloud_css_xp2p_detail_info_list` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Data Source:** `tencentcloud_dlc_check_data_engine_config_pairs_validity` ([#2259](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2259))
* **New Data Source:** `tencentcloud_elasticsearch_diagnose` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Data Source:** `tencentcloud_elasticsearch_instance_logs` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Data Source:** `tencentcloud_elasticsearch_instance_operations` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Data Source:** `tencentcloud_elasticsearch_instance_plugin_list` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Data Source:** `tencentcloud_elasticsearch_logstash_instance_logs` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Data Source:** `tencentcloud_elasticsearch_logstash_instance_operations` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Data Source:** `tencentcloud_elasticsearch_views` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Data Source:** `tencentcloud_mps_media_meta_data` ([#2263](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2263))
* **New Data Source:** `tencentcloud_mps_parse_live_stream_process_notification` ([#2263](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2263))
* **New Data Source:** `tencentcloud_mps_parse_notification` ([#2263](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2263))
* **New Resource:** `tencentcloud_css_backup_stream` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_callback_rule_attachment` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_callback_template` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_domain_referer` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_enable_optimal_switching` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_pad_rule_attachment` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_pad_template` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_record_rule_attachment` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_snapshot_rule_attachment` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_snapshot_template` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_stream_monitor` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_timeshift_rule_attachment` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_css_timeshift_template` ([#2265](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2265))
* **New Resource:** `tencentcloud_dlc_bind_work_groups_to_user_attachment` ([#2259](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2259))
* **New Resource:** `tencentcloud_dlc_restart_data_engine_operation` ([#2257](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2257))
* **New Resource:** `tencentcloud_dlc_switch_data_engine_image_operation` ([#2257](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2257))
* **New Resource:** `tencentcloud_dlc_upgrade_data_engine_image_operation` ([#2257](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2257))
* **New Resource:** `tencentcloud_dlc_user_data_engine_config` ([#2259](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2259))
* **New Resource:** `tencentcloud_dnspod_record_group` ([#2249](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2249))
* **New Resource:** `tencentcloud_elasticsearch_diagnose` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Resource:** `tencentcloud_elasticsearch_diagnose_instance` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Resource:** `tencentcloud_elasticsearch_index` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Resource:** `tencentcloud_elasticsearch_restart_instance_operation` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Resource:** `tencentcloud_elasticsearch_restart_kibana_operation` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Resource:** `tencentcloud_elasticsearch_restart_nodes_operation` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Resource:** `tencentcloud_elasticsearch_update_plugins_operation` ([#2264](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2264))
* **New Resource:** `tencentcloud_gaap_global_domain` ([#2267](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2267))
* **New Resource:** `tencentcloud_gaap_global_domain_dns` ([#2267](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2267))
* **New Resource:** `tencentcloud_mps_process_media_operation` ([#2263](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2263))

ENHANCEMENTS:

* tencentcloud_dlc_store_location_config: Support api `ModifyAdvancedStoreLocation` and `DescribeAdvancedStoreLocation`, replace `CreateStoreLocation` and `DescribeStoreLocation` ([#2262](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2262))
* tencentcloud_vpn_ssl_server: Support modify attributes ([#2266](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2266))

BUG FIXES:

* tencentcloud_dlc_describe_user_info: Fix policy set type ([#2259](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2259))
* tencentcloud_dlc_describe_work_group_info: Fix policy set type ([#2259](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2259))

## 1.81.40 (October 27, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_dlc_describe_data_engine` ([#2250](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2250))
* **New Data Source:** `tencentcloud_dlc_describe_data_engine_image_versions` ([#2250](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2250))
* **New Data Source:** `tencentcloud_dlc_describe_data_engine_python_spark_images` ([#2250](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2250))
* **New Data Source:** `tencentcloud_dlc_describe_engine_usage_info` ([#2250](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2250))
* **New Data Source:** `tencentcloud_dlc_describe_work_group_info` ([#2250](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2250))
* **New Data Source:** `tencentcloud_oceanus_check_savepoint` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Data Source:** `tencentcloud_oceanus_clusters` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Data Source:** `tencentcloud_oceanus_job_submission_log` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Data Source:** `tencentcloud_oceanus_resource_related_job` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Data Source:** `tencentcloud_oceanus_savepoint_list` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Data Source:** `tencentcloud_oceanus_system_resource` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Data Source:** `tencentcloud_oceanus_tree_jobs` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Data Source:** `tencentcloud_oceanus_tree_resources` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Data Source:** `tencentcloud_oceanus_work_spaces` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Resource:** `tencentcloud_dlc_modify_data_engine_description_operation` ([#2251](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2251))
* **New Resource:** `tencentcloud_dlc_modify_user_typ_operation` ([#2251](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2251))
* **New Resource:** `tencentcloud_dlc_renew_data_engine_operation` ([#2251](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2251))
* **New Resource:** `tencentcloud_oceanus_job` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Resource:** `tencentcloud_oceanus_job_config` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Resource:** `tencentcloud_oceanus_job_copy` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Resource:** `tencentcloud_oceanus_resource` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Resource:** `tencentcloud_oceanus_resource_config` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Resource:** `tencentcloud_oceanus_run_job` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Resource:** `tencentcloud_oceanus_stop_job` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Resource:** `tencentcloud_oceanus_trigger_job_savepoint` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))
* **New Resource:** `tencentcloud_oceanus_work_space` ([#2224](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2224))

ENHANCEMENTS:

* resource/tencentcloud_kubernetes_scale_worker: set up the node and import it into Terraform ([#2246](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2246))
* tencentcloud_ssm_secret: Remove default version while `secret_type` is 0 ([#2256](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2256))
* tencentcloud_tse_cngw_gateway: Add computed `instance_port.tcp_port` and `instance_port.udp_port` ([#2255](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2255))

BUG FIXES:

* resource/tencentcloud_kubernetes_cluster_attachment: fix the null pointer issue in data_disk ([#2246](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2246))

## 1.81.39 (October 25, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_bi_project` ([#2239](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2239))
* **New Data Source:** `tencentcloud_bi_user_project` ([#2239](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2239))
* **New Data Source:** `tencentcloud_dlc_check_data_engine_image_can_be_upgrade` ([#2245](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2245))
* **New Data Source:** `tencentcloud_dlc_describe_user_info` ([#2244](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2244))
* **New Data Source:** `tencentcloud_dlc_describe_user_roles` ([#2244](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2244))
* **New Data Source:** `tencentcloud_dlc_describe_user_type` ([#2244](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2244))
* **New Resource:** `tencentcloud_bi_datasource` ([#2239](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2239))
* **New Resource:** `tencentcloud_bi_datasource_cloud` ([#2239](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2239))
* **New Resource:** `tencentcloud_bi_embed_interval_apply` ([#2239](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2239))
* **New Resource:** `tencentcloud_bi_embed_token_apply` ([#2239](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2239))
* **New Resource:** `tencentcloud_bi_project` ([#2239](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2239))
* **New Resource:** `tencentcloud_bi_project_user_role` ([#2239](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2239))
* **New Resource:** `tencentcloud_bi_user_role` ([#2239](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2239))
* **New Resource:** `tencentcloud_cdwpg_instance` ([#2248](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2248))
* **New Resource:** `tencentcloud_dlc_data_engine` ([#2245](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2245))
* **New Resource:** `tencentcloud_dlc_rollback_data_engine_image_operation` ([#2245](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2245))
* **New Resource:** `tencentcloud_organization_org_member_policy_attachment` ([#2243](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2243))

ENHANCEMENTS:

* resource/tencentcloud_waf_clb_instance: support set `qps_limit` ([#2234](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2234))
* resource/tencentcloud_waf_saas_instance: support set `qps_limit` ([#2234](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2234))
* tencentcloud_cam_service_linked_role: Set `custom_suffix` computed ([#2247](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2247))

BUG FIXES:

* data-source/tencentcloud_dnspod_record_analytics: Fix set info error. ([#2242](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2242))

## 1.81.38 (October 20, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_dlc_check_data_engine_image_can_be_rollback` ([#2233](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2233))
* **New Data Source:** `tencentcloud_dnspod_domain_log_list` ([#2235](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2235))
* **New Data Source:** `tencentcloud_dnspod_record_analytics` ([#2235](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2235))
* **New Data Source:** `tencentcloud_dnspod_record_line_list` ([#2235](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2235))
* **New Data Source:** `tencentcloud_dnspod_record_list` ([#2235](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2235))
* **New Data Source:** `tencentcloud_dnspod_record_type` ([#2235](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2235))
* **New Data Source:** `tencentcloud_gaap_check_proxy_create` ([#2238](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2238))
* **New Data Source:** `tencentcloud_gaap_domain_error_page_infos` ([#2238](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2238))
* **New Data Source:** `tencentcloud_gaap_group_and_statistics_proxy` ([#2238](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2238))
* **New Data Source:** `tencentcloud_gaap_listener_real_servers` ([#2238](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2238))
* **New Data Source:** `tencentcloud_gaap_listener_statistics` ([#2238](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2238))
* **New Data Source:** `tencentcloud_gaap_proxies_status` ([#2238](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2238))
* **New Data Source:** `tencentcloud_gaap_proxy_and_statistics_listeners` ([#2238](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2238))
* **New Data Source:** `tencentcloud_gaap_region_and_price` ([#2238](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2238))
* **New Data Source:** `tencentcloud_gaap_resources_by_tag` ([#2238](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2238))
* **New Data Source:** `tencentcloud_gaap_rule_real_servers` ([#2238](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2238))
* **New Data Source:** `tencentcloud_monitor_alarm_basic_alarms` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Data Source:** `tencentcloud_monitor_alarm_basic_metric` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Data Source:** `tencentcloud_monitor_alarm_conditions_template` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Data Source:** `tencentcloud_monitor_alarm_history` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Data Source:** `tencentcloud_monitor_alarm_metric` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Data Source:** `tencentcloud_monitor_grafana_plugin_overviews` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Resource:** `tencentcloud_dlc_add_users_to_work_group_attachment` ([#2233](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2233))
* **New Resource:** `tencentcloud_dlc_store_location_config` ([#2233](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2233))
* **New Resource:** `tencentcloud_dlc_suspend_resume_data_engine` ([#2240](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2240))
* **New Resource:** `tencentcloud_monitor_alarm_policy_set_default` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Resource:** `tencentcloud_monitor_grafana_dns_config` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Resource:** `tencentcloud_monitor_grafana_env_config` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Resource:** `tencentcloud_monitor_grafana_sso_cam_config` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Resource:** `tencentcloud_monitor_grafana_sso_config` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Resource:** `tencentcloud_monitor_grafana_version_upgrade` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Resource:** `tencentcloud_monitor_grafana_whitelist_config` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))
* **New Resource:** `tencentcloud_organization_org_identity` ([#2237](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2237))

ENHANCEMENTS:

* tencentcloud_cam_service_linked_role: Support `import` function and Fix issues when update service role ([#2236](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2236))
* tencentcloud_monitor_grafana_instance: Support modifying `enable_internet` ([#2212](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2212))

BUG FIXES:

* resource/tencentcloud_cam_group: fix cam group modify error ([#2232](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2232))

## 1.81.37 (October 18, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_cam_account_summary` ([#2225](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2225))
* **New Data Source:** `tencentcloud_cam_list_attached_user_policy` ([#2220](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2220))
* **New Data Source:** `tencentcloud_cam_oidc_config` ([#2227](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2227))
* **New Data Source:** `tencentcloud_cam_policy_granting_service_access` ([#2225](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2225))
* **New Data Source:** `tencentcloud_cam_secret_last_used_time` ([#2225](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2225))
* **New Data Source:** `tencentcloud_dnspod_domain_analytics` ([#2221](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2221))
* **New Data Source:** `tencentcloud_dnspod_domain_list` ([#2218](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2218))
* **New Data Source:** `tencentcloud_kms_list_algorithms` ([#2222](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2222))
* **New Data Source:** `tencentcloud_kms_white_box_decrypt_key` ([#2222](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2222))
* **New Data Source:** `tencentcloud_kms_white_box_device_fingerprints` ([#2222](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2222))
* **New Resource:** `tencentcloud_cam_role_permission_boundary_attachment` ([#2226](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2226))
* **New Resource:** `tencentcloud_cam_set_policy_version_config` ([#2225](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2225))
* **New Resource:** `tencentcloud_cam_tag_role_attachment` ([#2220](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2220))
* **New Resource:** `tencentcloud_elasticsearch_logstash` ([#2228](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2228))
* **New Resource:** `tencentcloud_elasticsearch_logstash_pipeline` ([#2228](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2228))
* **New Resource:** `tencentcloud_elasticsearch_restart_logstash_instance_operation` ([#2228](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2228))
* **New Resource:** `tencentcloud_elasticsearch_start_logstash_pipeline_operation` ([#2228](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2228))
* **New Resource:** `tencentcloud_elasticsearch_stop_logstash_pipeline_operation` ([#2228](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2228))
* **New Resource:** `tencentcloud_kms_cloud_resource_attachment` ([#2222](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2222))
* **New Resource:** `tencentcloud_kms_overwrite_white_box_device_fingerprints` ([#2222](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2222))
* **New Resource:** `tencentcloud_organization_quit_organization_operation` ([#2229](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2229))

ENHANCEMENTS:

* resource/tencentcloud_kubernetes_auth_attachment: Support OIDC config. ([#2227](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2227))
* tencentcloud_ckafka_instance: Support api `UpdateRoleConsoleLogin` ([#2220](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2220))

## 1.81.36 (October 13, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_api_gateway_api_app_api` ([#2204](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2204))
* **New Data Source:** `tencentcloud_api_gateway_api_plugins` ([#2204](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2204))
* **New Data Source:** `tencentcloud_api_gateway_bind_api_apps_status` ([#2204](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2204))
* **New Data Source:** `tencentcloud_api_gateway_service_environment_list` ([#2204](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2204))
* **New Data Source:** `tencentcloud_api_gateway_service_release_versions` ([#2204](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2204))
* **New Data Source:** `tencentcloud_cam_list_entities_for_policy` ([#2211](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2211))
* **New Data Source:** `tencentcloud_mps_tasks` ([#2214](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2214))
* **New Data Source:** `tencentcloud_organization_org_financial_by_month` ([#2206](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2206))
* **New Data Source:** `tencentcloud_organization_org_financial_by_product` ([#2206](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2206))
* **New Resource:** `tencentcloud_api_gateway_import_open_api` ([#2204](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2204))
* **New Resource:** `tencentcloud_api_gateway_update_api_app_key` ([#2204](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2204))
* **New Resource:** `tencentcloud_cam_policy_version` ([#2211](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2211))
* **New Resource:** `tencentcloud_emr_user_manager` ([#2208](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2208))
* **New Resource:** `tencentcloud_mps_content_review_template` ([#2214](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2214))
* **New Resource:** `tencentcloud_mps_input` ([#2199](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2199))
* **New Resource:** `tencentcloud_mps_output` ([#2199](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2199))
* **New Resource:** `tencentcloud_mps_start_flow_operation` ([#2199](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2199))

ENHANCEMENTS:

* resource/tencentcloud_cfw_edge_policy: Update code logic ([#2207](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2207))
* resource/tencentcloud_cfw_nat_policy: Update code logic ([#2207](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2207))
* resource/tencentcloud_dts_sync_config: Support `database_net_env` field when the access type is ccn. ([#2201](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2201))
* tencentcloud_ckafka_instance: support param `upgrade_strategy` and postpaid scaling down ([#2209](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2209))
* tencentcloud_organization_org_member: support api `UpdateOrganizationMember` ([#2206](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2206))
* tencentcloud_tse_cngw_gateway: Add computed `instance_port` and `public_ip_addresses` ([#2210](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2210))

## 1.81.35 (October 11, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_cwp_machines_simple` ([#2186](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2186))
* **New Data Source:** `tencentcloud_mps_schedules` ([#2185](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2185))
* **New Data Source:** `tencentcloud_organization_org_auth_node` ([#2191](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2191))
* **New Data Source:** `tencentcloud_organization_org_financial_by_member` ([#2198](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2198))
* **New Data Source:** `tencentcloud_pts_scenario_with_jobs` ([#2193](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2193))
* **New Data Source:** `tencentcloud_ssl_describe_certificate` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_companies` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_api_gateway_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_cdn_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_clb_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_cos_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_ddos_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_deploy_record` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_deploy_record_detail` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_lighthouse_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_live_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_teo_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_tke_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_update_record` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_update_record_detail` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_vod_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_host_waf_instance_list` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_manager_detail` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Data Source:** `tencentcloud_ssl_describe_managers` ([#2173](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2173))
* **New Resource:** `tencentcloud_cam_access_key` ([#2182](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2182))
* **New Resource:** `tencentcloud_cwp_license_bind_attachment` ([#2186](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2186))
* **New Resource:** `tencentcloud_cwp_license_order` ([#2186](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2186))
* **New Resource:** `tencentcloud_mps_edit_media_operation` ([#2188](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2188))
* **New Resource:** `tencentcloud_mps_event` ([#2192](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2192))
* **New Resource:** `tencentcloud_mps_execute_function_operation` ([#2196](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2196))
* **New Resource:** `tencentcloud_mps_flow` ([#2192](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2192))
* **New Resource:** `tencentcloud_mps_manage_task_operation` ([#2196](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2196))
* **New Resource:** `tencentcloud_mps_process_live_stream_operation` ([#2188](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2188))
* **New Resource:** `tencentcloud_organization_instance` ([#2198](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2198))
* **New Resource:** `tencentcloud_organization_org_member_auth_identity_attachment` ([#2191](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2191))
* **New Resource:** `tencentcloud_organization_org_member_email` ([#2198](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2198))
* **New Resource:** `tencentcloud_pts_cron_job_abort` ([#2193](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2193))
* **New Resource:** `tencentcloud_pts_cron_job_restart` ([#2193](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2193))
* **New Resource:** `tencentcloud_pts_job_abort` ([#2193](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2193))
* **New Resource:** `tencentcloud_rum_instance_status_config` ([#2117](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2117))
* **New Resource:** `tencentcloud_rum_project_status_config` ([#2117](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2117))
* **New Resource:** `tencentcloud_ssl_check_certificate_chain_operation` ([#2174](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2174))
* **New Resource:** `tencentcloud_ssl_complete_certificate_operation` ([#2174](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2174))
* **New Resource:** `tencentcloud_ssl_deploy_certificate_instance_operation` ([#2174](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2174))
* **New Resource:** `tencentcloud_ssl_deploy_certificate_record_retry_operation` ([#2174](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2174))
* **New Resource:** `tencentcloud_ssl_deploy_certificate_record_rollback_operation` ([#2174](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2174))
* **New Resource:** `tencentcloud_ssl_download_certificate_operation` ([#2174](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2174))
* **New Resource:** `tencentcloud_ssl_replace_certificate_operation` ([#2175](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2175))
* **New Resource:** `tencentcloud_ssl_revoke_certificate_operation` ([#2175](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2175))
* **New Resource:** `tencentcloud_ssl_update_certificate_instance_operation` ([#2175](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2175))
* **New Resource:** `tencentcloud_ssl_update_certificate_record_retry_operation` ([#2175](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2175))
* **New Resource:** `tencentcloud_ssl_update_certificate_record_rollback_operation` ([#2175](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2175))
* **New Resource:** `tencentcloud_ssl_upload_revoke_letter_operation` ([#2175](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2175))

ENHANCEMENTS:

* resource/tencentcloud_instance: Add more data disk types ([#2189](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2189))
* resource/tencentcloud_kubernetes_cluster: Adjust the TKE doc example which refers to the security group of `auto_scaling_config`. ([#2197](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2197))
* resource/tencentcloud_ssm_product_secret: Optimize code structure ([#2190](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2190))
* resource/tencentcloud_ssm_rotate_product_secret: Optimize code structure ([#2190](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2190))
* resource/tencentcloud_ssm_secret: Optimize code structure ([#2190](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2190))
* resource/tencentcloud_ssm_ssh_key_pair_secret: Optimize code structure ([#2190](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2190))

## 1.81.34 (October 9, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_eb_plateform_event_template` ([#2181](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2181))
* **New Data Source:** `tencentcloud_eb_platform_event_names` ([#2180](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2180))
* **New Data Source:** `tencentcloud_eb_platform_event_patterns` ([#2180](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2180))
* **New Data Source:** `tencentcloud_eb_platform_products` ([#2180](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2180))
* **New Data Source:** `tencentcloud_gaap_access_regions` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Data Source:** `tencentcloud_gaap_access_regions_by_dest_region` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Data Source:** `tencentcloud_gaap_black_header` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Data Source:** `tencentcloud_gaap_country_area_mapping` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Data Source:** `tencentcloud_gaap_custom_header` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Data Source:** `tencentcloud_gaap_dest_regions` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Data Source:** `tencentcloud_gaap_proxy_detail` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Data Source:** `tencentcloud_gaap_proxy_group_statistics` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Data Source:** `tencentcloud_gaap_proxy_groups` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Data Source:** `tencentcloud_gaap_proxy_statistics` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Data Source:** `tencentcloud_gaap_real_servers_status` ([#2183](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2183))
* **New Resource:** `tencentcloud_cam_user_permission_boundary_attachment` ([#2177](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2177))
* **New Resource:** `tencentcloud_mps_enable_schedule_config` ([#2179](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2179))
* **New Resource:** `tencentcloud_mps_schedule` ([#2179](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2179))

## 1.81.33 (October 8, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_cat_metric_data` ([#2143](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2143))
* **New Data Source:** `tencentcloud_kms_describe_keys` ([#2171](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2171))
* **New Data Source:** `tencentcloud_kms_list_keys` ([#2171](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2171))
* **New Data Source:** `tencentcloud_kms_white_box_key_details` ([#2171](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2171))
* **New Data Source:** `tencentcloud_ssm_rotation_detail` ([#2172](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2172))
* **New Data Source:** `tencentcloud_ssm_rotation_history` ([#2172](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2172))
* **New Data Source:** `tencentcloud_ssm_service_status` ([#2172](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2172))
* **New Data Source:** `tencentcloud_ssm_ssh_key_pair_value` ([#2172](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2172))
* **New Resource:** `tencentcloud_kms_white_box_key` ([#2171](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2171))
* **New Resource:** `tencentcloud_mps_word_sample` ([#2176](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2176))
* **New Resource:** `tencentcloud_ssm_rotate_product_secret` ([#2172](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2172))

ENHANCEMENTS:

* datasource/tencentcloud_cat_node: Add computed `task_types` ([#2143](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2143))
* resource/tencentcloud_cat_task_set: Support `suspend` and `resume` for dial test tasks ([#2143](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2143))
* resource/tencentcloud_cfw_sync_asset: Update operation status query interface ([#2169](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2169))

BUG FIXES:

* resource/tencentcloud_clb_attachment: fix http backend delete problem. ([#2148](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2148))
* resource/tencentcloud_kubernetes_node_pool: fix bug in node pool creation ([#2170](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2170))

## 1.81.32 (September 28, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_cfw_edge_fw_switches` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Data Source:** `tencentcloud_cfw_nat_fw_switches` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Data Source:** `tencentcloud_cfw_vpc_fw_switches` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Data Source:** `tencentcloud_waf_instance_qps_limit` ([#2160](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2160))
* **New Resource:** `tencentcloud_cfw_address_template` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_block_ignore` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_edge_firewall_switch` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_edge_policy` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_nat_firewall_switch` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_nat_instance` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_nat_policy` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_sync_asset` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_sync_route` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_vpc_firewall_switch` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_vpc_instance` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_cfw_vpc_policy` ([#2149](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2149))
* **New Resource:** `tencentcloud_mps_withdraws_watermark_operation` ([#2163](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2163))
* **New Resource:** `tencentcloud_pts_tmp_key_generate` ([#2139](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2139))
* **New Resource:** `tencentcloud_teo_acceleration_domain` ([#2154](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2154))
* **New Resource:** `tencentcloud_teo_certificate_config` ([#2154](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2154))
* **New Resource:** `tencentcloud_teo_ownership_verify` ([#2154](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2154))
* **New Resource:** `tencentcloud_waf_anti_fake` ([#2160](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2160))
* **New Resource:** `tencentcloud_waf_anti_info_leak` ([#2160](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2160))

ENHANCEMENTS:

* resource/tencentcloud_kubernetes_node_pool: set up the node pool and import it into Terraform ([#2164](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2164))
* resource/tencentcloud_teo_zone: Support create zone with existing plan ([#2154](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2154))
* resource/tencentcloud_waf_clb_domain: alb_type support apisix, tsegw ([#2160](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2160))

## 1.81.31 (September 27, 2023)

FEATURES:

* **New Resource:** `tencentcloud_css_watermark_rule_attachment` ([#2156](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2156))
* **New Resource:** `tencentcloud_mps_enable_workflow_config` ([#2159](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2159))
* **New Resource:** `tencentcloud_trocket_rocketmq_role` ([#2153](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2153))

ENHANCEMENTS:

* resource/tencentcloud_trocket_rocketmq_role: Add computed `access_key` and `secret_key` ([#2158](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2158))

BUG FIXES:

* resource/tencentcloud_as_scaling_config: fix `key_ids` can not modify ([#2155](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2155))

## 1.81.30 (September 25, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_tse_gateway_certificates` ([#2147](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2147))
* **New Resource:** `tencentcloud_tse_cngw_certificate` ([#2147](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2147))

ENHANCEMENTS:

* resource/tencentcloud_vpc_bandwidth_package: support set `egress` ([#2146](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2146))

## 1.81.29 (September 22, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_waf_attack_log_histogram` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Data Source:** `tencentcloud_waf_attack_log_list` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Data Source:** `tencentcloud_waf_attack_overview` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Data Source:** `tencentcloud_waf_attack_total_count` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Data Source:** `tencentcloud_waf_ciphers` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Data Source:** `tencentcloud_waf_domains` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Data Source:** `tencentcloud_waf_find_domains` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Data Source:** `tencentcloud_waf_peak_points` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Data Source:** `tencentcloud_waf_ports` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Data Source:** `tencentcloud_waf_tls_versions` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Data Source:** `tencentcloud_waf_user_domains` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Resource:** `tencentcloud_ses_black_list_delete` ([#2131](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2131))
* **New Resource:** `tencentcloud_waf_clb_domain` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Resource:** `tencentcloud_waf_clb_instance` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Resource:** `tencentcloud_waf_custom_rule` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Resource:** `tencentcloud_waf_custom_white_rule` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Resource:** `tencentcloud_waf_saas_domain` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))
* **New Resource:** `tencentcloud_waf_saas_instance` ([#2111](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2111))

## 1.81.28 (September 20, 2023)

FEATURES:

* **New Resource:** `tencentcloud_private_dns_zone_vpc_attachment` ([#2136](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2136))

BUG FIXES:

* resource/tencentcloud_tdmq_rocketmq_group: Fix certificate import issue ([#2133](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2133))
* resource/tencentcloud_trocket_rocketmq_instance: fix `enable_public` read problem ([#2134](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2134))
* resource/tencentcloud_vpc_acl: Fix vpc acl entry inconsistent problem while port is `ALL`. ([#2135](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2135))

## 1.81.27 (September 18, 2023)

FEATURES:

* **New Resource:** `tencentcloud_trocket_rocketmq_consumer_group` ([#2123](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2123))
* **New Resource:** `tencentcloud_trocket_rocketmq_topic` ([#2123](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2123))
* **New Resource:** `tencentcloud_tse_cngw_group` ([#2126](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2126))

ENHANCEMENTS:

* resource/tencentcloud_ckafka_instance: support create postpaid instance ([#2130](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2130))
* resource/tencentcloud_eip: support set network egress ([#2129](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2129))
* resource/tencentcloud_clb_listener: support create port range listener ([#2127](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2127))
* resource/tencentcloud_cynosdb_readonly_instance: Support `vpc_id`, `subnet_id` fields. ([#2125](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2125))
* resource/tencentcloud_ses_receiver: Support import. ([#2124](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2124))
* resource/tencentcloud_tse_cngw_gateway: Support modifying `node_config`. ([#2126](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2126))

## 1.81.26 (September 13, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_clickhouse_backup_tables` ([#2120](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2120))
* **New Resource:** `tencentcloud_clickhouse_account` ([#2120](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2120))
* **New Resource:** `tencentcloud_clickhouse_account_permission` ([#2120](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2120))
* **New Resource:** `tencentcloud_trocket_rocketmq_instance` ([#2119](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2119))

## 1.81.25 (September 8, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_ses_black_email_address` ([#2112](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2112))
* **New Data Source:** `tencentcloud_ses_email_identities` ([#2112](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2112))
* **New Data Source:** `tencentcloud_ses_receivers` ([#2112](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2112))
* **New Data Source:** `tencentcloud_ses_send_email_status` ([#2112](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2112))
* **New Data Source:** `tencentcloud_ses_send_tasks` ([#2112](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2112))
* **New Data Source:** `tencentcloud_ses_statistics_report` ([#2112](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2112))
* **New Resource:** `tencentcloud_cam_mfa_flag` ([#2115](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2115))
* **New Resource:** `tencentcloud_ses_batch_send_email` ([#2112](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2112))
* **New Resource:** `tencentcloud_ses_receiver` ([#2112](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2112))
* **New Resource:** `tencentcloud_ses_send_email` ([#2112](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2112))
* **New Resource:** `tencentcloud_ses_verify_domain` ([#2112](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2112))

ENHANCEMENTS:

* resource/tencentcloud_ssl_pay_certificate: Supports `csr_type` modification ([#2115](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2115))
* resource/tencentcloud_tcr_service_account: update the description of the `permissions.actions` field. ([#2116](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2116))

## 1.81.24 (September 6, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_kms_get_parameters_for_import` ([#2109](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2109))
* **New Resource:** `tencentcloud_ssl_commit_certificate_information` ([#2105](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2105))

ENHANCEMENTS:

* resource/tencentcloud_ssl_pay_certificate: support parameter 'wait_commit_flag'. ([#2105](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2105))

## 1.81.23 (September 1, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_api_gateway_api_app_service` ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* **New Data Source:** `tencentcloud_api_gateway_api_usage_plans` ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* **New Data Source:** `tencentcloud_api_gateway_plugins` ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* **New Data Source:** `tencentcloud_api_gateway_upstreams` ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* **New Data Source:** `tencentcloud_clickhouse_backup_job_detail` ([#2098](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2098))
* **New Data Source:** `tencentcloud_clickhouse_backup_jobs` ([#2098](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2098))
* **New Data Source:** `tencentcloud_eb_event_rules` ([#2094](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2094))
* **New Data Source:** `tencentcloud_kms_public_key` ([#2093](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2093))
* **New Data Source:** `tencentcloud_private_dns_records` ([#2095](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2095))
* **New Data Source:** `tencentcloud_tse_gateway_routes` ([#2096](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2096))
* **New Data Source:** `tencentcloud_tse_gateways` ([#2096](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2096))
* **New Data Source:** `tencentcloud_tse_groups` ([#2096](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2096))
* **New Resource:** `tencentcloud_api_gateway_api_app_attachment` ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* **New Resource:** `tencentcloud_api_gateway_upstream` ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* **New Resource:** `tencentcloud_clickhouse_backup` ([#2098](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2098))
* **New Resource:** `tencentcloud_clickhouse_backup_strategy` ([#2098](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2098))
* **New Resource:** `tencentcloud_clickhouse_delete_backup_data` ([#2098](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2098))
* **New Resource:** `tencentcloud_clickhouse_recover_backup_job` ([#2098](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2098))
* **New Resource:** `tencentcloud_tse_cngw_gateway` ([#2096](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2096))
* **New Resource:** `tencentcloud_tse_cngw_route` ([#2096](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2096))
* **New Resource:** `tencentcloud_tse_cngw_route_rate_limit` ([#2096](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2096))
* **New Resource:** `tencentcloud_tse_cngw_service_rate_limit` ([#2096](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2096))

ENHANCEMENTS:

* data_source/tencentcloud_ssm_secrets: Support `secret_type`, `product_name`, Update return value ([#2075](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2075))
* resource/tencentcloud_api_gateway_api: Support Some new fields ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* resource/tencentcloud_api_gateway_api_app: Support `tag` ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* resource/tencentcloud_api_gateway_api_key: Support input `access_key_id`, `access_key_secret` ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* resource/tencentcloud_api_gateway_plugin_attachment: Optimization binding failure blocking issue ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* resource/tencentcloud_api_gateway_service: Support input `uniq_vpc_id`` ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* resource/tencentcloud_api_gateway_usage_plan_attachment: Support input `access_key_ids`` ([#2077](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2077))
* resource/tencentcloud_kubernetes_cluster: support for cluster configuration of CiliumOverlay network ([#2099](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2099))
* resource/tencentcloud_monitor_alarm_notice: support set 'is_valid' and 'validation_code'. ([#2102](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2102))
* resource/tencentcloud_private_dns_zone: support input 'cname_speedup_status'. ([#2101](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2101))
* resource/tencentcloud_scf_function: support `handler` and `runtime` optional ([#2091](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2091))
* resource/tencentcloud_scf_trigger_config: support new para, update api UpdateTrigger ([#2071](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2071))
* resource/tencentcloud_ssm_product_secret: Support `tags`, Update return value ([#2075](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2075))
* resource/tencentcloud_ssm_secret: Support `tags` ([#2075](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2075))
* resource/tencentcloud_ssm_ssh_key_pair_secret: Support `tags` ([#2075](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2075))
* resource/tencentcloud_tke_tmp_alert_policy: support modify 'alert_rule'. ([#2103](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2103))
* resource/tencentcloud_tse_cngw_service: Deprecate ineffective tags ([#2096](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2096))

BUG FIXES:

* resource/tencentcloud_kubernetes_cluster: fix DiffSuppressFunc would ignore docker_graph_path defaults when computing diffs ([#2090](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2090))
* resource/tencentcloud_tdmq_rocketmq_group: Fix the issue of inconsistent return values. ([#2087](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2087))

## 1.81.22 (August 25, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_wedata_rule_templates` ([#2067](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2067))
* **New Resource:** `tencentcloud_lighthouse_firewall_template` ([#2076](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2076))
* **New Resource:** `tencentcloud_wedata_rule_template` ([#2067](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2067))

ENHANCEMENTS:

* datasource/tencentcloud_monitor_alarm_notices: Support `amp_consumer_id` field ([#2081](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2081))
* resource/tencentcloud_kubernetes_node_pool: add `orderly_security_group_ids` field to set the security group orderly and deprecated the `security_group_ids` field. ([#2070](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2070))
* resource/tencentcloud_lighthouse_instance: support param `firewall_template_id` ([#2076](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2076))
* resource/tencentcloud_tdmq_rocketmq_environment_role: change `cluster_id`, `environment_name` and `role_name` to `ForceNew` ([#2078](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2078))
* resource/tencentcloud_tdmq_rocketmq_group: change `cluster_id`, `namespaces` and `group_id` to `ForceNew` ([#2078](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2078))
* resource/tencentcloud_teo_zone: `cname_speed_up` cannot be modified ([#2074](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2074))
* resource/tencentcloud_vpn_gateway: change `zone` to `optional` ([#2080](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2080))

BUG FIXES:

* resource/tencentcloud_vpn_connection: support change `ike_version` ([#2082](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2082))

## 1.81.21 (August 18, 2023)

FEATURES:

* **New Resource:** `tencentcloud_dlc_user` ([#2061](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2061))
* **New Resource:** `tencentcloud_dlc_work_group` ([#2061](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2061))
* **New Resource:** `tencentcloud_eb_event_connector` ([#2060](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2060))

ENHANCEMENTS:

* resource/tencentcloud_as_scaling_config: Support field `host_name` ([#2052](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2052))
* resource/tencentcloud_ccn_attachment: support `route_ids` fields ([#2053](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2053))
* resource/tencentcloud_cvm_reboot_instance: Deprecated `force_reboot` ([#2054](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2054))
* resource/tencentcloud_dayu_l7_rules_v2: fix query results. ([#2062](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2062))
* resource/tencentcloud_emr_cluster: Update `tags` logic ([#2059](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2059))
* resource/tencentcloud_eni_sg_attachment: Update Document Description ([#2049](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2049))
* resource/tencentcloud_nat_gateway: Optimize execution time ([#2058](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2058))

## 1.81.20 (August 11, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_eb_bus` ([#2039](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2039))
* **New Resource:** `tencentcloud_cls_data_transform` ([#2034](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2034))
* **New Resource:** `tencentcloud_cls_kafka_recharge` ([#2034](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2034))
* **New Resource:** `tencentcloud_cls_scheduled_sql` ([#2034](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2034))
* **New Resource:** `tencentcloud_eb_event_bus` ([#2039](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2039))
* **New Resource:** `tencentcloud_eb_event_rule` ([#2039](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2039))
* **New Resource:** `tencentcloud_eb_event_target` ([#2039](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2039))
* **New Resource:** `tencentcloud_eb_event_transform` ([#2039](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2039))

ENHANCEMENTS:

* data_source/tencentcloud_sqlserver_backups: Supports filtering queries by `backup_name` parameter ([#2023](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2023))
* resource/tencentcloud_clb_instance: `master_zone_id`, `slave_zone_id`, `project_id`, `vpc_id`, `subnet_id`, `address_ip_version`, `bandwidth_package_id`, `snat_pro`, `zone_id` fields cannot be modified and an error message is displayed ([#2036](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2036))
* resource/tencentcloud_clb_instance: support param `delete_protect`. ([#2037](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2037))
* resource/tencentcloud_eni: Update `security_groups` field ([#2047](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2047))
* resource/tencentcloud_sqlserver_full_backup_migration: return `backup_migration_id` parameter ([#2023](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2023))
* resource/tencentcloud_sqlserver_publish_subscribe: Update `database_tuples` field ([#2040](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2040))
* resource/tencentcloud_sqlserver_readonly_instance: Extension creation parameters ([#2023](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2023))
* resource/tencentcloud_tdmq_namespace: Update `retention_policy` Parameters ([#2033](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2033))
* resource/tencentcloud_tdmq_rocketmq_namespace: Abandoned parameters: 'ttl', 'retention_time`. ([#2029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2029))
* resource/tencentcloud_tdmq_rocketmq_vip_instance: Cancel import ([#2033](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2033))
* resource/tencentcloud_tdmq_rocketmq_vip_instance: Optimize `storage_size` parameter issue ([#2029](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2029))

## 1.81.19 (August 04, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_ssm_products` ([#2027](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2027))
* **New Resource:** `tencentcloud_ssm_product_secret` ([#2027](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2027))
* **New Resource:** `tencentcloud_tdmq_rocketmq_vip_instance` ([#2011](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2011))

ENHANCEMENTS:

* data_source/tencentcloud_clb_attachments: Optimize the problem display of target parameters ([#2003](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2003))
* datasource/tencentcloud_mysql_instance: Support query `group_id` ([#2002](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2002))
* resource/tencentcloud_cynosdb_cluster_slave_zone: Optimize timeout ([#1998](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1998))
* resource/tencentcloud_cynosdb_proxy: Add return value `ro_instances` ([#2008](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2008))
* resource/tencentcloud_emr_cluster: support `tags` param ([#2026](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2026))
* resource/tencentcloud_mongodb_instance: support params `node_num`, `availability_zone_list` and `hidden_zone` ([#2013](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2013))
* resource/tencentcloud_security_group_rule_set: Optimize Rule Delete Logic ([#2018](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2018))
* resource/tencentcloud_sqlserver_general_cloud_ro_instance: Support setting timeout ([#2016](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2016))

## 1.81.18 (July 31, 2023)

FEATURES:

* **New Resource:** `tencentcloud_tag` ([#1991](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1991))
* **New Resource:** `tencentcloud_tag_attachment` ([#1991](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1991))
* **New Resource:** `tencentcloud_tdmq_rabbitmq_vip_instance` ([#1984](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1984))

ENHANCEMENTS:

* datasource/tencentcloud_mysql_zone_config: Support querying `cpu`, `info`, `device_type` ([#1996](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1996))
* datasource/tencentcloud_ssl_certificates: Support querying `dv_auths` and `order_id` ([#1993](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1993))
* resource/tencentcloud_eip_association: Optimize the problem while binding EIP timeout. ([#1990](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1990))
* resource/tencentcloud_tcr_service_account: update example ([#1995](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1995))
* resource/tencentcloud_tcr_tag_retention_rule: allow `cron_setting` to be modified ([#1995](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1995))
* resource/tencentcloud_vpc_bandwidth_package: support set `time_span` ([#2001](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/2001))

## 1.81.17 (July 27, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_tse_gateway_services` ([#1975](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1975))
* **New Resource:** `tencentcloud_ckafka_route` ([#1983](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1983))
* **New Resource:** `tencentcloud_redis_connection_config` ([#1982](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1982))
* **New Resource:** `tencentcloud_tse_cngw_canary_rule` ([#1975](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1975))
* **New Resource:** `tencentcloud_tse_cngw_service` ([#1975](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1975))

ENHANCEMENTS:

* datasource/tencentcloud_mysql_parameter_list: support mysql8.0 ([#1982](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1982))
* resource/tencentcloud_dcx: support set `dc_owner_account` ([#1988](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1988))

BUG FIXES:

* resource/tencentcloud_instance: fix order `data_disks` ([#1985](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1985))

## 1.81.16 (July 21, 2023)

ENHANCEMENTS:

* resource/tencentcloud_ckafka_acl: remove `host` limit, update `principal` description ([#1973](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1973))
* resource/tencentcloud_eip: support create AntiDDoS EIP ([#1977](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1977))
* resource/tencentcloud_monitor_grafana_instance: Deprecated `is_distroy`, use `is_destroy` instead of `is_distroy` ([#1976](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1976))
* resource/tencentcloud_vpc_bandwidth_package: support set more `charge_type` ([#1981](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1981))

BUG FIXES:

* resource/tencentcloud_cynosdb_cluster: fix panic error while parameter not found in template ([#1981](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1981))

## 1.81.15 (July 19, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_tse_gateway_canary_rules` ([#1964](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1964))
* **New Resource:** `tencentcloud_clickhouse_instance` ([#1642](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1642))
* **New Resource:** `tencentcloud_mysql_isolate_instance` ([#1961](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1961))

ENHANCEMENTS:

* resource/tencentcloud_monitor_alarm_notice: Support return field `amp_consumer_id` ([#1971](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1971))
* resource/tencentcloud_monitor_grafana_instance: Support for cleaning up deactivated instances ([#1971](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1971))
* resource/tencentcloud_redis_instance: Non-multi-AZ supports modifying the number of replicas ([#1969](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1969))
* resource/tencentcloud_security_group_rule_set: Optimize `service_template_id` usage issue ([#1970](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1970))

BUG FIXES:

* resource/tencentcloud_instance: fix data-disks order issue ([#1968](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1968))

## 1.81.14 (July 12, 2023)

FEATURES:

* **New Resource:** `tencentcloud_cynosdb_upgrade_proxy_version` ([#1936](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1936))
* **New Resource:** `tencentcloud_elasticsearch_security_group` ([#1956](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1956))

ENHANCEMENTS:

* resource/tencentcloud_clb_instance: support set `dynamic_vip` and export `domain` ([#1954](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1954))
* resource/tencentcloud_cos_bucket: remove `multi_az` limit ([#1957](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1957))
* resource/tencentcloud_mysql_readonly_instance: support import ([#1958](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1958))
* resource/tencentcloud_nat_gateway: support set `zone` ([#1959](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1959))

## 1.81.13 (July 10, 2023)

FEATURES:

* **New Resource:** `tencentcloud_kubernetes_encryption_protection` ([#1949](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1949))
* **New Resource:** `tencentcloud_tcr_service_account` ([#1943](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1943))

## 1.81.12 (July 07, 2023)

FEATURES:

* **New Resource:** `tencentcloud_ciam_user_group` ([#1946](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1946))
* **New Resource:** `tencentcloud_ciam_user_store` ([#1946](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1946))
* **New Resource:** `tencentcloud_cos_bucket_generate_inventory_immediately_operation` ([#1941](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1941))
* **New Resource:** `tencentcloud_cos_object_abort_multipart_upload_operation` ([#1941](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1941))
* **New Resource:** `tencentcloud_cos_object_copy_operation` ([#1941](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1941))
* **New Resource:** `tencentcloud_cos_object_download_operation` ([#1947](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1947))
* **New Resource:** `tencentcloud_cos_object_restore_operation` ([#1941](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1941))
* **New Resource:** `tencentcloud_sqlserver_general_cloud_ro_instance` ([#1942](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1942))

## 1.81.11 (July 05, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_cos_batchs` ([#1928](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1928))
* **New Data Source:** `tencentcloud_cos_bucket_inventorys` ([#1928](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1928))
* **New Data Source:** `tencentcloud_cos_bucket_multipart_uploads` ([#1928](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1928))
* **New Resource:** `tencentcloud_as_complete_lifecycle` ([#1937](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1937))
* **New Resource:** `tencentcloud_cos_batch` ([#1928](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1928))
* **New Resource:** `tencentcloud_cynosdb_proxy_end_point` ([#1927](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1927))
* **New Resource:** `tencentcloud_ssm_ssh_key_pair_secret` ([#1929](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1929))

ENHANCEMENTS:

* resource/tencentcloud_cfs_file_system: support create turbo file system ([#1934](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1934))
* resource/tencentcloud_ci_media_concat_template: adapt to new cos sdk ([#1928](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1928))
* resource/tencentcloud_ci_media_transcode_template: adapt to new cos sdk ([#1928](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1928))
* resource/tencentcloud_ci_media_video_montage_template: adapt to new cos sdk ([#1928](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1928))
* resource/tencentcloud_cos_bucket: adjust `acl_body` to fit COS AccessControlPolicy sequence ([#1924](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1924))
* resource/tencentcloud_ssm_secret: support `service_type` and `additional_config` ([#1931](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1931))

## 1.81.10 (June 30, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_cynosdb_instance_slow_queries` ([#1915](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1915))
* **New Data Source:** `tencentcloud_cynosdb_proxy_node` ([#1912](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1912))
* **New Data Source:** `tencentcloud_cynosdb_proxy_version` ([#1912](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1912))
* **New Data Source:** `tencentcloud_sqlserver_ins_attribute` ([#1917](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1917))
* **New Data Source:** `tencentcloud_sqlserver_query_xevent` ([#1917](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1917))
* **New Data Source:** `tencentcloud_tse_gateway_nodes` ([#1913](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1913))
* **New Resource:** `tencentcloud_cynosdb_cluster_slave_zone` ([#1915](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1915))
* **New Resource:** `tencentcloud_cynosdb_proxy` ([#1912](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1912))
* **New Resource:** `tencentcloud_cynosdb_read_only_instance_exclusive_access` ([#1915](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1915))
* **New Resource:** `tencentcloud_cynosdb_reload_proxy_node` ([#1912](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1912))
* **New Resource:** `tencentcloud_mysql_proxy` ([#1911](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1911))
* **New Resource:** `tencentcloud_mysql_reload_balance_proxy_node` ([#1911](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1911))
* **New Resource:** `tencentcloud_mysql_reset_root_account` ([#1911](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1911))
* **New Resource:** `tencentcloud_mysql_ro_start_replication` ([#1911](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1911))
* **New Resource:** `tencentcloud_mysql_ro_stop_replication` ([#1911](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1911))
* **New Resource:** `tencentcloud_mysql_verify_root_account` ([#1911](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1911))
* **New Resource:** `tencentcloud_redis_security_group_attachment` ([#1914](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1914))
* **New Resource:** `tencentcloud_sqlserver_database_tde` ([#1917](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1917))
* **New Resource:** `tencentcloud_sqlserver_instance_tde` ([#1917](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1917))
* **New Resource:** `tencentcloud_sqlserver_start_xevent` ([#1917](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1917))
* **New Resource:** `tencentcloud_tdmq_rabbitmq_user` ([#1922](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1922))
* **New Resource:** `tencentcloud_tdmq_rabbitmq_virtual_host` ([#1922](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1922))
* **New Resource:** `tencentcloud_tdmq_send_rocketmq_message` ([#1922](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1922))

ENHANCEMENTS:

* data-source/tencentcloud_api_gateway_service: deprecated the `exclusive_set_name` field. ([#1917](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1917))
* resource/tencentcloud_api_gateway_service: support create `instance_id` field(Exclusive instance ID); deprecated the `exclusive_set_name` field. ([#1917](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1917))
* resource/tencentcloud_cynosdb_cluster: Support `cluster_name, storage_limit, vpc_id, subnet_id, old_ip_reserve_hours` modification ([#1909](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1909))
* resource/tencentcloud_mysql_account: support update `host` ([#1911](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1911))
* resource/tencentcloud_mysql_instance: Support modifying field `slave_deploy_mode`, `first_slave_zone`, `second_slave_zone`, `slave_sync_mode`. ([#1873](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1873))
* resource/tencentcloud_private_dns_record: catch error while record is not exist ([#1920](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1920))
* resource/tencentcloud_sqlserver_restore_instance: support create `encryption` ([#1917](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1917))
* resource/tencentcloud_sqlserver_rollback_instance: support create `encryption` ([#1917](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1917))
* resource/tencentcloud_vpn_gateway: support create `SSL_CCN` vpn gateway ([#1919](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1919))

## 1.81.9 (June 21, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_cynosdb_account_all_grant_privileges` ([#1904](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1904))
* **New Data Source:** `tencentcloud_cynosdb_audit_logs` ([#1903](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1903))
* **New Data Source:** `tencentcloud_cynosdb_binlog_download_url` ([#1903](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1903))
* **New Data Source:** `tencentcloud_cynosdb_cluster` ([#1903](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1903))
* **New Data Source:** `tencentcloud_cynosdb_cluster_detail_databases` ([#1903](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1903))
* **New Data Source:** `tencentcloud_cynosdb_cluster_param_logs` ([#1903](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1903))
* **New Data Source:** `tencentcloud_cynosdb_describe_instance_error_logs` ([#1903](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1903))
* **New Data Source:** `tencentcloud_cynosdb_describe_instance_slow_queries` ([#1903](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1903))
* **New Data Source:** `tencentcloud_cynosdb_project_security_groups` ([#1904](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1904))
* **New Data Source:** `tencentcloud_cynosdb_resource_package_list` ([#1904](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1904))
* **New Data Source:** `tencentcloud_cynosdb_resource_package_sale_specs` ([#1904](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1904))
* **New Data Source:** `tencentcloud_cynosdb_rollback_time_range` ([#1904](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1904))
* **New Data Source:** `tencentcloud_cynosdb_zone` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Data Source:** `tencentcloud_tdmq_pro_instance_detail` ([#1902](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1902))
* **New Data Source:** `tencentcloud_tdmq_pro_instances` ([#1902](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1902))
* **New Data Source:** `tencentcloud_tdmq_rocketmq_messages` ([#1902](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1902))
* **New Resource:** `tencentcloud_cos_bucket_inventory` ([#1907](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1907))
* **New Resource:** `tencentcloud_cynosdb_account` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_account_privileges` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_binlog_save_days` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_cluster_databases` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_cluster_password_complexity` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_cluster_resource_packages_attachment` ([#1906](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1906))
* **New Resource:** `tencentcloud_cynosdb_export_instance_error_logs` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_export_instance_slow_queries` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_instance_param` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_isolate_instance` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_param_template` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_restart_instance` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_roll_back_cluster` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))
* **New Resource:** `tencentcloud_cynosdb_wan` ([#1905](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1905))

ENHANCEMENTS:

* resource/tencentcloud_cos_bucket: support `enable_intelligent_tiering` and format client called ([#1907](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1907))
* resource/tencentcloud_cos_bucket_domain_certificate_attachment: format client called ([#1907](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1907))
* resource/tencentcloud_cos_bucket_object: format client called ([#1907](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1907))

## 1.81.8 (June 16, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_lighthouse_disks` ([#1895](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1895))
* **New Data Source:** `tencentcloud_tcr_replication_instance_create_tasks` ([#1900](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1900))
* **New Data Source:** `tencentcloud_tcr_replication_instance_sync_status` ([#1900](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1900))
* **New Data Source:** `tencentcloud_tcr_tag_retention_executions` ([#1900](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1900))
* **New Resource:** `tencentcloud_as_execute_scaling_policy` ([#1897](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1897))
* **New Resource:** `tencentcloud_as_scale_in_instances` ([#1897](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1897))
* **New Resource:** `tencentcloud_as_scale_out_instances` ([#1897](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1897))
* **New Resource:** `tencentcloud_as_scaling_group_status` ([#1897](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1897))
* **New Resource:** `tencentcloud_cos_bucket_referer` ([#1894](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1894))
* **New Resource:** `tencentcloud_cos_bucket_version` ([#1894](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1894))
* **New Resource:** `tencentcloud_tcr_tag_retention_execution_config` ([#1896](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1896))

ENHANCEMENTS:

* resource/tencentcloud_api_gateway_service: support param `tags`. ([#1898](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1898))
* resource/tencentcloud_lighthouse_instance: offline param `permit_default_key_pair_login` ([#1899](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1899))
* resource/tencentcloud_tcm_mesh: support `inject` and `sidecar_resources`. ([#1893](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1893))
* resource/tencentcloud_tcr_instance: support to update `instance_charge_type_prepaid_period` field. ([#1896](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1896))
* resource/tencentcloud_tcr_tag_retention_rule: add a computed field `retention_id`. ([#1896](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1896))

## 1.81.7 (June 14, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_dbbrain_diag_db_instances` ([#1886](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1886))
* **New Data Source:** `tencentcloud_dbbrain_mysql_process_list` ([#1886](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1886))
* **New Data Source:** `tencentcloud_dbbrain_no_primary_key_tables` ([#1886](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1886))
* **New Data Source:** `tencentcloud_dbbrain_redis_top_big_keys` ([#1886](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1886))
* **New Data Source:** `tencentcloud_dbbrain_redis_top_key_prefix_list` ([#1886](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1886))
* **New Data Source:** `tencentcloud_dcdb_file_download_url` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_instance_node_info` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_log_files` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_orders` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_price` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_project_security_groups` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_projects` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_renewal_price` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_sale_info` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_shard_spec` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_slow_logs` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_dcdb_upgrade_price` ([#1882](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1882))
* **New Data Source:** `tencentcloud_projects` ([#1891](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1891))
* **New Data Source:** `tencentcloud_tsf_container_group` ([#1889](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1889))
* **New Data Source:** `tencentcloud_tsf_groups` ([#1889](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1889))
* **New Data Source:** `tencentcloud_tsf_ms_api_list` ([#1889](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1889))
* **New Resource:** `tencentcloud_project` ([#1891](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1891))
* **New Resource:** `tencentcloud_redis_backup_operation` ([#1875](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1875))
* **New Resource:** `tencentcloud_redis_replicate_attachment` ([#1875](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1875))
* **New Resource:** `tencentcloud_tsf_deploy_container_group` ([#1868](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1868))
* **New Resource:** `tencentcloud_tsf_deploy_vm_group` ([#1868](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1868))
* **New Resource:** `tencentcloud_tsf_operate_container_group` ([#1868](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1868))
* **New Resource:** `tencentcloud_tsf_operate_group` ([#1868](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1868))
* **New Resource:** `tencentcloud_tsf_release_api_group` ([#1868](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1868))
* **New Resource:** `tencentcloud_tsf_unit_namespace` ([#1868](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1868))

ENHANCEMENTS:

* resource/tencentcloud_cos_bucket: Update the validation rules for the `bucket` parameter ([#1885](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1885))
* resource/tencentcloud_mysql_instance: adjust `slave_deploy_mode`, `first_slave_zone`, `second_slave_zone` and `slave_sync_mode` to unchangeable. ([#1892](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1892))
* resource/tencentcloud_postgresql_instance: adjust `availability_zone` to unchangeable. ([#1888](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1888))
* resource/tencentcloud_security_group_rule_set: Add default values for the protocol and port fields. ([#1885](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1885))

BUG FIXES:

* resource/tencentcloud_instance: change key_ids null value ([#1883](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1883))

## 1.81.6 (June 09, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_as_advices` ([#1867](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1867))
* **New Data Source:** `tencentcloud_as_last_activity` ([#1867](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1867))
* **New Data Source:** `tencentcloud_as_limits` ([#1867](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1867))
* **New Data Source:** `tencentcloud_chdfs_file_systems` ([#1870](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1870))
* **New Data Source:** `tencentcloud_mariadb_dcn_detail` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_file_download_url` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_flow` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_instance_specs` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_log_files` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_orders` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_price` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_project_security_groups` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_renewal_price` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_sale_info` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_slow_logs` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mariadb_upgrade_price` ([#1877](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1877))
* **New Data Source:** `tencentcloud_mysql_databases` ([#1871](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1871))
* **New Data Source:** `tencentcloud_mysql_error_log` ([#1871](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1871))
* **New Data Source:** `tencentcloud_mysql_project_security_group` ([#1871](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1871))
* **New Data Source:** `tencentcloud_mysql_ro_min_scale` ([#1871](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1871))
* **New Data Source:** `tencentcloud_vpc_limits` ([#1862](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1862))
* **New Data Source:** `tencentcloud_vpc_net_detect_state_check` ([#1866](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1866))
* **New Data Source:** `tencentcloud_vpc_security_group_limits` ([#1862](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1862))
* **New Data Source:** `tencentcloud_vpc_security_group_references` ([#1862](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1862))
* **New Data Source:** `tencentcloud_vpc_sg_snapshot_file_content` ([#1862](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1862))
* **New Data Source:** `tencentcloud_vpc_snapshot_files` ([#1862](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1862))
* **New Data Source:** `tencentcloud_vpc_subnet_resource_dashboard` ([#1862](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1862))
* **New Data Source:** `tencentcloud_vpc_template_limits` ([#1862](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1862))
* **New Data Source:** `tencentcloud_vpc_used_ip_address` ([#1862](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1862))
* **New Data Source:** `tencentcloud_vpn_default_health_check_ip` ([#1880](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1880))
* **New Resource:** `tencentcloud_ccn_instances_reject_attach` ([#1880](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1880))
* **New Resource:** `tencentcloud_dcdb_activate_hour_instance_operation` ([#1865](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1865))
* **New Resource:** `tencentcloud_dcdb_cancel_dcn_job_operation` ([#1865](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1865))
* **New Resource:** `tencentcloud_dcdb_db_sync_mode_config` ([#1863](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1863))
* **New Resource:** `tencentcloud_dcdb_encrypt_attributes_config` ([#1863](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1863))
* **New Resource:** `tencentcloud_dcdb_flush_binlog_operation` ([#1872](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1872))
* **New Resource:** `tencentcloud_dcdb_instance_config` ([#1872](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1872))
* **New Resource:** `tencentcloud_dcdb_isolate_hour_instance_operation` ([#1865](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1865))
* **New Resource:** `tencentcloud_dcdb_switch_db_instance_ha_operation` ([#1872](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1872))
* **New Resource:** `tencentcloud_eni_sg_attachment` ([#1866](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1866))
* **New Resource:** `tencentcloud_mariadb_account_privileges` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* **New Resource:** `tencentcloud_mariadb_backup_time` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* **New Resource:** `tencentcloud_mariadb_cancel_dcn_job` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* **New Resource:** `tencentcloud_mariadb_flush_binlog` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* **New Resource:** `tencentcloud_mariadb_instance_config` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* **New Resource:** `tencentcloud_mariadb_renew_instance` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* **New Resource:** `tencentcloud_mariadb_restart_instance` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* **New Resource:** `tencentcloud_postgresql_base_backup` ([#1803](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1803))
* **New Resource:** `tencentcloud_sqlserver_start_backup_full_migration` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* **New Resource:** `tencentcloud_sqlserver_start_backup_incremental_migration` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))

ENHANCEMENTS:

* resource/tencentcloud_ckafka_instance: deprecated the `public_network` field. ([#1876](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1876))
* resource/tencentcloud_dcdb_account: support to update the `password` field. ([#1869](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1869))
* resource/tencentcloud_dcdb_db_instance: optimize the processing logic of dcn. ([#1865](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1865))
* resource/tencentcloud_dcdb_hourdb_instance: support `dcn_region` and `dcn_instance_id` fields. ([#1865](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1865))
* resource/tencentcloud_dcdb_hourdb_instance: support `rs_access_strategy`, `vip` and `vipv6` field. ([#1872](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1872))
* resource/tencentcloud_dcdb_hourdb_instance: support the `extranet_access` field. ([#1869](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1869))
* resource/tencentcloud_dcdb_hourdb_instance: support to update the `project_id` field. ([#1869](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1869))
* resource/tencentcloud_lighthouse_instance: support param `isolate_data_disk`. ([#1879](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1879))
* resource/tencentcloud_mariadb_account: support to update `password` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* resource/tencentcloud_mariadb_dedicatedcluster_db_instance: support to update `project_id`, `vip` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* resource/tencentcloud_mariadb_hour_db_instance: support to update `project_id`, `vip` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* resource/tencentcloud_mariadb_instance: support to update `project_id`, `vip` ([#1864](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1864))
* resource/tencentcloud_mysql_account: support the `max_user_connections` field. ([#1874](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1874))

## 1.81.5 (June 02, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_dts_migrate_db_instances` ([#1845](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1845))
* **New Data Source:** `tencentcloud_tdmq_environment_attributes` ([#1844](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1844))
* **New Data Source:** `tencentcloud_tdmq_publisher_summary` ([#1844](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1844))
* **New Data Source:** `tencentcloud_tdmq_publishers` ([#1844](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1844))
* **New Data Source:** `tencentcloud_tdmq_rabbitmq_node_list` ([#1844](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1844))
* **New Data Source:** `tencentcloud_tdmq_rabbitmq_vip_instance` ([#1844](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1844))
* **New Data Source:** `tencentcloud_tdmq_vip_instance` ([#1844](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1844))
* **New Data Source:** `tencentcloud_tsf_gateway_all_group_apis` ([#1855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1855))
* **New Data Source:** `tencentcloud_tsf_group_config_release` ([#1855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1855))
* **New Data Source:** `tencentcloud_tsf_group_gateways` ([#1855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1855))
* **New Data Source:** `tencentcloud_tsf_group_instances` ([#1855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1855))
* **New Data Source:** `tencentcloud_tsf_pod_instances` ([#1855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1855))
* **New Data Source:** `tencentcloud_tsf_usable_unit_namespaces` ([#1855](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1855))
* **New Data Source:** `tencentcloud_vpc_account_attributes` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Data Source:** `tencentcloud_vpc_classic_link_instances` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Data Source:** `tencentcloud_vpc_cvm_instances` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Data Source:** `tencentcloud_vpc_gateway_flow_monitor_detail` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Data Source:** `tencentcloud_vpc_gateway_flow_qos` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Data Source:** `tencentcloud_vpc_net_detect_states` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Data Source:** `tencentcloud_vpc_network_interface_limit` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Data Source:** `tencentcloud_vpc_private_ip_addresses` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Data Source:** `tencentcloud_vpc_product_quota` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Data Source:** `tencentcloud_vpc_resource_dashboard` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Data Source:** `tencentcloud_vpc_route_conflicts` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Resource:** `tencentcloud_dts_sync_config` ([#1845](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1845))
* **New Resource:** `tencentcloud_dts_sync_job_continue_operation` ([#1857](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1857))
* **New Resource:** `tencentcloud_dts_sync_job_isolate_operation` ([#1857](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1857))
* **New Resource:** `tencentcloud_dts_sync_job_pause_operation` ([#1857](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1857))
* **New Resource:** `tencentcloud_dts_sync_job_recover_operation` ([#1857](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1857))
* **New Resource:** `tencentcloud_dts_sync_job_resize_operation` ([#1857](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1857))
* **New Resource:** `tencentcloud_dts_sync_job_start_operation` ([#1857](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1857))
* **New Resource:** `tencentcloud_dts_sync_job_stop_operation` ([#1857](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1857))
* **New Resource:** `tencentcloud_mysql_ro_group` ([#1856](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1856))
* **New Resource:** `tencentcloud_mysql_ro_group_load_operation` ([#1856](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1856))
* **New Resource:** `tencentcloud_mysql_ro_instance_ip` ([#1856](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1856))
* **New Resource:** `tencentcloud_mysql_rollback` ([#1856](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1856))
* **New Resource:** `tencentcloud_mysql_rollback_stop` ([#1856](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1856))
* **New Resource:** `tencentcloud_mysql_switch_for_upgrade` ([#1856](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1856))
* **New Resource:** `tencentcloud_mysql_switch_master_slave_operation` ([#1856](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1856))
* **New Resource:** `tencentcloud_scf_trigger_config` ([#1853](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1853))
* **New Resource:** `tencentcloud_vpc_enable_end_point_connect` ([#1861](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1861))
* **New Resource:** `tencentcloud_vpc_ipv6_cidr_block` ([#1860](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1860))
* **New Resource:** `tencentcloud_vpc_ipv6_eni_address` ([#1860](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1860))
* **New Resource:** `tencentcloud_vpc_ipv6_subnet_cidr_block` ([#1860](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1860))
* **New Resource:** `tencentcloud_vpc_local_gateway` ([#1860](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1860))
* **New Resource:** `tencentcloud_vpc_resume_snapshot_instance` ([#1860](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1860))

ENHANCEMENTS:

* resource/tencentcloud_cynosdb_cluster: adjust the delete logic ([#1845](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1845))

BUG FIXES:

* resource/tencentcloud_lighthouse_disk: add create ready status ([#1858](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1858))

## 1.81.4 (May 31, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_ckafka_zone` ([#1850](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1850))
* **New Resource:** `tencentcloud_ckafka_acl_rule` ([#1850](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1850))
* **New Resource:** `tencentcloud_ckafka_consumer_group` ([#1850](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1850))
* **New Resource:** `tencentcloud_ckafka_consumer_group_modify_offset` ([#1850](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1850))
* **New Resource:** `tencentcloud_ckafka_datahub_task` ([#1850](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1850))
* **New Resource:** `tencentcloud_ckafka_renew_instance` ([#1850](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1850))
* **New Resource:** `tencentcloud_cvm_export_images` ([#1828](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1828))
* **New Resource:** `tencentcloud_cvm_image_share_permission` ([#1828](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1828))

## 1.81.3 (May 31, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_ckafka_group_info` ([#1835](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1835))
* **New Data Source:** `tencentcloud_ckafka_task_status` ([#1835](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1835))
* **New Data Source:** `tencentcloud_ckafka_topic_flow_ranking` ([#1835](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1835))
* **New Data Source:** `tencentcloud_ckafka_topic_produce_connection` ([#1835](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1835))
* **New Data Source:** `tencentcloud_ckafka_topic_subscribe_group` ([#1835](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1835))
* **New Data Source:** `tencentcloud_ckafka_topic_sync_replica` ([#1835](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1835))
* **New Data Source:** `tencentcloud_clb_cluster_resources` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_clb_cross_targets` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_clb_exclusive_clusters` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_clb_idle_instances` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_clb_instance_by_cert_id` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_clb_instance_detail` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_clb_instance_traffic` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_clb_listeners_by_targets` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_clb_resources` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_clb_target_group_list` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_clb_target_health` ([#1838](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1838))
* **New Data Source:** `tencentcloud_postgresql_backup_download_urls` ([#1839](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1839))
* **New Data Source:** `tencentcloud_postgresql_base_backups` ([#1839](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1839))
* **New Data Source:** `tencentcloud_postgresql_db_instance_classes` ([#1839](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1839))
* **New Data Source:** `tencentcloud_postgresql_db_instance_versions` ([#1839](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1839))
* **New Data Source:** `tencentcloud_postgresql_default_parameters` ([#1839](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1839))
* **New Data Source:** `tencentcloud_postgresql_log_backups` ([#1839](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1839))
* **New Data Source:** `tencentcloud_postgresql_recovery_time` ([#1839](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1839))
* **New Data Source:** `tencentcloud_postgresql_regions` ([#1839](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1839))
* **New Data Source:** `tencentcloud_postgresql_zones` ([#1839](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1839))
* **New Data Source:** `tencentcloud_scf_account_info` ([#1836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1836))
* **New Data Source:** `tencentcloud_scf_async_event_management` ([#1836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1836))
* **New Data Source:** `tencentcloud_scf_async_event_status` ([#1836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1836))
* **New Data Source:** `tencentcloud_scf_function_address` ([#1836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1836))
* **New Data Source:** `tencentcloud_scf_function_aliases` ([#1836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1836))
* **New Data Source:** `tencentcloud_scf_function_versions` ([#1836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1836))
* **New Data Source:** `tencentcloud_scf_layer_versions` ([#1836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1836))
* **New Data Source:** `tencentcloud_scf_layers` ([#1836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1836))
* **New Data Source:** `tencentcloud_scf_request_status` ([#1836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1836))
* **New Data Source:** `tencentcloud_scf_triggers` ([#1836](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1836))
* **New Resource:** `tencentcloud_clb_security_group_attachment` ([#1840](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1840))
* **New Resource:** `tencentcloud_mysql_backup_encryption_status` ([#1830](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1830))
* **New Resource:** `tencentcloud_mysql_dr_instance_to_mater` ([#1830](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1830))
* **New Resource:** `tencentcloud_mysql_instance_encryption_operation` ([#1830](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1830))
* **New Resource:** `tencentcloud_mysql_password_complexity` ([#1830](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1830))
* **New Resource:** `tencentcloud_mysql_remote_backup_config` ([#1830](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1830))
* **New Resource:** `tencentcloud_mysql_restart_db_instances_operation` ([#1830](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1830))
* **New Resource:** `tencentcloud_postgresql_modify_account_remark_operation` ([#1837](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1837))
* **New Resource:** `tencentcloud_postgresql_modify_switch_time_period_operation` ([#1837](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1837))
* **New Resource:** `tencentcloud_scf_invoke_function` ([#1833](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1833))
* **New Resource:** `tencentcloud_scf_provisioned_concurrency_config` ([#1833](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1833))
* **New Resource:** `tencentcloud_scf_sync_invoke_function` ([#1833](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1833))
* **New Resource:** `tencentcloud_scf_terminate_async_event` ([#1833](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1833))
* **New Resource:** `tencentcloud_vpc_dhcp_ip` ([#1846](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1846))

ENHANCEMENTS:

* resource/tencentcloud_cam_user: support update `password` ([#1847](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1847))
* resource/tencentcloud_elasticsearch_instance: disk_type supports `CLOUD_HSSD` type ([#1849](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1849))
* resource/tencentcloud_postgresql_instance: support to update `charge_type` and `period` ([#1837](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1837))
* resource/tencentcloud_postgresql_instance: support to update `vpc_id` and `subnet_id` ([#1813](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1813))
* resource/tencentcloud_postgresql_readonly_group: support to update `vpc_id` and `subnet_id` ([#1813](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1813))
* resource/tencentcloud_postgresql_readonly_instance: support `private_access_ip` and `private_access_port` fields ([#1813](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1813))
* resource/tencentcloud_postgresql_readonly_instance: support the `read_only_group_id` field ([#1837](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1837))
* resource/tencentcloud_scf_function: support set `func_type` ([#1851](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1851))

## 1.81.2 (May 26, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_ckafka_connect_resource` ([#1809](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1809))
* **New Data Source:** `tencentcloud_ckafka_datahub_group_offsets` ([#1827](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1827))
* **New Data Source:** `tencentcloud_ckafka_datahub_task` ([#1827](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1827))
* **New Data Source:** `tencentcloud_ckafka_datahub_topic` ([#1827](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1827))
* **New Data Source:** `tencentcloud_ckafka_group` ([#1827](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1827))
* **New Data Source:** `tencentcloud_ckafka_group_offsets` ([#1827](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1827))
* **New Data Source:** `tencentcloud_ckafka_region` ([#1809](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1809))
* **New Data Source:** `tencentcloud_cls_machine_group_configs` ([#1818](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1818))
* **New Data Source:** `tencentcloud_cls_machines` ([#1818](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1818))
* **New Data Source:** `tencentcloud_cls_shipper_tasks` ([#1817](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1817))
* **New Data Source:** `tencentcloud_cvm_image_quota` ([#1802](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1802))
* **New Data Source:** `tencentcloud_cvm_image_share_permission` ([#1802](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1802))
* **New Data Source:** `tencentcloud_cvm_import_image_os` ([#1802](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1802))
* **New Data Source:** `tencentcloud_dbbrain_top_space_schema_time_series` ([#1796](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1796))
* **New Data Source:** `tencentcloud_dbbrain_top_space_schemas` ([#1796](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1796))
* **New Data Source:** `tencentcloud_dbbrain_top_space_table_time_series` ([#1796](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1796))
* **New Data Source:** `tencentcloud_dbbrain_top_space_tables` ([#1796](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1796))
* **New Data Source:** `tencentcloud_dc_access_points` ([#1826](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1826))
* **New Data Source:** `tencentcloud_dc_gateway_attachment` ([#1822](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1822))
* **New Data Source:** `tencentcloud_dc_internet_address_config` ([#1822](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1822))
* **New Data Source:** `tencentcloud_dc_internet_address_quota` ([#1823](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1823))
* **New Data Source:** `tencentcloud_dc_internet_address_statistics` ([#1823](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1823))
* **New Data Source:** `tencentcloud_dc_public_direct_connect_tunnel_routes` ([#1823](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1823))
* **New Data Source:** `tencentcloud_lighthouse_all_scene` ([#1797](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1797))
* **New Data Source:** `tencentcloud_lighthouse_modify_instance_bundle` ([#1797](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1797))
* **New Data Source:** `tencentcloud_mysql_db_features` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_inst_tables` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_instance_charset` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_instance_info` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_instance_param_record` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_instance_reboot_time` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_rollback_range_time` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_slow_log` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_slow_log_data` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_supported_privileges` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_switch_record` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_mysql_user_task` ([#1793](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1793))
* **New Data Source:** `tencentcloud_postgresql_readonly_groups` ([#1798](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1798))
* **New Data Source:** `tencentcloud_sqlserver_backup_by_flow_id` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_sqlserver_backup_upload_size` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_sqlserver_cross_region_zone` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_sqlserver_db_charsets` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_sqlserver_instance_param_records` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_sqlserver_project_security_groups` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_sqlserver_regions` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_sqlserver_rollback_time` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_sqlserver_slowlogs` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_sqlserver_upload_backup_info` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_sqlserver_upload_incremental_info` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Data Source:** `tencentcloud_tsf_repository` ([#1804](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1804))
* **New Resource:** `tencentcloud_clb_instance_mix_ip_target_config` ([#1790](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1790))
* **New Resource:** `tencentcloud_clb_instance_sla_config` ([#1790](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1790))
* **New Resource:** `tencentcloud_clb_replace_cert_for_lbs` ([#1790](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1790))
* **New Resource:** `tencentcloud_cls_alarm` ([#1808](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1808))
* **New Resource:** `tencentcloud_cls_alarm_notice` ([#1808](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1808))
* **New Resource:** `tencentcloud_cls_ckafka_consumer` ([#1817](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1817))
* **New Resource:** `tencentcloud_cls_cos_recharge` ([#1817](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1817))
* **New Resource:** `tencentcloud_cls_export` ([#1817](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1817))
* **New Resource:** `tencentcloud_cvm_renew_instance` ([#1802](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1802))
* **New Resource:** `tencentcloud_cvm_sync_image` ([#1802](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1802))
* **New Resource:** `tencentcloud_dc_instance` ([#1826](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1826))
* **New Resource:** `tencentcloud_dc_internet_address` ([#1822](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1822))
* **New Resource:** `tencentcloud_dc_share_dcx_config` ([#1822](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1822))
* **New Resource:** `tencentcloud_dcx_extra_config` ([#1822](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1822))
* **New Resource:** `tencentcloud_lighthouse_renew_disk` ([#1797](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1797))
* **New Resource:** `tencentcloud_lighthouse_renew_instance` ([#1797](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1797))
* **New Resource:** `tencentcloud_mysql_backup_download_restriction` ([#1812](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1812))
* **New Resource:** `tencentcloud_mysql_renew_db_instance_operation` ([#1812](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1812))
* **New Resource:** `tencentcloud_postgresql_backup_download_restriction_config` ([#1806](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1806))
* **New Resource:** `tencentcloud_postgresql_backup_plan_config` ([#1798](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1798))
* **New Resource:** `tencentcloud_postgresql_delete_log_backup_operation` ([#1820](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1820))
* **New Resource:** `tencentcloud_postgresql_disisolate_db_instance_operation` ([#1820](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1820))
* **New Resource:** `tencentcloud_postgresql_isolate_db_instance_operation` ([#1820](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1820))
* **New Resource:** `tencentcloud_postgresql_rebalance_readonly_group_operation` ([#1820](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1820))
* **New Resource:** `tencentcloud_postgresql_renew_db_instance_operation` ([#1820](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1820))
* **New Resource:** `tencentcloud_postgresql_restart_db_instance_operation` ([#1820](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1820))
* **New Resource:** `tencentcloud_postgresql_security_group_config` ([#1806](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1806))
* **New Resource:** `tencentcloud_scf_function_event_invoke_config` ([#1829](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1829))
* **New Resource:** `tencentcloud_scf_function_version` ([#1829](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1829))
* **New Resource:** `tencentcloud_sqlserver_business_intelligence_file` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Resource:** `tencentcloud_sqlserver_business_intelligence_instance` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Resource:** `tencentcloud_sqlserver_complete_expansion` ([#1811](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1811))
* **New Resource:** `tencentcloud_sqlserver_config_database_cdc` ([#1811](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1811))
* **New Resource:** `tencentcloud_sqlserver_config_database_ct` ([#1811](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1811))
* **New Resource:** `tencentcloud_sqlserver_config_database_mdf` ([#1811](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1811))
* **New Resource:** `tencentcloud_sqlserver_config_instance_param` ([#1811](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1811))
* **New Resource:** `tencentcloud_sqlserver_config_instance_ro_group` ([#1811](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1811))
* **New Resource:** `tencentcloud_sqlserver_config_terminate_db_instance` ([#1816](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1816))
* **New Resource:** `tencentcloud_sqlserver_general_cloud_instance` ([#1816](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1816))
* **New Resource:** `tencentcloud_sqlserver_general_communication` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Resource:** `tencentcloud_sqlserver_incre_backup_migration` ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* **New Resource:** `tencentcloud_sqlserver_renew_db_instance` ([#1816](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1816))
* **New Resource:** `tencentcloud_sqlserver_renew_postpaid_db_instance` ([#1816](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1816))
* **New Resource:** `tencentcloud_sqlserver_renew_postpaid_db_instance` ([#1816](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1816))
* **New Resource:** `tencentcloud_sqlserver_restart_db_instance` ([#1816](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1816))
* **New Resource:** `tencentcloud_sqlserver_restore_instance` ([#1816](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1816))
* **New Resource:** `tencentcloud_sqlserver_rollback_instance` ([#1816](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1816))
* **New Resource:** `tencentcloud_tsf_application_file_config` ([#1804](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1804))
* **New Resource:** `tencentcloud_tsf_enable_unit_rule` ([#1804](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1804))
* **New Resource:** `tencentcloud_vpc_flow_log_config` ([#1805](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1805))
* **New Resource:** `tencentcloud_vpc_net_detect` ([#1805](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1805))
* **New Resource:** `tencentcloud_vpc_snapshot_policy` ([#1799](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1799))
* **New Resource:** `tencentcloud_vpc_snapshot_policy_attachment` ([#1799](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1799))
* **New Resource:** `tencentcloud_vpc_snapshot_policy_config` ([#1801](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1801))
* **New Resource:** `tencentcloud_vpc_traffic_package` ([#1792](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1792))

ENHANCEMENTS:

* resource/tencentcloud_cdn_domain: fix compression problem ([#1814](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1814))
* resource/tencentcloud_ckafka_instance: support set `max_message_byte` ([#1819](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1819))
* resource/tencentcloud_clb_listener: `health_check_http_code` support multi code ([#1824](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1824))
* resource/tencentcloud_cls_config: support import function ([#1815](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1815))
* resource/tencentcloud_cls_config_attachment: support import function ([#1815](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1815))
* resource/tencentcloud_cls_config_extra: support import function ([#1815](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1815))
* resource/tencentcloud_kubernetes_cluster: fix tke doc typo ([#1788](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1788))
* resource/tencentcloud_sqlserver_full_backup_migration: update code ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* resource/tencentcloud_sqlserver_general_clone: update code ([#1791](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1791))
* resource/tencentcloud_sqlserver_instance: Support wait_switch ([#1811](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1811))
* resource/tencentcloud_sqlserver_instance: support update instance 'vpc-id', 'subnet-id' parameters ([#1816](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1816))

## 1.81.1 (May 17, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_dbbrain_db_space_status` ([#1776](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1776))
* **New Data Source:** `tencentcloud_dbbrain_slow_logs` ([#1776](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1776))
* **New Data Source:** `tencentcloud_dbbrain_sql_templates` ([#1776](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1776))
* **New Data Source:** `tencentcloud_eip_address_quota` ([#1780](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1780))
* **New Data Source:** `tencentcloud_eip_network_account_type` ([#1780](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1780))
* **New Data Source:** `tencentcloud_nat_dc_route` ([#1775](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1775))
* **New Data Source:** `tencentcloud_redis_instance_node_info` ([#1784](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1784))
* **New Data Source:** `tencentcloud_tse_nacos_server_interfaces` ([#1774](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1774))
* **New Data Source:** `tencentcloud_vpc_bandwidth_package_bill_usage` ([#1780](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1780))
* **New Data Source:** `tencentcloud_vpc_bandwidth_package_quota` ([#1780](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1780))
* **New Resource:** `tencentcloud_eip_normal_address_return` ([#1780](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1780))
* **New Resource:** `tencentcloud_eip_public_address_adjust` ([#1780](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1780))
* **New Resource:** `tencentcloud_lighthouse_disk` ([#1782](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1782))
* **New Resource:** `tencentcloud_lighthouse_key_pair_attachment` ([#1782](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1782))
* **New Resource:** `tencentcloud_lighthouse_reboot_instance` ([#1782](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1782))
* **New Resource:** `tencentcloud_lighthouse_start_instance` ([#1782](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1782))
* **New Resource:** `tencentcloud_lighthouse_stop_instance` ([#1782](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1782))
* **New Resource:** `tencentcloud_nat_refresh_nat_dc_route` ([#1775](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1775))
* **New Resource:** `tencentcloud_redis_backup_download_restriction` ([#1680](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1680))
* **New Resource:** `tencentcloud_redis_clear_instance_operation` ([#1680](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1680))
* **New Resource:** `tencentcloud_redis_renew_instance_operation` ([#1680](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1680))
* **New Resource:** `tencentcloud_redis_startup_instance_operation` ([#1680](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1680))
* **New Resource:** `tencentcloud_redis_switch_master` ([#1785](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1785))
* **New Resource:** `tencentcloud_redis_upgrade_cache_version_operation` ([#1680](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1680))
* **New Resource:** `tencentcloud_redis_upgrade_multi_zone_operation` ([#1680](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1680))
* **New Resource:** `tencentcloud_redis_upgrade_proxy_version_operation` ([#1680](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1680))

ENHANCEMENTS:

* resource/tencentcloud_gaap_proxy: support param network_type ([#1783](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1783))
* resource/tencentcloud_kubernetes_addon_attachment: The instance supports instance ip, vpc_id, subnet_id, port modification ([#1784](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1784))
* resource/tencentcloud_kubernetes_addon_attachment: adjust the example for installing tcr addon ([#1770](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1770))
* resource/tencentcloud_kubernetes_cluster: adjust notice for the cluster usage. ([#1786](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1786))
* resource/tencentcloud_kubernetes_cluster_endpoint: add parameter extensive_parameters ([#1762](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1762))
* resource/tencentcloud_lighthouse_instance: support update `bundle_id`, `blueprint_id`, `renew_flag`, `password`, support params `permit_default_key_pair_login`, `is_update_bundle_id_auto_voucher` ([#1782](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1782))
* resource/tencentcloud_mysql_param_template: support set `param_list` ([#1777](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1777))
* resource/tencentcloud_postgresql_instance: support to update minor kernel versions by field `db_kernel_version` ([#1712](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1712))
* resource/tencentcloud_redis_connection_config: Deprecate the attribute bandwidth, add the attributes total_bandwidth, add_bandwidth, min_add_bandwidth, max_add_bandwidth ([#1784](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1784))
* resource/tencentcloud_redis_instance: Support `typeid` immediate switching ([#1680](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1680))

## 1.81.0 (May 12, 2023)

Upgrade terraform plugin sdk from v1 to v2 ([#1714](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1714))

FEATURES:

* **New Data Source:** `tencentcloud_api_gateway_api_apps` ([#1726](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1726))
* **New Data Source:** `tencentcloud_api_gateway_api_docs` ([#1726](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1726))
* **New Data Source:** `tencentcloud_ccn_cross_border_region_bandwidth_limits` ([#1703](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1703))
* **New Data Source:** `tencentcloud_dbbrain_health_scores` ([#1755](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1755))
* **New Data Source:** `tencentcloud_kubernetes_cluster_authentication_options` ([#1767](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1767))
* **New Data Source:** `tencentcloud_lighthouse_bundle` ([#1704](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1704))
* **New Data Source:** `tencentcloud_lighthouse_disk_config` ([#1759](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1759))
* **New Data Source:** `tencentcloud_lighthouse_instance_blueprint` ([#1759](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1759))
* **New Data Source:** `tencentcloud_lighthouse_instance_disk_num` ([#1759](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1759))
* **New Data Source:** `tencentcloud_lighthouse_instance_traffic_package` ([#1759](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1759))
* **New Data Source:** `tencentcloud_lighthouse_instance_vnc_url` ([#1759](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1759))
* **New Data Source:** `tencentcloud_lighthouse_region` ([#1759](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1759))
* **New Data Source:** `tencentcloud_lighthouse_reset_instance_blueprint` ([#1759](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1759))
* **New Data Source:** `tencentcloud_lighthouse_scene` ([#1759](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1759))
* **New Data Source:** `tencentcloud_lighthouse_zone` ([#1759](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1759))
* **New Data Source:** `tencentcloud_mariadb_database_objects` ([#1690](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1690))
* **New Data Source:** `tencentcloud_mariadb_database_table` ([#1690](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1690))
* **New Data Source:** `tencentcloud_mariadb_databases` ([#1690](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1690))
* **New Data Source:** `tencentcloud_mongodb_instance_backups` ([#1754](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1754))
* **New Data Source:** `tencentcloud_mongodb_instance_connections` ([#1754](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1754))
* **New Data Source:** `tencentcloud_mongodb_instance_current_op` ([#1754](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1754))
* **New Data Source:** `tencentcloud_mongodb_instance_params` ([#1754](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1754))
* **New Data Source:** `tencentcloud_mongodb_instance_slow_log` ([#1754](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1754))
* **New Data Source:** `tencentcloud_sqlserver_backup_commands` ([#1765](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1765))
* **New Data Source:** `tencentcloud_tcr_image_manifests` ([#1760](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1760))
* **New Data Source:** `tencentcloud_tcr_images` ([#1707](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1707))
* **New Data Source:** `tencentcloud_tcr_tag_retention_execution_tasks` ([#1760](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1760))
* **New Data Source:** `tencentcloud_tse_access_address` ([#1736](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1736))
* **New Data Source:** `tencentcloud_tse_nacos_replicas` ([#1736](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1736))
* **New Data Source:** `tencentcloud_tse_zookeeper_replicas` ([#1736](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1736))
* **New Data Source:** `tencentcloud_tse_zookeeper_server_interfaces` ([#1736](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1736))
* **New Data Source:** `tencentcloud_tsf_api_detail` ([#1761](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1761))
* **New Data Source:** `tencentcloud_tsf_api_group` ([#1761](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1761))
* **New Data Source:** `tencentcloud_tsf_application` ([#1694](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1694))
* **New Data Source:** `tencentcloud_tsf_application_attribute` ([#1761](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1761))
* **New Data Source:** `tencentcloud_tsf_application_config` ([#1694](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1694))
* **New Data Source:** `tencentcloud_tsf_application_file_config` ([#1694](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1694))
* **New Data Source:** `tencentcloud_tsf_application_public_config` ([#1694](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1694))
* **New Data Source:** `tencentcloud_tsf_business_log_configs` ([#1761](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1761))
* **New Data Source:** `tencentcloud_tsf_cluster` ([#1694](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1694))
* **New Data Source:** `tencentcloud_tsf_config_summary` ([#1764](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1764))
* **New Data Source:** `tencentcloud_tsf_delivery_config_by_group_id` ([#1764](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1764))
* **New Data Source:** `tencentcloud_tsf_delivery_configs` ([#1764](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1764))
* **New Data Source:** `tencentcloud_tsf_microservice` ([#1694](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1694))
* **New Data Source:** `tencentcloud_tsf_microservice_api_version` ([#1761](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1761))
* **New Data Source:** `tencentcloud_tsf_public_config_summary` ([#1764](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1764))
* **New Data Source:** `tencentcloud_tsf_unit_rules` ([#1694](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1694))
* **New Resource:** `tencentcloud_api_gateway_api_app` ([#1726](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1726))
* **New Resource:** `tencentcloud_api_gateway_api_doc` ([#1726](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1726))
* **New Resource:** `tencentcloud_dbbrain_tdsql_audit_log` ([#1755](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1755))
* **New Resource:** `tencentcloud_eip_address_transform` ([#1769](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1769))
* **New Resource:** `tencentcloud_mariadb_encrypt_attributes_operation` ([#1690](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1690))
* **New Resource:** `tencentcloud_redis_replica_readonly` ([#1664](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1664))
* **New Resource:** `tencentcloud_scf_function_alias` ([#1706](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1706))
* **New Resource:** `tencentcloud_sqlserver_full_backup_migration` ([#1768](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1768))
* **New Resource:** `tencentcloud_sqlserver_general_backup` ([#1765](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1765))
* **New Resource:** `tencentcloud_sqlserver_general_clone` ([#1765](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1765))
* **New Resource:** `tencentcloud_tcr_create_image_signature_operation` ([#1707](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1707))
* **New Resource:** `tencentcloud_tcr_customized_domain` ([#1707](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1707))
* **New Resource:** `tencentcloud_tcr_delete_image_operation` ([#1707](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1707))
* **New Resource:** `tencentcloud_tcr_immutable_tag_rule` ([#1707](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1707))
* **New Resource:** `tencentcloud_tse_instance` ([#1736](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1736))
* **New Resource:** `tencentcloud_tsf_application` ([#1698](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1698))
* **New Resource:** `tencentcloud_tsf_application_file_config_release` ([#1698](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1698))
* **New Resource:** `tencentcloud_tsf_application_public_config` ([#1698](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1698))
* **New Resource:** `tencentcloud_tsf_application_public_config_release` ([#1698](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1698))
* **New Resource:** `tencentcloud_tsf_bind_api_group` ([#1761](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1761))
* **New Resource:** `tencentcloud_tsf_instances_attachment` ([#1698](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1698))

ENHANCEMENTS:

* resource/tencentcloud_api_gateway_custom_domain: support add `is_forced_https` params ([#1726](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1726))
* resource/tencentcloud_key_pair: support created by `CreateKeyPair` ([#1742](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1742))

BUG FIXES:

* tencentcloud_clb_instance: fix the nil exception caused by the security_groups ([#1763](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1763))
* tencentcloud_security_group_rule: fix the nil exception caused by the policy index changing ([#1705](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1705))

## 1.80.6 (April 27 2023)

FEATURES:

* **New Data Source:** `tencentcloud_tat_invocation_task` ([#1692](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1692))
* **New Data Source:** `tencentcloud_tcr_describe_webhook_trigger_logs` ([#1695](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1695))
* **New Resource:** `tencentcloud_ipv6_address_bandwidth` ([#1702](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1702))
* **New Resource:** `tencentcloud_kubernetes_backup_storage_location` ([#1691](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1691))
* **New Resource:** `tencentcloud_lighthouse_apply_instance_snapshot` ([#1693](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1693))
* **New Resource:** `tencentcloud_lighthouse_snapshot` ([#1693](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1693))
* **New Resource:** `tencentcloud_monitor_tmp_grafana_config` ([#1699](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1699))
* **New Resource:** `tencentcloud_tcr_manage_replication_operation` ([#1700](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1700))
* **New Resource:** `tencentcloud_tcr_webhook_trigger` ([#1695](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1695))

ENHANCEMENTS:

* resource/tencentcloud_eip: support update `internet_charge_type` ([#1702](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1702))
* resource/tencentcloud_elasticsearch_instance: affect `es_acl` when first apply ([#1697](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1697))
* resource/tencentcloud_monitor_grafana_instance: Query results support `internet_url`, `internal_url` ([#1699](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1699))

## 1.80.5 (April 21 2023)

FEATURES:

* **New Resource:** `tencentcloud_lighthouse_apply_disk_backup` ([#1687](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1687))
* **New Resource:** `tencentcloud_lighthouse_disk_attachment` ([#1687](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1687))
* **New Resource:** `tencentcloud_lighthouse_disk_backup` ([#1687](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1687))
* **New Resource:** `tencentcloud_lighthouse_key_pair` ([#1687](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1687))

ENHANCEMENTS:

* resource/tencentcloud_ckafka_topic: support update `replica_num` ([#1686](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1686))
* resource/tencentcloud_sqlserver_basic_instance: support set `CLOUD_BSSD` and `CLOUD_HSSD` ([#1685](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1685))

BUG FIXES:

* resource/tencentcloud_instance: fix key_ids update exception ([#1688](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1688))

## 1.80.4 (April 19 2023)

FEATURES:

* **New Data Source:** `tencentcloud_ccn_cross_border_compliance` ([#1675](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1675))
* **New Data Source:** `tencentcloud_ccn_cross_border_flow_monitor` ([#1675](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1675))
* **New Data Source:** `tencentcloud_ccn_tenant_instances` ([#1675](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1675))
* **New Data Source:** `tencentcloud_cvm_chc_denied_actions` ([#1669](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1669))
* **New Data Source:** `tencentcloud_cvm_chc_hosts` ([#1669](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1669))
* **New Data Source:** `tencentcloud_vpn_customer_gateway_vendors` ([#1675](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1675))
* **New Resource:** `tencentcloud_cvm_chc_config` ([#1669](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1669))
* **New Resource:** `tencentcloud_mongodb_instance_backup` ([#1679](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1679))
* **New Resource:** `tencentcloud_sqlserver_config_backup_strategy` ([#1678](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1678))

ENHANCEMENTS:

* resource/tencentcloud_ccn_attachment: support set `description` ([#1674](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1674))
* resource/tencentcloud_cfs_file_system: support create with tag ([#1681](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1681))
* resource/tencentcloud_cfs_snapshot: support create with tag ([#1681](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1681))
* resource/tencentcloud_instance: support set `disable_api_termination` ([#1682](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1682))

## 1.80.3 (April 14 2023)

FEATURES:

* **New Data Source:** `tencentcloud_ckafka_connect_resource` ([#1666](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1666))
* **New Data Source:** `tencentcloud_cvm_disaster_recover_group_quota` ([#1663](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1663))
* **New Data Source:** `tencentcloud_cvm_instance_vnc_url` ([#1663](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1663))
* **New Data Source:** `tencentcloud_mysql_backup_overview` ([#1670](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1670))
* **New Data Source:** `tencentcloud_mysql_backup_summaries` ([#1670](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1670))
* **New Data Source:** `tencentcloud_mysql_bin_log` ([#1670](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1670))
* **New Data Source:** `tencentcloud_mysql_binlog_backup_overview` ([#1670](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1670))
* **New Data Source:** `tencentcloud_mysql_clone_list` ([#1670](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1670))
* **New Data Source:** `tencentcloud_mysql_data_backup_overview` ([#1670](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1670))
* **New Data Source:** `tencentcloud_redis_instance_shards` ([#1660](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1660))
* **New Data Source:** `tencentcloud_redis_instance_task_list` ([#1660](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1660))
* **New Data Source:** `tencentcloud_redis_instance_zone_info` ([#1660](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1660))
* **New Data Source:** `tencentcloud_redis_param_records` ([#1660](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1660))
* **New Resource:** `tencentcloud_ckafka_datahub_topic` ([#1661](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1661))
* **New Resource:** `tencentcloud_cvm_reboot_instance` ([#1663](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1663))
* **New Resource:** `tencentcloud_dbbrain_db_diag_report_task` ([#1656](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1656))
* **New Resource:** `tencentcloud_dbbrain_modify_diag_db_instance_operation` ([#1662](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1662))
* **New Resource:** `tencentcloud_mongodb_instance_account` ([#1673](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1673))
* **New Resource:** `tencentcloud_tcr_tag_retention_rule` ([#1668](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1668))

ENHANCEMENTS:

* resource/tencentcloud_kubernetes_auth_attachment: change para issuer in tke auth attachment from required to optional ([#1667](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1667))
* resource/tencentcloud_kubernetes_cluster: update document ([#1667](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1667))
* resource/tencentcloud_tem_gateway: Add binding certificate instructions to deal with binding certificate issues ([#1672](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1672))

## 1.80.2 (April 10 2023)

FEATURES:

* **New Data Source:** `tencentcloud_redis_backup` ([#1652](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1652))
* **New Data Source:** `tencentcloud_redis_backup_download_info` ([#1652](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1652))
* **New Resource:** `tencentcloud_ccn_instances_accept_attach` ([#1653](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1653))
* **New Resource:** `tencentcloud_ccn_instances_reset_attach` ([#1653](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1653))
* **New Resource:** `tencentcloud_ccn_routes` ([#1653](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1653))
* **New Resource:** `tencentcloud_dts_compare_task_stop_operation` ([#1638](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1638))
* **New Resource:** `tencentcloud_dts_migrate_job_config` ([#1638](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1638))
* **New Resource:** `tencentcloud_dts_migrate_job_resume_operation` ([#1638](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1638))
* **New Resource:** `tencentcloud_dts_sync_check_job_operation` ([#1638](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1638))
* **New Resource:** `tencentcloud_dts_sync_job_resume_operation` ([#1638](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1638))
* **New Resource:** `tencentcloud_redis_account` ([#1652](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1652))
* **New Resource:** `tencentcloud_redis_maintenance_window` ([#1652](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1652))
* **New Resource:** `tencentcloud_redis_read_only` ([#1652](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1652))
* **New Resource:** `tencentcloud_redis_ssl` ([#1652](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1652))
* **New Resource:** `tencentcloud_vpn_gateway_ssl_client_cert` ([#1650](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1650))

ENHANCEMENTS:

* resource/tencentcloud_dts_migrate_job: adjust to avoid the unexpected diff when the `password` field is not modified ([#1638](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1638))
* resource/tencentcloud_elasticsearch_instance: set `basic_security_type` default 1 when import instance ([#1655](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1655))
* resource/tencentcloud_kubernetes_auth_attachment: add oidc parameters ([#1651](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1651))
* resource/tencentcloud_kubernetes_cluster: add oidc parameters ([#1651](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1651))

## 1.80.1 (April 6 2023)

FEATURES:

* **New Resource:** `tencentcloud_tsf_api_rate_limit_rule` ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))
* **New Resource:** `tencentcloud_tsf_application_release_config` ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))
* **New Resource:** `tencentcloud_tsf_cluster` ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))
* **New Resource:** `tencentcloud_tsf_config_template` ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))
* **New Resource:** `tencentcloud_tsf_group` ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))
* **New Resource:** `tencentcloud_tsf_lane` ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))
* **New Resource:** `tencentcloud_tsf_lane_rule` ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))
* **New Resource:** `tencentcloud_tsf_path_rewrite` ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))
* **New Resource:** `tencentcloud_tsf_task` ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))
* **New Resource:** `tencentcloud_tsf_unit_rule` ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))

ENHANCEMENTS:

* resource/tencentcloud_clb_attachment: remove update targets limit ([#1644](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1644))
* resource/tencentcloud_redis_instance:  `no_auth` support modify. ([#1641](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1641))
* resource/tencentcloud_tsf_namespace: Cancel support for import ([#1632](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1632))
* resource/tencentcloud_vpn_gateway: support update `prepaid_period` ([#1649](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1649))## 1.80.0 (March 31 2023)

## 1.80.0 (March 31 2023)

FEATURES:

* **New Resource:** `tencentcloud_cvm_security_group_attachment` ([#1633](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1633))
* **New Resource:** `tencentcloud_monitor_tmp_tke_basic_config` ([#1635](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1635))
* **New Resource:** `tencentcloud_vpn_connection_reset` ([#1636](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1636))
* **New Resource:** `tencentcloud_vpn_customer_gateway_configuration_download` ([#1636](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1636))

ENHANCEMENTS:

* resource/tencentcloud_cdn_domain:  `service_type` support set `hybrid`, `dynamic` ([#1637](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1637))

## 1.79.19 (March 29 2023)

FEATURES:

* **New Data Source:** `tencentcloud_lighthouse_firewall_rules_template` ([#1624](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1624))
* **New Data Source:** `tencentcloud_postgresql_parameter_templates` ([#1625](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1625))
* **New Resource:** `tencentcloud_cvm_launch_template_default_version` ([#1626](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1626))
* **New Resource:** `tencentcloud_lighthouse_firewall_rule` ([#1624](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1624))
* **New Resource:** `tencentcloud_mysql_audit_log_file` ([#1629](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1629))
* **New Resource:** `tencentcloud_postgresql_parameter_template` ([#1625](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1625))

ENHANCEMENTS:

* resource/tencentcloud_cfs_auto_snapshot_policy:  support set day_of_month and interval_days ([#1631](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1631))
* data-source/tencentcloud_ssm_secret_versions: return null when a resource is not found ([#1627](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1627))
* resource/tencentcloud_monitor_tmp_exporter_integration: update doc to introduce how to upgrade exporter version ([#1630](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1630))

## 1.79.18 (March 24 2023)

FEATURES:

* **New Data Source:** `tencentcloud_dbbrain_diag_event` ([#1610](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1610))
* **New Data Source:** `tencentcloud_dbbrain_diag_events` ([#1610](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1610))
* **New Data Source:** `tencentcloud_dbbrain_diag_history` ([#1610](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1610))
* **New Data Source:** `tencentcloud_dbbrain_security_audit_log_download_urls` ([#1610](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1610))
* **New Data Source:** `tencentcloud_dbbrain_slow_log_time_series_stats` ([#1610](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1610))
* **New Data Source:** `tencentcloud_dbbrain_slow_log_top_sqls` ([#1610](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1610))
* **New Data Source:** `tencentcloud_dbbrain_slow_log_user_host_stats` ([#1610](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1610))
* **New Data Source:** `tencentcloud_dbbrain_slow_log_user_sql_advice` ([#1610](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1610))
* **New Resource:** `tencentcloud_apm_instance` ([#1619](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1619))
* **New Resource:** `tencentcloud_cvm_launch_template_version` ([#1617](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1617))
* **New Resource:** `tencentcloud_lighthouse_blueprint` ([#1613](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1613))
* **New Resource:** `tencentcloud_monitor_tmp_manage_grafana_attachment` ([#1611](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1611))

BUG FIXES:

* resource/tencentcloud_instance: fix npe ([#1618](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1618))

## 1.79.17 (March 17 2023)

FEATURES:

* **New Data Source:** `tencentcloud_kubernetes_available_cluster_versions` ([#1608](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1608))

ENHANCEMENTS:

* resource/tencentcloud_kubernetes_cluster_attachment: support setting tke gpu args ([#1593](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1593))
* resource/tencentcloud_kubernetes_node_pool: support setting tke gpu args ([#1593](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1593))
* resource/tencentcloud_kubernetes_scale_worker: support setting tke gpu args ([#1593](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1593))

## 1.79.16 (March 15, 2023)

ENHANCEMENTS:

* resource/tencentcloud_mysql_instance: import support vpc_id and subnet_id ([#1605](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1605))

## 1.79.15 (March 14, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_tcm_mesh` ([#1600](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1600))
* **New Resource:** `tencentcloud_mariadb_instance` ([#1525](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1525))
* **New Resource:** `tencentcloud_mps_person_sample` ([#1601](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1601))

ENHANCEMENTS:

* resource/tencentcloud_mysql_account: support import ([#1598](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1598))

BUG FIXES:

* resource/tencentcloud_vpn_connection: fix dpd_timeout read error ([#1597](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1597))

## 1.79.14 (March 08, 2023)

FEATURES:

* **New Resource:** `tencentcloud_css_play_auth_key_config` ([#1587](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1587))
* **New Resource:** `tencentcloud_css_play_domain_cert_attachment` ([#1587](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1587))
* **New Resource:** `tencentcloud_css_push_auth_key_config` ([#1587](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1587))
* **New Resource:** `tencentcloud_mdl_stream_live_input` ([#1588](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1588))

ENHANCEMENTS:

* resource/tencentcloud_ckafka_instance: support set `instance_type` ([#1589](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1589))
* resource/tencentcloud_clb_listener_rule: support  `health_check_type` and `health_check_time_out` ([#1590](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1590))

## 1.79.13 (March 03, 2023)

FEATURES:

* **New Resource:** `tencentcloud_mps_adaptive_dynamic_streaming_template` ([#1585](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1585))
* **New Resource:** `tencentcloud_mps_ai_analysis_template` ([#1585](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1585))
* **New Resource:** `tencentcloud_mps_ai_recognition_template` ([#1585](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1585))
* **New Resource:** `tencentcloud_mps_animated_graphics_template` ([#1585](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1585))
* **New Resource:** `tencentcloud_mps_image_sprite_template` ([#1585](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1585))
* **New Resource:** `tencentcloud_mps_sample_snapshot_template` ([#1585](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1585))
* **New Resource:** `tencentcloud_mps_snapshot_by_timeoffset_template` ([#1585](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1585))

ENHANCEMENTS:

* resource/tencentcloud_elasticsearch_instance: adjust the scope of `license_type` field. ([#1574](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1574))
* resource/tencentcloud_mongodb_sharding_instance: support `availability_zone_list` and `hidden_zone`. ([#1582](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1582))

## 1.79.12 (February 27, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_chdfs_access_groups` ([#1572](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1572))
* **New Data Source:** `tencentcloud_chdfs_mount_points` ([#1572](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1572))
* **New Resource:** `tencentcloud_chdfs_life_cycle_rule` ([#1572](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1572))
* **New Resource:** `tencentcloud_chdfs_mount_point` ([#1572](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1572))
* **New Resource:** `tencentcloud_chdfs_mount_point_attachment` ([#1572](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1572))

ENHANCEMENTS:

* resource/tencentcloud_kubernetes_cluster: support setting tke cluster internet/intranet domain ([#1564](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1564))
* resource/tencentcloud_kubernetes_cluster: update tke cluster resource doc detail ([#1560](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1560))
* resource/tencentcloud_kubernetes_cluster_endpoint: support setting tke cluster internet/intranet domain ([#1564](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1564))

BUG FIXES:

* resource/tencentcloud_mongodb_instance: clean mongos params ([#1573](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1573))
* resource/tencentcloud_mongodb_sharding_instance: fix mongos_memory unit ([#1573](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1573))
* resource/tencentcloud_mongodb_standby_instance: clean mongos params ([#1573](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1573))

## 1.79.11 (February 23, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_css_domains` ([#1568](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1568))
* **New Resource:** `tencentcloud_chdfs_access_group` ([#1567](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1567))
* **New Resource:** `tencentcloud_chdfs_access_rule` ([#1570](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1570))
* **New Resource:** `tencentcloud_chdfs_file_system` ([#1570](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1570))
* **New Resource:** `tencentcloud_css_authenticate_domain_owner_operation` ([#1568](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1568))
* **New Resource:** `tencentcloud_css_domain` ([#1568](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1568))
* **New Resource:** `tencentcloud_css_domain_config` ([#1568](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1568))
* **New Resource:** `tencentcloud_mps_watermark_template` ([#1559](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1559))

ENHANCEMENTS:

* resource/tencentcloud_clb_function_targets_attachment: support update function targets. ([#1561](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1561))
* resource/tencentcloud_eip: support create prepaid eip. ([#1563](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1563))
* resource/tencentcloud_mongodb_sharding_instance: support import. ([#1566](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1566))

## 1.79.10 (February 15, 2023)

ENHANCEMENTS:

* resource/tencentcloud_clb_function_targets_attachment: set weight default `10`. ([#1554](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1554))
* resource/tencentcloud_tcr_namespace: Support is_auto_scan, is_prevent_vul, severity, cve_whitelist_items field ([#1552](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1552))

BUG FIXES:

* resource/tencentcloud_dayu_ddos_ip_attachment_v2: fix delete failed ([#1558](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1558))

## 1.79.9 (February 10, 2023)

FEATURES:

* **New Resource:** `tencentcloud_cbs_disk_backup` ([#1548](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1548))
* **New Resource:** `tencentcloud_cbs_disk_backup_rollback_operation` ([#1548](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1548))
* **New Resource:** `tencentcloud_cbs_snapshot_share_permission` ([#1548](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1548))
* **New Resource:** `tencentcloud_clb_function_targets_attachment` ([#1549](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1549))
* **New Resource:** `tencentcloud_mps_transcode_template` ([#1550](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1550))

ENHANCEMENTS:

* resource/tencentcloud_cbs_storage: support disk_backup_quota param ([#1548](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1548))
* resource/tencentcloud_kubernetes_cluster: support creating tke cluster endpoint even if cluster only have serverless node ([#1546](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1546))
* resource/tencentcloud_kubernetes_cluster_endpoint: support creating tke cluster endpoint even if cluster only have serverless node ([#1546](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1546))
* resource/tencentcloud_kubernetes_node_pool: adjust enhanced_security_service to not forceNew ([#1545](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1545))
* resource/tencentcloud_vpc_bandwidth_package: Support internet_max_bandwidth field ([#1551](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1551))

## 1.79.8 (February 7, 2023)

ENHANCEMENTS:

* resource/tencentcloud_mongodb_sharding_instance: support mongos params ([#1543](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1543))

BUG FIXES:

* resource/tencentcloud_tem_application: Make the `description` parameter of the TEM Application required ([#1527](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1527))

## 1.79.7 (February 3, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_cfs_available_zone` ([#1522](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1522))
* **New Data Source:** `tencentcloud_cfs_file_system_clients` ([#1522](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1522))
* **New Data Source:** `tencentcloud_cfs_mount_targets` ([#1522](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1522))
* **New Data Source:** `tencentcloud_cvm_instances_modification` ([#1521](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1521))
* **New Resource:** `tencentcloud_cfs_sign_up_cfs_service` ([#1522](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1522))
* **New Resource:** `tencentcloud_cfs_snapshot` ([#1522](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1522))
* **New Resource:** `tencentcloud_cvm_launch_template` ([#1521](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1521))
* **New Resource:** `tencentcloud_kubernetes_serverless_node_pool` ([#1519](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1519))
* **New Resource:** `tencentcloud_mps_workflow` ([#1541](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1541))
* **New Resource:** `tencentcloud_sqlserver_migration` ([#1523](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1523))
* **New Resource:** `tencentcloud_tsf_api_group` ([#1524](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1524))
* **New Resource:** `tencentcloud_tsf_namespace` ([#1524](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1524))

ENHANCEMENTS:

* resource/tencentcloud_ckafka_instance: support IP query when modifying the type ([#1520](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1520))
* resource/tencentcloud_cos_bucket: add `acceleration_enable` field ([#1537](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1537))
* resource/tencentcloud_cos_bucket: add `endpoint` field for static website ([#1535](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1535))
* resource/tencentcloud_monitor_tmp_instance: add computed field: `ipv4_address`, `proxy_address`, `remote_write`, `api_root_path` ([#1529](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1529))
* resource/tencentcloud_scf_function: keep state consistent when vpc_id is not set ([#1531](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1531))
* resource/tencentcloud_teo_rule_engine: Support `sub_rules` field, support multi-layer `if` ([#1538](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1538))

## 1.79.6 (January 13, 2023)

FEATURES:

* **New Resource:** `tencentcloud_tsf_application_config` ([#1513](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1513))
* **New Resource:** `tencentcloud_tsf_microservice` ([#1513](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1513))

ENHANCEMENTS:

* datasource/tencentcloud_cfs_file_systems: add computed fs_id ([#1517](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1517))
* resource/tencentcloud_cfs_file_system: add computed fs_id ([#1517](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1517))
* resource/tencentcloud_ckafka_instance: support create standard instance ([#1518](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1518))
* resource/tencentcloud_tcr_instance: support prepaid arguments ([#1515](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1515))
* resource/tencentcloud_tem_application_service: Support for inputting vpc_id and subnet_id, and adding IP for output ([#1514](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1514))

## 1.79.5 (January 11, 2023)

FEATURES:

* **New Data Source:** `tencentcloud_cynosdb_accounts` ([#1508](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1508))
* **New Data Source:** `tencentcloud_cynosdb_cluster_instance_groups` ([#1508](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1508))
* **New Data Source:** `tencentcloud_cynosdb_cluster_params` ([#1508](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1508))
* **New Data Source:** `tencentcloud_cynosdb_param_templates` ([#1508](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1508))
* **New Data Source:** `tencentcloud_dcdb_database_objects` ([#1500](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1500))
* **New Data Source:** `tencentcloud_dcdb_database_tables` ([#1500](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1500))
* **New Resource:** `tencentcloud_cynosdb_audit_log_file` ([#1508](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1508))
* **New Resource:** `tencentcloud_cynosdb_security_group` ([#1508](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1508))
* **New Resource:** `tencentcloud_dayu_ddos_ip_attachment_v2` ([#1511](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1511))
* **New Resource:** `tencentcloud_dcdb_db_parameters` ([#1500](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1500))
* **New Resource:** `tencentcloud_dts_migrate_job` ([#1490](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1490))
* **New Resource:** `tencentcloud_dts_migrate_job_start_operation` ([#1490](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1490))
* **New Resource:** `tencentcloud_dts_migrate_service` ([#1490](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1490))
* **New Resource:** `tencentcloud_mysql_deploy_group` ([#1512](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1512))
* **New Resource:** `tencentcloud_mysql_local_binlog_config` ([#1512](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1512))
* **New Resource:** `tencentcloud_mysql_param_template` ([#1512](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1512))
* **New Resource:** `tencentcloud_mysql_security_groups_attachment` ([#1512](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1512))
* **New Resource:** `tencentcloud_mysql_time_window` ([#1509](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1509))

ENHANCEMENTS:

* datasource/tencentcloud_mysql_zone_config: update the called api ([#1512](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1512))
* resource/tencentcloud_dcdb_account_privileges: optimize the logic after modifing the AccountPrivileges ([#1500](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1500))
* resource/tencentcloud_postgresql_instance: feat: pg - support prepaid arguments ([#1502](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1502))
* resource/tencentcloud_vpc_bandwidth_package_attachment: catch the error when bgp attach is nil ([#1504](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1504))

## 1.79.4 (January 5, 2023)

FEATURES:

* **New Resource:** `tencentcloud_api_gateway_plugin` ([#1496](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1496))
* **New Resource:** `tencentcloud_api_gateway_plugin_attachment` ([#1498](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1498))
* **New Resource:** `tencentcloud_ci_guetzli` ([#1489](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1489))
* **New Resource:** `tencentcloud_ci_original_image_protection` ([#1489](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1489))
* **New Resource:** `tencentcloud_dcdb_account_privileges` ([#1493](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1493))

ENHANCEMENTS:

* resource/tencentcloud_cdn_domain: support cache_key ([#1475](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1475))
* resource/tencentcloud_eip: support set `BandwidthPackageId` ([#1499](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1499))
* resource/tencentcloud_kubernetes_scale_worker: support adding hpc cluster id in worker config para ([#1477](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1477))
* resource/tencentcloud_tem_application: Support tags field ([#1494](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1494))
* resource/tencentcloud_tem_application_service: Support service modification ([#1494](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1494))
* resource/tencentcloud_tem_environment: Support service modification and support tags field ([#1494](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1494))

## 1.79.3 (December 29, 2022)

FEATURES:

* **New Data Source:** `tencentcloud_as_instances` ([#1482](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1482))
* **New Data Source:** `tencentcloud_tcmq_queue` ([#1480](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1480))
* **New Data Source:** `tencentcloud_tcmq_subscribe` ([#1480](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1480))
* **New Data Source:** `tencentcloud_tcmq_topic` ([#1480](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1480))
* **New Resource:** `tencentcloud_cfs_auto_snapshot_policy` ([#1481](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1481))
* **New Resource:** `tencentcloud_cfs_auto_snapshot_policy_attachment` ([#1481](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1481))
* **New Resource:** `tencentcloud_ci_bucket_attachment` ([#1471](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1471))
* **New Resource:** `tencentcloud_ci_bucket_attachment` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_bucket_pic_style` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_hot_link` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_animation_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_concat_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_pic_process_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_smart_cover_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_snapshot_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_speech_recognition_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_super_resolution_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_transcode_pro_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_transcode_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_tts_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_video_montage_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_video_process_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_voice_separate_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_ci_media_watermark_template` ([#1483](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1483))
* **New Resource:** `tencentcloud_tcmq_queue` ([#1480](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1480))
* **New Resource:** `tencentcloud_tcmq_subscribe` ([#1480](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1480))
* **New Resource:** `tencentcloud_tcmq_topic` ([#1480](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1480))
* **New Resource:** `tencentcloud_vpc_end_point_service_white_list` ([#1485](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1485))

ENHANCEMENTS:

* resource/tencentcloud_ckafka_instance: change `vpc_id` and `subnet_id` to optional, and optimize doc ([#1486](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1486))
* resource/tencentcloud_scf_function: optimize the retry logic when deleting the scf function ([#1478](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1478))
* resource/tencentcloud_ssl_pay_certificate: update the `product_id` value, support `confirm_letter` and `dv_auths` ([#1472](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1472))
* resource/tencentcloud_vpc_end_point_service: support service type. ([#1485](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1485))

## 1.79.2 (December 13, 2022)

FEATURES:

* **New Resource:** `tencentcloud_cvm_hpc_cluster` ([#1462](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1462))
* **New Resource:** `tencentcloud_tem_application_service` ([#1467](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1467))
* **New Resource:** `tencentcloud_vpc_end_point` ([#1466](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1466))
* **New Resource:** `tencentcloud_vpc_end_point_service` ([#1466](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1466))
* **New Resource:** `tencentcloud_vpc_flow_log` ([#1463](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1463))

ENHANCEMENTS:

* resource/tencentcloud_cam_service_linked_role: support for multiple qcs_service_name ([#1452](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1452))
* resource/tencentcloud_kubernetes_cluster: support tke cluster addon modify ([#1363](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1363))
* resource/tencentcloud_route_table_entry: update the `next_type` value ([#1461](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1461))

BUG FIXES:

* resource/tencentcloud_mysql_instance: fix: mysql - retry npe err ([#1460](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1460))

## 1.79.1 (December 7, 2022)

FEATURES:

* **New Data Source:** `tencentcloud_dts_compare_tasks` ([#1440](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1440))
* **New Data Source:** `tencentcloud_dts_migrate_jobs` ([#1440](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1440))
* **New Resource:** `tencentcloud_as_protect_instances` ([#1451](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1451))
* **New Resource:** `tencentcloud_as_remove_instances` ([#1451](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1451))
* **New Resource:** `tencentcloud_as_start_instances` ([#1451](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1451))
* **New Resource:** `tencentcloud_as_stop_instances` ([#1451](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1451))
* **New Resource:** `tencentcloud_dts_compare_task` ([#1440](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1440))
* **New Resource:** `tencentcloud_dts_migrate_job` ([#1440](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1440))
* **New Resource:** `tencentcloud_security_group_rule_set` ([#1453](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1453))

ENHANCEMENTS:

* resource/tencentcloud_gaap_security_rule: keep consistency when cidr is 1.1.1.1/32 and attr support change ([#1448](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1448))

## 1.79.0 (December 1, 2022)

FEATURES:

* **New Data Source:** `tencentcloud_tdmq_rocketmq_cluster` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))
* **New Data Source:** `tencentcloud_tdmq_rocketmq_group` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))
* **New Data Source:** `tencentcloud_tdmq_rocketmq_namespace` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))
* **New Data Source:** `tencentcloud_tdmq_rocketmq_role` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))
* **New Data Source:** `tencentcloud_tdmq_rocketmq_topic` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))
* **New Resource:** `tencentcloud_cam_service_linked_role` ([#1436](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1436))
* **New Resource:** `tencentcloud_tcm_access_log_config` ([#1444](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1444))
* **New Resource:** `tencentcloud_tcm_prometheus_attachment` ([#1444](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1444))
* **New Resource:** `tencentcloud_tcm_tracing_config` ([#1444](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1444))
* **New Resource:** `tencentcloud_tdmq_rocketmq_cluster` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))
* **New Resource:** `tencentcloud_tdmq_rocketmq_environment_role` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))
* **New Resource:** `tencentcloud_tdmq_rocketmq_group` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))
* **New Resource:** `tencentcloud_tdmq_rocketmq_namespace` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))
* **New Resource:** `tencentcloud_tdmq_rocketmq_role` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))
* **New Resource:** `tencentcloud_tdmq_rocketmq_topic` ([#1445](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1445))

ENHANCEMENTS:

* resource/tencentcloud_cos_bucket: enhance the storage_class of lifecycle ([#1434](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1434))
* resource/tencentcloud_cos_bucket_object: enhance the storage_class of lifecycle ([#1434](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1434))
* resource/tencentcloud_kubernetes_node_pool: optimize scaling_group_project_id default value ([#1439](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1439))
* resource/tencentcloud_kubernetes_node_pool: support host_name_style field ([#1435](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1435))
* resource/tencentcloud_mysql_instance: Support retry creating with client token ([#1330](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1330))
* resource/tencentcloud_tcm_mesh: support tracing, prometheus ([#1444](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1444))

BUG FIXES:

* resource/tencentcloud_mysql_instance: fix: mysql - no longer init ([#1443](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1443))

## 1.78.16 (November 29, 2022)

FEATURES:

* **New Data Source:** `tencentcloud_dts_sync_jobs` ([#1433](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1433))
* **New Resource:** `tencentcloud_audit_track` ([#1431](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1431))
* **New Resource:** `tencentcloud_dts_sync_job` ([#1433](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1433))
* **New Resource:** `tencentcloud_redis_param_template` ([#1432](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1432))

ENHANCEMENTS:

* resource/tencentcloud_cynosdb_cluster: fix: cynos - support serverless cluster pause/resume ([#1429](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1429))
* resource/tencentcloud_instance: support orderly security group ([#1430](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1430))

## 1.78.15 (November 24, 2022)

FEATURES:

* **New Data Source:** `tencentcloud_rum_offline_log_config` ([#1422](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1422))
* **New Data Source:** `tencentcloud_rum_project` ([#1422](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1422))
* **New Data Source:** `tencentcloud_rum_taw_instance` ([#1422](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1422))
* **New Data Source:** `tencentcloud_rum_whitelist` ([#1422](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1422))
* **New Resource:** `tencentcloud_rum_offline_log_config_attachment` ([#1422](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1422))
* **New Resource:** `tencentcloud_rum_project` ([#1422](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1422))
* **New Resource:** `tencentcloud_rum_taw_instance` ([#1422](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1422))
* **New Resource:** `tencentcloud_rum_whitelist` ([#1422](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1422))

## 1.78.14 (November 23, 2022)

FEATURES:

* **New Resource:** `tencentcloud_cam_policy_by_name` ([#1415](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1415))
* **New Resource:** `tencentcloud_cam_role_by_name` ([#1415](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1415))
* **New Resource:** `tencentcloud_cam_role_policy_attachment_by_name` ([#1415](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1415))
* **New Resource:** `tencentcloud_dbbrain_security_audit_log_export_task` ([#1417](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1417))
* **New Resource:** `tencentcloud_dbbrain_sql_filter` ([#1417](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1417))

ENHANCEMENTS:

* resource/tencentcloud_cynosdb_cluster: feat: tdsql-c - support serverless creation ([#1406](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1406))
* resource/tencentcloud_instance: update the `key_ids` of cvm smoothly ([#1403](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1403))

## 1.78.13 (November 22, 2022)

FEATURES:

* **New Data Source:** `tencentcloud_tdcpg_clusters` ([#1382](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1382))
* **New Data Source:** `tencentcloud_tdcpg_instances` ([#1382](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1382))
* **New Resource:** `tencentcloud_organization_org_member` ([#1412](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1412))
* **New Resource:** `tencentcloud_organization_org_node` ([#1408](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1408))
* **New Resource:** `tencentcloud_organization_policy_sub_account_attachment` ([#1412](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1412))
* **New Resource:** `tencentcloud_tdcpg_cluster` ([#1382](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1382))
* **New Resource:** `tencentcloud_tdcpg_instance` ([#1382](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1382))

ENHANCEMENTS:

* Provider resource increases the full name of the product ([#1395](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1395))
* resource/tencentcloud_cdn_domain: feat: cdn - support force_redirect.carry_headers ([#1394](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1394))
* resource/tencentcloud_monitor_tmp_alert_rule: Parameter check ([#1407](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1407))

BUG FIXES:

* resource/tencentcloud_monitor_tmp_exporter_integration: Optimized initialization ([#1410](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1410))

## 1.78.12 (November 17, 2022)

FEATURES:

* **New Data Source:** `tencentcloud_dnspod_records` ([#1398](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1398))
* **New Data Source:** `tencentcloud_tat_command` ([#1393](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1393))
* **New Data Source:** `tencentcloud_tat_invoker` ([#1393](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1393))
* **New Resource:** `tencentcloud_css_live_transcode_rule_attachment` ([#1401](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1401))
* **New Resource:** `tencentcloud_css_live_transcode_template` ([#1401](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1401))
* **New Resource:** `tencentcloud_css_pull_stream_task` ([#1401](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1401))
* **New Resource:** `tencentcloud_css_watermark` ([#1401](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1401))
* **New Resource:** `tencentcloud_tat_command` ([#1393](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1393))
* **New Resource:** `tencentcloud_tat_invoker` ([#1393](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1393))

BUG FIXES:

* resource/tencentcloud_gaap_proxy: fix gaap proxy retry error `DuplicatedRequest` ([#1402](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1402))

## 1.78.11 (November 15, 2022)

FEATURES:

* **New Resource:** `tencentcloud_pts_alert_channel` ([#1383](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1383))
* **New Resource:** `tencentcloud_pts_cron_job` ([#1383](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1383))
* **New Resource:** `tencentcloud_pts_file` ([#1383](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1383))
* **New Resource:** `tencentcloud_pts_job` ([#1383](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1383))
* **New Resource:** `tencentcloud_pts_project` ([#1383](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1383))
* **New Resource:** `tencentcloud_pts_scenario` ([#1383](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1383))

ENHANCEMENTS:

* resource/tencentcloud_instance: Remove maximum item limit on data disks ([#1390](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1390))

## 1.78.10 (November 11, 2022)

BUG FIXES:

* resource/tencentcloud_tcr_instance: fix: tcr - block creating while replication region same with current ([#1385](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1385))
* resource/tencentcloud_vpc_bandwidth_package: Handling api interface return delay problem ([#1379](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1379))

## 1.78.9 (November 9, 2022)

FEATURES:

* **New Data Source:** `tencentcloud_cat_node` ([#1378](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1378))
* **New Data Source:** `tencentcloud_cat_probe_data` ([#1378](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1378))
* **New Data Source:** `tencentcloud_dcdb_accounts` ([#1365](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1365))
* **New Data Source:** `tencentcloud_dcdb_databases` ([#1365](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1365))
* **New Data Source:** `tencentcloud_dcdb_instances` ([#1365](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1365))
* **New Data Source:** `tencentcloud_dcdb_parameters` ([#1365](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1365))
* **New Data Source:** `tencentcloud_dcdb_security_groups` ([#1365](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1365))
* **New Data Source:** `tencentcloud_dcdb_shards` ([#1365](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1365))
* **New Data Source:** `tencentcloud_mariadb_accounts` ([#1375](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1375))
* **New Data Source:** `tencentcloud_mariadb_security_groups` ([#1375](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1375))
* **New Resource:** `tencentcloud_cat_task_set` ([#1369](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1369))
* **New Resource:** `tencentcloud_mariadb_account` ([#1375](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1375))
* **New Resource:** `tencentcloud_mariadb_db_instances` ([#1370](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1370))
* **New Resource:** `tencentcloud_mariadb_dedicatedcluster_db_instance` ([#1370](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1370))
* **New Resource:** `tencentcloud_mariadb_hour_db_instance` ([#1370](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1370))
* **New Resource:** `tencentcloud_mariadb_log_file_retention_period` ([#1375](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1375))
* **New Resource:** `tencentcloud_mariadb_parameters` ([#1375](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1375))
* **New Resource:** `tencentcloud_mariadb_security_groups` ([#1375](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1375))

ENHANCEMENTS:

* resource/tencentcloud_ssl_free_certificate: support computed param dv_auths ([#1371](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1371))
* resource/tencentcloud_tcr_instance: fix: skip fail operation capture while add same region replica ([#1374](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1374))
* resource/tencentcloud_vpc_bandwidth_package: Handling api interface return is incorrect ([#1367](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1367))

BUG FIXES:

* resource/tencentcloud_postgresql_instance: network switching polling status check ([#1338](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1338))

## 1.78.8 (November 4, 2022)

FEATURES:

* **New Resource:** `tencentcloud_dcdb_account` ([#1351](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1351))
* **New Resource:** `tencentcloud_dcdb_hourdb_instance` ([#1351](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1351))
* **New Resource:** `tencentcloud_dcdb_security_group_attachment` ([#1351](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1351))
* **New Resource:** `tencentcloud_ses_domain` ([#1360](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1360))
* **New Resource:** `tencentcloud_ses_email_address` ([#1360](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1360))
* **New Resource:** `tencentcloud_ses_template` ([#1360](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1360))

ENHANCEMENTS:

* resource/tencentcloud_teo_application_proxy_rule: Add origin_port filed ([#1358](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1358))

## 1.78.7 (November 3, 2022)

FEATURES:

* **New Resource:** `tencentcloud_sms_sign` ([#1352](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1352))
* **New Resource:** `tencentcloud_sms_template` ([#1352](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1352))
* **New Resource:** `tencentcloud_sts_caller_identity` ([#1340](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1340))
* **New Resource:** `tencentcloud_vpc_bandwidth_package` ([#1343](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1343))
* **New Resource:** `tencentcloud_vpc_bandwidth_package_attachment` ([#1343](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1343))

ENHANCEMENTS:

* provider: Specify req client header with -ldflag ([#1318](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1318))
* resource/tencentcloud_tem_workload: support set tcr repo ([#1350](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1350))

BUG FIXES:

* resource/tencentcloud_elasticsearch_instance: status polling fix ([#1353](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1353))

## 1.78.6 (October 27, 2022)

FEATURES:

* **New Resource:** `tencentcloud_tcm_mesh` ([#1328](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1328))
* **New Resource:** `tencentcloud_tcm_cluster_attachment` ([#1328](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1328))

ENHANCEMENTS:

* resource/tencentcloud_cdn_domain: support PostMaxSize Params. ([#1329](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1329))
* resource/tencentcloud_kubernetes_node_pool: Support tag specifications ([#1317](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1317))
* resource/tencentcloud_monitor_tmp_exporter_integration: support cluster initialization ([#1320](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1320))
* resource/tencentcloud_redis_instance: support update `security_groups` ([#1336](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1336))
* resource/tencentcloud_security_group_lite_rule: support protocol template creating/updating ([#1315](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1315))

BUG FIXES:

* data-source/tencentcloud_security_groups: support protocol template readings ([#1315](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1315))
* resource/tencentcloud_cdn_domain: Fix testing domain verification failed. ([#1329](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1329))
* resource/tencentcloud_kubernetes_node_pool: adjust `security_group_ids` type to unordered ([#1321](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1321))
* resource/tencentcloud_teo_zone: change tag description from zoneName to zoneId ([#1326](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/1326))

## 1.78.5 (October 17, 2022)

COMMON:
* support env `PROVIDER_ASSUME_ROLE_SESSION_DURATION`

ENHANCEMENTS:
* resource `tencentcloud_ckafka_instance` support tag change
* resource `tencentcloud_gaap_layer4_listener` support read proxy_id
* resource `tencentcloud_gaap_layer7_listener` support read proxy_id
* resource `tencentcloud_vpn_gateway` support change prepaid_renew_flag
* data source `tencentcloud_gaap_layer4_listener` support read proxy_id
* data source `tencentcloud_gaap_layer7_listener` support read proxy_id

BUGFIXES:
* fix ckafka backend change
* fix clb unit test

## 1.78.4 (October 12, 2022)

COMMON:

* Add default teo variable for testcases
* Fix teo, dnspod, cdn testcases

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_node_pool` support hostname and instance name

BUGFIXES:

* resource `tencentcloud_ssl_free_certificate` optimize dv auth method import
* resource `tencentcloud_tem_application` support  set tcr instance id
* resource `tencentcloud_vpc` distinguish normal and docker assistant cidr

## 1.78.3 (September 29, 2022)

ENHANCEMENTS:
* resource `tencentcloud_monitor_tmp_cvm_agent` support `agent_id` param

BUGFIXES:
* resource `tencentcloud_security_group_rule` fix rule read error
* resource `tencentcloud_security_group_rule` fix one resource delete multi-rule which only description difference
* fix tke unit test

## 1.78.2 (September 27, 2022)

BUGFIXES:

* resource `tencentcloud_tcr_instance` extend replications retry timeout

## 1.78.1 (September 27, 2022)

ENHANCEMENTS:

* resource `tencentcloud_tcr_instance` support replications
* resource `tencentcloud_tcr_vpc_attachment` support cross region
* resource `tencentcloud_cos_bucket` add `force_clean` to support cleaning up all objects when a bucket is deleted

## 1.78.0 (September 21, 2022)

FEATURES:

* resource `tencentcloud_teo_xxxx` support api v2

ENHANCEMENTS:

* resource `tencentcloud_elasticsearch_instance` add `es_acl` to support whitelist and blacklist
* resource `tencentcloud_instance` add `key_ids`, deprecate `key_name`
* resource `tencentcloud_key_pair` change servcie params
* resource `tencentcloud_vpc_acl` support tag

## 1.77.11 (September 19, 2022)

ENHANCEMENTS:

* resource `tencentcloud_key_pair` support cam policy with restrict tag permission

## 1.77.10 (September 16, 2022)

ENHANCEMENTS:

* resource `tencentcloud_image` and `tencentcloud_cbs_storage` etc support cam policy with restrict tag permission (#1275)

## 1.77.9 (September 14, 2022)

FEATURES:

* New resource `tencentcloud_monitor_grafana_instance`
* New resource `tencentcloud_monitor_grafana_integration`
* New resource `tencentcloud_monitor_grafana_notification_channel`
* New resource `tencentcloud_monitor_grafana_plugin`
* New resource `tencentcloud_monitor_grafana_sso_account`

ENHANCEMENTS:

* resource `tencentcloud_vpc` and `tencentcloud_instance` support cam policy with tag permission

## 1.77.8 (August 24, 2022)

FEATURES:

* New resource `tencentcloud_cos_bucket_domain_certificate_attachment` (#1258)

ENHANCEMENTS:

* resource `tencentcloud_dnspod_record` support import (#1260)
* resource `tencentcloud_tem_application` update param desc (#1264)

BUGFIXES:

* resource `tencentcloud_monitor_tmp_instance` fix cluster agent status handle (#1261)
* resource `tencentcloud_tcaplus_idl` tcaplusdb support tdr (#1263)

COMMON:

* fix common testcases(#1259)

## 1.77.7 (August 24, 2022)

FEATURES:

* New data source `tencentcloud_cynosdb_zone_config`

ENHANCEMENTS:

feat: redis - support field: param template (#1249)
* resource `tencentcloud_instance` system support bssd (#1253)
* resource `tencentcloud_kubernetes_cluster` support internet security group modify (#1248)
* resource `tencentcloud_kubernetes_cluster_endpoint` support internet security group modify (#1248)

BUGFIXES:

* resource `tencentcloud_cynosdb_cluster` param_items not effect in create cluster (#1238)
* resource `tencentcloud_monitor_tmp_tke_cluster_agent` update irregular type (#1241)
* resource `tencentcloud_monitor_tmp_tke_global_notification` update irregular type (#1241)

COMMON:

* feat: pr lable (#1250)
* fix: auto label (#1256)
* fix: update irregular type (#1241)

## 1.77.6 (September 1, 2022)

ENHANCEMENTS:
* update teo/tem doc

## 1.77.5 (September 1, 2022)

FEATURES:

* new resource `tencentcloud_teo_zone`
* new resource `tencentcloud_teo_zone_setting`
* new resource `tencentcloud_teo_dns_record`
* new resource `tencentcloud_teo_dns_sec`
* new resource `tencentcloud_teo_load_balancing`
* new resource `tencentcloud_teo_origin_group`
* new resource `tencentcloud_teo_rule_engine`

## 1.77.4 (August 24, 2022)

ENHANCEMENTS:

* support TencentCloud Prometheus

## 1.77.3 (August 16, 2022)

ENHANCEMENTS:

* resource `tencentcloud_tem_gateway` output `vip` and `clb_id`

DEPRECATED:

* resource `tencentcloud_clb_target_group` deprecated `target_group_instances`

## 1.77.2 (August 15, 2022)

FEATURES:

* new resource `tencentcloud_cdn_domain_verifier`

ENHANCEMENTS:

* resource `tencentcloud_tem_workload` discard `reserved` config change
* resource `tencentcloud_tem_log_config` and `tencentcloud_tem_scale_rule` add param `workload_id`
* resource `tencentcloud_instances_set` add hint for not support change `instance_count`

## 1.77.1 (August 12, 2022)

ENHANCEMENTS:

* resource `tencentcloud_tem_workload` support update
* resource `tencentcloud_security_group_rule` support policy index

## 1.77.0 (August 11, 2022)

FEATURES:

* new resource `tencentcloud_tem_environment`
* new resource `tencentcloud_tem_application`
* new resource `tencentcloud_tem_workload`
* new resource `tencentcloud_tem_app_config`
* new resource `ttencentcloud_tem_log_config`
* new resource `tencentcloud_tem_scale_rule`
* new resource `tencentcloud_tem_gateway`

ENHANCEMENTS:

* resource `tencentcloud_cdn_domain` support `tls_versions`
* resource `tencentcloud_cbs_storage` skip prepaid upgrade testcase
* resource `tencentcloud_cos_bucket` make policy header to schema set
* resource `tencentcloud_cos_bucket_policy` fix testcase
* resource `tencentcloud_kubernetes_cluster` add node ready checking
* resource `tencentcloud_kubernetes_cluster_endpoint` add node ready checking
* resource `tencentcloud_kubernetes_node_pool` support multi_zone_subnet_policy modify
* resource `tencentcloud_postgresql_instance` increase waiting timeout
* resource `tencentcloud_sqlserver_publish_subscribe` optimize testcase sweepe
* service `service_tencentcloud_as` cancel deprecated asg relative api

BUGFIXES:

* resource `tencentcloud_postgresql_instance` fix npe

## 1.76.4 (August 1, 2022)

ENHANCEMENTS:

* datasource `tencentcloud_cfs_file_systems` output support `mount_ip`

BUGFIXES:

* resource`tencentcloud_as_scaling_group` clear as retry capture

## 1.76.3 (July 25, 2022)

ENHANCEMENTS:

* resource `tencentcloud_clb_listener` support `target_type`
* datasource `tencentcloud_images ` support `instance_type`

DEPRECATED:

* resource `tencentcloud_kubernetes_cluster` deprecated `as_enabled`

BUGFIXES:

* resource `tencentcloud_kubernetes_cluster` fix `cluster_level` desc

## 1.76.2 (July 22, 2022)

FEATURES:

* New resource `tencentcloud_monitor_tmp_scrape_job`
* New resource `tencentcloud_monitor_tmp_tke_alert_policy`
* New resource `tencentcloud_monitor_tmp_exporter_integration`

DEPRECATED:

* resource `tencentcloud_kubernetes_cluster_endpoint`, `tencentcloud_kubernetes_cluster` the
  argument `managed_cluster_internet_security_policies` was deprecated

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_cluster` and data source `tencentcloud_kubernetes_clusters` support output private
  kube config

BUGFIXES:

* resource `tencentcloud_kubernetes_cluster` fix cluster intranet subnet modification limit
* resource `tencentcloud_security_group_lite_rule`, `tencentcloud_monitor_tmp_cvm_agent`
  , `tencentcloud_monitor_tmp_instance` fix example codes.

COMMON:

* resource `tencentcloud_postgresql_instance` fix testcases.

## 1.76.1 (July 21, 2022)

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_cluster`, `tencentcloud_kubernetes_cluster_endpoint`
  support `cluster_internet_security_group`
* resource `tencentcloud_cdn_domain` extend multiple unsupport params.

## 1.76.0 (July 13, 2022)

ENHANCEMENTS:

* resource `tencentcloud_clb_attachment` support `eni_ip` for `target`

BUGFIXES:

* resource `tencentcloud_gaap_realserver` fix tag for loop mismatch

COMMON:

* document - add data type for every argument description
* fix clb, cls, vod, key-pair testcases
* service `service_tencentcloud_cvm` add max thread concurrent num

## 1.75.7 (July 08, 2022)

ENHANCEMENTS:

* resource `tencentcloud_monitor_tmp_alert_rule` update docs

## 1.75.6 (July 08, 2022)

FEATURES:

* new resource `tencentcloud_monitor_tmp_instance`
* new resource `tencentcloud_monitor_tmp_cvm_agent`
* new resource `tencentcloud_monitor_tmp_alert_rule`
* new resource `tencentcloud_monitor_tmp_recording_rule`
* new resource `tencentcloud_monitor_tmp_tke_template`
* new resource `tencentcloud_kubernetes_cluster_endpoint`

ENHANCEMENTS:

* resource `tencentcloud_ssl_certificate` support tags

## 1.75.5 (July 01, 2022)

ENHANCEMENTS:

* datasource `tencentcloud_instances_set`  hard check `instance_count` either equal
* datasource `tencentcloud_cbs_storages_set` hard check `disk_count` either equal

## 1.75.4 (July 01, 2022)

FEATURES:

* new datasource `tencentcloud_instances_set`
* new datasource `tencentcloud_cbs_storages_set`

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_cluster` support encrypt disk

## 1.75.3 (June 30, 2022)

FEATURES:

* new datasource `tencentcloud_ckafka_instances`

ENHANCEMENTS:

* resource `tencentcloud_instance_set` add timeout
* datasource `tencentcloud_instance` add field `instances_ids`
* resource `tencentcloud_ckafka_instance` support vips

## 1.75.2 (June 29, 2022)

FEATURES:

* new resource `tencentcloud_cbs_storage_set`

## 1.75.1 (June 28, 2022)

ENHANCEMENTS:

* resource `tencentcloud_instance_set` update docs

## 1.75.0 (June 28, 2022)

FEATURES:

* new resource `tencentcloud_instance_set`

ENHANCEMENTS:

* resource `tencentcloud_monitor_alarm_policy` optimize errors

## 1.74.3 (June 24, 2022)

ENHANCEMENTS:

* resource `tencentcloud_cls_machine_group` update tags
* resource `tencentcloud_kubernetes_node_pool` fix node count

## 1.74.2 (June 23, 2022)

FEATURES:

* new resource `tencentcloud_ssl_free_certificate`
* new resource `tencentcloud_clb_snat_ip`

ENHANCEMENTS:

* resource `tencentcloud_vpn_connection` support healthcheck
* resource `tencentcloud_postgresql_instance` add plan modify retry
* resource `tencentcloud_cdn_domain` support follow redirect and authentication
* resource `tencentcloud_redis_instance` update redis_shard_num
* resource `tencentcloud_clb_instance` support snat pro and snat ips
* resource `tencentcloud_cls_topic` update storage_type, auto_split, partition_count
* resource `tencentcloud_cls_logset` update tags

## 1.74.1 (June 20, 2022)

ENHANCEMENTS:

* resource `tencentcloud_kms_key` optimize kms key disable logic

## 1.74.0 (June 20, 2022)

FEATURES:

* New resource `tencentcloud_cdn_url_purge`
* New resource `tencentcloud_cdn_url_push`

ENHANCEMENTS:

* resource `tencentcloud_cam_role` support tags
* resource `tencentcloud_cam_user` support tags
* resource `tencentcloud_kms_key` optimize delete logic
* resource `tencentcloud_kubernetes_cluster` optimize cluster_level logic

COMMON:

* fix vod testcases

## 1.73.3 (June 16, 2022)

FEATURES:

* datasourcce `tencentcloud_user_info` support get userinfo

BUGFIXES:

* resource `resource_tc_cls_index` fix crash error
* resource `resource_tc_clb_instance` fix tag error
* resource `resource_tc_kubernetes_cluster` fix cluster_level error

COMMON:

* fix cls, cfs testcases

## 1.73.2 (June 13, 2022)

BUGFIXES:

* resource `tencentcloud_clb_attachment` disable target diff if using TCP-SSL listener
* resource `tencentcloud_instance` support prepaid/postpaid charge type dual modify

COMMON:

* fix gaap, cvm datasource testcases

## 1.73.1 (June 09, 2022)

ENHANCEMENTS:

* resource `tencentcloud_instance` remove unnecessary retries and optimize state management on launch failed cvm

## 1.73.0 (June 08, 2022)

FEATURES:

* new resource `tencentcloud_cls_index`
* new resource `tencentcloud_lighthouse_instance`

ENHANCEMENTS:

* resource `tencentcloud_dc_gateway` add protect sleep
* resource `tencentcloud_kubernetes_cluster` fix docs

COMMON:

* fix testcases

## 1.72.8 (June 02, 2022)

ENHANCEMENTS:

* resource `tencentcloud_mysql_instance` support device type
* resource `tencentcloud_mysql_readonly_instance` support device type

COMMON:

* resource `tencentcloud_tcaplus_table` add test sweeper
* chore: add changelog draft script

## 1.72.7 (May 31, 2022)

BUGFIXES:

* resource `tencentcloud_mysql_readonly_instance` make zone arguments optional.

## 1.72.6 (May 31, 2022)

FEATURES:

* new resource `tencentcloud_cam_role_sso`

ENHANCEMENTS:

* resource `tencentcloud_mysql_instance`  support fast upgrade and param template id

BUGFIXES:

* resource `tencentcloud_sqlserver_instance` disable recycle
* resource `tencentcloud_mysql_instance` / `tencentcloud_mongodb_instance` cancel validate required db pwd.

COMMON:

* testcases fix: dnspod, cvm, kafka, as, tcr, sqlserver, cbs, cos, cfs, mongo, tke

## 1.72.5 (May 20, 2022)

ENHANCEMENTS:

* resource `tencentcloud_ckafka_instance` support upgrade config
* resource `tencentcloud_kubernetes_cluster` support log agent/audit/event persistence

## 1.72.4 (May 18, 2022)

ENHANCEMENTS:

* resource `tencentcloud_vpn_gateway` add `cdc_id` and `max_connection`
* resource `tencentcloud_private_dns_zone` support tag change
* resource `tencentcloud_route_table_entry` add disabled argument
* datasource `tencentcloud_cbs_storages` support more query filters

BUGFIXES:

* resource `tencentcloud_private_dns_record` fix delete error

COMMON:

* fix testcases

## 1.72.3 (May 13, 2022)

ENHANCEMENTS:

* resource `tencentcloud_gaap_http_rule` add `sni` and `sni_switch`
* resource `tencentcloud_scf_function` remove memSize validate

## 1.72.2 (May 11, 2022)

BUGFIXES:

* resource `tencentcloud_mysql_readonly_instance` skip monitor check.
* resource `tencentcloud_clb_listener_rule` fix domain and port update.
* resource `tencentcloud_kubernetes_cluster` add cluster level modify retry.

COMMON:

* fix testcases

## 1.72.1 (May 6, 2022)

ENHANCEMENTS:

* resource `tencentcloud_scf_function` support java handler
* resource `tencentcloud_mongodb_instance` fix engine version doc

## 1.72.0 (May 5, 2022)

FEATURES:

* New datasource `tencentcloud_kubernetes_cluster_common_names`
* New resource `tencentcloud_cam_oidc_sso`

ENHANCEMENTS:

* resource `tencentcloud_postgresql_instance` support data transparent encryption

BUGFIXES:

* resource `tencentcloud_vpn_ssl_client` fix delete failed
* resource `tencentcloud_vpn_ssl_server` fix create duplicate instance

## 1.71.0 (April 24, 2022)

FEATURES:

* New datasource `tencentcloud_postgresql_xlogs`

ENHANCEMENTS:

* resource `tencentcloud_postgresql_instance` support backup plan
* resource `tencentcloud_redis_instance` support `replica_zone_ids` modify

BUGFIXES:

* resource `tencentcloud_eks_cluster` fix LB modified errors

## 1.70.3 (April 22, 2022)

ENHANCEMENTS:

* resource `tencentcloud_ckafka_instance` remove validate check
* resource `tencentcloud_tcr_vpc_attachment` modify doc

## 1.70.2 (April 21, 2022)

FEATURES:

* New resource `tencentcloud_cls_config_extra`

BUGFIXES:

* resource `tencentcloud_cls_topic` plan change
* resource `tencentcloud_cls_config` create failed when log_type is full regex

## 1.70.1 (April 19, 2022)

BUGFIXES:

* resource `tencentcloud_ckafka_instance` plan change
* resource `tencentcloud_mysql_instance`  ignore lowercase when 8.0

## 1.70.0 (April 19, 2022)

BUGFIXES:

* New resource `tencentcloud_cls_config`
* New resource `tencentcloud_cls_config_attachment`

## 1.69.0 (April 19, 2022)

FEATURES:

* New resource `tencentcloud_cls_config`
* New resource `tencentcloud_cls_config_attachment`

## 1.68.0 (April 12, 2022)

FEATURES:

* New resource `tencentcloud_cls_cos_shipper`
* New resource `tencentcloud_postgresql_readonly_instance`
* New resource `tencentcloud_postgresql_readonly_group`
* New resource `tencentcloud_postgresql_readonly_attachment`

BUGFIXES:

* resource `tencentcloud_elasticsearch_instance` support `web_node_type_info` modify

COMMON

* update es, eks(ci), myqsl, redis, sqlserver sweepers

## 1.67.0 (April 8, 2022)

FEATURES:

* New resource `tencentcloud_cls_logset`
* New resource `tencentcloud_cls_topic`
* New resource `tencentcloud_cls_machine_group`

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_cluster` support cluster level and auto upgrade level settings.
* resource `tencentcloud_emr_cluster` support `extend_fs_field`

BUGFIXES:

* resource `tencentcloud_emr_cluster` clear metadb
* resource `tencentcloud_instance` ignore local disk describe status.

COMMON:

* fix partial testcases gaap, tke, clb, cam, emr, e.g.

## 1.66.3 (March 31, 2022)

BUGFIXES:

* resource `tencentcloud_sqlserver_instance` block name mistake in example block.
* resource `tencentcloud_kubernetes_node_pool`, `tencentcloud_kubernetes_cluster_attachment`
  fix `data_disk.disk_partition` usage description.

ENHANCEMENTS:

* resource `tencentcloud_as_scaling_config` support prepaid and spot charge type.
* resource `tencentcloud_ckafka_config` support multi zones config.

## 1.66.2 (March 31, 2022)

BUGFIXES:

* resource `tencentcloud_instance` optimize charge type modification.
* resource `tencentcloud_postgresql_instance` fix instance upgrade polling error.
* resource `tencentcloud_vpn_ssl_client`
* fix vpc, vpngw, clb and eni testcases.

ENHANCEMENTS:

* resource `tencentcloud_sqlserver_instance`,`tencentcloud_sqlserver_readonly_instance` support prepaid charge type

## 1.66.1 (March 24, 2022)

BUGFIXES:

* resource `tencentcloud_elasticsearch_instance` testcases.

ENHANCEMENTS:

* resource `tencentcloud_mysql_readonly_instance` support cross-zone purchase and query remoteRoZones
* resource `tencentcloud_scf_function` support cfs config and fix testcases
* resource `tencentcloud_eks_cluster` support output kubeconfig
* resource `tencentcloud_emr_cluster` support `root_size` and `sg_id` params.

## 1.66.0 (March 22, 2022)

FEATURES:

* new datasource `tencentcloud_eks_cluster_credential`

ENHANCEMENTS:

* resource `tencentcloud_eks_cluster` support public and internal load balancer config

BUGFIXES:

* resource `tencentcloud_instance` charge type update then polling status mismatch

## 1.65.2 (March 18, 2022)

ENHANCEMENTS:

* resource `tencentcloud_route_entry` extend next type
* resource `tencentcloud_clb_instance` create with tag params
* resource `tencentcloud_emr_cluster` support need_master_wan param
* resource `tencentcloud_cos_bucket` support lifecycle rule id and delete marker
* resource `tencentcloud_mysql_readonly_instance` support zone param

BUGFIXES:

* testcases: cvm, cynosdb, eks, mysql, vpn
* resource `tencentcloud_emr_cluster` destroy
* resource `tencentcloud_mysql_instance` init
* resource `tencentcloud_redis_instance` reset pwd retry and comment

## 1.65.1 (March 14, 2022)

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_node_pool` data_disk support `delete_with_instance` option
* resource `tencentcloud_as_scaling_configs` data_disk support `delete_with_instance` option
* datasource `tencentcloud_as_scaling_configs` data_disk add `delete_with_instance` field

BUGFIXES:

* testcases fixes: CVM,SG,MySQL,EKS
* testcases tke use preset tke and tlinux image instead

## 1.65.0 (March 7, 2022)

FEATURES:

* new datasource `tencentcloud_dnspod_record`
* new resource `tencentcloud_vpn_ssl_server`
* new resource `tencentcloud_vpn_ssl_client`

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_node_pool` support `cam_role_name` as parameter

## 1.64.0 (March 7, 2022)

FEATURES:

* new datasource `tencentcloud_mysql_default_params`
* new resource `tencentcloud_dayu_cc_policy_v2`

BUGFIXES:

* resource `tencentcloud_mysql_instance` init database after create
* resource `tencentcloud_emr_cluster` pay_mode can not use 0
* testcases - add vpngw sweeper and fix cvm sweeper

ENHANCEMENTS:

* resource `tencentcloud_cynosdb_cluster` support `param_items` as parameter

## 1.63.0 (March 2, 2022)

FEATURES:

* new resource `tencentcloud_user_info`

BUGFIXES:

* resource `tencentcloud_instance` support data disk `delete_with_instance` on cvm remove
* resource `tencentcloud_redis_instance` support `replicas_read_only

## 1.62.0 (March 1, 2022)

FEATURES:

* new resource `tencentcloud_dayu_ddos_policy_v2`

BUGFIXES:

* datasource `tencentcloud_postgresql_instances` fix query id parameter
* resource `tencentcloud_cynosdb_cluster` make `storage_limit` optional and update description

DEPRECATED:

* datasource `tencentcloud_redis_zone_config`: The argument `mem_sizes` was deprecated, use `shard_memories` instead.

## 1.61.13 (February 24, 2022)

ENHANCEMENTS:

* resource `tencentcloud_redis_instance` support `auto_renew_flag`.

## 1.61.12 (February 23, 2022)

ENHANCEMENTS:

* resource `tencentcloud_instance` fix testcase default image

## 1.61.11 (February 23, 2022)

ENHANCEMENTS:

* resource `tencentcloud_instance` support System Disk size and DiskType modification

## 1.61.10 (February 22, 2022)

ENHANCEMENTS:

* resource `tencentcloud_instance` support Data Disk size modification
* resource generic gaap support client_ip_method

## 1.61.9 (February 21, 2022)

ENHANCEMENTS:

* resource `tencentcloud_cos_bucket` support Multi-AZ bucket import, Non-Current Version Lifecycle and ACL body

## 1.61.8 (February 17, 2022)

ENHANCEMENTS:

* resource `tencentcloud_as_scaling_group` support Service Settings

BUGFIXES:

* resource `tencentcloud_mysql_privilege` update global privileges check
* Fix testcases COS Bucket

## 1.61.7 (February 15, 2022)

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_node_pool` support CVM instance charge type

BUGFIXES:

* Fix testcases CamPolicies
* Fix testcases CamRole
* Fix testcases TcaplusDB
* Fix testcases EKSCluster

## 1.61.6 (February 5, 2022)

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_node_pool` support set `multi_zone_subnet_policy`

## 1.61.5 (January 29, 2022)

ENHANCEMENTS:

* resource `tencentcloud_redis_instance` remove change redis mem_size limit

## 1.61.4 (January 27, 2022)

ENHANCEMENTS:

* resource `tencentcloud_dayu_eip` support create new ddos eip rules
* resource `tencentcloud_dayu_l4_rule_v2` support create new ddos l4 rules
* resource `tencentcloud_dayu_l7_rule_v2` support create new ddos l7 rules
* data source `tencentcloud_dayu_eip` support query new ddos eip rules
* data source `tencentcloud_dayu_l4_rules_v2` support query new ddos l4 rules
* data source `tencentcloud_dayu_l7_rules_v2` support query new ddos l7 rules

## 1.61.3 (January 25, 2022)

ENHANCEMENTS:

* resource `tencentcloud_as_scaling_config` support `instance_name_settings` and `cam_role_name`.

## 1.61.2 (January 21, 2022)

COMMON:

* testcase fix `defaultGaapProxyId`

ENHANCEMENTS:

* resource `tencentcloud_postgresql_instance` support `db_node_set` for multiple available zone.

## 1.61.1 (January 14, 2022)

ENHANCEMENTS:

* resource `tencentcloud_monitor_alarm_policy` support policy tag

## 1.61.0 (January 11, 2022)

COMMON:

* add env variable `TENCENTCLOUD_READ_RETRY_TIMEOUT`  and `TENCENTCLOUD_WRITE_RETRY_TIMEOUT`

BUGFIXES:

* resource `tencentcloud_gaap_proxy` fix destroy error

## 1.60.27 (January 11, 2022)

COMMON:

* add `terraform {}` block to provider doc example
* add `pre-commit` status check

BUGFIXES:

* resource `tencentcloud_scf_function` add COS trigger `bucketUrl` field
* resource `tencentcloud_eip` fix NAT gateway detach error
* service `tencentcloud_gaap` destroy error

ENHANCEMENTS:

* resource `tencentcloud_vpc` add assistant CIDR
* resource `tencentcloud_cos_bucket` support bucket replication

## 1.60.26 (January 6, 2022)

COMMON:

* fix testcases

ENHANCEMENTS:

* resource `tencentcloud_vpc` support vpc assistant CIDR
* resource `tencentcloud_eip` support create HighQualityEIP
* resource `tencentcloud_instance` support charge type update

## 1.60.25 (December 27, 2021)

BUGFIXES:

* Resource `tencentcloud_instance` rollback instance_charge_type_prepaid_period

## 1.60.24 (December 27, 2021)

COMMON:

* add testcases basic required resource

BUGFIXES:

* fix service mongodb status query

## 1.60.23 (December 24, 2021)

ENHANCEMENTS:

* resource `tencentcloud_vpn_gateway` improve timeout min

BUGFIXES:

* Resource `tencentcloud_postgresql_instance` fix computed attribute uid get 0
* Resource `tencentcloud_monitor_alarm_policy` fix filter when has event condition

## 1.60.22 (December 21, 2021)

ENHANCEMENTS:

* resource `tencentcloud_mysql_readonly_instance` fix create prepaid instance

## 1.60.21 (December 20, 2021)

FEATURES:

* **New Resource**: `tencentcloud_private_dns_record`

ENHANCEMENTS:

* resource `tencentcloud_postgresql_instance` add computed attribute uid

## 1.60.20 (December 16, 2021)

ENHANCEMENTS:

* resource `tencentcloud_monitor_alarm_policy` fix filter bug
* resource `tencentcloud_clb_customized_config` remove content length validate
* resource `tencentcloud_ckafka_instance` modify docs

## 1.60.19 (December 16, 2021)

FEATURES:

* **New Resource**: `tencentcloud_ckafka_instance`

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_addon_attachment` add tcr example values
* resource  `tencentcloud_instance` support cvm import params
* datasource `tencentcloud_kubernetes_charts`fix domain 404

## 1.60.18 (December 10, 2021)

ENHANCEMENTS:

* resource `tencentcloud_redis_instance` support no-auth access

## 1.60.17 (December 9, 2021)

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_cluster` modify docs.
* resource `tencentcloud_kubernetes_cluster` ignore `ImageId` JSON if user make it empty.
* resource `tencentcloud_cbs_storage` remove `CLOUD_BASIC` type
* resource `tencentcloud_postgresql_instance` change available versions
* resource `tencentcloud_postgresql_instance` support `max_standby_*` params

## 1.60.16 (December 3, 2021)

ENHANCEMENTS:

* resource `tencentcloud_kubernetes_addon_attachment` fix tke_addon - npe error

## 1.60.15 (December 3, 2021)

FEATURES:

* **New Datasource** `tencentcloud_kubernetes_charts`
* **New Resource** `tencentcloud_kubernetes_addon_attachment`

## 1.60.14 (December 2, 2021)

ENHANCEMENTS:

* resource `resource_tc_monitor_alarm_policy` support change name and remark

## 1.60.13 (December 1, 2021)

ENHANCEMENTS:

* resource `resource_tc_security_group_rule`  remove cidr validate
* resource `resource_tc_monitor_alarm_policy` support manage tke rules
* datasource `tencentcloud_instance_types ` support filter by charge type

## 1.60.12 (November 30, 2021)

FEATURES:

* **New Resource**: `tencentcloud_private_dns_zone`

ENHANCEMENTS:

* resource `tencentcloud_monitor_policy_binding_object`  support import
* resource `tencentcloud_kubernetes_cluster` support modify enableClusterDeletionProtection

## 1.60.11 (November 25, 2021)

FEATURES:

* **New Resource**: `tencentcloud_emr_cluster`
* **New Resource**: `tencentcloud_clb_customized_config`

BUGFIXES:

* Resource `tencentcloud_kubernetes_node_pool` support launch config attributes modify

## 1.60.10 (November 22, 2021)

FEATURES:

* **New Resource**: `tencentcloud_eks_container_instance`
* **New Resource**: `tencentcloud_dnspod_domain_instance`

BUGFIXES:

* Resource `tencentcloud_monitor_alarm_policy` remove alarm policy binding check

DEPRECATED:

* Disk type `CLOUD_BASIC` which referenced by CVM/TKE/CBS was no longer available

## 1.60.9 (November 16, 2021)

BUGFIXES:

* Resource `tencentcloud_instance` omit `InstanceMarketOptions` field if spot arguments empty

## 1.60.8 (November 16, 2021)

BUGFIXES:

* Resource `tencentcloud_clb_log_set` remove doc

## 1.60.7 (November 16, 2021)

FEATURES:

* **New Resource**: `tencentcloud_clb_log_set`
* **New Resource**: `tencentcloud_clb_log_topic`

ENHANCEMENTS:

* resource `tencentcloud_clb_instance`  support set log
* resource `tencentcloud_instance` unlimited spot charge argument

## 1.60.6 (November 12, 2021)

FEATURES:

* **New Resource**: `tencentcloud_monitor_policy_binding_object`

ENHANCEMENTS:

* fix resource `tencentcloud_kubernetes_node_pool` plan change
* fix resource `tencentcloud_monitor_alarm_policy` import lacking eventditions

DEPRECATED:

* Resource: `tencentcloud_monitor_binding_object` replaced by tencentcloud_monitor_policy_binding_object

## 1.60.5 (November 8, 2021)

BUGFIXES:

*Resource `tencentcloud_availability_zones_by_product` add dependency

## 1.60.4 (November 8, 2021)

FEATURES:

* **New Resource**: `tencentcloud_vod_sub_application`
* **New Resource**: `tencentcloud_availability_zones_by_product`

ENHANCEMENTS:

* resource `tencentcloud_clb_instance` support set load_balancer_pass_to_target

DEPRECATED:

* Resource: `tencentcloud_availability_zones` replaced by `tencentcloud_availability_zones_by_product`

## 1.60.3 (November 3, 2021)

BUGFIXES:

* Resource `tencentcloud_tcr_repository` fix inaccurate document and example usage

ENHANCEMENTS:

* Resource `tencentcloud_postgresql_instance` support modifying `security_groups`

## 1.60.2 (November 1, 2021)

BUGFIXES:

* resource `tencentcloud_tcr_instance` fix document format error

## 1.60.1 (October 29, 2021)

ENHANCEMENTS:

* resource `tencentcloud_tcr_instance` support security policies

## 1.60.0 (October 28, 2021)

FEATURES:

* **New Resource**: `tencentcloud_scf_layer`

ENHANCEMENTS:

* resource/tencentcloud_scf_function: Add `layers` argument

## 1.59.20 (October 27, 2021)

ENHANCEMENTS:

* resource `tencentcloud_redis_instance` support multi replica zone ids

## 1.59.19 (October 27, 2021)

FEATURES:

* **New Resource**: `tencentcloud_monitor_alarm_policy`

DEPRECATED:

* Resource: `tencentcloud_monitor_policy_group` replaced by `tencentcloud_monitor_alarm_policy`

## 1.59.18 (October 25, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_instance` support instance reset and pay as you go
* Resource `tencentcloud_vpc` support `default_route_table_id`
* DataSource `tencentcloud_vpc_route_tables` add example usage for fetching default route table

DEPRECATED:

* Resource `tencentcloud_instance` argument `instance_count` was deprecated, replace by built-in `count`

## 1.59.17 (October 20, 2021)

BUGFIXES:

* Resource `tencentcloud_vod_xxx` fix resource related to vod, while user create vod resource use sub app id.

## 1.59.16 (October 19, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_clb_instance` support `bandwidth_package_id`

COMMON:

* Resource `tencentcloud_eks_cluster` make essential argument required

## 1.59.15 (October 19, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_scf_function` support scf create by image
* Resource `tencentcloud_security_group_lite_rule` ingress/egress policy support security group ID, address template as
  source

BUGFIXES:

* Resource `tencentcloud_clb_listener` skip empty import set
* Resource `tencentcloud_address_template_group` fix incorrect field `addresses` to `template_ids` in example usage.

## 1.59.14 (October 19, 2021)

BUGFIXES:

* Resource `tencentcloud_kubernetes_auth_attachment` fix official document synchronous error.
* Resource `tencentcloud_elasticsearch_instance` make zone and subnet optional for multi az case.

COMMON:

* Remove outdated go.sum

## 1.59.13 (October 18, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_kubernetes_cluster` support argument `extension_addon`

## 1.59.12 (October 15, 2021)

FEATURES:

* **New Resource**: `tencentcloud_kubernetes_auth_attachment`
* **New Resource**: `tencentcloud_tdmq_instance`
* **New Resource**: `tencentcloud_tdmq_namespace`
* **New Resource**: `tencentcloud_tdmq_topic`
* **New Resource**: `tencentcloud_tdmq_role`
* **New Resource**: `tencentcloud_tdmq_namespace_role_attachment`

## 1.59.11 (October 12, 2021)

FEATURES:

* **New Resource**: `tencentcloud_eks_cluster`
* **New Data Source**: `tencentcloud_eks_clusters`

ENHANCEMENTS:

* Resource `tencentcloud_tcr_vpc_attachment` support argument `region_name`

## 1.59.10 (October, 9, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_cos_bucket_object` support cos object tags.

BUGFIXES:

* Resource `tencentcloud_kubernetes_cluster` update authentication options immediately after create

## 1.59.9 (October, 9, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_kubernetes_cluster` cluster support authentication options.

## 1.59.8 (October, 6, 2021)

BUGFIXES:

* Resource `tencentcloud_tcr_vpc_attachment` pass region_id for delete if provided.

## 1.59.7 (October, 5, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_tcr_vpc_attachment` tcr vpc attachment support region.

## 1.59.6 (September 28, 2021)

BUGFIXES:

* Resource `tencentcloud_kubernetes_node_pool` fix termination_policies argument mistake

DEPRECATED:

* DataSource `data_source_tc_cam_user_policy_attachments` deprecate `user_id` and replaced by `user_name`

## 1.59.5 (September 28, 2021)

ENHANCEMENTS:

* TestCase `TestAccTencentCloudTkeNodePoolResource` extend relative auto-scaling group arguments in node pool resource

DEPRECATED:

* Resource `tencentcloud_cam_group_membership` argument `user_ids` was deprecated, replace by `user_names`
* Resource `tencentcloud_cam_user_policy_attachment` argument `user_id` was deprecated, replace by `user_name`

## 1.59.4 (September 24, 2021)

BUGFIXES:

* Resource `tencentcloud_tcr_instance` support modify tags
* Resource `service_tencentcloud_postgresql` support security group
* Resource `service_tencentcloud_monitor` fix binding policy query limit
* Resource `tencentcloud_api_gateway_api` fix destroy limitNumber

## 1.59.2 (September 18, 2021)

BUGFIXES:

* Add missing AuthorizationTransport.Token field

ENHANCEMENTS:

* Resource `tencentcloud_cdn_domain` support `ipv6_access_switch` config

## 1.59.1 (September 15, 2021)

BUGFIXES:

* Resource `tencentcloud_instance` remove last_update_status judge
* Resource `tencentcloud_instance` fix DescribeInstanceById return LatestOperationState
* Resource `tencentcloud_clb_attachment` clb_attachment check instances before unbind target groups

CHORE:

* Define `TENCENTCLOUD_APPID` Environment variable for testing appid.
* Format code style

## 1.58.5 (September 7, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_redis_backup_config ` change backup_period to optional
* Resource `tencentcloud_scf_function` enable public net config and eip config
* Resource `tencentcloud_cos_bucket` support MAZ, ACL XML body, Origin-Pull rules and origin domain rules

## 1.58.4 (Aug 24, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_kubernetes_node_pool` support `backup_instance_type` for `auto_scaling_config`

## 1.58.3 (Aug 18, 2021)

ENHANCEMENTS:

* Extend kubernetes node instance disk allow types
* Resource `tencentcloud_kubernetes_cluster_attachment` add `disk_partition` field

BUGFIXES:

* Resource `tencentcloud_kubernetes_cluster` fix `auto_format_and_mount` default value to `false`

## 1.58.2 (Aug 18, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_kubernetes_cluster` tke instance creation - add RuntimeVersion param
* Resource `tencentcloud_kubernetes_cluster` extend worker instance data disk mount settings

## 1.58.1 (Aug 17, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_mysql_instance` add cpu params for mysql

BUG FIXES:

* Resource `tencentcloud_instance` fix read cvm data_disks bug

## 1.58.0 (Aug 11, 2021)

FEATURES:

* **New Resource**: `tencentcloud_nat_gateway_snat`
* **New Data Source**: `tencentcloud_nat_gateway_snat`

## 1.57.3 (Aug 10, 2021)

BUG FIXES:

* DataSource `data_source_tc_elasticsearch_instances` skip kibana node info record after elastic search instance create
* Resource `tencentcloud_postgresql_instance` skip kibana node info record after elastic search instance create

## 1.57.2 (Aug 7, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_postgresql_instance` root_user setting support

## 1.57.1 (Aug 5, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_vpn_gateway_route` add example docs

## 1.57.0 (July 28, 2021)

FEATURES:

* **New Resource**: `tencentcloud_vpn_gateway_route`
* **New Data Source**: `tencentcloud_vpn_gateway_route`

## 1.56.15 (July 07, 2021)

BUG FIXES:

* Resource `tencentcloud_tc_kubernetes_cluster` filter the request field of *bandwidth_package_id* when it is null
* Resource `tencentcloud_tc_kubernetes_node_pool` filter the request field of *bandwidth_package_id* when it is null

## 1.56.14 (July 06, 2021)

BUG FIXES:

* Resource `tencentcloud_tc_clb_listener` exec the plan will lead the resource rebuild.

ENHANCEMENTS:

* Resource `tencentcloud_elasticsearch_instance` create **ES** cluster add new parametes of *web_node_type_info*.
* Resource `tencentcloud_tc_instance` add *instance_count* to support create multiple consecutive name of instance
* Resource `tencentcloud_tc_kubernetes_cluster` supports change *internet_max_bandwidth_out*
* Resource `tencentcloud_tc_instance` create cvm instance add *bandwidth_package_id*

## 1.56.13 (July 02, 2021)

BUG FIXES

* Resource `TkeCvmCreateInfo.data_disk.disk_type` support CLOUD_HSSD and CLOUD_TSSD

## 1.56.12 (July 02, 2021)

BUG FIXES

* Resource `TkeCvmCreateInfo.data_disk.disk_type` support CLOUD_HSSD

## 1.56.11

BUG FIXES

* Resource `tencentcloud_kubernetes_cluster` fix create cluster without *desired_pod_num* in tf, then crash
* Resource `tencentcloud_kubernetes_cluster` fix when upgrade terraform-provider-tencentclod from v1.56.1 to newer,
  cluster_os force replacement
* Resource `tencentcloud_kubernetes_cluster` fix when upgrade terraform-provider-tencentclod from v1.56.1 to newer,
  enable_customized_pod_cidr force replace

## 1.56.10

BUG FIXES

* Resource `tencentcloud_tcr_namespace` fix create two namespace and one name is substring of another, then got an error
  about more than 1
* Resource `tencentcloud_tcr_namespace` fix create two repositories and one name is substring of another, then got an
  error about more than 1

## 1.56.9 (Jun 09, 2021)

BUG FIXES:

* Resource `tencentcloud_instance` fix words spell, in tencendcloud/tencentcloud_instance.go L45,
  data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name change to
  data.tencentcloud_availability_zones.my_favorite_zones.zones.0.name".
* Resource `tencentcloud_kubernetes_clusters` fix the description of is_non_static_ip_mode

ENHANCEMENTS:

* Resource `tencentcloud_clb_target_group` add create target group.
* Resource `tencentcloud_clb_instance` add internal CLB supports security group.
* Resource `tencentcloud_clb_instance` add supports open and close CLB security group, default is open.
* Resource `tencentcloud_clb_instance` add external CLB create multi AZ instance.
* Resource `tencentcloud_kubernetes_cluster` add supports params of img_id to assign image.
* Resource `tencentcloud_as_scaling_group` add MultiZoneSubnetPolicy.

## 1.56.8 (May 26, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_kubernetes_cluster_attachment.worker_config` add `desired_pod_num`.
* Resource `tencentcloud_kubernetes_cluster_attachment` add `worker_config_overrides`.
* Resource `tencentcloud_kubernetes_scale_worker` add `desired_pod_num`.
* Resource `tencentcloud_kubernetes_cluster` add `enable_customized_pod_cidr`, `base_pod_num`, `globe_desired_pod_num`,
  and `exist_instance`.
* Resource `tencentcloud_kubernetes_cluster` update available value of `cluster_os`.
* Resource `tencentcloud_as_lifecycle_hook` update `heartbeat_timeout` value ranges.

## 1.56.7 (May 12, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_as_scaling_config` add `disk_type_policy`.
* Data Source `tencentcloud_as_scaling_configs` add `disk_type_policy` as result.

## 1.56.6 (May 7, 2021)

BUG FIXES:

* Resource: `tencentcloud_scf_function` filed `cls_logset_id` and `cls_logset_id` change to Computed.

## 1.56.5 (April 26, 2021)

BUG FIXES:

* Resource: `tencentcloud_kubernetes_cluster` upgrade cluster timeout from 3 to 9 minutes.

## 1.56.4 (April 26, 2021)

BUG FIXES:

* Resource: `tencentcloud_kubernetes_cluster` upgrade instances timeout depend on instance number.

## 1.56.3 (April 25, 2021)

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` add `upgrade_instances_follow_cluster` for upgrade all instances of
  cluster.

## 1.56.2 (April 19, 2021)

BUG FIXES:

* Remove `ResourceInsufficient` from `retryableErrorCode`.

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` upgrade `cluster_version` will send old `cluster_extra_args` to tke.

## 1.56.1 (April 6，2021)

BUG FIXES:

* Fix release permission denied.

## 1.56.0 (April 2，2021)

FEATURES:

* **New Resource**: `tencentcloud_cdh_instance`
* **New Data Source**: `tencentcloud_cdh_instances`

ENHANCEMENTS:

* Resource: `tencentcloud_instance` add `cdh_instance_type` and `cdh_host_id` to support create instance based on cdh.

## 1.55.2 (March 29, 2021)

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` add `node_pool_global_config` to support node pool global config setting.

## 1.55.1 (March 26, 2021)

ENHANCEMENTS:

* Resource: `tencentcloud_tcr_vpc_attachment` add more time for retry.

## 1.55.0 (March 26, 2021)

FEATURES:

* **New Resource**: `tencentcloud_ssm_secret`
* **New Resource**: `tencentcloud_ssm_secret_version`
* **New Data Source**: `tencentcloud_ssm_secrets`
* **New Data Source**: `tencentcloud_ssm_secret_versions`

ENHANCEMENTS:

* Resource: `tencentcloud_ssl_certificate` refactor logic with api3.0 .
* Data Source: `tencentcloud_ssl_certificates` refactor logic with api3.0 .
* Resource `tencentcloud_kubernetes_cluster` add `disaster_recover_group_ids` to set disaster recover group ID.
* Resource `tencentcloud_kubernetes_scale_worker` add `disaster_recover_group_ids` to set disaster recover group ID.

## 1.54.1 (March 24, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_tcr_vpc_attachment` add `enable_public_domain_dns`, `enable_vpc_domain_dns` to set whether to
  enable dns.
* Data Source `tencentcloud_tcr_vpc_attachments` add `enable_public_domain_dns`, `enable_vpc_domain_dns`.

## 1.54.0 (March 22, 2021)

FEATURES:

* **New Resource**: `tencentcloud_kms_key`
* **New Resource**: `tencentcloud_kms_external_key`
* **New Data Source**: `tencentcloud_kms_keys`

ENHANCEMENTS:

* Resource `tencentcloud_kubernetes_cluster_attachment` add `unschedulable` to set whether the joining node participates
  in the schedule.
* Resource `tencentcloud_kubernetes_cluster` add `unschedulable` to set whether the joining node participates in the
  schedule.
* Resource `tencentcloud_kubernetes_node_pool` add `unschedulable` to set whether the joining node participates in the
  schedule.
* Resource `tencentcloud_kubernetes_scale_worker` add `unschedulable` to set whether the joining node participates in
  the schedule.

## 1.53.9 (March 19, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_tcr_instance` add `open_public_network` to control public network access.
* Resource `tencentcloud_cfs_file_system` add `storage_type` to change file service StorageType.

## 1.53.8 (March 15, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_instance` add `cam_role_name` to support binding role to cvm instance.

BUG FIXES:

* Resource `tencentcloud_instance` fix bug that waiting 5 minutes when cloud disk sold out.
* Resource: `tencentcloud_tcr_instance` fix bug that only one tag is effective when setting multiple tags.

## 1.53.7 (March 10, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_kubernetes_node_pool` add `internet_max_bandwidth_out`, `public_ip_assigned` to support
  internet traffic setting.
* Resource `tencentcloud_instance` remove limit of `data_disk_size`.

## 1.53.6 (March 09, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_eip` support `internet_max_bandwidth_out` modification.
* Resource `tencentcloud_kubernetes_cluster` add `hostname` to support node hostname setting.
* Resource `tencentcloud_kubernetes_scale_worker` add `hostname` to support node hostname setting.

## 1.53.5 (March 01, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_clb_instance` add `internet_charge_type`, `internet_bandwidth_max_out` to support internet
  traffic setting with OPEN CLB instance.
* Resource `tencentcloud_clb_rule` add `http2_switch` to support HTTP2 protocol setting.
* Resource `tencentcloud_kubernetes_cluster` add `lan_ip` to show node LAN IP.
* Resource `tencentcloud_kubernetes_scale_worker` add `lan_ip` to show node LAN IP.
* Resource `tencentcloud_kubernetes_cluster_attachment` add `state` to show node state.
* Resource `tencentcloud_clb_rule` support certificate modifying.
* Data Source `tencentcloud_clb_instances` add `internet_charge_type`, `internet_bandwidth_max_out`.
* Data Source `tencentcloud_clb_rules` add `http2_switch`.

BUG FIXES:

* Resource: `tencentcloud_clb_attachment` fix bug that attach more than 20 targets will failed.

## 1.53.4 (February 08, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_kubernetes_scale_worker` add `data_disk`, `docker_graph_path` to support advanced instance
  setting.
* Resource `tencentcloud_instance` add tags to the disks created with the instance.

BUG FIXES:

* Resource: `tencentcloud_kubernetes_cluster_attachment` fix bug that only one extra argument set successfully.
* Resource: `tencentcloud_as_scaling_policy` fix bug that missing required parameters error happened when update metric
  parameters.

## 1.53.3 (February 02, 2021)

ENHANCEMENTS:

* Data Source `tencentcloud_cbs_storages` add `throughput_performance` to support adding extra performance to the cbs
  resources.
* Resource `tencentcloud_kubernetes_cluster_attachment` add `hostname` to support setting hostname with the attached
  instance.

## 1.53.2 (February 01, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_cbs_storage` add `throughput_performance` to support adding extra performance to the cbs
  resources.

BUG FIXES:

* Resource: `tencentcloud_cos_bucket` fix bug that error happens when applying unsupported logging region.
* Resource: `tencentcloud_as_scaling_policy` fix bug that missing required parameters error happened when update metric
  parameters.

## 1.53.1 (January 23, 2021)

ENHANCEMENTS:

* Resource `tencentcloud_instance` add `throughput_performance` to support adding extra performance to the data disks.
* Resource `tencentcloud_kubernetes_cluster_attachment` add `file_system`, `auto_format_and_mount` and `mount_target` to
  support advanced instance setting.
* Resource `tencentcloud_kubernetes_node_pool` add `file_system`, `auto_format_and_mount` and `mount_target` to support
  advanced instance setting.
* Resource `tencentcloud_kubernetes_node_pool` add `scaling_mode` to support scaling mode setting.
* Resource `tencentcloud_kubernetes` support version upgrade.

BUG FIXES:

* Resource: `tencentcloud_gaap_http_rule` fix bug that exception happens when create more than one rule.

## 1.53.0 (January 15, 2021)

FEATURES:

* **New Resource**: `tencentcloud_ssl_pay_certificate` to support ssl pay certificate.

ENHANCEMENTS:

* Resource `tencentcloud_ccn` add `charge_type` to support billing mode setting.
* Resource `tencentcloud_ccn` add `bandwidth_limit_type` to support the speed limit type setting.
* Resource `tencentcloud_ccn_bandwidth_limit` add `dst_region` to support destination area restriction setting.
* Resource `tencentcloud_cdn_domain` add `range_origin_switch` to support range back to source configuration.
* Resource `tencentcloud_cdn_domain` add `rule_cache` to support advanced path cache configuration.
* Resource `tencentcloud_cdn_domain` add `request_header` to support request header configuration.
* Data Source `tencentcloud_ccn_instances` add `charge_type` to support billing mode.
* Data Source `tencentcloud_ccn_instances` add `bandwidth_limit_type` to support the speed limit type.
* Data Source `tencentcloud_ccn_bandwidth_limit` add `dst_region` to support destination area restriction.
* Data Source `tencentcloud_cdn_domains` add `range_origin_switch` to support range back to source configuration.
* Data Source `tencentcloud_cdn_domains` add `rule_cache` to support advanced path cache configuration.
* Data Source `tencentcloud_cdn_domains` add `request_header` to support request header configuration.

## 1.52.0 (December 28, 2020)

FEATURES:

* **New Resource**: `tencentcloud_kubernetes_node_pool` to support node management.

DEPRECATED:

* Resource: `tencentcloud_kubernetes_as_scaling_group` replaced by `tencentcloud_kubernetes_node_pool`.

## 1.51.1 (December 22, 2020)

ENHANCEMENTS:

* Resource `tencentcloud_kubernetes_cluster_attachment` add `extra_args` to support node extra arguments setting.
* Resource `tencentcloud_cos_bucket` add `log_enbale`, `log_target_bucket` and `log_prefix` to support log status
  setting.

## 1.51.0 (December 15, 2020)

FEATURES:

* **New Resource**: `tencentcloud_tcr_vpc_attachment`
* **New Data Source**: `tencentcloud_tcr_vpc_attachments`

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` support `name`, `project_id` and `description` modification.
* Doc: optimize document.

## 1.50.0 (December 08, 2020)

FEATURES:

* **New Resource**: `tencentcloud_address_template`
* **New Resource**: `tencentcloud_address_template_group`
* **New Resource**: `tencentcloud_protocol_template`
* **New Resource**: `tencentcloud_protocol_template_group`
* **New Data Source**: `tencentcloud_address_templates`
* **New Data Source**: `tencentcloud_address_template_groups`
* **New Data Source**: `tencentcloud_protocol_templates`
* **New Data Source**: `tencentcloud_protocol_template_groups`

ENHANCEMENTS:

* Resource `tencentcloud_sercurity_group_rule` add `address_template` and `protocol_template` to support building new
  security group rule with resource `tencentcloud_address_template` and `tencentcloud_protocol_template`.
* Doc: optimize document directory.

BUG FIXES:

* Resource: `tencentcloud_cos_bucket` fix bucket name validator.

## 1.49.1 (December 01, 2020)

ENHANCEMENTS:

* Doc: Update directory of document.

## 1.49.0 (November 27, 2020)

FEATURES:

* **New Resource**: `tencentcloud_tcr_instance`
* **New Resource**: `tencentcloud_tcr_token`
* **New Resource**: `tencentcloud_tcr_namespace`
* **New Resource**: `tencentcloud_tcr_repository`
* **New Data Source**: `tencentcloud_tcr_instances`
* **New Data Source**: `tencentcloud_tcr_tokens`
* **New Data Source**: `tencentcloud_tcr_namespaces`
* **New Data Source**: `tencentcloud_tcr_repositories`
* **New Resource**: `tencentcloud_cos_bucket_policy`

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_as_scaling_group` support `max_size` and `min_size` modification.

## 1.48.0 (November 20, 2020)

FEATURES:

* **New Resource**: `tencentcloud_sqlserver_basic_instance`
* **New Data Source**: `tencentcloud_sqlserver_basic_instances`

ENHANCEMENTS:

* Resource: `tencentcloud_clb_listener` support configure HTTP health check for TCP
  listener([#539](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/539)).
* Resource: `tencentcloud_clb_listener` add computed argument `target_type`.
* Data Source: `tencentcloud_clb_listeners` support getting HTTP health check config for TCP listener.

DEPRECATED:

* Resource: `tencentcloud_clb_target_group_attachment`: optional argument `targrt_group_id` is no longer supported,
  replace by `target_group_id`.

## 1.47.0 (November 13, 2020)

ENHANCEMENTS:

* Resource: `tencentcloud_clb_listener` support import.
* Resource: `tencentcloud_clb_listener` add computed argument `listener_id`.
* Resource: `tencentcloud_clb_listener_rule` support import.
* Resource: `tencentcloud_cdn_domain` add example that use COS bucket url as origin.
* Resource: `tencentcloud_sqlserver_instance` add new argument `tags`.
* Resource: `tencentcloud_sqlserver_readonly_instance` add new argument `tags`.
* Resource: `tencentcloud_elasticsearch_instance` support `node_type` and `disk_size` modification.
* Data Source: `tencentcloud_instance_types` add argument `exclude_sold_out` to support filtering sold out instance
  types.
* Data Source: `tencentcloud_sqlserver_instances` add new argument `tags`.
* Data Source: `tencentcloud_instance_types` add argument `exclude_sold_out` to support filtering sold out instance
  types.

BUG FIXES:

* Resource: `tencentcloud_elasticsearch_instance` fix inconsistent bug.
* Resource: `tencentcloud_redis_instance` fix incorrect number when updating `mem_size`.
* Data Source: `tencentcloud_redis_instances` fix incorrect number for `mem_size`.

## 1.46.4 (November 6, 2020)

BUG FIXES:

* Resource: `tencentcloud_kubernetes_cluster` fix force replacement when updating `docker_graph_path`.

## 1.46.3 (November 6, 2020)

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` add more values with argument `cluster_os` to support linux OS system.

## 1.46.2 (November 5, 2020)

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` add new argument `kube_config`.
* Resource: `tencentcloud_kubernetes_cluster` add value `tlinux2.4x86_64` with argument `cluster_os` to support linux OS
  system.
* Resource: `tencentcloud_kubernetes_cluster` add new argument `mount_target` to support set disk mount path.
* Resource: `tencentcloud_kubernetes_cluster` add new argument `docker_graph_path` to support set docker graph path.
* Resource: `tencentcloud_clb_redirection` add new argument `delete_all_auto_rewrite` to delete all auto-associated
  redirection when destroying the resource.
* Resource: `tencentcloud_kubernetes_scale_worker` add new argument `labels` to support scale worker labels.
* Data Source: `tencentcloud_kubernetes_clusters` add new argument `kube_config`.
* Data Source: `tencentcloud_availability_regions` support getting local region info by setting argument `name` with
  value `default`.
* Docs: update argument description.

BUG FIXES:

* Resource: `tencentcloud_clb_redirection` fix inconsistent bug when creating more than one auto redirection.
* Resource: `tencentcloud_redis_instance` fix updating issue when redis `type_id` is set `5`.

## 1.46.1 (October 29, 2020)

ENHANCEMENTS:

* Resource: `tencentcloud_cos_bucket` add new argument `cos_bucket_url`.
* Resource: `tencentcloud_nat_gateway` add new argument `tags`.
* Resource: `tencentcloud_postgresql_instance` add new argument `tags`.
* Data Source: `tencentcloud_cos_buckets` add new argument `cos_bucket_url`.
* Data Source: `tencentcloud_nat_gateways` add new argument `tags`.
* Data Source: `tencentcloud_postgresql_instances` add new argument `tags`.

## 1.46.0 (October 26, 2020)

FEATURES:

* **New Resource**: `tencentcloud_api_gateway_api`
* **New Resource**: `tencentcloud_api_gateway_service`
* **New Resource**: `tencentcloud_api_gateway_custom_domain`
* **New Resource**: `tencentcloud_api_gateway_usage_plan`
* **New Resource**: `tencentcloud_api_gateway_usage_plan_attachment`
* **New Resource**: `tencentcloud_api_gateway_ip_strategy`
* **New Resource**: `tencentcloud_api_gateway_strategy_attachment`
* **New Resource**: `tencentcloud_api_gateway_api_key`
* **New Resource**: `tencentcloud_api_gateway_api_key_attachment`
* **New Resource**: `tencentcloud_api_gateway_service_release`
* **New Data Source**: `tencentcloud_api_gateway_apis`
* **New Data Source**: `tencentcloud_api_gateway_services`
* **New Data Source**: `tencentcloud_api_gateway_throttling_apis`
* **New Data Source**: `tencentcloud_api_gateway_throttling_services`
* **New Data Source**: `tencentcloud_api_gateway_usage_plans`
* **New Data Source**: `tencentcloud_api_gateway_ip_strategies`
* **New Data Source**: `tencentcloud_api_gateway_customer_domains`
* **New Data Source**: `tencentcloud_api_gateway_usage_plan_environments`
* **New Data Source**: `tencentcloud_api_gateway_api_keys`

## 1.45.3 (October 21, 2020)

BUG FIXES:

* Resource: `tencentcloud_sqlserver_instance` Fix the error of releasing associated resources when destroying sqlserver
  postpaid instance.
* Resource: `tencentcloud_sqlserver_readonly_instance` Fix the bug that the instance cannot be recycled when destroying
  sqlserver postpaid instance.
* Resource: `tencentcloud_clb_instance` fix force new when updating tags.
* Resource: `tencentcloud_redis_backup_config` fix doc issues.
* Resource: `tencentcloud_instance` fix `keep_image_login` force new issue when updating terraform version.
* Resource: `tencentcloud_clb_instance` fix tag creation bug.

## 1.45.2 (October 19, 2020)

BUG FIXES:

* Resource: `tencentcloud_mysql_instance` fix creating prepaid instance error.

## 1.45.1 (October 16, 2020)

ENHANCEMENTS:

* Resource: `tencentcloud_clb_target_group_instance_attachment` update doc.
* Resource: `tencentcloud_clb_target_group_attachment` update doc.

## 1.45.0 (October 15, 2020)

FEATURES:

* **New Resource**: `tencentcloud_clb_target_group_attachment`
* **New Resource**: `tencentcloud_clb_target_group`
* **New Resource**: `tencentcloud_clb_target_group_instance_attachment`
* **New Resource**: `tencentcloud_sqlserver_publish_subscribe`
* **New Resource**: `tencentcloud_vod_adaptive_dynamic_streaming_template`
* **New Resource**: `tencentcloud_vod_procedure_template`
* **New Resource**: `tencentcloud_vod_snapshot_by_time_offset_template`
* **New Resource**: `tencentcloud_vod_image_sprite_template`
* **New Resource**: `tencentcloud_vod_super_player_config`
* **New Data Source**: `tencentcloud_clb_target_groups`
* **New Data Source**: `tencentcloud_sqlserver_publish_subscribes`
* **New Data Source**: `tencentcloud_vod_adaptive_dynamic_streaming_templates`
* **New Data Source**: `tencentcloud_vod_image_sprite_templates`
* **New Data Source**: `tencentcloud_vod_procedure_templates`
* **New Data Source**: `tencentcloud_vod_snapshot_by_time_offset_templates`
* **New Data Source**: `tencentcloud_vod_super_player_configs`

ENHANCEMENTS:

* Resource: `tencentcloud_clb_listener_rule` add new argument `target_type` to support backend target type with rule.
* Resource: `tencentcloud_mysql_instance` modify argument `engine_version` to support mysql 8.0.
* Resource: `tencentcloud_clb_listener_rule` add new argument `forward_type` to support backend
  protocol([#522](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/522)).
* Resource: `tencentcloud_instance` add new argument `keep_image_login` to support keeping image login.
* Resource: `tencentcloud_kubernetes_cluster` add new argument `extra_args` to support Kubelet.
* Resource: `tencentcloud_kubernetes_scale_worker` add new argument `extra_args` to support Kubelet.
* Resource: `tencentcloud_kubernetes_as_scaling_group` add new argument `extra_args` to support Kubelet.

## 1.44.0 (September 25, 2020)

FEATURES:

* **New Resource**: `tencentcloud_cynosdb_cluster`
* **New Resource**: `tencentcloud_cynosdb_readonly_instance`.
* **New Data Source**: `tencentcloud_cynosdb_clusters`
* **New Data Source**: `tencentcloud_cynosdb_readonly_instances`.

ENHANCEMENTS:

* Resource: `tencentcloud_mongodb_standby_instance` change example type to `POSTPAID`.
* Resource: `tencentcloud_instance` add new argument `encrypt` to support data disk with encrypt.
* Resource: `tencentcloud_elasticsearch` add new argument `encrypt` to support disk with encrypt.
* Resource: `tencentcloud_kubernetes_cluster` add new argument `cam_role_name` to support authorization with instances.

## 1.43.0 (September 18, 2020)

FEATURES:

* **New Resource**: `tencentcloud_image`
* **New Resource**: `tencentcloud_audit`
* **New Data Source**: `tencentcloud_audits`
* **New Data Source**: `tencentcloud_audit_cos_regions`
* **New Data Source**: `tencentcloud_audit_key_alias`

ENHANCEMENTS:

* Resource: `tencentcloud_instance` add new argument `data_disk_snapshot_id` to support data disk
  with `SnapshotId`([#469](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/469))
* Data Source: `tencentcloud_instances` support filter by tags.

## 1.42.2 (September 14, 2020)

BUG FIXES:

* Resource: `tencentcloud_instance` fix `key_name` update
  error([#515](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/515)).

## 1.42.1 (September 10, 2020)

BUG FIXES:

* Resource: `tencentcloud_mongodb_instance` Fix the error of releasing associated resources when destroying mongodb
  postpaid instance.
* Resource: `tencentcloud_mongodb_sharding_instance` Fix the error of releasing associated resources when destroying
  mongodb postpaid sharding instance.
* Resource: `tencentcloud_mongodb_standby_instance` Fix the error of releasing associated resources when destroying
  mongodb postpaid standby instance.

## 1.42.0 (September 8, 2020)

FEATURES:

* **New Resource**: `tencentcloud_ckafka_topic`
* **New Data Source**: `tencentcloud_ckafka_topics`

ENHANCEMENTS:

* Doc: optimize document directory.
* Resource: `tencentcloud_mongodb_instance`, `tencentcloud_mongodb_sharding_instance`
  and `tencentcloud_mongodb_standby_instance` remove system reserved tag `project`.

## 1.41.3 (September 3, 2020)

ENHANCEMENTS:

* Resource: `tencentcloud_vpc_acl_attachment` perfect example field `subnet_ids`
  to `subnet_id`([#505](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/505)).
* Resource: `tencentcloud_cbs_storage_attachment` support import.
* Resource: `tencentcloud_eip_association` support import.
* Resource: `tencentcloud_route_table_entry` support import.
* Resource: `tencentcloud_acl_attachment` support import.

## 1.41.2 (August 28, 2020)

BUG FIXES:

* Resource: `tencentcloud_vpn_connection` fix `security_group_policy` update issue when apply repeatedly.
* Resource: `tencentcloud_vpn_connection` fix inconsistent state when deleted on console.

## 1.41.1 (August 27, 2020)

BUG FIXES:

* Resource: `tencentcloud_vpn_gateway` fix force new issue when apply repeatedly.
* Resource: `tencentcloud_vpn_connection` fix force new issue when apply repeatedly.
* Resource: `tencentcloud_instance` support for adjusting `internet_max_bandwidth_out` without forceNew when
  attribute `internet_charge_type` within `TRAFFIC_POSTPAID_BY_HOUR`,`BANDWIDTH_POSTPAID_BY_HOUR`
  ,`BANDWIDTH_PACKAGE` ([#498](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/498)).

## 1.41.0 (August 17, 2020)

FEATURES:

* **New Resource**: `tencentcloud_sqlserver_instance`
* **New Resource**: `tencentcloud_sqlserver_readonly_instance`
* **New Resource**: `tencentcloud_sqlserver_db`
* **New Resource**: `tencentcloud_sqlserver_account`
* **New Resource**: `tencentcloud_sqlserver_db_account_attachment`
* **New Resource**: `tencentcloud_vpc_acl`
* **New Resource**: `tencentcloud_vpc_acl_attachment`
* **New Resource**: `tencentcloud_ckafka_acl`
* **New Resource**: `tencentcloud_ckafka_user`
* **New Data Source**: `tencentcloud_sqlserver_instance`
* **New Data Source**: `tencentcloud_sqlserver_readonly_groups`
* **New Data Source**: `tencentcloud_vpc_acls`
* **New Data Source**: `tencentcloud_ckafka_acls`
* **New Data Source**: `tencentcloud_ckafka_users`

DEPRECATED:

* Data Source: `tencentcloud_cdn_domains` optional argument `offset` is no longer supported.

ENHANCEMENTS:

* Resource: `tencentcloud_mongodb_instance`, `tencentcloud_mongodb_sharding_instance`
  and `tencentcloud_mongodb_standby_instance` remove spec update validation.

## 1.40.3 (August 11, 2020)

ENHANCEMENTS:

* Data Source: `tencentcloud_kubernetes_clusters`add new attributes `cluster_as_enabled`,`node_name_type`
  ,`cluster_extra_args`,`network_type`,`is_non_static_ip_mode`,`kube_proxy_mode`,`service_cidr`,`eni_subnet_ids`
  ,`claim_expired_seconds` and `deletion_protection`.

BUG FIXES:

* Resource: `tencentcloud_vpn_gateway` fix creation of instance when `vpc_id` is specified.
* Resource: `tencentcloud_vpn_connection` fix creation of instance when `vpc_id` is specified.
* Resource: `tencentcloud_instance` fix `internet_charge_type` inconsistency when public ip is not allocated.

## 1.40.2 (August 08, 2020)

BUG FIXES:

* Resource: `tencentcloud_instance` fix accidentally fail to delete prepaid
  instance ([#485](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/485)).

## 1.40.1 (August 05, 2020)

BUG FIXES:

* Resource: `tencentcloud_vpn_connection` fix mulit `security_group_policy` is not
  supported ([#487](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/487)).

## 1.40.0 (July 31, 2020)

FEATURES:

* **New Resource**: `tencentcloud_mongodb_standby_instance`

ENHANCEMENTS:

* Resource: `tencentcloud_gaap_http_rule` argument `realservers` now is optional.
* Resource: `tencentcloud_kubernetes_cluster` supports multiple `availability_zone`.
* Data Source: `tencentcloud_mongodb_instances` add new argument `charge_type` and `auto_renew_flag` to support prepaid
  type.
* Resource: `tencentcloud_mongodb_instance` supports prepaid type, new mongodb SDK version `2019-07-25` and standby
  instance.
* Resource: `tencentcloud_mongodb_sharding_instance` supports prepaid type, new mongodb SDK version `2019-07-25` and
  standby instance.
* Resource: `tencentcloud_security_group_lite_rule` refine update process and doc.

BUG FIXES:

* Resource: `tencentcloud_instance` fix set `key_name` error.

## 1.39.0 (July 18, 2020)

ENHANCEMENTS:

* upgrade terraform 0.13
* update readme to new repository

## 1.38.3 (July 13, 2020)

ENHANCEMENTS:

* Data Source: `tencentcloud_images` supports list of snapshots.
* Resource: `tencentcloud_kubernetes_cluster_attachment` add new argument `worker_config` to support config with
  existing instances.
* Resource: `tencentcloud_ccn` add new argument `tags` to support tags settings.
* Resource: `tencentcloud_cfs_file_system` add new argument `tags` to support tags settings.

BUG FIXES:

* Resource: `tencentcloud_gaap_layer4_listener` fix error InvalidParameter when destroy resource.
* Resource: `tencentcloud_gaap_layer7_listener` fix error InvalidParameter when destroy resource.
* Resource: `tencentcloud_cdn_domain` fix incorrect setting `server_certificate_config`, `client_certificate_config`
  caused the program to crash.

## 1.38.2 (July 03, 2020)

BUG FIXES:

* Resource: `tencentcloud_instance` fix `allocate_public_ip` inconsistency when eip is attached to the cvm.
* Resource: `tencentcloud_mysql_instance` fix auto-forcenew on `charge_type` and `pay_type` when upgrading terraform
  version. ([#459](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/459)).

## 1.38.1 (June 30, 2020)

BUG FIXES:

* Resource: `tencentcloud_cos_bucket` fix creation failure.

## 1.38.0 (June 29, 2020)

FEATURES:

* **New Data Source**: `tencentcloud_cdn_domains`

BUG FIXES:

* Resource: `tencentcloud_gaap_http_domain` fix a condition for setting client certificate
  ids([#454](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/pull/454)).

## 1.37.0 (June 23, 2020)

FEATURES:

* **New Resource**: `tencentcloud_postgresql_instance`
* **New Data Source**: `tencentcloud_postgresql_instances`
* **New Data Source**: `tencentcloud_postgresql_speccodes`
* **New Data Source**: `tencentcloud_sqlserver_zone_config`

ENHANCEMENTS:

* Resource: `tencentcloud_mongodb_instance` support more machine type.

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
* Resource: `tencentcloud_redis_instance` add new argument `charge_type`, `prepaid_period` and `force_delete` to support
  prepaid type.
* Resource: `tencentcloud_mysql_instance` add new argument `force_delete` to support soft deletion.
* Resource: `tencentcloud_mysql_readonly_instance` add new argument `force_delete` to support soft deletion.

BUG FIXES:

* Resource: `tencentcloud_instance` fix `allocate_public_ip` inconsistency when eip is attached to the cvm.

DEPRECATED:

* Data Source: `tencentcloud_mysql_instances`: optional argument `pay_type` is no longer supported, replace
  by `charge_type`.
* Resource: `tencentcloud_mysql_instance`: optional arguments `pay_type` and `period` are no longer supported, replace
  by `charge_type` and `prepaid_period`.
* Resource: `tencentcloud_mysql_readonly_instance`: optional arguments `pay_type` and `period` are no longer supported,
  replace by `charge_type` and `prepaid_period`.
* Resource: `tencentcloud_tcaplus_group` replace by `tencentcloud_tcaplus_tablegroup`
* Data Source: `tencentcloud_tcaplus_groups` replace by `tencentcloud_tcaplus_tablegroups`
* Resource: `tencentcloud_tcaplus_tablegroup`,`tencentcloud_tcaplus_idl` and `tencentcloud_tcaplus_table`
  arguments `group_id`/`group_name`  replace by `tablegroup_id`/`tablegroup_name`
* Data Source: `tencentcloud_tcaplus_groups`,`tencentcloud_tcaplus_idls` and `tencentcloud_tcaplus_tables`
  arguments `group_id`/`group_name`  replace by `tablegroup_id`/`tablegroup_name`

## 1.35.1 (June 02, 2020)

ENHANCEMENTS:

* Resource: `tencentcloud_as_scaling_config`, `tencentcloud_eip` and `tencentcloud_kubernetes_cluster` remove the
  validate function of `internet_max_bandwidth_out`.
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

* Data Source: `tencentcloud_tcaplus_applications` replace by `tencentcloud_tcaplus_clusters`,optional
  arguments `app_id` and `app_name` are no longer supported, replace by `cluster_id` and `cluster_name`
* Data Source: `tencentcloud_tcaplus_zones` replace by `tencentcloud_tcaplus_groups`,optional arguments `app_id`
  ,`zone_id` and `zone_name` are no longer supported, replace by `cluster_id`,`group_id` and `cluster_name`
* Data Source: `tencentcloud_tcaplus_tables` optional arguments `app_id` and `zone_id` are no longer supported, replace
  by `cluster_id` and `group_id`
* Data Source: `tencentcloud_tcaplus_idls`: optional argument `app_id` is no longer supported, replace by `cluster_id`.
* Resource: `tencentcloud_tcaplus_application` replace by `tencentcloud_tcaplus_cluster`,input argument `app_name` is no
  longer supported, replace by `cluster_name`
* Resource: `tencentcloud_tcaplus_zone` replace by `tencentcloud_tcaplus_group`, input arguments `app_id`
  and `zone_name` are no longer supported, replace by `cluster_id` and `group_name`
* Resource: `tencentcloud_tcaplus_idl` input arguments `app_id` and `zone_id` are no longer supported, replace
  by `cluster_id` and `group_id`
* Resource: `tencentcloud_tcaplus_table` input arguments `app_id`and `zone_id` are no longer supported, replace
  by `cluster_id` and `group_id`
* Resource: `tencentcloud_redis_instance`: optional argument `type` is no longer supported, replace by `type_id`.
* Data Source: `tencentcloud_redis_instances`: output argument `type` is no longer supported, replace by `type_id`.
* Data Source: `tencentcloud_redis_zone_config`: output argument `type` is no longer supported, replace by `type_id`.

## 1.33.1 (May 22, 2020)

ENHANCEMENTS:

* Data Source: `tencentcloud_redis_instances` add new argument `type_id`, `redis_shard_num`, `redis_replicas_num`
* Data Source: `tencentcloud_redis_zone_config` add output argument `type_id` and new output argument `type_id`
  , `redis_shard_nums`, `redis_replicas_nums`
* Data Source: `tencentcloud_ccn_instances` add new type `VPNGW` for field `instance_type`
* Data Source: `tencentcloud_vpn_gateways` add new type `CCN` for field `type`
* Resource: `tencentcloud_redis_instance` add new argument `type_id`, `redis_shard_num`, `redis_replicas_num`
* Resource: `tencentcloud_ccn_attachment` add new type `CNN_INSTANCE_TYPE_VPNGW` for field `instance_type`
* Resource: `tencentcloud_vpn_gateway` add new type `CCN` for field `type`

BUG FIXES:

* Resource: `tencentcloud_cdn_domain` fix `https_config` inconsistency after
  apply([#413](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/413)).

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
* Resource: `tencentcloud_cdn_domain` add new
  argument `full_url_cache`([#405](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/405)).

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

* **New
  Resource**: `tencentcloud_kubernetes_cluster_attachment`([#285](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/285))
  .

ENHANCEMENTS:

* Resource: `tencentcloud_cdn_domain` add new
  attribute `cname`([#395](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/395)).

BUG FIXES:

* Resource: `tencentcloud_cos_bucket_object` mark the object as destroyed when the object not exist.

## 1.31.2 (April 17, 2020)

ENHANCEMENTS:

* Resource: `tencentcloud_cbs_storage` support modify `tags`.

## 1.31.1 (April 14, 2020)

BUG FIXES:

* Resource: `tencentcloud_keypair` fix bug when trying to destroy resources containing CVM and key
  pair([#375](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/375)).
* Resource: `tencentcloud_clb_attachment` fix bug when trying to destroy multiple attachments in the array.
* Resource: `tencentcloud_cam_group_membership` fix bug when trying to destroy multiple users in the array.

ENHANCEMENTS:

* Resource: `tencentcloud_mysql_account` add new
  argument `host`([#372](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/372)).
* Resource: `tencentcloud_mysql_account_privilege` add new
  argument `account_host`([#372](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/372)).
* Resource: `tencentcloud_mysql_privilege` add new
  argument `account_host`([#372](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/372)).
* Resource: `tencentcloud_mysql_readonly_instance` check master monitor data before
  create([#379](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/379)).
* Resource: `tencentcloud_tcaplus_application` remove the pull password from server.
* Resource: `tencentcloud_instance` support
  import `allocate_public_ip`([#382](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/382)).
* Resource: `tencentcloud_redis_instance` add two redis types.
* Data Source: `tencentcloud_vpc_instances` add new argument `cidr_block`
  ,`tag_key` ([#378](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/378)).
* Data Source: `tencentcloud_vpc_route_tables` add new argument `tag_key`,`vpc_id`
  ,`association_main` ([#378](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/378)).
* Data Source: `tencentcloud_vpc_subnets` add new argument `cidr_block`,`tag_key`
  ,`is_remote_vpc_snat` ([#378](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/378)).
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
* Resource: `tencentcloud_cam_user` add new argument `force_delete`
  .([#354](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/354))
* Data Source: `tencentcloud_vpc_subnets` add new argument `vpc_id`.

## 1.30.5 (March 19, 2020)

BUG FIXES:

* Resource: `tencentcloud_key_pair` will be replaced when `public_key` contains comment.
* Resource: `tencentcloud_scf_function` upload local file error.

ENHANCEMENTS:

* Resource: `tencentcloud_scf_function` runtime support nodejs8.9 and nodejs10.15.

## 1.30.4 (March 10, 2020)

BUG FIXES:

* Resource: `tencentcloud_cam_policy` fix read nil issue when the resource is not
  exist.([#344](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/#344)).
* Resource: `tencentcloud_key_pair` will be replaced when the end of `public_key` contains
  spaces([#343](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/343)).
* Resource: `tencentcloud_scf_function` fix trigger does not support cos_region.

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` add new attributes `cluster_os_type`,`cluster_internet`,`cluster_intranet`
  ,`managed_cluster_internet_security_policies` and `cluster_intranet_subnet_id`.

## 1.30.3 (February 24, 2020)

BUG FIXES:

* Resource: `tencentcloud_instance` fix that classic network does not
  support([#339](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/339)).

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

* gaap: optimize gaap describe: when describe resource by id but get more than 1 resources, return error directly
  instead of using the first result
* Resource: `tencentcloud_eni_attachment` fix detach may failed.
* Resource: `tencentcloud_instance` remove the tag that be added by as attachment
  automatically([#300](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/300)).
* Resource: `tencentcloud_clb_listener` fix `sni_switch` type
  error([#297](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/297)).
* Resource: `tencentcloud_vpn_gateway` shows argument `prepaid_renew_flag` has changed when applied
  again([#298](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/298)).
* Resource: `tencentcloud_clb_instance` fix the bug that instance id is not set in state
  file([#303](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/303)).
* Resource: `tencentcloud_vpn_gateway` that is postpaid charge type cannot be deleted
  normally([#312](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/312)).
* Resource: `tencentcloud_vpn_gateway` add `InternalError` SDK error to triggering the retry process.
* Resource: `tencentcloud_vpn_gateway` fix read nil issue when the resource is not exist.
* Resource: `tencentcloud_clb_listener_rule` fix unclear error message of SSL type error.
* Resource: `tencentcloud_ha_vip_attachment` fix read nil issue when the resource is not exist.
* Data Source: `tencentcloud_security_group` fix `project_id` type error.
* Data Source: `tencentcloud_security_groups` fix `project_id` filter not
  works([#303](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/314)).

## 1.29.0 (January 06, 2020)

FEATURES:

* **New Data Source**: `tencentcloud_gaap_domain_error_pages`
* **New Resource**: `tencentcloud_gaap_domain_error_page`

ENHANCEMENTS:

* Data Source: `tencentcloud_vpc_instances` add new optional argument `is_default`.
* Data Source: `tencentcloud_vpc_subnets` add new optional argument `availability_zone`,`is_default`.

BUG FIXES:

* Resource: `tencentcloud_redis_instance` field security_groups are id list, not name
  list([#291](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/291)).

## 1.28.0 (December 25, 2019)

FEATURES:

* **New Data Source**: `tencentcloud_cbs_snapshot_policies`
* **New Resource**: `tencentcloud_cbs_snapshot_policy_attachment`

ENHANCEMENTS:

* doc: rewrite website index
* Resource: `tencentcloud_instance` support modifying instance
  type([#251](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/251)).
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

* Resource: `tencentcloud_mongodb_instance` support more instance
  type([#241](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/241)).
* Resource: `tencentcloud_kubernetes_cluster` support more instance
  type([#237](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/237)).

BUG FIXES:

* Fix bug that resource `tencentcloud_instance` delete error when instance launch failed.
* Fix bug that resource `tencentcloud_security_group` read error when response is InternalError.
* Fix bug that the type of `cluster_type` is wrong in data
  source `tencentcloud_mongodb_instances`([#242](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/242))
  .
* Fix bug that resource `tencentcloud_eip` unattach
  error([#233](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/233)).
* Fix bug that terraform read nil attachment resource when the attached resource of attachment resource is removed of
  resource CLB and CAM.
* Fix doc example error of resource `tencentcloud_nat_gateway`.

DEPRECATED:

* Resource: `tencentcloud_eip`: optional argument `applicable_for_clb` is no longer supported.

## 1.26.0 (December 09, 2019)

FEATURES:

* **New
  Resource**: `tencentcloud_mysql_privilege`([#223](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/223))
  .
* **New
  Resource**: `tencentcloud_kubernetes_as_scaling_group`([#202](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/202))
  .

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

* Fix bug that resource `tencentcloud_clb_listener` 's
  unchangeable `health_check_switch`([#235](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/235))
  .
* Fix bug that resource `tencentcloud_clb_instance` read nil and report error.
* Fix example errors of resource `tencentcloud_cbs_snapshot_policy` and data source `tencentcloud_dnats`.

## 1.25.2 (December 04, 2019)

BUG FIXES:

* Fixed bug that the validator of cvm instance type is incorrect.

## 1.25.1 (December 03, 2019)

ENHANCEMENTS:

* Optimized error message of validators.

BUG FIXES:

* Fixed bug that the type of `state` is incorrect in data
  source `tencentcloud_nat_gateways`([#226](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/226))
  .
* Fixed bug that the value of `cluster_max_pod_num` is incorrect in
  resource `tencentcloud_kubernetes_cluster`([#228](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/228))
  .

## 1.25.0 (December 02, 2019)

ENHANCEMENTS:

* Resource: `tencentcloud_instance` support `SPOTPAID` instance. Thanks to
  @LipingMao ([#209](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/209)).
* Resource: `tencentcloud_vpn_gateway` add argument `prepaid_renew_flag` and `prepaid_period` to support prepaid VPN
  gateway instance creation.

BUG FIXES:

* Fixed bugs that update operations on `tencentcloud_cam_policy` do not work.
* Fixed bugs that filters on `tencentcloud_cam_users` do not work.

DEPRECATED:

* Data Source: `tencentcloud_cam_user_policy_attachments`:`policy_type` is no longer supported.
* Data Source: `tencentcloud_cam_group_policy_attachments`:`policy_type` is no longer supported.

## 1.24.1 (November 26, 2019)

ENHANCEMENTS:

* Resource: `tencentcloud_kubernetes_cluster` add support for `PREPAID` instance type. Thanks to
  @woodylic ([#204](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/204)).
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

* Resource: `tencentcloud_kubernetes_cluster` cluster_os add new support: `centos7.6x86_64`
  and `ubuntu18.04.1 LTSx86_64`
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
* Data Source: `tencentcloud_as_scaling_groups` add optional argument `tags` and attribute `tags`
  for `scaling_group_list`.
* Resource: `tencentcloud_eip` add optional argument `type`, `anycast_zone`, `internet_service_provider`, etc.
* Resource: `tencentcloud_as_scaling_group` add optional argument `tags`.

BUG FIXES:

* Data Source: `tencentcloud_gaap_http_domains` set response `certificate_id`, `client_certificate_id`
  , `realserver_auth`, `basic_auth` and `gaap_auth` default value when they are nil.
* Resource: `tencentcloud_gaap_http_domain` set response `certificate_id`, `client_certificate_id`, `realserver_auth`
  , `basic_auth` and `gaap_auth` default value when they are nil.

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
* Data Source: `tencentcloud_gaap_http_rules` add optional argument `forward_host` and attributes `forward_host`
  in `rules`.
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
* Data Source: `tencentcloud_security_groups`: add optional argument `result_output_file` and new attributes `ingress`
  , `egress` for list `security_groups`.
* Resource: `tencentcloud_security_group`: add optional argument `tags`.
* Resource: `tencentcloud_as_scaling_config`: internet charge type support `BANDWIDTH_PREPAID`
  , `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.

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
* `tencentcloud_clb_attachements`: rename `rewrite_source_rule_id` with `source_rule_id` and
  rename `rewrite_target_rule_id` with `target_rule_id`.

## 1.15.1 (August 13, 2019)

ENHANCEMENTS:

* `tencentcloud_instance`: changed `image_id` property to
  ForceNew ([#78](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/78))
* `tencentcloud_instance`: improved with
  retry ([#82](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/82))
* `tencentcloud_cbs_storages`: improved with
  retry ([#82](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/82))
* `tencentcloud_clb_instance`: bug fixed and improved with
  retry ([#37](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/37))

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
* **UPDATE Resource**:`tencentcloud_mysql_instance` and `tencentcloud_mysql_readonly_instance` completely delete
  instance.

BUG FIXIES:

* resource/tencentcloud_instance: fixed issue when data disks set as delete_with_instance not works.
* resource/tencentcloud_instance: if managed public_ip manually, please don't
  define `allocate_public_ip` ([#62](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/62)).
* resource/tencentcloud_eip_association: fixed issue when instances were manually
  deleted ([#60](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/60)).
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

* data/tencentcloud_instance: fixed vpc ip is in use error when re-creating with private
  ip ([#46](https://github.com/tencentcloudstack/terraform-provider-tencentcloud/issues/46)).

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

* **New
  Resource**: `tencentcloud_lb` ([#3](https://github.com/terraform-providers/terraform-provider-scaffolding/issues/3))

ENHANCEMENTS:

* resource/tencentcloud_instance: Add `user_data_raw`
  argument ([#4](https://github.com/terraform-providers/terraform-provider-scaffolding/issues/4))

## 1.2.2 (September 28, 2018)

BUG FIXES:

* resource/tencentcloud_cbs_storage: make name to be
  required ([#25](https://github.com/tencentyun/terraform-provider-tencentcloud/issues/25))
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
