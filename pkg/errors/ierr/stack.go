package ierr

import (
	"runtime"
	"strconv"
	"strings"
)

const unknownFilePath = "unknown"

// Frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type Frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return unknownFilePath
	}
	file, _ := fn.FileLine(f.pc())

	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())

	return line
}

// name returns the name of this function, if known.
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}

	return fn.Name()
}

// String - representation of frame.
func (f Frame) String() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(f.name())
	sb.WriteString("\n\t")
	sb.WriteString(f.file())
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(f.line()))

	return sb.String()
}

// StackTrace is stack of Frames from innermost (newest) to outermost (oldest).
type StackTrace []Frame

// stack represents a stack of program counters.
type stack []uintptr

// String - representation of stack trace.
func (s *stack) String() string {
	var sb strings.Builder
	for _, pc := range *s {
		sb.WriteString(Frame(pc).String())
	}

	return sb.String()
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:]) //nolint:gomnd
	var st stack = pcs[0:n]

	return &st
}
