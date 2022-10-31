/*
Provides a resource to create a vpc bandwidth_package_resources

Example Usage

```hcl
resource "tencentcloud_vpc_bandwidth_package_resources" "bandwidth_package_resources" {
  resource_ids          = "lb-dv1ai6ma"
  bandwidth_package_id  = "bwp-atmf0p9g"
  network_type          = "BGP"
  resource_type         = "LoadBalance"
  protocol              = ""
}

```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcBandwidthPackageResources() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudVpcBandwidthPackageResourcesRead,
		Create: resourceTencentCloudVpcBandwidthPackageResourcesCreate,
		Delete: resourceTencentCloudVpcBandwidthPackageResourcesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_ids": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique ID of the resource, currently supports EIP resources and LB resources, such as `eip-xxxx`, `lb-xxxx`.",
			},

			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Bandwidth package unique ID, in the form of `bwp-xxxx`.",
			},

			"network_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Bandwidth packet type, currently supports `BGP` type, indicating that the internal resource is BGP IP.",
			},

			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Resource types, including `Address`, `LoadBalance`.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Bandwidth packet protocol type. Currently `ipv4` and `ipv6` protocol types are supported.",
			},
		},
	}
}

func resourceTencentCloudVpcBandwidthPackageResourcesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_bandwidth_package_resources.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request            = vpc.NewAddBandwidthPackageResourcesRequest()
		resourceId         string
		bandwidthPackageId string
		resourceType       string
	)

	if v, ok := d.GetOk("resource_ids"); ok {
		resourceId = v.(string)
		request.ResourceIds = append(request.ResourceIds, helper.String(v.(string)))
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		bandwidthPackageId = v.(string)
		request.BandwidthPackageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok {
		request.NetworkType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		resourceType = v.(string)
		request.ResourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {

		request.Protocol = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AddBandwidthPackageResources(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc bandwidthPackageResources failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{bandwidthPackageId, resourceId, resourceType}, FILED_SP))
	return resourceTencentCloudVpcBandwidthPackageResourcesRead(d, meta)
}

func resourceTencentCloudVpcBandwidthPackageResourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_bandwidth_package_resources.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	bandwidthPackageId := ids[0]
	resourceId := ids[1]
	resourceType := ids[2]

	bandwidthPackageResources, err := service.DescribeVpcBandwidthPackageResources(ctx, bandwidthPackageId, resourceId, resourceType)

	if err != nil {
		return err
	}

	if bandwidthPackageResources == nil {
		d.SetId("")
		return fmt.Errorf("resource `bandwidthPackageResources` %s does not exist", resourceId)
	}

	_ = d.Set("bandwidth_package_id", bandwidthPackageId)

	if bandwidthPackageResources.ResourceId != nil {
		_ = d.Set("resource_ids", bandwidthPackageResources.ResourceId)
	}

	//if bandwidthPackageResources.NetworkType != nil {
	//	_ = d.Set("network_type", bandwidthPackageResources.NetworkType)
	//}

	if bandwidthPackageResources.ResourceType != nil {
		_ = d.Set("resource_type", bandwidthPackageResources.ResourceType)
	}

	//if bandwidthPackageResources.Protocol != nil {
	//	_ = d.Set("protocol", bandwidthPackageResources.Protocol)
	//}

	return nil
}

func resourceTencentCloudVpcBandwidthPackageResourcesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_bandwidth_package_resources.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	bandwidthPackageId := ids[0]
	resourceId := ids[1]
	resourceType := ids[2]

	if err := service.DeleteVpcBandwidthPackageResourcesById(ctx, bandwidthPackageId, resourceId, resourceType); err != nil {
		return err
	}

	return nil
}
