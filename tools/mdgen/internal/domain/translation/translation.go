package translation

// Details describes the translations for all the required fields.
type Details map[string]Language

// Language describes the details of the translations for a language.
type Language struct {
	Labels Labels `yaml:"labels"`
	Toc    Toc    `yaml:"toc"`
}

// Labels describes the details of the translation labels.
type Labels struct {
	GeneralDetailsTitle string `yaml:"generalDetailsTitle"`
	RecipePartsTitle    string `yaml:"recipePartsTitle"`
	ForPersons          string `yaml:"forPersons"`
	Pieces              string `yaml:"pieces"`
	Difficulty          string `yaml:"difficulty"`
	ExecutionTime       string `yaml:"executionTime"`
	BakingTime          string `yaml:"bakingTime"`
	Minutes             string `yaml:"minutes"`
	LabelsTextAlign     string `yaml:"labelsTextAlign"`
	TextAlign           string `yaml:"textAlign"`
}

// Toc describes the details of the translation toc.
type Toc struct {
	Order      int               `yaml:"order"`
	Language   string            `yaml:"language"`
	Categories map[string]string `yaml:"categories"`
	Recipes    map[string]string `yaml:"recipes"`
}
