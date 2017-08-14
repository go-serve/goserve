import React from 'react';
import { Link } from 'react-router-dom';

const renderNormalDirectoryLink = (dir) => <a href={`${dir.path}`}>{dir.name}/</a>;
const renderDirectoryLink = (dir) => <Link to={`${dir.path}`}>{dir.name}/</Link>;
const renderVideoLink = (file) => <Link to={`${file.path}`}>{file.name}</Link>;
const renderNormalLink = (file) => <a href={`${file.path}`}>{file.name}</a>;
const renderLink = (item) => {
  if (item.type === "directory" && item.hasIndex === true) return renderNormalDirectoryLink(item);
  if (item.type === "directory") return renderDirectoryLink(item);
  if (item.mime === "video/mp4") return renderVideoLink(item);
  return renderNormalLink(item);
}

const FileList = function(props) {
  const { className="filelist", path, self, children=[] } = props;
  if (typeof self === 'undefined' || self === null) return null;
  return (
    <div className={ className }>
      <h1>{ `Index of ${(self.name === '/') ? '' : self.name}/` }</h1>
      <ul className="listing">
        {children.map((child) => (
          <li key={child.path}>
            {renderLink(child)}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default FileList;
