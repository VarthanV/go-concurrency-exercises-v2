# Names

- When a package is imported, the package name becomes an accessor for the contents. After

```go
import bytes
```

- The package name should be good , short, concise , evocative.

- By convention packages are given lower-case, single-word names,there is no need for underscore or mixedCaps.

- The package name is only the default name for imports , It needn't be unique across the codebase.

- Another convention is that the package name is the base name of its source directory; the package in ``src/encoding/base64 ``is imported as ``encoding/base64`` but has name base64, not encoding_base64 and not encodingBase64.

- The importer of a package will use the name to refer to its contents, so exported names in the package can use that fact to avoid repetition.

-  For instance, the buffered reader type in the ``bufio`` package is called ``Reader``, not ``BufReader``, because users see it as ``bufio.Reader``, which is a clear, concise name.

- Imported entities are always refered with their package name , so there is no conflict between ``io.Reader`` and ``bufio.Reader``.

- Similary to ``ring`` package , we use ``ring.New`` to instantiate instead of ``ring.NewRing`` since the Ring is the only type exported.

- Similary in the once package , the function is named as ``once.Do(setup)`` instead of  ``once.DoOrWaitUntilDone(setup)`` long names doesn't justify verbosity , a simple doc comment is enough.

## Getters

- Go doesn't provide automatic support for setters and getters.

- There is noting wrong in providing getters and setters by ourself,It is necessary to do.

- But it is not necessary to put ``Get`` into getters name.

- If we have a field called ``owner``, the getter method can be called ``Owner`` 

- The use of upper-case names for export provides the hook to discriminate the field from the method

-  A setter function, if needed, will likely be called SetOwner

```go
owner := obj.Owner()
if owner != user {
    obj.SetOwner(user)
}
```
## Interfaces

- By convention , one-method interfaces are named by method name plus an -er suffix or similar modifications to construct an agent noun: Reader,Writer,Formatter, CloseNotifier etc.

- To avoid confusion, don't give your method one of those names unless it has the same signature and meaning.

- Call the string converter method to ``String`` instead of ``ToString``.

## MixedCaps

Finally, the convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names.