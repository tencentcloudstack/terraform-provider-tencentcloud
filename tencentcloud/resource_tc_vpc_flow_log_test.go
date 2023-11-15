package tencentcloud

import (
	"context"
	"testing"
	"time"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=ap-guangzhou
	resource.AddTestSweepers("ap-guangzhou", &resource.Sweeper{
		Name: "ap-guangzhou",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn

			service := VpcService{client}

			request := vpc.NewDescribeFlowLogsRequest()
			result, err := service.DescribeFlowLogs(ctx, request)
			if err != nil {
				return err
			}

			for i := range result {
				fl := result[i]
				created, err := time.Parse(TENCENTCLOUD_COMMON_TIME_LAYOUT, "*fl.CreatedTime")
				if err != nil {
					created = time.Time{}
				}
				if isResourcePersist(*fl.FlowLogName, &created) {
					continue
				}
				vpcId := ""
				if fl.VpcId != nil {
					vpcId = *fl.VpcId
				}
				_ = service.DeleteVpcFlowLogById(ctx, *fl.FlowLogId, vpcId)
			}
			return nil
		},
	})
}

func TestAccTencentCloudVpcFlowLogResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcFlowLog,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_flow_log.flow_log", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_flow_log.flow_log", "flow_log_name", "iac-test-1"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_flow_log.flow_log", "flow_log_description", "this is a testing flow log"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_flow_log.flow_log",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cloud_log_region",
					"flow_log_storage",
				},
			},
			{
				Config: testAccVpcFlowLogUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_flow_log.flow_log", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_flow_log.flow_log", "flow_log_name", "iac-test-2"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_flow_log.flow_log", "flow_log_description", "updated"),
				),
			},
		},
	})
}

const testAccVpcFlowLog = `

resource "tencentcloud_vpc_flow_log" "flow_log" {
  flow_log_name = "iac-test-1"
  resource_type = "NETWORKINTERFACE"
  resource_id = "eni-qz9wxgmd"
  traffic_type = "ACCEPT"
  vpc_id = "vpc-humgpppd"
  flow_log_description = "this is a testing flow log"
  cloud_log_id = "e6acd27c-365c-4959-8257-751d86657439" # FIXME use data.logsets (not supported) instead
  storage_type = "cls"
}
`

const testAccVpcFlowLogUpdate = `
resource "tencentcloud_vpc_flow_log" "flow_log" {
  flow_log_name = "iac-test-2"
  resource_type = "NETWORKINTERFACE"
  resource_id = "eni-qz9wxgmd"
  traffic_type = "ACCEPT"
  vpc_id = "vpc-humgpppd"
  flow_log_description = "updated"
  cloud_log_id = "e6acd27c-365c-4959-8257-751d86657439" # FIXME use data.logsets (not supported) instead
  storage_type = "cls"
}
`
