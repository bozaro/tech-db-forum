// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/bozaro/tech-db-forum/generated/models"
)

// NewPostUpdateParams creates a new PostUpdateParams object
// with the default values initialized.
func NewPostUpdateParams() *PostUpdateParams {
	var ()
	return &PostUpdateParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostUpdateParamsWithTimeout creates a new PostUpdateParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostUpdateParamsWithTimeout(timeout time.Duration) *PostUpdateParams {
	var ()
	return &PostUpdateParams{

		timeout: timeout,
	}
}

// NewPostUpdateParamsWithContext creates a new PostUpdateParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostUpdateParamsWithContext(ctx context.Context) *PostUpdateParams {
	var ()
	return &PostUpdateParams{

		Context: ctx,
	}
}

// NewPostUpdateParamsWithHTTPClient creates a new PostUpdateParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostUpdateParamsWithHTTPClient(client *http.Client) *PostUpdateParams {
	var ()
	return &PostUpdateParams{
		HTTPClient: client,
	}
}

/*PostUpdateParams contains all the parameters to send to the API endpoint
for the post update operation typically these are written to a http.Request
*/
type PostUpdateParams struct {

	/*ID
	  Идентификатор сообщения.

	*/
	ID int64
	/*Post
	  Изменения сообщения.

	*/
	Post *models.PostUpdate

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post update params
func (o *PostUpdateParams) WithTimeout(timeout time.Duration) *PostUpdateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post update params
func (o *PostUpdateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post update params
func (o *PostUpdateParams) WithContext(ctx context.Context) *PostUpdateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post update params
func (o *PostUpdateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post update params
func (o *PostUpdateParams) WithHTTPClient(client *http.Client) *PostUpdateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post update params
func (o *PostUpdateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the post update params
func (o *PostUpdateParams) WithID(id int64) *PostUpdateParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the post update params
func (o *PostUpdateParams) SetID(id int64) {
	o.ID = id
}

// WithPost adds the post to the post update params
func (o *PostUpdateParams) WithPost(post *models.PostUpdate) *PostUpdateParams {
	o.SetPost(post)
	return o
}

// SetPost adds the post to the post update params
func (o *PostUpdateParams) SetPost(post *models.PostUpdate) {
	o.Post = post
}

// WriteToRequest writes these params to a swagger request
func (o *PostUpdateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt64(o.ID)); err != nil {
		return err
	}

	if o.Post == nil {
		o.Post = new(models.PostUpdate)
	}

	if err := r.SetBodyParam(o.Post); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
