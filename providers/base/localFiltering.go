package base

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/MichaelPalmer1/scoutr-go/models"
)

type LocalFiltering struct {
	Filtering
	data          map[string]interface{}
	failedFilters []string
}

func (f *LocalFiltering) userFilters(filterFields []models.FilterField) (interface{}, error) {
	// Merge all possible values for this filter key together
	filters := make(map[string][]models.FilterField)
	for _, item := range filterFields {
		filters[item.Field] = append(filters[item.Field], item)
	}

	// Perform the filter
	var conditions interface{} = nil
	var err error
	for key, filterItems := range filters {
		if len(filterItems) == 1 {
			// Perform a single query
			item := filterItems[0]
			existingConditions := conditions
			conditions, err = f.performFilter(conditions, fmt.Sprintf("%s__%s", key, item.Operator), item.Value)
			if err != nil {
				return nil, err
			}
			if (existingConditions == nil || existingConditions == true) && conditions == false {
				f.failedFilters = append(f.failedFilters, key)
			}
		} else if len(filterItems) > 1 {
			// Perform an OR query against all possible values for this key
			var filterConds interface{}
			for _, item := range filterItems {
				result, err := f.performFilter(nil, fmt.Sprintf("%s__%s", key, item.Operator), item.Value)
				if err != nil {
					return nil, err
				}
				filterConds = f.Or(conditions, result)
			}
			existingConditions := conditions
			conditions = f.And(conditions, filterConds)
			if (existingConditions == nil || existingConditions == true) && conditions == false {
				f.failedFilters = append(f.failedFilters, key)
			}
		}
	}

	return conditions, nil
}

func (f *LocalFiltering) And(condition1, condition2 interface{}) bool {
	if condition1 == nil {
		condition1 = true
	}
	if condition2 == nil {
		condition2 = true
	}

	c1 := condition1.(bool)
	c2 := condition2.(bool)

	return c1 && c2
}

func (f *LocalFiltering) Or(condition1, condition2 interface{}) bool {
	if condition1 == nil {
		condition1 = true
	}
	if condition2 == nil {
		condition2 = true
	}

	c1 := condition1.(bool)
	c2 := condition2.(bool)

	return c1 || c2
}

func (f *LocalFiltering) Equals(attr string, value interface{}) bool {
	if val, ok := f.data[attr]; ok {
		return val == value
	} else {
		return false
	}
}

func (f *LocalFiltering) NotEqual(attr string, value interface{}) bool {
	if val, ok := f.data[attr]; ok {
		return val != value
	} else {
		return false
	}
}

func (f *LocalFiltering) Contains(attr string, value interface{}) bool {
	if val, ok := f.data[attr]; ok {
		return strings.Contains(val.(string), value.(string))
	} else {
		return false
	}
}

func (f *LocalFiltering) NotContains(attr string, value interface{}) bool {
	if val, ok := f.data[attr]; ok {
		return !strings.Contains(val.(string), value.(string))
	} else {
		return false
	}
}

func (f *LocalFiltering) StartsWith(attr string, value interface{}) bool {
	if val, ok := f.data[attr]; ok {
		return strings.HasPrefix(val.(string), value.(string))
	} else {
		return false
	}
}

func (f *LocalFiltering) Exists(attr string, value interface{}) bool {
	_, exists := f.data[attr]
	if value == "true" {
		return exists
	} else if value == "false" {
		return !exists
	}

	return false
}

func (f *LocalFiltering) GreaterThan(attr string, value interface{}) bool {
	if val, ok := f.data[attr]; ok {
		val2 := value.(string)

		switch val1 := val.(type) {
		case int64:
			val2, err := strconv.ParseInt(val2, 10, 64)
			if err != nil {
				return false
			}

			return val1 > val2
		case string:
			return val1 > val2
		}
	}

	return false
}

func (f *LocalFiltering) LessThan(attr string, value interface{}) bool {
	if val, ok := f.data[attr]; ok {
		val2 := value.(string)

		switch val1 := val.(type) {
		case int64:
			val2, err := strconv.ParseInt(val2, 10, 64)
			if err != nil {
				return false
			}

			return val1 < val2
		case string:
			return val1 < val2
		}
	}

	return false
}

func (f *LocalFiltering) GreaterThanEqual(attr string, value interface{}) bool {
	if val, ok := f.data[attr]; ok {
		val2 := value.(string)

		switch val1 := val.(type) {
		case int64:
			val2, err := strconv.ParseInt(val2, 10, 64)
			if err != nil {
				return false
			}

			return val1 >= val2
		case string:
			return val1 >= val2
		}
	}

	return false
}

func (f *LocalFiltering) LessThanEqual(attr string, value interface{}) bool {
	if val, ok := f.data[attr]; ok {
		val2 := value.(string)

		switch val1 := val.(type) {
		case int64:
			val2, err := strconv.ParseInt(val2, 10, 64)
			if err != nil {
				return false
			}

			return val1 <= val2
		case string:
			return val1 <= val2
		}
	}

	return false
}

func (f *LocalFiltering) In(attr string, values interface{}) bool {
	var valueList []string
	err := json.Unmarshal([]byte(values.(string)), &valueList)
	if err != nil {
		return false
	}

	if val, ok := f.data[attr]; ok {
		for _, item := range valueList {
			if item == val {
				return true
			}
		}
	}

	return false
}

func (f *LocalFiltering) NotIn(attr string, values interface{}) bool {
	var valueList []string
	err := json.Unmarshal([]byte(values.(string)), &valueList)
	if err != nil {
		return false
	}

	if val, ok := f.data[attr]; ok {
		for _, item := range valueList {
			if item == val {
				return false
			}
		}
	}

	return true
}
