# Test Case 1: Target group without port field (tests drift fix)
resource "tencentcloud_vpc" "test_vpc" {
  name       = "gwlb-test-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_gwlb_target_group" "without_port" {
  target_group_name = "test-tg-no-port"
  vpc_id            = tencentcloud_vpc.test_vpc.id
  
  health_check {
    health_switch = true
    protocol      = "PING"
  }
}

# Test Case 2: Target group with explicit port value
resource "tencentcloud_gwlb_target_group" "with_port" {
  target_group_name = "test-tg-with-port"
  vpc_id            = tencentcloud_vpc.test_vpc.id
  port              = 6081
  
  health_check {
    health_switch = true
    protocol      = "TCP"
    port          = 6081
    timeout       = 2
    interval_time = 5
    health_num    = 3
    un_health_num = 3
  }
}

# Output values to verify behavior
output "without_port_id" {
  value       = tencentcloud_gwlb_target_group.without_port.id
  description = "Target group ID without port specified"
}

output "without_port_value" {
  value       = tencentcloud_gwlb_target_group.without_port.port
  description = "Port value from API (should be default)"
}

output "with_port_id" {
  value       = tencentcloud_gwlb_target_group.with_port.id
  description = "Target group ID with explicit port"
}

output "with_port_value" {
  value       = tencentcloud_gwlb_target_group.with_port.port
  description = "Port value (should be 6081)"
}
