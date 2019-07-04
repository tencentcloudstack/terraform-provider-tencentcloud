package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type CvmService struct {
	client *connectivity.TencentCloudClient
}

func (me *CvmService) DescribeAssociateSecurityGroups(ctx context.Context, cvmId string) ([]string, error) {
	logId := GetLogId(ctx)

	request := cvm.NewDescribeInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{cvmId})

	descResponse, err := me.client.UseCvmClient().DescribeInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s get cvm by id %s error, reason: %v", logId, cvmId, err)
		return nil, err
	}

	if *descResponse.Response.TotalCount != 1 {
		err := fmt.Errorf("cvm id %s does not exist", cvmId)
		log.Printf("[CRITAL]%s %v", logId, err)
		return nil, err
	}

	cvmInstance := descResponse.Response.InstanceSet[0]
	sgIds := make([]string, 0, len(cvmInstance.SecurityGroupIds))
	for _, sgId := range cvmInstance.SecurityGroupIds {
		sgIds = append(sgIds, *sgId)
	}

	return sgIds, nil
}

func (me *CvmService) DescribeBySecurityGroups(ctx context.Context, sgId string) ([]*cvm.Instance, error) {
	logId := GetLogId(ctx)

	request := cvm.NewDescribeInstancesRequest()
	request.Limit = common.Int64Ptr(100)
	request.Filters = append(request.Filters, &cvm.Filter{
		Name:   common.StringPtr("security-group-id"),
		Values: common.StringPtrs([]string{sgId}),
	})

	response, err := me.client.UseCvmClient().DescribeInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s get cvms by security group id %s error, reason: %v", logId, sgId, err)
		return nil, err
	}

	return response.Response.InstanceSet, nil
}

func (me *CvmService) ModifySecurityGroups(ctx context.Context, cvmId string, sgIds []string) error {
	logId := GetLogId(ctx)

	request := cvm.NewModifyInstancesAttributeRequest()
	request.InstanceIds = common.StringPtrs([]string{cvmId})
	request.SecurityGroups = common.StringPtrs(sgIds)

	if _, err := me.client.UseCvmClient().ModifyInstancesAttribute(request); err != nil {
		log.Printf("[CRITAL]%s modify cvm %s security groups error, reason: %v", logId, cvmId, err)
		return err
	}

	return nil
}
