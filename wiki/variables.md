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