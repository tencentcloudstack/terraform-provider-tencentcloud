package tencentcloud

import (
	"context"
	"log"

	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type ElasticsearchService struct {
	client *connectivity.TencentCloudClient
}

func (me *ElasticsearchService) DescribeInstanceById(ctx context.Context, instanceId string) (instance *es.InstanceInfo, errRet error) {
	logId := getLogId(ctx)
	request := es.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseEsClient().DescribeInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	if len(response.Response.InstanceList) < 1 {
		return
	}
	instance = response.Response.InstanceList[0]
	return
}

func (me *ElasticsearchService) DescribeInstancesByFilter(ctx context.Context, instanceId, instanceName string,
	tags map[string]string) (instances []*es.InstanceInfo, errRet error) {

	logId := getLogId(ctx)
	request := es.NewDescribeInstancesRequest()
	if instanceId != "" {
		request.InstanceIds = []*string{&instanceId}
	}
	if instanceName != "" {
		request.InstanceNames = []*string{&instanceName}
	}
	for k, v := range tags {
		tag := es.TagInfo{
			TagKey:   helper.String(k),
			TagValue: helper.String(v),
		}
		request.TagList = append(request.TagList, &tag)
	}

	offset := 0
	pageSize := 100
	instances = make([]*es.InstanceInfo, 0)
	for {
		request.Offset = helper.IntUint64(offset)
		request.Limit = helper.IntUint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseEsClient().DescribeInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		if response == nil || len(response.Response.InstanceList) < 1 {
			break
		}
		instances = append(instances, response.Response.InstanceList...)
		if len(response.Response.InstanceList) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *ElasticsearchService) DeleteInstance(ctx context.Context, instanceId string) error {
	logId := getLogId(ctx)
	request := es.NewDeleteInstanceRequest()
	request.InstanceId = &instanceId

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseEsClient().DeleteInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *ElasticsearchService) UpdateInstance(ctx context.Context, instanceId, instanceName, password string, basicSecurityType int64) error {
	logId := getLogId(ctx)
	request := es.NewUpdateInstanceRequest()
	request.InstanceId = &instanceId
	if instanceName != "" {
		request.InstanceName = &instanceName
	}
	if password != "" {
		request.Password = &password
	}
	if basicSecurityType > 0 {
		request.BasicSecurityType = &basicSecurityType
	}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseEsClient().UpdateInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *ElasticsearchService) UpdateInstanceVersion(ctx context.Context, instanceId, version string) error {
	logId := getLogId(ctx)
	request := es.NewUpgradeInstanceRequest()
	request.InstanceId = &instanceId
	request.EsVersion = &version

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseEsClient().UpgradeInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *ElasticsearchService) UpdateInstanceLicense(ctx context.Context, instanceId, licenseType string) error {
	logId := getLogId(ctx)
	request := es.NewUpgradeLicenseRequest()
	request.InstanceId = &instanceId
	request.LicenseType = &licenseType

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseEsClient().UpgradeLicense(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}
