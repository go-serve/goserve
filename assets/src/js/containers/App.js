/* eslint
react/forbid-prop-types: 'warn'
*/
import React from 'react';
import FileList from './FileList';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import basename from 'basename';

const App = () => (
  <Router>
    <div>
      <Route
        path="/_goserve"
        render={() => <div>Incorrect link</div>}
      />
      <Route
        path="/"
        render={() => {
          const search = new URLSearchParams(window.location.search.substring(1));
          const sort = (search.has("sort") && (search.get("sort") !== "")) ?
            search.get("sort") : "-mtime";
          const selfName = basename(window.location.pathname) + "/";
          return <div className="path-display-wrapper">
            <Helmet>
              <title>{`Index of ${selfName}`}</title>
            </Helmet>
            <FileList
              sort={sort}
              path={window.location.pathname}
            />
          </div>;
        }}
      />
    </div>
  </Router>
);

export default App;
