/*
Provides a resource to create a CLB instance.

Example Usage

INTERNAL CLB

```hcl
resource "tencentcloud_clb_instance" "internal_clb" {
  network_type = "INTERNAL"
  clb_name     = "myclb"
  project_id   = 0
  vpc_id       = "vpc-7007ll7q"
  subnet_id    = "subnet-12rastkr"

  tags = {
    test = "tf"
  }
}
```

OPEN CLB

```hcl
resource "tencentcloud_clb_instance" "open_clb" {
  network_type              = "OPEN"
  clb_name                  = "myclb"
  project_id                = 0
  vpc_id                    = "vpc-da7ffa61"
  security_groups           = ["sg-o0ek7r93"]
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = "vpc-da7ffa61"

  tags = {
    test = "tf"
  }
}
```

Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_instance.foo lb-7a0t6zqb
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var clbActionMu = &sync.Mutex{}

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
				Description:  "Name of the CLB. The name can only contain Chinese characters, English letters, numbers, underscore and hyphen '-'.",
			},
			"clb_vips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The virtual service address table of the CLB.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "Id of the project within the CLB instance, '0' - Default Project.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "VPC id of the CLB.",
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
				Description:  "Subnet id of the CLB. Effective only for CLB within the VPC. Only supports 'INTERNAL' CLBs.",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security groups of the CLB instance. Only supports 'OPEN' CLBs.",
			},
			"target_region_info_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Region information of backend services are attached the CLB instance. Only supports 'OPEN' CLBs.",
			},
			"target_region_info_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Vpc information of backend services are attached the CLB instance. Only supports 'OPEN' CLBs.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "The available tags within this CLB.",
			},
		},
	}
}

func resourceTencentCloudClbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)

	networkType := d.Get("network_type").(string)
	clbName := d.Get("clb_name").(string)
	flag, e := checkSameName(clbName, meta)
	if e != nil {
		return e
	}
	if flag {
		return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB instance][Create] check: Same CLB name %s exists!", clbName)
	}
	targetRegionInfoRegion := ""
	targetRegionInfoVpcId := ""
	if v, ok := d.GetOk("target_region_info_region"); ok {
		targetRegionInfoRegion = v.(string)
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB instance][Create] check: INTERNAL network_type do not support this operation with target_region_info")
		}
	}
	if v, ok := d.GetOk("target_region_info_vpc_id"); ok {
		targetRegionInfoVpcId = v.(string)
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB instance][Create] check: INTERNAL network_type do not support this operation with target_region_info")
		}
	}
	if (targetRegionInfoRegion != "" && targetRegionInfoVpcId == "") || (targetRegionInfoRegion == "" && targetRegionInfoVpcId != "") {
		return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB instance][Create] check: region and vpc_id must be set at same time")
	}
	if _, ok := d.GetOk("security_groups"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB instance][Create] check: INTERNAL network_type do not support this operation with sercurity_groups")
		}
	}
	request := clb.NewCreateLoadBalancerRequest()
	request.LoadBalancerType = helper.String(networkType)
	request.LoadBalancerName = helper.String(clbName)
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.ProjectId = &projectId
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		if networkType == CLB_NETWORK_TYPE_OPEN {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB instance][Create] check: OPEN network_type do not support this operation with subnet_id")
		}
		request.SubnetId = helper.String(v.(string))
	}
	clbId := ""
	var response *clb.CreateLoadBalancerResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().CreateLoadBalancer(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
			if retryErr != nil {
				return retryError(errors.WithStack(retryErr))
			}
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
		return err
	}
	if len(response.Response.LoadBalancerIds) < 1 {
		return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB instance][Create] check: response error, load balancer id is nil")
	}
	d.SetId(*response.Response.LoadBalancerIds[0])
	clbId = *response.Response.LoadBalancerIds[0]

	if v, ok := d.GetOk("security_groups"); ok {
		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = helper.String(clbId)
		securityGroups := v.([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			sgResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
				requestId := *sgResponse.Response.RequestId

				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return retryError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create CLB instance security_groups failed, reason:%+v", logId, err)
			return err
		}
	}

	if targetRegionInfoRegion != "" {
		targetRegionInfo := clb.TargetRegionInfo{
			Region: &targetRegionInfoRegion,
			VpcId:  &targetRegionInfoVpcId,
		}
		mRequest := clb.NewModifyLoadBalancerAttributesRequest()
		mRequest.LoadBalancerId = helper.String(clbId)
		mRequest.TargetRegionInfo = &targetRegionInfo
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			mResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerAttributes(mRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, mRequest.GetAction(), mRequest.ToJsonString(), mResponse.ToJsonString())
				requestId := *mResponse.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return retryError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
			return err
		}
	}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("clb", "loadbalancerid", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudClbInstanceRead(d, meta)
}

func resourceTencentCloudClbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance.read")()
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
		log.Printf("[CRITAL]%s read CLB instance failed, reason:%+v", logId, err)
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("network_type", instance.LoadBalancerType)
	_ = d.Set("clb_name", instance.LoadBalancerName)
	_ = d.Set("clb_vips", helper.StringsInterfaces(instance.LoadBalancerVips))
	_ = d.Set("subnet_id", instance.SubnetId)
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("target_region_info_region", instance.TargetRegionInfo.Region)
	_ = d.Set("target_region_info_vpc_id", instance.TargetRegionInfo.VpcId)
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("security_groups", helper.StringsInterfaces(instance.SecureGroups))

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "clb", "loadbalancerid", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("tags", tags); err != nil {
		log.Printf("[CRITAL]%s provider set tags fail, reason:%s\n ", logId, err.Error())
		return err
	}
	return nil
}

func resourceTencentCloudClbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)

	d.Partial(true)

	clbId := d.Id()
	clbName := ""
	targetRegionInfo := clb.TargetRegionInfo{}
	changed := false

	if d.HasChange("clb_name") {
		changed = true
		clbName = d.Get("clb_name").(string)
		flag, err := checkSameName(clbName, meta)
		if err != nil {
			return err
		}
		if flag {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB instance][Update] check: Same CLB name %s exists!", clbName)
		}
	}

	if d.HasChange("target_region_info_region") || d.HasChange("target_region_info_vpc_id") {
		if d.Get("network_type") == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB instance %s][Update] check: INTERNAL network_type do not support this operation with target_region_info", clbId)
		}
		changed = true
		region := d.Get("target_region_info_region").(string)
		vpcId := d.Get("target_region_info_vpc_id").(string)
		targetRegionInfo = clb.TargetRegionInfo{
			Region: &region,
			VpcId:  &vpcId,
		}
	}

	if changed {
		request := clb.NewModifyLoadBalancerAttributesRequest()
		request.LoadBalancerId = helper.String(clbId)
		if d.HasChange("clb_name") {
			request.LoadBalancerName = helper.String(clbName)
		}
		if d.HasChange("target_region_info_region") || d.HasChange("target_region_info_vpc_id") {
			request.TargetRegionInfo = &targetRegionInfo
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerAttributes(request)
			if e != nil {
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
			log.Printf("[CRITAL]%s update CLB instance failed, reason:%+v", logId, err)
			return err
		}
		if d.HasChange("clb_name") {
			d.SetPartial("clb_name")
		}
		if d.HasChange("clb_vips") {
			d.SetPartial("clb_vips")
		}
		if d.HasChange("target_region_info_region") {
			d.SetPartial("target_region_info_region")
		}
		if d.HasChange("target_region_info_vpc_id") {
			d.SetPartial("target_region_info_vpc_id")
		}
	}

	if d.HasChange("security_groups") {
		if d.Get("network_type") == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][CLB instance %s][Update] check: INTERNAL network_type do not support this operation with sercurity_groups", clbId)
		}
		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = helper.String(clbId)
		securityGroups := d.Get("security_groups").([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			sgResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
				requestId := *sgResponse.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return retryError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB instance security_group failed, reason:%+v", logId, err)
			return err
		}
		d.SetPartial("security_groups")
	}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("clb", "loadbalancerid", region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		d.SetPartial("tags")
	}
	d.Partial(false)

	return nil
}

func resourceTencentCloudClbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance.delete")()

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
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CLB instance failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func checkSameName(name string, meta interface{}) (flag bool, errRet error) {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	flag = false
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	params := make(map[string]interface{})
	params["clb_name"] = name
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		clbs, e := clbService.DescribeLoadBalancerByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		if len(clbs) > 0 {
			flag = true
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB instance failed, reason:%+v", logId, err)
	}
	errRet = err
	return
}
