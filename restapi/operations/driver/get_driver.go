// Code generated by go-swagger; DO NOT EDIT.

package driver

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetDriverHandlerFunc turns a function with the right signature into a get driver handler
type GetDriverHandlerFunc func(GetDriverParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetDriverHandlerFunc) Handle(params GetDriverParams) middleware.Responder {
	return fn(params)
}

// GetDriverHandler interface for that can handle valid get driver params
type GetDriverHandler interface {
	Handle(GetDriverParams) middleware.Responder
}

// NewGetDriver creates a new http.Handler for the get driver operation
func NewGetDriver(ctx *middleware.Context, handler GetDriverHandler) *GetDriver {
	return &GetDriver{Context: ctx, Handler: handler}
}

/*
	GetDriver swagger:route GET /driver driver getDriver

GetDriver get driver API
*/
type GetDriver struct {
	Context *middleware.Context
	Handler GetDriverHandler
}

func (o *GetDriver) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetDriverParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetDriverBadRequestBody get driver bad request body
//
// swagger:model GetDriverBadRequestBody
type GetDriverBadRequestBody struct {

	// error
	Error string `json:"error,omitempty"`
}

// Validate validates this get driver bad request body
func (o *GetDriverBadRequestBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get driver bad request body based on context it is used
func (o *GetDriverBadRequestBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetDriverBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDriverBadRequestBody) UnmarshalBinary(b []byte) error {
	var res GetDriverBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetDriverForbiddenBody get driver forbidden body
//
// swagger:model GetDriverForbiddenBody
type GetDriverForbiddenBody struct {

	// error
	Error string `json:"error,omitempty"`
}

// Validate validates this get driver forbidden body
func (o *GetDriverForbiddenBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get driver forbidden body based on context it is used
func (o *GetDriverForbiddenBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetDriverForbiddenBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDriverForbiddenBody) UnmarshalBinary(b []byte) error {
	var res GetDriverForbiddenBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetDriverInternalServerErrorBody get driver internal server error body
//
// swagger:model GetDriverInternalServerErrorBody
type GetDriverInternalServerErrorBody struct {

	// error
	Error string `json:"error,omitempty"`
}

// Validate validates this get driver internal server error body
func (o *GetDriverInternalServerErrorBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get driver internal server error body based on context it is used
func (o *GetDriverInternalServerErrorBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetDriverInternalServerErrorBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDriverInternalServerErrorBody) UnmarshalBinary(b []byte) error {
	var res GetDriverInternalServerErrorBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetDriverOKBody get driver o k body
//
// swagger:model GetDriverOKBody
type GetDriverOKBody struct {

	// email
	Email string `json:"email,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// phone number
	PhoneNumber string `json:"phone_number,omitempty"`

	// raiting
	Raiting float64 `json:"raiting,omitempty"`
}

// Validate validates this get driver o k body
func (o *GetDriverOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get driver o k body based on context it is used
func (o *GetDriverOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetDriverOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDriverOKBody) UnmarshalBinary(b []byte) error {
	var res GetDriverOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetDriverUnauthorizedBody get driver unauthorized body
//
// swagger:model GetDriverUnauthorizedBody
type GetDriverUnauthorizedBody struct {

	// error
	Error string `json:"error,omitempty"`
}

// Validate validates this get driver unauthorized body
func (o *GetDriverUnauthorizedBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get driver unauthorized body based on context it is used
func (o *GetDriverUnauthorizedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetDriverUnauthorizedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDriverUnauthorizedBody) UnmarshalBinary(b []byte) error {
	var res GetDriverUnauthorizedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
