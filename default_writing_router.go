package flog


import (
	"github.com/reiver/go-dotquote"
	"github.com/reiver/go-oi"

	"fmt"
	"io"
	"time"
)


// NewDefaultWritingRouter returns an initialized DefaultWritingRouter
func NewDefaultWritingRouter(writer io.Writer) *DefaultWritingRouter {
	return NewDefaultWritingRouterWithPrefix(writer, nil)
}

// NewDefaultWritingRouterWithPrefix returns an initialized DefaultWritingRouter
func NewDefaultWritingRouterWithPrefix(writer io.Writer, prefix map[string]interface{}) *DefaultWritingRouter {
	var prefixBuffer []byte
	if 0 < len(prefix) {
		prefixBuffer = dotquote.AppendMap(prefixBuffer, prefix)
	}

	router := DefaultWritingRouter{
		writer:writer,
		prefix:prefixBuffer,
	}

	return &router
}


// DefaultWritingRouter is a router that writes the log in a default format.
//
// A DefaultWritingRouter is appropriate for a production (i.e., "PROD")
// deployment enviornment.
type DefaultWritingRouter struct {
	writer io.Writer
	prefix []byte
}



func (router *DefaultWritingRouter) Route(message string, context map[string]interface{}) error {
	if nil == router {
		return errNilReceiver
	}


	writer := router.writer
	if nil == writer {
		// Nothing to do, so just return silently.
		return nil
	}


	var buffer [512]byte
	p := buffer[:0]

	if prefix := router.prefix; 0 < len(prefix) {
		p = append(p, prefix...)
	}


	p = dotquote.AppendString(p, message, "text")
	p = append(p, ' ')
	p = dotquote.AppendString(p, time.Now().String(), "when")

	// If we have an error, then get the error.Error() into the log too.
	if errorFieldValue, ok := context["~error"]; ok {
		if err, ok := errorFieldValue.(error); ok {
			p = append(p, ' ')
			p = dotquote.AppendString(p, fmt.Sprintf("%T", err), "error", "type")
			p = append(p, ' ')
			p = dotquote.AppendString(p, err.Error(), "error", "text")
		}
	}


//@TODO: This is a potential heavy operation. Is there a better way
//       to get the ultimate result this is trying to archive?
//
	if 0 < len(context) {
		p = append(p, ' ')
		p = dotquote.AppendMap(p, context, "ctx")
	}


	_,_ = oi.LongWrite(router.writer, p)

//@TODO: Should this be checking for errors from oi.LongWrite()?
	return nil
}