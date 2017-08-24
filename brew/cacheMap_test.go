package brew_test

import (
	. "github.com/lgug2z/bfm/brew"

	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map", func() {
	var (
		cacheMap             CacheMap
		cache                Cache
		dbFile               = fmt.Sprintf("%s/src/github.com/lgug2z/bfm/testData/testDB.bolt", os.Getenv("GOPATH"))
		infoWithDependencies = []Info{
			Info{FullName: "vim",
				Dependencies:            []string{"python"},
				OptionalDependencies:    []string{"node"},
				BuildDependencies:       []string{"go"},
				RecommendedDependencies: []string{"ruby"}},
			Info{FullName: "python"},
			Info{FullName: "ruby"},
			Info{FullName: "go"},
			Info{FullName: "node"},
		}
	)

	BeforeEach(func() {
		cacheMap = CacheMap{
			Cache: &cache,
			Map:   make(Map),
		}
	})

	Describe("Initialising with a list of package names", func() {
		It("Should create an entry in the map for every package which has info in the cache", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			Expect(testDB.AddTestBrews("vim", "emacs")).To(Succeed())
			packages := []string{"brew 'vim'", "brew 'emacs'"}

			cacheMap.FromPackages(packages)
			vimEntry := Entry{Name: "vim", Info: Info{FullName: "vim"}}
			emacsEntry := Entry{Name: "emacs", Info: Info{FullName: "emacs"}}

			Expect(cacheMap.Map).To(HaveKeyWithValue("vim", vimEntry))
			Expect(cacheMap.Map).To(HaveKeyWithValue("emacs", emacsEntry))
		})

		It("Should not create entries in the map for packages which have no info in the cache", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			Expect(testDB.AddTestBrews("vim")).To(Succeed())
			packages := []string{"brew 'vim'", "brew 'emacs'"}
			cacheMap.FromPackages(packages)

			vimEntry := Entry{Name: "vim", Info: Info{FullName: "vim"}}
			Expect(cacheMap.Map).To(HaveKeyWithValue("vim", vimEntry))
			Expect(cacheMap.Map).ToNot(HaveKey("emacs"))
		})
	})

	Describe("Populated with packages and with a Cache", func() {
		It("Should update all entries in the map with whichever other packages that entry is required by", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			Expect(testDB.AddTestBrewsFromInfo(
				Info{FullName: "vim", Dependencies: []string{"python"}},
				Info{FullName: "python"},
			)).To(Succeed())

			packages := []string{"brew 'vim'", "brew 'python'"}
			cacheMap.FromPackages(packages)

			cacheMap.ResolveRequiredDependencyMap()

			Expect(cacheMap.Map["python"].RequiredBy).To(ContainElement("vim"))
		})
	})

	Describe("With a functioning bolt db", func() {
		It("Should add a new package without any of its dependencies", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			Expect(testDB.AddTestBrewsFromInfo(infoWithDependencies...)).To(Succeed())

			cacheMap.Add(Entry{Name: "vim", RestartService: "true", Args: []string{"with-override-system-vim"}}, AddPackageOnly)

			Expect(cacheMap.Map).To(HaveKey("vim"))
			notHave := []string{"ruby", "python", "go", "node"}
			for _, e := range notHave {
				Expect(cacheMap.Map).ToNot(HaveKey(e))
			}
		})

		It("Should add a new package with its required dependencies", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			Expect(testDB.AddTestBrewsFromInfo(infoWithDependencies...)).To(Succeed())

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
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			Expect(testDB.AddTestBrewsFromInfo(infoWithDependencies...)).To(Succeed())

			cacheMap.Add(Entry{Name: "vim", RestartService: "true", Args: []string{"with-override-system-vim"}}, AddAll)

			have := []string{"vim", "python", "ruby", "go", "node"}
			for _, e := range have {
				Expect(cacheMap.Map).To(HaveKey(e))
			}
		})
	})

	Describe("Initialised with a Cache and populated with packages", func() {
		It("Should remove a package without removing any of its dependencies", func() {
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			Expect(testDB.AddTestBrewsFromInfo(infoWithDependencies...)).To(Succeed())

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
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			Expect(testDB.AddTestBrewsFromInfo(infoWithDependencies...)).To(Succeed())

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
			testDB, err := NewTestDB(dbFile)
			Expect(err).ToNot(HaveOccurred())
			defer testDB.Close()
			cache.DB = testDB.DB

			Expect(testDB.AddTestBrewsFromInfo(infoWithDependencies...)).To(Succeed())

			cacheMap.Add(Entry{Name: "vim", RestartService: "true", Args: []string{"with-override-system-vim"}}, AddAll)
			cacheMap.Remove("vim", RemoveAll)

			notHave := []string{"vim", "python", "ruby", "go", "node"}
			for _, e := range notHave {
				Expect(cacheMap.Map).ToNot(HaveKey(e))
			}
		})
	})
})
