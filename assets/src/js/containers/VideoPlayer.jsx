import React from 'react';

const VideoPlayer = function(props) {
  const { path, mime, subtitles } = props;
  var hasDefaultSubtitle = false;
  return (
    <video controls>
      <source key="mp4" src={path} type={mime} />
      {subtitles.map((subtitle) => {
        const def = !hasDefaultSubtitle;
        // TODO: human readable subtitle lang / label
        if (!hasDefaultSubtitle) hasDefaultSubtitle = true;
        return (<track
          key={subtitle.path}
          kind="subtitles"
          src={`${subtitle.path}?mode=vtt`}
          srcLang="Subtitle"
          label="Subtitle"
          default={def}
        />);
      })}
    </video>
  );
}

export default VideoPlayer;
