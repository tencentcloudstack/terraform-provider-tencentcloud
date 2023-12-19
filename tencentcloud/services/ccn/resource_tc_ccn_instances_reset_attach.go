package ccn

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCcnInstancesResetAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnInstancesResetAttachCreate,
		Read:   resourceTencentCloudCcnInstancesResetAttachRead,
		Delete: resourceTencentCloudCcnInstancesResetAttachDelete,
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CCN Instance ID.",
			},

			"ccn_uin": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CCN Uin (root account).",
			},

			"instances": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "List Of Attachment Instances.",
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

func resourceTencentCloudCcnInstancesResetAttachCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_ccn_instances_reset_attach.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = vpc.NewResetAttachCcnInstancesRequest()
		ccnId   string
	)
	if v, ok := d.GetOk("ccn_id"); ok {
		ccnId = v.(string)
		request.CcnId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ccn_uin"); ok {
		request.CcnUin = helper.String(v.(string))
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ResetAttachCcnInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("[CRITAL]%s operate vpc ccnInstancesResetAttach failed, reason:%+v", logId, err)
	}

	d.SetId(ccnId)

	return resourceTencentCloudCcnInstancesResetAttachRead(d, meta)
}

func resourceTencentCloudCcnInstancesResetAttachRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_instances_reset_attach.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCcnInstancesResetAttachDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloudccn_instances_reset_attach.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
