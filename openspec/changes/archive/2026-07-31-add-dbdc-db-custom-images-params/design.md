## Context

数据源 `tencentcloud_dbdc_db_custom_images` 用于查询 DB Custom 可用的操作系统镜像列表，底层调用云 API `DescribeDBCustomImages`（dbdc v20201029）。当前数据源未暴露任何过滤入参，也未输出镜像的 `OsType`（操作系统类型）字段，仅返回 `image_id`、`os_name`、`image_type`、`architecture`。

云 API `DescribeDBCustomImages` 已支持：
- 入参 `Filters`（`[]*Filter`），其中 `Filter` 包含 `Name`（string）和 `Values`（[]*string）两个字段。过滤键可取 `image-id`、`os-type`、`image-type`、`architecture`。
- 出参 `ImageSet`（`[]*DBCustomImage`），其中 `DBCustomImage` 包含 `OsType`（string）字段。

本次变更新增入参 `Name`/`Values` 用于构造单个 `Filter`，并新增出参 `OsType`。

## Goals / Non-Goals

**Goals:**
- 在数据源 schema 中新增可选入参 `Name`（string）和 `Values`（list），供用户传入过滤条件。
- 在 `image_set` 出参 schema 中新增 `os_type`（string）字段。
- 在 Read 方法与服务层方法中，将 `Name`/`Values` 组装为 `[]*Filter` 并传入 `DescribeDBCustomImages` 请求。
- 在 Read 方法中将返回的 `OsType` 写入 state。
- 补充对应的单元测试。

**Non-Goals:**
- 不暴露 `Offset`/`Limit` 分页参数（由服务层内部自动分页，符合既有模式）。
- 不修改其它 dbdc 资源或数据源。
- 不改变现有的 `image_id`/`os_name`/`image_type`/`architecture` 字段行为。

## Decisions

### 决策一：`Name`/`Values` 作为顶层入参而非嵌套 Filter 块
云 API 的 `Filters` 是 `[]*Filter`，但本次仅新增单组 `Name`/`Values`。为保持 Terraform 用户配置简洁，将 `Name` 与 `Values` 作为数据源顶层 Optional 入参，在 Read 方法中组装为单个 `Filter` 传入请求。

**备选方案**：将 `filters` 定义为 `TypeList` 嵌套块（含 `name`/`values`）。该方案更贴近云 API 结构，但本次需求只新增一个过滤组，嵌套块会增加配置复杂度且不符合"单一新增参数"的约束。

### 决策二：`OsType` 作为 `image_set` 元素的 Computed 字段
`OsType` 属于云 API `DBCustomImage` 的出参字段，遵循现有 `image_set` 列表展开结构，将 `os_type` 加入 `image_set` 的 `Elem.Schema` 中，与其他出参字段平铺，符合既有代码风格。

### 决策三：服务层 `DescribeDBCustomImagesByFilter` 通过 paramMap 传递 Filters
现有 `DescribeDBCustomImagesByFilter(ctx, paramMap)` 已通过 `paramMap` 传参。在 Read 方法中将 `name`、`values`（以 `helper.StringPtr`/`[]*string` 处理）写入 `paramMap`，服务层读取后构造 `[]*Filter` 赋值给 `request.Filters`。这样改动最小且保持服务层签名不变。

## Risks / Trade-offs

- **[单组 Filter 限制]** → 本次仅支持一个过滤组（`Name`+`Values`）。若未来需要多组过滤条件，可再扩展为嵌套块。当前需求仅要求新增该参数，符合最小变更原则。
- **[向后兼容]** → `Name`/`Values` 为 Optional 且默认不传 Filters 时云 API 返回全部镜像，行为与现状一致；`os_type` 为新增 Computed 字段，不影响现有 state。无破坏性变更。
- **[空值处理]** → 服务层仅在 `Name` 非空时才构造 `Filter`，避免传入空 Filter 导致 API 报错。Read 方法设置 `os_type` 前判断 `image.OsType != nil`。
