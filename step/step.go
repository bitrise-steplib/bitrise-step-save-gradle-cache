package step

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/cache"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

const (
	stepId = "save-gradle-cache"

	// Cache key template
	// OS + Arch: to guarantee that stack-specific content (absolute paths, binaries) are stored separately
	// checksum values:
	// - `**/*.gradle*`: Gradle build files in any submodule, including ones written in Kotlin (*.gradle.kts)
	// - `**/gradle-wrapper.properties`: contains exact Gradle version
	// - `**/gradle.properties`: contains Gradle config values
	// - `**/libs.versions.toml`: version catalog file, contains dependencies and their versions
	key = `{{ .OS }}-{{ .Arch }}-gradle-cache-{{ checksum "**/*.gradle*" "**/gradle-wrapper.properties" "**/gradle.properties" "**/libs.versions.toml" }}`
)

func gradleUserHome(envRepo env.Repository) string {
	if v := strings.TrimSpace(envRepo.Get("GRADLE_USER_HOME")); v != "" {
		return v
	}

	return "~/.gradle"
}

func konanDataDir(envRepo env.Repository) string {
	if v := strings.TrimSpace(envRepo.Get("KONAN_DATA_DIR")); v != "" {
		return v
	}

	return "~/.konan"
}

func cachePaths(gradleHome, konanDir string, saveTransforms bool) []string {
	paths := []string{

		// Dependency JARs
		gradleHome + "/caches/jars-*",

		// Dependency AARs
		gradleHome + "/caches/modules-*/files-*",
		gradleHome + "/caches/modules-*/metadata-*",

		// Generated JARs for plugins and build scripts
		// The `**` segment matches the version-specific folder, such as `7.6`.
		gradleHome + "/caches/**/generated-gradle-jars/*.jar",

		// Kotlin build script cache
		// The `**` segment matches the version-specific folder, such as `7.6`.
		gradleHome + "/caches/**/kotlin-dsl",

		// Cache of downloaded Gradle binary
		gradleHome + "/wrapper",

		// Configuration cache is saved by separate step: save-gradle-configuration-cache

		// JDKs downloaded by the toolchain support
		gradleHome + "/jdks",

		// Kotlin/Native, relocated by KONAN_DATA_DIR (not GRADLE_USER_HOME).
		konanDir + "/kotlin-*",
		konanDir + "/dependencies",
	}

	if saveTransforms {
		// Save artifact transforms
		// The `**` segment matches the version-specific folder, such as `7.6`.
		paths = append(paths,
			gradleHome+"/caches/**/transforms",
			gradleHome+"/caches/transforms-*",
		)
	}

	return paths
}

type Input struct {
	Verbose          bool `env:"verbose,required"`
	CompressionLevel int  `env:"compression_level,range[1..19]"`
	SaveTransforms   bool `env:"save_transforms"`
}

type SaveCacheStep struct {
	logger       log.Logger
	inputParser  stepconf.InputParser
	pathChecker  pathutil.PathChecker
	pathProvider pathutil.PathProvider
	pathModifier pathutil.PathModifier
	envRepo      env.Repository
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	pathChecker pathutil.PathChecker,
	pathProvider pathutil.PathProvider,
	pathModifier pathutil.PathModifier,
	envRepo env.Repository,
) SaveCacheStep {
	return SaveCacheStep{
		logger:       logger,
		inputParser:  inputParser,
		pathChecker:  pathChecker,
		pathProvider: pathProvider,
		pathModifier: pathModifier,
		envRepo:      envRepo,
	}
}

func (step SaveCacheStep) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)

	paths := cachePaths(gradleUserHome(step.envRepo), konanDataDir(step.envRepo), input.SaveTransforms)

	step.logger.Println()
	step.logger.Printf("Cache key: %s", key)
	step.logger.Printf("Cache paths:")
	step.logger.Printf(strings.Join(paths, "\n"))
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	saver := cache.NewSaver(step.envRepo, step.logger, step.pathProvider, step.pathModifier, step.pathChecker, nil)
	return saver.Save(cache.SaveCacheInput{
		StepId:           stepId,
		Verbose:          input.Verbose,
		Key:              key,
		Paths:            paths,
		IsKeyUnique:      true,
		CompressionLevel: input.CompressionLevel,
	})
}
