import React, { Component } from 'react';
import axios from 'axios';

import MediaItem from './MediaItem'

export default class MediaList extends Component {
  constructor() {
    super()

    this.state = {
      selected: -1
    }
  }

  select(delta) {
    if (this.state.selected >= 0) {
      let idx = this.state.selected + delta
      if (idx < 0) {
        return 0
      } else {
        return Math.min(idx, this.props.list.length - 1)
      }
    } else {
      return 0
    }
  }

  selectNext() {
    this.setState({selected: this.select(1)})
  }

  selectPrev() {
      this.setState({selected: this.select(-1)})
  }

  clearSelected() {
    this.setState({selected: -1})
  }

  render() {
    const media = this.props.list.map((m, i) => {
      let active = this.state.selected === i
      return (<MediaItem item={m} key={m.name} active={active} />)
    })

    return (
      // Add your component markup and other subcomponent references here.
      <ul className="media-list">
        {media}
      </ul>
    );
  }
}
