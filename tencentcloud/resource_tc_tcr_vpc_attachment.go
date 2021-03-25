/*
Use this resource to create tcr vpc attachment to manage access of internal endpoint.

Example Usage

```hcl
resource "tencentcloud_tcr_vpc_attachment" "foo" {
  instance_id		= "cls-satg5125"
  vpc_id			= "vpc-asg3sfa3"
  subnet_id		 	= "subnet-1uwh63so"
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
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
)

func resourceTencentCloudTcrVpcAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrVpcAttachmentCreate,
		Read:   resourceTencentCloudTcrVpcAttachmentRead,
		Update: resourceTencentCloudTcrVpcAttachmentUpdate,
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
			"enable_public_domain_dns": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable public domain dns. Default value is `false`.",
			},
			"enable_vpc_domain_dns": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable vpc domain dns. Default value is `false`.",
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

	if enablePublicDomainDns := d.Get("enable_public_domain_dns").(bool); enablePublicDomainDns {
		err := EnableTcrVpcDns(ctx, tcrService, instanceId, vpcId, subnetId, true)
		if err != nil {
			return err
		}
	}

	if enableVpcDomainDns := d.Get("enable_vpc_domain_dns").(bool); enableVpcDomainDns {
		err := EnableTcrVpcDns(ctx, tcrService, instanceId, vpcId, subnetId, false)
		if err != nil {
			return err
		}
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

	if *vpcAccess.AccessIp != "" {
		publicDomainDnsStatus, err := GetDnsStatus(ctx, tcrService, instanceId, vpcId, *vpcAccess.AccessIp, true)
		if err != nil {
			return err
		}
		_ = d.Set("enable_public_domain_dns", *publicDomainDnsStatus.Status == TCR_VPC_DNS_STATUS_ENABLED)

		vpcDomainDnsStatus, err := GetDnsStatus(ctx, tcrService, instanceId, vpcId, *vpcAccess.AccessIp, false)
		if err != nil {
			return err
		}
		_ = d.Set("enable_vpc_domain_dns", *vpcDomainDnsStatus.Status == TCR_VPC_DNS_STATUS_ENABLED)
	}

	return nil
}

func resourceTencentCloudTcrVpcAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_vpc_attachment.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		instanceId = d.Get("instance_id").(string)
		vpcId      = d.Get("vpc_id").(string)
		subnetId   = d.Get("subnet_id").(string)
	)

	d.Partial(true)
	if d.HasChange("enable_public_domain_dns") {
		if isEnabled := d.Get("enable_public_domain_dns").(bool); isEnabled {
			err := EnableTcrVpcDns(ctx, tcrService, instanceId, vpcId, subnetId, true)
			if err != nil {
				return err
			}
		} else {
			err := DisableTcrVpcDns(ctx, tcrService, instanceId, vpcId, subnetId, true)
			if err != nil {
				return err
			}
		}
		d.SetPartial("enable_public_domain_dns")
	}

	if d.HasChange("enable_vpc_domain_dns") {
		if isEnabled := d.Get("enable_vpc_domain_dns").(bool); isEnabled {
			err := EnableTcrVpcDns(ctx, tcrService, instanceId, vpcId, subnetId, false)
			if err != nil {
				return err
			}
		} else {
			err := DisableTcrVpcDns(ctx, tcrService, instanceId, vpcId, subnetId, false)
			if err != nil {
				return err
			}
		}
		d.SetPartial("enable_vpc_domain_dns")
	}
	d.Partial(false)

	return resourceTencentCloudTcrVpcAttachmentRead(d, meta)
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

func WaitForAccessIpExists(ctx context.Context, tcrService TCRService, instanceId string, vpcId string, subnetId string) (accessIp string, errRet error) {
	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, has, inErr := tcrService.DescribeTCRVPCAttachmentById(ctx, instanceId, vpcId, subnetId)
		if inErr != nil {
			return retryError(inErr)
		}
		if !has {
			inErr = fmt.Errorf("%s create tcr vpcAccess %s fail, vpcAccess is not exists from SDK DescribeTcrVpcAttachmentById", instanceId, vpcId)
			return resource.RetryableError(inErr)
		}

		if *result.AccessIp == "" {
			inErr = fmt.Errorf("%s get tcr accessIp fail, accessIp is not exists from SDK DescribeTcrVpcAttachmentById", vpcId)
			return resource.RetryableError(inErr)
		}
		accessIp = *result.AccessIp
		return nil
	})
	return
}

func EnableTcrVpcDns(ctx context.Context, tcrService TCRService, instanceId string, vpcId string, subnetId string, usePublicDomain bool) error {
	accessIp, err := WaitForAccessIpExists(ctx, tcrService, instanceId, vpcId, subnetId)
	if err != nil {
		return err
	}

	outErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr := tcrService.CreateTcrVpcDns(ctx, instanceId, vpcId, accessIp, usePublicDomain)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})

	return outErr
}

func DisableTcrVpcDns(ctx context.Context, tcrService TCRService, instanceId string, vpcId string, subnetId string, usePublicDomain bool) error {
	accessIp, err := WaitForAccessIpExists(ctx, tcrService, instanceId, vpcId, subnetId)
	if err != nil {
		return err
	}

	outErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr := tcrService.DeleteTcrVpcDns(ctx, instanceId, vpcId, accessIp, usePublicDomain)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})

	return outErr
}

func GetDnsStatus(ctx context.Context, tcrService TCRService, instanceId string, vpcId string, accessIp string, usePublicDomain bool) (status *tcr.VpcPrivateDomainStatus, err error) {
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, has, inErr := tcrService.DescribeTcrVpcDnsById(ctx, instanceId, vpcId, accessIp, usePublicDomain)
		if inErr != nil {
			return retryError(inErr)
		}
		if !has {
			inErr = fmt.Errorf("%s get tcr vpc dns status fail, vpc dns is not exists from SDK DescribeTcrVpcDnsById", instanceId)
			return resource.RetryableError(inErr)
		}
		status = result
		return nil
	})

	return
}
