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
		schemeContent := `<?xml version="1.0" encoding="UTF-8"?>
<Scheme
   LastUpgradeVersion = "0700"
   version = "1.3">
   <BuildAction
      parallelizeBuildables = "YES"
      buildImplicitDependencies = "YES">
      <BuildActionEntries>
         <BuildActionEntry
            buildForTesting = "YES"
            buildForRunning = "YES"
            buildForProfiling = "YES"
            buildForArchiving = "YES"
            buildForAnalyzing = "YES">
            <BuildableReference
               BuildableIdentifier = "primary"
               BlueprintIdentifier = "BAC384091BA9F569005CFE20"
               BuildableName = "BitriseXcode7Sample.app"
               BlueprintName = "BitriseXcode7Sample"
               ReferencedContainer = "container:BitriseXcode7Sample.xcodeproj">
            </BuildableReference>
         </BuildActionEntry>
      </BuildActionEntries>
   </BuildAction>
   <TestAction
      buildConfiguration = "Debug"
      selectedDebuggerIdentifier = "Xcode.DebuggerFoundation.Debugger.LLDB"
      selectedLauncherIdentifier = "Xcode.DebuggerFoundation.Launcher.LLDB"
      shouldUseLaunchSchemeArgsEnv = "YES">
      <Testables>
         <TestableReference
            skipped = "NO">
            <BuildableReference
               BuildableIdentifier = "primary"
               BlueprintIdentifier = "BAC384221BA9F569005CFE20"
               BuildableName = "BitriseXcode7SampleTests.xctest"
               BlueprintName = "BitriseXcode7SampleTests"
               ReferencedContainer = "container:BitriseXcode7Sample.xcodeproj">
            </BuildableReference>
         </TestableReference>
         <TestableReference
            skipped = "NO">
            <BuildableReference
               BuildableIdentifier = "primary"
               BlueprintIdentifier = "BAC3842D1BA9F569005CFE20"
               BuildableName = "BitriseXcode7SampleUITests.xctest"
               BlueprintName = "BitriseXcode7SampleUITests"
               ReferencedContainer = "container:BitriseXcode7Sample.xcodeproj">
            </BuildableReference>
         </TestableReference>
      </Testables>
      <MacroExpansion>
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "BAC384091BA9F569005CFE20"
            BuildableName = "BitriseXcode7Sample.app"
            BlueprintName = "BitriseXcode7Sample"
            ReferencedContainer = "container:BitriseXcode7Sample.xcodeproj">
         </BuildableReference>
      </MacroExpansion>
      <AdditionalOptions>
      </AdditionalOptions>
   </TestAction>
   <LaunchAction
      buildConfiguration = "Debug"
      selectedDebuggerIdentifier = "Xcode.DebuggerFoundation.Debugger.LLDB"
      selectedLauncherIdentifier = "Xcode.DebuggerFoundation.Launcher.LLDB"
      launchStyle = "0"
      useCustomWorkingDirectory = "NO"
      ignoresPersistentStateOnLaunch = "NO"
      debugDocumentVersioning = "YES"
      debugServiceExtension = "internal"
      allowLocationSimulation = "YES">
      <BuildableProductRunnable
         runnableDebuggingMode = "0">
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "BAC384091BA9F569005CFE20"
            BuildableName = "BitriseXcode7Sample.app"
            BlueprintName = "BitriseXcode7Sample"
            ReferencedContainer = "container:BitriseXcode7Sample.xcodeproj">
         </BuildableReference>
      </BuildableProductRunnable>
      <AdditionalOptions>
      </AdditionalOptions>
   </LaunchAction>
   <ProfileAction
      buildConfiguration = "Release"
      shouldUseLaunchSchemeArgsEnv = "YES"
      savedToolIdentifier = ""
      useCustomWorkingDirectory = "NO"
      debugDocumentVersioning = "YES">
      <BuildableProductRunnable
         runnableDebuggingMode = "0">
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "BAC384091BA9F569005CFE20"
            BuildableName = "BitriseXcode7Sample.app"
            BlueprintName = "BitriseXcode7Sample"
            ReferencedContainer = "container:BitriseXcode7Sample.xcodeproj">
         </BuildableReference>
      </BuildableProductRunnable>
   </ProfileAction>
   <AnalyzeAction
      buildConfiguration = "Debug">
   </AnalyzeAction>
   <ArchiveAction
      buildConfiguration = "Release"
      revealArchiveInOrganizer = "YES">
   </ArchiveAction>
</Scheme>`

		contains, err := schemeFileContentContainsXCTestBuildAction(schemeContent)
		require.NoError(t, err)
		require.Equal(t, true, contains)
	}

	t.Log("Do NOT contains XCTestBuildAction")
	{
		schemeContent := `<?xml version="1.0" encoding="UTF-8"?>
<Scheme
   LastUpgradeVersion = "0730"
   version = "1.3">
   <BuildAction
      parallelizeBuildables = "YES"
      buildImplicitDependencies = "YES">
      <BuildActionEntries>
         <BuildActionEntry
            buildForTesting = "YES"
            buildForRunning = "YES"
            buildForProfiling = "YES"
            buildForArchiving = "YES"
            buildForAnalyzing = "YES">
            <BuildableReference
               BuildableIdentifier = "primary"
               BlueprintIdentifier = "BAC384091BA9F569005CFE20"
               BuildableName = "BitriseXcode7Sample.app"
               BlueprintName = "BitriseXcode7Sample"
               ReferencedContainer = "container:BitriseXcode7Sample.xcodeproj">
            </BuildableReference>
         </BuildActionEntry>
      </BuildActionEntries>
   </BuildAction>
   <TestAction
      buildConfiguration = "Debug"
      selectedDebuggerIdentifier = "Xcode.DebuggerFoundation.Debugger.LLDB"
      selectedLauncherIdentifier = "Xcode.DebuggerFoundation.Launcher.LLDB"
      shouldUseLaunchSchemeArgsEnv = "YES">
      <Testables>
      </Testables>
      <MacroExpansion>
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "BAC384091BA9F569005CFE20"
            BuildableName = "BitriseXcode7Sample.app"
            BlueprintName = "BitriseXcode7Sample"
            ReferencedContainer = "container:BitriseXcode7Sample.xcodeproj">
         </BuildableReference>
      </MacroExpansion>
      <AdditionalOptions>
      </AdditionalOptions>
   </TestAction>
   <LaunchAction
      buildConfiguration = "Debug"
      selectedDebuggerIdentifier = "Xcode.DebuggerFoundation.Debugger.LLDB"
      selectedLauncherIdentifier = "Xcode.DebuggerFoundation.Launcher.LLDB"
      launchStyle = "0"
      useCustomWorkingDirectory = "NO"
      ignoresPersistentStateOnLaunch = "NO"
      debugDocumentVersioning = "YES"
      debugServiceExtension = "internal"
      allowLocationSimulation = "YES">
      <BuildableProductRunnable
         runnableDebuggingMode = "0">
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "BAC384091BA9F569005CFE20"
            BuildableName = "BitriseXcode7Sample.app"
            BlueprintName = "BitriseXcode7Sample"
            ReferencedContainer = "container:BitriseXcode7Sample.xcodeproj">
         </BuildableReference>
      </BuildableProductRunnable>
      <AdditionalOptions>
      </AdditionalOptions>
   </LaunchAction>
   <ProfileAction
      buildConfiguration = "Release"
      shouldUseLaunchSchemeArgsEnv = "YES"
      savedToolIdentifier = ""
      useCustomWorkingDirectory = "NO"
      debugDocumentVersioning = "YES">
      <BuildableProductRunnable
         runnableDebuggingMode = "0">
         <BuildableReference
            BuildableIdentifier = "primary"
            BlueprintIdentifier = "BAC384091BA9F569005CFE20"
            BuildableName = "BitriseXcode7Sample.app"
            BlueprintName = "BitriseXcode7Sample"
            ReferencedContainer = "container:BitriseXcode7Sample.xcodeproj">
         </BuildableReference>
      </BuildableProductRunnable>
   </ProfileAction>
   <AnalyzeAction
      buildConfiguration = "Debug">
   </AnalyzeAction>
   <ArchiveAction
      buildConfiguration = "Release"
      revealArchiveInOrganizer = "YES">
   </ArchiveAction>
</Scheme>`

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
	content := `// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 46;
	objects = {

/* Begin PBXNativeTarget section */
		BAAFFED019EE788800F3AC91 /* SampleAppWithCocoapods */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = BAAFFEF719EE788800F3AC91 /* Build configuration list for PBXNativeTarget "SampleAppWithCocoapods" */;
			buildPhases = (
				BAAFFECD19EE788800F3AC91 /* Sources */,
				BAAFFECE19EE788800F3AC91 /* Frameworks */,
				BAAFFECF19EE788800F3AC91 /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
			);
			name = SampleAppWithCocoapods;
			productName = SampleAppWithCocoapods;
			productReference = BAAFFED119EE788800F3AC91 /* SampleAppWithCocoapods.app */;
			productType = "com.apple.product-type.application";
		};
		BAAFFEEC19EE788800F3AC91 /* SampleAppWithCocoapodsTests */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = BAAFFEFA19EE788800F3AC91 /* Build configuration list for PBXNativeTarget "SampleAppWithCocoapodsTests" */;
			buildPhases = (
				75ACE584234D974D15C5CAE9 /* Check Pods Manifest.lock */,
				BAAFFEE919EE788800F3AC91 /* Sources */,
				BAAFFEEA19EE788800F3AC91 /* Frameworks */,
				BAAFFEEB19EE788800F3AC91 /* Resources */,
				D0F06DBF2FED4262AA6DE7DB /* Copy Pods Resources */,
			);
			buildRules = (
			);
			dependencies = (
				BAAFFEEF19EE788800F3AC91 /* PBXTargetDependency */,
			);
			name = SampleAppWithCocoapodsTests;
			productName = SampleAppWithCocoapodsTests;
			productReference = BAAFFEED19EE788800F3AC91 /* SampleAppWithCocoapodsTests.xctest */;
			productType = "com.apple.product-type.bundle.unit-test";
		};
/* End PBXNativeTarget section */

/* Begin PBXVariantGroup section */
		BAAFFEE119EE788800F3AC91 /* Main.storyboard */ = {
			isa = PBXVariantGroup;
			children = (
				BAAFFEE219EE788800F3AC91 /* Base */,
			);
			name = Main.storyboard;
			sourceTree = "<group>";
		};
		BAAFFEE619EE788800F3AC91 /* LaunchScreen.xib */ = {
			isa = PBXVariantGroup;
			children = (
				BAAFFEE719EE788800F3AC91 /* Base */,
			);
			name = LaunchScreen.xib;
			sourceTree = "<group>";
		};
/* End PBXVariantGroup section */

	rootObject = BAAFFEC919EE788800F3AC91 /* Project object */;
}
`

	targets, err := pbxprojContentTartgets(content)
	require.NoError(t, err)
	require.Equal(t, 2, len(targets))
	require.Equal(t, "SampleAppWithCocoapods", targets[0])
	require.Equal(t, "SampleAppWithCocoapodsTests", targets[1])
}
