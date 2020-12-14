/*
Use this resource to create tcr vpc attachment to manage access of internal endpoint.

Example Usage

```hcl
resource "tencentcloud_tcr_vpc_attachment" "foo" {
  instance_id		= ""
  name              = "example"
  is_public		 	= true
}
```

Import

tcr vpc attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_vpc_attachment.foo cls-cda1iex1#vpcAccess
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudTcrVpcAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrVpcAttachmentCreate,
		Read:   resourceTencentCloudTcrVpcAttachmentRead,
		Delete: resourceTencentCLoudTcrVpcAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the TCR instance.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of VPC.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of subnet.",
			},
			//computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the internal access.",
			},
			"access_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the internal access.",
			},
		},
	}
}

func resourceTencentCloudTcrVpcAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_vpc_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		instanceId    = d.Get("instance_id").(string)
		vpcId         = d.Get("vpc_id").(string)
		subnetId      = d.Get("subnet_id").(string)
		outErr, inErr error
		has           bool
	)

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr = tcrService.CreateTCRVPCAttachment(ctx, instanceId, vpcId, subnetId)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId + FILED_SP + vpcId + FILED_SP + subnetId)

	//check exist
	//the attachment takes effect with a minute
	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr = tcrService.DescribeTCRVPCAttachmentById(ctx, instanceId, vpcId, subnetId)
		if inErr != nil {
			return retryError(inErr)
		}
		if !has {
			inErr = fmt.Errorf("create tcr vpcAccess %s fail, vpcAccess is not exists from SDK DescribeTcrVpcAttachmentById", instanceId)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	return resourceTencentCloudTcrVpcAttachmentRead(d, meta)
}

func resourceTencentCloudTcrVpcAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_vpc_attachment.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	vpcId := items[1]
	subnetId := items[2]

	var outErr, inErr error
	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	vpcAccess, has, outErr := tcrService.DescribeTCRVPCAttachmentById(ctx, instanceId, vpcId, subnetId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			vpcAccess, has, inErr = tcrService.DescribeTCRVPCAttachmentById(ctx, instanceId, vpcId, subnetId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("status", vpcAccess.Status)
	_ = d.Set("access_ip", vpcAccess.AccessIp)
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("vpc_id", vpcId)
	_ = d.Set("subnet_id", subnetId)

	return nil
}

func resourceTencentCLoudTcrVpcAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_vpc_attachment.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	vpcId := items[1]
	subnetId := items[2]

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var inErr, outErr error
	var has bool

	outErr = tcrService.DeleteTCRVPCAttachment(ctx, instanceId, vpcId, subnetId)
	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = tcrService.DeleteTCRVPCAttachment(ctx, instanceId, vpcId, subnetId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, inErr = tcrService.DescribeTCRVPCAttachmentById(ctx, instanceId, vpcId, subnetId)
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete tcr vpcAccess %s fail, vpcAccess still exists from SDK DescribeTcrVpcAttachmentById", resourceId)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	return nil
}
