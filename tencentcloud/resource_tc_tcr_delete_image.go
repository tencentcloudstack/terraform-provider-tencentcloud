/*
Provides a resource to create a tcr delete_image

Example Usage

```hcl
resource "tencentcloud_tcr_delete_image" "delete_image" {
  registry_id = "tcr-xxx"
  repository_name = "repo"
  image_version = "v1"
  namespace_name = "ns"
}
```

Import

tcr delete_image can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_delete_image.delete_image delete_image_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTcrDeleteImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrDeleteImageCreate,
		Read:   resourceTencentCloudTcrDeleteImageRead,
		Delete: resourceTencentCloudTcrDeleteImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"repository_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Repository name.",
			},

			"image_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image version name.",
			},

			"namespace_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace name.",
			},
		},
	}
}

func resourceTencentCloudTcrDeleteImageCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_delete_image.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = tcr.NewDeleteImageRequest()
		response       = tcr.NewDeleteImageResponse()
		registryId     string
		namespaceName  string
		repositoryName string
		imageVersion   string
	)
	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
		request.RegistryId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("repository_name"); ok {
		repositoryName = v.(string)
		request.RepositoryName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_version"); ok {
		imageVersion = v.(string)
		request.ImageVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		namespaceName = v.(string)
		request.NamespaceName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().DeleteImage(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate tcr DeleteImage failed, reason:%+v", logId, err)
		return err
	}

	registryId = *response.Response.RegistryId
	d.SetId(strings.Join([]string{registryId, namespaceName, repositoryName, imageVersion}, FILED_SP))

	return resourceTencentCloudTcrDeleteImageRead(d, meta)
}

func resourceTencentCloudTcrDeleteImageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_delete_image.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTcrDeleteImageDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_delete_image.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
