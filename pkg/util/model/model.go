package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrOperatorNotSupported -
	ErrOperatorNotSupported = errors.New("operator not supported")
)

// SQLFromParamsAndOperator -
func SQLFromParamsAndOperator(initialSQL string, params map[string]interface{},
	paramOperators map[string]string, alias string) (string, error) {

	// NOTE: If this is going to be used outside of our model methods
	//       that currently append `WHERE deleted_at IS NULL`, we may
	//       want to start this by checking if there is currently a WHERE
	//       clause or not because if there isn't the first AND clause will
	//       become a WHERE clause instead.

	sqlStmt := initialSQL

	// Copy the params into a new map
	newParams := map[string]interface{}{}
	for k, v := range params {
		newParams[k] = v
	}

	for param, val := range newParams {
		operator, found := paramOperators[param]
		if !found {
			if alias != "" {
				sqlStmt += fmt.Sprintf("AND %s.%s = :%s\n", alias, param, param)
			} else {
				sqlStmt += fmt.Sprintf("AND %s = :%s\n", param, param)
			}
			continue
		}

		// NOTE: If we want later to be able to support multiple
		//       operators at a time for a single param, we could
		//       do something like split operators by a comma and for each
		//       operator, work out what needs to happen to the param since
		//       at the moment, only 1 operator call per param will work.

		// NOTE: Until we implement the above not, we may still want to
		//       be able to use something like order by created_at desc,
		//       so there is a workaround for now under this loop.

		switch operator {
		case "between":
			valStr, ok := val.(string)
			if !ok {
				return "", fmt.Errorf("value for param %s is not a string", param)
			}
			split := strings.Split(valStr, ",")
			if len(split) != 2 {
				return "", fmt.Errorf("Param %s should have 2 values separated by a comma", param)
			}

			firstParamName := param + "_1"
			secondParamName := param + "_2"

			if alias != "" {
				sqlStmt += fmt.Sprintf("AND %s.%s >= :%s\n", alias, param, firstParamName)
				sqlStmt += fmt.Sprintf("AND %s.%s <= :%s\n", alias, param, secondParamName)
			} else {
				sqlStmt += fmt.Sprintf("AND %s >= :%s\n", param, firstParamName)
				sqlStmt += fmt.Sprintf("AND %s <= :%s\n", param, secondParamName)
			}

			// Delete the old param from the params.
			delete(params, param)
			// Add the new params to params.
			params[firstParamName] = split[0]
			params[secondParamName] = split[1]
		default:
			return "", ErrOperatorNotSupported
		}
	}

	// Apply whatever operators are left.
	//
	// NOTE: For now this will cause an issue if any params
	//       contain these operators as names so for now prefix
	//       these ones with a double underscore (__).

	// If an order by clause is present, add it now.
	if val, ok := paramOperators["__order_by_asc"]; ok {
		fieldVal := val
		if alias != "" {
			fieldVal = fmt.Sprintf("%s.%s", alias, val)
		}
		sqlStmt += fmt.Sprintf("ORDER BY %s ASC\n", fieldVal)
	}
	if val, ok := paramOperators["__order_by_desc"]; ok {
		fieldVal := val
		if alias != "" {
			fieldVal = fmt.Sprintf("%s.%s", alias, val)
		}
		sqlStmt += fmt.Sprintf("ORDER BY %s DESC\n", fieldVal)
	}

	for po, val := range paramOperators {
		switch po {
		case "__limit":
			sqlStmt += fmt.Sprintf("LIMIT %s\n", val)
		case "__offset":
			sqlStmt += fmt.Sprintf("OFFSET %s\n", val)
		}
	}

	return sqlStmt, nil
}
