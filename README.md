# gapt ![wip](https://img.shields.io/badge/-work%20in%20progress-red) ![draft](https://img.shields.io/badge/-draft-red)

**gapt** is a *go asset packaging tool* to bundle and preprocess resources before compilation, e.g.
by using `go generate`. It is more comparable to *aapt*, the android asset packaging tool. 

At first, this sounds very familiar to other bundle tools like [vfsgen](https://github.com/shurcooL/vfsgen),
  [go.rice](https://github.com/GeertJohan/go.rice), [statik](https://github.com/rakyll/statik),
   [packr](https://github.com/gobuffalo/packr) or [go-bindata](https://github.com/gnoso/go-bindata) but
*gapt* is more dynamic and allows efficient customization and replacements of asset files at runtime
using [fsnotify](https://fsnotify.org/) mechanism, both while developing or in release mode. 
This allows to overload embedded files dynamically at runtime, either by directly using the 
original files from your module while
developing or by placing them in a folder in release mode. This can be applied to configuration files
or templates as well and is very useful for whitelabel software which must be highly customizable without
forking or recompilation. You even don't need to restart the process at all.

## Milestones and features

- [x] find a cool name and define project goals  
- [ ] embedding files into byte slices    
- [ ] parse templates at processing time and fail early  
- [ ] using non-embedded files in development mode  
- [ ] overlay resources from different origins  
- [ ] fsnotify in release and development mode  
- [ ] localized resources  
- [ ] multi-client/multi-tenant resources  
- [ ] build variants and exclude resources by compilation flag  
- [ ] basic localization using android strings format  
- [ ] android string plurals and non-translatable texts  
- [ ] generate type safe interfaces to use parameterized strings safely and let compiler do the type checks


## FAQ

### How does it work with modules?
Modules will contain the generated type safe accessors and will register them self at the AssetManager
which provides a central *fnotify* infrastructure and aggregates an aggregated view over all resources. 
Each virtual tree may have an overlay by a local file system tree. 

Resources of GAPT are organized and accessed within a single virtual filesystem of 
[http.FileSystem](https://golang.org/pkg/net/http/#FileSystem). The convention is the following: 
```
/ 
│
└───www
│   │   app.wasm
│   │   style.css
│   │   ...
│   │
│   └───js
│       │   app.js
│       │   libXY.js
│       │   ...
│   
└───tmpl
|   │   header.gohtml
|   │   footer.gohtml
│   │   ...
│   │
│   └───user
│           show.gohtml
│           list.gohtml
│           ...
└───etc
|   │   myproject-config.json
│           ...
|
└───values
│   │   hello_world
│   │   ...
|   |   
└───values-de-DE
│   │   hello_world
│   └───github.com
│   |   └─── mycompany
│   |        └─── myproject
│   |                 userName
│   |                 password
│   │                 ...
|   |   
``` 
The `www`-folder is always published over http as root. All other files and folders should never
be made accessible through the web server. Localized values are also handled as files, however the according `File`
implementation will provide an efficient shortcut. In general, as a convention, a module should fan out their data
with their actual module name. Strings within the values folder are just plain simple utf-8 byte
sequences, with optional positional placeholders. Plurals and gender variants need still to be defined
but will likely be represented by a json format, e.g. indicated by *.pl* extension.

DSL proposal
```bash
# ADD <namespace>:<src>
# $local denotes the local module name 
ADD $local:*.gotmpl :.
ADD config/etc/colors.xml .
```

JSON proposal
```json
{
  ".": {
    "add": ["*.gotmpl","config/etc/colors.xml"],
    "ignore": ["users/test.gotmpl"]
  },
  "github.com/mycompany/privateprj": {
    "remove": ["groups/layout.xml"],
    "add": ["customization/stuff/*"]
  }
}
```