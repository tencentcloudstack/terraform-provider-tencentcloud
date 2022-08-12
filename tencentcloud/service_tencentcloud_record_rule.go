package tencentcloud

import (
	"context"
	"fmt"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"

	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type RecordRuleService struct {
	client *connectivity.TencentCloudClient
}

func (rrs *RecordRuleService) DeletePrometheusRecordRuleYaml(ctx context.Context, id, name string) (errRet error) {
	logId := getLogId(ctx)
	request := tke.NewDeletePrometheusRecordRuleYamlRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &id
	request.Names = []*string{&name}

	ratelimit.Check(request.GetAction())
	response, err := rrs.client.UseTkeClient().DeletePrometheusRecordRuleYaml(request)
	if err != nil {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		return err
	}

	return
}

func (rrs *RecordRuleService) DescribePrometheusRecordRuleByName(ctx context.Context, id, name string) (
	ret *tke.DescribePrometheusRecordRulesResponse, errRet error) {

	logId := getLogId(ctx)
	request := tke.NewDescribePrometheusRecordRulesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &id
	request.Filters = []*tke.Filter{
		{
			Name:   helper.String("name"),
			Values: []*string{&name},
		},
	}

	response, err := rrs.client.UseTkeClient().DescribePrometheusRecordRules(request)

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	return response, nil
}
