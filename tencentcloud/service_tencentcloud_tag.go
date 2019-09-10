package tencentcloud

import (
	"context"
	"log"

	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
)

type TagService struct {
	client *connectivity.TencentCloudClient
}

func (me *TagService) ModifyTags(ctx context.Context, resource string, replaceTags map[string]string, deleteKeys []string) error {
	logId := getLogId(ctx)
	request := tag.NewModifyResourceTagsRequest()
	request.Resource = &resource
	if len(replaceTags) > 0 {
		request.ReplaceTags = make([]*tag.Tag, 0, len(replaceTags))
		for k, v := range replaceTags {
			key := k
			value := v
			replaceTag := &tag.Tag{
				TagKey:   &key,
				TagValue: &value,
			}
			request.ReplaceTags = append(request.ReplaceTags, replaceTag)
		}
	}
	if len(deleteKeys) > 0 {
		request.DeleteTags = make([]*tag.TagKeyObject, 0, len(deleteKeys))
		for _, v := range deleteKeys {
			key := v
			deleteKey := &tag.TagKeyObject{
				TagKey: &key,
			}
			request.DeleteTags = append(request.DeleteTags, deleteKey)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTagClient().ModifyResourceTags(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func diffTags(oldTags, newTags map[string]interface{}) (replaceTags map[string]string, deleteTags []string) {
	replaceTags = make(map[string]string)
	deleteTags = make([]string, 0)
	for k, v := range newTags {
		_, ok := oldTags[k]
		if !ok || oldTags[k].(string) != v.(string) {
			replaceTags[k] = v.(string)
		}
	}
	for k := range oldTags {
		_, ok := newTags[k]
		if !ok {
			deleteTags = append(deleteTags, k)
		}
	}
	return
}
