package tencentcloud

import (
	"context"
	"fmt"
	"log"

	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type ClbService struct {
	client *connectivity.TencentCloudClient
}

func (me *ClbService) DescribeLoadBalancerById(ctx context.Context, clbId string) (clbInstance *clb.LoadBalancer, errRet error) {
	logId := GetLogId(ctx)
	request := clb.NewDescribeLoadBalancersRequest()
	request.LoadBalancerIds = []*string{&clbId}

	response, err := me.client.UseClbClient().DescribeLoadBalancers(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LoadBalancerSet) < 1 {
		errRet = fmt.Errorf("loadBalancer id is not found")
		return
	}
	clbInstance = response.Response.LoadBalancerSet[0]
	return
}

func (me *ClbService) DescribeLoadBalancerByFilter(ctx context.Context, params map[string]interface{}) (clbs []*clb.LoadBalancer, errRet error) {
	logId := GetLogId(ctx)
	request := clb.NewDescribeLoadBalancersRequest()

	for k, v := range params {
		if k == "clb_id" {
			request.LoadBalancerIds = []*string{stringToPointer(v.(string))}
		}

		if k == "network_type" {
			request.LoadBalancerType = stringToPointer(v.(string))
		}
		if k == "clb_name" {
			request.LoadBalancerName = stringToPointer(v.(string))
		}
		if k == "project_id" {
			projectId := int64(v.(int))
			request.ProjectId = &projectId
		}

	}

	offset := 0
	pageSize := 100
	clbs = make([]*clb.LoadBalancer, 0)
	for {
		offset64 := int64(offset)
		pageSize64 := int64(pageSize)
		request.Offset = &(offset64)
		request.Limit = &(pageSize64)
		response, err := me.client.UseClbClient().DescribeLoadBalancers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.LoadBalancerSet) < 1 {
			break
		}

		clbs = append(clbs, response.Response.LoadBalancerSet...)

		if len(response.Response.LoadBalancerSet) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *ClbService) DeleteLoadBalancerById(ctx context.Context, clbId string) error {
	logId := GetLogId(ctx)
	request := clb.NewDeleteLoadBalancerRequest()
	request.LoadBalancerIds = []*string{&clbId}
	response, err := me.client.UseClbClient().DeleteLoadBalancer(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func flattenClbTagsMapping(tags []*clb.TagInfo) (mapping map[string]string) {
	mapping = make(map[string]string)
	for _, tag := range tags {
		mapping[*tag.TagKey] = *tag.TagValue
	}
	return
}
