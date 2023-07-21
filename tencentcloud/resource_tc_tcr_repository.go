/*
Use this resource to create tcr repository.

Example Usage

Create a tcr repository instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "premium"
  delete_bucket = true
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_example_ns"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_repository" "example" {
  instance_id	 = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  name 	         = "test"
  brief_desc 	 = "111"
  description	 = "111111111111111111111111111111111111"
}
```

Import

tcr repository can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_repository.foo instance_id#namespace_name#repository_name
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudTcrRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrRepositoryCreate,
		Read:   resourceTencentCloudTcrRepositoryRead,
		Update: resourceTencentCloudTcrRepositoryUpdate,
		Delete: resourceTencentCLoudTcrRepositoryDelete,
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
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the TCR namespace.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the TCR repository. Valid length is [2~200]. It can only contain lowercase letters, numbers and separators (`.`, `_`, `-`, `/`), and cannot start, end or continue with separators. Support the use of multi-level address formats, such as `sub1/sub2/repo`.",
			},
			"brief_desc": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 100),
				Description:  "Brief description of the repository. Valid length is [1~100].",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 1000),
				Description:  "Description of the repository. Valid length is [1~1000].",
			},
			//computed
			"is_public": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicate the repository is public or not.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the repository.",
			},
		},
	}
}

func resourceTencentCloudTcrRepositoryCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_repository.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name          = d.Get("name").(string)
		instanceId    = d.Get("instance_id").(string)
		namespaceName = d.Get("namespace_name").(string)
		briefDesc     = d.Get("brief_desc").(string)
		description   = d.Get("description").(string)
		outErr, inErr error
	)

	outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr = tcrService.CreateTCRRepository(ctx, instanceId, namespaceName, name, briefDesc, description)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId + FILED_SP + namespaceName + FILED_SP + name)

	return resourceTencentCloudTcrRepositoryRead(d, meta)
}

func resourceTencentCloudTcrRepositoryUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_repository.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	namespaceName := items[1]
	repositoryName := items[2]

	if d.HasChange("brief_desc") || d.HasChange("description") {
		briefDesc := d.Get("brief_desc").(string)
		description := d.Get("description").(string)
		var outErr, inErr error
		tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
		outErr = tcrService.ModifyTCRRepository(ctx, instanceId, namespaceName, repositoryName, briefDesc, description)
		if outErr != nil {
			outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				inErr = tcrService.ModifyTCRRepository(ctx, instanceId, namespaceName, repositoryName, briefDesc, description)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
		}
		if outErr != nil {
			return outErr
		}
	}

	return resourceTencentCloudTcrRepositoryRead(d, meta)
}

func resourceTencentCloudTcrRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_repository.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	namespaceName := items[1]
	repositoryName := items[2]

	var outErr, inErr error
	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	repository, has, outErr := tcrService.DescribeTCRRepositoryById(ctx, instanceId, namespaceName, repositoryName)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			repository, has, inErr = tcrService.DescribeTCRRepositoryById(ctx, instanceId, namespaceName, repositoryName)
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

	_ = d.Set("name", repositoryName)
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("namespace_name", namespaceName)
	_ = d.Set("create_time", repository.CreationTime)
	_ = d.Set("update_time", repository.UpdateTime)
	_ = d.Set("brief_desc", repository.BriefDescription)
	_ = d.Set("description", repository.Description)
	_ = d.Set("is_public", repository.Public)

	//get public domain
	instance, has, outErr := tcrService.DescribeTCRInstanceById(ctx, instanceId)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			instance, has, inErr = tcrService.DescribeTCRInstanceById(ctx, instanceId)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}

	if has {
		_ = d.Set("url", fmt.Sprintf("%s/%s/%s", *instance.PublicDomain, namespaceName, repositoryName))
	} else {
		return fmt.Errorf("cannot find instance %s", instanceId)
	}

	return nil
}

func resourceTencentCLoudTcrRepositoryDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_repository.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	namespaceName := items[1]
	repositoryName := items[2]

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var inErr, outErr error
	var has bool

	outErr = tcrService.DeleteTCRRepository(ctx, instanceId, namespaceName, repositoryName)
	if outErr != nil {
		outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr = tcrService.DeleteTCRRepository(ctx, instanceId, namespaceName, repositoryName)
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
		_, has, inErr = tcrService.DescribeTCRRepositoryById(ctx, instanceId, namespaceName, repositoryName)
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
