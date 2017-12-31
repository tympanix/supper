import { EventEmitter } from 'events'
import API from '../api'

class ConfigStore extends EventEmitter {
  constructor() {
    super()
    this.config = {}
    this.update()
  }

  update() {
    API.getConfig().then(config => {
      this.config = config
      this.emit("change")
    })
  }

  getAll() {
    return this.config
  }

}

const configStore = new ConfigStore

export default configStore
