package ga2

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewGa2Service(client *connectivity.TencentCloudClient) Ga2Service {
	return Ga2Service{client: client}
}

type Ga2Service struct {
	client *connectivity.TencentCloudClient
}

func (me *Ga2Service) CreateAccelerateAreas(ctx context.Context, globalAcceleratorId string, areas []*ga2v20250115.AcceleratorAreas) (taskId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ga2v20250115.NewCreateAccelerateAreasRequest()
	request.GlobalAcceleratorId = &globalAcceleratorId
	request.AcceleratorAreas = areas

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *ga2v20250115.CreateAccelerateAreasResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseGa2V20250115Client().CreateAccelerateAreas(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.TaskId == nil {
		errRet = fmt.Errorf("CreateAccelerateAreas response is nil or TaskId is nil")
		return
	}

	taskId = *response.Response.TaskId
	return
}

func (me *Ga2Service) DescribeAccelerateAreas(ctx context.Context, globalAcceleratorId string) (areas []*ga2v20250115.AcceleratorAreas, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ga2v20250115.NewDescribeAccelerateAreasRequest()
	request.GlobalAcceleratorId = &globalAcceleratorId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset uint64 = 0
		limit  uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		var response *ga2v20250115.DescribeAccelerateAreasResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseGa2V20250115Client().DescribeAccelerateAreas(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || response.Response == nil {
			break
		}

		areas = append(areas, response.Response.AccelerateAreaSet...)

		if len(response.Response.AccelerateAreaSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *Ga2Service) ModifyAccelerateAreas(ctx context.Context, globalAcceleratorId string, areas []*ga2v20250115.AcceleratorAreas) (taskId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ga2v20250115.NewModifyAccelerateAreasRequest()
	request.GlobalAcceleratorId = &globalAcceleratorId
	request.AcceleratorAreas = areas

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *ga2v20250115.ModifyAccelerateAreasResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseGa2V20250115Client().ModifyAccelerateAreas(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.TaskId == nil {
		errRet = fmt.Errorf("ModifyAccelerateAreas response is nil or TaskId is nil")
		return
	}

	taskId = *response.Response.TaskId
	return
}

func (me *Ga2Service) DeleteAccelerateAreas(ctx context.Context, globalAcceleratorId string, areaIds []*string) (taskId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ga2v20250115.NewDeleteAccelerateAreasRequest()
	request.GlobalAcceleratorId = &globalAcceleratorId
	request.AcceleratorAreaIds = areaIds

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *ga2v20250115.DeleteAccelerateAreasResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseGa2V20250115Client().DeleteAccelerateAreas(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.TaskId == nil {
		errRet = fmt.Errorf("DeleteAccelerateAreas response is nil or TaskId is nil")
		return
	}

	taskId = *response.Response.TaskId
	return
}

func (me *Ga2Service) DescribeTaskResult(ctx context.Context, taskId string) (status string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := ga2v20250115.NewDescribeTaskResultRequest()
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *ga2v20250115.DescribeTaskResultResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseGa2V20250115Client().DescribeTaskResult(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.Status == nil {
		errRet = fmt.Errorf("DescribeTaskResult response is nil or Status is nil")
		return
	}

	status = *response.Response.Status
	return
}
