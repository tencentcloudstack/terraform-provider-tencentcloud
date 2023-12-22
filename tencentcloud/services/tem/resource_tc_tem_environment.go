package tem

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTemEnvironment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTemEnvironmentRead,
		Create: resourceTencentCloudTemEnvironmentCreate,
		Update: resourceTencentCloudTemEnvironmentUpdate,
		Delete: resourceTencentCloudTemEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "environment name.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "environment description.",
			},

			"vpc": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vpc ID.",
			},

			"subnet_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "subnet IDs.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "environment tag list.",
			},
		},
	}
}

func resourceTencentCloudTemEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_environment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = tem.NewCreateEnvironmentRequest()
		response *tem.CreateEnvironmentResponse
	)

	if v, ok := d.GetOk("environment_name"); ok {
		request.EnvironmentName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc"); ok {
		request.Vpc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIdsSet := v.(*schema.Set).List()
		for i := range subnetIdsSet {
			subnetIds := subnetIdsSet[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetIds)
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		for key, value := range v.(map[string]interface{}) {
			tag := tem.Tag{
				TagKey:   helper.String(key),
				TagValue: helper.String(value.(string)),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTemClient().CreateEnvironment(request)
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
		log.Printf("[CRITAL]%s create tem environment failed, reason:%+v", logId, err)
		return err
	}

	environmentId := *response.Response.Result

	service := TemService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTemEnvironmentStatus(ctx, environmentId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if *instance.ClusterStatus == "NORMAL" {
			return nil
		}
		if *instance.ClusterStatus == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("environment status is %v, operate failed.", *instance.ClusterStatus))
		}
		return resource.RetryableError(fmt.Errorf("environment status is %v, retry...", *instance.ClusterStatus))
	})
	if err != nil {
		return err
	}

	d.SetId(environmentId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::tem:%s:uin/:environment/%s", region, environmentId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTemEnvironmentRead(d, meta)
}

func resourceTencentCloudTemEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_environment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TemService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	environmentId := d.Id()

	environments, err := service.DescribeTemEnvironment(ctx, environmentId)

	if err != nil {
		return err
	}
	environment := environments.Result
	if environment == nil {
		d.SetId("")
		return fmt.Errorf("resource `environment` %s does not exist", environmentId)
	}

	if environment.EnvironmentName != nil {
		_ = d.Set("environment_name", environment.EnvironmentName)
	}

	if environment.Description != nil {
		_ = d.Set("description", environment.Description)
	}

	if environment.VpcId != nil {
		_ = d.Set("vpc", environment.VpcId)
	}

	if environment.SubnetIds != nil {
		_ = d.Set("subnet_ids", environment.SubnetIds)
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)
	region := client.Region
	tags, err := tagService.DescribeResourceTags(ctx, "tem", "environment", region, environmentId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTemEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_environment.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := tem.NewModifyEnvironmentRequest()

	request.EnvironmentId = helper.String(d.Id())

	if d.HasChange("environment_name") {
		if v, ok := d.GetOk("environment_name"); ok {
			request.EnvironmentName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("vpc") {
		return fmt.Errorf("`vpc` do not support change now.")
	}

	if d.HasChange("subnet_ids") {
		if v, ok := d.GetOk("subnet_ids"); ok {
			subnetIdsSet := v.(*schema.Set).List()
			for i := range subnetIdsSet {
				subnetIds := subnetIdsSet[i].(string)
				request.SubnetIds = append(request.SubnetIds, &subnetIds)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTemClient().ModifyEnvironment(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("tem", "environment", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTemEnvironmentRead(d, meta)
}

func resourceTencentCloudTemEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_environment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TemService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	environmentId := d.Id()

	if err := service.DeleteTemEnvironmentById(ctx, environmentId); err != nil {
		return err
	}

	return nil
}
