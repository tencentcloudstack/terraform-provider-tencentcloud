package cos

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCosBucketLifecycleRuleSchema(t *testing.T) {
	rules := ResourceTencentCloudCosBucket().Schema["lifecycle_rules"].Elem.(*schema.Resource).Schema

	if got := rules["status"].Default; got != s3.ExpirationStatusEnabled {
		t.Fatalf("status default = %v, want %q", got, s3.ExpirationStatusEnabled)
	}
	if _, errors := rules["status"].ValidateFunc("invalid", "status"); len(errors) == 0 {
		t.Fatal("status validation accepted an invalid value")
	}
	if got := rules["filter_tags"].Type; got != schema.TypeMap {
		t.Fatalf("filter_tags type = %v, want TypeMap", got)
	}
}

func TestExpandCosBucketLifecycleRuleFilter(t *testing.T) {
	tests := map[string]struct {
		prefix string
		tags   map[string]interface{}
		want   *s3.LifecycleRuleFilter
	}{
		"all objects": {
			want: &s3.LifecycleRuleFilter{Prefix: aws.String("")},
		},
		"prefix only": {
			prefix: "logs/",
			want:   &s3.LifecycleRuleFilter{Prefix: aws.String("logs/")},
		},
		"single tag only": {
			tags: map[string]interface{}{"env": "prod"},
			want: &s3.LifecycleRuleFilter{Tag: &s3.Tag{Key: aws.String("env"), Value: aws.String("prod")}},
		},
		"prefix and tag": {
			prefix: "logs/",
			tags:   map[string]interface{}{"env": "prod"},
			want: &s3.LifecycleRuleFilter{And: &s3.LifecycleRuleAndOperator{
				Prefix: aws.String("logs/"),
				Tags:   []*s3.Tag{{Key: aws.String("env"), Value: aws.String("prod")}},
			}},
		},
		"multiple tags are stable": {
			tags: map[string]interface{}{"team": "db", "env": "prod"},
			want: &s3.LifecycleRuleFilter{And: &s3.LifecycleRuleAndOperator{Tags: []*s3.Tag{
				{Key: aws.String("env"), Value: aws.String("prod")},
				{Key: aws.String("team"), Value: aws.String("db")},
			}}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := expandCosBucketLifecycleRuleFilter(test.prefix, test.tags); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("expandCosBucketLifecycleRuleFilter() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestFlattenCosBucketLifecycleRuleFilter(t *testing.T) {
	tests := map[string]struct {
		filter     *s3.LifecycleRuleFilter
		wantPrefix string
		wantTags   map[string]interface{}
	}{
		"nil": {
			wantTags: map[string]interface{}{},
		},
		"prefix": {
			filter:     &s3.LifecycleRuleFilter{Prefix: aws.String("logs/")},
			wantPrefix: "logs/",
			wantTags:   map[string]interface{}{},
		},
		"tag": {
			filter:   &s3.LifecycleRuleFilter{Tag: &s3.Tag{Key: aws.String("env"), Value: aws.String("prod")}},
			wantTags: map[string]interface{}{"env": "prod"},
		},
		"and": {
			filter: &s3.LifecycleRuleFilter{And: &s3.LifecycleRuleAndOperator{
				Prefix: aws.String("logs/"),
				Tags: []*s3.Tag{
					{Key: aws.String("env"), Value: aws.String("prod")},
					{Key: aws.String("team"), Value: aws.String("db")},
				},
			}},
			wantPrefix: "logs/",
			wantTags:   map[string]interface{}{"env": "prod", "team": "db"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			prefix, tags := flattenCosBucketLifecycleRuleFilter(test.filter)
			if prefix != test.wantPrefix || !reflect.DeepEqual(tags, test.wantTags) {
				t.Fatalf("flattenCosBucketLifecycleRuleFilter() = (%q, %#v), want (%q, %#v)", prefix, tags, test.wantPrefix, test.wantTags)
			}
		})
	}
}
