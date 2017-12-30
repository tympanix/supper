import React, { Component } from 'react';
import axios from 'axios';
import { Redirect } from 'react-router'

import Action from './Action'
import App from './App'

export default class MediaItem extends Component {
  constructor() {
    super()

    this.state = {
      redirect: false
    }
  }

  render() {
    if (this.state.redirect) {
      return (
        <Redirect push to={this.state.redirect} />
      )
    }

    let active = this.props.active ? 'active' : ''

    return (
      <li onClick={this.gotoMedia.bind(this)} className={active}>
        {this.props.item.name}
      </li>
    );
  }

  gotoMedia(e) {
    axios.post('./api/media', this.props.item).then((res) => {
      this.setState({redirect: {
        pathname: '/details',
        state: { media: res.data }
      }})
    })
  }
}
