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

package config

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const (
	bazelIndent = 4
)

type bazelVarExporter interface {
	asBazel(exportedStringVariables, exportedStringListVariables) []bazelConstant
}

// Helpers for exporting cc configuration information to Bazel.
var (
	// Map containing toolchain variables that are independent of the
	// environment variables of the build.
	exportedStringListVars     = exportedStringListVariables{}
	exportedStringVars         = exportedStringVariables{}
	exportedStringListDictVars = exportedStringListDictVariables{}
)

// Ensure that string s has no invalid characters to be generated into the bzl file.
func validateCharacters(s string) string {
	for _, c := range []string{`\n`, `"`, `\`} {
		if strings.Contains(s, c) {
			panic(fmt.Errorf("%s contains illegal character %s", s, c))
		}
	}
	return s
}

type bazelConstant struct {
	variableName       string
	internalDefinition string
}

type exportedStringVariables map[string]string

func (m exportedStringVariables) Set(k string, v string) {
	m[k] = v
}

func bazelIndention(level int) string {
	return strings.Repeat(" ", level*bazelIndent)
}

func printBazelList(items []string, indentLevel int) string {
	list := make([]string, 0, len(items)+2)
	list = append(list, "[")
	innerIndent := bazelIndention(indentLevel + 1)
	for _, item := range items {
		list = append(list, fmt.Sprintf(`%s"%s",`, innerIndent, item))
	}
	list = append(list, bazelIndention(indentLevel)+"]")
	return strings.Join(list, "\n")
}

func (m exportedStringVariables) asBazel(stringScope exportedStringVariables, stringListScope exportedStringListVariables) []bazelConstant {
	ret := make([]bazelConstant, 0, len(m))
	for k, variableValue := range m {
		expandedVar := expandVar(variableValue, exportedStringVars, exportedStringListVars)
		if len(expandedVar) > 1 {
			panic(fmt.Errorf("%s expands to more than one string value: %s", variableValue, expandedVar))
		}
		ret = append(ret, bazelConstant{
			variableName:       k,
			internalDefinition: fmt.Sprintf(`"%s"`, validateCharacters(expandedVar[0])),
		})
	}
	return ret
}

// Convenience function to declare a static variable and export it to Bazel's cc_toolchain.
func exportStringStaticVariable(name string, value string) {
	pctx.StaticVariable(name, value)
	exportedStringVars.Set(name, value)
}

type exportedStringListVariables map[string][]string

func (m exportedStringListVariables) Set(k string, v []string) {
	m[k] = v
}

func (m exportedStringListVariables) asBazel(stringScope exportedStringVariables, stringListScope exportedStringListVariables) []bazelConstant {
	ret := make([]bazelConstant, 0, len(m))
	// For each exported variable, recursively expand elements in the variableValue
	// list to ensure that interpolated variables are expanded according to their values
	// in the variable scope.
	for k, variableValue := range m {
		var expandedVars []string
		for _, v := range variableValue {
			expandedVars = append(expandedVars, expandVar(v, stringScope, stringListScope)...)
		}
		// Assign the list as a bzl-private variable; this variable will be exported
		// out through a constants struct later.
		ret = append(ret, bazelConstant{
			variableName:       k,
			internalDefinition: printBazelList(expandedVars, 0),
		})
	}
	return ret
}

// Convenience function to declare a static variable and export it to Bazel's cc_toolchain.
func exportStringListStaticVariable(name string, value []string) {
	pctx.StaticVariable(name, strings.Join(value, " "))
	exportedStringListVars.Set(name, value)
}

type exportedStringListDictVariables map[string]map[string][]string

func (m exportedStringListDictVariables) Set(k string, v map[string][]string) {
	m[k] = v
}

func printBazelStringListDict(dict map[string][]string) string {
	bazelDict := make([]string, 0, len(dict)+2)
	bazelDict = append(bazelDict, "{")
	for k, v := range dict {
		bazelDict = append(bazelDict,
			fmt.Sprintf(`%s"%s": %s,`, bazelIndention(1), k, printBazelList(v, 1)))
	}
	bazelDict = append(bazelDict, "}")
	return strings.Join(bazelDict, "\n")
}

// Since dictionaries are not supported in Ninja, we do not expand variables for dictionaries
func (m exportedStringListDictVariables) asBazel(_ exportedStringVariables, _ exportedStringListVariables) []bazelConstant {
	ret := make([]bazelConstant, 0, len(m))
	for k, dict := range m {
		ret = append(ret, bazelConstant{
			variableName:       k,
			internalDefinition: printBazelStringListDict(dict),
		})
	}
	return ret
}

// BazelCcToolchainVars generates bzl file content containing variables for
// Bazel's cc_toolchain configuration.
func BazelCcToolchainVars() string {
	return bazelToolchainVars(
		exportedStringListDictVars,
		exportedStringListVars,
		exportedStringVars)
}

func bazelToolchainVars(vars ...bazelVarExporter) string {
	ret := "# GENERATED FOR BAZEL FROM SOONG. DO NOT EDIT.\n\n"

	results := []bazelConstant{}
	for _, v := range vars {
		results = append(results, v.asBazel(exportedStringVars, exportedStringListVars)...)
	}

	sort.Slice(results, func(i, j int) bool { return results[i].variableName < results[j].variableName })

	definitions := make([]string, 0, len(results))
	constants := make([]string, 0, len(results))
	for _, b := range results {
		definitions = append(definitions,
			fmt.Sprintf("_%s = %s", b.variableName, b.internalDefinition))
		constants = append(constants,
			fmt.Sprintf("%[1]s%[2]s = _%[2]s,", bazelIndention(1), b.variableName))
	}

	// Build the exported constants struct.
	ret += strings.Join(definitions, "\n\n")
	ret += "\n\n"
	ret += "constants = struct(\n"
	ret += strings.Join(constants, "\n")
	ret += "\n)"

	return ret
}

// expandVar recursively expand interpolated variables in the exportedVars scope.
//
// We're using a string slice to track the seen variables to avoid
// stackoverflow errors with infinite recursion. it's simpler to use a
// string slice than to handle a pass-by-referenced map, which would make it
// quite complex to track depth-first interpolations. It's also unlikely the
// interpolation stacks are deep (n > 1).
func expandVar(toExpand string, stringScope exportedStringVariables, stringListScope exportedStringListVariables) []string {
	// e.g. "${ExternalCflags}"
	r := regexp.MustCompile(`\${([a-zA-Z0-9_]+)}`)

	// Internal recursive function.
	var expandVarInternal func(string, map[string]bool) []string
	expandVarInternal = func(toExpand string, seenVars map[string]bool) []string {
		var ret []string
		for _, v := range strings.Split(toExpand, " ") {
			matches := r.FindStringSubmatch(v)
			if len(matches) == 0 {
				return []string{v}
			}

			if len(matches) != 2 {
				panic(fmt.Errorf(
					"Expected to only match 1 subexpression in %s, got %d",
					v,
					len(matches)-1))
			}

			// Index 1 of FindStringSubmatch contains the subexpression match
			// (variable name) of the capture group.
			variable := matches[1]
			// toExpand contains a variable.
			if _, ok := seenVars[variable]; ok {
				panic(fmt.Errorf(
					"Unbounded recursive interpolation of variable: %s", variable))
			}
			// A map is passed-by-reference. Create a new map for
			// this scope to prevent variables seen in one depth-first expansion
			// to be also treated as "seen" in other depth-first traversals.
			newSeenVars := map[string]bool{}
			for k := range seenVars {
				newSeenVars[k] = true
			}
			newSeenVars[variable] = true
			if unexpandedVars, ok := stringListScope[variable]; ok {
				for _, unexpandedVar := range unexpandedVars {
					ret = append(ret, expandVarInternal(unexpandedVar, newSeenVars)...)
				}
			} else if unexpandedVar, ok := stringScope[variable]; ok {
				ret = append(ret, expandVarInternal(unexpandedVar, newSeenVars)...)
			}
		}
		return ret
	}

	return expandVarInternal(toExpand, map[string]bool{})
}
