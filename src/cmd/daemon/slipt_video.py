#!/usr/bin/env python
from moviepy.video.io.ffmpeg_tools import ffmpeg_extract_subclip
from moviepy.editor import VideoFileClip, concatenate_videoclips, CompositeVideoClipClip
import time
import os

VideoFile = "video.mp4"
TimeFile = "times.txt"
TimeToSleep = 5

with open(TimeFile) as f:
  times = f.readlines()

times = [x.strip() for x in times] 
delete_list_after = []
for time in times:
    starttime = int(time.split("-")[0])
    endtime = int(time.split("-")[1])
    new_video_name = sre(times.index(time)+1)+".mp4"
    ffmpeg_extract_subclip(VideoFile, starttime, endtime, targetname=new_video_name)
    delete_list_after.append(new_video_name)

video_name = "1.mp4"
for i in range(1, len(times)):
    time.sleep(TimeToSleep)
    video1 = VideoFileClip(video_name)
    video2 = VideoFileClip(str(i+1)+".mp4")
    final_clip = concatenate_videoclips([video1,video2])
    final_clip.write_videofile(video_name)

for file in delete_list_after:
    os.remove(file)