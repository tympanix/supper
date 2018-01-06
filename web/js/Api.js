import React, { Component } from 'react';
import axios from 'axios';

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
  }

  getSubtitles(folder, media, lang) {
    let config = {
      params: {
        action: "list",
        lang: lang
      }
    }
    let data = Object.assign({}, folder, {filepath: media.filepath})
    return axios.get("./api/subtitles", data, config)
      .then(this.handleError)
  }

  handleError(res) {
    if (res.status !== 200) {
      throw new Error(res.data.message || 'Unknown error')
    } else {
      return res.data
    }
  }
}

export default new Api()
