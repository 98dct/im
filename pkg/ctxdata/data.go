package ctxdata

import "context"

func GetUId(c context.Context) string {
	if v, ok := c.Value(Identify).(string); ok {
		return v
	}
	return ""
}
