include Makefile-common.mk

CODESPACES?=false

ALL_COMMAND=cat projects.txt | xargs -I {} $(1)
ALL_RUN=$(call ALL_COMMAND,sh -c 'echo {} && cd ../{} && $(1)')
# Run an arbitrary command in every project (COMMAND=<cmd>, defaults to ls).
.PHONY: all-run
all-run:
	$(eval COMMAND?=ls)
	$(call ALL_RUN,$(COMMAND))

ifeq ($(CODESPACES),true)
GIT_REPOSITORY_PATTERN=https://github.com/pierrre/{}.git
else
GIT_REPOSITORY_PATTERN=git@github.com:pierrre/{}.git
endif
# Clone every project repository into the parent directory.
.PHONY: all-git-clone
all-git-clone:
	$(call ALL_COMMAND,sh -c "(ls ../{} > /dev/null 2>&1 || git -C .. clone $(GIT_REPOSITORY_PATTERN))")

# Copy the shared files into every project.
.PHONY: all-copy-common
all-copy-common:
	$(call ALL_COMMAND,cp -r Makefile-common.mk LICENSE CODEOWNERS .gitignore .github .golangci.yml ../{})

# Run make all in every project.
.PHONY: all-all
all-all:
	$(call ALL_RUN,make all)

# Run make build in every project.
.PHONY: all-build
all-build:
	$(call ALL_RUN,make build)

# Run make test in every project.
.PHONY: all-test
all-test:
	$(call ALL_RUN,make test)

# Run make generate in every project.
.PHONY: all-generate
all-generate:
	$(call ALL_RUN,make generate)

# Run make lint in every project. Does not modify files.
.PHONY: all-lint
all-lint:
	$(call ALL_RUN,make lint)

# Run make lint-fix in every project. Modifies files.
.PHONY: all-lint-fix
all-lint-fix:
	$(call ALL_RUN,make lint-fix)

# Run make golangci-lint in every project. Does not modify files.
.PHONY: all-golangci-lint
all-golangci-lint:
	$(call ALL_RUN,make golangci-lint)

# Run make golangci-lint-fix in every project. Modifies files.
.PHONY: all-golangci-lint-fix
all-golangci-lint-fix:
	$(call ALL_RUN,make golangci-lint-fix)

# Run make lint-rules in every project. Does not modify files.
.PHONY: all-lint-rules
all-lint-rules:
	$(call ALL_RUN,make lint-rules)

# Run make mod-update in every project. Modifies go.mod and go.sum.
.PHONY: all-mod-update
all-mod-update: all-copy-common
	$(call ALL_RUN,make mod-update)

# Run make mod-update-pierrre in every project. Modifies go.mod and go.sum.
.PHONY: all-mod-update-pierrre
all-mod-update-pierrre: all-copy-common
	$(call ALL_RUN,make mod-update-pierrre)

# Run make mod-tidy in every project. Modifies go.mod and go.sum.
.PHONY: all-mod-tidy
all-mod-tidy:
	$(call ALL_RUN,make mod-tidy)

# Run make mod-tidy-diff in every project. Does not modify files.
.PHONY: all-mod-tidy-diff
all-mod-tidy-diff:
	$(call ALL_RUN,make mod-tidy-diff)

# Run make clean in every project.
.PHONY: all-clean
all-clean:
	$(call ALL_RUN,make clean)
