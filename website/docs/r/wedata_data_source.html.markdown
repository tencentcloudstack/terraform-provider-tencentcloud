---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_data_source"
sidebar_current: "docs-tencentcloud-resource-wedata_data_source"
description: |-
  Provides a resource to create a WeData data source
---

# tencentcloud_wedata_data_source

Provides a resource to create a WeData data source

## Example Usage

```hcl
resource "tencentcloud_wedata_data_source" "example" {
  project_id = "2983848457986924544"
  name       = "tf_example"
  type       = "MYSQL"
  prod_con_properties = jsonencode({
    "deployType" : "CONNSTR_PUBLICDB",
    "url" : "jdbc:mysql://1.1.1.1:1111/database",
    "username" : "root",
    "password" : "root"
  })

  display_name = "display_name"
  description  = "description"

  lifecycle {
    ignore_changes = [
      prod_con_properties,
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Data source name.
* `prod_con_properties` - (Required, String) Data source configuration information, stored in JSON KV format, with different KV storage information for each data source type.

> deployType: 
CONNSTR_PUBLICDB(Public network instance) 
CONNSTR_CVMDB(Self-built instance)
INSTANCE(Cloud instance)

```
mysql: Self-built instance
{
    "deployType": "CONNSTR_CVMDB",
    "url": "jdbc:mysql://1.1.1.1:1111/database",
    "username": "root",
    "password": "root",
    "region": "ap-shanghai",
    "vpcId": "vpc-kprq42yo",
    "type": "MYSQL"
}
mysql: Cloud instance
{
    "instanceid": "cdb-12uxdo5e",
    "db": "db",
    "region": "ap-shanghai",
    "username": "msyql",
    "password": "mysql",
    "deployType": "INSTANCE",
    "type": "TENCENT_MYSQL"
}
sql_server: 
{
    "deployType": "CONNSTR_PUBLICDB",
    "url": "jdbc:sqlserver://1.1.1.1:223;DatabaseName=database",
    "username": "user_1",
    "password": "pass_2",
    "type": "SQLSERVER"
}
redis:
    redisType:
    -NO_ACCOUT(No account)
    -SELF_ACCOUNT(Custom account)
{
    "deployType": "CONNSTR_PUBLICDB",
    "username":""
    "password": "pass",
    "ip": "1.1.1.1",
    "port": "6379",
    "redisType": "NO_ACCOUT",
    "type": "REDIS"
}
oracle: 
{
    "deployType": "CONNSTR_CVMDB",
    "url": "jdbc:oracle:thin:@1.1.1.1:1521:prod",
    "username": "oracle",
    "password": "pass",
    "region": "ap-shanghai",
    "vpcId": "vpc-kprq42yo",
    "type": "ORACLE"
}
mongodb:
    advanceParams(Custom parameters, will be appended to the URL)
{
    "advanceParams": [
        {
            "key": "authSource",
            "value": "auth"
        }
    ],
    "db": "admin",
    "deployType": "CONNSTR_PUBLICDB",
    "username": "user",
    "password": "pass",
    "type": "MONGODB",
    "host": "1.1.1.1:9200"
}
postgresql:
{
    "deployType": "CONNSTR_PUBLICDB",
    "url": "jdbc:postgresql://1.1.1.1:1921/database",
    "username": "user",
    "password": "pass",
    "type": "POSTGRE"
}
kafka:
    authType:
        - sasl
        - jaas
        - sasl_plaintext
        - sasl_ssl
        - GSSAPI
    ssl:
        -PLAIN
        -GSSAPI
{
    "deployType": "CONNSTR_PUBLICDB",
    "host": "1.1.1.1:9092",
    "ssl": "GSSAPI",
    "authType": "sasl",
    "type": "KAFKA",
    "principal": "aaaa",
    "serviceName": "kafka"
}

cos:
{
    "region": "ap-shanghai",
    "deployType": "INSTANCE",
    "secretId": "aaaaa",
    "secretKey": "sssssss",
    "bucket": "aaa",
    "type": "COS"
}

```.
* `project_id` - (Required, String, ForceNew) Data source project ID.
* `type` - (Required, String, ForceNew) Data source type: enumeration values.

- MYSQL
- TENCENT_MYSQL
- POSTGRE
- ORACLE
- SQLSERVER
- FTP
- HIVE
- HUDI
- HDFS
- ICEBERG
- KAFKA
- DTS_KAFKA
- HBASE
- SPARK
- TBASE
- DB2
- DM
- GAUSSDB
- GBASE
- IMPALA
- ES
- TENCENT_ES
- GREENPLUM
- SAP_HANA
- SFTP
- OCEANBASE
- CLICKHOUSE
- KUDU
- VERTICA
- REDIS
- COS
- DLC
- DORIS
- CKAFKA
- S3_DATAINSIGHT
- TDSQL
- TDSQL_MYSQL
- MONGODB
- TENCENT_MONGODB
- REST_API
- TiDB
- StarRocks
- Trino
- Kyuubi
- TCHOUSE_X
- TCHOUSE_P
- TCHOUSE_C
- TCHOUSE_D
- INFLUXDB
- BIG_QUERY
- SSH
- BLOB
- TDSQL_POSTGRE
- GDB
- TDENGINE
- TDSQLC.
* `description` - (Optional, String) Data source description information.
* `dev_con_properties` - (Optional, String) Development environment data source configuration information, required if the project is in standard mode.
* `dev_file_upload` - (Optional, List) Development environment data source file upload.
* `display_name` - (Optional, String) Data source display name, for visual viewing.
* `prod_file_upload` - (Optional, List) Production environment data source file upload.

The `dev_file_upload` object supports the following:

* `core_site` - (Optional, String) core-site.xml file.
* `hbase_site` - (Optional, String) hbase-site file.
* `hdfs_site` - (Optional, String) hdfs-site.xml file.
* `hive_site` - (Optional, String) hive-site.xml file.
* `key_store` - (Optional, String) Keystore authentication file, default filename keystore.jks.
* `key_tab` - (Optional, String) keytab file, default filename [data source name].keytab.
* `krb5_conf` - (Optional, String) krb5.conf file.
* `private_key` - (Optional, String) Private key, default filename private_key.pem.
* `public_key` - (Optional, String) Public key, default filename public_key.pem.
* `trust_store` - (Optional, String) Truststore authentication file, default filename truststore.jks.

The `prod_file_upload` object supports the following:

* `core_site` - (Optional, String) core-site.xml file.
* `hbase_site` - (Optional, String) hbase-site file.
* `hdfs_site` - (Optional, String) hdfs-site.xml file.
* `hive_site` - (Optional, String) hive-site.xml file.
* `key_store` - (Optional, String) Keystore authentication file, default filename keystore.jks.
* `key_tab` - (Optional, String) keytab file, default filename [data source name].keytab.
* `krb5_conf` - (Optional, String) krb5.conf file.
* `private_key` - (Optional, String) Private key, default filename private_key.pem.
* `public_key` - (Optional, String) Public key, default filename public_key.pem.
* `trust_store` - (Optional, String) Truststore authentication file, default filename truststore.jks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `data_source_id` - Data source ID.


