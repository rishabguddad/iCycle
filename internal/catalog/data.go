package catalog

var categories = []Category{
	{ID: 1, Name: "MacBooks"},
	{ID: 2, Name: "Components"},
	{ID: 3, Name: "Accessories"},
}

var products = []Product{
	{
		ID:          1,
		Name:        "MacBook Pro 16-inch M2 Pro (2023)",
		Description: "Refurbished and fully tested. Battery health 95%. No scratches or dents.",
		CategoryID:  1,
		Condition:   conditionRefurbished,
		Price:       12999900,
		Stock:       1,
		Specs: map[string]string{
			"Processor": "M2 Pro",
			"Memory":    "16GB",
			"Storage":   "512GB",
		},
	},
	{
		ID:          2,
		Name:        "MacBook Air 13-inch M5 (2026)",
		Description: "Used for 6 months. Light scratches on the lid. Battery health 98%.",
		CategoryID:  1,
		Condition:   conditionUsed,
		Price:       9899900,
		Stock:       5,
		Specs: map[string]string{
			"Processor": "M5",
			"Memory":    "8GB",
			"Storage":   "256GB",
		},
	},
	{
		ID:          3,
		Name:        "40W USB-C Charger",
		Description: "Charges any USB-C MacBook.",
		CategoryID:  3,
		Condition:   conditionRefurbished,
		Price:       199900,
		Stock:       20,
		Specs: map[string]string{
			"Connector": "USB-C",
			"Output":    "40W",
		},
	},
	{
		ID:          4,
		Name:        "MacBook Pro case",
		Description: "A durable and stylish case to protect your MacBook from scratches and bumps.",
		CategoryID:  3,
		Condition:   conditionNew,
		Price:       299900,
		Stock:       15,
		Specs: map[string]string{
			"Material": "Polycarbonate",
			"Color":    "Black",
		},
		CompatibleModels: []string{"MacBook Pro 16-inch"},
	},
	{
		ID:          5,
		Name:        "Display Assembly for MacBook Pro 16-inch",
		Description: "Complete replacement display, tested with no dead pixels.",
		CategoryID:  2,
		Condition:   conditionNew,
		Price:       4999900,
		Stock:       10,
		Specs: map[string]string{
			"Type": "Liquid Retina XDR",
			"Size": "16-inch",
		},
		CompatibleModels: []string{"MacBook Pro 16-inch M2 Pro (2023)"},
	},
	{
		ID:          6,
		Name:        "Battery for MacBook Air 13-inch M2",
		Description: "Replacement battery. Installation screws included.",
		CategoryID:  2,
		Condition:   conditionNew,
		Price:       699900,
		Stock:       25,
		Specs: map[string]string{
			"Capacity":     "52.6 Wh",
			"Battery Type": "Lithium-polymer",
		},
		CompatibleModels: []string{"MacBook Air 13-inch M2 (2022)"},
	},
	{
		ID:          7,
		Name:        "USB-C to USB-C Cable (2 m)",
		Description: "Braided cable for charging and data transfer.",
		CategoryID:  3,
		Condition:   conditionNew,
		Price:       99900,
		Stock:       50,
	},
}
