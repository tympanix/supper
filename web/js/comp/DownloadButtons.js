import React, { Component } from 'react';

import { configStore } from '../stores'

import Flag from './Flag'

class DownloadButtons extends Component {
  constructor() {
    super()

    this.getLanguages = this.getLanguages.bind(this)

    this.state = {
      langs: [],
      clicked: null,
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
    this.setState({
      langs: config.languages || []
    })
  }

  downloadHandler(lang) {
    return (event) => {
      this.setState({clicked: lang})
      this.props.onDownload(lang)
    }
  }

  render() {
    let off = this.props.disabled
    let buttons = this.state.langs.map((l) => {
      let loading = l.code === this.state.clicked
      return (
        <button
          disabled={off}
          className={loading ? 'loading' : null}
          key={l.code}
          onClick={this.downloadHandler(l.code)}>
          {l.language}
        </button>
      )
    })

    return (
      <div className="float center">
        <button disabled={off} onClick={this.downloadHandler()}>Download All</button>
        {buttons}
      </div>
    )
  }
}

export default DownloadButtons
