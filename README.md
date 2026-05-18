# Save Gradle Cache

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-save-gradle-cache?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-save-gradle-cache/releases)

Saves Gradle caches. This Step needs to be used in combination with **Restore Gradle Cache**.

<details>
<summary>Description</summary>

Saves Gradle dependency caches. This Step needs to be used in combination with **Restore Gradle Cache**.

This Step is based on [key-based caching](https://devcenter.bitrise.io/en/builds/caching/key-based-caching.html) and sets up the cache key and path automatically for Gradle. If you'd like to change the cache key (or paths to cache), you might want to use the generic [Save cache](https://github.com/bitrise-steplib/bitrise-step-save-cache) Step instead.

#### Related steps

[Restore Gradle cache](https://github.com/bitrise-steplib/bitrise-step-restore-gradle-cache/)

[Save cache](https://github.com/bitrise-steplib/bitrise-step-save-cache/)

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

### Examples

Check out [Workflow Recipes](https://github.com/bitrise-io/workflow-recipes#-key-based-caching-beta) for other platform-specific examples!

#### Minimal example
```yaml
steps:
- restore-gradle-cache@1: {}
- android-build@1: {}
- save-gradle-cache@1: {}
```


## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `verbose` | Enable logging additional information for troubleshooting | required | `false` |
| `compression_level` | Zstd compression level to control speed / archive size. Set to 1 for fastest option. Valid values are between 1 and 19. Defaults to 3. |  | `3` |
| `save_transforms` | Save transformed classes and resources under `$HOME/.gradle/caches/**/transforms` | required | `false` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-save-gradle-cache/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-save-gradle-cache/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

**Note:** this step's end-to-end tests (defined in `e2e/bitrise.yml`) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.


Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
