package tdmysql

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewTdmysqlService(client *connectivity.TencentCloudClient) TdmysqlService {
	return TdmysqlService{client: client}
}

type TdmysqlService struct {
	client *connectivity.TencentCloudClient
}

func (me *TdmysqlService) DescribeTdmysqlDbInstanceById(ctx context.Context, instanceId string) (ret *tdmysqlv20211122.DescribeDBInstanceDetailResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tdmysqlv20211122.NewDescribeDBInstanceDetailRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeDBInstanceDetailWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db instance failed, Response is nil."))
		}

		ret = result.Response
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *TdmysqlService) IsolateTdmysqlDbInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tdmysqlv20211122.NewIsolateDBInstanceRequest()
	request.InstanceIds = []*string{helper.String(instanceId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().IsolateDBInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Isolate tdmysql db instance failed, Response is nil."))
		}

		if result.Response.SuccessInstanceIds == nil || len(result.Response.SuccessInstanceIds) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Isolate tdmysql db instance failed, SuccessInstanceIds is empty."))
		}

		contained := false
		for _, id := range result.Response.SuccessInstanceIds {
			if id != nil && *id == instanceId {
				contained = true
				break
			}
		}
		if !contained {
			return resource.NonRetryableError(fmt.Errorf("Isolate tdmysql db instance failed, instance id %s not in SuccessInstanceIds.", instanceId))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *TdmysqlService) DescribeTdmysqlFlow(ctx context.Context, flowId int64) (status string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tdmysqlv20211122.NewDescribeFlowRequest()
	request.FlowId = helper.Int64(flowId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeFlowWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql flow failed, Response is nil."))
		}

		if result.Response.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql flow failed, Status is nil."))
		}

		status = *result.Response.Status
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}
