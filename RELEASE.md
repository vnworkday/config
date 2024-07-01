# Release Guide

Execute the following commands to release a new version of the project:

```bash
sh scripts/go-release.sh -t $type
```

Where `$type` is one of the following:

- `major`: for a major release, e.g., v1.0.0, v2.0.0, etc. Mainly for breaking changes.
- `minor`: for a minor release, e.g., v0.1.0, v0.2.0, etc. Mainly for new features.
- `patch`: for a patch release, e.g., v0.0.1, v0.0.2, etc. Mainly for bug fixes.

The script automates the process of releasing the Go module with the following steps:

1. Makes sure you are on the `main` branch. If you are not, it exits.
2. Makes sure the version is matched in the `package.json` file. If it is not, it exits.
3. Makes sure the git repository is clean. If it is not, it exits.
4. Gets the latest tag from the package.json file.
5. Increments the version using `npm version $type`.
6. Generates Go code using `npm run gen:go`.
7. Stages the files and commits them.
8. Pushes the changes to the remote repository and creates a new tag.

The remaining steps are done by the GitHub Actions workflow (`.github/workflows/release.yml`). It will build the Go code, create a release draft, and upload the binaries to the release page.