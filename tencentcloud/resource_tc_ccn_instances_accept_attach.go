/*
Provides a resource to create a vpc ccn_instances_accept_attach, you can use this resource to approve cross-region attachment.

Example Usage

```hcl
resource "tencentcloud_ccn_instances_accept_attach" "ccn_instances_accept_attach" {
  ccn_id = "ccn-39lqkygf"
  instances {
    instance_id = "vpc-j9yhbzpn"
    instance_region = "ap-guangzhou"
    instance_type = "VPC"
  }
}
```
*/
package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCcnInstancesAcceptAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnInstancesAcceptAttachCreate,
		Read:   resourceTencentCloudCcnInstancesAcceptAttachRead,
		Delete: resourceTencentCloudCcnInstancesAcceptAttachDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CCN Instance ID.",
			},

			"instances": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Accept List Of Attachment Instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attachment Instance ID.",
						},
						"instance_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance Region.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "InstanceType: `VPC`, `DIRECTCONNECT`, `BMVPC`, `VPNGW`.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description.",
						},
						"route_table_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the routing table associated with the instance. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCcnInstancesAcceptAttachCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_instances_accept_attach.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = vpc.NewAcceptAttachCcnInstancesRequest()
		ccnId   string
	)
	if v, ok := d.GetOk("ccn_id"); ok {
		ccnId = v.(string)
		request.CcnId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instances"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			ccnInstance := vpc.CcnInstance{}
			if v, ok := dMap["instance_id"]; ok {
				ccnInstance.InstanceId = helper.String(v.(string))
			}
			if v, ok := dMap["instance_region"]; ok {
				ccnInstance.InstanceRegion = helper.String(v.(string))
			}
			if v, ok := dMap["instance_type"]; ok {
				ccnInstance.InstanceType = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				ccnInstance.Description = helper.String(v.(string))
			}
			if v, ok := dMap["route_table_id"]; ok {
				routeTableId := v.(string)
				if routeTableId != "" {
					ccnInstance.RouteTableId = helper.String(v.(string))
				}
			}
			request.Instances = append(request.Instances, &ccnInstance)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AcceptAttachCcnInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("[CRITAL]%s operate vpc ccnInstancesAcceptAttach failed, reason:%+v", logId, err)
	}

	d.SetId(ccnId)

	return resourceTencentCloudCcnInstancesAcceptAttachRead(d, meta)
}

func resourceTencentCloudCcnInstancesAcceptAttachRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_instances_accept_attach.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCcnInstancesAcceptAttachDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ccn_instances_accept_attach.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
