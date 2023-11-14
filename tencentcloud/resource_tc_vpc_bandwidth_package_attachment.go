/*
Provides a resource to create a vpc bandwidth_package_attachment

Example Usage

```hcl
resource "tencentcloud_vpc_bandwidth_package_attachment" "bandwidth_package_attachment" {
  resource_id = &lt;nil&gt;
  bandwidth_package_id = &lt;nil&gt;
  network_type = &lt;nil&gt;
  resource_type = &lt;nil&gt;
  protocol = &lt;nil&gt;
}
```

Import

vpc bandwidth_package_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_bandwidth_package_attachment.bandwidth_package_attachment bandwidth_package_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudVpcBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcBandwidthPackageAttachmentCreate,
		Read:   resourceTencentCloudVpcBandwidthPackageAttachmentRead,
		Delete: resourceTencentCloudVpcBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the resource, currently supports EIP resources and LB resources, such as `eip-xxxx`, `lb-xxxx`.",
			},

			"bandwidth_package_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Bandwidth package unique ID, in the form of `bwp-xxxx`.",
			},

			"network_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Bandwidth packet type, currently supports `BGP` type, indicating that the internal resource is BGP IP.",
			},

			"resource_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Resource types, including `Address`, `LoadBalance`.",
			},

			"protocol": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Bandwidth packet protocol type. Currently `ipv4` and `ipv6` protocol types are supported.",
			},
		},
	}
}

func resourceTencentCloudVpcBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_bandwidth_package_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request            = vpc.NewAddBandwidthPackageResourcesRequest()
		response           = vpc.NewAddBandwidthPackageResourcesResponse()
		bandwidthPackageId string
		resourceId         string
	)
	if v, ok := d.GetOk("resource_id"); ok {
		resourceId = v.(string)
		request.ResourceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		bandwidthPackageId = v.(string)
		request.BandwidthPackageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok {
		request.NetworkType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc bandwidthPackageAttachment failed, reason:%+v", logId, err)
		return err
	}

	bandwidthPackageId = *response.Response.BandwidthPackageId
	d.SetId(strings.Join([]string{bandwidthPackageId, resourceId}, FILED_SP))

	return resourceTencentCloudVpcBandwidthPackageAttachmentRead(d, meta)
}

func resourceTencentCloudVpcBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_bandwidth_package_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bandwidthPackageId := idSplit[0]
	resourceId := idSplit[1]

	bandwidthPackageAttachment, err := service.DescribeVpcBandwidthPackageAttachmentById(ctx, bandwidthPackageId, resourceId)
	if err != nil {
		return err
	}

	if bandwidthPackageAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcBandwidthPackageAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if bandwidthPackageAttachment.ResourceId != nil {
		_ = d.Set("resource_id", bandwidthPackageAttachment.ResourceId)
	}

	if bandwidthPackageAttachment.BandwidthPackageId != nil {
		_ = d.Set("bandwidth_package_id", bandwidthPackageAttachment.BandwidthPackageId)
	}

	if bandwidthPackageAttachment.NetworkType != nil {
		_ = d.Set("network_type", bandwidthPackageAttachment.NetworkType)
	}

	if bandwidthPackageAttachment.ResourceType != nil {
		_ = d.Set("resource_type", bandwidthPackageAttachment.ResourceType)
	}

	if bandwidthPackageAttachment.Protocol != nil {
		_ = d.Set("protocol", bandwidthPackageAttachment.Protocol)
	}

	return nil
}

func resourceTencentCloudVpcBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_bandwidth_package_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bandwidthPackageId := idSplit[0]
	resourceId := idSplit[1]

	if err := service.DeleteVpcBandwidthPackageAttachmentById(ctx, bandwidthPackageId, resourceId); err != nil {
		return err
	}

	return nil
}
