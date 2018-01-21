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
      downloading: null,
    }
  }

  update(folder, media, lang) {
    this.state.loading = true
    this.state.folder = folder
    this.state.lang = lang
    this.state.media = media
    this.emit("change")
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
    this.state.downloading = sub
    this.emit("change")
    let f = this.state.folder
    let m = this.state.media
    return API.downloadSingleSubtitle(f, m, sub)
      .then((res) => {
        this.state.downloading = null
        this.emit("change")
        return res
      })
  }

  enabled() {
    return !!(this.state.folder && this.state.media)
  }

  reset() {
    this.state.subtitles = []
    this.state.lang = null
    this.state.media = null
    this.state.folder = null
  }

  getLang() {
    return this.state.lang
  }

  getSubtitles() {
    return this.state.subtitles
  }
}

export default new SubtitleStore
