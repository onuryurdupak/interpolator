**DESCRIPTION:**

Interpolator is a regex compliant template processor. It can be used for editing text files via key value pairs.

**USAGE:**

Assuming interpolator (or interpolator.exe) binary is globally accessible from your terminal, running:

```
interpolator -r somefile.txt '=' abcd=1234 xyz=abc
```

will:

1-) Read the contents of somefile.txt,

2-) Use equal sign '=' as separator and find every instance of 'abcd' and replace them with '1234',

3-) Then it will find every instance of 'xyz' and replace them with 'abc',

4-) Argument -r represents "recursive" mode. Without it, interpolator will produce an error message when a key is defined multiple times in the file.

Another example, running:

```
interpolator ./embed/data.go ':=' 'Stamp_build_date\s+=\s+"\${build_date}":=Stamp_build_date = '\"$DATE\"
```

(for /embded/data.go file in the repository) will:

1-) Read the contents of ./embed/data.go,

2-) Use ':=' as seperator and find every instance of 'Stamp_build_date\s+=\s+"\${build_date}' regex statement,

3-) Since -r argument is not given and regex statement will match only a single instance in the file, it will change the matching content to 'Stamp_build_date = '\"$DATE\"'.

Note that backslash character is used for escaping purposes.
