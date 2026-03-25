## 1. Schema 定义 - 基础信息字段

- [x] 1.1 在 `clb_list` 元素 Schema 中添加 `forward` 字段(TypeInt, Computed, Description: "负载均衡类型标识,1:负载均衡,0:传统型负载均衡")
- [x] 1.2 添加 `domain` 字段(TypeString, Computed, Description: "负载均衡实例域名(仅公网传统型),逐步下线中")
- [x] 1.3 添加 `load_balancer_domain` 字段(TypeString, Computed, Description: "负载均衡实例的域名")

## 2. Schema 定义 - 网络配置字段

- [x] 2.1 添加 `address_ipv6` 字段(TypeString, Computed, Description: "负载均衡实例的IPv6地址")
- [x] 2.2 添加 `ipv6_mode` 字段(TypeString, Computed, Description: "IP地址版本为ipv6时的模式,IPv6Nat64 | IPv6FullChain")
- [x] 2.3 添加 `mix_ip_target` 字段(TypeBool, Computed, Description: "IPv6FullChain负载均衡7层监听器支持混绑IPv4/IPv6目标")
- [x] 2.4 添加 `anycast_zone` 字段(TypeString, Computed, Description: "anycast负载均衡的发布域,非anycast返回空字符串")
- [x] 2.5 添加 `egress` 字段(TypeString, Computed, Description: "网络出口")
- [x] 2.6 添加 `local_bgp` 字段(TypeBool, Computed, Description: "IP类型是否是本地BGP")

## 3. Schema 定义 - 计费与生命周期字段

- [x] 3.1 添加 `charge_type` 字段(TypeString, Computed, Description: "计费类型,PREPAID:包年包月,POSTPAID_BY_HOUR:按量计费")
- [x] 3.2 添加 `expire_time` 字段(TypeString, Computed, Description: "负载均衡实例的过期时间,仅对预付费负载均衡生效,格式:YYYY-MM-DD HH:mm:ss")
- [x] 3.3 添加 `prepaid_period` 字段(TypeInt, Computed, Description: "预付费购买时长,单位:月")
- [x] 3.4 添加 `prepaid_renew_flag` 字段(TypeString, Computed, Description: "预付费续费标识,NOTIFY_AND_AUTO_RENEW:通知并自动续费,NOTIFY_AND_MANUAL_RENEW:通知但不自动续费,DISABLE_NOTIFY_AND_MANUAL_RENEW:不通知且不自动续费")
- [x] 3.5 **已跳过** `prepaid_cur_instance_deadline` 字段(SDK中LBChargePrepaid结构体不包含此字段)

## 4. Schema 定义 - 日志配置字段

- [x] 4.1 添加 `log_set_id` 字段(TypeString, Computed, Description: "负载均衡日志服务(CLS)的日志集ID")
- [x] 4.2 添加 `log_topic_id` 字段(TypeString, Computed, Description: "负载均衡日志服务(CLS)的日志主题ID")
- [x] 4.3 添加 `health_log_set_id` 字段(TypeString, Computed, Description: "负载均衡日志服务(CLS)的健康检查日志集ID")
- [x] 4.4 添加 `health_log_topic_id` 字段(TypeString, Computed, Description: "负载均衡日志服务(CLS)的健康检查日志主题ID")

## 5. Schema 定义 - 安全与隔离字段

- [x] 5.1 添加 `open_bgp` 字段(TypeInt, Computed, Description: "高防LB的标识,1:高防负载均衡,0:非高防负载均衡")
- [x] 5.2 添加 `snat` 字段(TypeBool, Computed, Description: "是否开启SNAT")
- [x] 5.3 添加 `snat_pro` 字段(TypeBool, Computed, Description: "是否开启SnatPro")
- [x] 5.4 添加 `snat_ips` 字段(TypeString, Computed, Description: "开启SnatPro负载均衡后的SnatIp列表(JSON格式)")
- [x] 5.5 添加 `isolation` 字段(TypeInt, Computed, Description: "是否被隔离,0:表示未被隔离,1:表示被隔离")
- [x] 5.6 添加 `isolated_time` 字段(TypeString, Computed, Description: "负载均衡实例被隔离的时间,格式:YYYY-MM-DD HH:mm:ss")
- [x] 5.7 添加 `is_block` 字段(TypeBool, Computed, Description: "vip是否被封堵")
- [x] 5.8 添加 `is_block_time` 字段(TypeString, Computed, Description: "封堵或解封时间,格式:YYYY-MM-DD HH:mm:ss")
- [x] 5.9 添加 `is_ddos` 字段(TypeBool, Computed, Description: "是否可绑定高防包")

## 6. Schema 定义 - 性能与容量字段

- [x] 6.1 添加 `sla_type` 字段(TypeString, Computed, Description: "性能容量型规格(clb.c1.small/clb.c2.medium/clb.c3.small/clb.c3.medium/clb.c4.small/clb.c4.medium/clb.c4.large/clb.c4.xlarge或空字符串)")
- [x] 6.2 添加 `exclusive` 字段(TypeInt, Computed, Description: "实例类型是否为独占型,1:独占型实例,0:非独占型实例")
- [x] 6.3 添加 `target_count` 字段(TypeInt, Computed, Description: "已绑定的后端服务数量")

## 7. Schema 定义 - 集群与部署字段

- [x] 7.1 添加 `cluster_ids` 字段(TypeList of TypeString, Computed, Description: "集群ID列表")
- [x] 7.2 添加 `cluster_tag` 字段(TypeString, Computed, Description: "7层独占标签")
- [x] 7.3 添加 `nfv_info` 字段(TypeString, Computed, Description: "CLB是否为NFV,空:不是,l7nfv:七层是NFV")
- [x] 7.4 添加 `backup_zone_set` 字段(TypeList of TypeMap, Computed, Description: "备可用区列表,每个元素包含zone_id/zone/zone_name/zone_region/local_zone")
- [x] 7.5 添加 `available_zone_affinity_info` 字段(TypeString, Computed, Description: "可用区转发亲和信息(JSON格式)")

## 8. Schema 定义 - 高级配置字段

- [x] 8.1 添加 `config_id` 字段(TypeString, Computed, Description: "负载均衡维度的个性化配置ID")
- [x] 8.2 添加 `load_balancer_pass_to_target` 字段(TypeBool, Computed, Description: "后端服务是否放通来自LB的流量")
- [x] 8.3 添加 `attribute_flags` 字段(TypeList of TypeString, Computed, Description: "负载均衡的属性标志数组")
- [x] 8.4 添加 `exclusive_cluster` 字段(TypeString, Computed, Description: "内网独占集群信息(JSON格式)")
- [x] 8.5 添加 `extra_info` 字段(TypeString, Computed, Description: "暂做保留,一般用户无需关注(JSON格式)")
- [x] 8.6 添加 `associate_endpoint` 字段(TypeString, Computed, Description: "负载均衡实例关联的Endpoint id")

## 9. 数据映射 - 基础信息字段

- [x] 9.1 在 `dataSourceTencentCloudClbInstancesRead` 函数的 mapping 中添加 `forward` 字段映射:`"forward": clbInstance.Forward`
- [x] 9.2 添加 `domain` 字段映射(需nil检查)
- [x] 9.3 添加 `load_balancer_domain` 字段映射(需nil检查)

## 10. 数据映射 - 网络配置字段

- [x] 10.1 添加 `address_ipv6` 字段映射(需nil检查)
- [x] 10.2 添加 `ipv6_mode` 字段映射(需nil检查)
- [x] 10.3 添加 `mix_ip_target` 字段映射(需nil检查)
- [x] 10.4 添加 `anycast_zone` 字段映射(需nil检查)
- [x] 10.5 添加 `egress` 字段映射(需nil检查)
- [x] 10.6 添加 `local_bgp` 字段映射(需nil检查)

## 11. 数据映射 - 计费与生命周期字段

- [x] 11.1 添加 `charge_type` 字段映射(需nil检查)
- [x] 11.2 添加 `expire_time` 字段映射(需nil检查)
- [x] 11.3 添加 `PrepaidAttributes` 的嵌套字段映射:首先检查 `clbInstance.PrepaidAttributes != nil`,然后分别映射 `prepaid_period`, `prepaid_renew_flag` (注: SDK中不包含CurInstanceDeadline字段)

## 12. 数据映射 - 日志配置字段

- [x] 12.1 添加 `log_set_id` 字段映射(需nil检查)
- [x] 12.2 添加 `log_topic_id` 字段映射(需nil检查)
- [x] 12.3 添加 `health_log_set_id` 字段映射(需nil检查)
- [x] 12.4 添加 `health_log_topic_id` 字段映射(需nil检查)

## 13. 数据映射 - 安全与隔离字段

- [x] 13.1 添加 `open_bgp` 字段映射(需nil检查)
- [x] 13.2 添加 `snat` 字段映射(需nil检查)
- [x] 13.3 添加 `snat_pro` 字段映射(需nil检查)
- [x] 13.4 添加 `snat_ips` 字段映射:检查 `clbInstance.SnatIps != nil`,使用 `json.Marshal` 序列化为 JSON 字符串
- [x] 13.5 添加 `isolation` 字段映射(需nil检查)
- [x] 13.6 添加 `isolated_time` 字段映射(需nil检查)
- [x] 13.7 添加 `is_block` 字段映射(需nil检查)
- [x] 13.8 添加 `is_block_time` 字段映射(需nil检查)
- [x] 13.9 添加 `is_ddos` 字段映射(需nil检查)

## 14. 数据映射 - 性能与容量字段

- [x] 14.1 添加 `sla_type` 字段映射(需nil检查)
- [x] 14.2 添加 `exclusive` 字段映射(需nil检查)
- [x] 14.3 添加 `target_count` 字段映射(需nil检查)

## 15. 数据映射 - 集群与部署字段

- [x] 15.1 添加 `cluster_ids` 字段映射:使用 `helper.StringsInterfaces(clbInstance.ClusterIds)` (需nil检查)
- [x] 15.2 添加 `cluster_tag` 字段映射(需nil检查)
- [x] 15.3 添加 `nfv_info` 字段映射(需nil检查)
- [x] 15.4 添加 `backup_zone_set` 字段映射:检查 `clbInstance.BackupZoneSet != nil`,遍历构造 map 数组
- [x] 15.5 添加 `available_zone_affinity_info` 字段映射:检查nil后使用 `json.Marshal` 序列化

## 16. 数据映射 - 高级配置字段

- [x] 16.1 添加 `config_id` 字段映射(需nil检查)
- [x] 16.2 添加 `load_balancer_pass_to_target` 字段映射(需nil检查)
- [x] 16.3 添加 `attribute_flags` 字段映射:使用 `helper.StringsInterfaces(clbInstance.AttributeFlags)` (需nil检查)
- [x] 16.4 添加 `exclusive_cluster` 字段映射:检查nil后使用 `json.Marshal` 序列化
- [x] 16.5 添加 `extra_info` 字段映射:检查nil后使用 `json.Marshal` 序列化
- [x] 16.6 添加 `associate_endpoint` 字段映射(需nil检查)

## 17. 文档更新

- [x] 17.1 在 `data_source_tc_clb_instances.md` 的 Attributes Reference 部分按分组添加所有新增字段的说明
- [x] 17.2 为每个字段添加清晰的描述、可选值范围和注意事项(如可能返回null)
- [x] 17.3 对于 JSON 格式字段,提供 `jsondecode()` 使用示例
- [x] 17.4 运行 `make doc` 生成 `website/docs/d/clb_instances.html.markdown`

## 18. 代码验证

- [x] 18.1 运行 `go build ./tencentcloud/services/clb/...` 验证代码编译通过
- [x] 18.2 运行 `gofmt -l tencentcloud/services/clb/data_source_tc_clb_instances.go` 检查代码格式化
- [x] 18.3 运行 `go vet ./tencentcloud/services/clb/...` 检查代码问题

## 19. 测试

- [x] 19.1 运行现有测试用例验证向后兼容性:`TF_ACC=1 go test ./tencentcloud/services/clb -run TestAccTencentCloudClbInstancesDataSource`
- [x] 19.2 在测试环境创建不同类型的CLB实例(公网/内网、预付费/后付费、不同规格)
- [x] 19.3 验证所有新增字段能正确返回数据或处理null值
- [x] 19.4 验证嵌套对象和JSON序列化字段的正确性

## 20. 最终检查

- [x] 20.1 检查所有新增字段均为 Computed 属性
- [x] 20.2 检查所有指针访问前都有nil检查
- [x] 20.3 检查文档已更新并生成
- [x] 20.4 验证现有用户配置不受影响(向后兼容)
