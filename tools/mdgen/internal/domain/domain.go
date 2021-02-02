package domain

import (
	"github.com/eujoy/openrecipe/tools/mdgen/internal/domain/translation"
)

// RecipeTemplateValues describes the recipe template values to be used.
type RecipeTemplateValues struct {
	Summary     Summary
	Details     Details
	Translation translation.Labels
}

// TocTemplateValues describes the toc template values to be used.
type TocTemplateValues struct {
	Languages []TocTemplateLanguageValues
}

// TocTemplateLanguageValues describes the toc template language details to be used.
type TocTemplateLanguageValues struct {
	Order      int
	Name       string
	Categories []TocTemplateCategoryValues
}

// TocTemplateCategoryValues describes the toc template category details to be used.
type TocTemplateCategoryValues struct {
	Name    string
	Recipes []TocTemplateRecipeValues
}

// TocTemplateRecipeValues describes the toc teemplate recipe details to be used.
type TocTemplateRecipeValues struct {
	Name     string
	FilePath string
}

// Recipe describes the whole recipe object.
type Recipe struct {
	RecipeData RecipeData `yaml:"recipe"`
}

// RecipeData describes the whole structure of a recipe.
type RecipeData struct {
	Summary Summary   `yaml:"summary"`
	Details []Details `yaml:"details"`
}

// Summary describes the general information of a recipe.
type Summary struct {
	ForPersons    int `yaml:"forPersons"`
	Pieces        int `yaml:"pieces"`
	Difficulty    int `yaml:"difficulty"`
	ExecutionTime int `yaml:"executionTime"`
	BakingTime    int `yaml:"bakingTime"`
}

// Details describes the specific details of a recipe.
type Details struct {
	Lang        string       `yaml:"lang"`
	Title       string       `yaml:"title"`
	RecipeParts []RecipePart `yaml:"recipeParts"`
}

// RecipePart describes the details of a specific recipe part.
type RecipePart struct {
	Summary        string       `yaml:"summary"`
	ExecutionTime  int          `yaml:"executionTime"`
	BakingTime     int          `yaml:"bakingTime"`
	Ingredients    []Ingredient `yaml:"ingredients"`
	ExecutionSteps []string     `yaml:"executionSteps"`
}

// Ingredient describes the specfics of a recipe incredient.
type Ingredient struct {
	Name     string  `yaml:"name"`
	Quantity float64 `yaml:"quantity"`
	Metric   string  `yaml:"metric"`
}
