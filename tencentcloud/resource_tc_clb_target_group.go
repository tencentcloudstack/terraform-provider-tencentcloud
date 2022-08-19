/*
Provides a resource to create a CLB target group.

Example Usage

```hcl
resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test"
    port              = 33
}
```


Import

CLB target group can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group.test lbtg-3k3io0i0
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbTargetGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbTargetCreate,
		Read:   resourceTencentCloudClbTargetRead,
		Update: resourceTencentCloudClbTargetUpdate,
		Delete: resourceTencentCloudClbTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "TF_target_group",
				Description: "Target group name.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
				ForceNew:    true,
				Description: "VPC ID, default is based on the network.",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validatePort,
				Description:  "The default port of target group, add server after can use it.",
			},
			"target_group_instances": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The backend server of target group bind.",
				Deprecated: "It has been deprecated from version 1.77.3. " +
					"please use `tencentcloud_clb_target_group_instance_attachment` instead.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIp,
							Description:  "The internal ip of target group instance.",
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validatePort,
							Description:  "The port of target group instance.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The weight of target group instance.",
						},
						"new_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validatePort,
							Description:  "The new port of target group instance.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClbTargetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group.create")()

	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		clbService      = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		vpcId           = d.Get("vpc_id").(string)
		targetGroupName = d.Get("target_group_name").(string)
		port            = uint64(d.Get("port").(int))
		insAttachments  = make([]*clb.TargetGroupInstance, 0)
		targetGroupId   string
		err             error
	)

	if v, ok := d.GetOk("target_group_instances"); ok {
		targetGroupInstances := v.([]interface{})
		for _, v1 := range targetGroupInstances {
			value := v1.(map[string]interface{})
			bindIP := value["bind_ip"].(string)
			port := uint64(value["port"].(int))
			weight := uint64(value["weight"].(int))
			newPort := uint64(value["new_port"].(int))
			tgtGrp := &clb.TargetGroupInstance{
				BindIP:  &bindIP,
				Port:    &port,
				Weight:  &weight,
				NewPort: &newPort,
			}
			insAttachments = append(insAttachments, tgtGrp)
		}
	}

	targetGroupId, err = clbService.CreateTargetGroup(ctx, targetGroupName, vpcId, port, insAttachments)
	if err != nil {
		return err
	}
	d.SetId(targetGroupId)

	return resourceTencentCloudClbTargetRead(d, meta)

}

func resourceTencentCloudClbTargetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		clbService = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		id         = d.Id()
	)
	filters := make(map[string]string)
	targetGroupInfos, err := clbService.DescribeTargetGroups(ctx, id, filters)
	if err != nil {
		return err
	}
	if len(targetGroupInfos) < 1 {
		d.SetId("")
		return nil
	}
	_ = d.Set("target_group_name", targetGroupInfos[0].TargetGroupName)
	_ = d.Set("vpc_id", targetGroupInfos[0].VpcId)
	_ = d.Set("port", targetGroupInfos[0].Port)

	return nil
}

func resourceTencentCloudClbTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group.update")()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		clbService    = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		targetGroupId = d.Id()
		port          uint64
		tgtGroupName  string
	)

	isChanged := false
	if d.HasChange("port") || d.HasChange("target_group_name") {
		isChanged = true
		port = uint64(d.Get("port").(int))
		tgtGroupName = d.Get("target_group_name").(string)
	}

	if isChanged {
		err := clbService.ModifyTargetGroup(ctx, targetGroupId, tgtGroupName, port)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudClbTargetRead(d, meta)
}

func resourceTencentCloudClbTargetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_target_group.delete")()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		clbService    = ClbService{client: meta.(*TencentCloudClient).apiV3Conn}
		targetGroupId = d.Id()
	)

	err := clbService.DeleteTarget(ctx, targetGroupId)

	if err != nil {
		return err
	}
	return nil
}
