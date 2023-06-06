---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_resource_dashboard"
sidebar_current: "docs-tencentcloud-datasource-vpc_resource_dashboard"
description: |-
  Use this data source to query detailed information of vpc resource_dashboard
---

# tencentcloud_vpc_resource_dashboard

Use this data source to query detailed information of vpc resource_dashboard

## Example Usage

```hcl
data "tencentcloud_vpc_resource_dashboard" "resource_dashboard" {
  vpc_ids = ["vpc-4owdpnwr"]
}
```

## Argument Reference

The following arguments are supported:

* `vpc_ids` - (Required, Set: [`String`]) Vpc instance ID, e.g. vpc-f1xjkw1b.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `resource_dashboard_set` - List of resource objects.
  * `cdb` - Relational database.
  * `cfs` - Cloud file storage - CFS.
  * `ckafka` - Cloud Kafka (CKafka).
  * `classic_link` - Classic link.
  * `cmem` - TencentDB for Memcached.
  * `cnas` - Cnas.
  * `cts_db` - Cloud time series database.
  * `cvm` - Cloud Virtual Machine.
  * `cynos_db_mysql` - An enterprise-grade TencentDB - CynosDB for MySQL.
  * `cynos_db_postgres` - Enterprise TencentDB - CynosDB for Postgres.
  * `db_audit` - Cloud database audit.
  * `dcdb` - A distributed cloud database - TencentDB for TDSQL.
  * `dcg` - Direct Connect gateway.
  * `elastic_search` - ElasticSearch Service.
  * `emr` - EMR cluster.
  * `flow_log` - Flow log.
  * `greenplumn` - Snova data warehouse.
  * `grocery` - Grocery.
  * `hsm` - Data encryption service.
  * `ip` - Total number of used IPs except for CVM IP, EIP and network probe IP. The three IP types will be independently counted.
  * `itop` - Itop.
  * `lb` - Load balancer.
  * `maria_db` - TencentDB for MariaDB (TDSQL).
  * `mongo_db` - TencentDB for MongoDB.
  * `nas` - Network attached storage.
  * `nat` - NAT gateway.
  * `network_acl` - Network ACL.
  * `network_detect` - Network probing.
  * `oracle` - Oracle.
  * `pcx` - Peering connection.
  * `postgres` - TencentDB for PostgreSQL.
  * `redis` - TencentDB for Redis.
  * `route_table` - Route table.
  * `seal` - SEAL.
  * `sql_server` - TencentDB for SQL Server.
  * `subnet_id` - Subnet instance ID, such as subnet-bthucmmy.
  * `subnet` - Subnets.
  * `t_baas` - Blockchain service.
  * `tcaplus` - Game storage - Tcaplus.
  * `ti_db` - HTAP database - TiDB.
  * `vpc_id` - VPC instance ID, such as `vpc-bq4bzxpj`.
  * `vpngw` - VPN gateway.


