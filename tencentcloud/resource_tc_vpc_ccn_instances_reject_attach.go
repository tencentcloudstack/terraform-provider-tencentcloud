/*
Provides a resource to create a vpc ccn_instances_reject_attach

Example Usage

```hcl
resource "tencentcloud_vpc_ccn_instances_reject_attach" "ccn_instances_reject_attach" {
  ccn_id = "ccn-gree226l"
  instances {
		instance_id = "vpc-r1x53ogh"
		instance_region = "ap-guangzhou"
		instance_type = ""
		description = ""
		route_table_id = ""

  }
}
```

Import

vpc ccn_instances_reject_attach can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_ccn_instances_reject_attach.ccn_instances_reject_attach ccn_instances_reject_attach_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudVpcCcnInstancesRejectAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcCcnInstancesRejectAttachCreate,
		Read:   resourceTencentCloudVpcCcnInstancesRejectAttachRead,
		Delete: resourceTencentCloudVpcCcnInstancesRejectAttachDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CCN ID.",
			},

			"instances": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: " .",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attach Instance ID.",
						},
						"instance_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attach Instance Region.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Associated instance type, optional values:&amp;lt;li&amp;gt;`VPC`: private network/li&amp;gt;&amp;lt;li&amp;gt;`DIRECTCONNECT`: private line gateway&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`BMVPC`: Blackstone private network&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`VPNGW`: VPNGW type&amp;lt;/li&amp;gt;.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Remark.",
						},
						"route_table_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The routing table ID associated with the instance.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcCcnInstancesRejectAttachCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ccn_instances_reject_attach.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = vpc.NewRejectAttachCcnInstancesRequest()
		response = vpc.NewRejectAttachCcnInstancesResponse()
		ccnId    string
	)
	if v, ok := d.GetOk("ccn_id"); ok {
		ccnId = v.(string)
		request.CcnId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instances"); ok {
		for _, item := range v.([]interface{}) {
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
				ccnInstance.RouteTableId = helper.String(v.(string))
			}
			request.Instances = append(request.Instances, &ccnInstance)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().RejectAttachCcnInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate vpc ccnInstancesRejectAttach failed, reason:%+v", logId, err)
		return err
	}

	ccnId = *response.Response.CcnId
	d.SetId(ccnId)

	return resourceTencentCloudVpcCcnInstancesRejectAttachRead(d, meta)
}

func resourceTencentCloudVpcCcnInstancesRejectAttachRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ccn_instances_reject_attach.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudVpcCcnInstancesRejectAttachDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ccn_instances_reject_attach.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
