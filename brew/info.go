package brew

type Info struct {
	Name     string   `json:"name"`
	FullName string   `json:"full_name"`
	Desc     string   `json:"desc"`
	Homepage string   `json:"homepage"`
	Oldname  string   `json:"oldname"`
	Aliases  []string `json:"aliases"`
	Versions struct {
		Stable string `json:"stable"`
		Bottle bool   `json:"bottle"`
		Devel  string `json:"devel"`
		Head   string `json:"head"`
	} `json:"versions"`
	Revision      int `json:"revision"`
	VersionScheme int `json:"version_scheme"`
	Installed     []struct {
		Version             string   `json:"version"`
		UsedOptions         []string `json:"used_options"`
		BuiltAsBottle       bool     `json:"built_as_bottle"`
		PouredFromBottle    bool     `json:"poured_from_bottle"`
		RuntimeDependencies []struct {
			FullName string `json:"full_name"`
			Version  string `json:"version"`
		} `json:"runtime_dependencies"`
		InstalledAsDependency bool `json:"installed_as_dependency"`
		InstalledOnRequest    bool `json:"installed_on_request"`
	} `json:"installed"`
	LinkedKeg               string   `json:"linked_keg"`
	Pinned                  bool     `json:"pinned"`
	Outdated                bool     `json:"outdated"`
	KegOnly                 bool     `json:"keg_only"`
	Dependencies            []string `json:"dependencies"`
	RecommendedDependencies []string `json:"recommended_dependencies"`
	OptionalDependencies    []string `json:"optional_dependencies"`
	BuildDependencies       []string `json:"build_dependencies"`
	ConflictsWith           []string `json:"conflicts_with"`
	Caveats                 string   `json:"caveats"`
	Requirements            []struct {
		Name           string `json:"name"`
		DefaultFormula string `json:"default_formula"`
		Cask           string `json:"cask"`
		Download       string `json:"download"`
	} `json:"requirements"`
	Options []struct {
		Option      string `json:"option"`
		Description string `json:"description"`
	} `json:"options"`
	Bottle struct {
		Stable struct {
			Rebuild int    `json:"rebuild"`
			Cellar  string `json:"cellar"`
			Prefix  string `json:"prefix"`
			RootURL string `json:"root_url"`
			Files   struct {
				Sierra struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"sierra"`
				ElCapitan struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"el_capitan"`
				Yosemite struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"yosemite"`
				Mavericks struct {
					URL    string `json:"url"`
					Sha256 string `json:"sha256"`
				} `json:"mavericks"`
			} `json:"files"`
		} `json:"stable"`
	} `json:"bottle"`
}
