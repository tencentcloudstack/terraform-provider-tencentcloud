---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_datasource_house_attachment"
sidebar_current: "docs-tencentcloud-resource-dlc_datasource_house_attachment"
description: |-
  Provides a resource to create a DLC datasource house attachment
---

# tencentcloud_dlc_datasource_house_attachment

Provides a resource to create a DLC datasource house attachment

## Example Usage

```hcl
resource "tencentcloud_dlc_datasource_house_attachment" "example" {
  datasource_connection_name = "tf-example"
  datasource_connection_type = "Mysql"
  datasource_connection_config {
    mysql {
      location {
        vpc_id            = "vpc-khkyabcd"
        vpc_cidr_block    = "192.168.0.0/16"
        subnet_id         = "subnet-o7n9eg12"
        subnet_cidr_block = "192.168.0.0/24"
      }
    }
  }

  data_engine_names       = ["engine_demo"]
  network_connection_type = 4
  network_connection_desc = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_names` - (Required, Set: [`String`], ForceNew) Engine name, only one engine can be bound.
* `datasource_connection_config` - (Required, List, ForceNew) Data source network configuration.
* `datasource_connection_name` - (Required, String, ForceNew) Network configuration name.
* `datasource_connection_type` - (Required, String, ForceNew) Data source type. Allow value: Mysql, HiveCos, HiveHdfs, HiveCHdfs, Kafka, OtherDatasourceConnection, PostgreSql, SqlServer, ClickHouse, Elasticsearch, TDSQLPostgreSql, TCHouseD, TccHive.
* `network_connection_type` - (Required, Int, ForceNew) Network type, 2-cross-source type, 4-enhanced type.
* `network_connection_desc` - (Optional, String) Network configuration description.

The `click_house` object of `datasource_connection_config` supports the following:

* `db_name` - (Optional, String, ForceNew) Default database name.
* `instance_id` - (Optional, String, ForceNew) Unique ID of the data source instance.
* `instance_name` - (Optional, String, ForceNew) Name of the data source.
* `location` - (Optional, List, ForceNew) VPC and subnet information for the data source.

The `datasource_connection_config` object supports the following:

* `click_house` - (Optional, List, ForceNew) Properties of ClickHouse data source connection.
* `elasticsearch` - (Optional, List, ForceNew) Properties of Elasticsearch data source connection.
* `hive` - (Optional, List, ForceNew) Properties of Hive data source connection.
* `kafka` - (Optional, List, ForceNew) Properties of Kafka data source connection.
* `mysql` - (Optional, List, ForceNew) Properties of MySQL data source connection.
* `other_datasource_connection` - (Optional, List, ForceNew) Properties of other data source connection.
* `postgre_sql` - (Optional, List, ForceNew) Properties of PostgreSQL data source connection.
* `sql_server` - (Optional, List, ForceNew) Properties of SQLServer data source connection.
* `tc_house_d` - (Optional, List, ForceNew) Properties of Doris data source connection.
* `tcc_hive` - (Optional, List, ForceNew) TccHive data catalog connection information.
* `tdsql_postgre_sql` - (Optional, List, ForceNew) Properties of TDSQL-PostgreSQL data source connection.

The `elasticsearch` object of `datasource_connection_config` supports the following:

* `db_name` - (Optional, String, ForceNew) Default database name.
* `instance_id` - (Optional, String, ForceNew) Data source ID.
* `instance_name` - (Optional, String, ForceNew) Data source name.
* `location` - (Optional, List, ForceNew) VPC and subnet information for the data source.
* `service_info` - (Optional, List, ForceNew) IP and port information for accessing Elasticsearch.

The `hive` object of `datasource_connection_config` supports the following:

* `location` - (Required, List, ForceNew) Private network information where the data source is located.
* `meta_store_url` - (Required, String, ForceNew) Address of Hive metastore.
* `type` - (Required, String, ForceNew) Hive data source type, representing data storage location, COS or HDFS.
* `bucket_url` - (Optional, String, ForceNew) If the type is COS, COS bucket connection needs to be filled in.
* `hdfs_properties` - (Optional, String, ForceNew) JSON string. If the type is HDFS, this field needs to be filled in.
* `high_availability` - (Optional, Bool, ForceNew) If the type is HDFS, high availability needs to be selected.
* `hive_version` - (Optional, String, ForceNew) Version number of Hive component in EMR cluster.
* `instance_id` - (Optional, String, ForceNew) EMR cluster ID.
* `instance_name` - (Optional, String, ForceNew) EMR cluster name.
* `kerberos_enable` - (Optional, Bool, ForceNew) Whether to enable Kerberos.
* `kerberos_info` - (Optional, List, ForceNew) Kerberos details.
* `mysql` - (Optional, List, ForceNew) Metadata database information for Hive.

The `kafka` object of `datasource_connection_config` supports the following:

* `instance_id` - (Required, String, ForceNew) Kafka instance ID.
* `location` - (Required, List, ForceNew) Network information for Kafka data source.

The `kerberos_info` object of `hive` supports the following:

* `key_tab` - (Optional, String, ForceNew) KeyTab file value.
* `krb5_conf` - (Optional, String, ForceNew) Krb5Conf file value.
* `service_principal` - (Optional, String, ForceNew) Service principal.

The `location` object of `click_house` supports the following:

* `subnet_cidr_block` - (Required, String, ForceNew) Subnet IPv4 CIDR.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.
* `vpc_cidr_block` - (Required, String, ForceNew) VPC IPv4 CIDR.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.

The `location` object of `elasticsearch` supports the following:

* `subnet_cidr_block` - (Required, String, ForceNew) Subnet IPv4 CIDR.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.
* `vpc_cidr_block` - (Required, String, ForceNew) VPC IPv4 CIDR.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.

The `location` object of `hive` supports the following:

* `subnet_cidr_block` - (Required, String, ForceNew) Subnet IPv4 CIDR.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.
* `vpc_cidr_block` - (Required, String, ForceNew) VPC IPv4 CIDR.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.

The `location` object of `kafka` supports the following:

* `subnet_cidr_block` - (Required, String, ForceNew) Subnet IPv4 CIDR.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.
* `vpc_cidr_block` - (Required, String, ForceNew) VPC IPv4 CIDR.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.

The `location` object of `mysql` supports the following:

* `subnet_cidr_block` - (Required, String, ForceNew) Subnet IPv4 CIDR.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.
* `vpc_cidr_block` - (Required, String, ForceNew) VPC IPv4 CIDR.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.

The `location` object of `other_datasource_connection` supports the following:

* `subnet_cidr_block` - (Required, String, ForceNew) Subnet IPv4 CIDR.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.
* `vpc_cidr_block` - (Required, String, ForceNew) VPC IPv4 CIDR.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.

The `location` object of `postgre_sql` supports the following:

* `subnet_cidr_block` - (Required, String, ForceNew) Subnet IPv4 CIDR.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.
* `vpc_cidr_block` - (Required, String, ForceNew) VPC IPv4 CIDR.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.

The `location` object of `sql_server` supports the following:

* `subnet_cidr_block` - (Required, String, ForceNew) Subnet IPv4 CIDR.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.
* `vpc_cidr_block` - (Required, String, ForceNew) VPC IPv4 CIDR.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.

The `location` object of `tc_house_d` supports the following:

* `subnet_cidr_block` - (Required, String, ForceNew) Subnet IPv4 CIDR.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.
* `vpc_cidr_block` - (Required, String, ForceNew) VPC IPv4 CIDR.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.

The `location` object of `tdsql_postgre_sql` supports the following:

* `subnet_cidr_block` - (Required, String, ForceNew) Subnet IPv4 CIDR.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID where the data connection is located, such as 'subnet-bthucmmy'.
* `vpc_cidr_block` - (Required, String, ForceNew) VPC IPv4 CIDR.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID where the data connection is located, such as 'vpc-azd4dt1c'.

The `mysql` object of `datasource_connection_config` supports the following:

* `location` - (Required, List, ForceNew) Network information for MySQL data source.
* `db_name` - (Optional, String, ForceNew) Database name.
* `instance_id` - (Optional, String, ForceNew) Database instance ID, consistent with the database side.
* `instance_name` - (Optional, String, ForceNew) Database instance name, consistent with the database side.

The `mysql` object of `hive` supports the following:

* `location` - (Required, List, ForceNew) Network information for MySQL data source.
* `db_name` - (Optional, String, ForceNew) Database name.
* `instance_id` - (Optional, String, ForceNew) Database instance ID, consistent with the database side.
* `instance_name` - (Optional, String, ForceNew) Database instance name, consistent with the database side.

The `other_datasource_connection` object of `datasource_connection_config` supports the following:

* `location` - (Required, List, ForceNew) Network parameters.

The `postgre_sql` object of `datasource_connection_config` supports the following:

* `db_name` - (Optional, String, ForceNew) Default database name.
* `instance_id` - (Optional, String, ForceNew) Unique ID of the data source instance.
* `instance_name` - (Optional, String, ForceNew) Name of the data source.
* `location` - (Optional, List, ForceNew) VPC and subnet information for the data source.

The `service_info` object of `elasticsearch` supports the following:

* `ip` - (Optional, String, ForceNew) IP information.
* `port` - (Optional, Int, ForceNew) Port information.

The `sql_server` object of `datasource_connection_config` supports the following:

* `db_name` - (Optional, String, ForceNew) Default database name.
* `instance_id` - (Optional, String, ForceNew) Unique ID of the data source instance.
* `instance_name` - (Optional, String, ForceNew) Name of the data source.
* `location` - (Optional, List, ForceNew) VPC and subnet information for the data source.

The `tc_house_d` object of `datasource_connection_config` supports the following:

* `access_info` - (Optional, String, ForceNew) Access information.
* `db_name` - (Optional, String, ForceNew) Default database name.
* `instance_id` - (Optional, String, ForceNew) Unique ID of the data source instance.
* `instance_name` - (Optional, String, ForceNew) Data source name.
* `location` - (Optional, List, ForceNew) VPC and subnet information for the data source.

The `tcc_connection` object of `tcc_hive` supports the following:

* `clb_ip` - (Optional, String, ForceNew) Service CLB IP.
* `clb_port` - (Optional, String, ForceNew) Service CLB port.
* `subnet_cidr_block` - (Optional, String, ForceNew) Subnet CIDR.
* `subnet_id` - (Optional, String, ForceNew) Subnet instance ID.
* `vpc_cidr_block` - (Optional, String, ForceNew) VPC CIDR.
* `vpc_id` - (Optional, String, ForceNew) VPC instance ID.

The `tcc_hive` object of `datasource_connection_config` supports the following:

* `endpoint_service_id` - (Optional, String, ForceNew) Endpoint service ID.
* `hive_version` - (Optional, String, ForceNew) Hive version.
* `hms_endpoint_service_id` - (Optional, String, ForceNew) HMS endpoint service ID.
* `instance_id` - (Optional, String, ForceNew) Instance ID.
* `instance_name` - (Optional, String, ForceNew) Instance name.
* `meta_store_url` - (Optional, String, ForceNew) Thrift connection address.
* `tcc_connection` - (Optional, List, ForceNew) Network information.

The `tdsql_postgre_sql` object of `datasource_connection_config` supports the following:

* `db_name` - (Optional, String, ForceNew) Default database name.
* `instance_id` - (Optional, String, ForceNew) Unique ID of the data source instance.
* `instance_name` - (Optional, String, ForceNew) Name of the data source.
* `location` - (Optional, List, ForceNew) VPC and subnet information for the data source.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



