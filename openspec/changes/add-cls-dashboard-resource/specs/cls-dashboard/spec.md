# CLS Dashboard Resource Spec

## ADDED Requirements

### Req-1: 支持创建仪表盘

#### Scenario: 创建基础仪表盘
**Given** 用户提供了仪表盘名称  
**When** 调用 CreateDashboard API  
**Then** 返回新创建的 dashboard_id  
**And** 仪表盘名称与输入一致  
**And** 配置数据为空（默认值）

#### Scenario: 创建带配置数据的仪表盘
**Given** 用户提供了仪表盘名称和配置数据（JSON 字符串）  
**When** 调用 CreateDashboard API  
**Then** 返回新创建的 dashboard_id  
**And** 配置数据正确保存

#### Scenario: 创建带标签的仪表盘
**Given** 用户提供了仪表盘名称和标签（最多 10 个）  
**When** 调用 CreateDashboard API  
**Then** 返回新创建的 dashboard_id  
**And** 标签正确绑定到仪表盘

#### Scenario: 创建同名仪表盘失败
**Given** 已存在名为 "dashboard-A" 的仪表盘  
**When** 尝试创建名为 "dashboard-A" 的新仪表盘  
**Then** API 返回 InvalidParameter.DashboardNameConflict 错误  
**And** 向用户返回友好的错误提示

#### Scenario: 标签数量超限
**Given** 用户提供了 11 个标签  
**When** 调用 CreateDashboard API  
**Then** API 返回 LimitExceeded.Tag 错误  
**And** 提示用户减少标签数量到 10 个以内

---

### Req-2: 支持查询仪表盘信息

#### Scenario: 查询存在的仪表盘
**Given** 仪表盘 ID 为 "dashboard-xxxx"  
**When** 调用 Service 层的 DescribeClsDashboardById 方法  
**Then** 返回完整的仪表盘信息  
**And** 包含 dashboard_name, data, tags, create_time, update_time

#### Scenario: 查询不存在的仪表盘
**Given** 仪表盘 ID 为 "non-existent-id"  
**When** 调用 Service 层的 DescribeClsDashboardById 方法  
**Then** 返回 nil（不报错）  
**And** Resource Read 函数清空资源 ID

#### Scenario: 分页查询大量仪表盘
**Given** 账户中有 50 个仪表盘  
**And** 目标仪表盘在第 3 页（第 41-60 条）  
**When** 调用 DescribeDashboards API 分页查询  
**Then** 系统遍历多页数据  
**And** 成功找到目标仪表盘  
**And** 返回正确的仪表盘信息

#### Scenario: API 调用失败时重试
**Given** DescribeDashboards API 第一次调用返回 InternalError  
**When** 触发重试机制  
**Then** 在 ReadRetryTimeout 时间内重试  
**And** 如果重试成功则返回结果  
**And** 如果重试失败则返回错误

---

### Req-3: 支持修改仪表盘

#### Scenario: 修改仪表盘名称
**Given** 仪表盘 ID 为 "dashboard-xxxx"  
**When** 用户更新 dashboard_name 字段  
**And** 调用 ModifyDashboard API  
**Then** 仪表盘名称成功更新  
**And** 其他字段保持不变

#### Scenario: 修改仪表盘配置数据
**Given** 仪表盘 ID 为 "dashboard-xxxx"  
**When** 用户更新 data 字段  
**And** 调用 ModifyDashboard API  
**Then** 配置数据成功更新  
**And** 其他字段保持不变

#### Scenario: 修改仪表盘标签
**Given** 仪表盘 ID 为 "dashboard-xxxx"  
**And** 原有标签为 {"env": "dev"}  
**When** 用户更新 tags 为 {"env": "prod", "team": "ops"}  
**And** 调用 ModifyDashboard API  
**Then** 标签成功更新  
**And** 新标签替换旧标签

#### Scenario: 同时修改多个字段
**Given** 仪表盘 ID 为 "dashboard-xxxx"  
**When** 用户同时更新 dashboard_name, data 和 tags  
**And** 调用 ModifyDashboard API  
**Then** 所有字段成功更新

#### Scenario: 修改时名称冲突
**Given** 已存在名为 "dashboard-B" 的仪表盘  
**When** 尝试将另一个仪表盘名称修改为 "dashboard-B"  
**Then** API 返回 InvalidParameter.DashboardNameConflict 错误  
**And** 向用户返回友好的错误提示

---

### Req-4: 支持删除仪表盘

#### Scenario: 删除存在的仪表盘
**Given** 仪表盘 ID 为 "dashboard-xxxx"  
**When** 调用 DeleteDashboard API  
**Then** 仪表盘成功删除  
**And** 后续查询返回不存在

#### Scenario: 删除不存在的仪表盘（幂等性）
**Given** 仪表盘 ID 为 "dashboard-xxxx" 已被删除  
**When** 再次调用 DeleteDashboard API  
**Then** 操作视为成功（不报错）  
**And** 记录日志但不返回错误给用户

#### Scenario: 删除操作重试
**Given** DeleteDashboard API 第一次调用返回 InternalError  
**When** 触发重试机制  
**Then** 在 WriteRetryTimeout 时间内重试  
**And** 如果重试成功则完成删除

---

### Req-5: 支持导入现有仪表盘

#### Scenario: 通过 dashboard_id 导入
**Given** 已存在仪表盘 ID 为 "dashboard-xxxx"  
**When** 执行 terraform import 命令  
**And** 提供 dashboard_id 作为参数  
**Then** 成功导入仪表盘到 Terraform 状态  
**And** 所有字段正确填充

#### Scenario: 导入不存在的仪表盘
**Given** 仪表盘 ID 为 "non-existent-id" 不存在  
**When** 执行 terraform import 命令  
**Then** 导入失败  
**And** 返回友好的错误提示

#### Scenario: 导入后状态一致性
**Given** 通过 import 导入了仪表盘  
**When** 运行 terraform plan  
**Then** 显示无变更（状态一致）

---

### Req-6: 支持标签管理

#### Scenario: 添加标签到新仪表盘
**Given** 创建新仪表盘时提供标签  
**When** 调用 CreateDashboard API  
**Then** 标签成功绑定  
**And** 查询时返回正确的标签

#### Scenario: 更新现有仪表盘的标签
**Given** 仪表盘原有标签为 {"env": "dev"}  
**When** 更新标签为 {"env": "prod", "owner": "admin"}  
**And** 调用 ModifyDashboard API  
**Then** 标签成功更新  
**And** 旧标签被新标签替换

#### Scenario: 删除所有标签
**Given** 仪表盘有标签 {"env": "dev", "team": "ops"}  
**When** 更新 tags 为空 map {}  
**And** 调用 ModifyDashboard API  
**Then** 所有标签被移除

#### Scenario: 标签格式转换
**Given** Terraform 中 tags 定义为 map[string]string  
**And** API 接受 Tag 数组（[{Key, Value}]）  
**When** 进行 Create 或 Update 操作  
**Then** 正确转换 map 为 Tag 数组  
**And** API 调用成功

---

### Req-7: 错误处理和重试机制

#### Scenario: 处理 API 频率限制
**Given** API 返回频率限制错误  
**When** 触发重试机制  
**Then** 等待并重试请求  
**And** 在超时前成功完成

#### Scenario: 处理账户欠费错误
**Given** API 返回 OperationDenied.AccountIsolate 错误  
**When** 执行任何操作  
**Then** 立即返回错误  
**And** 提示用户账户欠费，需要充值

#### Scenario: 处理内部错误重试
**Given** API 返回 InternalError  
**When** 触发重试机制  
**Then** 在配置的超时时间内重试  
**And** 记录重试日志

#### Scenario: 处理参数错误不重试
**Given** API 返回 InvalidParameter 错误  
**When** 检测到参数错误  
**Then** 不进行重试  
**And** 立即返回错误给用户

#### Scenario: Read 操作重试超时
**Given** DescribeDashboards API 持续返回错误  
**When** 超过 ReadRetryTimeout（1 分钟）  
**Then** 停止重试  
**And** 返回最后的错误信息

#### Scenario: Write 操作重试超时
**Given** CreateDashboard/ModifyDashboard/DeleteDashboard API 持续返回错误  
**When** 超过 WriteRetryTimeout（5 分钟）  
**Then** 停止重试  
**And** 返回最后的错误信息
