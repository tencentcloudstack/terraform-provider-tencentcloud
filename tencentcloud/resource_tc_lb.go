/*
Provides a Load Balancer resource.

~> **NOTE:** It has been deprecated and replaced by `tencentcloud_clb_instance`.

Example Usage

```hcl
resource "tencentcloud_lb" "classic" {
  type       = "OPEN"
  forward    = "APPLICATION"
  name       = "tf-test-classic"
  project_id = 0
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const (
	lbNetworkTypeOpen        = "OPEN"
	lbNetworkTypeInternal    = "INTERNAL"
	lbForwardTypeClassic     = "CLASSIC"
	lbForwardTypeApplication = "APPLICATION"
)

func resourceTencentCloudLB() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.15.0. Please use 'tencentcloud_clb_instance' instead.",
		Create:             resourceTencentCloudLBCreate,
		Read:               resourceTencentCloudLBRead,
		Update:             resourceTencentCloudLBUpdate,
		Delete:             resourceTencentCloudLBDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{lbNetworkTypeOpen, lbNetworkTypeInternal}),
				Description:  "The network type of the LB. Valid value: 'OPEN', 'INTERNAL'.",
			},
			"forward": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{lbForwardTypeClassic, lbForwardTypeApplication}),
				Description:  "The type of the LB. Valid value: 'CLASSIC', 'APPLICATION'.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the LB.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The VPC ID of the LB, unspecified or 0 stands for CVM basic network.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The project id of the LB, unspecified or 0 stands for default project.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the LB.",
			},
		},
	}
}

func resourceTencentCloudLBCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	request := clb.NewCreateLoadBalancerRequest()

	networkType := d.Get("type").(string)
	request.LoadBalancerType = helper.String(networkType)

	if v, ok := d.GetOk("forward"); ok {
		if v == lbForwardTypeClassic {
			request.Forward = helper.IntInt64(0)
		} else {
			request.Forward = helper.IntInt64(1)
		}
	}

	if v, ok := d.GetOk("name"); ok {
		request.LoadBalancerName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.ProjectId = &projectId
	}

	var response *clb.CreateLoadBalancerResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().CreateLoadBalancer(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			requestId := *result.Response.RequestId

			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return retryError(retryErr)
			}
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create clb instance failed, reason:%s\n ", logId, err.Error())
		return err
	}
	if len(response.Response.LoadBalancerIds) < 1 {
		return fmt.Errorf("load balancer id is nil")
	}
	d.SetId(*response.Response.LoadBalancerIds[0])

	return resourceTencentCloudLBRead(d, meta)
}

func resourceTencentCloudLBRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbId := d.Id()
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *clb.LoadBalancer
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeLoadBalancerById(ctx, clbId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read clb instance failed, reason:%s\n ", logId, err.Error())
		return err
	}

	_ = d.Set("type", instance.LoadBalancerType)
	_ = d.Set("name", instance.LoadBalancerName)
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("project_id", instance.ProjectId)

	if *instance.Forward == 0 {
		_ = d.Set("forward", lbForwardTypeClassic)
	} else {
		_ = d.Set("forward", lbForwardTypeApplication)
	}

	if *instance.Status == 0 {
		_ = d.Set("status", "CREATING")
	} else {
		_ = d.Set("status", "NORMAL")
	}

	return nil
}

func resourceTencentCloudLBUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)

	clbId := d.Id()
	if !d.HasChange("name") {
		return nil
	}

	d.Partial(true)

	request := clb.NewModifyLoadBalancerAttributesRequest()
	request.LoadBalancerId = helper.String(clbId)
	request.LoadBalancerName = helper.String(d.Get("name").(string))

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerAttributes(request)

		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			requestId := *response.Response.RequestId

			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return retryError(retryErr)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update clb instance failed, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetPartial("name")
	d.Partial(false)

	return resourceTencentCloudLBRead(d, meta)
}

func resourceTencentCloudLBDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbId := d.Id()
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteLoadBalancerById(ctx, clbId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete clb instance failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
