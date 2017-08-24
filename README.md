# Brewfile Manager
Keep your Brewfile tidy and organised.

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

Once installed, add `export BFM_BREWFILE=/path/to/your/Brewfile` to the rc file of the shell you are using.

## Usage
### Overview
The main commands of bfm are `refresh` `add`, `remove`, `check` and `clean`.

`add`, `remove` and `clean` are destructive commmands that will permanently modify your Brewfile without creating a backup.

It is recommended that if you are using bfm for the first time, you run these commands with the `--dry-run` flag.

When adding to the Brewfile, a flag must be used to specify what is being added:

```
bfm add --tap homebrew/dupes
bfm add --brew vim --args HEAD,with-override-system-vi
bfm add --brew crisidev/chunkwm/chunkwm --restart-service changed
bfm add --cask macvim
bfm add --mas Xcode --mas-id 497799835
```

Additional arguments for brew dependencies can be specified with the `--args` flag and service restart behaviour [always, changed] can be specified with the `--restart-service` flag.

The same flags must also be used with the `remove` and `check` commands.

```
bfm remove --tap homebrew/dupes
bfm check --brew vim
bfm remove --cask macvim
bfm check --mas Xcode
```

The `clean` command will organise your Brewfile and sort it into sections:

Before: 
```
mas 'Xcode', id: 497799835
cask 'macvim'
brew 'vim'
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
brew 'vim'
 
cask 'firefox'
cask 'macvim'
 
mas 'Xcode', id: 497799835
```

The Brewfile is automatically cleaned after every `add` and `remove` operation.

### Automatic Dependency Addition and Removal

The`--required` and `--all` flags can be used with the `add` and `remove` commands for
brew packages.

When used with `add`, all required dependencies or all required, recommended,
optional and build dependencies of a package can also be added to the Brewfile at the same time.

When used with the `remove` command, all required dependencies or all required, recommended,
optional and build dependencies of a package can also be removed from the Brewfile if they are
not required by any other package in the Brewfile.

The default behaviour for both the add and remove commands is not add or remove any of the package's dependencies.

Consider the NeoVim brew package:
```
neovim: stable 0.2.0 (bottled), HEAD
Ambitious Vim-fork focused on extensibility and agility
https://neovim.io/
/usr/local/Cellar/neovim/0.2.0_1 (1,352 files, 17.2MB) *
  Poured from bottle on 2017-08-04 at 15:25:02
From: https://github.com/Homebrew/homebrew-core/blob/master/Formula/neovim.rb
==> Dependencies
Build: luajit, cmake, lua@5.1, pkg-config
Required: gettext, jemalloc, libtermkey, libuv, libvterm, msgpack, unibilium
```

#### Add Examples

`bfm add -b neovim` will result in:

```
brew 'neovim'
```

`bfm add -b neovim --required` will result in:

```
brew 'gettext' # required by: neovim
brew 'jemalloc' # required by: neovim
brew 'libtermkey' # required by: neovim
brew 'libuv' # required by: neovim
brew 'libvterm' # required by: neovim
brew 'msgpack # required by: neovim
brew 'neovim'
brew 'unibilium' # required by: neovim
```

`bfm add -b neovim --all` will result in:

```
brew 'cmake'
brew 'gettext' # required by: neovim
brew 'jemalloc' # required by: neovim
brew 'libtermkey' # required by: neovim
brew 'libuv' # required by: neovim
brew 'libvterm' # required by: neovim
brew 'lua@5.1'
brew 'luajit'
brew 'msgpack # required by: neovim
brew 'neovim'
brew 'pkg-config'
brew 'unibilium' # required by: neovim
```

#### Remove Examples

If NeoVim has been added to the Brewfile with the `--all` flag, `bfm remove -b neovim` will result in:

```
brew 'cmake'
brew 'gettext'
brew 'jemalloc'
brew 'libtermkey'
brew 'libuv'
brew 'libvterm'
brew 'lua@5.1'
brew 'luajit'
brew 'msgpack
brew 'pkg-config'
brew 'unibilium'
```

If we also have the GnuPG package in our Brewfile, which also requires GNU gettext,
`bfm remove -b neovim --required` will result in:

```
brew 'cmake'
brew 'gettext' # required by: gnupg
brew 'gnupg'
brew 'lua@5.1'
brew 'luajit'
brew 'pkg-config'
```

Similarly `bfm remove -b neovim --all` will result in:

```
brew 'gettext' # required by: gnupg
brew 'gnupg'
```
