# Contributing

## Coding Guidelines

1. Update the [README](README.md) with any changes to the interface if applicable.
2. Include unit tests where applicable (your code should not reduce test coverage).
3. All commit messages must follow [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0/) format.
4. Please submit pull requests from a feature branch in your fork, not from `main` or `master`.

## Pre-Commit

If you want the best chance of having your pull request approved, install [pre-commit](https://pre-commit.com) (if you use [asdf](https://asdf-vm.com) there's already a `.tool-versions` file for you) and then run the following command in your local project directory:

```sh
pre-commit install
```

This will ensure that your code is likely to pass all checks.
