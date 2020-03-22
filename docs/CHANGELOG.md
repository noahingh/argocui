# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0]
### Added
- Add mocking command to test([#12](https://github.com/hanjunlee/argocui/pull/12)) UI stuff easily.
- Confirm before delete([#16](https://github.com/hanjunlee/argocui/pull/16)).
- Add UI messenger to show the error message([#20](https://github.com/hanjunlee/argocui/pull/20)).

### Fixed
- Fix the position of namespace([#10](https://github.com/hanjunlee/argocui/pull/10)).
- Refactoring the structure of codes to exand features easily([#13](https://github.com/hanjunlee/argocui/pull/13)).

### Deleted
- Delete read only option.

## [0.0.3]
### Added
- Add the `version` command option([#9](https://github.com/hanjunlee/argocui/pull/9)).

### Fixed
- Fix the version in the info view.

## [0.0.2]
### Added
- Switch to another namespace, but it doesn't change the context([#8](https://github.com/hanjunlee/argocui/pull/8)).
- Toggle namespace the current namespace and the global namespace(`*`)([#7](https://github.com/hanjunlee/argocui/pull/7)).

## [0.0.1]
### Added
- List up Argo workflows, same as `argo list`.
- Get the tree of a Argo workflow, same as `argo get`.
- Follow logs of a Argo workflow, same as `argo logs`.
- Delete a Argo workflow.
- Search Argo workflow.

### Changed

### Removed
