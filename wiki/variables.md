# Variables

For the convenience of use, gmake2 has some built-in available variables, which will be continuously updated.


# Menu
- [Variables](#variables)
- [Menu](#menu)
    - [time.now](#timenow)
    - [time.utc](#timeutc)
    - [time.unix](#timeunix)
    - [time.utc\_unix](#timeutc_unix)
    - [runtime.os](#runtimeos)
    - [runtime.arch](#runtimearch)
    - [gmake2.version](#gmake2version)
    - [gmake2.code](#gmake2code)
    - [gmake2.time](#gmake2time)

### time.now
Current time

```
@echo {{.time.now}}
```

### time.utc
Current UTC time

```
@echo {{.time.utc}}
```

### time.unix
Current Unix Time

```
@echo {{.time.unix}}
```

### time.utc_unix
Current UTC Unix time

```
@echo {{.time.utc_unix}}
```

### runtime.os
Current system name

```
@echo {{.runtime.os}}
```

### runtime.arch
Current System Architecture

```
@echo {{.runtime.arch}}
```

### gmake2.version
Current GMake2 Version

```
@echo {{.gmake2.version}}
```

### gmake2.code
Current GMake2 Version Code

```
@echo {{.gmake2.code}}
```

### gmake2.time
Current GMake2 Build Time

```
@echo {{.gmake2.time}}
```