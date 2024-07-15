package common

import "io"

type WrappedReader interface {
	io.Reader
	UpstreamReader() io.Reader
}

type WrappedWriter interface {
	io.Writer
	UpstreamWriter() io.Writer
}

func UnwrapReader(reader io.Reader) io.Reader {
	for {
		wrappedReader, isWrapped := reader.(WrappedReader)
		if !isWrapped {
			return reader
		}
		upstream := wrappedReader.UpstreamReader()
		if upstream == nil {
			return reader
		}
		reader = upstream
	}
}

func UnwrapWriter(writer io.Writer) io.Writer {
	for {
		wrappedWriter, isWrapped := writer.(WrappedWriter)
		if !isWrapped {
			return writer
		}
		upstream := wrappedWriter.UpstreamWriter()
		if upstream == nil {
			return writer
		}
		writer = upstream
	}
}
