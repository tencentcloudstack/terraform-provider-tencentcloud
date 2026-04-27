# Add tencentcloud_teo_load_balancer Resource

## What

Add a new Terraform resource `tencentcloud_teo_load_balancer` for managing Tencent Cloud EdgeOne (TEO) load balancer instances. This resource supports the full CRUD lifecycle for load balancers in the Edge Security Acceleration Platform.

## Why

TEO load balancers allow users to distribute traffic across multiple origin groups with health checking, steering policies, and failover strategies. Currently no Terraform resource exists for managing TEO load balancers, requiring manual portal operations. This resource enables infrastructure-as-code management of load balancer instances.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create | `CreateLoadBalancer` | Returns `InstanceId` (String, e.g. `lb-xxxxxxxx`) |
| Read | `DescribeLoadBalancerList` | Filter by `InstanceId`; returns `LoadBalancerList` |
| Update | `ModifyLoadBalancer` | Accepts `ZoneId` + `InstanceId` + changed fields |
| Delete | `DeleteLoadBalancer` | Accepts `ZoneId` + `InstanceId` |

## Resource ID

`InstanceId` (String, e.g. `lb-2my56s2a4fw7`). Since the resource belongs to a specific zone, the actual Terraform resource ID is composed as `<zone_id>#<instance_id>` to support import.
