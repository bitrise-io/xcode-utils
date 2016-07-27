package xcodeproj

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// Extensions
const (
	// XCWorkspaceExt ...
	XCWorkspaceExt = ".xcworkspace"
	// XCodeProjExt ...
	XCodeProjExt = ".xcodeproj"
	// XCSchemeExt ...
	XCSchemeExt = ".xcscheme"
)

// IsXCodeProj ...
func IsXCodeProj(pth string) bool {
	return strings.HasSuffix(pth, XCodeProjExt)
}

// IsXCWorkspace ...
func IsXCWorkspace(pth string) bool {
	return strings.HasSuffix(pth, XCWorkspaceExt)
}

// SchemeNameFromPath ...
func SchemeNameFromPath(schemePth string) string {
	basename := filepath.Base(schemePth)
	ext := filepath.Ext(schemePth)
	if ext != XCSchemeExt {
		return ""
	}
	return strings.TrimSuffix(basename, ext)
}

// SchemeFileContainsXCTestBuildAction ...
func SchemeFileContainsXCTestBuildAction(schemeFilePth string) (bool, error) {
	content, err := fileutil.ReadStringFromFile(schemeFilePth)
	if err != nil {
		return false, err
	}

	return schemeFileContentContainsXCTestBuildAction(content)
}

// ProjectSharedSchemeFilePaths ...
func ProjectSharedSchemeFilePaths(projectPth string) ([]string, error) {
	return sharedSchemeFilePaths(projectPth)
}

// WorkspaceSharedSchemeFilePaths ...
func WorkspaceSharedSchemeFilePaths(workspacePth string) ([]string, error) {
	workspaceSchemeFilePaths, err := sharedSchemeFilePaths(workspacePth)
	if err != nil {
		return []string{}, err
	}

	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		projectSchemeFilePaths, err := sharedSchemeFilePaths(project)
		if err != nil {
			return []string{}, err
		}
		workspaceSchemeFilePaths = append(workspaceSchemeFilePaths, projectSchemeFilePaths...)
	}

	sort.Strings(workspaceSchemeFilePaths)

	return workspaceSchemeFilePaths, nil
}

// ProjectSharedSchemes ...
func ProjectSharedSchemes(projectPth string) ([]string, error) {
	return sharedSchemes(projectPth)
}

// WorkspaceSharedSchemes ...
func WorkspaceSharedSchemes(workspacePth string) ([]string, error) {
	workspaceSchemes, err := sharedSchemes(workspacePth)
	if err != nil {
		return []string{}, err
	}

	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		projectSchemes, err := sharedSchemes(project)
		if err != nil {
			return []string{}, err
		}
		workspaceSchemes = append(workspaceSchemes, projectSchemes...)
	}

	sort.Strings(workspaceSchemes)

	return workspaceSchemes, nil
}

// ProjectUserSchemeFilePaths ...
func ProjectUserSchemeFilePaths(projectPth string) ([]string, error) {
	return userSchemeFilePaths(projectPth)
}

// WorkspaceUserSchemeFilePaths ...
func WorkspaceUserSchemeFilePaths(workspacePth string) ([]string, error) {
	workspaceSchemeFilePaths, err := userSchemeFilePaths(workspacePth)
	if err != nil {
		return []string{}, err
	}

	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		projectSchemeFilePaths, err := userSchemeFilePaths(project)
		if err != nil {
			return []string{}, err
		}
		workspaceSchemeFilePaths = append(workspaceSchemeFilePaths, projectSchemeFilePaths...)
	}

	sort.Strings(workspaceSchemeFilePaths)

	return workspaceSchemeFilePaths, nil
}

// ProjectUserSchemes ...
func ProjectUserSchemes(projectPth string) ([]string, error) {
	return userSchemes(projectPth)
}

// WorkspaceUserSchemes ...
func WorkspaceUserSchemes(workspacePth string) ([]string, error) {
	workspaceSchemes, err := userSchemes(workspacePth)
	if err != nil {
		return []string{}, err
	}

	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		projectSchemes, err := userSchemes(project)
		if err != nil {
			return []string{}, err
		}
		workspaceSchemes = append(workspaceSchemes, projectSchemes...)
	}

	sort.Strings(workspaceSchemes)

	return workspaceSchemes, nil
}

// ReCreateProjectUserSchemes ...
func ReCreateProjectUserSchemes(projectPth string) error {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("bitrise")
	if err != nil {
		return err
	}

	projectDir := filepath.Dir(projectPth)

	gemfileContent := `source 'https://rubygems.org'

gem 'xcodeproj'`

	gemfilePth := path.Join(tmpDir, "Gemfile")
	if err := fileutil.WriteStringToFile(gemfilePth, gemfileContent); err != nil {
		return err
	}

	envs := append(os.Environ(), "BUNDLE_GEMFILE="+gemfilePth)

	out, err := cmdex.NewCommand("bundle", "install").SetDir(projectDir).SetEnvs(envs).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return err
	}

	rubyScriptContent := `require 'xcodeproj'

project_path = ENV['project_path']

begin
  raise 'empty path' if project_path.empty?

  project = Xcodeproj::Project.open(project_path)
  project.recreate_user_schemes
  project.save
rescue => ex
  puts(ex.inspect.to_s)
  puts('--- Stack trace: ---')
  puts(ex.backtrace.to_s)
  exit(1)
end
`

	rubyScriptPth := path.Join(tmpDir, "recreate_user_schemes.rb")
	if err := fileutil.WriteStringToFile(rubyScriptPth, rubyScriptContent); err != nil {
		return err
	}

	projectBase := filepath.Base(projectPth)
	envs = append(os.Environ(), "project_path="+projectBase, "LC_ALL=en_US.UTF-8", "BUNDLE_GEMFILE="+gemfilePth)

	out, err = cmdex.NewCommand("bundle", "exec", "ruby", rubyScriptPth).SetDir(projectDir).SetEnvs(envs).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		if errorutil.IsExitStatusError(err) && out != "" {
			return errors.New(out)
		}
		return err
	}

	return nil
}

// ReCreateWorkspaceUserSchemes ...
func ReCreateWorkspaceUserSchemes(workspacePth string) error {
	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return err
	}

	for _, project := range projects {
		if err := ReCreateProjectUserSchemes(project); err != nil {
			return err
		}
	}

	return nil
}

// ProjectTargets ...
func ProjectTargets(projectPth string) ([]string, error) {
	pbxProjPth := filepath.Join(projectPth, "project.pbxproj")
	if exist, err := pathutil.IsPathExists(pbxProjPth); err != nil {
		return []string{}, err
	} else if !exist {
		return []string{}, fmt.Errorf("project.pbxproj does not exist at: %s", pbxProjPth)
	}

	content, err := fileutil.ReadStringFromFile(pbxProjPth)
	if err != nil {
		return []string{}, err
	}

	return pbxprojContentTartgets(content)
}

// WorkspaceTargets ...
func WorkspaceTargets(workspacePth string) ([]string, error) {
	projects, err := WorkspaceProjectReferences(workspacePth)
	if err != nil {
		return nil, err
	}

	targets := []string{}
	for _, project := range projects {
		projectTargets, err := ProjectTargets(project)
		if err != nil {
			return []string{}, err
		}
		targets = append(targets, projectTargets...)
	}

	sort.Strings(targets)

	return targets, nil
}

// WorkspaceProjectReferences ...
func WorkspaceProjectReferences(workspace string) ([]string, error) {
	projects := []string{}

	workspaceDir := filepath.Dir(workspace)

	xcworkspacedataPth := path.Join(workspace, "contents.xcworkspacedata")
	if exist, err := pathutil.IsPathExists(xcworkspacedataPth); err != nil {
		return []string{}, err
	} else if !exist {
		return []string{}, fmt.Errorf("contents.xcworkspacedata does not exist at: %s", xcworkspacedataPth)
	}

	xcworkspacedataStr, err := fileutil.ReadStringFromFile(xcworkspacedataPth)
	if err != nil {
		return []string{}, err
	}

	xcworkspacedataLines := strings.Split(xcworkspacedataStr, "\n")
	fileRefStart := false
	regexp := regexp.MustCompile(`location = "(.+):(.+).xcodeproj"`)

	for _, line := range xcworkspacedataLines {
		if strings.Contains(line, "<FileRef") {
			fileRefStart = true
			continue
		}

		if fileRefStart {
			fileRefStart = false
			matches := regexp.FindStringSubmatch(line)
			if len(matches) == 3 {
				projectName := matches[2]
				project := filepath.Join(workspaceDir, projectName+".xcodeproj")
				projects = append(projects, project)
			}
		}
	}

	sort.Strings(projects)

	return projects, nil
}

// ------------------------------
// Utility

func filesInDir(dir string) ([]string, error) {
	files := []string{}
	if err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	}); err != nil {
		return []string{}, err
	}
	return files, nil
}

func isUserSchemeFilePath(pth string) bool {
	regexpPattern := filepath.Join(".*[/]?xcuserdata", ".*[.]xcuserdatad", "xcschemes", ".+[.]xcscheme")
	regexp := regexp.MustCompile(regexpPattern)
	return (regexp.FindString(pth) != "")
}

func filterUserSchemeFilePaths(paths []string) []string {
	filteredPaths := []string{}
	for _, pth := range paths {
		if isUserSchemeFilePath(pth) {
			filteredPaths = append(filteredPaths, pth)
		}
	}

	sort.Strings(filteredPaths)

	return filteredPaths
}

func userSchemeFilePaths(projectOrWorkspacePth string) ([]string, error) {
	paths, err := filesInDir(projectOrWorkspacePth)
	if err != nil {
		return []string{}, err
	}
	return filterUserSchemeFilePaths(paths), nil
}

func userSchemes(projectOrWorkspacePth string) ([]string, error) {
	schemePaths, err := userSchemeFilePaths(projectOrWorkspacePth)
	if err != nil {
		return []string{}, err
	}

	schemes := []string{}
	for _, schemePth := range schemePaths {
		schemes = append(schemes, SchemeNameFromPath(schemePth))
	}

	sort.Strings(schemes)

	return schemes, nil
}

func isSharedSchemeFilePath(pth string) bool {
	regexpPattern := filepath.Join(".*[/]?xcshareddata", "xcschemes", ".+[.]xcscheme")
	regexp := regexp.MustCompile(regexpPattern)
	return (regexp.FindString(pth) != "")
}

func filterSharedSchemeFilePaths(paths []string) []string {
	filteredPaths := []string{}
	for _, pth := range paths {
		if isSharedSchemeFilePath(pth) {
			filteredPaths = append(filteredPaths, pth)
		}
	}

	sort.Strings(filteredPaths)

	return filteredPaths
}

func sharedSchemeFilePaths(projectOrWorkspacePth string) ([]string, error) {
	paths, err := filesInDir(projectOrWorkspacePth)
	if err != nil {
		return []string{}, err
	}
	return filterSharedSchemeFilePaths(paths), nil
}

func sharedSchemes(projectOrWorkspacePth string) ([]string, error) {
	schemePaths, err := sharedSchemeFilePaths(projectOrWorkspacePth)
	if err != nil {
		return []string{}, err
	}

	schemes := []string{}
	for _, schemePth := range schemePaths {
		schemes = append(schemes, SchemeNameFromPath(schemePth))
	}

	sort.Strings(schemes)

	return schemes, nil
}

func schemeFileContentContainsXCTestBuildAction(schemeFileContent string) (bool, error) {
	regexpPattern := `BuildableName = ".+.xctest"`
	regexp := regexp.MustCompile(regexpPattern)

	scanner := bufio.NewScanner(strings.NewReader(schemeFileContent))
	for scanner.Scan() {
		line := scanner.Text()
		if regexp.FindString(line) != "" {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func pbxprojContentTartgets(pbxprojContent string) ([]string, error) {
	nativeTargetSectionStart := "/* Begin PBXNativeTarget section */"
	nativeTargetSectionEnd := "/* End PBXNativeTarget section */"

	regexpPattern := `\s*name = (?P<name>.+);`
	regexp := regexp.MustCompile(regexpPattern)

	targets := []string{}
	isTargetSection := false

	scanner := bufio.NewScanner(strings.NewReader(pbxprojContent))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == nativeTargetSectionEnd {
			break
		}

		if strings.TrimSpace(line) == nativeTargetSectionStart {
			isTargetSection = true
		}

		if !isTargetSection {
			continue
		}

		if match := regexp.FindStringSubmatch(line); len(match) == 2 {
			target := match[1]
			targets = append(targets, target)
		}
	}

	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	sort.Strings(targets)

	return targets, nil
}
