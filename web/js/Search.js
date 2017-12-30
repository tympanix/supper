import React, { Component } from 'react'
import axios from 'axios'

import MediaList from './MediaList'

import { folderStore } from './stores'

class Search extends Component {
  constructor() {
    super()

    this.getFolders = this.getFolders.bind(this)

    this.state = {
      media: folderStore.getAll(),
      search: folderStore.getSearch(),
    }
  }

  componentDidMount(){
    this.searchInput.focus();
  }

  componentWillMount() {
    folderStore.on("change", this.getFolders)
  }

  componentWillUnmount() {
    folderStore.removeListener("change", this.getFolders)
  }

  getFolders() {
    this.setState({media: folderStore.getAll()})
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
          onBlur={this.handleBlur.bind(this)}
          placeholder="Search Media">
        </input>
        <MediaList list={media} ref={(m) => {this.mediaList = m}} />
      </div>
    );
  }

  handleBlur() {
    this.mediaList.clearSelected()
  }

  search(event) {
    this.mediaList.clearSelected()
    folderStore.setSearch(event.target.value)
    this.setState({
      search: event.target.value
    })
  }
}

export default Search
