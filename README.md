**What is it?**

interpolator is a template processor. It can be used for editing text files by providing key value pairs.



**Usage Example:**   

Running the following command line statement:

`interpolator ./main.go ':=' 'stamp_build_date\s+=\s+"\${build_date}":=stamp_build_date = '\"$DATE\"`

for `main.go` file in the repository will replace the given pattern `stamp_build_date\s+=\s+"\${build_date}"` using seperator `:=` and convert the pattern in file to `stamp_build_date = '\"$DATE\"`


**Breakdown:**

- First command line argument is the path of the file to change.
- Second argument is the user defined key=value seperator. 
- Each key definition (including regex in given example( must exist in a document exactly one time.

**Use case:**

Editing configuration files in a CI/CD pipeline.
