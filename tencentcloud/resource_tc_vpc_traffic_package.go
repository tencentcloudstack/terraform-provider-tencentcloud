/*
Provides a resource to create a vpc traffic_package

Example Usage

```hcl
resource "tencentcloud_vpc_traffic_package" "traffic_package" {
  traffic_amount = 10
      }
```

Import

vpc traffic_package can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_traffic_package.traffic_package traffic_package_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudVpcTrafficPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcTrafficPackageCreate,
		Read:   resourceTencentCloudVpcTrafficPackageRead,
		Update: resourceTencentCloudVpcTrafficPackageUpdate,
		Delete: resourceTencentCloudVpcTrafficPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"traffic_amount": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Traffic Package Amount, eg: 10,20,50,512,1024,5120,51200,60,300,600,3072,6144,30720,61440,307200.",
			},

			"remaining_amount": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "Remaining amount.",
			},

			"used_amount": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "Used amount.",
			},

			"created_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Created time.",
			},
		},
	}
}

func resourceTencentCloudVpcTrafficPackageCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_traffic_package.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = vpc.NewCreateTrafficPackagesRequest()
		response         = vpc.NewCreateTrafficPackagesResponse()
		trafficPackageId string
	)
	if v, ok := d.GetOkExists("traffic_amount"); ok {
		request.TrafficAmount = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateTrafficPackages(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc TrafficPackage failed, reason:%+v", logId, err)
		return err
	}

	trafficPackageId = *response.Response.TrafficPackageId
	d.SetId(trafficPackageId)

	return resourceTencentCloudVpcTrafficPackageRead(d, meta)
}

func resourceTencentCloudVpcTrafficPackageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_traffic_package.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	trafficPackageId := d.Id()

	TrafficPackage, err := service.DescribeVpcTrafficPackageById(ctx, trafficPackageId)
	if err != nil {
		return err
	}

	if TrafficPackage == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcTrafficPackage` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if TrafficPackage.TrafficAmount != nil {
		_ = d.Set("traffic_amount", TrafficPackage.TrafficAmount)
	}

	if TrafficPackage.RemainingAmount != nil {
		_ = d.Set("remaining_amount", TrafficPackage.RemainingAmount)
	}

	if TrafficPackage.UsedAmount != nil {
		_ = d.Set("used_amount", TrafficPackage.UsedAmount)
	}

	if TrafficPackage.CreatedTime != nil {
		_ = d.Set("created_time", TrafficPackage.CreatedTime)
	}

	return nil
}

func resourceTencentCloudVpcTrafficPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_traffic_package.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"traffic_amount", "remaining_amount", "used_amount", "created_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudVpcTrafficPackageRead(d, meta)
}

func resourceTencentCloudVpcTrafficPackageDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_traffic_package.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	trafficPackageId := d.Id()

	if err := service.DeleteVpcTrafficPackageById(ctx, trafficPackageId); err != nil {
		return err
	}

	return nil
}
