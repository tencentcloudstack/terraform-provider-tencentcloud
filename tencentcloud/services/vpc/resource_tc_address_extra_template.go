package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

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
							Optional:    true,
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
		for _, item := range v.(*schema.Set).List() {
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

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(map[string]interface{})
		request.Tags = make([]*vpc.Tag, 0, len(tags))
		for k, t := range tags {
			key := k
			value := t.(string)
			tag := vpc.Tag{
				Key:   &key,
				Value: &value,
			}
			request.Tags = append(request.Tags, &tag)
		}
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

	if len(template.AddressExtraSet) > 0 {
		addressExtraSets := make([]map[string]interface{}, 0, len(template.AddressExtraSet))
		for _, v := range template.AddressExtraSet {
			addressExtraSet := map[string]interface{}{}
			if v.Address != nil {
				addressExtraSet["address"] = *v.Address
			}
			if v.Description != nil {
				addressExtraSet["description"] = *v.Description
			}
			if v.UpdatedTime != nil {
				addressExtraSet["updated_time"] = *v.UpdatedTime
			}

			addressExtraSets = append(addressExtraSets, addressExtraSet)
		}
		_ = d.Set("addresses_extra", addressExtraSets)
	}

	if len(template.TagSet) > 0 {
		tags := make(map[string]string)
		for _, tag := range template.TagSet {
			if tag.Key != nil && tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			}
		}
		_ = d.Set("tags", tags)
	}
	return nil
}

func resourceTencentCloudAddressExtraTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_address_extra_template.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateId := d.Id()

	if d.HasChange("name") || d.HasChange("addresses_extra") {
		var (
			request = vpc.NewModifyAddressTemplateAttributeRequest()
		)
		request.AddressTemplateId = helper.String(templateId)

		if v, ok := d.GetOk("name"); ok {
			request.AddressTemplateName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("addresses_extra"); ok {
			addressInfos := make([]*vpc.AddressInfo, 0, 10)
			for _, item := range v.(*schema.Set).List() {
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
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyAddressTemplateAttribute(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s modify vpc address template failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(client)
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		region := client.Region

		resourceName := tccommon.BuildTagResourceName("vpc", "address", region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudAddressExtraTemplateRead(d, meta)
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
