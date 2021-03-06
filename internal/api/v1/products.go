package v1

import (
	"encoding/json"
	"market4/internal/cache"
	"market4/internal/model"
	"market4/internal/repository"
	"market4/internal/views"
	"net/http"
	"strconv"

	"github.com/unrolled/render"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

type ProductDTO struct {
	SKU         string    `json:"sku"`
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active,string,omitempty"`
	Shop_ID     int       `json:"shop_id,string,omitempty"`
	Category_ID int       `json:"category_id,string,omitempty"`
	Price       *PriceDTO `json:"price,omitempty"`
}
type Product struct {
	productRepo repository.Product
	priceRepo   repository.Price
	stock       cache.Cache
	lg          *zap.Logger
	renderer    *render.Render
}

func NewProduct(productRepo repository.Product,
	priceRepo repository.Price,
	stock cache.Cache,
	lg *zap.Logger,
	renderer *render.Render) *Product {
	return &Product{productRepo: productRepo,
		priceRepo: priceRepo,
		stock:     stock,
		lg:        lg,
		renderer:  renderer}
}
func (p *Product) AddProduct(writer http.ResponseWriter, request *http.Request) {
	var data *ProductDTO
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		p.lg.Error("addProduct", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	err = checkMandatoryFields(data.SKU, data.Name, data.Type, data.Description)
	if err != nil {
		p.lg.Error("addProduct: field id empty")
		err = p.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	var product = model.Product{
		SKU:         data.SKU,
		Name:        data.Name,
		Type:        data.Type,
		Description: data.Description,
	}

	addedProduct, err := p.productRepo.AddProduct(request.Context(), product, data.Shop_ID, data.Category_ID)
	if err != nil {
		p.lg.Error("addProduct", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	var editedPrice model.Price
	if data.Price != nil {
		var price = model.Price{
			SalePrice:     data.Price.SalePrice,
			FactoryPrice:  data.Price.FactoryPrice,
			DiscountPrice: data.Price.DiscountPrice,
			IsActive:      data.Price.IsActive,
			ProductID:     addedProduct.ID,
		}
		editedPrice, err = p.priceRepo.AddPrice(request.Context(), &price)
		if err != nil {
			p.lg.Error("addProduct", zap.Error(err))
			err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
			if err != nil {
				p.lg.Error("Auth", zap.Error(err))
			}
			return
		}
	}

	var productList = make([]model.Product, 0)
	productList = append(productList, addedProduct)

	var priceList = make([]model.Price, 0)
	priceList = append(priceList, editedPrice)

	result, err := views.MakeProductsListWithPrices(productList, priceList)
	if err != nil {
		p.lg.Error("addProduct", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(result)
	if err != nil {
		p.lg.Error("addProduct", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}

func (p *Product) EditProduct(writer http.ResponseWriter, request *http.Request) {
	var data *ProductDTO
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		p.lg.Error("EditProduct", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	err = checkMandatoryFields(data.SKU)
	if err != nil {
		p.lg.Error("EditProduct: SKU field is empty")
		err = p.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	shopID, categoryID := 0, 0

	var product = model.Product{
		SKU:         data.SKU,
		Name:        data.Name,
		Type:        data.Type,
		Description: data.Description,
		IsActive:    data.IsActive,
	}

	shopID = data.Shop_ID
	categoryID = data.Category_ID

	editedProduct, err := p.productRepo.EditProduct(request.Context(), product, shopID, categoryID)
	if err != nil {
		p.lg.Error("EditProduct", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	var editedPrice model.Price
	if data.Price != nil {
		var price = model.Price{
			SalePrice:     data.Price.SalePrice,
			FactoryPrice:  data.Price.FactoryPrice,
			DiscountPrice: data.Price.DiscountPrice,
			IsActive:      data.Price.IsActive,
			ProductID:     editedProduct.ID,
		}
		editedPrice, err = p.priceRepo.EditPriceByProductID(request.Context(), &price)
		if err != nil {
			p.lg.Error("EditProduct", zap.Error(err))
			err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
			if err != nil {
				p.lg.Error("Auth", zap.Error(err))
			}
			return
		}
	}

	var productList = make([]model.Product, 0)
	productList = append(productList, editedProduct)

	var priceList = make([]model.Price, 0)
	priceList = append(priceList, editedPrice)

	result, err := views.MakeProductsListWithPrices(productList, priceList)
	if err != nil {
		p.lg.Error("EditProduct", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(result)
	if err != nil {
		p.lg.Error("EditProduct", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}

func (p *Product) ListAllProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := p.productRepo.ListAllProducts(request.Context())
	if err != nil {
		p.lg.Error("ListAllProducts", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	prices, err := p.priceRepo.ListAllPrices(request.Context())
	if err != nil {
		p.lg.Error("ListAllProducts", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	productsList, err := views.MakeProductsListWithPrices(products, prices)
	if err != nil {
		p.lg.Error("ListAllProducts", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(productsList)
	if err != nil {
		p.lg.Error("ListAllProducts", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}
func (p *Product) SearchProductsByCategory(writer http.ResponseWriter, request *http.Request) {
	productsList, err := p.stock.FromCache(request.Context(), request.RequestURI)
	if err != nil {
		p.lg.Error("SearchProductsByCategory", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	if productsList != nil {
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(productsList)
		if err != nil {
			p.lg.Error("SearchProductsByCategory", zap.Error(err))
			err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
			if err != nil {
				p.lg.Error("Auth", zap.Error(err))
			}
			return
		}
		return
	}

	category, err := strconv.Atoi(chi.URLParam(request, "categoryID"))
	if err != nil {
		p.lg.Error("SearchProductsByCategory", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	products, err := p.productRepo.SearchProductsByCategory(request.Context(), category)
	if err != nil {
		p.lg.Error("SearchProductsByCategory", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	if len(products) == 0 {
		return
	}

	pList, err := views.MakeProductsList(products)
	if err != nil {
		p.lg.Error("SearchProductsByCategory", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(pList)
	if err != nil {
		p.lg.Error("SearchProductsByCategory", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	_, err = writer.Write(body)
	if err != nil {
		p.lg.Error("SearchProductsByCategory", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	err = p.stock.ToCache(request.Context(), request.RequestURI, body)
	if err != nil {
		p.lg.Error("SearchProductsByCategory", zap.Error(err))
	}
}
func (p *Product) SearchProductByName(writer http.ResponseWriter, request *http.Request) {
	if result, _ := p.stock.FromCache(request.Context(), request.RequestURI); result != nil {
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write(result)
		if err != nil {
			p.lg.Error("SearchProductByName", zap.Error(err))
			err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
			if err != nil {
				p.lg.Error("Auth", zap.Error(err))
			}
			return
		}
		return
	}

	productName := chi.URLParam(request, "product_name")
	product, err := p.productRepo.SearchProductsByName(request.Context(), productName)
	if err != nil {
		p.lg.Error("SearchProductByName", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	if product.ID == "" {
		return
	}
	price, err := p.priceRepo.SearchPriceByProductID(request.Context(), product.ID)
	if err != nil {
		p.lg.Error("SearchProductByName", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	var productList = make([]model.Product, 0)
	productList = append(productList, product)

	var priceList = make([]model.Price, 0)
	priceList = append(priceList, price)

	result, err := views.MakeProductsListWithPrices(productList, priceList)
	if err != nil {
		p.lg.Error("SearchProductByName", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(result)
	if err != nil {
		p.lg.Error("SearchProductByName", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	_, err = writer.Write(body)
	if err != nil {
		p.lg.Error("SearchProductByName", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	err = p.stock.ToCache(request.Context(), request.RequestURI, body)
	if err != nil {
		p.lg.Error("SearchProductByName", zap.Error(err))
	}
}
func (p *Product) SearchActiveProductsOfShop(writer http.ResponseWriter, request *http.Request) {
	shopID, err := strconv.Atoi(chi.URLParam(request, "shopID"))
	if err != nil {
		p.lg.Error("SearchActiveProductsOfShop", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusBadRequest, map[string]string{"Error": "BadRequest"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}
	products, err := p.productRepo.SearchProductsByShop(request.Context(), shopID)
	if err != nil {
		p.lg.Error("SearchActiveProductsOfShop", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	if len(products) == 0 {
		return
	}

	var prices = make([]model.Price, 0)
	for _, product := range products {
		if product.IsActive {
			price, cerr := p.priceRepo.SearchPriceByProductID(request.Context(), product.ID)
			if cerr != nil {
				p.lg.Error("SearchActiveProductsOfShop", zap.Error(err))
				err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
				if err != nil {
					p.lg.Error("Auth", zap.Error(err))
				}
				return
			}
			prices = append(prices, price)
		}
	}

	productsList, err := views.MakeProductsListWithPrices(products, prices)
	if err != nil {
		p.lg.Error("SearchActiveProductsOfShop", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(productsList)
	if err != nil {
		p.lg.Error("SearchActiveProductsOfShop", zap.Error(err))
		err = p.renderer.JSON(writer, http.StatusInternalServerError, map[string]string{"Error": "InternalServerError"})
		if err != nil {
			p.lg.Error("Auth", zap.Error(err))
		}
		return
	}
}
