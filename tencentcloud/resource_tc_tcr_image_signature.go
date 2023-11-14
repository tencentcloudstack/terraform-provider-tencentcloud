/*
Provides a resource to create a tcr image_signature

Example Usage

```hcl
resource "tencentcloud_tcr_image_signature" "image_signature" {
  registry_id = "tcr-xxx"
  namespace_name = "ns"
  repository_name = "repo"
  image_version = "v1"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr image_signature can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_image_signature.image_signature image_signature_id
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

func resourceTencentCloudTcrImageSignature() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrImageSignatureCreate,
		Read:   resourceTencentCloudTcrImageSignatureRead,
		Delete: resourceTencentCloudTcrImageSignatureDelete,
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

			"namespace_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace name.",
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

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTcrImageSignatureCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_image_signature.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = tcr.NewCreateSignatureRequest()
		response       = tcr.NewCreateSignatureResponse()
		registryId     string
		namespaceName  string
		repositoryName string
		imageVersion   string
	)
	if v, ok := d.GetOk("registry_id"); ok {
		registryId = v.(string)
		request.RegistryId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		namespaceName = v.(string)
		request.NamespaceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("repository_name"); ok {
		repositoryName = v.(string)
		request.RepositoryName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_version"); ok {
		imageVersion = v.(string)
		request.ImageVersion = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTcrClient().CreateSignature(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate tcr ImageSignature failed, reason:%+v", logId, err)
		return err
	}

	registryId = *response.Response.RegistryId
	d.SetId(strings.Join([]string{registryId, namespaceName, repositoryName, imageVersion}, FILED_SP))

	return resourceTencentCloudTcrImageSignatureRead(d, meta)
}

func resourceTencentCloudTcrImageSignatureRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_image_signature.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTcrImageSignatureDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_image_signature.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
