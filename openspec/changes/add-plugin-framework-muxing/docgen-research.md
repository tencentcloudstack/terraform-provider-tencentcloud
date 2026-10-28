# Documentation Generator Research

## 现状

`make doc` 调用 `gendoc/` 下的自研 Go 程序：

- 入口：`gendoc/main.go`
- 依赖：直接 import `terraform-plugin-sdk/v2/helper/schema`
- 数据源：遍历 `cloud.Provider().ResourcesMap` / `DataSourcesMap`（SDKv2）
- 模板：`gendoc/template.go` 与项目内每个资源同目录的 `*.md` 业务模板（如
  `tencentcloud/services/<svc>/resource_tc_<svc>_<name>.md`）
- 产物：`website/docs/r/<name>.html.markdown` / `website/docs/d/<name>.html.markdown`
- 索引：`gendoc/index.go` 维护服务侧栏分组

## 结论

`gendoc` **无法识别 framework schema**。它通过反射读取 SDKv2 schema 的字段
和元数据，对 `terraform-plugin-framework` 的 Resource/DataSource 不知情。

## 备选方案

### A. 切换到 `tfplugindocs`（HashiCorp 官方）

- 优点：原生支持 SDKv2 + framework；社区主流方案；维护成本低。
- 代价（高）：
  - 700+ 现有资源的 `*.md` 业务模板需要全部迁移到 `templates/` 目录的
    `tfplugindocs` 风格。
  - 服务侧栏 `gendoc/index.go` 的分组逻辑需要重写为
    `templates/index.md.tmpl`。
  - acceptance test、CI、`make doc-faster` 等依赖 `gendoc` 二进制的环节
    需要全部修改。

### B. 在 `gendoc` 旁补一个 `gendoc-fw` 子流程（混合方案）

- 优点：增量、零回归。
- 缺点：长期看需要维护两套生成器，最终也需要靠 A 收敛。

### C. 暂不切换，框架侧资源手写文档（本次 change 选用）

- 框架资源数量目前为 0（本变更新增第一个示范数据源）；手写一份 markdown
  的成本远小于改造 `gendoc`。
- `make doc` 命令保持原有产物不变，零回归。
- 待 framework 资源数量超过约 5 个、或第一次出现复杂嵌套字段时，再以
  独立 change 升级到 A。

## 本次行动

1. 为示范数据源 `tencentcloud_provider_runtime` 手写
   `website/docs/d/provider_runtime.html.markdown`。
2. 不动 `gendoc/`、不动 `make doc`。
3. 在 design.md 的"Open Questions"段落记录"待 framework 资源 ≥ 5 个再
   切 tfplugindocs"决议。
