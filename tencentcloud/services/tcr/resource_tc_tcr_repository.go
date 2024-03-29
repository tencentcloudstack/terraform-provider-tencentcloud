package tcr

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTcrRepository() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 100),
				Description:  "Brief description of the repository. Valid length is [1~100].",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 1000),
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
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_repository.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcrService := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		name          = d.Get("name").(string)
		instanceId    = d.Get("instance_id").(string)
		namespaceName = d.Get("namespace_name").(string)
		briefDesc     = d.Get("brief_desc").(string)
		description   = d.Get("description").(string)
		outErr, inErr error
	)

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr = tcrService.CreateTCRRepository(ctx, instanceId, namespaceName, name, briefDesc, description)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId + tccommon.FILED_SP + namespaceName + tccommon.FILED_SP + name)

	return resourceTencentCloudTcrRepositoryRead(d, meta)
}

func resourceTencentCloudTcrRepositoryUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_repository.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
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
		tcrService := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr = tcrService.ModifyTCRRepository(ctx, instanceId, namespaceName, repositoryName, briefDesc, description)
		if outErr != nil {
			outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				inErr = tcrService.ModifyTCRRepository(ctx, instanceId, namespaceName, repositoryName, briefDesc, description)
				if inErr != nil {
					return tccommon.RetryError(inErr)
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
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_repository.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	namespaceName := items[1]
	repositoryName := items[2]

	var outErr, inErr error
	tcrService := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	repository, has, outErr := tcrService.DescribeTCRRepositoryById(ctx, instanceId, namespaceName, repositoryName)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			repository, has, inErr = tcrService.DescribeTCRRepositoryById(ctx, instanceId, namespaceName, repositoryName)
			if inErr != nil {
				return tccommon.RetryError(inErr)
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
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, has, inErr = tcrService.DescribeTCRInstanceById(ctx, instanceId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
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
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_repository.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", resourceId)
	}

	instanceId := items[0]
	namespaceName := items[1]
	repositoryName := items[2]

	tcrService := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var inErr, outErr error
	var has bool

	outErr = tcrService.DeleteTCRRepository(ctx, instanceId, namespaceName, repositoryName)
	if outErr != nil {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = tcrService.DeleteTCRRepository(ctx, instanceId, namespaceName, repositoryName)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr = tcrService.DescribeTCRRepositoryById(ctx, instanceId, namespaceName, repositoryName)
		if inErr != nil {
			return tccommon.RetryError(inErr)
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
