package containers

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
)

// NewGetContainerLogsParams creates a new GetContainerLogsParams object
// with the default values initialized.
func NewGetContainerLogsParams() *GetContainerLogsParams {
	var (
		followDefault    = bool(false)
		timestampDefault = bool(false)
	)
	return &GetContainerLogsParams{
		Follow:    &followDefault,
		Timestamp: &timestampDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewGetContainerLogsParamsWithTimeout creates a new GetContainerLogsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetContainerLogsParamsWithTimeout(timeout time.Duration) *GetContainerLogsParams {
	var (
		followDefault    = bool(false)
		timestampDefault = bool(false)
	)
	return &GetContainerLogsParams{
		Follow:    &followDefault,
		Timestamp: &timestampDefault,

		timeout: timeout,
	}
}

// NewGetContainerLogsParamsWithContext creates a new GetContainerLogsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetContainerLogsParamsWithContext(ctx context.Context) *GetContainerLogsParams {
	var (
		followDefault    = bool(false)
		timestampDefault = bool(false)
	)
	return &GetContainerLogsParams{
		Follow:    &followDefault,
		Timestamp: &timestampDefault,

		Context: ctx,
	}
}

// NewGetContainerLogsParamsWithHTTPClient creates a new GetContainerLogsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetContainerLogsParamsWithHTTPClient(client *http.Client) *GetContainerLogsParams {
	var (
		followDefault    = bool(false)
		timestampDefault = bool(false)
	)
	return &GetContainerLogsParams{
		Follow:     &followDefault,
		Timestamp:  &timestampDefault,
		HTTPClient: client,
	}
}

/*GetContainerLogsParams contains all the parameters to send to the API endpoint
for the get container logs operation typically these are written to a http.Request
*/
type GetContainerLogsParams struct {

	/*OpID*/
	OpID *string
	/*Deadline*/
	Deadline *int64
	/*Follow*/
	Follow *bool
	/*ID*/
	ID string
	/*Since*/
	Since *int64
	/*Taillines*/
	Taillines *int64
	/*Timestamp*/
	Timestamp *bool

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get container logs params
func (o *GetContainerLogsParams) WithTimeout(timeout time.Duration) *GetContainerLogsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get container logs params
func (o *GetContainerLogsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get container logs params
func (o *GetContainerLogsParams) WithContext(ctx context.Context) *GetContainerLogsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get container logs params
func (o *GetContainerLogsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get container logs params
func (o *GetContainerLogsParams) WithHTTPClient(client *http.Client) *GetContainerLogsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get container logs params
func (o *GetContainerLogsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithOpID adds the opID to the get container logs params
func (o *GetContainerLogsParams) WithOpID(opID *string) *GetContainerLogsParams {
	o.SetOpID(opID)
	return o
}

// SetOpID adds the opId to the get container logs params
func (o *GetContainerLogsParams) SetOpID(opID *string) {
	o.OpID = opID
}

// WithDeadline adds the deadline to the get container logs params
func (o *GetContainerLogsParams) WithDeadline(deadline *int64) *GetContainerLogsParams {
	o.SetDeadline(deadline)
	return o
}

// SetDeadline adds the deadline to the get container logs params
func (o *GetContainerLogsParams) SetDeadline(deadline *int64) {
	o.Deadline = deadline
}

// WithFollow adds the follow to the get container logs params
func (o *GetContainerLogsParams) WithFollow(follow *bool) *GetContainerLogsParams {
	o.SetFollow(follow)
	return o
}

// SetFollow adds the follow to the get container logs params
func (o *GetContainerLogsParams) SetFollow(follow *bool) {
	o.Follow = follow
}

// WithID adds the id to the get container logs params
func (o *GetContainerLogsParams) WithID(id string) *GetContainerLogsParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get container logs params
func (o *GetContainerLogsParams) SetID(id string) {
	o.ID = id
}

// WithSince adds the since to the get container logs params
func (o *GetContainerLogsParams) WithSince(since *int64) *GetContainerLogsParams {
	o.SetSince(since)
	return o
}

// SetSince adds the since to the get container logs params
func (o *GetContainerLogsParams) SetSince(since *int64) {
	o.Since = since
}

// WithTaillines adds the taillines to the get container logs params
func (o *GetContainerLogsParams) WithTaillines(taillines *int64) *GetContainerLogsParams {
	o.SetTaillines(taillines)
	return o
}

// SetTaillines adds the taillines to the get container logs params
func (o *GetContainerLogsParams) SetTaillines(taillines *int64) {
	o.Taillines = taillines
}

// WithTimestamp adds the timestamp to the get container logs params
func (o *GetContainerLogsParams) WithTimestamp(timestamp *bool) *GetContainerLogsParams {
	o.SetTimestamp(timestamp)
	return o
}

// SetTimestamp adds the timestamp to the get container logs params
func (o *GetContainerLogsParams) SetTimestamp(timestamp *bool) {
	o.Timestamp = timestamp
}

// WriteToRequest writes these params to a swagger request
func (o *GetContainerLogsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	r.SetTimeout(o.timeout)
	var res []error

	if o.OpID != nil {

		// header param Op-ID
		if err := r.SetHeaderParam("Op-ID", *o.OpID); err != nil {
			return err
		}

	}

	if o.Deadline != nil {

		// query param deadline
		var qrDeadline int64
		if o.Deadline != nil {
			qrDeadline = *o.Deadline
		}
		qDeadline := swag.FormatInt64(qrDeadline)
		if qDeadline != "" {
			if err := r.SetQueryParam("deadline", qDeadline); err != nil {
				return err
			}
		}

	}

	if o.Follow != nil {

		// query param follow
		var qrFollow bool
		if o.Follow != nil {
			qrFollow = *o.Follow
		}
		qFollow := swag.FormatBool(qrFollow)
		if qFollow != "" {
			if err := r.SetQueryParam("follow", qFollow); err != nil {
				return err
			}
		}

	}

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if o.Since != nil {

		// query param since
		var qrSince int64
		if o.Since != nil {
			qrSince = *o.Since
		}
		qSince := swag.FormatInt64(qrSince)
		if qSince != "" {
			if err := r.SetQueryParam("since", qSince); err != nil {
				return err
			}
		}

	}

	if o.Taillines != nil {

		// query param taillines
		var qrTaillines int64
		if o.Taillines != nil {
			qrTaillines = *o.Taillines
		}
		qTaillines := swag.FormatInt64(qrTaillines)
		if qTaillines != "" {
			if err := r.SetQueryParam("taillines", qTaillines); err != nil {
				return err
			}
		}

	}

	if o.Timestamp != nil {

		// query param timestamp
		var qrTimestamp bool
		if o.Timestamp != nil {
			qrTimestamp = *o.Timestamp
		}
		qTimestamp := swag.FormatBool(qrTimestamp)
		if qTimestamp != "" {
			if err := r.SetQueryParam("timestamp", qTimestamp); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
