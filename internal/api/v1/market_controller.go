package v1

//type marketController struct {
//	repo repository.MarketRepository
//}

//type MarketController interface {
//	ListAllShops(writer http.ResponseWriter, request *http.Request)
//	AddShop(writer http.ResponseWriter, request *http.Request)
//	EditShop(writer http.ResponseWriter, request *http.Request)
//
//	ListAllCategories(writer http.ResponseWriter, request *http.Request)
//	AddCategory(writer http.ResponseWriter, request *http.Request)
//	EditCategory(writer http.ResponseWriter, request *http.Request)
//	//
//	//ListAllProducts(writer http.ResponseWriter, request *http.Request)
//	AddProduct(writer http.ResponseWriter, request *http.Request)
//	//EditProduct(writer http.ResponseWriter, request *http.Request)
//}

//func NewMarketController(repo repository.MarketRepository) MarketController {
//	return &marketController{repo: repo}
//}
//
//func (m *marketController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
//	m.ServeHTTP(writer, request)
//}

func IsEmpty(field string) bool {
	return field == ""
}
