**🚧Please note that Cix is very new. I use it for my own projects, but take care using it for yours!🚧**

# Cix - a minimal CI for nix

Cix is a project to make a minimal useful CI for use with nix.
It watches repositories, runs any tests and builds listed in `nix flake check`, and reports the status back to your forge.

It is small, easy to setup (one static binary, and one json configuration file) but hopefully it is useful.

It is very early days, but if you wish to try it, create a `config.json` like below, and run `nix run github:steeleduncan/cix -- path/to/config.json`

```
{
    "var": "$HOME/.cache/cix-var",
    "repositories": [
        {
            "branch": "main",
            "github": {
                "user": "steeleduncan",
                "repository": "cix",
                "statuspat": "a PAT with status permissions (and ideally nothing else) goes here"
            }
        }
    ]
}
```


## Things Cix does

- Watches a repository
- Runs tests with `nix flake check`
- Pushes a commit status to Github so you can see if the tests are running, passed or failed

## Things Cix won't do

I'm doing my best to lean on nix in any way possible keep Cix simple so these are not, and will not be, supported

- **Store logs or artefacts** Nix is reproducible, you can generate these locally, or get them from a shared binary cache
- **Serve a web front end** Your code forge is used as the front end of cix, it will push statuses there

If you are looking for a fuller featured CI, I urge you to take a look at Hydra. It is tougher to setup, but it does everything you are likely to need.

## Roadmap - things I want to add to Cix

[ ] **Bitbucket support**
[ ] **Success actions** essentially a `nix run` that is called on succeeding tests. This could be used for deploys
[ ] **Non-status notifiers** Discord, email, some shell script. Any of these would be useful
[ ] **Binary cache option** Part the reason I don't want to serve artefacts is that nix can do this through aa binary cache, but a configuration option needs to be passed to the checks for this

## Things I would love a PR for

These are things I'd love to see in Cix, but that I am unlikely to need, and thus do myself

[ ] **Other code forges** I only have projects on Github and Bitbucket, but i'd gladly accept PRs for these
[ ] **Non-flake checks** Personally, I only ever use flakes with nix, but if anyone uses cix for non-flake checks, please let me know how

## Licence

Copyright 2024 Duncan Steele

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
