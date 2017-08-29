package cmd

const (
	DocsRoot = `
Brewfile Manager (bfm) is a command line tool for managing
a dependency whitelist in the form of a Brewfile in a less
active and more comprehensible way.

In order to use bfm the following environment variables
first need to be exported in your shell rc file:

BFM_BREWFILE=/path/to/your/Brewfile
BFM_LEVEL=[required, recommended, optional, build]

When adding a new package to a Brewfile whitelist, it is
not uncommon for that package to install other packages
which are required dependencies, and depending on the
arguments given, recommended and optional dependencies too.

These additional dependencies, if not also added to the
Brewfile, get marked for removal by the 'brew bundle cleanup'
command. So you add those dependencies to your Brewfile too,
and eventually you end up with a Brewfile you struggle to
make sense of for all of these additional dependencies that
have been added for the cleanup command to remain useful.

By setting a BFM_LEVEL, when performing any operation with
bfm, the level of dependencies to operate on can be kept
consistent.

When you add and remove packages using bfm, depending on the
level chosen, all of the required, recommended, optional or
build dependencies belonging to that package can also be
added, complete with annotations, or removed from the Brewfile
at the same time.

`

	DocsAdd = `
Adds the dependency given as an argument to the Brewfile.

This command will modify your Brewfile without creating a
backup. Consider running the command with the --dry-run flag
if using bfm for the first time.

The type must be specified using the appropriate flag.

Taps must conform to the format <user/repo>.

Brew packages can have arguments specified using the --args
flag (multiple arguments can be separated by using a comma),
and can specify service restart behaviour ('always' to
restart every time bundle is run, 'changed' to restart only
when updated or changed) with the --restart-service flag.

MAS apps must specify an id using the --mas-id flag which
can be found by running 'mas search <app>'.

Examples:

bfm add -t homebrew/dupes
bfm add -b vim --args HEAD,with-override-system-vi
bfm add -b crisidev/chunkwm/chunkwm --restart-service changed
bfm add -c macvim
bfm add -m Xcode -i 497799835

`
	DocsCheck = `
Checks for the presence of the argument as an entry in the
Brewfile.

If the arguments corresponds to a brew entry in the Brewfile,
the check command will provide information about both any
dependencies it has, or any other entries of which it is
itself a dependency.

The type must be specified using the appropriate flag.

Examples:

bfm check -t homebrew/dupes
bfm check -b vim
bfm check -c macvim
bfm check -m Xcode

`
	DocsClean = `
Cleans up your Brewfile, removing all comments and sorting
all dependencies into alphabetised groups with the order tap
-> brew (primary) -> brew (dependent) -> cask -> mas.

This command will modify your Brewfile without creating a
backup. Consider running the command with the --dry-run flag
if using bfm for the first time.

Examples:

bfm clean
bfm clean --dry-run

`
	DocsRefresh = `
Refreshes the bfm cache stored at '$HOME/.bfm.bolt'.

This command will get information about all brews and casks
it is possible for you to install given the repositories
that you have tapped, and store it in a Bolt DB file in the
home folder.

This command should be run after adding a new tap.

Examples:

bfm refresh

`
	DocsRemove = `
Removes from the Brewfile the entry corresponding to the
argument.

This command will modify your Brewfile without creating a
backup.  Consider running the command with the --dry-run
flag if using bfm for the first time.

The type must be specified using the appropriate flag.

Examples:

bfm remove -t homebrew/dupes
bfm remove -b vim
bfm remove -c macvim
bfm remove -m Xcode

`
)
