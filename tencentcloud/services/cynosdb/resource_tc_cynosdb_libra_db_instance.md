Provides a resource to create a CynosDB (TDSQL-C) LibraDB read-only analytics engine instance

Example Usage

```hcl
resource "tencentcloud_cynosdb_libra_db_instance" "example" {
  cluster_id         = "cynosdbmysql-5oo78wv9"
  zone               = "ap-guangzhou-7"
  cpu                = 8
  mem                = 32
  storage_size       = 100
  pay_mode           = 0
  port               = 2000
  instance_name      = "tf-example"
  instance_type      = "Common"
  storage_type       = "CLOUD_TCS"
  vpc_id             = "vpc-i5yyodl9"
  subnet_id          = "subnet-5rrirqyc"
  libra_db_version   = "2.2410.18.0"
  src_instance_id    = "cynosdbmysql-ins-84ja0ye0"
  security_group_ids = ["sg-4rd5741x"]
  force_delete       = true
  objects {
    database_tables {
      migrate_db_mode = "all"
    }
  }
}
```

or 

```hcl
resource "tencentcloud_cynosdb_libra_db_instance" "example" {
  cluster_id         = "cynosdbmysql-5oo78wv9"
  zone               = "ap-guangzhou-7"
  cpu                = 8
  mem                = 32
  storage_size       = 100
  pay_mode           = 0
  port               = 2000
  instance_name      = "tf-example"
  instance_type      = "Common"
  storage_type       = "CLOUD_TCS"
  vpc_id             = "vpc-i5yyodl9"
  subnet_id          = "subnet-5rrirqyc"
  libra_db_version   = "2.2410.18.0"
  src_instance_id    = "cynosdbmysql-ins-84ja0ye0"
  security_group_ids = ["sg-4rd5741x"]
  force_delete       = true
  objects {
    database_tables {
      migrate_db_mode = "partial"
      databases {
        db_name            = "test"
        migrate_table_mode = "all"
      }
    }
  }
}
```

Import

CynosDB LibraDB instance can be imported using the cluster_id#instance_id, e.g.

```
terraform import tencentcloud_cynosdb_libra_db_instance.example cynosdbmysql-5oo78wv9#libradb-ins-irehx3rm
```
