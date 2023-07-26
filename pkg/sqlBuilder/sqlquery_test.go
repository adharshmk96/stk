package sqlBuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlQueryBuilder(t *testing.T) {
	t.Run("LongInsert", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.InsertInto("users").
			Fields(
				"name",
				"email",
				"age",
				"address",
				"phone",
			).
			Values(
				"'John Doe'",
				"'john@example.com'",
				"30",
				"'123 St, City, State, Country'",
				"'1234567890'",
			).Build()
		expected := "INSERT INTO users (name, email, age, address, phone) VALUES ('John Doe', 'john@example.com', 30, '123 St, City, State, Country', '1234567890')"
		assert.Equal(t, expected, result)
	})

	t.Run("LongUpdate", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Update("users").
			Set(
				"name='John Doe'",
				"email='john@example.com'",
				"age=30",
				"address='123 St, City, State, Country'",
				"phone='1234567890'",
			).
			Where("id=1").Build()
		expected := "UPDATE users SET name='John Doe', email='john@example.com', age=30, address='123 St, City, State, Country', phone='1234567890' WHERE id=1"
		assert.Equal(t, expected, result)
	})

	t.Run("LongSelect", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Select(
			"users.id",
			"users.name",
			"orders.order_id",
		).From("users").
			Join("orders").
			On("users.id=orders.user_id").
			Where("users.age > 18").OrderBy("users.name").
			Build()
		expected := "SELECT users.id, users.name, orders.order_id FROM users JOIN orders ON users.id=orders.user_id WHERE users.age > 18 ORDER BY users.name"
		assert.Equal(t, expected, result)
	})

	t.Run("build resets the query", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Select(
			"users.id",
			"users.name",
			"orders.order_id",
		).From("users").
			Join("orders").
			On("users.id=orders.user_id").
			Where("users.age > 18").OrderBy("users.name").
			Build()
		expected := "SELECT users.id, users.name, orders.order_id FROM users JOIN orders ON users.id=orders.user_id WHERE users.age > 18 ORDER BY users.name"
		assert.Equal(t, expected, result)

		result = builder.Build()
		assert.Equal(t, "", result)
	})

	t.Run("LongDelete", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.DeleteFrom("users").
			Where("id=1").Build()
		expected := "DELETE FROM users WHERE id=1"
		assert.Equal(t, expected, result)
	})

	t.Run("LongSelectWithLimitAndOffset", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Select(
			"users.id",
			"users.name",
			"orders.order_id",
		).From("users").
			Join("orders").
			On("users.id=orders.user_id").
			Where("users.age > 18").OrderBy("users.name").
			Limit("10").Offset("10").
			Build()
		expected := "SELECT users.id, users.name, orders.order_id FROM users JOIN orders ON users.id=orders.user_id WHERE users.age > 18 ORDER BY users.name LIMIT 10 OFFSET 10"
		assert.Equal(t, expected, result)
	})
}

func TestSqlBuilderShort(t *testing.T) {
	t.Run("InsertInto", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.InsertInto("users").Build()
		assert.Equal(t, "INSERT INTO users", result)
	})

	t.Run("Update", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Update("users").Build()
		assert.Equal(t, "UPDATE users", result)
	})

	t.Run("Fields", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Fields("name", "email").Build()
		assert.Equal(t, "(name, email)", result)
	})

	t.Run("Values", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Values("'John'", "'john@example.com'").Build()
		assert.Equal(t, "VALUES ('John', 'john@example.com')", result)
	})

	t.Run("Set", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Set("name='John'").Build()
		assert.Equal(t, "SET name='John'", result)
	})

	t.Run("Select", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Select("name", "email").From("users").Build()
		assert.Equal(t, "SELECT name, email FROM users", result)
	})

	t.Run("From", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.From("users").Build()
		assert.Equal(t, "FROM users", result)
	})

	t.Run("Where", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Where("id=1").Build()
		assert.Equal(t, "WHERE id=1", result)
	})

	t.Run("OrderBy", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.OrderBy("name").Build()
		assert.Equal(t, "ORDER BY name", result)
	})

	t.Run("Join", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Join("orders").Build()
		assert.Equal(t, "JOIN orders", result)
	})

	t.Run("On", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.On("users.id=orders.user_id").Build()
		assert.Equal(t, "ON users.id=orders.user_id", result)
	})

	t.Run("Limit", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Limit("10").Build()
		assert.Equal(t, "LIMIT 10", result)
	})

	t.Run("Offset", func(t *testing.T) {
		builder := NewSqlQuery()
		result := builder.Offset("10").Build()
		assert.Equal(t, "OFFSET 10", result)
	})
}
