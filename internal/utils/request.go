package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

// GetIntParam получает int параметр из URL, конвертируя строку в цисло
// но в случае ошибок, отправляет ответ в виде JSON
// Например,
//
// money := 0
// money, ok := utils.GetIntParam(r, "money")(w)
//
//	if !ok {
//	  return
//	}
//
// или
//
//	getMoneyParam := utils.GetIntParam(r, "money")
//	money, ok := getMoneyParam(w)
//
//	if !ok {
//	 return
//	}
//
// Функция найдет в URL параметр "money", но
// * если его не будет, то возвращаемая функция
// вызовет внутри себя MarshalResponse с кодом ошибки и сообщением,
// что параметр  "money" должен быть числом,
// * такая же ошибка вернется, если "money" будет иметь любое
// не числовое значение (строка, символ, закарючка)
func GetIntParam(r *http.Request, param string) func(w http.ResponseWriter) (int, bool) {
	urlParam := r.URL.Query().Get(param)
	buff, err := strconv.Atoi(urlParam)

	return func(w http.ResponseWriter) (int, bool) {
		if err != nil {
			customErrText := fmt.Sprintf("параметр  \"%s\" должен быть числом", param)

			message := JSONMessage{
				Code:    REQUEST_ERROR_CODE,
				Message: customErrText,
			}
			MarshalResponse(w, http.StatusBadRequest, &message)
			return 0, false
		}
		return buff, true
	}
}
