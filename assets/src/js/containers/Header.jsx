import React from 'react';
import { Link } from 'react-router-dom';

const renderNormalDirectoryLink = ({path}, className) => <a className={className} href={`${path}`}>&lt;</a>;
const renderDirectoryLink = ({path}, className) => <Link className={className} to={`${path}`}>&lt;</Link>;

const Prev = ({self: {name, parent}, className=""}) => {
  const link = (parent.hasIndex === true) ?
    renderNormalDirectoryLink(parent, className) :
    renderDirectoryLink(parent, className);
  return (name !== '/') ? link : <span className={className} />;
}

const Placeholder = ({ className }) => {
  return <span className={className}></span>
}

const Header = ({ name, path, type, className, parent }) => {
  console.log("parent:", parent);
  return (
    <header className={className}>
      <Prev className="prev" self={{ name, parent }} />
      <h1>{ `${(name === '/') ? '' : name}` }{ (type === 'directory') ? '/' : '' }</h1>
      <Placeholder className="user" />
    </header>
  );
}

export default Header;
