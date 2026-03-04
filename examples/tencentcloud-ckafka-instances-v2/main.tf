# Example 1: Query by instance_id
data "tencentcloud_ckafka_instances_v2" "by_instance_id" {
  instance_id = "ckafka-f9ife4zz"
}

output "instance_by_id" {
  value = data.tencentcloud_ckafka_instances_v2.by_instance_id.instance_list
}

# Example 2: Query by search_word (fuzzy search)
data "tencentcloud_ckafka_instances_v2" "by_search_word" {
  search_word = "test"
}

output "instances_by_search" {
  value = data.tencentcloud_ckafka_instances_v2.by_search_word.instance_list
}

# Example 3: Query by status
data "tencentcloud_ckafka_instances_v2" "by_status" {
  status = [1, 5] # 1: running, 5: isolated
}

output "instances_by_status" {
  value = data.tencentcloud_ckafka_instances_v2.by_status.instance_list
}

# Example 4: Query by tag_key
data "tencentcloud_ckafka_instances_v2" "by_tag_key" {
  tag_key = "env"
}

output "instances_by_tag_key" {
  value = data.tencentcloud_ckafka_instances_v2.by_tag_key.instance_list
}

# Example 5: Query by instance_id_list
data "tencentcloud_ckafka_instances_v2" "by_id_list" {
  instance_id_list = ["ckafka-f9ife4zz", "ckafka-abc123"]
}

output "instances_by_id_list" {
  value = data.tencentcloud_ckafka_instances_v2.by_id_list.instance_list
}

# Example 6: Query by tag_list (intersection)
data "tencentcloud_ckafka_instances_v2" "by_tag_list" {
  tag_list {
    tag_key   = "env"
    tag_value = "prod"
  }
  tag_list {
    tag_key   = "project"
    tag_value = "demo"
  }
}

output "instances_by_tag_list" {
  value = data.tencentcloud_ckafka_instances_v2.by_tag_list.instance_list
}

# Example 7: Query by filters - VPC
data "tencentcloud_ckafka_instances_v2" "by_vpc" {
  filters {
    name   = "VpcId"
    values = ["vpc-abc123"]
  }
}

output "instances_by_vpc" {
  value = data.tencentcloud_ckafka_instances_v2.by_vpc.instance_list
}

# Example 8: Query by filters - Multiple conditions
data "tencentcloud_ckafka_instances_v2" "by_multiple_filters" {
  filters {
    name   = "VpcId"
    values = ["vpc-abc123"]
  }
  filters {
    name   = "SubNetId"
    values = ["subnet-xyz789"]
  }
  filters {
    name   = "InstanceType"
    values = ["1"] # 1: dedicated
  }
}

output "instances_by_multiple_filters" {
  value = data.tencentcloud_ckafka_instances_v2.by_multiple_filters.instance_list
}

# Example 9: Combined query
data "tencentcloud_ckafka_instances_v2" "combined" {
  search_word = "prod"
  status      = [1] # running only

  filters {
    name   = "VpcId"
    values = ["vpc-abc123"]
  }

  tag_list {
    tag_key   = "env"
    tag_value = "prod"
  }
}

output "instances_combined" {
  value = data.tencentcloud_ckafka_instances_v2.combined.instance_list
}

# Example 10: Query all instances
data "tencentcloud_ckafka_instances_v2" "all" {}

output "all_instances" {
  value = data.tencentcloud_ckafka_instances_v2.all.instance_list
}

# Output detailed information
output "instance_details" {
  value = length(data.tencentcloud_ckafka_instances_v2.all.instance_list) > 0 ? {
    id        = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].instance_id
    name      = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].instance_name
    vip       = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].vip
    vport     = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].vport
    status    = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].status
    bandwidth = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].bandwidth
    disk_size = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].disk_size
    vpc_id    = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].vpc_id
    subnet_id = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].subnet_id
    zone_id   = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].zone_id
    healthy   = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].healthy
    topic_num = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].topic_num
    version   = data.tencentcloud_ckafka_instances_v2.all.instance_list[0].version
  } : null
}
