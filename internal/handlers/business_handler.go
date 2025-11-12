package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/services"
	"kilocode.dev/whatsapp-bot/pkg/utils"
)

type BusinessHandler struct {
	businessService *services.BusinessService
}

func NewBusinessHandler(businessService *services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
	}
}

// CreateProduct creates a product
func (h *BusinessHandler) CreateProduct(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,min=0"`
		Category    string  `json:"category" binding:"required"`
		ImageURL    string  `json:"image_url"`
		Stock       int     `json:"stock" binding:"min=0"`
		UserID      string  `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	product, err := h.businessService.CreateProduct(req.Name, req.Description, req.Price, req.Category, req.ImageURL, req.Stock, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, product)
}

// GetProducts gets products
func (h *BusinessHandler) GetProducts(c *gin.Context) {
	category := c.Query("category")
	search := c.Query("search")
	page := utils.Paginate(0, 1, 10) // Default pagination

	products, total, err := h.businessService.GetProducts(category, search, 1, 10)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, 1, 10)
	utils.ResponsePaginated(c, products, pagination)
}

// GetProduct gets a product
func (h *BusinessHandler) GetProduct(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	productUUID, err := uuid.Parse(productID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.businessService.GetProduct(productUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "Product not found")
		return
	}

	utils.ResponseSuccess(c, product)
}

// UpdateProduct updates a product
func (h *BusinessHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"min=0"`
		Category    string  `json:"category"`
		ImageURL    string  `json:"image_url"`
		Stock       int     `json:"stock" binding:"min=0"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	productUUID, err := uuid.Parse(productID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.businessService.UpdateProduct(productUUID, req.Name, req.Description, req.Price, req.Category, req.ImageURL, req.Stock)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, product)
}

// DeleteProduct deletes a product
func (h *BusinessHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	productUUID, err := uuid.Parse(productID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = h.businessService.DeleteProduct(productUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Product deleted successfully"})
}

// CreateOrder creates an order
func (h *BusinessHandler) CreateOrder(c *gin.Context) {
	var req struct {
		CustomerName  string `json:"customer_name" binding:"required"`
		CustomerPhone string `json:"customer_phone" binding:"required"`
		Items         []struct {
			ProductID string `json:"product_id" binding:"required"`
			Quantity  int    `json:"quantity" binding:"required,min=1"`
		} `json:"items" binding:"required,min=1"`
		UserID string `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Convert product IDs
	items := make([]models.OrderItem, len(req.Items))
	for i, item := range req.Items {
		productID, err := uuid.Parse(item.ProductID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid product ID")
			return
		}
		items[i] = models.OrderItem{
			ProductID: productID,
			Quantity:  item.Quantity,
		}
	}

	order, err := h.businessService.CreateOrder(req.CustomerName, req.CustomerPhone, items, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, order)
}

// GetOrders gets orders
func (h *BusinessHandler) GetOrders(c *gin.Context) {
	status := c.Query("status")
	userID := c.Query("user_id")
	page := utils.Paginate(0, 1, 10) // Default pagination

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	orders, total, err := h.businessService.GetOrders(status, userUUID, 1, 10)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, 1, 10)
	utils.ResponsePaginated(c, orders, pagination)
}

// GetOrder gets an order
func (h *BusinessHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.businessService.GetOrder(orderUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "Order not found")
		return
	}

	utils.ResponseSuccess(c, order)
}

// UpdateOrderStatus updates order status
func (h *BusinessHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.businessService.UpdateOrderStatus(orderUUID, req.Status)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, order)
}

// CreateInvoice creates an invoice
func (h *BusinessHandler) CreateInvoice(c *gin.Context) {
	var req struct {
		OrderID     string  `json:"order_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,min=0"`
		Description string  `json:"description"`
		DueDate     string  `json:"due_date" binding:"required"`
		UserID      string  `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	dueDate, err := time.Parse("YYYY-MM-DD", req.DueDate)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid due date format")
		return
	}

	invoice, err := h.businessService.CreateInvoice(orderID, req.Amount, req.Description, dueDate, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, invoice)
}

// GetInvoices gets invoices
func (h *BusinessHandler) GetInvoices(c *gin.Context) {
	status := c.Query("status")
	userID := c.Query("user_id")
	page := utils.Paginate(0, 1, 10) // Default pagination

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	invoices, total, err := h.businessService.GetInvoices(status, userUUID, 1, 10)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, 1, 10)
	utils.ResponsePaginated(c, invoices, pagination)
}

// GetInvoice gets an invoice
func (h *BusinessHandler) GetInvoice(c *gin.Context) {
	invoiceID := c.Param("invoice_id")
	if invoiceID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Invoice ID is required")
		return
	}

	invoiceUUID, err := uuid.Parse(invoiceID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	invoice, err := h.businessService.GetInvoice(invoiceUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "Invoice not found")
		return
	}

	utils.ResponseSuccess(c, invoice)
}

// UpdateInvoiceStatus updates invoice status
func (h *BusinessHandler) UpdateInvoiceStatus(c *gin.Context) {
	invoiceID := c.Param("invoice_id")
	if invoiceID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Invoice ID is required")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending paid overdue cancelled"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	invoiceUUID, err := uuid.Parse(invoiceID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	invoice, err := h.businessService.UpdateInvoiceStatus(invoiceUUID, req.Status)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, invoice)
}

// CreateCustomer creates a customer
func (h *BusinessHandler) CreateCustomer(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Email   string `json:"email" binding:"required,email"`
		Phone   string `json:"phone" binding:"required"`
		Address string `json:"address"`
		UserID  string `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	customer, err := h.businessService.CreateCustomer(req.Name, req.Email, req.Phone, req.Address, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, customer)
}

// GetCustomers gets customers
func (h *BusinessHandler) GetCustomers(c *gin.Context) {
	userID := c.Query("user_id")
	search := c.Query("search")
	page := utils.Paginate(0, 1, 10) // Default pagination

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	customers, total, err := h.businessService.GetCustomers(userUUID, search, 1, 10)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, 1, 10)
	utils.ResponsePaginated(c, customers, pagination)
}

// GetCustomer gets a customer
func (h *BusinessHandler) GetCustomer(c *gin.Context) {
	customerID := c.Param("customer_id")
	if customerID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Customer ID is required")
		return
	}

	customerUUID, err := uuid.Parse(customerID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.businessService.GetCustomer(customerUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "Customer not found")
		return
	}

	utils.ResponseSuccess(c, customer)
}

// UpdateCustomer updates a customer
func (h *BusinessHandler) UpdateCustomer(c *gin.Context) {
	customerID := c.Param("customer_id")
	if customerID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Customer ID is required")
		return
	}

	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email" binding:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	customerUUID, err := uuid.Parse(customerID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.businessService.UpdateCustomer(customerUUID, req.Name, req.Email, req.Phone, req.Address)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, customer)
}

// DeleteCustomer deletes a customer
func (h *BusinessHandler) DeleteCustomer(c *gin.Context) {
	customerID := c.Param("customer_id")
	if customerID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Customer ID is required")
		return
	}

	customerUUID, err := uuid.Parse(customerID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	err = h.businessService.DeleteCustomer(customerUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, gin.H{"message": "Customer deleted successfully"})
}

// CreatePayment creates a payment
func (h *BusinessHandler) CreatePayment(c *gin.Context) {
	var req struct {
		InvoiceID     string  `json:"invoice_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required,min=0"`
		PaymentMethod string  `json:"payment_method" binding:"required,oneof=cash transfer credit_card digital_wallet"`
		Reference     string  `json:"reference"`
		UserID        string  `json:"user_id" binding:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	invoiceID, err := uuid.Parse(req.InvoiceID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	payment, err := h.businessService.CreatePayment(invoiceID, req.Amount, req.PaymentMethod, req.Reference, userID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, payment)
}

// GetPayments gets payments
func (h *BusinessHandler) GetPayments(c *gin.Context) {
	invoiceID := c.Query("invoice_id")
	userID := c.Query("user_id")
	page := utils.Paginate(0, 1, 10) // Default pagination

	var invoiceUUID *uuid.UUID
	if invoiceID != "" {
		uuid, err := uuid.Parse(invoiceID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid invoice ID")
			return
		}
		invoiceUUID = &uuid
	}

	var userUUID *uuid.UUID
	if userID != "" {
		uuid, err := uuid.Parse(userID)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
			return
		}
		userUUID = &uuid
	}

	payments, total, err := h.businessService.GetPayments(invoiceUUID, userUUID, 1, 10)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := utils.Paginate(total, 1, 10)
	utils.ResponsePaginated(c, payments, pagination)
}

// GetPayment gets a payment
func (h *BusinessHandler) GetPayment(c *gin.Context) {
	paymentID := c.Param("payment_id")
	if paymentID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "Payment ID is required")
		return
	}

	paymentUUID, err := uuid.Parse(paymentID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	payment, err := h.businessService.GetPayment(paymentUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusNotFound, "Payment not found")
		return
	}

	utils.ResponseSuccess(c, payment)
}

// GetBusinessStats gets business statistics
func (h *BusinessHandler) GetBusinessStats(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	stats, err := h.businessService.GetBusinessStats(userUUID)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, stats)
}

// GetSalesReport gets sales report
func (h *BusinessHandler) GetSalesReport(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		utils.ResponseError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	start, err := time.Parse("YYYY-MM-DD", startDate)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid start date format")
		return
	}

	end, err := time.Parse("YYYY-MM-DD", endDate)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid end date format")
		return
	}

	report, err := h.businessService.GetSalesReport(userUUID, start, end)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(c, report)
}