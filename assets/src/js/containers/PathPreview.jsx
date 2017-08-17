import React from 'react';
import { graphql, gql } from 'react-apollo';
import BodyClassName from 'react-body-classname';
import { Helmet } from 'react-helmet';
import PropTypes from 'prop-types';

import Header from './Header';
import FileList from './FileList';
import VideoPlayer from './VideoPlayer';

export const Query = gql`
  query FileListQuery ($path: String = "/", $sort: String = "-mtime") {
    self: stat(path:$path){
      name
      path
      type
      mime
      parent {
        name
        path
        hasIndex
      }
      subtitles: siblings(nameLikeMe: true, nameLike: "*.srt") {
        path
      }
    }
    containing: list(path:$path, sort:$sort){
      name
      path
      type
      mime
      hasIndex
    }
  }
`;

const PathPreview = (props) => {
  const { path = '/', data: { self = null, containing = [] } } = props;
  if (self === null) return null;
  if (self.type === 'file' && self.mime === 'video/mp4') {
    return (
      <BodyClassName className="page-video">
        <div>
          <Helmet>
            <title>{`${self.name}`}</title>
          </Helmet>
          <Header {...self} />
          <div className="video-container">
            <VideoPlayer {...self} />
          </div>
        </div>
      </BodyClassName>
    );
  }
  return (
    <BodyClassName className="page-directory">
      <section>
        <Helmet>
          <title>{`Index of ${(self.name === '/') ? '' : self.name}/`}</title>
        </Helmet>
        <Header {...self} />
        <FileList
          path={path}
          self={self}
          containing={containing}
        />
      </section>
    </BodyClassName>
  );
};

PathPreview.propTypes = {
  path: PropTypes.string,
  data: PropTypes.shape({
    self: PropTypes.shape({
      name: PropTypes.string,
      mime: PropTypes.string,
    }),
    containing: PropTypes.arrayOf(PropTypes.object),
  }),
};

PathPreview.defaultProps = {
  path: '/',
  data: {
    self: null,
    containing: [],
  },
};

export default graphql(Query, {
  options: ({ path = '/', sort = '-mtime' }) => {
    const options = {
      variables: {
        path,
        sort,
      },
    };
    return options;
  },
})(PathPreview);
