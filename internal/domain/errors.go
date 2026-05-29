package domain

import (
	"github.com/netbill/ape"
)

var ErrorProductNotFound = ape.DeclareError("PRODUCT_NOT_FOUND")
var ErrorNotValidInput = ape.DeclareError("NOT_VALID_INPUT")
