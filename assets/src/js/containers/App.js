/* eslint
react/forbid-prop-types: 'warn'
*/
import React from 'react';
import FileList from './FileList';
import { BrowserRouter as Router, Route } from 'react-router-dom';

//import { connect } from 'react-redux';

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
          return <FileList path={window.location.pathname} />;
        }}
      />
    </div>
  </Router>
);

export default App;
