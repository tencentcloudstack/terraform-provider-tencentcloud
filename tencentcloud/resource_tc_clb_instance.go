/*
Provide a resource to create a CLB instance.

Example Usage

```hcl
resource "tencentcloud_clb_instance" "clblab" {
        network_type         = "OPEN"
        clb_name         = "myclb"
        project_id       = "Default Project"
        vpc_id           = "vpc-abcd1234"
        subnet_id        = "subnet-0agspqdn"
        tags             = "mytags"
        security_groups = ["sg-o0ek7r93"]
        target_region_info_region= "ap-guangzhou"
        target_region_info_vpc_id= "vpc-abcd1234"
		}
}
```

Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb.instance
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbInstanceCreate,
		Read:   resourceTencentCloudClbInstanceRead,
		Update: resourceTencentCloudClbInstanceUpdate,
		Delete: resourceTencentCloudClbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CLB_NETWORK_TYPE),
				Description:  "Type of CLB instance, and available values include 'OPEN' and 'INTERNAL'.",
			},
			"clb_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the CLB to be queried. The name can only contain Chinese characters, English letters, numbers, underscore and hyphen '-'",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "ID of the project to which the instance belongs.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "ID of the subnet within this VPC. The VIP of the intranet CLB instance will be generated from this subnet",
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
				Description:  "ID of the subnet within this VPC. The VIP of the intranet CLB instance will be generated from this subnet",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security groups to which a CLB instance belongs.",
			},
			"target_region_info_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Region information of backend service are attached the CLB instance.",
			},
			"target_region_info_vpc": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Vpc Id information of backend service are attached the CLB instance.",
			},
		},
	}
}

func resourceTencentCloudClbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_clb_instance.create")()
	request := clb.NewCreateLoadBalancerRequest()
	network_type := d.Get("network_type").(string)
	request.LoadBalancerType = stringToPointer(network_type)
	request.LoadBalancerName = stringToPointer(d.Get("clb_name").(string))
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.ProjectId = &projectId
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		if network_type == CLB_NETWORK_TYPE_OPEN {
			return fmt.Errorf("OPEN network_type do not support this operation with subnet_id")
		}
		request.SubnetId = stringToPointer(v.(string))
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().CreateLoadBalancer(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	if len(response.Response.LoadBalancerIds) < 1 {
		return fmt.Errorf("load balancer id is nil")
	}
	d.SetId(*response.Response.LoadBalancerIds[0])
	clbId := *response.Response.LoadBalancerIds[0]

	time.Sleep(3 * time.Second)
	if v, ok := d.GetOk("security_groups"); ok {
		if network_type == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("INTERNAL network_type do not support this operation with sercurity_groups")
		}
		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = stringToPointer(clbId)
		securityGroups := v.([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
		}
		sgResponse, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, sgRequest.GetAction(), sgRequest.ToJsonString(), err.Error())
			return err
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
		}
	}
	if v, ok := d.GetOk("target_region_info_region"); ok {
		region := v.(string)
		vpcId := ""
		if network_type == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("INTERNAL network_type do not support this operation with target_region_info")
		}
		if vv, ok := d.GetOk("target_region_info_vpc"); ok {
			vpcId = vv.(string)
		}

		targetRegionInfo := clb.TargetRegionInfo{
			Region: &region,
			VpcId:  &vpcId,
		}
		mRequest := clb.NewModifyLoadBalancerAttributesRequest()
		mRequest.LoadBalancerId = stringToPointer(clbId)
		mRequest.TargetRegionInfo = &targetRegionInfo
		mResponse, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerAttributes(mRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, mRequest.GetAction(), mRequest.ToJsonString(), err.Error())
			return err
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, mRequest.GetAction(), mRequest.ToJsonString(), mResponse.ToJsonString())
		}

	}
	return resourceTencentCloudClbInstanceRead(d, meta)
}

func resourceTencentCloudClbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_clb_instance.read")()
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbId := d.Id()
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instance, err := clbService.DescribeLoadBalancerById(ctx, clbId)
	if err != nil {
		return err
	}

	d.Set("network_type", instance.LoadBalancerType)
	d.Set("clb_name", instance.LoadBalancerName)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("vpc_id", instance.VpcId)
	d.Set("target_region_info_region", instance.TargetRegionInfo.Region)
	d.Set("target_region_info_vpc", instance.TargetRegionInfo.VpcId)
	d.Set("project_id", instance.ProjectId)
	d.Set("security_groups", flattenStringList(instance.SecureGroups))

	return nil
}

func resourceTencentCloudClbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_clb_instance.update")()
	d.Partial(true)

	clbId := d.Id()
	clbName := ""
	targetRegionInfo := clb.TargetRegionInfo{}
	changed := false

	if d.HasChange("clb_name") {
		changed = true
		clbName = d.Get("clb_name").(string)
	}

	if d.HasChange("target_region_info_region") || d.HasChange("target_region_info_vpc") {
		if d.Get("network_type") == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("INTERNAL network_type do not support this operation with target_region_info")
		}
		changed = true
		region := d.Get("target_region_info_region").(string)
		vpcId := d.Get("target_region_info_vpc").(string)
		targetRegionInfo = clb.TargetRegionInfo{
			Region: &region,
			VpcId:  &vpcId,
		}

	}
	if changed {
		request := clb.NewModifyLoadBalancerAttributesRequest()
		request.LoadBalancerId = stringToPointer(clbId)
		if d.HasChange("clb_name") {
			request.LoadBalancerName = stringToPointer(clbName)
		}
		if d.HasChange("target_region_info_region") || d.HasChange("target_region_info_vpc") {
			request.TargetRegionInfo = &targetRegionInfo
		}
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerAttributes(request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return err
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}
		if d.HasChange("clb_name") {
			d.SetPartial("clb_name")
		}

		if d.HasChange("target_region_info") {
			d.SetPartial("target_region_info")
		}
	}

	if d.HasChange("security_groups") {
		if d.Get("network_type") == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("INTERNAL network_type do not support this operation with sercurity_groups")
		}
		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = stringToPointer(clbId)
		securityGroups := d.Get("security_groups").([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
		}
		sgResponse, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, sgRequest.GetAction(), sgRequest.ToJsonString(), err.Error())
			return err
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
		}
		d.SetPartial("security_groups")
	}
	d.Partial(false)

	return nil
}

func resourceTencentCloudClbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_clb_instance.delete")()
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbId := d.Id()
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := clbService.DeleteLoadBalancerById(ctx, clbId)
	if err != nil {
		log.Printf("[CRITAL]%s reason[%s]\n", logId, err.Error())
		return err
	}
	return nil
}
