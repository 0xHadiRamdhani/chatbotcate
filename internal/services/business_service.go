package services

import (
	"time"

	"github.com/google/uuid"
	"kilocode.dev/whatsapp-bot/internal/models"
	"kilocode.dev/whatsapp-bot/pkg/logger"
)

type BusinessService struct {
	db *Database
}

func NewBusinessService(db *Database) *BusinessService {
	return &BusinessService{
		db: db,
	}
}

// CreateProduct creates a product
func (s *BusinessService) CreateProduct(name string, description string, price float64, category string, imageURL string, stock int, userID uuid.UUID) (*models.Product, error) {
	product := &models.Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		ImageURL:    imageURL,
		Stock:       stock,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.DB.Create(product).Error; err != nil {
		logger.Error("Failed to create product", err)
		return nil, err
	}

	return product, nil
}

// GetProducts gets products
func (s *BusinessService) GetProducts(category string, search string, page int, limit int) ([]models.Product, int, error) {
	var products []models.Product
	var total int64

	query := s.db.DB.Model(&models.Product{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		logger.Error("Failed to count products", err)
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		logger.Error("Failed to get products", err)
		return nil, 0, err
	}

	return products, int(total), nil
}

// GetProduct gets a product by ID
func (s *BusinessService) GetProduct(productID uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := s.db.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		logger.Error("Failed to get product", err)
		return nil, err
	}

	return &product, nil
}

// UpdateProduct updates a product
func (s *BusinessService) UpdateProduct(productID uuid.UUID, name string, description string, price float64, category string, imageURL string, stock int) (*models.Product, error) {
	var product models.Product
	if err := s.db.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		logger.Error("Failed to get product for update", err)
		return nil, err
	}

	// Update fields
	if name != "" {
		product.Name = name
	}
	if description != "" {
		product.Description = description
	}
	if price > 0 {
		product.Price = price
	}
	if category != "" {
		product.Category = category
	}
	if imageURL != "" {
		product.ImageURL = imageURL
	}
	if stock >= 0 {
		product.Stock = stock
	}
	product.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&product).Error; err != nil {
		logger.Error("Failed to update product", err)
		return nil, err
	}

	return &product, nil
}

// DeleteProduct deletes a product
func (s *BusinessService) DeleteProduct(productID uuid.UUID) error {
	if err := s.db.DB.Where("id = ?", productID).Delete(&models.Product{}).Error; err != nil {
		logger.Error("Failed to delete product", err)
		return err
	}

	return nil
}

// CreateOrder creates an order
func (s *BusinessService) CreateOrder(customerName string, customerPhone string, items []models.OrderItem, userID uuid.UUID) (*models.Order, error) {
	// Calculate total amount
	totalAmount := 0.0
	for _, item := range items {
		product, err := s.GetProduct(item.ProductID)
		if err != nil {
			logger.Error("Failed to get product for order", err)
			return nil, err
		}
		totalAmount += product.Price * float64(item.Quantity)
	}

	order := &models.Order{
		ID:            uuid.New(),
		CustomerName:  customerName,
		CustomerPhone: customerPhone,
		Items:         items,
		TotalAmount:   totalAmount,
		Status:        "pending",
		UserID:        userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.DB.Create(order).Error; err != nil {
		logger.Error("Failed to create order", err)
		return nil, err
	}

	return order, nil
}

// GetOrders gets orders
func (s *BusinessService) GetOrders(status string, userID *uuid.UUID, page int, limit int) ([]models.Order, int, error) {
	var orders []models.Order
	var total int64

	query := s.db.DB.Model(&models.Order{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		logger.Error("Failed to count orders", err)
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		logger.Error("Failed to get orders", err)
		return nil, 0, err
	}

	return orders, int(total), nil
}

// GetOrder gets an order by ID
func (s *BusinessService) GetOrder(orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := s.db.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		logger.Error("Failed to get order", err)
		return nil, err
	}

	return &order, nil
}

// UpdateOrderStatus updates order status
func (s *BusinessService) UpdateOrderStatus(orderID uuid.UUID, status string) (*models.Order, error) {
	var order models.Order
	if err := s.db.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		logger.Error("Failed to get order for status update", err)
		return nil, err
	}

	order.Status = status
	order.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&order).Error; err != nil {
		logger.Error("Failed to update order status", err)
		return nil, err
	}

	return &order, nil
}

// CreateInvoice creates an invoice
func (s *BusinessService) CreateInvoice(orderID uuid.UUID, amount float64, description string, dueDate time.Time, userID uuid.UUID) (*models.Invoice, error) {
	invoice := &models.Invoice{
		ID:          uuid.New(),
		OrderID:     orderID,
		Amount:      amount,
		Description: description,
		DueDate:     dueDate,
		Status:      "pending",
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.DB.Create(invoice).Error; err != nil {
		logger.Error("Failed to create invoice", err)
		return nil, err
	}

	return invoice, nil
}

// GetInvoices gets invoices
func (s *BusinessService) GetInvoices(status string, userID *uuid.UUID, page int, limit int) ([]models.Invoice, int, error) {
	var invoices []models.Invoice
	var total int64

	query := s.db.DB.Model(&models.Invoice{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		logger.Error("Failed to count invoices", err)
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&invoices).Error; err != nil {
		logger.Error("Failed to get invoices", err)
		return nil, 0, err
	}

	return invoices, int(total), nil
}

// GetInvoice gets an invoice by ID
func (s *BusinessService) GetInvoice(invoiceID uuid.UUID) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := s.db.DB.Where("id = ?", invoiceID).First(&invoice).Error; err != nil {
		logger.Error("Failed to get invoice", err)
		return nil, err
	}

	return &invoice, nil
}

// UpdateInvoiceStatus updates invoice status
func (s *BusinessService) UpdateInvoiceStatus(invoiceID uuid.UUID, status string) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := s.db.DB.Where("id = ?", invoiceID).First(&invoice).Error; err != nil {
		logger.Error("Failed to get invoice for status update", err)
		return nil, err
	}

	invoice.Status = status
	invoice.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&invoice).Error; err != nil {
		logger.Error("Failed to update invoice status", err)
		return nil, err
	}

	return &invoice, nil
}

// CreateCustomer creates a customer
func (s *BusinessService) CreateCustomer(name string, email string, phone string, address string, userID uuid.UUID) (*models.Customer, error) {
	customer := &models.Customer{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Address:   address,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.DB.Create(customer).Error; err != nil {
		logger.Error("Failed to create customer", err)
		return nil, err
	}

	return customer, nil
}

// GetCustomers gets customers
func (s *BusinessService) GetCustomers(userID *uuid.UUID, search string, page int, limit int) ([]models.Customer, int, error) {
	var customers []models.Customer
	var total int64

	query := s.db.DB.Model(&models.Customer{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		logger.Error("Failed to count customers", err)
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&customers).Error; err != nil {
		logger.Error("Failed to get customers", err)
		return nil, 0, err
	}

	return customers, int(total), nil
}

// GetCustomer gets a customer by ID
func (s *BusinessService) GetCustomer(customerID uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := s.db.DB.Where("id = ?", customerID).First(&customer).Error; err != nil {
		logger.Error("Failed to get customer", err)
		return nil, err
	}

	return &customer, nil
}

// UpdateCustomer updates a customer
func (s *BusinessService) UpdateCustomer(customerID uuid.UUID, name string, email string, phone string, address string) (*models.Customer, error) {
	var customer models.Customer
	if err := s.db.DB.Where("id = ?", customerID).First(&customer).Error; err != nil {
		logger.Error("Failed to get customer for update", err)
		return nil, err
	}

	// Update fields
	if name != "" {
		customer.Name = name
	}
	if email != "" {
		customer.Email = email
	}
	if phone != "" {
		customer.Phone = phone
	}
	if address != "" {
		customer.Address = address
	}
	customer.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&customer).Error; err != nil {
		logger.Error("Failed to update customer", err)
		return nil, err
	}

	return &customer, nil
}

// DeleteCustomer deletes a customer
func (s *BusinessService) DeleteCustomer(customerID uuid.UUID) error {
	if err := s.db.DB.Where("id = ?", customerID).Delete(&models.Customer{}).Error; err != nil {
		logger.Error("Failed to delete customer", err)
		return err
	}

	return nil
}

// CreatePayment creates a payment
func (s *BusinessService) CreatePayment(invoiceID uuid.UUID, amount float64, paymentMethod string, reference string, userID uuid.UUID) (*models.Payment, error) {
	payment := &models.Payment{
		ID:            uuid.New(),
		InvoiceID:     invoiceID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Reference:     reference,
		Status:        "completed",
		UserID:        userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.DB.Create(payment).Error; err != nil {
		logger.Error("Failed to create payment", err)
		return nil, err
	}

	// Update invoice status
	if err := s.db.DB.Model(&models.Invoice{}).Where("id = ?", invoiceID).Update("status", "paid").Error; err != nil {
		logger.Error("Failed to update invoice status after payment", err)
	}

	return payment, nil
}

// GetPayments gets payments
func (s *BusinessService) GetPayments(invoiceID *uuid.UUID, userID *uuid.UUID, page int, limit int) ([]models.Payment, int, error) {
	var payments []models.Payment
	var total int64

	query := s.db.DB.Model(&models.Payment{})

	if invoiceID != nil {
		query = query.Where("invoice_id = ?", *invoiceID)
	}

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		logger.Error("Failed to count payments", err)
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&payments).Error; err != nil {
		logger.Error("Failed to get payments", err)
		return nil, 0, err
	}

	return payments, int(total), nil
}

// GetPayment gets a payment by ID
func (s *BusinessService) GetPayment(paymentID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := s.db.DB.Where("id = ?", paymentID).First(&payment).Error; err != nil {
		logger.Error("Failed to get payment", err)
		return nil, err
	}

	return &payment, nil
}

// GetBusinessStats gets business statistics
func (s *BusinessService) GetBusinessStats(userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Product stats
	var productCount int64
	s.db.DB.Model(&models.Product{}).Where("user_id = ?", userID).Count(&productCount)
	stats["total_products"] = productCount

	// Order stats
	var orderCount int64
	var totalRevenue float64
	s.db.DB.Model(&models.Order{}).Where("user_id = ?", userID).Count(&orderCount)
	s.db.DB.Model(&models.Order{}).Where("user_id = ? AND status = ?", userID, "delivered").Select("SUM(total_amount)").Scan(&totalRevenue)
	stats["total_orders"] = orderCount
	stats["total_revenue"] = totalRevenue

	// Customer stats
	var customerCount int64
	s.db.DB.Model(&models.Customer{}).Where("user_id = ?", userID).Count(&customerCount)
	stats["total_customers"] = customerCount

	// Invoice stats
	var invoiceCount int64
	var pendingInvoices int64
	s.db.DB.Model(&models.Invoice{}).Where("user_id = ?", userID).Count(&invoiceCount)
	s.db.DB.Model(&models.Invoice{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&pendingInvoices)
	stats["total_invoices"] = invoiceCount
	stats["pending_invoices"] = pendingInvoices

	return stats, nil
}

// GetSalesReport gets sales report
func (s *BusinessService) GetSalesReport(userID uuid.UUID, startDate time.Time, endDate time.Time) (map[string]interface{}, error) {
	report := make(map[string]interface{})

	// Sales by date
	var salesByDate []map[string]interface{}
	s.db.DB.Model(&models.Order{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Select("DATE(created_at) as date, COUNT(*) as orders, SUM(total_amount) as revenue").
		Group("DATE(created_at)").
		Order("date").
		Scan(&salesByDate)
	report["sales_by_date"] = salesByDate

	// Top products
	var topProducts []map[string]interface{}
	s.db.DB.Model(&models.Product{}).
		Where("user_id = ?", userID).
		Select("name, price, stock").
		Order("price desc").
		Limit(10).
		Scan(&topProducts)
	report["top_products"] = topProducts

	// Order status breakdown
	var statusBreakdown []map[string]interface{}
	s.db.DB.Model(&models.Order{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusBreakdown)
	report["order_status_breakdown"] = statusBreakdown

	// Total summary
	var totalOrders int64
	var totalRevenue float64
	s.db.DB.Model(&models.Order{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Count(&totalOrders)
	s.db.DB.Model(&models.Order{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ? AND status = ?", userID, startDate, endDate, "delivered").
		Select("SUM(total_amount)").Scan(&totalRevenue)

	report["total_orders"] = totalOrders
	report["total_revenue"] = totalRevenue
	report["report_period"] = map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
	}

	return report, nil
}