package firsterror

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"runtime/debug"
)

/*
Allow to do actions while internal Err is nil.
If err is not nil - return error in every call without real work.

It isn't internally sync for concurrency - method must be call ordered by one method at time.
*/
type FirstError struct {
	err      error           // first error
	ErrStack []byte          // stacktrace to first error
	Context  context.Context // all methods return error if context error (for fast finish all works). It isn't interrupt method in middle of action while context cancelled.
}

func New()*FirstError{
	return &FirstError{
	}
}

func (fe *FirstError) Err() error {
	if fe.err != nil {
		return fe.err
	}
	if fe.Context != nil {
		return fe.setError(fe.Context.Err())
	}
	return nil
}

func (fe *FirstError) Error() string {
	if fe.err == nil {
		return "<nil>"
	} else {
		return fe.err.Error()
	}
}

func (fe *FirstError) BinaryRead(r io.Reader, order binary.ByteOrder, data interface{}) error {
	return fe.Do(func() error { return binary.Read(r, order, data) })
}

func (fe *FirstError) BinaryWriter(w io.Writer, order binary.ByteOrder, data interface{}) error {
	return fe.Do(func() error { return binary.Write(w, order, data) })
}

func (fe *FirstError) Close(force bool, c io.Closer) error {
	return fe.DoForce(c.Close)
}

func (fe *FirstError) Copy(dst io.Writer, src io.Reader)(written int64, err error){
	err = fe.Do(func()error{
		written, err = io.Copy(dst, src)
		return err
	})
	return
}

// Call function with error return.
// It handle panic and convert it in error.
// example:
// fe.Do(func()error{ fmt.Println("asd"); return nil})
func (fe *FirstError) Do(f func() error) error {
	return fe.do(false, f)
}

// Call function if it have no previous error
func (fe *FirstError) DoIt(f func()) error {
	return fe.do(false, func() error { f(); return nil })
}

// Call function with error return.
// It handle panic and convert it in error.
// function call independent of previous error, but
// if call has new error - save it
// example:
// fe.Do(func()error{ fmt.Println("asd"); return nil})
func (fe *FirstError) DoForce(f func() error) error {
	return fe.do(true, f)
}

// Call function with error return.
// It handle panic and convert it in error.
// Force - mean do action if have prev errors
func (fe *FirstError) do(force bool, f func() error) (err error) {
	defer func() {
		p := recover()
		if p == nil {
			return
		}

		if pErr, ok := p.(error); ok {
			err = fe.setError(pErr)
			return
		}

		err = fe.setError(fmt.Errorf("Panic handled: %v", p))
	}()

	// first - check context
	if fe.Context != nil && fe.Context.Err() != nil {
		return fe.setError(fe.Context.Err())
	}

	if fe.err != nil {
		return fe.err
	}

	return fe.setError(f())
}

func (fe *FirstError) GetCloser(c io.Closer, force bool) io.Closer {
	return internalCloser{
		fe:         fe,
		c:          c,
		forceClose: force,
	}
}

func (fe *FirstError) GetReader(r io.Reader) io.Reader {
	return internalReader{
		fe: fe,
		r:  r,
	}
}

func (fe *FirstError) GetWriter(w io.Writer) io.Writer {
	return internalWriter{
		fe: fe,
		w:  w,
	}
}

func (fe *FirstError) Read(r io.Reader, buf []byte) (readBytes int, err error) {
	err = fe.Do(func() error {
		readBytes, err = r.Read(buf)
		return err
	})
	return
}

func (fe *FirstError) Reset() {
	fe.err = nil
	fe.ErrStack = nil
}

func (fe *FirstError) Write(w io.Writer, buf []byte) (writeBytes int, err error) {
	err = fe.Do(func()error {
		writeBytes, err = w.Write(buf)
		return err
	})
	return
}

func (fe *FirstError) setError(err error) error {
	if fe.err != nil {
		return fe.err
	}

	if err == nil {
		return nil
	}

	fe.err = err
	fe.ErrStack = debug.Stack()

	return err
}
