**🚧 Cix is very new. I use it for my own projects, but take care using it for yours! 🚧**

# Cix - A minimal CI for Nix

Cix is a minimal, but useful, CI for use with Nix.
It watches repositories, runs any tests and builds listed in `nix flake check`, and reports the status back to your code forge (e.g. Github or Bitbucket).

Cix is not designed to be the best, or the most featureful, Nix CI.
It is designed to be the smallest and easiest to setup, while still being useful, something you can run in the background and see the status of each commit in Github/Bitbucket.
With one static binary, and one JSON configuration file, it will run tests for all commits you push, and report a status back, which show as a green tick or red cross against the commit.

A lot of simplicity is found by not storing logs and artefacts, or serving them through a web interface.
Cix relies on Nix's reproducibility instead, you can always rerun the full test command to see the results with perfect reproducibility.
For this reason the command is included in the status line on your code forge, if you see a red cross you can run it locally and see the problem.
If you share a binary cache with the runner, you will share its results without needing to calculate them, so this be efficient.

It is early days for Cix, but if you wish to try it, create a `config.json` similar to the following and run `nix run github:steeleduncan/cix -- config.json`

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

The `statuspat` is technically optional, but the commit tick in Github/Bitbucket is currently the only way to view results.
See below for the Bitbucket equivalent, and details of what permissions to allow on the token.

Cix will use git to pull the repositories over SSH, using whatever permissions are available in that context.
It is also designed so that Cix doesn't need to run continuously, if it is running your laptop when it goes to sleep, it will pull and run tests again on wake.

Cix was inspired by [nix-simple-ci](https://github.com/ElvishJerricco/nix-simple-ci)

## Things Cix does

- **Watches repositories**
- **Runs tests** with `nix flake check`
- **Pushes a commit status to Github//Bitbucket** so you can see if the tests are running, passed or failed
- **Catches up** Cix doesn't need to be online when the commit is made, so if you only have your machine on part the time, when it first checks it will enumerate and test all commits made since it was last on

Cix will run tests for every commit, not just the latest commit pushed.
However it won't run tests for commits before it was activated

## Things Cix won't do

Nix and your code forge have almost everything needed for a useful CI system, so with Cix I am doing my best to keep it minimal and rely on Nix wherever possible.
There are currently no plans to implement the features below

- **Store logs or artefacts** Nix is reproducible, you can generate these locally, or get them from a shared binary cache
- **Serve a web front end** Your code forge is used as the front end of cix, it will push statuses there

If you are looking for a fuller featured CI, I urge you to take a look at [Hydra](https://nixos.wiki/wiki/Hydra).
It is tougher to setup, but it does everything you are likely to need.

## Roadmap - things I want to add to Cix

- [ ] **Success actions** essentially a `nix run` that is called on success. This could be used for deploys
- [ ] **Non-status notifiers** e.g. a Discord, Slack or email message on success and/or failure
- [ ] **Binary cache option** Part the reason I don't want to serve artefacts is that Nix can do this through aa binary cache, but a configuration option needs to be passed to the checks for this
- [ ] **Repository maintenance** GC, prune, etc. Cix works by keeping a local copy of the repository in the var folder specified in the config. Most likely this would need the occasional GC
- [ ] **Leave logs as a comment** It would be helpful if logs were left as a comment on the commit when tests fail
- [ ] **Parallel tests** I imagine Cix being used in situations where you want some CPU left spare (e.g. if it runs on your dev machine), but it would be nice to have an option to parallelise and run multiple tests/builds in parallel

## Things I would love a PR for

These are things I'd love to see in Cix, but that I am unlikely to need, and thus do myself, but I'd gladly accept PRs for these

- [ ] **Other code forges** I only have projects on Github and Bitbucket, but htere are many other code forges it would be great if Cix supported
- [ ] **Non-flake checks** Personally, I only ever use flakes with Nix, but there are non-flake approaches I am not familiar with
- [ ] **Non-SSH access** Currently Cix uses the git binary and any SSH credentials available to it to pull commits. There are other approaches, and it would be useful to include these

## Alternatives

Depending on your needs the following might be useful

- **Hydra** The classic Nix CI, full featured, and probably what you are looking for if you have demanding requirements, and the time to maintain it
- **[github-nix-ci](https://github.com/juspay/github-nix-ci)** Run your own Github Actions self hosted runners on NixOS to use the GHA UI and your own hardware

## `Config.json` format

- `var` (required) A path to a work folder where cix may store copies of the repositories
- `name` (optional) A name for this runner, reported in the comment on code forge commit
- `timeout` (optional) Job timeout in seconds (defaults to 15 mins)
- `pollinginterval` (optional) Polling interval in seconds (defaults to 180s)
- `repositories` (required) A list of repositories
    - `branch` (required) The branch to test
    - `github` (optional)
        - `user` (required) User name on Github
        - `repository` (required) Repository name for that users account
        - `statuspat` (optional) A Personal Access Token with commit status read/write
    - `bitbucket` (optional)
        - `workspace` (required) The Bitbucket workspace name
        - `repository` (required) The Bitbucket repository slug
        - `token` (optional) An Access Token with write permission for the repository
    - `ssh` (optional)
        - `remote` (required) An ssh git url to pull commits from

Although `github`, `bitbucket`, `ssh` fields are all optional, you must have at least one per repository specified.
If more than one is specified the outcome is undefined.

## Licence

Copyright 2024 Duncan Steele

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
