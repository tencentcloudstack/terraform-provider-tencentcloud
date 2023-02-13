/*
Use this resource to create tcr namespace.

Example Usage

```hcl
resource "tencentcloud_tcr_namespace" "foo" {
  instance_id		= ""
  name          	= "example"
  is_public		 	= true
  is_auto_scan		= true
  is_prevent_vul	= true
  severity			= "medium"
  cve_whitelist_items	{
    cve_id = "cve-xxxxx"
  }
}
```

Import

tcr namespace can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_namespace.foo cls-cda1iex1#namespace
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

func resourceTencentCloudTcrNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrNamespaceCreate,
		Read:   resourceTencentCloudTcrNamespaceRead,
		Update: resourceTencentCloudTcrNamespaceUpdate,
		Delete: resourceTencentCLoudTcrNamespaceDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the TCR namespace. Valid length is [2~30]. It can only contain lowercase letters, numbers and separators (`.`, `_`, `-`), and cannot start, end or continue with separators.",
			},
			"is_public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate that the namespace is public or not. Default is `false`.",
			},
			"is_auto_scan": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Scanning level, `True` is automatic, `False` is manual. Default is `false`.",
			},

			"is_prevent_vul": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Blocking switch, `True` is open, `False` is closed. Default is `false`.",
			},

			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Block vulnerability level, currently only supports `low`, `medium`, `high`.",
			},

			"cve_whitelist_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Vulnerability Whitelist.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cve_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Vulnerability Whitelist ID.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTcrNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_namespace.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name           = d.Get("name").(string)
		instanceId     = d.Get("instance_id").(string)
		isPublic       = d.Get("is_public").(bool)
		isAutoScan     = d.Get("is_auto_scan").(bool)
		isPreventVUL   = d.Get("is_prevent_vul").(bool)
		severity       = d.Get("severity").(string)
		whitelistItems = d.Get("cve_whitelist_items").([]interface{})
		outErr, inErr  error
	)

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr = tcrService.CreateTCRNameSpace(ctx, instanceId, name, isPublic, isAutoScan, isPreventVUL, severity, whitelistItems)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId + FILED_SP + name)

	return resourceTencentCloudTcrNamespaceRead(d, meta)
}

func resourceTencentCloudTcrNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_namespace.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	namespaceName := items[1]

	var (
		isPublic       = d.Get("is_public").(bool)
		isAutoScan     = d.Get("is_auto_scan").(bool)
		isPreventVUL   = d.Get("is_prevent_vul").(bool)
		severity       = d.Get("severity").(string)
		whitelistItems = d.Get("cve_whitelist_items").([]interface{})
		outErr, inErr  error
	)

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	outErr = tcrService.ModifyTCRNameSpace(ctx, instanceId, namespaceName, isPublic, isAutoScan, isPreventVUL, severity, whitelistItems)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			inErr = tcrService.ModifyTCRNameSpace(ctx, instanceId, namespaceName, isPublic, isAutoScan, isPreventVUL, severity, whitelistItems)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}

	return resourceTencentCloudTcrNamespaceRead(d, meta)
}

func resourceTencentCloudTcrNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_namespace.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	namespaceName := items[1]

	var outErr, inErr error
	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	namespace, has, outErr := tcrService.DescribeTCRNameSpaceById(ctx, instanceId, namespaceName)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			namespace, has, inErr = tcrService.DescribeTCRNameSpaceById(ctx, instanceId, namespaceName)
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

	_ = d.Set("name", namespace.Name)
	_ = d.Set("is_public", namespace.Public)
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("is_auto_scan", namespace.AutoScan)
	_ = d.Set("is_prevent_vul", namespace.PreventVUL)
	_ = d.Set("severity", namespace.Severity)

	whiteList := []interface{}{}
	if namespace.CVEWhitelistItems != nil {
		for _, v := range namespace.CVEWhitelistItems {
			cveMap := map[string]interface{}{}
			if v.CVEID != nil {
				cveMap["cve_id"] = v.CVEID
			}

			whiteList = append(whiteList, cveMap)
		}
	}
	_ = d.Set("cve_whitelist_items", whiteList)

	return nil
}

func resourceTencentCLoudTcrNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_namespace.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	namespaceName := items[1]

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var inErr, outErr error
	var has bool

	outErr = tcrService.DeleteTCRNameSpace(ctx, instanceId, namespaceName)
	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = tcrService.DeleteTCRNameSpace(ctx, instanceId, namespaceName)
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
		_, has, inErr = tcrService.DescribeTCRNameSpaceById(ctx, instanceId, namespaceName)
		if inErr != nil {
			return retryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete tcr namespace %s fail, namespace still exists from SDK DescribeTcrNamespaceById", resourceId)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	return nil
}
