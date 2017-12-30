import React, { Component } from 'react'
import axios from 'axios'

import MediaList from './MediaList'

class Search extends Component {
  constructor() {
    super()

    this.state = {
      media: [],
      search: "",
      selected: -1
    }
    this.update()
  }

  componentDidMount(){
    this.searchInput.focus();
  }

  handleKey(event) {
    if (event.key === "ArrowDown") {
      this.mediaList.selectNext()
    } else if (event.key === "ArrowUp") {
      this.mediaList.selectPrev()
    } else if (event.key === "Enter") {

    }
  }

  render() {
    let search = this.state.search.toLowerCase()
    const media = this.state.media.filter((m) => {
      let name = m.name.toLowerCase()
      return name.includes(search)
    })

    return (
      // Add your component markup and other subcomponent references here.
      <div>
        <input type="text" spellCheck="false"
          onKeyUp={this.handleKey.bind(this)}
          ref={(i) => {this.searchInput = i}}
          value={this.state.search}
          onChange={this.search.bind(this)}
          placeholder="Search Media">
        </input>
        <MediaList list={media} ref={(m) => {this.mediaList = m}} />
      </div>
    );
  }

  search(event) {
    this.mediaList.clearSelected()
    this.setState({
      search: event.target.value
    })
  }

  update() {
    axios.get("./api/media").then(res => {
      this.setState({media: res.data});
    })
  }
}

export default Search
