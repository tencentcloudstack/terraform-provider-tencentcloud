package cdh

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewCdhService(client *connectivity.TencentCloudClient) CdhService {
	return CdhService{client: client}
}

type CdhService struct {
	client *connectivity.TencentCloudClient
}

func (me *CdhService) DescribeCdhInstanceById(ctx context.Context, hostId string) (instance *cvm.HostItem, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cvm.NewDescribeHostsRequest()
	filter := cvm.Filter{
		Name:   helper.String("host-id"),
		Values: []*string{helper.String(hostId)},
	}
	request.Filters = []*cvm.Filter{&filter}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DescribeHosts(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.HostSet) < 1 {
		return
	}
	instance = response.Response.HostSet[0]
	return
}

func (me *CdhService) DescribeCdhInstanceByFilter(ctx context.Context, filters map[string]string) (instances []*cvm.HostItem, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cvm.NewDescribeHostsRequest()
	request.Filters = make([]*cvm.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cvm.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}
	var offset uint64 = 0
	var pageSize uint64 = 100
	instances = make([]*cvm.HostItem, 0)
	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCvmClient().DescribeHosts(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.HostSet) < 1 {
			break
		}
		instances = append(instances, response.Response.HostSet...)
		if len(response.Response.HostSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *CdhService) CreateCdhInstance(ctx context.Context, placement *cvm.Placement, hostChargePrepaid *cvm.ChargePrepaid, hostChargeType, hostType string) (hostId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cvm.NewAllocateHostsRequest()
	request.Placement = placement
	request.HostChargePrepaid = hostChargePrepaid
	request.HostChargeType = helper.String(hostChargeType)
	request.HostType = helper.String(hostType)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().AllocateHosts(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.HostIdSet) < 1 {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	hostId = *response.Response.HostIdSet[0]
	return
}

func (me *CdhService) ModifyHostName(ctx context.Context, hostId, hostName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cvm.NewModifyHostsAttributeRequest()
	request.HostIds = []*string{helper.String(hostId)}
	request.HostName = helper.String(hostName)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyHostsAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdhService) ModifyProject(ctx context.Context, hostId string, projectId int) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cvm.NewModifyHostsAttributeRequest()
	request.HostIds = []*string{helper.String(hostId)}
	request.ProjectId = helper.IntUint64(projectId)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyHostsAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdhService) ModifyPrepaidRenewFlag(ctx context.Context, hostId, renewFlag string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cvm.NewModifyHostsAttributeRequest()
	request.HostIds = []*string{helper.String(hostId)}
	request.RenewFlag = helper.String(renewFlag)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyHostsAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
