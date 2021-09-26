// Copyright 2021 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cc

import (
	"fmt"
	"path/filepath"
	"strings"

	"android/soong/android"
	"android/soong/bazel"

	"github.com/google/blueprint/proptools"
)

// staticOrSharedAttributes are the Bazel-ified versions of StaticOrSharedProperties --
// properties which apply to either the shared or static version of a cc_library module.
type staticOrSharedAttributes struct {
	Srcs    bazel.LabelListAttribute
	Srcs_c  bazel.LabelListAttribute
	Srcs_as bazel.LabelListAttribute
	Copts   bazel.StringListAttribute

	Static_deps        bazel.LabelListAttribute
	Dynamic_deps       bazel.LabelListAttribute
	Whole_archive_deps bazel.LabelListAttribute

	System_dynamic_deps bazel.LabelListAttribute
}

func groupSrcsByExtension(ctx android.TopDownMutatorContext, srcs bazel.LabelListAttribute) (cppSrcs, cSrcs, asSrcs bazel.LabelListAttribute) {
	// Branch srcs into three language-specific groups.
	// C++ is the "catch-all" group, and comprises generated sources because we don't
	// know the language of these sources until the genrule is executed.
	// TODO(b/190006308): Handle language detection of sources in a Bazel rule.
	isCSrcOrFilegroup := func(s string) bool {
		return strings.HasSuffix(s, ".c") || strings.HasSuffix(s, "_c_srcs")
	}

	isAsmSrcOrFilegroup := func(s string) bool {
		return strings.HasSuffix(s, ".S") || strings.HasSuffix(s, ".s") || strings.HasSuffix(s, "_as_srcs")
	}

	// Check that a module is a filegroup type named <label>.
	isFilegroupNamed := func(m android.Module, fullLabel string) bool {
		if ctx.OtherModuleType(m) != "filegroup" {
			return false
		}
		labelParts := strings.Split(fullLabel, ":")
		if len(labelParts) > 2 {
			// There should not be more than one colon in a label.
			panic(fmt.Errorf("%s is not a valid Bazel label for a filegroup", fullLabel))
		} else {
			return m.Name() == labelParts[len(labelParts)-1]
		}
	}

	// Convert the filegroup dependencies into the extension-specific filegroups
	// filtered in the filegroup.bzl macro.
	cppFilegroup := func(label string) string {
		m, exists := ctx.ModuleFromName(label)
		if exists {
			aModule, _ := m.(android.Module)
			if isFilegroupNamed(aModule, label) {
				label = label + "_cpp_srcs"
			}
		}
		return label
	}
	cFilegroup := func(label string) string {
		m, exists := ctx.ModuleFromName(label)
		if exists {
			aModule, _ := m.(android.Module)
			if isFilegroupNamed(aModule, label) {
				label = label + "_c_srcs"
			}
		}
		return label
	}
	asFilegroup := func(label string) string {
		m, exists := ctx.ModuleFromName(label)
		if exists {
			aModule, _ := m.(android.Module)
			if isFilegroupNamed(aModule, label) {
				label = label + "_as_srcs"
			}
		}
		return label
	}

	cSrcs = bazel.MapLabelListAttribute(srcs, cFilegroup)
	cSrcs = bazel.FilterLabelListAttribute(cSrcs, isCSrcOrFilegroup)

	asSrcs = bazel.MapLabelListAttribute(srcs, asFilegroup)
	asSrcs = bazel.FilterLabelListAttribute(asSrcs, isAsmSrcOrFilegroup)

	cppSrcs = bazel.MapLabelListAttribute(srcs, cppFilegroup)
	cppSrcs = bazel.SubtractBazelLabelListAttribute(cppSrcs, cSrcs)
	cppSrcs = bazel.SubtractBazelLabelListAttribute(cppSrcs, asSrcs)
	return
}

// bp2buildParseSharedProps returns the attributes for the shared variant of a cc_library.
func bp2BuildParseSharedProps(ctx android.TopDownMutatorContext, module *Module) staticOrSharedAttributes {
	lib, ok := module.compiler.(*libraryDecorator)
	if !ok {
		return staticOrSharedAttributes{}
	}

	return bp2buildParseStaticOrSharedProps(ctx, module, lib, false)
}

// bp2buildParseStaticProps returns the attributes for the static variant of a cc_library.
func bp2BuildParseStaticProps(ctx android.TopDownMutatorContext, module *Module) staticOrSharedAttributes {
	lib, ok := module.compiler.(*libraryDecorator)
	if !ok {
		return staticOrSharedAttributes{}
	}

	return bp2buildParseStaticOrSharedProps(ctx, module, lib, true)
}

func bp2buildParseStaticOrSharedProps(ctx android.TopDownMutatorContext, module *Module, lib *libraryDecorator, isStatic bool) staticOrSharedAttributes {
	attrs := staticOrSharedAttributes{}

	setAttrs := func(axis bazel.ConfigurationAxis, config string, props StaticOrSharedProperties) {
		attrs.Copts.SetSelectValue(axis, config, props.Cflags)
		attrs.Srcs.SetSelectValue(axis, config, android.BazelLabelForModuleSrc(ctx, props.Srcs))
		attrs.Static_deps.SetSelectValue(axis, config, android.BazelLabelForModuleDeps(ctx, props.Static_libs))
		attrs.Dynamic_deps.SetSelectValue(axis, config, android.BazelLabelForModuleDeps(ctx, props.Shared_libs))
		attrs.Whole_archive_deps.SetSelectValue(axis, config, android.BazelLabelForModuleWholeDeps(ctx, props.Whole_static_libs))
		attrs.System_dynamic_deps.SetSelectValue(axis, config, android.BazelLabelForModuleDeps(ctx, props.System_shared_libs))
	}
	// system_dynamic_deps distinguishes between nil/empty list behavior:
	//    nil -> use default values
	//    empty list -> no values specified
	attrs.System_dynamic_deps.ForceSpecifyEmptyList = true

	if isStatic {
		for axis, configToProps := range module.GetArchVariantProperties(ctx, &StaticProperties{}) {
			for config, props := range configToProps {
				if staticOrSharedProps, ok := props.(*StaticProperties); ok {
					setAttrs(axis, config, staticOrSharedProps.Static)
				}
			}
		}
	} else {
		for axis, configToProps := range module.GetArchVariantProperties(ctx, &SharedProperties{}) {
			for config, props := range configToProps {
				if staticOrSharedProps, ok := props.(*SharedProperties); ok {
					setAttrs(axis, config, staticOrSharedProps.Shared)
				}
			}
		}
	}

	cppSrcs, cSrcs, asSrcs := groupSrcsByExtension(ctx, attrs.Srcs)
	attrs.Srcs = cppSrcs
	attrs.Srcs_c = cSrcs
	attrs.Srcs_as = asSrcs

	return attrs
}

// Convenience struct to hold all attributes parsed from prebuilt properties.
type prebuiltAttributes struct {
	Src bazel.LabelAttribute
}

func Bp2BuildParsePrebuiltLibraryProps(ctx android.TopDownMutatorContext, module *Module) prebuiltAttributes {
	var srcLabelAttribute bazel.LabelAttribute

	for axis, configToProps := range module.GetArchVariantProperties(ctx, &prebuiltLinkerProperties{}) {
		for config, props := range configToProps {
			if prebuiltLinkerProperties, ok := props.(*prebuiltLinkerProperties); ok {
				if len(prebuiltLinkerProperties.Srcs) > 1 {
					ctx.ModuleErrorf("Bp2BuildParsePrebuiltLibraryProps: Expected at most once source file for %s %s\n", axis, config)
					continue
				} else if len(prebuiltLinkerProperties.Srcs) == 0 {
					continue
				}
				src := android.BazelLabelForModuleSrcSingle(ctx, prebuiltLinkerProperties.Srcs[0])
				srcLabelAttribute.SetSelectValue(axis, config, src)
			}
		}
	}

	return prebuiltAttributes{
		Src: srcLabelAttribute,
	}
}

// Convenience struct to hold all attributes parsed from compiler properties.
type compilerAttributes struct {
	// Options for all languages
	copts bazel.StringListAttribute
	// Assembly options and sources
	asFlags bazel.StringListAttribute
	asSrcs  bazel.LabelListAttribute
	// C options and sources
	conlyFlags bazel.StringListAttribute
	cSrcs      bazel.LabelListAttribute
	// C++ options and sources
	cppFlags bazel.StringListAttribute
	srcs     bazel.LabelListAttribute

	rtti bazel.BoolAttribute
}

// bp2BuildParseCompilerProps returns copts, srcs and hdrs and other attributes.
func bp2BuildParseCompilerProps(ctx android.TopDownMutatorContext, module *Module) compilerAttributes {
	var srcs bazel.LabelListAttribute
	var copts bazel.StringListAttribute
	var asFlags bazel.StringListAttribute
	var conlyFlags bazel.StringListAttribute
	var cppFlags bazel.StringListAttribute
	var rtti bazel.BoolAttribute

	// Creates the -I flags for a directory, while making the directory relative
	// to the exec root for Bazel to work.
	includeFlags := func(dir string) []string {
		// filepath.Join canonicalizes the path, i.e. it takes care of . or .. elements.
		moduleDirRootedPath := filepath.Join(ctx.ModuleDir(), dir)
		return []string{
			"-I" + moduleDirRootedPath,
			// Include the bindir-rooted path (using make variable substitution). This most
			// closely matches Bazel's native include path handling, which allows for dependency
			// on generated headers in these directories.
			// TODO(b/188084383): Handle local include directories in Bazel.
			"-I$(BINDIR)/" + moduleDirRootedPath,
		}
	}

	// Parse the list of module-relative include directories (-I).
	parseLocalIncludeDirs := func(baseCompilerProps *BaseCompilerProperties) []string {
		// include_dirs are root-relative, not module-relative.
		includeDirs := bp2BuildMakePathsRelativeToModule(ctx, baseCompilerProps.Include_dirs)
		return append(includeDirs, baseCompilerProps.Local_include_dirs...)
	}

	parseCommandLineFlags := func(soongFlags []string) []string {
		var result []string
		for _, flag := range soongFlags {
			// Soong's cflags can contain spaces, like `-include header.h`. For
			// Bazel's copts, split them up to be compatible with the
			// no_copts_tokenization feature.
			result = append(result, strings.Split(flag, " ")...)
		}
		return result
	}

	// Parse srcs from an arch or OS's props value.
	parseSrcs := func(baseCompilerProps *BaseCompilerProperties) bazel.LabelList {
		// Add srcs-like dependencies such as generated files.
		// First create a LabelList containing these dependencies, then merge the values with srcs.
		generatedHdrsAndSrcs := baseCompilerProps.Generated_headers
		generatedHdrsAndSrcs = append(generatedHdrsAndSrcs, baseCompilerProps.Generated_sources...)
		generatedHdrsAndSrcsLabelList := android.BazelLabelForModuleDeps(ctx, generatedHdrsAndSrcs)

		allSrcsLabelList := android.BazelLabelForModuleSrcExcludes(ctx, baseCompilerProps.Srcs, baseCompilerProps.Exclude_srcs)
		return bazel.AppendBazelLabelLists(allSrcsLabelList, generatedHdrsAndSrcsLabelList)
	}

	archVariantCompilerProps := module.GetArchVariantProperties(ctx, &BaseCompilerProperties{})
	for axis, configToProps := range archVariantCompilerProps {
		for config, props := range configToProps {
			if baseCompilerProps, ok := props.(*BaseCompilerProperties); ok {
				// If there's arch specific srcs or exclude_srcs, generate a select entry for it.
				// TODO(b/186153868): do this for OS specific srcs and exclude_srcs too.
				if len(baseCompilerProps.Srcs) > 0 || len(baseCompilerProps.Exclude_srcs) > 0 {
					srcsList := parseSrcs(baseCompilerProps)
					srcs.SetSelectValue(axis, config, srcsList)
				}

				archVariantCopts := parseCommandLineFlags(baseCompilerProps.Cflags)
				archVariantAsflags := parseCommandLineFlags(baseCompilerProps.Asflags)
				for _, dir := range parseLocalIncludeDirs(baseCompilerProps) {
					archVariantCopts = append(archVariantCopts, includeFlags(dir)...)
					archVariantAsflags = append(archVariantAsflags, includeFlags(dir)...)
				}

				if axis == bazel.NoConfigAxis {
					if includeBuildDirectory(baseCompilerProps.Include_build_directory) {
						flags := includeFlags(".")
						archVariantCopts = append(archVariantCopts, flags...)
						archVariantAsflags = append(archVariantAsflags, flags...)
					}
				}

				copts.SetSelectValue(axis, config, archVariantCopts)
				asFlags.SetSelectValue(axis, config, archVariantAsflags)
				conlyFlags.SetSelectValue(axis, config, parseCommandLineFlags(baseCompilerProps.Conlyflags))
				cppFlags.SetSelectValue(axis, config, parseCommandLineFlags(baseCompilerProps.Cppflags))
				rtti.SetSelectValue(axis, config, baseCompilerProps.Rtti)
			}
		}
	}

	srcs.ResolveExcludes()

	productVarPropNameToAttribute := map[string]*bazel.StringListAttribute{
		"Cflags":   &copts,
		"Asflags":  &asFlags,
		"CppFlags": &cppFlags,
	}
	productVariableProps := android.ProductVariableProperties(ctx)
	for propName, attr := range productVarPropNameToAttribute {
		if props, exists := productVariableProps[propName]; exists {
			for _, prop := range props {
				flags, ok := prop.Property.([]string)
				if !ok {
					ctx.ModuleErrorf("Could not convert product variable %s property", proptools.PropertyNameForField(propName))
				}
				newFlags, _ := bazel.TryVariableSubstitutions(flags, prop.ProductConfigVariable)
				attr.SetSelectValue(bazel.ProductVariableConfigurationAxis(prop.FullConfig), prop.FullConfig, newFlags)
			}
		}
	}

	srcs, cSrcs, asSrcs := groupSrcsByExtension(ctx, srcs)

	return compilerAttributes{
		copts:      copts,
		srcs:       srcs,
		asFlags:    asFlags,
		asSrcs:     asSrcs,
		cSrcs:      cSrcs,
		conlyFlags: conlyFlags,
		cppFlags:   cppFlags,
		rtti:       rtti,
	}
}

// Convenience struct to hold all attributes parsed from linker properties.
type linkerAttributes struct {
	deps                          bazel.LabelListAttribute
	dynamicDeps                   bazel.LabelListAttribute
	systemDynamicDeps             bazel.LabelListAttribute
	wholeArchiveDeps              bazel.LabelListAttribute
	exportedDeps                  bazel.LabelListAttribute
	useLibcrt                     bazel.BoolAttribute
	linkopts                      bazel.StringListAttribute
	versionScript                 bazel.LabelAttribute
	stripKeepSymbols              bazel.BoolAttribute
	stripKeepSymbolsAndDebugFrame bazel.BoolAttribute
	stripKeepSymbolsList          bazel.StringListAttribute
	stripAll                      bazel.BoolAttribute
	stripNone                     bazel.BoolAttribute
}

// FIXME(b/187655838): Use the existing linkerFlags() function instead of duplicating logic here
func getBp2BuildLinkerFlags(linkerProperties *BaseLinkerProperties) []string {
	flags := linkerProperties.Ldflags
	if !BoolDefault(linkerProperties.Pack_relocations, true) {
		flags = append(flags, "-Wl,--pack-dyn-relocs=none")
	}
	return flags
}

// bp2BuildParseLinkerProps parses the linker properties of a module, including
// configurable attribute values.
func bp2BuildParseLinkerProps(ctx android.TopDownMutatorContext, module *Module) linkerAttributes {
	var headerDeps bazel.LabelListAttribute
	var staticDeps bazel.LabelListAttribute
	var exportedDeps bazel.LabelListAttribute
	var dynamicDeps bazel.LabelListAttribute
	var wholeArchiveDeps bazel.LabelListAttribute
	systemSharedDeps := bazel.LabelListAttribute{ForceSpecifyEmptyList: true}
	var linkopts bazel.StringListAttribute
	var versionScript bazel.LabelAttribute
	var useLibcrt bazel.BoolAttribute

	var stripKeepSymbols bazel.BoolAttribute
	var stripKeepSymbolsAndDebugFrame bazel.BoolAttribute
	var stripKeepSymbolsList bazel.StringListAttribute
	var stripAll bazel.BoolAttribute
	var stripNone bazel.BoolAttribute

	for axis, configToProps := range module.GetArchVariantProperties(ctx, &StripProperties{}) {
		for config, props := range configToProps {
			if stripProperties, ok := props.(*StripProperties); ok {
				stripKeepSymbols.SetSelectValue(axis, config, stripProperties.Strip.Keep_symbols)
				stripKeepSymbolsList.SetSelectValue(axis, config, stripProperties.Strip.Keep_symbols_list)
				stripKeepSymbolsAndDebugFrame.SetSelectValue(axis, config, stripProperties.Strip.Keep_symbols_and_debug_frame)
				stripAll.SetSelectValue(axis, config, stripProperties.Strip.All)
				stripNone.SetSelectValue(axis, config, stripProperties.Strip.None)
			}
		}
	}

	for axis, configToProps := range module.GetArchVariantProperties(ctx, &BaseLinkerProperties{}) {
		for config, props := range configToProps {
			if baseLinkerProps, ok := props.(*BaseLinkerProperties); ok {
				// Excludes to parallel Soong:
				// https://cs.android.com/android/platform/superproject/+/master:build/soong/cc/linker.go;l=247-249;drc=088b53577dde6e40085ffd737a1ae96ad82fc4b0
				staticLibs := android.FirstUniqueStrings(baseLinkerProps.Static_libs)
				staticDeps.SetSelectValue(axis, config, android.BazelLabelForModuleDepsExcludes(ctx, staticLibs, baseLinkerProps.Exclude_static_libs))
				wholeArchiveLibs := android.FirstUniqueStrings(baseLinkerProps.Whole_static_libs)
				wholeArchiveDeps.SetSelectValue(axis, config, android.BazelLabelForModuleWholeDepsExcludes(ctx, wholeArchiveLibs, baseLinkerProps.Exclude_static_libs))

				systemSharedLibs := baseLinkerProps.System_shared_libs
				// systemSharedLibs distinguishes between nil/empty list behavior:
				//    nil -> use default values
				//    empty list -> no values specified
				if len(systemSharedLibs) > 0 {
					systemSharedLibs = android.FirstUniqueStrings(systemSharedLibs)
				}
				systemSharedDeps.SetSelectValue(axis, config, android.BazelLabelForModuleDeps(ctx, systemSharedLibs))

				sharedLibs := android.FirstUniqueStrings(baseLinkerProps.Shared_libs)
				dynamicDeps.SetSelectValue(axis, config, android.BazelLabelForModuleDepsExcludes(ctx, sharedLibs, baseLinkerProps.Exclude_shared_libs))

				headerLibs := android.FirstUniqueStrings(baseLinkerProps.Header_libs)
				headerDeps.SetSelectValue(axis, config, android.BazelLabelForModuleDeps(ctx, headerLibs))
				exportedLibs := android.FirstUniqueStrings(baseLinkerProps.Export_header_lib_headers)
				exportedDeps.SetSelectValue(axis, config, android.BazelLabelForModuleDeps(ctx, exportedLibs))

				linkopts.SetSelectValue(axis, config, getBp2BuildLinkerFlags(baseLinkerProps))
				if baseLinkerProps.Version_script != nil {
					versionScript.SetSelectValue(axis, config, android.BazelLabelForModuleSrcSingle(ctx, *baseLinkerProps.Version_script))
				}
				useLibcrt.SetSelectValue(axis, config, baseLinkerProps.libCrt())
			}
		}
	}

	type productVarDep struct {
		// the name of the corresponding excludes field, if one exists
		excludesField string
		// reference to the bazel attribute that should be set for the given product variable config
		attribute *bazel.LabelListAttribute

		depResolutionFunc func(ctx android.BazelConversionPathContext, modules, excludes []string) bazel.LabelList
	}

	productVarToDepFields := map[string]productVarDep{
		// product variables do not support exclude_shared_libs
		"Shared_libs":       productVarDep{attribute: &dynamicDeps, depResolutionFunc: android.BazelLabelForModuleDepsExcludes},
		"Static_libs":       productVarDep{"Exclude_static_libs", &staticDeps, android.BazelLabelForModuleDepsExcludes},
		"Whole_static_libs": productVarDep{"Exclude_static_libs", &wholeArchiveDeps, android.BazelLabelForModuleWholeDepsExcludes},
	}

	productVariableProps := android.ProductVariableProperties(ctx)
	for name, dep := range productVarToDepFields {
		props, exists := productVariableProps[name]
		excludeProps, excludesExists := productVariableProps[dep.excludesField]
		// if neither an include or excludes property exists, then skip it
		if !exists && !excludesExists {
			continue
		}
		// collect all the configurations that an include or exclude property exists for.
		// we want to iterate all configurations rather than either the include or exclude because for a
		// particular configuration we may have only and include or only an exclude to handle
		configs := make(map[string]bool, len(props)+len(excludeProps))
		for config := range props {
			configs[config] = true
		}
		for config := range excludeProps {
			configs[config] = true
		}

		for config := range configs {
			prop, includesExists := props[config]
			excludesProp, excludesExists := excludeProps[config]
			var includes, excludes []string
			var ok bool
			// if there was no includes/excludes property, casting fails and that's expected
			if includes, ok = prop.Property.([]string); includesExists && !ok {
				ctx.ModuleErrorf("Could not convert product variable %s property", name)
			}
			if excludes, ok = excludesProp.Property.([]string); excludesExists && !ok {
				ctx.ModuleErrorf("Could not convert product variable %s property", dep.excludesField)
			}

			dep.attribute.SetSelectValue(bazel.ProductVariableConfigurationAxis(config), config, dep.depResolutionFunc(ctx, android.FirstUniqueStrings(includes), excludes))
		}
	}

	staticDeps.ResolveExcludes()
	dynamicDeps.ResolveExcludes()
	wholeArchiveDeps.ResolveExcludes()

	headerDeps.Append(staticDeps)

	return linkerAttributes{
		deps:              headerDeps,
		exportedDeps:      exportedDeps,
		dynamicDeps:       dynamicDeps,
		systemDynamicDeps: systemSharedDeps,
		wholeArchiveDeps:  wholeArchiveDeps,
		linkopts:          linkopts,
		useLibcrt:         useLibcrt,
		versionScript:     versionScript,

		// Strip properties
		stripKeepSymbols:              stripKeepSymbols,
		stripKeepSymbolsAndDebugFrame: stripKeepSymbolsAndDebugFrame,
		stripKeepSymbolsList:          stripKeepSymbolsList,
		stripAll:                      stripAll,
		stripNone:                     stripNone,
	}
}

// Relativize a list of root-relative paths with respect to the module's
// directory.
//
// include_dirs Soong prop are root-relative (b/183742505), but
// local_include_dirs, export_include_dirs and export_system_include_dirs are
// module dir relative. This function makes a list of paths entirely module dir
// relative.
//
// For the `include` attribute, Bazel wants the paths to be relative to the
// module.
func bp2BuildMakePathsRelativeToModule(ctx android.BazelConversionPathContext, paths []string) []string {
	var relativePaths []string
	for _, path := range paths {
		// Semantics of filepath.Rel: join(ModuleDir, rel(ModuleDir, path)) == path
		relativePath, err := filepath.Rel(ctx.ModuleDir(), path)
		if err != nil {
			panic(err)
		}
		relativePaths = append(relativePaths, relativePath)
	}
	return relativePaths
}

func bp2BuildParseExportedIncludes(ctx android.TopDownMutatorContext, module *Module) bazel.StringListAttribute {
	libraryDecorator := module.linker.(*libraryDecorator)
	return bp2BuildParseExportedIncludesHelper(ctx, module, libraryDecorator)
}

func Bp2BuildParseExportedIncludesForPrebuiltLibrary(ctx android.TopDownMutatorContext, module *Module) bazel.StringListAttribute {
	prebuiltLibraryLinker := module.linker.(*prebuiltLibraryLinker)
	libraryDecorator := prebuiltLibraryLinker.libraryDecorator
	return bp2BuildParseExportedIncludesHelper(ctx, module, libraryDecorator)
}

// bp2BuildParseExportedIncludes creates a string list attribute contains the
// exported included directories of a module.
func bp2BuildParseExportedIncludesHelper(ctx android.TopDownMutatorContext, module *Module, libraryDecorator *libraryDecorator) bazel.StringListAttribute {
	// Export_system_include_dirs and export_include_dirs are already module dir
	// relative, so they don't need to be relativized like include_dirs, which
	// are root-relative.
	includeDirs := libraryDecorator.flagExporter.Properties.Export_system_include_dirs
	includeDirs = append(includeDirs, libraryDecorator.flagExporter.Properties.Export_include_dirs...)
	var includeDirsAttribute bazel.StringListAttribute

	getVariantIncludeDirs := func(includeDirs []string, flagExporterProperties *FlagExporterProperties, subtract bool) []string {
		variantIncludeDirs := flagExporterProperties.Export_system_include_dirs
		variantIncludeDirs = append(variantIncludeDirs, flagExporterProperties.Export_include_dirs...)

		if subtract {
			// To avoid duplicate includes when base includes + arch includes are combined
			// TODO: Add something similar to ResolveExcludes() in bazel/properties.go
			variantIncludeDirs = bazel.SubtractStrings(variantIncludeDirs, includeDirs)
		}
		return variantIncludeDirs
	}

	for axis, configToProps := range module.GetArchVariantProperties(ctx, &FlagExporterProperties{}) {
		for config, props := range configToProps {
			if flagExporterProperties, ok := props.(*FlagExporterProperties); ok {
				archVariantIncludeDirs := getVariantIncludeDirs(includeDirs, flagExporterProperties, axis != bazel.NoConfigAxis)
				if len(archVariantIncludeDirs) > 0 {
					includeDirsAttribute.SetSelectValue(axis, config, archVariantIncludeDirs)
				}
			}
		}
	}

	return includeDirsAttribute
}
