## 1. Schema 定义与资源框架

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_load_balancer.go` 文件，定义 `ResourceTencentCloudTeoLoadBalancer()` 函数及完整 schema，包含：zone_id(Required, ForceNew)、name(Required)、type(Required, ForceNew)、origin_groups(Required, TypeList)、health_checker(Optional, TypeList MaxItems:1)、steering_policy(Optional)、failover_policy(Optional)，以及 computed 属性：instance_id、status、origin_group_health_status、l4_used_list、l7_used_list、references
- [x] 1.2 定义 origin_groups 嵌套 block schema：priority(Required, string)、origin_group_id(Required, string)
- [x] 1.3 定义 health_checker 嵌套 block schema：type(Required)、port(Optional)、interval(Optional)、timeout(Optional)、health_threshold(Optional)、critical_threshold(Optional)、path(Optional)、method(Optional)、expected_codes(Optional, TypeList)、headers(Optional, TypeList 嵌套 key/value)、follow_redirect(Optional)、send_context(Optional)、recv_context(Optional)
- [x] 1.4 定义 origin_group_health_status、references 等 computed 嵌套 block schema

## 2. CRUD 函数实现

- [x] 2.1 实现 `resourceTencentCloudTeoLoadBalancerCreate` 函数：构建 CreateLoadBalancerRequest，设置 zone_id/name/type/origin_groups/health_checker/steering_policy/failover_policy，调用 API 并使用 WriteRetryTimeout 重试，验证响应 InstanceId 非空，设置联合 ID（zone_id:instance_id），Create 完成后调用 Read 刷新状态
- [x] 2.2 实现 `resourceTencentCloudTeoLoadBalancerRead` 函数：从 d.Id() 解析 zone_id 和 instance_id，构建 DescribeLoadBalancerListRequest，设置 ZoneId、Filter(InstanceId)、Limit=1，使用 ReadRetryTimeout 重试，根据响应 LoadBalancerList 设置所有 schema 字段（注意判断 nil），若资源不存在则清空 ID
- [x] 2.3 实现 `resourceTencentCloudTeoLoadBalancerUpdate` 函数：从 d.Id() 解析 zone_id 和 instance_id，构建 ModifyLoadBalancerRequest，设置可更新字段（name/origin_groups/health_checker/steering_policy/failover_policy），使用 WriteRetryTimeout 重试，Update 完成后调用 Read 刷新状态
- [x] 2.4 实现 `resourceTencentCloudTeoLoadBalancerDelete` 函数：从 d.Id() 解析 zone_id 和 instance_id，构建 DeleteLoadBalancerRequest，使用 WriteRetryTimeout 重试

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_load_balancer` 资源（参考 tencentcloud_igtm_strategy 资源的注册方式）
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_load_balancer` 资源文档条目

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_load_balancer_test.go` 文件，使用 gomonkey mock 方式编写单元测试，覆盖 Create/Read/Update/Delete 函数的业务逻辑

## 5. 资源文档

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_load_balancer.md` 文件，包含一句话描述（带云产品名称 TEO）、Example Usage（使用 jsonencode() 处理 JSON 字段）、Import 部分（说明使用联合 ID 格式 zone_id#instance_id）
