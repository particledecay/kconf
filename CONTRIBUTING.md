# Contributing

## Coding Guidelines

1. Update the [README](README.md) with any changes to the interface if applicable.
2. Include unit tests where applicable (your code should not reduce test coverage).
3. All commit messages must follow [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0/) format.
4. Please submit Pull Requests from a feature branch in your fork, not from master.
5. Pull Requests with more than one commit will be rejected (squash your commits).

## pre-commit
If you want the best chance of having your Pull Request approved, install [pre-commit](https://pre-commit.com) and then run the following command in your local project directory:

```sh
pre-commit install -t pre-commit -t commit-msg
```

This will ensure that your code is likely to pass all checks.
