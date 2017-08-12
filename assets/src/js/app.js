import React from 'react';
import ReactDOM from 'react-dom';
import { ApolloClient, ApolloProvider, createNetworkInterface } from 'react-apollo';
import { createStore, combineReducers, applyMiddleware, compose } from 'redux';

import App from './containers/App';

import '../css/style.scss';

const client = new ApolloClient({
  networkInterface: createNetworkInterface({
    uri: '/_goserve/api/graphql',
  }),
});

const store = createStore(
  combineReducers({
    apollo: client.reducer(),
  }),
  {}, // initial state
  compose(
      applyMiddleware(client.middleware()),
      // If you are using the devToolsExtension, you can add it here also
      (typeof window.__REDUX_DEVTOOLS_EXTENSION__ !== 'undefined') ? window.__REDUX_DEVTOOLS_EXTENSION__() : f => f,
  ),
);

ReactDOM.render(
  <ApolloProvider client={client}>
    <App loading={true} />
  </ApolloProvider>,
  document.getElementById('app'),
);
