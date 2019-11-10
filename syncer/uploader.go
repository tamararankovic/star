package syncer

import fPb "github.com/c12s/scheme/flusher"

type Uploader interface {
	Upload(data *fPb.Update)
	Error(topic string, data []byte)
}
