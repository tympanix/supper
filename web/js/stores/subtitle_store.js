import { EventEmitter } from 'events'
import API from '../api'

class SubtitleStore extends EventEmitter {
  constructor() {
    super()

    this.state = {
      subtitles: [],
      lang: null,
      media: null,
      folder: null,
      loading: false,
    }
  }

  update(folder, media, lang) {
    this.state.loading = true
    this.emit("change")
    this.state.folder = folder
    this.state.lang = lang
    this.state.media = media
    API.getSubtitles(folder, media).then(subs => {
      this.state.loading = false
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
