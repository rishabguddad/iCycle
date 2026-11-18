package catalog

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// fileData is the shape of the YAML file.
//
// The yaml package, like encoding/json, can only fill exported fields.
// Product.stock is private, so the file is first read into these plain
// "record" structs, then checked and turned into real Categories and Products.
type fileData struct {
	Categories []categoryRecord `yaml:"categories"`
	Products   []productRecord  `yaml:"products"`
}

type categoryRecord struct {
	ID   int    `yaml:"id"`
	Name string `yaml:"name"`
}

type productRecord struct {
	ID               int               `yaml:"id"`
	Name             string            `yaml:"name"`
	Description      string            `yaml:"description"`
	CategoryID       int               `yaml:"categoryId"`
	Condition        string            `yaml:"condition"`
	Price            int64             `yaml:"price"`
	Stock            int               `yaml:"stock"`
	Specs            map[string]string `yaml:"specs"`
	CompatibleModels []string          `yaml:"compatibleModels"`
}

// LoadFile reads categories and products from a YAML file, checks them,
// and returns a store holding them.
func LoadFile(path string) (*Store, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read data file: %w", err)
	}

	var data fileData
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true) // a typo such as "stok:" is an error, not silently ignored
	err = decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("parse data file %s: %w", path, err)
	}

	categories, err := toCategories(data.Categories)
	if err != nil {
		return nil, fmt.Errorf("data file %s: %w", path, err)
	}

	products, err := toProducts(data.Products, categories)
	if err != nil {
		return nil, fmt.Errorf("data file %s: %w", path, err)
	}

	return NewStore(categories, products), nil
}

// toCategories checks the category records and converts them.
func toCategories(records []categoryRecord) ([]Category, error) {
	seen := make(map[int]bool) // a map used as a set: which IDs we've had
	var categories []Category

	for _, r := range records {
		if r.ID <= 0 {
			return nil, fmt.Errorf("category %q: id must be a positive number", r.Name)
		}
		if seen[r.ID] {
			return nil, fmt.Errorf("category id %d is used twice", r.ID)
		}
		if r.Name == "" {
			return nil, fmt.Errorf("category %d: name is required", r.ID)
		}
		seen[r.ID] = true
		categories = append(categories, Category{ID: r.ID, Name: r.Name})
	}

	if len(categories) == 0 {
		return nil, errors.New("no categories found")
	}
	return categories, nil
}

// toProducts checks the product records and converts them.
// Every product must belong to one of the given categories.
func toProducts(records []productRecord, categories []Category) ([]Product, error) {
	categoryIDs := make(map[int]bool)
	for _, c := range categories {
		categoryIDs[c.ID] = true
	}

	seen := make(map[int]bool)
	var products []Product

	for _, r := range records {
		switch {
		case r.ID <= 0:
			return nil, fmt.Errorf("product %q: id must be a positive number", r.Name)
		case seen[r.ID]:
			return nil, fmt.Errorf("product id %d is used twice", r.ID)
		case r.Stock < 0:
			return nil, fmt.Errorf("product %d: stock can't be negative", r.ID)
		}

		details := ProductDetails{
			Name:             r.Name,
			Description:      r.Description,
			CategoryID:       r.CategoryID,
			Condition:        Condition(r.Condition),
			Price:            r.Price,
			Specs:            r.Specs,
			CompatibleModels: r.CompatibleModels,
		}
		// The same checks the admin API uses, so the rules live in one place.
		err := details.validate(categoryIDs)
		if err != nil {
			return nil, fmt.Errorf("product %d: %w", r.ID, err)
		}

		seen[r.ID] = true
		products = append(products, Product{
			ID:             r.ID,
			ProductDetails: details,
			stock:          r.Stock, // allowed: we are inside package catalog
		})
	}
	return products, nil
}
