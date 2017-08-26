package brew_test

import (
	. "github.com/lgug2z/bfm/brew"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Entry", func() {
	var (
		info   Info
		actual Entry
	)

	BeforeEach(func() {
		info = Info{
			FullName:                "a",
			Dependencies:            []string{"b", "c", "d", "e", "f"},
			OptionalDependencies:    []string{"c"},
			BuildDependencies:       []string{"d"},
			RecommendedDependencies: []string{"e", "f"},
		}
	})

	Describe("With a valid Info object returned from the Cache", func() {
		It("Populates all the fields of an Entry", func() {
			expected := Entry{
				Name:                    "a",
				RequiredDependencies:    []string{"b"},
				OptionalDependencies:    []string{"c"},
				BuildDependencies:       []string{"d"},
				RecommendedDependencies: []string{"e", "f"},
			}

			actual.FromInfo(info)
			Expect(actual).To(Equal(expected))
		})
	})

	It("Determines all the dependencies of an Entry from its Info", func() {
		expected := Entry{
			Name:                    "a",
			RequiredDependencies:    []string{"b"},
			OptionalDependencies:    []string{"c"},
			BuildDependencies:       []string{"d"},
			RecommendedDependencies: []string{"e", "f"},
		}

		actual = Entry{
			Name: "a",
		}

		actual.DetermineDependencies(info)
		Expect(actual).To(Equal(expected))
	})

	Describe("With a populated Entry", func() {
		It("Formats the Entry with a package name as a Brewfile-compliant line entry", func() {
			expected := `brew 'vim'`
			entry := Entry{Name: "vim"}

			actual, err := entry.Format()
			Expect(err).To(BeNil())

			Expect(actual).To(Equal(expected))
		})

		It("Formats the Entry with a package name, args and restart-service as a Brewfile-compliant line entry", func() {
			expected := `brew 'vim', args: ['HEAD'], restart_service: :changed`
			entry := Entry{Name: "vim", RestartService: ":changed", Args: []string{"HEAD"}}

			actual, err := entry.Format()
			Expect(err).To(BeNil())

			Expect(actual).To(Equal(expected))
		})

		It("Formats the Entry with a comment specifying which other packages it is required by", func() {
			expected := `brew 'vim' # [required by: developers]`
			entry := Entry{Name: "vim", RequiredBy: []string{"developers"}}

			actual, err := entry.Format()
			Expect(err).To(BeNil())

			Expect(actual).To(Equal(expected))
		})
	})
})
