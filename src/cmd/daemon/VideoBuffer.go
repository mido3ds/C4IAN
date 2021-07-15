package main

import (
	"time"

	"github.com/cevaris/ordered_map"
	"github.com/cornelk/hashmap"
	"github.com/mido3ds/C4IAN/src/models"
)

const BufferingTime = 5 * time.Second

type VideoBuffer struct {
	Srcs    *hashmap.HashMap
	Manager *VideoFilesManager
}

type VideoBufferSrc struct {
	Buffer *ordered_map.OrderedMap
}

func NewVideoBuffer(manager *VideoFilesManager) *VideoBuffer {
	return &VideoBuffer{
		Srcs:    &hashmap.HashMap{},
		Manager: manager,
	}
}

func (vb *VideoBuffer) Insert(frag *models.VideoFragment) {
	val, ok := vb.Srcs.Get(frag.Src)
	if !ok {
		entry := &VideoBufferSrc{
			Buffer: ordered_map.NewOrderedMap(),
		}
		vb.Srcs.Set(frag.Src, entry)
		val = entry
	}
	buffer := val.(*VideoBufferSrc).Buffer
	if buffer.Len() == 0 {
		time.AfterFunc(BufferingTime, func() {
			iter := val.(*VideoBufferSrc).Buffer.IterFunc()
			for videoFrag, ok := iter(); ok; videoFrag, ok = iter() {
				vb.Manager.AppendVideoFragment(videoFrag.Value.(*models.VideoFragment))
			}
			// clear buffer of this src
			entry := &VideoBufferSrc{
				Buffer: ordered_map.NewOrderedMap(),
			}
			vb.Srcs.Set(frag.Src, entry)
		})
	}
	buffer.Set(frag.ID, frag)
}
