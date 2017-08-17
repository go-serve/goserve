import React from 'react';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';

const renderNormalDirectoryLink = dir => <a href={`${dir.path}`}>{dir.name}/</a>;
const renderDirectoryLink = dir => <Link to={`${dir.path}`}>{dir.name}/</Link>;
const renderVideoLink = file => <Link to={`${file.path}`}>{file.name}</Link>;
const renderNormalLink = file => <a href={`${file.path}`}>{file.name}</a>;
const renderLink = (item) => {
  if (item.type === 'directory' && item.hasIndex === true) return renderNormalDirectoryLink(item);
  if (item.type === 'directory') return renderDirectoryLink(item);
  if (item.mime === 'video/mp4') return renderVideoLink(item);
  return renderNormalLink(item);
};

const FileList = (props) => {
  const { className = 'filelist', self, containing = [] } = props;
  if (typeof self === 'undefined' || self === null) return null;
  return (
    <div className={className}>
      <ul className="listing">
        {containing.map(child => (
          <li key={child.path}>
            {renderLink(child)}
          </li>
        ))}
      </ul>
    </div>
  );
};

FileList.propTypes = {
  className: PropTypes.string,
  self: PropTypes.shape({
    name: PropTypes.string,
    mime: PropTypes.string,
  }),
  containing: PropTypes.arrayOf(PropTypes.object),
};

FileList.defaultProps = {
  className: '',
  self: null,
  containing: [],
};

export default FileList;
