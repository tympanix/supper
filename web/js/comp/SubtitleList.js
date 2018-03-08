import React, { Component } from 'react';

import subtitleStore from '../stores/subtitle_store'

import FileTags from './FileTags'
import ScoreTag from './ScoreTag'
import Spinner from './Spinner'

class SubtitleList extends Component {
  constructor() {
    super()
    this.state = {
      downloading: null
    }

    Object.assign(this.state, subtitleStore.getState())
    this.updateSubtitles = this.updateSubtitles.bind(this)
  }

  updateSubtitles() {
    console.log("Update subtitles", subtitleStore.getSubtitles())
    this.setState(subtitleStore.getState())
  }

  componentWillMount() {
    subtitleStore.on("change", this.updateSubtitles)
  }

  componentWillUnmount() {
    subtitleStore.removeListener("change", this.updateSubtitles)
  }

  downloadSubtitle(sub) {
    console.log("download sub", sub)
    subtitleStore.download(sub)
  }

  render() {
    if (this.state.loading) {
      return (
        <Spinner />
      )
    }

    if (!this.state.subtitles || this.state.subtitles.length === 0) {
      return (
        <h3 className="meta center">Select a file to display subtitles</h3>
      )
    }

    let subs = this.state.subtitles

    subs = subs.filter((s) => s.language === this.state.lang)

    subs = subs.map((s) => {
      let classes = ["small"]
      s===this.state.downloading && classes.push('loading')

      return (
        <li key={s.hash} className="flex center collapse">
          <div className="col inline spaced flex center nowrap">
            <div className="col">
              <ScoreTag score={s.score}/>
            </div>
            <div className="col name ellipsis">{s.media.name}</div>
          </div>
          <div className="">
            <FileTags media={s}/>
          </div>
          <div className="right">
            <button className={classes.join(" ")} disabled={this.state.downloading}
              onClick={() => this.downloadSubtitle(s)}>
              Download
            </button>
          </div>
        </li>
      )
    })

    return (
      <ul className="subtitle-list">{subs}</ul>
    )
  }
}

export default SubtitleList
