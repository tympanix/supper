import React, { Component } from 'react';
import axios from 'axios';

import Checkmark from './comp/Checkmark'
import Snackbar from './comp/Snackbar'

class Api {
  getMediaDetails(media) {
    return axios.post('./api/media', media)
      .then(this.handleError)
  }

  getFolders() {
    return axios.get("./api/media")
      .then(this.handleError)
  }

  getConfig() {
    return axios.get("./api/config")
      .then(this.handleError)
  }

  downloadSubtitles(folder, lang) {
    let config = {
      params: {
        action: "download",
        lang: lang
      }
    }
    return axios.post("./api/subtitles", folder, config)
      .then(this.handleError)
      .then(this.showSuccess)
  }

  downloadSingleSubtitle(folder, media, subtitle) {
    let config = {
      params: {
        action: "single"
      }
    }
    let data = Object.assign({}, folder,
      {filepath: media.filepath},
      {link: subtitle.link},
      {language: subtitle.language},
    )
    return axios.post('./api/subtitles', data, config)
      .then(this.handleError)
      .then(this.showSuccess)
  }

  getSubtitles(folder, media) {
    let config = {
      params: {
        action: "list"
      }
    }
    let data = Object.assign({}, folder, {filepath: media.filepath})
    return axios.post("./api/subtitles", data, config)
      .then(this.handleError)
  }

  showSuccess(data) {
    Checkmark.show()
    return data
  }

  handleError(res) {
    if (res.status !== 200) {
      if (res.message && typeof res.message === 'string') {
        Snackbar.error("Error", res.message)
      } else {
        Snackbar.error("Error", "An unexprected error occurred")
      }
      throw new Error(res.data.message || 'Unknown error')
    } else {
      return res.data
    }
  }
}

export default new Api()
