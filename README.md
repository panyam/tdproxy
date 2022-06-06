## TDProxy

A golang server that provides an auth and API proxy to the TD Ameritrade developer api.

Why is this needed?

Developers using the TDA Dev API are faced with the following constraints:

* Limited to a single auth.  A usecase for me was to have multiple friends leverage the API with their own auth-tokens without the need for managing their own clients to TDA and setting up their own callback endpoints (for oauth2).
* No caching.   When trying out different strategies, instruments need to be fetched from TDA but these do not need to be upto date.  Different caching and timeout mechanisms need to be reimplemented.  This proxy hides all that way and provides timeout SLOs transparently to the developer.
* Lack of persistence.   With trades and strategies modelled, the proxy also provides persistence of the trades into a DB that can be queried and ordered based on developer/modeller needs.
* Seperate Streaming API.  TDA also provides a streaming API however that is hard to work with and needs other customer services.  This proxy ensures (in development) a consistent interface to the streaming protocol so that the instruments (Options, Stock prices etc) have the same semantics (including caching).

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
