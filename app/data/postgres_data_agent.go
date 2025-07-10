package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"drm-app/app/db"
)

type PostgresDataAgent struct {
	db *db.Database
}

func NewPostgresDataAgent(database *db.Database) *PostgresDataAgent {
	return &PostgresDataAgent{
		db: database,
	}
}

func (p *PostgresDataAgent) ExecuteCommand(ctx context.Context, command *Command) (interface{}, error) {
	switch command.Action {
	case "create":
		return p.create(ctx, command.Entity, command.Data)
	case "read":
		return p.read(ctx, command.Entity, command.Data)
	case "update":
		return p.update(ctx, command.Entity, command.Data)
	case "delete":
		return p.delete(ctx, command.Entity, command.Data)
	default:
		return nil, fmt.Errorf("unsupported action: %s", command.Action)
	}
}

func (p *PostgresDataAgent) create(ctx context.Context, entity string, data map[string]interface{}) (interface{}, error) {
	switch entity {
	case "user":
		return p.createUser(ctx, data)
	case "product":
		return p.createProduct(ctx, data)
	case "order":
		return p.createOrder(ctx, data)
	default:
		return nil, fmt.Errorf("unsupported entity: %s", entity)
	}
}

func (p *PostgresDataAgent) read(ctx context.Context, entity string, data map[string]interface{}) (interface{}, error) {
	switch entity {
	case "user":
		return p.readUser(ctx, data)
	case "product":
		return p.readProduct(ctx, data)
	case "order":
		return p.readOrder(ctx, data)
	default:
		return nil, fmt.Errorf("unsupported entity: %s", entity)
	}
}

func (p *PostgresDataAgent) update(ctx context.Context, entity string, data map[string]interface{}) (interface{}, error) {
	switch entity {
	case "user":
		return p.updateUser(ctx, data)
	case "product":
		return p.updateProduct(ctx, data)
	case "order":
		return p.updateOrder(ctx, data)
	default:
		return nil, fmt.Errorf("unsupported entity: %s", entity)
	}
}

func (p *PostgresDataAgent) delete(ctx context.Context, entity string, data map[string]interface{}) (interface{}, error) {
	switch entity {
	case "user":
		return p.deleteUser(ctx, data)
	case "product":
		return p.deleteProduct(ctx, data)
	case "order":
		return p.deleteOrder(ctx, data)
	default:
		return nil, fmt.Errorf("unsupported entity: %s", entity)
	}
}

// User operations
func (p *PostgresDataAgent) createUser(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	name, ok := data["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("user name is required")
	}

	email, ok := data["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("user email is required")
	}

	var user User
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at, updated_at`
	err := p.db.DB.QueryRowContext(ctx, query, name, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (p *PostgresDataAgent) readUser(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	if idStr, ok := data["id"].(string); ok {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID: %w", err)
		}

		var user User
		query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
		err = p.db.DB.QueryRowContext(ctx, query, id).Scan(
			&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}

		return user, nil
	}

	var users []User
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id`
	err := p.db.DB.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to read users: %w", err)
	}

	return users, nil
}

func (p *PostgresDataAgent) updateUser(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	idStr, ok := data["id"].(string)
	if !ok {
		return nil, fmt.Errorf("user ID is required for update")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if name, ok := data["name"].(string); ok && name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, name)
		argIndex++
	}

	if email, ok := data["email"].(string); ok && email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, email)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf(`UPDATE users SET %s WHERE id = $%d RETURNING id, name, email, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex)

	var user User
	err = p.db.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (p *PostgresDataAgent) deleteUser(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	idStr, ok := data["id"].(string)
	if !ok {
		return nil, fmt.Errorf("user ID is required for delete")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	query := `DELETE FROM users WHERE id = $1`
	result, err := p.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return map[string]string{"message": "user deleted successfully"}, nil
}

// Product operations
func (p *PostgresDataAgent) createProduct(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	name, ok := data["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("product name is required")
	}

	price, ok := data["price"].(float64)
	if !ok || price <= 0 {
		return nil, fmt.Errorf("product price must be positive")
	}

	description, _ := data["description"].(string)

	var product Product
	query := `INSERT INTO products (name, price, description) VALUES ($1, $2, $3) RETURNING id, name, price, description, created_at, updated_at`
	err := p.db.DB.QueryRowContext(ctx, query, name, price, description).Scan(
		&product.ID, &product.Name, &product.Price, &product.Description, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (p *PostgresDataAgent) readProduct(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	if idStr, ok := data["id"].(string); ok {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %w", err)
		}

		var product Product
		query := `SELECT id, name, price, description, created_at, updated_at FROM products WHERE id = $1`
		err = p.db.DB.QueryRowContext(ctx, query, id).Scan(
			&product.ID, &product.Name, &product.Price, &product.Description, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}

		return product, nil
	}

	var products []Product
	query := `SELECT id, name, price, description, created_at, updated_at FROM products ORDER BY id`
	err := p.db.DB.SelectContext(ctx, &products, query)
	if err != nil {
		return nil, fmt.Errorf("failed to read products: %w", err)
	}

	return products, nil
}

func (p *PostgresDataAgent) updateProduct(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	idStr, ok := data["id"].(string)
	if !ok {
		return nil, fmt.Errorf("product ID is required for update")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if name, ok := data["name"].(string); ok && name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, name)
		argIndex++
	}

	if price, ok := data["price"].(float64); ok && price > 0 {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, price)
		argIndex++
	}

	if description, ok := data["description"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, description)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf(`UPDATE products SET %s WHERE id = $%d RETURNING id, name, price, description, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex)

	var product Product
	err = p.db.DB.QueryRowContext(ctx, query, args...).Scan(
		&product.ID, &product.Name, &product.Price, &product.Description, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (p *PostgresDataAgent) deleteProduct(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	idStr, ok := data["id"].(string)
	if !ok {
		return nil, fmt.Errorf("product ID is required for delete")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	query := `DELETE FROM products WHERE id = $1`
	result, err := p.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("product not found")
	}

	return map[string]string{"message": "product deleted successfully"}, nil
}

// Order operations
func (p *PostgresDataAgent) createOrder(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	items, ok := data["items"]
	if !ok {
		return nil, fmt.Errorf("order items are required")
	}

	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal items: %w", err)
	}

	var userID *int
	if userIDStr, ok := data["user_id"].(string); ok {
		uid, err := strconv.Atoi(userIDStr)
		if err == nil {
			userID = &uid
		}
	}

	var totalAmount *float64
	if ta, ok := data["total_amount"].(float64); ok {
		totalAmount = &ta
	}

	status, _ := data["status"].(string)
	if status == "" {
		status = "pending"
	}

	var order Order
	query := `INSERT INTO orders (user_id, items, total_amount, status) VALUES ($1, $2, $3, $4) RETURNING id, user_id, items, total_amount, status, created_at, updated_at`
	err = p.db.DB.QueryRowContext(ctx, query, userID, string(itemsJSON), totalAmount, status).Scan(
		&order.ID, &order.UserID, &order.Items, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

func (p *PostgresDataAgent) readOrder(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	if idStr, ok := data["id"].(string); ok {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid order ID: %w", err)
		}

		var order Order
		query := `SELECT id, user_id, items, total_amount, status, created_at, updated_at FROM orders WHERE id = $1`
		err = p.db.DB.QueryRowContext(ctx, query, id).Scan(
			&order.ID, &order.UserID, &order.Items, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("order not found: %w", err)
		}

		return order, nil
	}

	var orders []Order
	query := `SELECT id, user_id, items, total_amount, status, created_at, updated_at FROM orders ORDER BY id`
	err := p.db.DB.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, fmt.Errorf("failed to read orders: %w", err)
	}

	return orders, nil
}

func (p *PostgresDataAgent) updateOrder(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	idStr, ok := data["id"].(string)
	if !ok {
		return nil, fmt.Errorf("order ID is required for update")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if items, ok := data["items"]; ok {
		itemsJSON, err := json.Marshal(items)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal items: %w", err)
		}
		setParts = append(setParts, fmt.Sprintf("items = $%d", argIndex))
		args = append(args, string(itemsJSON))
		argIndex++
	}

	if totalAmount, ok := data["total_amount"].(float64); ok {
		setParts = append(setParts, fmt.Sprintf("total_amount = $%d", argIndex))
		args = append(args, totalAmount)
		argIndex++
	}

	if status, ok := data["status"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf(`UPDATE orders SET %s WHERE id = $%d RETURNING id, user_id, items, total_amount, status, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex)

	var order Order
	err = p.db.DB.QueryRowContext(ctx, query, args...).Scan(
		&order.ID, &order.UserID, &order.Items, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}

func (p *PostgresDataAgent) deleteOrder(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	idStr, ok := data["id"].(string)
	if !ok {
		return nil, fmt.Errorf("order ID is required for delete")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	query := `DELETE FROM orders WHERE id = $1`
	result, err := p.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("order not found")
	}

	return map[string]string{"message": "order deleted successfully"}, nil
}
