---
subcategory: "SSL Certificates(ssl)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssl_certificate_bind_resource_task_detail"
sidebar_current: "docs-tencentcloud-datasource-ssl_certificate_bind_resource_task_detail"
description: |-
  Use this data source to query detailed information of SSL certificate bind resource task detail, including the bind details of the certificate across various cloud resources (CLB, CDN, WAF, DDoS, Live, VOD, TKE, APIGateway, TCB, TEO, TSE, COS, TDMQ, MQTT, GAAP, SCF).
---

# tencentcloud_ssl_certificate_bind_resource_task_detail

Use this data source to query detailed information of SSL certificate bind resource task detail, including the bind details of the certificate across various cloud resources (CLB, CDN, WAF, DDoS, Live, VOD, TKE, APIGateway, TCB, TEO, TSE, COS, TDMQ, MQTT, GAAP, SCF).

## Example Usage

### Query all resource types bind details by task id

```hcl
data "tencentcloud_ssl_certificate_bind_resource_task_detail" "example" {
  task_id = "task-12345"
}
```

### Query specific resource types bind details

```hcl
data "tencentcloud_ssl_certificate_bind_resource_task_detail" "example" {
  task_id        = "task-12345"
  resource_types = ["clb", "cdn"]
}
```

### Query bind details with regions filter

```hcl
data "tencentcloud_ssl_certificate_bind_resource_task_detail" "example" {
  task_id        = "task-12345"
  resource_types = ["clb"]
  regions        = ["ap-guangzhou", "ap-beijing"]
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String) Task ID, query the bind cloud resource result according to the task ID obtained by CreateCertificateBindResourceSyncTask.
* `regions` - (Optional, Set: [`String`]) List of regions to query. clb, tke, waf, apigateway, tcb, cos, tse support region query.
* `resource_types` - (Optional, Set: [`String`]) Query the result details of the resource type. If not passed, query all. Supported values: clb, cdn, ddos, live, vod, waf, apigateway, teo, tke, cos, tse, tcb, tdmq, mqtt, gaap, scf.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `apigateway` - Associated APIGateway resource detail list.
  * `error` - Query error message.
  * `instance_list` - APIGateway instance detail list.
    * `cert_id` - Certificate ID.
    * `domain` - Domain name.
    * `protocol` - Protocol.
    * `service_id` - Service ID.
    * `service_name` - Service name.
  * `region` - Region.
  * `total_count` - Total number of APIGateway instances in the region.
* `cache_time` - Current result cache time.
* `cdn` - Associated CDN resource detail list.
  * `error` - Query error message.
  * `instance_list` - CDN domain detail list.
    * `cert_id` - Deployed certificate ID.
    * `domain` - Domain name.
    * `https_billing_switch` - Domain billing status, on means enabled, off means disabled.
    * `status` - Domain status. rejected: domain review not passed, processing: deploying, online: started, offline: closed.
  * `total_count` - Total number of CDN domains in the region.
* `clb` - Associated CLB resource detail list. Note: This field may return null, indicating that no valid value can be obtained.
  * `error` - Query error message.
  * `instance_list` - CLB instance detail list.
    * `forward` - Load balancer type, 0 traditional load balancer; 1 application load balancer.
    * `listeners` - CLB listener list.
      * `certificate` - Certificate data bound to the listener.
        * `cert_ca_id` - Root certificate ID.
        * `cert_id` - Certificate ID.
        * `dns_names` - Domain names bound to the certificate.
        * `s_s_l_mode` - Certificate authentication mode: UNIDIRECTIONAL one-way authentication, MUTUAL two-way authentication.
      * `listener_id` - Listener ID.
      * `listener_name` - Listener name.
      * `no_match_domains` - List of non-matching domains.
      * `protocol` - Listener protocol type, HTTPS|TCP_SSL.
      * `rules` - Listener rule list.
        * `certificate` - Certificate data bound to the rule.
          * `cert_ca_id` - Root certificate ID.
          * `cert_id` - Certificate ID.
          * `dns_names` - Domain names bound to the certificate.
          * `s_s_l_mode` - Certificate authentication mode: UNIDIRECTIONAL one-way authentication, MUTUAL two-way authentication.
        * `domain` - Domain name bound to the rule.
        * `is_match` - Whether the rule matches the domain name of the certificate to be bound.
        * `location_id` - Rule ID.
        * `no_match_domains` - List of non-matching domains.
        * `url` - Path bound to the rule.
      * `sni_switch` - Whether to enable SNI, 1 for enable, 0 for disable.
    * `load_balancer_id` - CLB instance ID.
    * `load_balancer_name` - CLB instance name.
  * `region` - Region.
  * `total_count` - Total number of CLB instances in the region.
* `cos` - Associated COS resource detail list.
  * `error` - Query error message.
  * `instance_list` - COS instance detail list.
    * `bucket` - Bucket name.
    * `cert_id` - Bound certificate ID.
    * `domain` - Domain name.
    * `region` - Bucket region.
    * `status` - ENABLED: domain online status; DISABLED: domain offline status.
  * `region` - Region.
  * `total_count` - Total number of COS instances in the region.
* `ddos` - Associated DDoS resource detail list.
  * `error` - Query error message.
  * `instance_list` - DDoS instance detail list.
    * `cert_id` - Certificate ID.
    * `domain` - Domain name.
    * `instance_id` - Instance ID.
    * `protocol` - Protocol type.
    * `virtual_port` - Forwarding port.
  * `total_count` - Total number of DDoS domains in the region.
* `gaap` - Associated GAAP resource detail list.
  * `error` - Query error message.
  * `instance_list` - GAAP instance detail list.
    * `instance_id` - Instance ID.
    * `instance_name` - Instance name.
    * `listener_list` - Listener list.
      * `cert_id_list` - Certificate ID list bound to the instance.
      * `listener_id` - Listener ID.
      * `listener_name` - Listener name.
      * `listener_status` - Listener status.
      * `no_match_domains` - List of non-matching domains.
      * `protocol` - Listener protocol.
  * `total_count` - Total number of GAAP instances in the region.
* `live` - Associated Live resource detail list.
  * `error` - Query error message.
  * `instance_list` - Live instance detail list.
    * `cert_id` - Bound certificate ID.
    * `domain` - Domain name.
    * `status` - -1: domain not associated with certificate. 1: domain https enabled. 0: domain https disabled.
  * `total_count` - Total number of Live instances in the region.
* `mqtt` - Associated MQTT resource detail list.
  * `error` - Query error message.
  * `instance_list` - MQTT instance detail list.
    * `ca_cert_id_list` - CA certificate list.
    * `instance_id` - Instance ID.
    * `instance_name` - Instance name.
    * `instance_status` - Instance status.
    * `no_match_domains` - List of non-matching domains.
    * `server_cert_id_list` - Server certificate list.
  * `region` - Region.
  * `total_count` - Total number of MQTT instances in the region.
* `scf` - Associated SCF resource detail list.
  * `error` - Query error message.
  * `instance_list` - SCF instance detail list.
    * `certificate_id` - Certificate ID.
    * `domain` - Domain name.
    * `protocol` - Protocol.
    * `region` - Region.
  * `region` - Region.
  * `total_count` - Total number of SCF instances in the region.
* `status` - Async query result status of associated cloud resources: 0 means querying in progress, 1 means query success, 2 means query exception.
* `tcb` - Associated TCB resource detail list.
  * `environments` - TCB environment instance detail list.
    * `access_service` - Access service.
      * `instance_list` - Access service instance list.
        * `domain` - Domain name.
        * `i_c_p_status` - ICP blacklist status, 0-not banned, 1-banned.
        * `is_preempted` - Whether preempted.
        * `old_certificate_id` - Bound certificate ID.
        * `status` - Status.
        * `union_status` - Unified domain status.
      * `total_count` - Total count.
    * `environment` - TCB environment.
      * `id` - Unique ID.
      * `name` - Name.
      * `source` - Source.
      * `status` - Status.
    * `host_service` - Host service.
      * `instance_list` - Host service instance list.
        * `d_n_s_status` - DNS status.
        * `domain` - Domain name.
        * `old_certificate_id` - Bound certificate ID.
        * `status` - Status.
      * `total_count` - Total count.
  * `error` - Query error message.
  * `region` - Region.
* `tdmq` - Associated TDMQ resource detail list.
  * `error` - Query error message.
  * `instance_list` - TDMQ instance detail list.
    * `ca_cert_id` - CA certificate ID.
    * `cert_id` - Server certificate ID.
    * `instance_id` - Instance ID.
    * `instance_name` - Instance name.
    * `instance_status` - Instance status.
    * `no_match_domains` - List of non-matching domains.
  * `region` - Region.
  * `total_count` - Total number of TDMQ instances in the region.
* `teo` - Associated TEO resource detail list.
  * `error` - Query error message.
  * `instance_list` - TEO instance detail list.
    * `algorithm` - Certificate encryption algorithm.
    * `cert_id` - Certificate ID.
    * `host` - Domain name.
    * `status` - Domain status. deployed: deployed; processing: deploying; applying: applying; failed: apply failed; issued: bind failed.
    * `zone_id` - Zone ID.
  * `total_count` - Total number of TEO instances in the region.
* `tke` - Associated TKE resource detail list.
  * `error` - Query error message.
  * `instance_list` - TKE instance detail list.
    * `cluster_id` - Cluster ID.
    * `cluster_name` - Cluster name.
    * `cluster_type` - Cluster type.
    * `cluster_version` - Cluster version.
    * `namespace_list` - Cluster namespace list.
      * `name` - Namespace name.
      * `secret_list` - Secret list.
        * `cert_id` - Certificate ID.
        * `ingress_list` - Ingress list.
          * `domains` - Ingress domain list.
          * `ingress_name` - Ingress name.
          * `tls_domains` - TLS domain list.
        * `name` - Secret name.
        * `no_match_domains` - List of domains that do not match the new certificate.
  * `region` - Region.
  * `total_count` - Total number of TKE instances in the region.
* `tse` - Associated TSE resource detail list.
  * `error` - Query error message.
  * `instance_list` - TSE instance detail list.
    * `certificate_list` - Gateway certificate list.
      * `bind_domains` - Bound domains.
      * `cert_id` - Currently bound SSL certificate ID.
      * `cert_source` - Certificate source.
      * `id` - Gateway certificate ID.
      * `name` - Gateway certificate name.
    * `gateway_id` - Gateway ID.
    * `gateway_name` - Gateway name.
  * `region` - Region.
  * `total_count` - Total number of TSE instances in the region.
* `vod` - Associated VOD resource detail list.
  * `error` - Query error message.
  * `instance_list` - VOD instance detail list.
    * `cert_id` - Certificate ID.
    * `domain` - Domain name.
  * `total_count` - Total number of VOD instances in the region.
* `waf` - Associated WAF resource detail list.
  * `error` - Query error message.
  * `instance_list` - WAF instance detail list.
    * `cert_id` - Certificate ID.
    * `domain` - Domain name.
    * `keepalive` - Whether to keep long connection, 1 yes, 0 no.
  * `region` - Region.
  * `total_count` - Total number of WAF instances in the region.


