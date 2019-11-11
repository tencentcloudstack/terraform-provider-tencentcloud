package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/pkg/errors"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TagService struct {
	client *connectivity.TencentCloudClient
}

func (me *TagService) ModifyTags(ctx context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
	request := tag.NewModifyResourceTagsRequest()
	request.Resource = &resourceName
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

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseTagClient().ModifyResourceTags(request); err != nil {
			return retryError(errors.WithStack(err))
		}

		return nil
	})
}

func (me *TagService) DescribeResourceTags(ctx context.Context, serviceType, resourceType, region, resourceId string) (tags map[string]string, err error) {
	request := tag.NewDescribeResourceTagsByResourceIdsRequest()
	request.ServiceType = &serviceType
	request.ResourcePrefix = &resourceType
	request.ResourceRegion = &region
	request.ResourceIds = []*string{&resourceId}
	request.Limit = intToPointer(DESCRIBE_TAGS_LIMIT)

	var offset uint64
	request.Offset = &offset

	// for run loop at least once
	count := DESCRIBE_TAGS_LIMIT
	for count == DESCRIBE_TAGS_LIMIT {
		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseTagClient().DescribeResourceTagsByResourceIds(request)
			if err != nil {
				count = 0

				return retryError(errors.WithStack(err))
			}

			allTags := response.Response.Tags
			count = len(allTags)

			for _, t := range allTags {
				if *t.ResourceId != resourceId {
					continue
				}
				if tags == nil {
					tags = make(map[string]string)
				}

				tags[*t.TagKey] = *t.TagValue
			}

			return nil
		}); err != nil {
			return nil, err
		}

		offset += uint64(count)
	}

	return
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
