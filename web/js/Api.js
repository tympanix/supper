import React, { Component } from 'react';
import axios from 'axios';

class Api {
  getMediaDetails(media) {
    return axios.post('./api/media', media).then((res) => {
      return res.data
    })
  }
}

export default new Api()
