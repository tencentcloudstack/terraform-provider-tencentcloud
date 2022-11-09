package tencentcloud

import (
	"context"
	"log"

	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CatService struct {
	client *connectivity.TencentCloudClient
}

func (me *CatService) DescribeCatTaskSet(ctx context.Context, taskId string) (taskSet *cat.ProbeTask, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cat.NewDescribeProbeTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.TaskIDs = []*string{helper.String(taskId)}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*cat.ProbeTask, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCatClient().DescribeProbeTasks(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TaskSet) < 1 {
			break
		}
		instances = append(instances, response.Response.TaskSet...)
		if len(response.Response.TaskSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	taskSet = instances[0]
	return
}

func (me *CatService) DeleteCatTaskSetById(ctx context.Context, taskId string) (errRet error) {
	logId := getLogId(ctx)

	request := cat.NewDeleteProbeTaskRequest()

	request.TaskIds = []*string{helper.String(taskId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCatClient().DeleteProbeTask(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CatService) DescribeCatNodeByFilter(ctx context.Context, param map[string]interface{}) (node []*cat.NodeDefine, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cat.NewDescribeProbeNodesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "node_type" {
			request.NodeType = v.(*int64)
		}

		if k == "location" {
			request.Location = v.(*int64)
		}

		if k == "is_ipv6" {
			request.IsIPv6 = v.(*bool)
		}

		if k == "node_name" {
			request.NodeName = v.(*string)
		}

		if k == "pay_mode" {
			request.PayMode = v.(*int64)
		}

	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCatClient().DescribeProbeNodes(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	node = append(node, response.Response.NodeSet...)

	return
}

func (me *CatService) DescribeCatProbeDataByFilter(ctx context.Context, param map[string]interface{}) (probeData []*cat.DetailedSingleDataDefine, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cat.NewDescribeDetailedSingleProbeDataRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "begin_time" {
			request.BeginTime = v.(*uint64)
		}

		if k == "end_time" {
			request.EndTime = v.(*uint64)
		}

		if k == "task_type" {
			request.TaskType = v.(*string)
		}

		if k == "sort_field" {
			request.SortField = v.(*string)
		}

		if k == "ascending" {
			request.Ascending = v.(*bool)
		}

		if k == "selected_fields" {
			var selectedFields []*string
			selectedFields = append(selectedFields, v.([]*string)...)
			request.SelectedFields = selectedFields
		}

		if k == "offset" {
			request.Offset = v.(*int64)
		}

		if k == "limit" {
			request.Limit = v.(*int64)
		}

		if k == "task_id" {
			var taskId []*string
			taskId = append(taskId, v.([]*string)...)
			request.TaskID = taskId
		}

		if k == "operators" {
			var operators []*string
			operators = append(operators, v.([]*string)...)
			request.Operators = operators
		}

		if k == "districts" {
			var districts []*string
			districts = append(districts, v.([]*string)...)
			request.Districts = districts
		}

		if k == "error_types" {
			var errorTypes []*string
			errorTypes = append(errorTypes, v.([]*string)...)
			request.ErrorTypes = errorTypes
		}

		if k == "city" {
			var city []*string
			city = append(city, v.([]*string)...)
			request.City = city
		}

	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCatClient().DescribeDetailedSingleProbeData(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	probeData = append(probeData, response.Response.DataSet...)

	return
}
