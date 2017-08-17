/* eslint
jsx-a11y/media-has-caption: 'warn',
*/

import React from 'react';
import PropTypes from 'prop-types';

const VideoPlayer = (props) => {
  const { path, mime, subtitles } = props;
  let hasDefaultSubtitle = false;
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
};

VideoPlayer.propTypes = {
  path: PropTypes.string,
  mime: PropTypes.string,
  subtitles: PropTypes.arrayOf(PropTypes.shape({
    name: PropTypes.path,
  })),
};

VideoPlayer.defaultProps = {
  path: '/',
  mime: '',
  subtitles: [],
};


export default VideoPlayer;
