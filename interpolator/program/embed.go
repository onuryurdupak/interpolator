package program

import (
	"fmt"
	"interpolator/utils/stdout_utils"
)

const (
	stamp_build_date  = "${build_date}"
	stamp_commit_hash = "${commit_hash}"
	stamp_source      = "${source}"

	ErrSuccess  = 0
	ErrInput    = 1
	ErrInternal = 2
	ErrUnknown  = 3

	helpPrompt = `Run 'interpolator -h' for help.`

	helpMessage = `
If your terminal does not render styles properly, run 'interpolator -hr' to view in style-free mode.

<b><u><yellow>PARAMETERS:</yellow></u></b>
<b><yellow>-v</yellow></b>: Show version info.
<b><yellow>-r</yellow></b>: Execute in recursive mode.

<b><u><yellow>DESCRIPTION</yellow></u></b>
Interpolator is a regex compliant template processor. It can be used for editing text files via key value pairs.

<b><u><yellow>USAGE</yellow></u></b>
Assuming interpolator (or interpolator.exe) binary is globally accessible from your terminal, running:

<green>interpolator -r somefile.txt '=' abcd=1234 xyz=abc</green>

will:

1-) Read the contents of somefile.txt,
2-) Use equal sign '=' as separator and find every instance of 'abcd' and replace them with '1234',
3-) Then it will find every instance of 'xyz' and replace them with 'abc'
4-) Argument -r represents "recursive" mode. Without it, interpolator will produce an error message when a key is defined multiple times in the file.

Another example, running:

<green>interpolator ./embed/data.go ':=' 'Stamp_build_date\s+=\s+"\${build_date}":=Stamp_build_date = '\"$DATE\"</green>

(for /embed/data.go file in the repository) will:

1-) Read the contents of ./embed/data.go,
2-) Use ':=' as separator and find every instance of 'Stamp_build_date\s+=\s+"\${build_date}' regex statement,
3-) Since -r argument is not given and regex statement will match only a single instance in the file, it will change the matching content to 'Stamp_build_date = '\"$DATE\"'.

Note that backslash character is used for escaping purposes.
`
)

func versionInfo() string {
	return fmt.Sprintf(`Build Date: %s | Commit: %s
Source: %s`, stamp_build_date, stamp_commit_hash, stamp_source)
}

func helpMessageStyled() string {
	msg, _ := stdout_utils.ProcessStyle(helpMessage)
	return msg
}

func helpMessageUnstyled() string {
	return stdout_utils.RemoveStyle(helpMessage)
}
