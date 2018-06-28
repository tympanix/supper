import React, { Component } from 'react';
import axios from 'axios';

import Checkmark from './comp/Checkmark'
import Snackbar from './comp/Snackbar'

class Api {
  constructor() {
    this.handleError = this.handleError.bind(this)
    this.handleException = this.handleException.bind(this)
  }

  getMediaDetails(media) {
    return axios.post('./api/media', media)
      .catch(this.handleException)
      .then(this.handleError)
  }

  getFolders() {
    return axios.get("./api/media")
      .catch(this.handleException)
      .then(this.handleError)
  }

  getConfig() {
    return axios.get("./api/config")
      .catch(this.handleException)
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
      .catch(this.handleException)
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
      .catch(this.handleException)
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
      .catch(this.handleException)
      .then(this.handleError)
  }

  showSuccess(data) {
    Checkmark.show()
    return data
  }

  handleException(error) {
    this.showError(error.response)
    throw error
  }

  showError(res) {
    if (res.data.error && typeof res.data.error === 'string') {
      if (res.status < 300) {
        Snackbar.success("Success", res.data.error)
      } else if (res.status < 400) {
        Snackbar.warning("Warning", res.data.error)
      } else {
        Snackbar.error("Error", res.data.error)
      }
    } else {
      Snackbar.error("Error", "An unexprected error occurred")
    }
  }

  handleError(res) {
    if (res.status !== 200) {
      this.showError(res)
      throw new Error(res.data.message || 'Unknown error')
    } else {
      return res.data
    }
  }
}

export default new Api()
