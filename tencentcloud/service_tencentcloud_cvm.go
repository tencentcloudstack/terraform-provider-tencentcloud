package tencentcloud

import (
	"context"
	"fmt"
	"log"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CvmService struct {
	client *connectivity.TencentCloudClient
}

func (me *CvmService) DescribeInstanceById(ctx context.Context, instanceId string) (instance *cvm.Instance, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DescribeInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceSet) < 1 {
		errRet = fmt.Errorf("instance id is not found")
		return
	}
	instance = response.Response.InstanceSet[0]
	return
}

func (me *CvmService) DescribeInstanceByFilter(ctx context.Context, filters map[string]string) (instances []*cvm.Instance, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeInstancesRequest()
	request.Filters = make([]*cvm.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cvm.Filter{
			Name:   stringToPointer(k),
			Values: []*string{stringToPointer(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DescribeInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceSet) < 1 {
		errRet = fmt.Errorf("instance id is not found")
		return
	}
	instances = response.Response.InstanceSet
	return
}

func (me *CvmService) ModifyInstanceName(ctx context.Context, instanceId, instanceName string) error {
	logId := getLogId(ctx)
	request := cvm.NewModifyInstancesAttributeRequest()
	request.InstanceIds = []*string{&instanceId}
	request.InstanceName = &instanceName

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyInstancesAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) ModifySecurityGroups(ctx context.Context, instanceId string, securityGroups []*string) error {
	logId := getLogId(ctx)
	request := cvm.NewModifyInstancesAttributeRequest()
	request.InstanceIds = []*string{&instanceId}
	request.SecurityGroups = securityGroups

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyInstancesAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) ModifyProjectId(ctx context.Context, instanceId string, projectId int64) error {
	logId := getLogId(ctx)
	request := cvm.NewModifyInstancesProjectRequest()
	request.InstanceIds = []*string{&instanceId}
	request.ProjectId = &projectId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyInstancesProject(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) ModifyInstanceType(ctx context.Context, instanceId, instanceType string) error {
	logId := getLogId(ctx)
	request := cvm.NewResetInstancesTypeRequest()
	request.InstanceIds = []*string{&instanceId}
	request.InstanceType = &instanceType

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ResetInstancesType(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) ModifyPassword(ctx context.Context, instanceId, password string) error {
	logId := getLogId(ctx)
	request := cvm.NewResetInstancesPasswordRequest()
	request.InstanceIds = []*string{&instanceId}
	request.Password = &password
	forceStop := true
	request.ForceStop = &forceStop

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ResetInstancesPassword(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) ModifyVpc(ctx context.Context, instanceId, vpcId, subnetId, privateIp string) error {
	logId := getLogId(ctx)
	request := cvm.NewModifyInstancesVpcAttributeRequest()
	request.InstanceIds = []*string{&instanceId}
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		VpcId:              &vpcId,
		SubnetId:           &subnetId,
		PrivateIpAddresses: []*string{&privateIp},
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyInstancesVpcAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) StopInstance(ctx context.Context, instanceId string) error {
	logId := getLogId(ctx)
	request := cvm.NewStopInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().StopInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) StartInstance(ctx context.Context, instanceId string) error {
	logId := getLogId(ctx)
	request := cvm.NewStartInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().StartInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) DeleteInstance(ctx context.Context, instanceId string) error {
	logId := getLogId(ctx)
	request := cvm.NewTerminateInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().TerminateInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func flattenCvmTagsMapping(tags []*cvm.Tag) (mapping map[string]string) {
	mapping = make(map[string]string)
	for _, tag := range tags {
		mapping[*tag.Key] = *tag.Value
	}
	return
}
