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

Once installed, add `export BFM_BREWFILE=/path/to/your/Brewfile` to your `.bashrc` or `.zshrc` or `.${shell}rc` file.

## Usage

The main commands of bfm are `add`, `remove`, `check` and `clean`.

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
