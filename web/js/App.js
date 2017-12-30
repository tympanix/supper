import React, { Component } from 'react';
import axios from 'axios';

import Action from './Action'
import Search from './Search'
import Details from './Details'

class App extends Component {
  constructor() {
    super()

    this.state = {
      media: undefined
    }

    Action.onShowMedia((media) => {
      console.log(media)
      this.setState({media})
    })
  }

  render() {
    if (this.state.media) {
      return (<Details />);
    } else {
      return (<Search />);
    }
  }
}

export default App
