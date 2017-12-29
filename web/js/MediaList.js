import React, { Component } from 'react';
import axios from 'axios';

import MediaItem from './MediaItem'

export default class MediaList extends Component {
  render() {
    const media = this.props.list.map((m) =>
      <MediaItem item={m} key={m.name} />
    )

    return (
      // Add your component markup and other subcomponent references here.
      <ul className="media-list">
        {media}
      </ul>
    );
  }
}
