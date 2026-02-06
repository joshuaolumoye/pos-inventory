package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joshuaolumoye/pos-backend/internal/domain"
)

type ProductRepo struct {
	DB *sqlx.DB
}

// GetLowStockCount returns the count of low stock products (optionally filtered by branch)
func (r *ProductRepo) GetLowStockCount(businessID, branchID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM products WHERE business_id = ? AND (deleted_at IS NULL OR deleted_at = 0) AND quantity_in_stock <= low_stock_threshold`
	args := []interface{}{businessID}
	if branchID != "" {
		query += " AND branch_id = ?"
		args = append(args, branchID)
	}
	err := r.DB.QueryRowx(query, args...).Scan(&count)
	return count, err
}

func (r *ProductRepo) CreateProduct(p *domain.Product) error {
	query := `INSERT INTO products (
		id, product_name, product_category, business_id, branch_id, 
		barcode_value, nafdac_reg_number, selling_price, cost_price, 
		quantity_in_stock, low_stock_threshold, expiry_date, product_image_url, 
		created_at, updated_at, created_by
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.DB.Exec(query,
		p.ID, p.ProductName, p.ProductCategory, p.BusinessID, p.BranchID,
		p.BarcodeValue, p.NAFDACRegNumber, p.SellingPrice, p.CostPrice,
		p.QuantityInStock, p.LowStockThreshold, p.ExpiryDate, p.ProductImageURL,
		p.CreatedAt, p.UpdatedAt, p.CreatedBy)
	return err
}

func (r *ProductRepo) GetProductByID(productID string) (*domain.Product, error) {
	var (
		p                    domain.Product
		createdAt, updatedAt int64
		deletedAt            *int64
		updatedBy            *string
	)
	query := `SELECT 
		       id, product_name, product_category, business_id, branch_id,
		       barcode_value, nafdac_reg_number, selling_price, cost_price,
		       quantity_in_stock, low_stock_threshold, expiry_date, product_image_url,
		       created_at, updated_at, deleted_at, created_by, updated_by
	       FROM products 
	       WHERE id = ? AND (deleted_at IS NULL OR deleted_at = 0)`

	row := r.DB.QueryRowx(query, productID)
	err := row.Scan(
		&p.ID, &p.ProductName, &p.ProductCategory, &p.BusinessID, &p.BranchID,
		&p.BarcodeValue, &p.NAFDACRegNumber, &p.SellingPrice, &p.CostPrice,
		&p.QuantityInStock, &p.LowStockThreshold, &p.ExpiryDate, &p.ProductImageURL,
		&createdAt, &updatedAt, &deletedAt, &p.CreatedBy, &p.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}
	p.CreatedAt = createdAt
	p.UpdatedAt = updatedAt
	p.DeletedAt = deletedAt
	p.UpdatedBy = updatedBy
	return &p, nil
}

func (r *ProductRepo) GetProductsByBusinessID(businessID string) ([]*domain.Product, error) {
	query := `SELECT 
		       id, product_name, product_category, business_id, branch_id,
		       barcode_value, nafdac_reg_number, selling_price, cost_price,
		       quantity_in_stock, low_stock_threshold, expiry_date, product_image_url,
		       created_at, updated_at, deleted_at, created_by, updated_by
	       FROM products 
	       WHERE business_id = ? AND (deleted_at IS NULL OR deleted_at = 0)
	       ORDER BY created_at DESC`

	rows, err := r.DB.Queryx(query, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*domain.Product
	for rows.Next() {
		var (
			p                    domain.Product
			createdAt, updatedAt int64
			deletedAt            *int64
			updatedBy            *string
		)
		err := rows.Scan(
			&p.ID, &p.ProductName, &p.ProductCategory, &p.BusinessID, &p.BranchID,
			&p.BarcodeValue, &p.NAFDACRegNumber, &p.SellingPrice, &p.CostPrice,
			&p.QuantityInStock, &p.LowStockThreshold, &p.ExpiryDate, &p.ProductImageURL,
			&createdAt, &updatedAt, &deletedAt, &p.CreatedBy, &p.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		p.CreatedAt = createdAt
		p.UpdatedAt = updatedAt
		p.DeletedAt = deletedAt
		p.UpdatedBy = updatedBy
		products = append(products, &p)
	}
	return products, nil
}

func (r *ProductRepo) GetProductsByBranchID(businessID, branchID string) ([]*domain.Product, error) {
	query := `SELECT 
		       id, product_name, product_category, business_id, branch_id,
		       barcode_value, nafdac_reg_number, selling_price, cost_price,
		       quantity_in_stock, low_stock_threshold, expiry_date, product_image_url,
		       created_at, updated_at, deleted_at, created_by, updated_by
	       FROM products 
	       WHERE business_id = ? AND branch_id = ? AND (deleted_at IS NULL OR deleted_at = 0)
	       ORDER BY created_at DESC`

	rows, err := r.DB.Queryx(query, businessID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*domain.Product
	for rows.Next() {
		var (
			p                    domain.Product
			createdAt, updatedAt int64
			deletedAt            *int64
			updatedBy            *string
		)
		err := rows.Scan(
			&p.ID, &p.ProductName, &p.ProductCategory, &p.BusinessID, &p.BranchID,
			&p.BarcodeValue, &p.NAFDACRegNumber, &p.SellingPrice, &p.CostPrice,
			&p.QuantityInStock, &p.LowStockThreshold, &p.ExpiryDate, &p.ProductImageURL,
			&createdAt, &updatedAt, &deletedAt, &p.CreatedBy, &updatedBy,
		)
		if err != nil {
			return nil, err
		}
		p.CreatedAt = createdAt
		p.UpdatedAt = updatedAt
		p.DeletedAt = deletedAt
		p.UpdatedBy = updatedBy
		products = append(products, &p)
	}
	return products, nil
}

func (r *ProductRepo) UpdateProduct(p *domain.Product) error {
	query := `UPDATE products SET 
		product_name = ?, product_category = ?, selling_price = ?, 
		cost_price = ?, quantity_in_stock = ?, low_stock_threshold = ?, 
		barcode_value = ?, nafdac_reg_number = ?, expiry_date = ?, 
		product_image_url = ?, branch_id = ?, updated_at = ?, updated_by = ?
	WHERE id = ? AND business_id = ? AND (deleted_at IS NULL OR deleted_at = 0)`

	_, err := r.DB.Exec(query,
		p.ProductName, p.ProductCategory, p.SellingPrice, p.CostPrice,
		p.QuantityInStock, p.LowStockThreshold, p.BarcodeValue,
		p.NAFDACRegNumber, p.ExpiryDate, p.ProductImageURL, p.BranchID,
		p.UpdatedAt, p.UpdatedBy, p.ID, p.BusinessID)
	return err
}

func (r *ProductRepo) DeleteProduct(productID string) error {
	query := `UPDATE products SET deleted_at = ? WHERE id = ?`
	_, err := r.DB.Exec(query, time.Now().Unix(), productID)
	return err
}

func (r *ProductRepo) SearchProducts(businessID, searchTerm string) ([]*domain.Product, error) {
	query := `SELECT 
		       id, product_name, product_category, business_id, branch_id,
		       barcode_value, nafdac_reg_number, selling_price, cost_price,
		       quantity_in_stock, low_stock_threshold, expiry_date, product_image_url,
		       created_at, updated_at, deleted_at, created_by, updated_by
	       FROM products 
	       WHERE business_id = ? 
		       AND (deleted_at IS NULL OR deleted_at = 0)
		       AND (
			       product_name LIKE ? 
			       OR product_category LIKE ?
			       OR barcode_value LIKE ?
			       OR nafdac_reg_number LIKE ?
		       )
	       ORDER BY created_at DESC`

	searchPattern := "%" + searchTerm + "%"
	rows, err := r.DB.Queryx(query, businessID, searchPattern, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*domain.Product
	for rows.Next() {
		var (
			p                    domain.Product
			createdAt, updatedAt int64
			deletedAt            *int64
			updatedBy            *string
		)
		err := rows.Scan(
			&p.ID, &p.ProductName, &p.ProductCategory, &p.BusinessID, &p.BranchID,
			&p.BarcodeValue, &p.NAFDACRegNumber, &p.SellingPrice, &p.CostPrice,
			&p.QuantityInStock, &p.LowStockThreshold, &p.ExpiryDate, &p.ProductImageURL,
			&createdAt, &updatedAt, &deletedAt, &p.CreatedBy, &updatedBy,
		)
		if err != nil {
			return nil, err
		}
		p.CreatedAt = createdAt
		p.UpdatedAt = updatedAt
		p.DeletedAt = deletedAt
		p.UpdatedBy = updatedBy
		products = append(products, &p)
	}
	return products, nil
}

func (r *ProductRepo) GetLowStockProducts(businessID string) ([]*domain.Product, error) {
	query := `SELECT 
		       id, product_name, product_category, business_id, branch_id,
		       barcode_value, nafdac_reg_number, selling_price, cost_price,
		       quantity_in_stock, low_stock_threshold, expiry_date, product_image_url,
		       created_at, updated_at, deleted_at, created_by, updated_by
	       FROM products 
	       WHERE business_id = ? 
		       AND (deleted_at IS NULL OR deleted_at = 0)
		       AND quantity_in_stock <= low_stock_threshold
	       ORDER BY quantity_in_stock ASC`

	rows, err := r.DB.Queryx(query, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*domain.Product
	for rows.Next() {
		var (
			p                    domain.Product
			createdAt, updatedAt int64
			deletedAt            *int64
			updatedBy            *string
		)
		err := rows.Scan(
			&p.ID, &p.ProductName, &p.ProductCategory, &p.BusinessID, &p.BranchID,
			&p.BarcodeValue, &p.NAFDACRegNumber, &p.SellingPrice, &p.CostPrice,
			&p.QuantityInStock, &p.LowStockThreshold, &p.ExpiryDate, &p.ProductImageURL,
			&createdAt, &updatedAt, &deletedAt, &p.CreatedBy, &updatedBy,
		)
		if err != nil {
			return nil, err
		}
		p.CreatedAt = createdAt
		p.UpdatedAt = updatedAt
		p.DeletedAt = deletedAt
		p.UpdatedBy = updatedBy
		products = append(products, &p)
	}
	return products, nil
}

func (r *ProductRepo) UpdateProductStock(productID string, quantity int) error {
	query := `UPDATE products SET 
		quantity_in_stock = ?, 
		updated_at = ? 
	WHERE id = ? AND (deleted_at IS NULL OR deleted_at = 0)`

	_, err := r.DB.Exec(query, quantity, time.Now().Unix(), productID)
	return err
}

// GetProductsByBranch - backward compatibility method
func (r *ProductRepo) GetProductsByBranch(branchID string) ([]*domain.Product, error) {
	query := `SELECT 
		       id, product_name, product_category, business_id, branch_id,
		       barcode_value, nafdac_reg_number, selling_price, cost_price,
		       quantity_in_stock, low_stock_threshold, expiry_date, product_image_url,
		       created_at, updated_at, deleted_at, created_by, updated_by
	       FROM products 
	       WHERE branch_id = ? AND (deleted_at IS NULL OR deleted_at = 0)
	       ORDER BY created_at DESC`

	rows, err := r.DB.Queryx(query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*domain.Product
	for rows.Next() {
		var (
			p                    domain.Product
			createdAt, updatedAt int64
			deletedAt            *int64
			updatedBy            *string
		)
		err := rows.Scan(
			&p.ID, &p.ProductName, &p.ProductCategory, &p.BusinessID, &p.BranchID,
			&p.BarcodeValue, &p.NAFDACRegNumber, &p.SellingPrice, &p.CostPrice,
			&p.QuantityInStock, &p.LowStockThreshold, &p.ExpiryDate, &p.ProductImageURL,
			&createdAt, &updatedAt, &deletedAt, &p.CreatedBy, &updatedBy,
		)
		if err != nil {
			return nil, err
		}
		p.CreatedAt = createdAt
		p.UpdatedAt = updatedAt
		p.DeletedAt = deletedAt
		p.UpdatedBy = updatedBy
		products = append(products, &p)
	}
	return products, nil
}

// QueryProductsNotification returns products for notification (in_stock, low_stock, expired) with pagination
func (r *ProductRepo) QueryProductsNotification(businessID, op string, stock int, expiry int64, lowStock int, limit, offset int, expired bool) ([]*domain.Product, error) {
	var (
		products []*domain.Product
		rows     *sqlx.Rows
		err      error
	)
	base := `SELECT id, product_name, barcode_value, selling_price, quantity_in_stock, expiry_date FROM products WHERE business_id = ? AND (deleted_at IS NULL OR deleted_at = 0)`
	var args []interface{}
	args = append(args, businessID)
	if expired {
		base += " AND expiry_date IS NOT NULL AND expiry_date < ?"
		args = append(args, expiry)
	} else if op == ">" || op == "<" {
		base += " AND quantity_in_stock " + op + " ?"
		args = append(args, stock)
	}
	base += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	rows, err = r.DB.Queryx(base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.ProductName, &p.BarcodeValue, &p.SellingPrice, &p.QuantityInStock, &p.ExpiryDate)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

// GetAllProductsPaginated returns all products for a business, paginated
func (r *ProductRepo) GetAllProductsPaginated(businessID string, limit, offset int) ([]*domain.Product, error) {
	query := `SELECT id, product_name, barcode_value, selling_price, quantity_in_stock, expiry_date FROM products WHERE business_id = ? AND (deleted_at IS NULL OR deleted_at = 0) ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.DB.Queryx(query, businessID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.ProductName, &p.BarcodeValue, &p.SellingPrice, &p.QuantityInStock, &p.ExpiryDate)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

// SearchProductsPaginated returns products matching search (fuzzy), paginated
func (r *ProductRepo) SearchProductsPaginated(businessID, search string, limit, offset int) ([]*domain.Product, error) {
	pattern := "%" + search + "%"
	query := `SELECT id, product_name, barcode_value, selling_price, quantity_in_stock, expiry_date FROM products WHERE business_id = ? AND (deleted_at IS NULL OR deleted_at = 0) AND product_name LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.DB.Queryx(query, businessID, pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.ProductName, &p.BarcodeValue, &p.SellingPrice, &p.QuantityInStock, &p.ExpiryDate)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}
