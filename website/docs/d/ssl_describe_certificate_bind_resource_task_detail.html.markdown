---
subcategory: "SSL Certificates"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_describe_certificate_bind_resource_task_detail"
sidebar_current: "docs-tencentcloud-datasource-ssl_describe_certificate_bind_resource_task_detail"
description: |-
  Use this data source to query detailed information of ssl describe_certificate_bind_resource_task_detail
---

# tencentcloud_ssl_describe_certificate_bind_resource_task_detail

Use this data source to query detailed information of ssl describe_certificate_bind_resource_task_detail

## Example Usage

```hcl
data "tencentcloud_ssl_describe_certificate_bind_resource_task_detail" "describe_certificate_bind_resource_task_detail" {
  task_id        = ""
  resource_types =
  regions        =
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String) Task ID, query the results of cloud resources based on the task ID query.
* `regions` - (Optional, Set: [`String`]) Query the data of the regional list, CLB, TKE, WAF, Apigateway, TCB support regional inquiries, other resource types do not support.
* `resource_types` - (Optional, Set: [`String`]) Query the results of the type of resource type.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `apigateway` - Related ApIgateway Resources DetailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `instance_list` - Apiguateway example detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `cert_id` - Certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `domain` - domain name.
    * `protocol` - Use Agreement.
    * `service_id` - Instance ID.
    * `service_name` - Example name.
  * `region` - area.
  * `total_count` - The total number of Apigateway instances under this region.
* `cache_time` - Current result cache time.
* `cdn` - Related CDN resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `instance_list` - CDN domain name detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `cert_id` - Deployment certificate ID.
    * `domain` - domain name.
    * `https_billing_switch` - Domain name billing status.
    * `status` - Domain name.
  * `total_count` - The total number of domain names in the region.
* `clb` - Related CLB resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `instance_list` - CLB instance detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `listeners` - CLB listener listNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `certificate` - Certificate data binding of listenersNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `cert_ca_id` - Root certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `cert_id` - Certificate ID.
        * `dns_names` - Domain name binding of certificates.
        * `s_s_l_mode` - Certificate certification mode: unidirectional unidirectional authentication, Mutual two -way certificationNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `listener_id` - Listener ID.
      * `listener_name` - Name of listeners.
      * `no_match_domains` - List of non -matching fieldsNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `protocol` - Type of listener protocol, https | TCP_SSL.
      * `rules` - List of listeners&#39; rulesNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `certificate` - Certificate data that has been bound to the rulesNote: This field may return NULL, indicating that the valid value cannot be obtained.
          * `cert_ca_id` - Root certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
          * `cert_id` - Certificate ID.
          * `dns_names` - Domain name binding of certificates.
          * `s_s_l_mode` - Certificate certification mode: unidirectional unidirectional authentication, Mutual two -way certificationNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `domain` - Domain name binding.
        * `is_match` - Whether the rules match the domain name to be bound to the certificate.
        * `location_id` - Rule ID.
        * `no_match_domains` - List of non -matching fieldsNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `sni_switch` - Whether to turn on SNI, 1 to open, 0 to close.
    * `load_balancer_id` - CLB instance ID.
    * `load_balancer_name` - CLB instance name name.
  * `region` - area.
  * `total_count` - The total number of CLB instances under this region.
* `ddos` - Related DDOS resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `instance_list` - DDOS example detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `cert_id` - Certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `domain` - domain name.
    * `instance_id` - Instance ID.
    * `protocol` - agreement type.
    * `virtual_port` - Forwarding port.
  * `total_count` - The total number of DDOS domain names under this region.
* `live` - Related live resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `instance_list` - Live instance detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `cert_id` - Binded certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `domain` - domain name.
    * `status` - -1: Unrelated certificate of domain name.1: The domain name HTTPS has been opened.0: The domain name HTTPS has been closed.
  * `total_count` - The total number of LIVE instances in this region.
* `status` - Related Cloud Resources Inquiry results: 0 indicates that in the query, 1 means the query is successful.2 means the query is abnormal; if the status is 1, check the results of bindResourceResult; if the state is 2, check the reason for ERROR.
* `tcb` - Related TCB resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `environments` - TCB environment example detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `access_service` - Access serviceNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `instance_list` - Instance listNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `domain` - domain nameNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `i_c_p_status` - ICP blacklist is banned, 0-unblocks, 1-banNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `is_preempted` - Whether it is seized, being seized means that the domain name is bound by other environments, and needs to be unbinded or re -bound.Note: This field may return NULL, indicating that the valid value cannot be obtained.
        * `old_certificate_id` - Binding certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `status` - stateNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `union_status` - Unified domain name stateNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `total_count` - quantityNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `environment` - TCB environmentNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `i_d` - Unique IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `name` - nameNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `source` - sourceNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `status` - stateNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `host_service` - Static custodyNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `instance_list` - Instance listNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `d_n_s_status` - Analytical statusNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `domain` - domain nameNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `old_certificate_id` - Binding certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
        * `status` - stateNote: This field may return NULL, indicating that the valid value cannot be obtained.
      * `total_count` - quantityNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `region` - area.
* `teo` - Related Teo resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `instance_list` - Edgeone example detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `cert_id` - Certificate ID.
    * `host` - domain name.
    * `status` - Domain name.
    * `zone_id` - Regional IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `total_count` - Edgeone total number.
* `tke` - Related TKE resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `instance_list` - TKE instance detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `cluster_id` - Cluster ID.
    * `cluster_name` - Cluster name.
    * `cluster_type` - Cluster.
    * `cluster_version` - Cluster.
    * `namespace_list` - Cluster Naming Space List.
      * `name` - namespace name.
      * `secret_list` - Secret list.
        * `cert_id` - Certificate ID.
        * `ingress_list` - Ingress list.
          * `domains` - Ingress domain name list.
          * `ingress_name` - Ingress name.
          * `tls_domains` - TLS domain name list.
        * `name` - Secret name.
        * `no_match_domains` - List of domain names that are not matched with the new certificateNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `region` - area.
  * `total_count` - Total TKE instance under this region.
* `vod` - Related VOD resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `instance_list` - VOD example detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `cert_id` - Certificate ID.
    * `domain` - domain name.
  * `total_count` - The total number of VOD instances under this region.
* `waf` - Related WAF resource detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `instance_list` - WAF instance detailsNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `cert_id` - Certificate IDNote: This field may return NULL, indicating that the valid value cannot be obtained.
    * `domain` - domain name.
    * `keepalive` - Whether to keep a long connection, 1 is, 0 NoNote: This field may return NULL, indicating that the valid value cannot be obtained.
  * `region` - area.
  * `total_count` - The total number of WAF instances under this region.


