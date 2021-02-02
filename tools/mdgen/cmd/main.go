package main

//
// @todo
// 1. Add tags on recipes.
// 2. Add photos in recipes (i.e. carousel).
//

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/eujoy/openrecipe/tools/mdgen/internal/domain"
	templateSchema "github.com/eujoy/openrecipe/tools/mdgen/internal/domain/template"
	translationDomain "github.com/eujoy/openrecipe/tools/mdgen/internal/domain/translation"
	"github.com/eujoy/openrecipe/tools/mdgen/internal/service/translation"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type tocFileDetails struct {
	name     string
	filePath string
}

type initSt struct {
	size    int64
	modTime time.Time
}

func main() {
	xLationsService := translation.New("translations.yaml")
	xLations := xLationsService.Prepare()

	templates := templateSchema.Prepare()

	var app = cli.NewApp()
	info(app)

	app.Commands = []*cli.Command{
		&cli.Command{
			Name:    "watch",
			Aliases: []string{"w"},
			Usage:   "Runs indefinetely and watches for changes on yaml files. Once identified, then the respective markdown files are rebuilt.",
			Flags:   []cli.Flag{},
			Action: func(c *cli.Context) error {
				doneChan := make(chan bool)

				go watchAndGenerateFromModified(xLations, templates.Recipe, doneChan)

				<-doneChan

				return nil
			},
		},
		&cli.Command{
			Name:    "generate",
			Aliases: []string{"w"},
			Usage:   "Generate the markdown files for all the existing yaml files.",
			Flags:   []cli.Flag{},
			Action: func(c *cli.Context) error {
				generateRecipes(xLations, templates.Recipe)

				return nil
			},
		},
		&cli.Command{
			Name:    "toc",
			Aliases: []string{"t"},
			Usage:   "Generate the table of contents based on the markdown files.",
			Flags:   []cli.Flag{},
			Action: func(c *cli.Context) error {
				generateToc(xLations, templates.Toc)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error running the service : %v\n", err)
		os.Exit(1)
	}
}

// info sets up the information of the tool.
func info(app *cli.App) {
	app.Authors = []*cli.Author{
		{
			Name:  "Angelos Giannis",
			Email: "angelos.giannis@gmail.com",
		},
	}
	app.Name = "recipes-generator"
	app.Usage = "Use this tool to automatically convert yaml files to markdown."
	app.Version = "v0.0.1"
}

func watchAndGenerateFromModified(xLations translationDomain.Details, templateSchema string, doneChan chan bool) {
	fileList, _ := walkMatch("./pages", "*.yaml")

	fmt.Println("Watching files for changes...")

	for {
		modifiedFile, err := watchFileList(fileList)
		if err != nil {
			log.Fatalf("Failed to watch list of files: %v", err)
		}

		if modifiedFile != "" {
			fmt.Printf("File modified : %v\n", modifiedFile)

			fileLocSplit := strings.Split(modifiedFile, "_data/")

			yamlFile, err := ioutil.ReadFile(modifiedFile)
			if err != nil {
				log.Printf("yamlFile.Get err   #%v ", err)
			}

			var recipe domain.Recipe
			err = yaml.Unmarshal(yamlFile, &recipe)
			if err != nil {
				log.Fatalf("Unmarshal: %v", err)
			}

			markdownContent, err := writeTemplate(xLations, templateSchema, recipe)
			if err != nil {
				log.Fatalf("Failed to create content from template: %v", err)
			}

			for lang, mkContent := range markdownContent {
				fileDir := fmt.Sprintf("%v%v", fileLocSplit[0], lang)
				if _, err := os.Stat(fileDir); os.IsNotExist(err) {
					os.Mkdir(fileDir, 0755)
				}

				markdownFile := strings.Replace(fileLocSplit[1], ".yaml", ".md", 1)
				err = ioutil.WriteFile(fmt.Sprintf("%v/%v", fileDir, markdownFile), []byte(mkContent), 0644)
				if err != nil {
					log.Fatalf("Failed to write data to markdown file: %v", err)
				}
			}
		}
	}
}

func generateRecipes(xLations translationDomain.Details, recipeTemplate string) {
	fileList, _ := walkMatch("./pages", "*.yaml")

	for _, fl := range fileList {
		fmt.Printf("Parsing file : %v\n", fl)

		fileLocSplit := strings.Split(fl, "_data/")

		yamlFile, err := ioutil.ReadFile(fl)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}

		var recipe domain.Recipe
		err = yaml.Unmarshal(yamlFile, &recipe)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}

		markdownContent, err := writeTemplate(xLations, recipeTemplate, recipe)
		if err != nil {
			log.Fatalf("Failed to create content from template: %v", err)
		}

		for lang, mkContent := range markdownContent {
			fileDir := fmt.Sprintf("%v%v", fileLocSplit[0], lang)
			if _, err := os.Stat(fileDir); os.IsNotExist(err) {
				os.Mkdir(fileDir, 0755)
			}

			markdownFile := strings.Replace(fileLocSplit[1], ".yaml", ".md", 1)
			err = ioutil.WriteFile(fmt.Sprintf("%v/%v", fileDir, markdownFile), []byte(mkContent), 0644)
			if err != nil {
				log.Fatalf("Failed to write data to markdown file: %v", err)
			}
		}
	}
}

func generateToc(xLations translationDomain.Details, tocTemplate string) {
	fileList, _ := walkMatch("./pages", "*.md")

	tocMap := make(map[string]map[string][]tocFileDetails)

	for _, fl := range fileList {
		if !strings.Contains(fl, "recipe") {
			continue
		}

		fileLocSplit := strings.Split(fl, "/")
		if _, ok := tocMap[fileLocSplit[3]]; !ok {
			tocMap[fileLocSplit[3]] = make(map[string][]tocFileDetails)
		}
		tocMap[fileLocSplit[3]][fileLocSplit[2]] = append(
			tocMap[fileLocSplit[3]][fileLocSplit[2]],
			tocFileDetails{
				name:     strings.Replace(fileLocSplit[4], ".md", "", 1),
				filePath: strings.Replace(fl, "pages/", "", 1),
			},
		)
	}

	var tmplLangData []domain.TocTemplateLanguageValues
	for lang, categoryDetails := range tocMap {
		var tmplCategoryData []domain.TocTemplateCategoryValues
		for category, recipeDetails := range categoryDetails {
			var tmplRecipeData []domain.TocTemplateRecipeValues
			for _, recipe := range recipeDetails {
				recipeData := domain.TocTemplateRecipeValues{
					Name:     xLations[lang].Toc.Recipes[recipe.name],
					FilePath: recipe.filePath,
				}

				tmplRecipeData = append(tmplRecipeData, recipeData)
			}

			categoryData := domain.TocTemplateCategoryValues{
				Name:    xLations[lang].Toc.Categories[category],
				Recipes: tmplRecipeData,
			}

			tmplCategoryData = append(tmplCategoryData, categoryData)
		}

		langData := domain.TocTemplateLanguageValues{
			Order:      xLations[lang].Toc.Order,
			Name:       xLations[lang].Toc.Language,
			Categories: tmplCategoryData,
		}

		tmplLangData = append(tmplLangData, langData)
	}

	tocData := domain.TocTemplateValues{
		Languages: tmplLangData,
	}

	tocContent := writeTpcTemplate(tocTemplate, tocData)
	err := ioutil.WriteFile("pages/README.md", []byte(tocContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write data to markdown file: %v", err)
	}
}

func watchFileList(fileList []string) (string, error) {
	initialValues := make(map[string]initSt)

	for _, fl := range fileList {
		initialState, err := os.Stat(fl)
		if err != nil {
			return "", err
		}

		initialValues[fl] = initSt{
			size:    initialState.Size(),
			modTime: initialState.ModTime(),
		}
	}

	for {
		for _, fl := range fileList {
			stat, err := os.Stat(fl)
			if err != nil {
				return "", err
			}

			if stat.Size() != initialValues[fl].size || stat.ModTime() != initialValues[fl].modTime {
				return fl, nil
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func walkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}

func writeTemplate(xLations translationDomain.Details, templateSchema string, recipe domain.Recipe) (map[string]string, error) {
	recipesContents := make(map[string]string)
	for _, det := range recipe.RecipeData.Details {
		templateValues := domain.RecipeTemplateValues{
			Summary:     recipe.RecipeData.Summary,
			Details:     det,
			Translation: xLations[det.Lang].Labels,
		}

		t, err := template.New("outputTemplate").Parse(templateSchema)
		if err != nil {
			fmt.Printf("Failed to prepare template with error : %v\n", err)
			return map[string]string{}, err
		}

		var tpl bytes.Buffer
		err = t.Execute(&tpl, templateValues)
		if err != nil {
			fmt.Printf("Failed to print text with error : %v\n", err)
			return map[string]string{}, err
		}

		recipesContents[det.Lang] = tpl.String()
	}

	return recipesContents, nil
}

func writeTpcTemplate(templateSchema string, tocTmplData domain.TocTemplateValues) string {
	sortTocData(&tocTmplData)

	t, err := template.New("outputTemplate").Parse(templateSchema)
	if err != nil {
		fmt.Printf("Failed to prepare template with error : %v\n", err)
		return ""
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, tocTmplData)
	if err != nil {
		fmt.Printf("Failed to print text with error : %v\n", err)
		return ""
	}

	return tpl.String()
}

func sortTocData(tocTmplData *domain.TocTemplateValues) {
	sort.Slice(tocTmplData.Languages, func(i, j int) bool {
		return tocTmplData.Languages[i].Order < tocTmplData.Languages[j].Order
	})

	for _, lang := range tocTmplData.Languages {
		for _, cat := range lang.Categories {
			sort.Slice(cat.Recipes, func(i, j int) bool {
				return cat.Recipes[i].Name < cat.Recipes[j].Name
			})
		}

		sort.Slice(lang.Categories, func(i, j int) bool {
			return lang.Categories[i].Name < lang.Categories[j].Name
		})
	}
}
