// Copyright © 2022 Meroxa, Inc.
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

package example

type GlobalConfig struct {
	GlobalString string
}

type SourceConfig struct {
	GlobalConfig
	MyInt  int
	Nested struct {
		MyFloat32 float32
		MyFloat64 float64
	}
}

type DestinationConfig struct {
	GlobalConfig
	MyBool bool
}

// Specs contains the specifications.
//spec:summary This is a test summary.
//spec:description This is a test description. It can be in multiple lines long
// if needed.
//spec:author Example Inc.
//spec:version v0.1.1
type Specs struct {
	//spec:sourceParams
	SourceConfig
	//spec:destinationParams
	DestinationConfig
}
