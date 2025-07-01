package books

import "errors"

var (
	ErrInvalidDiscount = errors.New("descuento inválido")
	ErrBookNotFound    = errors.New("libro no encontrado")
	ErrBookCreation    = errors.New("error al crear el libro")
)
