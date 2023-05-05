/*
Provides a resource to create a tcr image_signature_operation

Example Usage

```hcl
resource "tencentcloud_tcr_create_image_signature_operation" "image_signature_operation" {
  registry_id = "tcr-xxx"
  namespace_name = "ns"
  repository_name = "repo"
  image_version = "v1"

}
```

Import

tcr image_signature_operation can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_create_image_signature_operation.image_signature_operation image_signature_operation_id
```
*/
package tencentcloud

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTcrCreateImageSignatureOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrCreateImageSignatureOperationCreate,
		Read:   resourceTencentCloudTcrCreateImageSignatureOperationRead,
		Delete: resourceTencentCloudTcrCreateImageSignatureOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"namespace_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "namespace name.",
			},

			"repository_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "repository name.",
			},

			"image_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "image version name.",
			},
		},
	}
}

func resourceTencentCloudTcrCreateImageSignatureOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_create_image_signature_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = tcr.NewCreateSignatureRequest()
		registryId     string
		namespaceName  string
		repositoryName string
		imageVersion   string
	)
	if v, ok := d.GetOk("registry_id"); ok {
		request.RegistryId = helper.String(v.(string))
		registryId = v.(string)
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		request.NamespaceName = helper.String(v.(string))
		namespaceName = v.(string)
	}

	if v, ok := d.GetOk("repository_name"); ok {
		request.RepositoryName = helper.String(v.(string))
		repositoryName = v.(string)
	}

	if v, ok := d.GetOk("image_version"); ok {
		request.ImageVersion = helper.String(v.(string))
		imageVersion = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTCRClient().CreateSignature(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate tcr ImageSignatureOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{registryId, namespaceName, repositoryName, imageVersion}, FILED_SP))

	return resourceTencentCloudTcrCreateImageSignatureOperationRead(d, meta)
}

func resourceTencentCloudTcrCreateImageSignatureOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_create_image_signature_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTcrCreateImageSignatureOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcr_create_image_signature_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
