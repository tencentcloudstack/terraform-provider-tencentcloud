package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	vpc "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
)

var (
	dnatNotFound = errors.New("DNAT Not found")
)

func (client *TencentCloudClient) PollingVpcTaskResult(taskId *int) (status bool, err error) {
	taskReq := vpc.NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = taskId
	status = false
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		taskResp, err := client.vpcConn.DescribeVpcTaskResult(taskReq)
		b, _ := json.Marshal(taskResp)
		log.Printf("[DEBUG] client.vpcConn.DescribeVpcTaskResult response: %s", b)
		if _, ok := err.(*common.APIError); ok {
			return resource.NonRetryableError(fmt.Errorf("client.vpcConn.CreateNatGateway error: %v", err))
		}
		if *taskResp.Data.Status == 0 {
			status = true
			return nil
		}
		return resource.RetryableError(fmt.Errorf("taskId %v, not ready, status: %v", taskId, *taskResp.Data.Status))
	})
	return
}

func (client *TencentCloudClient) PollingVpcBillResult(billId *string) (status bool, err error) {
	queryReq := vpc.NewQueryNatGatewayProductionStatusRequest()
	queryReq.BillId = billId
	status = false
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		queryResp, err := client.vpcConn.QueryNatGatewayProductionStatus(queryReq)
		b, _ := json.Marshal(queryResp)
		log.Printf("[DEBUG] client.vpcConn.QueryNatGatewayProductionStatus response: %s", b)
		if _, ok := err.(*common.APIError); ok {
			return resource.NonRetryableError(fmt.Errorf("client.vpcConn.QueryNatGatewayProductionStatus error: %v", err))
		}
		if *queryResp.Data.Status == vpc.BillStatusSuccess {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("billId %v, not ready, status: %v", billId, *queryResp.Data.Status))
	})
	return
}

func (client *TencentCloudClient) DescribeDnat(entry *vpc.DnaptRule) (response *vpc.DnaptRule, err error) {
	descReq := vpc.NewGetDnaptRuleRequest()
	descReq.NatId = entry.UniqNatId
	descReq.VpcId = entry.UniqVpcId
	descResp, descErr := client.vpcConn.GetDnaptRule(descReq)
	b, _ := json.Marshal(descResp)
	log.Printf("[DEBUG] client.vpcConn.GetDnaptRule response: %s", b)
	if _, ok := descErr.(*common.APIError); ok {
		err = fmt.Errorf("client.vpcConn.GetDnaptRule error: %v", descErr)
		return
	}
	if *descResp.Data.TotalNum == 0 || len(descResp.Data.Detail) == 0 {
		err = dnatNotFound
		return
	}
	for _, e := range descResp.Data.Detail {
		if *entry.Proto == *e.Proto && *entry.Eip == *e.Eip && *entry.Eport == *e.Eport {
			response = e
			return
		}
	}
	err = dnatNotFound
	return
}
