// Code generated by go-swagger; DO NOT EDIT.

package driver

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PutDriverOKCode is the HTTP code returned for type PutDriverOK
const PutDriverOKCode int = 200

/*
PutDriverOK OK

swagger:response putDriverOK
*/
type PutDriverOK struct {
}

// NewPutDriverOK creates PutDriverOK with default headers values
func NewPutDriverOK() *PutDriverOK {

	return &PutDriverOK{}
}

// WriteResponse to the client
func (o *PutDriverOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// PutDriverBadRequestCode is the HTTP code returned for type PutDriverBadRequest
const PutDriverBadRequestCode int = 400

/*
PutDriverBadRequest Driver does not exist

swagger:response putDriverBadRequest
*/
type PutDriverBadRequest struct {

	/*
	  In: Body
	*/
	Payload *PutDriverBadRequestBody `json:"body,omitempty"`
}

// NewPutDriverBadRequest creates PutDriverBadRequest with default headers values
func NewPutDriverBadRequest() *PutDriverBadRequest {

	return &PutDriverBadRequest{}
}

// WithPayload adds the payload to the put driver bad request response
func (o *PutDriverBadRequest) WithPayload(payload *PutDriverBadRequestBody) *PutDriverBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put driver bad request response
func (o *PutDriverBadRequest) SetPayload(payload *PutDriverBadRequestBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutDriverBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutDriverUnauthorizedCode is the HTTP code returned for type PutDriverUnauthorized
const PutDriverUnauthorizedCode int = 401

/*
PutDriverUnauthorized Token expired

swagger:response putDriverUnauthorized
*/
type PutDriverUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *PutDriverUnauthorizedBody `json:"body,omitempty"`
}

// NewPutDriverUnauthorized creates PutDriverUnauthorized with default headers values
func NewPutDriverUnauthorized() *PutDriverUnauthorized {

	return &PutDriverUnauthorized{}
}

// WithPayload adds the payload to the put driver unauthorized response
func (o *PutDriverUnauthorized) WithPayload(payload *PutDriverUnauthorizedBody) *PutDriverUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put driver unauthorized response
func (o *PutDriverUnauthorized) SetPayload(payload *PutDriverUnauthorizedBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutDriverUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutDriverForbiddenCode is the HTTP code returned for type PutDriverForbidden
const PutDriverForbiddenCode int = 403

/*
PutDriverForbidden Incorrect token

swagger:response putDriverForbidden
*/
type PutDriverForbidden struct {

	/*
	  In: Body
	*/
	Payload *PutDriverForbiddenBody `json:"body,omitempty"`
}

// NewPutDriverForbidden creates PutDriverForbidden with default headers values
func NewPutDriverForbidden() *PutDriverForbidden {

	return &PutDriverForbidden{}
}

// WithPayload adds the payload to the put driver forbidden response
func (o *PutDriverForbidden) WithPayload(payload *PutDriverForbiddenBody) *PutDriverForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put driver forbidden response
func (o *PutDriverForbidden) SetPayload(payload *PutDriverForbiddenBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutDriverForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutDriverInternalServerErrorCode is the HTTP code returned for type PutDriverInternalServerError
const PutDriverInternalServerErrorCode int = 500

/*
PutDriverInternalServerError Unexpected server error

swagger:response putDriverInternalServerError
*/
type PutDriverInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *PutDriverInternalServerErrorBody `json:"body,omitempty"`
}

// NewPutDriverInternalServerError creates PutDriverInternalServerError with default headers values
func NewPutDriverInternalServerError() *PutDriverInternalServerError {

	return &PutDriverInternalServerError{}
}

// WithPayload adds the payload to the put driver internal server error response
func (o *PutDriverInternalServerError) WithPayload(payload *PutDriverInternalServerErrorBody) *PutDriverInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put driver internal server error response
func (o *PutDriverInternalServerError) SetPayload(payload *PutDriverInternalServerErrorBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutDriverInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
