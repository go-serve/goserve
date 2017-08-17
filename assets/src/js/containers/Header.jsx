import React from 'react';
import PropTypes from 'prop-types';

import { DirectLink, RouteLink } from '../components/Link';

const Prev = ({ className = '', self: { name, parent } }) => {
  const link = (parent.hasIndex === true) ?
    <DirectLink className={className} link={parent}>&lt;</DirectLink> :
    <RouteLink className={className} link={parent}>&lt;</RouteLink>;
  return (name !== '/') ? link : <span className={className} />;
};
Prev.propTypes = {
  className: PropTypes.string,
  self: PropTypes.shape({
    name: PropTypes.string,
    parent: PropTypes.shape({
      name: PropTypes.string,
      path: PropTypes.string,
      type: PropTypes.string,
    }),
  }),
};
Prev.defaultProps = {
  className: '',
  self: null,
};

const Placeholder = ({ className }) => <span className={className} />;
Placeholder.propTypes = {
  className: PropTypes.string,
};
Placeholder.defaultProps = {
  className: '',
};

const Header = ({ name, type, className, parent }) => (
  <header className={className}>
    <Prev className="prev" self={{ name, parent }} />
    <h1>{ `${(name === '/') ? '' : name}` }{ (type === 'directory') ? '/' : '' }</h1>
    <Placeholder className="user" />
  </header>
);
Header.propTypes = {
  className: PropTypes.string,
  name: PropTypes.string,
  type: PropTypes.string,
  parent: PropTypes.shape({
    name: PropTypes.string,
    path: PropTypes.string,
    type: PropTypes.string,
  }),
};
Header.defaultProps = {
  className: '',
  name: '',
  type: '',
  parent: null,
};

export default Header;
