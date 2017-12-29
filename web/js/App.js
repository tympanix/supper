import React, { Component } from 'react'
import axios from 'axios'

import MediaList from './MediaList'

class App extends Component {
  constructor() {
    super()

    this.state = {
      media: [],
      search: ""
    }
    this.update()
  }

  render() {
    const media = this.state.media.filter((m) => {
      let name = m.name.toLowerCase()
      return name.includes(this.state.search)
    })

    return (
      // Add your component markup and other subcomponent references here.
      <div>
        <input type="text" spellCheck="false"
          value={this.state.search}
          onChange={this.search.bind(this)}
          placeholder="Search Media">
        </input>
        <MediaList list={media} />
      </div>
    );
  }

  search(event) {
    this.setState({search: event.target.value.toLowerCase()})
  }

  update() {
    axios.get("./api/media").then(res => {
      this.setState({media: res.data});
    })
  }
}

export default App
