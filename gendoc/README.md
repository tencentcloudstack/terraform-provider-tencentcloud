# Terraform docs generator

## Why

经过观察，大部分友商的 Terraform Plugins 文档都是人肉编写的，它们或多或少都有这样的问题：

* 风格难以统一，比如：章节顺序、空行数量、命名风格在不同产品之间有明显差异
* 细节存在问题，比如：空格数量、缩进多少、中下划线出现不统一，甚至还夹带中文符号
* 内容存在问题，比如：参数必须或选填跟代码不一致，列表中漏写参数或属性

然而，最大的问题是写文档需要消耗大量的时间和精力去整理内容整理格式，最后还发现总有这样那样的问题，甚至有时还会文档更新不及时。
机器可以一如始终，完全无误差地完成各项有规律的重复性工作，而且不会出错，自动地生成文档也是 golang 所推动的标准做法。

## How

Terraform Plugins 文档，不管是 resource 还是 data_source，主要都分以下这几个主题：

* name
* description
* example usage
* argument reference
* attributes reference

### name

name 是 resource 及 data_source 的完整命名，它来源于 Provider 的 DataSourcesMap(data_source) 及 ResourcesMap(resource) 定义。
例如以下 DataSourcesMap 中的 tencentcloud_vpc 与 tencentcloud_mysql_instance 就是一个标准的 name(for resource or data_source)：

```go
DataSourcesMap: map[string]*schema.Resource{
    "tencentcloud_vpc": dataSourceTencentCloudVpc(),
    "tencentcloud_mysql_instance": dataSourceTencentCloudMysqlInstance(),
}
```

### description & example usage

description 包括一个用于表头的一句话描述，与一个用于正文的详细说明。
example usage 则是一个或几个使用示例。

description & example usage & usage description 需要在对应 resource 及 data_source 定义的文件中出现，它是符合 golang 标准文档注释的写法。例如：

    /*
    Use this data source to get information about a MySQL instance.
    \n
    ~> **NOTE:** The terminate operation of mysql does NOT take effect immediately，maybe takes for several hours.
    \n
    Example Usage
    \n
    Scenario1 title
    \n
    Description of the Scenario1
    \n
    ```hcl
    data "tencentcloud_mysql_instance" "database"{
      mysql_id = "my-test-database"
      result_output_file = "mytestpath"
    }
    ```
    \n
    Scenario2 title
    \n
    Description of the Scenario2
    \n
    ```hcl
    data "tencentcloud_mysql_instance" "database"{
      mysql_id = "my-test-database"
      result_output_file = "mytestpath"
    }
    ```
    */
    package tencentcloud

以上注释的格式要求如下：

    /*
    一句话描述
    \n
    在一句话描述基础上的补充描述，可以比较详细地说明各项内容，可以有多个段落。
    \n
    Example Usage
    \n
    Example Usage 是必须的，在 Example Usage 以下的内容都会填充到文档中。
    Example Usage 由一个到多个Scenario(场景)构成。
    每个Scenario 由 Scenario title 和 Scenario description 构成。
    \n
    Usage1 title
    \n
    Description of the Usage1
    \n
    Scenario title 是必须的。
    Scenario description 是可选的，可以根据情况填写。
    */
    package tencentcloud

符合以上要求的注释将会自动提取并填写到文档中的对应位置。

### argument reference & attributes reference

Terraform 用 schema.Schema 来描述 argument reference & attributes reference，每个 schema.Schema 都会有一个 Description 字段。
如果 Description 的内容不为空，那么这个 schema.Schema 将会被认为是需要写到文档里面的，如果 Optional 或 Required 设置了，它会被认为是一个参数，如果 Computed 为 true 则认为是一个属性。例如：

#### argument

```go
map[string]*schema.Schema{
    "instance_name": {
        Type:         schema.TypeString,
        Required:     true,
        ValidateFunc: validateStringLengthInRange(1, 100),
        Description:  "The name of a mysql instance.",
    },
}
```

#### attributes

```go
map[string]*schema.Schema{
    "mysql_id": {
        Type:     schema.TypeString,
        Computed: true,
        Description:  "Instance ID, such as cdb-c1nl9rpv. It is identical to the instance ID displayed in the database console page.",
    },
}
```

#### attributes list

属性中 Type 为 schema.TypeList 的 schema.Schema 也是支持的，它会被认为是一个列表，里面的子 schema.Schema 会依次列出填充到文档中。

```go
map[string]*schema.Schema{
    "instance_list": {
        Type:     schema.TypeList,
        Computed: true,
        Description: "A list of instances. Each element contains the following attributes:",
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "mysql_id": {
                    Type:     schema.TypeString,
                    Computed: true,
                    Description:  "Instance ID, such as cdb-c1nl9rpv. It is identical to the instance ID displayed in the database console page.",
                },
                "instance_name": {
                    Type:     schema.TypeString,
                    Computed: true,
                    Description:  "Name of mysql instance.",
                },
            }
        }
    }
}
```

## API Reference（资源/数据源关联的云 API）

resource、data source、action 与 ephemeral resource 的文档均包含一个可选的 "API Reference" 章节，列出该引用类型实际调用的云 API 及其官方文档链接：

- **resource**: Create / Read / Update / Delete 各阶段
- **data source**: Read 阶段
- **action**: Invoke 阶段
- **ephemeral resource**: Open 阶段

数据由两部分合成：

1. **静态分析**（`apiref.go`）：用 `go/ast` 从生命周期函数出发，追踪 service 方法调用链（含递归展开、goroutine 委托、包级辅助函数、局部 client 变量），收集 `UseXxxClient().Action(...)` 调用；client accessor 名通过 `connectivity/client.go` 的 SDK import 路径映射为目录产品名。
2. **国际站 API 目录**（`apiref_intl.go`）：`make doc` 每次运行实时抓取国际站（www.tencentcloud.com）——文档首页 `__NEXT_DATA__` 内嵌的全量目录树（129 产品 / 10000+ 接口页）与 `getCategoryByLang` 的产品 slug 表，提供 `(产品, action) -> 国际站文档 URL`。SDK 直接使用 vendor 目录下的 `github.com/tencentcloud/tencentcloud-sdk-go`，不依赖任何外部仓库。

**完全以国际站为准**：产品的名称与 ID 均取自国际站（slug 如 `cvm`/`teo`/`csip`），不使用国内站的任何信息。SDK 侧的产品名与国际站 slug 直接匹配；命名不一致的产品（如 SDK `redis` vs slug `crs`）通过动作集合的 Jaccard 相似度自动对齐。

**仅使用国际站链接**：渲染时只输出国际站（www.tencentcloud.com）链接。目录中未提供文档链接的接口回退到国际站搜索页：

```
https://www.tencentcloud.com/search/?wd=<ActionName>&type=document&lang=en
```

**渲染格式**：API Reference 只列出接口名称与对应的文档链接，不含接口描述，例如：

```
* Create: [CreateVpc](https://www.tencentcloud.com/document/api/215/15774)
```

生成的 API Reference 中不会出现任何 `cloud.tencent.com`（中国站）链接。抓取失败不影响构建，条目降级为搜索链接。

覆盖率：resource 95.7%（1302/1360）、data source 97.3%（923/949）、action 与 ephemeral 全覆盖；未覆盖项均为 S3 兼容 REST 接口（cos/ci）、已下线产品（dayu）、目录未收录的新版本 API，或**同名方法跨服务歧义**（分析器宁缺毋滥，跳过不确定的映射），可用下述人工段兜底。


人工兜底：当静态分析无法覆盖（如 COS/CI 的 S3 兼容 REST API、目录未收录的老版本产品）或需要修正时，可在资源描述文件 `tencentcloud/services/<product>/resource_tc_<name>.md` 中编写 "API" 段（位于正文与 "Example Usage" 之间），生成器会优先使用人工内容：

```
Use this resource to create dayu DDoS policy

API

* Create: CreateDDoSPolicy - https://www.tencentcloud.com/document/product/1022/37825
* Read: DescribeDDoSPolicy - https://www.tencentcloud.com/document/product/1022/37826

Example Usage
...
```

目录快照：在线拉取成功时，`make doc` 会把最新目录原子写入快照 `gendoc/data/api-intl-snapshot.json` 并随提交入库；在线拉取失败（如离线环境）时自动回退到该快照，生成结果与在线时保持一致，仅快照日期可能滞后。两种来源都不可用时，全部条目降级为搜索链接。

## 文档索引更新

文档索引文件，即 website/tencentcloud.erb 的更新数据来源于 provider.go 的文件注释。

完成了新的 Data Sources 或 Resources 后，需要更新 provider.go 的文件注释，格式可参考已有的 Data Sources 或 Resources。

### Data Source

在注释中找到对应产品的 `Data Source`，在它的下面填写新的 Data Source 名称。如果是新的产品，则先添加新的产品类，例如 `CVM`，产品名称的简写如果容易使人迷惑，则先写产品名称详写，再写缩写，例如 `Direct Connect(DC)`。

例如：

```go
CVM
  Data Source
    tencentcloud_image
```

如果是通用的 Data Source，则添加到 `Provider Data Sources` 这个类下面。

### Resource

在注释中找到对应产品的 `Resource`，在它的下面填写新的 Resource 名称。如果是新的产品，则先添加新的产品类，例如 `CVM`，产品名称的简写如果容易使人迷惑，则先写产品名称详写，再写缩写，例如 `Direct Connect(DC)`。

例如：

```go
CVM
  Data Source
    tencentcloud_image
    ...

  Resource
    tencentcloud_instance
```
