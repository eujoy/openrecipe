package template

// Schema describes the template schemas for all the required parts.
type Schema struct {
	Recipe string
	Toc    string
}

var recipeMarkdown = `# {{ .Details.Title }}

## {{ .Translation.GeneralDetailsTitle }}

| {{ .Translation.ForPersons }} | {{ .Translation.Pieces }} | {{ .Translation.Difficulty }} | {{ .Translation.ExecutionTime }} | {{ .Translation.BakingTime }} |
| {{ .Translation.LabelsTextAlign }} | {{ .Translation.LabelsTextAlign }} | {{ .Translation.LabelsTextAlign }} | {{ .Translation.LabelsTextAlign }} | {{ .Translation.LabelsTextAlign }} |
| {{- if gt .Summary.ForPersons 0 }}{{ .Summary.ForPersons }}{{- else }}~{{- end }} | {{- if gt .Summary.Pieces 0 }}{{ .Summary.Pieces }}{{- else }}~{{- end }} | {{ .Summary.Difficulty }} | {{ .Summary.ExecutionTime }} {{ .Translation.Minutes }} | {{ .Summary.BakingTime }} {{ .Translation.Minutes }} |

## {{ .Translation.RecipePartsTitle }}

| | |
| {{ .Translation.TextAlign }} | {{ .Translation.TextAlign }} |
{{- range .Details.RecipeParts }}
| **{{ .Summary }}** <br/> <ul> {{- range .ExecutionSteps }} <li> {{ . }} </li> {{- end }} </ul> | {{ if gt (len .Ingredients) 0 }}<ul>{{ range .Ingredients}}<li>{{ .Name }}: {{ .Quantity }}{{ .Metric }}</li> {{ end }} </ul> {{ end }}
{{- end }} |`

var tocMarkdown = `{{- range .Languages }}## {{ .Name }}
{{ range .Categories }}
#### {{ .Name }}
{{ range .Recipes }}
* [{{ .Name }}]({{ .FilePath }})
{{ end }}
{{- end }}
----

{{ end }}`

// Prepare and reutn the existing template schemas.
func Prepare() Schema {
	return Schema{
		Recipe: recipeMarkdown,
		Toc:    tocMarkdown,
	}
}
