package main

import (
	"time"

	"github.com/cevaris/ordered_map"
	"github.com/cornelk/hashmap"
	"github.com/mido3ds/C4IAN/src/models"
)

const BufferingTime = 5 * time.Second

type VideosBuffer struct {
	Srcs    *hashmap.HashMap
	Manager *VideoFilesManager
}

type SourceVideoBuffer struct {
	Videos *hashmap.HashMap
}

type OneVideoBuffer struct {
	Video     *ordered_map.OrderedMap
	LastSeqNo uint64
}

func NewVideoBuffer(manager *VideoFilesManager) *VideosBuffer {
	return &VideosBuffer{
		Srcs:    &hashmap.HashMap{},
		Manager: manager,
	}
}

// insert new video fragment with src, id, seqNo
// where src is the source identifier and id is the id of the video and seqNo is the seqNo of the frame
// when recieving a fragment with the same src and id it s
func (vb *VideosBuffer) Insert(frag *models.VideoFragment) {
	srcVideos, ok := vb.Srcs.Get(frag.Src)
	if !ok {
		entry := &SourceVideoBuffer{
			Videos: &hashmap.HashMap{},
		}
		vb.Srcs.Set(frag.Src, entry)
		srcVideos = entry
	}
	videos, ok := srcVideos.(*SourceVideoBuffer).Videos.Get(frag.ID)
	if !ok {
		entry := &OneVideoBuffer{
			Video: ordered_map.NewOrderedMap(),
		}
		srcVideos.(*SourceVideoBuffer).Videos.Set(frag.ID, entry)
		videos = entry
	}
	video := videos.(*OneVideoBuffer).Video
	if video.Len() == 0 {
		time.AfterFunc(BufferingTime, func() {
			iter := video.IterFunc()
			for videoFrag, ok := iter(); ok; videoFrag, ok = iter() {
				vf := videoFrag.Value.(*models.VideoFragment)
				videos.(*OneVideoBuffer).LastSeqNo = vf.SeqNo
				vb.Manager.AppendVideoFragment(vf)
			}
			// clear buffer of this video
			entry := &OneVideoBuffer{
				Video:     ordered_map.NewOrderedMap(),
				LastSeqNo: 0,
			}
			srcVideos.(*SourceVideoBuffer).Videos.Set(frag.ID, entry)
		})
	}
	if videos.(*OneVideoBuffer).LastSeqNo <= frag.SeqNo {
		video.Set(frag.SeqNo, frag)
	}
}
