package gmvc

import "fmt"

var (
	RequiredChecker = func(ctx GmvcContext, fieldMeta *ParamMeta, value interface{}) error {
		if value == nil {
			return fmt.Errorf("field %s is required", fieldMeta.GetName())
		}

		return nil
	}
)
