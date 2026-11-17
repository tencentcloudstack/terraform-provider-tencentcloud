package cosvectors

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// cosVectorConfigurationObject 解包单项配置，保留 null 和 unknown。
func cosVectorConfigurationObject(value types.List) (types.Object, error) {
	element := value.ElementType(context.Background()).(types.ObjectType)
	if value.IsNull() {
		return types.ObjectNull(element.AttrTypes), nil
	}
	if value.IsUnknown() {
		return types.ObjectUnknown(element.AttrTypes), nil
	}
	if len(value.Elements()) != 1 {
		return types.ObjectNull(element.AttrTypes), errors.New("configuration must contain exactly one block when specified")
	}
	object, ok := value.Elements()[0].(types.Object)
	if !ok || object.IsNull() {
		return types.ObjectNull(element.AttrTypes), errors.New("configuration block must not be null")
	}
	return object, nil
}
