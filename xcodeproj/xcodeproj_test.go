package xcodeproj

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSharedSchemeFilePath(t *testing.T) {
	// regexpPattern := filepath.Join(".*[/]?xcshareddata", "xcschemes", ".+[.]xcscheme")
	require.Equal(t, true, isSharedSchemeFilePath("/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme"))
	require.Equal(t, true, isSharedSchemeFilePath("./BitriseXcode7Sample.xcodeproj/xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme"))
	require.Equal(t, true, isSharedSchemeFilePath("./xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme"))
	require.Equal(t, true, isSharedSchemeFilePath("xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme"))

	// incorrect paths
	require.Equal(t, false, isSharedSchemeFilePath("./xcschemes/BitriseXcode7Sample.xcscheme"))
	require.Equal(t, false, isSharedSchemeFilePath("./xcshareddata/BitriseXcode7Sample.xcscheme"))
	require.Equal(t, false, isSharedSchemeFilePath("./BitriseXcode7Sample.xcscheme"))
	require.Equal(t, false, isSharedSchemeFilePath("BitriseXcode7Sample.xcscheme"))
	require.Equal(t, false, isSharedSchemeFilePath("xcshareddata/xcschemes/.xcscheme"))

	// user scheme
	require.Equal(t, false, isSharedSchemeFilePath("/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme"))
}

func TestFilterSharedSchemeFilePaths(t *testing.T) {
	t.Log("it finds shared schemes")
	{
		paths := []string{"/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme"}
		filteredPaths := filterSharedSchemeFilePaths(paths)
		require.Equal(t, 1, len(filteredPaths))
		require.Equal(t, "/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme", filteredPaths[0])
	}

	t.Log("it omitts user schemes")
	{
		paths := []string{
			"/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme",
			"/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme",
		}
		filteredPaths := filterSharedSchemeFilePaths(paths)
		require.Equal(t, 1, len(filteredPaths))
		require.Equal(t, "/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme", filteredPaths[0])
	}

	t.Log("it works for relative paths")
	{
		paths := []string{
			"./xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme",
			"./xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme",
		}
		filteredPaths := filterSharedSchemeFilePaths(paths)
		require.Equal(t, 1, len(filteredPaths))
		require.Equal(t, "./xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme", filteredPaths[0])
	}

	t.Log("it works for relative paths")
	{
		paths := []string{
			"xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme",
			"xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme",
		}
		filteredPaths := filterSharedSchemeFilePaths(paths)
		require.Equal(t, 1, len(filteredPaths))
		require.Equal(t, "xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme", filteredPaths[0])
	}

	t.Log("it filters base file paths")
	{
		paths := []string{
			"BitriseXcode7Sample.xcscheme",
			"BitriseXcode7SampleTest.xcscheme",
		}
		filteredPaths := filterSharedSchemeFilePaths(paths)
		require.Equal(t, 0, len(filteredPaths))
	}
}

func TestSchemeNameFromPath(t *testing.T) {
	require.Equal(t, "BitriseXcode7Sample", SchemeNameFromPath("/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme"))
	require.Equal(t, "BitriseXcode7SampleTest", SchemeNameFromPath("/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme"))

	require.Equal(t, "BitriseXcode7Sample", SchemeNameFromPath("./BitriseXcode7Sample.xcscheme"))
	require.Equal(t, "BitriseXcode7SampleTest", SchemeNameFromPath("./BitriseXcode7SampleTest.xcscheme"))

	require.Equal(t, "", SchemeNameFromPath(".xcscheme"))
	require.Equal(t, "", SchemeNameFromPath("xcscheme"))
}

func TestIsUserSchemeFilePath(t *testing.T) {
	// regexpPattern := filepath.Join(".*[/]?xcuserdata", ".*[.]xcuserdatad", "xcschemes", ".+[.]xcscheme")
	require.Equal(t, true, isUserSchemeFilePath("/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme"))
	require.Equal(t, true, isUserSchemeFilePath("./BitriseXcode7Sample.xcodeproj/xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme"))
	require.Equal(t, true, isUserSchemeFilePath("./xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme"))
	require.Equal(t, true, isUserSchemeFilePath("xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme"))

	// unknown user
	require.Equal(t, true, isUserSchemeFilePath("xcuserdata/.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme"))

	// incorrect paths
	require.Equal(t, false, isUserSchemeFilePath("./bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme"))
	require.Equal(t, false, isUserSchemeFilePath("./xcuserdata/xcschemes/BitriseXcode7SampleTest.xcscheme"))
	require.Equal(t, false, isUserSchemeFilePath("./xcuserdata/bitrise.xcuserdatad/BitriseXcode7SampleTest.xcscheme"))
	require.Equal(t, false, isUserSchemeFilePath("BitriseXcode7SampleTest.xcscheme"))
	require.Equal(t, false, isUserSchemeFilePath("xcuserdata/bitrise.xcuserdatad/xcschemes/.xcscheme"))

	// shared scheme
	require.Equal(t, false, isUserSchemeFilePath("/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme"))
}

func TestFilterUserSchemeFilePaths(t *testing.T) {
	t.Log("it finds user schemes")
	{
		paths := []string{"/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme"}
		filteredPaths := filterUserSchemeFilePaths(paths)
		require.Equal(t, 1, len(filteredPaths))
		require.Equal(t, "/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme", filteredPaths[0])
	}

	t.Log("it omitts shared schemes")
	{
		paths := []string{
			"/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme",
			"/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme",
		}
		filteredPaths := filterUserSchemeFilePaths(paths)
		require.Equal(t, 1, len(filteredPaths))
		require.Equal(t, "/Users/bitrise/BitriseXcode7Sample.xcodeproj/xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme", filteredPaths[0])
	}

	t.Log("it works for relative paths")
	{
		paths := []string{
			"./xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme",
			"./xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme",
		}
		filteredPaths := filterUserSchemeFilePaths(paths)
		require.Equal(t, 1, len(filteredPaths))
		require.Equal(t, "./xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme", filteredPaths[0])
	}

	t.Log("it works for relative paths")
	{
		paths := []string{
			"xcshareddata/xcschemes/BitriseXcode7Sample.xcscheme",
			"xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme",
		}
		filteredPaths := filterUserSchemeFilePaths(paths)
		require.Equal(t, 1, len(filteredPaths))
		require.Equal(t, "xcuserdata/bitrise.xcuserdatad/xcschemes/BitriseXcode7SampleTest.xcscheme", filteredPaths[0])
	}

	t.Log("it filters base file paths")
	{
		paths := []string{
			"BitriseXcode7Sample.xcscheme",
			"BitriseXcode7SampleTest.xcscheme",
		}
		filteredPaths := filterUserSchemeFilePaths(paths)
		require.Equal(t, 0, len(filteredPaths))
	}
}

func TestSchemeFileContentContainsXCTestBuildAction(t *testing.T) {
	t.Log("Contains XCTestBuildAction")
	{
		schemeContent := schemeContentWithXCTestBuildAction

		contains, err := schemeFileContentContainsXCTestBuildAction(schemeContent)
		require.NoError(t, err)
		require.Equal(t, true, contains)
	}

	t.Log("Do NOT contains XCTestBuildAction")
	{
		schemeContent := schemeContentWithoutXCTestBuildAction

		contains, err := schemeFileContentContainsXCTestBuildAction(schemeContent)
		require.NoError(t, err)
		require.Equal(t, false, contains)
	}
}

func TestIsXCodeProj(t *testing.T) {
	require.Equal(t, true, IsXCodeProj(XCodeProjExt))
	require.Equal(t, false, IsXCodeProj(XCWorkspaceExt))

	require.Equal(t, true, IsXCodeProj("a"+XCodeProjExt))
	require.Equal(t, false, IsXCodeProj("a"+XCWorkspaceExt))

	require.Equal(t, true, IsXCodeProj("./SampleAppWithCocoapods.xcodeproj"))
	require.Equal(t, false, IsXCodeProj("./SampleAppWithCocoapods.xcworkspace"))

	require.Equal(t, true, IsXCodeProj("/Users/bitrise/SampleAppWithCocoapods/SampleAppWithCocoapods.xcodeproj"))
	require.Equal(t, false, IsXCodeProj("/Users/bitrise/SampleAppWithCocoapods/SampleAppWithCocoapods.xcworkspace"))

	require.Equal(t, false, IsXCodeProj("xcworkspace"))
	require.Equal(t, false, IsXCodeProj("xcodeproj"))
}

func TestIsXCWorkspace(t *testing.T) {
	require.Equal(t, false, IsXCWorkspace(XCodeProjExt))
	require.Equal(t, true, IsXCWorkspace(XCWorkspaceExt))

	require.Equal(t, false, IsXCWorkspace("a"+XCodeProjExt))
	require.Equal(t, true, IsXCWorkspace("a"+XCWorkspaceExt))

	require.Equal(t, false, IsXCWorkspace("./SampleAppWithCocoapods.xcodeproj"))
	require.Equal(t, true, IsXCWorkspace("./SampleAppWithCocoapods.xcworkspace"))

	require.Equal(t, false, IsXCWorkspace("/Users/bitrise/SampleAppWithCocoapods/SampleAppWithCocoapods.xcodeproj"))
	require.Equal(t, true, IsXCWorkspace("/Users/bitrise/SampleAppWithCocoapods/SampleAppWithCocoapods.xcworkspace"))

	require.Equal(t, false, IsXCWorkspace("xcworkspace"))
	require.Equal(t, false, IsXCWorkspace("xcodeproj"))
}

func TestPBXProjContentTartgets(t *testing.T) {
	t.Log("target with space")
	{
		content := pbxNativeTargetSectionWithSpace

		targetMap, err := pbxprojContentTartgets(content)
		require.NoError(t, err)
		require.Equal(t, 1, len(targetMap))

		hasXCTest, found := targetMap["BitriseSampleAppsiOS With Spaces"]
		require.Equal(t, true, found)
		require.Equal(t, true, hasXCTest)
	}

	t.Log("simple targets")
	{
		content := pbxProjContentChunk

		targetMap, err := pbxprojContentTartgets(content)
		require.NoError(t, err)
		require.Equal(t, 1, len(targetMap))

		hasXCTest, found := targetMap["SampleAppWithCocoapods"]
		require.Equal(t, true, found)
		require.Equal(t, true, hasXCTest)
	}
}

func TestParsePBXTargetDependencies(t *testing.T) {
	desiredDependencieMap := map[string]PBXTargetDependency{
		"BAAFFEEF19EE788800F3AC91": PBXTargetDependency{
			id:     "BAAFFEEF19EE788800F3AC91",
			isa:    "PBXTargetDependency",
			target: "BAAFFED019EE788800F3AC91",
		},
	}

	dependencies, err := parsePBXTargetDependencies(pbxTargetDependencies)
	require.NoError(t, err)
	require.Equal(t, len(desiredDependencieMap), len(dependencies))

	matchingDependencieMap := map[string]bool{}
	for _, dependencie := range dependencies {
		desiredDependencie, found := desiredDependencieMap[dependencie.id]
		require.Equal(t, true, found)

		require.Equal(t, desiredDependencie.id, dependencie.id)
		require.Equal(t, desiredDependencie.isa, dependencie.isa)
		require.Equal(t, desiredDependencie.target, dependencie.target)

		matchingDependencieMap[dependencie.id] = true
	}

	require.Equal(t, len(desiredDependencieMap), len(matchingDependencieMap))
}

func TestParsePBXNativeTargets(t *testing.T) {
	desiredTargetMap := map[string]PBXNativeTarget{
		"BADDF9E61A703F87004C3526": PBXNativeTarget{
			id:           "BADDF9E61A703F87004C3526",
			isa:          "PBXNativeTarget",
			dependencies: []string{},
			name:         "BitriseSampleAppsiOS With Spaces",
			productPath:  "BitriseSampleAppsiOS With Spaces.app",
			productType:  "com.apple.product-type.application",
		},
		"BADDFA021A703F87004C3526": PBXNativeTarget{
			id:  "BADDFA021A703F87004C3526",
			isa: "PBXNativeTarget",
			dependencies: []string{
				"BADDFA051A703F87004C3526",
			},
			name:        "BitriseSampleAppsiOS With SpacesTests",
			productPath: "BitriseSampleAppsiOS With SpacesTests.xctest",
			productType: "com.apple.product-type.bundle.unit-test",
		},
	}

	targets, err := parsePBXNativeTargets(pbxNativeTargetSectionWithSpace)
	require.NoError(t, err)
	require.Equal(t, len(desiredTargetMap), len(targets))

	matchingTargetMap := map[string]bool{}
	for _, target := range targets {
		desiredTarget, found := desiredTargetMap[target.id]
		require.Equal(t, true, found)

		require.Equal(t, desiredTarget.id, target.id)
		require.Equal(t, desiredTarget.isa, target.isa)
		require.Equal(t, desiredTarget.name, target.name)
		require.Equal(t, desiredTarget.productPath, target.productPath)
		require.Equal(t, desiredTarget.productType, target.productType)
		require.Equal(t, len(desiredTarget.dependencies), len(target.dependencies))

		foundDependencieMap := map[string]bool{}
		for _, dep := range target.dependencies {
			for _, desiredDep := range desiredTarget.dependencies {
				if dep == desiredDep {
					foundDependencieMap[dep] = true
				}
			}
		}
		require.Equal(t, len(desiredTarget.dependencies), len(foundDependencieMap))

		matchingTargetMap[target.id] = true
	}

	require.Equal(t, len(desiredTargetMap), len(matchingTargetMap))
}
