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
	Video *ordered_map.OrderedMap
}

func NewVideoBuffer(manager *VideoFilesManager) *VideosBuffer {
	return &VideosBuffer{
		Srcs:    &hashmap.HashMap{},
		Manager: manager,
	}
}

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
				vb.Manager.AppendVideoFragment(videoFrag.Value.(*models.VideoFragment))
			}
			// clear buffer of this video
			entry := &OneVideoBuffer{
				Video: ordered_map.NewOrderedMap(),
			}
			srcVideos.(*SourceVideoBuffer).Videos.Set(frag.ID, entry)
		})
	}
	video.Set(frag.SeqNo, frag)
}
