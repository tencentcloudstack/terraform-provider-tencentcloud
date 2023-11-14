package tencentcloud

import (
	kubernetes "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type KubernetesService struct {
	client *connectivity.TencentCloudClient
}

func (me *KubernetesService) DescribeKubernetesClusterAuthenticationOptionsByFilter(ctx context.Context, param map[string]interface{}) (clusterAuthenticationOptions []*kubernetes.ServiceAccountAuthenticationOptions, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = kubernetes.NewDescribeClusterAuthenticationOptionsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseKubernetesClient().DescribeClusterAuthenticationOptions(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ServiceAccounts) < 1 {
			break
		}
		clusterAuthenticationOptions = append(clusterAuthenticationOptions, response.Response.ServiceAccounts...)
		if len(response.Response.ServiceAccounts) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
