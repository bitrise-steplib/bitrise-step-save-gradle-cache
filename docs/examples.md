### Examples

Check out [Workflow Recipes](https://github.com/bitrise-io/workflow-recipes#-key-based-caching-beta) for other platform-specific examples!

#### Minimal example
```yaml
steps:
- restore-gradle-cache@1: {}
- android-build@1: {}
- save-gradle-cache@1: {}
```
