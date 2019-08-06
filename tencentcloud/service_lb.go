package tencentcloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	lb "github.com/zqfan/tencentcloud-sdk-go/services/lb/unversioned"
)

func waitForLBReady(client *lb.Client, lbid *string) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		req := lb.NewDescribeLoadBalancersRequest()
		req.LoadBalancerIds = []*string{lbid}
		resp, err := client.DescribeLoadBalancers(req)
		if err != nil {
			return resource.RetryableError(err)
		}
		if *resp.LoadBalancerSet[0].Status == lb.LBStatusCreating {
			return resource.RetryableError(fmt.Errorf("LB %s is still creating...", *lbid))
		} else if *resp.LoadBalancerSet[0].Status == lb.LBStatusReady {
			return nil
		} else {
			return resource.NonRetryableError(fmt.Errorf("LB %s status unknown...", *lbid))
		}
	})
}

func waitForLBTaskFinish(client *lb.Client, taskid *int) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		req := lb.NewDescribeLoadBalancersTaskResultRequest()
		req.RequestId = taskid
		resp, err := client.DescribeLoadBalancersTaskResult(req)
		if err != nil {
			return resource.RetryableError(err)
		}
		if *resp.Data.Status == lb.LBTaskSuccess {
			return nil
		} else if *resp.Data.Status == lb.LBTaskFail {
			return resource.NonRetryableError(fmt.Errorf("LB task %d fail...", *taskid))
		} else {
			return resource.RetryableError(fmt.Errorf("LB task %d is still waiting...", *taskid))
		}
	})
}
