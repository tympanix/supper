import React, { Component } from 'react';
import { render } from 'react-dom';
import App from './App';

import '../sass/styles.scss';

import { BrowserRouter, Route, Switch, Link } from 'react-router-dom'

import Search from './Search'
import Details from './Details'

class Root extends Component {
  render() {

    return (
      <BrowserRouter>
        <Switch>
          <Route exact path="/" component={Search} />
          <Route path="/details" component={Details} />

          // Default route
          <Route component={Search} />
        </Switch>
      </BrowserRouter>
    )
  }
}


render(<Root/>, document.getElementById('root'));
