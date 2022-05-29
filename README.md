
## Dev Setup

### Install golang

OSX:

```
brew install golang
```

### Install python (for python bindings)

#### Virtual Env and Requirements

```
python3 -m venv env
source env/bin/activate
pip install -r dev_requirements
```

### Generate bindings (not required if you are not changing golang code)

```
# Needed to get the pypslite package to generate files into
# git submodule init
git submodule update --init
sh build.sh
```
