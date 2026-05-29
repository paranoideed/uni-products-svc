package responses

import (
	"net/http"

	"github.com/netbill/restkit/pagi"
	"github.com/paranoideed/uni-products-svc/internal/domain"
	"github.com/paranoideed/uni-products-svc/pkg/resources"
)

func Product(m domain.Product) resources.Product {
	return resources.Product{
		Data: resources.ProductData{
			Id:   m.ID,
			Type: "product",
			Attributes: resources.ProductDataAttributes{
				Name:      m.Name,
				Price:     m.Price,
				CreatedAt: m.CreatedAt,
			},
		},
	}
}

func ProductsCollection(r *http.Request, page pagi.Page[[]domain.Product]) resources.ProductsCollection {
	data := make([]resources.ProductData, 0, len(page.Data))

	for _, p := range page.Data {
		data = append(data, Product(p).Data)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	return resources.ProductsCollection{
		Data: data,
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}
}
