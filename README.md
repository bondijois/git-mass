# Git Mass 

Mass download all github repositories(public & private) of an organization, ideally in a few seconds.

> Writing this as a simple bash script would've taken me a few minutes, but this took me a few hours to learn and write.
> _Expect stupid mistakes._

## Installation

Download the binary from the [releases](https://github.com/bondijois/git-mass/releases) page and run it.

Or, you could clone the repository and build it yourself.

## Usage

Assuming you have the binary installed, you can do the following:

### Configure / Re-Configure

* Configure/re-configure your credentials:

```shell
git-mass config -u <username> -t <token>
```

* Verify if your configured credentials are valid:

```shell
git-mass config
```

> You can generate a token from [here](https://github.com/settings/tokens). Highly recommend you set token-expiry.

### List all (public)organizations

```shell
git-mass orgs
```

> Go to `https://github.com/orgs/<organization>/people` if you want to switch between public/private.

### Clone all repositories from an organization

```shell
git-mass clone -o <organization>
```

---

# TODOs

* Might add GitLab integration
