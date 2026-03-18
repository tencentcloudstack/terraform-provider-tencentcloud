# Change: Add Missing Fields to CLS Data Transform Resource

## Why
用户反馈 `tencentcloud_cls_data_transform` 资源当前支持的字段不完整，缺少腾讯云 CLS CreateDataTransform API 中的多个重要字段，导致用户无法通过 Terraform 配置完整的数据加工任务功能。根据[创建数据加工任务 API 文档](https://cloud.tencent.com/document/api/614/72184)，需要补齐缺失的字段以提供完整的功能支持。

## What Changes
补齐 `ResourceTencentCloudClsDataTransform` 资源中缺失的以下字段：

1. **preview_log_statistics** (可选) - 用于预览加工结果的测试数据
   - 类型：列表对象
   - 字段：log_content, line_num, dst_topic_id

2. **backup_give_up_data** (可选) - 当 FuncType 为 2 时，动态创建的日志集、日志主题的个数超出产品规格限制是否丢弃数据
   - 类型：布尔值
   - 默认值：false

3. **has_services_log** (可选) - 是否开启投递服务日志
   - 类型：整数
   - 可选值：1（关闭）、2（开启）

4. **data_transform_type** (可选) - 数据加工类型
   - 类型：整数
   - 可选值：0（标准加工任务）、1（前置加工任务）

5. **keep_failure_log** (可选) - 保留失败日志状态
   - 类型：整数
   - 可选值：1（不保留，默认）、2（保留）

6. **failure_log_key** (可选) - 失败日志的字段名称
   - 类型：字符串

7. **process_from_timestamp** (可选) - 指定加工数据的开始时间（秒级时间戳）
   - 类型：整数

8. **process_to_timestamp** (可选) - 指定加工数据的结束时间（秒级时间戳）
   - 类型：整数

9. **data_transform_sql_data_sources** (可选) - 关联的数据源信息
   - 类型：列表对象
   - 字段：data_source, region, instance_id, user, alias_name, password

10. **env_infos** (可选) - 设置的环境变量
    - 类型：列表对象
    - 字段：key, value

## Impact
- **Affected specs**: cls-data-transform
- **Affected code**: 
  - `tencentcloud/services/cls/resource_tc_cls_data_transform.go` - Schema 定义和 Create 函数
  - 需要同时更新 Read 和 Update 函数以支持新字段
- **Breaking changes**: 无，所有新增字段均为可选字段
- **Benefits**: 用户可以通过 Terraform 完整配置 CLS 数据加工任务的所有功能
