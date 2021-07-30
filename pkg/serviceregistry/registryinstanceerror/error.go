package registryinstanceerror

import (
	"errors"
	"fmt"

	"github.com/redhat-developer/app-services-sdk-go/registryinstance/apiv1internal/client"
)

type Error struct {
	Err error
}

func (e *Error) Error() string {
	return fmt.Sprint(e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// GetAPIError gets a strongly typed error from an error
func GetAPIError(err error) (e registryinstanceclient.Error, ok bool) {
	var apiError registryinstanceclient.GenericOpenAPIError

	if ok = errors.As(err, &apiError); ok {
		errModel := apiError.Model()

		e, ok = errModel.(registryinstanceclient.Error)
	}

	return e, ok
}

// Error code contains message that can be returned to the user
// TODO - this method doesn't seem to work and map common errors
func TransformError(err error) error {
	mappedErr, ok := GetAPIError(err)
	if !ok {
		return err
	}

	return errors.New("Error: " + mappedErr.GetMessage() + " Detail " + mappedErr.GetDetail())
}
