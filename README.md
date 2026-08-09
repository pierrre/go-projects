# Go projects

This repository manages a collection of Go projects hosted under `github.com/pierrre/`.

It provides a `projects.txt` file listing the project names and a Makefile that fans out commands across all sibling repositories.

## projects.txt

The `projects.txt` file lists the Go project names managed by this repository.

Each line is a project name, one per line.
The names correspond to repository names under `github.com/pierrre/`.
They are sorted in ascending order and must not contain duplicates.

## Makefile

The Makefile defines targets that operate across every project listed in `projects.txt`.
Each target iterates over the project names and runs a command in the corresponding sibling directory (`../<project>`).

The main targets are:

- `all-git-clone`: clones every project repository into the parent directory.
- `all-copy-common`: copies shared files (`Makefile-common.mk`, `LICENSE`, `CODEOWNERS`, `.gitignore`, `.github`, `.golangci.yml`) into every project.
- `all-all`: runs `make all` (build, test, lint) in every project.
- `all-build`: runs `make build` in every project.
- `all-test`: runs `make test` in every project.
- `all-lint`: runs `make lint` in every project.
- `all-clean`: runs `make clean` in every project.
- `all-run COMMAND=<cmd>`: runs an arbitrary command in every project directory (defaults to `ls`).

The `GIT_REPOSITORY_PATTERN` variable controls the clone URL: SSH by default, HTTPS when `CODESPACES=true`.
