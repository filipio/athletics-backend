package utils

import (
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
)

func TestBuildErrorPath(t *testing.T) {
	testData := []struct {
		name           string
		errorNamespace string
		expectedPath   []string
	}{
		{
			name:           "with two dots",
			errorNamespace: "utils.http.buildErrorPath",
			expectedPath:   []string{"http", "build_error_path"},
		},
		{
			name:           "with one dot",
			errorNamespace: "utils.buildErrorPath",
			expectedPath:   []string{"build_error_path"},
		},
		{
			name:           "without dots",
			errorNamespace: "buildErrorPath",
			expectedPath:   []string{},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			t.Parallel()
			actualPath := buildErrorPath(data.errorNamespace)

			if !cmp.Equal(actualPath, data.expectedPath) {
				t.Errorf("Error path does not match expected path: %s", cmp.Diff(actualPath, data.expectedPath))
			}

		})
	}
}

func TestValidationErrorResponse(t *testing.T) {
	err1 := FieldErrorMock{
		_Tag:       "required",
		_Namespace: "utils.http.buildErrorPath",
		_Field:     "buildErrorPath",
		_Param:     "",
		_Kind:      reflect.String,
		_Type:      reflect.TypeOf(""),
		_Error:     "This field is required",
		_Value:     "",
	}

	err2 := FieldErrorMock{
		_Tag:       "email",
		_Namespace: "some.veryDeep.error.path",
		_Field:     "path",
		_Param:     "",
		_Kind:      reflect.String,
		_Type:      reflect.TypeOf(""),
		_Error:     "some error",
		_Value:     "",
	}

	err := validator.ValidationErrors{
		err1,
		err2,
	}

	expectedResponse := &ErrorsResponse{
		Errors: AnyMap{
			"http": AnyMap{
				"build_error_path": ValidationResponseItem{
					Type:    "required",
					Message: "This field is required",
				},
			},
			"very_deep": AnyMap{
				"error": AnyMap{
					"path": ValidationResponseItem{
						Type:    "email",
						Message: "This field must be an email",
					},
				},
			},
		},
	}

	response := &ErrorsResponse{ValidationErrorMap(err)}

	if !cmp.Equal(response, expectedResponse) {
		t.Errorf("Response does not match expected response: %s", cmp.Diff(response, expectedResponse))
	}
}

type FieldErrorMock struct {
	_Tag       string
	_Namespace string
	_Field     string
	_Param     string
	_Kind      reflect.Kind
	_Type      reflect.Type
	_Error     string
	_Value     interface{}
}

func (f FieldErrorMock) Tag() string {
	return f._Tag
}

func (f FieldErrorMock) ActualTag() string {
	return f._Tag
}

func (f FieldErrorMock) Namespace() string {
	return f._Namespace
}

func (f FieldErrorMock) StructNamespace() string {
	return f._Namespace
}

func (f FieldErrorMock) Field() string {
	return f._Field
}

func (f FieldErrorMock) StructField() string {
	return f._Field
}

func (f FieldErrorMock) Param() string {
	return f._Param
}

func (f FieldErrorMock) Kind() reflect.Kind {
	return reflect.String
}

func (f FieldErrorMock) Type() reflect.Type {
	return f._Type
}

func (f FieldErrorMock) Translate(ut ut.Translator) string {
	return f._Error
}

func (f FieldErrorMock) Error() string {
	return f._Error
}

func (f FieldErrorMock) Value() interface{} {
	return f._Value
}
