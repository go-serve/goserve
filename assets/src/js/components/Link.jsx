import React from 'react';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';

const linkPropTypes = {
  className: PropTypes.string,
  link: PropTypes.shape({
    name: PropTypes.string,
    path: PropTypes.string,
    type: PropTypes.string,
  }),
  children: React.PropTypes.oneOfType([
    React.PropTypes.arrayOf(React.PropTypes.node),
    React.PropTypes.node,
  ]),
};

const linkDefaultProps = {
  className: '',
  link: {
    name: '',
    path: '',
    type: 'file',
  },
  children: null,
};

const DirectLink = ({ className, link: { path }, children }) => <a className={className} href={`${path}`}>{children}</a>;
DirectLink.propTypes = linkPropTypes;
DirectLink.defaultProps = linkDefaultProps;

const RouteLink = ({ className, link: { path }, children }) => <Link className={className} to={`${path}`}>{children}</Link>;
RouteLink.propTypes = linkPropTypes;
RouteLink.defaultProps = linkDefaultProps;

export {
  DirectLink,
  RouteLink,
};

export default RouteLink;
