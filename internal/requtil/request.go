package requtil

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetParamsByName(r *http.Request, param string) (int, error) {
	valueStr := r.URL.Query().Get(param)
	value, err := strconv.Atoi(valueStr)

	if err != nil {
		str := fmt.Sprintf("параметр  \"%s\" должен быть числом", param)
		return 0, errors.New(str)
	}

	return value, nil
}
