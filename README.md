Flood-gogl
====================
Rewrite of https://github.com/dwrolvink/flood using OpenGL (go-gl & glfw).

This project uses the [go-gl wrapper package](https://github.com/dwrolvink/gogl/tree/main) that I'm writing as an exercise, based on this lecture series: https://www.youtube.com/watch?v=EJz71vpNhSU&list=PLDZujg-VgQlZUy1iCqBbe5faZLMkA3g2x&index=42 

https://github.com/dwrolvink/flood-gogl/assets/30291690/47318633-f58a-4ba4-81b2-dd930a943b4b


Support 
--------------------
I have built this package on linux with go 1.16, it's not guaranteed to work on any other system, but it probably will.
- Recording will not work on non linux systems as of yet!
- Recording requires ffmpeg to be installed.


Good to know
--------------------
The gl packages don't have module support (as of 2021/04), so we need to disable the module system to be able to use them.

We can set this setting permanently by executing:
``` bash
go env -w GO111MODULE=auto
```
This will not apply module mode when the code is located in $GOPATH/src, and no go.mod file is present.

> On Linux, if $GOPATH is empty. Packages are stored in /home/{user}/go/src/. Modules are stored under /home/{user}/go/pkg/mod

If you are working outside of the package folder, you can set GO111MODULE=off:
``` bash 
go env -w GO111MODULE=off
```

Alternatively, you can pass the setting as a first argument:
```bash
GO111MODULE=off go get github.com/go-gl/gl/v4.5-core/gl
GO111MODULE=off go get github.com/go-gl/glfw/v3.2/glfw

```

Installation
====================
```bash
# Disable Module mode
go env -w GO111MODULE=off

# Download packages
go get github.com/go-gl/gl/v4.1-core/gl
go get github.com/go-gl/glfw/v3.2/glfw
go get github.com/dwrolvink/gogl
go get gopkg.in/yaml.v3
go get golang.org/x/exp/constraints
```

Run
====================
```bash
go run ./*.go 
```

Build
====================
``` bash
# Compile
go build -v

# Run
./go_gl
```

Args
====================

--fps
-------------------
Change the speed. Note that fps of 50 is the max that will work with gifs.
```bash
go run ./*.go --fps 20
```

Combinations allowed:
-------------------
```bash
go run ./*.go --record 2 --set 1 --fps 30
```

Dev Notes
====================
VS Code Plugins
--------------------
When working with Go, the following plugins have been tremendously helpful:
- the `Go` plugin, by *Go Team at Google*
- the `Go Doc` plugin, by *Minhaz Ahmed Syrus*

Especially the last one, which allows you to hover over a method and see the description, is one that I use profusely.
If you are staring at the code, not sure what is happening, install the plugin, and hover over some methods.
