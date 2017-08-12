import React from 'react';
import { graphql, gql } from 'react-apollo';
import { Link } from 'react-router-dom';

export const Query = gql`
  query FileListQuery ($path: String = "/") {
    self: stat(path:$path){
      name
      path
      type
    }
    children: list(path:$path){
      name
      path
      type
    }
  }
`;

const FileList = function(props) {
  const { path, data: { self, children=[] } } = props;
  console.log("props.path", path, "self", self, "children", children);
  if (typeof self === 'undefined') return null;
  return (
    <div>
      <h1>{ `Index of ${self.path}` }</h1>
      <ul className="listing">
        {children.map((child) => (
          <li key={child.path}>
            <Link to={`${child.path}`}>
              {child.name}{ (child.type == "directory") ? "/" : "" }
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default graphql(Query, {
  options: ({path = "/"}) => {
    const options = {
      variables: {
        path,
      },
    };
    return options;
  },
})(FileList);
