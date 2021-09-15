package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"market4/internal/cache"
	"market4/internal/model"
	"market4/internal/repository"
	"market4/internal/views"
	"net/http"
	"strconv"

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
}

func NewProduct(productRepo repository.Product, priceRepo repository.Price, stock cache.Cache) *Product {
	return &Product{productRepo: productRepo, priceRepo: priceRepo, stock: stock}
}
func (p *Product) AddProduct(writer http.ResponseWriter, request *http.Request) {
	var data *ProductDTO
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("addProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if IsEmpty(data.SKU) || IsEmpty(data.Name) || IsEmpty(data.Type) || IsEmpty(data.Description) {
		log.Println("field is empty")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var product = model.Product{
		SKU:         data.SKU,
		Name:        data.Name,
		Type:        data.Type,
		Description: data.Description,
	}

	addedProduct, err := p.productRepo.AddProduct(request.Context(), &product, data.Shop_ID, data.Category_ID)
	if err != nil {
		log.Println(fmt.Errorf("addProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	var editedPrice *model.Price
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
			log.Println(fmt.Errorf("addProduct: %w", err))
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	var productList = make([]*model.Product, 0)
	productList = append(productList, addedProduct)

	var priceList = make([]*model.Price, 0)
	priceList = append(priceList, editedPrice)

	result, err := views.ProductsListWithPrices(productList, priceList)
	if err != nil {
		log.Println(fmt.Errorf("ListAllProducts: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(result)
	if err != nil {
		log.Println(fmt.Errorf("addCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (p *Product) EditProduct(writer http.ResponseWriter, request *http.Request) {
	var data *ProductDTO
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(fmt.Errorf("editProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(data.SKU) {
		log.Println(fmt.Errorf("editProduct: SKU field is empty"))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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

	editedProduct, err := p.productRepo.EditProduct(request.Context(), &product, shopID, categoryID)
	if err != nil {
		log.Println(fmt.Errorf("editProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	var editedPrice *model.Price
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
			log.Println(fmt.Errorf("editProduct: %w", err))
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	var productList = make([]*model.Product, 0)
	productList = append(productList, editedProduct)

	var priceList = make([]*model.Price, 0)
	priceList = append(priceList, editedPrice)

	result, err := views.ProductsListWithPrices(productList, priceList)
	if err != nil {
		log.Println(fmt.Errorf("ListAllProducts: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(result)
	if err != nil {
		log.Println(fmt.Errorf("editProduct: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (p *Product) ListAllProducts(writer http.ResponseWriter, request *http.Request) {
	products, err := p.productRepo.ListAllProducts(request.Context())
	if err != nil {
		log.Println(fmt.Errorf("ListAllProducts: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	prices, err := p.priceRepo.ListAllPrices(request.Context())
	if err != nil {
		log.Println(fmt.Errorf("ListAllPrices: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	productsList, err := views.ProductsListWithPrices(products, prices)
	if err != nil {
		log.Println(fmt.Errorf("ListAllProducts: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(productsList)
	if err != nil {
		log.Println(fmt.Errorf("ListAllProducts: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
func (p *Product) SearchProductsByCategory(writer http.ResponseWriter, request *http.Request) {
	if productsList, _ := p.stock.FromCache(request.Context(), request.RequestURI); productsList != nil {
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write(productsList)
		if err != nil {
			log.Println(fmt.Errorf("SearchProductsByCategory: %w", err))
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	category, err := strconv.Atoi(chi.URLParam(request, "categoryID"))
	if err != nil {
		log.Println(fmt.Errorf("SearchProductsByCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	products, err := p.productRepo.SearchProductsByCategory(request.Context(), category)
	if err != nil {
		log.Println(fmt.Errorf("SearchProductsByCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if len(products) == 0 {
		return
	}

	productsList, err := views.ProductsList(products)
	if err != nil {
		log.Println(fmt.Errorf("SearchProductsByCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(productsList)
	if err != nil {
		log.Println(fmt.Errorf("SearchProductsByCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	_, err = writer.Write(body)
	if err != nil {
		log.Println(fmt.Errorf("SearchProductsByCategory: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	err = p.stock.ToCache(request.Context(), request.RequestURI, body)
	if err != nil {
		log.Println(fmt.Errorf("SearchProductsByCategory: %w", err))
	}
}
func (p *Product) SearchProductByName(writer http.ResponseWriter, request *http.Request) {
	if result, _ := p.stock.FromCache(request.Context(), request.RequestURI); result != nil {
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write(result)
		if err != nil {
			log.Println(fmt.Errorf("SearchProductByName: %w", err))
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	productName := chi.URLParam(request, "product_name")
	product, err := p.productRepo.SearchProductsByName(request.Context(), productName)
	if err != nil {
		log.Println(fmt.Errorf("SearchProductByName: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if product.ID == "" {
		return
	}
	price, err := p.priceRepo.SearchPriceByProductID(request.Context(), product.ID)
	if err != nil {
		log.Println(fmt.Errorf("SearchProductByName: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	var productList = make([]*model.Product, 0)
	productList = append(productList, product)

	var priceList = make([]*model.Price, 0)
	priceList = append(priceList, price)

	result, err := views.ProductsListWithPrices(productList, priceList)
	if err != nil {
		log.Println(fmt.Errorf("ListAllProducts: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(result)
	if err != nil {
		log.Println(fmt.Errorf("SearchPriceByProductID: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	_, err = writer.Write(body)
	if err != nil {
		log.Println(fmt.Errorf("SearchPriceByProductID: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	err = p.stock.ToCache(request.Context(), request.RequestURI, body)
	if err != nil {
		log.Println(fmt.Errorf("SearchPriceByProductID: %w", err))
	}
}
func (p *Product) SearchActiveProductsOfShop(writer http.ResponseWriter, request *http.Request) {
	shopID, err := strconv.Atoi(chi.URLParam(request, "shopID"))
	if err != nil {
		log.Println(fmt.Errorf("SearchActiveProductsOfShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	products, err := p.productRepo.SearchProductsByShop(request.Context(), shopID)
	if err != nil {
		log.Println(fmt.Errorf("SearchActiveProductsOfShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if len(products) == 0 {
		return
	}

	var prices = make([]*model.Price, 0)
	for _, product := range products {
		if product.IsActive {
			price, сerr := p.priceRepo.SearchPriceByProductID(request.Context(), product.ID)
			if сerr != nil {
				log.Println(fmt.Errorf("SearchActiveProductsOfShop: %w", err))
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			prices = append(prices, price)
		}
	}

	productsList, err := views.ProductsListWithPrices(products, prices)
	if err != nil {
		log.Println(fmt.Errorf("SearchActiveProductsOfShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(productsList)
	if err != nil {
		log.Println(fmt.Errorf("SearchActiveProductsOfShop: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
