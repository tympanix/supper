import React, { Component } from 'react';

import { BrowserRouter, Route, Switch, Link } from 'react-router-dom'

import Search from './Search'
import Details from './Details'
import Checkmark from './Checkmark'
import Snackbar from './Snackbar'

class App extends Component {
  constructor() {
    super()
  }

  componentDidCatch(error, info) {
    console.error(error, info)
    Snackbar.error("Exception", error)
  }

  render() {
    return (
      <div>
        <BrowserRouter>
          <Switch>
            <Route exact path="/" component={Search} />
            <Route path="/details" component={Details} />

            // Default route
            <Route component={Search} />
          </Switch>
        </BrowserRouter>
        <Snackbar/>
        <Checkmark/>
      </div>
    )
  }

}

export default App
