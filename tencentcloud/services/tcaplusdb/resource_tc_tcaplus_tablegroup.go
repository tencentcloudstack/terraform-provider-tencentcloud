package tcaplusdb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	tcaplusdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTcaplusTableGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusTableGroupCreate,
		Read:   resourceTencentCloudTcaplusTableGroupRead,
		Update: resourceTencentCloudTcaplusTableGroupUpdate,
		Delete: resourceTencentCloudTcaplusTableGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the TcaplusDB cluster to which the table group belongs.",
			},
			"tablegroup_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Table group name; may consist of Chinese characters, English letters, or numeric characters, with a maximum length of 32 characters.",
			},
			"table_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "ID of the TcaplusDB table group, can be user-specified (must be unique within the cluster) or auto-incremented by the API when not set. Immutable after creation.",
			},
			"resource_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of table group tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			// Computed values.
			"table_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tables.",
			},
			"total_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total storage size (MB).",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the TcaplusDB table group.",
			},
		},
	}
}

func resourceTencentCloudTcaplusTableGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_tablegroup.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		clusterId    = d.Get("cluster_id").(string)
		groupName    = d.Get("tablegroup_name").(string)
		tableGroupId = d.Get("table_group_id").(string)
		resourceTags []*tcaplusdb.TagInfoUnit
	)

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			resourceTag := tcaplusdb.TagInfoUnit{}
			if v, ok := dMap["tag_key"]; ok {
				resourceTag.TagKey = helper.String(v.(string))
			}

			if v, ok := dMap["tag_value"]; ok {
				resourceTag.TagValue = helper.String(v.(string))
			}

			resourceTags = append(resourceTags, &resourceTag)
		}
	}

	groupId, err := tcaplusService.CreateGroup(ctx, clusterId, groupName, tableGroupId, resourceTags)
	if err != nil {
		return err
	}

	log.Printf("[CRUD] tencentcloud_tcaplus_tablegroup create success, logId=%s, d.Id()=%s", logId, fmt.Sprintf("%s:%s", clusterId, groupId))
	d.SetId(fmt.Sprintf("%s:%s", clusterId, groupId))
	return resourceTencentCloudTcaplusTableGroupRead(d, meta)
}

func resourceTencentCloudTcaplusTableGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_tablegroup.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), ":")
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	tableGroupId := idSplit[1]

	info, has, err := tcaplusService.DescribeGroup(ctx, clusterId, tableGroupId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			info, has, err = tcaplusService.DescribeGroup(ctx, clusterId, tableGroupId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		log.Printf("[CRUD] tencentcloud_tcaplus_tablegroup id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", clusterId)

	if info.TableGroupName != nil {
		_ = d.Set("tablegroup_name", info.TableGroupName)
	}

	if info.TableGroupId != nil {
		_ = d.Set("table_group_id", info.TableGroupId)
	}

	if info.TableCount != nil {
		_ = d.Set("table_count", info.TableCount)
	}

	if info.TotalSize != nil {
		_ = d.Set("total_size", info.TotalSize)
	}

	if info.CreatedTime != nil {
		_ = d.Set("create_time", info.CreatedTime)
	}

	tags, err := tcaplusService.DescribeTableGroupTags(ctx, clusterId, tableGroupId)
	if err != nil {
		return err
	}

	if len(tags) > 0 {
		tagsList := make([]map[string]interface{}, 0, len(tags))
		for _, tag := range tags {
			tagMap := map[string]interface{}{}
			if tag.TagKey != nil {
				tagMap["tag_key"] = *tag.TagKey
			}
			if tag.TagValue != nil {
				tagMap["tag_value"] = *tag.TagValue
			}
			tagsList = append(tagsList, tagMap)
		}

		_ = d.Set("resource_tags", tagsList)
	} else {
		_ = d.Set("resource_tags", nil)
	}

	return nil
}

func resourceTencentCloudTcaplusTableGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_tablegroup.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), ":")
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	tableGroupId := idSplit[1]

	if d.HasChange("tablegroup_name") {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err := tcaplusService.ModifyGroupName(ctx, clusterId, tableGroupId, d.Get("tablegroup_name").(string))
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("resource_tags") {
		o, n := d.GetChange("resource_tags")
		oldTags := o.(*schema.Set).List()
		newTags := n.(*schema.Set).List()

		oldMap := make(map[string]string)
		for _, item := range oldTags {
			m := item.(map[string]interface{})
			oldMap[m["tag_key"].(string)] = m["tag_value"].(string)
		}

		newMap := make(map[string]string)
		for _, item := range newTags {
			m := item.(map[string]interface{})
			newMap[m["tag_key"].(string)] = m["tag_value"].(string)
		}

		var replaceTags []*tcaplusdb.TagInfoUnit
		var deleteTags []*tcaplusdb.TagInfoUnit

		for key, value := range newMap {
			if oldValue, ok := oldMap[key]; !ok || oldValue != value {
				replaceTags = append(replaceTags, &tcaplusdb.TagInfoUnit{
					TagKey:   helper.String(key),
					TagValue: helper.String(value),
				})
			}
		}

		for key := range oldMap {
			if _, ok := newMap[key]; !ok {
				deleteTags = append(deleteTags, &tcaplusdb.TagInfoUnit{
					TagKey: helper.String(key),
				})
			}
		}

		if len(replaceTags) > 0 || len(deleteTags) > 0 {
			var taskId string
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				resp, err := tcaplusService.ModifyTableGroupTags(ctx, clusterId, tableGroupId, replaceTags, deleteTags)
				if err != nil {
					return tccommon.RetryError(err)
				}

				taskId = resp
				return nil
			})

			if err != nil {
				return err
			}

			// wait
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				info, has, err := tcaplusService.DescribeTask(ctx, clusterId, taskId)
				if err != nil {
					return tccommon.RetryError(err)
				}
				if !has {
					return resource.NonRetryableError(fmt.Errorf("Modify table group tags task has been deleted"))
				}

				if *info.Progress == 100 {
					return nil
				}

				if *info.Progress >= 0 {
					return resource.RetryableError(fmt.Errorf("the table group tags modify is in progress, and our wait has timed out"))
				}
				if *info.Progress < 0 {
					return resource.NonRetryableError(fmt.Errorf("TencentCloud SDK return %d task status, modify table group tags task failed", *info.Progress))
				}

				return nil
			})

			if err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudTcaplusTableGroupRead(d, meta)
}

func resourceTencentCloudTcaplusTableGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_tablegroup.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), ":")
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	tableGroupId := idSplit[1]

	err := tcaplusService.DeleteGroup(ctx, clusterId, tableGroupId)
	if err != nil {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err = tcaplusService.DeleteGroup(ctx, clusterId, tableGroupId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	_, has, err := tcaplusService.DescribeGroup(ctx, clusterId, tableGroupId)
	if err != nil || has {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = tcaplusService.DescribeGroup(ctx, clusterId, tableGroupId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			if has {
				err = fmt.Errorf("delete group fail, group still exist from sdk DescribeGroup")
				return resource.RetryableError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete group fail, group still exist from sdk DescribeGroup")
	}
}
