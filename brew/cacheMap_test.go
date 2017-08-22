package brew_test

import (
	. "github.com/lgug2z/bfm/brew"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map", func() {
	var (
		infoOne, infoTwo, infoWithDependencies []Info
		cacheMap                               CacheMap
		cache                                  InfoCache
	)

	BeforeEach(func() {
		infoOne = []Info{Info{Name: "vim", FullName: "vim"}}
		infoTwo = []Info{Info{Name: "vim", FullName: "vim"}, Info{Name: "emacs", FullName: "emacs"}}
		infoWithDependencies = []Info{
			Info{
				Name:                    "vim",
				FullName:                "vim",
				Dependencies:            []string{"python"},
				OptionalDependencies:    []string{"node"},
				BuildDependencies:       []string{"ruby"},
				RecommendedDependencies: []string{"go"},
			},
			Info{Name: "ruby", FullName: "ruby"},
			Info{Name: "python", FullName: "python"},
			Info{Name: "node", FullName: "node"},
			Info{Name: "go", FullName: "go"},
		}
	})

	Describe("Initialising with a list of package names", func() {
		It("Should create an entry in the map for every package which has info in the cache", func() {
			cache = InfoCache(infoTwo)
			cacheMap = CacheMap{
				Cache: &cache,
				Map:   make(Map),
			}

			packages := []string{"brew 'vim'", "brew 'emacs'"}
			cacheMap.FromPackages(packages)
			vimEntry := Entry{Name: "vim", Info: Info{Name: "vim", FullName: "vim"}}
			emacsEntry := Entry{Name: "emacs", Info: Info{Name: "emacs", FullName: "emacs"}}

			Expect(cacheMap.Map).To(HaveKeyWithValue("vim", vimEntry))
			Expect(cacheMap.Map).To(HaveKeyWithValue("emacs", emacsEntry))
		})

		It("Should not create entries in the map for packages which have no info in the cache", func() {
			cache = InfoCache(infoOne)
			cacheMap = CacheMap{
				Cache: &cache,
				Map:   make(Map),
			}

			packages := []string{"brew 'vim'", "brew 'emacs'"}

			cacheMap.FromPackages(packages)

			vimEntry := Entry{Name: "vim", Info: Info{Name: "vim", FullName: "vim"}}
			Expect(cacheMap.Map).To(HaveKeyWithValue("vim", vimEntry))
			Expect(cacheMap.Map).ToNot(HaveKey("emacs"))
		})
	})

	Describe("Populated with packages and with an InfoCache", func() {
		It("Should update all entries in the map with whichever other packages that entry is required by", func() {
			cache = InfoCache(infoWithDependencies)
			cacheMap = CacheMap{
				Cache: &cache,
				Map:   make(Map),
			}

			packages := []string{"brew 'vim'", "brew 'python'"}
			cacheMap.FromPackages(packages)

			cacheMap.ResolveRequiredDependencyMap()

			Expect(cacheMap.Map["python"].RequiredBy).To(ContainElement("vim"))
		})
	})

	Describe("Initialised with an InfoCache", func() {
		It("Should add a new package without any of its dependencies", func() {
			cache = InfoCache(infoWithDependencies)
			cacheMap = CacheMap{
				Cache: &cache,
				Map:   make(Map),
			}

			cacheMap.Add(Entry{Name: "vim", RestartService: "true", Args: []string{"with-override-system-vim"}}, AddPackageOnly)

			Expect(cacheMap.Map).To(HaveKey("vim"))
			notHave := []string{"ruby", "python", "go", "node"}
			for _, e := range notHave {
				Expect(cacheMap.Map).ToNot(HaveKey(e))
			}
		})

		It("Should add a new package with its required dependencies", func() {
			cache = InfoCache(infoWithDependencies)
			cacheMap = CacheMap{
				Cache: &cache,
				Map:   make(Map),
			}

			cacheMap.Add(Entry{Name: "vim", RestartService: "true", Args: []string{"with-override-system-vim"}}, AddPackageAndRequired)

			have := []string{"vim", "python"}
			for _, e := range have {
				Expect(cacheMap.Map).To(HaveKey(e))
			}

			notHave := []string{"ruby", "go", "node"}
			for _, e := range notHave {
				Expect(cacheMap.Map).ToNot(HaveKey(e))
			}
		})

		It("Should add a new package with all of its required, recommended, optional and build dependencies", func() {
			cache = InfoCache(infoWithDependencies)
			cacheMap = CacheMap{
				Cache: &cache,
				Map:   make(Map),
			}

			cacheMap.Add(Entry{Name: "vim", RestartService: "true", Args: []string{"with-override-system-vim"}}, AddAll)

			have := []string{"vim", "python", "ruby", "go", "node"}
			for _, e := range have {
				Expect(cacheMap.Map).To(HaveKey(e))
			}
		})
	})

	Describe("Initialised with an InfoCache and populated with packages", func() {
		It("Should remove a package without removing any of its dependencies", func() {
			cache = InfoCache(infoWithDependencies)
			cacheMap = CacheMap{
				Cache: &cache,
				Map:   make(Map),
			}

			cacheMap.Add(Entry{Name: "vim", RestartService: "true", Args: []string{"with-override-system-vim"}}, AddAll)
			cacheMap.Remove("vim", RemovePackageOnly)

			have := []string{"python", "ruby", "go", "node"}
			for _, e := range have {
				Expect(cacheMap.Map).To(HaveKey(e))
			}

			notHave := []string{"vim"}
			for _, e := range notHave {
				Expect(cacheMap.Map).ToNot(HaveKey(e))
			}
		})

		It("Should remove a package and its required dependencies if they are not required by other packages", func() {
			cache = InfoCache(infoWithDependencies)
			cacheMap = CacheMap{
				Cache: &cache,
				Map:   make(Map),
			}

			cacheMap.Add(Entry{Name: "vim", RestartService: "true", Args: []string{"with-override-system-vim"}}, AddAll)
			cacheMap.Remove("vim", RemovePackageAndRequired)

			have := []string{"ruby", "go", "node"}
			for _, e := range have {
				Expect(cacheMap.Map).To(HaveKey(e))
			}

			notHave := []string{"vim", "python"}
			for _, e := range notHave {
				Expect(cacheMap.Map).ToNot(HaveKey(e))
			}
		})

		It("Should remove a package and all of its required, recommended, build and optional dependencies if they are not required by other packages", func() {
			cache = InfoCache(infoWithDependencies)
			cacheMap = CacheMap{
				Cache: &cache,
				Map:   make(Map),
			}

			cacheMap.Add(Entry{Name: "vim", RestartService: "true", Args: []string{"with-override-system-vim"}}, AddAll)
			cacheMap.Remove("vim", RemoveAll)

			notHave := []string{"vim", "python", "ruby", "go", "node"}
			for _, e := range notHave {
				Expect(cacheMap.Map).ToNot(HaveKey(e))
			}
		})
	})
})
