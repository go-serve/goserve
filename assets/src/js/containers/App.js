/* eslint
react/forbid-prop-types: 'warn'
*/
import React from 'react';
import PathPreview from './PathPreview';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import basename from 'basename';

const App = () => (
  <Router>
    <Switch>
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
          return <main className="path-display-wrapper">
            <PathPreview
              sort={sort}
              path={decodeURIComponent(window.location.pathname)}
            />
          </main>;
        }}
      />
    </Switch>
  </Router>
);

export default App;
