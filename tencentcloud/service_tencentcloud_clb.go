package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

	offset := int64(0)
	pageSize := int64(100)
	clbs = make([]*clb.LoadBalancer, 0)
	for {
		request.Offset = &(offset)
		request.Limit = &(pageSize)
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

		if int64(len(response.Response.LoadBalancerSet)) < pageSize {
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
	request_id := *response.Response.RequestId
	retryErr := retrySet(request_id, me.client.UseClbClient())
	if retryErr != nil {
		return retryErr
	}
	return nil
}

func (me *ClbService) DescribeListenerById(ctx context.Context, id string) (clbListener *clb.Listener, errRet error) {
	logId := GetLogId(ctx)
	request := clb.NewDescribeListenersRequest()
	items := strings.Split(id, "#")
	if len(items) != 2 {
		errRet = fmt.Errorf("id of resource.tencentcloud_clb_listener is wrong")
		return
	}

	listenerId := items[0]
	clbId := items[1]
	request.ListenerIds = []*string{&listenerId}
	request.LoadBalancerId = stringToPointer(clbId)

	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Listeners) < 1 {
		errRet = fmt.Errorf("Listener id is not found")
		return
	}
	clbListener = response.Response.Listeners[0]
	return
}

func (me *ClbService) DescribeListenersByFilter(ctx context.Context, params map[string]interface{}) (listeners []*clb.Listener, errRet error) {
	logId := GetLogId(ctx)
	request := clb.NewDescribeListenersRequest()
	clbId := ""
	for k, v := range params {
		if k == "listener_id" {
			items := strings.Split(v.(string), "#")
			if len(items) != 2 {
				errRet = fmt.Errorf("id of resource.tencentcloud_clb_listener is wrong")
				return
			}

			listenerId := items[0]
			clbId = items[1]
			request.ListenerIds = []*string{stringToPointer(listenerId)}
			request.LoadBalancerId = stringToPointer(clbId)
		}
		if k == "clb_id" {
			if clbId == "" {
				clbId = v.(string)
				request.LoadBalancerId = stringToPointer(clbId)
			}
		}
		if k == "protocol" {
			request.Protocol = stringToPointer(v.(string))
		}
		if k == "port" {
			port := int64(v.(int))
			request.Port = &port
		}

	}

	listeners = make([]*clb.Listener, 0)

	response, err := me.client.UseClbClient().DescribeListeners(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	listeners = append(listeners, response.Response.Listeners...)

	return
}

func (me *ClbService) DeleteListenerById(ctx context.Context, clbId string, listenerId string) error {
	logId := GetLogId(ctx)
	request := clb.NewDeleteListenerRequest()
	request.ListenerId = stringToPointer(listenerId)
	request.LoadBalancerId = stringToPointer(clbId)
	response, err := me.client.UseClbClient().DeleteListener(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	request_id := *response.Response.RequestId
	retryErr := retrySet(request_id, me.client.UseClbClient())
	if retryErr != nil {
		return retryErr
	}
	return nil
}

func checkHealthCheckPara(ctx context.Context, d *schema.ResourceData, protocol string) (healthSetFlag bool, healthCheckPara *clb.HealthCheck, errRet error) {
	var healthCheck clb.HealthCheck
	healthSetFlag = false
	healthCheckPara = &healthCheck
	if v, ok := d.GetOk("health_check_switch"); ok {
		healthSetFlag = true
		vv := int64(v.(int))
		healthCheck.HealthSwitch = &vv
	}
	if v, ok := d.GetOk("health_check_time_out"); ok {
		healthSetFlag = true
		vv := int64(v.(int))
		healthCheck.TimeOut = &vv
	}

	if v, ok := d.GetOk("health_check_interval_time"); ok {
		healthSetFlag = true
		vv := int64(v.(int))
		healthCheck.IntervalTime = &vv
	}

	if v, ok := d.GetOk("health_check_health_num"); ok {
		healthSetFlag = true
		vv := int64(v.(int))
		healthCheck.HealthNum = &vv
	}
	if v, ok := d.GetOk("health_check_unhealth_num"); ok {
		healthSetFlag = true
		vv := int64(v.(int))
		healthCheck.UnHealthNum = &vv
	}

	if v, ok := d.GetOk("health_check_http_code"); ok {
		//仅适用于HTTP/HTTPS转发规则、TCP监听器的HTTP健康检查方式
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_code can only be set with protocol TCP %s", protocol)
			return
		} else {
			healthSetFlag = true
			vv := int64(v.(int))
			healthCheck.HttpCode = &vv
		}
	}

	if v, ok := d.GetOk("health_check_http_path"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_path can only be set with protocol TCP")
			return
		} else {
			healthSetFlag = true
			healthCheck.HttpCheckPath = stringToPointer(v.(string))
		}
	}

	if v, ok := d.GetOk("health_check_http_domain"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_domain can only be set with protocol TCP")
			return
		} else {
			healthSetFlag = true
			healthCheck.HttpCheckDomain = stringToPointer(v.(string))
		}
	}

	if v, ok := d.GetOk("health_check_http_method"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP) {
			healthSetFlag = false
			errRet = fmt.Errorf("health_check_http_method can only be set with protocol TCP")
			return
		} else {
			healthSetFlag = true
			healthCheck.HttpCheckMethod = stringToPointer(v.(string))
		}

	}

	if healthSetFlag == true {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			healthSetFlag = false
			errRet = fmt.Errorf("health para can only be set with protocol TCP/UDP")
			return
		}
		healthCheckPara = &healthCheck
	}
	return

}

func checkCertificateInputPara(ctx context.Context, d *schema.ResourceData) (certificateSetFlag bool, certPara *clb.CertificateInput, errRet error) {
	certificateSetFlag = false
	var certificateInput clb.CertificateInput
	certificateSSLMode := ""
	certificateId := ""
	certificateCaId := ""

	if v, ok := d.GetOk("certificate_ssl_mode"); ok {
		certificateSetFlag = true
		certificateSSLMode = v.(string)
		certificateInput.SSLMode = stringToPointer(v.(string))
	}

	if v, ok := d.GetOk("certificate_id"); ok {
		certificateSetFlag = true
		certificateId = v.(string)
		certificateInput.CertId = stringToPointer(v.(string))
	}

	if v, ok := d.GetOk("certificate_ca_id"); ok {
		certificateSetFlag = true
		certificateCaId = v.(string)
		certificateInput.CertCaId = stringToPointer(v.(string))
	}

	if certificateSetFlag == true && certificateId == "" {
		certificateSetFlag = false
		errRet = fmt.Errorf("certificatedId is null")
		return
	}

	if certificateSetFlag == true && certificateSSLMode == CERT_SSL_MODE_MUT && certificateCaId == "" {
		certificateSetFlag = false
		errRet = fmt.Errorf("Certificate_ca_key is null and the ssl mode is 'MUTUAL' ")
		return
	}

	certPara = &certificateInput

	return
}
func retrySet(request_id string, meta *clb.Client) (err error) {
	taskQueryRequest := clb.NewDescribeTaskStatusRequest()
	taskQueryRequest.TaskId = &request_id
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		taskResponse, e := meta.DescribeTaskStatus(taskQueryRequest)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if *taskResponse.Response.Status == int64(CLB_TASK_EXPANDING) {
			return resource.RetryableError(fmt.Errorf("clb task status is %d, request_id is %s", *taskResponse.Response.Status, *taskResponse.Response.RequestId))
		}
		return nil
	})
	return
}

func flattenClbTagsMapping(tags []*clb.TagInfo) (mapping map[string]string) {
	mapping = make(map[string]string)
	for _, tag := range tags {
		mapping[*tag.TagKey] = *tag.TagValue
	}
	return
}
