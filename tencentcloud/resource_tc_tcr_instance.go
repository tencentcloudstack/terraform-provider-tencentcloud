/*
Use this resource to create tcr instance.

Example Usage

```hcl
resource "tencentcloud_tcr_instance" "foo" {
  name              = "example"
  instance_type		= "basic"

  tags = {
    test = "tf"
  }
}
```

Import

tcr instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_instance.foo cls-cda1iex1
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcrInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrInstanceCreate,
		Read:   resourceTencentCloudTcrInstanceRead,
		Update: resourceTencentCloudTcrInstanceUpdate,
		Delete: resourceTencentCLoudTcrInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the TCR instance.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "TCR types. Valid values are: `standard`, `basic`, `premium`.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "The available tags within this TCR instance.",
			},
			"public_operation": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Control public network access. Valid values are:`Create`, `Delete`.",
			},
			//Computed values
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the TCR instance.",
			},
			"public_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the TCR instance public network access.",
			},
			"public_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public address for access of the TCR instance.",
			},
			"internal_end_point": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal address for access of the TCR instance.",
			},
			"delete_bucket": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate to delete the COS bucket which is auto-created with the instance or not.",
			},
		},
	}
}

func resourceTencentCloudTcrInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name                  = d.Get("name").(string)
		insType               = d.Get("instance_type").(string)
		tags                  = helper.GetTags(d, "tags")
		outErr, inErr         error
		instanceId, operation string
	)

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = tcrService.CreateTCRInstance(ctx, name, insType, tags)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId)

	//check creation done
	err := resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		instance, has, err := tcrService.DescribeTCRInstanceById(ctx, instanceId)
		if err != nil {
			return retryError(err)
		} else if has && *instance.Status == "Running" {
			return nil
		} else if !has {
			return resource.NonRetryableError(fmt.Errorf("create tcr instance fail"))
		} else {
			return resource.RetryableError(fmt.Errorf("creating tcr instance %s , status %s ", instanceId, *instance.Status))
		}
	})

	if err != nil {
		return err
	}

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if v, ok := d.GetOk("public_operation"); ok {
			operation = v.(string)
			inErr = tcrService.ManageTCRExternalEndpoint(ctx, instanceId, operation)
			if inErr != nil {
				return retryError(inErr)
			}
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if v, ok := d.GetOk("public_operation"); ok {
			operation = v.(string)
			inErr = tcrService.ManageTCRExternalEndpoint(ctx, instanceId, operation)
			if inErr != nil {
				return retryError(inErr)
			}
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	return resourceTencentCloudTcrInstanceRead(d, meta)
}

func resourceTencentCloudTcrInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_instance.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var outErr, inErr error
	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	instance, has, outErr := tcrService.DescribeTCRInstanceById(ctx, d.Id())
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			instance, has, inErr = tcrService.DescribeTCRInstanceById(ctx, d.Id())
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

	publicStatus, has, outErr := tcrService.DescribeExternalEndpointStatus(ctx, d.Id())
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			publicStatus, has, inErr = tcrService.DescribeExternalEndpointStatus(ctx, d.Id())
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
	if publicStatus == "Opening" || publicStatus == "Opened" {
		_ = d.Set("public_operation", "Create")
	} else if publicStatus == "Closed" {
		_ = d.Set("public_operation", "Delete")
	}

	_ = d.Set("name", instance.RegistryName)
	_ = d.Set("instance_type", instance.RegistryType)
	_ = d.Set("status", instance.Status)
	_ = d.Set("public_domain", instance.PublicDomain)
	_ = d.Set("internal_end_point", instance.InternalEndpoint)
	_ = d.Set("public_status", publicStatus)

	tags := make(map[string]string, len(instance.TagSpecification.Tags))
	for _, tag := range instance.TagSpecification.Tags {
		tags[*tag.Key] = *tag.Value
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTcrInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_instance.update")()
	//delete_bucket
	return resourceTencentCloudTcrInstanceRead(d, meta)
}

func resourceTencentCLoudTcrInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()
	deleteBucket := d.Get("delete_bucket").(bool)
	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var inErr, outErr error
	var has bool

	outErr = tcrService.DeleteTCRInstance(ctx, instanceId, deleteBucket)
	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = tcrService.DeleteTCRInstance(ctx, instanceId, deleteBucket)
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
		_, has, inErr = tcrService.DescribeTCRInstanceById(ctx, d.Id())
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete tcr instance %s fail, instance still exists from SDK DescribeTcrInstanceById", instanceId)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	return nil
}
