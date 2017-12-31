import React, { Component } from 'react';

import { configStore } from '../stores'

import Flag from './Flag'

class DownloadButtons extends Component {
  constructor() {
    super()

    this.getLanguages = this.getLanguages.bind(this)

    this.state = {
      langs: []
    }
  }

  componentWillMount() {
    configStore.on("change", this.getLanguages)
    this.getLanguages()
  }

  componentWillUnmount() {
    configStore.removeListener("change", this.getLanguages)
  }

  getLanguages() {
    let config = configStore.getAll()
    console.log("Got config", config)
    this.setState({
      langs: config.languages || []
    })
  }

  render() {
    let buttons = this.state.langs.map((l) => {
      return (
        <button key={l.code} className="icon">
          {l.language}
          <Flag lang={l.code}/>
        </button>
      )
    })

    return (
      <div className="float center">
        <button>Download All</button>
        {buttons}
      </div>
    )
  }
}

export default DownloadButtons
