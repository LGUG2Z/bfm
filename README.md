# Brewfile Manager
Brewfile Manager (bfm) is a command line tool for managing
a dependency whitelist in the form of a Brewfile in a less
active and more comprehensible way.

## Requirements
* [Homebrew](https://github.com/homebrew/brew)
* [Homebrew Bundle](https://github.com/Homebrew/homebrew-bundle)
* [Go](https://github.com/golang/go)

## Install
The latest version of bfm can be installed using `go get`.

```
go get -u github.com/LGUG2Z/bfm
```

Make sure `$GOPATH` is set correctly that and that `$GOPATH/bin` is in your `$PATH`.

The `bfm` executable will be installed under the `$GOPATH/bin` directory.

In order to use bfm the following environment variables
first need to be exported in your shell rc file:

```
BFM_BREWFILE=/path/to/your/Brewfile
BFM_LEVEL=[required, recommended, optional, build]
```

## Overview
When adding a new package to a Brewfile whitelist, it is
not uncommon for that package to install other packages
which are required dependencies, and depending on the
arguments given, recommended and optional dependencies too.

These additional dependencies, if not also added to the
Brewfile, get marked for removal by the `brew bundle cleanup`
command. So you add those dependencies to your Brewfile too,
and eventually you end up with a Brewfile you struggle to
make sense of for all of these additional dependencies that
have been added for the cleanup command to remain useful.

Take as an example the neovim package:

```
neovim: stable 0.2.0 (bottled), HEAD
Ambitious Vim-fork focused on extensibility and agility
https://neovim.io/
/usr/local/Cellar/neovim/0.2.0_1 (1,352 files, 17.2MB) *
  Poured from bottle on 2017-08-04 at 15:25:02
From: https://github.com/Homebrew/homebrew-core/blob/master/Formula/neovim.rb
==> Dependencies
Build: luajit ✘, cmake ✔, lua@5.1 ✘, pkg-config ✔
Required: gettext ✔, jemalloc ✔, libtermkey ✔, libuv ✔, libvterm ✔, msgpack ✔, unibilium ✔
```

If only neovim is whitelisted in the Brewfile, then the
`brew bundle cleanup` command will suggest the removal of the
required dependencies that are not directly whitelisted in the
Brewfile but necessary to be able to use neovim. These can be
added manually to the Brewfile, but this approach very quickly
leads to a Brewfile full of tersely named packages which may
or may not all be required as packages are added and removed
from the Brewfile over a period of time.

Ultimately, a whitelist where you can't look at the
packages included and know exactly what is whitelisted and
why, is not a very useful whitelist at all.

By setting a `BFM_LEVEL`, when performing any operation with
bfm, the level of dependencies to operate on can be kept
consistent.

When you add and remove packages using bfm, depending on the
level chosen, all of the required, recommended, optional or
build dependencies belonging to that package can also be
added, complete with annotations, or removed from the Brewfile
at the same time.

### Use Cases
I have developed this tool primarily for my own personal use.
I am a consultant and can find myself working on multiple
different software projects throughout the course of a year,
often with quite a short turn around time between them. While
it would be nice to have a fresh install of macOS every time I
join a new project or to use Vagrant for every project to ensure
a clean development environment, this is not always practical or
possible.

By having a well managed and comprehensible Brewfile whitelist,
firstly, I can run `brew bundle cleanup --global --force` at the
end of a project and be sure that I have no lingering packages
that are unneeded going forward into a new project, and secondly,
I can look at my whitelist at any time and have great certainty
about what every brew installed package on my system is being
used for.

## Usage
### Basics
The main commands of bfm are `add`, `remove`, `clean`, `check` and `refresh`. Information
for any of these commands can be found by running `bfm [cmd] --help`.

#### Add, Remove, Clean
`add`, `remove` and `clean` are destructive commands that will permanently modify your Brewfile without creating a backup.

It is recommended that if you are using bfm for the first time, you run these commands with the `--dry-run` flag.

When adding to the Brewfile, a flag must be used to specify what is being added:

```
bfm add --tap homebrew/dupes
bfm add --brew vim --args HEAD,with-override-system-vi
bfm add --brew crisidev/chunkwm/chunkwm --restart-service changed
bfm add --cask macvim
bfm add --mas Xcode --mas-id 497799835
```

Additional arguments for brew dependencies can be specified with the `--args` flag and service restart behaviour (`always`, `changed`) can be specified with the `--restart-service` flag.

The same flags must also be used with the `remove` and `check` commands.

```
bfm remove --tap homebrew/dupes
bfm check --brew vim
bfm remove --cask macvim
bfm check --mas Xcode
```

The `clean` command will organise your Brewfile and sort it into sections in
the following order: taps -> primary brews -> dependent brews -> -> casks -> mas apps.

Before: 
```
mas 'Xcode', id: 497799835
cask 'macvim'
brew 'neovim'
brew 'gettext'
brew 'jemalloc'
brew 'libtermkey'
brew 'libuv'
brew 'libvterm'
brew 'msgpack'
brew 'unibilium'
# brew 'emacs'
cask 'firefox'
tap 'caskroom/cask'
tap 'homebrew/dupes'
# cask 'google-chrome'
brew 'cmus'
```

After:
```
tap 'caskroom/cask'
tap 'homebrew/dupes'
 
brew 'cmus'
brew 'neovim'

brew 'gettext' # [required by: neovim]
brew 'jemalloc' # [required by: neovim]
brew 'libtermkey' # [required by: neovim]
brew 'libuv' # [required by: neovim]
brew 'libvterm' # [required by: neovim]
brew 'msgpack' # [required by: neovim]
brew 'unibilium' # [required by: neovim]

cask 'firefox'
cask 'macvim'
 
mas 'Xcode', id: 497799835
```

By splitting up the brews into primary and dependent sections helps to separate the signal
from the noise. Essentially, you should have a clear understanding of what every package
in the primary brews section does and why it is there. If you don't, it is worth rethinking
its place in your whitelist.

The Brewfile is automatically cleaned after every `add` and `remove` operation.

#### Check
The `check` command is a quick way to get feedback about the presence of a package
in the Brewfile and, if it is a brew package, to get feedback about what its
dependencies are and if the package itself is a dependency of another package.

```
❯ bfm check -b neovim
neovim is present in the Brewfile.

Required dependencies: gettext, jemalloc, libtermkey, libuv, libvterm, msgpack, unibilium
Build Dependencies: luajit, cmake, lua@5.1, pkg-config

Not a required, recommended, optional or build dependency of any other package.
```

```
❯ bfm check -b gettext
gettext is present in the Brewfile.

No required, recommended, optional or build dependencies.

Required dependency of: glib, gnupg, libmp3splt, neovim, weechat
```

#### Refresh
The `refresh` command will get information about all installable brews and casks
given the repositories that have been tapped on the system, and stores it in a
BoltDB file in the home folder.

This command should be run after adding a new tap and periodically to stay up
to date with new information added after every `brew update`.

