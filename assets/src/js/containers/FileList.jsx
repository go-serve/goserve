import React from 'react';
import { graphql, gql } from 'react-apollo';
import { Link } from 'react-router-dom';

export const Query = gql`
  query FileListQuery ($path: String = "/", $sort: String = "-mtime") {
    self: stat(path:$path){
      name
      path
      type
    }
    children: list(path:$path, sort:$sort){
      name
      path
      type
    }
  }
`;

const FileList = function(props) {
  const { className="filelist", path, data: { self, children=[] } } = props;
  console.log("props.path", path, "self", self, "children", children);
  if (typeof self === 'undefined') return null;
  return (
    <div className={ className }>
      <h1>{ `Index of ${self.path}` }</h1>
      <ul className="listing">
        {children.map((child) => (
          (child.type === "directory") ?
            (<li key={child.path}>
              <Link to={`${child.path}`}>{child.name}/</Link>
            </li>) :
            (<li key={child.path}>
              <a href={`${child.path}`}>{child.name}</a>
            </li>)
        ))}
      </ul>
    </div>
  );
}

export default graphql(Query, {
  options: ({path = "/", sort="-mtime"}) => {
    const options = {
      variables: {
        path,
        sort,
      },
    };
    return options;
  },
})(FileList);
