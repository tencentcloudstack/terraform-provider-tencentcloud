/*
Provides a resource to create a CLB instance.

Example Usage

INTERNAL CLB

```hcl
resource "tencentcloud_clb_instance" "internal_clb" {
  network_type              = "INTERNAL"
  clb_name                  = "myclb"
  project_id                = 0
  vpc_id                    = "vpc-7007ll7q"
  subnet_id                 = "subnet-12rastkr"

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

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
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
		return fmt.Errorf("Same clb name exists!")
	}
	targetRegionInfoRegion := ""
	targetRegionInfoVpcId := ""
	if v, ok := d.GetOk("target_region_info_region"); ok {
		targetRegionInfoRegion = v.(string)
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("INTERNAL network_type do not support this operation with target_region_info")
		}
	}
	if v, ok := d.GetOk("target_region_info_vpc_id"); ok {
		targetRegionInfoVpcId = v.(string)
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("INTERNAL network_type do not support this operation with target_region_info")
		}
	}
	if (targetRegionInfoRegion != "" && targetRegionInfoVpcId == "") || (targetRegionInfoRegion == "" && targetRegionInfoVpcId != "") {
		return fmt.Errorf("region and vpc_id must be set at same time")
	}
	if _, ok := d.GetOk("security_groups"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("INTERNAL network_type do not support this operation with sercurity_groups")
		}
	}
	request := clb.NewCreateLoadBalancerRequest()
	request.LoadBalancerType = stringToPointer(networkType)
	request.LoadBalancerName = stringToPointer(clbName)
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.ProjectId = &projectId
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		if networkType == CLB_NETWORK_TYPE_OPEN {
			return fmt.Errorf("OPEN network_type do not support this operation with subnet_id")
		}
		request.SubnetId = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := v.(map[string]interface{})
		request.Tags = make([]*clb.TagInfo, 0, len(tags))
		for key, value := range tags {
			tag := clb.TagInfo{
				TagKey:   stringToPointer(key),
				TagValue: stringToPointer(value.(string)),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}
	clbId := ""
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
	clbId = *response.Response.LoadBalancerIds[0]

	if v, ok := d.GetOk("security_groups"); ok {
		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = stringToPointer(clbId)
		securityGroups := v.([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			sgResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
				requestId := *sgResponse.Response.RequestId

				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return retryError(retryErr)
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create clb instance failed, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	if targetRegionInfoRegion != "" {

		targetRegionInfo := clb.TargetRegionInfo{
			Region: &targetRegionInfoRegion,
			VpcId:  &targetRegionInfoVpcId,
		}
		mRequest := clb.NewModifyLoadBalancerAttributesRequest()
		mRequest.LoadBalancerId = stringToPointer(clbId)
		mRequest.TargetRegionInfo = &targetRegionInfo
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			mResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ModifyLoadBalancerAttributes(mRequest)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, mRequest.GetAction(), mRequest.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, mRequest.GetAction(), mRequest.ToJsonString(), mResponse.ToJsonString())
				requestId := *mResponse.Response.RequestId

				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return retryError(retryErr)
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create clb instance failed, reason:%s\n ", logId, err.Error())
			return err
		}
	}
	return resourceTencentCloudClbInstanceRead(d, meta)
}

func resourceTencentCloudClbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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

	d.Set("network_type", instance.LoadBalancerType)
	d.Set("clb_name", instance.LoadBalancerName)
	d.Set("clb_vips", flattenStringList(instance.LoadBalancerVips))
	d.Set("subnet_id", instance.SubnetId)
	d.Set("vpc_id", instance.VpcId)
	d.Set("target_region_info_region", instance.TargetRegionInfo.Region)
	d.Set("target_region_info_vpc_id", instance.TargetRegionInfo.VpcId)
	d.Set("project_id", instance.ProjectId)
	d.Set("security_groups", flattenStringList(instance.SecureGroups))
	d.Set("tags", flattenClbTagsMapping(instance.Tags))
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
			return fmt.Errorf("Same clb name exists!")
		}
	}

	if d.HasChange("target_region_info_region") || d.HasChange("target_region_info_vpc_id") {
		if d.Get("network_type") == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("INTERNAL network_type do not support this operation with target_region_info")
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
		request.LoadBalancerId = stringToPointer(clbId)
		if d.HasChange("clb_name") {
			request.LoadBalancerName = stringToPointer(clbName)
		}
		if d.HasChange("target_region_info_region") || d.HasChange("target_region_info_vpc_id") {
			request.TargetRegionInfo = &targetRegionInfo
		}
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
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			sgResponse, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
				requestId := *sgResponse.Response.RequestId

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
		d.SetPartial("security_groups")
	}
	d.Partial(false)

	return nil
}

func resourceTencentCloudClbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_instance.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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

func checkSameName(name string, meta interface{}) (flag bool, errRet error) {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
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
		log.Printf("[CRITAL]%s read clb instance failed, reason:%s\n ", logId, err.Error())
	}
	errRet = err
	return
}
