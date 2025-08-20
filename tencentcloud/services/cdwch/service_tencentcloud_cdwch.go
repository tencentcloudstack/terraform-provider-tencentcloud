package cdwch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	clickhouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewCdwchService(client *connectivity.TencentCloudClient) CdwchService {
	return CdwchService{client: client}
}

type CdwchService struct {
	client *connectivity.TencentCloudClient
}

func (me *CdwchService) DescribeInstance(ctx context.Context, instanceId string) (InstanceInfo *cdwch.InstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

	request := cdwch.NewDescribeInstancesNewRequest()
	response := cdwch.NewDescribeInstancesNewResponse()
	request.SearchInstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCdwchClient().DescribeInstancesNew(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe instances failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s describe clickhouse instances failed, reason:%+v", logId, errRet)
		return
	}

	instancesList = response.Response.InstancesList
	return
}

func (me *CdwchService) DescribeBackUpScheduleById(ctx context.Context, instanceId string) (backup *cdwch.DescribeBackUpScheduleResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCdwchClient().CreateBackUpSchedule(request)
		if e != nil {
			return tccommon.RetryError(e)
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
		logId   = tccommon.GetLogId(ctx)
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

func (me *CdwchService) DescribeClickhouseAccountByUserName(ctx context.Context, instanceId, userName string) (accounts []*AccountInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := clickhouse.NewDescribeCkSqlApisRequest()
	request.InstanceId = helper.String(instanceId)
	request.UserName = helper.String(userName)
	request.ApiType = helper.String(DESCRIBE_CK_SQL_APIS_GET_SYSTEM_USERS)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeCkSqlApis(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.ReturnData == nil {
		errRet = fmt.Errorf("DescribeCkSqlApis response is null")
		return
	}
	accounts = make([]*AccountInfo, 0)
	err = json.Unmarshal([]byte(*response.Response.ReturnData), &accounts)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *CdwchService) DescribeCkSqlApis(ctx context.Context, instanceId, cluster, userName, apiType string) error {
	logId := tccommon.GetLogId(ctx)

	request := clickhouse.NewDescribeCkSqlApisRequest()
	request.InstanceId = helper.String(instanceId)
	request.UserName = helper.String(userName)
	request.ApiType = helper.String(apiType)
	if cluster != "" {
		request.Cluster = helper.String(cluster)
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeCkSqlApis(request)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CdwchService) ActionAlterCkUser(ctx context.Context, apiType string, userInfo map[string]interface{}) error {
	logId := tccommon.GetLogId(ctx)

	request := clickhouse.NewActionAlterCkUserRequest()
	ckUserAlterInfo := clickhouse.CkUserAlterInfo{}
	if v, ok := userInfo["instance_id"]; ok {
		ckUserAlterInfo.InstanceId = helper.String(v.(string))
	}
	if v, ok := userInfo["user_name"]; ok {
		ckUserAlterInfo.UserName = helper.String(v.(string))
	}
	if v, ok := userInfo["password"]; ok {
		ckUserAlterInfo.PassWord = helper.String(v.(string))
	}
	if v, ok := userInfo["describe"]; ok {
		ckUserAlterInfo.Describe = helper.String(v.(string))
	}

	request.UserInfo = &ckUserAlterInfo
	request.ApiType = helper.String(apiType)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseCdwchClient().ActionAlterCkUser(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate clickhouse account failed, reason:%+v", logId, err)
		return err
	}
	return nil
}

func (me *CdwchService) DescribeCdwchAccountPermission(ctx context.Context, instanceId, cluster, username string) (userNewPrivilege *cdwch.ModifyUserNewPrivilegeRequestParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cdwch.NewDescribeCkSqlApisRequest()
	request.InstanceId = &instanceId
	request.Cluster = &cluster
	request.UserName = &username
	request.ApiType = helper.String("GetUserClusterNewPrivileges")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeCkSqlApis(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.ReturnData == nil {
		errRet = fmt.Errorf("DescribeCkSqlApis response is null")
		return
	}
	returnDate := *response.Response.ReturnData
	userNewPrivilege = &cdwch.ModifyUserNewPrivilegeRequestParams{}
	err = json.Unmarshal([]byte(returnDate), userNewPrivilege)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *CdwchService) DescribeClickhouseBackupTablesByFilter(ctx context.Context, instanceId string) (backupTables []*clickhouse.BackupTableContent, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdwch.NewDescribeBackUpTablesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeBackUpTables(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response != nil {
		backupTables = response.Response.AvailableTables
	}

	return
}

func (me *CdwchService) DescribeClickhouseKeyvalConfigById(ctx context.Context, instanceId string) (config []*cdwch.InstanceConfigInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cdwch.NewDescribeInstanceKeyValConfigsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeInstanceKeyValConfigs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.ConfigItems) < 1 {
		return
	}

	config = response.Response.ConfigItems
	return
}

func (me *CdwchService) DescribeClickhouseXmlConfigById(ctx context.Context, instanceId string) (xmlConfig []*cdwch.ClusterConfigsInfoFromEMR, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cdwch.NewDescribeClusterConfigsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeClusterConfigs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.ClusterConfList) < 1 {
		return
	}

	xmlConfig = response.Response.ClusterConfList
	return
}

func (me *CdwchService) InstanceStateRefreshFunc(instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		request := cdwch.NewDescribeInstanceStateRequest()
		request.InstanceId = &instanceId
		ratelimit.Check(request.GetAction())
		object, err := me.client.UseCdwchClient().DescribeInstanceState(request)

		if err != nil {
			return nil, "", err
		}
		if object == nil || object.Response == nil || object.Response.InstanceState == nil {
			return nil, "", nil
		}

		return object, *object.Response.InstanceState, nil
	}
}

func (me *CdwchService) DescribeClickhouseSpecByFilter(ctx context.Context, param map[string]interface{}) (spec *cdwch.DescribeSpecResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdwch.NewDescribeSpecRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Zone" {
			request.Zone = v.(*string)
		}
		if k == "PayMode" {
			request.PayMode = v.(*string)
		}
		if k == "IsElastic" {
			request.IsElastic = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeSpec(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	spec = response.Response
	return
}

func (me *CdwchService) DescribeClickhouseInstanceShardsByFilter(ctx context.Context, param map[string]interface{}) (instanceShards *cdwch.DescribeInstanceShardsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdwch.NewDescribeInstanceShardsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeInstanceShards(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	instanceShards = response.Response
	return
}

func (me *CdwchService) DescribeClickhouseInstanceNodesByFilter(ctx context.Context, param map[string]interface{}) (instanceNodes []*cdwch.InstanceNode, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cdwch.NewDescribeInstanceNodesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "NodeRole" {
			request.NodeRole = v.(*string)
		}
		if k == "DisplayPolicy" {
			request.DisplayPolicy = v.(*string)
		}
		if k == "ForceAll" {
			request.ForceAll = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdwchClient().DescribeInstanceNodes(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.InstanceNodesList) < 1 {
			break
		}
		instanceNodes = append(instanceNodes, response.Response.InstanceNodesList...)
		if len(response.Response.InstanceNodesList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
