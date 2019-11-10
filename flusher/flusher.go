package flusher

// import fPb "github.com/c12s/scheme/flusher"

type Fn func(event []byte)

type Sync interface {
	Subscribe(topic string, f Fn)
}
