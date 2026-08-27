// Copyright (c) 2017-2025 Tencent. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20170320

import (
	tcerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/json"
)

type Account struct {
	// <p>账号名，可输入1 - 32个字符。</p>
	User *string `json:"User,omitnil,omitempty" name:"User"`

	// <p>账号的主机。</p><p>IP 形式，支持填入%。</p>
	Host *string `json:"Host,omitnil,omitempty" name:"Host"`
}

type AccountInfo struct {
	// <p>账号备注信息</p>
	Notes *string `json:"Notes,omitnil,omitempty" name:"Notes"`

	// <p>账号的域名</p>
	Host *string `json:"Host,omitnil,omitempty" name:"Host"`

	// <p>账号的名称</p>
	User *string `json:"User,omitnil,omitempty" name:"User"`

	// <p>账号信息修改时间</p>
	ModifyTime *string `json:"ModifyTime,omitnil,omitempty" name:"ModifyTime"`

	// <p>修改密码的时间</p>
	ModifyPasswordTime *string `json:"ModifyPasswordTime,omitnil,omitempty" name:"ModifyPasswordTime"`

	// <p>该值已废弃</p>
	//
	// Deprecated: CreateTime is deprecated.
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`

	// <p>用户最大可用实例连接数</p>
	MaxUserConnections *int64 `json:"MaxUserConnections,omitnil,omitempty" name:"MaxUserConnections"`

	// <p>用户账号是否开启了密码轮转</p>
	OpenCam *bool `json:"OpenCam,omitnil,omitempty" name:"OpenCam"`
}

// Predefined struct for user
type AddTimeWindowRequestParams struct {
	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 星期一的可维护时间段，其中每一个时间段的格式形如：10:00-12:00；起始时间按半个小时对齐；最短半个小时，最长三个小时；可设置多个时间段。 一周中应至少设置一天的时间窗。下同。
	Monday []*string `json:"Monday,omitnil,omitempty" name:"Monday"`

	// 星期二的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Tuesday []*string `json:"Tuesday,omitnil,omitempty" name:"Tuesday"`

	// 星期三的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Wednesday []*string `json:"Wednesday,omitnil,omitempty" name:"Wednesday"`

	// 星期四的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Thursday []*string `json:"Thursday,omitnil,omitempty" name:"Thursday"`

	// 星期五的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Friday []*string `json:"Friday,omitnil,omitempty" name:"Friday"`

	// 星期六的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Saturday []*string `json:"Saturday,omitnil,omitempty" name:"Saturday"`

	// 星期日的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Sunday []*string `json:"Sunday,omitnil,omitempty" name:"Sunday"`

	// 最大延迟阈值（秒），仅对主实例和灾备实例有效。默认值：10，取值范围：1-10的整数。
	MaxDelayTime *uint64 `json:"MaxDelayTime,omitnil,omitempty" name:"MaxDelayTime"`
}

type AddTimeWindowRequest struct {
	*tchttp.BaseRequest

	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 星期一的可维护时间段，其中每一个时间段的格式形如：10:00-12:00；起始时间按半个小时对齐；最短半个小时，最长三个小时；可设置多个时间段。 一周中应至少设置一天的时间窗。下同。
	Monday []*string `json:"Monday,omitnil,omitempty" name:"Monday"`

	// 星期二的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Tuesday []*string `json:"Tuesday,omitnil,omitempty" name:"Tuesday"`

	// 星期三的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Wednesday []*string `json:"Wednesday,omitnil,omitempty" name:"Wednesday"`

	// 星期四的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Thursday []*string `json:"Thursday,omitnil,omitempty" name:"Thursday"`

	// 星期五的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Friday []*string `json:"Friday,omitnil,omitempty" name:"Friday"`

	// 星期六的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Saturday []*string `json:"Saturday,omitnil,omitempty" name:"Saturday"`

	// 星期日的可维护时间窗口。 一周中应至少设置一天的时间窗。
	Sunday []*string `json:"Sunday,omitnil,omitempty" name:"Sunday"`

	// 最大延迟阈值（秒），仅对主实例和灾备实例有效。默认值：10，取值范围：1-10的整数。
	MaxDelayTime *uint64 `json:"MaxDelayTime,omitnil,omitempty" name:"MaxDelayTime"`
}

func (r *AddTimeWindowRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AddTimeWindowRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "Monday")
	delete(f, "Tuesday")
	delete(f, "Wednesday")
	delete(f, "Thursday")
	delete(f, "Friday")
	delete(f, "Saturday")
	delete(f, "Sunday")
	delete(f, "MaxDelayTime")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AddTimeWindowRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AddTimeWindowResponseParams struct {
	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type AddTimeWindowResponse struct {
	*tchttp.BaseResponse
	Response *AddTimeWindowResponseParams `json:"Response"`
}

func (r *AddTimeWindowResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AddTimeWindowResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type AddressInfo struct {
	// 地址的资源id标识。
	ResourceId *string `json:"ResourceId,omitnil,omitempty" name:"ResourceId"`

	// 地址所在的vpc。
	UniqVpcId *string `json:"UniqVpcId,omitnil,omitempty" name:"UniqVpcId"`

	// 地址所在的子网。
	UniqSubnetId *string `json:"UniqSubnetId,omitnil,omitempty" name:"UniqSubnetId"`

	// 地址的vip。
	Vip *string `json:"Vip,omitnil,omitempty" name:"Vip"`

	// 地址的端口。
	VPort *int64 `json:"VPort,omitnil,omitempty" name:"VPort"`

	// 外网地址域名。
	WanDomain *string `json:"WanDomain,omitnil,omitempty" name:"WanDomain"`

	// 外网地址端口。
	WanPort *int64 `json:"WanPort,omitnil,omitempty" name:"WanPort"`
}

// Predefined struct for user
type AdjustCdbProxyAddressRequestParams struct {
	// <p>代理组 ID。可通过 <a href="https://cloud.tencent.com/document/api/236/90585">DescribeCdbProxyInfo</a> 接口获取。</p>
	ProxyGroupId *string `json:"ProxyGroupId,omitnil,omitempty" name:"ProxyGroupId"`

	// <p>权重分配模式，<br>系统自动分配：&quot;system&quot;， 自定义：&quot;custom&quot;</p>
	WeightMode *string `json:"WeightMode,omitnil,omitempty" name:"WeightMode"`

	// <p>是否开启延迟剔除，取值：&quot;true&quot; | &quot;false&quot;</p>
	IsKickOut *bool `json:"IsKickOut,omitnil,omitempty" name:"IsKickOut"`

	// <p>最小保留数量，最小取值：0。<br>说明：当 IsKickOut 为 true 时才有效。</p>
	MinCount *uint64 `json:"MinCount,omitnil,omitempty" name:"MinCount"`

	// <p>延迟剔除阈值，最小取值：1，取值范围：[1,10000]，整数。</p>
	MaxDelay *uint64 `json:"MaxDelay,omitnil,omitempty" name:"MaxDelay"`

	// <p>是否开启故障转移，取值：&quot;true&quot; | &quot;false&quot;</p>
	FailOver *bool `json:"FailOver,omitnil,omitempty" name:"FailOver"`

	// <p>是否自动添加RO，取值：&quot;true&quot; | &quot;false&quot;</p>
	AutoAddRo *bool `json:"AutoAddRo,omitnil,omitempty" name:"AutoAddRo"`

	// <p>是否是只读，取值：&quot;true&quot; | &quot;false&quot;</p>
	ReadOnly *bool `json:"ReadOnly,omitnil,omitempty" name:"ReadOnly"`

	// <p>代理组地址 ID。可通过 <a href="https://cloud.tencent.com/document/api/236/90585">DescribeCdbProxyInfo</a> 接口获取。</p>
	ProxyAddressId *string `json:"ProxyAddressId,omitnil,omitempty" name:"ProxyAddressId"`

	// <p>是否开启事务分离，取值：&quot;true&quot; | &quot;false&quot;，默认值 false。</p>
	TransSplit *bool `json:"TransSplit,omitnil,omitempty" name:"TransSplit"`

	// <p>是否开启连接池。默认关闭。<br>注意：如需使用数据库代理连接池能力，MySQL 8.0 主实例的内核小版本要大于等于 MySQL 8.0 20230630。</p>
	ConnectionPool *bool `json:"ConnectionPool,omitnil,omitempty" name:"ConnectionPool"`

	// <p>读写权重分配。如果 WeightMode 传的是 system ，则传入的权重不生效，由系统分配默认权重。</p>
	ProxyAllocation []*ProxyAllocation `json:"ProxyAllocation,omitnil,omitempty" name:"ProxyAllocation"`

	// <p>是否开启自适应负载均衡。默认关闭。</p>
	AutoLoadBalance *bool `json:"AutoLoadBalance,omitnil,omitempty" name:"AutoLoadBalance"`

	// <p>访问模式。</p><p>枚举值：</p><ul><li>nearby： 就近访问</li><li>balance： 均衡分配</li><li>direct_nearby： 纯网络转发就近访问</li><li>direct_balance： 纯网络转发均衡分配</li></ul>
	AccessMode *string `json:"AccessMode,omitnil,omitempty" name:"AccessMode"`

	// <p>是否将libra节点当作普通RO节点</p>
	ApNodeAsRoNode *bool `json:"ApNodeAsRoNode,omitnil,omitempty" name:"ApNodeAsRoNode"`

	// <p>libra节点故障，是否转发给其他节点</p>
	ApQueryToOtherNode *bool `json:"ApQueryToOtherNode,omitnil,omitempty" name:"ApQueryToOtherNode"`
}

type AdjustCdbProxyAddressRequest struct {
	*tchttp.BaseRequest

	// <p>代理组 ID。可通过 <a href="https://cloud.tencent.com/document/api/236/90585">DescribeCdbProxyInfo</a> 接口获取。</p>
	ProxyGroupId *string `json:"ProxyGroupId,omitnil,omitempty" name:"ProxyGroupId"`

	// <p>权重分配模式，<br>系统自动分配：&quot;system&quot;， 自定义：&quot;custom&quot;</p>
	WeightMode *string `json:"WeightMode,omitnil,omitempty" name:"WeightMode"`

	// <p>是否开启延迟剔除，取值：&quot;true&quot; | &quot;false&quot;</p>
	IsKickOut *bool `json:"IsKickOut,omitnil,omitempty" name:"IsKickOut"`

	// <p>最小保留数量，最小取值：0。<br>说明：当 IsKickOut 为 true 时才有效。</p>
	MinCount *uint64 `json:"MinCount,omitnil,omitempty" name:"MinCount"`

	// <p>延迟剔除阈值，最小取值：1，取值范围：[1,10000]，整数。</p>
	MaxDelay *uint64 `json:"MaxDelay,omitnil,omitempty" name:"MaxDelay"`

	// <p>是否开启故障转移，取值：&quot;true&quot; | &quot;false&quot;</p>
	FailOver *bool `json:"FailOver,omitnil,omitempty" name:"FailOver"`

	// <p>是否自动添加RO，取值：&quot;true&quot; | &quot;false&quot;</p>
	AutoAddRo *bool `json:"AutoAddRo,omitnil,omitempty" name:"AutoAddRo"`

	// <p>是否是只读，取值：&quot;true&quot; | &quot;false&quot;</p>
	ReadOnly *bool `json:"ReadOnly,omitnil,omitempty" name:"ReadOnly"`

	// <p>代理组地址 ID。可通过 <a href="https://cloud.tencent.com/document/api/236/90585">DescribeCdbProxyInfo</a> 接口获取。</p>
	ProxyAddressId *string `json:"ProxyAddressId,omitnil,omitempty" name:"ProxyAddressId"`

	// <p>是否开启事务分离，取值：&quot;true&quot; | &quot;false&quot;，默认值 false。</p>
	TransSplit *bool `json:"TransSplit,omitnil,omitempty" name:"TransSplit"`

	// <p>是否开启连接池。默认关闭。<br>注意：如需使用数据库代理连接池能力，MySQL 8.0 主实例的内核小版本要大于等于 MySQL 8.0 20230630。</p>
	ConnectionPool *bool `json:"ConnectionPool,omitnil,omitempty" name:"ConnectionPool"`

	// <p>读写权重分配。如果 WeightMode 传的是 system ，则传入的权重不生效，由系统分配默认权重。</p>
	ProxyAllocation []*ProxyAllocation `json:"ProxyAllocation,omitnil,omitempty" name:"ProxyAllocation"`

	// <p>是否开启自适应负载均衡。默认关闭。</p>
	AutoLoadBalance *bool `json:"AutoLoadBalance,omitnil,omitempty" name:"AutoLoadBalance"`

	// <p>访问模式。</p><p>枚举值：</p><ul><li>nearby： 就近访问</li><li>balance： 均衡分配</li><li>direct_nearby： 纯网络转发就近访问</li><li>direct_balance： 纯网络转发均衡分配</li></ul>
	AccessMode *string `json:"AccessMode,omitnil,omitempty" name:"AccessMode"`

	// <p>是否将libra节点当作普通RO节点</p>
	ApNodeAsRoNode *bool `json:"ApNodeAsRoNode,omitnil,omitempty" name:"ApNodeAsRoNode"`

	// <p>libra节点故障，是否转发给其他节点</p>
	ApQueryToOtherNode *bool `json:"ApQueryToOtherNode,omitnil,omitempty" name:"ApQueryToOtherNode"`
}

func (r *AdjustCdbProxyAddressRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AdjustCdbProxyAddressRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProxyGroupId")
	delete(f, "WeightMode")
	delete(f, "IsKickOut")
	delete(f, "MinCount")
	delete(f, "MaxDelay")
	delete(f, "FailOver")
	delete(f, "AutoAddRo")
	delete(f, "ReadOnly")
	delete(f, "ProxyAddressId")
	delete(f, "TransSplit")
	delete(f, "ConnectionPool")
	delete(f, "ProxyAllocation")
	delete(f, "AutoLoadBalance")
	delete(f, "AccessMode")
	delete(f, "ApNodeAsRoNode")
	delete(f, "ApQueryToOtherNode")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AdjustCdbProxyAddressRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AdjustCdbProxyAddressResponseParams struct {
	// <p>异步任务ID</p>
	AsyncRequestId *string `json:"AsyncRequestId,omitnil,omitempty" name:"AsyncRequestId"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type AdjustCdbProxyAddressResponse struct {
	*tchttp.BaseResponse
	Response *AdjustCdbProxyAddressResponseParams `json:"Response"`
}

func (r *AdjustCdbProxyAddressResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AdjustCdbProxyAddressResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AdjustCdbProxyRequestParams struct {
	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 代理组 ID。可通过 [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) 接口获取。
	ProxyGroupId *string `json:"ProxyGroupId,omitnil,omitempty" name:"ProxyGroupId"`

	// 节点规格配置
	// 备注：数据库代理支持的节点规格为：2C4000MB、4C8000MB、8C16000MB。
	// 示例中参数说明：
	// NodeCount：节点个数
	// Region：节点地域
	// Zone：节点可用区
	// Cpu：单个代理节点核数（单位：核）
	// Mem：单个代理节点内存数（单位：MB）
	ProxyNodeCustom []*ProxyNodeCustom `json:"ProxyNodeCustom,omitnil,omitempty" name:"ProxyNodeCustom"`

	// 重新负载均衡：auto(自动),manual(手动)
	ReloadBalance *string `json:"ReloadBalance,omitnil,omitempty" name:"ReloadBalance"`

	// 升级切换时间：nowTime(升级完成时),timeWindow(维护时间内)
	UpgradeTime *string `json:"UpgradeTime,omitnil,omitempty" name:"UpgradeTime"`
}

type AdjustCdbProxyRequest struct {
	*tchttp.BaseRequest

	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 代理组 ID。可通过 [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) 接口获取。
	ProxyGroupId *string `json:"ProxyGroupId,omitnil,omitempty" name:"ProxyGroupId"`

	// 节点规格配置
	// 备注：数据库代理支持的节点规格为：2C4000MB、4C8000MB、8C16000MB。
	// 示例中参数说明：
	// NodeCount：节点个数
	// Region：节点地域
	// Zone：节点可用区
	// Cpu：单个代理节点核数（单位：核）
	// Mem：单个代理节点内存数（单位：MB）
	ProxyNodeCustom []*ProxyNodeCustom `json:"ProxyNodeCustom,omitnil,omitempty" name:"ProxyNodeCustom"`

	// 重新负载均衡：auto(自动),manual(手动)
	ReloadBalance *string `json:"ReloadBalance,omitnil,omitempty" name:"ReloadBalance"`

	// 升级切换时间：nowTime(升级完成时),timeWindow(维护时间内)
	UpgradeTime *string `json:"UpgradeTime,omitnil,omitempty" name:"UpgradeTime"`
}

func (r *AdjustCdbProxyRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AdjustCdbProxyRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "ProxyGroupId")
	delete(f, "ProxyNodeCustom")
	delete(f, "ReloadBalance")
	delete(f, "UpgradeTime")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AdjustCdbProxyRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AdjustCdbProxyResponseParams struct {
	// 异步任务ID
	AsyncRequestId *string `json:"AsyncRequestId,omitnil,omitempty" name:"AsyncRequestId"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type AdjustCdbProxyResponse struct {
	*tchttp.BaseResponse
	Response *AdjustCdbProxyResponseParams `json:"Response"`
}

func (r *AdjustCdbProxyResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AdjustCdbProxyResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type AggregationCondition struct {
	// 聚合字段。目前仅支持host-源IP、user-用户名、dbName-数据库名、sqlType-sql类型。
	AggregationField *string `json:"AggregationField,omitnil,omitempty" name:"AggregationField"`

	// 偏移量。
	Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`

	// 该聚合字段下要返回聚合桶的数量，最大100。
	Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}

type AnalysisNodeInfo struct {
	// 节点ID
	NodeId *string `json:"NodeId,omitnil,omitempty" name:"NodeId"`

	// 节点状态
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// 数据加载状态
	DataStatus *string `json:"DataStatus,omitnil,omitempty" name:"DataStatus"`

	// cpu核数，单位：核
	Cpu *uint64 `json:"Cpu,omitnil,omitempty" name:"Cpu"`

	// 内存大小，单位: MB
	Memory *uint64 `json:"Memory,omitnil,omitempty" name:"Memory"`

	// 磁盘大小，单位：GB
	Storage *uint64 `json:"Storage,omitnil,omitempty" name:"Storage"`

	// 节点所在可用区
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`

	// 数据同步错误信息
	Message *string `json:"Message,omitnil,omitempty" name:"Message"`
}

// Predefined struct for user
type AnalyzeAuditLogsRequestParams struct {
	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 要分析的日志开始时间，格式为："2023-02-16 00:00:20"。
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// 要分析的日志结束时间，格式为："2023-02-16 00:10:20"。
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// 聚合维度的排序条件。
	AggregationConditions []*AggregationCondition `json:"AggregationConditions,omitnil,omitempty" name:"AggregationConditions"`

	// 已废弃。
	//
	// Deprecated: AuditLogFilter is deprecated.
	AuditLogFilter *AuditLogFilter `json:"AuditLogFilter,omitnil,omitempty" name:"AuditLogFilter"`

	// 该过滤条件下的审计日志结果集作为分析日志。
	LogFilter []*InstanceAuditLogFilters `json:"LogFilter,omitnil,omitempty" name:"LogFilter"`
}

type AnalyzeAuditLogsRequest struct {
	*tchttp.BaseRequest

	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 要分析的日志开始时间，格式为："2023-02-16 00:00:20"。
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// 要分析的日志结束时间，格式为："2023-02-16 00:10:20"。
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// 聚合维度的排序条件。
	AggregationConditions []*AggregationCondition `json:"AggregationConditions,omitnil,omitempty" name:"AggregationConditions"`

	// 已废弃。
	AuditLogFilter *AuditLogFilter `json:"AuditLogFilter,omitnil,omitempty" name:"AuditLogFilter"`

	// 该过滤条件下的审计日志结果集作为分析日志。
	LogFilter []*InstanceAuditLogFilters `json:"LogFilter,omitnil,omitempty" name:"LogFilter"`
}

func (r *AnalyzeAuditLogsRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AnalyzeAuditLogsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "StartTime")
	delete(f, "EndTime")
	delete(f, "AggregationConditions")
	delete(f, "AuditLogFilter")
	delete(f, "LogFilter")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AnalyzeAuditLogsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AnalyzeAuditLogsResponseParams struct {
	// 返回的聚合桶信息集
	Items []*AuditLogAggregationResult `json:"Items,omitnil,omitempty" name:"Items"`

	// 扫描的日志条数
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type AnalyzeAuditLogsResponse struct {
	*tchttp.BaseResponse
	Response *AnalyzeAuditLogsResponseParams `json:"Response"`
}

func (r *AnalyzeAuditLogsResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AnalyzeAuditLogsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AssociateSecurityGroupsRequestParams struct {
	// 安全组 ID。可通过 [DescribeDBSecurityGroups](https://cloud.tencent.com/document/api/236/15854) 接口获取。
	SecurityGroupId *string `json:"SecurityGroupId,omitnil,omitempty" name:"SecurityGroupId"`

	// 实例 ID 列表，一个或者多个实例 ID 组成的数组。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// 当传入只读实例ID时，默认操作的是对应只读组的安全组。如果需要操作只读实例ID的安全组， 需要将该入参置为True
	ForReadonlyInstance *bool `json:"ForReadonlyInstance,omitnil,omitempty" name:"ForReadonlyInstance"`
}

type AssociateSecurityGroupsRequest struct {
	*tchttp.BaseRequest

	// 安全组 ID。可通过 [DescribeDBSecurityGroups](https://cloud.tencent.com/document/api/236/15854) 接口获取。
	SecurityGroupId *string `json:"SecurityGroupId,omitnil,omitempty" name:"SecurityGroupId"`

	// 实例 ID 列表，一个或者多个实例 ID 组成的数组。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceIds []*string `json:"InstanceIds,omitnil,omitempty" name:"InstanceIds"`

	// 当传入只读实例ID时，默认操作的是对应只读组的安全组。如果需要操作只读实例ID的安全组， 需要将该入参置为True
	ForReadonlyInstance *bool `json:"ForReadonlyInstance,omitnil,omitempty" name:"ForReadonlyInstance"`
}

func (r *AssociateSecurityGroupsRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AssociateSecurityGroupsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "SecurityGroupId")
	delete(f, "InstanceIds")
	delete(f, "ForReadonlyInstance")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AssociateSecurityGroupsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AssociateSecurityGroupsResponseParams struct {
	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type AssociateSecurityGroupsResponse struct {
	*tchttp.BaseResponse
	Response *AssociateSecurityGroupsResponseParams `json:"Response"`
}

func (r *AssociateSecurityGroupsResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AssociateSecurityGroupsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type AuditFilter struct {
	// 过滤条件参数名称。目前支持：
	// SrcIp – 客户端 IP；
	// User – 数据库账户；
	// DB – 数据库名称；
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// 过滤条件匹配类型。目前支持：
	// INC – 包含；
	// EXC – 不包含；
	// EQ – 等于；
	// NEQ – 不等于；
	Compare *string `json:"Compare,omitnil,omitempty" name:"Compare"`

	// 过滤条件匹配值。
	Value *string `json:"Value,omitnil,omitempty" name:"Value"`
}

type AuditInstanceFilters struct {
	// 过滤条件名。支持InstanceId-实例ID，InstanceName-实例名称，ProjectId-项目ID，TagKey-标签键，Tag-标签（以竖线分割，例：Tagkey|Tagvalue）。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// true表示精确查找，false表示模糊匹配。
	ExactMatch *bool `json:"ExactMatch,omitnil,omitempty" name:"ExactMatch"`

	// 筛选值
	Values []*string `json:"Values,omitnil,omitempty" name:"Values"`
}

type AuditInstanceInfo struct {
	// 项目ID
	ProjectId *int64 `json:"ProjectId,omitnil,omitempty" name:"ProjectId"`

	// 标签信息
	TagList []*TagInfoUnit `json:"TagList,omitnil,omitempty" name:"TagList"`

	// 数据库内核类型
	DbType *string `json:"DbType,omitnil,omitempty" name:"DbType"`

	// 数据库内核版本
	DbVersion *string `json:"DbVersion,omitnil,omitempty" name:"DbVersion"`
}

type AuditLog struct {
	// 影响行数。
	AffectRows *int64 `json:"AffectRows,omitnil,omitempty" name:"AffectRows"`

	// 错误码。
	ErrCode *int64 `json:"ErrCode,omitnil,omitempty" name:"ErrCode"`

	// SQL 类型。
	SqlType *string `json:"SqlType,omitnil,omitempty" name:"SqlType"`

	// 审计策略名称，逐步下线。
	PolicyName *string `json:"PolicyName,omitnil,omitempty" name:"PolicyName"`

	// 数据库名称。
	DBName *string `json:"DBName,omitnil,omitempty" name:"DBName"`

	// SQL 语句。
	Sql *string `json:"Sql,omitnil,omitempty" name:"Sql"`

	// 客户端地址。
	Host *string `json:"Host,omitnil,omitempty" name:"Host"`

	// 用户名。
	User *string `json:"User,omitnil,omitempty" name:"User"`

	// 执行时间，微秒。
	ExecTime *int64 `json:"ExecTime,omitnil,omitempty" name:"ExecTime"`

	// 时间。
	Timestamp *string `json:"Timestamp,omitnil,omitempty" name:"Timestamp"`

	// 返回行数。
	SentRows *int64 `json:"SentRows,omitnil,omitempty" name:"SentRows"`

	// 线程ID。
	ThreadId *int64 `json:"ThreadId,omitnil,omitempty" name:"ThreadId"`

	// 扫描行数。
	CheckRows *int64 `json:"CheckRows,omitnil,omitempty" name:"CheckRows"`

	// cpu执行时间，微秒。
	CpuTime *float64 `json:"CpuTime,omitnil,omitempty" name:"CpuTime"`

	// IO等待时间，微秒。
	IoWaitTime *uint64 `json:"IoWaitTime,omitnil,omitempty" name:"IoWaitTime"`

	// 锁等待时间，微秒。
	LockWaitTime *uint64 `json:"LockWaitTime,omitnil,omitempty" name:"LockWaitTime"`

	// 开始时间，与timestamp构成一个精确到纳秒的时间。
	NsTime *uint64 `json:"NsTime,omitnil,omitempty" name:"NsTime"`

	// 事物持续时间，微秒。
	TrxLivingTime *uint64 `json:"TrxLivingTime,omitnil,omitempty" name:"TrxLivingTime"`

	// 日志命中规则模板的基本信息
	TemplateInfo []*LogRuleTemplateInfo `json:"TemplateInfo,omitnil,omitempty" name:"TemplateInfo"`

	//  事务ID
	TrxId *int64 `json:"TrxId,omitnil,omitempty" name:"TrxId"`

	// 端口
	// 注意：此字段可能返回 null，表示取不到有效值。
	ClientPort *int64 `json:"ClientPort,omitnil,omitempty" name:"ClientPort"`
}

type AuditLogAggregationResult struct {
	// 聚合维度
	AggregationField *string `json:"AggregationField,omitnil,omitempty" name:"AggregationField"`

	// 聚合桶的结果集
	Buckets []*Bucket `json:"Buckets,omitnil,omitempty" name:"Buckets"`
}

type AuditLogFile struct {
	// 审计日志文件名称
	FileName *string `json:"FileName,omitnil,omitempty" name:"FileName"`

	// 审计日志文件创建时间。格式为 : "2019-03-20 17:09:13"。
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`

	// 文件状态值。可能返回的值为：
	// "creating" - 生成中;
	// "failed" - 创建失败;
	// "success" - 已生成;
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// 文件大小，单位为 KB。
	FileSize *int64 `json:"FileSize,omitnil,omitempty" name:"FileSize"`

	// 审计日志下载地址。
	DownloadUrl *string `json:"DownloadUrl,omitnil,omitempty" name:"DownloadUrl"`

	// 错误信息。
	ErrMsg *string `json:"ErrMsg,omitnil,omitempty" name:"ErrMsg"`
}

type AuditLogFilter struct {
	// 客户端地址。
	Host []*string `json:"Host,omitnil,omitempty" name:"Host"`

	// 用户名。
	User []*string `json:"User,omitnil,omitempty" name:"User"`

	// 数据库名称。
	DBName []*string `json:"DBName,omitnil,omitempty" name:"DBName"`

	// 表名称。
	TableName []*string `json:"TableName,omitnil,omitempty" name:"TableName"`

	// 审计策略名称。
	PolicyName []*string `json:"PolicyName,omitnil,omitempty" name:"PolicyName"`

	// SQL 语句。支持模糊匹配。
	Sql *string `json:"Sql,omitnil,omitempty" name:"Sql"`

	// SQL 类型。目前支持："SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "ALTER", "SET", "REPLACE", "EXECUTE"。
	SqlType *string `json:"SqlType,omitnil,omitempty" name:"SqlType"`

	// 执行时间。单位为：ms。表示筛选执行时间大于该值的审计日志。
	ExecTime *int64 `json:"ExecTime,omitnil,omitempty" name:"ExecTime"`

	// 影响行数。表示筛选影响行数大于该值的审计日志。
	AffectRows *int64 `json:"AffectRows,omitnil,omitempty" name:"AffectRows"`

	// SQL 类型。支持多个类型同时查询。目前支持："SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "ALTER", "SET", "REPLACE", "EXECUTE"。
	SqlTypes []*string `json:"SqlTypes,omitnil,omitempty" name:"SqlTypes"`

	// SQL 语句。支持传递多个sql语句。
	Sqls []*string `json:"Sqls,omitnil,omitempty" name:"Sqls"`

	// 影响行数，格式为M-N，例如：10-200
	AffectRowsSection *string `json:"AffectRowsSection,omitnil,omitempty" name:"AffectRowsSection"`

	// 返回行数，格式为M-N，例如：10-200
	SentRowsSection *string `json:"SentRowsSection,omitnil,omitempty" name:"SentRowsSection"`

	// 执行时间，格式为M-N，例如：10-200
	ExecTimeSection *string `json:"ExecTimeSection,omitnil,omitempty" name:"ExecTimeSection"`

	// 锁等待时间，格式为M-N，例如：10-200
	LockWaitTimeSection *string `json:"LockWaitTimeSection,omitnil,omitempty" name:"LockWaitTimeSection"`

	// IO等待时间，格式为M-N，例如：10-200
	IoWaitTimeSection *string `json:"IoWaitTimeSection,omitnil,omitempty" name:"IoWaitTimeSection"`

	// 事务持续时间，格式为M-N，例如：10-200
	TransactionLivingTimeSection *string `json:"TransactionLivingTimeSection,omitnil,omitempty" name:"TransactionLivingTimeSection"`

	// 线程ID
	ThreadId []*string `json:"ThreadId,omitnil,omitempty" name:"ThreadId"`

	// 返回行数。表示筛选返回行数大于该值的审计日志。
	SentRows *int64 `json:"SentRows,omitnil,omitempty" name:"SentRows"`

	// mysql错误码
	ErrCode []*int64 `json:"ErrCode,omitnil,omitempty" name:"ErrCode"`
}

type AuditPolicy struct {
	// 审计策略 ID。
	PolicyId *string `json:"PolicyId,omitnil,omitempty" name:"PolicyId"`

	// 审计策略的状态。可能返回的值为：
	// "creating" - 创建中;
	// "running" - 运行中;
	// "paused" - 暂停中;
	// "failed" - 创建失败。
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// 数据库实例 ID。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 审计策略创建时间。格式为 : "2019-03-20 17:09:13"。
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`

	// 审计策略最后修改时间。格式为 : "2019-03-20 17:09:13"。
	ModifyTime *string `json:"ModifyTime,omitnil,omitempty" name:"ModifyTime"`

	// 审计策略名称。
	PolicyName *string `json:"PolicyName,omitnil,omitempty" name:"PolicyName"`

	// 审计规则 ID。
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`

	// 审计规则名称。
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`

	// 数据库实例名称
	InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`
}

type AuditRule struct {
	// 审计规则 Id。
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`

	// 审计规则创建时间。格式为 : "2019-03-20 17:09:13"。
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`

	// 审计规则最后修改时间。格式为 : "2019-03-20 17:09:13"。
	ModifyTime *string `json:"ModifyTime,omitnil,omitempty" name:"ModifyTime"`

	// 审计规则名称。
	// 注意：此字段可能返回 null，表示取不到有效值。
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`

	// 审计规则描述。
	// 注意：此字段可能返回 null，表示取不到有效值。
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`

	// 审计规则过滤条件。
	// 注意：此字段可能返回 null，表示取不到有效值。
	RuleFilters []*AuditFilter `json:"RuleFilters,omitnil,omitempty" name:"RuleFilters"`

	// 是否开启全审计。
	AuditAll *bool `json:"AuditAll,omitnil,omitempty" name:"AuditAll"`
}

type AuditRuleFilters struct {
	// 单条审计规则。
	RuleFilters []*RuleFilters `json:"RuleFilters,omitnil,omitempty" name:"RuleFilters"`
}

type AuditRuleTemplateInfo struct {
	// 规则模板ID。
	RuleTemplateId *string `json:"RuleTemplateId,omitnil,omitempty" name:"RuleTemplateId"`

	// 规则模板名称。
	RuleTemplateName *string `json:"RuleTemplateName,omitnil,omitempty" name:"RuleTemplateName"`

	// 规则模板的过滤条件。
	RuleFilters []*RuleFilters `json:"RuleFilters,omitnil,omitempty" name:"RuleFilters"`

	// 规则模板描述。
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`

	// 规则模板创建时间。
	CreateAt *string `json:"CreateAt,omitnil,omitempty" name:"CreateAt"`

	// 告警等级。1-低风险，2-中风险，3-高风险。
	AlarmLevel *uint64 `json:"AlarmLevel,omitnil,omitempty" name:"AlarmLevel"`

	// 告警策略。0-不告警，1-告警。
	AlarmPolicy *uint64 `json:"AlarmPolicy,omitnil,omitempty" name:"AlarmPolicy"`

	// 规则模板应用在哪些在实例。
	AffectedInstances []*string `json:"AffectedInstances,omitnil,omitempty" name:"AffectedInstances"`

	// 模板状态。0-无任务 ，1-修改中。
	Status *uint64 `json:"Status,omitnil,omitempty" name:"Status"`

	// 模板更新时间。
	UpdateAt *string `json:"UpdateAt,omitnil,omitempty" name:"UpdateAt"`
}

type AutoStrategy struct {
	// 自动扩容阈值，可选值40、50、60、70、80、90，代表 CPU 利用率达到40%、50%、60%、70%、80%、90%时后台进行自动扩容。
	ExpandThreshold *int64 `json:"ExpandThreshold,omitnil,omitempty" name:"ExpandThreshold"`

	// 自动缩容阈值，可选值10、20、30，代表CPU利用率达到10%、20%、30%时后台进行自动缩容
	ShrinkThreshold *int64 `json:"ShrinkThreshold,omitnil,omitempty" name:"ShrinkThreshold"`

	// 自动扩容观测周期，单位是分钟，可选值1、3、5、10、15、30。后台会按照配置的周期进行扩容判断。
	// 注意：此字段可能返回 null，表示取不到有效值。
	//
	// Deprecated: ExpandPeriod is deprecated.
	ExpandPeriod *int64 `json:"ExpandPeriod,omitnil,omitempty" name:"ExpandPeriod"`

	// 自动缩容观测周期，单位是分钟，可选值5、10、15、30。后台会按照配置的周期进行缩容判断。
	// 注意：此字段可能返回 null，表示取不到有效值。
	//
	// Deprecated: ShrinkPeriod is deprecated.
	ShrinkPeriod *int64 `json:"ShrinkPeriod,omitnil,omitempty" name:"ShrinkPeriod"`

	// 弹性扩容观测周期（秒级），可取值为：15，30，45，60，180，300，600，900，1800。
	ExpandSecondPeriod *int64 `json:"ExpandSecondPeriod,omitnil,omitempty" name:"ExpandSecondPeriod"`

	// 缩容观测周期（秒级），可取值为：300、600、900、1800。
	ShrinkSecondPeriod *int64 `json:"ShrinkSecondPeriod,omitnil,omitempty" name:"ShrinkSecondPeriod"`
}

type BackupConfig struct {
	// 第二个从库复制方式，可能的返回值：async-异步，semisync-半同步
	ReplicationMode *string `json:"ReplicationMode,omitnil,omitempty" name:"ReplicationMode"`

	// 第二个从库可用区的正式名称，如 ap-shanghai-2
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`

	// 第二个从库内网IP地址
	Vip *string `json:"Vip,omitnil,omitempty" name:"Vip"`

	// 第二个从库访问端口
	Vport *uint64 `json:"Vport,omitnil,omitempty" name:"Vport"`
}

type BackupInfo struct {
	// <p>备份文件名</p>
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// <p>备份文件大小，单位：Byte</p>
	Size *int64 `json:"Size,omitnil,omitempty" name:"Size"`

	// <p>备份快照时间，时间格式：2016-03-17 02:10:37</p>
	Date *string `json:"Date,omitnil,omitempty" name:"Date"`

	// <p>下载地址</p>
	IntranetUrl *string `json:"IntranetUrl,omitnil,omitempty" name:"IntranetUrl"`

	// <p>下载地址</p>
	InternetUrl *string `json:"InternetUrl,omitnil,omitempty" name:"InternetUrl"`

	// <p>日志具体类型。可能的值有 &quot;logical&quot;: 逻辑冷备， &quot;physical&quot;: 物理冷备。</p>
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// <p>备份子任务的ID，删除备份文件时使用</p>
	BackupId *int64 `json:"BackupId,omitnil,omitempty" name:"BackupId"`

	// <p>备份任务状态。可能的值有 &quot;SUCCESS&quot;: 备份成功， &quot;FAILED&quot;: 备份失败， &quot;RUNNING&quot;: 备份进行中。</p>
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// <p>备份任务的完成时间</p>
	FinishTime *string `json:"FinishTime,omitnil,omitempty" name:"FinishTime"`

	// <p>（该值将废弃，不建议使用）备份的创建者，可能的值：SYSTEM - 系统创建，Uin - 发起者Uin值。</p>
	Creator *string `json:"Creator,omitnil,omitempty" name:"Creator"`

	// <p>备份任务的开始时间</p>
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// <p>备份方法。可能的值有 &quot;full&quot;: 全量备份， &quot;partial&quot;: 部分备份。</p>
	Method *string `json:"Method,omitnil,omitempty" name:"Method"`

	// <p>备份方式。可能的值有 &quot;manual&quot;: 手动备份， &quot;automatic&quot;: 自动备份。</p>
	Way *string `json:"Way,omitnil,omitempty" name:"Way"`

	// <p>手动备份别名</p>
	ManualBackupName *string `json:"ManualBackupName,omitnil,omitempty" name:"ManualBackupName"`

	// <p>备份保留类型，save_mode_regular - 常规保存备份，save_mode_period - 定期保存备份</p>
	SaveMode *string `json:"SaveMode,omitnil,omitempty" name:"SaveMode"`

	// <p>本地备份所在地域</p>
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// <p>异地备份详细信息</p>
	RemoteInfo []*RemoteBackupInfo `json:"RemoteInfo,omitnil,omitempty" name:"RemoteInfo"`

	// <p>存储方式，0-常规存储，1-归档存储，2-标准存储，默认为0</p>
	CosStorageType *int64 `json:"CosStorageType,omitnil,omitempty" name:"CosStorageType"`

	// <p>实例 ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例 ID 相同。</p>
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// <p>备份完成进度</p>
	Progress *int64 `json:"Progress,omitnil,omitempty" name:"Progress"`

	// <p>备份文件是否加密， on-加密， off-未加密</p>
	EncryptionFlag *string `json:"EncryptionFlag,omitnil,omitempty" name:"EncryptionFlag"`

	// <p>备份GTID点位</p>
	ExecutedGTIDSet *string `json:"ExecutedGTIDSet,omitnil,omitempty" name:"ExecutedGTIDSet"`

	// <p>备份文件MD5值</p>
	MD5 *string `json:"MD5,omitnil,omitempty" name:"MD5"`
}

type BackupItem struct {
	// 需要备份的库名
	Db *string `json:"Db,omitnil,omitempty" name:"Db"`

	// 需要备份的表名。 如果传该参数，表示备份该库中的指定表。如果不传该参数则备份该db库
	Table *string `json:"Table,omitnil,omitempty" name:"Table"`
}

type BackupLimitVpcItem struct {
	// 限制下载来源的地域。目前仅支持当前地域。
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// 限制下载的vpc列表。
	VpcList []*string `json:"VpcList,omitnil,omitempty" name:"VpcList"`
}

type BackupSummaryItem struct {
	// 实例ID。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 该实例自动数据备份的个数。
	AutoBackupCount *int64 `json:"AutoBackupCount,omitnil,omitempty" name:"AutoBackupCount"`

	// 该实例自动数据备份的容量。
	AutoBackupVolume *int64 `json:"AutoBackupVolume,omitnil,omitempty" name:"AutoBackupVolume"`

	// 该实例手动数据备份的个数。
	ManualBackupCount *int64 `json:"ManualBackupCount,omitnil,omitempty" name:"ManualBackupCount"`

	// 该实例手动数据备份的容量。
	ManualBackupVolume *int64 `json:"ManualBackupVolume,omitnil,omitempty" name:"ManualBackupVolume"`

	// 该实例总的数据备份（包含自动备份和手动备份）个数。
	DataBackupCount *int64 `json:"DataBackupCount,omitnil,omitempty" name:"DataBackupCount"`

	// 该实例总的数据备份容量。
	DataBackupVolume *int64 `json:"DataBackupVolume,omitnil,omitempty" name:"DataBackupVolume"`

	// 该实例日志备份的个数。
	BinlogBackupCount *int64 `json:"BinlogBackupCount,omitnil,omitempty" name:"BinlogBackupCount"`

	// 该实例日志备份的容量。
	BinlogBackupVolume *int64 `json:"BinlogBackupVolume,omitnil,omitempty" name:"BinlogBackupVolume"`

	// 该实例的总备份（包含数据备份和日志备份）占用容量。
	BackupVolume *int64 `json:"BackupVolume,omitnil,omitempty" name:"BackupVolume"`
}

// Predefined struct for user
type BalanceRoGroupLoadRequestParams struct {
	// RO 组的 ID，格式如：cdbrg-c1nl9rpv。可通过 [DescribeRoGroups](https://cloud.tencent.com/document/api/236/40939) 获取。
	RoGroupId *string `json:"RoGroupId,omitnil,omitempty" name:"RoGroupId"`
}

type BalanceRoGroupLoadRequest struct {
	*tchttp.BaseRequest

	// RO 组的 ID，格式如：cdbrg-c1nl9rpv。可通过 [DescribeRoGroups](https://cloud.tencent.com/document/api/236/40939) 获取。
	RoGroupId *string `json:"RoGroupId,omitnil,omitempty" name:"RoGroupId"`
}

func (r *BalanceRoGroupLoadRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *BalanceRoGroupLoadRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "RoGroupId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "BalanceRoGroupLoadRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type BalanceRoGroupLoadResponseParams struct {
	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type BalanceRoGroupLoadResponse struct {
	*tchttp.BaseResponse
	Response *BalanceRoGroupLoadResponseParams `json:"Response"`
}

func (r *BalanceRoGroupLoadResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *BalanceRoGroupLoadResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type BinlogInfo struct {
	// <p>binlog 日志备份文件名</p>
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// <p>备份文件大小，单位：Byte</p>
	Size *int64 `json:"Size,omitnil,omitempty" name:"Size"`

	// <p>文件存储时间，时间格式：2016-03-17 02:10:37</p>
	Date *string `json:"Date,omitnil,omitempty" name:"Date"`

	// <p>下载地址<br>说明：此下载地址和参数 InternetUrl 的下载地址一样。</p>
	IntranetUrl *string `json:"IntranetUrl,omitnil,omitempty" name:"IntranetUrl"`

	// <p>下载地址<br>说明：此下载地址和参数 IntranetUrl 的下载地址一样。</p>
	InternetUrl *string `json:"InternetUrl,omitnil,omitempty" name:"InternetUrl"`

	// <p>日志具体类型，可能的值有：binlog - 二进制日志</p>
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// <p>binlog 文件起始时间</p>
	BinlogStartTime *string `json:"BinlogStartTime,omitnil,omitempty" name:"BinlogStartTime"`

	// <p>binlog 文件截止时间</p>
	BinlogFinishTime *string `json:"BinlogFinishTime,omitnil,omitempty" name:"BinlogFinishTime"`

	// <p>本地binlog文件所在地域</p>
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// <p>备份任务状态。可能的值有 &quot;SUCCESS&quot;: 备份成功， &quot;FAILED&quot;: 备份失败， &quot;RUNNING&quot;: 备份进行中。</p>
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// <p>binlog异地备份详细信息</p>
	RemoteInfo []*RemoteBackupInfo `json:"RemoteInfo,omitnil,omitempty" name:"RemoteInfo"`

	// <p>存储方式，0-常规存储，1-归档存储，2-标准存储，默认为0</p>
	CosStorageType *int64 `json:"CosStorageType,omitnil,omitempty" name:"CosStorageType"`

	// <p>实例 ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例 ID 相同。</p>
	//
	// Deprecated: InstanceId is deprecated.
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// <p>备份完成进度</p>
	Progress *int64 `json:"Progress,omitnil,omitempty" name:"Progress"`
}

type Bucket struct {
	// 无
	Key *string `json:"Key,omitnil,omitempty" name:"Key"`

	// key值出现的次数。
	Count *uint64 `json:"Count,omitnil,omitempty" name:"Count"`
}

type CdbRegionSellConf struct {
	// 地域中文名称
	RegionName *string `json:"RegionName,omitnil,omitempty" name:"RegionName"`

	// 所属大区
	Area *string `json:"Area,omitnil,omitempty" name:"Area"`

	// 是否为默认地域
	IsDefaultRegion *int64 `json:"IsDefaultRegion,omitnil,omitempty" name:"IsDefaultRegion"`

	// 地域名称
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// 地域的可用区售卖配置
	RegionConfig []*CdbZoneSellConf `json:"RegionConfig,omitnil,omitempty" name:"RegionConfig"`
}

type CdbSellConfig struct {
	// 内存大小，单位为MB
	Memory *int64 `json:"Memory,omitnil,omitempty" name:"Memory"`

	// CPU核心数
	Cpu *int64 `json:"Cpu,omitnil,omitempty" name:"Cpu"`

	// 磁盘最小规格，单位为GB
	VolumeMin *int64 `json:"VolumeMin,omitnil,omitempty" name:"VolumeMin"`

	// 磁盘最大规格，单位为GB
	VolumeMax *int64 `json:"VolumeMax,omitnil,omitempty" name:"VolumeMax"`

	// 磁盘步长，单位为GB
	VolumeStep *int64 `json:"VolumeStep,omitnil,omitempty" name:"VolumeStep"`

	// 每秒IO数量
	Iops *int64 `json:"Iops,omitnil,omitempty" name:"Iops"`

	// 应用场景描述
	Info *string `json:"Info,omitnil,omitempty" name:"Info"`

	// 状态值，0 表示该规格对外售卖
	Status *int64 `json:"Status,omitnil,omitempty" name:"Status"`

	// 实例类型，可能的取值范围有：UNIVERSAL (通用型), EXCLUSIVE (独享型), BASIC (基础型), BASIC_V2 (基础型v2)
	DeviceType *string `json:"DeviceType,omitnil,omitempty" name:"DeviceType"`

	// 引擎类型描述，可能的取值范围有：Innodb，RocksDB
	EngineType *string `json:"EngineType,omitnil,omitempty" name:"EngineType"`

	// 售卖规格Id
	Id *int64 `json:"Id,omitnil,omitempty" name:"Id"`
}

type CdbSellType struct {
	// 售卖实例名称。
	// Z3：是高可用类型，对应规格中的 DeviceType，包含 UNIVERSAL，EXCLUSIVE。
	// CVM：是基础版类型，对应规格中的 DeviceType 是 BASIC（已下线）。
	// TKE：是基础版v2类型，对应规格中的 DeviceType 是 BASIC_V2。
	// CLOUD_NATIVE_CLUSTER：表示云盘版标准型。
	// CLOUD_NATIVE_CLUSTER_EXCLUSIVE：表示云盘版加强型。
	// ECONOMICAL：表示经济型。
	TypeName *string `json:"TypeName,omitnil,omitempty" name:"TypeName"`

	// 引擎版本号
	EngineVersion []*string `json:"EngineVersion,omitnil,omitempty" name:"EngineVersion"`

	// 售卖规格Id
	ConfigIds []*int64 `json:"ConfigIds,omitnil,omitempty" name:"ConfigIds"`
}

type CdbZoneDataResult struct {
	// 售卖规格所有集合
	Configs []*CdbSellConfig `json:"Configs,omitnil,omitempty" name:"Configs"`

	// 售卖地域可用区集合
	Regions []*CdbRegionSellConf `json:"Regions,omitnil,omitempty" name:"Regions"`
}

type CdbZoneSellConf struct {
	// 可用区状态。可能的返回值为：1-上线；3-停售；4-不展示
	Status *int64 `json:"Status,omitnil,omitempty" name:"Status"`

	// 可用区中文名称
	ZoneName *string `json:"ZoneName,omitnil,omitempty" name:"ZoneName"`

	// 实例类型是否为自定义类型
	IsCustom *bool `json:"IsCustom,omitnil,omitempty" name:"IsCustom"`

	// 是否支持灾备
	IsSupportDr *bool `json:"IsSupportDr,omitnil,omitempty" name:"IsSupportDr"`

	// 是否支持私有网络
	IsSupportVpc *bool `json:"IsSupportVpc,omitnil,omitempty" name:"IsSupportVpc"`

	// 小时计费实例最大售卖数量
	HourInstanceSaleMaxNum *int64 `json:"HourInstanceSaleMaxNum,omitnil,omitempty" name:"HourInstanceSaleMaxNum"`

	// 是否为默认可用区
	IsDefaultZone *bool `json:"IsDefaultZone,omitnil,omitempty" name:"IsDefaultZone"`

	// 是否为黑石区
	IsBm *bool `json:"IsBm,omitnil,omitempty" name:"IsBm"`

	// 支持的付费类型。可能的返回值为：0-包年包月；1-小时计费；2-后付费
	PayType []*string `json:"PayType,omitnil,omitempty" name:"PayType"`

	// 数据复制类型。0-异步复制；1-半同步复制；2-强同步复制
	ProtectMode []*string `json:"ProtectMode,omitnil,omitempty" name:"ProtectMode"`

	// 可用区名称
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`

	// 多可用区信息
	ZoneConf *ZoneConf `json:"ZoneConf,omitnil,omitempty" name:"ZoneConf"`

	// 可支持的灾备可用区信息
	DrZone []*string `json:"DrZone,omitnil,omitempty" name:"DrZone"`

	// 是否支持跨可用区只读
	IsSupportRemoteRo *bool `json:"IsSupportRemoteRo,omitnil,omitempty" name:"IsSupportRemoteRo"`

	// 可支持的跨可用区只读区信息
	RemoteRoZone []*string `json:"RemoteRoZone,omitnil,omitempty" name:"RemoteRoZone"`

	// 独享型可用区状态。可能的返回值为：1-上线；3-停售；4-不展示
	ExClusterStatus *int64 `json:"ExClusterStatus,omitnil,omitempty" name:"ExClusterStatus"`

	// 独享型可支持的跨可用区只读区信息
	ExClusterRemoteRoZone []*string `json:"ExClusterRemoteRoZone,omitnil,omitempty" name:"ExClusterRemoteRoZone"`

	// 独享型多可用区信息
	ExClusterZoneConf *ZoneConf `json:"ExClusterZoneConf,omitnil,omitempty" name:"ExClusterZoneConf"`

	// 售卖实例类型数组，其中configIds的值与configs结构体中的id一一对应。
	SellType []*CdbSellType `json:"SellType,omitnil,omitempty" name:"SellType"`

	// 可用区id
	ZoneId *int64 `json:"ZoneId,omitnil,omitempty" name:"ZoneId"`

	// 是否支持ipv6
	IsSupportIpv6 *bool `json:"IsSupportIpv6,omitnil,omitempty" name:"IsSupportIpv6"`

	// 可支持的售卖数据库引擎类型
	EngineType []*string `json:"EngineType,omitnil,omitempty" name:"EngineType"`

	// 云盘版实例在当前可用区的售卖状态。可能的返回值为：1-上线；3-停售；4-不展示
	CloudNativeClusterStatus *int64 `json:"CloudNativeClusterStatus,omitnil,omitempty" name:"CloudNativeClusterStatus"`

	// 云盘版或者单节点基础型支持的磁盘类型。
	DiskTypeConf []*DiskTypeConfigItem `json:"DiskTypeConf,omitnil,omitempty" name:"DiskTypeConf"`
}

// Predefined struct for user
type CheckMigrateClusterRequestParams struct {
	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 实例 CPU 核数。当 InstanceId 为主实例时必传。
	Cpu *int64 `json:"Cpu,omitnil,omitempty" name:"Cpu"`

	// 实例内存大小，单位：MB。当 InstanceId 为主实例时必传。
	Memory *int64 `json:"Memory,omitnil,omitempty" name:"Memory"`

	// 实例硬盘大小，单位：GB。
	Volume *int64 `json:"Volume,omitnil,omitempty" name:"Volume"`

	// 磁盘类型。 CLOUD_SSD: SSD 云硬盘; CLOUD_HSSD: 增强型 SSD 云硬盘。
	DiskType *string `json:"DiskType,omitnil,omitempty" name:"DiskType"`

	// 云盘版节点拓扑配置。当 InstanceId 为主实例时必传。
	ClusterTopology *ClusterTopology `json:"ClusterTopology,omitnil,omitempty" name:"ClusterTopology"`

	// 迁移实例类型。支持值包括： "CLOUD_NATIVE_CLUSTER" - 标准型云盘版实例， "CLOUD_NATIVE_CLUSTER_EXCLUSIVE" - 加强型云盘版实例。
	DeviceType *string `json:"DeviceType,omitnil,omitempty" name:"DeviceType"`

	// 只读实例信息。
	RoInfo []*MigrateClusterRoInfo `json:"RoInfo,omitnil,omitempty" name:"RoInfo"`
}

type CheckMigrateClusterRequest struct {
	*tchttp.BaseRequest

	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 实例 CPU 核数。当 InstanceId 为主实例时必传。
	Cpu *int64 `json:"Cpu,omitnil,omitempty" name:"Cpu"`

	// 实例内存大小，单位：MB。当 InstanceId 为主实例时必传。
	Memory *int64 `json:"Memory,omitnil,omitempty" name:"Memory"`

	// 实例硬盘大小，单位：GB。
	Volume *int64 `json:"Volume,omitnil,omitempty" name:"Volume"`

	// 磁盘类型。 CLOUD_SSD: SSD 云硬盘; CLOUD_HSSD: 增强型 SSD 云硬盘。
	DiskType *string `json:"DiskType,omitnil,omitempty" name:"DiskType"`

	// 云盘版节点拓扑配置。当 InstanceId 为主实例时必传。
	ClusterTopology *ClusterTopology `json:"ClusterTopology,omitnil,omitempty" name:"ClusterTopology"`

	// 迁移实例类型。支持值包括： "CLOUD_NATIVE_CLUSTER" - 标准型云盘版实例， "CLOUD_NATIVE_CLUSTER_EXCLUSIVE" - 加强型云盘版实例。
	DeviceType *string `json:"DeviceType,omitnil,omitempty" name:"DeviceType"`

	// 只读实例信息。
	RoInfo []*MigrateClusterRoInfo `json:"RoInfo,omitnil,omitempty" name:"RoInfo"`
}

func (r *CheckMigrateClusterRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CheckMigrateClusterRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "Cpu")
	delete(f, "Memory")
	delete(f, "Volume")
	delete(f, "DiskType")
	delete(f, "ClusterTopology")
	delete(f, "DeviceType")
	delete(f, "RoInfo")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CheckMigrateClusterRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CheckMigrateClusterResponseParams struct {
	// 校验是否通过，通过为pass，失败为fail
	CheckResult *string `json:"CheckResult,omitnil,omitempty" name:"CheckResult"`

	// 校验项
	Items []*CheckMigrateResult `json:"Items,omitnil,omitempty" name:"Items"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CheckMigrateClusterResponse struct {
	*tchttp.BaseResponse
	Response *CheckMigrateClusterResponseParams `json:"Response"`
}

func (r *CheckMigrateClusterResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CheckMigrateClusterResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type CheckMigrateResult struct {
	// 校验名称
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// 校验结果，通过为pass，失败为fail
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`

	// 校验结果描述
	Desc *string `json:"Desc,omitnil,omitempty" name:"Desc"`
}

type CloneItem struct {
	// 克隆任务的源实例Id。
	SrcInstanceId *string `json:"SrcInstanceId,omitnil,omitempty" name:"SrcInstanceId"`

	// 克隆任务的新产生实例Id。
	DstInstanceId *string `json:"DstInstanceId,omitnil,omitempty" name:"DstInstanceId"`

	// 克隆任务对应的任务列表Id。
	CloneJobId *int64 `json:"CloneJobId,omitnil,omitempty" name:"CloneJobId"`

	// 克隆实例使用的策略， 包括以下类型：  timepoint:指定时间点回档，  backupset: 指定备份文件回档。
	RollbackStrategy *string `json:"RollbackStrategy,omitnil,omitempty" name:"RollbackStrategy"`

	// 克隆实例回档的时间点。
	RollbackTargetTime *string `json:"RollbackTargetTime,omitnil,omitempty" name:"RollbackTargetTime"`

	// 任务开始时间。
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// 任务结束时间。
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// 任务状态，包括以下状态：initial,running,wait_complete,success,failed
	TaskStatus *string `json:"TaskStatus,omitnil,omitempty" name:"TaskStatus"`

	// 克隆实例所在地域Id
	NewRegionId *int64 `json:"NewRegionId,omitnil,omitempty" name:"NewRegionId"`

	// 源实例所在地域Id
	SrcRegionId *int64 `json:"SrcRegionId,omitnil,omitempty" name:"SrcRegionId"`
}

// Predefined struct for user
type CloseAuditServiceRequestParams struct {
	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
}

type CloseAuditServiceRequest struct {
	*tchttp.BaseRequest

	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
}

func (r *CloseAuditServiceRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CloseAuditServiceRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CloseAuditServiceRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CloseAuditServiceResponseParams struct {
	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CloseAuditServiceResponse struct {
	*tchttp.BaseResponse
	Response *CloseAuditServiceResponseParams `json:"Response"`
}

func (r *CloseAuditServiceResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CloseAuditServiceResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CloseCDBProxyRequestParams struct {
	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 代理组 ID。可通过 [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) 接口获取。
	ProxyGroupId *string `json:"ProxyGroupId,omitnil,omitempty" name:"ProxyGroupId"`

	// 是否只关闭读写分离，取值："true" | "false"，默认为"false"
	OnlyCloseRW *bool `json:"OnlyCloseRW,omitnil,omitempty" name:"OnlyCloseRW"`
}

type CloseCDBProxyRequest struct {
	*tchttp.BaseRequest

	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 代理组 ID。可通过 [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) 接口获取。
	ProxyGroupId *string `json:"ProxyGroupId,omitnil,omitempty" name:"ProxyGroupId"`

	// 是否只关闭读写分离，取值："true" | "false"，默认为"false"
	OnlyCloseRW *bool `json:"OnlyCloseRW,omitnil,omitempty" name:"OnlyCloseRW"`
}

func (r *CloseCDBProxyRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CloseCDBProxyRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "ProxyGroupId")
	delete(f, "OnlyCloseRW")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CloseCDBProxyRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CloseCDBProxyResponseParams struct {
	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CloseCDBProxyResponse struct {
	*tchttp.BaseResponse
	Response *CloseCDBProxyResponseParams `json:"Response"`
}

func (r *CloseCDBProxyResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CloseCDBProxyResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CloseCdbProxyAddressRequestParams struct {
	// 代理组 ID。可通过 [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) 接口获取。
	ProxyGroupId *string `json:"ProxyGroupId,omitnil,omitempty" name:"ProxyGroupId"`

	// 代理组地址 ID。可通过 [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) 接口获取。
	ProxyAddressId *string `json:"ProxyAddressId,omitnil,omitempty" name:"ProxyAddressId"`
}

type CloseCdbProxyAddressRequest struct {
	*tchttp.BaseRequest

	// 代理组 ID。可通过 [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) 接口获取。
	ProxyGroupId *string `json:"ProxyGroupId,omitnil,omitempty" name:"ProxyGroupId"`

	// 代理组地址 ID。可通过 [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) 接口获取。
	ProxyAddressId *string `json:"ProxyAddressId,omitnil,omitempty" name:"ProxyAddressId"`
}

func (r *CloseCdbProxyAddressRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CloseCdbProxyAddressRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProxyGroupId")
	delete(f, "ProxyAddressId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CloseCdbProxyAddressRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CloseCdbProxyAddressResponseParams struct {
	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CloseCdbProxyAddressResponse struct {
	*tchttp.BaseResponse
	Response *CloseCdbProxyAddressResponseParams `json:"Response"`
}

func (r *CloseCdbProxyAddressResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CloseCdbProxyAddressResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CloseSSLRequestParams struct {
	// 实例 ID。只读组 ID 为空时必填。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 只读组 ID。实例 ID 为空时必填。可通过 [DescribeRoGroups](https://cloud.tencent.com/document/api/236/40939) 接口获取。
	RoGroupId *string `json:"RoGroupId,omitnil,omitempty" name:"RoGroupId"`
}

type CloseSSLRequest struct {
	*tchttp.BaseRequest

	// 实例 ID。只读组 ID 为空时必填。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 只读组 ID。实例 ID 为空时必填。可通过 [DescribeRoGroups](https://cloud.tencent.com/document/api/236/40939) 接口获取。
	RoGroupId *string `json:"RoGroupId,omitnil,omitempty" name:"RoGroupId"`
}

func (r *CloseSSLRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CloseSSLRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "RoGroupId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CloseSSLRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CloseSSLResponseParams struct {
	// 异步请求 ID。
	AsyncRequestId *string `json:"AsyncRequestId,omitnil,omitempty" name:"AsyncRequestId"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CloseSSLResponse struct {
	*tchttp.BaseResponse
	Response *CloseSSLResponseParams `json:"Response"`
}

func (r *CloseSSLResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CloseSSLResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CloseWanServiceRequestParams struct {
	// 实例 ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例 ID 相同，可使用 [查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值。可传入只读组 ID 关闭只读组外网访问。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 变更云盘版实例只读组时，InstanceId 传实例 ID，需要额外指定该参数表示操作只读组。如果操作读写节点则不需指定该参数。
	OpResourceId *string `json:"OpResourceId,omitnil,omitempty" name:"OpResourceId"`
}

type CloseWanServiceRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例 ID 相同，可使用 [查询实例列表](https://cloud.tencent.com/document/api/236/15872) 接口获取，其值为输出参数中字段 InstanceId 的值。可传入只读组 ID 关闭只读组外网访问。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 变更云盘版实例只读组时，InstanceId 传实例 ID，需要额外指定该参数表示操作只读组。如果操作读写节点则不需指定该参数。
	OpResourceId *string `json:"OpResourceId,omitnil,omitempty" name:"OpResourceId"`
}

func (r *CloseWanServiceRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CloseWanServiceRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "OpResourceId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CloseWanServiceRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CloseWanServiceResponseParams struct {
	// 异步任务的请求 ID，可使用此 ID 查询异步任务的执行结果。
	AsyncRequestId *string `json:"AsyncRequestId,omitnil,omitempty" name:"AsyncRequestId"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CloseWanServiceResponse struct {
	*tchttp.BaseResponse
	Response *CloseWanServiceResponseParams `json:"Response"`
}

func (r *CloseWanServiceResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CloseWanServiceResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type ClusterInfo struct {
	// 节点id
	NodeId *string `json:"NodeId,omitnil,omitempty" name:"NodeId"`

	// 节点类型：主节点，从节点
	Role *string `json:"Role,omitnil,omitempty" name:"Role"`

	// 地域
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`
}

type ClusterNodeInfo struct {
	// 节点id。
	NodeId *string `json:"NodeId,omitnil,omitempty" name:"NodeId"`

	// 节点的角色。
	Role *string `json:"Role,omitnil,omitempty" name:"Role"`

	// 节点所在可用区。
	Zone *string `json:"Zone,omitnil,omitempty" name:"Zone"`

	// 节点的权重
	Weight *uint64 `json:"Weight,omitnil,omitempty" name:"Weight"`

	// 节点状态。
	Status *string `json:"Status,omitnil,omitempty" name:"Status"`
}

type ClusterTopology struct {
	// RW 节点拓扑。
	// 说明：NodeId 可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 获取。
	ReadWriteNode *ReadWriteNode `json:"ReadWriteNode,omitnil,omitempty" name:"ReadWriteNode"`

	// RO 节点拓扑。
	// 说明：NodeId 可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 获取。
	ReadOnlyNodes []*ReadonlyNode `json:"ReadOnlyNodes,omitnil,omitempty" name:"ReadOnlyNodes"`
}

type ColumnPrivilege struct {
	// 数据库名
	Database *string `json:"Database,omitnil,omitempty" name:"Database"`

	// 数据库表名
	Table *string `json:"Table,omitnil,omitempty" name:"Table"`

	// 数据库列名
	Column *string `json:"Column,omitnil,omitempty" name:"Column"`

	// 权限信息
	Privileges []*string `json:"Privileges,omitnil,omitempty" name:"Privileges"`
}

type CommonTimeWindow struct {
	// 周一的时间窗，格式如： 02:00-06:00
	Monday *string `json:"Monday,omitnil,omitempty" name:"Monday"`

	// 周二的时间窗，格式如： 02:00-06:00
	Tuesday *string `json:"Tuesday,omitnil,omitempty" name:"Tuesday"`

	// 周三的时间窗，格式如： 02:00-06:00
	Wednesday *string `json:"Wednesday,omitnil,omitempty" name:"Wednesday"`

	// 周四的时间窗，格式如： 02:00-06:00
	Thursday *string `json:"Thursday,omitnil,omitempty" name:"Thursday"`

	// 周五的时间窗，格式如： 02:00-06:00
	Friday *string `json:"Friday,omitnil,omitempty" name:"Friday"`

	// 周六的时间窗，格式如： 02:00-06:00
	Saturday *string `json:"Saturday,omitnil,omitempty" name:"Saturday"`

	// 周日的时间窗，格式如： 02:00-06:00
	Sunday *string `json:"Sunday,omitnil,omitempty" name:"Sunday"`

	// 常规备份保留策略，weekly-按周备份，monthly-按月备份，默认为weekly
	BackupPeriodStrategy *string `json:"BackupPeriodStrategy,omitnil,omitempty" name:"BackupPeriodStrategy"`

	// 如果设置为按月备份，需填入每月具体备份日期，相邻备份天数不得超过两天。例[1,4,7,9,11,14,17,19,22,25,28,30,31]
	Days []*int64 `json:"Days,omitnil,omitempty" name:"Days"`

	// 月度备份时间窗，BackupPeriodStrategy为monthly时必填。格式如： 02:00-06:00
	BackupPeriodTime *string `json:"BackupPeriodTime,omitnil,omitempty" name:"BackupPeriodTime"`
}

// Predefined struct for user
type CreateAccountsRequestParams struct {
	// <p>实例 ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例 ID 相同。</p>
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// <p>云数据库账号。</p>
	Accounts []*Account `json:"Accounts,omitnil,omitempty" name:"Accounts"`

	// <p>新账户的密码。<br>说明：</p><ol><li>在8 ～ 64位字符数以内（推荐12位以上）。</li><li>至少包含其中两项：小写字母 a ~ z 或 大写字母 A ～ Z。数字0 ～ 9。_+-,&amp;=!@#$%^*().|。</li><li>不能包含非法字符。</li></ol>
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// <p>备注信息。最多支持输入255个字符。</p>
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`

	// <p>新账户最大可用连接数，默认值为10240，最大可设置值为10240。</p>
	MaxUserConnections *int64 `json:"MaxUserConnections,omitnil,omitempty" name:"MaxUserConnections"`
}

type CreateAccountsRequest struct {
	*tchttp.BaseRequest

	// <p>实例 ID，格式如：cdb-c1nl9rpv，与云数据库控制台页面中显示的实例 ID 相同。</p>
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// <p>云数据库账号。</p>
	Accounts []*Account `json:"Accounts,omitnil,omitempty" name:"Accounts"`

	// <p>新账户的密码。<br>说明：</p><ol><li>在8 ～ 64位字符数以内（推荐12位以上）。</li><li>至少包含其中两项：小写字母 a ~ z 或 大写字母 A ～ Z。数字0 ～ 9。_+-,&amp;=!@#$%^*().|。</li><li>不能包含非法字符。</li></ol>
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// <p>备注信息。最多支持输入255个字符。</p>
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`

	// <p>新账户最大可用连接数，默认值为10240，最大可设置值为10240。</p>
	MaxUserConnections *int64 `json:"MaxUserConnections,omitnil,omitempty" name:"MaxUserConnections"`
}

func (r *CreateAccountsRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAccountsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "Accounts")
	delete(f, "Password")
	delete(f, "Description")
	delete(f, "MaxUserConnections")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateAccountsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAccountsResponseParams struct {
	// <p>异步任务的请求 ID，可使用此 ID 查询异步任务的执行结果。</p>
	AsyncRequestId *string `json:"AsyncRequestId,omitnil,omitempty" name:"AsyncRequestId"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateAccountsResponse struct {
	*tchttp.BaseResponse
	Response *CreateAccountsResponseParams `json:"Response"`
}

func (r *CreateAccountsResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAccountsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAuditLogFileRequestParams struct {
	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 开始时间(建议开始到结束时间区间最大7天)。
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// 结束时间(建议开始到结束时间区间最大7天）。
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// 排序方式。支持值包括："ASC" - 升序，"DESC" - 降序，默认降序排序。
	Order *string `json:"Order,omitnil,omitempty" name:"Order"`

	// 排序字段。支持值包括(默认按照时间戳排序)： "timestamp" - 时间戳； "affectRows" - 影响行数； "execTime" - 执行时间。
	OrderBy *string `json:"OrderBy,omitnil,omitempty" name:"OrderBy"`

	// 已废弃。
	//
	// Deprecated: Filter is deprecated.
	Filter *AuditLogFilter `json:"Filter,omitnil,omitempty" name:"Filter"`

	// 过滤条件。可按设置的过滤条件过滤日志。
	LogFilter []*InstanceAuditLogFilters `json:"LogFilter,omitnil,omitempty" name:"LogFilter"`

	// 下载筛选列
	ColumnFilter []*string `json:"ColumnFilter,omitnil,omitempty" name:"ColumnFilter"`
}

type CreateAuditLogFileRequest struct {
	*tchttp.BaseRequest

	// 实例 ID。可通过 [DescribeDBInstances](https://cloud.tencent.com/document/product/236/15872) 接口获取。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 开始时间(建议开始到结束时间区间最大7天)。
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// 结束时间(建议开始到结束时间区间最大7天）。
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// 排序方式。支持值包括："ASC" - 升序，"DESC" - 降序，默认降序排序。
	Order *string `json:"Order,omitnil,omitempty" name:"Order"`

	// 排序字段。支持值包括(默认按照时间戳排序)： "timestamp" - 时间戳； "affectRows" - 影响行数； "execTime" - 执行时间。
	OrderBy *string `json:"OrderBy,omitnil,omitempty" name:"OrderBy"`

	// 已废弃。
	Filter *AuditLogFilter `json:"Filter,omitnil,omitempty" name:"Filter"`

	// 过滤条件。可按设置的过滤条件过滤日志。
	LogFilter []*InstanceAuditLogFilters `json:"LogFilter,omitnil,omitempty" name:"LogFilter"`

	// 下载筛选列
	ColumnFilter []*string `json:"ColumnFilter,omitnil,omitempty" name:"ColumnFilter"`
}

func (r *CreateAuditLogFileRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAuditLogFileRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "InstanceId")
	delete(f, "StartTime")
	delete(f, "EndTime")
	delete(f, "Order")
	delete(f, "OrderBy")
	delete(f, "Filter")
	delete(f, "LogFilter")
	delete(f, "ColumnFilter")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateAuditLogFileRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAuditLogFileResponseParams struct {
	// 审计日志文件名称。
	FileName *string `json:"FileName,omitnil,omitempty" name:"FileName"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateAuditLogFileResponse struct {
	*tchttp.BaseResponse
	Response *CreateAuditLogFileResponseParams `json:"Response"`
}

func (r *CreateAuditLogFileResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAuditLogFileResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAuditPolicyRequestParams struct {
	// 审计策略名称。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// 审计规则 ID。
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`

	// 实例 ID，格式如：cdb-c1nl9rpv 或者 cdbro-c1nl9rpv，与云数据库控制台页面中显示的实例 ID 相同。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 审计日志保存时长。支持值包括：
	// 7 - 一周
	// 30 - 一个月；
	// 180 - 六个月；
	// 365 - 一年；
	// 1095 - 三年；
	// 1825 - 五年；
	// 实例首次开通审计策略时，可传该值，用于设置存储日志保存天数，默认为 30 天。若实例已存在审计策略，则此参数无效，可使用 更改审计服务配置 接口修改日志存储时长。
	LogExpireDay *int64 `json:"LogExpireDay,omitnil,omitempty" name:"LogExpireDay"`
}

type CreateAuditPolicyRequest struct {
	*tchttp.BaseRequest

	// 审计策略名称。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// 审计规则 ID。
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`

	// 实例 ID，格式如：cdb-c1nl9rpv 或者 cdbro-c1nl9rpv，与云数据库控制台页面中显示的实例 ID 相同。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 审计日志保存时长。支持值包括：
	// 7 - 一周
	// 30 - 一个月；
	// 180 - 六个月；
	// 365 - 一年；
	// 1095 - 三年；
	// 1825 - 五年；
	// 实例首次开通审计策略时，可传该值，用于设置存储日志保存天数，默认为 30 天。若实例已存在审计策略，则此参数无效，可使用 更改审计服务配置 接口修改日志存储时长。
	LogExpireDay *int64 `json:"LogExpireDay,omitnil,omitempty" name:"LogExpireDay"`
}

func (r *CreateAuditPolicyRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAuditPolicyRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Name")
	delete(f, "RuleId")
	delete(f, "InstanceId")
	delete(f, "LogExpireDay")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateAuditPolicyRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAuditPolicyResponseParams struct {
	// 审计策略 ID。
	PolicyId *string `json:"PolicyId,omitnil,omitempty" name:"PolicyId"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateAuditPolicyResponse struct {
	*tchttp.BaseResponse
	Response *CreateAuditPolicyResponseParams `json:"Response"`
}

func (r *CreateAuditPolicyResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAuditPolicyResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAuditRuleRequestParams struct {
	// 审计规则名称。
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`

	// 审计规则描述。
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`

	// 审计规则过滤条件。若设置了过滤条件，则不会开启全审计。
	RuleFilters []*AuditFilter `json:"RuleFilters,omitnil,omitempty" name:"RuleFilters"`

	// 是否开启全审计。支持值包括：false – 不开启全审计，true – 开启全审计。用户未设置审计规则过滤条件时，默认开启全审计。
	AuditAll *bool `json:"AuditAll,omitnil,omitempty" name:"AuditAll"`
}

type CreateAuditRuleRequest struct {
	*tchttp.BaseRequest

	// 审计规则名称。
	RuleName *string `json:"RuleName,omitnil,omitempty" name:"RuleName"`

	// 审计规则描述。
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`

	// 审计规则过滤条件。若设置了过滤条件，则不会开启全审计。
	RuleFilters []*AuditFilter `json:"RuleFilters,omitnil,omitempty" name:"RuleFilters"`

	// 是否开启全审计。支持值包括：false – 不开启全审计，true – 开启全审计。用户未设置审计规则过滤条件时，默认开启全审计。
	AuditAll *bool `json:"AuditAll,omitnil,omitempty" name:"AuditAll"`
}

func (r *CreateAuditRuleRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAuditRuleRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "RuleName")
	delete(f, "Description")
	delete(f, "RuleFilters")
	delete(f, "AuditAll")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateAuditRuleRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAuditRuleResponseParams struct {
	// 审计规则 ID。
	RuleId *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateAuditRuleResponse struct {
	*tchttp.BaseResponse
	Response *CreateAuditRuleResponseParams `json:"Response"`
}

func (r *CreateAuditRuleResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAuditRuleResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAuditRuleTemplateRequestParams struct {
	// 审计规则。
	RuleFilters []*RuleFilters `json:"RuleFilters,omitnil,omitempty" name:"RuleFilters"`

	// 规则模板名称。最多支持输入30个字符。
	RuleTemplateName *string `json:"RuleTemplateName,omitnil,omitempty" name:"RuleTemplateName"`

	// 规则模板描述。最多支持输入200个字符。
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`

	// 告警等级。1 - 低风险，2 - 中风险，3 - 高风险。默认值为1。
	AlarmLevel *uint64 `json:"AlarmLevel,omitnil,omitempty" name:"AlarmLevel"`

	// 告警策略。0 - 不告警，1 - 告警。默认值为0。
	AlarmPolicy *uint64 `json:"AlarmPolicy,omitnil,omitempty" name:"AlarmPolicy"`
}

type CreateAuditRuleTemplateRequest struct {
	*tchttp.BaseRequest

	// 审计规则。
	RuleFilters []*RuleFilters `json:"RuleFilters,omitnil,omitempty" name:"RuleFilters"`

	// 规则模板名称。最多支持输入30个字符。
	RuleTemplateName *string `json:"RuleTemplateName,omitnil,omitempty" name:"RuleTemplateName"`

	// 规则模板描述。最多支持输入200个字符。
	Description *string `json:"Description,omitnil,omitempty" name:"Description"`

	// 告警等级。1 - 低风险，2 - 中风险，3 - 高风险。默认值为1。
	AlarmLevel *uint64 `json:"AlarmLevel,omitnil,omitempty" name:"AlarmLevel"`

	// 告警策略。0 - 不告警，1 - 告警。默认值为0。
	AlarmPolicy *uint64 `json:"AlarmPolicy,omitnil,omitempty" name:"AlarmPolicy"`
}

func (r *CreateAuditRuleTemplateRequest) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAuditRuleTemplateRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "RuleFilters")
	delete(f, "RuleTemplateName")
	delete(f, "Description")
	delete(f, "AlarmLevel")
	delete(f, "AlarmPolicy")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateAuditRuleTemplateRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAuditRuleTemplateResponseParams struct {
	// 生成的规则模板ID。
	RuleTemplateId *string `json:"RuleTemplateId,omitnil,omitempty" name:"RuleTemplateId"`

	// 唯一请求 ID，由服务端生成，每次请求都会返回（若请求因其他原因未能抵达服务端，则该次请求不会获得 RequestId）。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateAuditRuleTemplateResponse struct {
	*tchttp.BaseResponse
	Response *CreateAuditRuleTemplateResponseParams `json:"Response"`
}

func (r *CreateAuditRuleTemplateResponse) ToJsonString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAuditRuleTemplateResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateBackupRequestParams struct {
	// 实例 ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例 ID 相同。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 目标备份方法，可选的值：logical - 逻辑冷备，physical - 物理冷备，snapshot - 快照备份。基础版实例仅支持快照备份。
	BackupMethod *string `json:"BackupMethod,omitnil,omitempty" name:"BackupMethod"`

	// 需要备份的库表信息，如果不设置该参数，则默认整实例备份。在 BackupMethod=logical 逻辑备份中才可设置该参数。指定的库表必须存在，否则可能导致备份失败。
	// 例：如果需要备份 db1 库的 tb1、tb2 表 和 db2 库。则该参数设置为 [{"Db": "db1", "Table": "tb1"}, {"Db": "db1", "Table": "tb2"}, {"Db": "db2"}]。
	BackupDBTableList []*BackupItem `json:"BackupDBTableList,omitnil,omitempty" name:"BackupDBTableList"`

	// 手动备份别名，输入长度请在60个字符内。
	ManualBackupName *string `json:"ManualBackupName,omitnil,omitempty" name:"ManualBackupName"`

	// 是否需要加密物理备份，可选值为：on - 是，off - 否。当 BackupMethod 为 physical 时，该值才有意义。不指定则使用实例备份默认加密策略，这里的默认加密策略指通过 [DescribeBackupEncryptionStatus](https://cloud.tencent.com/document/product/236/86508) 接口查询出的实例当前加密策略。
	EncryptionFlag *string `json:"EncryptionFlag,omitnil,omitempty" name:"EncryptionFlag"`
}

type CreateBackupRequest struct {
	*tchttp.BaseRequest

	// 实例 ID，格式如：cdb-c1nl9rpv。与云数据库控制台页面中显示的实例 ID 相同。
	InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`

	// 目标备份方法，可选的值：logical - 逻辑冷备，physical - 物理冷备，snapshot - 快照备份。基础版实例仅支持快照�