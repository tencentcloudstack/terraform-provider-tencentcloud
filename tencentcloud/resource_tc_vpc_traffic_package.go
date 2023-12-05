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

func resourceTencentCloudVpcTrafficPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcTrafficPackageCreate,
		Read:   resourceTencentCloudVpcTrafficPackageRead,
		Delete: resourceTencentCloudVpcTrafficPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"traffic_amount": {
				Required:    true,
				ForceNew:    true,
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
	request.TrafficPackageCount = helper.IntUint64(1)

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

	trafficPackageSet := response.Response.TrafficPackageSet

	if len(trafficPackageSet) < 1 {
		return fmt.Errorf("create traffic package failed.")
	}
	trafficPackageId = *trafficPackageSet[0]

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

	if TrafficPackage.TotalAmount != nil {
		_ = d.Set("traffic_amount", TrafficPackage.TotalAmount)
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
