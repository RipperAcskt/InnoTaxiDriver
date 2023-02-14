// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PostDriverSingInHandlerFunc turns a function with the right signature into a post driver sing in handler
type PostDriverSingInHandlerFunc func(PostDriverSingInParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostDriverSingInHandlerFunc) Handle(params PostDriverSingInParams) middleware.Responder {
	return fn(params)
}

// PostDriverSingInHandler interface for that can handle valid post driver sing in params
type PostDriverSingInHandler interface {
	Handle(PostDriverSingInParams) middleware.Responder
}

// NewPostDriverSingIn creates a new http.Handler for the post driver sing in operation
func NewPostDriverSingIn(ctx *middleware.Context, handler PostDriverSingInHandler) *PostDriverSingIn {
	return &PostDriverSingIn{Context: ctx, Handler: handler}
}

/*
	PostDriverSingIn swagger:route POST /driver/sing-in postDriverSingIn

PostDriverSingIn post driver sing in API
*/
type PostDriverSingIn struct {
	Context *middleware.Context
	Handler PostDriverSingInHandler
}

func (o *PostDriverSingIn) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostDriverSingInParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PostDriverSingInBadRequestBody post driver sing in bad request body
//
// swagger:model PostDriverSingInBadRequestBody
type PostDriverSingInBadRequestBody struct {

	// error
	Error string `json:"error,omitempty"`
}

// Validate validates this post driver sing in bad request body
func (o *PostDriverSingInBadRequestBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post driver sing in bad request body based on context it is used
func (o *PostDriverSingInBadRequestBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostDriverSingInBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDriverSingInBadRequestBody) UnmarshalBinary(b []byte) error {
	var res PostDriverSingInBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostDriverSingInBody post driver sing in body
//
// swagger:model PostDriverSingInBody
type PostDriverSingInBody struct {

	// password
	Password string `json:"password,omitempty"`

	// phone number
	PhoneNumber string `json:"phone_number,omitempty"`
}

// Validate validates this post driver sing in body
func (o *PostDriverSingInBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post driver sing in body based on context it is used
func (o *PostDriverSingInBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostDriverSingInBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDriverSingInBody) UnmarshalBinary(b []byte) error {
	var res PostDriverSingInBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostDriverSingInForbiddenBody post driver sing in forbidden body
//
// swagger:model PostDriverSingInForbiddenBody
type PostDriverSingInForbiddenBody struct {

	// error
	Error string `json:"error,omitempty"`
}

// Validate validates this post driver sing in forbidden body
func (o *PostDriverSingInForbiddenBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post driver sing in forbidden body based on context it is used
func (o *PostDriverSingInForbiddenBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostDriverSingInForbiddenBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDriverSingInForbiddenBody) UnmarshalBinary(b []byte) error {
	var res PostDriverSingInForbiddenBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostDriverSingInInternalServerErrorBody post driver sing in internal server error body
//
// swagger:model PostDriverSingInInternalServerErrorBody
type PostDriverSingInInternalServerErrorBody struct {

	// error
	Error string `json:"error,omitempty"`
}

// Validate validates this post driver sing in internal server error body
func (o *PostDriverSingInInternalServerErrorBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post driver sing in internal server error body based on context it is used
func (o *PostDriverSingInInternalServerErrorBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostDriverSingInInternalServerErrorBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDriverSingInInternalServerErrorBody) UnmarshalBinary(b []byte) error {
	var res PostDriverSingInInternalServerErrorBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// PostDriverSingInOKBody post driver sing in o k body
//
// swagger:model PostDriverSingInOKBody
type PostDriverSingInOKBody struct {

	// access token
	AccessToken int64 `json:"access_token,omitempty"`

	// refresh token
	RefreshToken string `json:"refresh_token,omitempty"`
}

// Validate validates this post driver sing in o k body
func (o *PostDriverSingInOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post driver sing in o k body based on context it is used
func (o *PostDriverSingInOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostDriverSingInOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostDriverSingInOKBody) UnmarshalBinary(b []byte) error {
	var res PostDriverSingInOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
