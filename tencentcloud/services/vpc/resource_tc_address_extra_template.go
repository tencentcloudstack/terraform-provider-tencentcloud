package vpc

import (
	"context"
	"fmt"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudAddressExtraTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAddressExtraTemplateCreate,
		Read:   resourceTencentCloudAddressExtraTemplateRead,
		Update: resourceTencentCloudAddressExtraTemplateUpdate,
		Delete: resourceTencentCloudAddressExtraTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "IP address template name.",
			},
			"addresses_extra": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address.",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Remarks.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update Time.",
						},
					},
				},
				Description: "The address information can contain remarks and be presented by the IP, CIDR block or IP address range.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the Addresses.",
			},
		},
	}
}

func resourceTencentCloudAddressExtraTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_address_extra_template.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = vpc.NewCreateAddressTemplateRequest()
		response *vpc.CreateAddressTemplateResponse
	)

	if v, ok := d.GetOk("name"); ok {
		request.AddressTemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("addresses_extra"); ok {
		addressInfos := make([]*vpc.AddressInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			addressInfo := vpc.AddressInfo{}
			if v, ok := dMap["address"]; ok {
				addressInfo.Address = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				addressInfo.Description = helper.String(v.(string))
			}
			addressInfos = append(addressInfos, &addressInfo)
		}
		request.AddressesExtra = addressInfos
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateAddressTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc address template failed, reason:%+v", logId, err)
		return err
	}

	templateId := *response.Response.AddressTemplate.AddressTemplateId
	d.SetId(templateId)

	return resourceTencentCloudAddressExtraTemplateRead(d, meta)
}

func resourceTencentCloudAddressExtraTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_address_extra_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	templateId := d.Id()
	var outErr, inErr error
	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	template, has, outErr := vpcService.DescribeAddressTemplateById(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			template, has, inErr = vpcService.DescribeAddressTemplateById(ctx, templateId)
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

	_ = d.Set("name", template.AddressTemplateName)
	_ = d.Set("addresses", template.AddressSet)

	return nil
}

func resourceTencentCloudAddressExtraTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_address_extra_template.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateId := d.Id()

	if d.HasChange("name") || d.HasChange("addresses") {
		var outErr, inErr error
		name := d.Get("name").(string)
		addresses := d.Get("addresses").(*schema.Set).List()
		vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.ModifyAddressTemplate(ctx, templateId, name, addresses)
			if inErr != nil {
				return tccommon.RetryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}

	return resourceTencentCloudAddressTemplateRead(d, meta)
}

func resourceTencentCloudAddressExtraTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_address_extra_template.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateId := d.Id()
	vpcService := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var outErr, inErr error

	outErr = vpcService.DeleteAddressTemplate(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = vpcService.DeleteAddressTemplate(ctx, templateId)
			if inErr != nil {
				return tccommon.RetryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	//check not exist
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr := vpcService.DescribeAddressTemplateById(ctx, templateId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("address template %s is still exists, retry...", templateId))
		} else {
			return nil
		}
	})

	return outErr
}
