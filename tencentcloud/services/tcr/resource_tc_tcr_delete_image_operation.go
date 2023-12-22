package tcr

import (
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcrDeleteImageOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrDeleteImageOperationCreate,
		Read:   resourceTencentCloudTcrDeleteImageOperationRead,
		Delete: resourceTencentCloudTcrDeleteImageOperationDelete,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
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

			"namespace_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "namespace name.",
			},
		},
	}
}

func resourceTencentCloudTcrDeleteImageOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_delete_image_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = tcr.NewDeleteImageRequest()
		registryId     string
		namespaceName  string
		repositoryName string
		imageVersion   string
	)
	if v, ok := d.GetOk("registry_id"); ok {
		request.RegistryId = helper.String(v.(string))
		registryId = v.(string)
	}

	if v, ok := d.GetOk("repository_name"); ok {
		request.RepositoryName = helper.String(v.(string))
		repositoryName = v.(string)
	}

	if v, ok := d.GetOk("image_version"); ok {
		request.ImageVersion = helper.String(v.(string))
		imageVersion = v.(string)
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		request.NamespaceName = helper.String(v.(string))
		namespaceName = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTCRClient().DeleteImage(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate tcr DeleteImageOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{registryId, namespaceName, repositoryName, imageVersion}, tccommon.FILED_SP))

	return resourceTencentCloudTcrDeleteImageOperationRead(d, meta)
}

func resourceTencentCloudTcrDeleteImageOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_delete_image_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTcrDeleteImageOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_delete_image_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
