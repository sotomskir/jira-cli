// Copyright © 2019 Robert Sotomski <sotomski@gmail.com>
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

package logger

import (
	"github.com/fatih/color"
)

var Warn = color.New(color.FgHiYellow).PrintFunc()
var WarnLn = color.New(color.FgHiYellow).PrintlnFunc()
var WarnF = color.New(color.FgHiYellow).PrintfFunc()
var Info = color.New(color.FgHiBlue).PrintFunc()
var InfoLn = color.New(color.FgHiBlue).PrintlnFunc()
var InfoF = color.New(color.FgHiBlue).PrintfFunc()
var Error = color.New(color.FgHiRed).PrintFunc()
var ErrorLn = color.New(color.FgHiRed).PrintlnFunc()
var ErrorF = color.New(color.FgHiRed).PrintfFunc()
var Success = color.New(color.FgHiGreen).PrintFunc()
var SuccessLn = color.New(color.FgHiGreen).PrintlnFunc()
var SuccessF = color.New(color.FgHiGreen).PrintfFunc()
