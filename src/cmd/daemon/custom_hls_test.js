var CustomHlsPlayer = function () {
  var currentLoadingFrag = null;
  var inSegmentErrorRetry = false;

  function retryWithNextSegment(video, hlsPlayer) {
    if (currentLoadingFrag !== null) {
      inSegmentErrorRetry = true;
      var nextStartPoint =
        currentLoadingFrag.start + currentLoadingFrag.duration;
      hlsPlayer.startLoad(nextStartPoint);
      hlsPlayer.recoverMediaError();
      video.play();
    } else {
      console.log("Current fragment is null!!");
    }
  }
  this.switchToHls = function (video, videoUrl) {
    if (Hls.isSupported()) {
      var hls = new Hls();
      hls.loadSource(videoUrl);
      if (videoUrl) {
        hls.attachMedia(video);
      }
      hls.on(Hls.Events.MANIFEST_PARSED, function (event, data) {
        video.play();
      });

      hls.on(Hls.Events.ERROR, function (event, data) {
        switch (data.type) {
          case Hls.ErrorTypes.NETWORK_ERROR:
            switch (data.details) {
              case Hls.ErrorDetails.FRAG_LOAD_ERROR:
                if (inSegmentErrorRetry) {
                  retryWithNextSegment(video, hls);
                }
                break;
              default:
                break;
            }
            break;
          case Hls.ErrorTypes.MEDIA_ERROR:
            switch (data.details) {
              case Hls.ErrorDetails.BUFFER_STALLED_ERROR:
                retryWithNextSegment(video, hls);
                break;
              default:
                break;
            }
            break;
          default:
            hls.destroy();
            break;
        }
      });

      hls.on(Hls.Events.FRAG_LOADING, function (event, data) {
        currentLoadingFrag = data.frag;
      });

      hls.on(Hls.Events.FRAG_LOADED, function (event, data) {
        if (inSegmentErrorRetry === true) {
          inSegmentErrorRetry = false;
        }
      });
      hls.on("hlsError", function (event, data) {});
    }
  };
};
