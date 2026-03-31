# Tasks: Update tencentcloud_waf_cc options_arr Description

## Task 1: Update options_arr Description in resource_tc_waf_cc.go

**文件**: `tencentcloud/services/waf/resource_tc_waf_cc.go`

**位置**: Lines 82-103 (options_arr field in Schema)

### 实施步骤

#### 1.1 定位 options_arr 字段

找到 Schema 中的 `options_arr` 字段定义（约第 82 行）：
```go
"options_arr": {
    Optional: true,
    Type:     schema.TypeString,
    Description: `...`,  // 当前描述
},
```

#### 1.2 替换 Description 内容

将当前的 Description 字符串（使用反引号 \`\`）替换为以下内容（使用双引号 ""）：

```go
Description: "CC matching conditions JSON serialized string. " +
    "Example: [{\"key\":\"Method\",\"args\":[\"=R0VU\"],\"match\":\"0\",\"encodeflag\":true}]. " +
    "\n\nSupported key types: URL, Method, Post, Referer, Cookie, User-Agent, CustomHeader, IPLocation, CaptchaRisk, CaptchaDeviceRisk, CaptchaScore. " +
    "\n\nMatch operators by key type:\n" +
    "- When Key is URL: 0 (equal), 3 (not equal), 1 (prefix), 6 (suffix), 2 (contains), 7 (not contains)\n" +
    "- When Key is Method: 0 (equal), 3 (not equal)\n" +
    "- When Key is Post: 0 (equal), 3 (not equal), 2 (contains), 7 (not contains)\n" +
    "- When Key is Cookie: 0 (equal), 3 (not equal), 2 (contains), 7 (not contains)\n" +
    "- When Key is Referer: 0 (equal), 3 (not equal), 1 (prefix), 6 (suffix), 2 (contains), 7 (not contains), 12 (exists), 5 (not exists), 4 (empty)\n" +
    "- When Key is User-Agent: 0 (equal), 3 (not equal), 1 (prefix), 6 (suffix), 2 (contains), 7 (not contains), 12 (exists), 5 (not exists), 4 (empty)\n" +
    "- When Key is CustomHeader: 0 (equal), 3 (not equal), 2 (contains), 7 (not contains), 4 (empty), 5 (not exists)\n" +
    "- When Key is IPLocation: 13 (belongs to), 14 (not belongs to)\n" +
    "- When Key is CaptchaRisk: 15 (numerically equal), 16 (numerically not equal), 13 (belongs to), 14 (not belongs to), 12 (exists), 5 (not exists)\n" +
    "- When Key is CaptchaDeviceRisk: 13 (belongs to), 14 (not belongs to), 12 (exists), 5 (not exists)\n" +
    "- When Key is CaptchaScore: 15 (numerically equal), 17 (numerically greater than), 18 (numerically less than), 19 (numerically greater than or equal), 20 (numerically less than or equal), 12 (exists), 5 (not exists)\n" +
    "\n" +
    "Encoding rules: The args parameter requires encodeflag to be set to true. " +
    "For Post, Cookie, or CustomHeader keys, Base64 encode both parameter name and value (remove trailing =), then concatenate with = sign (e.g., Base64(name)=Base64(value)). " +
    "For Referer or User-Agent keys, Base64 encode the value (remove trailing =) and prefix with = sign (e.g., =Base64(value)).",
```

**注意**:
- 使用双引号 `"` 包裹整个描述，而不是反引号 \`
- 使用 `+` 连接多行字符串以提高可读性
- 在 JSON 示例中，内部的双引号需要转义为 `\"`
- 使用 `\n` 表示换行符

#### 1.3 格式化代码

修改完成后，运行以下命令格式化代码：

```bash
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
gofmt -w tencentcloud/services/waf/resource_tc_waf_cc.go
```

#### 1.4 保存文件

确保文件已保存并格式化正确。

### 验收标准

- [x] `options_arr` 字段的 Description 已更新
- [x] Description 使用双引号 `""` 而非反引号 \`\`
- [x] 包含所有 11 种 key 类型：URL, Method, Post, Referer, Cookie, User-Agent, CustomHeader, IPLocation, CaptchaRisk, CaptchaDeviceRisk, CaptchaScore
- [x] 每种 key 类型的 match 操作符都已正确列出
- [x] 包含 JSON 示例
- [x] 包含编码规则说明
- [x] 代码已使用 `go fmt` 格式化
- [x] 文件已保存
- [x] 无语法错误

---

## Task 2: 验证修改

### 2.1 语法检查

运行以下命令检查语法错误：

```bash
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
go build ./tencentcloud/services/waf/
```

### 2.2 检查格式

确认 `go fmt` 已成功运行，代码格式正确。

### 验收标准

- [x] `go build` 编译成功，无语法错误
- [x] 代码格式符合 Go 规范
- [x] 无 linter 错误

---

## 验收清单

### 文档准确性

- [x] Description 内容与 API 文档（https://cloud.tencent.com/document/api/627/97646）一致
- [x] 所有 key 类型均已包含
- [x] 所有 match 操作符准确无误
- [x] 编码规则说明清晰

### 代码质量

- [x] 使用双引号而非反引号
- [x] 字符串拼接格式清晰易读
- [x] JSON 示例中的双引号已正确转义
- [x] 代码已格式化（go fmt）
- [x] 无语法错误

### 完整性

- [x] 新增的 key 类型已添加：URL, IPLocation
- [x] 原有 key 类型的 match 操作符已更新
- [x] 编码规则已更新为最新版本

---

## 实施顺序

1. **Task 1**: 更新 `options_arr` Description（核心任务）
2. **Task 2**: 验证修改（质量保证）

---

## 参考信息

### API 文档
- **URL**: https://cloud.tencent.com/document/api/627/97646
- **接口名称**: Waf CC V2 Upsert
- **字段名称**: OptionsArr
- **最后更新**: 2026-03-06 03:48:41

### Key 类型完整列表
1. URL
2. Method
3. Post
4. Referer
5. Cookie
6. User-Agent
7. CustomHeader
8. IPLocation
9. CaptchaRisk
10. CaptchaDeviceRisk
11. CaptchaScore

### 文件路径
- **修改文件**: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/waf/resource_tc_waf_cc.go`
- **修改位置**: Lines 82-103 (options_arr Schema definition)
