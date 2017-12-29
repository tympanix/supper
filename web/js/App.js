import React, { Component } from 'react';
import axios from 'axios';
import Media from './Media'

class App extends Component {
  constructor() {
    super()
    this.media = []

    let self = this
    axios.get("./api/media").then(res => {
      self.media = res.data
      self.setState({media: res.data});
      console.log(self.media)
    })
  }
  render() {
    const media = this.media.map((m) =>
      <Media item={m} key={m.name} />
    )

    return (
      // Add your component markup and other subcomponent references here.
      <div>
      <input type="text" spellcheck="false" placeholder="Search Media"></input>
        <ul className="media-list">
          {media}
        </ul>
      </div>
    );
  }
}

export default App
