package web

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// * Handle otomatis berdasarkan content type
type RequestDecoder struct {
	MaxMemory int64 // * Maximum besar form data
}

func NewRequestDecoder() *RequestDecoder {
	return &RequestDecoder{
		MaxMemory: 32 << 20,
	}
}

// * Decode otomatis berdasarkan content type
func (rd *RequestDecoder) Decode(r *http.Request, dst any) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return fmt.Errorf("content-type header is required")
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("invalid content-type: %w", err)
	}

	switch mediaType {
	case "application/json":
		return rd.decodeJSON(r, dst)
	case "application/x-www-form-urlencoded":
		return rd.decodeFormURLEncoded(r, dst)
	case "multipart/form-data":
		return rd.decodeMultipartForm(r, dst)
	default:
		return fmt.Errorf("unsupported content-type: %s", mediaType)
	}
}

func (rd *RequestDecoder) decodeJSON(r *http.Request, dst any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	defer r.Body.Close()

	if len(body) == 0 {
		return fmt.Errorf("request body is empty")
	}

	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

func (rd *RequestDecoder) decodeFormURLEncoded(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	return rd.mapFormValues(r.Form, dst)
}

func (rd *RequestDecoder) decodeMultipartForm(r *http.Request, dst any) error {
	if err := r.ParseMultipartForm(rd.MaxMemory); err != nil {
		return fmt.Errorf("failed to parse multipart form: %w", err)
	}

	if err := rd.mapFormValues(r.MultipartForm.Value, dst); err != nil {
		return err
	}

	return rd.mapFileHeaders(r.MultipartForm.File, dst)
}

func (rd *RequestDecoder) mapFormValues(values url.Values, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		fieldName := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			if parts := strings.Split(jsonTag, ","); parts[0] != "" && parts[0] != "-" {
				fieldName = parts[0]
			}
		}

		if formTag := fieldType.Tag.Get("form"); formTag != "" {
			if parts := strings.Split(formTag, ","); parts[0] != "" && parts[0] != "-" {
				fieldName = parts[0]
			}
		}

		formValue := values.Get(fieldName)
		if formValue == "" {
			continue
		}

		if err := rd.setFieldValue(field, formValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", fieldName, err)
		}
	}

	return nil
}

func (rd *RequestDecoder) mapFileHeaders(files map[string][]*multipart.FileHeader, dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		if field.Type() != reflect.TypeOf(&multipart.FileHeader{}) &&
			field.Type() != reflect.TypeOf([]*multipart.FileHeader{}) {
			continue
		}

		fieldName := fieldType.Name
		if formTag := fieldType.Tag.Get("form"); formTag != "" {
			if parts := strings.Split(formTag, ","); parts[0] != "" && parts[0] != "-" {
				fieldName = parts[0]
			}
		} else if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			if parts := strings.Split(jsonTag, ","); parts[0] != "" && parts[0] != "-" {
				fieldName = parts[0]
			}
		}

		fileHeaders, exists := files[fieldName]
		if !exists || len(fileHeaders) == 0 {
			continue
		}

		if field.Type() == reflect.TypeOf(&multipart.FileHeader{}) {
			field.Set(reflect.ValueOf(fileHeaders[0]))
		} else if field.Type() == reflect.TypeOf([]*multipart.FileHeader{}) {
			field.Set(reflect.ValueOf(fileHeaders))
		}
	}

	return nil
}

func (rd *RequestDecoder) setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, err := strconv.ParseInt(value, 10, 64); err != nil {
			return err
		} else {
			field.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintVal, err := strconv.ParseUint(value, 10, 64); err != nil {
			return err
		} else {
			field.SetUint(uintVal)
		}
	case reflect.Float32, reflect.Float64:
		if floatVal, err := strconv.ParseFloat(value, 64); err != nil {
			return err
		} else {
			field.SetFloat(floatVal)
		}
	case reflect.Bool:
		if boolVal, err := strconv.ParseBool(value); err != nil {
			return err
		} else {
			field.SetBool(boolVal)
		}
	case reflect.Ptr:
		if field.Type().Elem().Kind() == reflect.String {
			field.Set(reflect.ValueOf(&value))
		} else {
			newVal := reflect.New(field.Type().Elem())
			if err := rd.setFieldValue(newVal.Elem(), value); err != nil {
				return err
			}
			field.Set(newVal)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %s", ve[0].Message)
}

// * Validasi pakai playground validator
func Validate(s any) error {
	if err := validate.Struct(s); err != nil {
		var validationErrors ValidationErrors

		for _, err := range err.(validator.ValidationErrors) {
			validationError := ValidationError{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: fmt.Sprintf("%v", err.Value()),
			}

			validationError.Message = generateValidationMessage(err)
			validationErrors = append(validationErrors, validationError)
		}

		return validationErrors
	}
	return nil
}

func generateValidationMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	default:
		return fmt.Sprintf("%s failed validation for tag '%s'", field, tag)
	}
}

// * Helpers
func DecodeAndValidate(r *http.Request, dst any) error {
	decoder := NewRequestDecoder()

	if err := decoder.Decode(r, dst); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	if err := Validate(dst); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func DecodeJSON(r *http.Request, dst any) error {
	decoder := NewRequestDecoder()
	return decoder.decodeJSON(r, dst)
}

func DecodeForm(r *http.Request, dst any) error {
	decoder := NewRequestDecoder()
	return decoder.decodeFormURLEncoded(r, dst)
}

func ValidateAndRespond(w http.ResponseWriter, s any) bool {
	if err := Validate(s); err != nil {
		if validationErrors, ok := err.(ValidationErrors); ok {
			Error(w, http.StatusBadRequest, "Validation failed", validationErrors)
			return false
		}
		Error(w, http.StatusBadRequest, "Validation failed", err.Error())
		return false
	}
	return true
}

func DecodeValidateAndRespond(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := DecodeAndValidate(r, dst); err != nil {
		if validationErrors, ok := err.(ValidationErrors); ok {
			Error(w, http.StatusBadRequest, "Request validation failed", validationErrors)
			return false
		}
		Error(w, http.StatusBadRequest, "Request processing failed", err.Error())
		return false
	}
	return true
}
