# Capability: CLB Target Group Resource

## Overview

This specification defines the requirements for the `tencentcloud_clb_target_group` Terraform resource, which manages Cloud Load Balancer (CLB) target groups in Tencent Cloud.

## ADDED Requirements

### Requirement: Support Health Check Configuration

The resource MUST support configuring health checks for target groups through a `health_check` nested block.

**Rationale**: Health checks are critical for ensuring backend service availability. Different protocols and scenarios require different health check strategies.

**Acceptance Criteria**:
- Health check configuration is optional
- All health check parameters map to the CreateTargetGroup API HealthCheck field
- Health check settings are properly read back from the API
- Health check configuration supports protocol-specific parameters

#### Scenario: Configure TCP Health Check

**Given** a user wants to create a target group with TCP health check  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "tcp_hc" {
  target_group_name = "tcp-target-group"
  type              = "v2"
  protocol          = "TCP"
  
  health_check {
    health_switch = true
    protocol      = "TCP"
    port          = 8080
    timeout       = 5
    gap_time      = 10
    good_limit    = 3
    bad_limit     = 2
  }
}
```
**Then** the target group is created with TCP health check enabled  
**And** the health check parameters are applied correctly  
**And** terraform plan shows no changes after apply

#### Scenario: Configure HTTP Health Check with Domain and Path

**Given** a user wants to create a target group with HTTP health check  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "http_hc" {
  target_group_name = "http-target-group"
  type              = "v2"
  protocol          = "HTTP"
  
  health_check {
    health_switch      = true
    protocol           = "HTTP"
    http_check_domain  = "example.com"
    http_check_path    = "/health"
    http_check_method  = "GET"
    http_code          = 31  # 1xx, 2xx, 3xx, 4xx, 5xx all healthy
    http_version       = "HTTP/1.1"
    timeout            = 2
    gap_time           = 5
  }
}
```
**Then** the target group is created with HTTP health check  
**And** health checks use the specified domain, path, method, and accepted status codes  
**And** the configuration persists after creation

#### Scenario: Disable Health Check

**Given** a user wants to create a target group without health checks  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "no_hc" {
  target_group_name = "no-health-check"
  type              = "v2"
  protocol          = "TCP"
  
  health_check {
    health_switch = false
  }
}
```
**Then** the target group is created with health checks disabled  
**And** other health check parameters are ignored

---

### Requirement: Support Scheduling Algorithm Configuration

The resource MUST support specifying a scheduling algorithm for v2 target groups with HTTP/HTTPS/GRPC protocols.

**Rationale**: Different workloads benefit from different load balancing algorithms. Weighted round robin (WRR) is good for evenly distributed traffic, least connections (LEAST_CONN) for varying request durations, and IP hash (IP_HASH) for session affinity.

**Acceptance Criteria**:
- `schedule_algorithm` parameter is optional with default value "WRR"
- Only valid for v2 target groups
- Only applies when protocol is HTTP, HTTPS, or GRPC
- Parameter cannot be modified after creation (ForceNew)
- Valid values: WRR, LEAST_CONN, IP_HASH

#### Scenario: Create Target Group with Least Connections Algorithm

**Given** a user wants to distribute traffic based on backend server load  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "least_conn" {
  target_group_name  = "least-conn-tg"
  type               = "v2"
  protocol           = "HTTP"
  schedule_algorithm = "LEAST_CONN"
}
```
**Then** the target group is created with LEAST_CONN algorithm  
**And** traffic is distributed to servers with fewest active connections  
**And** attempting to change schedule_algorithm forces resource recreation

#### Scenario: Create Target Group with IP Hash for Session Affinity

**Given** a user needs session affinity based on client IP  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "ip_hash" {
  target_group_name  = "ip-hash-tg"
  type               = "v2"
  protocol           = "HTTPS"
  schedule_algorithm = "IP_HASH"
}
```
**Then** the target group uses IP hash algorithm  
**And** requests from the same client IP always route to the same backend server

#### Scenario: Default to WRR Algorithm

**Given** a user creates a v2 HTTP target group without specifying algorithm  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "default_algo" {
  target_group_name = "default-algo-tg"
  type              = "v2"
  protocol          = "HTTP"
}
```
**Then** the target group uses WRR (weighted round robin) algorithm by default  
**And** the state reflects schedule_algorithm = "WRR"

---

### Requirement: Support Resource Tagging

The resource MUST support adding tags to target groups for resource management and billing.

**Rationale**: Tags enable resource organization, cost allocation, and access control in cloud environments.

**Acceptance Criteria**:
- `tags` parameter accepts a map of key-value string pairs
- Tags are created with the target group
- Tags are readable from state
- Tags can be updated without recreating the resource
- Tags integrate with Tencent Cloud's unified tagging system

#### Scenario: Create Target Group with Tags

**Given** a user wants to tag a target group for organization  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "tagged" {
  target_group_name = "production-api"
  type              = "v2"
  protocol          = "HTTP"
  
  tags = {
    Environment = "production"
    Team        = "backend"
    CostCenter  = "engineering"
  }
}
```
**Then** the target group is created with the specified tags  
**And** tags are visible in the Tencent Cloud console  
**And** terraform state contains the tags

#### Scenario: Update Target Group Tags

**Given** an existing target group with tags  
**When** a user modifies the tags:
```hcl
tags = {
  Environment = "production"
  Team        = "platform"      # changed
  CostCenter  = "engineering"
  Owner       = "john@example"  # added
}
```
**Then** tags are updated without recreating the target group  
**And** old tags are removed and new tags are added

---

### Requirement: Support Default Backend Server Weight

The resource MUST support setting a default weight for backend servers in v2 target groups.

**Rationale**: Default weights simplify backend server management when most servers should have the same weight. Servers added without explicit weight use this default.

**Acceptance Criteria**:
- `weight` parameter is optional
- Valid range: [0, 100]
- Only applies to v2 target groups
- v1 target groups ignore this parameter
- Validation rejects values outside the valid range

#### Scenario: Set Default Weight for Backend Servers

**Given** a user wants backend servers to default to weight 50  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "weighted" {
  target_group_name = "weighted-tg"
  type              = "v2"
  protocol          = "TCP"
  weight            = 50
}
```
**Then** the target group is created with default weight 50  
**And** servers added without explicit weight use 50  
**And** servers with explicit weight override the default

#### Scenario: Reject Invalid Weight Value

**Given** a user specifies an invalid weight  
**When** they set weight to 150 (outside [0, 100] range)  
**Then** terraform validate fails with a clear error message  
**And** the resource is not created

---

### Requirement: Support Full Listener Target Group Mode

The resource MUST support creating full listener target groups for v2 target groups.

**Rationale**: Full listener target groups bind to all ports of a listener, eliminating the need for explicit port specification and simplifying multi-port service management.

**Acceptance Criteria**:
- `full_listen_switch` parameter is optional (default false)
- Only valid for v2 target groups
- When true, `port` parameter should not be specified
- Parameter cannot be modified after creation (ForceNew)

#### Scenario: Create Full Listener Target Group

**Given** a user wants a target group for all listener ports  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "full_listen" {
  target_group_name  = "full-listen-tg"
  type               = "v2"
  protocol           = "TCP"
  full_listen_switch = true
  # Note: port is not specified
}
```
**Then** the target group is created as a full listener target group  
**And** it can be attached to listeners without port constraints  
**And** changing full_listen_switch requires resource recreation

---

### Requirement: Support Keep-Alive Connection Control

The resource MUST support enabling keep-alive connections for HTTP/HTTPS target groups.

**Rationale**: Keep-alive connections reduce latency and overhead by reusing TCP connections for multiple HTTP requests, improving performance for high-traffic APIs.

**Acceptance Criteria**:
- `keepalive_enable` parameter is optional (default false)
- Only applies to HTTP and HTTPS target groups
- When enabled, connections are kept alive between requests
- Can be modified after creation (if API supports)

#### Scenario: Enable Keep-Alive for HTTP API

**Given** a user wants to optimize HTTP API performance  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "keepalive" {
  target_group_name = "api-tg"
  type              = "v2"
  protocol          = "HTTP"
  keepalive_enable  = true
}
```
**Then** the target group enables keep-alive connections  
**And** HTTP connections are reused across requests  
**And** connection overhead is reduced

#### Scenario: Keep-Alive Has No Effect on TCP Target Group

**Given** a user specifies keepalive_enable for a TCP target group  
**When** they create:
```hcl
resource "tencentcloud_clb_target_group" "tcp_keepalive" {
  type              = "v2"
  protocol          = "TCP"
  keepalive_enable  = true  # Should be ignored or rejected
}
```
**Then** the API ignores keepalive_enable (protocol not HTTP/HTTPS)  
**Or** terraform provides a warning about incompatible protocol

---

### Requirement: Support Session Persistence Configuration

The resource MUST support configuring session persistence (sticky sessions) for v2 HTTP/HTTPS/GRPC target groups.

**Rationale**: Session persistence ensures requests from the same client are routed to the same backend server, which is essential for applications that maintain session state.

**Acceptance Criteria**:
- `session_expire_time` parameter is optional (default 0 = disabled)
- Valid range: [30, 3600] seconds, or 0 to disable
- Only applies to v2 target groups with HTTP, HTTPS, or GRPC protocols
- Session persistence is based on client identifier (likely cookie or IP)
- Validation rejects values outside valid range

#### Scenario: Enable 30-Minute Session Persistence

**Given** a user needs sticky sessions for a web application  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "sticky" {
  target_group_name   = "web-app-tg"
  type                = "v2"
  protocol            = "HTTP"
  session_expire_time = 1800  # 30 minutes
}
```
**Then** session persistence is enabled with 1800 second timeout  
**And** requests from the same client route to the same server for 30 minutes  
**And** after 30 minutes of inactivity, session expires

#### Scenario: Disable Session Persistence

**Given** a stateless application that doesn't need sticky sessions  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "stateless" {
  target_group_name   = "stateless-api-tg"
  type                = "v2"
  protocol            = "HTTP"
  session_expire_time = 0
}
```
**Then** session persistence is disabled  
**And** each request can be routed to any backend server

#### Scenario: Validate Session Persistence Time Range

**Given** a user sets invalid session persistence time  
**When** they set session_expire_time to 5000 (exceeds maximum 3600)  
**Then** terraform validate fails with a clear error message  
**And** valid range [30-3600] or 0 is documented

---

### Requirement: Support IP Version Configuration

The resource MUST support specifying the IP version type for target groups.

**Rationale**: Target groups may need to support IPv4, IPv6, or dual-stack configurations depending on network architecture and client requirements.

**Acceptance Criteria**:
- `ip_version` parameter is optional
- Common values: IPv4, IPv6, IPv6FullChain
- Parameter cannot be modified after creation (ForceNew)
- API determines default value if not specified

#### Scenario: Create IPv6 Target Group

**Given** a user needs IPv6-only backend services  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "ipv6" {
  target_group_name = "ipv6-tg"
  type              = "v2"
  protocol          = "HTTP"
  ip_version        = "IPv6"
}
```
**Then** the target group only accepts IPv6 addresses  
**And** backend servers must have IPv6 addresses  
**And** changing ip_version requires resource recreation

#### Scenario: Create Dual-Stack Target Group

**Given** a user needs to support both IPv4 and IPv6 clients  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "dual_stack" {
  target_group_name = "dual-stack-tg"
  type              = "v2"
  protocol          = "TCP"
  ip_version        = "IPv6FullChain"
}
```
**Then** the target group supports both IPv4 and IPv6  
**And** backend servers can be reached via either protocol

---

### Requirement: Support SNAT (Source IP Replacement) Configuration

The resource MUST support enabling SNAT (Source Network Address Translation) for target groups.

**Rationale**: SNAT allows the load balancer to replace client source IPs with its own IP when forwarding requests to backend servers. This is useful when backend servers need simplified network routing or when client IP transparency is not required.

**Acceptance Criteria**:
- `snat_enable` parameter is optional (default false)
- When enabled, client source IPs are replaced with the load balancer's IP
- When SNAT is enabled, "transparent client IP" option is disabled, and vice versa
- Can be modified after creation (if API supports)

#### Scenario: Enable SNAT for Simplified Routing

**Given** a user wants to simplify backend network routing  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "snat_enabled" {
  target_group_name = "snat-tg"
  type              = "v2"
  protocol          = "TCP"
  snat_enable       = true
}
```
**Then** the target group enables SNAT  
**And** client source IPs are replaced with load balancer IP  
**And** backend servers see all traffic coming from the load balancer

#### Scenario: Disable SNAT to Preserve Client IP

**Given** a user needs to preserve client source IP for logging or access control  
**When** they specify:
```hcl
resource "tencentcloud_clb_target_group" "no_snat" {
  target_group_name = "no-snat-tg"
  type              = "v2"
  protocol          = "HTTP"
  snat_enable       = false
}
```
**Then** client source IPs are preserved  
**And** backend servers can see the original client IP addresses  
**And** transparent client IP forwarding is enabled

---

## Implementation Notes

### Type and Protocol Constraints

The new parameters have version and protocol constraints:

| Parameter | v1 Support | v2 Support | Protocol Constraints |
|-----------|------------|------------|---------------------|
| health_check | Yes | Yes | All protocols |
| schedule_algorithm | No | Yes | HTTP/HTTPS/GRPC only |
| tags | Yes | Yes | All protocols |
| weight | No | Yes | All protocols |
| full_listen_switch | No | Yes | All protocols |
| keepalive_enable | Yes* | Yes* | HTTP/HTTPS only |
| session_expire_time | No | Yes | HTTP/HTTPS/GRPC only |
| ip_version | Yes | Yes | All protocols |
| snat_enable | Yes | Yes | All protocols |

*Note: keepalive_enable depends on protocol, not version

### ForceNew Parameters

The following parameters cannot be modified after creation and must trigger resource recreation:
- `schedule_algorithm`
- `full_listen_switch`
- `ip_version`

### Health Check Protocol Matrix

Different target group protocols support different health check protocols:

- **TCP target groups**: Health check protocols: TCP, HTTP, CUSTOM
- **UDP target groups**: Health check protocols: PING, CUSTOM
- **HTTP target groups**: Health check protocols: HTTP, TCP
- **HTTPS target groups**: Health check protocols: HTTPS, TCP
- **GRPC target groups**: Health check protocols: GRPC, TCP

### Parameter Interaction Rules

1. **Full Listener Mode**: When `full_listen_switch = true`, the `port` parameter should not be specified
2. **Schedule Algorithm**: Only effective when `type = "v2"` and `protocol` is HTTP/HTTPS/GRPC
3. **Session Persistence**: Only effective when `type = "v2"` and `protocol` is HTTP/HTTPS/GRPC
4. **Keep-Alive**: Only effective when `protocol` is HTTP or HTTPS (independent of version)
5. **Weight**: Only effective for v2 target groups; v1 ignores this parameter

### API Mapping

The new parameters map to the following API request fields in `CreateTargetGroupRequest`:

- `health_check` → `HealthCheck` (TargetGroupHealthCheck object)
- `schedule_algorithm` → `ScheduleAlgorithm` (string)
- `tags` → `Tags` ([]TagInfo)
- `weight` → `Weight` (uint64)
- `full_listen_switch` → `FullListenSwitch` (bool)
- `keepalive_enable` → `KeepaliveEnable` (bool)
- `session_expire_time` → `SessionExpireTime` (uint64)
- `ip_version` → `IpVersion` (string)
- `snat_enable` → `SnatEnable` (bool)

### Validation Rules

Schema-level validation should enforce:
- `schedule_algorithm` in [WRR, LEAST_CONN, IP_HASH]
- `weight` in [0, 100]
- `session_expire_time` in [30, 3600] or 0
- `health_check.timeout` in [2, 30]
- `health_check.gap_time` in [2, 300]
- `health_check.good_limit` in [2, 10]
- `health_check.bad_limit` in [2, 10]

API-level validation will enforce protocol/version constraints.

### Update Support

Research needed to determine which parameters support in-place updates via `ModifyTargetGroupAttribute` API:
- Parameters marked ForceNew are confirmed non-updatable
- Other parameters need API documentation review
- Testing should verify update behavior

### Backward Compatibility

All new parameters are optional. Existing configurations without these parameters will continue to work without modification. The API provides sensible defaults for all optional parameters.
