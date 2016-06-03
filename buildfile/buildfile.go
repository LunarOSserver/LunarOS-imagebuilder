package buildfile

import (
    "bytes"
    "fmt"
)

type Buildfile struct {
    bytes.Buffer
}

func New() *Buildfile {
    b := Buildfile{}
    b.WriteString(base)
    return &b
}

// WriteCmd writes a command to the build file. The
// command will be echoed back as a base16 encoded
// command so that it can be parsed and appended to
// the build output
func (b *Buildfile) WriteCmd(command string) {
    // echo the command as an encoded value
    b.WriteString(fmt.Sprintf("echo '#DOME:%s'\n", command))
    // and then run the command
    b.WriteString(fmt.Sprintf("%s\n", command))
}

// WriteCmdSilent writes a command to the build file
// but does not echo the command.
func (b *Buildfile) WriteCmdSilent(command string) {
    b.WriteString(fmt.Sprintf("%s\n", command))
}

// WriteComment adds a comment to the build file. This
// is really only used internally for debugging purposes.
func (b *Buildfile) WriteComment(comment string) {
    b.WriteString(fmt.Sprintf("#%s\n", comment))
}

// WriteEnv exports the environment variable as
// part of the script. The environment variables
// are not echoed back to the console, and are
// kept private by default.
func (b *Buildfile) WriteEnv(key, value string) {
    b.WriteString(fmt.Sprintf("export %s=%q\n", key, value))
}

// WriteHost adds an entry to the /etc/hosts file.
func (b *Buildfile) WriteHost(mapping string) {
    b.WriteCmdSilent(fmt.Sprintf("[ -f /usr/bin/sudo ] || echo %q | tee -a /etc/hosts", mapping))
    b.WriteCmdSilent(fmt.Sprintf("[ -f /usr/bin/sudo ] && echo %q | sudo tee -a /etc/hosts", mapping))
}

// WriteFile add files as part of the script.
func (b *Buildfile) WriteFile(path string, file []byte, i int) {
    b.WriteString(fmt.Sprintf("echo '%s' | tee %s > /dev/null\n", string(file), path))
    b.WriteCmdSilent(fmt.Sprintf("chmod %d %s", i, path))
}

// WriteFile add files as part of the script.
func (b *Buildfile) WriteFileSilent(path string, file []byte, i int) {
    b.WriteCmdSilent(fmt.Sprintf("chmod %d %s", i, path))
}

// every build script starts with the following
// code at the start.
var base = `
#!/bin/sh
set -e

if [ ! -d $HOME/.ssh ]; then
  mkdir -p $HOME/.ssh
fi

chmod 0700 $HOME/.ssh
echo 'StrictHostKeyChecking no' | tee $HOME/.ssh/config > /dev/null

# be sure to exit on error and print out
# our sh commands, so we can see which commands
# are executing and troubleshoot failures.
# set -e

# user-defined commands below ##############################
`
