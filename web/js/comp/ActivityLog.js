import React, { Component } from 'react';

const types = {
  "debug": "notify",
  "info": "success",
  "warn": "warning",
  "error": "error",
  "fatal": "error",
}

const known = {
  "media": m => {
    if (m.year) {
      return `${m.name} (${m.year})`
    }
    if (m.season && m.episode) {
      return `${m.name} S${m.season}E${m.episode}`
    }
  }
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}

class ActivityLog extends Component {
  constructor() {
    super()
  }

  renderTags(data) {
    let tags = []
    for (var key in data) {
      let value = data[key]

      if (!data.hasOwnProperty(key)) {
          continue
      }

      if (known[key]) {
        value = known[key](value)
      }

      if (typeof value !== "string" || !value.length) continue

      tags.push((
        <span key={key}>
          <span className="tag pad">{key}</span>{value}
        </span>
      ))
    }
    return tags
  }

  render() {
    if (!this.props.log || !this.props.log.length) {
      return <h3 className="center meta">No logs to show :(</h3>
    }

    let logs = this.props.log.map((s) => {
      return (
        <li className={types[s.level] || "error"} key={s.wsid}>
          <span className="title">
            {capitalizeFirstLetter(s.level)}
          </span>
          <span className="message">{s.message}</span>
          {this.renderTags(s.data)}
        </li>
      )
    })

    return (
      <ul className="snackbar log">
        {logs}
      </ul>
    )
  }
}

export default ActivityLog
