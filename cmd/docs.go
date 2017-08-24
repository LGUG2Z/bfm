package cmd

const (
	DocsAdd = `
Adds the dependency given as an argument to the Brewfile.

This command will modify your Brewfile without creating
a backup. Consider running the command with the --dry-run
flag if using bfm for the first time.

The type must be specified using the appropriate flag.

Taps must conform to the format <user/repo>.

Brew packages can have arguments specified using the --arg
flag (multiple arguments can be separated by using a comma),
and can specify service restart behaviour ('always' to restart
every time bundle is run, 'changed' to restart only when updated
or changed) with the --restart-service flag.

The --required flag will add a brew entry and its required
dependencies.

The --all flag will add a brew entry along with all of its
required, recommended, optional and build dependencies.

MAS apps must specify an id using the --mas-id flag which
can be found by running 'mas search <app>'.

Examples:

bfm add -t homebrew/dupes
bfm add -b vim -a HEAD,with-override-system-vi
bfm add -b crisidev/chunkwm/chunkwm -r changed
bfm add -c macvim
bfm add -m Xcode -i 497799835

`
	// TODO: Update this depending on whether templating will be included
	DocsCheck = `
Checks for the presence of the argument as an entry
in the Brewfile.

The type must be specified using the appropriate flag.

Examples:

bfm check -t homebrew/dupes
bfm check -b vim
bfm check -c macvim
bfm check -m Xcode

`
	DocsClean = `
Cleans up your Brewfile, removing all comments and
sorting all dependencies into alphabetised groups
with the order tap -> brew -> cask -> mas.

This command will modify your Brewfile without creating
a backup. Consider running the command with the --dry-run
flag if using bfm for the first time.

Examples:

bfm clean
bfm clean --dry-run

`
	DocsRefresh = `
Refreshes the bfm cache stored at ''$HOME/.bfm.bolt'.

This command will get information about all brews and casks
it is possible for you to install given the repositories that
you have tapped, and store it in a Bolt DB file in the home
folder.

This command should be run after adding a new tap.

Examples:

bfm refresh

`
	DocsRemove = `
Removes from the Brewfile the entry corresponding to the argument.

This command will modify your Brewfile without creating a backup.
Consider running the command with the --dry-run  flag if using
bfm for the first time.

The default behaviour is to remove only the entry corresponding
to the given argument.

The --required flag will remove a brew entry along with all of
the required dependencies of that entry which are no longer
required by any other brew entry.

The --all flag will remove a brew entry along with all of the
required, recommended, optional and build dependencies of that
entry which are no longer required by any other brew entry.

The type must be specified using the appropriate flag.

Examples:

bfm remove -t homebrew/dupes
bfm remove -b neovim --required
bfm remove -b vim --all
bfm remove -c macvim
bfm remove -m Xcode

`
)
