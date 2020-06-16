package cmd

import (
	"fmt"

	"strings"
	"text/template"

	"bytes"

	"github.com/LGUG2Z/bfm/brew"
	"github.com/LGUG2Z/bfm/brewfile"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var checkFlags Flags

func init() {
	RootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolVarP(&checkFlags.Tap, "tap", "t", false, "check a tap")
	checkCmd.Flags().BoolVarP(&checkFlags.Brew, "brew", "b", false, "check a brew package")
	checkCmd.Flags().BoolVarP(&checkFlags.Cask, "cask", "c", false, "check a cask")
	checkCmd.Flags().BoolVarP(&checkFlags.Mas, "mas", "m", false, "check a mas app")
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if a dependency is in your Brewfile",
	Long:  DocsCheck,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var packages brewfile.Packages

		db, err := bolt.Open(boltPath, 0600, nil)
		if err != nil {
			errorExit(err)
		}

		cache := brew.Cache{DB: db}

		err = Check(args, &packages, cache, brewfilePath, checkFlags, level)
		errorExit(err)
	},
}

func Check(args []string, packages *brewfile.Packages, cache brew.Cache, brewfilePath string, flags Flags, level int) error {
	if !flagProvided(checkFlags) {
		return ErrNoPackageType("check")
	}

	if err := packages.FromBrewfile(brewfilePath); err != nil {
		return err
	}

	toCheck := args[0]
	packageType := getPackageType(checkFlags)

	cacheMap := brew.CacheMap{Cache: &cache, Map: make(brew.Map)}

	if err := cacheMap.FromPackages(packages.Brew); err != nil {
		return err
	}

	if err := cacheMap.ResolveDependencyMap(level); err != nil {
		return err
	}

	b, err := packages.Bytes()
	if err != nil {
		return err
	}

	if entryExists(string(b), packageType, toCheck) {
		switch packageType {
		case "brew":

			presence := `{{ .Name }} is present in the Brewfile.`
			dependencies := `{{- if (or .RequiredDependencies .RecommendedDependencies .OptionalDependencies .BuildDependencies) }}

{{- if .RequiredDependencies }}
Required dependencies: {{ StringsJoin .RequiredDependencies ", " }}
{{- end -}}

{{- if .RecommendedDependencies }}
Recommended dependencies: {{ StringsJoin .RecommendedDependencies ", " }}
{{- end -}}

{{- if .OptionalDependencies }}
Optional dependencies: {{ StringsJoin .OptionalDependencies ", " }}
{{- end -}}

{{- if .BuildDependencies }}
Build Dependencies: {{ StringsJoin .BuildDependencies ", " }}
{{- end -}}

{{- else }}
No required, recommended, optional or build dependencies.
{{- end -}}`

			dependencyOf := `{{- if (or .RequiredBy .RecommendedFor .OptionalFor .BuildOf) -}}
{{- if .RequiredBy }}
Required dependency of: {{ StringsJoin .RequiredBy ", " }}
{{- end -}}

{{- if .RecommendedFor }}
Recommended dependency of: {{ StringsJoin .RecommendedFor ", " }}
{{- end -}}

{{- if .OptionalFor }}
Optional dependency of: {{ StringsJoin .OptionalFor ", " }}
{{- end -}}

{{- if .BuildOf }}
Build dependency of: {{ StringsJoin .BuildOf ", " }}
{{- end -}}
{{- else }}
Not a required, recommended, optional or build dependency of any other package.
{{- end -}}`

			brew := cacheMap.Map[toCheck]

			var presenceBytes bytes.Buffer
			var dependenciesBytes bytes.Buffer
			var dependencyOfBytes bytes.Buffer

			funcMap := template.FuncMap{"StringsJoin": strings.Join}
			tmpl := template.Must(template.New("presence").Funcs(funcMap).Parse(presence))
			if err := tmpl.Execute(&presenceBytes, brew); err != nil {
				return err
			}

			tmpl = template.Must(template.New("dependencies").Funcs(funcMap).Parse(dependencies))
			if err := tmpl.Execute(&dependenciesBytes, brew); err != nil {
				return err
			}

			tmpl = template.Must(template.New("dependencyOf").Funcs(funcMap).Parse(dependencyOf))
			if err := tmpl.Execute(&dependencyOfBytes, brew); err != nil {
				return err
			}

			fmt.Println(presenceBytes.String())
			fmt.Println(dependenciesBytes.String())
			fmt.Println(dependencyOfBytes.String())
		default:
			fmt.Printf("%s is present in the Brewfile.\n", toCheck)
		}
	} else {
		fmt.Printf("%s is not present in the Brewfile.\n", toCheck)
	}

	return nil
}
