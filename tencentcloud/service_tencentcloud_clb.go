package tencentcloud

import (
	"context"
	"log"

	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type ClbService struct {
	client *connectivity.TencentCloudClient
}

func (me *ClbService) ModifyLoadBalanceSecurityGroups(ctx context.Context, clbId string, sgIds []string) error {
	logId := GetLogId(ctx)

	request := clb.NewSetLoadBalancerSecurityGroupsRequest()
	request.LoadBalancerId = &clbId
	request.SecurityGroups = common.StringPtrs(sgIds)

	_, err := me.client.UseClbClient().SetLoadBalancerSecurityGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
	}

	return nil
}

func (me *ClbService) DescribeLoadBalances(ctx context.Context, ids []string, sgId *string) (lb []*clb.LoadBalancer, err error) {
	logId := GetLogId(ctx)

	request := clb.NewDescribeLoadBalancersRequest()
	request.Limit = common.Int64Ptr(20)

	if len(ids) > 0 {
		request.LoadBalancerIds = common.StringPtrs(ids)
	} else if sgId != nil {
		request.SecurityGroup = sgId
	}

	response, err := me.client.UseClbClient().DescribeLoadBalancers(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return nil, err
	}

	return response.Response.LoadBalancerSet, nil
}
