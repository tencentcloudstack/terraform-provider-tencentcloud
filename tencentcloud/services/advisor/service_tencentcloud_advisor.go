package advisor

import "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"

func NewAdvisorService(client *connectivity.TencentCloudClient) AdvisorService {
	return AdvisorService{client: client}
}

type AdvisorService struct {
	client *connectivity.TencentCloudClient
}
