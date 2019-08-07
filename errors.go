package main

import "fmt"

type FileNotFound struct {
	Filename string
}

func (f FileNotFound) Error() string {
	return fmt.Sprintf("File not found: %s", f.Filename)
}

type BadRequest struct {
	Message string
}

func (b BadRequest) Error() string {
	return b.Message
}

func NewBadRequest(message string) BadRequest {
	return BadRequest{message}
}
