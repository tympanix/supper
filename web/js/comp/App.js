import React, { Component } from 'react';

import { BrowserRouter, Route, Switch, Link } from 'react-router-dom'

import Search from './Search'
import Details from './Details'
import Checkmark from './Checkmark'

class App extends Component {
  constructor() {
    super()
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
        <button onClick={() => {Checkmark.show()}}>Click Me!</button>
        <Checkmark/>
      </div>
    )
  }

}

export default App
