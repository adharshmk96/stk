package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlBuilder(t *testing.T) {
	t.Run("should build a simple query", func(t *testing.T) {
		query := NewSqlQuery().Select("id", "name").From("users").Build()
		expected := "SELECT (id, name) FROM users"
		assert.Equal(t, expected, query)
	})

	t.Run("should build a query with where clause", func(t *testing.T) {
		query := NewSqlQuery().Select("id", "name").From("users").Where("id = 1").Build()
		expected := "SELECT (id, name) FROM users WHERE id = 1"
		assert.Equal(t, expected, query)
	})

	t.Run("should build a query with order by clause", func(t *testing.T) {
		query := NewSqlQuery().Select("id", "name").From("users").OrderBy("name").Build()
		expected := "SELECT (id, name) FROM users ORDER BY name"
		assert.Equal(t, expected, query)
	})

	t.Run("should build a query with multiple where and order by clause", func(t *testing.T) {
		query := NewSqlQuery().Select("id", "name").From("users").Where("id = 1", "name = 'John'").OrderBy("name").Build()
		expected := "SELECT (id, name) FROM users WHERE id = 1 AND name = 'John' ORDER BY name"
		assert.Equal(t, expected, query)
	})

	t.Run("should build a query with join clause", func(t *testing.T) {
		query := NewSqlQuery().Select("id", "name").From("users").Join("roles").On("users.role_id = roles.id").Build()
		expected := "SELECT (id, name) FROM users JOIN roles ON users.role_id = roles.id"
		assert.Equal(t, expected, query)
	})
}
