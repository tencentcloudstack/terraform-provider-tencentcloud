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

Using public network access whitelist
```hcl
resource "tencentcloud_tcr_instance" "foo" {
  name                  = "example"
  instance_type		    = "basic"
  open_public_operation = true
  security_policy {
    cidr_block = "10.0.0.1/24"
  }
  security_policy {
    cidr_block = "192.168.1.1"
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
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcrInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrInstanceCreate,
		Read:   resourceTencentCloudTcrInstanceRead,
		Update: resourceTencentCloudTcrInstanceUpdate,
		Delete: resourceTencentCloudTcrInstanceDelete,
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
				Description: "The available tags within this TCR instance.",
			},
			"open_public_operation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Control public network access.",
			},
			"security_policy": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Public network access allowlist policies of the TCR instance. Only available when `open_public_operation` is `true`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The public network IP address of the access source.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Remarks of policy.",
						},
						"index": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Index of policy.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of policy.",
						},
					},
				},
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

	client := meta.(*TencentCloudClient).apiV3Conn
	tcrService := TCRService{client: client}

	var (
		name           = d.Get("name").(string)
		insType        = d.Get("instance_type").(string)
		outErr, inErr  error
		instanceId     string
		instanceStatus string
		operation      = d.Get("open_public_operation").(bool)
	)

	// Check if security_policy but open_public_operation is false
	if _, ok := d.GetOk("security_policy"); ok && !operation {
		return fmt.Errorf("`open_public_operation` must be `true` if `security_policy` set")
	}

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = tcrService.CreateTCRInstance(ctx, name, insType, map[string]string{})
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
			instanceStatus = "Running"
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
	if instanceStatus == "Running" {
		openPublicOperation, ok := d.GetOk("open_public_operation")
		operation = openPublicOperation.(bool)

		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if ok {
				if operation {
					inErr = tcrService.ManageTCRExternalEndpoint(ctx, instanceId, "Create")
				} else {
					inErr = tcrService.ManageTCRExternalEndpoint(ctx, instanceId, "Delete")
				}
				if inErr != nil {
					return retryError(inErr)
				}
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

		if raw, ok := d.GetOk("security_policy"); ok && operation {
			// Waiting for External EndPoint opened
			err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
				var (
					status string
				)
				status, _, err = tcrService.DescribeExternalEndpointStatus(ctx, instanceId)
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("an error occured during DescribeExternalEndpointStatus: %s", err.Error()))
				}

				if status == "Opened" {
					return nil
				}

				if status == "Opening" {
					return resource.RetryableError(fmt.Errorf("external endpoint status is `%s`, retrying", status))
				}

				return resource.NonRetryableError(fmt.Errorf("unexpected external endpoint status: `%s`", status))
			})

			if err != nil {
				return err
			}
			if err := resourceTencentCloudTcrSecurityPolicyAdd(d, meta, raw.(*schema.Set).List()); err != nil {
				return err
			}
		} else if !operation {
			log.Printf("[WARN] `open_public_operation` was not opened, skip `security_policy` set.")
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("tcr", "instance", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTcrInstanceRead(d, meta)
}

func resourceTencentCloudTcrInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_instance.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var outErr, inErr error
	client := meta.(*TencentCloudClient).apiV3Conn
	tcrService := TCRService{client: client}
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
		_ = d.Set("open_public_operation", true)
	} else if publicStatus == "Closed" {
		_ = d.Set("open_public_operation", false)
	}

	_ = d.Set("name", instance.RegistryName)
	_ = d.Set("instance_type", instance.RegistryType)
	_ = d.Set("status", instance.Status)
	_ = d.Set("public_domain", instance.PublicDomain)
	_ = d.Set("internal_end_point", instance.InternalEndpoint)
	_ = d.Set("public_status", publicStatus)

	request := tcr.NewDescribeSecurityPoliciesRequest()
	request.RegistryId = helper.String(d.Id())
	response, err := client.UseTCRClient().DescribeSecurityPolicies(request)
	if err == nil {
		if response.Response.SecurityPolicySet != nil {
			securityPolicySet := response.Response.SecurityPolicySet
			policies := make([]interface{}, 0, len(securityPolicySet))
			for i := range securityPolicySet {
				item := securityPolicySet[i]
				policy := make(map[string]interface{})
				policy["cidr_block"] = *item.CidrBlock
				policy["description"] = *item.Description
				policy["index"] = *item.PolicyIndex
				policy["version"] = *item.PolicyVersion
				policies = append(policies, policy)
			}
			if err := d.Set("security_policy", policies); err != nil {
				return err
			}
		}
	} else {
		_ = d.Set("security_policy", make([]interface{}, 0))
		log.Printf("[WARN] %s error: %s", request.GetAction(), err.Error())
	}

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
	var (
		outErr, inErr error
		instanceId    string
		operation     bool
	)

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId = d.Id()
	if d.HasChange("open_public_operation") {
		operation = d.Get("open_public_operation").(bool)
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if operation {
				inErr = tcrService.ManageTCRExternalEndpoint(ctx, instanceId, "Create")
			} else {
				inErr = tcrService.ManageTCRExternalEndpoint(ctx, instanceId, "Delete")
			}
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
	}

	if d.HasChange("security_policy") {
		var err error
		// Waiting for External EndPoint opened
		err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
			var (
				status string
			)
			status, _, err = tcrService.DescribeExternalEndpointStatus(ctx, instanceId)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("an error occured during DescribeExternalEndpointStatus: %s", err.Error()))
			}

			if status == "Opened" {
				return nil
			}

			if status == "Opening" {
				return resource.RetryableError(fmt.Errorf("external endpoint status is `%s`, retrying", status))
			}

			return resource.NonRetryableError(fmt.Errorf("unexpected external endpoint status: `%s`", status))
		})

		if err != nil {
			return err
		}

		o, n := d.GetChange("security_policy")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()
		if len(remove) > 0 {
			err := resourceTencentCloudTcrSecurityPolicyRemove(d, meta, remove)
			if err != nil {
				return err
			}
		}
		if len(add) > 0 {
			err := resourceTencentCloudTcrSecurityPolicyAdd(d, meta, add)
			if err != nil {
				return err
			}
		}
		d.SetPartial("security_policy")
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("tcr", "instance", region, d.Id())
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		d.SetPartial("tags")
	}
	return resourceTencentCloudTcrInstanceRead(d, meta)
}

func resourceTencentCloudTcrInstanceDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTencentCloudTcrSecurityPolicyAdd(d *schema.ResourceData, meta interface{}, add []interface{}) error {
	client := meta.(*TencentCloudClient).apiV3Conn
	request := tcr.NewCreateMultipleSecurityPolicyRequest()
	request.RegistryId = helper.String(d.Id())

	for _, i := range add {
		dMap := i.(map[string]interface{})
		policy := &tcr.SecurityPolicy{}
		if cidr, ok := dMap["cidr_block"]; ok {
			policy.CidrBlock = helper.String(cidr.(string))
		}
		if desc, ok := dMap["description"]; ok {
			policy.Description = helper.String(desc.(string))
		}
		if index, ok := dMap["index"]; ok {
			policy.PolicyIndex = helper.IntInt64(index.(int))
		}
		if version, ok := dMap["version"]; ok {
			policy.PolicyVersion = helper.String(version.(string))
		}
		request.SecurityGroupPolicySet = append(request.SecurityGroupPolicySet, policy)
	}

	_, err := client.UseTCRClient().CreateMultipleSecurityPolicy(request)
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudTcrSecurityPolicyRemove(d *schema.ResourceData, meta interface{}, remove []interface{}) error {
	client := meta.(*TencentCloudClient).apiV3Conn
	request := tcr.NewDeleteMultipleSecurityPolicyRequest()
	request.RegistryId = helper.String(d.Id())

	for _, i := range remove {
		dMap := i.(map[string]interface{})
		policy := &tcr.SecurityPolicy{}
		if cidr, ok := dMap["cidr_block"]; ok {
			policy.CidrBlock = helper.String(cidr.(string))
		}
		if desc, ok := dMap["description"]; ok {
			policy.Description = helper.String(desc.(string))
		}
		if index, ok := dMap["index"]; ok {
			policy.PolicyIndex = helper.IntInt64(index.(int))
		}
		if version, ok := dMap["version"]; ok {
			policy.PolicyVersion = helper.String(version.(string))
		}
		request.SecurityGroupPolicySet = append(request.SecurityGroupPolicySet, policy)
	}

	_, err := client.UseTCRClient().DeleteMultipleSecurityPolicy(request)
	if err != nil {
		return err
	}
	return nil
}
