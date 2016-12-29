package firsterror

import "io"

type internalCloser struct {
	fe         *FirstError
	c          io.Closer
	forceClose bool
}

func (ic internalCloser) Close() error {
	return ic.fe.Close(ic.forceClose, ic.c)
}

type internalReader struct {
	fe *FirstError
	r  io.Reader
}

func (ir internalReader) Read(buf []byte) (int, error) {
	return ir.fe.Read(ir.r, buf)
}

type internalWriter struct {
	fe *FirstError
	w  io.Writer
}

func (iw internalWriter) Write(buf []byte) (int, error) {
	return iw.fe.Write(iw.w, buf)
}

type internalReadWriter struct {
	fe *FirstError
	rw io.ReadWriter
}

func (irw internalReadWriter) Read(buf []byte) (int, error) {
	return irw.fe.Read(irw.rw, buf)
}

func (irw internalReadWriter) Write(buf []byte) (int, error) {
	return irw.fe.Write(irw.rw, buf)
}
