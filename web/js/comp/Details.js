import React, { Component } from 'react';

import subtitleStore from '../stores/subtitle_store'

import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';

import Search from './Search'
import Spinner from './Spinner'
import FileList from './FileList'
import Flag from './Flag'
import DownloadButtons from './DownloadButtons'
import SubtitleList from './SubtitleList'

import API from '../api'

class Details extends Component {
  constructor() {
    super()

    this.subHotkey = this.subHotkey.bind(this)

    this.state = {
      tabIndex: 0,
      media: undefined,
      folder: undefined,
      busy: false,
      loading: true,
      failed: false,
    }
  }

  componentDidMount() {
    subtitleStore.reset()
    window.addEventListener("keyup", this.subHotkey)
  }

  subHotkey(e) {
    if (e.key === "s") {
      this.downloadSubtitles()
    }
  }

  componentWillUnmount() {
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
    API.downloadSubtitles(folder, lang).then((data) => {
      this.setState({media: data})
    }).catch((err) => {
      console.log(err)
    }).then(() => {
      this.setState({busy: false})
    })
  }

  render() {
    if (this.state.failed) {
      return <h1>No media found</h1>
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
            </Tabs>
          </section>
        </section>
      )
    }

  }

}

export default Details
