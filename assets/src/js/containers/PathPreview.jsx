import React from 'react';
import { graphql, gql } from 'react-apollo';
import BodyClassName from 'react-body-classname';
import { Helmet } from 'react-helmet';

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
    children: list(path:$path, sort:$sort){
      name
      path
      type
      mime
      hasIndex
    }
  }
`;

const PathPreview = function(props) {
  const { path="/", data: { self=null, children=[] } } = props;
  if (self === null) return null;
  if (self.type === "file" && self.mime === "video/mp4") {
    return (
      <BodyClassName className="page-video">
        <div>
          <Helmet>
            <title>{`${self.name}`}</title>
          </Helmet>
          <Header {...self} />
          <div className="video-container">
            <VideoPlayer {...self}/>
          </div>
        </div>
      </BodyClassName>
    );
  }
  const randClass = `class-${Math.random()}`;
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
          children={children}
        />
      </section>
    </BodyClassName>
  );
}

export default graphql(Query, {
  options: ({path = "/", sort="-mtime"}) => {
    console.log(`query: ${path} ${sort}`)
    const options = {
      variables: {
        path,
        sort,
      },
    };
    return options;
  },
})(PathPreview);
