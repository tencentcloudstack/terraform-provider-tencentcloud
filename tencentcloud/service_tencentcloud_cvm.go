package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CvmService struct {
	client *connectivity.TencentCloudClient
}

func (me *CvmService) DescribeInstanceSetByIds(ctx context.Context, instanceSetIds string) (instance []*cvm.Instance, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeInstancesRequest()

	instanceSet, err := helper.StrToStrList(instanceSetIds)
	if err != nil {
		return
	}

	for _, v := range instanceSet {
		ins := v
		instanceId := &ins
		request.InstanceIds = append(request.InstanceIds, instanceId)
	}

	request.Limit = helper.IntInt64(100)

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
		return
	}
	instance = response.Response.InstanceSet
	return
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
		return
	}
	instance = response.Response.InstanceSet[0]
	return
}

func (me *CvmService) DescribeInstanceByFilter(ctx context.Context, instancesId []*string, filters map[string]string) (instances []*cvm.Instance, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeInstancesRequest()
	if instancesId != nil {
		request.InstanceIds = instancesId
	} else {
		request.Filters = make([]*cvm.Filter, 0, len(filters))
		for k, v := range filters {
			filter := cvm.Filter{
				Name:   helper.String(k),
				Values: []*string{helper.String(v)},
			}
			request.Filters = append(request.Filters, &filter)
		}
	}

	var offset int64 = 0
	var pageSize int64 = 100
	instances = make([]*cvm.Instance, 0)
	for {
		request.Offset = &offset
		request.Limit = &pageSize
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

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}
		instances = append(instances, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *CvmService) DescribeInstanceInParallelByFilter(ctx context.Context, filters map[string]string) (instances []*cvm.Instance, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeInstancesRequest()

	request.Filters = make([]*cvm.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cvm.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	response, err := me.client.UseCvmClient().DescribeInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	if response == nil || len(response.Response.InstanceSet) < 1 {
		return
	}

	total := response.Response.TotalCount

	var limit = 100

	num := int(*total) / limit

	maxConcurrentNum := 50
	//g := NewGoRoutine(num + 1)
	g := NewGoRoutine(maxConcurrentNum)
	wg := sync.WaitGroup{}

	var instanceSetList = make([]interface{}, num+1)

	for i := 0; i <= num; i++ {
		wg.Add(1)
		value := i
		goFunc := func() {
			offset := value * limit
			request := cvm.NewDescribeInstancesRequest()
			request.Filters = make([]*cvm.Filter, 0, len(filters))
			for k, v := range filters {
				filter := cvm.Filter{
					Name:   helper.String(k),
					Values: []*string{helper.String(v)},
				}
				request.Filters = append(request.Filters, &filter)
			}

			request.Offset = helper.IntInt64(offset)
			request.Limit = helper.IntInt64(limit)

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

			instanceSetList[value] = response.Response.InstanceSet

			wg.Done()
			log.Printf("[DEBUG]%s thread %d finished", logId, value)
		}
		g.Run(goFunc)
	}
	wg.Wait()

	log.Printf("[DEBUG]%s DescribeInstance requet finished", logId)
	for _, v := range instanceSetList {
		instances = append(instances, v.([]*cvm.Instance)...)
	}
	log.Printf("[DEBUG]%s transfer Instance finished", logId)
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
	request.ForceStop = helper.Bool(true)

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

func (me *CvmService) ModifyInternetMaxBandwidthOut(ctx context.Context, instanceId, internetChargeType string, internetMaxBandWidthOut int64) error {
	logId := getLogId(ctx)
	request := cvm.NewResetInstancesInternetMaxBandwidthRequest()
	request.InstanceIds = []*string{&instanceId}
	request.InternetAccessible = &cvm.InternetAccessible{
		InternetChargeType:      &internetChargeType,
		InternetMaxBandwidthOut: &internetMaxBandWidthOut,
	}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCvmClient().ResetInstancesInternetMaxBandwidth(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CvmService) ModifyVpc(ctx context.Context, instanceId, vpcId, subnetId, privateIp string) error {
	logId := getLogId(ctx)
	request := cvm.NewModifyInstancesVpcAttributeRequest()
	request.InstanceIds = []*string{&instanceId}
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		VpcId:    &vpcId,
		SubnetId: &subnetId,
	}
	if privateIp != "" {
		request.VirtualPrivateCloud.PrivateIpAddresses = []*string{&privateIp}
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

func (me *CvmService) StopInstance(ctx context.Context, instanceId string, stoppedMode string) error {
	logId := getLogId(ctx)
	request := cvm.NewStopInstancesRequest()
	request.InstanceIds = []*string{&instanceId}
	if stoppedMode != "" {
		request.StoppedMode = &stoppedMode
	}

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

func (me *CvmService) DeleteInstanceSetByIds(ctx context.Context, instanceSetIds string) error {
	logId := getLogId(ctx)
	request := cvm.NewTerminateInstancesRequest()

	instanceSet, err := helper.StrToStrList(instanceSetIds)
	if err != nil {
		return err
	}

	for _, v := range instanceSet {
		ins := v
		instanceId := &ins
		request.InstanceIds = append(request.InstanceIds, instanceId)
	}

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

func (me *CvmService) ResetInstance(ctx context.Context, request *cvm.ResetInstanceRequest) (errRet error) {
	logId := getLogId(ctx)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCvmClient().ResetInstance(request)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) DescribeInstanceTypes(ctx context.Context, zone string) (instanceTypes []*cvm.InstanceTypeConfig, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeInstanceTypeConfigsRequest()
	if zone != "" {
		request.Filters = make([]*cvm.Filter, 0, 1)
		filter := &cvm.Filter{
			Name:   helper.String("zone"),
			Values: []*string{&zone},
		}
		request.Filters = append(request.Filters, filter)
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DescribeInstanceTypeConfigs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instanceTypes = response.Response.InstanceTypeConfigSet
	return
}

func (me *CvmService) DescribeInstanceTypesByFilter(ctx context.Context, filters map[string][]string) (instanceTypes []*cvm.InstanceTypeConfig, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeInstanceTypeConfigsRequest()
	request.Filters = make([]*cvm.Filter, 0, len(filters))
	for k, v := range filters {
		values := make([]*string, 0, len(v))
		for _, value := range v {
			values = append(values, helper.String(value))
		}
		filter := &cvm.Filter{
			Name:   helper.String(k),
			Values: values,
		}
		request.Filters = append(request.Filters, filter)
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DescribeInstanceTypeConfigs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instanceTypes = response.Response.InstanceTypeConfigSet
	return
}

func (me *CvmService) DescribeInstancesSellTypeByFilter(ctx context.Context, filters map[string][]string) (instanceTypes []*cvm.InstanceTypeQuotaItem, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeZoneInstanceConfigInfosRequest()
	request.Filters = make([]*cvm.Filter, 0, len(filters))
	for k, v := range filters {
		values := make([]*string, 0, len(v))
		for _, value := range v {
			values = append(values, helper.String(value))
		}
		filter := &cvm.Filter{
			Name:   helper.String(k),
			Values: values,
		}
		request.Filters = append(request.Filters, filter)
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DescribeZoneInstanceConfigInfos(request)
	if err != nil {
		//deal with not supported error
		e, ok := err.(*sdkErrors.TencentCloudSDKError)
		if ok && e.Code == CVM_ZONE_NOT_SUPPORT_ERROR {
			return
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	instanceTypes = response.Response.InstanceTypeQuotaSet
	return
}

func (me *CvmService) DescribeKeyPairById(ctx context.Context, keyId string) (keyPair *cvm.KeyPair, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeKeyPairsRequest()
	request.KeyIds = []*string{&keyId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DescribeKeyPairs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.KeyPairSet) > 0 {
		keyPair = response.Response.KeyPairSet[0]
	}
	return
}

func (me *CvmService) DescribeKeyPairByFilter(ctx context.Context, id, name string, projectId *int) (keyPairs []*cvm.KeyPair, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeKeyPairsRequest()
	if id != "" {
		request.KeyIds = []*string{&id}
	}
	request.Filters = make([]*cvm.Filter, 0)
	if name != "" {
		filter := &cvm.Filter{
			Name:   helper.String("key-name"),
			Values: []*string{&name},
		}
		request.Filters = append(request.Filters, filter)
	}
	if projectId != nil {
		filter := &cvm.Filter{
			Name:   helper.String("project-id"),
			Values: []*string{helper.String(fmt.Sprintf("%d", *projectId))},
		}
		request.Filters = append(request.Filters, filter)
	}

	var offset int64 = 0
	var pageSize int64 = 100
	keyPairs = make([]*cvm.KeyPair, 0)
	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCvmClient().DescribeKeyPairs(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.KeyPairSet) < 1 {
			break
		}
		keyPairs = append(keyPairs, response.Response.KeyPairSet...)
		if len(response.Response.KeyPairSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *CvmService) CreateKeyPair(ctx context.Context, keyName, publicKey string, projectId int64, tags map[string]string) (keyId string, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewImportKeyPairRequest()
	request.KeyName = &keyName
	request.ProjectId = &projectId
	request.PublicKey = &publicKey

	if len(tags) > 0 {
		tagsSpec := make([]*cvm.Tag, 0)
		for tagKey, tagValue := range tags {
			tag := cvm.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			tagsSpec = append(tagsSpec, &tag)
		}
		tagSpecification := cvm.TagSpecification{
			ResourceType: helper.String("keypair"),
			Tags:         tagsSpec,
		}
		request.TagSpecification = append(request.TagSpecification, &tagSpecification)
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ImportKeyPair(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.KeyId == nil || len(*response.Response.KeyId) < 1 {
		errRet = fmt.Errorf("key pair id is nil")
		return
	}
	keyId = *response.Response.KeyId
	return
}

func (me *CvmService) ModifyKeyPairName(ctx context.Context, keyId, keyName string) error {
	logId := getLogId(ctx)
	request := cvm.NewModifyKeyPairAttributeRequest()
	request.KeyId = &keyId
	request.KeyName = &keyName

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyKeyPairAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) DeleteKeyPair(ctx context.Context, keyId string) error {
	logId := getLogId(ctx)
	request := cvm.NewDeleteKeyPairsRequest()
	request.KeyIds = []*string{&keyId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DeleteKeyPairs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) UnbindKeyPair(ctx context.Context, keyId string, instanceIds []*string) error {
	logId := getLogId(ctx)
	request := cvm.NewDisassociateInstancesKeyPairsRequest()
	request.KeyIds = []*string{&keyId}
	request.InstanceIds = instanceIds
	request.ForceStop = helper.Bool(true)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DisassociateInstancesKeyPairs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) BindKeyPair(ctx context.Context, keyId, instanceId string) error {
	logId := getLogId(ctx)
	request := cvm.NewAssociateInstancesKeyPairsRequest()
	request.KeyIds = []*string{&keyId}
	request.InstanceIds = []*string{&instanceId}
	request.ForceStop = helper.Bool(true)

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCvmClient().AssociateInstancesKeyPairs(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}

	return nil
}

func (me *CvmService) CreatePlacementGroup(ctx context.Context, placementName, placementType string) (placementId string, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewCreateDisasterRecoverGroupRequest()
	request.Name = &placementName
	request.Type = &placementType

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().CreateDisasterRecoverGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.DisasterRecoverGroupId == nil {
		errRet = fmt.Errorf("placement group id is nil")
		return
	}
	placementId = *response.Response.DisasterRecoverGroupId
	return
}

func (me *CvmService) DescribePlacementGroupById(ctx context.Context, placementId string) (placementGroup *cvm.DisasterRecoverGroup, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeDisasterRecoverGroupsRequest()
	request.DisasterRecoverGroupIds = []*string{&placementId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DescribeDisasterRecoverGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DisasterRecoverGroupSet) < 1 {
		return
	}
	placementGroup = response.Response.DisasterRecoverGroupSet[0]
	return
}

func (me *CvmService) DescribePlacementGroupByFilter(ctx context.Context, id, name string) (placementGroups []*cvm.DisasterRecoverGroup, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeDisasterRecoverGroupsRequest()
	if id != "" {
		request.DisasterRecoverGroupIds = []*string{&id}
	}
	if name != "" {
		request.Name = &name
	}

	var offset int64 = 0
	var pageSize int64 = 100
	placementGroups = make([]*cvm.DisasterRecoverGroup, 0)
	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCvmClient().DescribeDisasterRecoverGroups(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DisasterRecoverGroupSet) < 1 {
			break
		}
		placementGroups = append(placementGroups, response.Response.DisasterRecoverGroupSet...)
		if len(response.Response.DisasterRecoverGroupSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *CvmService) ModifyPlacementGroup(ctx context.Context, placementId, name string) error {
	logId := getLogId(ctx)
	request := cvm.NewModifyDisasterRecoverGroupAttributeRequest()
	request.DisasterRecoverGroupId = &placementId
	request.Name = &name

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyDisasterRecoverGroupAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) DeletePlacementGroup(ctx context.Context, placementId string) error {
	logId := getLogId(ctx)
	request := cvm.NewDeleteDisasterRecoverGroupsRequest()
	request.DisasterRecoverGroupIds = []*string{&placementId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DeleteDisasterRecoverGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) DescribeRegions(ctx context.Context) (zones []*cvm.RegionInfo, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeRegionsRequest()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DescribeRegions(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	zones = response.Response.RegionSet
	return
}

func (me *CvmService) DescribeZones(ctx context.Context) (zones []*cvm.ZoneInfo, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeZonesRequest()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().DescribeZones(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	zones = response.Response.ZoneSet
	return
}

func (me *CvmService) CreateReservedInstance(ctx context.Context, configId string, count int64, extendParams map[string]interface{}) (instanceId string, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewPurchaseReservedInstancesOfferingRequest()
	request.ReservedInstancesOfferingId = &configId
	request.InstanceCount = &count
	if v, ok := extendParams["reserved_instance_name"]; ok {
		reservedInstanceName := v.(string)
		request.ReservedInstanceName = &reservedInstanceName
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().PurchaseReservedInstancesOffering(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ReservedInstanceId == nil {
		errRet = fmt.Errorf("reserved instance id is nil")
		return
	}
	instanceId = *response.Response.ReservedInstanceId
	return
}

func (me *CvmService) DescribeReservedInstanceByFilter(ctx context.Context, filters map[string]string) (instances []*cvm.ReservedInstances, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeReservedInstancesRequest()
	request.Filters = make([]*cvm.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cvm.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset int64 = 0
	var pageSize int64 = 100
	instances = make([]*cvm.ReservedInstances, 0)
	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCvmClient().DescribeReservedInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ReservedInstancesSet) < 1 {
			break
		}
		instances = append(instances, response.Response.ReservedInstancesSet...)
		if len(response.Response.ReservedInstancesSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *CvmService) DescribeReservedInstanceConfigs(ctx context.Context, filters map[string]string) (configs []*cvm.ReservedInstancesOffering, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeReservedInstancesOfferingsRequest()
	request.Filters = make([]*cvm.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cvm.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, &filter)
	}

	var offset int64 = 0
	var pageSize int64 = 100
	configs = make([]*cvm.ReservedInstancesOffering, 0)
	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCvmClient().DescribeReservedInstancesOfferings(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ReservedInstancesOfferingsSet) < 1 {
			break
		}
		configs = append(configs, response.Response.ReservedInstancesOfferingsSet...)
		if len(response.Response.ReservedInstancesOfferingsSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func flattenCvmTagsMapping(tags []*cvm.Tag) (mapping map[string]string) {
	mapping = make(map[string]string)
	for _, tag := range tags {
		mapping[*tag.Key] = *tag.Value
	}
	return
}

type cvmImages []*cvm.Image

func (a cvmImages) Len() int {
	return len(a)
}

func (a cvmImages) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a cvmImages) Less(i, j int) bool {
	if a[i].CreatedTime == nil || a[j].CreatedTime == nil {
		return false
	}

	itime, _ := time.Parse(time.RFC3339, *a[i].CreatedTime)
	jtime, _ := time.Parse(time.RFC3339, *a[j].CreatedTime)

	return itime.Unix() < jtime.Unix()
}

// Sort images by creation date, in descending order.
func sortImages(images cvmImages) cvmImages {
	sortedImages := images
	sort.Sort(sort.Reverse(sortedImages))
	return sortedImages
}

func (me *CvmService) ModifyImage(ctx context.Context, instanceId, imageName, imageDesc string) (errRet error) {
	logId := getLogId(ctx)
	var request = cvm.NewModifyImageAttributeRequest()
	request.ImageId = helper.String(instanceId)
	request.ImageName = helper.String(imageName)
	request.ImageDescription = helper.String(imageDesc)

	err := resource.Retry(6*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, e := me.client.UseCvmClient().ModifyImageAttribute(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "InvalidImageId.Malformed" || ee.Code == "InvalidImageId.NotFound" ||
					ee.Code == "InvalidImageName.Duplicate" || ee.Code == "InvalidParameter.ValueTooLarge" {
					return resource.NonRetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.RetryableError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (me *CvmService) DeleteImage(ctx context.Context, imageId string) error {
	logId := getLogId(ctx)
	request := cvm.NewDeleteImagesRequest()
	request.ImageIds = []*string{&imageId}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCvmClient().DeleteImages(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}

	return nil
}

func (me *CvmService) DescribeImageById(ctx context.Context, keyId string, isDelete bool) (image *cvm.Image, has bool, errRet error) {
	logId := getLogId(ctx)
	request := cvm.NewDescribeImagesRequest()
	request.ImageIds = []*string{&keyId}

	var imgRsp *cvm.DescribeImagesResponse
	err := resource.Retry(20*readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCvmClient().DescribeImages(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		if response != nil && response.Response != nil {
			if len(response.Response.ImageSet) == 0 && !isDelete {
				return resource.RetryableError(fmt.Errorf("iamge instance status is processing"))
			}
			if len(response.Response.ImageSet) > 0 {
				if *response.Response.ImageSet[0].ImageState == "CREATEFAILED" {
					return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s Create Image is failed", logId))
				}
				if *response.Response.ImageSet[0].ImageState != "NORMAL" {
					return resource.RetryableError(fmt.Errorf("iamge instance status is processing"))
				}
			}

			imgRsp = response
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("response is null"))
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	if imgRsp == nil {
		return
	}

	if len(imgRsp.Response.ImageSet) > 0 && len(imgRsp.Response.ImageSet[0].SnapshotSet) != 0 {
		has = true
		image = imgRsp.Response.ImageSet[0]
	}
	return
}

func (me *CvmService) DescribeImagesByFilter(ctx context.Context, filters map[string][]string, instanceType string) (images []*cvm.Image, errRet error) {
	logId := getLogId(ctx)

	request := cvm.NewDescribeImagesRequest()
	request.Filters = make([]*cvm.Filter, 0, len(filters))
	for k, v := range filters {
		filter := cvm.Filter{
			Name:   helper.String(k),
			Values: []*string{},
		}
		for _, vv := range v {
			filter.Values = append(filter.Values, helper.String(vv))
		}
		request.Filters = append(request.Filters, &filter)
	}
	if instanceType != "" {
		request.InstanceType = helper.String(instanceType)
	}
	var offset uint64 = 0
	var pageSize uint64 = 100
	images = make([]*cvm.Image, 0)
	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCvmClient().DescribeImages(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ImageSet) < 1 {
			break
		}
		images = append(images, response.Response.ImageSet...)
		if len(response.Response.ImageSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	return
}

func (me *CvmService) ModifyRenewParam(ctx context.Context, instanceId string, renewFlag string) error {
	logId := getLogId(ctx)
	request := cvm.NewModifyInstancesRenewFlagRequest()
	request.InstanceIds = []*string{&instanceId}
	request.RenewFlag = &renewFlag

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyInstancesRenewFlag(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) ModifyInstanceChargeType(ctx context.Context, instanceId string, chargeType string, period int,
	renewFlag string) error {
	logId := getLogId(ctx)
	request := cvm.NewModifyInstancesChargeTypeRequest()
	request.InstanceIds = []*string{&instanceId}
	request.InstanceChargeType = &chargeType
	if chargeType == CVM_CHARGE_TYPE_PREPAID {
		request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{}
		if period != -1 {
			request.InstanceChargePrepaid.Period = helper.IntInt64(period)
		}
		if renewFlag != "" {
			request.InstanceChargePrepaid.RenewFlag = &renewFlag
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ModifyInstancesChargeType(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CvmService) ResizeInstanceDisks(ctx context.Context, request *cvm.ResizeInstanceDisksRequest) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCvmClient().ResizeInstanceDisks(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
