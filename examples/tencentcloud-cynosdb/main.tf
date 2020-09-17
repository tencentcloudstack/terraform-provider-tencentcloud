provider "tencentcloud" {
  region = "ap-guangzhou"
}

resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = "vpc-h70b6b49"
  subnet_id                    = "subnet-q6fhy1mi"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb"
  password                     = "cynos@123"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2

  tags = {
    test = "test"
  }

  force_delete = false

  rw_group_sg = [
    "sg-ibyjkl6r",
  ]
  ro_group_sg = [
    "sg-ibyjkl6r",
  ]
}

resource "tencentcloud_cynosdb_readonly_instance" "foo" {
  cluster_id           = tencentcloud_cynosdb_cluster.foo.id
  instance_name        = "tf-cynosdb-readonly-instance"
  force_delete         = true
  instance_cpu_core    = 1
  instance_memory_size = 2

  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]
}

data "tencentcloud_cynosdb_clusters" "foo" {
  cluster_id   = "cynosdbmysql-dzj5l8gz"
  project_id   = 0
  db_type      = "MYSQL"
  cluster_name = "test"
}

data "tencentcloud_cynosdb_instances" "foo" {
  instance_id   = "cynosdbmysql-ins-0wln9u6w"
  project_id    = 0
  db_type       = "MYSQL"
  instance_name = "test"
}

