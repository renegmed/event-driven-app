package repository

import (
	"fmt"
	"job-runner-app/internal/config"
)

func BuildDataSourceName(c config.DBConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
	// return fmt.Sprintf(
	// 	"host=%s port=%s password=%s dbname=%s sslmode=disable",
	// 	c.Host,
	// 	c.Port,
	// 	c.User,
	// 	c.Password,
	// 	c.DBName,
	// )
}
