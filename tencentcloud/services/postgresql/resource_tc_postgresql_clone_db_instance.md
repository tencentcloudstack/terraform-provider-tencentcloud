Provides a resource to create a postgresql clone db instance

Example Usage

Clone db instance by recovery_target_time

```hcl
resource "tencentcloud_postgresql_clone_db_instance" "example" {
  db_instance_id       = "postgres-ckwcgdf1"
  name                 = "tf-example"
  spec_code            = "pg.it.medium4"
  storage              = 100
  period               = 1
  auto_renew_flag      = 0
  vpc_id               = "vpc-i5yyodl9"
  subnet_id            = "subnet-hhi88a58"
  instance_charge_type = "POSTPAID_BY_HOUR"
  security_group_ids   = ["sg-rs32zv1r", "sg-37tigqat"]
  project_id           = 0
  recovery_target_time = "2026-07-10 01:00:06"
  deletion_protection  = true
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }

  db_node_set {
    role = "Standby"
    zone = "ap-guangzhou-7"
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

Clone db instance by backup_set_id

```hcl
data "tencentcloud_postgresql_base_backups" "base_backups" {
  filters {
    name   = "db-instance-id"
    values = ["postgres-evsqpyap"]
  }

  order_by      = "Size"
  order_by_type = "asc"
}

resource "tencentcloud_postgresql_clone_db_instance" "example" {
  db_instance_id       = "postgres-evsqpyap"
  name                 = "tf-example-clone"
  spec_code            = "pg.it.medium4"
  storage              = 200
  period               = 1
  auto_renew_flag      = 0
  vpc_id               = "vpc-a6zec4mf"
  subnet_id            = "subnet-b8hintyy"
  instance_charge_type = "POSTPAID_BY_HOUR"
  security_group_ids   = ["sg-8stavs03"]
  project_id           = 0
  backup_set_id        = data.tencentcloud_postgresql_base_backups.base_backups.base_backup_set.0.id
  deletion_protection  = true
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }

  db_node_set {
    role = "Standby"
    zone = "ap-guangzhou-6"
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

Clone db instance from CDC

```hcl
resource "tencentcloud_postgresql_clone_db_instance" "example" {
  db_instance_id       = "postgres-evsqpyap"
  name                 = "tf-example-clone"
  spec_code            = "pg.it.medium4"
  storage              = 200
  period               = 1
  auto_renew_flag      = 0
  vpc_id               = "vpc-a6zec4mf"
  subnet_id            = "subnet-b8hintyy"
  instance_charge_type = "POSTPAID_BY_HOUR"
  security_group_ids   = ["sg-8stavs03"]
  project_id           = 0
  recovery_target_time = "2024-10-12 18:17:00"
  deletion_protection  = true
  db_node_set {
    role                 = "Primary"
    zone                 = "ap-guangzhou-6"
    dedicated_cluster_id = "cluster-262n63e8"
  }

  db_node_set {
    role                 = "Standby"
    zone                 = "ap-guangzhou-6"
    dedicated_cluster_id = "cluster-262n63e8"
  }

  tags = {
    tagKey = "tagValue"
  }
}
```
