package prompt

import _ "embed"

//go:embed format_reference.md
var FormatReference string

//go:embed default_guidelines.md
var DefaultGuidelines string

//go:embed readme_template.md
var ReadmeTemplate string

//go:embed config_template.yaml
var ConfigTemplate string

//go:embed prompt_user_template.md
var PromptUserTemplate string
