// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/bozaro/tech-db-forum/generated/models"
)

// PostGetOneReader is a Reader for the PostGetOne structure.
type PostGetOneReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostGetOneReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPostGetOneOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 404:
		result := NewPostGetOneNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPostGetOneOK creates a PostGetOneOK with default headers values
func NewPostGetOneOK() *PostGetOneOK {
	return &PostGetOneOK{}
}

/*PostGetOneOK handles this case with default header values.

Информация о ветке обсуждения.

*/
type PostGetOneOK struct {
	Payload *models.PostFull
}

func (o *PostGetOneOK) Error() string {
	return fmt.Sprintf("[GET /post/{id}/details][%d] postGetOneOK  %+v", 200, o.Payload)
}

func (o *PostGetOneOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PostFull)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostGetOneNotFound creates a PostGetOneNotFound with default headers values
func NewPostGetOneNotFound() *PostGetOneNotFound {
	return &PostGetOneNotFound{}
}

/*PostGetOneNotFound handles this case with default header values.

Ветка обсуждения отсутсвует в форуме.

*/
type PostGetOneNotFound struct {
}

func (o *PostGetOneNotFound) Error() string {
	return fmt.Sprintf("[GET /post/{id}/details][%d] postGetOneNotFound ", 404)
}

func (o *PostGetOneNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
