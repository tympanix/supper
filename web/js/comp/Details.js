import React, { Component } from 'react';

import subtitleStore from '../stores/subtitle_store'

import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';

import Search from './Search'
import Spinner from './Spinner'
import FileList from './FileList'
import Flag from './Flag'
import DownloadButtons from './DownloadButtons'
import SubtitleList from './SubtitleList'
import ActivityLog from './ActivityLog'

import API from '../api'
import websocket from '../websocket'

class Details extends Component {
  constructor() {
    super()

    this.subHotkey = this.subHotkey.bind(this)
    this.subtitleFromWebsocket = this.subtitleFromWebsocket.bind(this)

    this.state = {
      tabIndex: 0,
      media: undefined,
      folder: undefined,
      busy: false,
      loading: true,
      failed: false,
      logs: [],
    }
  }

  componentDidMount() {
    subtitleStore.reset()
    websocket.subscribe(this.subtitleFromWebsocket)
    window.addEventListener("keyup", this.subHotkey)
  }

  subHotkey(e) {
    if (e.key === "s") {
      this.downloadSubtitles()
    }
  }

  subtitleFromWebsocket(msg) {
    if (!msg.data.media || !msg.extra.sub) {
      return
    }

    let media = msg.data.media
    let sub = msg.extra.sub

    this.setState((prev) => {
      var found = prev.media.find(m => m.media.id === media.id)

      if (found) {
        found.subtitles.push(sub)
      }

      return {
        "media": prev.media,
        "logs": prev.logs.concat([msg]),
      }
    })
  }

  componentWillUnmount() {
    websocket.remove(this.subtitleFromWebsocket)
    window.removeEventListener("keyup", this.subHotkey)
  }

  componentWillMount() {
    let folder = this.getLocationState()
    API.getMediaDetails(folder)
      .then((media) => this.setState({media: media}))
      .catch(() => this.setState({failed: true}))
      .then(() => this.setState({loading: false}))
  }

  getLocationState() {
    try {
      let folder = this.props.location.state.folder
      if (folder) {
        this.setState({folder: folder})
        return folder
      }
    } catch (e) {
      this.setState({failed: true})
    }
  }

  languageClicked(event, media, lang) {
    let folder = this.state.folder
    this.setState({tabIndex: 1})
    console.log("Clicked", media, lang)
    subtitleStore.update(folder, media, lang)
  }

  downloadSubtitles(lang) {
    this.setState({busy: true})
    let folder = this.state.folder
    API.downloadSubtitles(folder, lang).catch((err) => {
      console.log(err)
    }).then(() => {
      this.setState({busy: false})
    })
  }

  render() {
    if (this.state.failed) {
      return <h1 className="center">No media found</h1>
    }

    if (this.state.loading) {
      return <Spinner/>
    }

    if (this.state.media) {
      return (
        <section>
          <header>
            <h1>{this.state.folder.name}</h1>
          </header>

          <section className="dark">
            <header>
              <h3 className="center">Download Subtitles</h3>
            </header>
            <DownloadButtons
              disabled={this.state.busy}
              onDownload={this.downloadSubtitles.bind(this)}/>
          </section>

          <section>
            <Tabs className="tabs"
              selectedIndex={this.state.tabIndex}
              onSelect={tabIndex => this.setState({ tabIndex })}>
              <TabList className="tablist">
                <Tab selectedClassName="active">Files</Tab>
                <Tab selectedClassName="active">Subtitles</Tab>
                <Tab selectedClassName="active">
                  Activity
                  {this.state.logs.length
                    ? <span className="tag pad">{this.state.logs.length}</span>
                    : null
                  }
                </Tab>
              </TabList>

              <TabPanel className="tab-panel">
                <section>
                  <h2>Files</h2>
                  <FileList files={this.state.media}
                    languageClicked={this.languageClicked.bind(this)}/>
                  <Spinner visible={this.state.busy}/>
                </section>
              </TabPanel>
              <TabPanel className="tab-panel">
                <section>
                  <h2>Subtitles</h2>
                  <SubtitleList/>
                </section>
              </TabPanel>
              <TabPanel className="tab-panel">
                <section>
                  <h2>Activity</h2>
                  <ActivityLog log={this.state.logs}/>
                </section>
              </TabPanel>
            </Tabs>
          </section>
        </section>
      )
    }

  }

}

export default Details
