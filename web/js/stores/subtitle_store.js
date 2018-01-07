import { EventEmitter } from 'events'
import API from '../api'

import mock from '../mock/mock_subtitles'

class SubtitleStore extends EventEmitter {
  constructor() {
    super()

    this.state = {
      subtitles: [],
      lang: null,
      media: null,
      folder: null
    }

    Object.assign(this.state, mock)
  }

  update(folder, media, lang) {
    this.state.folder = folder
    this.state.lang = lang
    this.state.media = media
    API.getSubtitles(folder, media).then(subs => {
      this.state.subtitles = subs.sort((a,b) => b.score - a.score)
      this.emit("change")
    })
  }

  getState() {
    return Object.assign({}, this.state)
  }

  download(sub) {
    let f = this.state.folder
    let m = this.state.media
    return API.downloadSingleSubtitle(f, m, sub)
  }

  reset() {
    this.state.subtitles = []
    this.state.lang = null
    this.state.media = null
  }

  getLang() {
    return this.state.lang
  }

  getSubtitles() {
    return this.state.subtitles
  }
}

export default new SubtitleStore
