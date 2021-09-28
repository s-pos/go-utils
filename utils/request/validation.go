package request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/s-pos/go-utils/logger"
)

const (
	REQUIRED = "required"
	JSON     = "json"
	FORM     = "form"
	QUERY    = "query"
	PARAM    = "param"
	FILE     = "file"
)

// BodyValidation is validation for request body
//
// How to use:
//   1. On struct, please add tag 'required' with type body you want use value.
//   	e.g 'required:"form"'. `form` means for Content-Type x-www-form-urlencoded / form-data
// 		e.g 'required:"json"'. `json` means for Content-Type application/json
// 		e.g 'required:"query"'. `query` means for QueryString url. `https://url?trxCode=abc`
// 		e.g 'required:"param"'. `param` means for Parameter on url section
// 		e.g 'required:"file"'. `file` means for Content-Type multipart/form-data
// 	 2. Please add too tag with value of 'required'.
// 		e.g if your 'required' is 'form' then you need add tag like this 'form:"field_name"'
// 		e.g if your 'required' is 'json' then you need add tag like this 'json:"field_name"'
// 		e.g if your 'required' is 'query' then you need add tag like this 'query:"field_name"'
// 		e.g if your 'required' is 'param' then you need add tag like this 'param:"field_name"'
// 		e.g if your 'required' is 'file' just using like 'form' required tag. e.g. 'form:"file"'
//   3. On your controller or Endpoint handler, please add this function for make
//      any request body validation
func BodyValidation(ctx context.Context, e echo.Context, req interface{}, reqType string) error {
	var (
		totalRequired int
	)

	if reqType == JSON {
		if err := json.NewDecoder(e.Request().Body).Decode(req); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				logger.FieldMandatory(ute.Field, fmt.Sprintf("%s harus berupa %s bukan %s", convertFieldName(ute.Field), convertTypeData(ute.Type.String()), convertTypeData(ute.Value))).Body(ctx)
				return err
			}
			return err
		}
	}

	to := reflect.TypeOf(req).Elem()
	vo := reflect.ValueOf(req).Elem()

	for i := 0; i < to.NumField(); i++ {
		var (
			field string
		)
		fieldTo := to.Field(i)
		valueRequired, reqField := fieldTo.Tag.Lookup(REQUIRED)

		if reqType == JSON {
			field = fieldTo.Tag.Get(JSON)
		}

		fieldVo := vo.Field(i)
		if reqType == FORM && reqField {
			field = fieldTo.Tag.Get(FORM)
			if !fieldVo.CanSet() {
				continue
			}

			if valueRequired != FILE {
				val := reflect.ValueOf(e.FormValue(field))
				if !val.IsZero() {
					v, err := convertKindValue(ctx, fieldTo, val, field, valueRequired)
					if err != nil {
						logger.Message("error convert kind %w", err).To(ctx)
						totalRequired++
						continue
					}

					val = v

					fieldVo.Set(val)
				}
			} else if valueRequired == FILE {
				v, err := e.FormFile(field)
				if err != nil {
					if !errors.Is(err, http.ErrMissingFile) {
						logger.FieldMandatory(strings.ToLower(fieldTo.Name), fmt.Sprintf("Field %s harus berbentuk file", convertFieldName(fieldTo.Name))).Body(ctx)
						totalRequired++
						continue
					}
				}
				val := reflect.ValueOf(v)
				if !val.IsZero() {
					fieldVo.Set(val)
				}
			}
		} else if reqType == FORM {
			field = fieldTo.Tag.Get(FORM)
			val := reflect.ValueOf(e.FormValue(field))
			if !val.IsZero() {
				v, err := convertKindValue(ctx, fieldTo, val, field, FORM)
				if err != nil {
					logger.Message("error convert kind %w", err).To(ctx)
					totalRequired++
					continue
				}

				val = v

				fieldVo.Set(val)
			}
		}

		tQuery, ok := fieldTo.Tag.Lookup(QUERY)
		if ok {
			field = tQuery

			val := reflect.ValueOf(e.QueryParam(field))
			if !val.IsZero() {
				val, err := convertKindValue(ctx, fieldTo, val, field, QUERY)
				if err != nil {
					logger.Message("error convert kind %w", err).To(ctx)
					totalRequired++
					continue
				}

				fieldVo.Set(val)
			}
		}

		tParam, paramOk := fieldTo.Tag.Lookup(PARAM)
		if paramOk {
			field = tParam

			val := reflect.ValueOf(e.Param(field))
			if !val.IsZero() {
				val, err := convertKindValue(ctx, fieldTo, val, field, PARAM)
				if err != nil {
					totalRequired++
					logger.Message("error convert kind %w", err).To(ctx)
					continue
				}

				fieldVo.Set(val)
			}
		}

		if reqField {
			if fieldVo.IsZero() {
				// log.Println(valueRequired, field, fieldVo.Kind(), fieldTo.Type.Kind())
				switch valueRequired {
				case JSON, FORM, FILE:
					logger.FieldMandatory(field, fmt.Sprintf("%s is required", convertFieldName(field))).Body(ctx)
				case QUERY:
					logger.FieldMandatory(tQuery, fmt.Sprintf("%s is required", convertFieldName(tQuery))).QueryString(ctx)
				case PARAM:
					logger.FieldMandatory(field, fmt.Sprintf("%s is required", convertFieldName(field))).Path(ctx)
				default:
				}

				totalRequired++
			}
		}
	}

	if totalRequired > 0 {
		err := fmt.Errorf("%d field required", totalRequired)

		return err
	}

	return nil
}

func convertKindValue(ctx context.Context, rStruct reflect.StructField, value reflect.Value, tag, dest string) (reflect.Value, error) {
	var (
		err  error
		logg logger.LogMessage
	)
	kind := rStruct.Type.Kind()
	fieldName := tag

	switch kind {
	case reflect.Float32, reflect.Float64:
		val, Err := strconv.ParseFloat(value.String(), 64)
		if Err != nil {
			err = Err
			logg = logger.FieldMandatory(fieldName, fmt.Sprintf("%s harus berupa %s bukan %s", convertFieldName(rStruct.Name), convertTypeData(kind.String()), convertTypeData(value.Type().String())))
			break
		}

		return reflect.ValueOf(val), nil
	case reflect.Int, reflect.Int32, reflect.Int64:
		val, Err := strconv.Atoi(value.String())
		if Err != nil {
			err = Err
			logg = logger.FieldMandatory(fieldName, fmt.Sprintf("%s harus berupa %s bukan %s", convertFieldName(rStruct.Name), convertTypeData(kind.String()), convertTypeData(value.Type().String())))
			break
		}

		return reflect.ValueOf(val), nil
	default:
		return value, nil
	}

	if err != nil {
		switch dest {
		case FORM, JSON:
			logg.Body(ctx)
			return reflect.Value{}, err
		case QUERY:
			logg.QueryString(ctx)
			return reflect.Value{}, err
		case PARAM:
			logg.Path(ctx)
			return reflect.Value{}, err
		default:
		}
	}
	return reflect.Value{}, err
}

func convertFieldName(field string) string {
	field = strings.ReplaceAll(field, "_", " ")

	return strings.Title(field)
}

func convertTypeData(data string) string {
	switch data {
	case reflect.Int.String(), reflect.Float64.String(), reflect.Float32.String(), reflect.Int8.String(), reflect.Int16.String(), reflect.Int64.String():
		return "number"
	case reflect.String.String():
		return "string"
	default:
		return data
	}
}
