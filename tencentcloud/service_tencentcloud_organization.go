package tencentcloud

import (
	"context"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type OrganizationService struct {
	client *connectivity.TencentCloudClient
}

func (me *OrganizationService) DescribeOrganizationOrgNode(ctx context.Context, nodeId string) (orgNode *organization.OrgNode, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = organization.NewDescribeOrganizationNodesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	var offset int64 = 0
	var pageSize int64 = 50
	instances := make([]*organization.OrgNode, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseOrganizationClient().DescribeOrganizationNodes(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		instances = append(instances, response.Response.Items...)
		if len(response.Response.Items) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}

	for _, instance := range instances {
		if helper.Int64ToStr(*instance.NodeId) == nodeId {
			orgNode = instance
		}
	}

	return
}

func (me *OrganizationService) DeleteOrganizationOrgNodeById(ctx context.Context, nodeId string) (errRet error) {
	logId := getLogId(ctx)

	request := organization.NewDeleteOrganizationNodesRequest()

	request.NodeId = []*int64{helper.Int64(helper.StrToInt64(nodeId))}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseOrganizationClient().DeleteOrganizationNodes(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
