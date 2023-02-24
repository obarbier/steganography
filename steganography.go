package main

import "io"

type StegWriter interface {
	Write(writer io.Writer) (int, error)
}
