/*
Provides a resource to create a tsf repository

Example Usage

```hcl
resource "tencentcloud_tsf_repository" "repository" {
  repository_name = ""
  repository_type = ""
  bucket_name = ""
  bucket_region = ""
  directory = ""
  repository_desc = ""
}
```

Import

tsf repository can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_repository.repository repository_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfRepositoryCreate,
		Read:   resourceTencentCloudTsfRepositoryRead,
		Update: resourceTencentCloudTsfRepositoryUpdate,
		Delete: resourceTencentCloudTsfRepositoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Warehouse ID.",
			},

			"repository_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "warehouse name.",
			},

			"repository_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "warehouse type (default warehouse: default, private warehouse: private).",
			},

			"bucket_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "the name of the bucket where the warehouse is located.",
			},

			"bucket_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Bucket region where the warehouse is located.",
			},

			"directory": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "directory.",
			},

			"repository_desc": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "warehouse description.",
			},

			"is_used": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "whether the repository is in use.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "warehouse creation time.",
			},
		},
	}
}

func resourceTencentCloudTsfRepositoryCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_repository.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = tsf.NewCreateRepositoryRequest()
		// response     = tsf.NewCreateRepositoryResponse()
		repositoryId string
	)
	if v, ok := d.GetOk("repository_name"); ok {
		request.RepositoryName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("repository_type"); ok {
		request.RepositoryType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bucket_name"); ok {
		request.BucketName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bucket_region"); ok {
		request.BucketRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("directory"); ok {
		request.Directory = helper.String(v.(string))
	}

	if v, ok := d.GetOk("repository_desc"); ok {
		request.RepositoryDesc = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateRepository(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		// response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf repository failed, reason:%+v", logId, err)
		return err
	}

	// repositoryId = *response.Response.RepositoryId
	d.SetId(repositoryId)

	return resourceTencentCloudTsfRepositoryRead(d, meta)
}

func resourceTencentCloudTsfRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_repository.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	repositoryId := d.Id()

	repository, err := service.DescribeTsfRepositoryById(ctx, repositoryId)
	if err != nil {
		return err
	}

	if repository == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfRepository` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if repository.RepositoryId != nil {
		_ = d.Set("repository_id", repository.RepositoryId)
	}

	if repository.RepositoryName != nil {
		_ = d.Set("repository_name", repository.RepositoryName)
	}

	if repository.RepositoryType != nil {
		_ = d.Set("repository_type", repository.RepositoryType)
	}

	if repository.BucketName != nil {
		_ = d.Set("bucket_name", repository.BucketName)
	}

	if repository.BucketRegion != nil {
		_ = d.Set("bucket_region", repository.BucketRegion)
	}

	if repository.Directory != nil {
		_ = d.Set("directory", repository.Directory)
	}

	if repository.RepositoryDesc != nil {
		_ = d.Set("repository_desc", repository.RepositoryDesc)
	}

	if repository.IsUsed != nil {
		_ = d.Set("is_used", repository.IsUsed)
	}

	if repository.CreateTime != nil {
		_ = d.Set("create_time", repository.CreateTime)
	}

	return nil
}

func resourceTencentCloudTsfRepositoryUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_repository.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewUpdateRepositoryRequest()

	repositoryId := d.Id()

	request.RepositoryId = &repositoryId

	immutableArgs := []string{"repository_id", "repository_name", "repository_type", "bucket_name", "bucket_region", "directory", "repository_desc", "is_used", "create_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("repository_desc") {
		if v, ok := d.GetOk("repository_desc"); ok {
			request.RepositoryDesc = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().UpdateRepository(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf repository failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfRepositoryRead(d, meta)
}

func resourceTencentCloudTsfRepositoryDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_repository.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	repositoryId := d.Id()

	if err := service.DeleteTsfRepositoryById(ctx, repositoryId); err != nil {
		return err
	}

	return nil
}
