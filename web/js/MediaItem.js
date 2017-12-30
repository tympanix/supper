import React, { Component } from 'react';
import axios from 'axios';
import { Redirect, Link } from 'react-router-dom'

import Api from './Api'
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
    let link = {pathname: "/details", state: {folder: this.props.item}}

    return (
      <li className={active}>
        <Link to={link}>
          {this.props.item.name}
        </Link>
      </li>
    );
  }

  gotoMedia(e) {
    Api.getMediaDetails(this.props.item).then((red) => {
      this.setState({redirect: red})
    })
  }
}
