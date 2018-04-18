import React, { Component } from 'react'
import axios from 'axios'

import { Redirect } from 'react-router'

import MediaList from './MediaList'
import Spinner from './Spinner'

import { folderStore } from '../stores'

class Search extends Component {
  constructor() {
    super()

    this.getFolders = this.getFolders.bind(this)

    this.state = {
      media: folderStore.getAll(),
      search: folderStore.getSearch(),
      loading: folderStore.isLoading(),
      redirect: false,
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
    this.setState({
      media: folderStore.getAll(),
      loading: folderStore.isLoading()
    })
  }

  handleKey(event) {
    if (event.key === "ArrowDown") {
      this.mediaList.selectNext()
    } else if (event.key === "ArrowUp") {
      this.mediaList.selectPrev()
    } else if (event.key === "Enter") {
      let selected = this.mediaList.getSelected()
      if (selected) {
        let redirect = {pathname: "./details", state: {folder: selected}}
        this.setState({redirect})
      }
    }
  }

  render() {
    if (this.state.redirect) {
      return <Redirect push to={this.state.redirect}/>
    }

    let search = this.state.search.toLowerCase()
    const media = this.state.media.filter((m) => {
      let name = m.name.toLowerCase()
      return name.indexOf(search) !== -1
    })

    return (
      <div>
        <input type="text" className="highlight big" spellCheck="false"
          onKeyUp={this.handleKey.bind(this)}
          ref={(i) => {this.searchInput = i}}
          value={this.state.search}
          onChange={this.search.bind(this)}
          onBlur={this.handleBlur.bind(this)}
          placeholder="Search Media">
        </input>
        {this.state.loading ?
          <Spinner/>
        :
          <MediaList list={media} ref={(m) => {this.mediaList = m}} />
        }
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
