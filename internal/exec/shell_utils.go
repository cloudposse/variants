package exec

import (
	"bytes"
	"context"
	"fmt"
	u "github.com/cloudposse/atmos/pkg/utils"
	"io"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ExecuteShellCommand prints and executes the provided command with args and flags
func ExecuteShellCommand(
	command string,
	args []string,
	dir string,
	env []string,
	dryRun bool,
	verbose bool,
	redirectStdError string,
) error {
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if redirectStdError == "" {
		cmd.Stderr = os.Stderr
	} else {
		f, err := os.OpenFile(redirectStdError, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		defer func(f *os.File) {
			err = f.Close()
			if err != nil {
				u.PrintError(err)
			}
		}(f)

		cmd.Stderr = f
	}

	if verbose {
		u.PrintInfo("\nExecuting command:")
		u.PrintMessage(cmd.String())
	}

	if dryRun {
		return nil
	}

	return cmd.Run()
}

// ExecuteShell runs a shell script
func ExecuteShell(command string, name string, dir string, env []string, dryRun bool, verbose bool) error {
	if verbose {
		u.PrintInfo("\nExecuting command:")
		u.PrintMessage(command)
	}

	if dryRun {
		return nil
	}

	return shellRunner(command, name, dir, env, os.Stdout)
}

// ExecuteShellAndReturnOutput runs a shell script and capture its standard output
func ExecuteShellAndReturnOutput(command string, name string, dir string, env []string, dryRun bool, verbose bool) (string, error) {
	var b bytes.Buffer

	if verbose {
		u.PrintInfo("\nExecuting command:")
		u.PrintMessage(command)
	}

	if dryRun {
		return "", nil
	}

	err := shellRunner(command, name, dir, env, &b)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// shellRunner uses mvdan.cc/sh/v3's parser and interpreter to run a shell script and divert its stdout
func shellRunner(command string, name string, dir string, env []string, out io.Writer) error {
	parser, err := syntax.NewParser().Parse(strings.NewReader(command), name)
	if err != nil {
		return err
	}

	environ := append(os.Environ(), env...)
	listEnviron := expand.ListEnviron(environ...)
	runner, err := interp.New(
		interp.Dir(dir),
		interp.Env(listEnviron),
		interp.StdIO(os.Stdin, out, os.Stderr),
	)
	if err != nil {
		return err
	}

	return runner.Run(context.TODO(), parser)
}

// ExecuteShellCommandAndReturnOutput prints and executes the provided command with args and flags and returns the command output
func ExecuteShellCommandAndReturnOutput(
	command string,
	args []string,
	dir string,
	env []string,
	dryRun bool,
	verbose bool,
	redirectStdError string,
) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin

	if redirectStdError == "" {
		cmd.Stderr = os.Stderr
	} else {
		f, err := os.OpenFile(redirectStdError, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return "", err
		}

		defer func(f *os.File) {
			err = f.Close()
			if err != nil {
				u.PrintError(err)
			}
		}(f)

		cmd.Stderr = f
	}

	if verbose {
		u.PrintInfo("\nExecuting command:")
		u.PrintMessage(cmd.String())
	}

	if dryRun {
		return "", nil
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// ExecuteShellCommands sequentially executes the provided list of commands
func ExecuteShellCommands(
	commands []string,
	dir string,
	env []string,
	dryRun bool,
	verbose bool,
	redirectStdError string,
) error {
	for _, command := range commands {
		args := strings.Fields(command)
		if len(args) > 0 {
			if err := ExecuteShellCommand(args[0], args[1:], dir, env, dryRun, verbose, redirectStdError); err != nil {
				return err
			}
		}
	}
	return nil
}

// execTerraformShellCommand executes `terraform shell` command by starting a new interactive shell
func execTerraformShellCommand(
	component string,
	stack string,
	componentEnvList []string,
	varFile string,
	workingDir string,
	workspaceName string,
	componentPath string) error {

	componentEnvList = append(componentEnvList, fmt.Sprintf("TF_CLI_ARGS_plan=-var-file=%s", varFile))
	componentEnvList = append(componentEnvList, fmt.Sprintf("TF_CLI_ARGS_apply=-var-file=%s", varFile))
	componentEnvList = append(componentEnvList, fmt.Sprintf("TF_CLI_ARGS_refresh=-var-file=%s", varFile))
	componentEnvList = append(componentEnvList, fmt.Sprintf("TF_CLI_ARGS_import=-var-file=%s", varFile))
	componentEnvList = append(componentEnvList, fmt.Sprintf("TF_CLI_ARGS_destroy=-var-file=%s", varFile))
	componentEnvList = append(componentEnvList, fmt.Sprintf("TF_CLI_ARGS_console=-var-file=%s", varFile))

	u.PrintInfo("\nStarting a new interactive shell where you can execute all native Terraform commands (type 'exit' to go back)")
	u.PrintMessage(fmt.Sprintf("Component: %s\n", component))
	u.PrintMessage(fmt.Sprintf("Stack: %s\n", stack))
	u.PrintMessage(fmt.Sprintf("Working directory: %s\n", workingDir))
	u.PrintMessage(fmt.Sprintf("Terraform workspace: %s\n", workspaceName))
	u.PrintInfo("\nSetting the ENV vars in the shell:\n")
	for _, v := range componentEnvList {
		u.PrintInfo(v)
	}

	// Transfer stdin, stdout, and stderr to the new process and also set the target directory for the shell to start in
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   componentPath,
		Env:   append(os.Environ(), componentEnvList...),
	}

	// Start a new shell
	var shellCommand string

	if runtime.GOOS == "windows" {
		shellCommand = "cmd.exe"
	} else {
		// If 'SHELL' ENV var is not defined, use 'bash' shell
		shellCommand = os.Getenv("SHELL")
		if len(shellCommand) == 0 {
			bashPath, err := exec.LookPath("bash")
			if err != nil {
				return err
			}
			shellCommand = bashPath
		}
		shellCommand = shellCommand + " -l"
	}

	u.PrintMessage(fmt.Sprintf("Starting process: %s\n", shellCommand))

	args := strings.Fields(shellCommand)

	proc, err := os.StartProcess(args[0], args[1:], &pa)
	if err != nil {
		return err
	}

	// Wait until user exits the shell
	state, err := proc.Wait()
	if err != nil {
		return err
	}

	u.PrintMessage(fmt.Sprintf("Exited shell: %s\n", state.String()))
	return nil
}
