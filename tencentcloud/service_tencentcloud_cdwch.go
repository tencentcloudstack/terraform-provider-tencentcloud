package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	clickhouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CdwchService struct {
	client *connectivity.TencentCloudClient
}

func (me *CdwchService) DescribeInstance(ctx context.Context, instanceId string) (InstanceInfo *cdwch.InstanceInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewDescribeInstanceRequest()
	request.IsOpenApi = helper.Bool(true)
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	InstanceInfo = response.Response.InstanceInfo

	return
}

func (me *CdwchService) DestroyInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewDestroyInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DestroyInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response.ErrorMsg != nil && *response.Response.ErrorMsg != "" {
		errRet = fmt.Errorf(*response.Response.ErrorMsg)
	}
	return
}

func (me *CdwchService) ResizeDisk(ctx context.Context, instanceId string, nodeType string, resizeDisk int) (errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewResizeDiskRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.Type = &nodeType
	request.DiskSize = helper.IntInt64(resizeDisk)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().ResizeDisk(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ErrorMsg != nil && *response.Response.ErrorMsg != "" {
		errRet = fmt.Errorf(*response.Response.ErrorMsg)
	}
	return
}

func (me *CdwchService) ScaleUpInstance(ctx context.Context, instanceId, nodeType, specName string) (errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewScaleUpInstanceRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.ScaleUpEnableRolling = helper.Bool(true)
	request.Type = &nodeType
	request.SpecName = &specName
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().ScaleUpInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ErrorMsg != nil && *response.Response.ErrorMsg != "" {
		errRet = fmt.Errorf(*response.Response.ErrorMsg)
	}
	return
}

func (me *CdwchService) ScaleOutInstance(ctx context.Context, instanceId string, nodeType string, scaleOutCluster string, nodeCount int, userSubnetIPNum int, shardIps []*string) (errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewScaleOutInstanceRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.Type = &nodeType
	request.NodeCount = helper.IntInt64(nodeCount)
	request.ScaleOutCluster = &scaleOutCluster
	request.UserSubnetIPNum = helper.IntInt64(userSubnetIPNum)
	if shardIps != nil {
		request.ReduceShardInfo = shardIps
	}
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().ScaleOutInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response.ErrorMsg != nil && *response.Response.ErrorMsg != "" {
		errRet = fmt.Errorf(*response.Response.ErrorMsg)
	}
	return
}

func (me *CdwchService) DescribeInstanceClusters(ctx context.Context, instanceId string) (clusterInfos []*cdwch.ClusterInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewDescribeInstanceClustersRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeInstanceClusters(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	clusterInfos = response.Response.Clusters

	return
}

func (me *CdwchService) DescribeInstancesNew(ctx context.Context, instanceId string) (instancesList []*cdwch.InstanceInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewDescribeInstancesNewRequest()
	request.SearchInstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeInstancesNew(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instancesList = response.Response.InstancesList

	return
}

func (me *CdwchService) DescribeBackUpScheduleById(ctx context.Context, instanceId string) (backup *cdwch.DescribeBackUpScheduleResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewDescribeBackUpScheduleRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeBackUpSchedule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	backup = response.Response
	return
}

func (me *CdwchService) CreateBackUpSchedule(ctx context.Context, instanceId string, paramMap map[string]interface{}) error {
	logId := getLogId(ctx)

	request := cdwch.NewCreateBackUpScheduleRequest()
	request.InstanceId = &instanceId
	for k, v := range paramMap {
		if k == "schedule_id" {
			value := v.(int64)
			request.ScheduleId = helper.Int64(value)
		}
		if k == "operation_type" {
			value := v.(string)
			request.OperationType = helper.String(value)
		}
		if k == "schedule_type" {
			value := v.(string)
			request.ScheduleType = helper.String(value)
		}
		if k == "week_days" {
			value := v.(string)
			request.WeekDays = helper.String(value)
		}
		if k == "execute_hour" {
			value := v.(int)
			request.ExecuteHour = helper.IntInt64(value)
		}
		if k == "retain_days" {
			value := v.(int)
			request.RetainDays = helper.IntInt64(value)
		}
		if k == "back_up_tables" {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				backupTableContent := cdwch.BackupTableContent{}
				if v, ok := dMap["database"]; ok {
					backupTableContent.Database = helper.String(v.(string))
				}
				if v, ok := dMap["table"]; ok {
					backupTableContent.Table = helper.String(v.(string))
				}
				if v, ok := dMap["total_bytes"]; ok {
					backupTableContent.TotalBytes = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["v_cluster"]; ok {
					backupTableContent.VCluster = helper.String(v.(string))
				}
				if v, ok := dMap["ips"]; ok {
					backupTableContent.Ips = helper.String(v.(string))
				}
				if v, ok := dMap["zoo_path"]; ok {
					backupTableContent.ZooPath = helper.String(v.(string))
				}
				if v, ok := dMap["rip"]; ok {
					backupTableContent.Rip = helper.String(v.(string))
				}
				request.BackUpTables = append(request.BackUpTables, &backupTableContent)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCdwchClient().CreateBackUpSchedule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clickhouse backUpSchedule failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func (me *CdwchService) DescribeClickhouseBackupJobsByFilter(ctx context.Context, param map[string]interface{}) (backupJobs []*clickhouse.BackUpJobDisplay, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdwch.NewDescribeBackUpJobRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "begin_time" {
			request.BeginTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.PageNum = &offset
		request.PageSize = &limit
		response, err := me.client.UseCdwchClient().DescribeBackUpJob(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.BackUpJobs) < 1 {
			break
		}
		backupJobs = append(backupJobs, response.Response.BackUpJobs...)
		if len(response.Response.BackUpJobs) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
