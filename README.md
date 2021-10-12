**What is it?**

interpolator is a template processor. It can be used for editing text files by providing key value pairs.



**Usage Example:**   

Running the following command line statement:

`interpolator a.json \"\${a}\"=123.45 \"\${b}\"=\"new_B\" false=true`

for `a.json` file with following content:

```
{
  "fieldA": "${a}",
  "fieldB": "${b}",
  "c": {
     "fieldOfC": false
  }
}
```

will convert it to:

```
{
  "fieldA": 123.45,
  "fieldB": "new_B",
  "c": {
     "fieldOfC": true
  }
}
```


**Breakdown:**

- First command line argument is the path of the file to change.
- Every following argument must be a key=value definition like in the example.
- Backslash character `\` can be used to escape special shell characters such as single quotes, double quotes and dollar signs.
- Each key definition must exist in a document exactly one time.


**Use case:**

Editing configuration files in a CI/CD pipeline.
