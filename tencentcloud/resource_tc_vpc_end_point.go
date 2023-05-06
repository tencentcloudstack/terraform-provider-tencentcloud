/*
Provides a resource to create a vpc end_point

Example Usage

```hcl
resource "tencentcloud_vpc_end_point" "end_point" {
  vpc_id = "vpc-391sv4w3"
  subnet_id = "subnet-ljyn7h30"
  end_point_name = "terraform-test"
  end_point_service_id = "vpcsvc-69y13tdb"
  end_point_vip = "10.0.2.1"
}
```

Import

vpc end_point can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_end_point.end_point end_point_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcEndPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcEndPointCreate,
		Read:   resourceTencentCloudVpcEndPointRead,
		Update: resourceTencentCloudVpcEndPointUpdate,
		Delete: resourceTencentCloudVpcEndPointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of vpc instance.",
			},

			"subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of subnet instance.",
			},

			"end_point_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of endpoint.",
			},

			"end_point_service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of endpoint service.",
			},

			"end_point_vip": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VIP of endpoint ip.",
			},

			"end_point_owner": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "APPID.",
			},

			"state": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "state of end point.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Create Time.",
			},
		},
	}
}

func resourceTencentCloudVpcEndPointCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_end_point.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = vpc.NewCreateVpcEndPointRequest()
		response   = vpc.NewCreateVpcEndPointResponse()
		endPointId string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_name"); ok {
		request.EndPointName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_service_id"); ok {
		request.EndPointServiceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_vip"); ok {
		request.EndPointVip = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateVpcEndPoint(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc endPoint failed, reason:%+v", logId, err)
		return err
	}

	endPointId = *response.Response.EndPoint.EndPointId
	d.SetId(endPointId)

	return resourceTencentCloudVpcEndPointRead(d, meta)
}

func resourceTencentCloudVpcEndPointRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_end_point.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	endPointId := d.Id()

	endPoint, err := service.DescribeVpcEndPointById(ctx, endPointId)
	if err != nil {
		return err
	}

	if endPoint == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if endPoint.VpcId != nil {
		_ = d.Set("vpc_id", endPoint.VpcId)
	}

	if endPoint.SubnetId != nil {
		_ = d.Set("subnet_id", endPoint.SubnetId)
	}

	if endPoint.EndPointName != nil {
		_ = d.Set("end_point_name", endPoint.EndPointName)
	}

	if endPoint.EndPointServiceId != nil {
		_ = d.Set("end_point_service_id", endPoint.EndPointServiceId)
	}

	if endPoint.EndPointVip != nil {
		_ = d.Set("end_point_vip", endPoint.EndPointVip)
	}

	if endPoint.EndPointOwner != nil {
		_ = d.Set("end_point_owner", endPoint.EndPointOwner)
	}

	if endPoint.State != nil {
		_ = d.Set("state", endPoint.State)
	}

	if endPoint.CreateTime != nil {
		_ = d.Set("create_time", endPoint.CreateTime)
	}

	return nil
}

func resourceTencentCloudVpcEndPointUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_end_point.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := vpc.NewModifyVpcEndPointAttributeRequest()

	endPointId := d.Id()

	request.EndPointId = &endPointId

	unsupportedUpdateFields := []string{
		"vpc_id",
		"subnet_id",
		"end_point_service_id",
		"end_point_vip",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_vpc_end_point update on %s is not support yet", field)
		}
	}

	if d.HasChange("end_point_name") {
		if v, ok := d.GetOk("end_point_name"); ok {
			request.EndPointName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyVpcEndPointAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc endPoint failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpcEndPointRead(d, meta)
}

func resourceTencentCloudVpcEndPointDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_end_point.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	endPointId := d.Id()

	if err := service.DeleteVpcEndPointById(ctx, endPointId); err != nil {
		return nil
	}

	return nil
}
