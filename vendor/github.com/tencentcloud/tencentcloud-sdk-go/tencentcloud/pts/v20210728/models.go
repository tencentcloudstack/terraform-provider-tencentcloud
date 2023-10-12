// Copyright (c) 2017-2018 THL A29 Limited, a Tencent company. All Rights Reserved.
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

package v20210728

import (
    tcerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/json"
)

// Predefined struct for user
type AbortCronJobsRequestParams struct {
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 定时任务ID数组
	CronJobIds []*string `json:"CronJobIds,omitnil" name:"CronJobIds"`
}

type AbortCronJobsRequest struct {
	*tchttp.BaseRequest
	
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 定时任务ID数组
	CronJobIds []*string `json:"CronJobIds,omitnil" name:"CronJobIds"`
}

func (r *AbortCronJobsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AbortCronJobsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "CronJobIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AbortCronJobsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AbortCronJobsResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type AbortCronJobsResponse struct {
	*tchttp.BaseResponse
	Response *AbortCronJobsResponseParams `json:"Response"`
}

func (r *AbortCronJobsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AbortCronJobsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AbortJobRequestParams struct {
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 中断原因
	AbortReason *int64 `json:"AbortReason,omitnil" name:"AbortReason"`
}

type AbortJobRequest struct {
	*tchttp.BaseRequest
	
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 中断原因
	AbortReason *int64 `json:"AbortReason,omitnil" name:"AbortReason"`
}

func (r *AbortJobRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AbortJobRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ProjectId")
	delete(f, "ScenarioId")
	delete(f, "AbortReason")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AbortJobRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AbortJobResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type AbortJobResponse struct {
	*tchttp.BaseResponse
	Response *AbortJobResponseParams `json:"Response"`
}

func (r *AbortJobResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AbortJobResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AdjustJobSpeedRequestParams struct {
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 目标RPS
	TargetRequestsPerSecond *int64 `json:"TargetRequestsPerSecond,omitnil" name:"TargetRequestsPerSecond"`
}

type AdjustJobSpeedRequest struct {
	*tchttp.BaseRequest
	
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 目标RPS
	TargetRequestsPerSecond *int64 `json:"TargetRequestsPerSecond,omitnil" name:"TargetRequestsPerSecond"`
}

func (r *AdjustJobSpeedRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AdjustJobSpeedRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "TargetRequestsPerSecond")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "AdjustJobSpeedRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type AdjustJobSpeedResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type AdjustJobSpeedResponse struct {
	*tchttp.BaseResponse
	Response *AdjustJobSpeedResponseParams `json:"Response"`
}

func (r *AdjustJobSpeedResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *AdjustJobSpeedResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type AggregationLegend struct {
	// 指标支持的聚合函数
	Aggregation *string `json:"Aggregation,omitnil" name:"Aggregation"`

	// 聚合函数作用于指标后对应的描述
	Legend *string `json:"Legend,omitnil" name:"Legend"`

	// 聚合之后的指标单位
	Unit *string `json:"Unit,omitnil" name:"Unit"`
}

type AlertChannel struct {
	// 通知模板ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`

	// AMP consumer ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	AMPConsumerId *string `json:"AMPConsumerId,omitnil" name:"AMPConsumerId"`
}

type AlertChannelRecord struct {
	// Notice ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`

	// Consumer ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	AMPConsumerId *string `json:"AMPConsumerId,omitnil" name:"AMPConsumerId"`

	// 项目 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 状态
	// 注意：此字段可能返回 null，表示取不到有效值。
	Status *uint64 `json:"Status,omitnil" name:"Status"`

	// 创建时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	CreatedAt *string `json:"CreatedAt,omitnil" name:"CreatedAt"`

	// 更新时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`

	// App ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	AppId *int64 `json:"AppId,omitnil" name:"AppId"`

	// 主账号
	// 注意：此字段可能返回 null，表示取不到有效值。
	Uin *string `json:"Uin,omitnil" name:"Uin"`

	// 子账号
	// 注意：此字段可能返回 null，表示取不到有效值。
	SubAccountUin *string `json:"SubAccountUin,omitnil" name:"SubAccountUin"`
}

type AlertRecord struct {
	// 告警历史记录项 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	AlertRecordId *string `json:"AlertRecordId,omitnil" name:"AlertRecordId"`

	// 项目 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 状态
	// 注意：此字段可能返回 null，表示取不到有效值。
	Status *AlertRecordStatus `json:"Status,omitnil" name:"Status"`

	// 创建时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	CreatedAt *string `json:"CreatedAt,omitnil" name:"CreatedAt"`

	// 修改时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`

	// 任务 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// App ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	AppId *int64 `json:"AppId,omitnil" name:"AppId"`

	// 主账号
	// 注意：此字段可能返回 null，表示取不到有效值。
	Uin *string `json:"Uin,omitnil" name:"Uin"`

	// 子账号
	// 注意：此字段可能返回 null，表示取不到有效值。
	SubAccountUin *string `json:"SubAccountUin,omitnil" name:"SubAccountUin"`

	// 场景名称
	// 注意：此字段可能返回 null，表示取不到有效值。
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`

	// 告警对象
	// 注意：此字段可能返回 null，表示取不到有效值。
	Target *string `json:"Target,omitnil" name:"Target"`

	// 告警规则 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	JobSLAId *string `json:"JobSLAId,omitnil" name:"JobSLAId"`

	// 告警规则描述
	// 注意：此字段可能返回 null，表示取不到有效值。
	JobSLADescription *string `json:"JobSLADescription,omitnil" name:"JobSLADescription"`
}

type AlertRecordStatus struct {
	// 停止压测任务成功与否
	// 注意：此字段可能返回 null，表示取不到有效值。
	AbortJob *uint64 `json:"AbortJob,omitnil" name:"AbortJob"`

	// 发送告警通知成功与否
	// 注意：此字段可能返回 null，表示取不到有效值。
	SendNotice *uint64 `json:"SendNotice,omitnil" name:"SendNotice"`
}

type Attributes struct {
	// 采用请求返回码
	// 注意：此字段可能返回 null，表示取不到有效值。
	Status *string `json:"Status,omitnil" name:"Status"`

	// 采样请求结果码
	// 注意：此字段可能返回 null，表示取不到有效值。
	Result *string `json:"Result,omitnil" name:"Result"`

	// 采样请求API
	// 注意：此字段可能返回 null，表示取不到有效值。
	Service *string `json:"Service,omitnil" name:"Service"`

	// 采样请求调用方法
	// 注意：此字段可能返回 null，表示取不到有效值。
	Method *string `json:"Method,omitnil" name:"Method"`

	// 采样请求延时时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	Duration *string `json:"Duration,omitnil" name:"Duration"`
}

type CheckSummary struct {
	// 检查点名字
	Name *string `json:"Name,omitnil" name:"Name"`

	// 检查点所在步骤名字
	Step *string `json:"Step,omitnil" name:"Step"`

	// 检查点成功次数
	SuccessCount *int64 `json:"SuccessCount,omitnil" name:"SuccessCount"`

	// 检查失败次数
	FailCount *int64 `json:"FailCount,omitnil" name:"FailCount"`

	// 错误比例
	ErrorRate *float64 `json:"ErrorRate,omitnil" name:"ErrorRate"`
}

type Concurrency struct {
	// 多阶段配置数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	Stages []*Stage `json:"Stages,omitnil" name:"Stages"`

	// 运行次数
	// 注意：此字段可能返回 null，表示取不到有效值。
	IterationCount *int64 `json:"IterationCount,omitnil" name:"IterationCount"`

	// 最大RPS
	// 注意：此字段可能返回 null，表示取不到有效值。
	MaxRequestsPerSecond *int64 `json:"MaxRequestsPerSecond,omitnil" name:"MaxRequestsPerSecond"`

	// 优雅终止任务的等待时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	GracefulStopSeconds *int64 `json:"GracefulStopSeconds,omitnil" name:"GracefulStopSeconds"`

	// 资源数
	// 注意：此字段可能返回 null，表示取不到有效值。
	Resources *int64 `json:"Resources,omitnil" name:"Resources"`
}

// Predefined struct for user
type CopyScenarioRequestParams struct {
	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景 ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`
}

type CopyScenarioRequest struct {
	*tchttp.BaseRequest
	
	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景 ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`
}

func (r *CopyScenarioRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CopyScenarioRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "ScenarioId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CopyScenarioRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CopyScenarioResponseParams struct {
	// 复制出的新场景 ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type CopyScenarioResponse struct {
	*tchttp.BaseResponse
	Response *CopyScenarioResponseParams `json:"Response"`
}

func (r *CopyScenarioResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CopyScenarioResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAlertChannelRequestParams struct {
	// Notice ID
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`

	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// AMP Consumer ID
	AMPConsumerId *string `json:"AMPConsumerId,omitnil" name:"AMPConsumerId"`
}

type CreateAlertChannelRequest struct {
	*tchttp.BaseRequest
	
	// Notice ID
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`

	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// AMP Consumer ID
	AMPConsumerId *string `json:"AMPConsumerId,omitnil" name:"AMPConsumerId"`
}

func (r *CreateAlertChannelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAlertChannelRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "NoticeId")
	delete(f, "ProjectId")
	delete(f, "AMPConsumerId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateAlertChannelRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateAlertChannelResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type CreateAlertChannelResponse struct {
	*tchttp.BaseResponse
	Response *CreateAlertChannelResponseParams `json:"Response"`
}

func (r *CreateAlertChannelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateAlertChannelResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateCronJobRequestParams struct {
	// 定时任务名字
	Name *string `json:"Name,omitnil" name:"Name"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 场景名称
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`

	// 执行频率类型，1:只执行一次; 2:日粒度; 3:周粒度; 4:高级
	FrequencyType *int64 `json:"FrequencyType,omitnil" name:"FrequencyType"`

	// cron表达式
	CronExpression *string `json:"CronExpression,omitnil" name:"CronExpression"`

	// 任务发起人
	JobOwner *string `json:"JobOwner,omitnil" name:"JobOwner"`

	// 结束时间
	EndTime *string `json:"EndTime,omitnil" name:"EndTime"`

	// Notice ID
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`

	// 备注
	Note *string `json:"Note,omitnil" name:"Note"`
}

type CreateCronJobRequest struct {
	*tchttp.BaseRequest
	
	// 定时任务名字
	Name *string `json:"Name,omitnil" name:"Name"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 场景名称
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`

	// 执行频率类型，1:只执行一次; 2:日粒度; 3:周粒度; 4:高级
	FrequencyType *int64 `json:"FrequencyType,omitnil" name:"FrequencyType"`

	// cron表达式
	CronExpression *string `json:"CronExpression,omitnil" name:"CronExpression"`

	// 任务发起人
	JobOwner *string `json:"JobOwner,omitnil" name:"JobOwner"`

	// 结束时间
	EndTime *string `json:"EndTime,omitnil" name:"EndTime"`

	// Notice ID
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`

	// 备注
	Note *string `json:"Note,omitnil" name:"Note"`
}

func (r *CreateCronJobRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateCronJobRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Name")
	delete(f, "ProjectId")
	delete(f, "ScenarioId")
	delete(f, "ScenarioName")
	delete(f, "FrequencyType")
	delete(f, "CronExpression")
	delete(f, "JobOwner")
	delete(f, "EndTime")
	delete(f, "NoticeId")
	delete(f, "Note")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateCronJobRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateCronJobResponseParams struct {
	// 定时任务ID
	CronJobId *string `json:"CronJobId,omitnil" name:"CronJobId"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type CreateCronJobResponse struct {
	*tchttp.BaseResponse
	Response *CreateCronJobResponseParams `json:"Response"`
}

func (r *CreateCronJobResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateCronJobResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateFileRequestParams struct {
	// 文件 ID
	FileId *string `json:"FileId,omitnil" name:"FileId"`

	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 文件种类，参数文件-1，协议文件-2，请求文件-3
	Kind *int64 `json:"Kind,omitnil" name:"Kind"`

	// 文件名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 文件大小
	Size *int64 `json:"Size,omitnil" name:"Size"`

	// 文件类型，文件夹-folder
	Type *string `json:"Type,omitnil" name:"Type"`

	// 行数
	LineCount *int64 `json:"LineCount,omitnil" name:"LineCount"`

	// 前几行数据
	HeadLines []*string `json:"HeadLines,omitnil" name:"HeadLines"`

	// 后几行数据
	TailLines []*string `json:"TailLines,omitnil" name:"TailLines"`

	// 表头是否在文件内
	HeaderInFile *bool `json:"HeaderInFile,omitnil" name:"HeaderInFile"`

	// 表头
	HeaderColumns []*string `json:"HeaderColumns,omitnil" name:"HeaderColumns"`

	// 文件夹中的文件
	FileInfos []*FileInfo `json:"FileInfos,omitnil" name:"FileInfos"`
}

type CreateFileRequest struct {
	*tchttp.BaseRequest
	
	// 文件 ID
	FileId *string `json:"FileId,omitnil" name:"FileId"`

	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 文件种类，参数文件-1，协议文件-2，请求文件-3
	Kind *int64 `json:"Kind,omitnil" name:"Kind"`

	// 文件名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 文件大小
	Size *int64 `json:"Size,omitnil" name:"Size"`

	// 文件类型，文件夹-folder
	Type *string `json:"Type,omitnil" name:"Type"`

	// 行数
	LineCount *int64 `json:"LineCount,omitnil" name:"LineCount"`

	// 前几行数据
	HeadLines []*string `json:"HeadLines,omitnil" name:"HeadLines"`

	// 后几行数据
	TailLines []*string `json:"TailLines,omitnil" name:"TailLines"`

	// 表头是否在文件内
	HeaderInFile *bool `json:"HeaderInFile,omitnil" name:"HeaderInFile"`

	// 表头
	HeaderColumns []*string `json:"HeaderColumns,omitnil" name:"HeaderColumns"`

	// 文件夹中的文件
	FileInfos []*FileInfo `json:"FileInfos,omitnil" name:"FileInfos"`
}

func (r *CreateFileRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateFileRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "FileId")
	delete(f, "ProjectId")
	delete(f, "Kind")
	delete(f, "Name")
	delete(f, "Size")
	delete(f, "Type")
	delete(f, "LineCount")
	delete(f, "HeadLines")
	delete(f, "TailLines")
	delete(f, "HeaderInFile")
	delete(f, "HeaderColumns")
	delete(f, "FileInfos")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateFileRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateFileResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type CreateFileResponse struct {
	*tchttp.BaseResponse
	Response *CreateFileResponseParams `json:"Response"`
}

func (r *CreateFileResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateFileResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateProjectRequestParams struct {
	// 项目名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 项目描述
	Description *string `json:"Description,omitnil" name:"Description"`

	// 标签数组
	Tags []*TagSpec `json:"Tags,omitnil" name:"Tags"`
}

type CreateProjectRequest struct {
	*tchttp.BaseRequest
	
	// 项目名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 项目描述
	Description *string `json:"Description,omitnil" name:"Description"`

	// 标签数组
	Tags []*TagSpec `json:"Tags,omitnil" name:"Tags"`
}

func (r *CreateProjectRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateProjectRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Name")
	delete(f, "Description")
	delete(f, "Tags")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateProjectRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateProjectResponseParams struct {
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type CreateProjectResponse struct {
	*tchttp.BaseResponse
	Response *CreateProjectResponseParams `json:"Response"`
}

func (r *CreateProjectResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateProjectResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateScenarioRequestParams struct {
	// 场景名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 压测引擎类型
	Type *string `json:"Type,omitnil" name:"Type"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景描述
	Description *string `json:"Description,omitnil" name:"Description"`

	// 施压配置
	Load *Load `json:"Load,omitnil" name:"Load"`

	// deprecated
	Configs []*string `json:"Configs,omitnil" name:"Configs"`

	// 测试数据集
	Datasets []*TestData `json:"Datasets,omitnil" name:"Datasets"`

	// deprecated
	Extensions []*string `json:"Extensions,omitnil" name:"Extensions"`

	// deprecated
	SLAId *string `json:"SLAId,omitnil" name:"SLAId"`

	// cron job ID
	CronId *string `json:"CronId,omitnil" name:"CronId"`

	// deprecated
	Scripts []*string `json:"Scripts,omitnil" name:"Scripts"`

	// 测试脚本文件信息
	TestScripts []*ScriptInfo `json:"TestScripts,omitnil" name:"TestScripts"`

	// 协议文件路径
	Protocols []*ProtocolInfo `json:"Protocols,omitnil" name:"Protocols"`

	// 请求文件路径
	RequestFiles []*FileInfo `json:"RequestFiles,omitnil" name:"RequestFiles"`

	// SLA 策略
	SLAPolicy *SLAPolicy `json:"SLAPolicy,omitnil" name:"SLAPolicy"`

	// 拓展包文件路径
	Plugins []*FileInfo `json:"Plugins,omitnil" name:"Plugins"`

	// 域名解析配置
	DomainNameConfig *DomainNameConfig `json:"DomainNameConfig,omitnil" name:"DomainNameConfig"`

	// 创建人名
	Owner *string `json:"Owner,omitnil" name:"Owner"`
}

type CreateScenarioRequest struct {
	*tchttp.BaseRequest
	
	// 场景名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 压测引擎类型
	Type *string `json:"Type,omitnil" name:"Type"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景描述
	Description *string `json:"Description,omitnil" name:"Description"`

	// 施压配置
	Load *Load `json:"Load,omitnil" name:"Load"`

	// deprecated
	Configs []*string `json:"Configs,omitnil" name:"Configs"`

	// 测试数据集
	Datasets []*TestData `json:"Datasets,omitnil" name:"Datasets"`

	// deprecated
	Extensions []*string `json:"Extensions,omitnil" name:"Extensions"`

	// deprecated
	SLAId *string `json:"SLAId,omitnil" name:"SLAId"`

	// cron job ID
	CronId *string `json:"CronId,omitnil" name:"CronId"`

	// deprecated
	Scripts []*string `json:"Scripts,omitnil" name:"Scripts"`

	// 测试脚本文件信息
	TestScripts []*ScriptInfo `json:"TestScripts,omitnil" name:"TestScripts"`

	// 协议文件路径
	Protocols []*ProtocolInfo `json:"Protocols,omitnil" name:"Protocols"`

	// 请求文件路径
	RequestFiles []*FileInfo `json:"RequestFiles,omitnil" name:"RequestFiles"`

	// SLA 策略
	SLAPolicy *SLAPolicy `json:"SLAPolicy,omitnil" name:"SLAPolicy"`

	// 拓展包文件路径
	Plugins []*FileInfo `json:"Plugins,omitnil" name:"Plugins"`

	// 域名解析配置
	DomainNameConfig *DomainNameConfig `json:"DomainNameConfig,omitnil" name:"DomainNameConfig"`

	// 创建人名
	Owner *string `json:"Owner,omitnil" name:"Owner"`
}

func (r *CreateScenarioRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateScenarioRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Name")
	delete(f, "Type")
	delete(f, "ProjectId")
	delete(f, "Description")
	delete(f, "Load")
	delete(f, "Configs")
	delete(f, "Datasets")
	delete(f, "Extensions")
	delete(f, "SLAId")
	delete(f, "CronId")
	delete(f, "Scripts")
	delete(f, "TestScripts")
	delete(f, "Protocols")
	delete(f, "RequestFiles")
	delete(f, "SLAPolicy")
	delete(f, "Plugins")
	delete(f, "DomainNameConfig")
	delete(f, "Owner")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateScenarioRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateScenarioResponseParams struct {
	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type CreateScenarioResponse struct {
	*tchttp.BaseResponse
	Response *CreateScenarioResponseParams `json:"Response"`
}

func (r *CreateScenarioResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateScenarioResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type Credentials struct {
	// 临时secret ID
	TmpSecretId *string `json:"TmpSecretId,omitnil" name:"TmpSecretId"`

	// 临时secret key
	TmpSecretKey *string `json:"TmpSecretKey,omitnil" name:"TmpSecretKey"`

	// 临时token
	Token *string `json:"Token,omitnil" name:"Token"`
}

type CronJob struct {
	// 定时任务ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	CronJobId *string `json:"CronJobId,omitnil" name:"CronJobId"`

	// 定时任务名字
	// 注意：此字段可能返回 null，表示取不到有效值。
	Name *string `json:"Name,omitnil" name:"Name"`

	// 项目ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 场景名称
	// 注意：此字段可能返回 null，表示取不到有效值。
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`

	// cron 表达式
	// 注意：此字段可能返回 null，表示取不到有效值。
	CronExpression *string `json:"CronExpression,omitnil" name:"CronExpression"`

	// 结束时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	EndTime *string `json:"EndTime,omitnil" name:"EndTime"`

	// 中止原因
	// 注意：此字段可能返回 null，表示取不到有效值。
	AbortReason *int64 `json:"AbortReason,omitnil" name:"AbortReason"`

	// 定时任务状态
	// 注意：此字段可能返回 null，表示取不到有效值。
	Status *int64 `json:"Status,omitnil" name:"Status"`

	// Notice ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`

	// 创建时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	CreatedAt *string `json:"CreatedAt,omitnil" name:"CreatedAt"`

	// 更新时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`

	// 执行频率类型
	// 注意：此字段可能返回 null，表示取不到有效值。
	FrequencyType *int64 `json:"FrequencyType,omitnil" name:"FrequencyType"`

	// 备注
	// 注意：此字段可能返回 null，表示取不到有效值。
	Note *string `json:"Note,omitnil" name:"Note"`

	// tom
	// 注意：此字段可能返回 null，表示取不到有效值。
	JobOwner *string `json:"JobOwner,omitnil" name:"JobOwner"`

	// App ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	AppId *int64 `json:"AppId,omitnil" name:"AppId"`

	// 主账号
	// 注意：此字段可能返回 null，表示取不到有效值。
	Uin *string `json:"Uin,omitnil" name:"Uin"`

	// 子账号
	// 注意：此字段可能返回 null，表示取不到有效值。
	SubAccountUin *string `json:"SubAccountUin,omitnil" name:"SubAccountUin"`
}

type CustomSample struct {
	// 指标名
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 聚合条件
	Aggregation *string `json:"Aggregation,omitnil" name:"Aggregation"`

	// 过滤条件
	// 注意：此字段可能返回 null，表示取不到有效值。
	Labels []*Label `json:"Labels,omitnil" name:"Labels"`

	// 查询值
	Value *float64 `json:"Value,omitnil" name:"Value"`

	// Time is the number of milliseconds since the epoch
	// // (1970-01-01 00:00 UTC) excluding leap seconds.
	Timestamp *int64 `json:"Timestamp,omitnil" name:"Timestamp"`

	// 指标对应的单位，当前单位有：s,bytes,bytes/s,reqs,reqs/s,checks,checks/s,iters,iters/s,VUs, %
	Unit *string `json:"Unit,omitnil" name:"Unit"`

	// 指标序列名字
	// 注意：此字段可能返回 null，表示取不到有效值。
	Name *string `json:"Name,omitnil" name:"Name"`
}

type CustomSampleMatrix struct {
	// 指标名字
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 聚合函数
	Aggregation *string `json:"Aggregation,omitnil" name:"Aggregation"`

	// 指标单位
	// 注意：此字段可能返回 null，表示取不到有效值。
	Unit *string `json:"Unit,omitnil" name:"Unit"`

	// 指标序列数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	Streams []*SampleStream `json:"Streams,omitnil" name:"Streams"`
}

type DNSConfig struct {
	// DNS IP 列表
	// 注意：此字段可能返回 null，表示取不到有效值。
	Nameservers []*string `json:"Nameservers,omitnil" name:"Nameservers"`
}

// Predefined struct for user
type DeleteAlertChannelRequestParams struct {
	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// Notice ID
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`
}

type DeleteAlertChannelRequest struct {
	*tchttp.BaseRequest
	
	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// Notice ID
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`
}

func (r *DeleteAlertChannelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteAlertChannelRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "NoticeId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteAlertChannelRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteAlertChannelResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DeleteAlertChannelResponse struct {
	*tchttp.BaseResponse
	Response *DeleteAlertChannelResponseParams `json:"Response"`
}

func (r *DeleteAlertChannelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteAlertChannelResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteCronJobsRequestParams struct {
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 定时任务ID数组
	CronJobIds []*string `json:"CronJobIds,omitnil" name:"CronJobIds"`
}

type DeleteCronJobsRequest struct {
	*tchttp.BaseRequest
	
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 定时任务ID数组
	CronJobIds []*string `json:"CronJobIds,omitnil" name:"CronJobIds"`
}

func (r *DeleteCronJobsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteCronJobsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "CronJobIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteCronJobsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteCronJobsResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DeleteCronJobsResponse struct {
	*tchttp.BaseResponse
	Response *DeleteCronJobsResponseParams `json:"Response"`
}

func (r *DeleteCronJobsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteCronJobsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteFilesRequestParams struct {
	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 文件 ID 数组
	FileIds []*string `json:"FileIds,omitnil" name:"FileIds"`
}

type DeleteFilesRequest struct {
	*tchttp.BaseRequest
	
	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 文件 ID 数组
	FileIds []*string `json:"FileIds,omitnil" name:"FileIds"`
}

func (r *DeleteFilesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteFilesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "FileIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteFilesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteFilesResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DeleteFilesResponse struct {
	*tchttp.BaseResponse
	Response *DeleteFilesResponseParams `json:"Response"`
}

func (r *DeleteFilesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteFilesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteJobsRequestParams struct {
	// 任务ID数组
	JobIds []*string `json:"JobIds,omitnil" name:"JobIds"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`
}

type DeleteJobsRequest struct {
	*tchttp.BaseRequest
	
	// 任务ID数组
	JobIds []*string `json:"JobIds,omitnil" name:"JobIds"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`
}

func (r *DeleteJobsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteJobsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobIds")
	delete(f, "ProjectId")
	delete(f, "ScenarioIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteJobsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteJobsResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DeleteJobsResponse struct {
	*tchttp.BaseResponse
	Response *DeleteJobsResponseParams `json:"Response"`
}

func (r *DeleteJobsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteJobsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteProjectsRequestParams struct {
	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 是否删除项目相关的场景。默认为否。
	DeleteScenarios *bool `json:"DeleteScenarios,omitnil" name:"DeleteScenarios"`

	// 是否删除项目相关的任务。默认为否。
	DeleteJobs *bool `json:"DeleteJobs,omitnil" name:"DeleteJobs"`
}

type DeleteProjectsRequest struct {
	*tchttp.BaseRequest
	
	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 是否删除项目相关的场景。默认为否。
	DeleteScenarios *bool `json:"DeleteScenarios,omitnil" name:"DeleteScenarios"`

	// 是否删除项目相关的任务。默认为否。
	DeleteJobs *bool `json:"DeleteJobs,omitnil" name:"DeleteJobs"`
}

func (r *DeleteProjectsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteProjectsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectIds")
	delete(f, "DeleteScenarios")
	delete(f, "DeleteJobs")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteProjectsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteProjectsResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DeleteProjectsResponse struct {
	*tchttp.BaseResponse
	Response *DeleteProjectsResponseParams `json:"Response"`
}

func (r *DeleteProjectsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteProjectsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteScenariosRequestParams struct {
	// 场景ID数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 是否删除场景相关的任务。默认为否。
	DeleteJobs *bool `json:"DeleteJobs,omitnil" name:"DeleteJobs"`
}

type DeleteScenariosRequest struct {
	*tchttp.BaseRequest
	
	// 场景ID数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 是否删除场景相关的任务。默认为否。
	DeleteJobs *bool `json:"DeleteJobs,omitnil" name:"DeleteJobs"`
}

func (r *DeleteScenariosRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteScenariosRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ScenarioIds")
	delete(f, "ProjectId")
	delete(f, "DeleteJobs")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteScenariosRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteScenariosResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DeleteScenariosResponse struct {
	*tchttp.BaseResponse
	Response *DeleteScenariosResponseParams `json:"Response"`
}

func (r *DeleteScenariosResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteScenariosResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeAlertChannelsRequestParams struct {
	// 项目 ID 列表
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 偏移量，默认为0
	Offset *uint64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量，默认为20，最大为100
	Limit *uint64 `json:"Limit,omitnil" name:"Limit"`

	// Notice ID 列表
	NoticeIds []*string `json:"NoticeIds,omitnil" name:"NoticeIds"`

	// 排序项
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 是否正序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`
}

type DescribeAlertChannelsRequest struct {
	*tchttp.BaseRequest
	
	// 项目 ID 列表
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 偏移量，默认为0
	Offset *uint64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量，默认为20，最大为100
	Limit *uint64 `json:"Limit,omitnil" name:"Limit"`

	// Notice ID 列表
	NoticeIds []*string `json:"NoticeIds,omitnil" name:"NoticeIds"`

	// 排序项
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 是否正序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`
}

func (r *DescribeAlertChannelsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeAlertChannelsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectIds")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "NoticeIds")
	delete(f, "OrderBy")
	delete(f, "Ascend")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeAlertChannelsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeAlertChannelsResponseParams struct {
	// 告警通知接收组列表
	// 注意：此字段可能返回 null，表示取不到有效值。
	AlertChannelSet []*AlertChannelRecord `json:"AlertChannelSet,omitnil" name:"AlertChannelSet"`

	// 告警通知接收组数目
	// 注意：此字段可能返回 null，表示取不到有效值。
	Total *uint64 `json:"Total,omitnil" name:"Total"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeAlertChannelsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeAlertChannelsResponseParams `json:"Response"`
}

func (r *DescribeAlertChannelsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeAlertChannelsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeAlertRecordsRequestParams struct {
	// 项目 ID 列表
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 场景 ID 列表
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`

	// 任务 ID 列表
	JobIds []*string `json:"JobIds,omitnil" name:"JobIds"`

	// 是否正序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`

	// 排序项
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 偏移量，默认为0
	Offset *uint64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量，默认为20，最大为100
	Limit *uint64 `json:"Limit,omitnil" name:"Limit"`

	// 按场景名筛选
	ScenarioNames []*string `json:"ScenarioNames,omitnil" name:"ScenarioNames"`
}

type DescribeAlertRecordsRequest struct {
	*tchttp.BaseRequest
	
	// 项目 ID 列表
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 场景 ID 列表
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`

	// 任务 ID 列表
	JobIds []*string `json:"JobIds,omitnil" name:"JobIds"`

	// 是否正序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`

	// 排序项
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 偏移量，默认为0
	Offset *uint64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量，默认为20，最大为100
	Limit *uint64 `json:"Limit,omitnil" name:"Limit"`

	// 按场景名筛选
	ScenarioNames []*string `json:"ScenarioNames,omitnil" name:"ScenarioNames"`
}

func (r *DescribeAlertRecordsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeAlertRecordsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectIds")
	delete(f, "ScenarioIds")
	delete(f, "JobIds")
	delete(f, "Ascend")
	delete(f, "OrderBy")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "ScenarioNames")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeAlertRecordsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeAlertRecordsResponseParams struct {
	// 告警历史
	// 注意：此字段可能返回 null，表示取不到有效值。
	AlertRecordSet []*AlertRecord `json:"AlertRecordSet,omitnil" name:"AlertRecordSet"`

	// 告警历史记录的总数
	// 注意：此字段可能返回 null，表示取不到有效值。
	Total *uint64 `json:"Total,omitnil" name:"Total"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeAlertRecordsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeAlertRecordsResponseParams `json:"Response"`
}

func (r *DescribeAlertRecordsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeAlertRecordsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeAvailableMetricsRequestParams struct {

}

type DescribeAvailableMetricsRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeAvailableMetricsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeAvailableMetricsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeAvailableMetricsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeAvailableMetricsResponseParams struct {
	// 系统支持的所有指标
	MetricSet []*MetricInfo `json:"MetricSet,omitnil" name:"MetricSet"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeAvailableMetricsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeAvailableMetricsResponseParams `json:"Response"`
}

func (r *DescribeAvailableMetricsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeAvailableMetricsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeCheckSummaryRequestParams struct {
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`
}

type DescribeCheckSummaryRequest struct {
	*tchttp.BaseRequest
	
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`
}

func (r *DescribeCheckSummaryRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeCheckSummaryRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ScenarioId")
	delete(f, "ProjectId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeCheckSummaryRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeCheckSummaryResponseParams struct {
	// 检查点汇总信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	CheckSummarySet []*CheckSummary `json:"CheckSummarySet,omitnil" name:"CheckSummarySet"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeCheckSummaryResponse struct {
	*tchttp.BaseResponse
	Response *DescribeCheckSummaryResponseParams `json:"Response"`
}

func (r *DescribeCheckSummaryResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeCheckSummaryResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeCronJobsRequestParams struct {
	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 偏移量，默认为0
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 定时任务ID数组
	CronJobIds []*string `json:"CronJobIds,omitnil" name:"CronJobIds"`

	// 定时任务名字，模糊查询
	CronJobName *string `json:"CronJobName,omitnil" name:"CronJobName"`

	// 定时任务状态数组
	CronJobStatus []*int64 `json:"CronJobStatus,omitnil" name:"CronJobStatus"`

	// 排序的列
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 是否正序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`
}

type DescribeCronJobsRequest struct {
	*tchttp.BaseRequest
	
	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 偏移量，默认为0
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 定时任务ID数组
	CronJobIds []*string `json:"CronJobIds,omitnil" name:"CronJobIds"`

	// 定时任务名字，模糊查询
	CronJobName *string `json:"CronJobName,omitnil" name:"CronJobName"`

	// 定时任务状态数组
	CronJobStatus []*int64 `json:"CronJobStatus,omitnil" name:"CronJobStatus"`

	// 排序的列
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 是否正序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`
}

func (r *DescribeCronJobsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeCronJobsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectIds")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "CronJobIds")
	delete(f, "CronJobName")
	delete(f, "CronJobStatus")
	delete(f, "OrderBy")
	delete(f, "Ascend")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeCronJobsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeCronJobsResponseParams struct {
	// 定时任务总数
	Total *int64 `json:"Total,omitnil" name:"Total"`

	// 定时任务信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	CronJobSet []*CronJob `json:"CronJobSet,omitnil" name:"CronJobSet"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeCronJobsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeCronJobsResponseParams `json:"Response"`
}

func (r *DescribeCronJobsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeCronJobsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeErrorSummaryRequestParams struct {
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 过滤参数
	Filters []*Filter `json:"Filters,omitnil" name:"Filters"`
}

type DescribeErrorSummaryRequest struct {
	*tchttp.BaseRequest
	
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 过滤参数
	Filters []*Filter `json:"Filters,omitnil" name:"Filters"`
}

func (r *DescribeErrorSummaryRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeErrorSummaryRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ScenarioId")
	delete(f, "ProjectId")
	delete(f, "Filters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeErrorSummaryRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeErrorSummaryResponseParams struct {
	// 错误汇总信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	ErrorSummarySet []*ErrorSummary `json:"ErrorSummarySet,omitnil" name:"ErrorSummarySet"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeErrorSummaryResponse struct {
	*tchttp.BaseResponse
	Response *DescribeErrorSummaryResponseParams `json:"Response"`
}

func (r *DescribeErrorSummaryResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeErrorSummaryResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeFilesRequestParams struct {
	// 项目 ID 数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 文件 ID 数组
	FileIds []*string `json:"FileIds,omitnil" name:"FileIds"`

	// 文件名
	FileName *string `json:"FileName,omitnil" name:"FileName"`

	// 偏移量，默认为 0
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量，默认为 20，最大为 100
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 文件种类，参数文件-1，协议文件-2，请求文件-3
	Kind *int64 `json:"Kind,omitnil" name:"Kind"`
}

type DescribeFilesRequest struct {
	*tchttp.BaseRequest
	
	// 项目 ID 数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 文件 ID 数组
	FileIds []*string `json:"FileIds,omitnil" name:"FileIds"`

	// 文件名
	FileName *string `json:"FileName,omitnil" name:"FileName"`

	// 偏移量，默认为 0
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量，默认为 20，最大为 100
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 文件种类，参数文件-1，协议文件-2，请求文件-3
	Kind *int64 `json:"Kind,omitnil" name:"Kind"`
}

func (r *DescribeFilesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeFilesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectIds")
	delete(f, "FileIds")
	delete(f, "FileName")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "Kind")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeFilesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeFilesResponseParams struct {
	// 文件列表
	// 注意：此字段可能返回 null，表示取不到有效值。
	FileSet []*File `json:"FileSet,omitnil" name:"FileSet"`

	// 文件总数
	// 注意：此字段可能返回 null，表示取不到有效值。
	Total *int64 `json:"Total,omitnil" name:"Total"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeFilesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeFilesResponseParams `json:"Response"`
}

func (r *DescribeFilesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeFilesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeJobsRequestParams struct {
	// 场景ID数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`

	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 分页起始位置
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 每页最大数目
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 任务ID数组
	JobIds []*string `json:"JobIds,omitnil" name:"JobIds"`

	// 按字段排序
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 升序/降序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`

	// 任务开始时间
	StartTime *string `json:"StartTime,omitnil" name:"StartTime"`

	// 任务结束时间
	EndTime *string `json:"EndTime,omitnil" name:"EndTime"`

	// 调试任务标记
	Debug *bool `json:"Debug,omitnil" name:"Debug"`

	// 任务的状态
	Status []*int64 `json:"Status,omitnil" name:"Status"`
}

type DescribeJobsRequest struct {
	*tchttp.BaseRequest
	
	// 场景ID数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`

	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 分页起始位置
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 每页最大数目
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 任务ID数组
	JobIds []*string `json:"JobIds,omitnil" name:"JobIds"`

	// 按字段排序
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 升序/降序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`

	// 任务开始时间
	StartTime *string `json:"StartTime,omitnil" name:"StartTime"`

	// 任务结束时间
	EndTime *string `json:"EndTime,omitnil" name:"EndTime"`

	// 调试任务标记
	Debug *bool `json:"Debug,omitnil" name:"Debug"`

	// 任务的状态
	Status []*int64 `json:"Status,omitnil" name:"Status"`
}

func (r *DescribeJobsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeJobsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ScenarioIds")
	delete(f, "ProjectIds")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "JobIds")
	delete(f, "OrderBy")
	delete(f, "Ascend")
	delete(f, "StartTime")
	delete(f, "EndTime")
	delete(f, "Debug")
	delete(f, "Status")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeJobsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeJobsResponseParams struct {
	// 任务列表
	// 注意：此字段可能返回 null，表示取不到有效值。
	JobSet []*Job `json:"JobSet,omitnil" name:"JobSet"`

	// 任务数量
	// 注意：此字段可能返回 null，表示取不到有效值。
	Total *int64 `json:"Total,omitnil" name:"Total"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeJobsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeJobsResponseParams `json:"Response"`
}

func (r *DescribeJobsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeJobsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeLabelValuesRequestParams struct {
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 指标名称
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 查询标签名称
	LabelName *string `json:"LabelName,omitnil" name:"LabelName"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`
}

type DescribeLabelValuesRequest struct {
	*tchttp.BaseRequest
	
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 指标名称
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 查询标签名称
	LabelName *string `json:"LabelName,omitnil" name:"LabelName"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`
}

func (r *DescribeLabelValuesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeLabelValuesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ScenarioId")
	delete(f, "Metric")
	delete(f, "LabelName")
	delete(f, "ProjectId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeLabelValuesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeLabelValuesResponseParams struct {
	// 标签值数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	LabelValueSet []*string `json:"LabelValueSet,omitnil" name:"LabelValueSet"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeLabelValuesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeLabelValuesResponseParams `json:"Response"`
}

func (r *DescribeLabelValuesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeLabelValuesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeMetricLabelWithValuesRequestParams struct {
	// job id
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// project id
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// scenario id
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`
}

type DescribeMetricLabelWithValuesRequest struct {
	*tchttp.BaseRequest
	
	// job id
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// project id
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// scenario id
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`
}

func (r *DescribeMetricLabelWithValuesRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeMetricLabelWithValuesRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ProjectId")
	delete(f, "ScenarioId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeMetricLabelWithValuesRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeMetricLabelWithValuesResponseParams struct {
	// 指标所有的label和values数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	MetricLabelWithValuesSet []*MetricLabelWithValues `json:"MetricLabelWithValuesSet,omitnil" name:"MetricLabelWithValuesSet"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeMetricLabelWithValuesResponse struct {
	*tchttp.BaseResponse
	Response *DescribeMetricLabelWithValuesResponseParams `json:"Response"`
}

func (r *DescribeMetricLabelWithValuesResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeMetricLabelWithValuesResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeNormalLogsRequestParams struct {
	// 压测项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 测试场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 压测任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 日志上下文，加载更多日志时使用，透传上次返回的Context值，获取后续的日志内容。过期时间1小时
	Context *string `json:"Context,omitnil" name:"Context"`

	// 日志开始时间
	From *string `json:"From,omitnil" name:"From"`

	// 日志结束时间
	To *string `json:"To,omitnil" name:"To"`

	// 日志级别，可取debug/info/error
	SeverityText *string `json:"SeverityText,omitnil" name:"SeverityText"`

	// 施压节点IP
	Instance *string `json:"Instance,omitnil" name:"Instance"`

	// 施压节点所在地域
	InstanceRegion *string `json:"InstanceRegion,omitnil" name:"InstanceRegion"`

	// 日志类型， console代表用户输出，engine代表引擎输出
	LogType *string `json:"LogType,omitnil" name:"LogType"`

	// 返回日志条数限制，最大100
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`
}

type DescribeNormalLogsRequest struct {
	*tchttp.BaseRequest
	
	// 压测项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 测试场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 压测任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 日志上下文，加载更多日志时使用，透传上次返回的Context值，获取后续的日志内容。过期时间1小时
	Context *string `json:"Context,omitnil" name:"Context"`

	// 日志开始时间
	From *string `json:"From,omitnil" name:"From"`

	// 日志结束时间
	To *string `json:"To,omitnil" name:"To"`

	// 日志级别，可取debug/info/error
	SeverityText *string `json:"SeverityText,omitnil" name:"SeverityText"`

	// 施压节点IP
	Instance *string `json:"Instance,omitnil" name:"Instance"`

	// 施压节点所在地域
	InstanceRegion *string `json:"InstanceRegion,omitnil" name:"InstanceRegion"`

	// 日志类型， console代表用户输出，engine代表引擎输出
	LogType *string `json:"LogType,omitnil" name:"LogType"`

	// 返回日志条数限制，最大100
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`
}

func (r *DescribeNormalLogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeNormalLogsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "ScenarioId")
	delete(f, "JobId")
	delete(f, "Context")
	delete(f, "From")
	delete(f, "To")
	delete(f, "SeverityText")
	delete(f, "Instance")
	delete(f, "InstanceRegion")
	delete(f, "LogType")
	delete(f, "Limit")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeNormalLogsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeNormalLogsResponseParams struct {
	// 日志上下文，加载更多日志时使用，透传上次返回的Context值，获取后续的日志内容。过期时间1小时
	// 注意：此字段可能返回 null，表示取不到有效值。
	Context *string `json:"Context,omitnil" name:"Context"`

	// 日志数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	NormalLogs []*NormalLog `json:"NormalLogs,omitnil" name:"NormalLogs"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeNormalLogsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeNormalLogsResponseParams `json:"Response"`
}

func (r *DescribeNormalLogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeNormalLogsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeProjectsRequestParams struct {
	// 分页offset
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 每页limit
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 项目名
	ProjectName *string `json:"ProjectName,omitnil" name:"ProjectName"`

	// 按字段排序
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 升序/降序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`

	// 标签数组
	TagFilters []*TagSpec `json:"TagFilters,omitnil" name:"TagFilters"`
}

type DescribeProjectsRequest struct {
	*tchttp.BaseRequest
	
	// 分页offset
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 每页limit
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 项目名
	ProjectName *string `json:"ProjectName,omitnil" name:"ProjectName"`

	// 按字段排序
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 升序/降序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`

	// 标签数组
	TagFilters []*TagSpec `json:"TagFilters,omitnil" name:"TagFilters"`
}

func (r *DescribeProjectsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeProjectsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "ProjectIds")
	delete(f, "ProjectName")
	delete(f, "OrderBy")
	delete(f, "Ascend")
	delete(f, "TagFilters")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeProjectsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeProjectsResponseParams struct {
	// 项目数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	ProjectSet []*Project `json:"ProjectSet,omitnil" name:"ProjectSet"`

	// 项目数量
	Total *int64 `json:"Total,omitnil" name:"Total"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeProjectsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeProjectsResponseParams `json:"Response"`
}

func (r *DescribeProjectsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeProjectsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRegionsRequestParams struct {

}

type DescribeRegionsRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeRegionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRegionsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRegionsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRegionsResponseParams struct {
	// 地域数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	RegionSet []*RegionDetail `json:"RegionSet,omitnil" name:"RegionSet"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeRegionsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRegionsResponseParams `json:"Response"`
}

func (r *DescribeRegionsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRegionsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRequestSummaryRequestParams struct {
	// 压测任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 压测场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 压测项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`
}

type DescribeRequestSummaryRequest struct {
	*tchttp.BaseRequest
	
	// 压测任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 压测场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 压测项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`
}

func (r *DescribeRequestSummaryRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRequestSummaryRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ScenarioId")
	delete(f, "ProjectId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeRequestSummaryRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeRequestSummaryResponseParams struct {
	// 请求汇总信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	RequestSummarySet []*RequestSummary `json:"RequestSummarySet,omitnil" name:"RequestSummarySet"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeRequestSummaryResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRequestSummaryResponseParams `json:"Response"`
}

func (r *DescribeRequestSummaryResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeRequestSummaryResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSampleBatchQueryRequestParams struct {
	// job id
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 场景id
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 查询指标数组
	Queries []*InternalMetricQuery `json:"Queries,omitnil" name:"Queries"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`
}

type DescribeSampleBatchQueryRequest struct {
	*tchttp.BaseRequest
	
	// job id
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 场景id
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 查询指标数组
	Queries []*InternalMetricQuery `json:"Queries,omitnil" name:"Queries"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`
}

func (r *DescribeSampleBatchQueryRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSampleBatchQueryRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ScenarioId")
	delete(f, "Queries")
	delete(f, "ProjectId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeSampleBatchQueryRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSampleBatchQueryResponseParams struct {
	// 返回指标内容
	// 注意：此字段可能返回 null，表示取不到有效值。
	MetricSampleSet []*CustomSample `json:"MetricSampleSet,omitnil" name:"MetricSampleSet"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeSampleBatchQueryResponse struct {
	*tchttp.BaseResponse
	Response *DescribeSampleBatchQueryResponseParams `json:"Response"`
}

func (r *DescribeSampleBatchQueryResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSampleBatchQueryResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSampleLogsRequestParams struct {
	// 测试项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 测试场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 测试任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 加载更多日志时使用，透传上次返回的Context值，获取后续的日志内容。过期时间1小时
	Context *string `json:"Context,omitnil" name:"Context"`

	// 日志开始时间
	From *string `json:"From,omitnil" name:"From"`

	// 日志结束时间
	To *string `json:"To,omitnil" name:"To"`

	// 日志级别debug,info,error
	SeverityText *string `json:"SeverityText,omitnil" name:"SeverityText"`

	// ap-shanghai, ap-guangzhou
	InstanceRegion *string `json:"InstanceRegion,omitnil" name:"InstanceRegion"`

	// 施压引擎节点IP
	Instance *string `json:"Instance,omitnil" name:"Instance"`

	// request 代表采样日志,可为不填
	LogType *string `json:"LogType,omitnil" name:"LogType"`

	// 返回日志条数，最大100
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 采样日志响应时间范围
	ReactionTimeRange *ReactionTimeRange `json:"ReactionTimeRange,omitnil" name:"ReactionTimeRange"`

	// 采样请求状态码
	Status *string `json:"Status,omitnil" name:"Status"`

	// 采样请求结果码
	Result *string `json:"Result,omitnil" name:"Result"`

	// 采样请求方法
	Method *string `json:"Method,omitnil" name:"Method"`

	// 采样服务API
	Service *string `json:"Service,omitnil" name:"Service"`
}

type DescribeSampleLogsRequest struct {
	*tchttp.BaseRequest
	
	// 测试项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 测试场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 测试任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 加载更多日志时使用，透传上次返回的Context值，获取后续的日志内容。过期时间1小时
	Context *string `json:"Context,omitnil" name:"Context"`

	// 日志开始时间
	From *string `json:"From,omitnil" name:"From"`

	// 日志结束时间
	To *string `json:"To,omitnil" name:"To"`

	// 日志级别debug,info,error
	SeverityText *string `json:"SeverityText,omitnil" name:"SeverityText"`

	// ap-shanghai, ap-guangzhou
	InstanceRegion *string `json:"InstanceRegion,omitnil" name:"InstanceRegion"`

	// 施压引擎节点IP
	Instance *string `json:"Instance,omitnil" name:"Instance"`

	// request 代表采样日志,可为不填
	LogType *string `json:"LogType,omitnil" name:"LogType"`

	// 返回日志条数，最大100
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 采样日志响应时间范围
	ReactionTimeRange *ReactionTimeRange `json:"ReactionTimeRange,omitnil" name:"ReactionTimeRange"`

	// 采样请求状态码
	Status *string `json:"Status,omitnil" name:"Status"`

	// 采样请求结果码
	Result *string `json:"Result,omitnil" name:"Result"`

	// 采样请求方法
	Method *string `json:"Method,omitnil" name:"Method"`

	// 采样服务API
	Service *string `json:"Service,omitnil" name:"Service"`
}

func (r *DescribeSampleLogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSampleLogsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "ScenarioId")
	delete(f, "JobId")
	delete(f, "Context")
	delete(f, "From")
	delete(f, "To")
	delete(f, "SeverityText")
	delete(f, "InstanceRegion")
	delete(f, "Instance")
	delete(f, "LogType")
	delete(f, "Limit")
	delete(f, "ReactionTimeRange")
	delete(f, "Status")
	delete(f, "Result")
	delete(f, "Method")
	delete(f, "Service")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeSampleLogsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSampleLogsResponseParams struct {
	// 日志总数
	Total *int64 `json:"Total,omitnil" name:"Total"`

	// 日志上下文，加载更多日志时使用，透传上次返回的Context值，获取后续的日志内容。过期时间1小时
	// 注意：此字段可能返回 null，表示取不到有效值。
	Context *string `json:"Context,omitnil" name:"Context"`

	// 采样日志数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	SampleLogs []*SampleLog `json:"SampleLogs,omitnil" name:"SampleLogs"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeSampleLogsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeSampleLogsResponseParams `json:"Response"`
}

func (r *DescribeSampleLogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSampleLogsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSampleMatrixBatchQueryRequestParams struct {
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 查询语句
	Queries []*InternalMetricQuery `json:"Queries,omitnil" name:"Queries"`
}

type DescribeSampleMatrixBatchQueryRequest struct {
	*tchttp.BaseRequest
	
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 查询语句
	Queries []*InternalMetricQuery `json:"Queries,omitnil" name:"Queries"`
}

func (r *DescribeSampleMatrixBatchQueryRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSampleMatrixBatchQueryRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ProjectId")
	delete(f, "ScenarioId")
	delete(f, "Queries")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeSampleMatrixBatchQueryRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSampleMatrixBatchQueryResponseParams struct {
	// 批量指标矩阵
	// 注意：此字段可能返回 null，表示取不到有效值。
	MetricSampleMatrixSet []*CustomSampleMatrix `json:"MetricSampleMatrixSet,omitnil" name:"MetricSampleMatrixSet"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeSampleMatrixBatchQueryResponse struct {
	*tchttp.BaseResponse
	Response *DescribeSampleMatrixBatchQueryResponseParams `json:"Response"`
}

func (r *DescribeSampleMatrixBatchQueryResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSampleMatrixBatchQueryResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSampleMatrixQueryRequestParams struct {
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 指标名字
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 聚合函数
	Aggregation *string `json:"Aggregation,omitnil" name:"Aggregation"`

	// 指标过滤
	Filters []*Filter `json:"Filters,omitnil" name:"Filters"`

	// 分组
	GroupBy []*string `json:"GroupBy,omitnil" name:"GroupBy"`
}

type DescribeSampleMatrixQueryRequest struct {
	*tchttp.BaseRequest
	
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 指标名字
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 聚合函数
	Aggregation *string `json:"Aggregation,omitnil" name:"Aggregation"`

	// 指标过滤
	Filters []*Filter `json:"Filters,omitnil" name:"Filters"`

	// 分组
	GroupBy []*string `json:"GroupBy,omitnil" name:"GroupBy"`
}

func (r *DescribeSampleMatrixQueryRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSampleMatrixQueryRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ProjectId")
	delete(f, "ScenarioId")
	delete(f, "Metric")
	delete(f, "Aggregation")
	delete(f, "Filters")
	delete(f, "GroupBy")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeSampleMatrixQueryRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSampleMatrixQueryResponseParams struct {
	// 指标矩阵
	// 注意：此字段可能返回 null，表示取不到有效值。
	MetricSampleMatrix *CustomSampleMatrix `json:"MetricSampleMatrix,omitnil" name:"MetricSampleMatrix"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeSampleMatrixQueryResponse struct {
	*tchttp.BaseResponse
	Response *DescribeSampleMatrixQueryResponseParams `json:"Response"`
}

func (r *DescribeSampleMatrixQueryResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSampleMatrixQueryResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSampleQueryRequestParams struct {
	// job id
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 场景Id
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 指标名
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 聚合条件
	Aggregation *string `json:"Aggregation,omitnil" name:"Aggregation"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 过滤条件
	Labels []*Label `json:"Labels,omitnil" name:"Labels"`
}

type DescribeSampleQueryRequest struct {
	*tchttp.BaseRequest
	
	// job id
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 场景Id
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 指标名
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 聚合条件
	Aggregation *string `json:"Aggregation,omitnil" name:"Aggregation"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 过滤条件
	Labels []*Label `json:"Labels,omitnil" name:"Labels"`
}

func (r *DescribeSampleQueryRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSampleQueryRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ScenarioId")
	delete(f, "Metric")
	delete(f, "Aggregation")
	delete(f, "ProjectId")
	delete(f, "Labels")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeSampleQueryRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeSampleQueryResponseParams struct {
	// 返回指标内容
	// 注意：此字段可能返回 null，表示取不到有效值。
	MetricSample *CustomSample `json:"MetricSample,omitnil" name:"MetricSample"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeSampleQueryResponse struct {
	*tchttp.BaseResponse
	Response *DescribeSampleQueryResponseParams `json:"Response"`
}

func (r *DescribeSampleQueryResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeSampleQueryResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeScenarioWithJobsRequestParams struct {
	// 偏移量，默认为0
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量，默认为20，最大为100
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 场景ID数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`

	// 场景名
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`

	// 场景状态数组
	ScenarioStatus *int64 `json:"ScenarioStatus,omitnil" name:"ScenarioStatus"`

	// 排序的列
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 是否正序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`

	// job相关参数
	ScenarioRelatedJobsParams *ScenarioRelatedJobsParams `json:"ScenarioRelatedJobsParams,omitnil" name:"ScenarioRelatedJobsParams"`

	// 是否需要返回场景的脚本内容
	IgnoreScript *bool `json:"IgnoreScript,omitnil" name:"IgnoreScript"`

	// 是否需要返回测试数据文件信息
	IgnoreDataset *bool `json:"IgnoreDataset,omitnil" name:"IgnoreDataset"`

	// 场景类型，如pts-http, pts-js, pts-trpc, pts-jmeter	
	ScenarioType *string `json:"ScenarioType,omitnil" name:"ScenarioType"`

	// 创建人员
	Owner *string `json:"Owner,omitnil" name:"Owner"`
}

type DescribeScenarioWithJobsRequest struct {
	*tchttp.BaseRequest
	
	// 偏移量，默认为0
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量，默认为20，最大为100
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 场景ID数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`

	// 场景名
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`

	// 场景状态数组
	ScenarioStatus *int64 `json:"ScenarioStatus,omitnil" name:"ScenarioStatus"`

	// 排序的列
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 是否正序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`

	// job相关参数
	ScenarioRelatedJobsParams *ScenarioRelatedJobsParams `json:"ScenarioRelatedJobsParams,omitnil" name:"ScenarioRelatedJobsParams"`

	// 是否需要返回场景的脚本内容
	IgnoreScript *bool `json:"IgnoreScript,omitnil" name:"IgnoreScript"`

	// 是否需要返回测试数据文件信息
	IgnoreDataset *bool `json:"IgnoreDataset,omitnil" name:"IgnoreDataset"`

	// 场景类型，如pts-http, pts-js, pts-trpc, pts-jmeter	
	ScenarioType *string `json:"ScenarioType,omitnil" name:"ScenarioType"`

	// 创建人员
	Owner *string `json:"Owner,omitnil" name:"Owner"`
}

func (r *DescribeScenarioWithJobsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeScenarioWithJobsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "ProjectIds")
	delete(f, "ScenarioIds")
	delete(f, "ScenarioName")
	delete(f, "ScenarioStatus")
	delete(f, "OrderBy")
	delete(f, "Ascend")
	delete(f, "ScenarioRelatedJobsParams")
	delete(f, "IgnoreScript")
	delete(f, "IgnoreDataset")
	delete(f, "ScenarioType")
	delete(f, "Owner")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeScenarioWithJobsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeScenarioWithJobsResponseParams struct {
	// 场景配置以及附带的job内容
	// 注意：此字段可能返回 null，表示取不到有效值。
	ScenarioWithJobsSet []*ScenarioWithJobs `json:"ScenarioWithJobsSet,omitnil" name:"ScenarioWithJobsSet"`

	// 场景总数
	Total *int64 `json:"Total,omitnil" name:"Total"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeScenarioWithJobsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeScenarioWithJobsResponseParams `json:"Response"`
}

func (r *DescribeScenarioWithJobsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeScenarioWithJobsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeScenariosRequestParams struct {
	// 场景ID数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`

	// 场景名
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`

	// 场景状态数组
	ScenarioStatus []*int64 `json:"ScenarioStatus,omitnil" name:"ScenarioStatus"`

	// 偏移量，默认为0
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量，默认为20，最大为100
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 排序的列
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 是否正序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`

	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 场景类型
	ScenarioType *string `json:"ScenarioType,omitnil" name:"ScenarioType"`
}

type DescribeScenariosRequest struct {
	*tchttp.BaseRequest
	
	// 场景ID数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`

	// 场景名
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`

	// 场景状态数组
	ScenarioStatus []*int64 `json:"ScenarioStatus,omitnil" name:"ScenarioStatus"`

	// 偏移量，默认为0
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 返回数量，默认为20，最大为100
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 排序的列
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 是否正序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`

	// 项目ID数组
	ProjectIds []*string `json:"ProjectIds,omitnil" name:"ProjectIds"`

	// 场景类型
	ScenarioType *string `json:"ScenarioType,omitnil" name:"ScenarioType"`
}

func (r *DescribeScenariosRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeScenariosRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ScenarioIds")
	delete(f, "ScenarioName")
	delete(f, "ScenarioStatus")
	delete(f, "Offset")
	delete(f, "Limit")
	delete(f, "OrderBy")
	delete(f, "Ascend")
	delete(f, "ProjectIds")
	delete(f, "ScenarioType")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeScenariosRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeScenariosResponseParams struct {
	// 场景列表
	// 注意：此字段可能返回 null，表示取不到有效值。
	ScenarioSet []*Scenario `json:"ScenarioSet,omitnil" name:"ScenarioSet"`

	// 场景总数
	Total *int64 `json:"Total,omitnil" name:"Total"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type DescribeScenariosResponse struct {
	*tchttp.BaseResponse
	Response *DescribeScenariosResponseParams `json:"Response"`
}

func (r *DescribeScenariosResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeScenariosResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DomainNameConfig struct {
	// 域名绑定配置
	// 注意：此字段可能返回 null，表示取不到有效值。
	HostAliases []*HostAlias `json:"HostAliases,omitnil" name:"HostAliases"`

	// DNS 配置
	// 注意：此字段可能返回 null，表示取不到有效值。
	DNSConfig *DNSConfig `json:"DNSConfig,omitnil" name:"DNSConfig"`
}

type ErrorSummary struct {
	// 状态码
	// 注意：此字段可能返回 null，表示取不到有效值。
	Status *string `json:"Status,omitnil" name:"Status"`

	// 结果码
	// 注意：此字段可能返回 null，表示取不到有效值。
	Result *string `json:"Result,omitnil" name:"Result"`

	// 错误出现次数
	// 注意：此字段可能返回 null，表示取不到有效值。
	Count *int64 `json:"Count,omitnil" name:"Count"`

	// 错误率
	// 注意：此字段可能返回 null，表示取不到有效值。
	Rate *float64 `json:"Rate,omitnil" name:"Rate"`

	// 错误信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	Message *string `json:"Message,omitnil" name:"Message"`

	// 请求协议类型
	// 注意：此字段可能返回 null，表示取不到有效值。
	Proto *string `json:"Proto,omitnil" name:"Proto"`
}

type File struct {
	// 文件 ID
	FileId *string `json:"FileId,omitnil" name:"FileId"`

	// 文件种类，参数文件-1，协议文件-2，请求文件-3
	Kind *int64 `json:"Kind,omitnil" name:"Kind"`

	// 文件名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 文件字节数
	Size *int64 `json:"Size,omitnil" name:"Size"`

	// 文件类型
	Type *string `json:"Type,omitnil" name:"Type"`

	// 更新时间
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`

	// 文件行数
	// 注意：此字段可能返回 null，表示取不到有效值。
	LineCount *int64 `json:"LineCount,omitnil" name:"LineCount"`

	// 头部数据行
	// 注意：此字段可能返回 null，表示取不到有效值。
	HeadLines []*string `json:"HeadLines,omitnil" name:"HeadLines"`

	// 尾部数据行
	// 注意：此字段可能返回 null，表示取不到有效值。
	TailLines []*string `json:"TailLines,omitnil" name:"TailLines"`

	// 首行是否为参数名
	// 注意：此字段可能返回 null，表示取不到有效值。
	HeaderInFile *bool `json:"HeaderInFile,omitnil" name:"HeaderInFile"`

	// 参数名数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	HeaderColumns []*string `json:"HeaderColumns,omitnil" name:"HeaderColumns"`

	// 文件夹中的文件
	// 注意：此字段可能返回 null，表示取不到有效值。
	FileInfos []*FileInfo `json:"FileInfos,omitnil" name:"FileInfos"`

	// 关联场景
	// 注意：此字段可能返回 null，表示取不到有效值。
	ScenarioSet []*Scenario `json:"ScenarioSet,omitnil" name:"ScenarioSet"`

	// 文件状态
	// 注意：此字段可能返回 null，表示取不到有效值。
	Status *int64 `json:"Status,omitnil" name:"Status"`

	// 创建时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	CreatedAt *string `json:"CreatedAt,omitnil" name:"CreatedAt"`

	// 项目 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 此字段不再使用
	// 注意：此字段可能返回 null，表示取不到有效值。
	AppID *int64 `json:"AppID,omitnil" name:"AppID"`

	// 用户主账号
	// 注意：此字段可能返回 null，表示取不到有效值。
	Uin *string `json:"Uin,omitnil" name:"Uin"`

	// 用户子账号
	// 注意：此字段可能返回 null，表示取不到有效值。
	SubAccountUin *string `json:"SubAccountUin,omitnil" name:"SubAccountUin"`

	// 用户账号的 App ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	AppId *int64 `json:"AppId,omitnil" name:"AppId"`
}

type FileInfo struct {
	// 文件名
	// 注意：此字段可能返回 null，表示取不到有效值。
	Name *string `json:"Name,omitnil" name:"Name"`

	// 文件大小
	// 注意：此字段可能返回 null，表示取不到有效值。
	Size *int64 `json:"Size,omitnil" name:"Size"`

	// 文件类型
	// 注意：此字段可能返回 null，表示取不到有效值。
	Type *string `json:"Type,omitnil" name:"Type"`

	// 更新时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`

	// 文件 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	FileId *string `json:"FileId,omitnil" name:"FileId"`
}

type Filter struct {
	// 等于：0，不等于：1
	Operator *int64 `json:"Operator,omitnil" name:"Operator"`

	// 标签名，可选值包括：
	// 1. method，请求方法名；
	// 2. proto：协议名；
	// 3. service：服务名；
	// 4. status：响应状态码；
	// 5. result：响应详情；
	// 6. check：检查名。
	LabelName *string `json:"LabelName,omitnil" name:"LabelName"`

	// 标签值：
	// 1. method：请求方法名，以 http 协议为例，method 为 GET、POST、PUT 等；
	// 2. proto：协议名，以 http 协议为例，proto 为 HTTP/1.1、HTTP/2 等；
	// 3. service：服务名，以 http 协议为例，service 为请求 url，如 http://httpbin.org/get 等；
	// 4. status：响应状态码，以 http 协议为例，状态码包括 200、404、500 等；
	// 5. result：响应详情，通过 result 判断请求成功或失败；请求正常，result 标签值为 ok；请求失败，result 标签携带错误码和描述；
	// 6. check：检查名，标签值为用户设置的检查点名称。
	LabelValue *string `json:"LabelValue,omitnil" name:"LabelValue"`
}

// Predefined struct for user
type GenerateTmpKeyRequestParams struct {
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`
}

type GenerateTmpKeyRequest struct {
	*tchttp.BaseRequest
	
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`
}

func (r *GenerateTmpKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *GenerateTmpKeyRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "ScenarioId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "GenerateTmpKeyRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type GenerateTmpKeyResponseParams struct {
	// 临时访问凭证获取时刻的时间戳（单位秒）
	StartTime *int64 `json:"StartTime,omitnil" name:"StartTime"`

	// 临时访问凭证超时 时刻的时间戳（单位秒）
	ExpiredTime *int64 `json:"ExpiredTime,omitnil" name:"ExpiredTime"`

	// 临时访问凭证
	Credentials *Credentials `json:"Credentials,omitnil" name:"Credentials"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type GenerateTmpKeyResponse struct {
	*tchttp.BaseResponse
	Response *GenerateTmpKeyResponseParams `json:"Response"`
}

func (r *GenerateTmpKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *GenerateTmpKeyResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type GeoRegionsLoadItem struct {
	// 地域ID
	RegionId *int64 `json:"RegionId,omitnil" name:"RegionId"`

	// 地域
	Region *string `json:"Region,omitnil" name:"Region"`

	// 百分比
	Percentage *int64 `json:"Percentage,omitnil" name:"Percentage"`
}

type HostAlias struct {
	// 需绑定的域名列表
	// 注意：此字段可能返回 null，表示取不到有效值。
	HostNames []*string `json:"HostNames,omitnil" name:"HostNames"`

	// 需绑定的 IP 地址
	// 注意：此字段可能返回 null，表示取不到有效值。
	IP *string `json:"IP,omitnil" name:"IP"`
}

type InternalMetricQuery struct {
	// 指标名
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 聚合函数
	Aggregation *string `json:"Aggregation,omitnil" name:"Aggregation"`

	// deprecated, 请使用Filters
	Labels []*Label `json:"Labels,omitnil" name:"Labels"`

	// 指标过滤
	Filters []*Filter `json:"Filters,omitnil" name:"Filters"`

	// 指标分组
	GroupBy []*string `json:"GroupBy,omitnil" name:"GroupBy"`
}

type Job struct {
	// 任务的JobID
	// 注意：此字段可能返回 null，表示取不到有效值。
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 任务的场景ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 任务的施压配置
	// 注意：此字段可能返回 null，表示取不到有效值。
	Load *Load `json:"Load,omitnil" name:"Load"`

	// 此字段不再使用
	// 注意：此字段可能返回 null，表示取不到有效值。
	Configs []*string `json:"Configs,omitnil" name:"Configs"`

	// 任务的数据集文件
	// 注意：此字段可能返回 null，表示取不到有效值。
	Datasets []*TestData `json:"Datasets,omitnil" name:"Datasets"`

	// 此字段不再使用
	// 注意：此字段可能返回 null，表示取不到有效值。
	Extensions []*string `json:"Extensions,omitnil" name:"Extensions"`

	// 任务的运行状态, JobUnknown: 0,JobCreated:1,JobPending:2, JobPreparing:3,JobSelectClustering:4,JobCreateTasking:5,JobSyncTasking:6
	// JobRunning:11,JobFinished:12,JobPrepareException:13,JobFinishException:14,JobAborting:15,JobAborted:16,JobAbortException:17,JobDeleted:18,
	// JobSelectClusterException:19,JobCreateTaskException:20,JobSyncTaskException:21
	// 注意：此字段可能返回 null，表示取不到有效值。
	Status *int64 `json:"Status,omitnil" name:"Status"`

	// 任务的开始时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	StartTime *string `json:"StartTime,omitnil" name:"StartTime"`

	// 任务的结束时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	EndTime *string `json:"EndTime,omitnil" name:"EndTime"`

	// 任务的最大VU数
	// 注意：此字段可能返回 null，表示取不到有效值。
	MaxVirtualUserCount *int64 `json:"MaxVirtualUserCount,omitnil" name:"MaxVirtualUserCount"`

	// 任务的备注信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	Note *string `json:"Note,omitnil" name:"Note"`

	// 错误率百分比
	// 注意：此字段可能返回 null，表示取不到有效值。
	ErrorRate *float64 `json:"ErrorRate,omitnil" name:"ErrorRate"`

	// 任务发起人
	// 注意：此字段可能返回 null，表示取不到有效值。
	JobOwner *string `json:"JobOwner,omitnil" name:"JobOwner"`

	// 此字段不再使用
	// 注意：此字段可能返回 null，表示取不到有效值。
	LoadSources *LoadSource `json:"LoadSources,omitnil" name:"LoadSources"`

	// 任务时长
	// 注意：此字段可能返回 null，表示取不到有效值。
	Duration *int64 `json:"Duration,omitnil" name:"Duration"`

	// 最大每秒请求数
	// 注意：此字段可能返回 null，表示取不到有效值。
	MaxRequestsPerSecond *int64 `json:"MaxRequestsPerSecond,omitnil" name:"MaxRequestsPerSecond"`

	// 总请求数
	// 注意：此字段可能返回 null，表示取不到有效值。
	RequestTotal *float64 `json:"RequestTotal,omitnil" name:"RequestTotal"`

	// 平均每秒请求数
	// 注意：此字段可能返回 null，表示取不到有效值。
	RequestsPerSecond *float64 `json:"RequestsPerSecond,omitnil" name:"RequestsPerSecond"`

	// 平均响应时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResponseTimeAverage *float64 `json:"ResponseTimeAverage,omitnil" name:"ResponseTimeAverage"`

	// 响应时间第99百分位
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResponseTimeP99 *float64 `json:"ResponseTimeP99,omitnil" name:"ResponseTimeP99"`

	// 响应时间第95百分位
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResponseTimeP95 *float64 `json:"ResponseTimeP95,omitnil" name:"ResponseTimeP95"`

	// 响应时间第90百分位
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResponseTimeP90 *float64 `json:"ResponseTimeP90,omitnil" name:"ResponseTimeP90"`

	// 此字段不再使用
	// 注意：此字段可能返回 null，表示取不到有效值。
	Scripts []*string `json:"Scripts,omitnil" name:"Scripts"`

	// 最大响应时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResponseTimeMax *float64 `json:"ResponseTimeMax,omitnil" name:"ResponseTimeMax"`

	// 最小响应时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResponseTimeMin *float64 `json:"ResponseTimeMin,omitnil" name:"ResponseTimeMin"`

	// 发压host信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	LoadSourceInfos []*LoadSource `json:"LoadSourceInfos,omitnil" name:"LoadSourceInfos"`

	// 测试脚本信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	TestScripts []*ScriptInfo `json:"TestScripts,omitnil" name:"TestScripts"`

	// 协议脚本信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	Protocols []*ProtocolInfo `json:"Protocols,omitnil" name:"Protocols"`

	// 请求文件信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	RequestFiles []*FileInfo `json:"RequestFiles,omitnil" name:"RequestFiles"`

	// 拓展包文件信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	Plugins []*FileInfo `json:"Plugins,omitnil" name:"Plugins"`

	// 定时任务ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	CronId *string `json:"CronId,omitnil" name:"CronId"`

	// 场景类型
	// 注意：此字段可能返回 null，表示取不到有效值。
	Type *string `json:"Type,omitnil" name:"Type"`

	// 域名绑定配置
	// 注意：此字段可能返回 null，表示取不到有效值。
	DomainNameConfig *DomainNameConfig `json:"DomainNameConfig,omitnil" name:"DomainNameConfig"`

	// false
	// 注意：此字段可能返回 null，表示取不到有效值。
	Debug *bool `json:"Debug,omitnil" name:"Debug"`

	// 中断原因
	// 注意：此字段可能返回 null，表示取不到有效值。
	AbortReason *int64 `json:"AbortReason,omitnil" name:"AbortReason"`

	// 任务的创建时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	CreatedAt *string `json:"CreatedAt,omitnil" name:"CreatedAt"`

	// 项目ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 通知事件回调
	// 注意：此字段可能返回 null，表示取不到有效值。
	NotificationHooks []*NotificationHook `json:"NotificationHooks,omitnil" name:"NotificationHooks"`

	// 每秒接收字节数
	// 注意：此字段可能返回 null，表示取不到有效值。
	NetworkReceiveRate *float64 `json:"NetworkReceiveRate,omitnil" name:"NetworkReceiveRate"`

	// 每秒发送字节数
	// 注意：此字段可能返回 null，表示取不到有效值。
	NetworkSendRate *float64 `json:"NetworkSendRate,omitnil" name:"NetworkSendRate"`

	// 任务状态描述
	// 注意：此字段可能返回 null，表示取不到有效值。
	Message *string `json:"Message,omitnil" name:"Message"`

	// test-project
	// 注意：此字段可能返回 null，表示取不到有效值。
	ProjectName *string `json:"ProjectName,omitnil" name:"ProjectName"`

	// test-scenario
	// 注意：此字段可能返回 null，表示取不到有效值。
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`
}

type Label struct {
	// label名字
	LabelName *string `json:"LabelName,omitnil" name:"LabelName"`

	// label值
	LabelValue *string `json:"LabelValue,omitnil" name:"LabelValue"`
}

type LabelWithValues struct {
	// 标签名称
	LabelName *string `json:"LabelName,omitnil" name:"LabelName"`

	// 标签值
	LabelValues []*string `json:"LabelValues,omitnil" name:"LabelValues"`
}

type Load struct {
	// 施压配置
	// 注意：此字段可能返回 null，表示取不到有效值。
	LoadSpec *LoadSpec `json:"LoadSpec,omitnil" name:"LoadSpec"`

	// 压力来源
	// 注意：此字段可能返回 null，表示取不到有效值。
	VpcLoadDistribution *VpcLoadDistribution `json:"VpcLoadDistribution,omitnil" name:"VpcLoadDistribution"`

	// 压力分布
	// 注意：此字段可能返回 null，表示取不到有效值。
	GeoRegionsLoadDistribution []*GeoRegionsLoadItem `json:"GeoRegionsLoadDistribution,omitnil" name:"GeoRegionsLoadDistribution"`
}

type LoadSource struct {
	// 发压host的IP
	// 注意：此字段可能返回 null，表示取不到有效值。
	IP *string `json:"IP,omitnil" name:"IP"`

	// 发压host所在的pod
	// 注意：此字段可能返回 null，表示取不到有效值。
	PodName *string `json:"PodName,omitnil" name:"PodName"`

	// 所属地域
	// 注意：此字段可能返回 null，表示取不到有效值。
	Region *string `json:"Region,omitnil" name:"Region"`
}

type LoadSpec struct {
	// 并发施压模式的配置
	// 注意：此字段可能返回 null，表示取不到有效值。
	Concurrency *Concurrency `json:"Concurrency,omitnil" name:"Concurrency"`

	// RPS施压模式的配置
	// 注意：此字段可能返回 null，表示取不到有效值。
	RequestsPerSecond *RequestsPerSecond `json:"RequestsPerSecond,omitnil" name:"RequestsPerSecond"`

	// 脚本内置压力模式
	// 注意：此字段可能返回 null，表示取不到有效值。
	ScriptOrigin *ScriptOrigin `json:"ScriptOrigin,omitnil" name:"ScriptOrigin"`
}

type MetricInfo struct {
	// 后台指标
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 前台展示指标名称
	Alias *string `json:"Alias,omitnil" name:"Alias"`

	// 指标描述
	// 注意：此字段可能返回 null，表示取不到有效值。
	Description *string `json:"Description,omitnil" name:"Description"`

	// 指标类型
	MetricType *string `json:"MetricType,omitnil" name:"MetricType"`

	// 默认指标单位
	Unit *string `json:"Unit,omitnil" name:"Unit"`

	// 指标支持的聚合函数
	Aggregations []*AggregationLegend `json:"Aggregations,omitnil" name:"Aggregations"`

	// 是否内部指标，内部指标不可在前台提供用户自由选择
	InnerMetric *bool `json:"InnerMetric,omitnil" name:"InnerMetric"`
}

type MetricLabelWithValues struct {
	// metric 名字
	MetricName *string `json:"MetricName,omitnil" name:"MetricName"`

	// label及values 数组
	LabelValuesSet []*LabelWithValues `json:"LabelValuesSet,omitnil" name:"LabelValuesSet"`
}

type NormalLog struct {
	// 毫秒时间戳
	// 注意：此字段可能返回 null，表示取不到有效值。
	Timestamp *string `json:"Timestamp,omitnil" name:"Timestamp"`

	// 日志级别
	// 注意：此字段可能返回 null，表示取不到有效值。
	SeverityText *string `json:"SeverityText,omitnil" name:"SeverityText"`

	// 日志输出内容
	// 注意：此字段可能返回 null，表示取不到有效值。
	Body *string `json:"Body,omitnil" name:"Body"`
}

type Notification struct {
	// 发生事件
	Events []*string `json:"Events,omitnil" name:"Events"`

	// webhook的网址
	URL *string `json:"URL,omitnil" name:"URL"`
}

type NotificationHook struct {
	// 通知事件
	// 注意：此字段可能返回 null，表示取不到有效值。
	Events []*string `json:"Events,omitnil" name:"Events"`

	// 回调 URL
	// 注意：此字段可能返回 null，表示取不到有效值。
	URL *string `json:"URL,omitnil" name:"URL"`
}

type Project struct {
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 项目名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 项目描述
	// 注意：此字段可能返回 null，表示取不到有效值。
	Description *string `json:"Description,omitnil" name:"Description"`

	// 标签数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	Tags []*TagSpec `json:"Tags,omitnil" name:"Tags"`

	// 项目状态
	Status *int64 `json:"Status,omitnil" name:"Status"`

	// 创建时间
	CreatedAt *string `json:"CreatedAt,omitnil" name:"CreatedAt"`

	// 修改时间
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`

	// App ID
	AppId *int64 `json:"AppId,omitnil" name:"AppId"`

	// 用户ID
	Uin *string `json:"Uin,omitnil" name:"Uin"`

	// 子用户ID
	SubAccountUin *string `json:"SubAccountUin,omitnil" name:"SubAccountUin"`
}

type ProtocolInfo struct {
	// 协议详情
	// 注意：此字段可能返回 null，表示取不到有效值。
	Name *string `json:"Name,omitnil" name:"Name"`

	// 文件大小
	// 注意：此字段可能返回 null，表示取不到有效值。
	Size *int64 `json:"Size,omitnil" name:"Size"`

	// 文件类型
	// 注意：此字段可能返回 null，表示取不到有效值。
	Type *string `json:"Type,omitnil" name:"Type"`

	// 更新时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`

	// 文件 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	FileId *string `json:"FileId,omitnil" name:"FileId"`
}

type ReactionTimeRange struct {
	// 最小响应时间，单位ms
	Min *string `json:"Min,omitnil" name:"Min"`

	// 最大响应时间，单位ms
	Max *string `json:"Max,omitnil" name:"Max"`
}

type RegionDetail struct {
	// 地域代码
	Region *string `json:"Region,omitnil" name:"Region"`

	// 地域ID
	RegionId *int64 `json:"RegionId,omitnil" name:"RegionId"`

	// 地域所在的地区
	Area *string `json:"Area,omitnil" name:"Area"`

	// 地域名称
	RegionName *string `json:"RegionName,omitnil" name:"RegionName"`

	// 地域状态
	RegionState *int64 `json:"RegionState,omitnil" name:"RegionState"`

	// 地域简称
	RegionShortName *string `json:"RegionShortName,omitnil" name:"RegionShortName"`

	// 创建时间
	CreatedAt *string `json:"CreatedAt,omitnil" name:"CreatedAt"`

	// 更新时间
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`
}

type RequestSummary struct {
	// 请求URL
	Service *string `json:"Service,omitnil" name:"Service"`

	// 请求方法
	Method *string `json:"Method,omitnil" name:"Method"`

	// 请求次数
	Count *int64 `json:"Count,omitnil" name:"Count"`

	// 请求响应平均耗时，单位秒
	Average *float64 `json:"Average,omitnil" name:"Average"`

	// 请求p90耗时，单位秒
	P90 *float64 `json:"P90,omitnil" name:"P90"`

	// 请求p95耗时，单位秒
	P95 *float64 `json:"P95,omitnil" name:"P95"`

	// 请求最小耗时，单位秒
	Min *float64 `json:"Min,omitnil" name:"Min"`

	// 请求最大耗时，单位秒
	Max *float64 `json:"Max,omitnil" name:"Max"`

	// 请求错误率
	ErrorPercentage *float64 `json:"ErrorPercentage,omitnil" name:"ErrorPercentage"`

	// 请求p99耗时，单位秒
	P99 *float64 `json:"P99,omitnil" name:"P99"`

	// 响应状态码
	Status *string `json:"Status,omitnil" name:"Status"`

	// 响应详情
	Result *string `json:"Result,omitnil" name:"Result"`
}

type RequestsPerSecond struct {
	// 最大RPS
	// 注意：此字段可能返回 null，表示取不到有效值。
	MaxRequestsPerSecond *int64 `json:"MaxRequestsPerSecond,omitnil" name:"MaxRequestsPerSecond"`

	// 施压时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	DurationSeconds *int64 `json:"DurationSeconds,omitnil" name:"DurationSeconds"`

	// deprecated
	// 注意：此字段可能返回 null，表示取不到有效值。
	TargetVirtualUsers *int64 `json:"TargetVirtualUsers,omitnil" name:"TargetVirtualUsers"`

	// 资源数
	// 注意：此字段可能返回 null，表示取不到有效值。
	Resources *int64 `json:"Resources,omitnil" name:"Resources"`

	// 起始RPS
	// 注意：此字段可能返回 null，表示取不到有效值。
	StartRequestsPerSecond *int64 `json:"StartRequestsPerSecond,omitnil" name:"StartRequestsPerSecond"`

	// 目标RPS，入参无效
	// 注意：此字段可能返回 null，表示取不到有效值。
	TargetRequestsPerSecond *int64 `json:"TargetRequestsPerSecond,omitnil" name:"TargetRequestsPerSecond"`

	// 优雅关停的等待时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	GracefulStopSeconds *int64 `json:"GracefulStopSeconds,omitnil" name:"GracefulStopSeconds"`
}

// Predefined struct for user
type RestartCronJobsRequestParams struct {
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 定时任务ID数组
	CronJobIds []*string `json:"CronJobIds,omitnil" name:"CronJobIds"`
}

type RestartCronJobsRequest struct {
	*tchttp.BaseRequest
	
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 定时任务ID数组
	CronJobIds []*string `json:"CronJobIds,omitnil" name:"CronJobIds"`
}

func (r *RestartCronJobsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RestartCronJobsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "CronJobIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "RestartCronJobsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type RestartCronJobsResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type RestartCronJobsResponse struct {
	*tchttp.BaseResponse
	Response *RestartCronJobsResponseParams `json:"Response"`
}

func (r *RestartCronJobsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *RestartCronJobsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type SLALabel struct {
	// 标签名
	// 注意：此字段可能返回 null，表示取不到有效值。
	LabelName *string `json:"LabelName,omitnil" name:"LabelName"`

	// 标签值
	// 注意：此字段可能返回 null，表示取不到有效值。
	LabelValue *string `json:"LabelValue,omitnil" name:"LabelValue"`
}

type SLAPolicy struct {
	// SLA 规则
	// 注意：此字段可能返回 null，表示取不到有效值。
	SLARules []*SLARule `json:"SLARules,omitnil" name:"SLARules"`

	// 告警通知渠道
	// 注意：此字段可能返回 null，表示取不到有效值。
	AlertChannel *AlertChannel `json:"AlertChannel,omitnil" name:"AlertChannel"`
}

type SLARule struct {
	// 压测指标
	// 注意：此字段可能返回 null，表示取不到有效值。
	Metric *string `json:"Metric,omitnil" name:"Metric"`

	// 压测指标聚合方法
	// 注意：此字段可能返回 null，表示取不到有效值。
	Aggregation *string `json:"Aggregation,omitnil" name:"Aggregation"`

	// 压测指标条件判断符号
	// 注意：此字段可能返回 null，表示取不到有效值。
	Condition *string `json:"Condition,omitnil" name:"Condition"`

	// 阈值
	// 注意：此字段可能返回 null，表示取不到有效值。
	Value *float64 `json:"Value,omitnil" name:"Value"`

	// 标签
	// 注意：此字段可能返回 null，表示取不到有效值。
	LabelFilter []*SLALabel `json:"LabelFilter,omitnil" name:"LabelFilter"`

	// 是否停止压测任务
	// 注意：此字段可能返回 null，表示取不到有效值。
	AbortFlag *bool `json:"AbortFlag,omitnil" name:"AbortFlag"`

	// 持续时长
	// 注意：此字段可能返回 null，表示取不到有效值。
	For *string `json:"For,omitnil" name:"For"`
}

type SampleLog struct {
	// 日志毫秒时间戳
	// 注意：此字段可能返回 null，表示取不到有效值。
	Timestamp *string `json:"Timestamp,omitnil" name:"Timestamp"`

	// 采样日志属性
	// 注意：此字段可能返回 null，表示取不到有效值。
	Attributes *Attributes `json:"Attributes,omitnil" name:"Attributes"`

	// har格式的采样请求
	// 注意：此字段可能返回 null，表示取不到有效值。
	Body *string `json:"Body,omitnil" name:"Body"`
}

type SamplePair struct {
	// is the number of milliseconds since the epoch (1970-01-01 00:00 UTC) excluding leap seconds.
	Timestamp *int64 `json:"Timestamp,omitnil" name:"Timestamp"`

	// is a representation of a value for a given sample at a given time.
	Value *float64 `json:"Value,omitnil" name:"Value"`
}

type SampleStream struct {
	// labels描述
	// 注意：此字段可能返回 null，表示取不到有效值。
	Labels []*Label `json:"Labels,omitnil" name:"Labels"`

	// 指标采样数组
	Values []*SamplePair `json:"Values,omitnil" name:"Values"`

	// 指标序列名字
	// 注意：此字段可能返回 null，表示取不到有效值。
	Name *string `json:"Name,omitnil" name:"Name"`
}

type Scenario struct {
	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 场景名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 场景描述
	// 注意：此字段可能返回 null，表示取不到有效值。
	Description *string `json:"Description,omitnil" name:"Description"`

	// 场景类型，如pts-http, pts-js, pts-trpc, pts-jmeter
	// 注意：此字段可能返回 null，表示取不到有效值。
	Type *string `json:"Type,omitnil" name:"Type"`

	// 场景状态
	// 注意：此字段可能返回 null，表示取不到有效值。
	Status *int64 `json:"Status,omitnil" name:"Status"`

	// 施压配置
	// 注意：此字段可能返回 null，表示取不到有效值。
	Load *Load `json:"Load,omitnil" name:"Load"`

	// deprecated
	// 注意：此字段可能返回 null，表示取不到有效值。
	EncodedScripts *string `json:"EncodedScripts,omitnil" name:"EncodedScripts"`

	// deprecated
	// 注意：此字段可能返回 null，表示取不到有效值。
	Configs []*string `json:"Configs,omitnil" name:"Configs"`

	// deprecated
	// 注意：此字段可能返回 null，表示取不到有效值。
	Extensions []*string `json:"Extensions,omitnil" name:"Extensions"`

	// 测试数据集
	// 注意：此字段可能返回 null，表示取不到有效值。
	Datasets []*TestData `json:"Datasets,omitnil" name:"Datasets"`

	// SLA规则的ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	SLAId *string `json:"SLAId,omitnil" name:"SLAId"`

	// Cron Job规则的ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	CronId *string `json:"CronId,omitnil" name:"CronId"`

	// 场景创建时间
	CreatedAt *string `json:"CreatedAt,omitnil" name:"CreatedAt"`

	// 场景修改时间
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`

	// 项目ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// App ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	AppId *int64 `json:"AppId,omitnil" name:"AppId"`

	// 用户ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	Uin *string `json:"Uin,omitnil" name:"Uin"`

	// 子用户ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	SubAccountUin *string `json:"SubAccountUin,omitnil" name:"SubAccountUin"`

	// 测试脚本信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	TestScripts []*ScriptInfo `json:"TestScripts,omitnil" name:"TestScripts"`

	// 协议文件信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	Protocols []*ProtocolInfo `json:"Protocols,omitnil" name:"Protocols"`

	// 请求文件信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	RequestFiles []*FileInfo `json:"RequestFiles,omitnil" name:"RequestFiles"`

	// SLA 策略
	// 注意：此字段可能返回 null，表示取不到有效值。
	SLAPolicy *SLAPolicy `json:"SLAPolicy,omitnil" name:"SLAPolicy"`

	// 扩展包信息
	// 注意：此字段可能返回 null，表示取不到有效值。
	Plugins []*FileInfo `json:"Plugins,omitnil" name:"Plugins"`

	// 域名解析配置
	// 注意：此字段可能返回 null，表示取不到有效值。
	DomainNameConfig *DomainNameConfig `json:"DomainNameConfig,omitnil" name:"DomainNameConfig"`

	// 通知事件回调
	// 注意：此字段可能返回 null，表示取不到有效值。
	NotificationHooks []*NotificationHook `json:"NotificationHooks,omitnil" name:"NotificationHooks"`

	// 创建人员
	// 注意：此字段可能返回 null，表示取不到有效值。
	Owner *string `json:"Owner,omitnil" name:"Owner"`

	// 场景所在的项目的名字
	// 注意：此字段可能返回 null，表示取不到有效值。
	ProjectName *string `json:"ProjectName,omitnil" name:"ProjectName"`
}

type ScenarioRelatedJobsParams struct {
	// job偏移量
	Offset *int64 `json:"Offset,omitnil" name:"Offset"`

	// 限制最多查询的job数
	Limit *int64 `json:"Limit,omitnil" name:"Limit"`

	// 排序字段
	OrderBy *string `json:"OrderBy,omitnil" name:"OrderBy"`

	// 是否升序
	Ascend *bool `json:"Ascend,omitnil" name:"Ascend"`
}

type ScenarioWithJobs struct {
	// scecario结果
	// 注意：此字段可能返回 null，表示取不到有效值。
	Scenario *Scenario `json:"Scenario,omitnil" name:"Scenario"`

	// job结果
	// 注意：此字段可能返回 null，表示取不到有效值。
	Jobs []*Job `json:"Jobs,omitnil" name:"Jobs"`
}

type ScriptInfo struct {
	// 文件名
	// 注意：此字段可能返回 null，表示取不到有效值。
	Name *string `json:"Name,omitnil" name:"Name"`

	// 文件大小
	// 注意：此字段可能返回 null，表示取不到有效值。
	Size *int64 `json:"Size,omitnil" name:"Size"`

	// 文件类型
	// 注意：此字段可能返回 null，表示取不到有效值。
	Type *string `json:"Type,omitnil" name:"Type"`

	// 更新时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`

	// base64编码后的文件内容
	// 注意：此字段可能返回 null，表示取不到有效值。
	EncodedContent *string `json:"EncodedContent,omitnil" name:"EncodedContent"`

	// base64编码后的har结构体
	// 注意：此字段可能返回 null，表示取不到有效值。
	EncodedHttpArchive *string `json:"EncodedHttpArchive,omitnil" name:"EncodedHttpArchive"`

	// 脚本权重，范围 1-100
	// 注意：此字段可能返回 null，表示取不到有效值。
	LoadWeight *int64 `json:"LoadWeight,omitnil" name:"LoadWeight"`

	// 文件 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	FileId *string `json:"FileId,omitnil" name:"FileId"`
}

type ScriptOrigin struct {
	// 机器数量
	MachineNumber *int64 `json:"MachineNumber,omitnil" name:"MachineNumber"`

	// 机器规格
	MachineSpecification *string `json:"MachineSpecification,omitnil" name:"MachineSpecification"`

	// 压测时长
	DurationSeconds *int64 `json:"DurationSeconds,omitnil" name:"DurationSeconds"`
}

type Stage struct {
	// 施压时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	DurationSeconds *int64 `json:"DurationSeconds,omitnil" name:"DurationSeconds"`

	// 虚拟用户数
	// 注意：此字段可能返回 null，表示取不到有效值。
	TargetVirtualUsers *int64 `json:"TargetVirtualUsers,omitnil" name:"TargetVirtualUsers"`
}

// Predefined struct for user
type StartJobRequestParams struct {
	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 任务发起人
	JobOwner *string `json:"JobOwner,omitnil" name:"JobOwner"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 是否调试
	Debug *bool `json:"Debug,omitnil" name:"Debug"`

	// 备注
	Note *string `json:"Note,omitnil" name:"Note"`
}

type StartJobRequest struct {
	*tchttp.BaseRequest
	
	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 任务发起人
	JobOwner *string `json:"JobOwner,omitnil" name:"JobOwner"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 是否调试
	Debug *bool `json:"Debug,omitnil" name:"Debug"`

	// 备注
	Note *string `json:"Note,omitnil" name:"Note"`
}

func (r *StartJobRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *StartJobRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ScenarioId")
	delete(f, "JobOwner")
	delete(f, "ProjectId")
	delete(f, "Debug")
	delete(f, "Note")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "StartJobRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type StartJobResponseParams struct {
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type StartJobResponse struct {
	*tchttp.BaseResponse
	Response *StartJobResponseParams `json:"Response"`
}

func (r *StartJobResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *StartJobResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type TagSpec struct {
	// 标签键
	// 注意：此字段可能返回 null，表示取不到有效值。
	TagKey *string `json:"TagKey,omitnil" name:"TagKey"`

	// 标签值
	// 注意：此字段可能返回 null，表示取不到有效值。
	TagValue *string `json:"TagValue,omitnil" name:"TagValue"`
}

type TestData struct {
	// 测试数据集所在的文件名
	// 注意：此字段可能返回 null，表示取不到有效值。
	Name *string `json:"Name,omitnil" name:"Name"`

	// 测试数据集是否分片
	// 注意：此字段可能返回 null，表示取不到有效值。
	Split *bool `json:"Split,omitnil" name:"Split"`

	// 首行是否为参数名
	// 注意：此字段可能返回 null，表示取不到有效值。
	HeaderInFile *bool `json:"HeaderInFile,omitnil" name:"HeaderInFile"`

	// 参数名数组
	// 注意：此字段可能返回 null，表示取不到有效值。
	HeaderColumns []*string `json:"HeaderColumns,omitnil" name:"HeaderColumns"`

	// 文件行数
	// 注意：此字段可能返回 null，表示取不到有效值。
	LineCount *int64 `json:"LineCount,omitnil" name:"LineCount"`

	// 更新时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	UpdatedAt *string `json:"UpdatedAt,omitnil" name:"UpdatedAt"`

	// 文件字节数
	// 注意：此字段可能返回 null，表示取不到有效值。
	Size *int64 `json:"Size,omitnil" name:"Size"`

	// 头部数据行
	// 注意：此字段可能返回 null，表示取不到有效值。
	HeadLines []*string `json:"HeadLines,omitnil" name:"HeadLines"`

	// 尾部数据行
	// 注意：此字段可能返回 null，表示取不到有效值。
	TailLines []*string `json:"TailLines,omitnil" name:"TailLines"`

	// 文件类型
	// 注意：此字段可能返回 null，表示取不到有效值。
	Type *string `json:"Type,omitnil" name:"Type"`

	// 文件 ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	FileId *string `json:"FileId,omitnil" name:"FileId"`
}

// Predefined struct for user
type UpdateCronJobRequestParams struct {
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 定时任务ID
	CronJobId *string `json:"CronJobId,omitnil" name:"CronJobId"`

	// 备注
	Note *string `json:"Note,omitnil" name:"Note"`

	// cron表达式
	CronExpression *string `json:"CronExpression,omitnil" name:"CronExpression"`

	// 执行频率类型，1:只执行一次; 2:日粒度; 3:周粒度; 4:高级
	FrequencyType *int64 `json:"FrequencyType,omitnil" name:"FrequencyType"`

	// 定时任务名字
	Name *string `json:"Name,omitnil" name:"Name"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 场景名称
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`

	// 任务发起人
	JobOwner *string `json:"JobOwner,omitnil" name:"JobOwner"`

	// 结束时间
	EndTime *string `json:"EndTime,omitnil" name:"EndTime"`

	// Notice ID
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`
}

type UpdateCronJobRequest struct {
	*tchttp.BaseRequest
	
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 定时任务ID
	CronJobId *string `json:"CronJobId,omitnil" name:"CronJobId"`

	// 备注
	Note *string `json:"Note,omitnil" name:"Note"`

	// cron表达式
	CronExpression *string `json:"CronExpression,omitnil" name:"CronExpression"`

	// 执行频率类型，1:只执行一次; 2:日粒度; 3:周粒度; 4:高级
	FrequencyType *int64 `json:"FrequencyType,omitnil" name:"FrequencyType"`

	// 定时任务名字
	Name *string `json:"Name,omitnil" name:"Name"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 场景名称
	ScenarioName *string `json:"ScenarioName,omitnil" name:"ScenarioName"`

	// 任务发起人
	JobOwner *string `json:"JobOwner,omitnil" name:"JobOwner"`

	// 结束时间
	EndTime *string `json:"EndTime,omitnil" name:"EndTime"`

	// Notice ID
	NoticeId *string `json:"NoticeId,omitnil" name:"NoticeId"`
}

func (r *UpdateCronJobRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateCronJobRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "CronJobId")
	delete(f, "Note")
	delete(f, "CronExpression")
	delete(f, "FrequencyType")
	delete(f, "Name")
	delete(f, "ScenarioId")
	delete(f, "ScenarioName")
	delete(f, "JobOwner")
	delete(f, "EndTime")
	delete(f, "NoticeId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "UpdateCronJobRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type UpdateCronJobResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type UpdateCronJobResponse struct {
	*tchttp.BaseResponse
	Response *UpdateCronJobResponseParams `json:"Response"`
}

func (r *UpdateCronJobResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateCronJobResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type UpdateFileScenarioRelationRequestParams struct {
	// 文件 ID
	FileId *string `json:"FileId,omitnil" name:"FileId"`

	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景 ID 数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`
}

type UpdateFileScenarioRelationRequest struct {
	*tchttp.BaseRequest
	
	// 文件 ID
	FileId *string `json:"FileId,omitnil" name:"FileId"`

	// 项目 ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景 ID 数组
	ScenarioIds []*string `json:"ScenarioIds,omitnil" name:"ScenarioIds"`
}

func (r *UpdateFileScenarioRelationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateFileScenarioRelationRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "FileId")
	delete(f, "ProjectId")
	delete(f, "ScenarioIds")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "UpdateFileScenarioRelationRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type UpdateFileScenarioRelationResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type UpdateFileScenarioRelationResponse struct {
	*tchttp.BaseResponse
	Response *UpdateFileScenarioRelationResponseParams `json:"Response"`
}

func (r *UpdateFileScenarioRelationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateFileScenarioRelationResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type UpdateJobRequestParams struct {
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 任务备注信息
	Note *string `json:"Note,omitnil" name:"Note"`
}

type UpdateJobRequest struct {
	*tchttp.BaseRequest
	
	// 任务ID
	JobId *string `json:"JobId,omitnil" name:"JobId"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 任务备注信息
	Note *string `json:"Note,omitnil" name:"Note"`
}

func (r *UpdateJobRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateJobRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "JobId")
	delete(f, "ProjectId")
	delete(f, "ScenarioId")
	delete(f, "Note")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "UpdateJobRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type UpdateJobResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type UpdateJobResponse struct {
	*tchttp.BaseResponse
	Response *UpdateJobResponseParams `json:"Response"`
}

func (r *UpdateJobResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateJobResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type UpdateProjectRequestParams struct {
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 项目名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 项目描述
	Description *string `json:"Description,omitnil" name:"Description"`

	// 项目状态，默认传递1
	Status *int64 `json:"Status,omitnil" name:"Status"`

	// 标签数组
	Tags []*TagSpec `json:"Tags,omitnil" name:"Tags"`
}

type UpdateProjectRequest struct {
	*tchttp.BaseRequest
	
	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 项目名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 项目描述
	Description *string `json:"Description,omitnil" name:"Description"`

	// 项目状态，默认传递1
	Status *int64 `json:"Status,omitnil" name:"Status"`

	// 标签数组
	Tags []*TagSpec `json:"Tags,omitnil" name:"Tags"`
}

func (r *UpdateProjectRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateProjectRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ProjectId")
	delete(f, "Name")
	delete(f, "Description")
	delete(f, "Status")
	delete(f, "Tags")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "UpdateProjectRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type UpdateProjectResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type UpdateProjectResponse struct {
	*tchttp.BaseResponse
	Response *UpdateProjectResponseParams `json:"Response"`
}

func (r *UpdateProjectResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateProjectResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type UpdateScenarioRequestParams struct {
	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 场景名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 场景描述
	Description *string `json:"Description,omitnil" name:"Description"`

	// 压测引擎类型
	Type *string `json:"Type,omitnil" name:"Type"`

	// 施压配置
	Load *Load `json:"Load,omitnil" name:"Load"`

	// deprecated
	EncodedScripts *string `json:"EncodedScripts,omitnil" name:"EncodedScripts"`

	// deprecated
	Configs []*string `json:"Configs,omitnil" name:"Configs"`

	// 测试数据集
	Datasets []*TestData `json:"Datasets,omitnil" name:"Datasets"`

	// deprecated
	Extensions []*string `json:"Extensions,omitnil" name:"Extensions"`

	// SLA规则ID
	SLAId *string `json:"SLAId,omitnil" name:"SLAId"`

	// cron job ID
	CronId *string `json:"CronId,omitnil" name:"CronId"`

	// 场景状态（注：现已无需传递该参数）
	Status *int64 `json:"Status,omitnil" name:"Status"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 测试脚本路径
	TestScripts []*ScriptInfo `json:"TestScripts,omitnil" name:"TestScripts"`

	// 协议文件路径
	Protocols []*ProtocolInfo `json:"Protocols,omitnil" name:"Protocols"`

	// 请求文件路径
	RequestFiles []*FileInfo `json:"RequestFiles,omitnil" name:"RequestFiles"`

	// SLA 策略
	SLAPolicy *SLAPolicy `json:"SLAPolicy,omitnil" name:"SLAPolicy"`

	// 拓展包文件路径
	Plugins []*FileInfo `json:"Plugins,omitnil" name:"Plugins"`

	// 域名解析配置
	DomainNameConfig *DomainNameConfig `json:"DomainNameConfig,omitnil" name:"DomainNameConfig"`

	// WebHook请求配置
	NotificationHooks []*Notification `json:"NotificationHooks,omitnil" name:"NotificationHooks"`

	// 创建人名
	Owner *string `json:"Owner,omitnil" name:"Owner"`

	// 环境ID
	EnvId *string `json:"EnvId,omitnil" name:"EnvId"`
}

type UpdateScenarioRequest struct {
	*tchttp.BaseRequest
	
	// 场景ID
	ScenarioId *string `json:"ScenarioId,omitnil" name:"ScenarioId"`

	// 场景名
	Name *string `json:"Name,omitnil" name:"Name"`

	// 场景描述
	Description *string `json:"Description,omitnil" name:"Description"`

	// 压测引擎类型
	Type *string `json:"Type,omitnil" name:"Type"`

	// 施压配置
	Load *Load `json:"Load,omitnil" name:"Load"`

	// deprecated
	EncodedScripts *string `json:"EncodedScripts,omitnil" name:"EncodedScripts"`

	// deprecated
	Configs []*string `json:"Configs,omitnil" name:"Configs"`

	// 测试数据集
	Datasets []*TestData `json:"Datasets,omitnil" name:"Datasets"`

	// deprecated
	Extensions []*string `json:"Extensions,omitnil" name:"Extensions"`

	// SLA规则ID
	SLAId *string `json:"SLAId,omitnil" name:"SLAId"`

	// cron job ID
	CronId *string `json:"CronId,omitnil" name:"CronId"`

	// 场景状态（注：现已无需传递该参数）
	Status *int64 `json:"Status,omitnil" name:"Status"`

	// 项目ID
	ProjectId *string `json:"ProjectId,omitnil" name:"ProjectId"`

	// 测试脚本路径
	TestScripts []*ScriptInfo `json:"TestScripts,omitnil" name:"TestScripts"`

	// 协议文件路径
	Protocols []*ProtocolInfo `json:"Protocols,omitnil" name:"Protocols"`

	// 请求文件路径
	RequestFiles []*FileInfo `json:"RequestFiles,omitnil" name:"RequestFiles"`

	// SLA 策略
	SLAPolicy *SLAPolicy `json:"SLAPolicy,omitnil" name:"SLAPolicy"`

	// 拓展包文件路径
	Plugins []*FileInfo `json:"Plugins,omitnil" name:"Plugins"`

	// 域名解析配置
	DomainNameConfig *DomainNameConfig `json:"DomainNameConfig,omitnil" name:"DomainNameConfig"`

	// WebHook请求配置
	NotificationHooks []*Notification `json:"NotificationHooks,omitnil" name:"NotificationHooks"`

	// 创建人名
	Owner *string `json:"Owner,omitnil" name:"Owner"`

	// 环境ID
	EnvId *string `json:"EnvId,omitnil" name:"EnvId"`
}

func (r *UpdateScenarioRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateScenarioRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ScenarioId")
	delete(f, "Name")
	delete(f, "Description")
	delete(f, "Type")
	delete(f, "Load")
	delete(f, "EncodedScripts")
	delete(f, "Configs")
	delete(f, "Datasets")
	delete(f, "Extensions")
	delete(f, "SLAId")
	delete(f, "CronId")
	delete(f, "Status")
	delete(f, "ProjectId")
	delete(f, "TestScripts")
	delete(f, "Protocols")
	delete(f, "RequestFiles")
	delete(f, "SLAPolicy")
	delete(f, "Plugins")
	delete(f, "DomainNameConfig")
	delete(f, "NotificationHooks")
	delete(f, "Owner")
	delete(f, "EnvId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "UpdateScenarioRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type UpdateScenarioResponseParams struct {
	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type UpdateScenarioResponse struct {
	*tchttp.BaseResponse
	Response *UpdateScenarioResponseParams `json:"Response"`
}

func (r *UpdateScenarioResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *UpdateScenarioResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type VpcLoadDistribution struct {
	// 地域ID
	RegionId *int64 `json:"RegionId,omitnil" name:"RegionId"`

	// 地域
	Region *string `json:"Region,omitnil" name:"Region"`

	// VPC ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	VpcId *string `json:"VpcId,omitnil" name:"VpcId"`

	// 子网ID列表
	// 注意：此字段可能返回 null，表示取不到有效值。
	SubnetIds []*string `json:"SubnetIds,omitnil" name:"SubnetIds"`
}