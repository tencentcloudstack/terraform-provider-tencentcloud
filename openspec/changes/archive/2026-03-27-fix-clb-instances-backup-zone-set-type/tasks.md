# Tasks: Fix CLB Instances Backup Zone Set Data Type

## Preparation
- [ ] 确认问题范围: 验证 `backup_zone_set` 是唯一存在此类型问题的字段
- [ ] 检查 SDK 版本: 确认 `tencentcloud-sdk-go/tencentcloud/clb/v20180317` 中 `ZoneInfo` 结构定义
- [ ] 备份当前代码: 记录修改前的代码状态

## Implementation

### 1. 修复 Schema 定义
- [ ] 打开文件 `tencentcloud/services/clb/data_source_tc_clb_instances.go`
- [ ] 定位到 `backup_zone_set` 字段定义 (约第360-365行)
- [ ] 将 `Elem: &schema.Schema{Type: schema.TypeMap}` 替换为:
  ```go
  Elem: &schema.Resource{
      Schema: map[string]*schema.Schema{
          "zone_id": {
              Type:        schema.TypeInt,
              Computed:    true,
              Description: "Availability zone ID (numerical representation).",
          },
          "zone": {
              Type:        schema.TypeString,
              Computed:    true,
              Description: "Availability zone unique identifier (string representation).",
          },
          "zone_name": {
              Type:        schema.TypeString,
              Computed:    true,
              Description: "Availability zone name.",
          },
          "zone_region": {
              Type:        schema.TypeString,
              Computed:    true,
              Description: "Region that this availability zone belongs to.",
          },
          "local_zone": {
              Type:        schema.TypeBool,
              Computed:    true,
              Description: "Whether this is a local zone.",
          },
      },
  },
  ```
- [ ] 保存文件

### 2. 验证数据处理逻辑
- [ ] 检查第623-645行的数据转换代码
- [ ] 确认每个字段的类型转换与 SDK 定义一致:
  - `zone.ZoneId` (*int64) → `backupZone["zone_id"]` (TypeInt) ✅
  - `zone.Zone` (*string) → `backupZone["zone"]` (TypeString) ✅
  - `zone.ZoneName` (*string) → `backupZone["zone_name"]` (TypeString) ✅
  - `zone.ZoneRegion` (*string) → `backupZone["zone_region"]` (TypeString) ✅
  - `zone.LocalZone` (*bool) → `backupZone["local_zone"]` (TypeBool) ✅
- [ ] 确认逻辑正确，无需修改

### 3. 代码格式化
- [ ] 在项目根目录执行: `go fmt tencentcloud/services/clb/data_source_tc_clb_instances.go`
- [ ] 检查格式化输出，确保无错误
- [ ] 再次查看文件，确认格式正确

## Testing

### 4. 编译验证
- [ ] 执行 `go build` 确保代码可以编译通过
- [ ] 检查是否有类型错误或语法错误

### 5. 单元测试 (可选)
- [ ] 如果存在相关测试，运行: `go test -v ./tencentcloud/services/clb/ -run TestAccTencentCloudClbInstancesDataSource`
- [ ] 验证测试通过

### 6. 手动验证 (推荐)
- [ ] 创建测试 Terraform 配置:
  ```hcl
  data "tencentcloud_clb_instances" "test" {
    clb_id = "lb-xxxxxx"  # 使用实际的 CLB ID
  }
  
  output "backup_zones" {
    value = data.tencentcloud_clb_instances.test.clb_list[0].backup_zone_set
  }
  ```
- [ ] 执行 `terraform init`
- [ ] 执行 `terraform plan`
- [ ] 验证 `backup_zone_set` 字段能正确读取且类型正确

## Documentation

### 7. 更新文档 (可选)
- [ ] 检查 `tencentcloud/services/clb/data_source_tc_clb_instances.md` 文件
- [ ] 如果文档中有 `backup_zone_set` 字段的描述，确保类型信息准确
- [ ] 示例代码中如果引用了此字段，确保示例正确

## Validation

### 8. 最终检查
- [ ] 确认所有文件已保存
- [ ] 确认代码已格式化
- [ ] 确认编译无错误
- [ ] 确认类型定义与 SDK 一致
- [ ] 记录本次修改的变更点

## Cleanup
- [ ] 删除测试配置文件 (如果创建了临时测试)
- [ ] 提交代码前再次运行 `go fmt`

## Notes
- 本次修改仅涉及 schema 类型定义，不修改数据处理逻辑
- 数据处理逻辑 (第623-645行) 已经正确实现，无需改动
- 修改完成后必须执行 `go fmt` 格式化代码
- SDK 参考: `tencentcloud-sdk-go/tencentcloud/clb/v20180317/models.go` 中的 `ZoneInfo` 结构
- API 文档: https://cloud.tencent.com/document/api/214/30685
