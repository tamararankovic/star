package syncer

// import fPb "github.com/c12s/scheme/flusher"

type Fn func(event []byte)

type Syncer interface {
	Subscribe(topic string, f Fn)
	Error(topic string, data []byte)
}
