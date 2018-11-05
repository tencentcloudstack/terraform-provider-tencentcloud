package unversioned

import (
	"errors"
	"fmt"
	"time"

	"github.com/zqfan/tencentcloud-sdk-go/common"
)

// wait interval 10 seconds
func WaitForLBReady(lbid *string, c *Client, retry int) error {
	if retry <= 0 {
		return errors.New("[ERROR] LB failed to be ready")
	}
	descReq := NewDescribeLoadBalancersRequest()
	descReq.LoadBalancerIds = []*string{lbid}
	descResp, err := c.DescribeLoadBalancers(descReq)
	if _, ok := err.(*common.APIError); ok {
		return fmt.Errorf("[ERROR] Failed to describe lb, err=%v", err)
	}
	if *descResp.LoadBalancerSet[0].Status == LBStatusCreating {
		time.Sleep(10 * time.Second)
		return WaitForLBReady(lbid, c, retry-1)
	}
	return nil
}

// wait interval 10 seconds
func WaitForTaskSuccess(reqid *int, c *Client, retry int) error {
	if retry <= 0 {
		return errors.New("[ERROR] task failed to success")
	}
	taskReq := NewDescribeLoadBalancersTaskResultRequest()
	taskReq.RequestId = reqid
	taskResp, err := c.DescribeLoadBalancersTaskResult(taskReq)
	if _, ok := err.(*common.APIError); ok {
		return fmt.Errorf("[ERROR] Failed to describe task, err=%v", err)
	}
	if *taskResp.Data.Status == LBTaskSuccess {
		return nil
	} else if *taskResp.Data.Status == LBTaskFail {
		return fmt.Errorf("[ERROR] Task %d failed", *reqid)
	}
	time.Sleep(10 * time.Second)
	return WaitForTaskSuccess(reqid, c, retry-1)
}
