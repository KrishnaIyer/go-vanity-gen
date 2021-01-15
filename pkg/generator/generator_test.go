// Copyright © 2021 Krishna Iyer Easwaran
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

package generator

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

var vanityCfg = []byte(`
host: go.example.com
paths:
  /mycoolproject:
    repo: https://github.com/user/mycoolproject
    packages:
      - pkg/package1
      - pkg/package2

  /myothercoolproject:
    repo: https://github.com/user/myothercoolproject

`)

var indexTemplate = `
<!DOCTYPE html>
<html>
<body>
<h1>Welcome to {{.Host}}</h1>
<ul>
{{range .Paths}}<li><a href="https://pkg.go.dev/{{.}}">{{.}}</a></li>{{end}}
</ul>
</body>
</html>
`

var packageTemplate = `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.Import}} {{.VCS}} {{.Repo}}">
<meta name="go-source" content="{{.Import}} {{.Display}}">
</head>
<body>
Nothing to see here folks!
</body>
</html>
`

func TestGenerate(t *testing.T) {
	a := assertions.New(t)
	ctx := context.Background()
	gen, err := New(ctx, vanityCfg)
	a.So(err, should.BeNil)
	index, err := gen.Index(ctx, indexTemplate)
	a.So(err, should.BeNil)
	fmt.Println(string(index))
	vanity, err := gen.Package(ctx, packageTemplate)
	a.So(err, should.BeNil)
	mcp := vanity.raw["/mycoolproject"]
	fmt.Println(string(mcp.content))
}
